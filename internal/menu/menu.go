package menu

import (
	"bytes"
	"image/color"
	"unicode"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

// Action identifies the behavior selected from the main menu.
type Action int

const (
	wizardNameMaxRunes = 16
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
)

// Button describes a rectangular menu target and the action it selects.
type Button struct {
	Label    string
	X        int
	Y        int
	W        int
	H        int
	Action   Action
	Disabled bool
}

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

var (
	backgroundColor = color.RGBA{R: 18, G: 19, B: 17, A: 255}
	panelColor      = color.RGBA{R: 45, G: 58, B: 49, A: 255}
	panelEdgeColor  = color.RGBA{R: 134, G: 114, B: 65, A: 255}
	textColor       = color.RGBA{R: 238, G: 224, B: 188, A: 255}
	mutedTextColor  = color.RGBA{R: 184, G: 172, B: 139, A: 255}
	hoverColor      = color.RGBA{R: 150, G: 124, B: 49, A: 255}
	buttonColor     = color.RGBA{R: 74, G: 83, B: 68, A: 255}
	disabledColor   = color.RGBA{R: 51, G: 57, B: 51, A: 255}
	accentColor     = color.RGBA{R: 98, G: 90, B: 145, A: 255}
)

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
}

// New creates the menu state and font faces.
func New(width, height int) (*Menu, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &Menu{
		width:  width,
		height: height,
		mainButtons: []Button{
			{Label: "New", X: width/2 - 110, Y: 252, W: 220, H: 44, Action: ActionNew},
			{Label: "Load", X: width/2 - 110, Y: 306, W: 220, H: 44, Disabled: true},
			{Label: "Settings", X: width/2 - 110, Y: 360, W: 220, H: 44, Action: ActionSettings},
			{Label: "Quit", X: width/2 - 110, Y: 414, W: 220, H: 44, Action: ActionQuit},
		},
		settingsBackButton: Button{
			Label:  "Back",
			X:      width/2 - 110,
			Y:      384,
			W:      220,
			H:      54,
			Action: ActionBack,
		},
		newGameButtons: []Button{
			{Label: "Cancel", X: width/2 - 230, Y: 386, W: 180, H: 52, Action: ActionBack},
			{Label: "Start", X: width/2 + 50, Y: 386, W: 180, H: 52, Disabled: true},
		},
		titleFace: &text.GoTextFace{
			Source: source,
			Size:   74,
		},
		bodyFace: &text.GoTextFace{
			Source: source,
			Size:   24,
		},
		buttonFace: &text.GoTextFace{
			Source: source,
			Size:   30,
		},
	}, nil
}

// Contains reports whether the point is inside the button bounds.
func (b Button) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.W && y >= b.Y && y < b.Y+b.H
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
	m.updateWizardName(input)

	if !input.Clicked {
		return ActionNone
	}

	m.updateWizardNameFocus(input.CursorX, input.CursorY)
	action := ActionAt(m.activeButtons(), input.CursorX, input.CursorY)
	m.handleAction(action)
	return action
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
	screen.Fill(backgroundColor)
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

// updateWizardName applies text edits while the new-game name field is focused.
func (m *Menu) updateWizardName(input Input) {
	if m.screen != ScreenNewGame || !m.wizardNameFocused {
		return
	}

	if input.Backspace {
		m.wizardName = removeLastRune(m.wizardName)
	}
	for _, value := range input.Typed {
		if !unicode.IsPrint(value) {
			continue
		}
		if utf8.RuneCountInString(m.wizardName) >= wizardNameMaxRunes {
			return
		}
		m.wizardName += string(value)
	}
}

// updateWizardNameFocus handles pointer focus for the new-game name field.
func (m *Menu) updateWizardNameFocus(cursorX, cursorY int) {
	if m.screen != ScreenNewGame {
		return
	}
	m.wizardNameFocused = m.wizardNameFieldContains(cursorX, cursorY)
}

// removeLastRune returns value without its final rune.
func removeLastRune(value string) string {
	if value == "" {
		return ""
	}
	_, size := utf8.DecodeLastRuneInString(value)
	return value[:len(value)-size]
}

// wizardNameFieldContains reports whether a point is inside the name field.
func (m *Menu) wizardNameFieldContains(x, y int) bool {
	fieldX, fieldY, fieldW, fieldH := m.wizardNameFieldBounds()
	return x >= fieldX && x < fieldX+fieldW && y >= fieldY && y < fieldY+fieldH
}

// wizardNameFieldBounds returns the New Game Wizard name field rectangle.
func (m *Menu) wizardNameFieldBounds() (int, int, int, int) {
	return m.width/2 - 220, 282, 440, 54
}

