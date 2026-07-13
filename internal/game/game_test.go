package game

import (
	"math"
	"testing"

	"td/internal/ui"
)

// TestNewStateStartsRunning verifies initial game state.
func TestNewStateStartsRunning(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if state.WizardName() != "Merlin" {
		t.Fatalf("wizard name = %q, want %q", state.WizardName(), "Merlin")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
	if state.Paused() {
		t.Fatal("expected new state to start unpaused")
	}
	if state.status.phase != phaseCalm {
		t.Fatalf("phase = %v, want %v", state.status.phase, phaseCalm)
	}
	if state.status.resources.wood != 100 || state.status.resources.stone != 50 || state.status.resources.metal != 20 {
		t.Fatalf("resources = %+v, want wood 100 stone 50 metal 20", state.status.resources)
	}
	if state.gameMap.Home.Tiles[homePlotCenter][homePlotCenter].Feature != featureSanctum {
		t.Fatal("expected new state to store the default home Plot")
	}
	if state.assetCatalog.Sprite.Structure.Sanctum == nil {
		t.Fatal("expected new state to store the Sanctum sprite")
	}
	if state.enemyCatalog.SkeletonSwordShield.Name == "" {
		t.Fatal("expected new state to store the skeleton sword-and-shield enemy template")
	}
	if state.enemyCatalog.SkeletonSwordShield.SpriteKey != "skeleton-sword-shield" {
		t.Fatalf("skeleton sprite key = %q, want %q", state.enemyCatalog.SkeletonSwordShield.SpriteKey, "skeleton-sword-shield")
	}
	if state.enemyCatalog.SkeletonSwordShield.Sprite == nil {
		t.Fatal("expected new state to store the skeleton sword-and-shield enemy sprite")
	}
	if state.structureCatalog.BowTower.Name != "Bow Tower" {
		t.Fatalf("Bow Tower name = %q, want %q", state.structureCatalog.BowTower.Name, "Bow Tower")
	}
	if state.structureCatalog.BowTower.Sprite == nil {
		t.Fatal("expected new state to store the Bow Tower structure template sprite")
	}
	if state.structureCatalog.House.Name != "House" {
		t.Fatalf("House name = %q, want %q", state.structureCatalog.House.Name, "House")
	}
	if state.structureCatalog.House.Sprite == nil {
		t.Fatal("expected new state to store the House structure template sprite")
	}
	if state.structureCatalog.Dorm.Name != "Dorm" {
		t.Fatalf("Dorm name = %q, want %q", state.structureCatalog.Dorm.Name, "Dorm")
	}
	if state.structureCatalog.Dorm.Sprite == nil {
		t.Fatal("expected new state to store the Dorm structure template sprite")
	}
	if state.structureCatalog.FlameBoltTower.Name != "Flame Bolt Tower" {
		t.Fatalf("Flame Bolt Tower name = %q, want %q", state.structureCatalog.FlameBoltTower.Name, "Flame Bolt Tower")
	}
	if state.structureCatalog.FlameBoltTower.Sprite == nil {
		t.Fatal("expected new state to store the Flame Bolt Tower structure template sprite")
	}
	if state.camera.zoom != cameraInitialZoom {
		t.Fatalf("camera zoom = %f, want %f", state.camera.zoom, cameraInitialZoom)
	}
	if state.camera.centerX != 0 || state.camera.centerY != 0 {
		t.Fatalf("camera center = (%f,%f), want (0,0)", state.camera.centerX, state.camera.centerY)
	}
}

// TestDefaultHomePlotShape verifies the prototype Plot dimensions and center.
func TestDefaultHomePlotShape(t *testing.T) {
	plot := NewDefaultHomePlot()

	if len(plot.Tiles) != plotSize {
		t.Fatalf("plot rows = %d, want %d", len(plot.Tiles), plotSize)
	}
	if len(plot.Tiles[0]) != plotSize {
		t.Fatalf("plot columns = %d, want %d", len(plot.Tiles[0]), plotSize)
	}
	if plot.Tiles[homePlotCenter][homePlotCenter].Feature != featureSanctum {
		t.Fatal("expected Sanctum at the center Tile")
	}
}

// TestDefaultHomePlotStartsWithoutTowers verifies the player receives no free defenses.
func TestDefaultHomePlotStartsWithoutTowers(t *testing.T) {
	plot := NewDefaultHomePlot()
	for y := range plot.Tiles {
		for x := range plot.Tiles[y] {
			feature := plot.Tiles[y][x].Feature
			if x == homePlotCenter && y == homePlotCenter {
				if feature != featureSanctum {
					t.Fatalf("center feature = %v, want Sanctum", feature)
				}
				continue
			}
			if feature != featureNone {
				t.Fatalf("tile (%d,%d) feature = %v, want none", x, y, feature)
			}
		}
	}
}

// TestDefaultHomePlotAssignsTileTweaks verifies Tile creation stores tweak values on map data.
func TestDefaultHomePlotAssignsTileTweaks(t *testing.T) {
	var next uint16
	plot := newDefaultHomePlotWithTweakSource(func() uint16 {
		value := next
		next++
		return value
	})

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			want := uint16(y*plotSize + x)
			if plot.Tiles[y][x].Tweak != want {
				t.Fatalf("tile (%d,%d) tweak = %d, want %d", x, y, plot.Tiles[y][x].Tweak, want)
			}
		}
	}
}

