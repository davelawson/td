package game

import (
	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type buildDragState struct {
	active  bool
	itemID  buildingBarItemID
	cursorX int
	cursorY int
}

// buildingBarVisible reports whether the current phase exposes construction controls.
func (s *State) buildingBarVisible() bool {
	return s.status.phase == phaseManagement && !s.raid.active && !s.raid.breached
}

// buildingBarBounds returns the screen-space building-bar rectangle.
func (s *State) buildingBarBounds() ui.Button[int] {
	return ui.BuildingBarBounds(topBarHeight, s.ui.height)
}

// buildingBarContains reports whether a point is inside the visible building bar.
func (s *State) buildingBarContains(x, y int) bool {
	return s.buildingBarVisible() && ui.BuildingBarContains(topBarHeight, s.ui.height, x, y)
}

// buildingBarModel adapts structure templates and host presentation state for the UI package.
func (s *State) buildingBarModel() ui.BuildingBarModel {
	actions := ui.BuildingBarActions()
	items := make([]ui.BuildingBarItem, 0, len(actions))
	for _, action := range actions {
		item, ok := s.buildingBarItem(action)
		if ok {
			items = append(items, item)
		}
	}
	return ui.BuildingBarModel{
		Items: items,
		Icons: ui.BuildingBarIcons{
			Apprentice: s.assetCatalog.Sprite.Icon.Apprentice,
			Soldier:    s.assetCatalog.Sprite.Icon.Soldier,
			Peasant:    s.assetCatalog.Sprite.Icon.Peasant,
		},
		SelectedCategory: s.ui.buildBarCategory,
		HoveredItem:      s.ui.buildBarHover,
		HoveredCategory:  s.ui.buildBarTabHover,
	}
}

// buildingBarItem adapts one structure template into UI-facing construction facts.
func (s *State) buildingBarItem(action buildingBarItemID) (ui.BuildingBarItem, bool) {
	template, ok := s.buildingTemplateForItemID(action)
	if !ok {
		return ui.BuildingBarItem{}, false
	}
	return ui.BuildingBarItem{
		Action:                        action,
		Name:                          template.Name,
		Description:                   template.Description,
		Sprite:                        template.Sprite,
		Cost:                          resourceAmounts(template.Cost),
		Staffing:                      staffingAmounts(template.Staffing),
		PopulationCost:                populationCostAmounts(template.PopulationCost),
		PopulationGrant:               populationGrantAmounts(template.PopulationGrant),
		ResourceYield:                 resourceAmounts(template.ResourceYield),
		RangeTiles:                    template.RangeTiles,
		Damage:                        template.Damage,
		FireIntervalSeconds:           template.FireIntervalSeconds,
		ProjectileSpeedTilesPerSecond: template.ProjectileSpeedTilesPerSecond,
		DamageAllEnemiesInTargetTile:  template.DamageAllEnemiesInTargetTile,
		Buildable:                     s.canConstructBuilding(action),
	}, true
}

// canConstructBuilding reports whether current resources and staff cover one action.
func (s *State) canConstructBuilding(action buildingBarItemID) bool {
	template, ok := s.buildingTemplateForItemID(action)
	return ok && s.canAffordBuildingCost(template.Cost) &&
		s.canPayPopulationCost(template.PopulationCost) &&
		s.canStaff(template.Staffing)
}

// updateBuildingBarHover records the item and category under the cursor.
func (s *State) updateBuildingBarHover(input Input) {
	if !s.buildingBarVisible() {
		s.clearBuildingBarHover()
		return
	}
	model := s.buildingBarModel()
	s.ui.buildBarHover = ui.BuildingBarItemIndexAt(topBarHeight, model, input.CursorX, input.CursorY)
	s.ui.buildBarTabHover = ui.BuildingBarCategoryAt(topBarHeight, input.CursorX, input.CursorY)
}

// clearBuildingBarHover resets transient building-bar hover state.
func (s *State) clearBuildingBarHover() {
	s.ui.buildBarHover = -1
	s.ui.buildBarTabHover = ui.BuildingBarNoCategory
}

// updateBuildDrag starts, tracks, completes, or cancels building drags.
func (s *State) updateBuildDrag(input Input) {
	if !s.buildingBarVisible() {
		s.buildDrag = buildDragState{}
		return
	}
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
	if category := ui.BuildingBarCategoryAt(topBarHeight, input.CursorX, input.CursorY); category != ui.BuildingBarNoCategory {
		s.ui.buildBarCategory = category
		s.ui.buildBarHover = -1
		return
	}

	item, ok := ui.BuildingBarItemAt(topBarHeight, s.buildingBarModel(), input.CursorX, input.CursorY)
	if !ok || !s.canConstructBuilding(item.Action) || !s.canBuildTowersNow() {
		return
	}
	s.buildDrag = buildDragState{
		active:  true,
		itemID:  item.Action,
		cursorX: input.CursorX,
		cursorY: input.CursorY,
	}
}

// canBuildTowersNow reports whether Management currently allows building placement.
func (s *State) canBuildTowersNow() bool {
	return s.buildingBarVisible()
}

// placeDraggedBuilding attempts to build the active dragged structure at a screen point.
func (s *State) placeDraggedBuilding(x, y int) {
	template, ok := s.draggedBuildingTemplate()
	if !ok || !s.canConstructBuilding(s.buildDrag.itemID) || !s.canBuildTowersNow() || s.buildDropBlockedByUI(x, y) {
		return
	}
	tile, ok := s.exploredTileAtScreenPosition(x, y)
	if !ok || !s.canBuildOnTile(tile) {
		return
	}
	feature, ok := buildingFeatureForItemID(s.buildDrag.itemID)
	if !ok {
		return
	}

	s.deductBuildingCost(template.Cost)
	s.deductPopulationCost(template.PopulationCost)
	s.reserveStaffing(template.Staffing)
	s.grantPopulation(template.PopulationGrant)
	if plot, ok := s.gameMap.plot(tile.Plot); ok {
		plot.Tiles[tile.Y][tile.X].Feature = feature
	}
}

// buildDropBlockedByUI reports whether a drop point is on screen-space game UI.
func (s *State) buildDropBlockedByUI(x, y int) bool {
	return s.buildingBarContains(x, y) ||
		s.nextRaidButtonContains(x, y) ||
		s.selectionPanelContains(x, y) ||
		s.exploreButtonContains(x, y)
}

// draggedBuildingTemplate returns the structure template attached to the cursor.
func (s *State) draggedBuildingTemplate() (StructureTemplate, bool) {
	if !s.buildDrag.active {
		return StructureTemplate{}, false
	}
	return s.buildingTemplateForItemID(s.buildDrag.itemID)
}

// exploredTileAtScreenPosition returns the explored Tile under a screen point.
func (s *State) exploredTileAtScreenPosition(x, y int) (tileCoordinate, bool) {
	viewport := s.sceneViewport()
	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		for tileY := 0; tileY < plotSize; tileY++ {
			for tileX := 0; tileX < plotSize; tileX++ {
				worldWest, worldNorth, worldW, worldH := plotTileWorldRect(plotCoord, tileX, tileY)
				rect := s.projectRect(viewport, worldWest, worldNorth, worldW, worldH)
				if rectContainsPoint(rect, x, y) {
					return tileCoordinate{Plot: plotCoord, X: tileX, Y: tileY}, true
				}
			}
		}
	}
	return tileCoordinate{}, false
}

