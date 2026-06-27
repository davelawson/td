package menu

import (
	"bytes"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

// Action identifies the behavior selected from the main menu.
type Action int

const (
	mainPanelWidth       = 520
	mainPanelHeight      = 398
	screenPanelWidth     = 840
	screenPanelHeight    = 340
	menuButtonWidth      = 220
	menuButtonHeight     = 44
	menuButtonGap        = 10
	settingsButtonHeight = 54
)

const (
	// ActionNone means no menu action was selected.
	ActionNone Action = iota
	// ActionNew means the app should show the new-game configuration screen.
	ActionNew
	// ActionSettings means the app should show the settings placeholder.
	ActionSettings
	// ActionBack means the app should return to the main menu.
	ActionBack
	// ActionQuit means the app should terminate cleanly.
	ActionQuit
	// ActionStart means the app should leave the menu and begin the game.
	ActionStart
)

// Button describes a rectangular menu target and the action it selects.
type Button = ui.Button[Action]

// Screen identifies the current menu-owned screen.
type Screen int

const (
	// ScreenMain is the top-level main menu.
	ScreenMain Screen = iota
	// ScreenNewGame is the configuration screen reached from New.
	ScreenNewGame
	// ScreenSettings is the placeholder reached from Settings.
	ScreenSettings
)

// Input describes one frame of menu input collected by the executable.
type Input struct {
	CursorX   int
	CursorY   int
	Clicked   bool
	Typed     []rune
	Backspace bool
}

// Menu owns the menu screen state, input routing, and rendering.
type Menu struct {
	width              int
	height             int
	screen             Screen
	mainButtons        []Button
	settingsBackButton Button
	newGameButtons     []Button
	hoverAction        Action
	wizardName         string
	wizardNameFocused  bool
	titleFace          *text.GoTextFace
	bodyFace           *text.GoTextFace
	buttonFace         *text.GoTextFace
	nameFace           *text.GoTextFace
}

// New creates the menu state and font faces.
func New(width, height int) (*Menu, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	menu := &Menu{
		width:  width,
		height: height,
		titleFace: &text.GoTextFace{
			Source: source,
			Size:   ui.MenuTitleFontSize,
		},
		bodyFace: &text.GoTextFace{
			Source: source,
			Size:   ui.MenuBodyFontSize,
		},
		buttonFace: &text.GoTextFace{
			Source: source,
			Size:   ui.MenuButtonFontSize,
		},
		nameFace: &text.GoTextFace{
			Source: source,
			Size:   ui.MenuNameFontSize,
		},
	}
	menu.layoutButtons()
	return menu, nil
}

// ActionAt returns the first button action containing the point.
func ActionAt(buttons []Button, x, y int) Action {
	for _, button := range buttons {
		if button.Disabled {
			continue
		}
		if button.Contains(x, y) {
			return button.Action
		}
	}
	return ActionNone
}

// Update applies pointer state to the menu and returns the selected action.
func (m *Menu) Update(input Input) Action {
	m.hoverAction = ActionAt(m.activeButtons(), input.CursorX, input.CursorY)
	if m.updateWizardName(input) {
		m.layoutButtons()
	}

	if !input.Clicked {
		return ActionNone
	}

	m.updateWizardNameFocus(input.CursorX, input.CursorY)
	action := ActionAt(m.activeButtons(), input.CursorX, input.CursorY)
	m.handleAction(action)
	return action
}

// Resize updates menu geometry for a changed drawable size.
func (m *Menu) Resize(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	if m.width == width && m.height == height {
		return
	}
	m.width = width
	m.height = height
	m.layoutButtons()
	m.hoverAction = ActionNone
}

// Screen returns the current menu screen.
func (m *Menu) Screen() Screen {
	return m.screen
}

// SetScreenForTest sets the current menu screen for screenshot capture and tests.
func (m *Menu) SetScreenForTest(screen Screen) {
	m.screen = screen
	m.wizardNameFocused = screen == ScreenNewGame
	m.hoverAction = ActionNone
}

// ResetToMain returns the menu to its top-level screen.
func (m *Menu) ResetToMain() {
	m.screen = ScreenMain
	m.wizardNameFocused = false
	m.hoverAction = ActionNone
}

// WizardName returns the current new-game Wizard name.
func (m *Menu) WizardName() string {
	return m.wizardName
}

// WizardNameFocused reports whether the Wizard name field is active.
func (m *Menu) WizardNameFocused() bool {
	return m.wizardNameFocused
}

// Draw renders the current menu screen.
func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(colors.background)
	m.drawBackdrop(screen)
	switch m.screen {
	case ScreenNewGame:
		m.drawNewGamePanel(screen)
		m.drawButtons(screen, m.activeButtons())
	case ScreenSettings:
		m.drawSettingsPanel(screen)
		m.drawButtons(screen, m.activeButtons())
	default:
		m.drawMenuPanel(screen)
		m.drawButtons(screen, m.activeButtons())
	}
}

