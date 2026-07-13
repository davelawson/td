package game

import "math/rand"

const (
	grasslandsForestTweakModulo  uint16 = 17
	grasslandsBoulderTweakModulo uint16 = 29
)

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
	return newGrasslandsPlotWithTweakSource(randomTileTweak)
}

// newGrasslandsPlotWithTweakSource creates a grasslands Plot with caller-provided Tile tweaks.
func newGrasslandsPlotWithTweakSource(nextTweak func() uint16) Plot {
	plot := newOpenGrasslandsPlotWithTweakSource(nextTweak)
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			tile := &plot.Tiles[y][x]
			if grasslandsTileIsBoulder(tile.Tweak) {
				tile.Terrain = terrainBoulder
				continue
			}
			if grasslandsTileIsForest(tile.Tweak) {
				tile.Terrain = terrainForest
			}
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

// grasslandsTileIsForest reports whether a grasslands Tile should contain sparse forest.
func grasslandsTileIsForest(tweak uint16) bool {
	return tweak%grasslandsForestTweakModulo == 0
}

// grasslandsTileIsBoulder reports whether a grasslands Tile should contain sparse Boulder.
func grasslandsTileIsBoulder(tweak uint16) bool {
	return tweak%grasslandsBoulderTweakModulo == 0
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
