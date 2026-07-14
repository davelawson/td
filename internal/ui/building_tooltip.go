package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	buildingTooltipWidth      = 500
	buildingTooltipPadding    = 12
	buildingTooltipGap        = 12
	buildingTooltipScreenGap  = 8
	buildingTooltipTitleStep  = 26
	buildingTooltipLineStep   = 20
	buildingTooltipSectionGap = 6
)

type buildingTooltip struct {
	Title  string
	Bounds Button[int]
	Lines  []buildingTooltipLine
}

type buildingTooltipLine struct {
	Value string
	Color color.Color
}

// DrawBuildingTooltip renders the tooltip for the currently hovered building icon.
func DrawBuildingTooltip(screen *ebiten.Image, regularFace, boldFace *text.GoTextFace, width, height, top int, model BuildingBarModel) {
	tooltip, ok := hoveredBuildingTooltip(width, height, top, model)
	if !ok {
		return
	}

	bounds := tooltip.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), SelectionPanelBackground, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, Bronze, false)

	x := float64(bounds.X + buildingTooltipPadding)
	y := float64(bounds.Y + buildingTooltipPadding)
	DrawText(screen, tooltip.Title, boldFace, x, y, LightBronze)
	y += buildingTooltipTitleStep + buildingTooltipSectionGap
	for _, line := range tooltip.Lines {
		DrawText(screen, line.Value, regularFace, x, y, line.Color)
		y += buildingTooltipLineStep
	}
}

// hoveredBuildingTooltip returns tooltip content for the visible hovered item.
func hoveredBuildingTooltip(width, height, top int, model BuildingBarModel) (buildingTooltip, bool) {
	items := buildingBarLayoutItems(top, model)
	if model.HoveredItem < 0 || model.HoveredItem >= len(items) {
		return buildingTooltip{}, false
	}
	item := items[model.HoveredItem]
	lines := buildingTooltipLines(item.BuildingBarItem)
	return buildingTooltip{
		Title:  item.Name,
		Bounds: buildingTooltipBounds(width, height, top, item.Bounds, len(lines)),
		Lines:  lines,
	}, true
}

// buildingTooltipLines returns player-facing details for one building choice.
func buildingTooltipLines(item BuildingBarItem) []buildingTooltipLine {
	lines := []buildingTooltipLine{
		{Value: item.Description, Color: Parchment},
		{Value: "Cost: " + buildingResourcesTooltipText(item.Cost), Color: MutedParchment},
		{Value: "Staffing: " + buildingPopulationTooltipText(item.Staffing), Color: MutedParchment},
	}
	populationEffect := buildingPopulationEffectTooltipText(item.PopulationCost, item.PopulationGrant)
	if populationEffect != "" {
		lines = append(lines, buildingTooltipLine{Value: "Effect: " + populationEffect, Color: Parchment})
	}
	if yield := buildingResourceYieldTooltipText(item.ResourceYield); yield != "" {
		lines = append(lines, buildingTooltipLine{Value: "Production: " + yield + " during each Labour phase", Color: Parchment})
	}
	if item.RangeTiles > 0 {
		lines = append(lines,
			buildingTooltipLine{Value: fmt.Sprintf("Range: %.1f Tiles", item.RangeTiles), Color: Parchment},
			buildingTooltipLine{Value: fmt.Sprintf("Damage: %d", item.Damage), Color: Parchment},
			buildingTooltipLine{Value: fmt.Sprintf("Fire: every %.1fs", item.FireIntervalSeconds), Color: Parchment},
			buildingTooltipLine{Value: fmt.Sprintf("Projectile: %.1f Tiles/s", item.ProjectileSpeedTilesPerSecond), Color: Parchment},
		)
		if item.DamageAllEnemiesInTargetTile {
			lines = append(lines, buildingTooltipLine{Value: "Area: damages every enemy in the target Tile", Color: Parchment})
		}
	}
	return lines
}

