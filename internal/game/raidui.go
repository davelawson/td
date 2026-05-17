package game

import (
	"image/color"
	"math"

	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	nextRaidButtonX      = 42
	nextRaidButtonW      = 190
	nextRaidButtonH      = 52
	nextRaidMarginY      = 42
	raidEnemyRadius      = 14
	raidEnemySpriteSize  = 44
	raidHealthBarHeight  = 5
	raidHealthBarGap     = 4
	projectileSpriteSize = 24
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
		s.drawRaidEnemy(screen, viewport, enemy)
	}
}

// drawRaidEnemy renders one active enemy at its projected world position.
func (s *State) drawRaidEnemy(screen *ebiten.Image, viewport sceneViewport, enemy raidEnemy) {
	var sprite *ebiten.Image
	if enemy.template != nil {
		sprite = enemy.template.Sprite
	}
	if sprite == nil {
		radius := raidEnemyRadius / plotBaseTileSize
		rect := s.projectRect(
			viewport,
			enemy.position.X-radius,
			enemy.position.Y+radius,
			radius*2,
			radius*2,
		)
		vector.FillCircle(screen, rect.x+rect.w/2, rect.y+rect.h/2, rect.w/2, raidEnemyColor, false)
		vector.StrokeCircle(screen, rect.x+rect.w/2, rect.y+rect.h/2, rect.w/2, 2, textColor, false)
		s.drawRaidEnemyHealthBar(screen, rect, enemy)
		return
	}

	size := raidEnemySpriteSize / plotBaseTileSize
	rect := s.projectRect(
		viewport,
		enemy.position.X-size/2,
		enemy.position.Y+size/2,
		size,
		size,
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
	s.drawRaidEnemyHealthBar(screen, rect, enemy)
}

// drawRaidEnemyHealthBar renders a proportional health bar above an enemy.
func (s *State) drawRaidEnemyHealthBar(screen *ebiten.Image, rect projectedRect, enemy raidEnemy) {
	if rect.w <= 0 {
		return
	}

	fraction := raidEnemyHealthFraction(enemy)
	barX := rect.x
	barY := rect.y - raidHealthBarGap - raidHealthBarHeight
	backing := color.RGBA{R: 18, G: 19, B: 17, A: 210}
	vector.FillRect(screen, barX, barY, rect.w, raidHealthBarHeight, backing, false)
	vector.FillRect(screen, barX, barY, rect.w*float32(fraction), raidHealthBarHeight, raidEnemyHealthBarColor(fraction), false)
}

// raidEnemyHealthFraction returns the enemy's current health as a clamped ratio.
func raidEnemyHealthFraction(enemy raidEnemy) float64 {
	maxHealth := raidEnemyMaxHealth(enemy)
	if maxHealth <= 0 {
		return 1
	}
	fraction := float64(enemy.health) / float64(maxHealth)
	if fraction < 0 {
		return 0
	}
	if fraction > 1 {
		return 1
	}
	return fraction
}

// raidEnemyMaxHealth returns the enemy template's maximum health when known.
func raidEnemyMaxHealth(enemy raidEnemy) int {
	if enemy.template == nil {
		return 0
	}
	return enemy.template.MaxHealth
}

// raidEnemyHealthBarColor returns the green-to-red health bar fill color.
func raidEnemyHealthBarColor(fraction float64) color.RGBA {
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}
	return color.RGBA{
		R: uint8(math.Round(255 * (1 - fraction))),
		G: uint8(math.Round(255 * fraction)),
		B: 0,
		A: 255,
	}
}

// drawProjectiles renders active combat projectiles on the camera-projected scene.
func (s *State) drawProjectiles(screen *ebiten.Image) {
	if len(s.combat.projectiles) == 0 {
		return
	}

	viewport := s.sceneViewport()
	for _, projectile := range s.combat.projectiles {
		s.drawProjectile(screen, viewport, projectile)
	}
}

// drawProjectile renders one active combat projectile at its projected world position.
func (s *State) drawProjectile(screen *ebiten.Image, viewport sceneViewport, projectile combatProjectile) {
	size := projectileSpriteSize / plotBaseTileSize
	rect := s.projectRect(
		viewport,
		projectile.position.X-size/2,
		projectile.position.Y+size/2,
		size,
		size,
	)
	if rect.w <= 0 || rect.h <= 0 {
		return
	}
	if projectile.sprite == nil {
		vector.FillCircle(screen, rect.x+rect.w/2, rect.y+rect.h/2, rect.w/3, textColor, false)
		return
	}

	spriteWidth := float64(projectile.sprite.Bounds().Dx())
	spriteHeight := float64(projectile.sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		return
	}

	scale := float64(rect.w) / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
	options.GeoM.Scale(scale, scale)
	options.GeoM.Rotate(s.projectileDrawRotation(projectile))
	options.GeoM.Translate(float64(rect.x)+float64(rect.w)/2, float64(rect.y)+float64(rect.h)/2)
	screen.DrawImage(projectile.sprite, options)
}

// projectileDrawRotation returns a screen-space rotation for projectile sprites.
func (s *State) projectileDrawRotation(projectile combatProjectile) float64 {
	enemyIndex, ok := s.enemyIndexByID(projectile.targetID)
	if !ok {
		return 0
	}
	target := s.raid.enemies[enemyIndex].position
	screenDX := target.X - projectile.position.X
	screenDY := projectile.position.Y - target.Y
	if screenDX == 0 && screenDY == 0 {
		return 0
	}
	const spriteNativeAngle = -math.Pi / 4
	return math.Atan2(screenDY, screenDX) - spriteNativeAngle
}
