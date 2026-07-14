package game

import (
	"testing"

	"td/internal/ui"
)

// TestBuildingBarModelAdaptsStructureCatalog verifies game facts cross the UI boundary intact.
func TestBuildingBarModelAdaptsStructureCatalog(t *testing.T) {
	state := newRaidTestState(t)
	model := state.buildingBarModel()

	if model.SelectedCategory != ui.BuildingBarCategoryHousing {
		t.Fatalf("selected category = %v, want Housing", model.SelectedCategory)
	}
	if len(model.Items) != 10 {
		t.Fatalf("items = %d, want 10", len(model.Items))
	}

	cases := []struct {
		action    ui.BuildingBarAction
		name      string
		cost      ui.ResourceAmounts
		staff     ui.PopulationAmounts
		popCost   ui.PopulationAmounts
		grant     ui.PopulationAmounts
		buildable bool
	}{
		{ui.BuildingBarHouse, "House", ui.ResourceAmounts{Wood: 20}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, ui.PopulationAmounts{Peasants: 2}, true},
		{ui.BuildingBarBarracks, "Barracks", ui.ResourceAmounts{Wood: 10, Stone: 10}, ui.PopulationAmounts{}, ui.PopulationAmounts{Peasants: 2}, ui.PopulationAmounts{Soldiers: 2}, false},
		{ui.BuildingBarDorm, "Dorm", ui.ResourceAmounts{Wood: 10, Stone: 10}, ui.PopulationAmounts{}, ui.PopulationAmounts{Peasants: 1}, ui.PopulationAmounts{Apprentices: 1}, false},
		{ui.BuildingBarWoodcutter, "Woodcutter", ui.ResourceAmounts{Wood: 10}, ui.PopulationAmounts{Peasants: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarStoneQuarry, "Stone Quarry", ui.ResourceAmounts{Wood: 10, Stone: 10}, ui.PopulationAmounts{Peasants: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarIronMine, "Iron Mine", ui.ResourceAmounts{Wood: 10, Stone: 10, Iron: 10}, ui.PopulationAmounts{Peasants: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarMarket, "Market", ui.ResourceAmounts{Wood: 30}, ui.PopulationAmounts{Soldiers: 1, Peasants: 2}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarBowTower, "Bow Tower", ui.ResourceAmounts{Wood: 20, Stone: 10}, ui.PopulationAmounts{Soldiers: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarFlameBoltTower, "Flame Bolt Tower", ui.ResourceAmounts{Stone: 30, Iron: 20}, ui.PopulationAmounts{Apprentices: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
		{ui.BuildingBarCatapultTower, "Catapult Tower", ui.ResourceAmounts{Wood: 40, Stone: 60, Iron: 25}, ui.PopulationAmounts{Soldiers: 1, Peasants: 1}, ui.PopulationAmounts{}, ui.PopulationAmounts{}, false},
	}

	for _, tc := range cases {
		item := findBuildingBarModelItem(t, model, tc.action)
		if item.Name != tc.name || item.Cost != tc.cost || item.Staffing != tc.staff ||
			item.PopulationCost != tc.popCost || item.PopulationGrant != tc.grant || item.Buildable != tc.buildable {
			t.Fatalf("action %v = %+v, want name=%q cost=%+v staff=%+v population cost=%+v grant=%+v buildable=%t", tc.action, item, tc.name, tc.cost, tc.staff, tc.popCost, tc.grant, tc.buildable)
		}
		if item.Sprite == nil || item.Description == "" {
			t.Fatalf("%s missing sprite or description", tc.name)
		}
	}
}

// TestBuildingBarModelRecomputesAvailability verifies gameplay capacity drives UI facts.
func TestBuildingBarModelRecomputesAvailability(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 1, 1, 2)
	state.status.resources = resourceCounts{wood: 80, stone: 80, iron: 30}

	model := state.buildingBarModel()
	for _, action := range ui.BuildingBarActions() {
		if !findBuildingBarModelItem(t, model, action).Buildable {
			t.Fatalf("action %v should be buildable with sufficient capacity", action)
		}
	}
}

// TestBuildingBarTabClickSwitchesCategory verifies UI tab actions update host state.
func TestBuildingBarTabClickSwitchesCategory(t *testing.T) {
	state := newRaidTestState(t)
	bounds, ok := ui.BuildingBarCategoryBounds(topBarHeight, ui.BuildingBarCategoryEconomic)
	if !ok {
		t.Fatal("missing Economic tab")
	}

	state.Update(Input{CursorX: bounds.X + bounds.W/2, CursorY: bounds.Y + bounds.H/2, Clicked: true})

	if state.ui.buildBarCategory != ui.BuildingBarCategoryEconomic {
		t.Fatalf("category = %v, want Economic", state.ui.buildBarCategory)
	}
	if state.buildDrag.active {
		t.Fatal("tab click should not start a building drag")
	}
}

// TestBuildingBarInputDoesNotClearSelection verifies widget clicks remain screen-space UI input.
func TestBuildingBarInputDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.Update(clickTileInput(state, homePlotCenter+1, 5))
	bounds, ok := ui.BuildingBarItemBounds(topBarHeight, state.buildingBarModel(), ui.BuildingBarHouse)
	if !ok {
		t.Fatal("missing House bounds")
	}

	state.Update(Input{CursorX: bounds.X + bounds.W/2, CursorY: bounds.Y + bounds.H/2, Clicked: true})

	if state.selection.kind != selectedItemStructure || state.selection.tile != (tileCoordinate{X: homePlotCenter + 1, Y: 5}) {
		t.Fatalf("selection = %+v, want Bow Tower tile", state.selection)
	}
}

// TestBuildingBarHoverClearsWhenIngameMenuOpens verifies overlay state drops widget hover.
func TestBuildingBarHoverClearsWhenIngameMenuOpens(t *testing.T) {
	state := newRaidTestState(t)
	bounds, ok := ui.BuildingBarItemBounds(topBarHeight, state.buildingBarModel(), ui.BuildingBarHouse)
	if !ok {
		t.Fatal("missing House bounds")
	}
	state.Update(Input{CursorX: bounds.X + bounds.W/2, CursorY: bounds.Y + bounds.H/2})
	if state.ui.buildBarHover != 0 {
		t.Fatalf("hover = %d, want first item", state.ui.buildBarHover)
	}

	state.Update(Input{ToggleMenu: true})
	if state.ui.buildBarHover != -1 || state.ui.buildBarTabHover != ui.BuildingBarNoCategory {
		t.Fatalf("hover = %d/%v, want cleared", state.ui.buildBarHover, state.ui.buildBarTabHover)
	}
}

// TestNextRaidButtonAvoidsBuildingBar verifies adjacent game UI still uses the exported bar width.
func TestNextRaidButtonAvoidsBuildingBar(t *testing.T) {
	state := newRaidTestState(t)
	bar := state.buildingBarBounds()
	button := state.nextRaidButton()
	if button.X < bar.X+bar.W {
		t.Fatalf("Next Raid X = %d, want at least %d", button.X, bar.X+bar.W)
	}
}

// findBuildingBarModelItem returns one adapted item or fails the test.
func findBuildingBarModelItem(t *testing.T, model ui.BuildingBarModel, action ui.BuildingBarAction) ui.BuildingBarItem {
	t.Helper()
	for _, item := range model.Items {
		if item.Action == action {
			return item
		}
	}
	t.Fatalf("missing building-bar action %v", action)
	return ui.BuildingBarItem{}
}
