package game

import (
	"image/color"
	"testing"
)

// TestBuildingBarBoundsFillPlayableLeftEdge verifies the bar covers the scene edge below the HUD.
func TestBuildingBarBoundsFillPlayableLeftEdge(t *testing.T) {
	state := newRaidTestState(t)

	bounds := state.buildingBarBounds()

	if bounds.X != 0 {
		t.Fatalf("bar X = %d, want 0", bounds.X)
	}
	if bounds.Y != topBarHeight {
		t.Fatalf("bar Y = %d, want top bar height %d", bounds.Y, topBarHeight)
	}
	if bounds.W != buildingBarWidth {
		t.Fatalf("bar width = %d, want %d", bounds.W, buildingBarWidth)
	}
	if bounds.W != 260 {
		t.Fatalf("bar width = %d, want expanded width 260", bounds.W)
	}
	if bounds.H != state.ui.height-topBarHeight {
		t.Fatalf("bar height = %d, want %d", bounds.H, state.ui.height-topBarHeight)
	}
}

// TestBuildingBarDefaultsToHousing verifies the first visible category shows population buildings.
func TestBuildingBarDefaultsToHousing(t *testing.T) {
	state := newRaidTestState(t)

	if state.ui.buildBarCategory != buildingBarCategoryHousing {
		t.Fatalf("default category = %v, want Housing", state.ui.buildBarCategory)
	}
	items := state.buildingBarItems()

	assertBuildingBarItems(t, state, []buildingBarItemID{buildingBarHouseIndex, buildingBarBarracksIndex, buildingBarDormIndex})
	if items[0].Sprite != state.structureCatalog.House.Sprite {
		t.Fatal("expected first item to use House sprite")
	}
	if items[0].Cost != (Resources{Wood: 20}) {
		t.Fatalf("House cost = %+v, want 20 wood", items[0].Cost)
	}
	if items[0].Staffing != (StaffingRequirements{}) {
		t.Fatalf("House staffing = %+v, want none", items[0].Staffing)
	}
	if items[0].PopulationGrant != (PopulationGrant{Peasants: 2}) {
		t.Fatalf("House population grant = %+v, want 2 Peasants", items[0].PopulationGrant)
	}
	if items[1].Sprite != state.structureCatalog.Barracks.Sprite {
		t.Fatal("expected second item to use Barracks sprite")
	}
	if items[1].Cost != (Resources{Wood: 10, Stone: 10}) {
		t.Fatalf("Barracks cost = %+v, want 10 wood 10 stone", items[1].Cost)
	}
	if items[1].PopulationCost != (PopulationCost{Peasants: 2}) {
		t.Fatalf("Barracks population cost = %+v, want 2 Peasants", items[1].PopulationCost)
	}
	if items[1].PopulationGrant != (PopulationGrant{Soldiers: 2}) {
		t.Fatalf("Barracks population grant = %+v, want 2 Soldiers", items[1].PopulationGrant)
	}
	if items[2].Sprite != state.structureCatalog.Dorm.Sprite {
		t.Fatal("expected third item to use Dorm sprite")
	}
	if items[2].Cost != (Resources{Wood: 10, Stone: 10}) {
		t.Fatalf("Dorm cost = %+v, want 10 wood 10 stone", items[2].Cost)
	}
	if items[2].PopulationCost != (PopulationCost{Peasants: 1}) {
		t.Fatalf("Dorm population cost = %+v, want 1 Peasant", items[2].PopulationCost)
	}
	if items[2].PopulationGrant != (PopulationGrant{Apprentices: 1}) {
		t.Fatalf("Dorm population grant = %+v, want 1 Apprentice", items[2].PopulationGrant)
	}
}

