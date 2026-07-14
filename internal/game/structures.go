package game

import (
	"td/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// StructureTemplate describes shared metadata for every instance of one structure type.
type StructureTemplate struct {
	Name                          string
	Description                   string
	Sprite                        *ebiten.Image
	Cost                          Resources
	Staffing                      StaffingRequirements
	PopulationCost                PopulationCost
	PopulationGrant               PopulationGrant
	ResourceYield                 Resources
	RangeTiles                    float64
	Damage                        int
	FireIntervalSeconds           float64
	ProjectileSpeedTilesPerSecond float64
	ProjectileSprite              *ebiten.Image
	DamageAllEnemiesInTargetTile  bool
}

// StaffingRequirements describes the inhabitants required by one structure.
type StaffingRequirements struct {
	Apprentices int
	Soldiers    int
	Peasants    int
}

// PopulationGrant describes inhabitants added by constructing one structure.
type PopulationGrant struct {
	Apprentices int
	Soldiers    int
	Peasants    int
}

// PopulationCost describes inhabitants removed by constructing one structure.
type PopulationCost struct {
	Apprentices int
	Soldiers    int
	Peasants    int
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
	House          StructureTemplate
	Barracks       StructureTemplate
	Dorm           StructureTemplate
	Woodcutter     StructureTemplate
	StoneQuarry    StructureTemplate
	IronMine       StructureTemplate
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
		House: StructureTemplate{
			Name:            "House",
			Description:     "Shelters new Peasants for the Domain.",
			Sprite:          assetCatalog.Sprite.Structure.House,
			Cost:            Resources{Wood: 20},
			PopulationGrant: PopulationGrant{Peasants: 2},
		},
		Barracks: StructureTemplate{
			Name:            "Barracks",
			Description:     "Trains Peasants into Soldiers for staffed defenses.",
			Sprite:          assetCatalog.Sprite.Structure.Barracks,
			Cost:            Resources{Wood: 10, Stone: 10},
			PopulationCost:  PopulationCost{Peasants: 2},
			PopulationGrant: PopulationGrant{Soldiers: 2},
		},
		Dorm: StructureTemplate{
			Name:            "Dorm",
			Description:     "Houses Peasants studying to become Apprentices.",
			Sprite:          assetCatalog.Sprite.Structure.Dorm,
			Cost:            Resources{Wood: 10, Stone: 10},
			PopulationCost:  PopulationCost{Peasants: 1},
			PopulationGrant: PopulationGrant{Apprentices: 1},
		},
		Woodcutter: StructureTemplate{
			Name:          "Woodcutter",
			Description:   "Assigns a Peasant to bring in Wood during Labour.",
			Sprite:        assetCatalog.Sprite.Structure.Woodcutter,
			Cost:          Resources{Wood: 10},
			Staffing:      StaffingRequirements{Peasants: 1},
			ResourceYield: Resources{Wood: 10},
		},
		StoneQuarry: StructureTemplate{
			Name:          "Stone Quarry",
			Description:   "Assigns a Peasant to quarry Stone during Labour.",
			Sprite:        assetCatalog.Sprite.Structure.StoneQuarry,
			Cost:          Resources{Wood: 10, Stone: 10},
			Staffing:      StaffingRequirements{Peasants: 1},
			ResourceYield: Resources{Stone: 10},
		},
		IronMine: StructureTemplate{
			Name:          "Iron Mine",
			Description:   "Assigns a Peasant to extract Metal during Labour.",
			Sprite:        assetCatalog.Sprite.Structure.IronMine,
			Cost:          Resources{Wood: 10, Stone: 10, Metal: 10},
			Staffing:      StaffingRequirements{Peasants: 1},
			ResourceYield: Resources{Metal: 10},
		},
		BowTower: StructureTemplate{
			Name:                          "Bow Tower",
			Description:                   "A staffed archer tower that fires quick arrows.",
			Sprite:                        assetCatalog.Sprite.Structure.BowTower,
			Cost:                          Resources{Wood: 20, Stone: 10},
			Staffing:                      StaffingRequirements{Soldiers: 1},
			RangeTiles:                    3.0,
			Damage:                        10,
			FireIntervalSeconds:           1.0,
			ProjectileSpeedTilesPerSecond: 9.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.BowTowerProjectile,
		},
		FlameBoltTower: StructureTemplate{
			Name:                          "Flame Bolt Tower",
			Description:                   "An apprentice-staffed tower that hurls focused fire.",
			Sprite:                        assetCatalog.Sprite.Structure.FlameBoltTower,
			Cost:                          Resources{Stone: 30, Metal: 20},
			Staffing:                      StaffingRequirements{Apprentices: 1},
			RangeTiles:                    2.5,
			Damage:                        20,
			FireIntervalSeconds:           1.5,
			ProjectileSpeedTilesPerSecond: 7.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.FlameBoltTowerProjectile,
		},
		CatapultTower: StructureTemplate{
			Name:                          "Catapult Tower",
			Description:                   "A heavy crewed tower that crushes enemies in one Tile.",
			Sprite:                        assetCatalog.Sprite.Structure.CatapultTower,
			Cost:                          Resources{Wood: 40, Stone: 60, Metal: 25},
			Staffing:                      StaffingRequirements{Soldiers: 1, Peasants: 1},
			RangeTiles:                    5.0,
			Damage:                        30,
			FireIntervalSeconds:           6.0,
			ProjectileSpeedTilesPerSecond: 3.0,
			ProjectileSprite:              assetCatalog.Sprite.Projectile.CatapultTowerProjectile,
			DamageAllEnemiesInTargetTile:  true,
		},
	}
}
