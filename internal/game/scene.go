package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	treeHorizontalFlipMask   uint16 = 0x8000
	selectedSpriteBrightness        = 1.65
)

// drawHomePlot renders the static home Plot from map state.
func (s *State) drawHomePlot(screen *ebiten.Image) {
	viewport := s.sceneViewport()
	margin := 18 / plotBaseTileSize
	size := float64(plotSize)
	backdrop := s.projectRect(viewport, -size/2-margin, size/2+margin, size+margin*2, size+margin*2)
	vector.FillRect(screen, backdrop.x, backdrop.y, backdrop.w, backdrop.h, colors.plotBackdrop, false)
	vector.StrokeRect(screen, backdrop.x, backdrop.y, backdrop.w, backdrop.h, 3, colors.fieldEdge, false)

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

	tileColor := colors.emptyTile
	switch tile.Terrain {
	case terrainRoad:
		tileColor = colors.roadTile
	case terrainForest:
		tileColor = colors.forestTile
	}
	vector.FillRect(screen, rect.x, rect.y, rect.w, rect.h, tileColor, false)
	vector.StrokeRect(screen, rect.x, rect.y, rect.w, rect.h, 1, colors.tileGrid, false)

	if tile.Terrain == terrainForest {
		s.drawPineTree(screen, rect.x, rect.y, rect.w, tile)
	}
	selected := s.selectedStructure(x, y)
	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureHouse {
		s.drawStructureSprite(screen, s.structureCatalog.House.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureBowTower {
		s.drawStructureSprite(screen, s.structureCatalog.BowTower.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureFlameBoltTower {
		s.drawStructureSprite(screen, s.structureCatalog.FlameBoltTower.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureCatapultTower {
		s.drawStructureSprite(screen, s.structureCatalog.CatapultTower.Sprite, rect.x, rect.y, rect.w, selected)
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
func (s *State) drawSanctum(screen *ebiten.Image, tileX, tileY, tileSize float32, selected bool) {
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
	if selected {
		brightenDrawOptions(options)
	}
	screen.DrawImage(sanctum, options)
}

// drawStructureSprite renders a centered authored structure feature.
func (s *State) drawStructureSprite(screen *ebiten.Image, sprite *ebiten.Image, tileX, tileY, tileSize float32, selected bool) {
	if sprite == nil || tileSize <= 0 {
		return
	}

	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())
	targetSize := float64(tileSize) * 0.76
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(tileX)+(float64(tileSize)-spriteWidth*scale)/2,
		float64(tileY)+(float64(tileSize)-spriteHeight*scale)/2,
	)
	if selected {
		brightenDrawOptions(options)
	}
	screen.DrawImage(sprite, options)
}

// brightenDrawOptions applies the selected-object brightness treatment to a sprite draw.
func brightenDrawOptions(options *ebiten.DrawImageOptions) {
	options.ColorScale.Scale(selectedSpriteBrightness, selectedSpriteBrightness, selectedSpriteBrightness, 1)
}
