package game

import (
	"testing"

	"td/internal/ui"
)

// TestCameraRightDragStartsOverGameView verifies map-space right presses begin dragging.
func TestCameraRightDragStartsOverGameView(t *testing.T) {
	state := newCameraDragTestState(t)

	state.Update(rightDragInput(stateCenterX(state), stateCenterY(state), true, true, false))

	if !state.cameraDrag.active {
		t.Fatal("expected camera drag to start over the game view")
	}
}

// TestCameraRightDragGrabsWorld verifies dragged map content follows the cursor.
func TestCameraRightDragGrabsWorld(t *testing.T) {
	state := newCameraDragTestState(t)
	startX := state.camera.centerX
	startY := state.camera.centerY
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(cursorX+108, cursorY+54, false, true, false))

	if got, want := state.camera.centerX-startX, -2.0; !almostEqual(got, want) {
		t.Fatalf("camera x delta = %f, want %f", got, want)
	}
	if got, want := state.camera.centerY-startY, 1.0; !almostEqual(got, want) {
		t.Fatalf("camera y delta = %f, want %f", got, want)
	}
}

// TestCameraRightDragScalesWithZoom verifies drag speed follows visible world scale.
func TestCameraRightDragScalesWithZoom(t *testing.T) {
	state := newCameraDragTestState(t)
	state.camera.zoom = 2
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(cursorX+108, cursorY, false, true, false))

	if got, want := state.camera.centerX, -1.0; !almostEqual(got, want) {
		t.Fatalf("camera center x = %f, want %f", got, want)
	}
}

// TestCameraRightDragStartBlockedByUI verifies screen-space UI does not start map drags.
func TestCameraRightDragStartBlockedByUI(t *testing.T) {
	tests := []struct {
		name string
		x    func(*State) int
		y    func(*State) int
	}{
		{
			name: "top bar",
			x:    func(s *State) int { return stateCenterX(s) },
			y:    func(*State) int { return topBarHeight / 2 },
		},
		{
			name: "building bar",
			x:    func(*State) int { return ui.BuildingBarWidth / 2 },
			y:    func(*State) int { return topBarHeight + 40 },
		},
		{
			name: "next raid button",
			x:    func(s *State) int { return s.nextRaidButton().X + s.nextRaidButton().W/2 },
			y:    func(s *State) int { return s.nextRaidButton().Y + s.nextRaidButton().H/2 },
		},
		{
			name: "selection panel",
			x: func(s *State) int {
				panel, _ := s.selectionPanelBounds()
				return panel.X + panel.W/2
			},
			y: func(s *State) int {
				panel, _ := s.selectionPanelBounds()
				return panel.Y + panel.H/2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := newCameraDragTestState(t)
			state.selection = selectedItem{
				kind: selectedItemStructure,
				tile: tileCoordinate{X: homePlotCenter, Y: homePlotCenter},
			}

			state.Update(rightDragInput(tt.x(state), tt.y(state), true, true, false))

			if state.cameraDrag.active {
				t.Fatal("expected UI press not to start camera drag")
			}
		})
	}
}

// TestCameraRightDragContinuesOverUI verifies UI only blocks the initial press.
func TestCameraRightDragContinuesOverUI(t *testing.T) {
	state := newCameraDragTestState(t)
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(ui.BuildingBarWidth/2, topBarHeight+40, false, true, false))

	if !state.cameraDrag.active {
		t.Fatal("expected drag to remain active over UI after a valid start")
	}
	if state.camera.centerX == 0 && state.camera.centerY == 0 {
		t.Fatal("expected camera to move while dragging over UI")
	}
}

// TestCameraRightDragReleaseClearsState verifies button release ends panning.
func TestCameraRightDragReleaseClearsState(t *testing.T) {
	state := newCameraDragTestState(t)
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(cursorX+10, cursorY, false, false, true))

	if state.cameraDrag.active {
		t.Fatal("expected right release to clear camera drag")
	}
}

// TestCameraRightDragWorksWhilePaused verifies inspection remains available while paused.
func TestCameraRightDragWorksWhilePaused(t *testing.T) {
	state := newCameraDragTestState(t)
	state.Update(Input{TogglePause: true})
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(cursorX+54, cursorY, false, true, false))

	if state.camera.centerX >= 0 {
		t.Fatalf("camera center x = %f, want negative after drag", state.camera.centerX)
	}
	if state.Updates() != 0 {
		t.Fatalf("updates = %d, want 0 while paused", state.Updates())
	}
}

// TestIngameMenuBlocksCameraRightDrag verifies overlay-open frames ignore drag input.
func TestIngameMenuBlocksCameraRightDrag(t *testing.T) {
	state := newCameraDragTestState(t)
	state.Update(Input{ToggleMenu: true})
	startCamera := state.camera

	state.Update(rightDragInput(stateCenterX(state), stateCenterY(state), true, true, false))
	state.Update(rightDragInput(stateCenterX(state)+54, stateCenterY(state), false, true, false))

	if state.camera != startCamera {
		t.Fatalf("camera = %+v, want unchanged %+v", state.camera, startCamera)
	}
	if state.cameraDrag.active {
		t.Fatal("expected overlay to block camera drag state")
	}
}

// TestCameraRightDragDoesNotChangeHomePlot verifies panning is inspection only.
func TestCameraRightDragDoesNotChangeHomePlot(t *testing.T) {
	state := newCameraDragTestState(t)
	initial := state.gameMap
	cursorX := stateCenterX(state)
	cursorY := stateCenterY(state)

	state.Update(rightDragInput(cursorX, cursorY, true, true, false))
	state.Update(rightDragInput(cursorX+54, cursorY+54, false, true, false))

	if !mapsEqual(state.gameMap, initial) {
		t.Fatal("expected right-drag camera movement to leave the map unchanged")
	}
}

// newCameraDragTestState creates a game state for camera drag tests.
func newCameraDragTestState(t *testing.T) *State {
	t.Helper()
	state, err := New("Merlin", 1920, 1080)
	if err != nil {
		t.Fatal(err)
	}
	return state
}

// stateCenterX returns the horizontal midpoint of the game drawable area.
func stateCenterX(state *State) int {
	return state.ui.width / 2
}

// stateCenterY returns the vertical midpoint of the camera scene.
func stateCenterY(state *State) int {
	return topBarHeight + (state.ui.height-topBarHeight)/2
}

// rightDragInput creates one frame of right-mouse camera drag input.
func rightDragInput(x, y int, pressed, down, released bool) Input {
	return Input{
		CursorX:       x,
		CursorY:       y,
		RightPressed:  pressed,
		RightDown:     down,
		RightReleased: released,
	}
}
