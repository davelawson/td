package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	selectionPanelMargin       = 42
	selectionPanelWidth        = 390
	selectionPanelPadding      = 18
	selectionPanelRowHeight    = 46
	selectionPanelBottomPad    = 16
	selectionPanelUnknownValue = "Unknown"
)

// SelectionPanelKind identifies the selected object category a panel describes.
type SelectionPanelKind int

const (
	SelectionPanelNone SelectionPanelKind = iota
	SelectionPanelRaider
	SelectionPanelStructure
	SelectionPanelPopulationBuilding
	SelectionPanelEconomicBuilding
	SelectionPanelMarket
	SelectionPanelTower
	SelectionPanelTerrain
)

// ResourceAmounts describes UI-facing construction and currency resource counts.
type ResourceAmounts struct {
	Wood  int
	Stone int
	Iron  int
	Gold  int
}

// MarketTradePrices describes Gold paid for one unit of each Market material.
type MarketTradePrices struct {
	Wood  int
	Stone int
	Iron  int
}

// PopulationAmounts describes UI-facing inhabitant counts by role.
type PopulationAmounts struct {
	Apprentices int
	Soldiers    int
	Peasants    int
}

// TODO: create individual structs for each selection panel subject type
// SelectionPanelData describes the selected object using presentation-neutral facts.
type SelectionPanelData struct {
	Kind                SelectionPanelKind
	Name                string
	Health              int
	MaxHealth           int
	SpeedTilesPerSecond float64
	SanctumDamage       int
	GoldDrop            int
	Cost                ResourceAmounts
	Staffing            PopulationAmounts
	PopulationCost      PopulationAmounts
	PopulationGrant     PopulationAmounts
	ResourceYield       ResourceAmounts
	MarketPrices        MarketTradePrices
	RangeTiles          float64
	FireIntervalSeconds float64
	Damage              int
	TerrainName         string
	BiomeName           string
}

type RaiderSelectionPanelData struct {
	Name                string
	Health              int
	MaxHealth           int
	SpeedTilesPerSecond float64
	SanctumDamage       int
}

type selectionPanelRow struct {
	Label string
	Value string
}

// SelectionPanelBounds returns the bottom-right selection-panel rectangle.
func SelectionPanelBounds(width, height int, data SelectionPanelData) (Button[int], bool) {
	rows, ok := selectionPanelRows(data)
	if !ok {
		return Button[int]{}, false
	}
	panelHeight := selectionPanelPadding + len(rows)*selectionPanelRowHeight + selectionPanelBottomPad
	return Button[int]{
		X: width - selectionPanelMargin - selectionPanelWidth,
		Y: height - selectionPanelMargin - panelHeight,
		W: selectionPanelWidth,
		H: panelHeight,
	}, true
}

// SelectionPanelContains reports whether a point is inside the visible selection panel.
func SelectionPanelContains(width, height int, data SelectionPanelData, x, y int) bool {
	bounds, ok := SelectionPanelBounds(width, height, data)
	return ok && bounds.Contains(x, y)
}

// DrawSelectionPanel renders a selected-object detail panel.
func DrawSelectionPanel(screen *ebiten.Image, face *text.GoTextFace, width, height int, data SelectionPanelData) {
	rows, ok := selectionPanelRows(data)
	if !ok {
		return
	}
	bounds, ok := SelectionPanelBounds(width, height, data)
	if !ok {
		return
	}

	x := float32(bounds.X)
	y := float32(bounds.Y)
	w := float32(bounds.W)
	h := float32(bounds.H)
	vector.FillRect(screen, x, y, w, h, SelectionPanelBackground, false)
	vector.StrokeRect(screen, x, y, w, h, 3, Bronze, false)

	rowX := float64(bounds.X + selectionPanelPadding)
	rowY := float64(bounds.Y + selectionPanelPadding)
	for _, row := range rows {
		DrawText(screen, row.Label, face, rowX, rowY, MutedParchment)
		DrawText(screen, row.Value, face, rowX+20, rowY+20, Parchment)
		rowY += selectionPanelRowHeight
	}
}

func selectionPanelRows(data SelectionPanelData) ([]selectionPanelRow, bool) {
	switch data.Kind {
	case SelectionPanelRaider:
		return raiderSelectionPanelRows(data), true
	case SelectionPanelStructure:
		return []selectionPanelRow{{Label: "Structure", Value: selectedName(data.Name)}}, true
	case SelectionPanelPopulationBuilding:
		return populationBuildingSelectionPanelRows(data), true
	case SelectionPanelEconomicBuilding:
		return economicBuildingSelectionPanelRows(data), true
	case SelectionPanelMarket:
		return marketSelectionPanelRows(data), true
	case SelectionPanelTower:
		return towerSelectionPanelRows(data), true
	case SelectionPanelTerrain:
		return []selectionPanelRow{
			{Label: "Terrain", Value: selectedName(data.TerrainName)},
			{Label: "Biome", Value: selectedName(data.BiomeName)},
		}, true
	default:
		return nil, false
	}
}

