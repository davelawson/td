package game

import (
	"image/color"

	"td/internal/ui"
)

var colors = struct {
	background     color.Color
	field          color.Color
	fieldEdge      color.Color
	text           color.Color
	mutedText      color.Color
	pause          color.Color
	fieldAccent    color.Color
	clearing       color.Color
	plotBackdrop   color.Color
	emptyTile      color.Color
	roadTile       color.Color
	forestTile     color.Color
	tileGrid       color.Color
	raidEnemy      color.Color
	topBar         color.Color
	topBarEdge     color.Color
	selectionPanel color.Color
	overlay        color.Color
}{
	background:     ui.CharcoalBlack,
	field:          ui.PineGreen,
	fieldEdge:      ui.Bronze,
	text:           ui.Parchment,
	mutedText:      ui.MutedParchment,
	pause:          ui.LightBronze,
	fieldAccent:    ui.Purple,
	clearing:       ui.MossGreen,
	plotBackdrop:   ui.DarkCharcoalGreen,
	emptyTile:      ui.PineGreen,
	roadTile:       ui.OliveBrown,
	forestTile:     ui.DarkMossGreen,
	tileGrid:       ui.DarkMossGreen,
	raidEnemy:      ui.Purple,
	topBar:         ui.DarkCharcoalGreen,
	topBarEdge:     ui.Bronze,
	selectionPanel: color.RGBA{R: 26, G: 31, B: 24, A: 232},
	overlay:        ui.TransparentBlack,
}
