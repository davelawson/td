package game

import (
	"td/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyTemplate describes shared stats for every instance of one enemy type.
type EnemyTemplate struct {
	Name                string
	MaxHealth           int
	SpeedTilesPerSecond float64
	SanctumDamage       int
	SpriteKey           string
	Sprite              *ebiten.Image
}

// EnemyCatalog groups every enemy template available to game systems.
type EnemyCatalog struct {
	SkeletonSwordShield EnemyTemplate
	Zombie              EnemyTemplate
}

// NewEnemyCatalog creates the default enemy template catalog.
func NewEnemyCatalog(assetCatalog assets.Catalog) EnemyCatalog {
	return EnemyCatalog{
		SkeletonSwordShield: EnemyTemplate{
			Name:                "Skeleton Sword-and-Shield",
			MaxHealth:           50,
			SpeedTilesPerSecond: 1.0,
			SanctumDamage:       1,
			SpriteKey:           "skeleton-sword-shield",
			Sprite:              assetCatalog.Sprite.Enemy.SkeletonSwordShield,
		},
		Zombie: EnemyTemplate{
			Name:                "Zombie",
			MaxHealth:           75,
			SpeedTilesPerSecond: 0.7,
			SanctumDamage:       1,
			SpriteKey:           "zombie",
			Sprite:              assetCatalog.Sprite.Enemy.Zombie,
		},
	}
}
