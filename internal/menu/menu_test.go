package menu

import "testing"

// TestButtonContainsIncludesTopLeft verifies inclusive top-left hit bounds.
func TestButtonContainsIncludesTopLeft(t *testing.T) {
	button := Button{X: 10, Y: 20, W: 100, H: 40}

	if !button.Contains(10, 20) {
		t.Fatal("expected top-left corner to be inside the button")
	}
}

// TestButtonContainsExcludesBottomRight verifies exclusive bottom-right hit bounds.
func TestButtonContainsExcludesBottomRight(t *testing.T) {
	button := Button{X: 10, Y: 20, W: 100, H: 40}

	if button.Contains(110, 60) {
		t.Fatal("expected bottom-right edge to be outside the button")
	}
}

// TestActionAtReturnsMatchingButtonAction verifies menu selection by point.
func TestActionAtReturnsMatchingButtonAction(t *testing.T) {
	buttons := []Button{
		{Label: "Quit", X: 10, Y: 20, W: 100, H: 40, Action: ActionQuit},
	}

	if action := ActionAt(buttons, 40, 30); action != ActionQuit {
		t.Fatalf("ActionAt() = %v, want %v", action, ActionQuit)
	}
}

// TestActionAtReturnsNoneOutsideButtons verifies empty menu selection.
func TestActionAtReturnsNoneOutsideButtons(t *testing.T) {
	buttons := []Button{
		{Label: "Quit", X: 10, Y: 20, W: 100, H: 40, Action: ActionQuit},
	}

	if action := ActionAt(buttons, 4, 30); action != ActionNone {
		t.Fatalf("ActionAt() = %v, want %v", action, ActionNone)
	}
}
