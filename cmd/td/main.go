package main

import (
	"bytes"
	"errors"
	"image/color"
	"log"

	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	screenWidth  = 960
	screenHeight = 540
)

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

type screenMode int

const (
	screenMainMenu screenMode = iota
	screenNewGame
	screenSettings
)

type game struct {
	screen      screenMode
	mainButtons []menu.Button
	backButton  menu.Button
	hoverAction menu.Action
	titleFace   *text.GoTextFace
	bodyFace    *text.GoTextFace
	buttonFace  *text.GoTextFace
}

// main starts the Ebitengine desktop application.
func main() {
	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game, err := newGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}

// newGame creates the menu state and font faces used by the app.
func newGame() (*game, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &game{
		mainButtons: []menu.Button{
			{Label: "New", X: screenWidth/2 - 110, Y: 252, W: 220, H: 44, Action: menu.ActionNew},
			{Label: "Load", X: screenWidth/2 - 110, Y: 306, W: 220, H: 44, Disabled: true},
			{Label: "Settings", X: screenWidth/2 - 110, Y: 360, W: 220, H: 44, Action: menu.ActionSettings},
			{Label: "Quit", X: screenWidth/2 - 110, Y: 414, W: 220, H: 44, Action: menu.ActionQuit},
		},
		backButton: menu.Button{
			Label:  "Back",
			X:      screenWidth/2 - 110,
			Y:      384,
			W:      220,
			H:      54,
			Action: menu.ActionBack,
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

// Update handles pointer input and returns a clean termination signal on quit.
func (g *game) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	g.hoverAction = menu.ActionAt(g.activeButtons(), cursorX, cursorY)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return g.handleAction(menu.ActionAt(g.activeButtons(), cursorX, cursorY))
	}
	return nil
}

// handleAction applies a selected menu action to the current game state.
func (g *game) handleAction(action menu.Action) error {
	switch action {
	case menu.ActionNew:
		g.screen = screenNewGame
	case menu.ActionSettings:
		g.screen = screenSettings
	case menu.ActionBack:
		g.screen = screenMainMenu
	case menu.ActionQuit:
		return ebiten.Termination
	}
	return nil
}

// Draw renders the first main-menu screen.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.drawBackdrop(screen)
	switch g.screen {
	case screenNewGame:
		g.drawPlaceholderPanel(screen, "New Game", "The first expedition is not prepared yet.")
		g.drawButtons(screen, g.activeButtons())
	case screenSettings:
		g.drawSettingsPanel(screen)
		g.drawButtons(screen, g.activeButtons())
	default:
		g.drawMenuPanel(screen)
		g.drawButtons(screen, g.activeButtons())
	}
}

// Layout returns the fixed logical resolution for the prototype.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// drawBackdrop paints simple fantasy accents behind the menu.
func (g *game) drawBackdrop(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, screenWidth, 82, color.RGBA{R: 26, G: 32, B: 28, A: 255}, false)
	vector.FillRect(screen, 0, 458, screenWidth, 82, color.RGBA{R: 26, G: 31, B: 27, A: 255}, false)

	for i := 0; i < 6; i++ {
		x := float32(110 + i*145)
		vector.StrokeRect(screen, x, 102, 46, 46, 2, color.RGBA{R: 65, G: 60, B: 94, A: 130}, true)
		vector.StrokeRect(screen, x+9, 111, 28, 28, 2, color.RGBA{R: 111, G: 96, B: 58, A: 115}, true)
	}
}

// activeButtons returns the buttons available on the current screen.
func (g *game) activeButtons() []menu.Button {
	switch g.screen {
	case screenNewGame, screenSettings:
		return []menu.Button{g.backButton}
	default:
		return g.mainButtons
	}
}

// drawMenuPanel renders the title area and menu copy.
func (g *game) drawMenuPanel(screen *ebiten.Image) {
	panelX := float32(220)
	panelY := float32(82)
	panelW := float32(520)
	panelH := float32(398)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	g.drawCenteredText(screen, "td", g.titleFace, 122, textColor)
	g.drawCenteredText(screen, "Arcane defenses await their first command.", g.bodyFace, 214, mutedTextColor)
}

// drawButtons renders menu buttons with hover feedback.
func (g *game) drawButtons(screen *ebiten.Image, buttons []menu.Button) {
	for _, button := range buttons {
		fill := buttonColor
		edge := panelEdgeColor
		labelColor := textColor
		if button.Disabled {
			fill = disabledColor
			edge = color.RGBA{R: 83, G: 84, B: 73, A: 255}
			labelColor = mutedTextColor
		} else if g.hoverAction != menu.ActionNone && g.hoverAction == button.Action {
			fill = hoverColor
			edge = textColor
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)
		g.drawCenteredText(screen, button.Label, g.buttonFace, float64(button.Y+9), labelColor)
	}
}

// drawPlaceholderPanel renders a temporary screen reached from the main menu.
func (g *game) drawPlaceholderPanel(screen *ebiten.Image, title, message string) {
	panelX := float32(180)
	panelY := float32(128)
	panelW := float32(600)
	panelH := float32(340)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	g.drawCenteredText(screen, title, g.titleFace, 174, textColor)
	g.drawCenteredText(screen, message, g.bodyFace, 286, mutedTextColor)
}

// drawSettingsPanel renders the temporary settings screen.
func (g *game) drawSettingsPanel(screen *ebiten.Image) {
	panelX := float32(180)
	panelY := float32(128)
	panelW := float32(600)
	panelH := float32(340)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	g.drawCenteredText(screen, "Settings", g.titleFace, 214, textColor)
}

// drawCenteredText draws one line centered horizontally at the given y coordinate.
func (g *game) drawCenteredText(screen *ebiten.Image, value string, face *text.GoTextFace, y float64, clr color.Color) {
	width, _ := text.Measure(value, face, face.Size)
	options := &text.DrawOptions{}
	options.GeoM.Translate((screenWidth-width)/2, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}
