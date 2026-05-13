package assets

import "testing"

// TestNewCatalogLoadsSanctumSprite verifies the required structure sprite is embedded.
func TestNewCatalogLoadsSanctumSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	sanctum := catalog.Sprite.Structure.Sanctum
	if sanctum == nil {
		t.Fatal("expected Sanctum sprite to load")
	}
	width, height := sanctum.Bounds().Dx(), sanctum.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("Sanctum sprite size = %dx%d, want 64x64", width, height)
	}
}

// TestNewCatalogLoadsPineTreeSprites verifies the required terrain sprites are embedded.
func TestNewCatalogLoadsPineTreeSprites(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	for i, tree := range catalog.Sprite.Terrain.PineTrees {
		if tree == nil {
			t.Fatalf("expected pine tree sprite %d to load", i+1)
		}
		width, height := tree.Bounds().Dx(), tree.Bounds().Dy()
		if width != 64 || height != 64 {
			t.Fatalf("pine tree sprite %d size = %dx%d, want 64x64", i+1, width, height)
		}
	}
}
