package game

import (
	"testing"

	"td/internal/ui"
)

// TestMarketTemplateDefinesAcceptedEconomyFacts verifies catalog metadata and sprite wiring.
func TestMarketTemplateDefinesAcceptedEconomyFacts(t *testing.T) {
	state := newRaidTestState(t)
	market := state.structureCatalog.Market
	if market.Name != "Market" || market.Description == "" {
		t.Fatalf("Market identity = %q/%q", market.Name, market.Description)
	}
	if market.Sprite == nil || market.Sprite != state.assetCatalog.Sprite.Structure.Market {
		t.Fatal("Market should use the loaded structure sprite")
	}
	if market.Cost != (Resources{Wood: 30}) {
		t.Fatalf("Market cost = %+v, want 30 Wood", market.Cost)
	}
	if market.Staffing != (StaffingRequirements{Soldiers: 1, Peasants: 2}) {
		t.Fatalf("Market staffing = %+v, want 1 Soldier and 2 Peasants", market.Staffing)
	}
	if market.ResourceYield != (Resources{}) {
		t.Fatalf("Market Labour yield = %+v, want none", market.ResourceYield)
	}
}

// TestMarketConstructionRequiresAllCapacity verifies Wood and both staff roles gate building.
func TestMarketConstructionRequiresAllCapacity(t *testing.T) {
	tests := []struct {
		name      string
		wood      int
		soldiers  int
		peasants  int
		buildable bool
	}{
		{name: "all requirements", wood: 30, soldiers: 1, peasants: 2, buildable: true},
		{name: "missing Wood", wood: 29, soldiers: 1, peasants: 2},
		{name: "missing Soldier", wood: 30, soldiers: 0, peasants: 2},
		{name: "missing Peasant", wood: 30, soldiers: 1, peasants: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := newRaidTestState(t)
			state.status.resources.wood = test.wood
			setAvailablePopulations(state, 0, test.soldiers, test.peasants)
			if got := state.canConstructBuilding(buildingBarMarketIndex); got != test.buildable {
				t.Fatalf("Market buildable = %v, want %v", got, test.buildable)
			}
		})
	}
}

// TestMarketPlacementDeductsWoodAndReservesStaff verifies atomic construction effects.
func TestMarketPlacementDeductsWoodAndReservesStaff(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 1, 2)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	setHomeTilesEmpty(state, tile)

	state.Update(pressBuildingBarItemInput(state, buildingBarMarketIndex))
	state.Update(releaseTileInput(state, tile.X, tile.Y))

	if state.gameMap.Home.Tiles[tile.Y][tile.X].Feature != featureMarket {
		t.Fatalf("feature = %v, want Market", state.gameMap.Home.Tiles[tile.Y][tile.X].Feature)
	}
	if state.status.resources.wood != 70 || state.status.resources.stone != 50 || state.status.resources.iron != 20 || state.status.resources.gold != 0 {
		t.Fatalf("resources = %+v, want 70 Wood and other initial totals", state.status.resources)
	}
	if state.status.populations.soldiers != (populationCount{available: 0, total: 1}) ||
		state.status.populations.peasants != (populationCount{available: 0, total: 2}) {
		t.Fatalf("populations = %+v, want reserved Soldier and Peasants", state.status.populations)
	}
}

// TestTwoMarketsReserveIndependentStaff verifies Market crews cannot be reused.
func TestTwoMarketsReserveIndependentStaff(t *testing.T) {
	state := newRaidTestState(t)
	setAvailablePopulations(state, 0, 2, 4)
	first := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	second := tileCoordinate{X: homePlotCenter + 3, Y: 5}
	setHomeTilesEmpty(state, first, second)

	for _, tile := range []tileCoordinate{first, second} {
		state.Update(pressBuildingBarItemInput(state, buildingBarMarketIndex))
		state.Update(releaseTileInput(state, tile.X, tile.Y))
	}

	if state.status.resources.wood != 40 || state.status.populations.soldiers.available != 0 || state.status.populations.peasants.available != 0 {
		t.Fatalf("two-Market capacity = resources %+v populations %+v", state.status.resources, state.status.populations)
	}
}

// TestMarketDoesNotProduceDuringLabour verifies trading is not a passive payout.
func TestMarketDoesNotProduceDuringLabour(t *testing.T) {
	state := selectedMarketTestState(t, 4)
	before := state.status.resources
	state.beginPostRaidDay()
	if state.status.resources != before {
		t.Fatalf("Market Labour resources = %+v, want unchanged %+v", state.status.resources, before)
	}
}

