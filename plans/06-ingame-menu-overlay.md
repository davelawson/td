# Add In-Game Menu Overlay

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/06-ingame-menu-overlay.md` because `plans/00-initial-ebitengine-menu.md` through `plans/05-main-game-update-loop.md` already exist.

## Purpose / Big Picture

The app can currently enter a placeholder game state and pause it with SPACE, but there is no in-game route back to the main menu. After this change, a player can press ESC from the game view to open an overlay menu, resume the current game, or surrender back to the main menu. The overlay keeps the game view visible underneath a roughly 50% dark layer so it reads as an in-game pause menu rather than a separate screen.

## Progress

- [x] (2026-05-08 02:25Z) Created this ExecPlan for the accepted in-game menu overlay plan.
- [x] (2026-05-08 02:28Z) Added in-game menu state, pause-on-open behavior, overlay rendering, and tests in `internal/game`.
- [x] (2026-05-08 02:30Z) Added app-level surrender routing and menu reset behavior.
- [x] (2026-05-08 02:34Z) Updated durable control documents for the new user-visible workflow and architecture.
- [x] (2026-05-08 02:36Z) Captured screenshot evidence for the in-game overlay.
- [x] (2026-05-08 02:38Z) Ran final validation and recorded outcomes.
- [x] (2026-05-08 02:39Z) Checked hand-written code-file line counts and found no files over the 600-line preference.

## Surprises & Discoveries

- Observation: The `Resume` button should not return an app-level action, so hover state cannot be represented only by the exported `game.Action` value.
  Evidence: `internal/game/ingamemenu.go` tracks `ingameMenuHover` as a button index while keeping exported actions limited to `ActionNone` and `ActionSurrender`.

- Observation: Screenshot capture produced a clear overlay result with the game scene visible beneath a dark layer and the centered panel readable.
  Evidence: `plans/06-ingame-menu-overlay/screenshots/ingame-menu.png` is a 1920x1080 PNG and was reviewed after capture.

## Decision Log

- Decision: Put in-game menu-specific code in `internal/game/ingamemenu.go`.
  Rationale: The overlay is tightly coupled to game rendering, pause state, and game-level actions, but separating it from `game.go` keeps the top-level game state file focused.
  Date/Author: 2026-05-08 / Codex

- Decision: Opening the in-game menu records the previous pause state, then sets `paused = true`.
  Rationale: ESC should always pause active gameplay while the overlay is visible, and Resume should preserve whether the player had already paused with SPACE before opening the overlay.
  Date/Author: 2026-05-08 / Codex

- Decision: Surrender discards the active game state and returns to the top-level main menu without confirmation.
  Rationale: There is no save/load system or meaningful game progress yet, so confirmation UI would add scope before it protects real user data.
  Date/Author: 2026-05-08 / Codex

## Outcomes & Retrospective

Implementation completed the ESC in-game menu overlay. Pressing ESC from the game view now opens a centered overlay, pauses the game, and draws the existing game scene beneath a roughly 50% dark layer. ESC while the overlay is open and clicking `Resume` both close the overlay and restore the previous pause state. Clicking `Surrender` returns a game-level action that `cmd/td` handles by clearing the active game state and returning to the top-level main menu.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.805s

    file plans/06-ingame-menu-overlay/screenshots/*.png
    plans/06-ingame-menu-overlay/screenshots/ingame-menu.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/06-ingame-menu-overlay/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/06-ingame-menu-overlay/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/06-ingame-menu-overlay/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/06-ingame-menu-overlay/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

Final hand-written Go file line-count review:

    19 internal/ui/colors.go
    151 cmd/td/main.go
    166 internal/game/ingamemenu.go
    181 cmd/td/main_test.go
    188 internal/game/game.go
    201 internal/game/game_test.go
    282 internal/menu/menu_test.go
    489 internal/menu/menu.go
    1677 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local Go/Ebitengine desktop prototype. Ebitengine calls `Update`, `Draw`, and `Layout` on the app in `cmd/td/main.go`. That executable owns process startup, input polling, app-mode routing between the main menu and game, and returning Ebitengine's termination signal for Quit.

The main menu lives in `internal/menu/menu.go`. It owns main menu screen state, New Game configuration, Wizard name input, button hit testing, and rendering. It should not own the in-game overlay because the overlay is drawn over the game scene and controls game pause behavior. A production `ResetToMain` method returns the main menu to its top-level screen when surrender leaves the game.

The current game state lives in `internal/game/game.go`. It owns the Wizard name, drawable size, logical update counter, SPACE pause state, and placeholder game rendering. This plan adds `internal/game/ingamemenu.go` as the home for in-game menu-specific code: button geometry, open/close behavior, hover and click handling, ESC resume handling, and overlay drawing. A logical update is one project-owned game-state step; it must not advance while paused or while the in-game menu is open.

The root control documents constrain this work. `PRODUCT.md` and `README.md` must describe ESC, Resume, and Surrender because they change current user-visible workflows. `ARCHITECTURE.md` must describe that the game package owns the in-game overlay while `cmd/td` owns app-level surrender routing. `ROADMAP.md` only needs a small status update so the current phase remains accurate. `DESIGN.md` already says UI should be readable, stable, restrained, and medieval wizardry themed; the overlay should reuse the existing palette rather than introduce a new design direction. `CODESTYLE.md` requires `gofmt`, doc comments on Go functions and methods, tests for pure behavior, and a final hand-written code-file line-count review.

## Plan of Work

First, update `internal/game`. Extend `Input` with ESC, cursor position, and click fields. Make `State.Update` return a game-level action so app routing can respond to Surrender. Add `ActionNone` and `ActionSurrender`. Keep the normal update counter and SPACE pause behavior in `game.go`, but move overlay-specific constants, button structs, layout, input handling, and drawing to `internal/game/ingamemenu.go`.

The in-game menu opens when ESC is received while the game overlay is closed. Opening stores the current `paused` value in `pausedBeforeMenu`, then sets `paused = true` and `ingameMenuOpen = true`. While the overlay is open, SPACE is ignored and update counts do not advance. ESC while the overlay is open closes it and restores `paused` from `pausedBeforeMenu`. Clicking `Resume` performs the same close-and-restore behavior. Clicking `Surrender` returns `ActionSurrender`.

Next, update `cmd/td`. In game mode, poll SPACE, ESC, cursor position, and left-click state, pass them to `game.State.Update`, and when it returns `game.ActionSurrender`, switch to menu mode, clear `gameState`, and reset the main menu to `ScreenMain`. Keep Ebitengine startup, layout policy, and main-menu input routing unchanged.

Finally, update tests, screenshot capture, and documentation. Add game tests for pause-on-open, resume restoring running state, resume restoring paused state, ESC-as-resume, update blocking while open, Surrender action, and resize geometry. Add an app-level test for surrender returning to main menu. Extend gated screenshot capture to write the in-game menu screenshot under this plan's screenshot directory. Update `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, `ROADMAP.md`, and the stale runtime-status line in `CODESTYLE.md`.

