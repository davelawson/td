package game

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Input describes one frame of game input collected by the executable.
type Input struct {
	TogglePause bool
	ToggleMenu  bool
	CursorX     int
	CursorY     int
	Clicked     bool
}

// State owns the current game state and logical update rules.
type State struct {
	wizardName string
	updates    int
	paused     bool
	gameMap    Map
	status     gameStatus
	ui         gameUI
}

type gameUI struct {
	width     int
	height    int
	menu      ingameMenu
	titleFace *text.GoTextFace
	bodyFace  *text.GoTextFace
	hudFace   *text.GoTextFace
}

// New creates the initial game state for a Wizard name.
func New(wizardName string, width, height int) (*State, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	state := &State{
		wizardName: wizardName,
		gameMap:    NewDefaultMap(),
		ui:         newGameUI(source, width, height),
	}
	state.setPrototypeGameStatus()
	state.layoutIngameMenu()
	return state, nil
}

// newGameUI creates render-facing state for the current drawable size.
func newGameUI(source *text.GoTextFaceSource, width, height int) gameUI {
	return gameUI{
		width:  width,
		height: height,
		titleFace: &text.GoTextFace{
			Source: source,
			Size:   34,
		},
		bodyFace: &text.GoTextFace{
			Source: source,
			Size:   24,
		},
		hudFace: &text.GoTextFace{
			Source: source,
			Size:   22,
		},
		menu: ingameMenu{
			titleFace: &text.GoTextFace{
				Source: source,
				Size:   46,
			},
			buttonFace: &text.GoTextFace{
				Source: source,
				Size:   28,
			},
		},
	}
}

// Resize updates drawable dimensions for game rendering.
func (s *State) Resize(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	s.ui.width = width
	s.ui.height = height
	s.layoutIngameMenu()
}

// Update applies game input and advances one logical update when unpaused.
func (s *State) Update(input Input) Action {
	if s.ui.menu.open {
		return s.updateIngameMenu(input)
	}
	if input.ToggleMenu {
		s.openIngameMenu()
		return ActionNone
	}
	if input.TogglePause {
		s.paused = !s.paused
		return ActionNone
	}
	if s.paused {
		return ActionNone
	}
	s.updates++
	return ActionNone
}

// Draw renders the current game screen.
func (s *State) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	s.drawHomePlot(screen)
	s.drawTopBar(screen)
	s.drawWizardName(screen)
	s.drawCounter(screen)
	s.drawIngameMenu(screen)
}

// Updates returns the number of logical game updates processed.
func (s *State) Updates() int {
	return s.updates
}

// Paused reports whether logical updates are currently paused.
func (s *State) Paused() bool {
	return s.paused
}

// IngameMenuOpen reports whether the in-game overlay menu is visible.
func (s *State) IngameMenuOpen() bool {
	return s.ui.menu.open
}

// WizardName returns the Wizard name used to start the game.
func (s *State) WizardName() string {
	return s.wizardName
}

// drawWizardName renders the active Wizard name.
func (s *State) drawWizardName(screen *ebiten.Image) {
	value := fmt.Sprintf("Wizard %s", s.wizardName)
	s.drawText(screen, value, s.ui.titleFace, 56, 112, textColor)
	s.drawText(screen, "The first defenses are waiting for orders.", s.ui.bodyFace, 56, 156, mutedTextColor)
}

// drawCounter renders update and pause status in the top-right corner.
func (s *State) drawCounter(screen *ebiten.Image) {
	value := fmt.Sprintf("Updates: %d", s.updates)
	width, _ := text.Measure(value, s.ui.bodyFace, s.ui.bodyFace.Size)
	x := float64(s.ui.width) - width - 48
	s.drawText(screen, value, s.ui.bodyFace, x, float64(s.ui.height)-58, mutedTextColor)

	if s.paused {
		pauseWidth, _ := text.Measure("PAUSED", s.ui.bodyFace, s.ui.bodyFace.Size)
		s.drawText(screen, "PAUSED", s.ui.bodyFace, float64(s.ui.width)-pauseWidth-48, float64(s.ui.height)-94, pauseColor)
	}
}

// drawText draws one line at the given coordinates.
func (s *State) drawText(screen *ebiten.Image, value string, face *text.GoTextFace, x, y float64, clr color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}