// TestMarketTradesBuyOneUnitAtAcceptedPrices verifies all material transactions.
func TestMarketTradesBuyOneUnitAtAcceptedPrices(t *testing.T) {
	tests := []struct {
		name   string
		action ui.MarketTradeAction
		cost   int
		check  func(resourceCounts) int
	}{
		{name: "Wood", action: ui.MarketTradeBuyWood, cost: 1, check: func(r resourceCounts) int { return r.wood }},
		{name: "Stone", action: ui.MarketTradeBuyStone, cost: 1, check: func(r resourceCounts) int { return r.stone }},
		{name: "Iron", action: ui.MarketTradeBuyIron, cost: 3, check: func(r resourceCounts) int { return r.iron }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := selectedMarketTestState(t, test.cost)
			before := state.status.resources
			state.Update(marketTradeClickInput(t, state, test.action))
			if got, want := test.check(state.status.resources), test.check(before)+1; got != want {
				t.Fatalf("%s = %d, want %d", test.name, got, want)
			}
			if state.status.resources.gold != 0 {
				t.Fatalf("Gold = %d, want exact-balance 0", state.status.resources.gold)
			}
			if state.selection.kind != selectedItemStructure {
				t.Fatalf("selection kind = %v, want selected structure", state.selection.kind)
			}
			if state.gameMap.Home.Tiles[state.selection.tile.Y][state.selection.tile.X].Feature != featureMarket {
				t.Fatal("Market trade click should preserve Market selection")
			}
			if state.raid.active {
				t.Fatal("Market trade click should not pass through to Next Raid")
			}
		})
	}
}

// TestDisabledMarketTradeConsumesClickWithoutMutation verifies insufficient Gold is a no-op.
func TestDisabledMarketTradeConsumesClickWithoutMutation(t *testing.T) {
	state := selectedMarketTestState(t, 2)
	before := state.status.resources
	state.Update(marketTradeClickInput(t, state, ui.MarketTradeBuyIron))
	if state.status.resources != before {
		t.Fatalf("resources = %+v, want unchanged %+v", state.status.resources, before)
	}
	if state.selection.kind != selectedItemStructure {
		t.Fatalf("selection = %+v, want preserved Market", state.selection)
	}
}

// TestMarketControlsFollowAcceptedPhaseRules verifies availability across game states.
func TestMarketControlsFollowAcceptedPhaseRules(t *testing.T) {
	state := selectedMarketTestState(t, 3)
	if !state.marketControlsVisible() {
		t.Fatal("expected Market controls during Management")
	}
	state.paused = true
	if !state.marketControlsVisible() {
		t.Fatal("expected Market controls during paused Management")
	}
	state.status.phase = phaseLabour
	if state.marketControlsVisible() {
		t.Fatal("Labour should hide Market controls")
	}
	state.status.phase = phaseRaid
	state.raid.active = true
	if state.marketControlsVisible() {
		t.Fatal("Raid should hide Market controls")
	}
	state.status.phase = phaseManagement
	state.raid.active = false
	state.raid.breached = true
	if state.marketControlsVisible() {
		t.Fatal("breach should hide Market controls")
	}
}

// TestPausedManagementMarketTradeSucceeds verifies pause freezes simulation but not trading.
func TestPausedManagementMarketTradeSucceeds(t *testing.T) {
	state := selectedMarketTestState(t, 1)
	state.paused = true
	state.Update(marketTradeClickInput(t, state, ui.MarketTradeBuyWood))
	if state.status.resources.gold != 0 || state.status.resources.wood != 101 {
		t.Fatalf("paused trade resources = %+v, want +1 Wood and 0 Gold", state.status.resources)
	}
}

// TestIngameMenuBlocksMarketTrade verifies overlay input cannot mutate the economy.
func TestIngameMenuBlocksMarketTrade(t *testing.T) {
	state := selectedMarketTestState(t, 3)
	click := marketTradeClickInput(t, state, ui.MarketTradeBuyIron)
	state.Update(Input{ToggleMenu: true})
	if state.marketControlsVisible() {
		t.Fatal("overlay should hide Market controls")
	}
	before := state.status.resources
	state.Update(click)
	if state.status.resources != before {
		t.Fatalf("overlay trade resources = %+v, want unchanged %+v", state.status.resources, before)
	}
}

// TestMarketControlsTrackCameraAndHideOffscreen verifies world anchoring and visibility.
func TestMarketControlsTrackCameraAndHideOffscreen(t *testing.T) {
	state := selectedMarketTestState(t, 3)
	firstX := state.marketControlAnchor().X
	state.camera.centerX++
	secondX := state.marketControlAnchor().X
	if secondX >= firstX {
		t.Fatalf("anchor X moved from %d to %d, want left after camera pans right", firstX, secondX)
	}
	state.camera.centerX = 1000
	if state.marketControlsVisible() {
		t.Fatal("off-screen selected Market should hide controls")
	}
}

// selectedMarketTestState returns a Management state with one selected Market.
func selectedMarketTestState(t *testing.T, gold int) *State {
	t.Helper()
	state := newRaidTestState(t)
	tile := tileCoordinate{X: homePlotCenter + 2, Y: 5}
	state.gameMap.Home.Tiles[tile.Y][tile.X] = Tile{Terrain: terrainEmpty, Feature: featureMarket}
	state.selection = selectedItem{kind: selectedItemStructure, tile: tile}
	state.status.resources.gold = gold
	return state
}

// marketTradeClickInput returns a click at the center of one Market action.
func marketTradeClickInput(t *testing.T, state *State, action ui.MarketTradeAction) Input {
	t.Helper()
	buttons := ui.MarketControlButtons(state.marketControlAnchor(), state.marketControlArea(), state.marketControlsModel())
	for _, button := range buttons {
		if button.Action == action {
			return Input{CursorX: button.X + button.W/2, CursorY: button.Y + button.H/2, Clicked: true}
		}
	}
	t.Fatalf("missing Market action %v", action)
	return Input{}
}