// TestBuildingBarEconomicTabShowsResourceBuildings verifies economic structures are grouped.
func TestBuildingBarEconomicTabShowsResourceBuildings(t *testing.T) {
	state := newRaidTestState(t)
	state.ui.buildBarCategory = buildingBarCategoryEconomic
	items := state.buildingBarItems()

	assertBuildingBarItems(t, state, []buildingBarItemID{buildingBarWoodcutterIndex, buildingBarStoneQuarryIndex, buildingBarIronMineIndex})
	if items[0].Sprite != state.structureCatalog.Woodcutter.Sprite {
		t.Fatal("expected third item to use Woodcutter sprite")
	}
	if items[0].Cost != (Resources{Wood: 10}) {
		t.Fatalf("Woodcutter cost = %+v, want 10 wood", items[0].Cost)
	}
	if items[0].Staffing != (StaffingRequirements{Peasants: 1}) {
		t.Fatalf("Woodcutter staffing = %+v, want 1 Peasant", items[0].Staffing)
	}
	if items[1].Sprite != state.structureCatalog.StoneQuarry.Sprite {
		t.Fatal("expected fourth item to use Stone Quarry sprite")
	}
	if items[1].Cost != (Resources{Wood: 10, Stone: 10}) {
		t.Fatalf("Stone Quarry cost = %+v, want 10 wood 10 stone", items[1].Cost)
	}
	if items[1].Staffing != (StaffingRequirements{Peasants: 1}) {
		t.Fatalf("Stone Quarry staffing = %+v, want 1 Peasant", items[1].Staffing)
	}
	if items[2].Sprite != state.structureCatalog.IronMine.Sprite {
		t.Fatal("expected fifth item to use Iron Mine sprite")
	}
	if items[2].Cost != (Resources{Wood: 10, Stone: 10, Metal: 10}) {
		t.Fatalf("Iron Mine cost = %+v, want 10 wood 10 stone 10 metal", items[2].Cost)
	}
	if items[2].Staffing != (StaffingRequirements{Peasants: 1}) {
		t.Fatalf("Iron Mine staffing = %+v, want 1 Peasant", items[2].Staffing)
	}
}

// TestBuildingBarDefensesTabShowsTowers verifies defensive structures are grouped.
func TestBuildingBarDefensesTabShowsTowers(t *testing.T) {
	state := newRaidTestState(t)
	state.ui.buildBarCategory = buildingBarCategoryDefenses
	items := state.buildingBarItems()

	assertBuildingBarItems(t, state, []buildingBarItemID{buildingBarBowTowerIndex, buildingBarFlameBoltTowerIndex, buildingBarCatapultTowerIndex})
	if items[0].Sprite != state.structureCatalog.BowTower.Sprite {
		t.Fatal("expected sixth item to use Bow Tower sprite")
	}
	if items[0].Cost != (Resources{Wood: 20, Stone: 10}) {
		t.Fatalf("Bow Tower cost = %+v, want 20 wood 10 stone", items[0].Cost)
	}
	if items[0].Staffing != (StaffingRequirements{Soldiers: 1}) {
		t.Fatalf("Bow Tower staffing = %+v, want 1 Soldier", items[0].Staffing)
	}
	if items[1].Sprite != state.structureCatalog.FlameBoltTower.Sprite {
		t.Fatal("expected seventh item to use Flame Bolt Tower sprite")
	}
	if items[1].Cost != (Resources{Stone: 30, Metal: 20}) {
		t.Fatalf("Flame Bolt Tower cost = %+v, want 30 stone 20 metal", items[1].Cost)
	}
	if items[1].Staffing != (StaffingRequirements{Apprentices: 1}) {
		t.Fatalf("Flame Bolt Tower staffing = %+v, want 1 Apprentice", items[1].Staffing)
	}
	if items[2].Sprite != state.structureCatalog.CatapultTower.Sprite {
		t.Fatal("expected eighth item to use Catapult Tower sprite")
	}
	if items[2].Cost != (Resources{Wood: 40, Stone: 60, Metal: 25}) {
		t.Fatalf("Catapult Tower cost = %+v, want 40 wood 60 stone 25 metal", items[2].Cost)
	}
	if items[2].Staffing != (StaffingRequirements{Soldiers: 1, Peasants: 1}) {
		t.Fatalf("Catapult Tower staffing = %+v, want 1 Soldier and 1 Peasant", items[2].Staffing)
	}
}

