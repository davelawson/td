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
	center := float64(plotSize) * plotBaseTileSize / 2
	return camera{
		zoom:    cameraInitialZoom,
		centerX: center,
		centerY: center,
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

	distance := cameraPanSpeed / s.camera.zoom
	if input.PanUp {
		s.camera.centerY -= distance
	}
	if input.PanDown {
		s.camera.centerY += distance
	}
	if input.PanLeft {
		s.camera.centerX -= distance
	}
	if input.PanRight {
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
func (s *State) projectRect(viewport sceneViewport, worldX, worldY, worldW, worldH float64) projectedRect {
	centerX := viewport.x + viewport.w/2
	centerY := viewport.y + viewport.h/2
	return projectedRect{
		x: float32(centerX + (worldX-s.camera.centerX)*s.camera.zoom),
		y: float32(centerY + (worldY-s.camera.centerY)*s.camera.zoom),
		w: float32(worldW * s.camera.zoom),
		h: float32(worldH * s.camera.zoom),
	}
}
