package game

import "testing"

// TestDefaultMapAssignsInitialFrontierBiomes verifies pre-exploration biome assignment order.
func TestDefaultMapAssignsInitialFrontierBiomes(t *testing.T) {
	gameMap := newDefaultMapWithBiomeSource(repeatingTerrainRolls(0, 50, 49, 99))
	want := map[plotCoordinate]plotBiome{
		{X: 0, Y: 1}:  biomeGrasslands,
		{X: 1, Y: 0}:  biomeHills,
		{X: 0, Y: -1}: biomeGrasslands,
		{X: -1, Y: 0}: biomeHills,
	}

	if len(gameMap.frontierBiomes) != len(want) {
		t.Fatalf("frontier biomes = %d, want %d", len(gameMap.frontierBiomes), len(want))
	}
	for coord, wantBiome := range want {
		if got, ok := gameMap.frontierBiome(coord); !ok || got != wantBiome {
			t.Errorf("frontier %+v biome = %v, assigned = %v; want %v", coord, got, ok, wantBiome)
		}
	}
}

// TestRevealUsesAssignedFrontierBiome verifies preview and generated Plot stay consistent.
func TestRevealUsesAssignedFrontierBiome(t *testing.T) {
	gameMap := newDefaultMapWithBiomeSource(constantTerrainRoll(50))
	target := plotCoordinate{X: 1, Y: 0}

	gameMap.revealPlot(target)

	plot, ok := gameMap.plot(target)
	if !ok {
		t.Fatal("expected assigned frontier Plot to be revealed")
	}
	if plot.Biome != biomeHills {
		t.Fatalf("revealed biome = %v, want assigned hills", plot.Biome)
	}
	if _, assigned := gameMap.frontierBiome(target); assigned {
		t.Fatal("expected explored Plot to leave frontier biome storage")
	}
}

// TestRevealAssignsNewFrontierBiomes verifies expansion previews are created once.
func TestRevealAssignsNewFrontierBiomes(t *testing.T) {
	rolls := repeatingTerrainRolls(0, 0, 0, 0, 50, 0, 50, 0, 0)
	gameMap := newDefaultMapWithBiomeSource(rolls)
	north := plotCoordinate{X: 0, Y: 1}
	gameMap.revealPlotWithBiomeSource(north, rolls)

	want := map[plotCoordinate]plotBiome{
		{X: 0, Y: 2}:  biomeHills,
		{X: 1, Y: 1}:  biomeGrasslands,
		{X: -1, Y: 1}: biomeHills,
	}
	for coord, wantBiome := range want {
		if got, ok := gameMap.frontierBiome(coord); !ok || got != wantBiome {
			t.Errorf("new frontier %+v biome = %v, assigned = %v; want %v", coord, got, ok, wantBiome)
		}
	}

	shared := plotCoordinate{X: 1, Y: 1}
	before, _ := gameMap.frontierBiome(shared)
	gameMap.revealPlotWithBiomeSource(plotCoordinate{X: 1, Y: 0}, rolls)
	after, ok := gameMap.frontierBiome(shared)
	if !ok || after != before {
		t.Fatalf("shared frontier biome = %v, assigned = %v; want retained %v", after, ok, before)
	}
}

// TestRevealWithoutFrontierAssignmentDoesNothing verifies reveal never decides a biome late.
func TestRevealWithoutFrontierAssignmentDoesNothing(t *testing.T) {
	gameMap := NewDefaultMap()
	target := plotCoordinate{X: 2, Y: 0}

	gameMap.revealPlot(target)

	if gameMap.explored(target) {
		t.Fatal("expected unassigned non-frontier Plot to remain unexplored")
	}
}
