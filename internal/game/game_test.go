package game

import "testing"

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
