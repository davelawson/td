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
	accentColor     = color.RGBA{R: 98, G: 90, B: 145, A: 255}
)

type game struct {
	buttons     []menu.Button
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

	quitButton := menu.Button{
		Label:  "Quit",
		X:      screenWidth/2 - 110,
		Y:      348,
		W:      220,
		H:      62,
		Action: menu.ActionQuit,
	}

	return &game{
		buttons: []menu.Button{quitButton},
		titleFace: &text.GoTextFace{
			Source: source,
			Size:   88,
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
	g.hoverAction = menu.ActionAt(g.buttons, cursorX, cursorY)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if menu.ActionAt(g.buttons, cursorX, cursorY) == menu.ActionQuit {
			return ebiten.Termination
		}
	}
	return nil
}

// Draw renders the first main-menu screen.
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.drawBackdrop(screen)
	g.drawMenuPanel(screen)
	g.drawButtons(screen)
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

// drawMenuPanel renders the title area and menu copy.
func (g *game) drawMenuPanel(screen *ebiten.Image) {
	panelX := float32(220)
	panelY := float32(132)
	panelW := float32(520)
	panelH := float32(300)

	vector.FillRect(screen, panelX, panelY, panelW, panelH, panelColor, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 4, panelEdgeColor, false)
	vector.StrokeRect(screen, panelX+12, panelY+12, panelW-24, panelH-24, 1.5, accentColor, false)

	g.drawCenteredText(screen, "td", g.titleFace, 185, textColor)
	g.drawCenteredText(screen, "Arcane defenses await their first command.", g.bodyFace, 286, mutedTextColor)
}

// drawButtons renders menu buttons with hover feedback.
func (g *game) drawButtons(screen *ebiten.Image) {
	for _, button := range g.buttons {
		fill := buttonColor
		edge := panelEdgeColor
		if g.hoverAction == button.Action {
			fill = hoverColor
			edge = textColor
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)
		g.drawCenteredText(screen, button.Label, g.buttonFace, float64(button.Y+15), textColor)
	}
}

// drawCenteredText draws one line centered horizontally at the given y coordinate.
func (g *game) drawCenteredText(screen *ebiten.Image, value string, face *text.GoTextFace, y float64, clr color.Color) {
	width, _ := text.Measure(value, face, face.Size)
	options := &text.DrawOptions{}
	options.GeoM.Translate((screenWidth-width)/2, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}
