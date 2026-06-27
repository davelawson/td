package game

import (
	"fmt"
	"image/color"
	"strings"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
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
	Bounds ui.Button[int]
	Lines  []buildingTooltipLine
}

type buildingTooltipLine struct {
	Value string
	Color color.Color
}

// hoveredBuildingTooltip returns the tooltip for the currently hovered building icon.
func (s *State) hoveredBuildingTooltip() (buildingTooltip, bool) {
	if s.buildDrag.active || s.ui.buildBarHover < 0 {
		return buildingTooltip{}, false
	}
	items := s.buildingBarItems()
	if s.ui.buildBarHover >= len(items) {
		return buildingTooltip{}, false
	}
	return s.buildingTooltipForItem(items[s.ui.buildBarHover])
}

// buildingTooltipForItem creates tooltip content and bounds for one visible building item.
func (s *State) buildingTooltipForItem(item buildingBarItem) (buildingTooltip, bool) {
	template, ok := s.buildingTemplateForItemID(item.ID)
	if !ok {
		return buildingTooltip{}, false
	}
	lines := buildingTooltipLines(template)
	return buildingTooltip{
		Title:  template.Name,
		Bounds: s.buildingTooltipBounds(item, len(lines)),
		Lines:  lines,
	}, true
}

// buildingTooltipLines returns the player-facing details for one structure template.
func buildingTooltipLines(template StructureTemplate) []buildingTooltipLine {
	lines := []buildingTooltipLine{
		{Value: template.Description, Color: colors.text},
		{Value: "Cost: " + resourcesTooltipText(template.Cost), Color: colors.mutedText},
		{Value: "Staffing: " + staffingTooltipText(template.Staffing), Color: colors.mutedText},
	}
	lines = append(lines, structureEffectTooltipLines(template)...)
	return lines
}

// structureEffectTooltipLines returns implemented effect and stat rows for one template.
func structureEffectTooltipLines(template StructureTemplate) []buildingTooltipLine {
	lines := []buildingTooltipLine{}
	populationEffect := populationEffectTooltipText(template.PopulationCost, template.PopulationGrant)
	if populationEffect != "" {
		lines = append(lines, buildingTooltipLine{Value: "Effect: " + populationEffect, Color: colors.text})
	}
	if yield := resourceYieldTooltipText(template.ResourceYield); yield != "" {
		lines = append(lines, buildingTooltipLine{Value: "Production: " + yield + " after each defeated Raid", Color: colors.text})
	}
	if template.RangeTiles > 0 {
		lines = append(
			lines,
			buildingTooltipLine{Value: fmt.Sprintf("Range: %.1f Tiles", template.RangeTiles), Color: colors.text},
			buildingTooltipLine{Value: fmt.Sprintf("Damage: %d", template.Damage), Color: colors.text},
			buildingTooltipLine{Value: fmt.Sprintf("Fire: every %.1fs", template.FireIntervalSeconds), Color: colors.text},
			buildingTooltipLine{Value: fmt.Sprintf("Projectile: %.1f Tiles/s", template.ProjectileSpeedTilesPerSecond), Color: colors.text},
		)
		if template.DamageAllEnemiesInTargetTile {
			lines = append(lines, buildingTooltipLine{Value: "Area: damages every enemy in the target Tile", Color: colors.text})
		}
	}
	return lines
}

// buildingTooltipBounds returns a tooltip rectangle clamped to the drawable area.
func (s *State) buildingTooltipBounds(item buildingBarItem, lineCount int) ui.Button[int] {
	height := buildingTooltipPadding*2 + buildingTooltipTitleStep + buildingTooltipSectionGap + lineCount*buildingTooltipLineStep
	x := s.buildingBarBounds().X + s.buildingBarBounds().W + buildingTooltipGap
	y := item.Bounds.Y

	if x+buildingTooltipWidth > s.ui.width-buildingTooltipScreenGap {
		x = s.ui.width - buildingTooltipWidth - buildingTooltipScreenGap
	}
	if x < buildingTooltipScreenGap {
		x = buildingTooltipScreenGap
	}
	if y+height > s.ui.height-buildingTooltipScreenGap {
		y = s.ui.height - height - buildingTooltipScreenGap
	}
	minY := topBarHeight + buildingTooltipScreenGap
	if y < minY {
		y = minY
	}

	return ui.Button[int]{
		X: x,
		Y: y,
		W: buildingTooltipWidth,
		H: height,
	}
}

