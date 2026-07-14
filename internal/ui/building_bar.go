package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BuildingBarWidth           = 260
	buildingBarPadding         = 16
	buildingBarTabHeight       = 28
	buildingBarTabGap          = 6
	buildingBarTabBottomGap    = 14
	buildingBarItemSize        = 64
	buildingBarItemGap         = 12
	buildingBarSpriteInset     = 8
	buildingBarDragIconSize    = buildingBarItemSize / 2
	buildingBarDisabledAlpha   = 0.70
	buildingBarHoverBrightness = 1.65
)

var (
	buildingBarBuildableColor = color.RGBA{R: 92, G: 220, B: 104, A: 255}
	buildingBarBlockedColor   = color.RGBA{R: 224, G: 76, B: 65, A: 255}
)

// BuildingBarAction identifies a construction choice returned to the game package.
type BuildingBarAction int

const (
	BuildingBarHouse BuildingBarAction = iota
	BuildingBarBarracks
	BuildingBarDorm
	BuildingBarWoodcutter
	BuildingBarStoneQuarry
	BuildingBarIronMine
	BuildingBarMarket
	BuildingBarBowTower
	BuildingBarFlameBoltTower
	BuildingBarCatapultTower
)

// BuildingBarCategory identifies one visible group of construction choices.
type BuildingBarCategory int

const (
	BuildingBarCategoryDefenses BuildingBarCategory = iota
	BuildingBarCategoryEconomic
	BuildingBarCategoryHousing
	BuildingBarNoCategory BuildingBarCategory = -1
)

// BuildingBarIcons contains the inhabitant portraits used by item metadata.
type BuildingBarIcons struct {
	Apprentice *ebiten.Image
	Soldier    *ebiten.Image
	Peasant    *ebiten.Image
}

// BuildingBarItem describes one construction choice using presentation-neutral facts.
type BuildingBarItem struct {
	Action                        BuildingBarAction
	Name                          string
	Description                   string
	Sprite                        *ebiten.Image
	Cost                          ResourceAmounts
	Staffing                      PopulationAmounts
	PopulationCost                PopulationAmounts
	PopulationGrant               PopulationAmounts
	ResourceYield                 ResourceAmounts
	RangeTiles                    float64
	Damage                        int
	FireIntervalSeconds           float64
	ProjectileSpeedTilesPerSecond float64
	DamageAllEnemiesInTargetTile  bool
	Buildable                     bool
}

// BuildingBarModel contains all data and host-owned presentation state for the widget.
type BuildingBarModel struct {
	Items            []BuildingBarItem
	Icons            BuildingBarIcons
	SelectedCategory BuildingBarCategory
	HoveredItem      int
	HoveredCategory  BuildingBarCategory
}

type buildingBarTab struct {
	Category BuildingBarCategory
	Label    string
	Bounds   Button[int]
}

type buildingBarLayoutItem struct {
	BuildingBarItem
	Bounds Button[int]
}

// BuildingBarBounds returns the screen-space building-bar rectangle.
func BuildingBarBounds(top, height int) Button[int] {
	return Button[int]{X: 0, Y: top, W: BuildingBarWidth, H: height - top}
}

// BuildingBarContains reports whether a point is inside the building bar.
func BuildingBarContains(top, height, x, y int) bool {
	return BuildingBarBounds(top, height).Contains(x, y)
}

// BuildingBarItemIndexAt returns the visible item index at a point, or -1.
func BuildingBarItemIndexAt(top int, model BuildingBarModel, x, y int) int {
	for index, item := range buildingBarLayoutItems(top, model) {
		if item.Bounds.Contains(x, y) {
			return index
		}
	}
	return -1
}

// BuildingBarItemAt returns the visible construction choice at a point.
func BuildingBarItemAt(top int, model BuildingBarModel, x, y int) (BuildingBarItem, bool) {
	items := buildingBarLayoutItems(top, model)
	for _, item := range items {
		if item.Bounds.Contains(x, y) {
			return item.BuildingBarItem, true
		}
	}
	return BuildingBarItem{}, false
}

// BuildingBarItemBounds returns the visible bounds for one construction action.
func BuildingBarItemBounds(top int, model BuildingBarModel, action BuildingBarAction) (Button[int], bool) {
	for _, item := range buildingBarLayoutItems(top, model) {
		if item.Action == action {
			return item.Bounds, true
		}
	}
	return Button[int]{}, false
}

// BuildingBarActions returns every construction action in stable catalog order.
func BuildingBarActions() []BuildingBarAction {
	return []BuildingBarAction{
		BuildingBarHouse,
		BuildingBarBarracks,
		BuildingBarDorm,
		BuildingBarWoodcutter,
		BuildingBarStoneQuarry,
		BuildingBarIronMine,
		BuildingBarMarket,
		BuildingBarBowTower,
		BuildingBarFlameBoltTower,
		BuildingBarCatapultTower,
	}
}

