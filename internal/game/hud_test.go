package game

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TestStateResourceHUDItems verifies the top bar resource icon order and counts.
func TestStateResourceHUDItems(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	items := state.resourceHUDItems()
	if len(items) != 3 {
		t.Fatalf("resource item count = %d, want 3", len(items))
	}
	tests := []struct {
		name  string
		count int
		color color.Color
	}{
		{name: "Wood", count: 100, color: colors.resourceWood},
		{name: "Stone", count: 50, color: colors.resourceStone},
		{name: "Metal", count: 20, color: colors.resourceMetal},
	}
	for i, test := range tests {
		item := items[i]
		if item.Name != test.name {
			t.Fatalf("resource item %d name = %q, want %q", i, item.Name, test.name)
		}
		if item.Count != test.count {
			t.Fatalf("resource item %s count = %d, want %d", item.Name, item.Count, test.count)
		}
		if item.Sprite == nil {
			t.Fatalf("resource item %s sprite is nil", item.Name)
		}
		if item.Color != test.color {
			t.Fatalf("resource item %s color = %#v, want %#v", item.Name, item.Color, test.color)
		}
	}
}

// TestStatePopulationHUDItems verifies the top bar population icon order and initial values.
func TestStatePopulationHUDItems(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	items := state.populationHUDItems()
	if len(items) != 3 {
		t.Fatalf("population item count = %d, want 3", len(items))
	}
	names := []string{"Apprentice", "Soldier", "Peasant"}
	for i, name := range names {
		item := items[i]
		if item.Name != name {
			t.Fatalf("population item %d name = %q, want %q", i, item.Name, name)
		}
		if item.Available != 0 || item.Total != 0 {
			t.Fatalf("population item %s = %d/%d, want 0/0", item.Name, item.Available, item.Total)
		}
		if item.Sprite == nil {
			t.Fatalf("population item %s sprite is nil", item.Name)
		}
		if item.Color != colors.text {
			t.Fatalf("population item %s color = %#v, want %#v", item.Name, item.Color, colors.text)
		}
	}
}

// TestPopulationHUDItemTextFormatsAvailableBeforeTotal verifies population value order.
func TestPopulationHUDItemTextFormatsAvailableBeforeTotal(t *testing.T) {
	item := populationHUDItem{Available: 3, Total: 8}
	if value := populationHUDItemText(item); value != "3/8" {
		t.Fatalf("populationHUDItemText = %q, want %q", value, "3/8")
	}
}

// TestStateDomainStatusWidthIncludesSeparateGroups verifies full status measurement.
func TestStateDomainStatusWidthIncludesSeparateGroups(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	resources := state.resourceHUDItems()
	populations := state.populationHUDItems()
	barricade := state.barricadeText()
	got := state.domainStatusWidth(resources, populations, barricade)
	barricadeWidth, _ := text.Measure(barricade, state.ui.hudFace, state.ui.hudFace.Size)
	want := state.resourceHUDGroupWidth(resources) +
		state.statusGroupSeparatorWidth() +
		state.populationHUDGroupWidth(populations) +
		state.statusGroupSeparatorWidth() +
		barricadeWidth
	if got != want {
		t.Fatalf("domainStatusWidth = %f, want %f", got, want)
	}
}

// TestStateFormatsBarricadeText verifies the top bar Sanctum defense text.
func TestStateFormatsBarricadeText(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if value := state.barricadeText(); value != "Barricade 3" {
		t.Fatalf("barricadeText = %q, want %q", value, "Barricade 3")
	}
}
