package game

import (
	"math"
	"testing"
)

// TestSpawnRaidEnemyAssignsHealthAndStableIDs verifies enemies enter combat with targetable state.
func TestSpawnRaidEnemyAssignsHealthAndStableIDs(t *testing.T) {
	state := newRaidTestState(t)

	state.startNextRaid()
	first := state.raid.enemies[0]
	state.spawnRaidEnemy()
	second := state.raid.enemies[1]

	if first.id != 0 || second.id != 1 {
		t.Fatalf("enemy IDs = %d, %d; want 0, 1", first.id, second.id)
	}
	if first.health != state.enemyCatalog.SkeletonSwordShield.MaxHealth {
		t.Fatalf("first enemy health = %d, want %d", first.health, state.enemyCatalog.SkeletonSwordShield.MaxHealth)
	}
	if second.health != second.template.MaxHealth {
		t.Fatalf("second enemy health = %d, want %d", second.health, second.template.MaxHealth)
	}
}

// TestBowTowerDoesNotFireOutsideRange verifies range gates projectile launch.
func TestBowTowerDoesNotFireOutsideRange(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingFlameBoltTower(state)
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 7}, 20)}

	state.updateCombat()

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// TestBowTowerFiresAtEnemyInRangeAndStartsCooldown verifies launch and second-based cooldown state.
func TestBowTowerFiresAtEnemyInRangeAndStartsCooldown(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingFlameBoltTower(state)
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 2}, 20)}

	state.updateCombat()

	if len(state.combat.projectiles) != 1 {
		t.Fatalf("projectiles = %d, want 1", len(state.combat.projectiles))
	}
	projectile := state.combat.projectiles[0]
	if projectile.targetID != 0 {
		t.Fatalf("projectile target ID = %d, want 0", projectile.targetID)
	}
	if projectile.damage != 10 {
		t.Fatalf("projectile damage = %d, want 10", projectile.damage)
	}
	if projectile.speedTilesPerSecond != 9.0 {
		t.Fatalf("projectile speed = %f, want 9.0", projectile.speedTilesPerSecond)
	}
	key := tileCoordinate{X: homePlotCenter + 1, Y: 5}
	if got, want := state.combat.towerCooldowns[key], 1.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("cooldown seconds = %f, want %f", got, want)
	}
}

// TestFlameBoltTowerFiresAtEnemyInRange verifies the authored flame tower combat stats.
func TestFlameBoltTowerFiresAtEnemyInRange(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingBowTower(state)
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 2}, 20)}

	state.updateCombat()

	if len(state.combat.projectiles) != 1 {
		t.Fatalf("projectiles = %d, want 1", len(state.combat.projectiles))
	}
	projectile := state.combat.projectiles[0]
	if projectile.damage != 20 {
		t.Fatalf("projectile damage = %d, want 20", projectile.damage)
	}
	if projectile.speedTilesPerSecond != 7.0 {
		t.Fatalf("projectile speed = %f, want 7.0", projectile.speedTilesPerSecond)
	}
	key := tileCoordinate{X: homePlotCenter - 1, Y: 5}
	if got, want := state.combat.towerCooldowns[key], 1.5; math.Abs(got-want) > 0.000001 {
		t.Fatalf("cooldown seconds = %f, want %f", got, want)
	}
}

// TestCatapultTowerFiresAtEnemyInRange verifies the authored Catapult Tower combat stats.
func TestCatapultTowerFiresAtEnemyInRange(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingBowTower(state)
	removeStartingFlameBoltTower(state)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureCatapultTower
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 2}, 100)}

	state.updateCombat()

	if len(state.combat.projectiles) != 1 {
		t.Fatalf("projectiles = %d, want 1", len(state.combat.projectiles))
	}
	projectile := state.combat.projectiles[0]
	if projectile.damage != 75 {
		t.Fatalf("projectile damage = %d, want 75", projectile.damage)
	}
	if projectile.speedTilesPerSecond != 3.0 {
		t.Fatalf("projectile speed = %f, want 3.0", projectile.speedTilesPerSecond)
	}
	if !projectile.damageAllEnemiesInTargetTile {
		t.Fatal("expected Catapult projectile to damage all enemies in its target Tile")
	}
	key := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	if got, want := state.combat.towerCooldowns[key], 3.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("cooldown seconds = %f, want %f", got, want)
	}
}

