package ui

import (
	"bytes"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const testBuildingBarTop = 86

// TestBuildingBarBoundsFillPlayableLeftEdge verifies the widget occupies the expected screen edge.
func TestBuildingBarBoundsFillPlayableLeftEdge(t *testing.T) {
	bounds := BuildingBarBounds(testBuildingBarTop, 1080)
	if bounds != (Button[int]{X: 0, Y: testBuildingBarTop, W: 260, H: 1080 - testBuildingBarTop}) {
		t.Fatalf("bounds = %+v, want playable left edge", bounds)
	}
	if !BuildingBarContains(testBuildingBarTop, 1080, 1, testBuildingBarTop+1) {
		t.Fatal("expected point inside building bar")
	}
}

// TestBuildingBarCategoriesAndActionsHaveStableOrder verifies UI-owned grouping and ordering.
func TestBuildingBarCategoriesAndActionsHaveStableOrder(t *testing.T) {
	tabs := buildingBarTabs(testBuildingBarTop)
	wantCategories := []BuildingBarCategory{BuildingBarCategoryDefenses, BuildingBarCategoryEconomic, BuildingBarCategoryHousing}
	wantLabels := []string{"Defenses", "Economic", "Housing"}
	if len(tabs) != len(wantCategories) {
		t.Fatalf("tabs = %d, want %d", len(tabs), len(wantCategories))
	}
	for index, tab := range tabs {
		if tab.Category != wantCategories[index] || tab.Label != wantLabels[index] {
			t.Fatalf("tab %d = %+v, want %v %q", index, tab, wantCategories[index], wantLabels[index])
		}
		if got := BuildingBarCategoryAt(testBuildingBarTop, tab.Bounds.X+1, tab.Bounds.Y+1); got != tab.Category {
			t.Fatalf("category hit = %v, want %v", got, tab.Category)
		}
	}

	cases := []struct {
		category BuildingBarCategory
		actions  []BuildingBarAction
	}{
		{BuildingBarCategoryHousing, []BuildingBarAction{BuildingBarHouse, BuildingBarBarracks, BuildingBarDorm}},
		{BuildingBarCategoryEconomic, []BuildingBarAction{BuildingBarWoodcutter, BuildingBarStoneQuarry, BuildingBarIronMine}},
		{BuildingBarCategoryDefenses, []BuildingBarAction{BuildingBarBowTower, BuildingBarFlameBoltTower, BuildingBarCatapultTower}},
	}
	model := testBuildingBarModel()
	for _, tc := range cases {
		model.SelectedCategory = tc.category
		items := buildingBarLayoutItems(testBuildingBarTop, model)
		if len(items) != len(tc.actions) {
			t.Fatalf("category %v items = %d, want %d", tc.category, len(items), len(tc.actions))
		}
		for index, action := range tc.actions {
			if items[index].Action != action {
				t.Fatalf("category %v item %d action = %v, want %v", tc.category, index, items[index].Action, action)
			}
			if BuildingBarCategoryForAction(action) != tc.category {
				t.Fatalf("action %v category = %v, want %v", action, BuildingBarCategoryForAction(action), tc.category)
			}
		}
	}
}

// TestBuildingBarLayoutKeepsItemsInsideBar verifies stable icon geometry and spacing.
func TestBuildingBarLayoutKeepsItemsInsideBar(t *testing.T) {
	model := testBuildingBarModel()
	bar := BuildingBarBounds(testBuildingBarTop, 1080)
	for _, category := range []BuildingBarCategory{BuildingBarCategoryHousing, BuildingBarCategoryEconomic, BuildingBarCategoryDefenses} {
		model.SelectedCategory = category
		items := buildingBarLayoutItems(testBuildingBarTop, model)
		for index, item := range items {
			if item.Bounds.X != buildingBarPadding || item.Bounds.W != buildingBarItemSize || item.Bounds.H != buildingBarItemSize {
				t.Fatalf("item bounds = %+v, want padded %dx%d icon", item.Bounds, buildingBarItemSize, buildingBarItemSize)
			}
			if !bar.Contains(item.Bounds.X, item.Bounds.Y) || !bar.Contains(item.Bounds.X+item.Bounds.W-1, item.Bounds.Y+item.Bounds.H-1) {
				t.Fatalf("item bounds %+v should fit inside bar %+v", item.Bounds, bar)
			}
			if index > 0 && item.Bounds.Y <= items[index-1].Bounds.Y+items[index-1].Bounds.H {
				t.Fatalf("item %d overlaps previous item", index)
			}
		}
	}
}

// TestBuildingBarHitTestingUsesIconBounds verifies metadata and empty bar space are non-interactive.
func TestBuildingBarHitTestingUsesIconBounds(t *testing.T) {
	model := testBuildingBarModel()
	items := buildingBarLayoutItems(testBuildingBarTop, model)
	first := items[0]
	if got := BuildingBarItemIndexAt(testBuildingBarTop, model, first.Bounds.X+1, first.Bounds.Y+1); got != 0 {
		t.Fatalf("icon hit index = %d, want 0", got)
	}
	if got := BuildingBarItemIndexAt(testBuildingBarTop, model, buildingBarMetadataX(first), first.Bounds.Y+buildingBarCostOffsetY); got != -1 {
		t.Fatalf("metadata hit index = %d, want -1", got)
	}
	if got := BuildingBarItemIndexAt(testBuildingBarTop, model, 1, testBuildingBarTop+1); got != -1 {
		t.Fatalf("empty bar hit index = %d, want -1", got)
	}
}

// TestBuildingBarPopulationMetadataUsesStableRoleOrder verifies staffing and conversions.
func TestBuildingBarPopulationMetadataUsesStableRoleOrder(t *testing.T) {
	icons := BuildingBarIcons{
		Apprentice: ebiten.NewImage(1, 1),
		Soldier:    ebiten.NewImage(1, 1),
		Peasant:    ebiten.NewImage(1, 1),
	}
	staffed := buildingBarPopulationMetadataItems(icons, BuildingBarItem{
		Staffing: PopulationAmounts{Apprentices: 3, Soldiers: 1, Peasants: 2},
	})
	if len(staffed) != 3 || staffed[0].Count != 3 || staffed[0].Sprite != icons.Apprentice ||
		staffed[1].Count != 1 || staffed[1].Sprite != icons.Soldier || staffed[2].Count != 2 || staffed[2].Sprite != icons.Peasant {
		t.Fatalf("staffing metadata = %+v, want Apprentice, Soldier, Peasant order", staffed)
	}

	conversion := buildingBarPopulationMetadataItems(icons, BuildingBarItem{
		PopulationCost:  PopulationAmounts{Peasants: 2},
		PopulationGrant: PopulationAmounts{Soldiers: 2},
	})
	if len(conversion) != 2 || conversion[0].Value != "-2" || conversion[0].Sprite != icons.Peasant ||
		conversion[1].Value != "+2" || conversion[1].Sprite != icons.Soldier {
		t.Fatalf("conversion metadata = %+v, want -2 Peasants then +2 Soldiers", conversion)
	}
}

// TestBuildingBarCostsUseStableResourceOrder verifies compact metadata values and colors.
func TestBuildingBarCostsUseStableResourceOrder(t *testing.T) {
	items := buildingBarCostItems(ResourceAmounts{Wood: 40, Stone: 60, Metal: 25})
	if len(items) != 3 || items[0].Value != "40" || items[0].Color != ResourceWood ||
		items[1].Value != "60" || items[1].Color != ResourceStone ||
		items[2].Value != "25" || items[2].Color != ResourceMetal {
		t.Fatalf("cost items = %+v, want Wood, Stone, Metal order", items)
	}
}

// TestBuildingBarAvailabilityControlsVisualState verifies presentation consumes the game flag only.
func TestBuildingBarAvailabilityControlsVisualState(t *testing.T) {
	if alpha := buildingBarIconAlpha(true); alpha != 1 {
		t.Fatalf("buildable alpha = %.2f, want 1", alpha)
	}
	if alpha := buildingBarIconAlpha(false); alpha != 0.70 {
		t.Fatalf("blocked alpha = %.2f, want 0.70", alpha)
	}
	if buildingBarOutlineColor(true) != buildingBarBuildableColor || buildingBarOutlineColor(false) != buildingBarBlockedColor {
		t.Fatal("availability outlines do not match UI palette")
	}
}

// TestBuildingBarHoveredCostFitsMetadataArea verifies bold values fit the compact row.
func TestBuildingBarHoveredCostFitsMetadataArea(t *testing.T) {
	face := testBuildingBarFace(t)
	model := testBuildingBarModel()
	for _, item := range buildingBarLayoutItems(testBuildingBarTop, model) {
		width := buildingBarCostWidth(face, buildingBarCostItems(item.Cost))
		available := buildingBarMetadataRight() - buildingBarMetadataX(item)
		if width > float64(available) {
			t.Fatalf("%s cost width = %.2f, want <= %d", item.Name, width, available)
		}
	}
}

// testBuildingBarModel returns complete presentation facts for UI unit tests.
func testBuildingBarModel() BuildingBarModel {
	items := []BuildingBarItem{
		{Action: BuildingBarHouse, Name: "House", Description: "Shelters new Peasants for the Domain.", Cost: ResourceAmounts{Wood: 20}, PopulationGrant: PopulationAmounts{Peasants: 2}, Buildable: true},
		{Action: BuildingBarBarracks, Name: "Barracks", Description: "Trains Peasants into Soldiers for staffed defenses.", Cost: ResourceAmounts{Wood: 10, Stone: 10}, PopulationCost: PopulationAmounts{Peasants: 2}, PopulationGrant: PopulationAmounts{Soldiers: 2}},
		{Action: BuildingBarDorm, Name: "Dorm", Description: "Houses Peasants studying to become Apprentices.", Cost: ResourceAmounts{Wood: 10, Stone: 10}, PopulationCost: PopulationAmounts{Peasants: 1}, PopulationGrant: PopulationAmounts{Apprentices: 1}},
		{Action: BuildingBarWoodcutter, Name: "Woodcutter", Description: "Assigns a Peasant to bring in Wood during Labour.", Cost: ResourceAmounts{Wood: 10}, Staffing: PopulationAmounts{Peasants: 1}, ResourceYield: ResourceAmounts{Wood: 10}},
		{Action: BuildingBarStoneQuarry, Name: "Stone Quarry", Description: "Assigns a Peasant to quarry Stone during Labour.", Cost: ResourceAmounts{Wood: 10, Stone: 10}, Staffing: PopulationAmounts{Peasants: 1}, ResourceYield: ResourceAmounts{Stone: 10}},
		{Action: BuildingBarIronMine, Name: "Iron Mine", Description: "Assigns a Peasant to extract Metal during Labour.", Cost: ResourceAmounts{Wood: 10, Stone: 10, Metal: 10}, Staffing: PopulationAmounts{Peasants: 1}, ResourceYield: ResourceAmounts{Metal: 10}},
		{Action: BuildingBarBowTower, Name: "Bow Tower", Description: "A staffed archer tower that fires quick arrows.", Cost: ResourceAmounts{Wood: 20, Stone: 10}, Staffing: PopulationAmounts{Soldiers: 1}, RangeTiles: 3, Damage: 10, FireIntervalSeconds: 1, ProjectileSpeedTilesPerSecond: 9},
		{Action: BuildingBarFlameBoltTower, Name: "Flame Bolt Tower", Description: "An apprentice-staffed tower that hurls focused fire.", Cost: ResourceAmounts{Stone: 30, Metal: 20}, Staffing: PopulationAmounts{Apprentices: 1}, RangeTiles: 2.5, Damage: 20, FireIntervalSeconds: 1.5, ProjectileSpeedTilesPerSecond: 7},
		{Action: BuildingBarCatapultTower, Name: "Catapult Tower", Description: "A heavy crewed tower that crushes enemies in one Tile.", Cost: ResourceAmounts{Wood: 40, Stone: 60, Metal: 25}, Staffing: PopulationAmounts{Soldiers: 1, Peasants: 1}, RangeTiles: 5, Damage: 30, FireIntervalSeconds: 6, ProjectileSpeedTilesPerSecond: 3, DamageAllEnemiesInTargetTile: true},
	}
	return BuildingBarModel{Items: items, SelectedCategory: BuildingBarCategoryHousing, HoveredItem: -1, HoveredCategory: BuildingBarNoCategory}
}

// testBuildingBarFace returns a real font face for measurement tests.
func testBuildingBarFace(t *testing.T) *text.GoTextFace {
	t.Helper()
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("font source: %v", err)
	}
	return &text.GoTextFace{Source: source, Size: GameCostFontSize}
}