// BuildingBarCategoryForAction returns the category containing one construction action.
func BuildingBarCategoryForAction(action BuildingBarAction) BuildingBarCategory {
	switch action {
	case BuildingBarHouse, BuildingBarBarracks, BuildingBarDorm:
		return BuildingBarCategoryHousing
	case BuildingBarWoodcutter, BuildingBarStoneQuarry, BuildingBarIronMine, BuildingBarMarket:
		return BuildingBarCategoryEconomic
	case BuildingBarBowTower, BuildingBarFlameBoltTower, BuildingBarCatapultTower:
		return BuildingBarCategoryDefenses
	default:
		return BuildingBarNoCategory
	}
}

// BuildingBarCategoryAt returns the category tab at a point, or no category.
func BuildingBarCategoryAt(top, x, y int) BuildingBarCategory {
	for _, tab := range buildingBarTabs(top) {
		if tab.Bounds.Contains(x, y) {
			return tab.Category
		}
	}
	return BuildingBarNoCategory
}

// BuildingBarCategoryBounds returns the bounds of one category tab.
func BuildingBarCategoryBounds(top int, category BuildingBarCategory) (Button[int], bool) {
	for _, tab := range buildingBarTabs(top) {
		if tab.Category == category {
			return tab.Bounds, true
		}
	}
	return Button[int]{}, false
}

// DrawBuildingBar renders the construction picker at the left edge of the scene.
func DrawBuildingBar(screen *ebiten.Image, regularFace, boldFace *text.GoTextFace, top, height int, model BuildingBarModel) {
	bar := BuildingBarBounds(top, height)
	if bar.H <= 0 {
		return
	}

	vector.FillRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), SelectionPanelBackground, false)
	vector.StrokeLine(screen, float32(bar.X+bar.W-2), float32(bar.Y), float32(bar.X+bar.W-2), float32(bar.Y+bar.H), 3, Bronze, false)

	for _, tab := range buildingBarTabs(top) {
		drawBuildingBarTab(screen, regularFace, model, tab)
	}
	for index, item := range buildingBarLayoutItems(top, model) {
		drawBuildingBarItem(screen, regularFace, boldFace, model, index, item)
	}
}

// DrawBuildingDrag renders a building sprite attached to the cursor.
func DrawBuildingDrag(screen *ebiten.Image, sprite *ebiten.Image, cursorX, cursorY int) {
	if sprite == nil {
		return
	}
	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())
	if spriteWidth <= 0 || spriteHeight <= 0 {
		return
	}

	scale := float64(buildingBarDragIconSize) / spriteWidth
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(float64(cursorX)-spriteWidth*scale/2, float64(cursorY)-spriteHeight*scale/2)
	screen.DrawImage(sprite, options)
}

// buildingBarLayoutItems returns visible items in the selected category's stable order.
func buildingBarLayoutItems(top int, model BuildingBarModel) []buildingBarLayoutItem {
	itemByAction := make(map[BuildingBarAction]BuildingBarItem, len(model.Items))
	for _, item := range model.Items {
		itemByAction[item.Action] = item
	}

	actions := buildingBarActionsForCategory(model.SelectedCategory)
	items := make([]buildingBarLayoutItem, 0, len(actions))
	startY := top + buildingBarPadding + buildingBarTabsHeight() + buildingBarTabBottomGap
	for index, action := range actions {
		item, ok := itemByAction[action]
		if !ok {
			continue
		}
		itemBounds := Button[int]{
			Label:  item.Name,
			X:      buildingBarPadding,
			Y:      startY + index*(buildingBarItemSize+buildingBarItemGap),
			W:      buildingBarItemSize,
			H:      buildingBarItemSize,
			Action: int(item.Action),
		}
		items = append(items, buildingBarLayoutItem{BuildingBarItem: item, Bounds: itemBounds})
	}
	return items
}

// buildingBarTabs returns category tabs in rendered order.
func buildingBarTabs(top int) []buildingBarTab {
	categories := []BuildingBarCategory{
		BuildingBarCategoryDefenses,
		BuildingBarCategoryEconomic,
		BuildingBarCategoryHousing,
	}
	tabs := make([]buildingBarTab, 0, len(categories))
	for index, category := range categories {
		tabs = append(tabs, buildingBarTab{
			Category: category,
			Label:    buildingBarCategoryLabel(category),
			Bounds: Button[int]{
				Label:  buildingBarCategoryLabel(category),
				X:      buildingBarTabGap,
				Y:      top + buildingBarPadding + index*(buildingBarTabHeight+buildingBarTabGap),
				W:      BuildingBarWidth - buildingBarTabGap*2,
				H:      buildingBarTabHeight,
				Action: int(category),
			},
		})
	}
	return tabs
}