// TestBuildingBarTabsExposeCategories verifies tabs are stable and inside the bar.
func TestBuildingBarTabsExposeCategories(t *testing.T) {
	state := newRaidTestState(t)
	tabs := state.buildingBarTabs()

	if len(tabs) != 3 {
		t.Fatalf("tabs = %d, want 3", len(tabs))
	}
	expected := []buildingBarCategory{buildingBarCategoryDefenses, buildingBarCategoryEconomic, buildingBarCategoryHousing}
	for i, category := range expected {
		if tabs[i].Category != category {
			t.Fatalf("tab %d category = %v, want %v", i, tabs[i].Category, category)
		}
		if tabs[i].Label != buildingBarCategoryLabel(category) {
			t.Fatalf("tab %d label = %q, want %q", i, tabs[i].Label, buildingBarCategoryLabel(category))
		}
		assertBuildingBarTabInsideBar(t, state, tabs[i])
	}
}

// TestBuildingBarStaffingItems verifies non-zero roles use stable display ordering.
func TestBuildingBarStaffingItems(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarStaffingItems(StaffingRequirements{
		Apprentices: 3,
		Soldiers:    1,
		Peasants:    2,
	})

	if len(items) != 3 {
		t.Fatalf("staffing items = %d, want 3", len(items))
	}
	if items[0].Count != 3 || items[0].Sprite != state.assetCatalog.Sprite.Icon.Apprentice {
		t.Fatalf("first staffing item = %+v, want 3 Apprentices", items[0])
	}
	if items[1].Count != 1 || items[1].Sprite != state.assetCatalog.Sprite.Icon.Soldier {
		t.Fatalf("second staffing item = %+v, want 1 Soldier", items[1])
	}
	if items[2].Count != 2 || items[2].Sprite != state.assetCatalog.Sprite.Icon.Peasant {
		t.Fatalf("third staffing item = %+v, want 2 Peasants", items[2])
	}
	if got := state.buildingBarStaffingItems(StaffingRequirements{}); len(got) != 0 {
		t.Fatalf("zero staffing items = %+v, want none", got)
	}
}

// TestBuildingBarPopulationGrantItems verifies non-zero grants use stable display ordering.
func TestBuildingBarPopulationGrantItems(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarPopulationGrantItems(PopulationGrant{Peasants: 2})

	if len(items) != 1 {
		t.Fatalf("grant items = %d, want 1", len(items))
	}
	if items[0].Count != 2 || items[0].Value != "+2" || items[0].Sprite != state.assetCatalog.Sprite.Icon.Peasant {
		t.Fatalf("grant item = %+v, want +2 Peasants", items[0])
	}
	if got := state.buildingBarPopulationGrantItems(PopulationGrant{}); len(got) != 0 {
		t.Fatalf("zero grant items = %+v, want none", got)
	}
}

// TestBuildingBarPopulationCostItems verifies non-zero population costs use stable display ordering.
func TestBuildingBarPopulationCostItems(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarPopulationCostItems(PopulationCost{Peasants: 2})

	if len(items) != 1 {
		t.Fatalf("cost items = %d, want 1", len(items))
	}
	if items[0].Count != 2 || items[0].Value != "-2" || items[0].Sprite != state.assetCatalog.Sprite.Icon.Peasant {
		t.Fatalf("cost item = %+v, want -2 Peasants", items[0])
	}
	if got := state.buildingBarPopulationCostItems(PopulationCost{}); len(got) != 0 {
		t.Fatalf("zero population-cost items = %+v, want none", got)
	}
}

