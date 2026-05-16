package game

import (
	"td/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// StructureTemplate describes shared metadata for every instance of one structure type.
type StructureTemplate struct {
	Name                          string
	Sprite                        *ebiten.Image
	RangeTiles                    float64
	Damage                        int
	FireIntervalSeconds           float64
	ProjectileSpeedTilesPerSecond float64
	ProjectileSprite              *ebiten.Image
}

// Structure describes one placed structure instance on the map.
type Structure struct {
	Template *StructureTemplate
	X        int
	Y        int
}

// StructureCatalog groups every structure template available to game systems.
type StructureCatalog struct {
	Sanctum        StructureTemplate
	BowTower       StructureTemplate
	FlameBoltTower StructureTemplate
}

// NewStructureCatalog creates the default structure template catalog.
func NewStructureCatalog(assetCatalog assets.Catalog) StructureCatalog {
	return StructureCatalog{
		Sanctum: StructureTemplate{
			Name:   "Sanctum",
			Sprite: assetCatalog.Sprite.Structure.Sanctum,
		},
		BowTower: StructureTemplate{
			Name:                          "Bow Tower",
			Sprite:                        assetCatalog.Sprite.Structure.BowTower,
			RangeTiles:                    3.0,
			Damage:                        10,
			FireIntervalSeconds:           1.0,
			ProjectileSpeedTilesPerSecond: 9.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.BowTowerProjectile,
		},
		FlameBoltTower: StructureTemplate{
			Name:                          "Flame Bolt Tower",
			Sprite:                        assetCatalog.Sprite.Structure.FlameBoltTower,
			RangeTiles:                    2.5,
			Damage:                        20,
			FireIntervalSeconds:           1.5,
			ProjectileSpeedTilesPerSecond: 7.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.FlameBoltTowerProjectile,
		},
	}
}
