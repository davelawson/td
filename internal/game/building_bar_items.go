package game

import "td/internal/ui"

// buildingBarItemID identifies a build action independent of visible tab position.
type buildingBarItemID = ui.BuildingBarAction

const (
	buildingBarHouseIndex          = ui.BuildingBarHouse
	buildingBarBarracksIndex       = ui.BuildingBarBarracks
	buildingBarDormIndex           = ui.BuildingBarDorm
	buildingBarWoodcutterIndex     = ui.BuildingBarWoodcutter
	buildingBarStoneQuarryIndex    = ui.BuildingBarStoneQuarry
	buildingBarIronMineIndex       = ui.BuildingBarIronMine
	buildingBarMarketIndex         = ui.BuildingBarMarket
	buildingBarBowTowerIndex       = ui.BuildingBarBowTower
	buildingBarFlameBoltTowerIndex = ui.BuildingBarFlameBoltTower
	buildingBarCatapultTowerIndex  = ui.BuildingBarCatapultTower
)

// buildingFeatureForItemID maps a stable building-bar action to the Tile feature it places.
func buildingFeatureForItemID(id buildingBarItemID) (tileFeature, bool) {
	switch id {
	case buildingBarHouseIndex:
		return featureHouse, true
	case buildingBarBarracksIndex:
		return featureBarracks, true
	case buildingBarDormIndex:
		return featureDorm, true
	case buildingBarWoodcutterIndex:
		return featureWoodcutter, true
	case buildingBarStoneQuarryIndex:
		return featureStoneQuarry, true
	case buildingBarIronMineIndex:
		return featureIronMine, true
	case buildingBarMarketIndex:
		return featureMarket, true
	case buildingBarBowTowerIndex:
		return featureBowTower, true
	case buildingBarFlameBoltTowerIndex:
		return featureFlameBoltTower, true
	case buildingBarCatapultTowerIndex:
		return featureCatapultTower, true
	default:
		return featureNone, false
	}
}

// buildingTemplateForItemID returns the structure template for a stable building-bar action.
func (s *State) buildingTemplateForItemID(id buildingBarItemID) (StructureTemplate, bool) {
	switch id {
	case buildingBarHouseIndex:
		return s.structureCatalog.House, true
	case buildingBarBarracksIndex:
		return s.structureCatalog.Barracks, true
	case buildingBarDormIndex:
		return s.structureCatalog.Dorm, true
	case buildingBarWoodcutterIndex:
		return s.structureCatalog.Woodcutter, true
	case buildingBarStoneQuarryIndex:
		return s.structureCatalog.StoneQuarry, true
	case buildingBarIronMineIndex:
		return s.structureCatalog.IronMine, true
	case buildingBarMarketIndex:
		return s.structureCatalog.Market, true
	case buildingBarBowTowerIndex:
		return s.structureCatalog.BowTower, true
	case buildingBarFlameBoltTowerIndex:
		return s.structureCatalog.FlameBoltTower, true
	case buildingBarCatapultTowerIndex:
		return s.structureCatalog.CatapultTower, true
	default:
		return StructureTemplate{}, false
	}
}
