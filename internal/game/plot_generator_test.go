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

// TestGeneratedHillsPlotUsesHillsBiome verifies hills Plot biome metadata.
func TestGeneratedHillsPlotUsesHillsBiome(t *testing.T) {
	plot := NewHillsPlot()

	if plot.Biome != biomeHills {
		t.Fatalf("generated biome = %v, want hills", plot.Biome)
	}
}

// TestBiomeForRollSplitsExploredPlotsEvenly verifies biome selection boundaries.
func TestBiomeForRollSplitsExploredPlotsEvenly(t *testing.T) {
	tests := []struct {
		roll int
		want plotBiome
	}{
		{roll: 0, want: biomeGrasslands},
		{roll: 49, want: biomeGrasslands},
		{roll: 50, want: biomeHills},
		{roll: 99, want: biomeHills},
	}

	for _, test := range tests {
		if got := biomeForRoll(test.roll); got != test.want {
			t.Errorf("roll %d biome = %v, want %v", test.roll, got, test.want)
		}
	}
}

// TestGeneratedPlotUsesAssignedBiome verifies random sources stay independently testable.
func TestGeneratedPlotUsesAssignedBiome(t *testing.T) {
	plot := newPlotForBiomeWithSources(biomeHills, func() uint16 {
		return 0
	}, constantTerrainRoll(3))

	if plot.Biome != biomeHills {
		t.Fatalf("generated biome = %v, want hills", plot.Biome)
	}
	if plot.Tiles[0][0].Terrain != terrainBoulder {
		t.Fatalf("first terrain = %v, want hills Boulder", plot.Tiles[0][0].Terrain)
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

// TestGeneratedHillsPlotCanContainObstacles verifies stone-biased hills terrain generation.
func TestGeneratedHillsPlotCanContainObstacles(t *testing.T) {
	terrainRolls := repeatingTerrainRolls(0, 3, 9)
	plot := newHillsPlotWithSources(func() uint16 {
		return 0
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
				t.Fatalf("tile (%d,%d) terrain = %v, want hills terrain", x, y, plot.Tiles[y][x].Terrain)
			}
		}
	}
	if trees == 0 || boulders == 0 || empty == 0 {
		t.Fatalf("hills terrain counts = Tree %d, Boulder %d, empty %d; want every terrain", trees, boulders, empty)
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

// TestHillsTerrainWeightsBiasBoulders verifies the hills percentage boundaries.
func TestHillsTerrainWeightsBiasBoulders(t *testing.T) {
	tests := []struct {
		roll int
		want tileTerrain
	}{
		{roll: 0, want: terrainTree},
		{roll: 2, want: terrainTree},
		{roll: 3, want: terrainBoulder},
		{roll: 8, want: terrainBoulder},
		{roll: 9, want: terrainEmpty},
		{roll: 99, want: terrainEmpty},
	}

	for _, test := range tests {
		if got := weightedTerrain(hillsTerrainWeights, test.roll); got != test.want {
			t.Errorf("roll %d terrain = %v, want %v", test.roll, got, test.want)
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
