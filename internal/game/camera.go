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
