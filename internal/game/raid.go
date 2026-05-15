package game

const (
	firstRaidEnemyCount = 5
	raidEnemyGrowth     = 2
	raidSpawnInterval   = 45
	raidEnemySpeed      = 3.0
)

type raidState struct {
	active         bool
	breached       bool
	number         int
	pendingEnemies int
	spawnCountdown int
	enemies        []raidEnemy
}

type raidEnemy struct {
	template *EnemyTemplate
	progress float64
}

// startNextRaid begins the next deterministic Raid when the game can accept one.
func (s *State) startNextRaid() {
	if !s.canStartRaid() {
		return
	}

	s.raid.number++
	s.raid.active = true
	s.raid.pendingEnemies = raidEnemyCount(s.raid.number)
	s.raid.enemies = nil
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

	s.raid.enemies = append(s.raid.enemies, raidEnemy{
		template: &s.enemyCatalog.SkeletonSwordShield,
	})
	s.raid.pendingEnemies--
	s.raid.spawnCountdown = raidSpawnInterval
}

// updateRaidEnemies moves active enemies and applies Sanctum contact rules.
func (s *State) updateRaidEnemies() {
	survivors := s.raid.enemies[:0]
	for _, enemy := range s.raid.enemies {
		enemy.progress += raidEnemySpeed
		if enemy.progress >= raidPathLength() {
			if s.applySanctumContact() {
				continue
			}
			return
		}
		survivors = append(survivors, enemy)
	}
	s.raid.enemies = survivors
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
	s.status.phase = phaseCalm
	return false
}

// completeRaid returns the game to calm state after all Raid enemies are gone.
func (s *State) completeRaid() {
	s.raid.active = false
	s.status.phase = phaseCalm
	s.status.day++
	s.status.calmTime = 120
}

// raidEnemiesRemaining returns active and pending enemies in the current Raid.
func (s *State) raidEnemiesRemaining() int {
	return s.raid.pendingEnemies + len(s.raid.enemies)
}

// raidPathLength returns the current straight north-road distance to the Sanctum.
func raidPathLength() float64 {
	return float64(homePlotCenter) * plotBaseTileSize
}

// raidEnemyWorldPosition returns an enemy's current world-space center.
func raidEnemyWorldPosition(enemy raidEnemy) (float64, float64) {
	x := (float64(homePlotCenter) + 0.5) * plotBaseTileSize
	y := 0.5*plotBaseTileSize + enemy.progress
	return x, y
}
