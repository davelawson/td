package ui

// Button describes a rectangular UI target and the action it selects.
type Button[T comparable] struct {
	Label    string
	X        int
	Y        int
	W        int
	H        int
	Action   T
	Disabled bool
}

// Contains reports whether the point is inside the button bounds.
func (b Button[T]) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.W && y >= b.Y && y < b.Y+b.H
}