// TestDefaultHomePlotRoadRunsNorth verifies the authored road layout.
func TestDefaultHomePlotRoadRunsNorth(t *testing.T) {
	plot := NewDefaultHomePlot()

	for y := 0; y <= homePlotCenter; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainRoad {
			t.Fatalf("tile (%d,%d) terrain = %v, want road", homePlotCenter, y, plot.Tiles[y][homePlotCenter].Terrain)
		}
	}
	for y := homePlotCenter + 1; y < plotSize; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain == terrainRoad {
			t.Fatalf("tile (%d,%d) is road below the Sanctum", homePlotCenter, y)
		}
	}
}

// TestDefaultHomePlotIsOtherwiseGrassland verifies every non-road Tile is empty grass.
func TestDefaultHomePlotIsOtherwiseGrassland(t *testing.T) {
	plot := NewDefaultHomePlot()

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			tile := plot.Tiles[y][x]
			onNorthRoad := x == homePlotCenter && y <= homePlotCenter

			if onNorthRoad && tile.Terrain != terrainRoad {
				t.Fatalf("tile (%d,%d) terrain = %v, want road", x, y, tile.Terrain)
			}
			if !onNorthRoad && tile.Terrain != terrainEmpty {
				t.Fatalf("tile (%d,%d) terrain = %v, want empty grass", x, y, tile.Terrain)
			}
		}
	}
}

// TestDefaultHomePlotInteriorIsOtherwiseEmpty verifies non-road interior Tiles stay empty.
func TestDefaultHomePlotInteriorIsOtherwiseEmpty(t *testing.T) {
	plot := NewDefaultHomePlot()

	for y := 1; y < plotSize-1; y++ {
		for x := 1; x < plotSize-1; x++ {
			tile := plot.Tiles[y][x]
			onNorthRoad := x == homePlotCenter && y <= homePlotCenter
			atSanctum := x == homePlotCenter && y == homePlotCenter
			atStartingTower := y == 5 && (x == homePlotCenter+1 || x == homePlotCenter-1)

			if !onNorthRoad && tile.Terrain != terrainEmpty {
				t.Fatalf("tile (%d,%d) terrain = %v, want empty interior", x, y, tile.Terrain)
			}
			if !atSanctum && !atStartingTower && tile.Feature != featureNone {
				t.Fatalf("tile (%d,%d) feature = %v, want none", x, y, tile.Feature)
			}
		}
	}
}

// TestTerrainSpriteIndexUsesLowTweakBits verifies terrain variant selection ignores the flip bit.
func TestTerrainSpriteIndexUsesLowTweakBits(t *testing.T) {
	tests := []struct {
		name  string
		tweak uint16
		want  int
	}{
		{name: "first", tweak: 0, want: 0},
		{name: "wrap", tweak: 5, want: 1},
		{name: "high bit ignored", tweak: 0x8005, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := terrainSpriteIndex(tt.tweak, 4); got != tt.want {
				t.Fatalf("terrainSpriteIndex(%d, 4) = %d, want %d", tt.tweak, got, tt.want)
			}
		})
	}
}

// TestTerrainSpriteFlippedUsesHighTweakBit verifies the high tweak bit controls mirroring.
func TestTerrainSpriteFlippedUsesHighTweakBit(t *testing.T) {
	if terrainSpriteFlipped(0x7fff) {
		t.Fatal("expected tweak below high bit to avoid horizontal flip")
	}
	if !terrainSpriteFlipped(0x8000) {
		t.Fatal("expected high bit to request horizontal flip")
	}
}

// TestUpdatesDoNotChangeHomePlot verifies time passing does not mutate the scene.
func TestUpdatesDoNotChangeHomePlot(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	initial := state.gameMap

	state.Update(Input{})
	state.Update(Input{})
	state.Update(Input{})

	if !mapsEqual(state.gameMap, initial) {
		t.Fatal("expected logical updates to leave the prototype map unchanged")
	}
}

