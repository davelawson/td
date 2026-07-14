package game

import (
	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	marketWoodGoldCost  = 1
	marketStoneGoldCost = 1
	marketIronGoldCost  = 3
	marketControlUIGap  = 10
)

// marketControlsVisible reports whether the selected Market can trade now.
func (s *State) marketControlsVisible() bool {
	if s.ui.menu.open || s.status.phase != phaseManagement || s.raid.active || s.raid.breached || s.selection.kind != selectedItemStructure {
		return false
	}
	plot, ok := s.gameMap.plot(s.selection.tile.Plot)
	if !ok || s.selection.tile.X < 0 || s.selection.tile.X >= plotSize || s.selection.tile.Y < 0 || s.selection.tile.Y >= plotSize {
		return false
	}
	if plot.Tiles[s.selection.tile.Y][s.selection.tile.X].Feature != featureMarket {
		return false
	}
	anchor := s.marketControlAnchor()
	centerX := anchor.X + anchor.W/2
	centerY := anchor.Y + anchor.H/2
	return viewportContainsPoint(s.sceneViewport(), centerX, centerY)
}

// marketControlAnchor returns the selected Market Tile in screen coordinates.
func (s *State) marketControlAnchor() ui.Button[int] {
	tile := s.selection.tile
	west, north, width, height := plotTileWorldRect(tile.Plot, tile.X, tile.Y)
	rect := s.projectRect(s.sceneViewport(), west, north, width, height)
	return ui.Button[int]{X: int(rect.x), Y: int(rect.y), W: int(rect.w), H: int(rect.h)}
}

// marketControlArea returns space clear of persistent Management UI.
func (s *State) marketControlArea() ui.Button[int] {
	right := s.ui.width
	if panel, ok := s.selectionPanelBounds(); ok {
		right = panel.X - marketControlUIGap
	}
	bottom := s.nextRaidButton().Y - marketControlUIGap
	return ui.Button[int]{
		X: ui.BuildingBarWidth,
		Y: topBarHeight,
		W: right - ui.BuildingBarWidth,
		H: bottom - topBarHeight,
	}
}

// marketControlsModel adapts current resources and icons for Market presentation.
func (s *State) marketControlsModel() ui.MarketControlsModel {
	return ui.MarketControlsModel{
		Items: []ui.MarketTradeItem{
			{
				Action:   ui.MarketTradeBuyWood,
				Label:    "+1 Wood · 1 Gold",
				Icon:     s.assetCatalog.Sprite.Icon.Wood,
				GoldCost: marketWoodGoldCost,
				Enabled:  s.status.resources.gold >= marketWoodGoldCost,
			},
			{
				Action:   ui.MarketTradeBuyStone,
				Label:    "+1 Stone · 1 Gold",
				Icon:     s.assetCatalog.Sprite.Icon.Stone,
				GoldCost: marketStoneGoldCost,
				Enabled:  s.status.resources.gold >= marketStoneGoldCost,
			},
			{
				Action:   ui.MarketTradeBuyIron,
				Label:    "+1 Iron · 3 Gold",
				Icon:     s.assetCatalog.Sprite.Icon.Iron,
				GoldCost: marketIronGoldCost,
				Enabled:  s.status.resources.gold >= marketIronGoldCost,
			},
		},
		Hovered: s.ui.marketTradeHover,
	}
}

// updateMarketControls updates hover state and applies one clicked purchase.
func (s *State) updateMarketControls(input Input) {
	s.ui.marketTradeHover = ui.MarketTradeNoAction
	if !s.marketControlsVisible() {
		return
	}
	anchor := s.marketControlAnchor()
	area := s.marketControlArea()
	model := s.marketControlsModel()
	item, ok := ui.MarketTradeAt(anchor, area, model, input.CursorX, input.CursorY)
	if !ok {
		return
	}
	s.ui.marketTradeHover = item.Action
	if input.Clicked && item.Enabled {
		s.buyMarketResource(item.Action, item.GoldCost)
	}
}

// buyMarketResource atomically spends Gold and grants one selected material.
func (s *State) buyMarketResource(action ui.MarketTradeAction, goldCost int) {
	if goldCost <= 0 || s.status.resources.gold < goldCost {
		return
	}
	s.status.resources.gold -= goldCost
	switch action {
	case ui.MarketTradeBuyWood:
		s.status.resources.wood++
	case ui.MarketTradeBuyStone:
		s.status.resources.stone++
	case ui.MarketTradeBuyIron:
		s.status.resources.iron++
	default:
		s.status.resources.gold += goldCost
	}
}

// marketControlsContains reports whether a point lies on visible Market UI.
func (s *State) marketControlsContains(x, y int) bool {
	if !s.marketControlsVisible() {
		return false
	}
	return ui.MarketControlsContains(
		s.marketControlAnchor(),
		s.marketControlArea(),
		s.marketControlsModel(),
		x,
		y,
	)
}

// drawMarketControls renders contextual purchases beside the selected Market.
func (s *State) drawMarketControls(screen *ebiten.Image) {
	if !s.marketControlsVisible() {
		return
	}
	ui.DrawMarketControls(
		screen,
		s.ui.costBoldFace,
		s.marketControlAnchor(),
		s.marketControlArea(),
		s.marketControlsModel(),
	)
}
