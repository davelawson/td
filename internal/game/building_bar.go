package game

import (
	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	buildingBarWidth       = 96
	buildingBarPadding     = 16
	buildingBarItemSize    = 64
	buildingBarItemGap     = 14
	buildingBarSpriteInset = 8
)

type buildingBarItem struct {
	Name   string
	Sprite *ebiten.Image
	Bounds ui.Button[int]
}

// buildingBarBounds returns the screen-space building bar rectangle.
func (s *State) buildingBarBounds() ui.Button[int] {
	return ui.Button[int]{
		X: 0,
		Y: topBarHeight,
		W: buildingBarWidth,
		H: s.ui.height - topBarHeight,
	}
}

// buildingBarItems returns the visual-only tower choices shown in the building bar.
func (s *State) buildingBarItems() []buildingBarItem {
	bar := s.buildingBarBounds()
	x := bar.X + (bar.W-buildingBarItemSize)/2
	y := bar.Y + buildingBarPadding
	return []buildingBarItem{
		{
			Name:   s.structureCatalog.BowTower.Name,
			Sprite: s.structureCatalog.BowTower.Sprite,
			Bounds: ui.Button[int]{
				Label:  s.structureCatalog.BowTower.Name,
				X:      x,
				Y:      y,
				W:      buildingBarItemSize,
				H:      buildingBarItemSize,
				Action: 0,
			},
		},
		{
			Name:   s.structureCatalog.FlameBoltTower.Name,
			Sprite: s.structureCatalog.FlameBoltTower.Sprite,
			Bounds: ui.Button[int]{
				Label:  s.structureCatalog.FlameBoltTower.Name,
				X:      x,
				Y:      y + buildingBarItemSize + buildingBarItemGap,
				W:      buildingBarItemSize,
				H:      buildingBarItemSize,
				Action: 1,
			},
		},
	}
}

// buildingBarContains reports whether a point is inside the visual building bar.
func (s *State) buildingBarContains(x, y int) bool {
	return s.buildingBarBounds().Contains(x, y)
}

// drawBuildingBar renders the visual-only tower picker at the left edge of the scene.
func (s *State) drawBuildingBar(screen *ebiten.Image) {
	bar := s.buildingBarBounds()
	if bar.H <= 0 {
		return
	}

	vector.FillRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), colors.selectionPanel, false)
	vector.StrokeLine(screen, float32(bar.X+bar.W-2), float32(bar.Y), float32(bar.X+bar.W-2), float32(bar.Y+bar.H), 3, colors.fieldEdge, false)

	for _, item := range s.buildingBarItems() {
		s.drawBuildingBarItem(screen, item)
	}
}

// drawBuildingBarItem renders one tower icon slot.
func (s *State) drawBuildingBarItem(screen *ebiten.Image, item buildingBarItem) {
	bounds := item.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), colors.plotBackdrop, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, colors.fieldEdge, false)

	if item.Sprite == nil {
		return
	}
	spriteWidth := float64(item.Sprite.Bounds().Dx())
	spriteHeight := float64(item.Sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		return
	}

	targetSize := float64(bounds.W - buildingBarSpriteInset*2)
	scale := targetSize / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(bounds.X)+(float64(bounds.W)-spriteWidth*scale)/2,
		float64(bounds.Y)+(float64(bounds.H)-spriteHeight*scale)/2,
	)
	screen.DrawImage(item.Sprite, options)
}
