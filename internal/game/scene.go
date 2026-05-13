package game

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	switch tile.Terrain {
	case terrainRoad:
		tileColor = roadTileColor
	case terrainForest:
		tileColor = forestTileColor
	}
	vector.FillRect(screen, rect.x, rect.y, rect.w, rect.h, tileColor, false)
	vector.StrokeRect(screen, rect.x, rect.y, rect.w, rect.h, 1, tileGridColor, false)

	if tile.Terrain == terrainForest {
		s.drawPineTree(screen, rect.x, rect.y, rect.w, x, y)
	}
	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, rect.x, rect.y, rect.w)
	}
}

// drawPineTree renders a deterministic tree sprite for a forest Tile.
func (s *State) drawPineTree(screen *ebiten.Image, tileX, tileY, tileSize float32, plotX, plotY int) {
	trees := s.assetCatalog.Sprite.Terrain.PineTrees
	tree := trees[(plotX*3+plotY*5)%len(trees)]
	if tree == nil || tileSize <= 0 {
		return
	}

	spriteWidth := float64(tree.Bounds().Dx())
	spriteHeight := float64(tree.Bounds().Dy())
	targetSize := float64(tileSize) * 0.92
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(tileX)+(float64(tileSize)-spriteWidth*scale)/2,
		float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
	)
	screen.DrawImage(tree, options)
}

// drawSanctum renders the centered Sanctum feature.
func (s *State) drawSanctum(screen *ebiten.Image, tileX, tileY, tileSize float32) {
	sanctum := s.assetCatalog.Sprite.Structure.Sanctum
	if sanctum == nil || tileSize <= 0 {
		return
	}

	spriteWidth := float64(sanctum.Bounds().Dx())
	spriteHeight := float64(sanctum.Bounds().Dy())
	targetSize := float64(tileSize) * 0.82
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(tileX)+(float64(tileSize)-spriteWidth*scale)/2,
		float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
	)
	screen.DrawImage(sanctum, options)
}
