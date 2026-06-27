package game

import (
	"fmt"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// drawBuildingBarCost renders non-zero resource costs to the right of one icon.
func (s *State) drawBuildingBarCost(screen *ebiten.Image, item buildingBarItem, hovered bool) {
	costItems := buildingBarCostItems(item.Cost)
	if len(costItems) == 0 {
		return
	}

	face := s.buildingBarCostFace(hovered)
	x := float64(s.buildingBarMetadataX(item))
	y := float64(item.Bounds.Y + buildingBarCostOffsetY)
	for i, costItem := range costItems {
		width, _ := text.Measure(costItem.Value, face, face.Size)
		if hovered {
			ui.DrawText(screen, costItem.Value, face, x+1, y+1, buildingBarCostShadow)
			ui.DrawText(screen, costItem.Value, face, x-1, y+1, buildingBarCostShadow)
		}
		ui.DrawText(screen, costItem.Value, face, x, y, costItem.Color)
		x += width
		if i < len(costItems)-1 {
			x += buildingBarCostItemGap
		}
	}
}

// drawBuildingBarPopulationMetadata renders staffing requirements or population grants beside one icon.
func (s *State) drawBuildingBarPopulationMetadata(screen *ebiten.Image, item buildingBarItem) {
	metadataItems := s.buildingBarPopulationMetadataItems(item)
	if len(metadataItems) == 0 {
		return
	}

	x := float64(s.buildingBarMetadataX(item))
	y := float64(item.Bounds.Y + buildingBarStaffingOffsetY)
	for i, staffingItem := range metadataItems {
		if staffingItem.Sprite != nil {
			spriteWidth := float64(staffingItem.Sprite.Bounds().Dx())
			spriteHeight := float64(staffingItem.Sprite.Bounds().Dy())
			if spriteWidth > 0 && spriteHeight > 0 {
				options := &ebiten.DrawImageOptions{}
				scale := float64(buildingBarStaffIconSize) / spriteWidth
				options.GeoM.Scale(scale, scale)
				options.GeoM.Translate(x, y)
				screen.DrawImage(staffingItem.Sprite, options)
			}
		}
		x += buildingBarStaffIconSize + buildingBarStaffIconGap
		value := staffingItem.Value
		if value == "" {
			value = fmt.Sprintf("%d", staffingItem.Count)
		}
		ui.DrawText(screen, value, s.ui.costFace, x, y-1, colors.text)
		valueWidth, _ := text.Measure(value, s.ui.costFace, s.ui.costFace.Size)
		x += valueWidth
		if i < len(metadataItems)-1 {
			x += buildingBarCostItemGap
		}
	}
}

// buildingBarMetadataX returns the x coordinate where right-side item values begin.
func (s *State) buildingBarMetadataX(item buildingBarItem) int {
	return item.Bounds.X + item.Bounds.W + buildingBarMetadataGap
}

// buildingBarMetadataRight returns the right edge available to item values.
func (s *State) buildingBarMetadataRight() int {
	bar := s.buildingBarBounds()
	return bar.X + bar.W - buildingBarPadding
}

// buildingBarPopulationMetadataItems returns the row shown beneath one structure cost.
func (s *State) buildingBarPopulationMetadataItems(item buildingBarItem) []buildingBarStaffingItem {
	staffingItems := s.buildingBarStaffingItems(item.Staffing)
	if len(staffingItems) > 0 {
		return staffingItems
	}
	return append(
		s.buildingBarPopulationCostItems(item.PopulationCost),
		s.buildingBarPopulationGrantItems(item.PopulationGrant)...,
	)
}

// buildingBarStaffingItems returns non-zero requirements in Apprentice, Soldier, Peasant order.
func (s *State) buildingBarStaffingItems(requirements StaffingRequirements) []buildingBarStaffingItem {
	items := []buildingBarStaffingItem{}
	if requirements.Apprentices > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: requirements.Apprentices, Sprite: s.assetCatalog.Sprite.Icon.Apprentice,
		})
	}
	if requirements.Soldiers > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: requirements.Soldiers, Sprite: s.assetCatalog.Sprite.Icon.Soldier,
		})
	}
	if requirements.Peasants > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: requirements.Peasants, Sprite: s.assetCatalog.Sprite.Icon.Peasant,
		})
	}
	return items
}

