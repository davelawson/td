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
	if bowTower.Cost != (ResourceCost{Wood: 30, Stone: 10, Metal: 10}) {
		t.Fatalf("Bow Tower cost = %+v, want 30 wood 10 stone 10 metal", bowTower.Cost)
	}
	if bowTower.RangeTiles != 3.0 {
		t.Fatalf("Bow Tower range = %f, want 3.0", bowTower.RangeTiles)
	}
	if bowTower.Damage != 10 {
		t.Fatalf("Bow Tower damage = %d, want 10", bowTower.Damage)
	}
	if bowTower.FireIntervalSeconds != 1.0 {
		t.Fatalf("Bow Tower fire interval = %f, want 1.0", bowTower.FireIntervalSeconds)
	}
	if bowTower.ProjectileSpeedTilesPerSecond != 9.0 {
		t.Fatalf("Bow Tower projectile speed = %f, want 9.0", bowTower.ProjectileSpeedTilesPerSecond)
	}
	if bowTower.ProjectileSprite == nil {
		t.Fatal("expected Bow Tower projectile sprite to be assigned")
	}
	if bowTower.ProjectileSprite != assetCatalog.Sprite.Projectile.BowTowerProjectile {
		t.Fatal("expected Bow Tower projectile sprite to reference the loaded asset catalog projectile sprite")
	}
}

// TestNewStructureCatalogIncludesFlameBoltTower verifies the Flame Bolt Tower template values.
func TestNewStructureCatalogIncludesFlameBoltTower(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	flameBoltTower := catalog.FlameBoltTower

	if flameBoltTower.Name != "Flame Bolt Tower" {
		t.Fatalf("Flame Bolt Tower name = %q, want %q", flameBoltTower.Name, "Flame Bolt Tower")
	}
	if flameBoltTower.Sprite != assetCatalog.Sprite.Structure.FlameBoltTower {
		t.Fatal("expected Flame Bolt Tower sprite to reference the loaded asset catalog sprite")
	}
	if flameBoltTower.Cost != (ResourceCost{Stone: 30, Metal: 20}) {
		t.Fatalf("Flame Bolt Tower cost = %+v, want 30 stone 20 metal", flameBoltTower.Cost)
	}
	if flameBoltTower.RangeTiles != 2.5 {
		t.Fatalf("Flame Bolt Tower range = %f, want 2.5", flameBoltTower.RangeTiles)
	}
	if flameBoltTower.Damage != 20 {
		t.Fatalf("Flame Bolt Tower damage = %d, want 20", flameBoltTower.Damage)
	}
	if flameBoltTower.FireIntervalSeconds != 1.5 {
		t.Fatalf("Flame Bolt Tower fire interval = %f, want 1.5", flameBoltTower.FireIntervalSeconds)
	}
	if flameBoltTower.ProjectileSpeedTilesPerSecond != 7.0 {
		t.Fatalf("Flame Bolt Tower projectile speed = %f, want 7.0", flameBoltTower.ProjectileSpeedTilesPerSecond)
	}
	if flameBoltTower.ProjectileSprite != assetCatalog.Sprite.Projectile.FlameBoltTowerProjectile {
		t.Fatal("expected Flame Bolt Tower projectile sprite to reference the loaded asset catalog projectile sprite")
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
	if structure.Template.ProjectileSprite != assetCatalog.Sprite.Projectile.BowTowerProjectile {
		t.Fatal("expected structure template projectile sprite to reference the loaded Bow Tower projectile sprite")
	}
	if structure.X != 3 || structure.Y != 4 {
		t.Fatalf("structure coordinates = (%d,%d), want (3,4)", structure.X, structure.Y)
	}
}
