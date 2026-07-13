package game

// coord is a Sanctum-centered position in Tile units, with positive Y north.
type coord struct {
	X float64
	Y float64
}

// tileCoordinate identifies one Tile in an explored Plot.
type tileCoordinate struct {
	Plot plotCoordinate
	X    int
	Y    int
}

// homeTileCoordinate identifies one Tile in the home Plot.
func homeTileCoordinate(x, y int) tileCoordinate {
	return tileCoordinate{X: x, Y: y}
}

// tileWorldCenter returns the world-space center of a home Plot Tile.
func tileWorldCenter(x, y int) coord {
	return plotTileWorldCenter(homePlotCoordinate, x, y)
}

// plotTileWorldCenter returns the world-space center of a Plot Tile.
func plotTileWorldCenter(plot plotCoordinate, x, y int) coord {
	return coord{
		X: float64(plot.X*plotSize + x - homePlotCenter),
		Y: float64(plot.Y*plotSize + homePlotCenter - y),
	}
}

// tileWorldRect returns the west edge, north edge, width, and height of a home Plot Tile.
func tileWorldRect(x, y int) (float64, float64, float64, float64) {
	return plotTileWorldRect(homePlotCoordinate, x, y)
}

// plotTileWorldRect returns the west edge, north edge, width, and height of a Tile.
func plotTileWorldRect(plot plotCoordinate, x, y int) (float64, float64, float64, float64) {
	center := plotTileWorldCenter(plot, x, y)
	return center.X - 0.5, center.Y + 0.5, 1, 1
}

// plotWorldRect returns the west edge, north edge, width, and height of a Plot.
func plotWorldRect(plot plotCoordinate) (float64, float64, float64, float64) {
	west := float64(plot.X*plotSize) - float64(homePlotCenter) - 0.5
	north := float64(plot.Y*plotSize) + float64(homePlotCenter) + 0.5
	return west, north, plotSize, plotSize
}
