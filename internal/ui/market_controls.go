package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	marketControlWidth       = 220
	marketControlHeight      = 42
	marketControlGap         = 8
	marketControlAnchorGap   = 10
	marketControlPadding     = 10
	marketControlIconSize    = 26
	marketControlIconTextGap = 9
)

var marketControlDisabledFill = color.RGBA{R: 47, G: 48, B: 43, A: 230}

// MarketTradeAction identifies one material purchase at a selected Market.
type MarketTradeAction int

const (
	MarketTradeBuyWood MarketTradeAction = iota
	MarketTradeBuyStone
	MarketTradeBuyIron
	MarketTradeNoAction MarketTradeAction = -1
)

// MarketTradeItem describes one Market purchase using presentation-neutral facts.
type MarketTradeItem struct {
	Action   MarketTradeAction
	Label    string
	Icon     *ebiten.Image
	GoldCost int
	Enabled  bool
}

// MarketControlsModel contains ordered purchases and host-owned hover state.
type MarketControlsModel struct {
	Items   []MarketTradeItem
	Hovered MarketTradeAction
}

// MarketControlButtons lays out fixed-size controls beside a projected Market.
func MarketControlButtons(anchor, area Button[int], model MarketControlsModel) []Button[MarketTradeAction] {
	if len(model.Items) == 0 || area.W <= 0 || area.H <= 0 {
		return nil
	}
	totalHeight := len(model.Items)*marketControlHeight + (len(model.Items)-1)*marketControlGap
	x := anchor.X + anchor.W + marketControlAnchorGap
	if x+marketControlWidth > area.X+area.W {
		x = anchor.X - marketControlAnchorGap - marketControlWidth
	}
	x = clampMarketControlCoordinate(x, area.X, area.X+area.W-marketControlWidth)
	y := clampMarketControlCoordinate(anchor.Y, area.Y, area.Y+area.H-totalHeight)

	buttons := make([]Button[MarketTradeAction], 0, len(model.Items))
	for index, item := range model.Items {
		buttons = append(buttons, Button[MarketTradeAction]{
			Label:    item.Label,
			X:        x,
			Y:        y + index*(marketControlHeight+marketControlGap),
			W:        marketControlWidth,
			H:        marketControlHeight,
			Action:   item.Action,
			Disabled: !item.Enabled,
		})
	}
	return buttons
}

// MarketTradeAt returns the Market item under a screen point, including disabled items.
func MarketTradeAt(anchor, area Button[int], model MarketControlsModel, x, y int) (MarketTradeItem, bool) {
	buttons := MarketControlButtons(anchor, area, model)
	for index, button := range buttons {
		if button.Contains(x, y) {
			return model.Items[index], true
		}
	}
	return MarketTradeItem{}, false
}

// MarketControlsContains reports whether a point is inside any visible Market button.
func MarketControlsContains(anchor, area Button[int], model MarketControlsModel, x, y int) bool {
	_, ok := MarketTradeAt(anchor, area, model, x, y)
	return ok
}

// DrawMarketControls renders material purchase buttons beside the selected Market.
func DrawMarketControls(screen *ebiten.Image, face *text.GoTextFace, anchor, area Button[int], model MarketControlsModel) {
	buttons := MarketControlButtons(anchor, area, model)
	for index, button := range buttons {
		item := model.Items[index]
		fill := DarkCharcoalGreen
		edge := Bronze
		labelColor := Parchment
		if button.Disabled {
			fill = marketControlDisabledFill
			edge = TransparentBronze
			labelColor = MutedParchment
		} else if model.Hovered == item.Action {
			fill = MossGreen
			edge = ResourceGold
		}

		vector.FillRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), fill, false)
		vector.StrokeRect(screen, float32(button.X), float32(button.Y), float32(button.W), float32(button.H), 3, edge, false)
		drawMarketControlIcon(screen, item.Icon, button)
		DrawText(
			screen,
			button.Label,
			face,
			float64(button.X+marketControlPadding+marketControlIconSize+marketControlIconTextGap),
			float64(button.Y+10),
			labelColor,
		)
	}
}

// drawMarketControlIcon renders one material icon inside a Market button.
func drawMarketControlIcon(screen *ebiten.Image, icon *ebiten.Image, button Button[MarketTradeAction]) {
	if icon == nil || icon.Bounds().Dx() <= 0 || icon.Bounds().Dy() <= 0 {
		return
	}
	scale := float64(marketControlIconSize) / float64(icon.Bounds().Dx())
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(
		float64(button.X+marketControlPadding),
		float64(button.Y+(button.H-marketControlIconSize)/2),
	)
	screen.DrawImage(icon, options)
}

// clampMarketControlCoordinate clamps a coordinate and tolerates undersized areas.
func clampMarketControlCoordinate(value, minimum, maximum int) int {
	if maximum < minimum {
		return minimum
	}
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
