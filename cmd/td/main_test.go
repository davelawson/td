package main

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"td/internal/game"
	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenshotApp struct {
	*app
	targets []screenshotTarget
	index   int
}

type screenshotTarget struct {
	screen     menu.Screen
	wizardName string
	paused     bool
	path       string
}

// TestStartGameSwitchesToGameMode verifies app-level game startup.
func TestStartGameSwitchesToGameMode(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	if err := app.startGame("Merlin"); err != nil {
		t.Fatal(err)
	}
	if app.mode != appModeGame {
		t.Fatalf("mode = %v, want %v", app.mode, appModeGame)
	}
	if app.gameState == nil {
		t.Fatal("expected game state after starting")
	}
	if app.gameState.WizardName() != "Merlin" {
		t.Fatalf("wizard name = %q, want %q", app.gameState.WizardName(), "Merlin")
	}
}

// TestCaptureMainMenuScreenshot writes visual evidence when explicitly enabled.
func TestCaptureMainMenuScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	basePath := filepath.Join("..", "..", "plans", "05-main-game-update-loop", "screenshots")
	capture := &screenshotApp{
		app: app,
		targets: []screenshotTarget{
			{screen: menu.ScreenMain, path: filepath.Join(basePath, "main-menu.png")},
			{screen: menu.ScreenNewGame, path: filepath.Join(basePath, "new-game-configuration.png")},
			{wizardName: "Merlin", path: filepath.Join(basePath, "running-game.png")},
			{wizardName: "Merlin", paused: true, path: filepath.Join(basePath, "paused-game.png")},
		},
	}

	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	if err := ebiten.RunGame(capture); err != nil && !errors.Is(err, ebiten.Termination) {
		t.Fatal(err)
	}
	if capture.index != len(capture.targets) {
		t.Fatalf("captured %d screenshots, want %d", capture.index, len(capture.targets))
	}
}

// Update sets the next screen to capture or terminates after all captures.
func (a *screenshotApp) Update() error {
	if a.index >= len(a.targets) {
		return ebiten.Termination
	}
	target := a.targets[a.index]
	if target.wizardName != "" {
		if err := a.app.startGame(target.wizardName); err != nil {
			return err
		}
		a.app.gameState.Update(game.Input{})
		if target.paused {
			a.app.gameState.Update(game.Input{TogglePause: true})
		}
		return nil
	}

	a.app.mode = appModeMenu
	a.app.gameState = nil
	a.app.mainMenu.SetScreenForTest(target.screen)
	return nil
}

// Draw renders the current target screen and writes the frame to disk.
func (a *screenshotApp) Draw(screen *ebiten.Image) {
	if a.index >= len(a.targets) {
		return
	}

	a.app.Draw(screen)
	target := a.targets[a.index]

	if err := os.MkdirAll(filepath.Dir(target.path), 0o755); err != nil {
		panic(err)
	}
	file, err := os.Create(target.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	frame := image.NewRGBA(image.Rect(0, 0, defaultWindowWidth, defaultWindowHeight))
	screen.ReadPixels(frame.Pix)
	if err := png.Encode(file, frame); err != nil {
		panic(err)
	}
	a.index++
}