## Concrete Steps

From the repository root, inspect the current state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Edit `internal/game/game.go` and create `internal/game/ingamemenu.go`. Run:

    gofmt -w internal/game/game.go internal/game/ingamemenu.go internal/game/game_test.go
    go test ./internal/game

Edit `internal/menu/menu.go` to add `ResetToMain`, and edit `cmd/td/main.go` for ESC, pointer input, and surrender routing. Run:

    gofmt -w internal/menu/menu.go cmd/td/main.go cmd/td/main_test.go
    go test ./cmd/td ./internal/menu

Update screenshot capture in `cmd/td/main_test.go` so it writes:

    plans/06-ingame-menu-overlay/screenshots/main-menu.png
    plans/06-ingame-menu-overlay/screenshots/new-game-configuration.png
    plans/06-ingame-menu-overlay/screenshots/running-game.png
    plans/06-ingame-menu-overlay/screenshots/paused-game.png
    plans/06-ingame-menu-overlay/screenshots/ingame-menu.png

Update `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, `ROADMAP.md`, and `CODESTYLE.md`. Then run:

    gofmt -w cmd/td/*.go internal/game/*.go internal/menu/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, screenshot capture writes all five PNG files under `plans/06-ingame-menu-overlay/screenshots/`, and `go run ./cmd/td` opens without startup errors. A human should observe that starting a game and pressing ESC opens a centered in-game overlay, darkens the rest of the game view by about 50%, shows `Resume` and `Surrender`, and stops the update counter while visible.

Resume acceptance has two cases. From a running game, ESC opens the overlay and ESC again or `Resume` returns to the running game. From a SPACE-paused game, ESC opens the overlay and ESC again or `Resume` returns to the game still paused. Surrender acceptance is that clicking `Surrender` leaves the game state and returns to the top-level main menu.

Documentation is accepted when `PRODUCT.md` and `README.md` describe the current workflow, `ARCHITECTURE.md` records the ownership split between `cmd/td` and `internal/game`, `ROADMAP.md` remains accurate about current foundation work, and `CODESTYLE.md` no longer says runtime code does not exist.

## Idempotence and Recovery

The changes are local to game state, menu reset, app routing, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when graphics are available.

If Surrender routing fails during implementation, keep the game-level `ActionSurrender` return value in `internal/game` and fix only the app-mode transition in `cmd/td`; do not move app-mode ownership into the game package.

## Artifacts and Notes

Important artifacts:

    plans/06-ingame-menu-overlay.md
    plans/06-ingame-menu-overlay/screenshots/ingame-menu.png
    cmd/td/main.go
    cmd/td/main_test.go
    internal/game/game.go
    internal/game/ingamemenu.go
    internal/game/game_test.go
    internal/menu/menu.go
    PRODUCT.md
    README.md
    ARCHITECTURE.md
    ROADMAP.md
    CODESTYLE.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

At the end of this plan, `internal/game` exposes:

    type Input struct {
        TogglePause bool
        ToggleMenu  bool
        CursorX     int
        CursorY     int
        Clicked     bool
    }

    type Action int

    const (
        ActionNone Action = iota
        ActionSurrender
    )

    func New(wizardName string, width, height int) (*State, error)
    func (s *State) Resize(width, height int)
    func (s *State) Update(input Input) Action
    func (s *State) Draw(screen *ebiten.Image)
    func (s *State) Updates() int
    func (s *State) Paused() bool
    func (s *State) IngameMenuOpen() bool
    func (s *State) WizardName() string

The `internal/menu` package exposes:

    func (m *Menu) ResetToMain()

`ResetToMain` is production code for returning from Surrender. The existing `SetScreenForTest` method remains test and screenshot setup support only.