// buildingBarPopulationGrantItems returns non-zero grants in Apprentice, Soldier, Peasant order.
func (s *State) buildingBarPopulationGrantItems(grant PopulationGrant) []buildingBarStaffingItem {
	items := []buildingBarStaffingItem{}
	if grant.Apprentices > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: grant.Apprentices, Value: fmt.Sprintf("+%d", grant.Apprentices), Sprite: s.assetCatalog.Sprite.Icon.Apprentice,
		})
	}
	if grant.Soldiers > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: grant.Soldiers, Value: fmt.Sprintf("+%d", grant.Soldiers), Sprite: s.assetCatalog.Sprite.Icon.Soldier,
		})
	}
	if grant.Peasants > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: grant.Peasants, Value: fmt.Sprintf("+%d", grant.Peasants), Sprite: s.assetCatalog.Sprite.Icon.Peasant,
		})
	}
	return items
}

// buildingBarPopulationCostItems returns non-zero population costs in Apprentice, Soldier, Peasant order.
func (s *State) buildingBarPopulationCostItems(cost PopulationCost) []buildingBarStaffingItem {
	items := []buildingBarStaffingItem{}
	if cost.Apprentices > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: cost.Apprentices, Value: fmt.Sprintf("-%d", cost.Apprentices), Sprite: s.assetCatalog.Sprite.Icon.Apprentice, Cost: true,
		})
	}
	if cost.Soldiers > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: cost.Soldiers, Value: fmt.Sprintf("-%d", cost.Soldiers), Sprite: s.assetCatalog.Sprite.Icon.Soldier, Cost: true,
		})
	}
	if cost.Peasants > 0 {
		items = append(items, buildingBarStaffingItem{
			Count: cost.Peasants, Value: fmt.Sprintf("-%d", cost.Peasants), Sprite: s.assetCatalog.Sprite.Icon.Peasant, Cost: true,
		})
	}
	return items
}

// buildingBarStaffingWidth measures one inline inhabitant-requirement row.
func (s *State) buildingBarStaffingWidth(items []buildingBarStaffingItem) float64 {
	total := 0.0
	for i, item := range items {
		value := item.Value
		if value == "" {
			value = fmt.Sprintf("%d", item.Count)
		}
		valueWidth, _ := text.Measure(value, s.ui.costFace, s.ui.costFace.Size)
		total += buildingBarStaffIconSize + buildingBarStaffIconGap + valueWidth
		if i < len(items)-1 {
			total += buildingBarCostItemGap
		}
	}
	return total
}

// buildingBarCostWidth measures the full inline cost row width.
func (s *State) buildingBarCostWidth(items []buildingBarCostItem, hovered bool) float64 {
	total := 0.0
	face := s.buildingBarCostFace(hovered)
	for i, item := range items {
		width, _ := text.Measure(item.Value, face, face.Size)
		total += width
		if i < len(items)-1 {
			total += buildingBarCostItemGap
		}
	}
	return total
}

// buildingBarCostFace returns the regular or hover-emphasis cost face.
func (s *State) buildingBarCostFace(hovered bool) *text.GoTextFace {
	if hovered && s.ui.costBoldFace != nil {
		return s.ui.costBoldFace
	}
	return s.ui.costFace
}

// buildingBarCostItems returns non-zero costs in Wood, Stone, Metal order.
func buildingBarCostItems(cost Resources) []buildingBarCostItem {
	items := []buildingBarCostItem{}
	if cost.Wood > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Wood), Color: colors.resourceWood})
	}
	if cost.Stone > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Stone), Color: colors.resourceStone})
	}
	if cost.Metal > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Metal), Color: colors.resourceMetal})
	}
	return items
}
