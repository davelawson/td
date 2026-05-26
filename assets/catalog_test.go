package assets

import (
	"bytes"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// TestNewCatalogLoadsSkeletonSwordShieldSprite verifies the required enemy sprite is embedded.
func TestNewCatalogLoadsSkeletonSwordShieldSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	skeleton := catalog.Sprite.Enemy.SkeletonSwordShield
	if skeleton == nil {
		t.Fatal("expected skeleton sword-and-shield sprite to load")
	}
	width, height := skeleton.Bounds().Dx(), skeleton.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("skeleton sword-and-shield sprite size = %dx%d, want 64x64", width, height)
	}
}

// TestNewCatalogLoadsZombieSprite verifies the required zombie sprite is embedded.
func TestNewCatalogLoadsZombieSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	zombie := catalog.Sprite.Enemy.Zombie
	if zombie == nil {
		t.Fatal("expected zombie sprite to load")
	}
	width, height := zombie.Bounds().Dx(), zombie.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("zombie sprite size = %dx%d, want 64x64", width, height)
	}
}

// TestNewCatalogLoadsResourceIconSprites verifies the required HUD icon sprites are embedded.
func TestNewCatalogLoadsResourceIconSprites(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		icon *ebiten.Image
	}{
		{name: "wood", icon: catalog.Sprite.Icon.Wood},
		{name: "stone", icon: catalog.Sprite.Icon.Stone},
		{name: "metal", icon: catalog.Sprite.Icon.Metal},
	}
	for _, test := range tests {
		if test.icon == nil {
			t.Fatalf("expected %s icon sprite to load", test.name)
		}
		width, height := test.icon.Bounds().Dx(), test.icon.Bounds().Dy()
		if width != 64 || height != 64 {
			t.Fatalf("%s icon sprite size = %dx%d, want 64x64", test.name, width, height)
		}
	}
}

// TestNewCatalogLoadsBowTowerProjectileSprite verifies the required projectile sprite is embedded.
func TestNewCatalogLoadsBowTowerProjectileSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	projectile := catalog.Sprite.Projectile.BowTowerProjectile
	if projectile == nil {
		t.Fatal("expected Bow Tower projectile sprite to load")
	}
	width, height := projectile.Bounds().Dx(), projectile.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("Bow Tower projectile sprite size = %dx%d, want 64x64", width, height)
	}
}

// TestNewCatalogLoadsFlameBoltTowerProjectileSprite verifies the required projectile sprite is embedded.
func TestNewCatalogLoadsFlameBoltTowerProjectileSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	projectile := catalog.Sprite.Projectile.FlameBoltTowerProjectile
	if projectile == nil {
		t.Fatal("expected Flame Bolt Tower projectile sprite to load")
	}
	width, height := projectile.Bounds().Dx(), projectile.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("Flame Bolt Tower projectile sprite size = %dx%d, want 64x64", width, height)
	}
}

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

// TestNewCatalogLoadsBowTowerSprite verifies the required Bow Tower sprite is embedded.
func TestNewCatalogLoadsBowTowerSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	bowTower := catalog.Sprite.Structure.BowTower
	if bowTower == nil {
		t.Fatal("expected Bow Tower sprite to load")
	}
	width, height := bowTower.Bounds().Dx(), bowTower.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("Bow Tower sprite size = %dx%d, want 64x64", width, height)
	}
}

// TestNewCatalogLoadsFlameBoltTowerSprite verifies the required Flame Bolt Tower sprite is embedded.
func TestNewCatalogLoadsFlameBoltTowerSprite(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	flameBoltTower := catalog.Sprite.Structure.FlameBoltTower
	if flameBoltTower == nil {
		t.Fatal("expected Flame Bolt Tower sprite to load")
	}
	width, height := flameBoltTower.Bounds().Dx(), flameBoltTower.Bounds().Dy()
	if width != 64 || height != 64 {
		t.Fatalf("Flame Bolt Tower sprite size = %dx%d, want 64x64", width, height)
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

// TestNewCatalogLoadsRaiderDefeatedAudio verifies the required defeat sound is embedded.
func TestNewCatalogLoadsRaiderDefeatedAudio(t *testing.T) {
	catalog, err := NewCatalog()
	if err != nil {
		t.Fatal(err)
	}

	if len(catalog.Audio.RaiderDefeated) == 0 {
		t.Fatal("expected raider defeated audio bytes to load")
	}
	if _, err := wav.DecodeF32(bytes.NewReader(catalog.Audio.RaiderDefeated)); err != nil {
		t.Fatalf("decode raider defeated audio: %v", err)
	}
}
