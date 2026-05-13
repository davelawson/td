package game

import (
	"testing"

	"td/internal/ui"
)

// TestNewStateStartsRunning verifies initial game state.
func TestNewStateStartsRunning(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if state.WizardName() != "Merlin" {
		t.Fatalf("wizard name = %q, want %q", state.WizardName(), "Merlin")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
	if state.Paused() {
		t.Fatal("expected new state to start unpaused")
	}
	if state.status.phase != phaseCalm {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseCalm)
	}
	if state.status.resources.wood != 80 || state.status.resources.stone != 45 || state.status.resources.metal != 12 {
		t.Fatalf("resources = %+v, want wood 80 stone 45 metal 12", state.status.resources)
	}
	if state.gameMap.Home.Tiles[homePlotCenter][homePlotCenter].Feature != featureSanctum {
		t.Fatal("expected new state to store the default home Plot")
	}
}

// TestDefaultHomePlotShape verifies the prototype Plot dimensions and center.
func TestDefaultHomePlotShape(t *testing.T) {
	plot := NewDefaultHomePlot()

	if len(plot.Tiles) != plotSize {
		t.Fatalf("plot rows = %d, want %d", len(plot.Tiles), plotSize)
	}
	if len(plot.Tiles[0]) != plotSize {
		t.Fatalf("plot columns = %d, want %d", len(plot.Tiles[0]), plotSize)
	}
	if plot.Tiles[homePlotCenter][homePlotCenter].Feature != featureSanctum {
		t.Fatal("expected Sanctum at the center Tile")
	}
}

// TestDefaultHomePlotRoadRunsNorth verifies the authored road layout.
func TestDefaultHomePlotRoadRunsNorth(t *testing.T) {
	plot := NewDefaultHomePlot()

	for y := 0; y <= homePlotCenter; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainRoad {
			t.Fatalf("tile (%d,%d) terrain = %v, want road", homePlotCenter, y, plot.Tiles[y][homePlotCenter].Terrain)
		}
	}
	for y := homePlotCenter + 1; y < plotSize; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain == terrainRoad {
			t.Fatalf("tile (%d,%d) is road below the Sanctum", homePlotCenter, y)
		}
	}
}

// TestDefaultHomePlotIsOtherwiseEmpty verifies no extra terrain or features exist yet.
func TestDefaultHomePlotIsOtherwiseEmpty(t *testing.T) {
	plot := NewDefaultHomePlot()

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			tile := plot.Tiles[y][x]
			onNorthRoad := x == homePlotCenter && y <= homePlotCenter
			atSanctum := x == homePlotCenter && y == homePlotCenter

			if !onNorthRoad && tile.Terrain != terrainEmpty {
				t.Fatalf("tile (%d,%d) terrain = %v, want empty", x, y, tile.Terrain)
			}
			if !atSanctum && tile.Feature != featureNone {
				t.Fatalf("tile (%d,%d) feature = %v, want none", x, y, tile.Feature)
			}
		}
	}
}

// TestUpdatesDoNotChangeHomePlot verifies time passing does not mutate the scene.
func TestUpdatesDoNotChangeHomePlot(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	initial := state.gameMap

	state.Update(Input{})
	state.Update(Input{})
	state.Update(Input{})

	if state.gameMap != initial {
		t.Fatal("expected logical updates to leave the prototype map unchanged")
	}
}

// TestStateFormatsCalmTopBar verifies initial top-bar text.
func TestStateFormatsCalmTopBar(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if value := state.chapterDayText(); value != "Chapter I: The Ashen Copse | Day 1" {
		t.Fatalf("chapterDayText = %q", value)
	}
	if value := state.phaseText(); value != "Raid in 02:00" {
		t.Fatalf("phaseText = %q", value)
	}
	if value := state.resourcesAndBarricadeText(); value != "Wood 80  Stone 45  Metal 12 | Barricade 3" {
		t.Fatalf("resourcesAndBarricadeText = %q", value)
	}
}

