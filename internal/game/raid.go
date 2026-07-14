package game

type raidState struct {
	active         bool
	breached       bool
	number         int
	template       raidTemplate
	progress       float64
	pendingEnemies int
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

	nextRaidNumber := s.raid.number + 1
	template := generateRaid(nextRaidNumber, s.settlementPopulation(), len(s.gameMap.exploredPlotCoordinates()))
	s.raid.number = nextRaidNumber
	s.raid.active = true
	s.raid.template = template
	s.raid.progress = 0
	s.raid.pendingEnemies = template.totalEnemies()
	s.raid.nextEnemyID = 0
	s.raid.enemies = nil
	s.resetCombatForRaid()
	s.status.phase = phaseRaid
}

// canStartRaid reports whether the player can start another Raid.
func (s *State) canStartRaid() bool {
	return s.status.phase == phaseManagement && !s.paused && !s.raid.active && !s.raid.breached
}

// updateRaid advances spawning, movement, Sanctum contact, and completion.
func (s *State) updateRaid() {
	if !s.raid.active {
		return
	}

	s.updateRaidProgress(gameUpdateSeconds)
	s.updateCombat()
	s.updateRaidEnemies()
	if s.raid.active && s.raid.progress >= 1 && s.raid.pendingEnemies == 0 && len(s.raid.enemies) == 0 {
		s.completeRaid()
	}
}

// updateRaidProgress advances the generated schedule and spawns newly reached enemies.
func (s *State) updateRaidProgress(deltaSeconds float64) {
	if !s.raid.active || s.raid.progress >= 1 || deltaSeconds <= 0 || s.raid.template.progressDurationSeconds <= 0 {
		return
	}

	previousProgress := s.raid.progress
	s.raid.progress += deltaSeconds / s.raid.template.progressDurationSeconds
	if s.raid.progress > 1 {
		s.raid.progress = 1
	}
	previousScore := previousProgress * s.raid.template.challengeRating
	currentScore := s.raid.progress * s.raid.template.challengeRating
	for _, rule := range s.raid.template.enemyRules {
		newEnemies := scheduledEnemyCount(currentScore, rule) - scheduledEnemyCount(previousScore, rule)
		for i := 0; i < newEnemies; i++ {
			s.spawnRaidEnemy(rule.kind)
		}
	}
}

// spawnRaidEnemy adds one generated enemy at the current road edge.
func (s *State) spawnRaidEnemy(kind raidEnemyKind) {
	if s.raid.pendingEnemies == 0 {
		return
	}

	template, ok := s.enemyTemplateForRaidKind(kind)
	if !ok {
		return
	}
	s.raid.enemies = append(s.raid.enemies, raidEnemy{
		id:       s.raid.nextEnemyID,
		template: template,
		position: s.raidEnemySpawnPosition(),
		health:   template.MaxHealth,
	})
	s.raid.nextEnemyID++
	s.raid.pendingEnemies--
}

// enemyTemplateForRaidKind maps a generated enemy kind to the active catalog.
func (s *State) enemyTemplateForRaidKind(kind raidEnemyKind) (*EnemyTemplate, bool) {
	switch kind {
	case raidEnemySkeletonSwordShield:
		return &s.enemyCatalog.SkeletonSwordShield, true
	case raidEnemyZombie:
		return &s.enemyCatalog.Zombie, true
	case raidEnemyGhoul:
		return &s.enemyCatalog.Ghoul, true
	case raidEnemyArmouredSkeleton:
		return &s.enemyCatalog.ArmouredSkeleton, true
	default:
		return nil, false
	}
}

// settlementPopulation returns total inhabitants across every settlement role.
func (s *State) settlementPopulation() int {
	return s.status.populations.apprentices.total +
		s.status.populations.soldiers.total +
		s.status.populations.peasants.total
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
	return false
}

// completeRaid begins the next Day after all Raid enemies are gone.
func (s *State) completeRaid() {
	s.raid.active = false
	s.beginPostRaidDay()
}

// raidEnemiesRemaining returns active and pending enemies in the current Raid.
func (s *State) raidEnemiesRemaining() int {
	return s.raid.pendingEnemies + len(s.raid.enemies)
}

// raidEnemySpawnPosition returns the current north-road enemy spawn point for tests.
func raidEnemySpawnPosition() coord {
	return coord{X: 0, Y: float64(homePlotCenter)}
}

// raidEnemySpawnPosition returns the north-road enemy spawn point for explored central north Plots.
func (s *State) raidEnemySpawnPosition() coord {
	northPlotY := 0
	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		if plotCoord.X == 0 && plotCoord.Y > northPlotY {
			northPlotY = plotCoord.Y
		}
	}
	return coord{X: 0, Y: float64(northPlotY*plotSize + homePlotCenter)}
}

// raidEnemyReachedSanctum reports whether the enemy has contacted the Sanctum.
func raidEnemyReachedSanctum(enemy raidEnemy) bool {
	return enemy.position.Y <= 0
}
