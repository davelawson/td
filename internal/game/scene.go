package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const treeHorizontalFlipMask uint16 = 0x8000

// drawHomePlot renders the static home Plot from map state.
func (s *State) drawHomePlot(screen *ebiten.Image) {
	viewport := s.sceneViewport()
	margin := 18 / plotBaseTileSize
	size := float64(plotSize)
	backdrop := s.projectRect(viewport, -size/2-margin, size/2+margin, size+margin*2, size+margin*2)
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
	worldWest, worldNorth, worldW, worldH := tileWorldRect(x, y)
	rect := s.projectRect(viewport, worldWest, worldNorth, worldW, worldH)

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
		s.drawPineTree(screen, rect.x, rect.y, rect.w, tile)
	}
	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, rect.x, rect.y, rect.w)
	}
	if tile.Feature == featureBowTower {
		s.drawBowTower(screen, rect.x, rect.y, rect.w)
	}
}

// drawPineTree renders a tree sprite variant chosen from the Tile tweak.
func (s *State) drawPineTree(screen *ebiten.Image, tileX, tileY, tileSize float32, tile Tile) {
	trees := s.assetCatalog.Sprite.Terrain.PineTrees
	tree := trees[pineTreeSpriteIndex(tile.Tweak, len(trees))]
	if tree == nil || tileSize <= 0 {
		return
	}

	spriteWidth := float64(tree.Bounds().Dx())
	spriteHeight := float64(tree.Bounds().Dy())
	targetSize := float64(tileSize) * 0.92
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	targetWidth := spriteWidth * scale
	if treeSpriteFlipped(tile.Tweak) {
		options.GeoM.Scale(-scale, scale)
		options.GeoM.Translate(
			float64(tileX)+(float64(tileSize)+targetWidth)/2,
			float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
		)
	} else {
		options.GeoM.Scale(scale, scale)
		options.GeoM.Translate(
			float64(tileX)+(float64(tileSize)-targetWidth)/2,
			float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
		)
	}
	screen.DrawImage(tree, options)
}

// pineTreeSpriteIndex returns the terrain sprite variant selected by the Tile tweak.
func pineTreeSpriteIndex(tweak uint16, variants int) int {
	if variants <= 0 {
		return 0
	}
	return int(tweak&^treeHorizontalFlipMask) % variants
}

// treeSpriteFlipped reports whether the Tile tweak requests horizontal tree flipping.
func treeSpriteFlipped(tweak uint16) bool {
	return tweak&treeHorizontalFlipMask != 0
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

// drawBowTower renders the authored Bow Tower feature.
func (s *State) drawBowTower(screen *ebiten.Image, tileX, tileY, tileSize float32) {
	bowTower := s.structureCatalog.BowTower.Sprite
	if bowTower == nil || tileSize <= 0 {
		return
	}

	spriteWidth := float64(bowTower.Bounds().Dx())
	spriteHeight := float64(bowTower.Bounds().Dy())
	targetSize := float64(tileSize) * 0.76
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(tileX)+(float64(tileSize)-spriteWidth*scale)/2,
		float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
	)
	screen.DrawImage(bowTower, options)
}
