package main

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenshotApp struct {
	*app
	targets []screenshotTarget
	index   int
}

type screenshotTarget struct {
	screen menu.Screen
	path   string
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

	basePath := filepath.Join("..", "..", "plans", "03-new-game-configuration", "screenshots")
	capture := &screenshotApp{
		app: app,
		targets: []screenshotTarget{
			{screen: menu.ScreenMain, path: filepath.Join(basePath, "main-menu.png")},
			{screen: menu.ScreenNewGame, path: filepath.Join(basePath, "new-game-configuration.png")},
			{screen: menu.ScreenSettings, path: filepath.Join(basePath, "settings-placeholder.png")},
		},
	}

	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(screenWidth, screenHeight)
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
	a.app.mainMenu.SetScreenForTest(a.targets[a.index].screen)
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

	frame := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	screen.ReadPixels(frame.Pix)
	if err := png.Encode(file, frame); err != nil {
		panic(err)
	}
	a.index++
}
