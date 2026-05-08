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
	if state.phase != phaseCalm {
		t.Fatalf("phase = %v, want %v", state.phase, phaseCalm)
	}
	if state.resources.wood != 80 || state.resources.stone != 45 || state.resources.metal != 12 {
		t.Fatalf("resources = %+v, want wood 80 stone 45 metal 12", state.resources)
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
	state.phase = phaseRaid
	state.raidCount = 7

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
	state.Update(clickInput(state.menu.buttons[0]))
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
	if action := state.Update(clickInput(state.menu.buttons[1])); action != ActionSurrender {
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
	resumeButton := state.menu.buttons[0]
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
