package game

const (
	firstRaidEnemyCount = 5
	raidEnemyGrowth     = 2
	raidSpawnInterval   = 45
)

type raidState struct {
	active         bool
	breached       bool
	number         int
	pendingEnemies int
	spawnCountdown int
	nextEnemyID    int
	enemies        []raidEnemy
}

type raidEnemy struct {
	id       int
	template *EnemyTemplate
	position coord
	health   int
}

// startNextRaid begins the next deterministic Raid when the game can accept one.
func (s *State) startNextRaid() {
	if !s.canStartRaid() {
		return
	}

	s.raid.number++
	s.raid.active = true
	s.raid.pendingEnemies = raidEnemyCount(s.raid.number)
	s.raid.nextEnemyID = 0
	s.raid.enemies = nil
	s.resetCombatForRaid()
	s.spawnRaidEnemy()
	s.status.phase = phaseRaid
}

// canStartRaid reports whether the player can start another Raid.
func (s *State) canStartRaid() bool {
	return !s.paused && !s.raid.active && !s.raid.breached
}

// raidEnemyCount returns the scripted enemy count for a Raid number.
func raidEnemyCount(number int) int {
	if number < 1 {
		return firstRaidEnemyCount
	}
	return firstRaidEnemyCount + (number-1)*raidEnemyGrowth
}

// updateRaid advances spawning, movement, Sanctum contact, and completion.
func (s *State) updateRaid() {
	if !s.raid.active {
		return
	}

	s.updateRaidSpawning()
	s.updateCombat()
	s.updateRaidEnemies()
	if s.raid.active && s.raid.pendingEnemies == 0 && len(s.raid.enemies) == 0 {
		s.completeRaid()
	}
}

// updateRaidSpawning advances the staggered enemy spawn timer.
func (s *State) updateRaidSpawning() {
	if s.raid.pendingEnemies == 0 {
		return
	}

	s.raid.spawnCountdown--
	if s.raid.spawnCountdown <= 0 {
		s.spawnRaidEnemy()
	}
}

// spawnRaidEnemy adds one enemy at the road edge and resets the spawn timer.
func (s *State) spawnRaidEnemy() {
	if s.raid.pendingEnemies == 0 {
		return
	}

	template := s.nextRaidEnemyTemplate()
	s.raid.enemies = append(s.raid.enemies, raidEnemy{
		id:       s.raid.nextEnemyID,
		template: template,
		position: raidEnemySpawnPosition(),
		health:   template.MaxHealth,
	})
	s.raid.nextEnemyID++
	s.raid.pendingEnemies--
	s.raid.spawnCountdown = raidSpawnInterval
}

// nextRaidEnemyTemplate returns the template for the next deterministic spawn.
func (s *State) nextRaidEnemyTemplate() *EnemyTemplate {
	if s.raid.number == 1 && s.raid.nextEnemyID%2 == 1 {
		return &s.enemyCatalog.Zombie
	}
	return &s.enemyCatalog.SkeletonSwordShield
}

// updateRaidEnemies moves active enemies and applies Sanctum contact rules.
func (s *State) updateRaidEnemies() {
	survivors := s.raid.enemies[:0]
	for _, enemy := range s.raid.enemies {
		enemy.position.Y -= raidEnemyMovementStep(enemy, gameUpdateSeconds)
		if raidEnemyReachedSanctum(enemy) {
			if s.applySanctumContact() {
				continue
			}
			return
		}
		survivors = append(survivors, enemy)
	}
	s.raid.enemies = survivors
}

// raidEnemyMovementStep returns the enemy's movement distance in Tiles.
func raidEnemyMovementStep(enemy raidEnemy, deltaSeconds float64) float64 {
	if enemy.template == nil || enemy.template.SpeedTilesPerSecond <= 0 || deltaSeconds <= 0 {
		return 0
	}
	return enemy.template.SpeedTilesPerSecond * deltaSeconds
}

// applySanctumContact removes a reaching enemy or breaches the Sanctum.
func (s *State) applySanctumContact() bool {
	if s.status.barricade > 0 {
		s.status.barricade--
		return true
	}

	s.raid.breached = true
	s.raid.active = false
	s.raid.pendingEnemies = 0
	s.raid.enemies = nil
	s.combat.projectiles = nil
	s.status.phase = phaseCalm
	return false
}

// completeRaid returns the game to calm state after all Raid enemies are gone.
func (s *State) completeRaid() {
	s.grantEconomicBuildingResources()
	s.raid.active = false
	s.status.phase = phaseCalm
	s.status.day++
	s.status.calmTime = 120
}

// raidEnemiesRemaining returns active and pending enemies in the current Raid.
func (s *State) raidEnemiesRemaining() int {
	return s.raid.pendingEnemies + len(s.raid.enemies)
}

// raidEnemySpawnPosition returns the current north-road enemy spawn point.
func raidEnemySpawnPosition() coord {
	return coord{X: 0, Y: float64(homePlotCenter)}
}

// raidEnemyReachedSanctum reports whether the enemy has contacted the Sanctum.
func raidEnemyReachedSanctum(enemy raidEnemy) bool {
	return enemy.position.Y <= 0
}
