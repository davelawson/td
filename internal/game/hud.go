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
	topBarHeight             = 86
	topBarMargin             = 42
	resourceIconDisplaySize  = 28
	resourceIconTextGap      = 7
	resourceHUDItemGap       = 22
	resourceBarricadeTextGap = 24
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

type gameStatus struct {
	phase     phase
	chapter   string
	day       int
	calmTime  int
	barricade int
	resources resourceCounts
}

type resourceHUDItem struct {
	Name   string
	Count  int
	Sprite *ebiten.Image
	Color  color.Color
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
			wood:  80,
			stone: 45,
			metal: 12,
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

// barricadeText formats Sanctum defense status for the top bar.
func (s *State) barricadeText() string {
	return fmt.Sprintf("| Barricade %d", s.status.barricade)
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

	s.drawResourceStatus(screen)
}

// drawResourceStatus renders resource icons, counts, and Barricade status.
func (s *State) drawResourceStatus(screen *ebiten.Image) {
	items := s.resourceHUDItems()
	barricade := s.barricadeText()
	totalWidth := s.resourceStatusWidth(items, barricade)
	x := float64(s.ui.width) - totalWidth - topBarMargin

	for i, item := range items {
		itemWidth := s.resourceHUDItemWidth(item)
		s.drawResourceHUDItem(screen, item, x)
		x += itemWidth
		if i < len(items)-1 {
			x += resourceHUDItemGap
		}
	}

	x += resourceBarricadeTextGap
	ui.DrawText(screen, barricade, s.ui.hudFace, x, 29, colors.text)
}

// resourceStatusWidth measures the full right-side HUD group.
func (s *State) resourceStatusWidth(items []resourceHUDItem, barricade string) float64 {
	total := 0.0
	for i, item := range items {
		total += s.resourceHUDItemWidth(item)
		if i < len(items)-1 {
			total += resourceHUDItemGap
		}
	}
	barricadeWidth, _ := text.Measure(barricade, s.ui.hudFace, s.ui.hudFace.Size)
	return total + resourceBarricadeTextGap + barricadeWidth
}

// resourceHUDItemWidth measures one resource icon and count pair.
func (s *State) resourceHUDItemWidth(item resourceHUDItem) float64 {
	count := fmt.Sprintf("%d", item.Count)
	countWidth, _ := text.Measure(count, s.ui.hudFace, s.ui.hudFace.Size)
	return resourceIconDisplaySize + resourceIconTextGap + countWidth
}

// drawResourceHUDItem renders one resource icon and count pair.
func (s *State) drawResourceHUDItem(screen *ebiten.Image, item resourceHUDItem, x float64) {
	if item.Sprite != nil {
		spriteWidth := float64(item.Sprite.Bounds().Dx())
		spriteHeight := float64(item.Sprite.Bounds().Dy())
		if spriteWidth > 0 && spriteHeight > 0 {
			scale := float64(resourceIconDisplaySize) / spriteWidth
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(scale, scale)
			options.GeoM.Translate(x, float64(topBarHeight-resourceIconDisplaySize)/2)
			screen.DrawImage(item.Sprite, options)
		}
	}

	count := fmt.Sprintf("%d", item.Count)
	ui.DrawText(screen, count, s.ui.hudFace, x+resourceIconDisplaySize+resourceIconTextGap, 29, item.Color)
}
