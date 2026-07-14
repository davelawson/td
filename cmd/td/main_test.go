package main

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"td/internal/game"
	"td/internal/menu"

	"github.com/hajimehoshi/ebiten/v2"
)

type screenshotApp struct {
	*app
	targets []screenshotTarget
	index   int
}

type screenshotTarget struct {
	screen           menu.Screen
	wizardName       string
	paused           bool
	ingameMenu       bool
	activeRaid       bool
	selectedSanctum  bool
	selectedRaider   bool
	hoverBuilding    bool
	hoverEconomic    bool
	hoverBarracks    bool
	hoverDorm        bool
	placedHouse      bool
	placedBarracks   bool
	placedDorm       bool
	placedWoodcutter bool
	exploredBiomes   bool
	path             string
}

// TestStartGameSwitchesToGameMode verifies app-level game startup.
func TestStartGameSwitchesToGameMode(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	if err := app.startGame("Merlin"); err != nil {
		t.Fatal(err)
	}
	if app.mode != appModeGame {
		t.Fatalf("mode = %v, want %v", app.mode, appModeGame)
	}
	if app.gameState == nil {
		t.Fatal("expected game state after starting")
	}
	if app.gameState.WizardName() != "Merlin" {
		t.Fatalf("wizard name = %q, want %q", app.gameState.WizardName(), "Merlin")
	}
}

// TestSurrenderReturnsToMainMenu verifies app-level surrender routing.
func TestSurrenderReturnsToMainMenu(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	if err := app.startGame("Merlin"); err != nil {
		t.Fatal(err)
	}

	app.gameState.Update(game.Input{ToggleMenu: true})
	if action := app.gameState.Update(gameClickInput(app.gameState, 1)); action != game.ActionSurrender {
		t.Fatalf("Update(surrender click) = %v, want %v", action, game.ActionSurrender)
	}
	app.returnToMainMenu()
	if app.mode != appModeMenu {
		t.Fatalf("mode = %v, want %v", app.mode, appModeMenu)
	}
	if app.gameState != nil {
		t.Fatal("expected surrender to clear game state")
	}
	if app.mainMenu.Screen() != menu.ScreenMain {
		t.Fatalf("screen = %v, want %v", app.mainMenu.Screen(), menu.ScreenMain)
	}
}

// TestCaptureMainMenuScreenshot writes visual evidence when explicitly enabled.
func TestCaptureMainMenuScreenshot(t *testing.T) {
	if os.Getenv("TD_CAPTURE_SCREENSHOT") == "" {
		t.Skip("set TD_CAPTURE_SCREENSHOT to capture visual evidence")
	}

	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	basePath := filepath.Join("..", "..", "plans", "51-peaceful-phase-split", "screenshots")
	capture := &screenshotApp{
		app: app,
		targets: []screenshotTarget{
			{screen: menu.ScreenMain, path: filepath.Join(basePath, "main-menu.png")},
			{screen: menu.ScreenNewGame, path: filepath.Join(basePath, "new-game-configuration.png")},
			{wizardName: "Merlin", path: filepath.Join(basePath, "running-game.png")},
			{wizardName: "Merlin", exploredBiomes: true, path: filepath.Join(basePath, "explored-biomes.png")},
			{wizardName: "Merlin", hoverBuilding: true, path: filepath.Join(basePath, "house-icon.png")},
			{wizardName: "Merlin", hoverEconomic: true, path: filepath.Join(basePath, "woodcutter-tooltip.png")},
			{wizardName: "Merlin", placedHouse: true, path: filepath.Join(basePath, "placed-house.png")},
			{wizardName: "Merlin", placedWoodcutter: true, path: filepath.Join(basePath, "placed-woodcutter.png")},
			{wizardName: "Merlin", hoverBarracks: true, placedHouse: true, path: filepath.Join(basePath, "barracks-icon.png")},
			{wizardName: "Merlin", placedBarracks: true, path: filepath.Join(basePath, "placed-barracks.png")},
			{wizardName: "Merlin", hoverDorm: true, placedHouse: true, path: filepath.Join(basePath, "dorm-icon.png")},
			{wizardName: "Merlin", placedDorm: true, path: filepath.Join(basePath, "placed-dorm.png")},
			{wizardName: "Merlin", selectedSanctum: true, path: filepath.Join(basePath, "selected-sanctum.png")},
			{wizardName: "Merlin", activeRaid: true, path: filepath.Join(basePath, "active-raid.png")},
			{wizardName: "Merlin", activeRaid: true, selectedRaider: true, path: filepath.Join(basePath, "selected-raider.png")},
			{wizardName: "Merlin", paused: true, path: filepath.Join(basePath, "paused-game.png")},
			{wizardName: "Merlin", ingameMenu: true, path: filepath.Join(basePath, "ingame-menu.png")},
		},
	}

	ebiten.SetWindowTitle("td")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	if err := ebiten.RunGame(capture); err != nil && !errors.Is(err, ebiten.Termination) {
		t.Fatal(err)
	}
	if capture.index != len(capture.targets) {
		t.Fatalf("captured %d screenshots, want %d", capture.index, len(capture.targets))
	}
}

