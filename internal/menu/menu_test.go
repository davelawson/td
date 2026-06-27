package menu

import (
	"strings"
	"testing"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	testWidth  = 1920
	testHeight = 1080
)

// TestMenuFontSizesUseUIConstants verifies menu faces use centralized sizes.
func TestMenuFontSizesUseUIConstants(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}

	if menu.titleFace.Size != ui.MenuTitleFontSize {
		t.Fatalf("title font size = %.1f, want %.1f", menu.titleFace.Size, ui.MenuTitleFontSize)
	}
	if menu.bodyFace.Size != ui.MenuBodyFontSize {
		t.Fatalf("body font size = %.1f, want %.1f", menu.bodyFace.Size, ui.MenuBodyFontSize)
	}
	if menu.buttonFace.Size != ui.MenuButtonFontSize {
		t.Fatalf("button font size = %.1f, want %.1f", menu.buttonFace.Size, ui.MenuButtonFontSize)
	}
	if menu.nameFace.Size != ui.MenuNameFontSize {
		t.Fatalf("name font size = %.1f, want %.1f", menu.nameFace.Size, ui.MenuNameFontSize)
	}
}

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
		{Label: "Start", X: 10, Y: 220, W: 100, H: 40, Action: ActionStart},
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
		{name: "start", x: 40, y: 230, want: ActionStart},
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
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}

	newButton := menu.buttonForAction(ActionNew)
	if action := menu.Update(clickInput(newButton)); action != ActionNew {
		t.Fatalf("Update(new click) = %v, want %v", action, ActionNew)
	}
	if menu.Screen() != ScreenNewGame {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenNewGame)
	}
	if !menu.WizardNameFocused() {
		t.Fatal("expected wizard name field to focus after entering new game")
	}

	cancelButton := menu.buttonForAction(ActionBack)
	if action := menu.Update(clickInput(cancelButton)); action != ActionBack {
		t.Fatalf("Update(cancel click) = %v, want %v", action, ActionBack)
	}
	if menu.Screen() != ScreenMain {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenMain)
	}

	settingsButton := menu.buttonForAction(ActionSettings)
	if action := menu.Update(clickInput(settingsButton)); action != ActionSettings {
		t.Fatalf("Update(settings click) = %v, want %v", action, ActionSettings)
	}
	if menu.Screen() != ScreenSettings {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenSettings)
	}

	if action := menu.Update(Input{CursorX: 4, CursorY: 30, Clicked: true}); action != ActionNone {
		t.Fatalf("Update(outside click) = %v, want %v", action, ActionNone)
	}
	if menu.Screen() != ScreenSettings {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenSettings)
	}

	menu.SetScreenForTest(ScreenMain)
	quitButton := menu.buttonForAction(ActionQuit)
	if action := menu.Update(clickInput(quitButton)); action != ActionQuit {
		t.Fatalf("Update(quit click) = %v, want %v", action, ActionQuit)
	}
	if menu.Screen() != ScreenMain {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenMain)
	}
}

// TestWizardNameInputEditsFocusedField verifies keyboard edits on the name field.
func TestWizardNameInputEditsFocusedField(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}
	menu.SetScreenForTest(ScreenNewGame)

	menu.Update(Input{Typed: []rune("Merlin")})
	if menu.WizardName() != "Merlin" {
		t.Fatalf("wizard name = %q, want %q", menu.WizardName(), "Merlin")
	}

	menu.Update(Input{Backspace: true})
	if menu.WizardName() != "Merli" {
		t.Fatalf("wizard name = %q, want %q", menu.WizardName(), "Merli")
	}
}

// TestWizardNameInputCapsLength verifies the name field keeps a stable width.
func TestWizardNameInputCapsLength(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}
	menu.SetScreenForTest(ScreenNewGame)

	menu.Update(Input{Typed: []rune(strings.Repeat("a", wizardNameMaxRunes+4))})
	if got := len([]rune(menu.WizardName())); got != wizardNameMaxRunes {
		t.Fatalf("wizard name length = %d, want %d", got, wizardNameMaxRunes)
	}
}

// TestWizardNameMaxLengthFitsField verifies the field can display 32 runes.
func TestWizardNameMaxLengthFitsField(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}

	fieldX, _, fieldW, _ := menu.wizardNameFieldBounds()
	width, _ := text.Measure(strings.Repeat("W", wizardNameMaxRunes), menu.nameFace, menu.nameFace.Size)
	if int(width) > fieldW-36 {
		t.Fatalf("max name width = %d, want at most %d inside field at x %d", int(width), fieldW-36, fieldX)
	}
}

// TestWizardNameFocusFollowsFieldClick verifies field clicks activate text entry.
func TestWizardNameFocusFollowsFieldClick(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}
	menu.SetScreenForTest(ScreenNewGame)
	menu.wizardNameFocused = false

	fieldX, fieldY, fieldW, fieldH := menu.wizardNameFieldBounds()
	menu.Update(Input{CursorX: fieldX + fieldW/2, CursorY: fieldY + fieldH/2, Clicked: true})
	if !menu.WizardNameFocused() {
		t.Fatal("expected wizard name field click to focus text entry")
	}
}

// TestMenuResizeRecentersButtonTargets verifies resizing updates hit geometry.
func TestMenuResizeRecentersButtonTargets(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}

	menu.Resize(2560, 1440)
	newButton := menu.buttonForAction(ActionNew)
	if newButton.X+newButton.W/2 != 1280 {
		t.Fatalf("new button center x = %d, want %d", newButton.X+newButton.W/2, 1280)
	}

	if action := menu.Update(clickInput(newButton)); action != ActionNew {
		t.Fatalf("Update(resized new click) = %v, want %v", action, ActionNew)
	}
}

// TestStartIsDisabledWithoutWizardName verifies Start is inert before name entry.
func TestStartIsDisabledWithoutWizardName(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}
	menu.SetScreenForTest(ScreenNewGame)

	startButton := menu.disabledButtonWithLabel("Start")
	if action := menu.Update(clickInput(startButton)); action != ActionNone {
		t.Fatalf("Update(start click) = %v, want %v", action, ActionNone)
	}
	if menu.Screen() != ScreenNewGame {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenNewGame)
	}
}

// TestStartBecomesActiveWithWizardName verifies a named Wizard can start.
func TestStartBecomesActiveWithWizardName(t *testing.T) {
	menu, err := New(testWidth, testHeight)
	if err != nil {
		t.Fatal(err)
	}
	menu.SetScreenForTest(ScreenNewGame)

	menu.Update(Input{Typed: []rune("Merlin")})
	startButton := menu.buttonForAction(ActionStart)
	if action := menu.Update(clickInput(startButton)); action != ActionStart {
		t.Fatalf("Update(start click) = %v, want %v", action, ActionStart)
	}
	if menu.Screen() != ScreenNewGame {
		t.Fatalf("screen = %v, want %v", menu.Screen(), ScreenNewGame)
	}
}

// buttonForAction returns the active button that reports action.
func (m *Menu) buttonForAction(action Action) Button {
	for _, button := range m.activeButtons() {
		if button.Action == action && !button.Disabled {
			return button
		}
	}
	return Button{}
}

// disabledButtonWithLabel returns the active disabled button with label.
func (m *Menu) disabledButtonWithLabel(label string) Button {
	for _, button := range m.activeButtons() {
		if button.Label == label && button.Disabled {
			return button
		}
	}
	return Button{}
}

// clickInput returns a click at the center of button.
func clickInput(button Button) Input {
	return Input{
		CursorX: button.X + button.W/2,
		CursorY: button.Y + button.H/2,
		Clicked: true,
	}
}
