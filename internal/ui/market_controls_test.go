package ui

import "testing"

// TestMarketControlButtonsPreferRightAndFallBackLeft verifies anchor-relative placement.
func TestMarketControlButtonsPreferRightAndFallBackLeft(t *testing.T) {
	area := Button[int]{X: 260, Y: 86, W: 1200, H: 800}
	model := testMarketControlsModel()
	rightAnchor := Button[int]{X: 600, Y: 300, W: 54, H: 54}
	right := MarketControlButtons(rightAnchor, area, model)
	if len(right) != 3 || right[0].X != 664 || right[0].Y != 300 {
		t.Fatalf("right-side buttons = %+v, want first at 664,300", right)
	}

	leftAnchor := Button[int]{X: 1400, Y: 300, W: 54, H: 54}
	left := MarketControlButtons(leftAnchor, area, model)
	if len(left) != 3 || left[0].X != 1170 || left[0].Y != 300 {
		t.Fatalf("left-side buttons = %+v, want first at 1170,300", left)
	}
}

// TestMarketControlButtonsClampInsideArea verifies controls stay clear of fixed UI.
func TestMarketControlButtonsClampInsideArea(t *testing.T) {
	area := Button[int]{X: 260, Y: 86, W: 1200, H: 800}
	model := testMarketControlsModel()
	buttons := MarketControlButtons(Button[int]{X: 200, Y: 900, W: 54, H: 54}, area, model)
	if len(buttons) != 3 {
		t.Fatalf("buttons = %d, want 3", len(buttons))
	}
	for _, button := range buttons {
		if !area.Contains(button.X, button.Y) || !area.Contains(button.X+button.W-1, button.Y+button.H-1) {
			t.Fatalf("button %+v should fit area %+v", button, area)
		}
	}
}

// TestMarketTradeAtIncludesDisabledButtons verifies disabled clicks remain owned by Market UI.
func TestMarketTradeAtIncludesDisabledButtons(t *testing.T) {
	area := Button[int]{X: 260, Y: 86, W: 1200, H: 800}
	anchor := Button[int]{X: 600, Y: 300, W: 54, H: 54}
	model := testMarketControlsModel()
	model.Items[2].Enabled = false
	buttons := MarketControlButtons(anchor, area, model)
	item, ok := MarketTradeAt(anchor, area, model, buttons[2].X+1, buttons[2].Y+1)
	if !ok || item.Action != MarketTradeBuyIron || item.Enabled {
		t.Fatalf("disabled hit = %+v, available %v, want disabled Iron", item, ok)
	}
	if !MarketControlsContains(anchor, area, model, buttons[2].X+1, buttons[2].Y+1) {
		t.Fatal("expected disabled button to remain inside Market UI")
	}
}

// testMarketControlsModel returns the stable three-action control model.
func testMarketControlsModel() MarketControlsModel {
	return MarketControlsModel{
		Items: []MarketTradeItem{
			{Action: MarketTradeBuyWood, Label: "+1 Wood · 1 Gold", GoldCost: 1, Enabled: true},
			{Action: MarketTradeBuyStone, Label: "+1 Stone · 1 Gold", GoldCost: 1, Enabled: true},
			{Action: MarketTradeBuyIron, Label: "+1 Iron · 3 Gold", GoldCost: 3, Enabled: true},
		},
		Hovered: MarketTradeNoAction,
	}
}
