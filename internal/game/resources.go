package game

type resourceCounts struct {
	wood  int
	stone int
	iron  int
	gold  int
}

// Resources describes a set of game resources.
type Resources struct {
	Wood  int
	Stone int
	Iron  int
	Gold  int
}

// canAffordBuildingCost reports whether current resources cover a structure cost.
func (s *State) canAffordBuildingCost(cost Resources) bool {
	return s.status.resources.wood >= cost.Wood &&
		s.status.resources.stone >= cost.Stone &&
		s.status.resources.iron >= cost.Iron &&
		s.status.resources.gold >= cost.Gold
}

// deductBuildingCost spends the resources required to build a structure.
func (s *State) deductBuildingCost(cost Resources) {
	s.status.resources.wood -= cost.Wood
	s.status.resources.stone -= cost.Stone
	s.status.resources.iron -= cost.Iron
	s.status.resources.gold -= cost.Gold
}

// grantEnemyGold awards the template Gold drop for a combat-defeated enemy.
func (s *State) grantEnemyGold(enemy raidEnemy) {
	if enemy.template == nil {
		return
	}
	s.status.resources.gold += enemy.template.GoldDrop
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
				production, ok := s.economicBuildingProduction(plot.Tiles[y][x].Feature)
				if !ok {
					continue
				}
				producer := tileCoordinate{Plot: plotCoord, X: x, Y: y}
				terrainTile, ok := s.nearestTerrainTile(producer, production.terrain)
				if !ok {
					continue
				}
				s.consumeTerrain(terrainTile)
				s.status.resources.wood += production.yield.Wood
				s.status.resources.stone += production.yield.Stone
				s.status.resources.iron += production.yield.Iron
				s.status.resources.gold += production.yield.Gold
			}
		}
	}
}

type economicProduction struct {
	yield   Resources
	terrain tileTerrain
}

// economicBuildingProduction returns the yield and terrain consumed by one placed feature.
func (s *State) economicBuildingProduction(feature tileFeature) (economicProduction, bool) {
	switch feature {
	case featureWoodcutter:
		return economicProduction{
			yield:   s.structureCatalog.Woodcutter.ResourceYield,
			terrain: terrainTree,
		}, true
	case featureStoneQuarry:
		return economicProduction{
			yield:   s.structureCatalog.StoneQuarry.ResourceYield,
			terrain: terrainBoulder,
		}, true
	case featureIronMine:
		return economicProduction{
			yield:   s.structureCatalog.IronMine.ResourceYield,
			terrain: terrainIronDeposit,
		}, true
	default:
		return economicProduction{}, false
	}
}

// nearestTerrainTile returns the closest matching terrain anywhere in the explored Domain.
func (s *State) nearestTerrainTile(origin tileCoordinate, terrain tileTerrain) (tileCoordinate, bool) {
	originCenter := plotTileWorldCenter(origin.Plot, origin.X, origin.Y)
	var nearest tileCoordinate
	var nearestDistance float64
	found := false

	for _, plotCoord := range s.gameMap.exploredPlotCoordinates() {
		plot, ok := s.gameMap.plot(plotCoord)
		if !ok {
			continue
		}
		for y := 0; y < plotSize; y++ {
			for x := 0; x < plotSize; x++ {
				if plot.Tiles[y][x].Terrain != terrain {
					continue
				}
				center := plotTileWorldCenter(plotCoord, x, y)
				dx := center.X - originCenter.X
				dy := center.Y - originCenter.Y
				distance := dx*dx + dy*dy
				if found && distance >= nearestDistance {
					continue
				}
				nearest = tileCoordinate{Plot: plotCoord, X: x, Y: y}
				nearestDistance = distance
				found = true
			}
		}
	}
	return nearest, found
}

// consumeTerrain replaces one natural Tile with its biome's default terrain.
func (s *State) consumeTerrain(tile tileCoordinate) {
	plot, ok := s.gameMap.plot(tile.Plot)
	if !ok || tile.X < 0 || tile.X >= plotSize || tile.Y < 0 || tile.Y >= plotSize {
		return
	}
	plot.Tiles[tile.Y][tile.X].Terrain = defaultTerrainForBiome(plot.Biome)
	if s.selection.kind == selectedItemTerrain && s.selection.tile == tile {
		s.selection = selectedItem{}
	}
}

// defaultTerrainForBiome returns the terrain left after a natural Tile is consumed.
func defaultTerrainForBiome(biome plotBiome) tileTerrain {
	switch biome {
	case biomeGrasslands, biomeHills:
		return terrainEmpty
	default:
		return terrainEmpty
	}
}
