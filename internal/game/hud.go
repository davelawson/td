package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	topBarHeight = 86
	topBarMargin = 42
)

type phase int

const (
	phaseCalm phase = iota
	phaseRaid
)

type resourceCounts struct {
	wood  int
	stone int
	metal int
}

var (
	topBarColor     = uiColor(26, 32, 28, 245)
	topBarEdgeColor = uiColor(134, 114, 65, 210)
)

// setPrototypeGameStatus initializes fixed state shown before gameplay systems exist.
func (s *State) setPrototypeGameStatus() {
	s.phase = phaseCalm
	s.chapter = "Chapter I: The Ashen Copse"
	s.day = 1
	s.calmTime = 120
	s.raidCount = 12
	s.barricade = 3
	s.resources = resourceCounts{
		wood:  80,
		stone: 45,
		metal: 12,
	}
}

// chapterDayText formats the Chapter and Day summary for the top bar.
func (s *State) chapterDayText() string {
	return fmt.Sprintf("%s | Day %d", s.chapter, s.day)
}

// phaseText formats the phase-specific top bar status.
func (s *State) phaseText() string {
	switch s.phase {
	case phaseRaid:
		return fmt.Sprintf("Enemies remaining: %d", s.raidCount)
	default:
		minutes := s.calmTime / 60
		seconds := s.calmTime % 60
		return fmt.Sprintf("Raid in %02d:%02d", minutes, seconds)
	}
}

// resourcesAndBarricadeText formats economy and Sanctum defense status.
func (s *State) resourcesAndBarricadeText() string {
	return fmt.Sprintf(
		"Wood %d  Stone %d  Metal %d | Barricade %d",
		s.resources.wood,
		s.resources.stone,
		s.resources.metal,
		s.barricade,
	)
}

// drawTopBar renders the game status bar at the top of the screen.
func (s *State) drawTopBar(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(s.width), topBarHeight, topBarColor, false)
	vector.StrokeLine(screen, 0, topBarHeight-2, float32(s.width), topBarHeight-2, 3, topBarEdgeColor, false)

	left := s.chapterDayText()
	center := s.phaseText()
	right := s.resourcesAndBarricadeText()

	s.drawText(screen, left, s.hudFace, topBarMargin, 29, textColor)

	centerWidth, _ := text.Measure(center, s.hudFace, s.hudFace.Size)
	s.drawText(screen, center, s.hudFace, (float64(s.width)-centerWidth)/2, 29, pauseColor)

	rightWidth, _ := text.Measure(right, s.hudFace, s.hudFace.Size)
	s.drawText(screen, right, s.hudFace, float64(s.width)-rightWidth-topBarMargin, 29, textColor)
}

// uiColor builds an RGBA color for HUD-specific translucent surfaces.
func uiColor(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}
