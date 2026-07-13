package game

import "math/rand"

const (
	plotSize       = 15
	homePlotCenter = plotSize / 2
)

// Map owns the prototype game map.
type Map struct {
	Home Plot
}

// Plot owns a fixed 15x15 group of Tiles.
type Plot struct {
	Tiles [plotSize][plotSize]Tile
}

// Tile describes one space in a Plot.
type Tile struct {
	Terrain tileTerrain
	Feature tileFeature
	Tweak   uint16
}

type tileTerrain int

const (
	terrainEmpty tileTerrain = iota
	terrainRoad
	terrainForest
)

type tileFeature int

const (
	featureNone tileFeature = iota
	featureSanctum
	featureHouse
	featureBarracks
	featureDorm
	featureWoodcutter
	featureStoneQuarry
	featureIronMine
	featureBowTower
	featureFlameBoltTower
	featureCatapultTower
)

// NewDefaultMap creates the prototype map used by a new game.
func NewDefaultMap() Map {
	return Map{
		Home: NewDefaultHomePlot(),
	}
}

// NewDefaultHomePlot creates the starting Plot with a Sanctum, north road, and tree border.
func NewDefaultHomePlot() Plot {
	return newDefaultHomePlotWithTweakSource(randomTileTweak)
}

// newDefaultHomePlotWithTweakSource creates a home Plot with caller-provided Tile tweaks.
func newDefaultHomePlotWithTweakSource(nextTweak func() uint16) Plot {
	var plot Plot
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			plot.Tiles[y][x] = newTile(nextTweak)
			if x == 0 || y == 0 || x == plotSize-1 || y == plotSize-1 {
				plot.Tiles[y][x].Terrain = terrainForest
			}
		}
	}
	for y := 0; y <= homePlotCenter; y++ {
		plot.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
	plot.Tiles[homePlotCenter][homePlotCenter].Feature = featureSanctum
	return plot
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
