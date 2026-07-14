package game

import (
	"testing"

	"td/internal/ui"
)

// TestSelectedRaiderPanelData verifies raider selection exposes current combat stats.
func TestSelectedRaiderPanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.raid.enemies = []raidEnemy{{
		id:       41,
		template: &state.enemyCatalog.SkeletonSwordShield,
		position: coord{X: 0, Y: 3},
		health:   25,
	}}
	state.selection = selectedItem{kind: selectedItemRaider, raiderID: 41}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected raider panel")
	}

	if panel.Kind != ui.SelectionPanelRaider {
		t.Fatalf("kind = %v, want raider", panel.Kind)
	}
	if panel.Name != "Skeleton Sword-and-Shield" {
		t.Fatalf("name = %q, want Skeleton Sword-and-Shield", panel.Name)
	}
	if panel.Health != 25 || panel.MaxHealth != 50 {
		t.Fatalf("health = %d/%d, want 25/50", panel.Health, panel.MaxHealth)
	}
	if panel.SpeedTilesPerSecond != 1.0 {
		t.Fatalf("speed = %.1f, want 1.0", panel.SpeedTilesPerSecond)
	}
	if panel.SanctumDamage != 1 {
		t.Fatalf("sanctum damage = %d, want 1", panel.SanctumDamage)
	}
}

// TestSelectedTerrainPanelData verifies terrain selection exposes its type and Plot biome.
func TestSelectedTerrainPanelData(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainTree}
	state.selection = selectedItem{kind: selectedItemTerrain, tile: tile}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected terrain panel")
	}
	if panel.Kind != ui.SelectionPanelTerrain {
		t.Fatalf("kind = %v, want terrain", panel.Kind)
	}
	if panel.TerrainName != "Tree" || panel.BiomeName != "Grasslands" {
		t.Fatalf("terrain panel = %q/%q, want Tree/Grasslands", panel.TerrainName, panel.BiomeName)
	}
}

// TestSelectedTerrainPanelUsesExploredPlotBiome verifies non-home biome context is retained.
func TestSelectedTerrainPanelUsesExploredPlotBiome(t *testing.T) {
	state := newRaidTestState(t)
	plotCoord := plotCoordinate{X: 1}
	plot := Plot{Biome: biomeHills}
	tile := tileCoordinate{Plot: plotCoord, X: 2, Y: 3}
	plot.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainBoulder}
	state.gameMap.ensurePlots()
	state.gameMap.Plots[plotCoord] = &plot
	state.selection = selectedItem{kind: selectedItemTerrain, tile: tile}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected explored-Plot terrain panel")
	}
	if panel.TerrainName != "Boulder" || panel.BiomeName != "Hills" {
		t.Fatalf("terrain panel = %q/%q, want Boulder/Hills", panel.TerrainName, panel.BiomeName)
	}
}

// TestSelectedIronDepositPanelData verifies deposit selection uses the player-facing name.
func TestSelectedIronDepositPanelData(t *testing.T) {
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainIronDeposit}
	state.selection = selectedItem{kind: selectedItemTerrain, tile: tile}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Iron Deposit panel")
	}
	if panel.TerrainName != "Iron Deposit" || panel.BiomeName != "Grasslands" {
		t.Fatalf("terrain panel = %q/%q, want Iron Deposit/Grasslands", panel.TerrainName, panel.BiomeName)
	}
}

// TestSelectedHousePanelData verifies House selection exposes its population effect.
func TestSelectedHousePanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureHouse
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 2, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected House panel")
	}

	if panel.Kind != ui.SelectionPanelPopulationBuilding {
		t.Fatalf("kind = %v, want population building", panel.Kind)
	}
	if panel.Name != "House" {
		t.Fatalf("name = %q, want House", panel.Name)
	}
	if panel.Cost != (ui.ResourceAmounts{Wood: 20}) {
		t.Fatalf("cost = %+v, want 20 Wood", panel.Cost)
	}
	if panel.PopulationGrant != (ui.PopulationAmounts{Peasants: 2}) {
		t.Fatalf("population grant = %+v, want 2 Peasants", panel.PopulationGrant)
	}
	if panel.RangeTiles != 0 || panel.Damage != 0 {
		t.Fatalf("combat stats = range %.1f damage %d, want empty", panel.RangeTiles, panel.Damage)
	}
}

