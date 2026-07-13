package game

import (
	"td/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// currentSelectionPanel returns UI-facing data for the currently selected object.
func (s *State) currentSelectionPanel() (ui.SelectionPanelData, bool) {
	switch s.selection.kind {
	case selectedItemRaider:
		return s.selectedRaiderPanel()
	case selectedItemStructure:
		return s.selectedStructurePanel()
	default:
		return ui.SelectionPanelData{}, false
	}
}

// selectedRaiderPanel returns panel data for the selected active raider.
func (s *State) selectedRaiderPanel() (ui.SelectionPanelData, bool) {
	for _, enemy := range s.raid.enemies {
		if enemy.id != s.selection.raiderID || enemy.health <= 0 {
			continue
		}

		data := ui.SelectionPanelData{
			Kind:   ui.SelectionPanelRaider,
			Health: enemy.health,
		}
		if enemy.template != nil {
			data.Name = enemy.template.Name
			data.MaxHealth = enemy.template.MaxHealth
			data.SpeedTilesPerSecond = enemy.template.SpeedTilesPerSecond
			data.SanctumDamage = enemy.template.SanctumDamage
		}

		return data, true
	}
	return ui.SelectionPanelData{}, false
}

// selectedStructurePanel returns panel data for the selected structure tile.
func (s *State) selectedStructurePanel() (ui.SelectionPanelData, bool) {
	tile := s.selection.tile
	if tile.X < 0 || tile.X >= plotSize || tile.Y < 0 || tile.Y >= plotSize {
		return ui.SelectionPanelData{}, false
	}
	plot, ok := s.gameMap.plot(tile.Plot)
	if !ok {
		return ui.SelectionPanelData{}, false
	}

	switch plot.Tiles[tile.Y][tile.X].Feature {
	case featureSanctum:
		return structureSelectionPanel(s.structureCatalog.Sanctum), true
	case featureHouse:
		return populationBuildingSelectionPanel(s.structureCatalog.House), true
	case featureBarracks:
		return populationBuildingSelectionPanel(s.structureCatalog.Barracks), true
	case featureDorm:
		return populationBuildingSelectionPanel(s.structureCatalog.Dorm), true
	case featureWoodcutter:
		return economicBuildingSelectionPanel(s.structureCatalog.Woodcutter), true
	case featureStoneQuarry:
		return economicBuildingSelectionPanel(s.structureCatalog.StoneQuarry), true
	case featureIronMine:
		return economicBuildingSelectionPanel(s.structureCatalog.IronMine), true
	case featureBowTower:
		return towerSelectionPanel(s.structureCatalog.BowTower), true
	case featureFlameBoltTower:
		return towerSelectionPanel(s.structureCatalog.FlameBoltTower), true
	case featureCatapultTower:
		return towerSelectionPanel(s.structureCatalog.CatapultTower), true
	default:
		return ui.SelectionPanelData{}, false
	}
}

// structureSelectionPanel returns panel data for a basic structure.
func structureSelectionPanel(template StructureTemplate) ui.SelectionPanelData {
	return ui.SelectionPanelData{
		Kind: ui.SelectionPanelStructure,
		Name: template.Name,
	}
}

// populationBuildingSelectionPanel returns panel data for population buildings.
func populationBuildingSelectionPanel(template StructureTemplate) ui.SelectionPanelData {
	return ui.SelectionPanelData{
		Kind:            ui.SelectionPanelPopulationBuilding,
		Name:            template.Name,
		Cost:            resourceAmounts(template.Cost),
		PopulationCost:  populationCostAmounts(template.PopulationCost),
		PopulationGrant: populationGrantAmounts(template.PopulationGrant),
	}
}

// economicBuildingSelectionPanel returns panel data for resource-producing buildings.
func economicBuildingSelectionPanel(template StructureTemplate) ui.SelectionPanelData {
	return ui.SelectionPanelData{
		Kind:          ui.SelectionPanelEconomicBuilding,
		Name:          template.Name,
		Cost:          resourceAmounts(template.Cost),
		Staffing:      staffingAmounts(template.Staffing),
		ResourceYield: resourceAmounts(template.ResourceYield),
	}
}

// towerSelectionPanel returns panel data for one combat tower template.
func towerSelectionPanel(template StructureTemplate) ui.SelectionPanelData {
	return ui.SelectionPanelData{
		Kind:                ui.SelectionPanelTower,
		Name:                template.Name,
		Staffing:            staffingAmounts(template.Staffing),
		RangeTiles:          template.RangeTiles,
		FireIntervalSeconds: template.FireIntervalSeconds,
		Damage:              template.Damage,
	}
}

// selectionPanelBounds returns the current bottom-right panel rectangle.
func (s *State) selectionPanelBounds() (ui.Button[int], bool) {
	panel, ok := s.currentSelectionPanel()
	if !ok {
		return ui.Button[int]{}, false
	}
	return ui.SelectionPanelBounds(s.ui.width, s.ui.height, panel)
}

// selectionPanelContains reports whether a screen point is inside the visible selection panel.
func (s *State) selectionPanelContains(x, y int) bool {
	panel, ok := s.currentSelectionPanel()
	return ok && ui.SelectionPanelContains(s.ui.width, s.ui.height, panel, x, y)
}

// drawSelectionPanel renders the current selected-object detail panel.
func (s *State) drawSelectionPanel(screen *ebiten.Image) {
	panel, ok := s.currentSelectionPanel()
	if !ok {
		return
	}
	ui.DrawSelectionPanel(screen, s.ui.hudFace, s.ui.width, s.ui.height, panel)
}

func resourceAmounts(resources Resources) ui.ResourceAmounts {
	return ui.ResourceAmounts{
		Wood:  resources.Wood,
		Stone: resources.Stone,
		Metal: resources.Metal,
	}
}

func staffingAmounts(staffing StaffingRequirements) ui.PopulationAmounts {
	return ui.PopulationAmounts{
		Apprentices: staffing.Apprentices,
		Soldiers:    staffing.Soldiers,
		Peasants:    staffing.Peasants,
	}
}

func populationCostAmounts(cost PopulationCost) ui.PopulationAmounts {
	return ui.PopulationAmounts{
		Apprentices: cost.Apprentices,
		Soldiers:    cost.Soldiers,
		Peasants:    cost.Peasants,
	}
}

func populationGrantAmounts(grant PopulationGrant) ui.PopulationAmounts {
	return ui.PopulationAmounts{
		Apprentices: grant.Apprentices,
		Soldiers:    grant.Soldiers,
		Peasants:    grant.Peasants,
	}
}