// canBuildOnTile reports whether a Tile can receive a new structure.
func (s *State) canBuildOnTile(tile tileCoordinate) bool {
	if tile.X < 0 || tile.Y < 0 || tile.X >= plotSize || tile.Y >= plotSize {
		return false
	}
	plot, ok := s.gameMap.plot(tile.Plot)
	if !ok {
		return false
	}
	target := plot.Tiles[tile.Y][tile.X]
	return target.Terrain == terrainEmpty && target.Feature == featureNone
}

// drawBuildingBar delegates visible construction presentation to the UI package.
func (s *State) drawBuildingBar(screen *ebiten.Image) {
	if !s.buildingBarVisible() {
		return
	}
	ui.DrawBuildingBar(screen, s.ui.costFace, s.ui.costBoldFace, topBarHeight, s.ui.height, s.buildingBarModel())
}

// drawBuildingTooltip delegates the current hover tooltip to the UI package.
func (s *State) drawBuildingTooltip(screen *ebiten.Image) {
	if !s.buildingBarVisible() || s.buildDrag.active {
		return
	}
	ui.DrawBuildingTooltip(screen, s.ui.costFace, s.ui.costBoldFace, s.ui.width, s.ui.height, topBarHeight, s.buildingBarModel())
}

// drawBuildDrag delegates the active drag sprite to the UI package.
func (s *State) drawBuildDrag(screen *ebiten.Image) {
	if !s.buildingBarVisible() {
		return
	}
	template, ok := s.draggedBuildingTemplate()
	if !ok {
		return
	}
	ui.DrawBuildingDrag(screen, template.Sprite, s.buildDrag.cursorX, s.buildDrag.cursorY)
}
