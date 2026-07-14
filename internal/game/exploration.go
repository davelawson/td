package game

type exploreButton struct {
	From   plotCoordinate
	Target plotCoordinate
	Center coord
	Biome  plotBiome
}

// updateExploration applies Management-phase Plot reveal clicks.
func (s *State) updateExploration(input Input) {
	if !input.Clicked || !s.canExploreNow() || s.exploreClickBlockedByUI(input.CursorX, input.CursorY) {
		return
	}
	target, ok := s.exploreTargetAtScreenPosition(input.CursorX, input.CursorY)
	if !ok {
		return
	}
	s.gameMap.revealPlot(target)
}

// canExploreNow reports whether the current game state allows revealing a Plot.
func (s *State) canExploreNow() bool {
	return s.status.phase == phaseManagement && !s.raid.active && !s.raid.breached
}

// exploreClickBlockedByUI reports whether screen-space UI owns the click.
func (s *State) exploreClickBlockedByUI(x, y int) bool {
	return s.buildingBarContains(x, y) ||
		s.nextRaidButtonContains(x, y) ||
		s.selectionPanelContains(x, y) ||
		s.marketControlsContains(x, y)
}

// exploreTargetAtScreenPosition returns the unexplored Plot targeted by an explore button click.
func (s *State) exploreTargetAtScreenPosition(x, y int) (plotCoordinate, bool) {
	viewport := s.sceneViewport()
	for _, button := range s.exploreButtons() {
		rect := s.projectRect(
			viewport,
			button.Center.X-exploreButtonSize/2,
			button.Center.Y+exploreButtonSize/2,
			exploreButtonSize,
			exploreButtonSize,
		)
		if rectContainsPoint(rect, x, y) {
			return button.Target, true
		}
	}
	return plotCoordinate{}, false
}

// exploreButtonContains reports whether a screen point is over a visible explore button.
func (s *State) exploreButtonContains(x, y int) bool {
	_, ok := s.exploreTargetAtScreenPosition(x, y)
	return ok
}

// exploreButtons returns border controls from explored Plots to unexplored orthogonal neighbors.
func (s *State) exploreButtons() []exploreButton {
	s.gameMap.ensurePlots()
	var buttons []exploreButton
	for _, from := range s.gameMap.exploredPlotCoordinates() {
		for _, target := range orthogonalPlotNeighbors(from) {
			if s.gameMap.explored(target) {
				continue
			}
			biome, ok := s.gameMap.frontierBiome(target)
			if !ok {
				continue
			}
			buttons = append(buttons, exploreButton{
				From:   from,
				Target: target,
				Center: exploreButtonCenter(from, target),
				Biome:  biome,
			})
		}
	}
	return buttons
}

// biomeLabel returns the player-facing name for a Plot biome.
func biomeLabel(biome plotBiome) string {
	switch biome {
	case biomeHills:
		return "Hills"
	case biomeForest:
		return "Forest"
	default:
		return "Grasslands"
	}
}

// exploreButtonCenter returns the shared border center for an exploration control.
func exploreButtonCenter(from, target plotCoordinate) coord {
	west, north, width, height := plotWorldRect(from)
	dx := target.X - from.X
	dy := target.Y - from.Y
	switch {
	case dx == 1:
		return coord{X: west + width, Y: north - height/2}
	case dx == -1:
		return coord{X: west, Y: north - height/2}
	case dy == 1:
		return coord{X: west + width/2, Y: north}
	case dy == -1:
		return coord{X: west + width/2, Y: north - height}
	default:
		return coord{}
	}
}
