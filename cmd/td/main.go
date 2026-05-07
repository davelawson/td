package main

import (
	"errors"
	"log"

	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	defaultWindowWidth  = 1920
	defaultWindowHeight = 1080
)

type app struct {
	mainMenu *menu.Menu
}

// main starts the Ebitengine desktop application.
func main() {
	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(app); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}

// newApp creates the app state used by Ebitengine callbacks.
func newApp() (*app, error) {
	mainMenu, err := menu.New(defaultWindowWidth, defaultWindowHeight)
	if err != nil {
		return nil, err
	}
	return &app{mainMenu: mainMenu}, nil
}

// Update handles pointer input and returns a clean termination signal on quit.
func (a *app) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	input := menu.Input{
		CursorX:   cursorX,
		CursorY:   cursorY,
		Clicked:   inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		Typed:     ebiten.AppendInputChars(nil),
		Backspace: inpututil.IsKeyJustPressed(ebiten.KeyBackspace),
	}
	if action := a.mainMenu.Update(input); action == menu.ActionQuit {
		return ebiten.Termination
	}
	return nil
}

// Draw renders the current game screen.
func (a *app) Draw(screen *ebiten.Image) {
	a.mainMenu.Draw(screen)
}

// Layout returns a pixel-sized drawable layout for the current window.
func (a *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		outsideWidth = defaultWindowWidth
		outsideHeight = defaultWindowHeight
	}
	a.mainMenu.Resize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
