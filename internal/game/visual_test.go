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

// TestCaptureTerrainProductionBeforeScreenshot records a producer and its selected source.
func TestCaptureTerrainProductionBeforeScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state := newTerrainProductionScreenshotState(t)
	path := filepath.Join(
		"..", "..", "plans", "57-terrain-consuming-resource-production", "screenshots", "terrain-production-before.png",
	)
	captureStateScreenshot(t, state, path)
}

// TestCaptureTerrainProductionAfterScreenshot records the terrain and resource result after Labour.
func TestCaptureTerrainProductionAfterScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state := newTerrainProductionScreenshotState(t)
	state.beginPostRaidDay()
	path := filepath.Join(
		"..", "..", "plans", "57-terrain-consuming-resource-production", "screenshots", "terrain-production-after.png",
	)
	captureStateScreenshot(t, state, path)
}

func newTerrainProductionScreenshotState(t *testing.T) *State {
	t.Helper()
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	clearNaturalTerrain(&state.gameMap.Home)
	producer := homeTileCoordinate(homePlotCenter+2, homePlotCenter)
	source := homeTileCoordinate(homePlotCenter+4, homePlotCenter)
	state.gameMap.Home.Tiles[producer.Y][producer.X].Feature = featureWoodcutter
	state.gameMap.Home.Tiles[source.Y][source.X] = Tile{Terrain: terrainTree, Tweak: 1}
	state.selection = selectedItem{kind: selectedItemTerrain, tile: source}
	return state
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

// TestCaptureRaidTempoBeforeScreenshot writes the legacy four-second challenge-16 frame when enabled.
func TestCaptureRaidTempoBeforeScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state := newRaidTempoScreenshotState(t)
	state.raid.template.progressDurationSeconds = raidProgressDurationBase + state.raid.template.challengeRating
	for i := 0; i < 240; i++ {
		state.Update(Input{})
	}
	if got, want := len(state.raid.enemies), 1; got != want {
		t.Fatalf("legacy active enemies after four seconds = %d, want %d", got, want)
	}

	path := filepath.Join(
		"..", "..", "plans", "58-accelerating-raid-tempo", "screenshots", "before", "raid-tempo.png",
	)
	captureStateScreenshot(t, state, path)
}

// TestCaptureRaidTempoAfterScreenshot writes the accelerated four-second challenge-16 frame when enabled.
func TestCaptureRaidTempoAfterScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state := newRaidTempoScreenshotState(t)
	for i := 0; i < 240; i++ {
		state.Update(Input{})
	}
	if got, want := len(state.raid.enemies), 3; got != want {
		t.Fatalf("active enemies after four seconds = %d, want %d", got, want)
	}

	path := filepath.Join(
		"..", "..", "plans", "58-accelerating-raid-tempo", "screenshots", "after", "raid-tempo.png",
	)
	captureStateScreenshot(t, state, path)
}

func newRaidTempoScreenshotState(t *testing.T) *State {
	t.Helper()
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	clearNaturalTerrain(&state.gameMap.Home)
	state.gameMap.frontierBiomes = map[plotCoordinate]plotBiome{
		{X: 0, Y: 1}:  biomeGrasslands,
		{X: 1, Y: 0}:  biomeGrasslands,
		{X: 0, Y: -1}: biomeGrasslands,
		{X: -1, Y: 0}: biomeGrasslands,
	}
	state.status.populations.peasants = populationCount{available: 120, total: 120}
	state.startNextRaid()
	return state
}

// TestCaptureGhoulAndArmouredSkeletonScreenshot writes focused new-enemy evidence when enabled.
func TestCaptureGhoulAndArmouredSkeletonScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	clearNaturalTerrain(&state.gameMap.Home)
	for y := 0; y < plotSize; y++ {
		state.gameMap.Home.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
	state.status.phase = phaseRaid
	state.camera.centerY = 2.5
	state.raid = raidState{
		active:   true,
		number:   1,
		template: raidTemplate{challengeRating: 8, progressDurationSeconds: raidProgressDuration(8)},
		enemies: []raidEnemy{
			{id: 41, template: &state.enemyCatalog.Ghoul, position: coord{X: 0, Y: 5.5}, health: 20},
			{id: 42, template: &state.enemyCatalog.ArmouredSkeleton, position: coord{X: 0, Y: 1.5}, health: 80},
		},
	}
	state.selection = selectedItem{kind: selectedItemRaider, raiderID: 42}

	panel, ok := state.currentSelectionPanel()
	if !ok || panel.Name != "Armoured Skeleton" || panel.Health != 80 || panel.MaxHealth != 125 || panel.SpeedTilesPerSecond != 0.9 {
		t.Fatalf("selected Armoured Skeleton panel = %+v, available %v", panel, ok)
	}
	path := filepath.Join(
		"..", "..", "plans", "59-ghouls-and-armoured-skeletons", "screenshots", "ghoul-and-armoured-skeleton.png",
	)
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
