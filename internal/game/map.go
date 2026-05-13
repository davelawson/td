package game

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
)

// NewDefaultMap creates the prototype map used by a new game.
func NewDefaultMap() Map {
	return Map{
		Home: NewDefaultHomePlot(),
	}
}

// NewDefaultHomePlot creates the starting Plot with a Sanctum, north road, and tree border.
func NewDefaultHomePlot() Plot {
	var plot Plot
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
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