// TestBuildingBarPopulationMetadataItemsShowsConversion verifies costs and grants share one row.
func TestBuildingBarPopulationMetadataItemsShowsConversion(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[1]

	items := state.buildingBarPopulationMetadataItems(item)

	if len(items) != 2 {
		t.Fatalf("metadata items = %d, want 2", len(items))
	}
	if items[0].Value != "-2" || items[0].Sprite != state.assetCatalog.Sprite.Icon.Peasant {
		t.Fatalf("first metadata item = %+v, want -2 Peasants", items[0])
	}
	if items[1].Value != "+2" || items[1].Sprite != state.assetCatalog.Sprite.Icon.Soldier {
		t.Fatalf("second metadata item = %+v, want +2 Soldiers", items[1])
	}

	item = state.buildingBarItems()[2]
	items = state.buildingBarPopulationMetadataItems(item)
	if len(items) != 2 {
		t.Fatalf("Dorm metadata items = %d, want 2", len(items))
	}
	if items[0].Value != "-1" || items[0].Sprite != state.assetCatalog.Sprite.Icon.Peasant {
		t.Fatalf("first Dorm metadata item = %+v, want -1 Peasant", items[0])
	}
	if items[1].Value != "+1" || items[1].Sprite != state.assetCatalog.Sprite.Icon.Apprentice {
		t.Fatalf("second Dorm metadata item = %+v, want +1 Apprentice", items[1])
	}
}

// TestBuildingBarCostItems verifies non-zero tower costs render in resource order.
func TestBuildingBarCostItems(t *testing.T) {
	items := buildingBarCostItems(Resources{Wood: 20})
	if len(items) != 1 {
		t.Fatalf("House cost items = %d, want 1", len(items))
	}
	assertCostItem(t, items[0], "20", colors.resourceWood)

	items = buildingBarCostItems(Resources{Wood: 10, Stone: 10})
	if len(items) != 2 {
		t.Fatalf("Barracks cost items = %d, want 2", len(items))
	}
	assertCostItem(t, items[0], "10", colors.resourceWood)
	assertCostItem(t, items[1], "10", colors.resourceStone)

	items = buildingBarCostItems(Resources{Wood: 20, Stone: 10})
	if len(items) != 2 {
		t.Fatalf("cost items = %d, want 2", len(items))
	}
	assertCostItem(t, items[0], "20", colors.resourceWood)
	assertCostItem(t, items[1], "10", colors.resourceStone)

	items = buildingBarCostItems(Resources{Wood: 10, Stone: 10, Metal: 10})
	if len(items) != 3 {
		t.Fatalf("Iron Mine cost items = %d, want 3", len(items))
	}
	assertCostItem(t, items[0], "10", colors.resourceWood)
	assertCostItem(t, items[1], "10", colors.resourceStone)
	assertCostItem(t, items[2], "10", colors.resourceMetal)

	items = buildingBarCostItems(Resources{Stone: 30, Metal: 20})
	if len(items) != 2 {
		t.Fatalf("cost items = %d, want 2", len(items))
	}
	assertCostItem(t, items[0], "30", colors.resourceStone)
	assertCostItem(t, items[1], "20", colors.resourceMetal)

	items = buildingBarCostItems(Resources{Wood: 40, Stone: 60, Metal: 25})
	if len(items) != 3 {
		t.Fatalf("cost items = %d, want 3", len(items))
	}
	assertCostItem(t, items[0], "40", colors.resourceWood)
	assertCostItem(t, items[1], "60", colors.resourceStone)
	assertCostItem(t, items[2], "25", colors.resourceMetal)
}

