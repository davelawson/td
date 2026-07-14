package ui

import "testing"

// TestBuildingTooltipContentForAllBuildingMenuItems verifies every choice has detailed hover text.
func TestBuildingTooltipContentForAllBuildingMenuItems(t *testing.T) {
	model := testBuildingBarModel()
	cases := []struct {
		category BuildingBarCategory
		index    int
		title    string
		lines    []string
	}{
		{BuildingBarCategoryHousing, 0, "House", []string{"Shelters new Peasants for the Domain.", "Cost: 20 Wood", "Staffing: None", "Effect: +2 Peasants"}},
		{BuildingBarCategoryHousing, 1, "Barracks", []string{"Trains Peasants into Soldiers for staffed defenses.", "Cost: 10 Wood, 10 Stone", "Staffing: None", "Effect: -2 Peasants, +2 Soldiers"}},
		{BuildingBarCategoryHousing, 2, "Dorm", []string{"Houses Peasants studying to become Apprentices.", "Cost: 10 Wood, 10 Stone", "Staffing: None", "Effect: -1 Peasant, +1 Apprentice"}},
		{BuildingBarCategoryEconomic, 0, "Woodcutter", []string{"Assigns a Peasant to bring in Wood during Labour.", "Cost: 10 Wood", "Staffing: 1 Peasant", "Production: +10 Wood during each Labour phase"}},
		{BuildingBarCategoryEconomic, 1, "Stone Quarry", []string{"Assigns a Peasant to quarry Stone during Labour.", "Cost: 10 Wood, 10 Stone", "Staffing: 1 Peasant", "Production: +10 Stone during each Labour phase"}},
		{BuildingBarCategoryEconomic, 2, "Iron Mine", []string{"Assigns a Peasant to extract Iron during Labour.", "Cost: 10 Wood, 10 Stone, 10 Iron", "Staffing: 1 Peasant", "Production: +10 Iron during each Labour phase"}},
		{BuildingBarCategoryDefenses, 0, "Bow Tower", []string{"A staffed archer tower that fires quick arrows.", "Cost: 20 Wood, 10 Stone", "Staffing: 1 Soldier", "Range: 3.0 Tiles", "Damage: 10", "Fire: every 1.0s", "Projectile: 9.0 Tiles/s"}},
		{BuildingBarCategoryDefenses, 1, "Flame Bolt Tower", []string{"An apprentice-staffed tower that hurls focused fire.", "Cost: 30 Stone, 20 Iron", "Staffing: 1 Apprentice", "Range: 2.5 Tiles", "Damage: 20", "Fire: every 1.5s", "Projectile: 7.0 Tiles/s"}},
		{BuildingBarCategoryDefenses, 2, "Catapult Tower", []string{"A heavy crewed tower that crushes enemies in one Tile.", "Cost: 40 Wood, 60 Stone, 25 Iron", "Staffing: 1 Soldier, 1 Peasant", "Range: 5.0 Tiles", "Damage: 30", "Fire: every 6.0s", "Projectile: 3.0 Tiles/s", "Area: damages every enemy in the target Tile"}},
	}

	for _, tc := range cases {
		model.SelectedCategory = tc.category
		model.HoveredItem = tc.index
		tooltip, ok := hoveredBuildingTooltip(1920, 1080, testBuildingBarTop, model)
		if !ok {
			t.Fatalf("%s did not produce a tooltip", tc.title)
		}
		if tooltip.Title != tc.title {
			t.Fatalf("title = %q, want %q", tooltip.Title, tc.title)
		}
		values := buildingTooltipLineValues(tooltip)
		for _, want := range tc.lines {
			if !containsBuildingTooltipLine(values, want) {
				t.Fatalf("%s lines = %+v, want %q", tc.title, values, want)
			}
		}
	}
}

// TestBuildingTooltipRequiresValidIconHover verifies non-icon state stays tooltip-free.
func TestBuildingTooltipRequiresValidIconHover(t *testing.T) {
	model := testBuildingBarModel()
	for _, hover := range []int{-1, 3, 99} {
		model.HoveredItem = hover
		if _, ok := hoveredBuildingTooltip(1920, 1080, testBuildingBarTop, model); ok {
			t.Fatalf("hover index %d unexpectedly produced a tooltip", hover)
		}
	}
}

// TestBuildingTooltipBoundsStayInDrawableArea verifies tooltips clamp to screen bounds.
func TestBuildingTooltipBoundsStayInDrawableArea(t *testing.T) {
	model := testBuildingBarModel()
	model.SelectedCategory = BuildingBarCategoryDefenses
	model.HoveredItem = 2
	tooltip, ok := hoveredBuildingTooltip(800, 500, testBuildingBarTop, model)
	if !ok {
		t.Fatal("expected Catapult Tower tooltip")
	}
	if tooltip.Bounds.X < BuildingBarWidth || tooltip.Bounds.X+tooltip.Bounds.W > 800 {
		t.Fatalf("horizontal bounds = %+v, want beside bar and inside screen", tooltip.Bounds)
	}
	if tooltip.Bounds.Y < testBuildingBarTop || tooltip.Bounds.Y+tooltip.Bounds.H > 500 {
		t.Fatalf("vertical bounds = %+v, want drawable area", tooltip.Bounds)
	}
}

// buildingTooltipLineValues returns tooltip text for assertions.
func buildingTooltipLineValues(tooltip buildingTooltip) []string {
	values := make([]string, 0, len(tooltip.Lines))
	for _, line := range tooltip.Lines {
		values = append(values, line.Value)
	}
	return values
}

// containsBuildingTooltipLine reports whether a tooltip contains an exact line.
func containsBuildingTooltipLine(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
