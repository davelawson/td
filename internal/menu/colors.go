package menu

import (
	"image/color"

	"td/internal/ui"
)

var colors = struct {
	background         color.Color
	backdropBand       color.Color
	panel              color.Color
	panelEdge          color.Color
	text               color.Color
	mutedText          color.Color
	hover              color.Color
	button             color.Color
	disabled           color.Color
	disabledButtonEdge color.Color
	accent             color.Color
	transparentAccent  color.Color
	transparentEdge    color.Color
	nameField          color.Color
}{
	background:         ui.CharcoalBlack,
	backdropBand:       ui.DarkCharcoalGreen,
	panel:              ui.PineGreen,
	panelEdge:          ui.Bronze,
	text:               ui.Parchment,
	mutedText:          ui.MutedParchment,
	hover:              ui.LightBronze,
	button:             ui.MossGreen,
	disabled:           ui.DarkMossGreen,
	disabledButtonEdge: ui.MossGreen,
	accent:             ui.Purple,
	transparentAccent:  ui.TransparentPurple,
	transparentEdge:    ui.TransparentBronze,
	nameField:          ui.DarkCharcoalGreen,
}
