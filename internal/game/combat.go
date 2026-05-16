package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// gameUpdateSeconds is the current fixed-step duration used to advance real-time combat stats.
const gameUpdateSeconds = 1.0 / 60.0

// combatState stores transient Raid combat state.
type combatState struct {
	towerCooldowns map[tileCoordinate]float64
	projectiles    []combatProjectile
}

// tileCoordinate identifies one Tile in the home Plot.
type tileCoordinate struct {
	X int
	Y int
}

// combatProjectile describes one active projectile tracking an original target.
type combatProjectile struct {
	targetID            int
	position            worldPosition
	damage              int
	speedTilesPerSecond float64
	sprite              *ebiten.Image
}

// resetCombatForRaid clears transient combat state for a newly started Raid.
func (s *State) resetCombatForRaid() {
	s.combat = combatState{
		towerCooldowns: make(map[tileCoordinate]float64),
	}
}

// updateCombat advances tower cooldowns, projectile hits, and tower firing.
func (s *State) updateCombat() {
	s.ensureCombatState()
	s.updateTowerCooldowns(gameUpdateSeconds)
	s.updateProjectiles(gameUpdateSeconds)
	s.fireBowTowers()
}

// ensureCombatState initializes lazy combat maps for tests that construct State values directly.
func (s *State) ensureCombatState() {
	if s.combat.towerCooldowns == nil {
		s.combat.towerCooldowns = make(map[tileCoordinate]float64)
	}
}

// updateTowerCooldowns reduces per-tower fire cooldowns in seconds.
func (s *State) updateTowerCooldowns(deltaSeconds float64) {
	for key, cooldown := range s.combat.towerCooldowns {
		cooldown -= deltaSeconds
		if cooldown < 0 {
			cooldown = 0
		}
		s.combat.towerCooldowns[key] = cooldown
	}
}

// fireBowTowers launches projectiles from ready Bow Towers with targets in range.
func (s *State) fireBowTowers() {
	template := s.structureCatalog.BowTower
	if template.RangeTiles <= 0 || template.Damage <= 0 || template.FireIntervalSeconds <= 0 || template.ProjectileSpeedTilesPerSecond <= 0 {
		return
	}

	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			if s.gameMap.Home.Tiles[y][x].Feature != featureBowTower {
				continue
			}
			key := tileCoordinate{X: x, Y: y}
			if s.combat.towerCooldowns[key] > 0 {
				continue
			}
			towerPosition := tileWorldCenter(x, y)
			target, ok := s.findBowTowerTarget(towerPosition, template.RangeTiles)
			if !ok {
				continue
			}
			s.combat.projectiles = append(s.combat.projectiles, combatProjectile{
				targetID:            target.id,
				position:            towerPosition,
				damage:              template.Damage,
				speedTilesPerSecond: template.ProjectileSpeedTilesPerSecond,
				sprite:              template.ProjectileSprite,
			})
			s.combat.towerCooldowns[key] = template.FireIntervalSeconds
		}
	}
}

// findBowTowerTarget chooses the in-range enemy closest to the Sanctum.
func (s *State) findBowTowerTarget(towerPosition worldPosition, rangeTiles float64) (raidEnemy, bool) {
	rangeSquared := rangeTiles * rangeTiles
	var target raidEnemy
	found := false
	bestSanctumDistance := 0.0
	for _, enemy := range s.raid.enemies {
		if enemy.health <= 0 || distanceSquared(towerPosition, enemy.position) > rangeSquared {
			continue
		}
		sanctumDistance := distanceSquared(worldPosition{}, enemy.position)
		if !found || sanctumDistance < bestSanctumDistance || (sanctumDistance == bestSanctumDistance && enemy.id < target.id) {
			target = enemy
			bestSanctumDistance = sanctumDistance
			found = true
		}
	}
	return target, found
}

// updateProjectiles advances active projectiles and applies damage on impact.
func (s *State) updateProjectiles(deltaSeconds float64) {
	survivors := s.combat.projectiles[:0]
	for _, projectile := range s.combat.projectiles {
		enemyIndex, ok := s.enemyIndexByID(projectile.targetID)
		if !ok {
			continue
		}

		target := s.raid.enemies[enemyIndex].position
		distance := distance(projectile.position, target)
		step := projectile.speedTilesPerSecond * deltaSeconds
		if step >= distance {
			s.damageEnemy(enemyIndex, projectile.damage)
			continue
		}

		projectile.position.X += (target.X - projectile.position.X) / distance * step
		projectile.position.Y += (target.Y - projectile.position.Y) / distance * step
		survivors = append(survivors, projectile)
	}
	s.combat.projectiles = survivors
}

// enemyIndexByID returns the active enemy slice index for an enemy ID.
func (s *State) enemyIndexByID(id int) (int, bool) {
	for i, enemy := range s.raid.enemies {
		if enemy.id == id && enemy.health > 0 {
			return i, true
		}
	}
	return 0, false
}

// damageEnemy applies damage and removes an enemy when its health reaches zero.
func (s *State) damageEnemy(index int, damage int) {
	if index < 0 || index >= len(s.raid.enemies) || damage <= 0 {
		return
	}
	s.raid.enemies[index].health -= damage
	if s.raid.enemies[index].health > 0 {
		return
	}
	copy(s.raid.enemies[index:], s.raid.enemies[index+1:])
	s.raid.enemies = s.raid.enemies[:len(s.raid.enemies)-1]
}

// distance returns the Euclidean distance between two world positions.
func distance(a, b worldPosition) float64 {
	return math.Sqrt(distanceSquared(a, b))
}

// distanceSquared returns the squared Euclidean distance between two world positions.
func distanceSquared(a, b worldPosition) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}
