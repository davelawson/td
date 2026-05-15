package game

import (
	"td/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// StructureTemplate describes shared metadata for every instance of one structure type.
type StructureTemplate struct {
	Name   string
	Sprite *ebiten.Image
}

// Structure describes one placed structure instance on the map.
type Structure struct {
	Template *StructureTemplate
	X        int
	Y        int
}

// StructureCatalog groups every structure template available to game systems.
type StructureCatalog struct {
	Sanctum  StructureTemplate
	BowTower StructureTemplate
}

// NewStructureCatalog creates the default structure template catalog.
func NewStructureCatalog(assetCatalog assets.Catalog) StructureCatalog {
	return StructureCatalog{
		Sanctum: StructureTemplate{
			Name:   "Sanctum",
			Sprite: assetCatalog.Sprite.Structure.Sanctum,
		},
		BowTower: StructureTemplate{
			Name:   "Bow Tower",
			Sprite: assetCatalog.Sprite.Structure.BowTower,
		},
	}
}
