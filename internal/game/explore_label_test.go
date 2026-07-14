package game

import "testing"

// TestExploreBiomeLabelRectUsesOutwardDirection verifies labels stay outside explored land.
func TestExploreBiomeLabelRectUsesOutwardDirection(t *testing.T) {
	buttonRect := projectedRect{x: 100, y: 100, w: 40, h: 40}
	tests := []struct {
		name   string
		target plotCoordinate
		check  func(projectedRect) bool
	}{
		{name: "north", target: plotCoordinate{X: 0, Y: 1}, check: func(label projectedRect) bool { return label.y+label.h < buttonRect.y }},
		{name: "east", target: plotCoordinate{X: 1, Y: 0}, check: func(label projectedRect) bool { return label.x > buttonRect.x+buttonRect.w }},
		{name: "south", target: plotCoordinate{X: 0, Y: -1}, check: func(label projectedRect) bool { return label.y > buttonRect.y+buttonRect.h }},
		{name: "west", target: plotCoordinate{X: -1, Y: 0}, check: func(label projectedRect) bool { return label.x+label.w < buttonRect.x }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			label := exploreBiomeLabelRect(exploreButton{Target: test.target}, buttonRect, 80, 16)
			if !test.check(label) {
				t.Fatalf("label rect = %+v, want outward from button %+v", label, buttonRect)
			}
		})
	}
}

// TestClickingExploreBiomeLabelDoesNotReveal verifies labels remain informational.
func TestClickingExploreBiomeLabelDoesNotReveal(t *testing.T) {
	state := newRaidTestState(t)
	target := plotCoordinate{X: 1, Y: 0}
	var button exploreButton
	found := false
	for _, candidate := range state.exploreButtons() {
		if candidate.Target == target {
			button = candidate
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected east explore button")
	}
	buttonRect := state.projectRect(
		state.sceneViewport(),
		button.Center.X-exploreButtonSize/2,
		button.Center.Y+exploreButtonSize/2,
		exploreButtonSize,
		exploreButtonSize,
	)
	labelRect := state.projectedExploreBiomeLabelRect(button, buttonRect)

	state.Update(Input{
		CursorX: int(labelRect.x + labelRect.w/2),
		CursorY: int(labelRect.y + labelRect.h/2),
		Clicked: true,
	})

	if state.gameMap.explored(target) {
		t.Fatal("expected biome label click not to reveal its Plot")
	}
}