// buildingTooltipBounds returns a tooltip rectangle clamped to the drawable area.
func buildingTooltipBounds(width, height, top int, itemBounds Button[int], lineCount int) Button[int] {
	tooltipHeight := buildingTooltipPadding*2 + buildingTooltipTitleStep + buildingTooltipSectionGap + lineCount*buildingTooltipLineStep
	x := BuildingBarWidth + buildingTooltipGap
	y := itemBounds.Y
	if x+buildingTooltipWidth > width-buildingTooltipScreenGap {
		x = width - buildingTooltipWidth - buildingTooltipScreenGap
	}
	if x < buildingTooltipScreenGap {
		x = buildingTooltipScreenGap
	}
	if y+tooltipHeight > height-buildingTooltipScreenGap {
		y = height - tooltipHeight - buildingTooltipScreenGap
	}
	minimumY := top + buildingTooltipScreenGap
	if y < minimumY {
		y = minimumY
	}
	return Button[int]{X: x, Y: y, W: buildingTooltipWidth, H: tooltipHeight}
}

// buildingResourcesTooltipText formats resources in stable Wood, Stone, Iron, Gold order.
func buildingResourcesTooltipText(resources ResourceAmounts) string {
	parts := []string{}
	if resources.Wood > 0 {
		parts = append(parts, buildingCountLabel(resources.Wood, "Wood", "Wood"))
	}
	if resources.Stone > 0 {
		parts = append(parts, buildingCountLabel(resources.Stone, "Stone", "Stone"))
	}
	if resources.Iron > 0 {
		parts = append(parts, buildingCountLabel(resources.Iron, "Iron", "Iron"))
	}
	if resources.Gold > 0 {
		parts = append(parts, buildingCountLabel(resources.Gold, "Gold", "Gold"))
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, ", ")
}

// buildingPopulationTooltipText formats roles in stable Apprentice, Soldier, Peasant order.
func buildingPopulationTooltipText(population PopulationAmounts) string {
	parts := []string{}
	if population.Apprentices > 0 {
		parts = append(parts, buildingCountLabel(population.Apprentices, "Apprentice", "Apprentices"))
	}
	if population.Soldiers > 0 {
		parts = append(parts, buildingCountLabel(population.Soldiers, "Soldier", "Soldiers"))
	}
	if population.Peasants > 0 {
		parts = append(parts, buildingCountLabel(population.Peasants, "Peasant", "Peasants"))
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, ", ")
}

// buildingPopulationEffectTooltipText formats immediate population costs and grants.
func buildingPopulationEffectTooltipText(cost, grant PopulationAmounts) string {
	parts := []string{}
	appendCounts := func(prefix string, population PopulationAmounts) {
		if population.Apprentices > 0 {
			parts = append(parts, prefix+buildingCountLabel(population.Apprentices, "Apprentice", "Apprentices"))
		}
		if population.Soldiers > 0 {
			parts = append(parts, prefix+buildingCountLabel(population.Soldiers, "Soldier", "Soldiers"))
		}
		if population.Peasants > 0 {
			parts = append(parts, prefix+buildingCountLabel(population.Peasants, "Peasant", "Peasants"))
		}
	}
	appendCounts("-", cost)
	appendCounts("+", grant)
	return strings.Join(parts, ", ")
}

// buildingResourceYieldTooltipText formats Labour resource production.
func buildingResourceYieldTooltipText(resources ResourceAmounts) string {
	parts := []string{}
	if resources.Wood > 0 {
		parts = append(parts, "+"+buildingCountLabel(resources.Wood, "Wood", "Wood"))
	}
	if resources.Stone > 0 {
		parts = append(parts, "+"+buildingCountLabel(resources.Stone, "Stone", "Stone"))
	}
	if resources.Iron > 0 {
		parts = append(parts, "+"+buildingCountLabel(resources.Iron, "Iron", "Iron"))
	}
	if resources.Gold > 0 {
		parts = append(parts, "+"+buildingCountLabel(resources.Gold, "Gold", "Gold"))
	}
	return strings.Join(parts, ", ")
}

// buildingCountLabel formats a count with singular or plural text.
func buildingCountLabel(count int, singular, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	return fmt.Sprintf("%d %s", count, plural)
}
