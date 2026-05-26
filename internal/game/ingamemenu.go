package game

import (
	"td/internal/ui"

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

type ingameMenu struct {
	pausedBeforeOpen bool
	open             bool
	buttons          []ui.Button[Action]
	hover            int
	titleFace        *text.GoTextFace
	buttonFace       *text.GoTextFace
}

// layoutIngameMenu rebuilds overlay button rectangles from the drawable size.
func (s *State) layoutIngameMenu() {
	panelX := s.ui.width/2 - ingameMenuPanelWidth/2
	panelY := s.ui.height/2 - ingameMenuPanelHeight/2
	buttonX := panelX + ingameMenuPanelWidth/2 - ingameMenuButtonWidth/2
	buttonY := panelY + 118

	s.ui.menu.buttons = []ui.Button[Action]{
		{Label: "Resume", X: buttonX, Y: buttonY, W: ingameMenuButtonWidth, H: ingameMenuButtonH, Action: ActionNone},
		{Label: "Surrender", X: buttonX, Y: buttonY + ingameMenuButtonH + ingameMenuButtonGap, W: ingameMenuButtonWidth, H: ingameMenuButtonH, Action: ActionSurrender},
	}
}

// openIngameMenu pauses the game and shows the overlay menu.
func (s *State) openIngameMenu() {
	s.ui.menu.pausedBeforeOpen = s.paused
	s.paused = true
	s.ui.menu.open = true
	s.ui.menu.hover = -1
	s.ui.buildBarHover = -1
}

// closeIngameMenu hides the overlay and restores the previous pause state.
func (s *State) closeIngameMenu() {
	s.ui.menu.open = false
	s.ui.menu.hover = -1
	s.paused = s.ui.menu.pausedBeforeOpen
}

// updateIngameMenu applies overlay input and returns selected game actions.
func (s *State) updateIngameMenu(input Input) Action {
	s.ui.menu.hover = s.ingameMenuButtonIndexAt(input.CursorX, input.CursorY)

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
	for i, button := range s.ui.menu.buttons {
		if button.Contains(x, y) {
			return i
		}
	}
	return -1
}

// ingameMenuActionAt returns the action for the first overlay button at a point.
func (s *State) ingameMenuActionAt(x, y int) Action {
	for _, button := range s.ui.menu.buttons {
		if button.Contains(x, y) {
			return button.Action
		}
	}
	return ActionNone
}

// ingameMenuResumeContains reports whether the point is inside the Resume button.
func (s *State) ingameMenuResumeContains(x, y int) bool {
	if len(s.ui.menu.buttons) == 0 {
		return false
	}
	return s.ui.menu.buttons[0].Contains(x, y)
}

// drawIngameMenu renders the overlay menu when it is open.
func (s *State) drawIngameMenu(screen *ebiten.Image) {
	if !s.ui.menu.open {
		return
	}

	vector.FillRect(screen, 0, 0, float32(s.ui.width), float32(s.ui.height), colors.overlay, false)
	s.drawIngameMenuPanel(screen)
	s.drawIngameMenuButtons(screen)
}

// drawIngameMenuPanel renders the centered overlay panel.
func (s *State) drawIngameMenuPanel(screen *ebiten.Image) {
	panelX := float32(s.ui.width/2 - ingameMenuPanelWidth/2)
	panelY := float32(s.ui.height/2 - ingameMenuPanelHeight/2)

	vector.FillRect(screen, panelX, panelY, ingameMenuPanelWidth, ingameMenuPanelHeight, colors.field, false)
	vector.StrokeRect(screen, panelX, panelY, ingameMenuPanelWidth, ingameMenuPanelHeight, 4, colors.fieldEdge, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, ingameMenuPanelWidth-24, ingameMenuPanelHeight-24, 1.5, colors.fieldAccent, false)

	titleWidth, _ := text.Measure("Paused", s.ui.menu.titleFace, s.ui.menu.titleFace.Size)
	titleX := (float64(s.ui.width) - titleWidth) / 2
	ui.DrawText(screen, "Paused", s.ui.menu.titleFace, titleX, float64(panelY)+42, colors.text)
}

// drawIngameMenuButtons renders overlay buttons with hover feedback.
func (s *State) drawIngameMenuButtons(screen *ebiten.Image) {
	for i, button := range s.ui.menu.buttons {
		fill := colors.clearing
		edge := colors.fieldEdge
		if s.ui.menu.hover == i {
			fill = colors.pause
			edge = colors.text
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)

		labelWidth, _ := text.Measure(button.Label, s.ui.menu.buttonFace, s.ui.menu.buttonFace.Size)
		labelX := float64(button.X) + (float64(button.W)-labelWidth)/2
		ui.DrawText(screen, button.Label, s.ui.menu.buttonFace, labelX, float64(button.Y+10), colors.text)
	}
}
