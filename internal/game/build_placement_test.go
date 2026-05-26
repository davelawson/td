package game

import "testing"

// TestAffordableBuildingDragStarts verifies affordable tower icons can leave the bar.
func TestAffordableBuildingDragStarts(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, 0))

	if !state.buildDrag.active {
		t.Fatal("expected affordable Bow Tower drag to start")
	}
	if state.buildDrag.itemIndex != 0 {
		t.Fatalf("drag item index = %d, want Bow Tower index 0", state.buildDrag.itemIndex)
	}
}

// TestUnaffordableBuildingDragDoesNotStart verifies expensive towers cannot be dragged.
func TestUnaffordableBuildingDragDoesNotStart(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, 1))

	if state.buildDrag.active {
		t.Fatal("expected unaffordable Flame Bolt Tower drag not to start")
	}
}

// TestBuildDragTracksCursor verifies the dragged icon follows held mouse input.
func TestBuildDragTracksCursor(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(Input{CursorX: 320, CursorY: 440, MouseDown: true})

	if state.buildDrag.cursorX != 320 || state.buildDrag.cursorY != 440 {
		t.Fatalf("drag cursor = (%d,%d), want (320,440)", state.buildDrag.cursorX, state.buildDrag.cursorY)
	}
}

// TestBuildDragPlacesTowerAndDeductsResources verifies a valid drop constructs a tower.
func TestBuildDragPlacesTowerAndDeductsResources(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 50 || state.status.resources.stone != 35 || state.status.resources.metal != 2 {
		t.Fatalf("resources = %+v, want wood 50 stone 35 metal 2", state.status.resources)
	}
}

// TestBuildDragPlacesCatapultTower verifies Catapult Tower placement maps to its feature and cost.
func TestBuildDragPlacesCatapultTower(t *testing.T) {
	state := newRaidTestState(t)
	state.status.resources = resourceCounts{wood: 100, stone: 100, metal: 50}
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, 2))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureCatapultTower {
		t.Fatalf("tile feature = %v, want Catapult Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 60 || state.status.resources.stone != 40 || state.status.resources.metal != 25 {
		t.Fatalf("resources = %+v, want wood 60 stone 40 metal 25", state.status.resources)
	}
}

// TestBuildDragDoesNotReplaceOccupiedTile verifies occupied feature Tiles reject placement.
func TestBuildDragDoesNotReplaceOccupiedTile(t *testing.T) {
	state := newRaidTestState(t)
	initialResources := state.status.resources
	tile := tileCoordinate{X: homePlotCenter + 1, Y: 5}

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want existing Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestBuildDragRejectsRoadTile verifies roads are not buildable in the first placement slice.
func TestBuildDragRejectsRoadTile(t *testing.T) {
	state := newRaidTestState(t)
	initialResources := state.status.resources
	tile := tileCoordinate{X: homePlotCenter, Y: 4}

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureNone {
		t.Fatalf("road tile feature = %v, want none", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestBuildDragRejectsForestTile verifies forest border Tiles are not buildable.
func TestBuildDragRejectsForestTile(t *testing.T) {
	state := newRaidTestState(t)
	initialResources := state.status.resources
	tile := tileCoordinate{X: 1, Y: 0}

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureNone {
		t.Fatalf("forest tile feature = %v, want none", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestBuildDragRejectsActiveRaid verifies tower placement is calm-phase only.
func TestBuildDragRejectsActiveRaid(t *testing.T) {
	state := newRaidTestState(t)
	initialResources := state.status.resources

	state.startNextRaid()
	state.Update(pressBuildingBarItemInput(state, 0))

	if state.buildDrag.active {
		t.Fatal("expected active Raid to block build dragging")
	}
	state.Update(releaseTileInput(state, homePlotCenter+2, 5))
	if state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature != featureNone {
		t.Fatalf("tile feature = %v, want none", state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestBuildDragAllowsPausedCalmPlacement verifies pause does not block calm building.
func TestBuildDragAllowsPausedCalmPlacement(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.Update(Input{TogglePause: true})

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
}

// TestBuildDragInvalidReleaseClearsDrag verifies invalid drops are canceled safely.
func TestBuildDragInvalidReleaseClearsDrag(t *testing.T) {
	state := newRaidTestState(t)
	initialResources := state.status.resources

	state.Update(pressBuildingBarItemInput(state, 0))
	state.Update(Input{CursorX: -100, CursorY: -100, Released: true})

	if state.buildDrag.active {
		t.Fatal("expected invalid release to clear drag")
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// pressBuildingBarItemInput returns a left press at the center of a building-bar icon.
func pressBuildingBarItemInput(state *State, index int) Input {
	item := state.buildingBarItems()[index]
	return Input{
		CursorX:   item.Bounds.X + item.Bounds.W/2,
		CursorY:   item.Bounds.Y + item.Bounds.H/2,
		Clicked:   true,
		MouseDown: true,
	}
}

// releaseTileInput returns a left release at the center of a projected home Plot Tile.
func releaseTileInput(state *State, x, y int) Input {
	input := clickTileInput(state, x, y)
	input.Clicked = false
	input.MouseDown = false
	input.Released = true
	return input
}
