package game

// beginPostRaidDay advances the Day, resolves its immediate Labour, and opens Management.
func (s *State) beginPostRaidDay() {
	s.status.day++
	s.status.phase = phaseLabour
	s.grantEconomicBuildingResources()
	s.status.phase = phaseManagement
}
