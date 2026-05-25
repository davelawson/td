package game

// coord is a Sanctum-centered position in Tile units, with positive Y north.
type coord struct {
	X float64
	Y float64
}

// tileWorldCenter returns the world-space center of a Plot Tile.
func tileWorldCenter(x, y int) coord {
	return coord{
		X: float64(x - homePlotCenter),
		Y: float64(homePlotCenter - y),
	}
}

// tileWorldRect returns the west edge, north edge, width, and height of a Tile.
func tileWorldRect(x, y int) (float64, float64, float64, float64) {
	center := tileWorldCenter(x, y)
	return center.X - 0.5, center.Y + 0.5, 1, 1
}
