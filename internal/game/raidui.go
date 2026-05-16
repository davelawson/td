package game

import (
	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	nextRaidButtonX     = 42
	nextRaidButtonW     = 190
	nextRaidButtonH     = 52
	nextRaidMarginY     = 42
	raidEnemyRadius     = 14
	raidEnemySpriteSize = 44
)

// updateRaidControls applies pointer hover and button clicks for Raid controls.
func (s *State) updateRaidControls(input Input) {
	s.ui.nextRaidHover = s.nextRaidButtonContains(input.CursorX, input.CursorY)
	if input.Clicked && s.ui.nextRaidHover {
		s.startNextRaid()
	}
}

// nextRaidButton returns the screen-space Next Raid button bounds.
func (s *State) nextRaidButton() ui.Button[int] {
	return ui.Button[int]{
		Label:    "Next Raid",
		X:        nextRaidButtonX,
		Y:        s.ui.height - nextRaidMarginY - nextRaidButtonH,
		W:        nextRaidButtonW,
		H:        nextRaidButtonH,
		Disabled: !s.canStartRaid(),
	}
}

// nextRaidButtonContains reports whether a point is inside the Next Raid button.
func (s *State) nextRaidButtonContains(x, y int) bool {
	return s.nextRaidButton().Contains(x, y)
}

// drawRaidControls renders the bottom-left Raid UI button.
func (s *State) drawRaidControls(screen *ebiten.Image) {
	button := s.nextRaidButton()
	fill := clearingColor
	edge := fieldEdgeColor
	labelColor := textColor
	if button.Disabled {
		fill = plotBackdropColor
		edge = tileGridColor
		labelColor = mutedTextColor
	} else if s.ui.nextRaidHover {
		fill = pauseColor
		edge = textColor
	}

	vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
	vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)

	labelWidth, _ := text.Measure(button.Label, s.ui.hudFace, s.ui.hudFace.Size)
	labelX := float64(button.X) + (float64(button.W)-labelWidth)/2
	ui.DrawText(screen, button.Label, s.ui.hudFace, labelX, float64(button.Y+13), labelColor)
}

// drawRaidEnemies renders active Raid enemies on the camera-projected road.
func (s *State) drawRaidEnemies(screen *ebiten.Image) {
	if len(s.raid.enemies) == 0 {
		return
	}

	viewport := s.sceneViewport()
	for _, enemy := range s.raid.enemies {
		worldX, worldY := raidEnemyWorldPosition(enemy)
		s.drawRaidEnemy(screen, viewport, enemy, worldX, worldY)
	}
}

// drawRaidEnemy renders one active enemy at its projected world position.
func (s *State) drawRaidEnemy(screen *ebiten.Image, viewport sceneViewport, enemy raidEnemy, worldX, worldY float64) {
	var sprite *ebiten.Image
	if enemy.template != nil {
		sprite = enemy.template.Sprite
	}
	if sprite == nil {
		rect := s.projectRect(
			viewport,
			worldX-raidEnemyRadius,
			worldY-raidEnemyRadius,
			raidEnemyRadius*2,
			raidEnemyRadius*2,
		)
		vector.FillCircle(screen, rect.x+rect.w/2, rect.y+rect.h/2, rect.w/2, raidEnemyColor, false)
		vector.StrokeCircle(screen, rect.x+rect.w/2, rect.y+rect.h/2, rect.w/2, 2, textColor, false)
		return
	}

	rect := s.projectRect(
		viewport,
		worldX-raidEnemySpriteSize/2,
		worldY-raidEnemySpriteSize/2,
		raidEnemySpriteSize,
		raidEnemySpriteSize,
	)
	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 || rect.w <= 0 || rect.h <= 0 {
		return
	}

	scale := float64(rect.w) / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(rect.x)+(float64(rect.w)-spriteWidth*scale)/2,
		float64(rect.y)+(float64(rect.h)-spriteHeight*scale)/2,
	)
	screen.DrawImage(sprite, options)
}
