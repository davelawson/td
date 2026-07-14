package game

import "math/rand"

var grasslandsTerrainWeights = terrainWeights{
	Tree:    6,
	Boulder: 3,
}

// terrainWeights describes percentage chances for generated terrain in a Tile.
type terrainWeights struct {
	Tree    int
	Boulder int
}

// NewDefaultHomePlot creates the grassland starting Plot with a Sanctum and north road.
func NewDefaultHomePlot() Plot {
	return newDefaultHomePlotWithTweakSource(randomTileTweak)
}

// newDefaultHomePlotWithTweakSource creates a home Plot with caller-provided Tile tweaks.
func newDefaultHomePlotWithTweakSource(nextTweak func() uint16) Plot {
	plot := newOpenGrasslandsPlotWithTweakSource(nextTweak)
	for y := 0; y <= homePlotCenter; y++ {
		plot.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
	plot.Tiles[homePlotCenter][homePlotCenter].Feature = featureSanctum
	return plot
}

// NewGrasslandsPlot creates a generated grasslands Plot for exploration.
func NewGrasslandsPlot() Plot {
	return newGrasslandsPlotWithSources(randomTileTweak, randomTerrainRoll)
}

// newGrasslandsPlotWithSources creates a grasslands Plot with caller-provided random sources.
func newGrasslandsPlotWithSources(nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	plot := newOpenGrasslandsPlotWithTweakSource(nextTweak)
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			tile := &plot.Tiles[y][x]
			tile.Terrain = weightedTerrain(grasslandsTerrainWeights, nextTerrainRoll())
		}
	}
	return plot
}

// newOpenGrasslandsPlotWithTweakSource creates a grasslands Plot with only empty grass Tiles.
func newOpenGrasslandsPlotWithTweakSource(nextTweak func() uint16) Plot {
	var plot Plot
	plot.Biome = biomeGrasslands
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			plot.Tiles[y][x] = newTile(nextTweak)
		}
	}
	return plot
}

// weightedTerrain returns the generated terrain selected by a percentage roll.
func weightedTerrain(weights terrainWeights, roll int) tileTerrain {
	if roll < weights.Tree {
		return terrainTree
	}
	roll -= weights.Tree
	if roll < weights.Boulder {
		return terrainBoulder
	}
	return terrainEmpty
}

// newTile creates one Tile with its prototype variation tweak assigned.
func newTile(nextTweak func() uint16) Tile {
	return Tile{
		Tweak: nextTweak(),
	}
}

// randomTileTweak returns a random unsigned 16-bit value for Tile variation.
func randomTileTweak() uint16 {
	return uint16(rand.Intn(1 << 16))
}

// randomTerrainRoll returns a random percentage roll for generated terrain.
func randomTerrainRoll() int {
	return rand.Intn(100)
}
