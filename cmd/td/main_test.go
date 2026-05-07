package main

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenshotGame struct {
	*game
	path     string
	captured bool
}

// TestCaptureMainMenuScreenshot writes visual evidence when explicitly enabled.
func TestCaptureMainMenuScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	game, err := newGame()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join("..", "..", "plans", "00-initial-ebitengine-menu", "screenshots", "main-menu.png")
	capture := &screenshotGame{game: game, path: path}

	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(capture); err != nil && !errors.Is(err, ebiten.Termination) {
		t.Fatal(err)
	}
	if !capture.captured {
		t.Fatal("expected screenshot capture before termination")
	}
}

// Update terminates after Draw captures the first rendered frame.
func (g *screenshotGame) Update() error {
	if g.captured {
		return ebiten.Termination
	}
	return nil
}

// Draw renders the menu and writes the first frame to disk.
func (g *screenshotGame) Draw(screen *ebiten.Image) {
	g.game.Draw(screen)
	if g.captured {
		return
	}

	if err := os.MkdirAll(filepath.Dir(g.path), 0o755); err != nil {
		panic(err)
	}
	file, err := os.Create(g.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	frame := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	screen.ReadPixels(frame.Pix)
	if err := png.Encode(file, frame); err != nil {
		panic(err)
	}
	g.captured = true
}
