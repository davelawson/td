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
	buildingBarWidth          = 96
	buildingBarPadding        = 16
	buildingBarItemSize       = 64
	buildingBarCostGap        = 4
	buildingBarCostTextHeight = 18
	buildingBarCostItemGap    = 6
	buildingBarItemGap        = 12
	buildingBarSpriteInset    = 8
)

var buildingBarCostShadow = color.RGBA{R: 8, G: 10, B: 8, A: 220}

type buildingBarItem struct {
	Name   string
	Sprite *ebiten.Image
	Cost   ResourceCost
	Bounds ui.Button[int]
}

type buildingBarCostItem struct {
	Value string
	Color color.Color
}

// buildingBarBounds returns the screen-space building bar rectangle.
func (s *State) buildingBarBounds() ui.Button[int] {
	return ui.Button[int]{
		X: 0,
		Y: topBarHeight,
		W: buildingBarWidth,
		H: s.ui.height - topBarHeight,
	}
}

// buildingBarItems returns the visual-only tower choices shown in the building bar.
func (s *State) buildingBarItems() []buildingBarItem {
	bar := s.buildingBarBounds()
	x := bar.X + (bar.W-buildingBarItemSize)/2
	y := bar.Y + buildingBarPadding
	nextY := y + buildingBarItemSize + buildingBarCostGap + buildingBarCostTextHeight + buildingBarItemGap
	return []buildingBarItem{
		{
			Name:   s.structureCatalog.BowTower.Name,
			Sprite: s.structureCatalog.BowTower.Sprite,
			Cost:   s.structureCatalog.BowTower.Cost,
			Bounds: ui.Button[int]{
				Label:  s.structureCatalog.BowTower.Name,
				X:      x,
				Y:      y,
				W:      buildingBarItemSize,
				H:      buildingBarItemSize,
				Action: 0,
			},
		},
		{
			Name:   s.structureCatalog.FlameBoltTower.Name,
			Sprite: s.structureCatalog.FlameBoltTower.Sprite,
			Cost:   s.structureCatalog.FlameBoltTower.Cost,
			Bounds: ui.Button[int]{
				Label:  s.structureCatalog.FlameBoltTower.Name,
				X:      x,
				Y:      nextY,
				W:      buildingBarItemSize,
				H:      buildingBarItemSize,
				Action: 1,
			},
		},
	}
}

// buildingBarContains reports whether a point is inside the visual building bar.
func (s *State) buildingBarContains(x, y int) bool {
	return s.buildingBarBounds().Contains(x, y)
}

// updateBuildingBarHover records which tower icon, if any, is under the cursor.
func (s *State) updateBuildingBarHover(input Input) {
	s.ui.buildBarHover = s.buildingBarItemIndexAt(input.CursorX, input.CursorY)
}

// buildingBarItemIndexAt returns the tower icon index at a point, or -1.
func (s *State) buildingBarItemIndexAt(x, y int) int {
	for i, item := range s.buildingBarItems() {
		if item.Bounds.Contains(x, y) {
			return i
		}
	}
	return -1
}

// drawBuildingBar renders the visual-only tower picker at the left edge of the scene.
func (s *State) drawBuildingBar(screen *ebiten.Image) {
	bar := s.buildingBarBounds()
	if bar.H <= 0 {
		return
	}

	vector.FillRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), colors.selectionPanel, false)
	vector.StrokeLine(screen, float32(bar.X+bar.W-2), float32(bar.Y), float32(bar.X+bar.W-2), float32(bar.Y+bar.H), 3, colors.fieldEdge, false)

	for i, item := range s.buildingBarItems() {
		s.drawBuildingBarItem(screen, item, s.ui.buildBarHover == i)
	}
}

// drawBuildingBarItem renders one tower icon slot.
func (s *State) drawBuildingBarItem(screen *ebiten.Image, item buildingBarItem, hovered bool) {
	bounds := item.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), colors.plotBackdrop, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, colors.fieldEdge, false)

	if item.Sprite == nil {
		s.drawBuildingBarCost(screen, item, hovered)
		return
	}
	spriteWidth := float64(item.Sprite.Bounds().Dx())
	spriteHeight := float64(item.Sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		s.drawBuildingBarCost(screen, item, hovered)
		return
	}

	targetSize := float64(bounds.W - buildingBarSpriteInset*2)
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(bounds.X)+(float64(bounds.W)-spriteWidth*scale)/2,
		float64(bounds.Y)+(float64(bounds.H)-spriteHeight*scale)/2,
	)
	if hovered {
		brightenDrawOptions(options)
	}
	screen.DrawImage(item.Sprite, options)
	s.drawBuildingBarCost(screen, item, hovered)
}

// drawBuildingBarCost renders non-zero resource costs below one tower icon.
func (s *State) drawBuildingBarCost(screen *ebiten.Image, item buildingBarItem, hovered bool) {
	costItems := buildingBarCostItems(item.Cost)
	if len(costItems) == 0 {
		return
	}

	face := s.buildingBarCostFace(hovered)
	totalWidth := s.buildingBarCostWidth(costItems, hovered)
	x := float64(item.Bounds.X) + (float64(item.Bounds.W)-totalWidth)/2
	y := float64(item.Bounds.Y + item.Bounds.H + buildingBarCostGap)
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
func buildingBarCostItems(cost ResourceCost) []buildingBarCostItem {
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
