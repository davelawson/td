package ui

import "testing"

// TestSelectionPanelRows verifies selected-object facts are formatted for display.
func TestSelectionPanelRows(t *testing.T) {
	tests := []struct {
		name string
		data SelectionPanelData
		rows map[string]string
	}{
		{
			name: "raider",
			data: SelectionPanelData{
				Kind:                SelectionPanelRaider,
				Name:                "Skeleton Sword-and-Shield",
				Health:              25,
				MaxHealth:           50,
				SpeedTilesPerSecond: 1.0,
				SanctumDamage:       1,
			},
			rows: map[string]string{
				"Raider Type":      "Skeleton Sword-and-Shield",
				"Health":           "25",
				"Max Health":       "50",
				"Health Remaining": "50%",
				"Speed":            "1.0 tiles/s",
				"Sanctum Damage":   "1",
			},
		},
		{
			name: "structure",
			data: SelectionPanelData{Kind: SelectionPanelStructure, Name: "Sanctum"},
			rows: map[string]string{"Structure": "Sanctum"},
		},
		{
			name: "population building",
			data: SelectionPanelData{
				Kind:            SelectionPanelPopulationBuilding,
				Name:            "Barracks",
				Cost:            ResourceAmounts{Wood: 10, Stone: 10},
				PopulationCost:  PopulationAmounts{Peasants: 2},
				PopulationGrant: PopulationAmounts{Soldiers: 2},
			},
			rows: map[string]string{
				"Structure":         "Barracks",
				"Cost":              "10 Wood, 10 Stone",
				"Consumes Peasants": "2",
				"Grants Soldiers":   "2",
			},
		},
		{
			name: "economic building",
			data: SelectionPanelData{
				Kind:          SelectionPanelEconomicBuilding,
				Name:          "Stone Quarry",
				Cost:          ResourceAmounts{Wood: 10, Stone: 10},
				Staffing:      PopulationAmounts{Peasants: 1},
				ResourceYield: ResourceAmounts{Stone: 10},
			},
			rows: map[string]string{
				"Structure":         "Stone Quarry",
				"Cost":              "10 Wood, 10 Stone",
				"Required Peasants": "1",
				"Produces":          "10 Stone after each Raid",
			},
		},
		{
			name: "tower",
			data: SelectionPanelData{
				Kind:                SelectionPanelTower,
				Name:                "Catapult Tower",
				RangeTiles:          5.0,
				FireIntervalSeconds: 6.0,
				Damage:              30,
				Staffing:            PopulationAmounts{Soldiers: 1, Peasants: 2},
			},
			rows: map[string]string{
				"Tower Type":        "Catapult Tower",
				"Range":             "5.0 tiles",
				"Attack Speed":      "every 6.0s",
				"Damage":            "30",
				"Required Soldiers": "1",
				"Required Peasants": "2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rows, ok := selectionPanelRows(test.data)
			if !ok {
				t.Fatal("expected selection panel rows")
			}
			for label, value := range test.rows {
				assertSelectionPanelRow(t, rows, label, value)
			}
		})
	}
}

// TestSelectionPanelFormatting verifies resource and health formatting edge cases.
func TestSelectionPanelFormatting(t *testing.T) {
	resourceTests := []struct {
		cost ResourceAmounts
		want string
	}{
		{cost: ResourceAmounts{}, want: "Free"},
		{cost: ResourceAmounts{Wood: 20}, want: "20 Wood"},
		{cost: ResourceAmounts{Wood: 10, Stone: 10}, want: "10 Wood, 10 Stone"},
		{cost: ResourceAmounts{Wood: 10, Stone: 10, Metal: 10}, want: "10 Wood, 10 Stone, 10 Metal"},
		{cost: ResourceAmounts{Wood: 30, Stone: 10, Metal: 10}, want: "30 Wood, 10 Stone, 10 Metal"},
	}
	for _, test := range resourceTests {
		if got := formatSelectionResourceCost(test.cost); got != test.want {
			t.Fatalf("formatSelectionResourceCost(%+v) = %q, want %q", test.cost, got, test.want)
		}
	}

	healthTests := []struct {
		health int
		max    int
		want   int
	}{
		{health: 25, max: 50, want: 50},
		{health: 1, max: 3, want: 33},
		{health: -1, max: 10, want: 0},
		{health: 12, max: 10, want: 100},
		{health: 5, max: 0, want: 0},
	}
	for _, test := range healthTests {
		if got := selectedHealthPercent(test.health, test.max); got != test.want {
			t.Fatalf("selectedHealthPercent(%d, %d) = %d, want %d", test.health, test.max, got, test.want)
		}
	}
}

// TestSelectionPanelBounds verifies panel bounds and hit testing use row count.
func TestSelectionPanelBounds(t *testing.T) {
	data := SelectionPanelData{
		Kind:                SelectionPanelTower,
		Name:                "Bow Tower",
		RangeTiles:          3.0,
		FireIntervalSeconds: 1.0,
		Damage:              10,
		Staffing:            PopulationAmounts{Soldiers: 1},
	}

	bounds, ok := SelectionPanelBounds(1920, 1080, data)
	if !ok {
		t.Fatal("expected selection panel bounds")
	}
	if bounds.X != 1488 || bounds.Y != 738 || bounds.W != 390 || bounds.H != 300 {
		t.Fatalf("bounds = %+v, want X:1488 Y:738 W:390 H:300", bounds)
	}
	if !SelectionPanelContains(1920, 1080, data, bounds.X+bounds.W/2, bounds.Y+bounds.H/2) {
		t.Fatal("expected center point inside selection panel")
	}
	if SelectionPanelContains(1920, 1080, data, bounds.X-1, bounds.Y) {
		t.Fatal("expected point outside selection panel")
	}
	if _, ok := SelectionPanelBounds(1920, 1080, SelectionPanelData{}); ok {
		t.Fatal("expected no bounds for empty selection panel")
	}
}

func assertSelectionPanelRow(t *testing.T, rows []selectionPanelRow, label, value string) {
	t.Helper()
	for _, row := range rows {
		if row.Label == label {
			if row.Value != value {
				t.Fatalf("%s value = %q, want %q", label, row.Value, value)
			}
			return
		}
	}
	t.Fatalf("missing panel row %q in %+v", label, rows)
}
