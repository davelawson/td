package game

import "testing"

// TestNewEnemyCatalogIncludesSkeletonSwordShield verifies the initial enemy template values.
func TestNewEnemyCatalogIncludesSkeletonSwordShield(t *testing.T) {
	catalog := NewEnemyCatalog()
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
}
