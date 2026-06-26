package game

import (
	"fmt"
	"math"
	"strings"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	selectionPanelMargin       = 42
	selectionPanelWidth        = 390
	selectionPanelPadding      = 18
	selectionPanelTitleGap     = 36
	selectionPanelRowHeight    = 46
	selectionPanelBottomPad    = 16
	selectionPanelTitle        = "Selection"
	selectionPanelUnknownValue = "Unknown"
)

type selectionPanelRow struct {
	Label string
	Value string
}

type selectionPanel struct {
	Rows []selectionPanelRow
}

// currentSelectionPanel returns display rows for the currently selected object.
func (s *State) currentSelectionPanel() (selectionPanel, bool) {
	switch s.selection.kind {
	case selectedItemRaider:
		return s.selectedRaiderPanel()
	case selectedItemStructure:
		return s.selectedStructurePanel()
	default:
		return selectionPanel{}, false
	}
}

// selectedRaiderPanel returns panel rows for the selected active raider.
func (s *State) selectedRaiderPanel() (selectionPanel, bool) {
	for _, enemy := range s.raid.enemies {
		if enemy.id != s.selection.raiderID || enemy.health <= 0 {
			continue
		}

		name := selectionPanelUnknownValue
		maxHealth := 0
		speed := 0.0
		sanctumDamage := 0
		if enemy.template != nil {
			name = enemy.template.Name
			maxHealth = enemy.template.MaxHealth
			speed = enemy.template.SpeedTilesPerSecond
			sanctumDamage = enemy.template.SanctumDamage
		}

		return selectionPanel{Rows: []selectionPanelRow{
			{Label: "Raider Type", Value: name},
			{Label: "Health", Value: fmt.Sprintf("%d", enemy.health)},
			{Label: "Max Health", Value: fmt.Sprintf("%d", maxHealth)},
			{Label: "Health Remaining", Value: fmt.Sprintf("%d%%", selectedHealthPercent(enemy.health, maxHealth))},
			{Label: "Speed", Value: fmt.Sprintf("%.1f tiles/s", speed)},
			{Label: "Sanctum Damage", Value: fmt.Sprintf("%d", sanctumDamage)},
		}}, true
	}
	return selectionPanel{}, false
}

// selectedStructurePanel returns panel rows for the selected structure tile.
func (s *State) selectedStructurePanel() (selectionPanel, bool) {
	tile := s.selection.tile
	if tile.X < 0 || tile.X >= plotSize || tile.Y < 0 || tile.Y >= plotSize {
		return selectionPanel{}, false
	}

	switch s.gameMap.Home.Tiles[tile.Y][tile.X].Feature {
	case featureSanctum:
		return selectionPanel{Rows: []selectionPanelRow{
			{Label: "Structure", Value: s.structureCatalog.Sanctum.Name},
		}}, true
	case featureHouse:
		return houseSelectionPanel(s.structureCatalog.House), true
	case featureBowTower:
		return towerSelectionPanel(s.structureCatalog.BowTower), true
	case featureFlameBoltTower:
		return towerSelectionPanel(s.structureCatalog.FlameBoltTower), true
	case featureCatapultTower:
		return towerSelectionPanel(s.structureCatalog.CatapultTower), true
	default:
		return selectionPanel{}, false
	}
}

// houseSelectionPanel returns display rows for the population-providing House.
func houseSelectionPanel(template StructureTemplate) selectionPanel {
	name := template.Name
	if name == "" {
		name = selectionPanelUnknownValue
	}
	return selectionPanel{Rows: []selectionPanelRow{
		{Label: "Structure", Value: name},
		{Label: "Cost", Value: formatResourceCost(template.Cost)},
		{Label: "Grants Peasants", Value: fmt.Sprintf("%d", template.PopulationGrant.Peasants)},
	}}
}

