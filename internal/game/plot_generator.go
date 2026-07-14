package game

import "math/rand"

var grasslandsTerrainWeights = terrainWeights{
	Tree:    6,
	Boulder: 3,
}

var hillsTerrainWeights = terrainWeights{
	Tree:    3,
	Boulder: 6,
}

// terrainWeights describes percentage chances for generated terrain in a Tile.
type terrainWeights struct {
	Tree    int
	Boulder int
}

// NewDefaultHomePlot creates the grassland starting Plot with a Sanctum and north road.
func NewDefaultHomePlot() Plot {
	return newDefaultHomePlotWithSources(randomTileTweak, randomPercentageRoll)
}

// newDefaultHomePlotWithSources creates a grasslands home Plot with caller-provided random sources.
func newDefaultHomePlotWithSources(nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	plot := newGeneratedPlotWithSources(biomeGrasslands, grasslandsTerrainWeights, nextTweak, nextTerrainRoll)
	for y := 0; y <= homePlotCenter; y++ {
		plot.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
	plot.Tiles[homePlotCenter][homePlotCenter].Feature = featureSanctum
	return plot
}

// NewGrasslandsPlot creates a generated grasslands Plot for exploration.
func NewGrasslandsPlot() Plot {
	return newGrasslandsPlotWithSources(randomTileTweak, randomPercentageRoll)
}

// newGrasslandsPlotWithSources creates a grasslands Plot with caller-provided random sources.
func newGrasslandsPlotWithSources(nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	return newGeneratedPlotWithSources(biomeGrasslands, grasslandsTerrainWeights, nextTweak, nextTerrainRoll)
}

// NewHillsPlot creates a generated hills Plot for exploration.
func NewHillsPlot() Plot {
	return newHillsPlotWithSources(randomTileTweak, randomPercentageRoll)
}

// newHillsPlotWithSources creates a hills Plot with caller-provided random sources.
func newHillsPlotWithSources(nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	return newGeneratedPlotWithSources(biomeHills, hillsTerrainWeights, nextTweak, nextTerrainRoll)
}

// newPlotForBiome creates a generated Plot for a previously assigned biome.
func newPlotForBiome(biome plotBiome) Plot {
	return newPlotForBiomeWithSources(biome, randomTileTweak, randomPercentageRoll)
}

// newPlotForBiomeWithSources creates an assigned-biome Plot with caller-provided random sources.
func newPlotForBiomeWithSources(biome plotBiome, nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	switch biome {
	case biomeHills:
		return newHillsPlotWithSources(nextTweak, nextTerrainRoll)
	default:
		return newGrasslandsPlotWithSources(nextTweak, nextTerrainRoll)
	}
}

// biomeForRoll selects an explored Plot biome from a percentage roll.
func biomeForRoll(roll int) plotBiome {
	if roll < 50 {
		return biomeGrasslands
	}
	return biomeHills
}

// newGeneratedPlotWithSources creates a biome Plot using explicit terrain weights.
func newGeneratedPlotWithSources(biome plotBiome, weights terrainWeights, nextTweak func() uint16, nextTerrainRoll func() int) Plot {
	var plot Plot
	plot.Biome = biome
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			plot.Tiles[y][x] = newTile(nextTweak)
			plot.Tiles[y][x].Terrain = weightedTerrain(weights, nextTerrainRoll())
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

// randomPercentageRoll returns a random value from zero through 99.
func randomPercentageRoll() int {
	return rand.Intn(100)
}
