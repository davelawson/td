package game

import "testing"

// TestCombatDefeatsGrantOnlyTieredGold verifies every raider's deterministic drop.
func TestCombatDefeatsGrantOnlyTieredGold(t *testing.T) {
	state := newRaidTestState(t)
	templates := []*EnemyTemplate{
		&state.enemyCatalog.SkeletonSwordShield,
		&state.enemyCatalog.Zombie,
		&state.enemyCatalog.Ghoul,
		&state.enemyCatalog.ArmouredSkeleton,
	}
	startingMaterials := resourceCounts{
		wood:  state.status.resources.wood,
		stone: state.status.resources.stone,
		iron:  state.status.resources.iron,
	}
	wantGold := state.status.resources.gold

	for index, template := range templates {
		state.raid.enemies = []raidEnemy{{id: index, template: template, health: 1}}
		state.damageEnemy(0, 1)
		wantGold += template.GoldDrop
		if state.status.resources.gold != wantGold {
			t.Fatalf("%s defeat Gold = %d, want %d", template.Name, state.status.resources.gold, wantGold)
		}
		if state.status.resources.wood != startingMaterials.wood ||
			state.status.resources.stone != startingMaterials.stone ||
			state.status.resources.iron != startingMaterials.iron {
			t.Fatalf("%s defeat changed materials: %+v", template.Name, state.status.resources)
		}
	}
}

// TestAreaDefeatSumsGoldDrops verifies one impact rewards every defeated raider once.
func TestAreaDefeatSumsGoldDrops(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{
		{id: 1, template: &state.enemyCatalog.Ghoul, position: coord{X: 0, Y: 2}, health: 10},
		{id: 2, template: &state.enemyCatalog.ArmouredSkeleton, position: coord{X: 0.2, Y: 2.2}, health: 10},
	}
	starting := state.status.resources

	state.damageEnemiesInTargetTile(coord{X: 0, Y: 2}, 10)

	if len(state.raid.enemies) != 0 {
		t.Fatalf("active enemies = %d, want 0", len(state.raid.enemies))
	}
	want := starting
	want.gold += 8
	if state.status.resources != want {
		t.Fatalf("resources = %+v, want only 8 added Gold: %+v", state.status.resources, want)
	}
}
