package game

import "testing"

// TestBuildingTooltipContentForAllBuildingMenuItems verifies every tab entry has detailed hover text.
func TestBuildingTooltipContentForAllBuildingMenuItems(t *testing.T) {
	state := newRaidTestState(t)
	cases := []struct {
		category buildingBarCategory
		index    int
		title    string
		lines    []string
	}{
		{
			category: buildingBarCategoryHousing,
			index:    0,
			title:    "House",
			lines: []string{
				"Shelters new Peasants for the Domain.",
				"Cost: 20 Wood",
				"Staffing: None",
				"Effect: +2 Peasants",
			},
		},
		{
			category: buildingBarCategoryHousing,
			index:    1,
			title:    "Barracks",
			lines: []string{
				"Trains Peasants into Soldiers for staffed defenses.",
				"Cost: 10 Wood, 10 Stone",
				"Staffing: None",
				"Effect: -2 Peasants, +2 Soldiers",
			},
		},
		{
			category: buildingBarCategoryHousing,
			index:    2,
			title:    "Dorm",
			lines: []string{
				"Houses Peasants studying to become Apprentices.",
				"Cost: 10 Wood, 10 Stone",
				"Staffing: None",
				"Effect: -1 Peasant, +1 Apprentice",
			},
		},
		{
			category: buildingBarCategoryEconomic,
			index:    0,
			title:    "Woodcutter",
			lines: []string{
				"Assigns a Peasant to bring in Wood after defeated Raids.",
				"Cost: 10 Wood",
				"Staffing: 1 Peasant",
				"Production: +10 Wood after each defeated Raid",
			},
		},
		{
			category: buildingBarCategoryEconomic,
			index:    1,
			title:    "Stone Quarry",
			lines: []string{
				"Assigns a Peasant to quarry Stone after defeated Raids.",
				"Cost: 10 Wood, 10 Stone",
				"Staffing: 1 Peasant",
				"Production: +10 Stone after each defeated Raid",
			},
		},
		{
			category: buildingBarCategoryEconomic,
			index:    2,
			title:    "Iron Mine",
			lines: []string{
				"Assigns a Peasant to extract Metal after defeated Raids.",
				"Cost: 10 Wood, 10 Stone, 10 Metal",
				"Staffing: 1 Peasant",
				"Production: +10 Metal after each defeated Raid",
			},
		},
		{
			category: buildingBarCategoryDefenses,
			index:    0,
			title:    "Bow Tower",
			lines: []string{
				"A staffed archer tower that fires quick arrows.",
				"Cost: 30 Wood, 10 Stone, 10 Metal",
				"Staffing: 1 Soldier",
				"Range: 3.0 Tiles",
				"Damage: 10",
				"Fire: every 1.0s",
				"Projectile: 9.0 Tiles/s",
			},
		},
		{
			category: buildingBarCategoryDefenses,
			index:    1,
			title:    "Flame Bolt Tower",
			lines: []string{
				"An apprentice-staffed tower that hurls focused fire.",
				"Cost: 30 Stone, 20 Metal",
				"Staffing: 1 Apprentice",
				"Range: 2.5 Tiles",
				"Damage: 20",
				"Fire: every 1.5s",
				"Projectile: 7.0 Tiles/s",
			},
		},
		{
			category: buildingBarCategoryDefenses,
			index:    2,
			title:    "Catapult Tower",
			lines: []string{
				"A heavy crewed tower that crushes enemies in one Tile.",
				"Cost: 40 Wood, 60 Stone, 25 Metal",
				"Staffing: 1 Soldier, 2 Peasants",
				"Range: 5.0 Tiles",
				"Damage: 75",
				"Fire: every 3.0s",
				"Projectile: 3.0 Tiles/s",
				"Area: damages every enemy in the target Tile",
			},
		},
	}

	for _, tc := range cases {
		state.ui.buildBarCategory = tc.category
		item := state.buildingBarItems()[tc.index]
		state.updateBuildingBarHover(Input{
			CursorX: item.Bounds.X + item.Bounds.W/2,
			CursorY: item.Bounds.Y + item.Bounds.H/2,
		})

		tooltip, ok := state.hoveredBuildingTooltip()
		if !ok {
			t.Fatalf("%s did not produce a tooltip", tc.title)
		}
		if tooltip.Title != tc.title {
			t.Fatalf("tooltip title = %q, want %q", tooltip.Title, tc.title)
		}
		values := tooltipLineValues(tooltip)
		for _, want := range tc.lines {
			if !containsString(values, want) {
				t.Fatalf("%s tooltip lines = %+v, want line %q", tc.title, values, want)
			}
		}
	}
}

