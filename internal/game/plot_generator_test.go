package game

import "testing"

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

// TestDefaultHomePlotUsesGrasslandsBiome verifies the starting Plot records its biome.
func TestDefaultHomePlotUsesGrasslandsBiome(t *testing.T) {
	plot := NewDefaultHomePlot()

	if plot.Biome != biomeGrasslands {
		t.Fatalf("home biome = %v, want grasslands", plot.Biome)
	}
}

// TestDefaultHomePlotStartsWithOnlyTheSanctum verifies the player receives no free structures.
func TestDefaultHomePlotStartsWithOnlyTheSanctum(t *testing.T) {
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

// TestDefaultHomePlotAssignsTileTweaks verifies Tile creation stores independent tweak values.
func TestDefaultHomePlotAssignsTileTweaks(t *testing.T) {
	var next uint16
	plot := newDefaultHomePlotWithSources(func() uint16 {
		value := next
		next++
		return value
	}, constantTerrainRoll(99))

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			want := uint16(y*plotSize + x)
			if plot.Tiles[y][x].Tweak != want {
				t.Fatalf("tile (%d,%d) tweak = %d, want %d", x, y, plot.Tiles[y][x].Tweak, want)
			}
		}
	}
}

// TestDefaultHomePlotUsesGrasslandsTerrainWeights verifies exact starting-Plot terrain boundaries.
func TestDefaultHomePlotUsesGrasslandsTerrainWeights(t *testing.T) {
	tests := []struct {
		roll int
		want tileTerrain
	}{
		{roll: 0, want: terrainTree},
		{roll: 5, want: terrainTree},
		{roll: 6, want: terrainBoulder},
		{roll: 8, want: terrainBoulder},
		{roll: 9, want: terrainEmpty},
		{roll: 99, want: terrainEmpty},
	}

	for _, test := range tests {
		plot := newDefaultHomePlotWithSources(func() uint16 {
			return 0
		}, constantTerrainRoll(test.roll))
		if got := plot.Tiles[0][0].Terrain; got != test.want {
			t.Errorf("roll %d home terrain = %v, want %v", test.roll, got, test.want)
		}
	}
}

// TestDefaultHomePlotRoadOverridesGrasslandsTerrain verifies generated obstacles cannot block the road.
func TestDefaultHomePlotRoadOverridesGrasslandsTerrain(t *testing.T) {
	plot := newDefaultHomePlotWithSources(func() uint16 {
		return 0
	}, constantTerrainRoll(6))

	for y := 0; y <= homePlotCenter; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainRoad {
			t.Fatalf("tile (%d,%d) terrain = %v, want road", homePlotCenter, y, plot.Tiles[y][homePlotCenter].Terrain)
		}
	}
	for y := homePlotCenter + 1; y < plotSize; y++ {
		if plot.Tiles[y][homePlotCenter].Terrain != terrainBoulder {
			t.Fatalf("tile (%d,%d) terrain = %v, want generated Boulder", homePlotCenter, y, plot.Tiles[y][homePlotCenter].Terrain)
		}
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
