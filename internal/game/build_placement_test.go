package game

import "testing"

// TestAffordableBuildingDragStarts verifies affordable building icons can leave the bar.
func TestAffordableBuildingDragStarts(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))

	if !state.buildDrag.active {
		t.Fatal("expected affordable House drag to start")
	}
	if state.buildDrag.itemID != buildingBarHouseIndex {
		t.Fatalf("drag item ID = %d, want House", state.buildDrag.itemID)
	}
}

// TestInsufficientStaffBuildingDragDoesNotStart verifies staffing gates dragging.
func TestInsufficientStaffBuildingDragDoesNotStart(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))

	if state.buildDrag.active {
		t.Fatal("expected unstaffed Bow Tower drag not to start")
	}
}

// TestBuildDragTracksCursor verifies the dragged icon follows held mouse input.
func TestBuildDragTracksCursor(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(Input{CursorX: 320, CursorY: 440, MouseDown: true})

	if state.buildDrag.cursorX != 320 || state.buildDrag.cursorY != 440 {
		t.Fatalf("drag cursor = (%d,%d), want (320,440)", state.buildDrag.cursorX, state.buildDrag.cursorY)
	}
}

// TestBuildDragPlacesTowerAndDeductsResources verifies a valid drop constructs a tower.
func TestBuildDragPlacesTowerAndDeductsResources(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 70 || state.status.resources.stone != 40 || state.status.resources.metal != 10 {
		t.Fatalf("resources = %+v, want wood 70 stone 40 metal 10", state.status.resources)
	}
	if state.status.populations.soldiers != (populationCount{available: 0, total: 1}) {
		t.Fatalf("soldiers = %+v, want 0/1 after staffing Bow Tower", state.status.populations.soldiers)
	}
}

// TestBuildDragPlacesHouseAndGrantsPeasants verifies House construction adds population.
func TestBuildDragPlacesHouseAndGrantsPeasants(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureHouse {
		t.Fatalf("tile feature = %v, want House", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 80 || state.status.resources.stone != 50 || state.status.resources.metal != 20 {
		t.Fatalf("resources = %+v, want wood 80 stone 50 metal 20", state.status.resources)
	}
	if state.status.populations.peasants != (populationCount{available: 2, total: 2}) {
		t.Fatalf("peasants = %+v, want 2/2 after building House", state.status.populations.peasants)
	}
}

// TestBuildDragPlacesEconomicBuildingAndReservesPeasant verifies worker assignment.
func TestBuildDragPlacesEconomicBuildingAndReservesPeasant(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 0, 1)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarWoodcutterIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureWoodcutter {
		t.Fatalf("tile feature = %v, want Woodcutter", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 90 || state.status.resources.stone != 50 || state.status.resources.metal != 20 {
		t.Fatalf("resources = %+v, want wood 90 stone 50 metal 20", state.status.resources)
	}
	if state.status.populations.peasants != (populationCount{available: 0, total: 1}) {
		t.Fatalf("peasants = %+v, want 0/1 after staffing Woodcutter", state.status.populations.peasants)
	}
}

// TestBarracksDragRequiresPeasants verifies conversion buildings need source population.
func TestBarracksDragRequiresPeasants(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarBarracksIndex))

	if state.buildDrag.active {
		t.Fatal("expected Barracks drag not to start without Peasants")
	}
}

// TestBuildDragPlacesBarracksAndConvertsPeasants verifies Barracks population conversion.
func TestBuildDragPlacesBarracksAndConvertsPeasants(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 0, 2)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarBarracksIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBarracks {
		t.Fatalf("tile feature = %v, want Barracks", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 90 || state.status.resources.stone != 40 || state.status.resources.metal != 20 {
		t.Fatalf("resources = %+v, want wood 90 stone 40 metal 20", state.status.resources)
	}
	if state.status.populations.peasants != (populationCount{available: 0, total: 0}) {
		t.Fatalf("peasants = %+v, want 0/0 after building Barracks", state.status.populations.peasants)
	}
	if state.status.populations.soldiers != (populationCount{available: 2, total: 2}) {
		t.Fatalf("soldiers = %+v, want 2/2 after building Barracks", state.status.populations.soldiers)
	}
}

// TestDormDragRequiresPeasant verifies Dorm conversion needs source population.
func TestDormDragRequiresPeasant(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarDormIndex))

	if state.buildDrag.active {
		t.Fatal("expected Dorm drag not to start without a Peasant")
	}
}

