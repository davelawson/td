package game

type selectedItemKind int

const (
	selectedItemNone selectedItemKind = iota
	selectedItemStructure
	selectedItemRaider
	selectedItemTerrain
)

type selectedItem struct {
	kind     selectedItemKind
	tile     tileCoordinate
	raiderID int
}

// updateSelection applies one left-click selection update for map objects.
func (s *State) updateSelection(input Input) {
	if !input.Clicked ||
		s.nextRaidButtonContains(input.CursorX, input.CursorY) ||
		s.buildingBarContains(input.CursorX, input.CursorY) ||
		s.selectionPanelContains(input.CursorX, input.CursorY) ||
		s.exploreButtonContains(input.CursorX, input.CursorY) {
		return
	}

	if raiderID, ok := s.raiderAtScreenPosition(input.CursorX, input.CursorY); ok {
		s.selection = selectedItem{
			kind:     selectedItemRaider,
			raiderID: raiderID,
		}
		return
	}

	if tile, ok := s.structureAtScreenPosition(input.CursorX, input.CursorY); ok {
		s.selection = selectedItem{
			kind: selectedItemStructure,
			tile: tile,
		}
		return
	}

	if tile, ok := s.terrainAtScreenPosition(input.CursorX, input.CursorY); ok {
		s.selection = selectedItem{
			kind: selectedItemTerrain,
			tile: tile,
		}
		return
	}

	s.selection = selectedItem{}
}

// raiderAtScreenPosition returns the topmost active raider sprite at a screen point.
func (s *State) raiderAtScreenPosition(x, y int) (int, bool) {
	viewport := s.sceneViewport()
	for i := len(s.raid.enemies) - 1; i >= 0; i-- {
		enemy := s.raid.enemies[i]
		if rectContainsPoint(s.raidEnemyProjectedRect(viewport, enemy), x, y) {
			return enemy.id, true
		}
	}
	return 0, false
}

// structureAtScreenPosition returns the explored structure tile at a screen point.
func (s *State) structureAtScreenPosition(x, y int) (tileCoordinate, bool) {
	viewport := s.sceneViewport()
	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		plot, ok := s.gameMap.plot(plotCoord)
		if !ok {
			continue
		}
		for tileY := 0; tileY < plotSize; tileY++ {
			for tileX := 0; tileX < plotSize; tileX++ {
				if plot.Tiles[tileY][tileX].Feature == featureNone {
					continue
				}
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

// terrainAtScreenPosition returns a selectable Tree or Boulder Tile at a screen point.
func (s *State) terrainAtScreenPosition(x, y int) (tileCoordinate, bool) {
	tile, ok := s.exploredTileAtScreenPosition(x, y)
	if !ok {
		return tileCoordinate{}, false
	}
	plot, ok := s.gameMap.plot(tile.Plot)
	if !ok || plot.Tiles[tile.Y][tile.X].Feature != featureNone {
		return tileCoordinate{}, false
	}
	switch plot.Tiles[tile.Y][tile.X].Terrain {
	case terrainTree, terrainBoulder:
		return tile, true
	default:
		return tileCoordinate{}, false
	}
}

// selectedStructure reports whether a structure tile is currently selected.
func (s *State) selectedStructure(tile tileCoordinate) bool {
	return s.selection.kind == selectedItemStructure && s.selection.tile == tile
}

// selectedRaider reports whether a raider is currently selected.
func (s *State) selectedRaider(id int) bool {
	return s.selection.kind == selectedItemRaider && s.selection.raiderID == id
}

// selectedTerrain reports whether a terrain Tile is currently selected.
func (s *State) selectedTerrain(tile tileCoordinate) bool {
	return s.selection.kind == selectedItemTerrain && s.selection.tile == tile
}

// clearMissingSelectedRaider clears raider selection after the selected raider leaves active state.
func (s *State) clearMissingSelectedRaider() {
	if s.selection.kind != selectedItemRaider {
		return
	}
	for _, enemy := range s.raid.enemies {
		if enemy.id == s.selection.raiderID && enemy.health > 0 {
			return
		}
	}
	s.selection = selectedItem{}
}

// raidEnemyProjectedRect returns the current screen-space bounds used to draw a raider.
func (s *State) raidEnemyProjectedRect(viewport sceneViewport, enemy raidEnemy) projectedRect {
	if enemy.template != nil && enemy.template.Sprite != nil {
		size := raidEnemySpriteSize / plotBaseTileSize
		return s.projectRect(
			viewport,
			enemy.position.X-size/2,
			enemy.position.Y+size/2,
			size,
			size,
		)
	}

	radius := raidEnemyRadius / plotBaseTileSize
	return s.projectRect(
		viewport,
		enemy.position.X-radius,
		enemy.position.Y+radius,
		radius*2,
		radius*2,
	)
}

// rectContainsPoint reports whether a projected rectangle contains a screen point.
func rectContainsPoint(rect projectedRect, x, y int) bool {
	screenX := float32(x)
	screenY := float32(y)
	return screenX >= rect.x && screenX <= rect.x+rect.w && screenY >= rect.y && screenY <= rect.y+rect.h
}
