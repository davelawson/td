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
		{Label: "New", X: 10, Y: 20, W: 100, H: 40, Action: ActionNew},
		{Label: "Settings", X: 10, Y: 70, W: 100, H: 40, Action: ActionSettings},
		{Label: "Back", X: 10, Y: 120, W: 100, H: 40, Action: ActionBack},
		{Label: "Quit", X: 10, Y: 170, W: 100, H: 40, Action: ActionQuit},
	}

	tests := []struct {
		name string
		x    int
		y    int
		want Action
	}{
		{name: "new", x: 40, y: 30, want: ActionNew},
		{name: "settings", x: 40, y: 80, want: ActionSettings},
		{name: "back", x: 40, y: 130, want: ActionBack},
		{name: "quit", x: 40, y: 180, want: ActionQuit},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if action := ActionAt(buttons, test.x, test.y); action != test.want {
				t.Fatalf("ActionAt() = %v, want %v", action, test.want)
			}
		})
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

// TestActionAtIgnoresDisabledButtons verifies disabled targets do not select actions.
func TestActionAtIgnoresDisabledButtons(t *testing.T) {
	buttons := []Button{
		{Label: "Load", X: 10, Y: 20, W: 100, H: 40, Action: ActionNone, Disabled: true},
		{Label: "Quit", X: 10, Y: 70, W: 100, H: 40, Action: ActionQuit},
	}

	if action := ActionAt(buttons, 40, 30); action != ActionNone {
		t.Fatalf("ActionAt() = %v, want %v", action, ActionNone)
	}
}