// towerSelectionPanel returns display rows for one combat tower template.
func towerSelectionPanel(template StructureTemplate) selectionPanel {
	name := template.Name
	if name == "" {
		name = selectionPanelUnknownValue
	}
	rows := []selectionPanelRow{
		{Label: "Tower Type", Value: name},
		{Label: "Range", Value: fmt.Sprintf("%.1f tiles", template.RangeTiles)},
		{Label: "Attack Speed", Value: fmt.Sprintf("every %.1fs", template.FireIntervalSeconds)},
		{Label: "Damage", Value: fmt.Sprintf("%d", template.Damage)},
	}
	if template.Staffing.Apprentices > 0 {
		rows = append(rows, selectionPanelRow{
			Label: "Required Apprentices", Value: fmt.Sprintf("%d", template.Staffing.Apprentices),
		})
	}
	if template.Staffing.Soldiers > 0 {
		rows = append(rows, selectionPanelRow{
			Label: "Required Soldiers", Value: fmt.Sprintf("%d", template.Staffing.Soldiers),
		})
	}
	if template.Staffing.Peasants > 0 {
		rows = append(rows, selectionPanelRow{
			Label: "Required Peasants", Value: fmt.Sprintf("%d", template.Staffing.Peasants),
		})
	}
	return selectionPanel{Rows: rows}
}

// formatResourceCost formats a compact human-readable construction cost.
func formatResourceCost(cost Resources) string {
	parts := []string{}
	if cost.Wood > 0 {
		parts = append(parts, fmt.Sprintf("%d Wood", cost.Wood))
	}
	if cost.Stone > 0 {
		parts = append(parts, fmt.Sprintf("%d Stone", cost.Stone))
	}
	if cost.Metal > 0 {
		parts = append(parts, fmt.Sprintf("%d Metal", cost.Metal))
	}
	if len(parts) == 0 {
		return "Free"
	}
	return strings.Join(parts, ", ")
}

// selectedHealthPercent returns rounded health percentage remaining.
func selectedHealthPercent(health, maxHealth int) int {
	if maxHealth <= 0 {
		return 0
	}
	percent := int(math.Round(float64(health) / float64(maxHealth) * 100))
	if percent < 0 {
		return 0
	}
	if percent > 100 {
		return 100
	}
	return percent
}

// selectionPanelBounds returns the current bottom-right panel rectangle.
func (s *State) selectionPanelBounds() (ui.Button[int], bool) {
	panel, ok := s.currentSelectionPanel()
	if !ok {
		return ui.Button[int]{}, false
	}

	height := selectionPanelPadding + selectionPanelTitleGap + len(panel.Rows)*selectionPanelRowHeight + selectionPanelBottomPad
	return ui.Button[int]{
		X: s.ui.width - selectionPanelMargin - selectionPanelWidth,
		Y: s.ui.height - selectionPanelMargin - height,
		W: selectionPanelWidth,
		H: height,
	}, true
}

// selectionPanelContains reports whether a screen point is inside the visible selection panel.
func (s *State) selectionPanelContains(x, y int) bool {
	bounds, ok := s.selectionPanelBounds()
	return ok && bounds.Contains(x, y)
}

// drawSelectionPanel renders the current selected-object detail panel.
func (s *State) drawSelectionPanel(screen *ebiten.Image) {
	panel, ok := s.currentSelectionPanel()
	if !ok {
		return
	}
	bounds, ok := s.selectionPanelBounds()
	if !ok {
		return
	}

	x := float32(bounds.X)
	y := float32(bounds.Y)
	w := float32(bounds.W)
	h := float32(bounds.H)
	vector.FillRect(screen, x, y, w, h, colors.selectionPanel, false)
	vector.StrokeRect(screen, x, y, w, h, 3, colors.fieldEdge, false)

	titleX := float64(bounds.X + selectionPanelPadding)
	titleY := float64(bounds.Y + selectionPanelPadding - 1)
	ui.DrawText(screen, selectionPanelTitle, s.ui.hudFace, titleX, titleY, colors.pause)

	rowY := float64(bounds.Y + selectionPanelPadding + selectionPanelTitleGap)
	for _, row := range panel.Rows {
		ui.DrawText(screen, row.Label, s.ui.hudFace, titleX, rowY, colors.mutedText)
		ui.DrawText(screen, row.Value, s.ui.hudFace, titleX+20, rowY+20, colors.text)
		rowY += selectionPanelRowHeight
	}
}
