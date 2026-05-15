package game

import (
	"testing"

	"td/internal/ui"
)

// TestNextRaidButtonStartsFirstRaid verifies the game UI starts a Raid immediately.
func TestNextRaidButtonStartsFirstRaid(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(clickNextRaidInput(state))

	if !state.raid.active {
		t.Fatal("expected Next Raid click to start an active Raid")
	}
	if state.raid.number != 1 {
		t.Fatalf("raid number = %d, want 1", state.raid.number)
	}
	if got, want := state.raidEnemiesRemaining(), firstRaidEnemyCount; got != want {
		t.Fatalf("remaining enemies = %d, want %d", got, want)
	}
	if len(state.raid.enemies) != 1 {
		t.Fatalf("active enemies = %d, want 1", len(state.raid.enemies))
	}
	if state.status.phase != phaseRaid {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseRaid)
	}
}

// TestNextRaidButtonDoesNotQueueWhileActive verifies active Raids block new starts.
func TestNextRaidButtonDoesNotQueueWhileActive(t *testing.T) {
	state := newRaidTestState(t)
	state.Update(clickNextRaidInput(state))

	state.Update(clickNextRaidInput(state))

	if state.raid.number != 1 {
		t.Fatalf("raid number = %d, want 1", state.raid.number)
	}
	if got, want := state.raidEnemiesRemaining(), firstRaidEnemyCount; got != want {
		t.Fatalf("remaining enemies = %d, want %d", got, want)
	}
}

// TestRaidSpawnsEnemiesOnStagger verifies enemies do not all spawn at once.
func TestRaidSpawnsEnemiesOnStagger(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()

	advanceRaidUpdates(state, raidSpawnInterval-1)
	if len(state.raid.enemies) != 1 {
		t.Fatalf("active enemies before stagger = %d, want 1", len(state.raid.enemies))
	}

	advanceRaidUpdates(state, 1)
	if len(state.raid.enemies) != 2 {
		t.Fatalf("active enemies after stagger = %d, want 2", len(state.raid.enemies))
	}
}

// TestRaidEnemiesMoveTowardSanctum verifies active enemies advance along the path.
func TestRaidEnemiesMoveTowardSanctum(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	start := state.raid.enemies[0].progress

	state.Update(Input{})

	if state.raid.enemies[0].progress <= start {
		t.Fatalf("enemy progress = %f, want greater than %f", state.raid.enemies[0].progress, start)
	}
}

// TestRaidCompletionReturnsToCalmAndAdvancesDay verifies successful Raid lifecycle completion.
func TestRaidCompletionReturnsToCalmAndAdvancesDay(t *testing.T) {
	state := newRaidTestState(t)
	state.status.barricade = 99
	state.startNextRaid()

	advanceUntilRaidEnds(t, state)

	if state.raid.breached {
		t.Fatal("expected Raid to complete without breach")
	}
	if state.raid.active {
		t.Fatal("expected Raid to be inactive after completion")
	}
	if state.status.phase != phaseCalm {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseCalm)
	}
	if state.status.day != 2 {
		t.Fatalf("day = %d, want 2", state.status.day)
	}
}

// TestRaidEnemyAtSanctumSpendsBarricade verifies reaching enemies consume charges.
func TestRaidEnemyAtSanctumSpendsBarricade(t *testing.T) {
	state := newRaidTestState(t)
	state.status.barricade = 2
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{progress: raidPathLength() - raidEnemySpeed}},
	}
	state.status.phase = phaseRaid

	state.Update(Input{})

	if state.status.barricade != 1 {
		t.Fatalf("barricade = %d, want 1", state.status.barricade)
	}
	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
}

// TestRaidBreachClearsRaidAndDisablesStarts verifies zero Barricade creates a terminal breach state.
func TestRaidBreachClearsRaidAndDisablesStarts(t *testing.T) {
	state := newRaidTestState(t)
	state.status.barricade = 0
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{progress: raidPathLength() - raidEnemySpeed}},
	}
	state.status.phase = phaseRaid

	state.Update(Input{})

	if !state.raid.breached {
		t.Fatal("expected Sanctum to be breached")
	}
	if state.raid.active {
		t.Fatal("expected Raid to be cleared after breach")
	}
	if state.raidEnemiesRemaining() != 0 {
		t.Fatalf("remaining enemies = %d, want 0", state.raidEnemiesRemaining())
	}
	if state.canStartRaid() {
		t.Fatal("expected breach to disable future Raid starts")
	}
	if value := state.phaseText(); value != "Sanctum breached" {
		t.Fatalf("phaseText = %q", value)
	}
}

// TestRaidDoesNotAdvanceWhilePaused verifies pause stops Raid logic.
func TestRaidDoesNotAdvanceWhilePaused(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	state.Update(Input{TogglePause: true})
	progress := state.raid.enemies[0].progress
	countdown := state.raid.spawnCountdown

	state.Update(Input{})

	if state.raid.enemies[0].progress != progress {
		t.Fatalf("enemy progress = %f, want %f", state.raid.enemies[0].progress, progress)
	}
	if state.raid.spawnCountdown != countdown {
		t.Fatalf("spawn countdown = %d, want %d", state.raid.spawnCountdown, countdown)
	}
}

// TestRaidDoesNotAdvanceWhileIngameMenuOpen verifies the overlay blocks Raid logic.
func TestRaidDoesNotAdvanceWhileIngameMenuOpen(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	state.Update(Input{ToggleMenu: true})
	progress := state.raid.enemies[0].progress
	countdown := state.raid.spawnCountdown

	state.Update(Input{})

	if state.raid.enemies[0].progress != progress {
		t.Fatalf("enemy progress = %f, want %f", state.raid.enemies[0].progress, progress)
	}
	if state.raid.spawnCountdown != countdown {
		t.Fatalf("spawn countdown = %d, want %d", state.raid.spawnCountdown, countdown)
	}
}

// newRaidTestState creates a game state for Raid behavior tests.
func newRaidTestState(t *testing.T) *State {
	t.Helper()
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	return state
}

// clickNextRaidInput returns an input click at the center of the Next Raid button.
func clickNextRaidInput(state *State) Input {
	return clickRaidButtonInput(state.nextRaidButton())
}

// clickRaidButtonInput returns an input click at the center of a Raid UI button.
func clickRaidButtonInput(button ui.Button[int]) Input {
	return Input{
		CursorX: button.X + button.W/2,
		CursorY: button.Y + button.H/2,
		Clicked: true,
	}
}

// advanceRaidUpdates applies empty update frames to Raid state.
func advanceRaidUpdates(state *State, updates int) {
	for i := 0; i < updates; i++ {
		state.Update(Input{})
	}
}

// advanceUntilRaidEnds updates until the active Raid ends or the test times out.
func advanceUntilRaidEnds(t *testing.T, state *State) {
	t.Helper()
	for i := 0; i < 2000; i++ {
		if !state.raid.active {
			return
		}
		state.Update(Input{})
	}
	t.Fatal("Raid did not end within 2000 updates")
}
