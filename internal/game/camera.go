package game

import "math"

const (
	cameraInitialZoom = 1.0
	cameraMinZoom     = 0.1
	cameraZoomFactor  = 1.1
	cameraPanSpeed    = 12.0
	plotBaseTileSize  = 54.0
)

type camera struct {
	zoom    float64
	centerX float64
	centerY float64
}

type cameraDragState struct {
	active      bool
	lastCursorX int
	lastCursorY int
}

type sceneViewport struct {
	x float64
	y float64
	w float64
	h float64
}

type projectedRect struct {
	x float32
	y float32
	w float32
	h float32
}

// newCamera creates the initial camera centered on the default home Plot.
func newCamera() camera {
	return camera{
		zoom:    cameraInitialZoom,
		centerX: 0,
		centerY: 0,
	}
}

// applyCameraInput updates zoom and pan without changing game simulation state.
func (s *State) applyCameraInput(input Input) {
	if input.WheelY != 0 {
		s.camera.zoom *= math.Pow(cameraZoomFactor, input.WheelY)
		if s.camera.zoom < cameraMinZoom {
			s.camera.zoom = cameraMinZoom
		}
	}

	distance := cameraPanSpeed / (plotBaseTileSize * s.camera.zoom)
	if input.Pan.Up {
		s.camera.centerY += distance
	}
	if input.Pan.Down {
		s.camera.centerY -= distance
	}
	if input.Pan.Left {
		s.camera.centerX -= distance
	}
	if input.Pan.Right {
		s.camera.centerX += distance
	}

	s.applyCameraDragInput(input)
}

// applyCameraDragInput updates map inspection panning from held right-drag input.
func (s *State) applyCameraDragInput(input Input) {
	if !s.cameraDrag.active && input.RightPressed && s.canStartCameraDrag(input.CursorX, input.CursorY) {
		s.cameraDrag = cameraDragState{
			active:      true,
			lastCursorX: input.CursorX,
			lastCursorY: input.CursorY,
		}
	}

	if !s.cameraDrag.active {
		return
	}
	if input.RightReleased || !input.RightDown {
		s.cameraDrag = cameraDragState{}
		return
	}

	deltaX := input.CursorX - s.cameraDrag.lastCursorX
	deltaY := input.CursorY - s.cameraDrag.lastCursorY
	scale := plotBaseTileSize * s.camera.zoom
	if scale > 0 {
		s.camera.centerX -= float64(deltaX) / scale
		s.camera.centerY += float64(deltaY) / scale
	}
	s.cameraDrag.lastCursorX = input.CursorX
	s.cameraDrag.lastCursorY = input.CursorY
}

// canStartCameraDrag reports whether a point begins over the map rather than UI.
func (s *State) canStartCameraDrag(x, y int) bool {
	viewport := s.sceneViewport()
	return viewportContainsPoint(viewport, x, y) &&
		!s.buildingBarContains(x, y) &&
		!s.nextRaidButtonContains(x, y) &&
		!s.selectionPanelContains(x, y) &&
		!s.marketControlsContains(x, y)
}

// sceneViewport returns the screen-space area used for camera-projected map rendering.
func (s *State) sceneViewport() sceneViewport {
	return sceneViewport{
		x: 0,
		y: topBarHeight,
		w: float64(s.ui.width),
		h: float64(s.ui.height - topBarHeight),
	}
}

// viewportContainsPoint reports whether a point is inside a scene viewport.
func viewportContainsPoint(viewport sceneViewport, x, y int) bool {
	screenX := float64(x)
	screenY := float64(y)
	return screenX >= viewport.x &&
		screenX <= viewport.x+viewport.w &&
		screenY >= viewport.y &&
		screenY <= viewport.y+viewport.h
}

// projectRect converts a world-space rectangle into a screen-space rectangle.
func (s *State) projectRect(viewport sceneViewport, worldWest, worldNorth, worldW, worldH float64) projectedRect {
	centerX := viewport.x + viewport.w/2
	centerY := viewport.y + viewport.h/2
	scale := plotBaseTileSize * s.camera.zoom
	return projectedRect{
		x: float32(centerX + (worldWest-s.camera.centerX)*scale),
		y: float32(centerY + (s.camera.centerY-worldNorth)*scale),
		w: float32(worldW * scale),
		h: float32(worldH * scale),
	}
}