// TestBuildingTooltipRequiresIconHover verifies non-icon building-bar areas stay tooltip-free.
func TestBuildingTooltipRequiresIconHover(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.updateBuildingBarHover(Input{
		CursorX: state.buildingBarMetadataX(item),
		CursorY: item.Bounds.Y + buildingBarCostOffsetY,
	})
	if _, ok := state.hoveredBuildingTooltip(); ok {
		t.Fatal("expected no tooltip over building metadata")
	}

	tab := state.buildingBarTabs()[0]
	state.updateBuildingBarHover(Input{
		CursorX: tab.Bounds.X + tab.Bounds.W/2,
		CursorY: tab.Bounds.Y + tab.Bounds.H/2,
	})
	if _, ok := state.hoveredBuildingTooltip(); ok {
		t.Fatal("expected no tooltip over building category tab")
	}

	state.updateBuildingBarHover(Input{
		CursorX: 1,
		CursorY: topBarHeight + 1,
	})
	if _, ok := state.hoveredBuildingTooltip(); ok {
		t.Fatal("expected no tooltip over empty building bar")
	}
}

// TestBuildingTooltipBoundsStayInDrawableArea verifies the tooltip clamps to screen bounds.
func TestBuildingTooltipBoundsStayInDrawableArea(t *testing.T) {
	state := newRaidTestState(t)
	state.ui.buildBarCategory = buildingBarCategoryDefenses
	item := state.buildingBarItems()[2]
	state.updateBuildingBarHover(Input{
		CursorX: item.Bounds.X + item.Bounds.W/2,
		CursorY: item.Bounds.Y + item.Bounds.H/2,
	})

	tooltip, ok := state.hoveredBuildingTooltip()
	if !ok {
		t.Fatal("expected Catapult Tower tooltip")
	}
	bar := state.buildingBarBounds()
	if tooltip.Bounds.X < bar.X+bar.W {
		t.Fatalf("tooltip X = %d, want at or right of bar edge %d", tooltip.Bounds.X, bar.X+bar.W)
	}
	if tooltip.Bounds.X+tooltip.Bounds.W > state.ui.width {
		t.Fatalf("tooltip right = %d, want <= screen width %d", tooltip.Bounds.X+tooltip.Bounds.W, state.ui.width)
	}
	if tooltip.Bounds.Y < topBarHeight {
		t.Fatalf("tooltip Y = %d, want >= top bar height %d", tooltip.Bounds.Y, topBarHeight)
	}
	if tooltip.Bounds.Y+tooltip.Bounds.H > state.ui.height {
		t.Fatalf("tooltip bottom = %d, want <= screen height %d", tooltip.Bounds.Y+tooltip.Bounds.H, state.ui.height)
	}
}

// TestBuildingTooltipHidesDuringDrag verifies dragging a building suppresses hover help.
func TestBuildingTooltipHidesDuringDrag(t *testing.T) {
	state := newRaidTestState(t)
	item := state.buildingBarItems()[0]

	state.Update(Input{
		CursorX:   item.Bounds.X + item.Bounds.W/2,
		CursorY:   item.Bounds.Y + item.Bounds.H/2,
		Clicked:   true,
		MouseDown: true,
	})

	if !state.buildDrag.active {
		t.Fatal("expected House drag to start")
	}
	if _, ok := state.hoveredBuildingTooltip(); ok {
		t.Fatal("expected no tooltip while dragging a building")
	}
}

func tooltipLineValues(tooltip buildingTooltip) []string {
	values := make([]string, 0, len(tooltip.Lines))
	for _, line := range tooltip.Lines {
		values = append(values, line.Value)
	}
	return values
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