// Update sets the next screen to capture or terminates after all captures.
func (a *screenshotApp) Update() error {
	if a.index >= len(a.targets) {
		return ebiten.Termination
	}
	target := a.targets[a.index]
	if target.wizardName != "" {
		if err := a.app.startGame(target.wizardName); err != nil {
			return err
		}
		a.app.gameState.Update(game.Input{})
		if target.activeRaid {
			a.app.gameState.Update(game.Input{
				CursorX: 397,
				CursorY: defaultWindowHeight - 68,
				Clicked: true,
			})
			for i := 0; i < 45; i++ {
				a.app.gameState.Update(game.Input{})
			}
		}
		if target.exploredBiomes {
			centerX := defaultWindowWidth / 2
			centerY := topOfGameScene() + defaultGameSceneHeight()/2
			for _, point := range [][2]int{
				{centerX, centerY - 405},
				{centerX + 405, centerY},
				{centerX, centerY + 405},
				{centerX - 405, centerY},
			} {
				a.app.gameState.Update(game.Input{
					CursorX: point[0],
					CursorY: point[1],
					Clicked: true,
				})
			}
			a.app.gameState.Update(game.Input{WheelY: -10})
		}
		if target.hoverBuilding {
			a.app.gameState.Update(game.Input{
				CursorX: 48,
				CursorY: topOfGameScene() + 158,
			})
		}
		if target.hoverEconomic {
			a.app.gameState.Update(game.Input{
				CursorX: 48,
				CursorY: topOfGameScene() + 64,
				Clicked: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX: 48,
				CursorY: topOfGameScene() + 158,
			})
		}
		if target.placedHouse {
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 158,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1068,
				CursorY:  475,
				Released: true,
			})
		}
		if target.hoverBarracks {
			a.app.gameState.Update(game.Input{
				CursorX: 48,
				CursorY: topOfGameScene() + 234,
			})
		}
		if target.hoverDorm {
			a.app.gameState.Update(game.Input{
				CursorX: 48,
				CursorY: topOfGameScene() + 310,
			})
		}
		if target.placedBarracks {
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 158,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1068,
				CursorY:  475,
				Released: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 234,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1130,
				CursorY:  475,
				Released: true,
			})
		}
		if target.placedDorm {
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 158,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1068,
				CursorY:  475,
				Released: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 310,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1130,
				CursorY:  475,
				Released: true,
			})
		}
		if target.placedWoodcutter {
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 158,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1068,
				CursorY:  475,
				Released: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 64,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{})
			a.app.gameState.Update(game.Input{
				CursorX:   48,
				CursorY:   topOfGameScene() + 158,
				Clicked:   true,
				MouseDown: true,
			})
			a.app.gameState.Update(game.Input{
				CursorX:  1130,
				CursorY:  475,
				Released: true,
			})
		}
		if target.selectedRaider {
			a.app.gameState.Update(game.Input{
				CursorX: defaultWindowWidth / 2,
				CursorY: topOfGameScene() + 160,
				Clicked: true,
			})
		}
		if target.selectedSanctum {
			a.app.gameState.Update(game.Input{
				CursorX: defaultWindowWidth / 2,
				CursorY: topOfGameScene() + (defaultWindowHeight-topOfGameScene())/2,
				Clicked: true,
			})
		}
		if target.paused {
			a.app.gameState.Update(game.Input{TogglePause: true})
		}
		if target.ingameMenu {
			a.app.gameState.Update(game.Input{ToggleMenu: true})
		}
		return nil
	}

	a.app.mode = appModeMenu
	a.app.gameState = nil
	a.app.mainMenu.SetScreenForTest(target.screen)
	return nil
}

// Draw renders the current target screen and writes the frame to disk.
func (a *screenshotApp) Draw(screen *ebiten.Image) {
	if a.index >= len(a.targets) {
		return
	}

	a.app.Draw(screen)
	target := a.targets[a.index]

	if err := os.MkdirAll(filepath.Dir(target.path), 0o755); err != nil {
		panic(err)
	}
	file, err := os.Create(target.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	frame := image.NewRGBA(image.Rect(0, 0, defaultWindowWidth, defaultWindowHeight))
	screen.ReadPixels(frame.Pix)
	if err := png.Encode(file, frame); err != nil {
		panic(err)
	}
	a.index++
}

// gameClickInput returns a click at the center of a known in-game menu button.
func gameClickInput(state *game.State, buttonIndex int) game.Input {
	switch buttonIndex {
	case 0:
		return game.Input{CursorX: stateCenterX(state), CursorY: stateCenterY(state), Clicked: true}
	default:
		return game.Input{CursorX: stateCenterX(state), CursorY: stateCenterY(state) + 68, Clicked: true}
	}
}

// stateCenterX returns the test window horizontal center.
func stateCenterX(_ *game.State) int {
	return defaultWindowWidth / 2
}

// stateCenterY returns the top in-game button vertical center.
func stateCenterY(_ *game.State) int {
	return defaultWindowHeight/2 + 8
}

// topOfGameScene returns the default screenshot top edge below the game HUD.
func topOfGameScene() int {
	return 86
}

// defaultGameSceneHeight returns the default screenshot map viewport height.
func defaultGameSceneHeight() int {
	return defaultWindowHeight - topOfGameScene()
}
