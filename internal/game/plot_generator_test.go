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
	terrainRolls := repeatingTerrainRolls(0, 6, 9)
	plot := newGrasslandsPlotWithSources(func() uint16 {
		value := next
		next++
		return value
	}, terrainRolls)

	trees := 0
	boulders := 0
	empty := 0
	for y := range plot.Tiles {
		for x := range plot.Tiles[y] {
			switch plot.Tiles[y][x].Terrain {
			case terrainTree:
				trees++
			case terrainBoulder:
				boulders++
			case terrainEmpty:
				empty++
			default:
				t.Fatalf("tile (%d,%d) terrain = %v, want grasslands terrain", x, y, plot.Tiles[y][x].Terrain)
			}
		}
	}
	if trees == 0 {
		t.Fatal("expected deterministic grasslands generation to include Tree")
	}
	if boulders == 0 {
		t.Fatal("expected deterministic grasslands generation to include Boulder")
	}
	if empty == 0 {
		t.Fatal("expected deterministic grasslands generation to keep buildable grass")
	}
}

// TestDefaultHomePlotDoesNotGenerateObstacles verifies the starting Plot stays forgiving.
func TestDefaultHomePlotDoesNotGenerateObstacles(t *testing.T) {
	plot := newDefaultHomePlotWithTweakSource(func() uint16 {
		return 0
	})

	for y := range plot.Tiles {
		for x := range plot.Tiles[y] {
			if plot.Tiles[y][x].Terrain == terrainTree || plot.Tiles[y][x].Terrain == terrainBoulder {
				t.Fatalf("home tile (%d,%d) terrain = %v, want no generated obstacle", x, y, plot.Tiles[y][x].Terrain)
			}
		}
	}
}

// TestWeightedTerrainSelectsTree verifies Tree uses the first weight range.
func TestWeightedTerrainSelectsTree(t *testing.T) {
	weights := terrainWeights{Tree: 6, Boulder: 3}

	for _, roll := range []int{0, 5} {
		if got := weightedTerrain(weights, roll); got != terrainTree {
			t.Fatalf("roll %d terrain = %v, want Tree", roll, got)
		}
	}
}

// TestWeightedTerrainSelectsBoulder verifies Boulder follows Tree in the weight range.
func TestWeightedTerrainSelectsBoulder(t *testing.T) {
	weights := terrainWeights{Tree: 6, Boulder: 3}

	for _, roll := range []int{6, 8} {
		if got := weightedTerrain(weights, roll); got != terrainBoulder {
			t.Fatalf("roll %d terrain = %v, want Boulder", roll, got)
		}
	}
}

// TestWeightedTerrainSelectsEmpty verifies unweighted percentages stay empty.
func TestWeightedTerrainSelectsEmpty(t *testing.T) {
	weights := terrainWeights{Tree: 6, Boulder: 3}

	for _, roll := range []int{9, 99} {
		if got := weightedTerrain(weights, roll); got != terrainEmpty {
			t.Fatalf("roll %d terrain = %v, want empty", roll, got)
		}
	}
}

// TestWeightedTerrainRejectsInvalidWeights verifies invalid percentages fall back to empty.
func TestWeightedTerrainRejectsInvalidWeights(t *testing.T) {
	for _, weights := range []terrainWeights{
		{Tree: -1, Boulder: 3},
		{Tree: 6, Boulder: -1},
		{Tree: 75, Boulder: 50},
	} {
		if got := weightedTerrain(weights, 0); got != terrainEmpty {
			t.Fatalf("weights %+v terrain = %v, want empty", weights, got)
		}
	}
	for _, roll := range []int{-1, 100} {
		if got := weightedTerrain(terrainWeights{Tree: 6, Boulder: 3}, roll); got != terrainEmpty {
			t.Fatalf("roll %d terrain = %v, want empty", roll, got)
		}
	}
}

// TestNorthRoadOverridesGeneratedBoulder verifies road generation protects Raid paths.
func TestNorthRoadOverridesGeneratedBoulder(t *testing.T) {
	plot := newGrasslandsPlotWithSources(func() uint16 {
		return 0
	}, constantTerrainRoll(6))

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
	plot := newGrasslandsPlotWithSources(func() uint16 {
		return 0
	}, constantTerrainRoll(6))
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

// constantTerrainRoll returns a terrain roll source that always returns one value.
func constantTerrainRoll(roll int) func() int {
	return func() int {
		return roll
	}
}

// repeatingTerrainRolls returns a terrain roll source that cycles through values.
func repeatingTerrainRolls(rolls ...int) func() int {
	var next int
	return func() int {
		roll := rolls[next%len(rolls)]
		next++
		return roll
	}
}
