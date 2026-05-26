package game

import "testing"

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

	if len(items) != 2 {
		t.Fatalf("building bar items = %d, want 2", len(items))
	}
	assertBuildingBarItem(t, state, items[0], "Bow Tower")
	assertBuildingBarItem(t, state, items[1], "Flame Bolt Tower")
	if items[0].Sprite != state.structureCatalog.BowTower.Sprite {
		t.Fatal("expected first item to use Bow Tower sprite")
	}
	if items[1].Sprite != state.structureCatalog.FlameBoltTower.Sprite {
		t.Fatal("expected second item to use Flame Bolt Tower sprite")
	}
	if items[1].Bounds.Y <= items[0].Bounds.Y {
		t.Fatalf("second item Y = %d, want below first item Y %d", items[1].Bounds.Y, items[0].Bounds.Y)
	}
}

// TestBuildingBarClickDoesNotClearSelection verifies bar clicks are blocked as UI input.
func TestBuildingBarClickDoesNotClearSelection(t *testing.T) {
	state := newRaidTestState(t)
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
	if !bar.Contains(item.Bounds.X, item.Bounds.Y) ||
		!bar.Contains(item.Bounds.X+item.Bounds.W-1, item.Bounds.Y+item.Bounds.H-1) {
		t.Fatalf("item bounds %+v should fit inside bar %+v", item.Bounds, bar)
	}
}