// TestStateFormatsRaidTopBar verifies raid-specific top-bar text.
func TestStateFormatsRaidTopBar(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	state.status.phase = phaseRaid
	state.status.raidCount = 7

	if value := state.phaseText(); value != "Enemies remaining: 7" {
		t.Fatalf("phaseText = %q", value)
	}
}

// TestUpdateIncrementsWhenRunning verifies logical updates advance while unpaused.
func TestUpdateIncrementsWhenRunning(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
}

// TestTogglePauseDoesNotIncrement verifies pause input is not a logical update.
func TestTogglePauseDoesNotIncrement(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	if !state.Paused() {
		t.Fatal("expected state to pause")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
}

// TestPausedUpdatesDoNotIncrement verifies logical updates stop while paused.
func TestPausedUpdatesDoNotIncrement(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{})
	state.Update(Input{})
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
}

// TestUnpauseThenUpdateIncrements verifies updates resume after unpausing.
func TestUnpauseThenUpdateIncrements(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{TogglePause: true})
	if state.Paused() {
		t.Fatal("expected state to unpause")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates after toggle = %d, want %d", state.Updates(), 0)
	}

	state.Update(Input{})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
}

// TestEscapeOpensIngameMenuAndPauses verifies ESC pauses and shows the overlay.
func TestEscapeOpensIngameMenuAndPauses(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if action := state.Update(Input{ToggleMenu: true}); action != ActionNone {
		t.Fatalf("Update(escape) = %v, want %v", action, ActionNone)
	}
	if !state.IngameMenuOpen() {
		t.Fatal("expected in-game menu to open")
	}
	if !state.Paused() {
		t.Fatal("expected in-game menu to pause the game")
	}
}

// TestResumeRestoresRunningState verifies Resume unpauses a previously running game.
func TestResumeRestoresRunningState(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{ToggleMenu: true})
	state.Update(clickInput(state.ui.menu.buttons[0]))
	if state.IngameMenuOpen() {
		t.Fatal("expected Resume to close the in-game menu")
	}
	if state.Paused() {
		t.Fatal("expected Resume to restore running state")
	}
}

// TestResumeRestoresPausedState verifies Resume preserves an existing pause.
func TestResumeRestoresPausedState(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{ToggleMenu: true})
	state.Update(Input{ToggleMenu: true})
	if state.IngameMenuOpen() {
		t.Fatal("expected ESC to close the in-game menu")
	}
	if !state.Paused() {
		t.Fatal("expected Resume to restore prior paused state")
	}
}

// TestIngameMenuBlocksUpdates verifies overlay-open frames do not advance logic.
func TestIngameMenuBlocksUpdates(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{})
	state.Update(Input{ToggleMenu: true})
	state.Update(Input{})
	state.Update(Input{TogglePause: true})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
	if !state.IngameMenuOpen() {
		t.Fatal("expected in-game menu to remain open")
	}
}

// TestSurrenderAction verifies Surrender reports a game-level action.
func TestSurrenderAction(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{ToggleMenu: true})
	if action := state.Update(clickInput(state.ui.menu.buttons[1])); action != ActionSurrender {
		t.Fatalf("Update(surrender click) = %v, want %v", action, ActionSurrender)
	}
}

// TestIngameMenuResizeRecentersButtons verifies resizing updates overlay hit geometry.
func TestIngameMenuResizeRecentersButtons(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Resize(2560, 1440)
	state.Update(Input{ToggleMenu: true})
	resumeButton := state.ui.menu.buttons[0]
	if resumeButton.X+resumeButton.W/2 != 1280 {
		t.Fatalf("resume center x = %d, want %d", resumeButton.X+resumeButton.W/2, 1280)
	}

	state.Update(clickInput(resumeButton))
	if state.IngameMenuOpen() {
		t.Fatal("expected resized Resume button to close the in-game menu")
	}
}

// clickInput returns a click at the center of button.
func clickInput(button ui.Button[Action]) Input {
	return Input{
		CursorX: button.X + button.W/2,
		CursorY: button.Y + button.H/2,
		Clicked: true,
	}
}