// drawBuildingTooltip renders the current building hover tooltip if one is active.
func (s *State) drawBuildingTooltip(screen *ebiten.Image) {
	tooltip, ok := s.hoveredBuildingTooltip()
	if !ok {
		return
	}

	bounds := tooltip.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), colors.selectionPanel, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, colors.fieldEdge, false)

	x := float64(bounds.X + buildingTooltipPadding)
	y := float64(bounds.Y + buildingTooltipPadding)
	ui.DrawText(screen, tooltip.Title, s.ui.costBoldFace, x, y, colors.pause)
	y += buildingTooltipTitleStep + buildingTooltipSectionGap
	for _, line := range tooltip.Lines {
		ui.DrawText(screen, line.Value, s.ui.costFace, x, y, line.Color)
		y += buildingTooltipLineStep
	}
}

// resourcesTooltipText formats construction resources in stable Wood, Stone, Metal order.
func resourcesTooltipText(resources Resources) string {
	parts := []string{}
	if resources.Wood > 0 {
		parts = append(parts, countLabel(resources.Wood, "Wood", "Wood"))
	}
	if resources.Stone > 0 {
		parts = append(parts, countLabel(resources.Stone, "Stone", "Stone"))
	}
	if resources.Metal > 0 {
		parts = append(parts, countLabel(resources.Metal, "Metal", "Metal"))
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, ", ")
}

// staffingTooltipText formats required staff in stable Apprentice, Soldier, Peasant order.
func staffingTooltipText(staffing StaffingRequirements) string {
	parts := []string{}
	if staffing.Apprentices > 0 {
		parts = append(parts, countLabel(staffing.Apprentices, "Apprentice", "Apprentices"))
	}
	if staffing.Soldiers > 0 {
		parts = append(parts, countLabel(staffing.Soldiers, "Soldier", "Soldiers"))
	}
	if staffing.Peasants > 0 {
		parts = append(parts, countLabel(staffing.Peasants, "Peasant", "Peasants"))
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, ", ")
}

// populationEffectTooltipText formats immediate population costs and grants.
func populationEffectTooltipText(cost PopulationCost, grant PopulationGrant) string {
	parts := []string{}
	if cost.Apprentices > 0 {
		parts = append(parts, "-"+countLabel(cost.Apprentices, "Apprentice", "Apprentices"))
	}
	if cost.Soldiers > 0 {
		parts = append(parts, "-"+countLabel(cost.Soldiers, "Soldier", "Soldiers"))
	}
	if cost.Peasants > 0 {
		parts = append(parts, "-"+countLabel(cost.Peasants, "Peasant", "Peasants"))
	}
	if grant.Apprentices > 0 {
		parts = append(parts, "+"+countLabel(grant.Apprentices, "Apprentice", "Apprentices"))
	}
	if grant.Soldiers > 0 {
		parts = append(parts, "+"+countLabel(grant.Soldiers, "Soldier", "Soldiers"))
	}
	if grant.Peasants > 0 {
		parts = append(parts, "+"+countLabel(grant.Peasants, "Peasant", "Peasants"))
	}
	return strings.Join(parts, ", ")
}

// resourceYieldTooltipText formats post-Raid resource production.
func resourceYieldTooltipText(resources Resources) string {
	parts := []string{}
	if resources.Wood > 0 {
		parts = append(parts, "+"+countLabel(resources.Wood, "Wood", "Wood"))
	}
	if resources.Stone > 0 {
		parts = append(parts, "+"+countLabel(resources.Stone, "Stone", "Stone"))
	}
	if resources.Metal > 0 {
		parts = append(parts, "+"+countLabel(resources.Metal, "Metal", "Metal"))
	}
	return strings.Join(parts, ", ")
}

// countLabel formats a count with singular or plural text.
func countLabel(count int, singular, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	return fmt.Sprintf("%d %s", count, plural)
}