// TestBuildingBarHoverTracksIconBounds verifies only icon rectangles receive hover state.
func TestBuildingBarHoverTracksIconBounds(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarItems()

	state.updateBuildingBarHover(Input{
		CursorX: items[0].Bounds.X + items[0].Bounds.W/2,
		CursorY: items[0].Bounds.Y + items[0].Bounds.H/2,
	})
	if state.ui.buildBarHover != 0 {
		t.Fatalf("building bar hover = %d, want first item", state.ui.buildBarHover)
	}

	state.updateBuildingBarHover(Input{
		CursorX: items[1].Bounds.X + items[1].Bounds.W/2,
		CursorY: items[1].Bounds.Y + items[1].Bounds.H/2,
	})
	if state.ui.buildBarHover != 1 {
		t.Fatalf("building bar hover = %d, want second item", state.ui.buildBarHover)
	}

	state.ui.buildBarCategory = buildingBarCategoryDefenses
	items = state.buildingBarItems()
	state.updateBuildingBarHover(Input{
		CursorX: items[2].Bounds.X + items[2].Bounds.W/2,
		CursorY: items[2].Bounds.Y + items[2].Bounds.H/2,
	})
	if state.ui.buildBarHover != 2 {
		t.Fatalf("building bar hover = %d, want Catapult Tower visible item", state.ui.buildBarHover)
	}
}

// TestBuildingBarHighlightRequiresResourcesAndStaff verifies all construction inputs gate emphasis.
func TestBuildingBarHighlightRequiresResourcesAndStaff(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarItems()

	state.ui.buildBarHover = 0
	if !state.buildingBarItemHighlighted(0, items[0]) {
		t.Fatal("expected zero-staff House to highlight when resources cover it")
	}
	state.ui.buildBarHover = 1
	if state.buildingBarItemHighlighted(1, items[1]) {
		t.Fatal("expected insufficient Peasants to suppress Barracks highlight")
	}
	setAvailablePopulations(state, 1, 1, 2)
	if !state.buildingBarItemHighlighted(1, items[1]) {
		t.Fatal("expected sufficient resources and Peasants to highlight Barracks")
	}

	state.ui.buildBarCategory = buildingBarCategoryEconomic
	items = state.buildingBarItems()
	state.ui.buildBarHover = 0
	state.status.populations.peasants.available = 0
	if state.buildingBarItemHighlighted(0, items[0]) {
		t.Fatal("expected insufficient Peasants to suppress Woodcutter highlight")
	}
	state.status.populations.peasants.available = 2
	if !state.buildingBarItemHighlighted(0, items[0]) {
		t.Fatal("expected sufficient resources and Peasant to highlight Woodcutter")
	}

	state.ui.buildBarCategory = buildingBarCategoryDefenses
	items = state.buildingBarItems()
	state.ui.buildBarHover = 0
	if !state.buildingBarItemHighlighted(0, items[0]) {
		t.Fatal("expected sufficient resources and staff to highlight Bow Tower")
	}

	state.ui.buildBarHover = 1
	if !state.buildingBarItemHighlighted(1, items[1]) {
		t.Fatal("expected sufficient resources and staff to highlight Flame Bolt Tower")
	}

	state.ui.buildBarHover = 2
	if state.buildingBarItemHighlighted(2, items[2]) {
		t.Fatal("expected insufficient resources to suppress Catapult Tower highlight")
	}
	state.status.resources = resourceCounts{wood: 80, stone: 80, metal: 30}
	if !state.buildingBarItemHighlighted(2, items[2]) {
		t.Fatal("expected Catapult Tower to highlight after resources and staff cover it")
	}
	state.status.populations.peasants.available = 0
	if state.buildingBarItemHighlighted(2, items[2]) {
		t.Fatal("expected a missing Peasant to suppress Catapult Tower highlight")
	}
}

// TestBuildingBarTabClickSwitchesCategory verifies tabs change visible structures.
func TestBuildingBarTabClickSwitchesCategory(t *testing.T) {
	state := newRaidTestState(t)
	tab := state.buildingBarTabs()[1]

	state.Update(Input{
		CursorX: tab.Bounds.X + tab.Bounds.W/2,
		CursorY: tab.Bounds.Y + tab.Bounds.H/2,
		Clicked: true,
	})

	if state.ui.buildBarCategory != buildingBarCategoryEconomic {
		t.Fatalf("category = %v, want Economic", state.ui.buildBarCategory)
	}
	assertBuildingBarItems(t, state, []buildingBarItemID{buildingBarWoodcutterIndex, buildingBarStoneQuarryIndex, buildingBarIronMineIndex})
	if state.buildDrag.active {
		t.Fatal("expected tab click not to start build drag")
	}
}