// TODO: create individual structs for each selection panel subject type
func raiderSelectionPanelRows(data SelectionPanelData) []selectionPanelRow {
	return []selectionPanelRow{
		{Label: "Raider Type", Value: selectedName(data.Name)},
		{Label: "Health", Value: fmt.Sprintf("%d", data.Health)},
		{Label: "Max Health", Value: fmt.Sprintf("%d", data.MaxHealth)},
		{Label: "Health Remaining", Value: fmt.Sprintf("%d%%", selectedHealthPercent(data.Health, data.MaxHealth))},
		{Label: "Speed", Value: fmt.Sprintf("%.1f tiles/s", data.SpeedTilesPerSecond)},
		{Label: "Sanctum Damage", Value: fmt.Sprintf("%d", data.SanctumDamage)},
		{Label: "Gold Drop", Value: fmt.Sprintf("%d", data.GoldDrop)},
	}
}

func populationBuildingSelectionPanelRows(data SelectionPanelData) []selectionPanelRow {
	rows := []selectionPanelRow{
		{Label: "Structure", Value: selectedName(data.Name)},
		{Label: "Cost", Value: formatSelectionResourceCost(data.Cost)},
	}
	rows = appendPopulationRows(rows, "Consumes", data.PopulationCost)
	rows = appendPopulationRows(rows, "Grants", data.PopulationGrant)
	return rows
}

func economicBuildingSelectionPanelRows(data SelectionPanelData) []selectionPanelRow {
	rows := []selectionPanelRow{
		{Label: "Structure", Value: selectedName(data.Name)},
		{Label: "Cost", Value: formatSelectionResourceCost(data.Cost)},
	}
	rows = appendPopulationRows(rows, "Required", data.Staffing)
	rows = append(rows, selectionPanelRow{Label: "Produces", Value: formatSelectionResourceYield(data.ResourceYield)})
	return rows
}

// marketSelectionPanelRows formats Market cost, staffing, and exchange-rate facts.
func marketSelectionPanelRows(data SelectionPanelData) []selectionPanelRow {
	rows := []selectionPanelRow{
		{Label: "Structure", Value: selectedName(data.Name)},
		{Label: "Cost", Value: formatSelectionResourceCost(data.Cost)},
	}
	rows = appendPopulationRows(rows, "Required", data.Staffing)
	return append(rows,
		selectionPanelRow{Label: "Buys Wood", Value: fmt.Sprintf("1 for %d Gold", data.MarketPrices.Wood)},
		selectionPanelRow{Label: "Buys Stone", Value: fmt.Sprintf("1 for %d Gold", data.MarketPrices.Stone)},
		selectionPanelRow{Label: "Buys Iron", Value: fmt.Sprintf("1 for %d Gold", data.MarketPrices.Iron)},
	)
}

func towerSelectionPanelRows(data SelectionPanelData) []selectionPanelRow {
	rows := []selectionPanelRow{
		{Label: "Tower Type", Value: selectedName(data.Name)},
		{Label: "Range", Value: fmt.Sprintf("%.1f tiles", data.RangeTiles)},
		{Label: "Attack Speed", Value: fmt.Sprintf("every %.1fs", data.FireIntervalSeconds)},
		{Label: "Damage", Value: fmt.Sprintf("%d", data.Damage)},
	}
	rows = appendPopulationRows(rows, "Required", data.Staffing)
	return rows
}

func appendPopulationRows(rows []selectionPanelRow, prefix string, counts PopulationAmounts) []selectionPanelRow {
	if counts.Apprentices > 0 {
		rows = append(rows, selectionPanelRow{Label: prefix + " Apprentices", Value: fmt.Sprintf("%d", counts.Apprentices)})
	}
	if counts.Soldiers > 0 {
		rows = append(rows, selectionPanelRow{Label: prefix + " Soldiers", Value: fmt.Sprintf("%d", counts.Soldiers)})
	}
	if counts.Peasants > 0 {
		rows = append(rows, selectionPanelRow{Label: prefix + " Peasants", Value: fmt.Sprintf("%d", counts.Peasants)})
	}
	return rows
}

func selectedName(name string) string {
	if name == "" {
		return selectionPanelUnknownValue
	}
	return name
}

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

func formatSelectionResourceCost(cost ResourceAmounts) string {
	parts := []string{}
	if cost.Wood > 0 {
		parts = append(parts, fmt.Sprintf("%d Wood", cost.Wood))
	}
	if cost.Stone > 0 {
		parts = append(parts, fmt.Sprintf("%d Stone", cost.Stone))
	}
	if cost.Iron > 0 {
		parts = append(parts, fmt.Sprintf("%d Iron", cost.Iron))
	}
	if cost.Gold > 0 {
		parts = append(parts, fmt.Sprintf("%d Gold", cost.Gold))
	}
	if len(parts) == 0 {
		return "Free"
	}
	return strings.Join(parts, ", ")
}

func formatSelectionResourceYield(yield ResourceAmounts) string {
	cost := formatSelectionResourceCost(yield)
	if cost == "Free" {
		return "Nothing"
	}
	return cost + " after each Raid"
}
