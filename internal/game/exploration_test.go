package game

import "testing"

// TestNewMapStartsWithOnlyHomePlotExplored verifies the starting Domain size.
func TestNewMapStartsWithOnlyHomePlotExplored(t *testing.T) {
	gameMap := NewDefaultMap()

	coords := gameMap.exploredPlotCoordinates()
	if len(coords) != 1 {
		t.Fatalf("explored plots = %d, want 1", len(coords))
	}
	if coords[0] != homePlotCoordinate {
		t.Fatalf("explored plot = %+v, want home", coords[0])
	}
}

// TestInitialExploreButtonsTargetOrthogonalNeighbors verifies the first reveal choices.
func TestInitialExploreButtonsTargetOrthogonalNeighbors(t *testing.T) {
	state := newRaidTestState(t)

	buttons := state.exploreButtons()
	if len(buttons) != 4 {
		t.Fatalf("explore buttons = %d, want 4", len(buttons))
	}
	want := map[plotCoordinate]bool{
		{X: 0, Y: 1}:  true,
		{X: 1, Y: 0}:  true,
		{X: 0, Y: -1}: true,
		{X: -1, Y: 0}: true,
	}
	for _, button := range buttons {
		if !want[button.Target] {
			t.Fatalf("unexpected explore target %+v", button.Target)
		}
		biome, ok := state.gameMap.frontierBiome(button.Target)
		if !ok || button.Biome != biome {
			t.Fatalf("button biome = %v, stored biome = %v, assigned = %v", button.Biome, biome, ok)
		}
		delete(want, button.Target)
	}
	if len(want) != 0 {
		t.Fatalf("missing explore targets: %+v", want)
	}
}

// TestExploreClickRevealsAdjacentBiomeAndClearsSharedEdge verifies the reveal action.
func TestExploreClickRevealsAdjacentBiomeAndClearsSharedEdge(t *testing.T) {
	state := newRaidTestState(t)
	target := plotCoordinate{X: 1, Y: 0}
	preview, ok := state.gameMap.frontierBiome(target)
	if !ok {
		t.Fatal("expected east Plot biome to be assigned before exploration")
	}

	state.Update(clickExploreButtonInput(t, state, target))

	plot, ok := state.gameMap.plot(target)
	if !ok {
		t.Fatal("expected east plot to be explored")
	}
	if plot.Biome != preview {
		t.Fatalf("east plot biome = %v, want previewed biome %v", plot.Biome, preview)
	}
	if plot.Tiles[homePlotCenter][homePlotCenter].Feature != featureNone {
		t.Fatalf("east plot center feature = %v, want none", plot.Tiles[homePlotCenter][homePlotCenter].Feature)
	}
	if state.gameMap.Home.Tiles[homePlotCenter][plotSize-1].Terrain != terrainEmpty {
		t.Fatal("expected home east shared edge to become grass")
	}
	if plot.Tiles[homePlotCenter][0].Terrain != terrainEmpty {
		t.Fatal("expected new plot west shared edge to become grass")
	}
}

// TestRevealExistingPlotPreservesGeneratedBiomeAndTerrain verifies reveal idempotence.
func TestRevealExistingPlotPreservesGeneratedBiomeAndTerrain(t *testing.T) {
	gameMap := NewDefaultMap()
	target := plotCoordinate{X: 1, Y: 0}
	gameMap.revealPlot(target)
	plot, _ := gameMap.plot(target)
	plot.Tiles[2][2].Terrain = terrainBoulder

	gameMap.revealPlot(target)
	again, _ := gameMap.plot(target)

	if again != plot {
		t.Fatal("expected repeated reveal to preserve the stored Plot")
	}
	if again.Tiles[2][2].Terrain != terrainBoulder {
		t.Fatalf("retained terrain = %v, want Boulder", again.Tiles[2][2].Terrain)
	}
}

// TestExploreWorksWhilePausedManagement verifies paused Management still allows preparation actions.
func TestExploreWorksWhilePausedManagement(t *testing.T) {
	state := newRaidTestState(t)
	target := plotCoordinate{X: -1, Y: 0}
	state.Update(Input{TogglePause: true})

	state.Update(clickExploreButtonInput(t, state, target))

	if !state.gameMap.explored(target) {
		t.Fatal("expected paused Management exploration to reveal the target plot")
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want 0 while paused", state.Updates())
	}
}

// TestExploreBlockedDuringRaid verifies Raid phase blocks Plot reveals.
func TestExploreBlockedDuringRaid(t *testing.T) {
	state := newRaidTestState(t)
	target := plotCoordinate{X: 0, Y: 1}
	state.startNextRaid()

	state.Update(clickExploreButtonInput(t, state, target))

	if state.gameMap.explored(target) {
		t.Fatal("expected active Raid to block exploration")
	}
}

// TestNorthExplorationAddsRoadAndExtendsRaidSpawn verifies the first path extension.
func TestNorthExplorationAddsRoadAndExtendsRaidSpawn(t *testing.T) {
	state := newRaidTestState(t)
	target := plotCoordinate{X: 0, Y: 1}

	state.Update(clickExploreButtonInput(t, state, target))
	plot, ok := state.gameMap.plot(target)
	if !ok {
		t.Fatal("expected north plot to be explored")
	}
	for y := 0; y < plotSize; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainRoad {
			t.Fatalf("north plot road tile y=%d terrain = %v, want road", y, plot.Tiles[y][homePlotCenter].Terrain)
		}
	}

	state.startNextRaid()
	state.spawnRaidEnemy(raidEnemySkeletonSwordShield)

	if got, want := state.raid.enemies[0].position, (coord{X: 0, Y: float64(plotSize + homePlotCenter)}); got != want {
		t.Fatalf("spawn position = %+v, want %+v", got, want)
	}
}

