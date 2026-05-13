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