// TestBuildDragPlacesDormAndConvertsPeasant verifies Dorm population conversion.
func TestBuildDragPlacesDormAndConvertsPeasant(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 0, 1)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarDormIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.buildDrag.active {
		t.Fatal("expected drag to clear after release")
	}
	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureDorm {
		t.Fatalf("tile feature = %v, want Dorm", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 90 || state.status.resources.stone != 40 || state.status.resources.metal != 20 {
		t.Fatalf("resources = %+v, want wood 90 stone 40 metal 20", state.status.resources)
	}
	if state.status.populations.peasants != (populationCount{available: 0, total: 0}) {
		t.Fatalf("peasants = %+v, want 0/0 after building Dorm", state.status.populations.peasants)
	}
	if state.status.populations.apprentices != (populationCount{available: 1, total: 1}) {
		t.Fatalf("apprentices = %+v, want 1/1 after building Dorm", state.status.populations.apprentices)
	}
}

// TestDormInvalidReleasePreservesPopulations verifies conversion is atomic.
func TestDormInvalidReleasePreservesPopulations(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 0, 1)
	initialResources := state.status.resources
	initialPopulations := state.status.populations

	state.Update(pressBuildingBarItemInput(state, buildingBarDormIndex))
	state.Update(Input{CursorX: -100, CursorY: -100, Released: true})

	if state.buildDrag.active {
		t.Fatal("expected invalid release to clear drag")
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
	if state.status.populations != initialPopulations {
		t.Fatalf("populations = %+v, want unchanged %+v", state.status.populations, initialPopulations)
	}
}

// TestHouseThenDormEnablesFlameBoltTower verifies the first Apprentice-producing workflow.
func TestHouseThenDormEnablesFlameBoltTower(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	state.Update(releaseTileInput(state, homePlotCenter+2, 5))
	state.Update(pressBuildingBarItemInput(state, buildingBarDormIndex))
	state.Update(releaseTileInput(state, homePlotCenter+3, 5))
	state.Update(pressBuildingBarItemInput(state, buildingBarFlameBoltTowerIndex))
	state.Update(releaseTileInput(state, homePlotCenter+4, 5))

	if state.gameMap.Home.Tiles[5][homePlotCenter+4].Feature != featureFlameBoltTower {
		t.Fatalf("tile feature = %v, want Flame Bolt Tower", state.gameMap.Home.Tiles[5][homePlotCenter+4].Feature)
	}
	if state.status.populations.apprentices != (populationCount{available: 0, total: 1}) {
		t.Fatalf("apprentices = %+v, want 0/1 after staffing Flame Bolt Tower", state.status.populations.apprentices)
	}
	if state.status.populations.peasants != (populationCount{available: 1, total: 1}) {
		t.Fatalf("peasants = %+v, want 1/1 after Dorm conversion", state.status.populations.peasants)
	}
}

// TestBarracksInvalidReleasePreservesPopulations verifies conversion is atomic.
func TestBarracksInvalidReleasePreservesPopulations(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 0, 2)
	initialResources := state.status.resources
	initialPopulations := state.status.populations

	state.Update(pressBuildingBarItemInput(state, buildingBarBarracksIndex))
	state.Update(Input{CursorX: -100, CursorY: -100, Released: true})

	if state.buildDrag.active {
		t.Fatal("expected invalid release to clear drag")
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
	if state.status.populations != initialPopulations {
		t.Fatalf("populations = %+v, want unchanged %+v", state.status.populations, initialPopulations)
	}
}

// TestHouseThenBarracksEnablesBowTower verifies the first Soldier-producing workflow.
func TestHouseThenBarracksEnablesBowTower(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	state.Update(releaseTileInput(state, homePlotCenter+2, 5))
	state.Update(pressBuildingBarItemInput(state, buildingBarBarracksIndex))
	state.Update(releaseTileInput(state, homePlotCenter+3, 5))
	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(releaseTileInput(state, homePlotCenter+4, 5))

	if state.gameMap.Home.Tiles[5][homePlotCenter+4].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want Bow Tower", state.gameMap.Home.Tiles[5][homePlotCenter+4].Feature)
	}
	if state.status.populations.soldiers != (populationCount{available: 1, total: 2}) {
		t.Fatalf("soldiers = %+v, want 1/2 after staffing Bow Tower", state.status.populations.soldiers)
	}
	if state.status.populations.peasants != (populationCount{available: 0, total: 0}) {
		t.Fatalf("peasants = %+v, want 0/0 after Barracks conversion", state.status.populations.peasants)
	}
}

