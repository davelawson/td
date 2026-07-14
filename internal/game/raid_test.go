package game

import (
	"math"
	"testing"

	"td/internal/ui"
)

// TestNextRaidButtonStartsFirstRaid verifies the game UI initializes a generated Raid.
func TestNextRaidButtonStartsFirstRaid(t *testing.T) {
	state := newRaidTestState(t)

	state.Update(clickNextRaidInput(state))

	if !state.raid.active {
		t.Fatal("expected Next Raid click to start an active Raid")
	}
	if state.raid.number != 1 {
		t.Fatalf("raid number = %d, want 1", state.raid.number)
	}
	if got, want := state.raidEnemiesRemaining(), 3; got != want {
		t.Fatalf("remaining enemies = %d, want %d", got, want)
	}
	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want none before the first threshold", len(state.raid.enemies))
	}
	if state.raid.progress <= 0 || state.raid.progress >= 1 {
		t.Fatalf("progress = %f, want one active update above zero", state.raid.progress)
	}
	if state.raid.template.challengeRating != 4 || state.raid.template.progressDurationSeconds != 9 {
		t.Fatalf("generated Raid = %+v, want challenge 4 and duration 9", state.raid.template)
	}
	if state.status.phase != phaseRaid {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseRaid)
	}
}

// TestRaidStartsAtZeroWithoutActiveEnemies verifies generation does not spawn immediately.
func TestRaidStartsAtZeroWithoutActiveEnemies(t *testing.T) {
	state := newRaidTestState(t)

	state.startNextRaid()

	if state.raid.progress != 0 {
		t.Fatalf("progress = %f, want 0", state.raid.progress)
	}
	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
	if got, want := state.raid.pendingEnemies, 3; got != want {
		t.Fatalf("pending enemies = %d, want %d", got, want)
	}

	state.updateRaid()
	if !state.raid.active {
		t.Fatal("expected an empty pre-threshold Raid interval to remain active")
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
	if got, want := state.raidEnemiesRemaining(), 3; got != want {
		t.Fatalf("remaining enemies = %d, want %d", got, want)
	}
}

// TestRaidProgressSpawnsAtFirstThreshold verifies exact equality schedules an enemy.
func TestRaidProgressSpawnsAtFirstThreshold(t *testing.T) {
	state := newRaidTestState(t)
	state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond = 0
	state.startNextRaid()

	state.updateRaidProgress(state.raid.template.progressDurationSeconds/2 - 0.01)
	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies before threshold = %d, want 0", len(state.raid.enemies))
	}

	state.updateRaidProgress(0.01)
	if len(state.raid.enemies) != 1 {
		t.Fatalf("active enemies at threshold = %d, want 1", len(state.raid.enemies))
	}
	if state.raid.enemies[0].template != &state.enemyCatalog.SkeletonSwordShield {
		t.Fatal("expected first threshold to spawn a Skeleton")
	}
	if got, want := state.raid.enemies[0].position, raidEnemySpawnPosition(); got != want {
		t.Fatalf("spawn position = %+v, want %+v", got, want)
	}
}

// TestRaidProgressSpawnsSimultaneousRulesInTemplateOrder verifies the full baseline roster.
func TestRaidProgressSpawnsSimultaneousRulesInTemplateOrder(t *testing.T) {
	state := newRaidTestState(t)
	state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond = 0
	state.enemyCatalog.Zombie.SpeedTilesPerSecond = 0

	state.startNextRaid()
	state.updateRaidProgress(state.raid.template.progressDurationSeconds)

	want := []*EnemyTemplate{
		&state.enemyCatalog.SkeletonSwordShield,
		&state.enemyCatalog.SkeletonSwordShield,
		&state.enemyCatalog.Zombie,
	}
	if len(state.raid.enemies) != len(want) {
		t.Fatalf("active enemies = %d, want %d", len(state.raid.enemies), len(want))
	}
	for i, template := range want {
		if state.raid.enemies[i].template != template {
			t.Fatalf("enemy %d template = %q, want %q", i, state.raid.enemies[i].template.Name, template.Name)
		}
		if state.raid.enemies[i].health != template.MaxHealth {
			t.Fatalf("enemy %d health = %d, want %d", i, state.raid.enemies[i].health, template.MaxHealth)
		}
	}
	if state.raid.progress != 1 || state.raid.pendingEnemies != 0 || !state.raid.active {
		t.Fatalf("completed schedule state = %+v, want active cleanup at 100%%", state.raid)
	}
}

// TestRaidGenerationUsesLiveSettlementInputs verifies State passes total population and all Plots.
func TestRaidGenerationUsesLiveSettlementInputs(t *testing.T) {
	state := newRaidTestState(t)
	state.status.populations.apprentices = populationCount{available: 0, total: 2}
	state.status.populations.soldiers = populationCount{available: 1, total: 3}
	state.status.populations.peasants = populationCount{available: 0, total: 4}
	state.gameMap.revealPlot(plotCoordinate{X: 1, Y: 0})
	state.raid.number = 1

	state.startNextRaid()

	want := generateRaid(2, 9, 2)
	if math.Abs(state.raid.template.challengeRating-want.challengeRating) > 0.000000001 {
		t.Fatalf("challenge = %f, want %f", state.raid.template.challengeRating, want.challengeRating)
	}
	if state.raid.pendingEnemies != want.totalEnemies() {
		t.Fatalf("pending enemies = %d, want %d", state.raid.pendingEnemies, want.totalEnemies())
	}
}