// handleAction applies a selected action to menu-owned screen state.
func (m *Menu) handleAction(action Action) {
	switch action {
	case ActionNew:
		m.screen = ScreenNewGame
		m.wizardNameFocused = true
	case ActionSettings:
		m.screen = ScreenSettings
		m.wizardNameFocused = false
	case ActionBack:
		m.screen = ScreenMain
		m.wizardNameFocused = false
	}
}

// activeButtons returns the buttons available on the current screen.
func (m *Menu) activeButtons() []Button {
	switch m.screen {
	case ScreenNewGame:
		return m.newGameButtons
	case ScreenSettings:
		return []Button{m.settingsBackButton}
	default:
		return m.mainButtons
	}
}

// layoutButtons rebuilds button rectangles from the current drawable size.
func (m *Menu) layoutButtons() {
	centerX := m.width / 2
	mainButtonX := centerX - menuButtonWidth/2
	mainButtonY := m.mainPanelY() + 170
	m.mainButtons = []Button{
		{Label: "New", X: mainButtonX, Y: mainButtonY, W: menuButtonWidth, H: menuButtonHeight, Action: ActionNew},
		{Label: "Load", X: mainButtonX, Y: mainButtonY + menuButtonHeight + menuButtonGap, W: menuButtonWidth, H: menuButtonHeight, Disabled: true},
		{Label: "Settings", X: mainButtonX, Y: mainButtonY + 2*(menuButtonHeight+menuButtonGap), W: menuButtonWidth, H: menuButtonHeight, Action: ActionSettings},
		{Label: "Quit", X: mainButtonX, Y: mainButtonY + 3*(menuButtonHeight+menuButtonGap), W: menuButtonWidth, H: menuButtonHeight, Action: ActionQuit},
	}

	screenButtonY := m.screenPanelY() + 258
	m.settingsBackButton = Button{
		Label:  "Back",
		X:      centerX - menuButtonWidth/2,
		Y:      screenButtonY,
		W:      menuButtonWidth,
		H:      settingsButtonHeight,
		Action: ActionBack,
	}
	m.layoutStartButtons(screenButtonY)
}

// drawBackdrop paints simple fantasy accents behind the menu.
func (m *Menu) drawBackdrop(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(m.width), 82, colors.backdropBand, false)
	vector.FillRect(screen, 0, float32(m.height-82), float32(m.width), 82, colors.backdropBand, false)

	for i := 0; i < 6; i++ {
		x := float32(m.width/2 - 395 + i*145)
		vector.StrokeRect(screen, x, 102, 46, 46, 2, colors.transparentAccent, true)
		vector.StrokeRect(screen, x+9, 111, 28, 28, 2, colors.transparentEdge, true)
	}
}

// mainPanelY returns the top edge for the main menu panel.
func (m *Menu) mainPanelY() int {
	return (m.height - mainPanelHeight) / 2
}

// screenPanelY returns the top edge for secondary menu panels.
func (m *Menu) screenPanelY() int {
	return (m.height - screenPanelHeight) / 2
}

// drawMenuPanel renders the title area and menu copy.
func (m *Menu) drawMenuPanel(screen *ebiten.Image) {
	panelX := float32(m.width/2 - mainPanelWidth/2)
	panelY := float32(m.mainPanelY())
	panelW := float32(mainPanelWidth)
	panelH := float32(mainPanelHeight)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, colors.panel, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, colors.panelEdge, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, colors.accent, false)

	ui.DrawCenteredText(screen, m.width, "td", m.titleFace, float64(m.mainPanelY()+40), colors.text)
	ui.DrawCenteredText(screen, m.width, "Arcane defenses await their first command.", m.bodyFace, float64(m.mainPanelY()+132), colors.mutedText)
}

// drawButtons renders menu buttons with hover feedback.
func (m *Menu) drawButtons(screen *ebiten.Image, buttons []Button) {
	for _, button := range buttons {
		fill := colors.button
		edge := colors.panelEdge
		labelColor := colors.text
		if button.Disabled {
			fill = colors.disabled
			edge = colors.disabledButtonEdge
			labelColor = colors.mutedText
		} else if m.hoverAction != ActionNone && m.hoverAction == button.Action {
			fill = colors.hover
			edge = colors.text
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)
		ui.DrawCenteredTextInRect(screen, button.Label, m.buttonFace, button.X, button.Y, button.W, 9, labelColor)
	}
}

// drawSettingsPanel renders the temporary settings screen.
func (m *Menu) drawSettingsPanel(screen *ebiten.Image) {
	panelX := float32(m.width/2 - screenPanelWidth/2)
	panelY := float32(m.screenPanelY())
	panelW := float32(screenPanelWidth)
	panelH := float32(screenPanelHeight)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, colors.panel, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, colors.panelEdge, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, colors.accent, false)

	ui.DrawCenteredText(screen, m.width, "Settings", m.titleFace, float64(m.screenPanelY()+86), colors.text)
}