// drawBackdrop paints simple fantasy accents behind the menu.
func (m *Menu) drawBackdrop(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(m.width), 82, color.RGBA{R: 26, G: 32, B: 28, A: 255}, false)
	vector.FillRect(screen, 0, float32(m.height-82), float32(m.width), 82, color.RGBA{R: 26, G: 31, B: 27, A: 255}, false)

	for i := 0; i < 6; i++ {
		x := float32(110 + i*145)
		vector.StrokeRect(screen, x, 102, 46, 46, 2, color.RGBA{R: 65, G: 60, B: 94, A: 130}, true)
		vector.StrokeRect(screen, x+9, 111, 28, 28, 2, color.RGBA{R: 111, G: 96, B: 58, A: 115}, true)
	}
}

// drawMenuPanel renders the title area and menu copy.
func (m *Menu) drawMenuPanel(screen *ebiten.Image) {
	panelX := float32(220)
	panelY := float32(82)
	panelW := float32(520)
	panelH := float32(398)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	m.drawCenteredText(screen, "td", m.titleFace, 122, textColor)
	m.drawCenteredText(screen, "Arcane defenses await their first command.", m.bodyFace, 214, mutedTextColor)
}

// drawButtons renders menu buttons with hover feedback.
func (m *Menu) drawButtons(screen *ebiten.Image, buttons []Button) {
	for _, button := range buttons {
		fill := buttonColor
		edge := panelEdgeColor
		labelColor := textColor
		if button.Disabled {
			fill = disabledColor
			edge = color.RGBA{R: 83, G: 84, B: 73, A: 255}
			labelColor = mutedTextColor
		} else if m.hoverAction != ActionNone && m.hoverAction == button.Action {
			fill = hoverColor
			edge = textColor
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)
		m.drawCenteredButtonText(screen, button, labelColor)
	}
}

// drawNewGamePanel renders the new-game configuration screen.
func (m *Menu) drawNewGamePanel(screen *ebiten.Image) {
	panelX := float32(180)
	panelY := float32(128)
	panelW := float32(600)
	panelH := float32(340)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	m.drawCenteredText(screen, "New Game", m.titleFace, 158, textColor)
	m.drawText(screen, "Wizard Name", m.bodyFace, float64(m.width/2-220), 248, mutedTextColor)
	m.drawWizardNameField(screen)
}

// drawSettingsPanel renders the temporary settings screen.
func (m *Menu) drawSettingsPanel(screen *ebiten.Image) {
	panelX := float32(180)
	panelY := float32(128)
	panelW := float32(600)
	panelH := float32(340)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	m.drawCenteredText(screen, "Settings", m.titleFace, 214, textColor)
}

// drawWizardNameField renders the editable Wizard name field.
func (m *Menu) drawWizardNameField(screen *ebiten.Image) {
	fieldX, fieldY, fieldW, fieldH := m.wizardNameFieldBounds()
	edge := panelEdgeColor
	if m.wizardNameFocused {
		edge = textColor
	}

	vector.FillRect(screen, float32(fieldX), float32(fieldY), float32(fieldW), float32(fieldH), color.RGBA{R: 31, G: 36, B: 32, A: 255}, false)
	vector.StrokeRect(screen, float32(fieldX), float32(fieldY), float32(fieldW), float32(fieldH), 3, edge, false)

	value := m.wizardName
	labelColor := textColor
	if value == "" {
		value = "Enter name"
		labelColor = mutedTextColor
	}
	m.drawText(screen, value, m.buttonFace, float64(fieldX+18), float64(fieldY+11), labelColor)

	if m.wizardNameFocused {
		textWidth, _ := text.Measure(m.wizardName, m.buttonFace, m.buttonFace.Size)
		caretX := float32(fieldX + 19 + int(textWidth))
		vector.StrokeLine(screen, caretX, float32(fieldY+12), caretX, float32(fieldY+42), 2, textColor, false)
	}
}

// drawCenteredText draws one line centered horizontally at the given y coordinate.
func (m *Menu) drawCenteredText(screen *ebiten.Image, value string, face *text.GoTextFace, y float64, clr color.Color) {
	width, _ := text.Measure(value, face, face.Size)
	m.drawText(screen, value, face, (float64(m.width)-width)/2, y, clr)
}

// drawCenteredButtonText draws a button label centered inside its button bounds.
func (m *Menu) drawCenteredButtonText(screen *ebiten.Image, button Button, clr color.Color) {
	width, _ := text.Measure(button.Label, m.buttonFace, m.buttonFace.Size)
	x := float64(button.X) + (float64(button.W)-width)/2
	m.drawText(screen, button.Label, m.buttonFace, x, float64(button.Y+9), clr)
}

// drawText draws one line at the given coordinates.
func (m *Menu) drawText(screen *ebiten.Image, value string, face *text.GoTextFace, x, y float64, clr color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}