// TestBuildingBarTabClickDoesNotClearSelection verifies category tabs are UI input.
func TestBuildingBarTabClickDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.Update(clickTileInput(state, homePlotCenter+1, 5))
	tab := state.buildingBarTabs()[0]

	state.Update(Input{
		CursorX: tab.Bounds.X + tab.Bounds.W/2,
		CursorY: tab.Bounds.Y + tab.Bounds.H/2,
		Clicked: true,
	})

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != (tileCoordinate{X: homePlotCenter + 1, Y: 5}) {
		t.Fatalf("selected tile = %+v, want Bow Tower", state.selection.tile)
	}
}

// TestBuildingBarHoverIgnoresCostsAndEmptyBar verifies costs do not trigger icon hover.
func TestBuildingBarHoverIgnoresCostsAndEmptyBar(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.updateBuildingBarHover(Input{
		CursorX: state.buildingBarMetadataX(item),
		CursorY: item.Bounds.Y + buildingBarCostOffsetY,
	})
	if state.ui.buildBarHover != -1 {
		t.Fatalf("building bar hover = %d, want none over cost row", state.ui.buildBarHover)
	}

	state.updateBuildingBarHover(Input{
		CursorX: 1,
		CursorY: topBarHeight + 1,
	})
	if state.ui.buildBarHover != -1 {
		t.Fatalf("building bar hover = %d, want none over empty bar", state.ui.buildBarHover)
	}
}

// TestBuildingBarHoverClearsWhenIngameMenuOpens verifies overlay state drops icon hover.
func TestBuildingBarHoverClearsWhenIngameMenuOpens(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.Update(Input{
		CursorX: item.Bounds.X + item.Bounds.W/2,
		CursorY: item.Bounds.Y + item.Bounds.H/2,
	})
	if state.ui.buildBarHover != 0 {
		t.Fatalf("building bar hover = %d, want first item", state.ui.buildBarHover)
	}

	state.Update(Input{ToggleMenu: true})
	if state.ui.buildBarHover != -1 {
		t.Fatalf("building bar hover = %d, want none after menu opens", state.ui.buildBarHover)
	}
}

// TestBuildingBarHoveredCostFitsMetadataArea verifies bold hover costs keep the compact row.
func TestBuildingBarHoveredCostFitsMetadataArea(t *testing.T) {
	state := newRaidTestState(t)

	if state.ui.costBoldFace.Size != state.ui.costFace.Size {
		t.Fatalf("hovered cost font size = %.1f, want normal cost font size %.1f", state.ui.costBoldFace.Size, state.ui.costFace.Size)
	}

	for _, item := range state.buildingBarItems() {
		costItems := buildingBarCostItems(item.Cost)
		width := state.buildingBarCostWidth(costItems, true)
		available := state.buildingBarMetadataRight() - state.buildingBarMetadataX(item)
		if width > float64(available) {
			t.Fatalf("%s hovered cost width = %.2f, want <= %d", item.Name, width, available)
		}
	}
}

// TestBuildingBarClickDoesNotClearSelection verifies bar clicks are blocked as UI input.
func TestBuildingBarClickDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
	state.gameMap.Home.Tiles[5][homePlotCenter+1].Feature = featureBowTower
	state.Update(clickTileInput(state, homePlotCenter+1, 5))
	item := state.buildingBarItems()[0]

	state.Update(Input{
		CursorX: item.Bounds.X + item.Bounds.W/2,
		CursorY: item.Bounds.Y + item.Bounds.H/2,
		Clicked: true,
	})

	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection kind = %v, want structure", state.selection.kind)
	}
	if state.selection.tile != (tileCoordinate{X: homePlotCenter + 1, Y: 5}) {
		t.Fatalf("selected tile = %+v, want Bow Tower", state.selection.tile)
	}
}

