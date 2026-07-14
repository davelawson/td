package game

import "testing"

// TestClickingStructureTilesSelectsStructures verifies every current structure can be selected.
func TestClickingStructureTilesSelectsStructures(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
	}{
		{name: "sanctum", x: homePlotCenter, y: homePlotCenter},
		{name: "house", x: homePlotCenter + 2, y: 5},
		{name: "barracks", x: homePlotCenter + 3, y: 5},
		{name: "dorm", x: homePlotCenter + 4, y: 5},
		{name: "bow tower", x: homePlotCenter + 1, y: 5},
		{name: "flame bolt tower", x: homePlotCenter - 1, y: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := newRaidTestState(t)
			if tt.name == "bow tower" {
				state.gameMap.Home.Tiles[tt.y][tt.x].Feature = featureBowTower
			}
			if tt.name == "house" {
				state.gameMap.Home.Tiles[tt.y][tt.x].Feature = featureHouse
			}
			if tt.name == "barracks" {
				state.gameMap.Home.Tiles[tt.y][tt.x].Feature = featureBarracks
			}
			if tt.name == "dorm" {
				state.gameMap.Home.Tiles[tt.y][tt.x].Feature = featureDorm
			}
			if tt.name == "flame bolt tower" {
				state.gameMap.Home.Tiles[tt.y][tt.x].Feature = featureFlameBoltTower
			}

			state.Update(clickTileInput(state, tt.x, tt.y))

			if state.selection.kind != selectedItemStructure {
				t.Fatalf("selection kind = %v, want structure", state.selection.kind)
			}
			if state.selection.tile != (tileCoordinate{X: tt.x, Y: tt.y}) {
				t.Fatalf("selected tile = %+v, want (%d,%d)", state.selection.tile, tt.x, tt.y)
			}
		})
	}
}

// TestClickingNaturalTerrainSelectsTreeAndBoulder verifies selectable terrain scope.
func TestClickingNaturalTerrainSelectsTreeAndBoulder(t *testing.T) {
	tests := []struct {
		name    string
		terrain tileTerrain
	}{
		{name: "tree", terrain: terrainTree},
		{name: "boulder", terrain: terrainBoulder},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newRaidTestState(t)
			tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
			state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: test.terrain}

			state.Update(clickTileInput(state, tile.X, tile.Y))

			if state.selection.kind != selectedItemTerrain || state.selection.tile != tile {
				t.Fatalf("selection = %+v, want terrain at %+v", state.selection, tile)
			}
		})
	}
}

// TestTerrainSelectionWorksDuringPausedRaid verifies inspection remains phase-independent.
func TestTerrainSelectionWorksDuringPausedRaid(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainTree}
	state.status.phase = phaseRaid
	state.raid.active = true
	state.paused = true

	state.Update(clickTileInput(state, tile.X, tile.Y))

	if state.selection.kind != selectedItemTerrain || state.selection.tile != tile {
		t.Fatalf("selection = %+v, want paused-Raid terrain at %+v", state.selection, tile)
	}
}

// TestClickingGrassOrRoadDoesNotSelectTerrain verifies unoccupied Tiles remain unselectable.
func TestClickingGrassOrRoadDoesNotSelectTerrain(t *testing.T) {
	tests := []struct {
		name    string
		terrain tileTerrain
	}{
		{name: "grass", terrain: terrainEmpty},
		{name: "road", terrain: terrainRoad},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newRaidTestState(t)
			tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
			state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: test.terrain}
			state.selection = selectedItem{kind: selectedItemTerrain, tile: tileCoordinate{X: 1, Y: 1}}

			state.Update(clickTileInput(state, tile.X, tile.Y))

			if state.selection.kind != selectedItemNone {
				t.Fatalf("selection kind = %v, want none", state.selection.kind)
			}
		})
	}
}

// TestStructureSelectionHasPriorityOverTerrain verifies features win Tile overlap.
func TestStructureSelectionHasPriorityOverTerrain(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainTree, Feature: featureHouse}

	state.Update(clickTileInput(state, tile.X, tile.Y))

	if state.selection.kind != selectedItemStructure || state.selection.tile != tile {
		t.Fatalf("selection = %+v, want structure at %+v", state.selection, tile)
	}
}

// TestClickingRaiderSelectsRaider verifies active raiders can be selected by sprite bounds.
func TestClickingRaiderSelectsRaider(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{{
		id:       17,
		template: &state.enemyCatalog.SkeletonSwordShield,
		position: coord{X: 0, Y: 3},
		health:   state.enemyCatalog.SkeletonSwordShield.MaxHealth,
	}}
	state.paused = true

	state.Update(clickRaiderInput(state, state.raid.enemies[0]))

	if state.selection.kind != selectedItemRaider {
		t.Fatalf("selection kind = %v, want raider", state.selection.kind)
	}
	if state.selection.raiderID != 17 {
		t.Fatalf("selected raider ID = %d, want 17", state.selection.raiderID)
	}
}

