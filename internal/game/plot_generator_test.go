package game

import "testing"

// TestDefaultHomePlotUsesGrasslandsBiome verifies the starting Plot records its biome.
func TestDefaultHomePlotUsesGrasslandsBiome(t *testing.T) {
	plot := NewDefaultHomePlot()

	if plot.Biome != biomeGrasslands {
		t.Fatalf("home biome = %v, want grasslands", plot.Biome)
	}
}

// TestGeneratedGrasslandsPlotUsesGrasslandsBiome verifies explored Plot biome metadata.
func TestGeneratedGrasslandsPlotUsesGrasslandsBiome(t *testing.T) {
	plot := NewGrasslandsPlot()

	if plot.Biome != biomeGrasslands {
		t.Fatalf("generated biome = %v, want grasslands", plot.Biome)
	}
}

// TestGeneratedGrasslandsPlotCanContainObstacles verifies sparse terrain generation.
func TestGeneratedGrasslandsPlotCanContainObstacles(t *testing.T) {
	var next uint16
	plot := newGrasslandsPlotWithTweakSource(func() uint16 {
		value := next
		next++
		return value
	})

	forests := 0
	boulders := 0
	empty := 0
	for y := range plot.Tiles {
		for x := range plot.Tiles[y] {
			switch plot.Tiles[y][x].Terrain {
			case terrainForest:
				forests++
			case terrainBoulder:
				boulders++
			case terrainEmpty:
				empty++
			default:
				t.Fatalf("tile (%d,%d) terrain = %v, want grasslands terrain", x, y, plot.Tiles[y][x].Terrain)
			}
		}
	}
	if forests == 0 {
		t.Fatal("expected deterministic grasslands generation to include forest")
	}
	if boulders == 0 {
		t.Fatal("expected deterministic grasslands generation to include Boulder")
	}
	if empty == 0 {
		t.Fatal("expected deterministic grasslands generation to keep buildable grass")
	}
}

// TestDefaultHomePlotDoesNotGenerateForest verifies the starting Plot stays forgiving.
func TestDefaultHomePlotDoesNotGenerateForest(t *testing.T) {
	plot := newDefaultHomePlotWithTweakSource(func() uint16 {
		return 0
	})

	for y := range plot.Tiles {
		for x := range plot.Tiles[y] {
			if plot.Tiles[y][x].Terrain == terrainForest || plot.Tiles[y][x].Terrain == terrainBoulder {
				t.Fatalf("home tile (%d,%d) terrain = %v, want no generated obstacle", x, y, plot.Tiles[y][x].Terrain)
			}
		}
	}
}

// TestGrasslandsForestGenerationUsesTweakModulo verifies the sparse forest rule.
func TestGrasslandsForestGenerationUsesTweakModulo(t *testing.T) {
	if !grasslandsTileIsForest(grasslandsForestTweakModulo) {
		t.Fatal("expected tweak matching forest modulo to generate forest")
	}
	if grasslandsTileIsForest(1) {
		t.Fatal("expected tweak 1 to remain empty grass")
	}
}

// TestGrasslandsBoulderGenerationUsesTweakModulo verifies the sparse Boulder rule.
func TestGrasslandsBoulderGenerationUsesTweakModulo(t *testing.T) {
	if !grasslandsTileIsBoulder(0) {
		t.Fatal("expected tweak 0 to generate Boulder")
	}
	if !grasslandsTileIsBoulder(grasslandsBoulderTweakModulo) {
		t.Fatal("expected tweak matching Boulder modulo to generate Boulder")
	}
	if grasslandsTileIsBoulder(1) {
		t.Fatal("expected tweak 1 to avoid Boulder")
	}
}

// TestGrasslandsBoulderWinsOverForest verifies overlapping generation is deterministic.
func TestGrasslandsBoulderWinsOverForest(t *testing.T) {
	overlap := grasslandsForestTweakModulo * grasslandsBoulderTweakModulo
	plot := newGrasslandsPlotWithTweakSource(func() uint16 {
		return overlap
	})

	if plot.Tiles[0][0].Terrain != terrainBoulder {
		t.Fatalf("overlap terrain = %v, want Boulder", plot.Tiles[0][0].Terrain)
	}
}

// TestNorthRoadOverridesGeneratedBoulder verifies road generation protects Raid paths.
func TestNorthRoadOverridesGeneratedBoulder(t *testing.T) {
	plot := newGrasslandsPlotWithTweakSource(func() uint16 {
		return 0
	})

	applyNorthRoadIfNeeded(plotCoordinate{X: 0, Y: 1}, &plot)

	for y := 0; y < plotSize; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainRoad {
			t.Fatalf("road tile y=%d terrain = %v, want road", y, plot.Tiles[y][homePlotCenter].Terrain)
		}
	}
}

// TestSharedEdgeClearingOverridesGeneratedBoulder verifies joined Plots stay open.
func TestSharedEdgeClearingOverridesGeneratedBoulder(t *testing.T) {
	gameMap := NewDefaultMap()
	coord := plotCoordinate{X: 1, Y: 0}
	plot := newGrasslandsPlotWithTweakSource(func() uint16 {
		return 0
	})
	gameMap.Plots[coord] = &plot

	gameMap.clearSharedEdges(coord)

	for y := 0; y < plotSize; y++ {
		if plot.Tiles[y][0].Terrain != terrainEmpty {
			t.Fatalf("new west edge y=%d terrain = %v, want empty", y, plot.Tiles[y][0].Terrain)
		}
		if gameMap.Home.Tiles[y][plotSize-1].Terrain != terrainEmpty {
			t.Fatalf("home east edge y=%d terrain = %v, want empty", y, gameMap.Home.Tiles[y][plotSize-1].Terrain)
		}
	}
}
