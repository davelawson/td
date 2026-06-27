package game

import "testing"

// TestSelectedRaiderPanelRows verifies raider selection exposes current combat stats.
func TestSelectedRaiderPanelRows(t *testing.T) {
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

	assertPanelRow(t, panel, "Raider Type", "Skeleton Sword-and-Shield")
	assertPanelRow(t, panel, "Health", "25")
	assertPanelRow(t, panel, "Max Health", "50")
	assertPanelRow(t, panel, "Health Remaining", "50%")
	assertPanelRow(t, panel, "Speed", "1.0 tiles/s")
	assertPanelRow(t, panel, "Sanctum Damage", "1")
}

// TestSelectedHousePanelRows verifies House selection exposes its population effect.
func TestSelectedHousePanelRows(t *testing.T) {
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

	assertPanelRow(t, panel, "Structure", "House")
	assertPanelRow(t, panel, "Cost", "20 Wood")
	assertPanelRow(t, panel, "Grants Peasants", "2")
	assertPanelRowAbsent(t, panel, "Range")
	assertPanelRowAbsent(t, panel, "Damage")
}

// TestSelectedBarracksPanelRows verifies Barracks selection exposes its conversion effect.
func TestSelectedBarracksPanelRows(t *testing.T) {
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

	assertPanelRow(t, panel, "Structure", "Barracks")
	assertPanelRow(t, panel, "Cost", "10 Wood, 10 Stone")
	assertPanelRow(t, panel, "Consumes Peasants", "2")
	assertPanelRow(t, panel, "Grants Soldiers", "2")
	assertPanelRowAbsent(t, panel, "Range")
	assertPanelRowAbsent(t, panel, "Damage")
}

// TestSelectedEconomicBuildingPanelRows verifies resource producers expose yield details.
func TestSelectedEconomicBuildingPanelRows(t *testing.T) {
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

	assertPanelRow(t, panel, "Structure", "Stone Quarry")
	assertPanelRow(t, panel, "Cost", "10 Wood, 10 Stone")
	assertPanelRow(t, panel, "Required Peasants", "1")
	assertPanelRow(t, panel, "Produces", "10 Stone after each Raid")
	assertPanelRowAbsent(t, panel, "Range")
	assertPanelRowAbsent(t, panel, "Damage")
}

// TestSelectedBowTowerPanelRows verifies Bow Tower selection exposes tower stats.
func TestSelectedBowTowerPanelRows(t *testing.T) {
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

	assertPanelRow(t, panel, "Tower Type", "Bow Tower")
	assertPanelRow(t, panel, "Range", "3.0 tiles")
	assertPanelRow(t, panel, "Attack Speed", "every 1.0s")
	assertPanelRow(t, panel, "Damage", "10")
	assertPanelRow(t, panel, "Required Soldiers", "1")
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

	assertPanelRow(t, panel, "Tower Type", "Flame Bolt Tower")
	assertPanelRow(t, panel, "Range", "2.5 tiles")
	assertPanelRow(t, panel, "Attack Speed", "every 1.5s")
	assertPanelRow(t, panel, "Damage", "20")
	assertPanelRow(t, panel, "Required Apprentices", "1")
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

	assertPanelRow(t, panel, "Tower Type", "Catapult Tower")
	assertPanelRow(t, panel, "Range", "5.0 tiles")
	assertPanelRow(t, panel, "Attack Speed", "every 3.0s")
	assertPanelRow(t, panel, "Damage", "75")
	assertPanelRow(t, panel, "Required Soldiers", "1")
	assertPanelRow(t, panel, "Required Peasants", "2")
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

	if len(panel.Rows) != 1 {
		t.Fatalf("panel rows = %d, want 1", len(panel.Rows))
	}
	assertPanelRow(t, panel, "Structure", "Sanctum")
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

// TestFormatResourceCost verifies selection panels format construction costs.
func TestFormatResourceCost(t *testing.T) {
	tests := []struct {
		cost Resources
		want string
	}{
		{cost: Resources{}, want: "Free"},
		{cost: Resources{Wood: 20}, want: "20 Wood"},
		{cost: Resources{Wood: 10, Stone: 10}, want: "10 Wood, 10 Stone"},
		{cost: Resources{Wood: 10, Stone: 10, Metal: 10}, want: "10 Wood, 10 Stone, 10 Metal"},
		{cost: Resources{Wood: 30, Stone: 10, Metal: 10}, want: "30 Wood, 10 Stone, 10 Metal"},
	}
	for _, test := range tests {
		if got := formatResourceCost(test.cost); got != test.want {
			t.Fatalf("formatResourceCost(%+v) = %q, want %q", test.cost, got, test.want)
		}
	}
}

// assertPanelRow verifies one label and value pair exists in a selection panel.
func assertPanelRow(t *testing.T, panel selectionPanel, label, value string) {
	t.Helper()
	for _, row := range panel.Rows {
		if row.Label == label {
			if row.Value != value {
				t.Fatalf("%s value = %q, want %q", label, row.Value, value)
			}
			return
		}
	}
	t.Fatalf("missing panel row %q in %+v", label, panel.Rows)
}

// assertPanelRowAbsent verifies one label does not exist in a selection panel.
func assertPanelRowAbsent(t *testing.T, panel selectionPanel, label string) {
	t.Helper()
	for _, row := range panel.Rows {
		if row.Label == label {
			t.Fatalf("unexpected panel row %q in %+v", label, panel.Rows)
		}
	}
}
