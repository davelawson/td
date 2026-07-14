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
	treeTile       color.Color
	boulderTile    color.Color
	tileGrid       color.Color
	raidEnemy      color.Color
	topBar         color.Color
	topBarEdge     color.Color
	resourceWood   color.Color
	resourceStone  color.Color
	resourceMetal  color.Color
	buildable      color.Color
	buildBlocked   color.Color
	exploreButton  color.Color
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
	treeTile:       ui.DarkMossGreen,
	boulderTile:    color.RGBA{R: 94, G: 97, B: 91, A: 255},
	tileGrid:       ui.DarkMossGreen,
	raidEnemy:      ui.Purple,
	topBar:         ui.DarkCharcoalGreen,
	topBarEdge:     ui.Bronze,
	resourceWood:   ui.ResourceWood,
	resourceStone:  ui.ResourceStone,
	resourceMetal:  ui.ResourceMetal,
	buildable:      color.RGBA{R: 92, G: 220, B: 104, A: 255},
	buildBlocked:   color.RGBA{R: 224, G: 76, B: 65, A: 255},
	exploreButton:  color.RGBA{R: 218, G: 198, B: 132, A: 245},
	selectionPanel: ui.SelectionPanelBackground,
	overlay:        ui.TransparentBlack,
}