// TestRaidEnemiesMoveTowardSanctum verifies active enemies advance along the path.
func TestRaidEnemiesMoveTowardSanctum(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	state.updateRaidProgress(state.raid.template.progressDurationSeconds / 2)
	start := state.raid.enemies[0].position

	state.Update(Input{})

	if state.raid.enemies[0].position.Y >= start.Y {
		t.Fatalf("enemy y = %f, want less than %f", state.raid.enemies[0].position.Y, start.Y)
	}
	if state.raid.enemies[0].position.X != start.X {
		t.Fatalf("enemy x = %f, want %f", state.raid.enemies[0].position.X, start.X)
	}
}

// TestRaidEnemyMovementUsesTemplateSpeed verifies enemy speed is second-based template data.
func TestRaidEnemyMovementUsesTemplateSpeed(t *testing.T) {
	state := newRaidTestState(t)
	template := &EnemyTemplate{SpeedTilesPerSecond: 2.5}
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{template: template, position: coord{X: 0, Y: 5}}},
	}

	state.updateRaidEnemies()

	if got, want := state.raid.enemies[0].position.Y, 5.0-2.5*gameUpdateSeconds; got != want {
		t.Fatalf("enemy y = %f, want %f", got, want)
	}
}

// TestRaidEnemyWithoutPositiveSpeedDoesNotMove verifies malformed enemies stay put.
func TestRaidEnemyWithoutPositiveSpeedDoesNotMove(t *testing.T) {
	state := newRaidTestState(t)
	state.raid = raidState{
		active: true,
		enemies: []raidEnemy{
			{position: coord{X: 0, Y: 5}},
			{template: &EnemyTemplate{SpeedTilesPerSecond: -1}, position: coord{X: 0, Y: 4}},
		},
	}

	state.updateRaidEnemies()

	if got, want := state.raid.enemies[0].position.Y, 5.0; got != want {
		t.Fatalf("nil-template enemy y = %f, want %f", got, want)
	}
	if got, want := state.raid.enemies[1].position.Y, 4.0; got != want {
		t.Fatalf("nonpositive-speed enemy y = %f, want %f", got, want)
	}
}

// TestRaidCompletionReturnsToManagementAndAdvancesDay verifies successful Raid lifecycle completion.
func TestRaidCompletionReturnsToManagementAndAdvancesDay(t *testing.T) {
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
	if state.status.phase != phaseManagement {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseManagement)
	}
	if state.status.day != 2 {
		t.Fatalf("day = %d, want 2", state.status.day)
	}
}

// TestPostRaidLabourGrantsEconomicBuildingResources verifies Labour pays producers.
func TestPostRaidLabourGrantsEconomicBuildingResources(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureWoodcutter
	state.gameMap.Home.Tiles[5][homePlotCenter+3].Feature = featureStoneQuarry
	state.gameMap.Home.Tiles[5][homePlotCenter+4].Feature = featureIronMine
	startingResources := state.status.resources

	state.completeRaid()

	if state.status.resources.wood != startingResources.wood+10 {
		t.Fatalf("wood = %d, want %d", state.status.resources.wood, startingResources.wood+10)
	}
	if state.status.resources.stone != startingResources.stone+10 {
		t.Fatalf("stone = %d, want %d", state.status.resources.stone, startingResources.stone+10)
	}
	if state.status.resources.metal != startingResources.metal+10 {
		t.Fatalf("metal = %d, want %d", state.status.resources.metal, startingResources.metal+10)
	}
}

// TestRaidBreachDoesNotResolveLabour verifies failed Raids do not begin the next Day.
func TestRaidBreachDoesNotResolveLabour(t *testing.T) {
	state := newRaidTestState(t)
	state.status.barricade = 0
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureWoodcutter
	step := state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond * gameUpdateSeconds
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{template: &state.enemyCatalog.SkeletonSwordShield, position: coord{X: 0, Y: step}}},
	}
	state.status.phase = phaseRaid
	startingResources := state.status.resources
	startingDay := state.status.day

	state.Update(Input{})

	if state.status.resources != startingResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, startingResources)
	}
	if state.status.day != startingDay {
		t.Fatalf("day = %d, want unchanged %d", state.status.day, startingDay)
	}
	if state.status.phase != phaseRaid {
		t.Fatalf("phase = %v, want terminal %v", state.status.phase, phaseRaid)
	}
}

// TestRaidEnemyAtSanctumSpendsBarricade verifies reaching enemies consume charges.
func TestRaidEnemyAtSanctumSpendsBarricade(t *testing.T) {
	state := newRaidTestState(t)
	state.status.barricade = 2
	step := state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond * gameUpdateSeconds
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{template: &state.enemyCatalog.SkeletonSwordShield, position: coord{X: 0, Y: step}}},
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
	step := state.enemyCatalog.SkeletonSwordShield.SpeedTilesPerSecond * gameUpdateSeconds
	state.raid = raidState{
		active:  true,
		enemies: []raidEnemy{{template: &state.enemyCatalog.SkeletonSwordShield, position: coord{X: 0, Y: step}}},
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
	progress := state.raid.progress

	state.Update(Input{})

	if state.raid.progress != progress {
		t.Fatalf("progress = %f, want paused %f", state.raid.progress, progress)
	}
}

// TestRaidDoesNotAdvanceWhileIngameMenuOpen verifies the overlay blocks Raid logic.
func TestRaidDoesNotAdvanceWhileIngameMenuOpen(t *testing.T) {
	state := newRaidTestState(t)
	state.startNextRaid()
	state.Update(Input{ToggleMenu: true})
	progress := state.raid.progress

	state.Update(Input{})

	if state.raid.progress != progress {
		t.Fatalf("progress = %f, want overlay-frozen %f", state.raid.progress, progress)
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
