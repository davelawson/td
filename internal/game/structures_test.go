package game

import (
	"testing"

	"td/assets"
)

// TestNewStructureCatalogIncludesSanctum verifies the Sanctum template values.
func TestNewStructureCatalogIncludesSanctum(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	sanctum := catalog.Sanctum

	if sanctum.Name != "Sanctum" {
		t.Fatalf("Sanctum name = %q, want %q", sanctum.Name, "Sanctum")
	}
	if sanctum.Sprite == nil {
		t.Fatal("expected Sanctum sprite to be assigned")
	}
	if sanctum.Sprite != assetCatalog.Sprite.Structure.Sanctum {
		t.Fatal("expected Sanctum sprite to reference the loaded asset catalog sprite")
	}
}

// TestNewStructureCatalogIncludesBowTower verifies the Bow Tower template values.
func TestNewStructureCatalogIncludesBowTower(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	bowTower := catalog.BowTower

	if bowTower.Name != "Bow Tower" {
		t.Fatalf("Bow Tower name = %q, want %q", bowTower.Name, "Bow Tower")
	}
	if bowTower.Sprite == nil {
		t.Fatal("expected Bow Tower sprite to be assigned")
	}
	if bowTower.Sprite != assetCatalog.Sprite.Structure.BowTower {
		t.Fatal("expected Bow Tower sprite to reference the loaded asset catalog sprite")
	}
}

// TestStructureStoresTemplateAndTileCoordinates verifies a placed structure instance.
func TestStructureStoresTemplateAndTileCoordinates(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	structure := Structure{
		Template: &catalog.BowTower,
		X:        3,
		Y:        4,
	}

	if structure.Template != &catalog.BowTower {
		t.Fatal("expected structure template to reference the Bow Tower template")
	}
	if structure.Template.Name != "Bow Tower" {
		t.Fatalf("structure template name = %q, want %q", structure.Template.Name, "Bow Tower")
	}
	if structure.Template.Sprite != assetCatalog.Sprite.Structure.BowTower {
		t.Fatal("expected structure template sprite to reference the loaded Bow Tower sprite")
	}
	if structure.X != 3 || structure.Y != 4 {
		t.Fatalf("structure coordinates = (%d,%d), want (3,4)", structure.X, structure.Y)
	}
}
