package game

import (
	"fmt"
	"strings"
)

type resourceCounts struct {
	wood  int
	stone int
	metal int
}

// Resources describes a set of game resources.
type Resources struct {
	Wood  int
	Stone int
	Metal int
}

// canAffordBuildingCost reports whether current resources cover a structure cost.
func (s *State) canAffordBuildingCost(cost Resources) bool {
	return s.status.resources.wood >= cost.Wood &&
		s.status.resources.stone >= cost.Stone &&
		s.status.resources.metal >= cost.Metal
}

// deductBuildingCost spends the resources required to build a structure.
func (s *State) deductBuildingCost(cost Resources) {
	s.status.resources.wood -= cost.Wood
	s.status.resources.stone -= cost.Stone
	s.status.resources.metal -= cost.Metal
}

// grantEnemyResources awards the template resources for a combat-defeated enemy.
func (s *State) grantEnemyResources(enemy raidEnemy) {
	if enemy.template == nil {
		return
	}
	s.status.resources.wood += enemy.template.Resources.Wood
	s.status.resources.stone += enemy.template.Resources.Stone
	s.status.resources.metal += enemy.template.Resources.Metal
}

// grantEconomicBuildingResources awards placed economic building yields after a defeated Raid.
func (s *State) grantEconomicBuildingResources() {
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			yield, ok := s.economicBuildingYield(s.gameMap.Home.Tiles[y][x].Feature)
			if !ok {
				continue
			}
			s.status.resources.wood += yield.Wood
			s.status.resources.stone += yield.Stone
			s.status.resources.metal += yield.Metal
		}
	}
}

// economicBuildingYield returns the Raid-completion yield for one placed feature.
func (s *State) economicBuildingYield(feature tileFeature) (Resources, bool) {
	switch feature {
	case featureWoodcutter:
		return s.structureCatalog.Woodcutter.ResourceYield, true
	case featureStoneQuarry:
		return s.structureCatalog.StoneQuarry.ResourceYield, true
	case featureIronMine:
		return s.structureCatalog.IronMine.ResourceYield, true
	default:
		return Resources{}, false
	}
}

// formatResourceCost formats a compact human-readable construction cost.
func formatResourceCost(cost Resources) string {
	parts := []string{}
	if cost.Wood > 0 {
		parts = append(parts, fmt.Sprintf("%d Wood", cost.Wood))
	}
	if cost.Stone > 0 {
		parts = append(parts, fmt.Sprintf("%d Stone", cost.Stone))
	}
	if cost.Metal > 0 {
		parts = append(parts, fmt.Sprintf("%d Metal", cost.Metal))
	}
	if len(parts) == 0 {
		return "Free"
	}
	return strings.Join(parts, ", ")
}

// formatResourceYield formats Raid-completion resource production.
func formatResourceYield(yield Resources) string {
	cost := formatResourceCost(yield)
	if cost == "Free" {
		return "Nothing"
	}
	return cost + " after each Raid"
}
