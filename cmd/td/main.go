package main

import (
	"errors"
	"log"

	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 960
	screenHeight = 540
)

type app struct {
	mainMenu *menu.Menu
}

// main starts the Ebitengine desktop application.
func main() {
	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(screenWidth, screenHeight)
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
	mainMenu, err := menu.New(screenWidth, screenHeight)
	if err != nil {
		return nil, err
	}
	return &app{mainMenu: mainMenu}, nil
}

// Update handles pointer input and returns a clean termination signal on quit.
func (a *app) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if action := a.mainMenu.Update(cursorX, cursorY, clicked); action == menu.ActionQuit {
		return ebiten.Termination
	}
	return nil
}

// Draw renders the current game screen.
func (a *app) Draw(screen *ebiten.Image) {
	a.mainMenu.Draw(screen)
}

// Layout returns the fixed logical resolution for the prototype.
func (a *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
