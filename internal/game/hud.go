package game

import (
	"fmt"
	"image/color"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	topBarHeight            = 86
	topBarMargin            = 42
	statusIconDisplaySize   = 28
	statusIconTextGap       = 7
	statusHUDItemGap        = 18
	statusGroupSeparatorGap = 16
)

type phase int

const (
	phaseCalm phase = iota
	phaseRaid
)

type resourceCounts struct {
	wood  int
	stone int
	metal int
}

type populationCount struct {
	available int
	total     int
}

type populationCounts struct {
	apprentices populationCount
	soldiers    populationCount
	peasants    populationCount
}

type gameStatus struct {
	phase       phase
	chapter     string
	day         int
	calmTime    int
	barricade   int
	resources   resourceCounts
	populations populationCounts
}

// canStaff reports whether every required inhabitant role is currently available.
func (s *State) canStaff(requirements StaffingRequirements) bool {
	return s.status.populations.apprentices.available >= requirements.Apprentices &&
		s.status.populations.soldiers.available >= requirements.Soldiers &&
		s.status.populations.peasants.available >= requirements.Peasants
}

// reserveStaffing commits available inhabitants to a newly built structure.
func (s *State) reserveStaffing(requirements StaffingRequirements) {
	s.status.populations.apprentices.available -= requirements.Apprentices
	s.status.populations.soldiers.available -= requirements.Soldiers
	s.status.populations.peasants.available -= requirements.Peasants
}

// grantPopulation adds new inhabitants made available by a constructed structure.
func (s *State) grantPopulation(grant PopulationGrant) {
	s.status.populations.apprentices.available += grant.Apprentices
	s.status.populations.apprentices.total += grant.Apprentices
	s.status.populations.soldiers.available += grant.Soldiers
	s.status.populations.soldiers.total += grant.Soldiers
	s.status.populations.peasants.available += grant.Peasants
	s.status.populations.peasants.total += grant.Peasants
}

type resourceHUDItem struct {
	Name   string
	Count  int
	Sprite *ebiten.Image
	Color  color.Color
}

type populationHUDItem struct {
	Name      string
	Available int
	Total     int
	Sprite    *ebiten.Image
	Color     color.Color
}

// setPrototypeGameStatus initializes fixed state shown before gameplay systems exist.
func (s *State) setPrototypeGameStatus() {
	s.status = gameStatus{
		phase:     phaseCalm,
		chapter:   "Chapter I: The Ashen Copse",
		day:       1,
		calmTime:  120,
		barricade: 3,
		resources: resourceCounts{
			wood:  100,
			stone: 50,
			metal: 20,
		},
		populations: populationCounts{
			apprentices: populationCount{},
			soldiers:    populationCount{},
			peasants:    populationCount{},
		},
	}
}

// chapterDayText formats the Chapter and Day summary for the top bar.
func (s *State) chapterDayText() string {
	return fmt.Sprintf("%s | Day %d", s.status.chapter, s.status.day)
}

// phaseText formats the phase-specific top bar status.
func (s *State) phaseText() string {
	if s.raid.breached {
		return "Sanctum breached"
	}
	switch s.status.phase {
	case phaseRaid:
		return fmt.Sprintf("Enemies remaining: %d", s.raidEnemiesRemaining())
	default:
		minutes := s.status.calmTime / 60
		seconds := s.status.calmTime % 60
		return fmt.Sprintf("Raid in %02d:%02d", minutes, seconds)
	}
}

// resourceHUDItems returns the resources shown in the top bar from left to right.
func (s *State) resourceHUDItems() []resourceHUDItem {
	return []resourceHUDItem{
		{Name: "Wood", Count: s.status.resources.wood, Sprite: s.assetCatalog.Sprite.Icon.Wood, Color: colors.resourceWood},
		{Name: "Stone", Count: s.status.resources.stone, Sprite: s.assetCatalog.Sprite.Icon.Stone, Color: colors.resourceStone},
		{Name: "Metal", Count: s.status.resources.metal, Sprite: s.assetCatalog.Sprite.Icon.Metal, Color: colors.resourceMetal},
	}
}

// populationHUDItems returns the inhabitant populations shown in the top bar from left to right.
func (s *State) populationHUDItems() []populationHUDItem {
	return []populationHUDItem{
		{
			Name:      "Apprentice",
			Available: s.status.populations.apprentices.available,
			Total:     s.status.populations.apprentices.total,
			Sprite:    s.assetCatalog.Sprite.Icon.Apprentice,
			Color:     colors.text,
		},
		{
			Name:      "Soldier",
			Available: s.status.populations.soldiers.available,
			Total:     s.status.populations.soldiers.total,
			Sprite:    s.assetCatalog.Sprite.Icon.Soldier,
			Color:     colors.text,
		},
		{
			Name:      "Peasant",
			Available: s.status.populations.peasants.available,
			Total:     s.status.populations.peasants.total,
			Sprite:    s.assetCatalog.Sprite.Icon.Peasant,
			Color:     colors.text,
		},
	}
}

// barricadeText formats Sanctum defense status for the top bar.
func (s *State) barricadeText() string {
	return fmt.Sprintf("Barricade %d", s.status.barricade)
}

// drawTopBar renders the game status bar at the top of the screen.
func (s *State) drawTopBar(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(s.ui.width), topBarHeight, colors.topBar, false)
	vector.StrokeLine(screen, 0, topBarHeight-2, float32(s.ui.width), topBarHeight-2, 3, colors.topBarEdge, false)

	left := s.chapterDayText()
	center := s.phaseText()

	ui.DrawText(screen, left, s.ui.hudFace, topBarMargin, 29, colors.text)

	centerWidth, _ := text.Measure(center, s.ui.hudFace, s.ui.hudFace.Size)
	ui.DrawText(screen, center, s.ui.hudFace, (float64(s.ui.width)-centerWidth)/2, 29, colors.pause)

	s.drawDomainStatus(screen)
}