// TestCameraWheelUpIncreasesZoom verifies scroll-up zooms in.
func TestCameraWheelUpIncreasesZoom(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{WheelY: 1})

	if state.camera.zoom <= cameraInitialZoom {
		t.Fatalf("camera zoom = %f, want greater than %f", state.camera.zoom, cameraInitialZoom)
	}
}

// TestCameraWheelDownDecreasesZoomToFloor verifies scroll-down zooms out safely.
func TestCameraWheelDownDecreasesZoomToFloor(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{WheelY: -1})
	if state.camera.zoom >= cameraInitialZoom {
		t.Fatalf("camera zoom = %f, want less than %f", state.camera.zoom, cameraInitialZoom)
	}

	state.Update(Input{WheelY: -1000})
	if state.camera.zoom != cameraMinZoom {
		t.Fatalf("camera zoom = %f, want floor %f", state.camera.zoom, cameraMinZoom)
	}
}

// TestCameraPanInputMovesCenter verifies WASD pan in the expected directions.
func TestCameraPanInputMovesCenter(t *testing.T) {
	tests := []struct {
		name  string
		input Input
		wantX float64
		wantY float64
	}{
		{name: "up", input: inputWithPan(true, false, false, false), wantX: 0, wantY: cameraPanSpeed / plotBaseTileSize},
		{name: "down", input: inputWithPan(false, true, false, false), wantX: 0, wantY: -cameraPanSpeed / plotBaseTileSize},
		{name: "left", input: inputWithPan(false, false, true, false), wantX: -cameraPanSpeed / plotBaseTileSize, wantY: 0},
		{name: "right", input: inputWithPan(false, false, false, true), wantX: cameraPanSpeed / plotBaseTileSize, wantY: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := New("Merlin", 1920, 1080)
			if err != nil {
				t.Fatal(err)
			}
			startX := state.camera.centerX
			startY := state.camera.centerY

			state.Update(tt.input)

			if !almostEqual(state.camera.centerX-startX, tt.wantX) || !almostEqual(state.camera.centerY-startY, tt.wantY) {
				t.Fatalf("camera delta = (%f,%f), want (%f,%f)", state.camera.centerX-startX, state.camera.centerY-startY, tt.wantX, tt.wantY)
			}
		})
	}
}

// TestCameraPanSpeedDividesByZoom verifies panning slows in world units when zoomed in.
func TestCameraPanSpeedDividesByZoom(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	state.camera.zoom = 2
	startX := state.camera.centerX

	state.Update(inputWithPan(false, false, false, true))

	if got, want := state.camera.centerX-startX, cameraPanSpeed/(plotBaseTileSize*2); !almostEqual(got, want) {
		t.Fatalf("camera x delta = %f, want %f", got, want)
	}
}

// TestCameraInputWorksWhilePaused verifies pause stops logic but not map inspection.
func TestCameraInputWorksWhilePaused(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	startX := state.camera.centerX
	startZoom := state.camera.zoom
	state.Update(inputWithWheelAndPan(1, false, false, false, true))

	if state.camera.centerX <= startX {
		t.Fatalf("camera center x = %f, want greater than %f", state.camera.centerX, startX)
	}
	if state.camera.zoom <= startZoom {
		t.Fatalf("camera zoom = %f, want greater than %f", state.camera.zoom, startZoom)
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
}

// TestIngameMenuBlocksCameraInput verifies overlay-open frames ignore camera controls.
func TestIngameMenuBlocksCameraInput(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{ToggleMenu: true})
	startCamera := state.camera
	state.Update(inputWithWheelAndPan(1, false, true, false, true))

	if state.camera != startCamera {
		t.Fatalf("camera = %+v, want unchanged %+v", state.camera, startCamera)
	}
}

// TestCameraChangesDoNotChangeHomePlot verifies inspection controls leave map data untouched.
func TestCameraChangesDoNotChangeHomePlot(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	initial := state.gameMap

	state.Update(inputWithWheelAndPan(2, true, false, true, false))
	state.Update(inputWithWheelAndPan(-4, false, true, false, true))

	if !mapsEqual(state.gameMap, initial) {
		t.Fatal("expected camera changes to leave the prototype map unchanged")
	}
}

// inputWithPan creates Input with directional camera pan flags.
func inputWithPan(up, down, left, right bool) Input {
	var input Input
	input.Pan.Up = up
	input.Pan.Down = down
	input.Pan.Left = left
	input.Pan.Right = right
	return input
}

// inputWithWheelAndPan creates Input with zoom and directional camera pan flags.
func inputWithWheelAndPan(wheelY float64, up, down, left, right bool) Input {
	input := inputWithPan(up, down, left, right)
	input.WheelY = wheelY
	return input
}

