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
	if skeleton.GoldDrop != 1 {
		t.Fatalf("skeleton Gold drop = %d, want 1", skeleton.GoldDrop)
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
	if zombie.GoldDrop != 2 {
		t.Fatalf("zombie Gold drop = %d, want 2", zombie.GoldDrop)
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

// TestNewEnemyCatalogIncludesGhoul verifies the Ghoul enemy template values.
func TestNewEnemyCatalogIncludesGhoul(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewEnemyCatalog(assetCatalog)
	ghoul := catalog.Ghoul

	if ghoul.Name != "Ghoul" {
		t.Fatalf("Ghoul name = %q, want %q", ghoul.Name, "Ghoul")
	}
	if ghoul.MaxHealth != 40 {
		t.Fatalf("Ghoul max health = %d, want 40", ghoul.MaxHealth)
	}
	if ghoul.SpeedTilesPerSecond != 1.5 {
		t.Fatalf("Ghoul speed = %f, want 1.5", ghoul.SpeedTilesPerSecond)
	}
	if ghoul.SanctumDamage != 1 {
		t.Fatalf("Ghoul Sanctum damage = %d, want 1", ghoul.SanctumDamage)
	}
	if ghoul.GoldDrop != 3 {
		t.Fatalf("Ghoul Gold drop = %d, want 3", ghoul.GoldDrop)
	}
	if ghoul.SpriteKey != "ghoul" {
		t.Fatalf("Ghoul sprite key = %q, want %q", ghoul.SpriteKey, "ghoul")
	}
	if ghoul.Sprite == nil {
		t.Fatal("expected Ghoul sprite to be assigned")
	}
	if ghoul.Sprite != assetCatalog.Sprite.Enemy.Ghoul {
		t.Fatal("expected Ghoul sprite to reference the loaded asset catalog sprite")
	}
}

// TestNewEnemyCatalogIncludesArmouredSkeleton verifies the Armoured Skeleton template values.
func TestNewEnemyCatalogIncludesArmouredSkeleton(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewEnemyCatalog(assetCatalog)
	armouredSkeleton := catalog.ArmouredSkeleton

	if armouredSkeleton.Name != "Armoured Skeleton" {
		t.Fatalf("Armoured Skeleton name = %q, want %q", armouredSkeleton.Name, "Armoured Skeleton")
	}
	if armouredSkeleton.MaxHealth != 125 {
		t.Fatalf("Armoured Skeleton max health = %d, want 125", armouredSkeleton.MaxHealth)
	}
	if armouredSkeleton.SpeedTilesPerSecond != 0.9 {
		t.Fatalf("Armoured Skeleton speed = %f, want 0.9", armouredSkeleton.SpeedTilesPerSecond)
	}
	if armouredSkeleton.SanctumDamage != 1 {
		t.Fatalf("Armoured Skeleton Sanctum damage = %d, want 1", armouredSkeleton.SanctumDamage)
	}
	if armouredSkeleton.GoldDrop != 5 {
		t.Fatalf("Armoured Skeleton Gold drop = %d, want 5", armouredSkeleton.GoldDrop)
	}
	if armouredSkeleton.SpriteKey != "armoured-skeleton" {
		t.Fatalf("Armoured Skeleton sprite key = %q, want %q", armouredSkeleton.SpriteKey, "armoured-skeleton")
	}
	if armouredSkeleton.Sprite == nil {
		t.Fatal("expected Armoured Skeleton sprite to be assigned")
	}
	if armouredSkeleton.Sprite != assetCatalog.Sprite.Enemy.ArmouredSkeleton {
		t.Fatal("expected Armoured Skeleton sprite to reference the loaded asset catalog sprite")
	}
}

// TestEnemyTemplateSpeedsRemainDistinct verifies every raider type differs in speed by at least three percent.
func TestEnemyTemplateSpeedsRemainDistinct(t *testing.T) {
	assetCatalog, err := assets.NewCatalog()
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewEnemyCatalog(assetCatalog)
	templates := []*EnemyTemplate{
		&catalog.SkeletonSwordShield,
		&catalog.Zombie,
		&catalog.Ghoul,
		&catalog.ArmouredSkeleton,
	}

	for i, first := range templates {
		for _, second := range templates[i+1:] {
			slower, faster := first, second
			if slower.SpeedTilesPerSecond > faster.SpeedTilesPerSecond {
				slower, faster = faster, slower
			}
			separation := (faster.SpeedTilesPerSecond - slower.SpeedTilesPerSecond) / slower.SpeedTilesPerSecond
			if separation < 0.03 {
				t.Fatalf(
					"%s and %s speed separation = %.2f%%, want at least 3%%",
					first.Name,
					second.Name,
					separation*100,
				)
			}
		}
	}
}
