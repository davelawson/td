package game

import "testing"

// TestPostRaidDayResolvesLabourBeforeManagement verifies production and the phase transition.
func TestPostRaidDayResolvesLabourBeforeManagement(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureWoodcutter
	startingWood := state.status.resources.wood

	state.beginPostRaidDay()

	if state.status.day != 2 {
		t.Fatalf("day = %d, want 2", state.status.day)
	}
	if state.status.resources.wood != startingWood+10 {
		t.Fatalf("wood = %d, want %d", state.status.resources.wood, startingWood+10)
	}
	if state.status.phase != phaseManagement {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseManagement)
	}
}

// TestManagementProducerWaitsForNextLabour verifies construction does not pay immediately.
func TestManagementProducerWaitsForNextLabour(t *testing.T) {
	state := newRaidTestState(t)
	state.ui.buildBarCategory = buildingBarCategoryEconomic
	setAvailablePopulations(state, 0, 0, 1)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	setHomeTilesEmpty(state, tile)
	startingWood := state.status.resources.wood

	state.Update(pressBuildingBarItemInput(state, buildingBarWoodcutterIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.status.resources.wood != startingWood-state.structureCatalog.Woodcutter.Cost.Wood {
		t.Fatalf("wood after construction = %d, want only the construction cost deducted", state.status.resources.wood)
	}
	state.beginPostRaidDay()
	if state.status.resources.wood != startingWood {
		t.Fatalf("wood after Labour = %d, want %d", state.status.resources.wood, startingWood)
	}
}

// TestManagementControlsRequireManagementPhase verifies peaceful actions are Management-only.
func TestManagementControlsRequireManagementPhase(t *testing.T) {
	state := newRaidTestState(t)

	state.status.phase = phaseLabour
	if state.canExploreNow() || state.canBuildTowersNow() || state.canStartRaid() {
		t.Fatal("expected Labour to block Management controls")
	}

	state.status.phase = phaseManagement
	if !state.canExploreNow() || !state.canBuildTowersNow() || !state.canStartRaid() {
		t.Fatal("expected Management to enable peaceful controls")
	}
}

// TestPausedManagementAllowsPreparationButNotRaidStart verifies the existing pause policy.
func TestPausedManagementAllowsPreparationButNotRaidStart(t *testing.T) {
	state := newRaidTestState(t)
	state.paused = true

	if !state.canExploreNow() || !state.canBuildTowersNow() {
		t.Fatal("expected paused Management to allow preparation")
	}
	if state.canStartRaid() {
		t.Fatal("expected paused Management to block Raid start")
	}
}

// TestLabourPhaseTextExistsForLifecycleDiagnostics verifies the atomic phase has a label.
func TestLabourPhaseTextExistsForLifecycleDiagnostics(t *testing.T) {
	state := newRaidTestState(t)
	state.status.phase = phaseLabour

	if got := state.phaseText(); got != "Labour phase | Challenge 4.0" {
		t.Fatalf("phaseText = %q, want %q", got, "Labour phase | Challenge 4.0")
	}
}