// TestSelectedBarracksPanelData verifies Barracks selection exposes its conversion effect.
func TestSelectedBarracksPanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+3].Feature = featureBarracks
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 3, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Barracks panel")
	}

	if panel.Kind != ui.SelectionPanelPopulationBuilding {
		t.Fatalf("kind = %v, want population building", panel.Kind)
	}
	if panel.Name != "Barracks" {
		t.Fatalf("name = %q, want Barracks", panel.Name)
	}
	if panel.Cost != (ui.ResourceAmounts{Wood: 10, Stone: 10}) {
		t.Fatalf("cost = %+v, want 10 Wood, 10 Stone", panel.Cost)
	}
	if panel.PopulationCost != (ui.PopulationAmounts{Peasants: 2}) {
		t.Fatalf("population cost = %+v, want 2 Peasants", panel.PopulationCost)
	}
	if panel.PopulationGrant != (ui.PopulationAmounts{Soldiers: 2}) {
		t.Fatalf("population grant = %+v, want 2 Soldiers", panel.PopulationGrant)
	}
}

// TestSelectedDormPanelData verifies Dorm selection exposes its conversion effect.
func TestSelectedDormPanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+3].Feature = featureDorm
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 3, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Dorm panel")
	}

	if panel.Kind != ui.SelectionPanelPopulationBuilding {
		t.Fatalf("kind = %v, want population building", panel.Kind)
	}
	if panel.Name != "Dorm" {
		t.Fatalf("name = %q, want Dorm", panel.Name)
	}
	if panel.Cost != (ui.ResourceAmounts{Wood: 10, Stone: 10}) {
		t.Fatalf("cost = %+v, want 10 Wood, 10 Stone", panel.Cost)
	}
	if panel.PopulationCost != (ui.PopulationAmounts{Peasants: 1}) {
		t.Fatalf("population cost = %+v, want 1 Peasant", panel.PopulationCost)
	}
	if panel.PopulationGrant != (ui.PopulationAmounts{Apprentices: 1}) {
		t.Fatalf("population grant = %+v, want 1 Apprentice", panel.PopulationGrant)
	}
}

// TestSelectedEconomicBuildingPanelData verifies resource producers expose yield details.
func TestSelectedEconomicBuildingPanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureStoneQuarry
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 2, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Stone Quarry panel")
	}

	if panel.Kind != ui.SelectionPanelEconomicBuilding {
		t.Fatalf("kind = %v, want economic building", panel.Kind)
	}
	if panel.Name != "Stone Quarry" {
		t.Fatalf("name = %q, want Stone Quarry", panel.Name)
	}
	if panel.Cost != (ui.ResourceAmounts{Wood: 10, Stone: 10}) {
		t.Fatalf("cost = %+v, want 10 Wood, 10 Stone", panel.Cost)
	}
	if panel.Staffing != (ui.PopulationAmounts{Peasants: 1}) {
		t.Fatalf("staffing = %+v, want 1 Peasant", panel.Staffing)
	}
	if panel.ResourceYield != (ui.ResourceAmounts{Stone: 10}) {
		t.Fatalf("yield = %+v, want 10 Stone", panel.ResourceYield)
	}
}

// TestSelectedBowTowerPanelData verifies Bow Tower selection exposes tower stats.
func TestSelectedBowTowerPanelData(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 1, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Bow Tower panel")
	}

	if panel.Kind != ui.SelectionPanelTower {
		t.Fatalf("kind = %v, want tower", panel.Kind)
	}
	if panel.Name != "Bow Tower" {
		t.Fatalf("name = %q, want Bow Tower", panel.Name)
	}
	if panel.RangeTiles != 3.0 || panel.FireIntervalSeconds != 1.0 || panel.Damage != 10 {
		t.Fatalf("tower stats = range %.1f speed %.1f damage %d, want 3.0/1.0/10", panel.RangeTiles, panel.FireIntervalSeconds, panel.Damage)
	}
	if panel.Staffing != (ui.PopulationAmounts{Soldiers: 1}) {
		t.Fatalf("staffing = %+v, want 1 Soldier", panel.Staffing)
	}
}

