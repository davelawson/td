package game

import (
	"fmt"

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

type gameStatus struct {
	phase     phase
	chapter   string
	day       int
	calmTime  int
	raidCount int
	barricade int
	resources resourceCounts
}

// setPrototypeGameStatus initializes fixed state shown before gameplay systems exist.
func (s *State) setPrototypeGameStatus() {
	s.status = gameStatus{
		phase:     phaseCalm,
		chapter:   "Chapter I: The Ashen Copse",
		day:       1,
		calmTime:  120,
		raidCount: 12,
		barricade: 3,
		resources: resourceCounts{
			wood:  80,
			stone: 45,
			metal: 12,
		},
	}
}

// chapterDayText formats the Chapter and Day summary for the top bar.
func (s *State) chapterDayText() string {
	return fmt.Sprintf("%s | Day %d", s.status.chapter, s.status.day)
}

// phaseText formats the phase-specific top bar status.
func (s *State) phaseText() string {
	switch s.status.phase {
	case phaseRaid:
		return fmt.Sprintf("Enemies remaining: %d", s.status.raidCount)
	default:
		minutes := s.status.calmTime / 60
		seconds := s.status.calmTime % 60
		return fmt.Sprintf("Raid in %02d:%02d", minutes, seconds)
	}
}

// resourcesAndBarricadeText formats economy and Sanctum defense status.
func (s *State) resourcesAndBarricadeText() string {
	return fmt.Sprintf(
		"Wood %d  Stone %d  Metal %d | Barricade %d",
		s.status.resources.wood,
		s.status.resources.stone,
		s.status.resources.metal,
		s.status.barricade,
	)
}

// drawTopBar renders the game status bar at the top of the screen.
func (s *State) drawTopBar(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(s.ui.width), topBarHeight, topBarColor, false)
	vector.StrokeLine(screen, 0, topBarHeight-2, float32(s.ui.width), topBarHeight-2, 3, topBarEdgeColor, false)

	left := s.chapterDayText()
	center := s.phaseText()
	right := s.resourcesAndBarricadeText()

	s.drawText(screen, left, s.ui.hudFace, topBarMargin, 29, textColor)

	centerWidth, _ := text.Measure(center, s.ui.hudFace, s.ui.hudFace.Size)
	s.drawText(screen, center, s.ui.hudFace, (float64(s.ui.width)-centerWidth)/2, 29, pauseColor)

	rightWidth, _ := text.Measure(right, s.ui.hudFace, s.ui.hudFace.Size)
	s.drawText(screen, right, s.ui.hudFace, float64(s.ui.width)-rightWidth-topBarMargin, 29, textColor)
}