// TestBuildDragPlacesStructureOnExploredPlot verifies revealed Plots behave like home for building.
func TestBuildDragPlacesStructureOnExploredPlot(t *testing.T) {
	state := newRaidTestState(t)
	targetPlot := plotCoordinate{X: 1, Y: 0}
	targetTile := tileCoordinate{Plot: targetPlot, X: 2, Y: 5}
	state.gameMap.revealPlot(targetPlot)
	plot, _ := state.gameMap.plot(targetPlot)
	plot.Tiles[targetTile.Y][targetTile.X].Terrain = terrainEmpty

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	state.Update(releasePlotTileInput(state, targetTile))

	if plot.Tiles[targetTile.Y][targetTile.X].Feature != featureHouse {
		t.Fatalf("tile feature = %v, want House", plot.Tiles[targetTile.Y][targetTile.X].Feature)
	}
}

// TestBuildDropOnExploreButtonDoesNotPlaceStructure prevents command overlap on grass borders.
func TestBuildDropOnExploreButtonDoesNotPlaceStructure(t *testing.T) {
	state := newRaidTestState(t)
	explored := plotCoordinate{X: 1, Y: 0}
	target := plotCoordinate{X: 2, Y: 0}
	state.gameMap.revealPlot(explored)
	initialResources := state.status.resources

	state.Update(pressBuildingBarItemInput(state, buildingBarHouseIndex))
	release := clickExploreButtonInput(t, state, target)
	release.Clicked = false
	release.MouseDown = false
	release.Released = true
	state.Update(release)

	if state.gameMap.explored(target) {
		t.Fatal("expected build release not to trigger exploration")
	}
	plot, _ := state.gameMap.plot(explored)
	if plot.Tiles[homePlotCenter][plotSize-1].Feature != featureNone {
		t.Fatalf("border feature = %v, want none", plot.Tiles[homePlotCenter][plotSize-1].Feature)
	}
	if state.status.resources != initialResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, initialResources)
	}
}

// TestSelectionWorksOnExploredPlot verifies revealed Plot structures can be inspected.
func TestSelectionWorksOnExploredPlot(t *testing.T) {
	state := newRaidTestState(t)
	targetPlot := plotCoordinate{X: 1, Y: 0}
	targetTile := tileCoordinate{Plot: targetPlot, X: 2, Y: 5}
	state.gameMap.revealPlot(targetPlot)
	plot, _ := state.gameMap.plot(targetPlot)
	plot.Tiles[targetTile.Y][targetTile.X].Feature = featureHouse

	state.Update(clickPlotTileInput(state, targetTile))

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != targetTile {
		t.Fatalf("selected tile = %+v, want %+v", state.selection.tile, targetTile)
	}
}

// TestEconomicBuildingOnExploredPlotWorksDuringLabour verifies yields include every explored Plot.
func TestEconomicBuildingOnExploredPlotWorksDuringLabour(t *testing.T) {
	state := newRaidTestState(t)
	targetPlot := plotCoordinate{X: 1, Y: 0}
	state.gameMap.revealPlot(targetPlot)
	plot, _ := state.gameMap.plot(targetPlot)
	plot.Tiles[5][2].Feature = featureWoodcutter
	startingResources := state.status.resources

	state.completeRaid()

	if state.status.resources.wood != startingResources.wood+10 {
		t.Fatalf("wood = %d, want %d", state.status.resources.wood, startingResources.wood+10)
	}
}

// TestCombatTowerOnExploredPlotFires verifies tower range checks use world coordinates across Plots.
func TestCombatTowerOnExploredPlotFires(t *testing.T) {
	state := newRaidTestState(t)
	targetPlot := plotCoordinate{X: 1, Y: 0}
	state.gameMap.revealPlot(targetPlot)
	plot, _ := state.gameMap.plot(targetPlot)
	plot.Tiles[homePlotCenter][0].Feature = featureBowTower
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 8, Y: 0}, 20)}

	state.updateCombat()

	if len(state.combat.projectiles) != 1 {
		t.Fatalf("projectiles = %d, want 1", len(state.combat.projectiles))
	}
}

func clickExploreButtonInput(t *testing.T, state *State, target plotCoordinate) Input {
	t.Helper()
	for _, button := range state.exploreButtons() {
		if button.Target != target {
			continue
		}
		rect := state.projectRect(
			state.sceneViewport(),
			button.Center.X-exploreButtonSize/2,
			button.Center.Y+exploreButtonSize/2,
			exploreButtonSize,
			exploreButtonSize,
		)
		return clickProjectedRectInput(rect)
	}
	t.Fatalf("no explore button for target %+v", target)
	return Input{}
}

func releasePlotTileInput(state *State, tile tileCoordinate) Input {
	input := clickPlotTileInput(state, tile)
	input.Clicked = false
	input.MouseDown = false
	input.Released = true
	return input
}

func clickPlotTileInput(state *State, tile tileCoordinate) Input {
	worldWest, worldNorth, worldW, worldH := plotTileWorldRect(tile.Plot, tile.X, tile.Y)
	rect := state.projectRect(state.sceneViewport(), worldWest, worldNorth, worldW, worldH)
	return clickProjectedRectInput(rect)
}
