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
	if bounds.H != state.ui.height-topBarHeight {
		t.Fatalf("bar height = %d, want %d", bounds.H, state.ui.height-topBarHeight)
	}
}

// TestBuildingBarItemsExposeTowerIcons verifies the visual build choices use the tower templates.
func TestBuildingBarItemsExposeTowerIcons(t *testing.T) {
	state := newRaidTestState(t)

	items := state.buildingBarItems()

	if len(items) != 3 {
		t.Fatalf("building bar items = %d, want 3", len(items))
	}
	assertBuildingBarItem(t, state, items[0], "Bow Tower")
	assertBuildingBarItem(t, state, items[1], "Flame Bolt Tower")
	assertBuildingBarItem(t, state, items[2], "Catapult Tower")
	if items[0].Sprite != state.structureCatalog.BowTower.Sprite {
		t.Fatal("expected first item to use Bow Tower sprite")
	}
	if items[0].Cost != (Resources{Wood: 30, Stone: 10, Metal: 10}) {
		t.Fatalf("Bow Tower cost = %+v, want 30 wood 10 stone 10 metal", items[0].Cost)
	}
	if items[0].Staffing != (StaffingRequirements{Soldiers: 1}) {
		t.Fatalf("Bow Tower staffing = %+v, want 1 Soldier", items[0].Staffing)
	}
	if items[1].Sprite != state.structureCatalog.FlameBoltTower.Sprite {
		t.Fatal("expected second item to use Flame Bolt Tower sprite")
	}
	if items[1].Cost != (Resources{Stone: 30, Metal: 20}) {
		t.Fatalf("Flame Bolt Tower cost = %+v, want 30 stone 20 metal", items[1].Cost)
	}
	if items[1].Staffing != (StaffingRequirements{Apprentices: 1}) {
		t.Fatalf("Flame Bolt Tower staffing = %+v, want 1 Apprentice", items[1].Staffing)
	}
	if items[2].Sprite != state.structureCatalog.CatapultTower.Sprite {
		t.Fatal("expected third item to use Catapult Tower sprite")
	}
	if items[2].Cost != (Resources{Wood: 40, Stone: 60, Metal: 25}) {
		t.Fatalf("Catapult Tower cost = %+v, want 40 wood 60 stone 25 metal", items[2].Cost)
	}
	if items[2].Staffing != (StaffingRequirements{Soldiers: 1, Peasants: 2}) {
		t.Fatalf("Catapult Tower staffing = %+v, want 1 Soldier and 2 Peasants", items[2].Staffing)
	}
	firstBlockBottom := buildingBarItemBottom(items[0])
	if items[1].Bounds.Y <= firstBlockBottom {
		t.Fatalf("second item Y = %d, want below first item cost bottom %d", items[1].Bounds.Y, firstBlockBottom)
	}
	secondBlockBottom := buildingBarItemBottom(items[1])
	if items[2].Bounds.Y <= secondBlockBottom {
		t.Fatalf("third item Y = %d, want below second item cost bottom %d", items[2].Bounds.Y, secondBlockBottom)
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

// TestBuildingBarCostItems verifies non-zero tower costs render in resource order.
func TestBuildingBarCostItems(t *testing.T) {
	items := buildingBarCostItems(Resources{Wood: 30, Stone: 10, Metal: 10})
	if len(items) != 3 {
		t.Fatalf("cost items = %d, want 3", len(items))
	}
	assertCostItem(t, items[0], "30", colors.resourceWood)
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

	state.updateBuildingBarHover(Input{
		CursorX: items[2].Bounds.X + items[2].Bounds.W/2,
		CursorY: items[2].Bounds.Y + items[2].Bounds.H/2,
	})
	if state.ui.buildBarHover != 2 {
		t.Fatalf("building bar hover = %d, want third item", state.ui.buildBarHover)
	}
}

// TestBuildingBarHighlightRequiresResourcesAndStaff verifies all construction inputs gate emphasis.
func TestBuildingBarHighlightRequiresResourcesAndStaff(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarItems()

	state.ui.buildBarHover = 0
	if state.buildingBarItemHighlighted(0, items[0]) {
		t.Fatal("expected zero Soldiers to suppress Bow Tower highlight")
	}
	setAvailablePopulations(state, 1, 1, 2)
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
	state.status.populations.peasants.available = 1
	if state.buildingBarItemHighlighted(2, items[2]) {
		t.Fatal("expected one missing Peasant to suppress Catapult Tower highlight")
	}
}

// TestBuildingBarHoverIgnoresCostsAndEmptyBar verifies costs do not trigger icon hover.
func TestBuildingBarHoverIgnoresCostsAndEmptyBar(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.updateBuildingBarHover(Input{
		CursorX: item.Bounds.X + item.Bounds.W/2,
		CursorY: item.Bounds.Y + item.Bounds.H + buildingBarCostGap + buildingBarCostTextHeight/2,
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

// TestBuildingBarHoveredCostFitsIcon verifies bold hover costs keep the compact row.
func TestBuildingBarHoveredCostFitsIcon(t *testing.T) {
	state := newRaidTestState(t)

	if state.ui.costBoldFace.Size != state.ui.costFace.Size {
		t.Fatalf("hovered cost font size = %.1f, want normal cost font size %.1f", state.ui.costBoldFace.Size, state.ui.costFace.Size)
	}

	for _, item := range state.buildingBarItems() {
		costItems := buildingBarCostItems(item.Cost)
		width := state.buildingBarCostWidth(costItems, true)
		if width > float64(item.Bounds.W) {
			t.Fatalf("%s hovered cost width = %.2f, want <= %d", item.Name, width, item.Bounds.W)
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

// assertBuildingBarItem verifies a bar item has stable bounds within the bar.
func assertBuildingBarItem(t *testing.T, state *State, item buildingBarItem, name string) {
	t.Helper()
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

// buildingBarItemBottom returns the bottom edge of one icon and metadata block.
func buildingBarItemBottom(item buildingBarItem) int {
	return item.Bounds.Y + item.Bounds.H +
		buildingBarCostGap + buildingBarCostTextHeight +
		buildingBarStaffingGap + buildingBarStaffingHeight
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
