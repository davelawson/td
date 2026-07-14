package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	treeHorizontalFlipMask   uint16 = 0x8000
	selectedSpriteBrightness        = 1.65
	exploreButtonSize               = 0.78
	exploreButtonStroke             = 0.08
)

// drawExploredPlots renders every explored Plot from map state.
func (s *State) drawExploredPlots(screen *ebiten.Image) {
	viewport := s.sceneViewport()
	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		plot, ok := s.gameMap.plot(plotCoord)
		if !ok {
			continue
		}

		for y := 0; y < plotSize; y++ {
			for x := 0; x < plotSize; x++ {
				s.drawPlotTile(screen, viewport, plotCoord, x, y, plot.Tiles[y][x])
			}
		}
	}
	s.drawExploreButtons(screen, viewport)
}

// drawPlotTile renders one Tile in an explored Plot.
func (s *State) drawPlotTile(screen *ebiten.Image, viewport sceneViewport, plot plotCoordinate, x, y int, tile Tile) {
	worldWest, worldNorth, worldW, worldH := plotTileWorldRect(plot, x, y)
	rect := s.projectRect(viewport, worldWest, worldNorth, worldW, worldH)

	tileColor := colors.emptyTile
	switch tile.Terrain {
	case terrainRoad:
		tileColor = colors.roadTile
	case terrainTree:
		tileColor = colors.treeTile
	case terrainBoulder:
		tileColor = colors.boulderTile
	}
	vector.FillRect(screen, rect.x, rect.y, rect.w, rect.h, tileColor, false)
	vector.StrokeRect(screen, rect.x, rect.y, rect.w, rect.h, 1, colors.tileGrid, false)

	if tile.Terrain == terrainTree {
		s.drawPineTree(screen, rect.x, rect.y, rect.w, tile)
	}
	if tile.Terrain == terrainBoulder {
		s.drawBoulder(screen, rect.x, rect.y, rect.w, tile)
	}
	selected := s.selectedStructure(tileCoordinate{Plot: plot, X: x, Y: y})
	if tile.Feature == featureSanctum {
		s.drawSanctum(screen, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureHouse {
		s.drawStructureSprite(screen, s.structureCatalog.House.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureBarracks {
		s.drawStructureSprite(screen, s.structureCatalog.Barracks.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureDorm {
		s.drawStructureSprite(screen, s.structureCatalog.Dorm.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureWoodcutter {
		s.drawStructureSprite(screen, s.structureCatalog.Woodcutter.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureStoneQuarry {
		s.drawStructureSprite(screen, s.structureCatalog.StoneQuarry.Sprite, rect.x, rect.y, rect.w, selected)
	}
	if tile.Feature == featureIronMine {
		s.drawStructureSprite(screen, s.structureCatalog.IronMine.Sprite, rect.x, rect.y, rect.w, selected)
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

// drawExploreButtons renders border reveal controls for unexplored adjacent Plots.
func (s *State) drawExploreButtons(screen *ebiten.Image, viewport sceneViewport) {
	for _, button := range s.exploreButtons() {
		rect := s.projectRect(
			viewport,
			button.Center.X-exploreButtonSize/2,
			button.Center.Y+exploreButtonSize/2,
			exploreButtonSize,
			exploreButtonSize,
		)
		if rect.w <= 0 || rect.h <= 0 {
			continue
		}
		cx := rect.x + rect.w/2
		cy := rect.y + rect.h/2
		radius := rect.w / 2
		vector.FillCircle(screen, cx, cy, radius, colors.plotBackdrop, false)
		vector.StrokeCircle(screen, cx, cy, radius, 3, colors.exploreButton, false)

		lensRadius := rect.w * 0.19
		vector.StrokeCircle(screen, cx-rect.w*0.08, cy-rect.h*0.08, lensRadius, 3, colors.exploreButton, false)
		vector.StrokeLine(
			screen,
			cx+rect.w*0.06,
			cy+rect.h*0.06,
			cx+rect.w*0.22,
			cy+rect.h*0.22,
			float32(exploreButtonStroke*float64(rect.w)),
			colors.exploreButton,
			false,
		)
	}
}

// drawPineTree renders a tree sprite variant chosen from the Tile tweak.
func (s *State) drawPineTree(screen *ebiten.Image, tileX, tileY, tileSize float32, tile Tile) {
	trees := s.assetCatalog.Sprite.Terrain.PineTrees
	tree := trees[terrainSpriteIndex(tile.Tweak, len(trees))]
	if tree == nil || tileSize <= 0 {
		return
	}

	targetSize := float64(tileSize) * 0.92
	drawTerrainSprite(screen, tree, tileX, tileY, tileSize, targetSize, tile.Tweak)
}

// drawBoulder renders a Boulder sprite variant chosen from the Tile tweak.
func (s *State) drawBoulder(screen *ebiten.Image, tileX, tileY, tileSize float32, tile Tile) {
	boulders := s.assetCatalog.Sprite.Terrain.Boulders
	boulder := boulders[terrainSpriteIndex(tile.Tweak, len(boulders))]
	if boulder == nil || tileSize <= 0 {
		return
	}

	targetSize := float64(tileSize) * 0.78
	drawTerrainSprite(screen, boulder, tileX, tileY, tileSize, targetSize, tile.Tweak)
}

// drawTerrainSprite renders one terrain sprite with optional tweak-driven mirroring.
func drawTerrainSprite(screen *ebiten.Image, sprite *ebiten.Image, tileX, tileY, tileSize float32, targetSize float64, tweak uint16) {
	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	targetWidth := spriteWidth * scale
	if terrainSpriteFlipped(tweak) {
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
	screen.DrawImage(sprite, options)
}

// terrainSpriteIndex returns the terrain sprite variant selected by the Tile tweak.
func terrainSpriteIndex(tweak uint16, variants int) int {
	if variants <= 0 {
		return 0
	}
	return int(tweak&^treeHorizontalFlipMask) % variants
}

// terrainSpriteFlipped reports whether the Tile tweak requests horizontal sprite flipping.
func terrainSpriteFlipped(tweak uint16) bool {
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
