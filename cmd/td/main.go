package main

import (
	"errors"
	"log"

	"td/internal/game"
	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	defaultWindowWidth  = 1920
	defaultWindowHeight = 1080
)

type app struct {
	mode      appMode
	width     int
	height    int
	mainMenu  *menu.Menu
	gameState *game.State
}

type appMode int

const (
	appModeMenu appMode = iota
	appModeGame
)

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
	return &app{
		mode:     appModeMenu,
		width:    defaultWindowWidth,
		height:   defaultWindowHeight,
		mainMenu: mainMenu,
	}, nil
}

// Update routes Ebitengine input to the active app mode.
func (a *app) Update() error {
	switch a.mode {
	case appModeGame:
		return a.updateGame()
	default:
		return a.updateMenu()
	}
}

// updateMenu handles menu input and mode transitions.
func (a *app) updateMenu() error {
	cursorX, cursorY := ebiten.CursorPosition()
	input := menu.Input{
		CursorX:   cursorX,
		CursorY:   cursorY,
		Clicked:   inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		Typed:     ebiten.AppendInputChars(nil),
		Backspace: inpututil.IsKeyJustPressed(ebiten.KeyBackspace),
	}
	switch action := a.mainMenu.Update(input); action {
	case menu.ActionQuit:
		return ebiten.Termination
	case menu.ActionStart:
		return a.startGame(a.mainMenu.WizardName())
	}
	return nil
}

// updateGame handles in-game input and logical updates.
func (a *app) updateGame() error {
	cursorX, cursorY := ebiten.CursorPosition()
	_, wheelY := ebiten.Wheel()
	switch action := a.gameState.Update(game.Input{
		TogglePause: inpututil.IsKeyJustPressed(ebiten.KeySpace),
		ToggleMenu:  inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		WheelY:      wheelY,
		PanUp:       ebiten.IsKeyPressed(ebiten.KeyW),
		PanDown:     ebiten.IsKeyPressed(ebiten.KeyS),
		PanLeft:     ebiten.IsKeyPressed(ebiten.KeyA),
		PanRight:    ebiten.IsKeyPressed(ebiten.KeyD),
		CursorX:     cursorX,
		CursorY:     cursorY,
		Clicked:     inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
	}); action {
	case game.ActionSurrender:
		a.returnToMainMenu()
	}
	return nil
}

// Draw renders the current game screen.
func (a *app) Draw(screen *ebiten.Image) {
	switch a.mode {
	case appModeGame:
		a.gameState.Draw(screen)
	default:
		a.mainMenu.Draw(screen)
	}
}

// Layout returns a pixel-sized drawable layout for the current window.
func (a *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		outsideWidth = defaultWindowWidth
		outsideHeight = defaultWindowHeight
	}
	a.width = outsideWidth
	a.height = outsideHeight
	if a.mainMenu != nil {
		a.mainMenu.Resize(outsideWidth, outsideHeight)
	}
	if a.gameState != nil {
		a.gameState.Resize(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

// startGame creates the first game state and closes the menu.
func (a *app) startGame(wizardName string) error {
	gameState, err := game.New(wizardName, a.width, a.height)
	if err != nil {
		return err
	}
	a.mode = appModeGame
	a.gameState = gameState
	return nil
}

// returnToMainMenu leaves the active game and shows the top-level menu.
func (a *app) returnToMainMenu() {
	a.mode = appModeMenu
	a.gameState = nil
	a.mainMenu.ResetToMain()
}
