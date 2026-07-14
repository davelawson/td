package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	buildingBarMetadataGap     = 12
	buildingBarCostOffsetY     = 11
	buildingBarStaffingOffsetY = 35
	buildingBarStaffIconSize   = 14
	buildingBarStaffIconGap    = 2
	buildingBarCostItemGap     = 5
)

var buildingBarCostShadow = color.RGBA{R: 8, G: 10, B: 8, A: 220}

type buildingBarCostItem struct {
	Value string
	Color color.Color
}

type buildingBarPopulationItem struct {
	Count  int
	Value  string
	Sprite *ebiten.Image
}

// drawBuildingBarCost renders non-zero resource costs beside one icon.
func drawBuildingBarCost(screen *ebiten.Image, regularFace, boldFace *text.GoTextFace, item buildingBarLayoutItem, hovered bool) {
	costItems := buildingBarCostItems(item.Cost)
	if len(costItems) == 0 {
		return
	}

	face := buildingBarCostFace(regularFace, boldFace, hovered)
	x := float64(buildingBarMetadataX(item))
	y := float64(item.Bounds.Y + buildingBarCostOffsetY)
	for index, costItem := range costItems {
		width, _ := text.Measure(costItem.Value, face, face.Size)
		if hovered {
			DrawText(screen, costItem.Value, face, x+1, y+1, buildingBarCostShadow)
			DrawText(screen, costItem.Value, face, x-1, y+1, buildingBarCostShadow)
		}
		DrawText(screen, costItem.Value, face, x, y, costItem.Color)
		x += width
		if index < len(costItems)-1 {
			x += buildingBarCostItemGap
		}
	}
}

// drawBuildingBarPopulationMetadata renders staffing, cost, or grant facts beside one icon.
func drawBuildingBarPopulationMetadata(screen *ebiten.Image, face *text.GoTextFace, icons BuildingBarIcons, item buildingBarLayoutItem) {
	items := buildingBarPopulationMetadataItems(icons, item.BuildingBarItem)
	if len(items) == 0 {
		return
	}

	x := float64(buildingBarMetadataX(item))
	y := float64(item.Bounds.Y + buildingBarStaffingOffsetY)
	for index, metadataItem := range items {
		if metadataItem.Sprite != nil {
			spriteWidth := float64(metadataItem.Sprite.Bounds().Dx())
			spriteHeight := float64(metadataItem.Sprite.Bounds().Dy())
			if spriteWidth > 0 && spriteHeight > 0 {
				options := &ebiten.DrawImageOptions{}
				scale := float64(buildingBarStaffIconSize) / spriteWidth
				options.GeoM.Scale(scale, scale)
				options.GeoM.Translate(x, y)
				screen.DrawImage(metadataItem.Sprite, options)
			}
		}
		x += buildingBarStaffIconSize + buildingBarStaffIconGap
		value := metadataItem.Value
		if value == "" {
			value = fmt.Sprintf("%d", metadataItem.Count)
		}
		DrawText(screen, value, face, x, y-1, Parchment)
		valueWidth, _ := text.Measure(value, face, face.Size)
		x += valueWidth
		if index < len(items)-1 {
			x += buildingBarCostItemGap
		}
	}
}

// buildingBarMetadataX returns the x coordinate where item values begin.
func buildingBarMetadataX(item buildingBarLayoutItem) int {
	return item.Bounds.X + item.Bounds.W + buildingBarMetadataGap
}

// buildingBarMetadataRight returns the right edge available to item values.
func buildingBarMetadataRight() int {
	return BuildingBarWidth - buildingBarPadding
}

// buildingBarPopulationMetadataItems returns the population row shown for an item.
func buildingBarPopulationMetadataItems(icons BuildingBarIcons, item BuildingBarItem) []buildingBarPopulationItem {
	staffingItems := buildingBarPopulationItems(icons, item.Staffing, "")
	if len(staffingItems) > 0 {
		return staffingItems
	}
	items := buildingBarPopulationItems(icons, item.PopulationCost, "-")
	return append(items, buildingBarPopulationItems(icons, item.PopulationGrant, "+")...)
}

// buildingBarPopulationItems returns non-zero roles in Apprentice, Soldier, Peasant order.
func buildingBarPopulationItems(icons BuildingBarIcons, amounts PopulationAmounts, prefix string) []buildingBarPopulationItem {
	items := []buildingBarPopulationItem{}
	appendItem := func(count int, sprite *ebiten.Image) {
		if count <= 0 {
			return
		}
		value := ""
		if prefix != "" {
			value = fmt.Sprintf("%s%d", prefix, count)
		}
		items = append(items, buildingBarPopulationItem{Count: count, Value: value, Sprite: sprite})
	}
	appendItem(amounts.Apprentices, icons.Apprentice)
	appendItem(amounts.Soldiers, icons.Soldier)
	appendItem(amounts.Peasants, icons.Peasant)
	return items
}

// buildingBarCostItems returns non-zero costs in Wood, Stone, Metal order.
func buildingBarCostItems(cost ResourceAmounts) []buildingBarCostItem {
	items := []buildingBarCostItem{}
	if cost.Wood > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Wood), Color: ResourceWood})
	}
	if cost.Stone > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Stone), Color: ResourceStone})
	}
	if cost.Metal > 0 {
		items = append(items, buildingBarCostItem{Value: fmt.Sprintf("%d", cost.Metal), Color: ResourceMetal})
	}
	return items
}

// buildingBarCostFace returns the regular or hover-emphasis face.
func buildingBarCostFace(regularFace, boldFace *text.GoTextFace, hovered bool) *text.GoTextFace {
	if hovered && boldFace != nil {
		return boldFace
	}
	return regularFace
}

// buildingBarCostWidth measures the full inline resource-cost row.
func buildingBarCostWidth(face *text.GoTextFace, items []buildingBarCostItem) float64 {
	total := 0.0
	for index, item := range items {
		width, _ := text.Measure(item.Value, face, face.Size)
		total += width
		if index < len(items)-1 {
			total += buildingBarCostItemGap
		}
	}
	return total
}
