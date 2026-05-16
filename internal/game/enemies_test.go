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
	if skeleton.MaxHealth != 20 {
		t.Fatalf("skeleton max health = %d, want 20", skeleton.MaxHealth)
	}
	if skeleton.Speed != 3.0 {
		t.Fatalf("skeleton speed = %f, want 3.0", skeleton.Speed)
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
