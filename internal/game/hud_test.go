package game

import "testing"

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
	}{
		{name: "Wood", count: 80},
		{name: "Stone", count: 45},
		{name: "Metal", count: 12},
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
	}
}

// TestStateFormatsBarricadeText verifies the top bar Sanctum defense text.
func TestStateFormatsBarricadeText(t *testing.T) {
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}

	if value := state.barricadeText(); value != "| Barricade 3" {
		t.Fatalf("barricadeText = %q, want %q", value, "| Barricade 3")
	}
}
