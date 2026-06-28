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
	targetID                     int
	position                     coord
	damage                       int
	speedTilesPerSecond          float64
	sprite                       *ebiten.Image
	damageAllEnemiesInTargetTile bool
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
	s.fireCombatTowers()
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

// fireCombatTowers launches projectiles from ready combat towers with targets in range.
func (s *State) fireCombatTowers() {
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			template, ok := s.combatTowerTemplate(s.gameMap.Home.Tiles[y][x].Feature)
			if !ok || !template.canFireProjectiles() {
				continue
			}
			key := tileCoordinate{X: x, Y: y}
			if s.combat.towerCooldowns[key] > 0 {
				continue
			}
			towerPosition := tileWorldCenter(x, y)
			target, ok := s.findTowerTarget(towerPosition, template.RangeTiles)
			if !ok {
				continue
			}
			s.combat.projectiles = append(s.combat.projectiles, combatProjectile{
				targetID:                     target.id,
				position:                     towerPosition,
				damage:                       template.Damage,
				speedTilesPerSecond:          template.ProjectileSpeedTilesPerSecond,
				sprite:                       template.ProjectileSprite,
				damageAllEnemiesInTargetTile: template.DamageAllEnemiesInTargetTile,
			})
			s.combat.towerCooldowns[key] = template.FireIntervalSeconds
		}
	}
}

// combatTowerTemplate returns the projectile-firing template for a placed feature.
func (s *State) combatTowerTemplate(feature tileFeature) (StructureTemplate, bool) {
	switch feature {
	case featureBowTower:
		return s.structureCatalog.BowTower, true
	case featureFlameBoltTower:
		return s.structureCatalog.FlameBoltTower, true
	case featureCatapultTower:
		return s.structureCatalog.CatapultTower, true
	default:
		return StructureTemplate{}, false
	}
}

// canFireProjectiles reports whether a structure template has complete projectile combat stats.
func (t StructureTemplate) canFireProjectiles() bool {
	return t.RangeTiles > 0 && t.Damage > 0 && t.FireIntervalSeconds > 0 && t.ProjectileSpeedTilesPerSecond > 0
}

// findTowerTarget chooses the in-range enemy closest to the Sanctum.
func (s *State) findTowerTarget(towerPosition coord, rangeTiles float64) (raidEnemy, bool) {
	rangeSquared := rangeTiles * rangeTiles
	var target raidEnemy
	found := false
	bestSanctumDistance := 0.0
	for _, enemy := range s.raid.enemies {
		if enemy.health <= 0 || distanceSquared(towerPosition, enemy.position) > rangeSquared {
			continue
		}
		sanctumDistance := distanceSquared(coord{}, enemy.position)
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
			if projectile.damageAllEnemiesInTargetTile {
				s.damageEnemiesInTargetTile(target, projectile.damage)
			} else {
				s.damageEnemy(enemyIndex, projectile.damage)
			}
			continue
		}

		projectile.position.X += (target.X - projectile.position.X) / distance * step
		projectile.position.Y += (target.Y - projectile.position.Y) / distance * step
		survivors = append(survivors, projectile)
	}
	s.combat.projectiles = survivors
}

// damageEnemiesInTargetTile applies damage to every living enemy in the target's current Tile.
func (s *State) damageEnemiesInTargetTile(targetPosition coord, damage int) {
	targetTile, ok := tileAtWorldPosition(targetPosition)
	if !ok {
		return
	}
	for i := 0; i < len(s.raid.enemies); {
		if enemyTile, ok := tileAtWorldPosition(s.raid.enemies[i].position); ok && enemyTile == targetTile {
			previousCount := len(s.raid.enemies)
			s.damageEnemy(i, damage)
			if len(s.raid.enemies) < previousCount {
				continue
			}
		}
		i++
	}
}

// tileAtWorldPosition returns the home Plot Tile containing a world position.
func tileAtWorldPosition(position coord) (tileCoordinate, bool) {
	x := int(math.Floor(position.X + homePlotCenter + 0.5))
	y := int(math.Floor(homePlotCenter - position.Y + 0.5))
	if x < 0 || y < 0 || x >= plotSize || y >= plotSize {
		return tileCoordinate{}, false
	}
	return tileCoordinate{X: x, Y: y}, true
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
	s.grantEnemyResources(s.raid.enemies[index])
	s.playRaiderDefeatedSound()
	copy(s.raid.enemies[index:], s.raid.enemies[index+1:])
	s.raid.enemies = s.raid.enemies[:len(s.raid.enemies)-1]
}

// distance returns the Euclidean distance between two world positions.
func distance(a, b coord) float64 {
	return math.Sqrt(distanceSquared(a, b))
}

// distanceSquared returns the squared Euclidean distance between two world positions.
func distanceSquared(a, b coord) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}
