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
	Resources           Resources
	SpriteKey           string
	Sprite              *ebiten.Image
}

// EnemyCatalog groups every enemy template available to game systems.
type EnemyCatalog struct {
	SkeletonSwordShield EnemyTemplate
	Zombie              EnemyTemplate
	Ghoul               EnemyTemplate
	ArmouredSkeleton    EnemyTemplate
}

// NewEnemyCatalog creates the default enemy template catalog.
func NewEnemyCatalog(assetCatalog assets.Catalog) EnemyCatalog {
	return EnemyCatalog{
		SkeletonSwordShield: EnemyTemplate{
			Name:                "Skeleton Sword-and-Shield",
			MaxHealth:           50,
			SpeedTilesPerSecond: 1.0,
			SanctumDamage:       1,
			Resources:           Resources{Wood: 5, Stone: 2},
			SpriteKey:           "skeleton-sword-shield",
			Sprite:              assetCatalog.Sprite.Enemy.SkeletonSwordShield,
		},
		Zombie: EnemyTemplate{
			Name:                "Zombie",
			MaxHealth:           75,
			SpeedTilesPerSecond: 0.7,
			SanctumDamage:       1,
			Resources:           Resources{Wood: 4, Stone: 3, Metal: 1},
			SpriteKey:           "zombie",
			Sprite:              assetCatalog.Sprite.Enemy.Zombie,
		},
		Ghoul: EnemyTemplate{
			Name:                "Ghoul",
			MaxHealth:           40,
			SpeedTilesPerSecond: 1.5,
			SanctumDamage:       1,
			Resources:           Resources{Wood: 4, Metal: 1},
			SpriteKey:           "ghoul",
			Sprite:              assetCatalog.Sprite.Enemy.Ghoul,
		},
		ArmouredSkeleton: EnemyTemplate{
			Name:                "Armoured Skeleton",
			MaxHealth:           125,
			SpeedTilesPerSecond: 0.9,
			SanctumDamage:       1,
			Resources:           Resources{Stone: 5, Metal: 2},
			SpriteKey:           "armoured-skeleton",
			Sprite:              assetCatalog.Sprite.Enemy.ArmouredSkeleton,
		},
	}
}