// drawDomainStatus renders resource, population, and Barricade status groups.
func (s *State) drawDomainStatus(screen *ebiten.Image) {
	resources := s.resourceHUDItems()
	populations := s.populationHUDItems()
	barricade := s.barricadeText()
	totalWidth := s.domainStatusWidth(resources, populations, barricade)
	x := float64(s.ui.width) - totalWidth - topBarMargin

	for i, item := range resources {
		itemWidth := s.resourceHUDItemWidth(item)
		s.drawResourceHUDItem(screen, item, x)
		x += itemWidth
		if i < len(resources)-1 {
			x += statusHUDItemGap
		}
	}

	x = s.drawStatusGroupSeparator(screen, x)
	for i, item := range populations {
		itemWidth := s.populationHUDItemWidth(item)
		s.drawPopulationHUDItem(screen, item, x)
		x += itemWidth
		if i < len(populations)-1 {
			x += statusHUDItemGap
		}
	}

	x = s.drawStatusGroupSeparator(screen, x)
	ui.DrawText(screen, barricade, s.ui.hudFace, x, 29, colors.text)
}

// domainStatusWidth measures the full right-side top-bar status.
func (s *State) domainStatusWidth(resources []resourceHUDItem, populations []populationHUDItem, barricade string) float64 {
	total := s.resourceHUDGroupWidth(resources)
	total += s.statusGroupSeparatorWidth()
	total += s.populationHUDGroupWidth(populations)
	total += s.statusGroupSeparatorWidth()
	barricadeWidth, _ := text.Measure(barricade, s.ui.hudFace, s.ui.hudFace.Size)
	return total + barricadeWidth
}

// resourceHUDGroupWidth measures a resource icon-and-count group.
func (s *State) resourceHUDGroupWidth(items []resourceHUDItem) float64 {
	total := 0.0
	for i, item := range items {
		total += s.resourceHUDItemWidth(item)
		if i < len(items)-1 {
			total += statusHUDItemGap
		}
	}
	return total
}

// populationHUDGroupWidth measures a population icon-and-value group.
func (s *State) populationHUDGroupWidth(items []populationHUDItem) float64 {
	total := 0.0
	for i, item := range items {
		total += s.populationHUDItemWidth(item)
		if i < len(items)-1 {
			total += statusHUDItemGap
		}
	}
	return total
}

// resourceHUDItemWidth measures one resource icon and count pair.
func (s *State) resourceHUDItemWidth(item resourceHUDItem) float64 {
	return s.statusHUDItemWidth(fmt.Sprintf("%d", item.Count))
}

// populationHUDItemWidth measures one population icon and available/total pair.
func (s *State) populationHUDItemWidth(item populationHUDItem) float64 {
	return s.statusHUDItemWidth(populationHUDItemText(item))
}

// statusHUDItemWidth measures one icon and value pair.
func (s *State) statusHUDItemWidth(value string) float64 {
	valueWidth, _ := text.Measure(value, s.ui.hudFace, s.ui.hudFace.Size)
	return statusIconDisplaySize + statusIconTextGap + valueWidth
}

// statusGroupSeparatorWidth measures one padded separator between status groups.
func (s *State) statusGroupSeparatorWidth() float64 {
	separatorWidth, _ := text.Measure("|", s.ui.hudFace, s.ui.hudFace.Size)
	return statusGroupSeparatorGap + separatorWidth + statusGroupSeparatorGap
}

// drawStatusGroupSeparator renders a padded separator and returns the next draw position.
func (s *State) drawStatusGroupSeparator(screen *ebiten.Image, x float64) float64 {
	x += statusGroupSeparatorGap
	ui.DrawText(screen, "|", s.ui.hudFace, x, 29, colors.mutedText)
	separatorWidth, _ := text.Measure("|", s.ui.hudFace, s.ui.hudFace.Size)
	return x + separatorWidth + statusGroupSeparatorGap
}

// drawResourceHUDItem renders one resource icon and count pair.
func (s *State) drawResourceHUDItem(screen *ebiten.Image, item resourceHUDItem, x float64) {
	count := fmt.Sprintf("%d", item.Count)
	s.drawStatusHUDItem(screen, item.Sprite, count, item.Color, x)
}

// drawPopulationHUDItem renders one population icon and available/total pair.
func (s *State) drawPopulationHUDItem(screen *ebiten.Image, item populationHUDItem, x float64) {
	s.drawStatusHUDItem(screen, item.Sprite, populationHUDItemText(item), item.Color, x)
}

// populationHUDItemText formats a population value with available before total.
func populationHUDItemText(item populationHUDItem) string {
	return fmt.Sprintf("%d/%d", item.Available, item.Total)
}

// drawStatusHUDItem renders one status icon and its value.
func (s *State) drawStatusHUDItem(screen *ebiten.Image, sprite *ebiten.Image, value string, valueColor color.Color, x float64) {
	if sprite != nil {
		spriteWidth := float64(sprite.Bounds().Dx())
		spriteHeight := float64(sprite.Bounds().Dy())
		if spriteWidth > 0 && spriteHeight > 0 {
			scale := float64(statusIconDisplaySize) / spriteWidth
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(scale, scale)
			options.GeoM.Translate(x, float64(topBarHeight-statusIconDisplaySize)/2)
			screen.DrawImage(sprite, options)
		}
	}

	ui.DrawText(screen, value, s.ui.hudFace, x+statusIconDisplaySize+statusIconTextGap, 29, valueColor)
}
