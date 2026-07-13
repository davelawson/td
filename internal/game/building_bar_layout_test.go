package game

import "testing"

// TestBuildingBarMetadataSitsRightOfIcon verifies values use the expanded row space.
func TestBuildingBarMetadataSitsRightOfIcon(t *testing.T) {
	state := newRaidTestState(t)

	for _, category := range buildingBarCategories() {
		state.ui.buildBarCategory = category
		for _, item := range state.buildingBarItems() {
			metadataX := state.buildingBarMetadataX(item)
			if metadataX <= item.Bounds.X+item.Bounds.W {
				t.Fatalf("%s metadata X = %d, want right of icon edge %d", item.Name, metadataX, item.Bounds.X+item.Bounds.W)
			}
			if metadataX >= state.buildingBarMetadataRight() {
				t.Fatalf("%s metadata X = %d, want left of metadata right edge %d", item.Name, metadataX, state.buildingBarMetadataRight())
			}
			if item.Bounds.X != buildingBarPadding {
				t.Fatalf("%s icon X = %d, want left padding %d", item.Name, item.Bounds.X, buildingBarPadding)
			}
		}
	}
}

// TestBuildingBarMetadataClickDoesNotStartDrag verifies values remain non-interactive.
func TestBuildingBarMetadataClickDoesNotStartDrag(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.Update(Input{
		CursorX:   state.buildingBarMetadataX(item),
		CursorY:   item.Bounds.Y + buildingBarCostOffsetY,
		Clicked:   true,
		MouseDown: true,
	})

	if state.buildDrag.active {
		t.Fatal("expected metadata click not to start build drag")
	}
	if state.ui.buildBarHover != -1 {
		t.Fatalf("building bar hover = %d, want none over metadata", state.ui.buildBarHover)
	}
}

// TestBuildingBarIconAlphaTracksConstructionCapacity verifies capacity-only opacity.
func TestBuildingBarIconAlphaTracksConstructionCapacity(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarItems()

	if alpha := state.buildingBarIconAlpha(items[0]); alpha != 1 {
		t.Fatalf("House alpha = %.2f, want 1.00", alpha)
	}
	if alpha := state.buildingBarIconAlpha(items[1]); alpha != 0.70 {
		t.Fatalf("Barracks alpha = %.2f, want 0.70 without Peasants", alpha)
	}
	if alpha := state.buildingBarIconAlpha(items[2]); alpha != 0.70 {
		t.Fatalf("Dorm alpha = %.2f, want 0.70 without Peasant", alpha)
	}

	setAvailablePopulations(state, 0, 0, 2)
	items = state.buildingBarItems()
	if alpha := state.buildingBarIconAlpha(items[1]); alpha != 1 {
		t.Fatalf("Barracks alpha = %.2f, want 1.00 with Peasants", alpha)
	}
	if alpha := state.buildingBarIconAlpha(items[2]); alpha != 1 {
		t.Fatalf("Dorm alpha = %.2f, want 1.00 with Peasant", alpha)
	}

	state.ui.buildBarCategory = buildingBarCategoryDefenses
	items = state.buildingBarItems()
	if alpha := state.buildingBarIconAlpha(items[0]); alpha != 0.70 {
		t.Fatalf("Bow Tower alpha = %.2f, want 0.70 without Soldier", alpha)
	}
}

// TestBuildingBarItemOutlineColorTracksConstructionCapacity verifies green/red outlines.
func TestBuildingBarItemOutlineColorTracksConstructionCapacity(t *testing.T) {
	state := newRaidTestState(t)
	items := state.buildingBarItems()

	if got := state.buildingBarItemOutlineColor(items[0]); got != colors.buildable {
		t.Fatalf("House outline = %#v, want buildable green %#v", got, colors.buildable)
	}
	if got := state.buildingBarItemOutlineColor(items[1]); got != colors.buildBlocked {
		t.Fatalf("Barracks outline = %#v, want blocked red %#v", got, colors.buildBlocked)
	}
	if got := state.buildingBarItemOutlineColor(items[2]); got != colors.buildBlocked {
		t.Fatalf("Dorm outline = %#v, want blocked red %#v", got, colors.buildBlocked)
	}

	setAvailablePopulations(state, 0, 0, 2)
	items = state.buildingBarItems()
	if got := state.buildingBarItemOutlineColor(items[1]); got != colors.buildable {
		t.Fatalf("Barracks outline = %#v, want buildable green with Peasants %#v", got, colors.buildable)
	}
	if got := state.buildingBarItemOutlineColor(items[2]); got != colors.buildable {
		t.Fatalf("Dorm outline = %#v, want buildable green with Peasant %#v", got, colors.buildable)
	}

	state.ui.buildBarCategory = buildingBarCategoryDefenses
	items = state.buildingBarItems()
	if got := state.buildingBarItemOutlineColor(items[0]); got != colors.buildBlocked {
		t.Fatalf("Bow Tower outline = %#v, want blocked red without Soldier %#v", got, colors.buildBlocked)
	}
	setAvailablePopulations(state, 0, 1, 0)
	items = state.buildingBarItems()
	if got := state.buildingBarItemOutlineColor(items[0]); got != colors.buildable {
		t.Fatalf("Bow Tower outline = %#v, want buildable green with Soldier %#v", got, colors.buildable)
	}
}
