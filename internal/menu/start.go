package menu

import (
	"unicode"
	"unicode/utf8"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	wizardNameMaxRunes = 32
)

const (
	screenButtonWidth  = 180
	screenButtonHeight = 52
	screenButtonGap    = 100
	nameFieldWidth     = 680
	nameFieldHeight    = 54
)

// layoutStartButtons rebuilds the New Game screen button rectangles.
func (m *Menu) layoutStartButtons(screenButtonY int) {
	centerX := m.width / 2
	totalButtonWidth := 2*screenButtonWidth + screenButtonGap
	m.newGameButtons = []Button{
		{Label: "Cancel", X: centerX - totalButtonWidth/2, Y: screenButtonY, W: screenButtonWidth, H: screenButtonHeight, Action: ActionBack},
		{Label: "Start", X: centerX - totalButtonWidth/2 + screenButtonWidth + screenButtonGap, Y: screenButtonY, W: screenButtonWidth, H: screenButtonHeight, Action: ActionStart, Disabled: m.wizardName == ""},
	}
}

// updateWizardName applies text edits and reports whether the name changed.
func (m *Menu) updateWizardName(input Input) bool {
	if m.screen != ScreenNewGame || !m.wizardNameFocused {
		return false
	}

	changed := false
	if input.Backspace {
		next := removeLastRune(m.wizardName)
		changed = next != m.wizardName
		m.wizardName = next
	}
	for _, value := range input.Typed {
		if !unicode.IsPrint(value) {
			continue
		}
		if utf8.RuneCountInString(m.wizardName) >= wizardNameMaxRunes {
			return changed
		}
		m.wizardName += string(value)
		changed = true
	}
	return changed
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
	return m.width/2 - nameFieldWidth/2, m.screenPanelY() + 154, nameFieldWidth, nameFieldHeight
}

// drawNewGamePanel renders the new-game configuration screen.
func (m *Menu) drawNewGamePanel(screen *ebiten.Image) {
	panelX := float32(m.width/2 - screenPanelWidth/2)
	panelY := float32(m.screenPanelY())
	panelW := float32(screenPanelWidth)
	panelH := float32(screenPanelHeight)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, colors.panel, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, colors.panelEdge, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, colors.accent, false)

	ui.DrawCenteredText(screen, m.width, "New Game", m.titleFace, float64(m.screenPanelY()+30), colors.text)
	ui.DrawText(screen, "Wizard Name", m.bodyFace, float64(m.width/2-nameFieldWidth/2), float64(m.screenPanelY()+120), colors.mutedText)
	m.drawWizardNameField(screen)
}

// drawWizardNameField renders the editable Wizard name field.
func (m *Menu) drawWizardNameField(screen *ebiten.Image) {
	fieldX, fieldY, fieldW, fieldH := m.wizardNameFieldBounds()
	edge := colors.panelEdge
	if m.wizardNameFocused {
		edge = colors.text
	}

	vector.FillRect(screen, float32(fieldX), float32(fieldY), float32(fieldW), float32(fieldH), colors.nameField, false)
	vector.StrokeRect(screen, float32(fieldX), float32(fieldY), float32(fieldW), float32(fieldH), 3, edge, false)

	value := m.wizardName
	labelColor := colors.text
	if value == "" {
		value = "Enter name"
		labelColor = colors.mutedText
	}
	ui.DrawText(screen, value, m.nameFace, float64(fieldX+18), float64(fieldY+15), labelColor)

	if m.wizardNameFocused {
		textWidth, _ := text.Measure(m.wizardName, m.nameFace, m.nameFace.Size)
		caretX := float32(fieldX + 19 + int(textWidth))
		vector.StrokeLine(screen, caretX, float32(fieldY+12), caretX, float32(fieldY+42), 2, colors.text, false)
	}
}
