package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// DrawText draws one line at the given coordinates.
func DrawText(screen *ebiten.Image, value string, face *text.GoTextFace, x, y float64, clr color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}

// DrawCenteredText draws one line centered within a container width.
func DrawCenteredText(screen *ebiten.Image, containerWidth int, value string, face *text.GoTextFace, y float64, clr color.Color) {
	width, _ := text.Measure(value, face, face.Size)
	DrawText(screen, value, face, (float64(containerWidth)-width)/2, y, clr)
}

// DrawCenteredTextInRect draws one line centered horizontally inside a rectangle.
func DrawCenteredTextInRect(screen *ebiten.Image, value string, face *text.GoTextFace, rectX, rectY, rectW int, textYOffset float64, clr color.Color) {
	width, _ := text.Measure(value, face, face.Size)
	x := float64(rectX) + (float64(rectW)-width)/2
	DrawText(screen, value, face, x, float64(rectY)+textYOffset, clr)
}
