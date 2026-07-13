package game

import "sort"

const (
	plotSize       = 15
	homePlotCenter = plotSize / 2
)

// Map owns the prototype game map.
type Map struct {
	Home  Plot
	Plots map[plotCoordinate]*Plot
}

// Plot owns a fixed 15x15 group of Tiles.
type Plot struct {
	Biome plotBiome
	Tiles [plotSize][plotSize]Tile
}

type plotBiome int

const (
	biomeGrasslands plotBiome = iota
)

// plotCoordinate identifies one Plot, with the home Plot at (0,0).
type plotCoordinate struct {
	X int
	Y int
}

var homePlotCoordinate = plotCoordinate{}

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
	terrainBoulder
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
	gameMap := Map{
		Home: NewDefaultHomePlot(),
	}
	gameMap.ensurePlots()
	return gameMap
}

// ensurePlots initializes the explored Plot map for tests that construct Map values directly.
func (m *Map) ensurePlots() {
	if m.Plots == nil {
		m.Plots = map[plotCoordinate]*Plot{}
	}
	m.Plots[homePlotCoordinate] = &m.Home
}

// plot returns the explored Plot at a coordinate.
func (m *Map) plot(coord plotCoordinate) (*Plot, bool) {
	m.ensurePlots()
	plot, ok := m.Plots[coord]
	return plot, ok
}

// explored reports whether a Plot coordinate is visible and usable.
func (m *Map) explored(coord plotCoordinate) bool {
	_, ok := m.plot(coord)
	return ok
}

// revealPlot marks an adjacent Plot explored without replacing existing content.
func (m *Map) revealPlot(coord plotCoordinate) {
	m.ensurePlots()
	if _, ok := m.Plots[coord]; ok {
		return
	}
	plot := NewGrasslandsPlot()
	applyNorthRoadIfNeeded(coord, &plot)
	m.Plots[coord] = &plot
	m.clearSharedEdges(coord)
}

// applyNorthRoadIfNeeded adds the visible straight Raid road to central north-chain Plots.
func applyNorthRoadIfNeeded(coord plotCoordinate, plot *Plot) {
	if coord.X != 0 || coord.Y <= 0 {
		return
	}
	for y := 0; y < plotSize; y++ {
		plot.Tiles[y][homePlotCenter].Terrain = terrainRoad
	}
}

// clearSharedEdges clears borders shared by two explored Plots, preserving road connectors.
func (m *Map) clearSharedEdges(coord plotCoordinate) {
	neighbors := orthogonalPlotNeighbors(coord)
	for _, neighbor := range neighbors {
		if !m.explored(neighbor) {
			continue
		}
		m.clearSharedEdge(coord, neighbor)
	}
}

// clearSharedEdge clears the border between two adjacent explored Plots.
func (m *Map) clearSharedEdge(a, b plotCoordinate) {
	plotA, okA := m.plot(a)
	plotB, okB := m.plot(b)
	if !okA || !okB {
		return
	}

	dx := b.X - a.X
	dy := b.Y - a.Y
	switch {
	case dx == 1 && dy == 0:
		clearVerticalSharedEdge(plotA, plotB, plotSize-1, 0)
	case dx == -1 && dy == 0:
		clearVerticalSharedEdge(plotA, plotB, 0, plotSize-1)
	case dx == 0 && dy == 1:
		clearHorizontalSharedEdge(plotA, plotB, 0, plotSize-1, isCentralNorthRoad(a, b))
	case dx == 0 && dy == -1:
		clearHorizontalSharedEdge(plotA, plotB, plotSize-1, 0, isCentralNorthRoad(a, b))
	}
}

func clearVerticalSharedEdge(a, b *Plot, ax, bx int) {
	for y := 0; y < plotSize; y++ {
		a.Tiles[y][ax].Terrain = terrainEmpty
		b.Tiles[y][bx].Terrain = terrainEmpty
	}
}

func clearHorizontalSharedEdge(a, b *Plot, ay, by int, roadConnector bool) {
	for x := 0; x < plotSize; x++ {
		a.Tiles[ay][x].Terrain = terrainEmpty
		b.Tiles[by][x].Terrain = terrainEmpty
	}
	if roadConnector {
		a.Tiles[ay][homePlotCenter].Terrain = terrainRoad
		b.Tiles[by][homePlotCenter].Terrain = terrainRoad
	}
}

func isCentralNorthRoad(a, b plotCoordinate) bool {
	return a.X == 0 && b.X == 0 && (a.Y >= 0 || b.Y >= 0)
}

func orthogonalPlotNeighbors(coord plotCoordinate) []plotCoordinate {
	return []plotCoordinate{
		{X: coord.X, Y: coord.Y + 1},
		{X: coord.X + 1, Y: coord.Y},
		{X: coord.X, Y: coord.Y - 1},
		{X: coord.X - 1, Y: coord.Y},
	}
}

// exploredPlotCoordinates returns explored Plot coordinates in deterministic draw order.
func (m *Map) exploredPlotCoordinates() []plotCoordinate {
	m.ensurePlots()
	coords := make([]plotCoordinate, 0, len(m.Plots))
	for coord := range m.Plots {
		coords = append(coords, coord)
	}
	sort.Slice(coords, func(i, j int) bool {
		if coords[i].Y != coords[j].Y {
			return coords[i].Y < coords[j].Y
		}
		return coords[i].X < coords[j].X
	})
	return coords
}
