package game

import "td/internal/ui"

// buildingBarItemID identifies a buildable structure independent of visible tab position.
type buildingBarItemID = int

const (
	buildingBarHouseIndex buildingBarItemID = iota
	buildingBarBarracksIndex
	buildingBarDormIndex
	buildingBarWoodcutterIndex
	buildingBarStoneQuarryIndex
	buildingBarIronMineIndex
	buildingBarBowTowerIndex
	buildingBarFlameBoltTowerIndex
	buildingBarCatapultTowerIndex
)

// buildingBarCategory identifies one visible group of buildable structures.
type buildingBarCategory int

const (
	buildingBarCategoryDefenses buildingBarCategory = iota
	buildingBarCategoryEconomic
	buildingBarCategoryHousing
	buildingBarNoCategory buildingBarCategory = -1
)

// buildingBarTab describes one screen-space building category tab.
type buildingBarTab struct {
	Category buildingBarCategory
	Label    string
	Bounds   ui.Button[int]
}

// buildingBarCategories returns tabs in their rendered order.
func buildingBarCategories() []buildingBarCategory {
	return []buildingBarCategory{
		buildingBarCategoryDefenses,
		buildingBarCategoryEconomic,
		buildingBarCategoryHousing,
	}
}

// buildingBarCategoryLabel returns the text shown for one category tab.
func buildingBarCategoryLabel(category buildingBarCategory) string {
	switch category {
	case buildingBarCategoryDefenses:
		return "Defenses"
	case buildingBarCategoryEconomic:
		return "Economic"
	case buildingBarCategoryHousing:
		return "Housing"
	default:
		return ""
	}
}

// buildingBarCategoryForItem returns the tab that contains a stable item.
func buildingBarCategoryForItem(id buildingBarItemID) buildingBarCategory {
	switch id {
	case buildingBarHouseIndex, buildingBarBarracksIndex, buildingBarDormIndex:
		return buildingBarCategoryHousing
	case buildingBarWoodcutterIndex, buildingBarStoneQuarryIndex, buildingBarIronMineIndex:
		return buildingBarCategoryEconomic
	case buildingBarBowTowerIndex, buildingBarFlameBoltTowerIndex, buildingBarCatapultTowerIndex:
		return buildingBarCategoryDefenses
	default:
		return buildingBarNoCategory
	}
}

// buildingBarItemIDsForCategory returns stable item IDs in visible category order.
func buildingBarItemIDsForCategory(category buildingBarCategory) []buildingBarItemID {
	switch category {
	case buildingBarCategoryHousing:
		return []buildingBarItemID{buildingBarHouseIndex, buildingBarBarracksIndex, buildingBarDormIndex}
	case buildingBarCategoryEconomic:
		return []buildingBarItemID{buildingBarWoodcutterIndex, buildingBarStoneQuarryIndex, buildingBarIronMineIndex}
	case buildingBarCategoryDefenses:
		return []buildingBarItemID{buildingBarBowTowerIndex, buildingBarFlameBoltTowerIndex, buildingBarCatapultTowerIndex}
	default:
		return nil
	}
}

// buildingFeatureForItemID maps a stable item ID to the Tile feature it places.
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

// buildingTemplateForItemID returns the structure template for a stable item ID.
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
