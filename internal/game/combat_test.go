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
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 7}, 20)}

	state.updateCombat()

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// TestBowTowerFiresAtEnemyInRangeAndStartsCooldown verifies launch and second-based cooldown state.
func TestBowTowerFiresAtEnemyInRangeAndStartsCooldown(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingFlameBoltTower(state)
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 20)}

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
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 20)}

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

// TestBowTowerTargetPriorityUsesClosestEnemyToSanctum verifies urgent targets are chosen first.
func TestBowTowerTargetPriorityUsesClosestEnemyToSanctum(t *testing.T) {
	state := newRaidTestState(t)
	removeStartingFlameBoltTower(state)
	state.raid.enemies = []raidEnemy{
		combatTestEnemy(0, worldPosition{X: 0, Y: 3}, 20),
		combatTestEnemy(1, worldPosition{X: 0, Y: 1}, 20),
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
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 20)}
	state.combat.projectiles = []combatProjectile{{
		targetID:            0,
		position:            worldPosition{X: 0, Y: 2},
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
}

// TestProjectileHitRemovesDefeatedEnemy verifies health reaching zero defeats the enemy.
func TestProjectileHitRemovesDefeatedEnemy(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 10)}
	state.combat.projectiles = []combatProjectile{{
		targetID:            0,
		position:            worldPosition{X: 0, Y: 2},
		damage:              10,
		speedTilesPerSecond: 9.0,
	}}

	state.updateProjectiles(gameUpdateSeconds)

	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
}

// TestProjectileDisappearsWhenTargetIsGone verifies projectiles do not retarget.
func TestProjectileDisappearsWhenTargetIsGone(t *testing.T) {
	state := newRaidTestState(t)
	state.combat.projectiles = []combatProjectile{{
		targetID:            99,
		position:            worldPosition{X: 0, Y: 2},
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
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 20)}
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
	state.raid.enemies = []raidEnemy{combatTestEnemy(0, worldPosition{X: 0, Y: 2}, 20)}
	state.Update(Input{ToggleMenu: true})

	state.Update(Input{})

	if len(state.combat.projectiles) != 0 {
		t.Fatalf("projectiles = %d, want 0", len(state.combat.projectiles))
	}
}

// combatTestEnemy creates a targetable enemy for focused combat tests.
func combatTestEnemy(id int, position worldPosition, health int) raidEnemy {
	return raidEnemy{
		id:       id,
		position: position,
		health:   health,
	}
}

// removeStartingBowTower removes the default Bow Tower from focused combat tests.
func removeStartingBowTower(state *State) {
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureNone
}

// removeStartingFlameBoltTower removes the default Flame Bolt Tower from focused combat tests.
func removeStartingFlameBoltTower(state *State) {
	state.gameMap.Home.Tiles[5][homePlotCenter-1].Feature = featureNone
}