// TestBowTowerTargetPriorityUsesClosestEnemyToSanctum verifies urgent targets are chosen first.
func TestBowTowerTargetPriorityUsesClosestEnemyToSanctum(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingFlameBoltTower(state)
	state.raid.enemies = []raidEnemy{
		combatTestEnemy(0, coord{X: 0, Y: 3}, 20),
		combatTestEnemy(1, coord{X: 0, Y: 1}, 20),
	}

	state.updateCombat()

	if len(state.combat.projectiles) != 1 {
		t.Fatalf("projectiles = %d, want 1", len(state.combat.projectiles))
	}
	if got, want := state.combat.projectiles[0].targetID, 1; got != want {
		t.Fatalf("projectile target ID = %d, want %d", got, want)
	}
}

// TestProjectileHitDamagesEnemy verifies impact applies Bow Tower damage.
func TestProjectileHitDamagesEnemy(t *testing.T) {
	state := newRaidTestState(t)
	sink := &recordingSoundSink{}
	state.SetSoundSink(sink)
	state.raid.enemies = []raidEnemy{combatTestEnemyWithTemplate(0, coord{X: 0, Y: 2}, 20, &state.enemyCatalog.SkeletonSwordShield)}
	startingResources := state.status.resources
	state.combat.projectiles = []combatProjectile{{
		targetID:            0,
		position:            coord{X: 0, Y: 2},
		damage:              10,
		speedTilesPerSecond: 9.0,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if got, want := state.raid.enemies[0].health, 10; got != want {
		t.Fatalf("enemy health = %d, want %d", got, want)
	}
	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
	if sink.raiderDefeated != 0 {
		t.Fatalf("raider defeated sounds = %d, want 0", sink.raiderDefeated)
	}
	if state.status.resources != startingResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, startingResources)
	}
}

