package game

import (
	"math"
	"testing"
)

// TestSelectedTowerRangeIndicatorUsesTemplateRange verifies every combat tower's exact authored radius.
func TestSelectedTowerRangeIndicatorUsesTemplateRange(t *testing.T) {
	tests := []struct {
		name         string
		feature      tileFeature
		rangeInTiles float64
	}{
		{name: "bow tower", feature: featureBowTower, rangeInTiles: 3},
		{name: "flame bolt tower", feature: featureFlameBoltTower, rangeInTiles: 2.5},
		{name: "catapult tower", feature: featureCatapultTower, rangeInTiles: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newTowerRangeTestState(t)
			tile := homeTileCoordinate(homePlotCenter+2, homePlotCenter-1)
			state.gameMap.Home.Tiles[tile.Y][tile.X].Feature = test.feature
			state.selection = selectedItem{kind: selectedItemStructure, tile: tile}

			indicator, ok := state.selectedTowerRangeIndicator(state.sceneViewport())
			if !ok {
				t.Fatal("expected selected tower range indicator")
			}
			assertTowerRangeIndicator(t, state, tile, indicator, test.rangeInTiles)
		})
	}
}

// TestSelectedTowerRangeIndicatorTracksCameraAndPlot verifies projection across zoomed explored Plots.
func TestSelectedTowerRangeIndicatorTracksCameraAndPlot(t *testing.T) {
	state := newTowerRangeTestState(t)
	state.camera = camera{zoom: 1.75, centerX: 3.25, centerY: -2.5}
	plotCoord := plotCoordinate{X: 1, Y: -1}
	plot := &Plot{Biome: biomeGrasslands}
	state.gameMap.Plots[plotCoord] = plot
	tile := tileCoordinate{Plot: plotCoord, X: 4, Y: 9}
	plot.Tiles[tile.Y][tile.X].Feature = featureCatapultTower
	state.selection = selectedItem{kind: selectedItemStructure, tile: tile}

	indicator, ok := state.selectedTowerRangeIndicator(state.sceneViewport())
	if !ok {
		t.Fatal("expected selected tower range indicator in explored Plot")
	}
	assertTowerRangeIndicator(t, state, tile, indicator, 5)
}

// TestSelectedTowerRangeIndicatorIsSelectionAndPhaseIndependent verifies only tower identity gates the overlay.
func TestSelectedTowerRangeIndicatorIsSelectionAndPhaseIndependent(t *testing.T) {
	states := []struct {
		name     string
		phase    phase
		paused   bool
		breached bool
	}{
		{name: "labour", phase: phaseLabour},
		{name: "management", phase: phaseManagement},
		{name: "paused management", phase: phaseManagement, paused: true},
		{name: "raid", phase: phaseRaid},
		{name: "breached raid", phase: phaseRaid, breached: true},
	}

	for _, test := range states {
		t.Run(test.name, func(t *testing.T) {
			state := newTowerRangeTestState(t)
			tile := homeTileCoordinate(homePlotCenter+1, homePlotCenter)
			state.gameMap.Home.Tiles[tile.Y][tile.X].Feature = featureBowTower
			state.selection = selectedItem{kind: selectedItemStructure, tile: tile}
			state.status.phase = test.phase
			state.paused = test.paused
			state.raid.breached = test.breached

			if _, ok := state.selectedTowerRangeIndicator(state.sceneViewport()); !ok {
				t.Fatal("expected selected tower range indicator")
			}
		})
	}
}

// TestSelectedTowerRangeIndicatorRejectsNonTowers verifies stale or unrelated selections draw no range.
func TestSelectedTowerRangeIndicatorRejectsNonTowers(t *testing.T) {
	tests := []struct {
		name      string
		selection selectedItem
		feature   tileFeature
		zeroRange bool
	}{
		{name: "no selection"},
		{name: "raider", selection: selectedItem{kind: selectedItemRaider, raiderID: 7}},
		{name: "terrain", selection: selectedItem{kind: selectedItemTerrain, tile: homeTileCoordinate(2, 2)}},
		{name: "ordinary building", selection: selectedItem{kind: selectedItemStructure, tile: homeTileCoordinate(2, 2)}, feature: featureHouse},
		{name: "empty tile", selection: selectedItem{kind: selectedItemStructure, tile: homeTileCoordinate(2, 2)}},
		{name: "missing plot", selection: selectedItem{kind: selectedItemStructure, tile: tileCoordinate{Plot: plotCoordinate{X: 8}, X: 2, Y: 2}}, feature: featureBowTower},
		{name: "negative tile", selection: selectedItem{kind: selectedItemStructure, tile: tileCoordinate{X: -1, Y: 2}}, feature: featureBowTower},
		{name: "tile past edge", selection: selectedItem{kind: selectedItemStructure, tile: tileCoordinate{X: plotSize, Y: 2}}, feature: featureBowTower},
		{name: "tower with no range", selection: selectedItem{kind: selectedItemStructure, tile: homeTileCoordinate(2, 2)}, feature: featureBowTower, zeroRange: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newTowerRangeTestState(t)
			state.selection = test.selection
			if test.selection.tile.Plot == homePlotCoordinate &&
				test.selection.tile.X >= 0 && test.selection.tile.X < plotSize &&
				test.selection.tile.Y >= 0 && test.selection.tile.Y < plotSize {
				state.gameMap.Home.Tiles[test.selection.tile.Y][test.selection.tile.X].Feature = test.feature
			}
			if test.zeroRange {
				state.structureCatalog.BowTower.RangeTiles = 0
			}

			if indicator, ok := state.selectedTowerRangeIndicator(state.sceneViewport()); ok {
				t.Fatalf("unexpected range indicator: %+v", indicator)
			}
		})
	}
}

func newTowerRangeTestState(t *testing.T) *State {
	t.Helper()
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func assertTowerRangeIndicator(
	t *testing.T,
	state *State,
	tile tileCoordinate,
	indicator towerRangeIndicator,
	rangeInTiles float64,
) {
	t.Helper()
	viewport := state.sceneViewport()
	center := plotTileWorldCenter(tile.Plot, tile.X, tile.Y)
	scale := plotBaseTileSize * state.camera.zoom
	wantX := viewport.x + viewport.w/2 + (center.X-state.camera.centerX)*scale
	wantY := viewport.y + viewport.h/2 + (state.camera.centerY-center.Y)*scale
	wantRadius := rangeInTiles * scale
	assertTowerRangeValue(t, "center X", float64(indicator.centerX), wantX)
	assertTowerRangeValue(t, "center Y", float64(indicator.centerY), wantY)
	assertTowerRangeValue(t, "radius", float64(indicator.radius), wantRadius)
}

func assertTowerRangeValue(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.01 {
		t.Fatalf("%s = %.3f, want %.3f", name, got, want)
	}
}
