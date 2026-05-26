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
	buildDragIconSize         = buildingBarItemSize / 2
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
	Cost   Resources
	Bounds ui.Button[int]
}

type buildingBarCostItem struct {
	Value string
	Color color.Color
}

type buildDragState struct {
	active    bool
	itemIndex int
	cursorX   int
	cursorY   int
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

// buildingBarItems returns the tower choices shown in the building bar.
func (s *State) buildingBarItems() []buildingBarItem {
	bar := s.buildingBarBounds()
	x := bar.X + (bar.W-buildingBarItemSize)/2
	startY := bar.Y + buildingBarPadding
	stepY := buildingBarItemSize + buildingBarCostGap + buildingBarCostTextHeight + buildingBarItemGap
	templates := []StructureTemplate{
		s.structureCatalog.BowTower,
		s.structureCatalog.FlameBoltTower,
		s.structureCatalog.CatapultTower,
	}
	items := make([]buildingBarItem, 0, len(templates))
	for i, template := range templates {
		y := startY + i*stepY
		items = append(items, buildingBarItem{
			Name:   template.Name,
			Sprite: template.Sprite,
			Cost:   template.Cost,
			Bounds: ui.Button[int]{
				Label:  template.Name,
				X:      x,
				Y:      y,
				W:      buildingBarItemSize,
				H:      buildingBarItemSize,
				Action: i,
			},
		})
	}
	return items
}

// buildingBarContains reports whether a point is inside the visual building bar.
func (s *State) buildingBarContains(x, y int) bool {
	return s.buildingBarBounds().Contains(x, y)
}

// canAffordBuildingCost reports whether current resources cover a structure cost.
func (s *State) canAffordBuildingCost(cost Resources) bool {
	return s.status.resources.wood >= cost.Wood &&
		s.status.resources.stone >= cost.Stone &&
		s.status.resources.metal >= cost.Metal
}

// updateBuildingBarHover records which tower icon, if any, is under the cursor.
func (s *State) updateBuildingBarHover(input Input) {
	s.ui.buildBarHover = s.buildingBarItemIndexAt(input.CursorX, input.CursorY)
}

// updateBuildDrag starts, tracks, completes, or cancels tower build drags.
func (s *State) updateBuildDrag(input Input) {
	if s.buildDrag.active {
		s.buildDrag.cursorX = input.CursorX
		s.buildDrag.cursorY = input.CursorY
		if input.Released {
			s.placeDraggedBuilding(input.CursorX, input.CursorY)
			s.buildDrag = buildDragState{}
			return
		}
		if !input.MouseDown {
			s.buildDrag = buildDragState{}
		}
		return
	}

	if !input.Clicked {
		return
	}

	index := s.buildingBarItemIndexAt(input.CursorX, input.CursorY)
	if index < 0 {
		return
	}
	item := s.buildingBarItems()[index]
	if !s.canAffordBuildingCost(item.Cost) || !s.canBuildTowersNow() {
		return
	}
	s.buildDrag = buildDragState{
		active:    true,
		itemIndex: index,
		cursorX:   input.CursorX,
		cursorY:   input.CursorY,
	}
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

// canBuildTowersNow reports whether the current game phase allows tower placement.
func (s *State) canBuildTowersNow() bool {
	return s.status.phase == phaseCalm && !s.raid.active && !s.raid.breached
}

// placeDraggedBuilding attempts to build the active dragged tower at a screen point.
func (s *State) placeDraggedBuilding(x, y int) {
	item, ok := s.draggedBuildingItem()
	if !ok || !s.canAffordBuildingCost(item.Cost) || !s.canBuildTowersNow() || s.buildDropBlockedByUI(x, y) {
		return
	}
	tile, ok := s.homePlotTileAtScreenPosition(x, y)
	if !ok || !s.canBuildOnTile(tile) {
		return
	}
	feature, ok := buildingFeatureForItemIndex(s.buildDrag.itemIndex)
	if !ok {
		return
	}

	s.deductBuildingCost(item.Cost)
	s.gameMap.Home.Tiles[tile.Y][tile.X].Feature = feature
}

// buildDropBlockedByUI reports whether a drop point is on screen-space game UI.
func (s *State) buildDropBlockedByUI(x, y int) bool {
	return s.buildingBarContains(x, y) ||
		s.nextRaidButtonContains(x, y) ||
		s.selectionPanelContains(x, y)
}

// draggedBuildingItem returns the building-bar item currently attached to the cursor.
func (s *State) draggedBuildingItem() (buildingBarItem, bool) {
	if !s.buildDrag.active {
		return buildingBarItem{}, false
	}
	items := s.buildingBarItems()
	if s.buildDrag.itemIndex < 0 || s.buildDrag.itemIndex >= len(items) {
		return buildingBarItem{}, false
	}
	return items[s.buildDrag.itemIndex], true
}

// homePlotTileAtScreenPosition returns the home Plot Tile under a screen point.
func (s *State) homePlotTileAtScreenPosition(x, y int) (tileCoordinate, bool) {
	viewport := s.sceneViewport()
	for tileY := 0; tileY < plotSize; tileY++ {
		for tileX := 0; tileX < plotSize; tileX++ {
			worldWest, worldNorth, worldW, worldH := tileWorldRect(tileX, tileY)
			rect := s.projectRect(viewport, worldWest, worldNorth, worldW, worldH)
			if rectContainsPoint(rect, x, y) {
				return tileCoordinate{X: tileX, Y: tileY}, true
			}
		}
	}
	return tileCoordinate{}, false
}

// canBuildOnTile reports whether a Tile can receive a new tower.
func (s *State) canBuildOnTile(tile tileCoordinate) bool {
	if tile.X < 0 || tile.Y < 0 || tile.X >= plotSize || tile.Y >= plotSize {
		return false
	}
	target := s.gameMap.Home.Tiles[tile.Y][tile.X]
	return target.Terrain == terrainEmpty && target.Feature == featureNone
}

// buildingFeatureForItemIndex maps building-bar choices to placed Tile features.
func buildingFeatureForItemIndex(index int) (tileFeature, bool) {
	switch index {
	case 0:
		return featureBowTower, true
	case 1:
		return featureFlameBoltTower, true
	case 2:
		return featureCatapultTower, true
	default:
		return featureNone, false
	}
}

// deductBuildingCost spends the resources required to build a tower.
func (s *State) deductBuildingCost(cost Resources) {
	s.status.resources.wood -= cost.Wood
	s.status.resources.stone -= cost.Stone
	s.status.resources.metal -= cost.Metal
}

// drawBuildingBar renders the tower picker at the left edge of the scene.
func (s *State) drawBuildingBar(screen *ebiten.Image) {
	bar := s.buildingBarBounds()
	if bar.H <= 0 {
		return
	}

	vector.FillRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), colors.selectionPanel, false)
	vector.StrokeLine(screen, float32(bar.X+bar.W-2), float32(bar.Y), float32(bar.X+bar.W-2), float32(bar.Y+bar.H), 3, colors.fieldEdge, false)

	for i, item := range s.buildingBarItems() {
		s.drawBuildingBarItem(screen, item, s.buildingBarItemHighlighted(i, item))
	}
}

// drawBuildDrag renders the active tower icon attached to the cursor.
func (s *State) drawBuildDrag(screen *ebiten.Image) {
	item, ok := s.draggedBuildingItem()
	if !ok || item.Sprite == nil {
		return
	}
	spriteWidth := float64(item.Sprite.Bounds().Dx())
	spriteHeight := float64(item.Sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		return
	}

	scale := float64(buildDragIconSize) / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(s.buildDrag.cursorX)-spriteWidth*scale/2,
		float64(s.buildDrag.cursorY)-spriteHeight*scale/2,
	)
	screen.DrawImage(item.Sprite, options)
}

// buildingBarItemHighlighted reports whether an item should receive hover emphasis.
func (s *State) buildingBarItemHighlighted(index int, item buildingBarItem) bool {
	return s.ui.buildBarHover == index && s.canAffordBuildingCost(item.Cost)
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
