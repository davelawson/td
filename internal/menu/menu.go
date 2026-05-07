package menu

// Action identifies the behavior selected from the main menu.
type Action int

const (
	// ActionNone means no menu action was selected.
	ActionNone Action = iota
	// ActionQuit means the app should terminate cleanly.
	ActionQuit
)

// Button describes a rectangular menu target and the action it selects.
type Button struct {
	Label  string
	X      int
	Y      int
	W      int
	H      int
	Action Action
}

// Contains reports whether the point is inside the button bounds.
func (b Button) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.W && y >= b.Y && y < b.Y+b.H
}

// ActionAt returns the first button action containing the point.
func ActionAt(buttons []Button, x, y int) Action {
	for _, button := range buttons {
		if button.Contains(x, y) {
			return button.Action
		}
	}
	return ActionNone
}
