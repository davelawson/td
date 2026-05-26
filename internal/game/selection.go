package game

type selectedItemKind int

const (
	selectedItemNone selectedItemKind = iota
	selectedItemStructure
	selectedItemRaider
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
		s.selectionPanelContains(input.CursorX, input.CursorY) {
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

// structureAtScreenPosition returns the home Plot structure tile at a screen point.
func (s *State) structureAtScreenPosition(x, y int) (tileCoordinate, bool) {
	viewport := s.sceneViewport()
	for tileY := 0; tileY < plotSize; tileY++ {
		for tileX := 0; tileX < plotSize; tileX++ {
			if s.gameMap.Home.Tiles[tileY][tileX].Feature == featureNone {
				continue
			}
			worldWest, worldNorth, worldW, worldH := tileWorldRect(tileX, tileY)
			rect := s.projectRect(viewport, worldWest, worldNorth, worldW, worldH)
			if rectContainsPoint(rect, x, y) {
				return tileCoordinate{X: tileX, Y: tileY}, true
			}
		}
	}
	return tileCoordinate{}, false
}

// selectedStructure reports whether a structure tile is currently selected.
func (s *State) selectedStructure(x, y int) bool {
	return s.selection.kind == selectedItemStructure && s.selection.tile.X == x && s.selection.tile.Y == y
}

// selectedRaider reports whether a raider is currently selected.
func (s *State) selectedRaider(id int) bool {
	return s.selection.kind == selectedItemRaider && s.selection.raiderID == id
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