// TestBuildDragPlacesCatapultTower verifies Catapult Tower placement maps to its feature and cost.
func TestBuildDragPlacesCatapultTower(t *testing.T) {
	state := newRaidTestState(t)
	state.status.resources = resourceCounts{wood: 100, stone: 100, metal: 50}
	setAvailablePopulations(state, 0, 1, 2)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarCatapultTowerIndex))
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
	if state.status.populations.soldiers != (populationCount{available: 0, total: 1}) {
		t.Fatalf("soldiers = %+v, want 0/1", state.status.populations.soldiers)
	}
	if state.status.populations.peasants != (populationCount{available: 0, total: 2}) {
		t.Fatalf("peasants = %+v, want 0/2", state.status.populations.peasants)
	}
}

// TestBuildDragRejectsCatapultWhenOneRoleIsShort verifies mixed requirements are atomic.
func TestBuildDragRejectsCatapultWhenOneRoleIsShort(t *testing.T) {
	state := newRaidTestState(t)
	state.status.resources = resourceCounts{wood: 100, stone: 100, metal: 50}
	setAvailablePopulations(state, 0, 1, 1)

	state.Update(pressBuildingBarItemInput(state, buildingBarCatapultTowerIndex))

	if state.buildDrag.active {
		t.Fatal("expected one missing Peasant to block Catapult drag")
	}
}

// TestSecondTowerBuildIsBlockedAfterStaffReservation verifies staff cannot be reused.
func TestSecondTowerBuildIsBlockedAfterStaffReservation(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(releaseTileInput(state, homePlotCenter+2, 5))
	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))

	if state.buildDrag.active {
		t.Fatal("expected reserved Soldier to block a second Bow Tower")
	}
}

// TestBuildReleaseRechecksStaffing verifies changed staffing cancels an active drag.
func TestBuildReleaseRechecksStaffing(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.status.populations.soldiers.available = 0
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureNone {
		t.Fatalf("tile feature = %v, want none", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestBuildDragDoesNotReplaceOccupiedTile verifies occupied feature Tiles reject placement.
func TestBuildDragDoesNotReplaceOccupiedTile(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources
	initialPopulations := state.status.populations
	tile := tileCoordinate{X: homePlotCenter + 1, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X].Feature = featureBowTower

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want existing Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
	if state.status.populations != initialPopulations {
		t.Fatalf("populations = %+v, want unchanged %+v", state.status.populations, initialPopulations)
	}
}

// TestBuildDragRejectsRoadTile verifies roads are not buildable in the first placement slice.
func TestBuildDragRejectsRoadTile(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources
	tile := tileCoordinate{X: homePlotCenter, Y: 4}

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
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
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources
	tile := tileCoordinate{X: 1, Y: 0}

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
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
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources

	state.startNextRaid()
	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))

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
	setAvailablePopulations(state, 0, 1, 0)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.Update(Input{TogglePause: true})

	state.Update(pressBuildingBarItemInput(state, buildingBarBowTowerIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureBowTower {
		t.Fatalf("tile feature = %v, want Bow Tower", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
}

// TestBuildDragInvalidReleaseClearsDrag verifies invalid drops are canceled safely.
func TestBuildDragInvalidReleaseClearsDrag(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 0)
	initialResources := state.status.resources
	initialPopulations := state.status.populations

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	state.Update(Input{CursorX: -100, CursorY: -100, Released: true})

	if state.buildDrag.active {
		t.Fatal("expected invalid release to clear drag")
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
	if state.status.populations != initialPopulations {
		t.Fatalf("populations = %+v, want unchanged %+v", state.status.populations, initialPopulations)
	}
}

// pressBuildingBarItemInput returns a left press at the center of a building-bar icon.
func pressBuildingBarItemInput(state *State, id buildingBarItemID) Input {
	state.ui.buildBarCategory = buildingBarCategoryForItem(id)
	item, ok := state.buildingBarItemByID(id)
	if !ok {
		return Input{}
	}
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

// setAvailablePopulations gives tests matching available and total inhabitant counts.
func setAvailablePopulations(state *State, apprentices, soldiers, peasants int) {
	state.status.populations = populationCounts{
		apprentices: populationCount{available: apprentices, total: apprentices},
		soldiers:    populationCount{available: soldiers, total: soldiers},
		peasants:    populationCount{available: peasants, total: peasants},
	}
}
