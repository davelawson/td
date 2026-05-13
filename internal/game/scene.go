package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawHomePlot renders the static home Plot from map state.
func (s *State) drawHomePlot(screen *ebiten.Image) {
	viewport := s.sceneViewport()
	size := float64(plotSize) * plotBaseTileSize
	backdrop := s.projectRect(viewport, -18, -18, size+36, size+36)
	vector.FillRect(screen, backdrop.x, backdrop.y, backdrop.w, backdrop.h, plotBackdropColor, false)
	vector.StrokeRect(screen, backdrop.x, backdrop.y, backdrop.w, backdrop.h, 3, fieldEdgeColor, false)

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			s.drawHomePlotTile(screen, viewport, x, y, s.gameMap.Home.Tiles[y][x])
		}
	}
}

// drawHomePlotTile renders one Tile in the static home Plot.
func (s *State) drawHomePlotTile(screen *ebiten.Image, viewport sceneViewport, x, y int, tile Tile) {
	worldX := float64(x) * plotBaseTileSize
	worldY := float64(y) * plotBaseTileSize
	rect := s.projectRect(viewport, worldX, worldY, plotBaseTileSize, plotBaseTileSize)

	tileColor := emptyTileColor
	if tile.Terrain == terrainRoad {
		tileColor = roadTileColor
	}
	vector.FillRect(screen, rect.x, rect.y, rect.w, rect.h, tileColor, false)
	vector.StrokeRect(screen, rect.x, rect.y, rect.w, rect.h, 1, tileGridColor, false)

	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, rect.x, rect.y, rect.w)
	}
}

// drawSanctum renders the centered Sanctum feature.
func (s *State) drawSanctum(screen *ebiten.Image, tileX, tileY, tileSize float32) {
	centerX := tileX + tileSize/2
	centerY := tileY + tileSize/2
	radius := tileSize * 0.34

	vector.FillCircle(screen, centerX, centerY, radius, sanctumColor, false)
	vector.StrokeCircle(screen, centerX, centerY, radius, 3, fieldEdgeColor, false)

	if tileSize < 18 {
		return
	}
	label := "S"
	labelWidth, _ := text.Measure(label, s.ui.titleFace, s.ui.titleFace.Size)
	labelScale := float64(tileSize) / plotBaseTileSize
	options := &text.DrawOptions{}
	options.GeoM.Scale(labelScale, labelScale)
	options.GeoM.Translate(float64(centerX)-labelWidth*labelScale/2, float64(centerY)-22*labelScale)
	options.ColorScale.ScaleWithColor(textColor)
	text.Draw(screen, label, s.ui.titleFace, options)
}
