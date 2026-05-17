package game

import (
	"testing"

	"td/assets"
)

// TestNewEnemyCatalogIncludesSkeletonSwordShield verifies the initial enemy template values.
func TestNewEnemyCatalogIncludesSkeletonSwordShield(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewEnemyCatalog(assetCatalog)
	skeleton := catalog.SkeletonSwordShield

	if skeleton.Name != "Skeleton Sword-and-Shield" {
		t.Fatalf("skeleton name = %q, want %q", skeleton.Name, "Skeleton Sword-and-Shield")
	}
	if skeleton.MaxHealth != 50 {
		t.Fatalf("skeleton max health = %d, want 50", skeleton.MaxHealth)
	}
	if skeleton.SpeedTilesPerSecond != 1.0 {
		t.Fatalf("skeleton speed = %f, want 1.0", skeleton.SpeedTilesPerSecond)
	}
	if skeleton.SanctumDamage != 1 {
		t.Fatalf("skeleton Sanctum damage = %d, want 1", skeleton.SanctumDamage)
	}
	if skeleton.SpriteKey != "skeleton-sword-shield" {
		t.Fatalf("skeleton sprite key = %q, want %q", skeleton.SpriteKey, "skeleton-sword-shield")
	}
	if skeleton.Sprite == nil {
		t.Fatal("expected skeleton sprite to be assigned")
	}
	if skeleton.Sprite != assetCatalog.Sprite.Enemy.SkeletonSwordShield {
		t.Fatal("expected skeleton sprite to reference the loaded asset catalog sprite")
	}
}

// TestNewEnemyCatalogIncludesZombie verifies the zombie enemy template values.
func TestNewEnemyCatalogIncludesZombie(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewEnemyCatalog(assetCatalog)
	zombie := catalog.Zombie

	if zombie.Name != "Zombie" {
		t.Fatalf("zombie name = %q, want %q", zombie.Name, "Zombie")
	}
	if zombie.MaxHealth != 75 {
		t.Fatalf("zombie max health = %d, want 75", zombie.MaxHealth)
	}
	if zombie.SpeedTilesPerSecond != 0.7 {
		t.Fatalf("zombie speed = %f, want 0.7", zombie.SpeedTilesPerSecond)
	}
	if zombie.SanctumDamage != 1 {
		t.Fatalf("zombie Sanctum damage = %d, want 1", zombie.SanctumDamage)
	}
	if zombie.SpriteKey != "zombie" {
		t.Fatalf("zombie sprite key = %q, want %q", zombie.SpriteKey, "zombie")
	}
	if zombie.Sprite == nil {
		t.Fatal("expected zombie sprite to be assigned")
	}
	if zombie.Sprite != assetCatalog.Sprite.Enemy.Zombie {
		t.Fatal("expected zombie sprite to reference the loaded asset catalog sprite")
	}
}