// TestProjectileHitRemovesDefeatedEnemy verifies health reaching zero defeats the enemy.
func TestProjectileHitRemovesDefeatedEnemy(t *testing.T) {
	state := newRaidTestState(t)
	sink := &recordingSoundSink{}
	state.SetSoundSink(sink)
	state.raid.enemies = []raidEnemy{combatTestEnemyWithTemplate(0, coord{X: 0, Y: 2}, 10, &state.enemyCatalog.SkeletonSwordShield)}
	startingResources := state.status.resources
	state.combat.projectiles = []combatProjectile{{
		targetID:            0,
		position:            coord{X: 0, Y: 2},
		damage:              10,
		speedTilesPerSecond: 9.0,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
	if sink.raiderDefeated != 1 {
		t.Fatalf("raider defeated sounds = %d, want 1", sink.raiderDefeated)
	}
	wantResources := resourceCounts{
		wood:  startingResources.wood + state.enemyCatalog.SkeletonSwordShield.Resources.Wood,
		stone: startingResources.stone + state.enemyCatalog.SkeletonSwordShield.Resources.Stone,
		metal: startingResources.metal + state.enemyCatalog.SkeletonSwordShield.Resources.Metal,
	}
	if state.status.resources != wantResources {
		t.Fatalf("resources = %+v, want %+v", state.status.resources, wantResources)
	}
}

// TestCatapultProjectileHitDamagesEnemiesInTargetTile verifies Catapult area impact.
func TestCatapultProjectileHitDamagesEnemiesInTargetTile(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{
		combatTestEnemy(0, coord{X: 0, Y: 2}, 100),
		combatTestEnemy(1, coord{X: 0.2, Y: 2.2}, 100),
		combatTestEnemy(2, coord{X: 0, Y: 1}, 100),
	}
	state.combat.projectiles = []combatProjectile{{
		targetID:                     0,
		position:                     coord{X: 0, Y: 2},
		damage:                       75,
		speedTilesPerSecond:          3.0,
		damageAllEnemiesInTargetTile: true,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if got, want := state.raid.enemies[0].health, 25; got != want {
		t.Fatalf("target enemy health = %d, want %d", got, want)
	}
	if got, want := state.raid.enemies[1].health, 25; got != want {
		t.Fatalf("same-Tile enemy health = %d, want %d", got, want)
	}
	if got, want := state.raid.enemies[2].health, 100; got != want {
		t.Fatalf("adjacent-Tile enemy health = %d, want %d", got, want)
	}
	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// TestCatapultProjectileHitRemovesDefeatedEnemiesInTargetTile verifies area defeats are removed.
func TestCatapultProjectileHitRemovesDefeatedEnemiesInTargetTile(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{
		combatTestEnemy(0, coord{X: 0, Y: 2}, 75),
		combatTestEnemy(1, coord{X: 0.2, Y: 2.2}, 75),
		combatTestEnemy(2, coord{X: 0, Y: 1}, 75),
	}
	state.combat.projectiles = []combatProjectile{{
		targetID:                     0,
		position:                     coord{X: 0, Y: 2},
		damage:                       75,
		speedTilesPerSecond:          3.0,
		damageAllEnemiesInTargetTile: true,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if len(state.raid.enemies) != 1 {
		t.Fatalf("active enemies = %d, want 1", len(state.raid.enemies))
	}
	if state.raid.enemies[0].id != 2 {
		t.Fatalf("remaining enemy ID = %d, want adjacent-Tile enemy 2", state.raid.enemies[0].id)
	}
}

// TestSanctumContactDoesNotGrantDefeatEffects verifies Barricade removal is not a combat defeat.
func TestSanctumContactDoesNotGrantDefeatEffects(t *testing.T) {
	state := newRaidTestState(t)
	sink := &recordingSoundSink{}
	state.SetSoundSink(sink)
	state.raid.enemies = []raidEnemy{combatTestEnemyWithTemplate(0, coord{X: 0, Y: 0}, 10, &state.enemyCatalog.SkeletonSwordShield)}
	startingResources := state.status.resources

	state.updateRaidEnemies()

	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
	if sink.raiderDefeated != 0 {
		t.Fatalf("raider defeated sounds = %d, want 0", sink.raiderDefeated)
	}
	if state.status.resources != startingResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, startingResources)
	}
}

// TestProjectileDisappearsWhenTargetIsGone verifies projectiles do not retarget.
func TestProjectileDisappearsWhenTargetIsGone(t *testing.T) {
	state := newRaidTestState(t)
	state.combat.projectiles = []combatProjectile{{
		targetID:            99,
		position:            coord{X: 0, Y: 2},
		damage:              10,
		speedTilesPerSecond: 9.0,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// TestCombatDoesNotAdvanceWhilePaused verifies pause blocks tower firing.
func TestCombatDoesNotAdvanceWhilePaused(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.active = true
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 2}, 20)}
	state.Update(Input{TogglePause: true})

	state.Update(Input{})

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// TestCombatDoesNotAdvanceWhileIngameMenuOpen verifies overlay-open frames block combat.
func TestCombatDoesNotAdvanceWhileIngameMenuOpen(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.active = true
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, coord{X: 0, Y: 2}, 20)}
	state.Update(Input{ToggleMenu: true})

	state.Update(Input{})

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// combatTestEnemy creates a targetable enemy for focused combat tests.
func combatTestEnemy(id int, position coord, health int) raidEnemy {
	return raidEnemy{
		id:       id,
		position: position,
		health:   health,
	}
}

// combatTestEnemyWithTemplate creates a targetable enemy with reward-bearing template data.
func combatTestEnemyWithTemplate(id int, position coord, health int, template *EnemyTemplate) raidEnemy {
	enemy := combatTestEnemy(id, position, health)
	enemy.template = template
	return enemy
}

// removeStartingBowTower removes the default Bow Tower from focused combat tests.
func removeStartingBowTower(state *State) {
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureNone
}

// removeStartingFlameBoltTower removes the default Flame Bolt Tower from focused combat tests.
func removeStartingFlameBoltTower(state *State) {
	state.gameMap.Home.Tiles[5][homePlotCenter-1].Feature = featureNone
}

type recordingSoundSink struct {
	raiderDefeated int
}

// PlayRaiderDefeated records one raider-defeated sound event.
func (s *recordingSoundSink) PlayRaiderDefeated() {
	s.raiderDefeated++
}
