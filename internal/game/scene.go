package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	plotTopPadding    = 150
	plotBottomPadding = 116
	plotMaxTileSize   = 64
	plotMinTileSize   = 24
)

type plotLayout struct {
	x    float32
	y    float32
	tile float32
	size float32
}

// homePlotLayout calculates the centered square used to draw the home Plot.
func (s *State) homePlotLayout() plotLayout {
	availableHeight := s.ui.height - plotTopPadding - plotBottomPadding
	tile := availableHeight / plotSize
	if tile > plotMaxTileSize {
		tile = plotMaxTileSize
	}
	if tile < plotMinTileSize {
		tile = plotMinTileSize
	}
	size := tile * plotSize
	return plotLayout{
		x:    float32(s.ui.width)/2 - float32(size)/2,
		y:    float32(plotTopPadding) + float32(availableHeight-size)/2,
		tile: float32(tile),
		size: float32(size),
	}
}

// drawHomePlot renders the static home Plot from map state.
func (s *State) drawHomePlot(screen *ebiten.Image) {
	layout := s.homePlotLayout()
	vector.FillRect(screen, layout.x-18, layout.y-18, layout.size+36, layout.size+36, plotBackdropColor, false)
	vector.StrokeRect(screen, layout.x-18, layout.y-18, layout.size+36, layout.size+36, 3, fieldEdgeColor, false)

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			s.drawHomePlotTile(screen, layout, x, y, s.gameMap.Home.Tiles[y][x])
		}
	}
}

// drawHomePlotTile renders one Tile in the static home Plot.
func (s *State) drawHomePlotTile(screen *ebiten.Image, layout plotLayout, x, y int, tile Tile) {
	tileX := layout.x + float32(x)*layout.tile
	tileY := layout.y + float32(y)*layout.tile

	tileColor := emptyTileColor
	if tile.Terrain == terrainRoad {
		tileColor = roadTileColor
	}
	vector.FillRect(screen, tileX, tileY, layout.tile, layout.tile, tileColor, false)
	vector.StrokeRect(screen, tileX, tileY, layout.tile, layout.tile, 1, tileGridColor, false)

	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, tileX, tileY, layout.tile)
	}
}

// drawSanctum renders the centered Sanctum feature.
func (s *State) drawSanctum(screen *ebiten.Image, tileX, tileY, tileSize float32) {
	centerX := tileX + tileSize/2
	centerY := tileY + tileSize/2
	radius := tileSize * 0.34

	vector.FillCircle(screen, centerX, centerY, radius, sanctumColor, false)
	vector.StrokeCircle(screen, centerX, centerY, radius, 3, fieldEdgeColor, false)

	label := "S"
	labelWidth, _ := text.Measure(label, s.ui.titleFace, s.ui.titleFace.Size)
	s.drawText(screen, label, s.ui.titleFace, float64(centerX)-labelWidth/2, float64(centerY)-22, textColor)
}