// TestRaiderSelectionHasPriorityOverStructure verifies raiders win overlapping clicks.
func TestRaiderSelectionHasPriorityOverStructure(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{{
		id:       23,
		template: &state.enemyCatalog.Zombie,
		position: tileWorldCenter(homePlotCenter, homePlotCenter),
		health:   state.enemyCatalog.Zombie.MaxHealth,
	}}
	state.paused = true

	state.Update(clickTileInput(state, homePlotCenter, homePlotCenter))

	if state.selection.kind != selectedItemRaider {
		t.Fatalf("selection kind = %v, want raider", state.selection.kind)
	}
	if state.selection.raiderID != 23 {
		t.Fatalf("selected raider ID = %d, want 23", state.selection.raiderID)
	}
}

// TestClickingEmptySpaceClearsSelection verifies non-object clicks unselect the current item.
func TestClickingEmptySpaceClearsSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.Update(clickTileInput(state, homePlotCenter, homePlotCenter))

	state.Update(Input{CursorX: state.ui.width - 12, CursorY: state.ui.height - 12, Clicked: true})

	if state.selection.kind != selectedItemNone {
		t.Fatalf("selection kind = %v, want none", state.selection.kind)
	}
}

// TestSelectionWorksWhilePaused verifies pause still allows object inspection clicks.
func TestSelectionWorksWhilePaused(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.Update(Input{TogglePause: true})

	state.Update(clickTileInput(state, homePlotCenter+1, 5))

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want 0", state.Updates())
	}
}

// TestIngameMenuBlocksSelection verifies overlay-open frames do not change selected items.
func TestIngameMenuBlocksSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.Update(clickTileInput(state, homePlotCenter, homePlotCenter))
	state.Update(Input{ToggleMenu: true})

	state.Update(clickTileInput(state, homePlotCenter+1, 5))

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != (tileCoordinate{X: homePlotCenter, Y: homePlotCenter}) {
		t.Fatalf("selected tile = %+v, want Sanctum", state.selection.tile)
	}
}

// TestNextRaidClickDoesNotClearSelection verifies Raid UI clicks are not map selection clicks.
func TestNextRaidClickDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.Update(clickTileInput(state, homePlotCenter, homePlotCenter))

	state.Update(clickNextRaidInput(state))

	if !state.raid.active {
		t.Fatal("expected Next Raid click to start a Raid")
	}
	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != (tileCoordinate{X: homePlotCenter, Y: homePlotCenter}) {
		t.Fatalf("selected tile = %+v, want Sanctum", state.selection.tile)
	}
}

// TestRemovedRaiderSelectionClears verifies stale raider selections do not remain active.
func TestRemovedRaiderSelectionClears(t *testing.T) {
	state := newRaidTestState(t)
	step := state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond * gameUpdateSeconds
	state.status.barricade = 1
	state.status.phase = phaseRaid
	state.raid = raidState{
		active: true,
		enemies: []raidEnemy{{
			id:       31,
			template: &state.enemyCatalog.SkeletonSwordShield,
			position: coord{X: 0, Y: step},
			health:   state.enemyCatalog.SkeletonSwordShield.MaxHealth,
		}},
	}
	state.selection = selectedItem{kind: selectedItemRaider, raiderID: 31}

	state.Update(Input{})

	if state.selection.kind != selectedItemNone {
		t.Fatalf("selection kind = %v, want none", state.selection.kind)
	}
}

// clickTileInput returns a click at the center of a projected home Plot tile.
func clickTileInput(state *State, x, y int) Input {
	worldWest, worldNorth, worldW, worldH := tileWorldRect(x, y)
	rect := state.projectRect(state.sceneViewport(), worldWest, worldNorth, worldW, worldH)
	return clickProjectedRectInput(rect)
}

// clickRaiderInput returns a click at the center of a projected raider sprite.
func clickRaiderInput(state *State, enemy raidEnemy) Input {
	return clickProjectedRectInput(state.raidEnemyProjectedRect(state.sceneViewport(), enemy))
}

// clickProjectedRectInput returns a click at the center of a projected rectangle.
func clickProjectedRectInput(rect projectedRect) Input {
	return Input{
		CursorX: int(rect.x + rect.w/2),
		CursorY: int(rect.y + rect.h/2),
		Clicked: true,
	}
}
