package game

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type visualScreenshotGame struct {
	state    *State
	path     string
	captured bool
	err      error
}

// TestCaptureSelectedTerrainScreenshot writes focused terrain-selection evidence when enabled.
func TestCaptureSelectedTerrainScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainTree}
	state.selection = selectedItem{kind: selectedItemTerrain, tile: tile}

	path := filepath.Join("..", "..", "plans", "53-phase-aware-game-ui", "screenshots", "selected-terrain.png")
	captureStateScreenshot(t, state, path)
}

// TestCaptureSelectedIronDepositScreenshot writes focused Iron Deposit evidence when enabled.
func TestCaptureSelectedIronDepositScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainIronDeposit, Tweak: 1}
	state.selection = selectedItem{kind: selectedItemTerrain, tile: tile}

	path := filepath.Join("..", "..", "plans", "56-iron-deposit-terrain", "screenshots", "selected-iron-deposit.png")
	captureStateScreenshot(t, state, path)
}

// TestCaptureActiveRaidScreenshot writes focused phase-aware Raid UI evidence when enabled.
func TestCaptureActiveRaidScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	state.startNextRaid()
	for i := 0; i < 330; i++ {
		state.Update(Input{})
	}

	path := filepath.Join("..", "..", "plans", "53-phase-aware-game-ui", "screenshots", "active-raid.png")
	captureStateScreenshot(t, state, path)
}

// captureStateScreenshot runs one state inside Ebitengine and verifies its frame was saved.
func captureStateScreenshot(t *testing.T, state *State, path string) {
	t.Helper()
	capture := &visualScreenshotGame{state: state, path: path}
	ebiten.SetWindowSize(1920, 1080)
	if err := ebiten.RunGame(capture); err != nil && !errors.Is(err, ebiten.Termination) {
		t.Fatal(err)
	}
	if capture.err != nil {
		t.Fatal(capture.err)
	}
	if !capture.captured {
		t.Fatal("expected selected-terrain screenshot capture")
	}
}

// Update terminates after the selected-terrain frame has been captured.
func (g *visualScreenshotGame) Update() error {
	if g.captured {
		return ebiten.Termination
	}
	return nil
}

// Draw renders and saves one selected-terrain frame from inside Ebitengine's game loop.
func (g *visualScreenshotGame) Draw(screen *ebiten.Image) {
	if g.captured || g.err != nil {
		return
	}
	rendered := ebiten.NewImage(1920, 1080)
	g.state.Draw(rendered)
	frame := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
	rendered.ReadPixels(frame.Pix)
	screen.DrawImage(rendered, nil)

	if err := os.MkdirAll(filepath.Dir(g.path), 0o755); err != nil {
		g.err = err
		return
	}
	file, err := os.Create(g.path)
	if err != nil {
		g.err = err
		return
	}
	defer file.Close()
	if err := png.Encode(file, frame); err != nil {
		g.err = err
		return
	}
	g.captured = true
}

// Layout preserves the fixed screenshot dimensions.
func (g *visualScreenshotGame) Layout(_, _ int) (int, int) {
	return 1920, 1080
}
