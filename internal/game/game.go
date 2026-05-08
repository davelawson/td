package game

import (
	"bytes"
	"fmt"
	"image/color"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	width      int
	height     int
	updates    int
	paused     bool
	menu       ingameMenu
	phase      phase
	chapter    string
	day        int
	calmTime   int
	raidCount  int
	barricade  int
	resources  resourceCounts
	titleFace  *text.GoTextFace
	bodyFace   *text.GoTextFace
	hudFace    *text.GoTextFace
}

var (
	backgroundColor  = ui.CharcoalBlack
	fieldColor       = ui.PineGreen
	fieldEdgeColor   = ui.Bronze
	textColor        = ui.Parchment
	mutedTextColor   = ui.MutedParchment
	pauseColor       = ui.LightBronze
	fieldAccentColor = ui.Purple
	pathColor        = ui.OliveBrown
	clearingColor    = ui.MossGreen
)

// New creates the initial game state for a Wizard name.
func New(wizardName string, width, height int) (*State, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	state := &State{
		wizardName: wizardName,
		width:      width,
		height:     height,
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
	state.setPrototypeGameStatus()
	state.layoutIngameMenu()
	return state, nil
}

// Resize updates drawable dimensions for game rendering.
func (s *State) Resize(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	s.width = width
	s.height = height
	s.layoutIngameMenu()
}

// Update applies game input and advances one logical update when unpaused.
func (s *State) Update(input Input) Action {
	if s.menu.open {
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
	s.drawPrototypeField(screen)
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
	return s.menu.open
}

// WizardName returns the Wizard name used to start the game.
func (s *State) WizardName() string {
	return s.wizardName
}

// drawPrototypeField renders the first placeholder game scene.
func (s *State) drawPrototypeField(screen *ebiten.Image) {
	fieldW := float32(820)
	fieldH := float32(460)
	fieldX := float32(s.width)/2 - fieldW/2
	fieldY := float32(s.height)/2 - fieldH/2

	vector.FillRect(screen, fieldX, fieldY, fieldW, fieldH, fieldColor, false)
	vector.StrokeRect(screen, fieldX, fieldY, fieldW, fieldH, 4, fieldEdgeColor, false)
	vector.StrokeRect(screen, fieldX+18, fieldY+18, fieldW-36, fieldH-36, 1.5, fieldAccentColor, false)

	pathY := fieldY + fieldH/2
	vector.StrokeLine(screen, fieldX+70, pathY, fieldX+fieldW-70, pathY, 10, pathColor, false)
	vector.FillCircle(screen, fieldX+fieldW/2, pathY, 42, clearingColor, false)
	vector.StrokeCircle(screen, fieldX+fieldW/2, pathY, 42, 3, fieldEdgeColor, false)
}

// drawWizardName renders the active Wizard name.
func (s *State) drawWizardName(screen *ebiten.Image) {
	value := fmt.Sprintf("Wizard %s", s.wizardName)
	s.drawText(screen, value, s.titleFace, 56, 112, textColor)
	s.drawText(screen, "The first defenses are waiting for orders.", s.bodyFace, 56, 156, mutedTextColor)
}

// drawCounter renders update and pause status in the top-right corner.
func (s *State) drawCounter(screen *ebiten.Image) {
	value := fmt.Sprintf("Updates: %d", s.updates)
	width, _ := text.Measure(value, s.bodyFace, s.bodyFace.Size)
	x := float64(s.width) - width - 48
	s.drawText(screen, value, s.bodyFace, x, float64(s.height)-58, mutedTextColor)

	if s.paused {
		pauseWidth, _ := text.Measure("PAUSED", s.bodyFace, s.bodyFace.Size)
		s.drawText(screen, "PAUSED", s.bodyFace, float64(s.width)-pauseWidth-48, float64(s.height)-94, pauseColor)
	}
}

// drawText draws one line at the given coordinates.
func (s *State) drawText(screen *ebiten.Image, value string, face *text.GoTextFace, x, y float64, clr color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}
