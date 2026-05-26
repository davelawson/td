package game

import (
	"td/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// StructureTemplate describes shared metadata for every instance of one structure type.
type StructureTemplate struct {
	Name                          string
	Sprite                        *ebiten.Image
	Cost                          Resources
	RangeTiles                    float64
	Damage                        int
	FireIntervalSeconds           float64
	ProjectileSpeedTilesPerSecond float64
	ProjectileSprite              *ebiten.Image
	DamageAllEnemiesInTargetTile  bool
}

// Resources describes the resources required to construct a structure.
type Resources struct {
	Wood  int
	Stone int
	Metal int
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
	CatapultTower  StructureTemplate
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
			Cost:                          Resources{Wood: 30, Stone: 10, Metal: 10},
			RangeTiles:                    3.0,
			Damage:                        10,
			FireIntervalSeconds:           1.0,
			ProjectileSpeedTilesPerSecond: 9.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.BowTowerProjectile,
		},
		FlameBoltTower: StructureTemplate{
			Name:                          "Flame Bolt Tower",
			Sprite:                        assetCatalog.Sprite.Structure.FlameBoltTower,
			Cost:                          Resources{Stone: 30, Metal: 20},
			RangeTiles:                    2.5,
			Damage:                        20,
			FireIntervalSeconds:           1.5,
			ProjectileSpeedTilesPerSecond: 7.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.FlameBoltTowerProjectile,
		},
		CatapultTower: StructureTemplate{
			Name:                          "Catapult Tower",
			Sprite:                        assetCatalog.Sprite.Structure.CatapultTower,
			Cost:                          Resources{Wood: 40, Stone: 60, Metal: 25},
			RangeTiles:                    5.0,
			Damage:                        75,
			FireIntervalSeconds:           3.0,
			ProjectileSpeedTilesPerSecond: 3.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.CatapultTowerProjectile,
			DamageAllEnemiesInTargetTile:  true,
		},
	}
}
