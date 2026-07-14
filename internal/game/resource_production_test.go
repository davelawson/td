package game

import "testing"

func TestEconomicBuildingsConsumeMatchingTerrainForYield(t *testing.T) {
	tests := []struct {
		name     string
		feature  tileFeature
		terrain  tileTerrain
		resource func(resourceCounts) int
	}{
		{name: "woodcutter", feature: featureWoodcutter, terrain: terrainTree, resource: func(r resourceCounts) int { return r.wood }},
		{name: "stone quarry", feature: featureStoneQuarry, terrain: terrainBoulder, resource: func(r resourceCounts) int { return r.stone }},
		{name: "iron mine", feature: featureIronMine, terrain: terrainIronDeposit, resource: func(r resourceCounts) int { return r.iron }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newRaidTestState(t)
			clearNaturalTerrain(&state.gameMap.Home)
			producer := homeTileCoordinate(5, 5)
			source := homeTileCoordinate(6, 5)
			state.gameMap.Home.Tiles[producer.Y][producer.X].Feature = test.feature
			state.gameMap.Home.Tiles[source.Y][source.X] = Tile{Terrain: test.terrain, Tweak: 31}
			starting := test.resource(state.status.resources)

			state.grantEconomicBuildingResources()

			if got := test.resource(state.status.resources); got != starting+10 {
				t.Fatalf("resource = %d, want %d", got, starting+10)
			}
			if got := state.gameMap.Home.Tiles[source.Y][source.X].Terrain; got != terrainEmpty {
				t.Fatalf("consumed terrain = %v, want grass", got)
			}
		})
	}
}

func TestEconomicBuildingUsesNearestTerrainAcrossExploredDomain(t *testing.T) {
	state := newRaidTestState(t)
	clearNaturalTerrain(&state.gameMap.Home)
	producer := homeTileCoordinate(plotSize-1, homePlotCenter)
	farSource := homeTileCoordinate(0, homePlotCenter)
	state.gameMap.Home.Tiles[producer.Y][producer.X].Feature = featureWoodcutter
	state.gameMap.Home.Tiles[farSource.Y][farSource.X].Terrain = terrainTree

	eastCoord := plotCoordinate{X: 1}
	east := Plot{Biome: biomeHills}
	nearSource := tileCoordinate{Plot: eastCoord, X: 0, Y: homePlotCenter}
	east.Tiles[nearSource.Y][nearSource.X].Terrain = terrainTree
	state.gameMap.ensurePlots()
	state.gameMap.Plots[eastCoord] = &east

	state.grantEconomicBuildingResources()

	if got := east.Tiles[nearSource.Y][nearSource.X].Terrain; got != terrainEmpty {
		t.Fatalf("nearest cross-Plot terrain = %v, want consumed grass", got)
	}
	if got := state.gameMap.Home.Tiles[farSource.Y][farSource.X].Terrain; got != terrainTree {
		t.Fatalf("far terrain = %v, want Tree retained", got)
	}
}

func TestEconomicBuildingTerrainDistanceTieUsesCanonicalMapOrder(t *testing.T) {
	state := newRaidTestState(t)
	clearNaturalTerrain(&state.gameMap.Home)
	state.gameMap.Home.Tiles[homePlotCenter][homePlotCenter].Feature = featureWoodcutter
	left := homeTileCoordinate(homePlotCenter-1, homePlotCenter)
	right := homeTileCoordinate(homePlotCenter+1, homePlotCenter)
	state.gameMap.Home.Tiles[left.Y][left.X].Terrain = terrainTree
	state.gameMap.Home.Tiles[right.Y][right.X].Terrain = terrainTree

	state.grantEconomicBuildingResources()

	if got := state.gameMap.Home.Tiles[left.Y][left.X].Terrain; got != terrainEmpty {
		t.Fatalf("first equal-distance terrain = %v, want consumed grass", got)
	}
	if got := state.gameMap.Home.Tiles[right.Y][right.X].Terrain; got != terrainTree {
		t.Fatalf("second equal-distance terrain = %v, want Tree retained", got)
	}
}

func TestEconomicBuildingsCannotConsumeOneTerrainTileTwice(t *testing.T) {
	state := newRaidTestState(t)
	clearNaturalTerrain(&state.gameMap.Home)
	state.gameMap.Home.Tiles[5][5].Feature = featureWoodcutter
	state.gameMap.Home.Tiles[5][9].Feature = featureWoodcutter
	state.gameMap.Home.Tiles[5][6].Terrain = terrainTree
	startingWood := state.status.resources.wood

	state.grantEconomicBuildingResources()

	if got := state.status.resources.wood; got != startingWood+10 {
		t.Fatalf("wood = %d, want one terrain-backed yield %d", got, startingWood+10)
	}
}

func TestEconomicBuildingWithoutMatchingTerrainProducesNothingUntilLaterLabour(t *testing.T) {
	state := newRaidTestState(t)
	clearNaturalTerrain(&state.gameMap.Home)
	state.gameMap.Home.Tiles[5][5].Feature = featureStoneQuarry
	state.gameMap.Home.Tiles[5][6].Terrain = terrainTree
	startingResources := state.status.resources

	state.grantEconomicBuildingResources()

	if state.status.resources != startingResources {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, startingResources)
	}
	if got := state.gameMap.Home.Tiles[5][6].Terrain; got != terrainTree {
		t.Fatalf("nonmatching terrain = %v, want Tree retained", got)
	}

	state.gameMap.Home.Tiles[5][7].Terrain = terrainBoulder
	state.grantEconomicBuildingResources()
	if got := state.status.resources.stone; got != startingResources.stone+10 {
		t.Fatalf("later stone = %d, want %d", got, startingResources.stone+10)
	}
}

func TestConsumedTerrainUsesBiomeDefaultAndPreservesOtherTileData(t *testing.T) {
	tests := []struct {
		name  string
		biome plotBiome
	}{
		{name: "grasslands", biome: biomeGrasslands},
		{name: "hills", biome: biomeHills},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newRaidTestState(t)
			clearNaturalTerrain(&state.gameMap.Home)
			state.gameMap.Home.Biome = test.biome
			state.gameMap.Home.Tiles[5][5].Feature = featureWoodcutter
			source := homeTileCoordinate(6, 5)
			state.gameMap.Home.Tiles[source.Y][source.X] = Tile{
				Terrain: terrainTree,
				Feature: featureHouse,
				Tweak:   47,
			}

			state.grantEconomicBuildingResources()

			got := state.gameMap.Home.Tiles[source.Y][source.X]
			if got.Terrain != terrainEmpty || got.Feature != featureHouse || got.Tweak != 47 {
				t.Fatalf("consumed Tile = %+v, want grass with feature %v and tweak 47", got, featureHouse)
			}
		})
	}
}

func TestConsumingSelectedTerrainClearsSelection(t *testing.T) {
	state := newRaidTestState(t)
	clearNaturalTerrain(&state.gameMap.Home)
	state.gameMap.Home.Tiles[5][5].Feature = featureIronMine
	source := homeTileCoordinate(6, 5)
	state.gameMap.Home.Tiles[source.Y][source.X].Terrain = terrainIronDeposit
	state.selection = selectedItem{kind: selectedItemTerrain, tile: source}

	state.grantEconomicBuildingResources()

	if state.selection.kind != selectedItemNone {
		t.Fatalf("selection = %+v, want cleared", state.selection)
	}
}

func clearNaturalTerrain(plot *Plot) {
	for y := 0; y < plotSize; y++ {
		for x := 0; x < plotSize; x++ {
			plot.Tiles[y][x].Terrain = terrainEmpty
		}
	}
}