// buildingBarTabsHeight returns the vertical space reserved for category tabs.
func buildingBarTabsHeight() int {
	return 3*buildingBarTabHeight + 2*buildingBarTabGap
}

// buildingBarCategoryLabel returns the visible label for one category.
func buildingBarCategoryLabel(category BuildingBarCategory) string {
	switch category {
	case BuildingBarCategoryDefenses:
		return "Defenses"
	case BuildingBarCategoryEconomic:
		return "Economic"
	case BuildingBarCategoryHousing:
		return "Housing"
	default:
		return ""
	}
}

// buildingBarActionsForCategory returns stable actions in visible order.
func buildingBarActionsForCategory(category BuildingBarCategory) []BuildingBarAction {
	switch category {
	case BuildingBarCategoryHousing:
		return []BuildingBarAction{BuildingBarHouse, BuildingBarBarracks, BuildingBarDorm}
	case BuildingBarCategoryEconomic:
		return []BuildingBarAction{BuildingBarWoodcutter, BuildingBarStoneQuarry, BuildingBarIronMine, BuildingBarMarket}
	case BuildingBarCategoryDefenses:
		return []BuildingBarAction{BuildingBarBowTower, BuildingBarFlameBoltTower, BuildingBarCatapultTower}
	default:
		return nil
	}
}

// drawBuildingBarTab renders one category tab.
func drawBuildingBarTab(screen *ebiten.Image, face *text.GoTextFace, model BuildingBarModel, tab buildingBarTab) {
	selected := model.SelectedCategory == tab.Category
	hovered := model.HoveredCategory == tab.Category
	fill := DarkCharcoalGreen
	if selected {
		fill = Bronze
	}
	vector.FillRect(screen, float32(tab.Bounds.X), float32(tab.Bounds.Y), float32(tab.Bounds.W), float32(tab.Bounds.H), fill, false)
	vector.StrokeRect(screen, float32(tab.Bounds.X), float32(tab.Bounds.Y), float32(tab.Bounds.W), float32(tab.Bounds.H), 1, Bronze, false)

	textColor := Parchment
	if hovered && !selected {
		textColor = LightBronze
	}
	width, height := text.Measure(tab.Label, face, face.Size)
	x := float64(tab.Bounds.X) + (float64(tab.Bounds.W)-width)/2
	y := float64(tab.Bounds.Y) + (float64(tab.Bounds.H)-height)/2 - 1
	DrawText(screen, tab.Label, face, x, y, textColor)
}

// drawBuildingBarItem renders one building icon slot and its metadata.
func drawBuildingBarItem(screen *ebiten.Image, regularFace, boldFace *text.GoTextFace, model BuildingBarModel, index int, item buildingBarLayoutItem) {
	hovered := model.HoveredItem == index && item.Buildable
	bounds := item.Bounds
	vector.FillRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), DarkCharcoalGreen, false)
	vector.StrokeRect(screen, float32(bounds.X), float32(bounds.Y), float32(bounds.W), float32(bounds.H), 2, buildingBarOutlineColor(item.Buildable), false)

	if item.Sprite != nil {
		spriteWidth := float64(item.Sprite.Bounds().Dx())
		spriteHeight := float64(item.Sprite.Bounds().Dy())
		if spriteWidth > 0 && spriteHeight > 0 {
			targetSize := float64(bounds.W - buildingBarSpriteInset*2)
			scale := targetSize / spriteWidth
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(scale, scale)
			options.GeoM.Translate(float64(bounds.X)+(float64(bounds.W)-spriteWidth*scale)/2, float64(bounds.Y)+(float64(bounds.H)-spriteHeight*scale)/2)
			options.ColorScale.Scale(1, 1, 1, buildingBarIconAlpha(item.Buildable))
			if hovered {
				options.ColorScale.Scale(buildingBarHoverBrightness, buildingBarHoverBrightness, buildingBarHoverBrightness, 1)
			}
			screen.DrawImage(item.Sprite, options)
		}
	}

	drawBuildingBarCost(screen, regularFace, boldFace, item, hovered)
	drawBuildingBarPopulationMetadata(screen, regularFace, model.Icons, item)
}

// buildingBarIconAlpha returns opacity for a construction choice.
func buildingBarIconAlpha(buildable bool) float32 {
	if buildable {
		return 1
	}
	return buildingBarDisabledAlpha
}

// buildingBarOutlineColor returns the availability outline color.
func buildingBarOutlineColor(buildable bool) color.Color {
	if buildable {
		return buildingBarBuildableColor
	}
	return buildingBarBlockedColor
}
