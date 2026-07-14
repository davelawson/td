package game

import (
	"testing"

	"td/internal/ui"
)

// TestBuildingBarVisibilityFollowsManagement verifies construction UI is phase-specific.
func TestBuildingBarVisibilityFollowsManagement(t *testing.T) {
	state := newRaidTestState(t)
	bar := state.buildingBarBounds()
	x := bar.X + bar.W/2
	y := bar.Y + 40

	if !state.buildingBarVisible() || !state.buildingBarContains(x, y) {
		t.Fatal("expected Management to expose the building bar and its input region")
	}

	state.status.phase = phaseLabour
	if state.buildingBarVisible() || state.buildingBarContains(x, y) {
		t.Fatal("expected instantaneous Labour to hide construction controls")
	}

	state.status.phase = phaseRaid
	state.raid.active = true
	if state.buildingBarVisible() || state.buildingBarContains(x, y) {
		t.Fatal("expected Raid to hide the building bar and release its input region")
	}

	state.raid.active = false
	state.raid.breached = true
	if state.buildingBarVisible() || state.buildingBarContains(x, y) {
		t.Fatal("expected breach to keep the building bar hidden")
	}
}

// TestHiddenBuildingBarClearsTransientState verifies hidden controls cannot retain hover or drag UI.
func TestHiddenBuildingBarClearsTransientState(t *testing.T) {
	state := newRaidTestState(t)
	state.ui.buildBarHover = 1
	state.ui.buildBarTabHover = ui.BuildingBarCategoryDefenses
	state.buildDrag = buildDragState{active: true, itemID: buildingBarBowTowerIndex}
	state.status.phase = phaseRaid
	state.raid.active = true

	state.updateBuildingBarHover(Input{})
	state.updateBuildDrag(Input{})

	if state.ui.buildBarHover != -1 || state.ui.buildBarTabHover != ui.BuildingBarNoCategory {
		t.Fatalf("hidden hover = %d/%v, want cleared", state.ui.buildBarHover, state.ui.buildBarTabHover)
	}
	if state.buildDrag.active {
		t.Fatal("expected hidden building bar to cancel build drag")
	}
}

// TestHiddenBuildingBarAreaAllowsCameraDrag verifies the former bar region becomes map input.
func TestHiddenBuildingBarAreaAllowsCameraDrag(t *testing.T) {
	state := newRaidTestState(t)
	state.status.phase = phaseRaid
	state.raid.active = true
	state.paused = true

	state.Update(rightDragInput(ui.BuildingBarWidth/2, topBarHeight+40, true, true, false))

	if !state.cameraDrag.active {
		t.Fatal("expected camera drag to start in the hidden building-bar region")
	}
}

// TestChallengeTextPreviewsNextRaid verifies peaceful HUD challenge uses live generator inputs.
func TestChallengeTextPreviewsNextRaid(t *testing.T) {
	state := newRaidTestState(t)

	if got, want := state.challengeText(), "Challenge 4.0"; got != want {
		t.Fatalf("challenge text = %q, want %q", got, want)
	}
	if got, want := state.phaseText(), "Management phase | Challenge 4.0"; got != want {
		t.Fatalf("phase text = %q, want %q", got, want)
	}

	state.status.populations.peasants = populationCount{available: 10, total: 10}
	state.gameMap.ensurePlots()
	state.gameMap.Plots[plotCoordinate{X: 1}] = &Plot{Biome: biomeHills}
	want := generateRaid(1, 10, 2).challengeRating
	if got := state.displayedChallengeRating(); got != want {
		t.Fatalf("displayed challenge = %f, want %f", got, want)
	}
}

// TestChallengeTextUsesFrozenRaidTemplate verifies active and breached Raids retain their rating.
func TestChallengeTextUsesFrozenRaidTemplate(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	frozen := state.raid.template.challengeRating
	state.status.populations.peasants = populationCount{available: 100, total: 100}

	if got := state.displayedChallengeRating(); got != frozen {
		t.Fatalf("active Raid challenge = %f, want frozen %f", got, frozen)
	}
	if got, want := state.phaseText(), "Enemies remaining: 3 | Challenge 4.0"; got != want {
		t.Fatalf("Raid phase text = %q, want %q", got, want)
	}

	state.raid.active = false
	state.raid.breached = true
	if got, want := state.phaseText(), "Sanctum breached | Challenge 4.0"; got != want {
		t.Fatalf("breach text = %q, want %q", got, want)
	}
}