// TestSelectedFlameBoltTowerPanelRows verifies Flame Bolt Tower selection exposes tower stats.
func TestSelectedFlameBoltTowerPanelRows(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter-1].Feature = featureFlameBoltTower
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter - 1, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Flame Bolt Tower panel")
	}

	if panel.Kind != ui.SelectionPanelTower {
		t.Fatalf("kind = %v, want tower", panel.Kind)
	}
	if panel.Name != "Flame Bolt Tower" {
		t.Fatalf("name = %q, want Flame Bolt Tower", panel.Name)
	}
	if panel.RangeTiles != 2.5 || panel.FireIntervalSeconds != 1.5 || panel.Damage != 20 {
		t.Fatalf("tower stats = range %.1f speed %.1f damage %d, want 2.5/1.5/20", panel.RangeTiles, panel.FireIntervalSeconds, panel.Damage)
	}
	if panel.Staffing != (ui.PopulationAmounts{Apprentices: 1}) {
		t.Fatalf("staffing = %+v, want 1 Apprentice", panel.Staffing)
	}
}

// TestSelectedCatapultTowerPanelRows verifies Catapult Tower selection exposes tower stats.
func TestSelectedCatapultTowerPanelRows(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+2].Feature = featureCatapultTower
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter + 2, Y: 5},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Catapult Tower panel")
	}

	if panel.Kind != ui.SelectionPanelTower {
		t.Fatalf("kind = %v, want tower", panel.Kind)
	}
	if panel.Name != "Catapult Tower" {
		t.Fatalf("name = %q, want Catapult Tower", panel.Name)
	}
	if panel.RangeTiles != 5.0 || panel.FireIntervalSeconds != 6.0 || panel.Damage != 30 {
		t.Fatalf("tower stats = range %.1f speed %.1f damage %d, want 5.0/6.0/30", panel.RangeTiles, panel.FireIntervalSeconds, panel.Damage)
	}
	if panel.Staffing != (ui.PopulationAmounts{Soldiers: 1, Peasants: 1}) {
		t.Fatalf("staffing = %+v, want 1 Soldier and 1 Peasant", panel.Staffing)
	}
}

// TestSelectedSanctumPanelRows verifies Sanctum selection exposes only basic structure identity.
func TestSelectedSanctumPanelRows(t *testing.T) {
	state := newRaidTestState(t)
	state.selection = selectedItem{
		kind: selectedItemStructure,
		tile: tileCoordinate{X: homePlotCenter, Y: homePlotCenter},
	}

	panel, ok := state.currentSelectionPanel()
	if !ok {
		t.Fatal("expected selected Sanctum panel")
	}

	if panel.Kind != ui.SelectionPanelStructure {
		t.Fatalf("kind = %v, want structure", panel.Kind)
	}
	if panel.Name != "Sanctum" {
		t.Fatalf("name = %q, want Sanctum", panel.Name)
	}
}

// TestSelectionPanelClickDoesNotClearSelection verifies panel clicks are blocked as UI input.
func TestSelectionPanelClickDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.Update(clickTileInput(state, homePlotCenter+1, 5))
	state.paused = true
	bounds, ok := state.selectionPanelBounds()
	if !ok {
		t.Fatal("expected selected Bow Tower panel bounds")
	}

	state.Update(Input{
		CursorX: bounds.X + bounds.W/2,
		CursorY: bounds.Y + bounds.H/2,
		Clicked: true,
	})

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != (tileCoordinate{X: homePlotCenter + 1, Y: 5}) {
		t.Fatalf("selected tile = %+v, want Bow Tower", state.selection.tile)
	}
}