// TestNextRaidButtonAvoidsBuildingBar verifies the Raid control remains accessible beside the bar.
func TestNextRaidButtonAvoidsBuildingBar(t *testing.T) {
	state := newRaidTestState(t)

	bar := state.buildingBarBounds()
	button := state.nextRaidButton()

	if button.X < bar.X+bar.W {
		t.Fatalf("Next Raid X = %d, want at least %d", button.X, bar.X+bar.W)
	}
}

// assertBuildingBarItems verifies visible bar items match stable IDs in order.
func assertBuildingBarItems(t *testing.T, state *State, ids []buildingBarItemID) {
	t.Helper()
	items := state.buildingBarItems()
	if len(items) != len(ids) {
		t.Fatalf("building bar items = %d, want %d", len(items), len(ids))
	}
	for i, id := range ids {
		template, ok := state.buildingTemplateForItemID(id)
		if !ok {
			t.Fatalf("missing template for item ID %d", id)
		}
		assertBuildingBarItem(t, state, items[i], id, template.Name)
		if i > 0 {
			previousBottom := buildingBarItemBottom(items[i-1])
			if items[i].Bounds.Y <= previousBottom {
				t.Fatalf("item %d Y = %d, want below previous item bottom %d", i, items[i].Bounds.Y, previousBottom)
			}
		}
	}
}

// assertBuildingBarItem verifies a bar item has stable bounds within the bar.
func assertBuildingBarItem(t *testing.T, state *State, item buildingBarItem, id buildingBarItemID, name string) {
	t.Helper()
	if item.ID != id {
		t.Fatalf("item ID = %d, want %d", item.ID, id)
	}
	if item.Name != name {
		t.Fatalf("item name = %q, want %q", item.Name, name)
	}
	if item.Bounds.Label != name {
		t.Fatalf("item label = %q, want %q", item.Bounds.Label, name)
	}
	if item.Bounds.W != buildingBarItemSize || item.Bounds.H != buildingBarItemSize {
		t.Fatalf("item size = %dx%d, want %dx%d", item.Bounds.W, item.Bounds.H, buildingBarItemSize, buildingBarItemSize)
	}
	bar := state.buildingBarBounds()
	itemBottom := buildingBarItemBottom(item)
	if !bar.Contains(item.Bounds.X, item.Bounds.Y) ||
		!bar.Contains(item.Bounds.X+item.Bounds.W-1, itemBottom-1) {
		t.Fatalf("item bounds %+v should fit inside bar %+v", item.Bounds, bar)
	}
}

// assertBuildingBarTabInsideBar verifies one category tab stays inside the bar.
func assertBuildingBarTabInsideBar(t *testing.T, state *State, tab buildingBarTab) {
	t.Helper()
	bar := state.buildingBarBounds()
	if !bar.Contains(tab.Bounds.X, tab.Bounds.Y) ||
		!bar.Contains(tab.Bounds.X+tab.Bounds.W-1, tab.Bounds.Y+tab.Bounds.H-1) {
		t.Fatalf("tab bounds %+v should fit inside bar %+v", tab.Bounds, bar)
	}
}

// buildingBarItemBottom returns the bottom edge of one icon and metadata block.
func buildingBarItemBottom(item buildingBarItem) int {
	return item.Bounds.Y + item.Bounds.H
}

// assertCostItem verifies one resource-cost display item.
func assertCostItem(t *testing.T, item buildingBarCostItem, value string, clr color.Color) {
	t.Helper()
	if item.Value != value {
		t.Fatalf("cost item value = %q, want %q", item.Value, value)
	}
	if item.Color != clr {
		t.Fatalf("cost item color = %#v, want %#v", item.Color, clr)
	}
}
