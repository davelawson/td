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

// NewDefaultHomePlot creates the starting Plot with only a Sanctum and north road.
func NewDefaultHomePlot() Plot {
	var plot Plot
	for y := 0; y <= homePlotCenter; y++ {
		plot.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
	plot.Tiles[homePlotCenter][homePlotCenter].Feature = featureSanctum
	return plot
}
