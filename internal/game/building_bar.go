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
	buildingBarStaffingGap    = 4
	buildingBarStaffingHeight = 18
	buildingBarStaffIconSize  = 14
	buildingBarStaffIconGap   = 2
	buildingBarCostItemGap    = 5
	buildingBarItemGap        = 12
	buildingBarSpriteInset    = 8
)

var buildingBarCostShadow = color.RGBA{R: 8, G: 10, B: 8, A: 220}

const (
	buildingBarHouseIndex = iota
	buildingBarBarracksIndex
	buildingBarWoodcutterIndex
	buildingBarStoneQuarryIndex
	buildingBarIronMineIndex
	buildingBarBowTowerIndex
	buildingBarFlameBoltTowerIndex
	buildingBarCatapultTowerIndex
)

type buildingBarItem struct {
	Name            string
	Sprite          *ebiten.Image
	Cost            Resources
	Staffing        StaffingRequirements
	PopulationCost  PopulationCost
	PopulationGrant PopulationGrant
	Bounds          ui.Button[int]
}

type buildingBarCostItem struct {
	Value string
	Color color.Color
}

type buildingBarStaffingItem struct {
	Count  int
	Value  string
	Sprite *ebiten.Image
	Cost   bool
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

// buildingBarItems returns the structure choices shown in the building bar.
func (s *State) buildingBarItems() []buildingBarItem {
	bar := s.buildingBarBounds()
	x := bar.X + (bar.W-buildingBarItemSize)/2
	startY := bar.Y + buildingBarPadding
	stepY := buildingBarItemSize +
		buildingBarCostGap + buildingBarCostTextHeight +
		buildingBarStaffingGap + buildingBarStaffingHeight +
		buildingBarItemGap
	templates := []StructureTemplate{
		s.structureCatalog.House,
		s.structureCatalog.Barracks,
		s.structureCatalog.Woodcutter,
		s.structureCatalog.StoneQuarry,
		s.structureCatalog.IronMine,
		s.structureCatalog.BowTower,
		s.structureCatalog.FlameBoltTower,
		s.structureCatalog.CatapultTower,
	}
	items := make([]buildingBarItem, 0, len(templates))
	for i, template := range templates {
		y := startY + i*stepY
		items = append(items, buildingBarItem{
			Name:            template.Name,
			Sprite:          template.Sprite,
			Cost:            template.Cost,
			Staffing:        template.Staffing,
			PopulationCost:  template.PopulationCost,
			PopulationGrant: template.PopulationGrant,
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

// canConstructBuilding reports whether current resources and staff cover one item.
func (s *State) canConstructBuilding(item buildingBarItem) bool {
	return s.canAffordBuildingCost(item.Cost) &&
		s.canPayPopulationCost(item.PopulationCost) &&
		s.canStaff(item.Staffing)
}

// updateBuildingBarHover records which tower icon, if any, is under the cursor.
func (s *State) updateBuildingBarHover(input Input) {
	s.ui.buildBarHover = s.buildingBarItemIndexAt(input.CursorX, input.CursorY)
}

// updateBuildDrag starts, tracks, completes, or cancels building drags.
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
	if !s.canConstructBuilding(item) || !s.canBuildTowersNow() {
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

// canBuildTowersNow reports whether the current game phase allows building placement.
func (s *State) canBuildTowersNow() bool {
	return s.status.phase == phaseCalm && !s.raid.active && !s.raid.breached
}

// placeDraggedBuilding attempts to build the active dragged structure at a screen point.
func (s *State) placeDraggedBuilding(x, y int) {
	item, ok := s.draggedBuildingItem()
	if !ok || !s.canConstructBuilding(item) || !s.canBuildTowersNow() || s.buildDropBlockedByUI(x, y) {
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
	s.deductPopulationCost(item.PopulationCost)
	s.reserveStaffing(item.Staffing)
	s.grantPopulation(item.PopulationGrant)
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
	case buildingBarHouseIndex:
		return featureHouse, true
	case buildingBarBarracksIndex:
		return featureBarracks, true
	case buildingBarWoodcutterIndex:
		return featureWoodcutter, true
	case buildingBarStoneQuarryIndex:
		return featureStoneQuarry, true
	case buildingBarIronMineIndex:
		return featureIronMine, true
	case buildingBarBowTowerIndex:
		return featureBowTower, true
	case buildingBarFlameBoltTowerIndex:
		return featureFlameBoltTower, true
	case buildingBarCatapultTowerIndex:
		return featureCatapultTower, true
	default:
		return featureNone, false
	}
}

// deductBuildingCost spends the resources required to build a structure.
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

// drawBuildDrag renders the active building icon attached to the cursor.
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
	return s.ui.buildBarHover == index && s.canConstructBuilding(item)
}

// drawBuildingBarItem renders one tower icon slot.
func (s *State) drawBuildingBarItem(screen *ebiten.Image, item buildingBarItem, hovered bool) {
	bounds := item.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), colors.plotBackdrop, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, colors.fieldEdge, false)

	if item.Sprite == nil {
		s.drawBuildingBarCost(screen, item, hovered)
		s.drawBuildingBarPopulationMetadata(screen, item)
		return
	}
	spriteWidth := float64(item.Sprite.Bounds().Dx())
	spriteHeight := float64(item.Sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		s.drawBuildingBarCost(screen, item, hovered)
		s.drawBuildingBarPopulationMetadata(screen, item)
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
	s.drawBuildingBarPopulationMetadata(screen, item)
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

// drawBuildingBarPopulationMetadata renders staffing requirements or population grants.
func (s *State) drawBuildingBarPopulationMetadata(screen *ebiten.Image, item buildingBarItem) {
	metadataItems := s.buildingBarPopulationMetadataItems(item)
	if len(metadataItems) == 0 {
		return
	}

	totalWidth := s.buildingBarStaffingWidth(metadataItems)
	x := float64(item.Bounds.X) + (float64(item.Bounds.W)-totalWidth)/2
	y := float64(item.Bounds.Y + item.Bounds.H + buildingBarCostGap + buildingBarCostTextHeight + buildingBarStaffingGap)
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