// TestStateFormatsCalmTopBar verifies initial top-bar text.
func TestStateFormatsCalmTopBar(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if value := state.chapterDayText(); value != "Chapter I: The Ashen Copse | Day 1" {
		t.Fatalf("chapterDayText = %q", value)
	}
	if value := state.phaseText(); value != "Raid in 02:00" {
		t.Fatalf("phaseText = %q", value)
	}
}

// TestStateFormatsRaidTopBar verifies raid-specific top-bar text.
func TestStateFormatsRaidTopBar(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	state.status.phase = phaseRaid
	state.raid.active = true
	state.raid.pendingEnemies = 5
	state.raid.enemies = []raidEnemy{{}, {}}

	if value := state.phaseText(); value != "Enemies remaining: 7" {
		t.Fatalf("phaseText = %q", value)
	}
}

// TestUpdateIncrementsWhenRunning verifies logical updates advance while unpaused.
func TestUpdateIncrementsWhenRunning(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
}

// TestTogglePauseDoesNotIncrement verifies pause input is not a logical update.
func TestTogglePauseDoesNotIncrement(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	if !state.Paused() {
		t.Fatal("expected state to pause")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
}

// TestPausedUpdatesDoNotIncrement verifies logical updates stop while paused.
func TestPausedUpdatesDoNotIncrement(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{})
	state.Update(Input{})
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 0)
	}
}

// TestUnpauseThenUpdateIncrements verifies updates resume after unpausing.
func TestUnpauseThenUpdateIncrements(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{TogglePause: true})
	if state.Paused() {
		t.Fatal("expected state to unpause")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates after toggle = %d, want %d", state.Updates(), 0)
	}

	state.Update(Input{})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
}

// TestEscapeOpensIngameMenuAndPauses verifies ESC pauses and shows the overlay.
func TestEscapeOpensIngameMenuAndPauses(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if action := state.Update(Input{ToggleMenu: true}); action != ActionNone {
		t.Fatalf("Update(escape) = %v, want %v", action, ActionNone)
	}
	if !state.IngameMenuOpen() {
		t.Fatal("expected in-game menu to open")
	}
	if !state.Paused() {
		t.Fatal("expected in-game menu to pause the game")
	}
}

// TestResumeRestoresRunningState verifies Resume unpauses a previously running game.
func TestResumeRestoresRunningState(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{ToggleMenu: true})
	state.Update(clickInput(state.ui.menu.buttons[0]))
	if state.IngameMenuOpen() {
		t.Fatal("expected Resume to close the in-game menu")
	}
	if state.Paused() {
		t.Fatal("expected Resume to restore running state")
	}
}

// TestResumeRestoresPausedState verifies Resume preserves an existing pause.
func TestResumeRestoresPausedState(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{TogglePause: true})
	state.Update(Input{ToggleMenu: true})
	state.Update(Input{ToggleMenu: true})
	if state.IngameMenuOpen() {
		t.Fatal("expected ESC to close the in-game menu")
	}
	if !state.Paused() {
		t.Fatal("expected Resume to restore prior paused state")
	}
}

// TestIngameMenuBlocksUpdates verifies overlay-open frames do not advance logic.
func TestIngameMenuBlocksUpdates(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{})
	state.Update(Input{ToggleMenu: true})
	state.Update(Input{})
	state.Update(Input{TogglePause: true})
	if state.Updates() != 1 {
		t.Fatalf("updates = %d, want %d", state.Updates(), 1)
	}
	if !state.IngameMenuOpen() {
		t.Fatal("expected in-game menu to remain open")
	}
}

// TestSurrenderAction verifies Surrender reports a game-level action.
func TestSurrenderAction(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Update(Input{ToggleMenu: true})
	if action := state.Update(clickInput(state.ui.menu.buttons[1])); action != ActionSurrender {
		t.Fatalf("Update(surrender click) = %v, want %v", action, ActionSurrender)
	}
}

// TestIngameMenuResizeRecentersButtons verifies resizing updates overlay hit geometry.
func TestIngameMenuResizeRecentersButtons(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	state.Resize(2560, 1440)
	state.Update(Input{ToggleMenu: true})
	resumeButton := state.ui.menu.buttons[0]
	if resumeButton.X+resumeButton.W/2 != 1280 {
		t.Fatalf("resume center x = %d, want %d", resumeButton.X+resumeButton.W/2, 1280)
	}

	state.Update(clickInput(resumeButton))
	if state.IngameMenuOpen() {
		t.Fatal("expected resized Resume button to close the in-game menu")
	}
}

// clickInput returns a click at the center of button.
func clickInput(button ui.Button[Action]) Input {
	return Input{
		CursorX: button.X + button.W/2,
		CursorY: button.Y + button.H/2,
		Clicked: true,
	}
}

// almostEqual reports whether two float values are close enough for camera tests.
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}
