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

// TestMenuUpdateRoutesMenuScreens verifies menu actions change menu screens or report quit.
func TestMenuUpdateRoutesMenuScreens(t *testing.T) {
	menu, err := New(960, 540)
	if err != nil {
		t.Fatal(err)
	}

	if action := menu.Update(480, 274, true); action != ActionNew {
		t.Fatalf("Update(new click) = %v, want %v", action, ActionNew)
	}
	if menu.Screen() != ScreenNewGame {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenNewGame)
	}

	if action := menu.Update(480, 411, true); action != ActionBack {
		t.Fatalf("Update(back click) = %v, want %v", action, ActionBack)
	}
	if menu.Screen() != ScreenMain {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenMain)
	}

	if action := menu.Update(480, 382, true); action != ActionSettings {
		t.Fatalf("Update(settings click) = %v, want %v", action, ActionSettings)
	}
	if menu.Screen() != ScreenSettings {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenSettings)
	}

	if action := menu.Update(4, 30, true); action != ActionNone {
		t.Fatalf("Update(outside click) = %v, want %v", action, ActionNone)
	}
	if menu.Screen() != ScreenSettings {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenSettings)
	}

	menu.SetScreenForTest(ScreenMain)
	if action := menu.Update(480, 436, true); action != ActionQuit {
		t.Fatalf("Update(quit click) = %v, want %v", action, ActionQuit)
	}
	if menu.Screen() != ScreenMain {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenMain)
	}
}
