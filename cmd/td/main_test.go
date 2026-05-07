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

type screenshotGame struct {
	*game
	targets []screenshotTarget
	index   int
}

type screenshotTarget struct {
	mode screenMode
	path string
}

// TestHandleActionRoutesMenuScreens verifies menu actions change screens or quit.
func TestHandleActionRoutesMenuScreens(t *testing.T) {
	game, err := newGame()
	if err != nil {
		t.Fatal(err)
	}

	if err := game.handleAction(menu.ActionNew); err != nil {
		t.Fatal(err)
	}
	if game.screen != screenNewGame {
		t.Fatalf("screen = %v, want %v", game.screen, screenNewGame)
	}

	if err := game.handleAction(menu.ActionBack); err != nil {
		t.Fatal(err)
	}
	if game.screen != screenMainMenu {
		t.Fatalf("screen = %v, want %v", game.screen, screenMainMenu)
	}

	if err := game.handleAction(menu.ActionSettings); err != nil {
		t.Fatal(err)
	}
	if game.screen != screenSettings {
		t.Fatalf("screen = %v, want %v", game.screen, screenSettings)
	}

	if err := game.handleAction(menu.ActionNone); err != nil {
		t.Fatal(err)
	}
	if game.screen != screenSettings {
		t.Fatalf("screen = %v, want %v", game.screen, screenSettings)
	}

	if err := game.handleAction(menu.ActionQuit); !errors.Is(err, ebiten.Termination) {
		t.Fatalf("handleAction(ActionQuit) = %v, want ebiten.Termination", err)
	}
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

	basePath := filepath.Join("..", "..", "plans", "01-expanded-main-menu", "screenshots")
	capture := &screenshotGame{
		game: game,
		targets: []screenshotTarget{
			{mode: screenMainMenu, path: filepath.Join(basePath, "main-menu.png")},
			{mode: screenNewGame, path: filepath.Join(basePath, "new-game-placeholder.png")},
			{mode: screenSettings, path: filepath.Join(basePath, "settings-placeholder.png")},
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
func (g *screenshotGame) Update() error {
	if g.index >= len(g.targets) {
		return ebiten.Termination
	}
	g.game.screen = g.targets[g.index].mode
	return nil
}

// Draw renders the current target screen and writes the frame to disk.
func (g *screenshotGame) Draw(screen *ebiten.Image) {
	if g.index >= len(g.targets) {
		return
	}

	g.game.Draw(screen)
	target := g.targets[g.index]

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
	g.index++
}
