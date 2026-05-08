package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ingameMenuPanelWidth  = 420
	ingameMenuPanelHeight = 270
	ingameMenuButtonWidth = 220
	ingameMenuButtonH     = 50
	ingameMenuButtonGap   = 18
)

// Action identifies game-level actions selected from game UI.
type Action int

const (
	// ActionNone means no game-level action was selected.
	ActionNone Action = iota
	// ActionSurrender means the app should leave the game and show the main menu.
	ActionSurrender
)

type ingameMenuButton struct {
	label  string
	x      int
	y      int
	w      int
	h      int
	action Action
}

var overlayColor = color.RGBA{R: 0, G: 0, B: 0, A: 128}

// contains reports whether the point is inside the in-game menu button bounds.
func (b ingameMenuButton) contains(x, y int) bool {
	return x >= b.x && x < b.x+b.w && y >= b.y && y < b.y+b.h
}

// layoutIngameMenu rebuilds overlay button rectangles from the drawable size.
func (s *State) layoutIngameMenu() {
	panelX := s.width/2 - ingameMenuPanelWidth/2
	panelY := s.height/2 - ingameMenuPanelHeight/2
	buttonX := panelX + ingameMenuPanelWidth/2 - ingameMenuButtonWidth/2
	buttonY := panelY + 118

	s.ingameMenuButtons = []ingameMenuButton{
		{label: "Resume", x: buttonX, y: buttonY, w: ingameMenuButtonWidth, h: ingameMenuButtonH, action: ActionNone},
		{label: "Surrender", x: buttonX, y: buttonY + ingameMenuButtonH + ingameMenuButtonGap, w: ingameMenuButtonWidth, h: ingameMenuButtonH, action: ActionSurrender},
	}
}

// openIngameMenu pauses the game and shows the overlay menu.
func (s *State) openIngameMenu() {
	s.pausedBeforeMenu = s.paused
	s.paused = true
	s.ingameMenuOpen = true
	s.ingameMenuHover = -1
}

// closeIngameMenu hides the overlay and restores the previous pause state.
func (s *State) closeIngameMenu() {
	s.ingameMenuOpen = false
	s.ingameMenuHover = -1
	s.paused = s.pausedBeforeMenu
}

// updateIngameMenu applies overlay input and returns selected game actions.
func (s *State) updateIngameMenu(input Input) Action {
	s.ingameMenuHover = s.ingameMenuButtonIndexAt(input.CursorX, input.CursorY)

	if input.ToggleMenu {
		s.closeIngameMenu()
		return ActionNone
	}
	if !input.Clicked {
		return ActionNone
	}

	switch action := s.ingameMenuActionAt(input.CursorX, input.CursorY); action {
	case ActionSurrender:
		return ActionSurrender
	default:
		if s.ingameMenuResumeContains(input.CursorX, input.CursorY) {
			s.closeIngameMenu()
		}
		return ActionNone
	}
}

// ingameMenuButtonIndexAt returns the index for the first overlay button at a point.
func (s *State) ingameMenuButtonIndexAt(x, y int) int {
	for i, button := range s.ingameMenuButtons {
		if button.contains(x, y) {
			return i
		}
	}
	return -1
}

// ingameMenuActionAt returns the action for the first overlay button at a point.
func (s *State) ingameMenuActionAt(x, y int) Action {
	for _, button := range s.ingameMenuButtons {
		if button.contains(x, y) {
			return button.action
		}
	}
	return ActionNone
}

// ingameMenuResumeContains reports whether the point is inside the Resume button.
func (s *State) ingameMenuResumeContains(x, y int) bool {
	if len(s.ingameMenuButtons) == 0 {
		return false
	}
	return s.ingameMenuButtons[0].contains(x, y)
}

// drawIngameMenu renders the overlay menu when it is open.
func (s *State) drawIngameMenu(screen *ebiten.Image) {
	if !s.ingameMenuOpen {
		return
	}

	vector.FillRect(screen, 0, 0, float32(s.width), float32(s.height), overlayColor, false)
	s.drawIngameMenuPanel(screen)
	s.drawIngameMenuButtons(screen)
}

// drawIngameMenuPanel renders the centered overlay panel.
func (s *State) drawIngameMenuPanel(screen *ebiten.Image) {
	panelX := float32(s.width/2 - ingameMenuPanelWidth/2)
	panelY := float32(s.height/2 - ingameMenuPanelHeight/2)

	vector.FillRect(screen, panelX, panelY, ingameMenuPanelWidth, ingameMenuPanelHeight, fieldColor, false)
	vector.StrokeRect(screen, panelX, panelY, ingameMenuPanelWidth, ingameMenuPanelHeight, 4, fieldEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, ingameMenuPanelWidth-24, ingameMenuPanelHeight-24, 1.5, fieldAccentColor, false)

	titleWidth, _ := text.Measure("Paused", s.ingameMenuTitleFace, s.ingameMenuTitleFace.Size)
	titleX := (float64(s.width) - titleWidth) / 2
	s.drawText(screen, "Paused", s.ingameMenuTitleFace, titleX, float64(panelY)+42, textColor)
}

// drawIngameMenuButtons renders overlay buttons with hover feedback.
func (s *State) drawIngameMenuButtons(screen *ebiten.Image) {
	for i, button := range s.ingameMenuButtons {
		fill := clearingColor
		edge := fieldEdgeColor
		if s.ingameMenuHover == i {
			fill = pauseColor
			edge = textColor
		}

		vector.FillRect(screen, float32(button.x), float32(button.y), float32(button.w), float32(button.h), fill, false)
		vector.StrokeRect(screen, float32(button.x), float32(button.y), float32(button.w), float32(button.h), 3, edge, false)

		labelWidth, _ := text.Measure(button.label, s.ingameMenuButtonFace, s.ingameMenuButtonFace.Size)
		labelX := float64(button.x) + (float64(button.w)-labelWidth)/2
		s.drawText(screen, button.label, s.ingameMenuButtonFace, labelX, float64(button.y+10), textColor)
	}
}
