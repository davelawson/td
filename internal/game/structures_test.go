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
	if bowTower.Cost != (Resources{Wood: 30, Stone: 10, Metal: 10}) {
		t.Fatalf("Bow Tower cost = %+v, want 30 wood 10 stone 10 metal", bowTower.Cost)
	}
	if bowTower.Staffing != (StaffingRequirements{Soldiers: 1}) {
		t.Fatalf("Bow Tower staffing = %+v, want 1 Soldier", bowTower.Staffing)
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

// TestNewStructureCatalogIncludesHouse verifies the House template values.
func TestNewStructureCatalogIncludesHouse(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	house := catalog.House

	if house.Name != "House" {
		t.Fatalf("House name = %q, want %q", house.Name, "House")
	}
	if house.Sprite == nil {
		t.Fatal("expected House sprite to be assigned")
	}
	if house.Sprite != assetCatalog.Sprite.Structure.House {
		t.Fatal("expected House sprite to reference the loaded asset catalog sprite")
	}
	if house.Cost != (Resources{Wood: 20}) {
		t.Fatalf("House cost = %+v, want 20 wood", house.Cost)
	}
	if house.Staffing != (StaffingRequirements{}) {
		t.Fatalf("House staffing = %+v, want none", house.Staffing)
	}
	if house.PopulationGrant != (PopulationGrant{Peasants: 2}) {
		t.Fatalf("House population grant = %+v, want 2 Peasants", house.PopulationGrant)
	}
	if house.canFireProjectiles() {
		t.Fatal("expected House not to have projectile combat stats")
	}
}

// TestNewStructureCatalogIncludesBarracks verifies the Barracks template values.
func TestNewStructureCatalogIncludesBarracks(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	barracks := catalog.Barracks

	if barracks.Name != "Barracks" {
		t.Fatalf("Barracks name = %q, want %q", barracks.Name, "Barracks")
	}
	if barracks.Sprite == nil {
		t.Fatal("expected Barracks sprite to be assigned")
	}
	if barracks.Sprite != assetCatalog.Sprite.Structure.Barracks {
		t.Fatal("expected Barracks sprite to reference the loaded asset catalog sprite")
	}
	if barracks.Cost != (Resources{Wood: 10, Stone: 10}) {
		t.Fatalf("Barracks cost = %+v, want 10 wood 10 stone", barracks.Cost)
	}
	if barracks.Staffing != (StaffingRequirements{}) {
		t.Fatalf("Barracks staffing = %+v, want none", barracks.Staffing)
	}
	if barracks.PopulationCost != (PopulationCost{Peasants: 2}) {
		t.Fatalf("Barracks population cost = %+v, want 2 Peasants", barracks.PopulationCost)
	}
	if barracks.PopulationGrant != (PopulationGrant{Soldiers: 2}) {
		t.Fatalf("Barracks population grant = %+v, want 2 Soldiers", barracks.PopulationGrant)
	}
	if barracks.canFireProjectiles() {
		t.Fatal("expected Barracks not to have projectile combat stats")
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
	if flameBoltTower.Cost != (Resources{Stone: 30, Metal: 20}) {
		t.Fatalf("Flame Bolt Tower cost = %+v, want 30 stone 20 metal", flameBoltTower.Cost)
	}
	if flameBoltTower.Staffing != (StaffingRequirements{Apprentices: 1}) {
		t.Fatalf("Flame Bolt Tower staffing = %+v, want 1 Apprentice", flameBoltTower.Staffing)
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

// TestNewStructureCatalogIncludesCatapultTower verifies the Catapult Tower template values.
func TestNewStructureCatalogIncludesCatapultTower(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	catalog := NewStructureCatalog(assetCatalog)
	catapultTower := catalog.CatapultTower

	if catapultTower.Name != "Catapult Tower" {
		t.Fatalf("Catapult Tower name = %q, want %q", catapultTower.Name, "Catapult Tower")
	}
	if catapultTower.Sprite != assetCatalog.Sprite.Structure.CatapultTower {
		t.Fatal("expected Catapult Tower sprite to reference the loaded asset catalog sprite")
	}
	if catapultTower.Cost != (Resources{Wood: 40, Stone: 60, Metal: 25}) {
		t.Fatalf("Catapult Tower cost = %+v, want 40 wood 60 stone 25 metal", catapultTower.Cost)
	}
	if catapultTower.Staffing != (StaffingRequirements{Soldiers: 1, Peasants: 2}) {
		t.Fatalf("Catapult Tower staffing = %+v, want 1 Soldier and 2 Peasants", catapultTower.Staffing)
	}
	if catapultTower.RangeTiles != 5.0 {
		t.Fatalf("Catapult Tower range = %f, want 5.0", catapultTower.RangeTiles)
	}
	if catapultTower.Damage != 75 {
		t.Fatalf("Catapult Tower damage = %d, want 75", catapultTower.Damage)
	}
	if catapultTower.FireIntervalSeconds != 3.0 {
		t.Fatalf("Catapult Tower fire interval = %f, want 3.0", catapultTower.FireIntervalSeconds)
	}
	if catapultTower.ProjectileSpeedTilesPerSecond != 3.0 {
		t.Fatalf("Catapult Tower projectile speed = %f, want 3.0", catapultTower.ProjectileSpeedTilesPerSecond)
	}
	if catapultTower.ProjectileSprite != assetCatalog.Sprite.Projectile.CatapultTowerProjectile {
		t.Fatal("expected Catapult Tower projectile sprite to reference the loaded asset catalog projectile sprite")
	}
	if !catapultTower.DamageAllEnemiesInTargetTile {
		t.Fatal("expected Catapult Tower to damage all enemies in its target Tile")
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
