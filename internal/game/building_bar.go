package game

import (
	"image/color"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	buildingBarWidth           = 260
	buildingBarPadding         = 16
	buildingBarTabHeight       = 28
	buildingBarTabGap          = 6
	buildingBarTabBottomGap    = 14
	buildingBarItemSize        = 64
	buildDragIconSize          = buildingBarItemSize / 2
	buildingBarMetadataGap     = 12
	buildingBarCostOffsetY     = 11
	buildingBarStaffingOffsetY = 35
	buildingBarStaffIconSize   = 14
	buildingBarStaffIconGap    = 2
	buildingBarCostItemGap     = 5
	buildingBarItemGap         = 12
	buildingBarSpriteInset     = 8
)

var buildingBarCostShadow = color.RGBA{R: 8, G: 10, B: 8, A: 220}

type buildingBarItem struct {
	ID              buildingBarItemID
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
	active  bool
	itemID  buildingBarItemID
	cursorX int
	cursorY int
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
	x := bar.X + buildingBarPadding
	startY := bar.Y + buildingBarPadding + buildingBarTabsHeight() + buildingBarTabBottomGap
	stepY := buildingBarItemSize + buildingBarItemGap
	ids := buildingBarItemIDsForCategory(s.ui.buildBarCategory)
	items := make([]buildingBarItem, 0, len(ids))
	for i, id := range ids {
		template, ok := s.buildingTemplateForItemID(id)
		if !ok {
			continue
		}
		y := startY + i*stepY
		items = append(items, buildingBarItem{
			ID:              id,
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

// buildingBarTabs returns the category tabs shown at the top of the building bar.
func (s *State) buildingBarTabs() []buildingBarTab {
	bar := s.buildingBarBounds()
	tabX := bar.X + buildingBarTabGap
	tabY := bar.Y + buildingBarPadding
	tabW := bar.W - buildingBarTabGap*2
	categories := buildingBarCategories()
	tabs := make([]buildingBarTab, 0, len(categories))
	for i, category := range categories {
		tabs = append(tabs, buildingBarTab{
			Category: category,
			Label:    buildingBarCategoryLabel(category),
			Bounds: ui.Button[int]{
				Label:  buildingBarCategoryLabel(category),
				X:      tabX,
				Y:      tabY + i*(buildingBarTabHeight+buildingBarTabGap),
				W:      tabW,
				H:      buildingBarTabHeight,
				Action: int(category),
			},
		})
	}
	return tabs
}

// buildingBarTabsHeight returns the vertical space reserved for category tabs.
func buildingBarTabsHeight() int {
	return len(buildingBarCategories())*buildingBarTabHeight +
		(len(buildingBarCategories())-1)*buildingBarTabGap
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
	s.ui.buildBarTabHover = s.buildingBarTabAt(input.CursorX, input.CursorY)
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

	if category := s.buildingBarTabAt(input.CursorX, input.CursorY); category != buildingBarNoCategory {
		s.ui.buildBarCategory = category
		s.ui.buildBarHover = -1
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
		active:  true,
		itemID:  item.ID,
		cursorX: input.CursorX,
		cursorY: input.CursorY,
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

// buildingBarItemByID returns the visible item with the requested stable ID.
func (s *State) buildingBarItemByID(id buildingBarItemID) (buildingBarItem, bool) {
	for _, item := range s.buildingBarItems() {
		if item.ID == id {
			return item, true
		}
	}
	return buildingBarItem{}, false
}

// buildingBarTabAt returns the category tab under a point, or no category.
func (s *State) buildingBarTabAt(x, y int) buildingBarCategory {
	for _, tab := range s.buildingBarTabs() {
		if tab.Bounds.Contains(x, y) {
			return tab.Category
		}
	}
	return buildingBarNoCategory
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
	feature, ok := buildingFeatureForItemID(item.ID)
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
	template, ok := s.buildingTemplateForItemID(s.buildDrag.itemID)
	if !ok {
		return buildingBarItem{}, false
	}
	return buildingBarItem{
		ID:              s.buildDrag.itemID,
		Name:            template.Name,
		Sprite:          template.Sprite,
		Cost:            template.Cost,
		Staffing:        template.Staffing,
		PopulationCost:  template.PopulationCost,
		PopulationGrant: template.PopulationGrant,
	}, true
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

// deductBuildingCost spends the resources required to build a structure.
func (s *State) deductBuildingCost(cost Resources) {
	s.status.resources.wood -= cost.Wood
	s.status.resources.stone -= cost.Stone
	s.status.resources.metal -= cost.Metal
}

// drawBuildingBar renders the building picker at the left edge of the scene.
func (s *State) drawBuildingBar(screen *ebiten.Image) {
	bar := s.buildingBarBounds()
	if bar.H <= 0 {
		return
	}

	vector.FillRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), colors.selectionPanel, false)
	vector.StrokeLine(screen, float32(bar.X+bar.W-2), float32(bar.Y), float32(bar.X+bar.W-2), float32(bar.Y+bar.H), 3, colors.fieldEdge, false)

	for _, tab := range s.buildingBarTabs() {
		s.drawBuildingBarTab(screen, tab)
	}

	for i, item := range s.buildingBarItems() {
		s.drawBuildingBarItem(screen, item, s.buildingBarItemHighlighted(i, item))
	}
}

// drawBuildingBarTab renders one build-category tab.
func (s *State) drawBuildingBarTab(screen *ebiten.Image, tab buildingBarTab) {
	bounds := tab.Bounds
	selected := s.ui.buildBarCategory == tab.Category
	hovered := s.ui.buildBarTabHover == tab.Category
	fill := colors.plotBackdrop
	if selected {
		fill = colors.fieldEdge
	}
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), fill, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 1, colors.fieldEdge, false)

	textColor := colors.text
	if hovered && !selected {
		textColor = colors.pause
	}
	width, height := text.Measure(tab.Label, s.ui.costFace, s.ui.costFace.Size)
	x := float64(bounds.X) + (float64(bounds.W)-width)/2
	y := float64(bounds.Y) + (float64(bounds.H)-height)/2 - 1
	ui.DrawText(screen, tab.Label, s.ui.costFace, x, y, textColor)
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

// drawBuildingBarItem renders one building icon slot and its right-side values.
func (s *State) drawBuildingBarItem(screen *ebiten.Image, item buildingBarItem, hovered bool) {
	bounds := item.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), colors.plotBackdrop, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, s.buildingBarItemOutlineColor(item), false)

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
	options.ColorScale.Scale(1, 1, 1, s.buildingBarIconAlpha(item))
	if hovered {
		brightenDrawOptions(options)
	}
	screen.DrawImage(item.Sprite, options)
	s.drawBuildingBarCost(screen, item, hovered)
	s.drawBuildingBarPopulationMetadata(screen, item)
}

// buildingBarIconAlpha returns the icon opacity for current construction capacity.
func (s *State) buildingBarIconAlpha(item buildingBarItem) float32 {
	if s.canConstructBuilding(item) {
		return 1
	}
	return 0.70
}

// buildingBarItemOutlineColor returns the icon slot outline color for construction capacity.
func (s *State) buildingBarItemOutlineColor(item buildingBarItem) color.Color {
	if s.canConstructBuilding(item) {
		return colors.buildable
	}
	return colors.buildBlocked
}
