package game

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

// grantEconomicBuildingResources resolves placed economic building work during Labour.
func (s *State) grantEconomicBuildingResources() {
	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		plot, ok := s.gameMap.plot(plotCoord)
		if !ok {
			continue
		}
		for y := 0; y < plotSize; y++ {
			for x := 0; x < plotSize; x++ {
				yield, ok := s.economicBuildingYield(plot.Tiles[y][x].Feature)
				if !ok {
					continue
				}
				s.status.resources.wood += yield.Wood
				s.status.resources.stone += yield.Stone
				s.status.resources.metal += yield.Metal
			}
		}
	}
}

// economicBuildingYield returns the Labour-phase yield for one placed feature.
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
