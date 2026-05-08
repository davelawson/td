# Add Main Game Update Loop

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/05-main-game-update-loop.md` because `plans/00-initial-ebitengine-menu.md`, `plans/01-expanded-main-menu.md`, `plans/02-menu-package-refactor.md`, `plans/03-new-game-configuration.md`, and `plans/04-resolution-and-pixel-text-scaling.md` already exist.

## Purpose / Big Picture

The app currently stops at menu navigation: the New Game screen accepts a Wizard name, but Start is disabled and no game state exists. After this change, a contributor can enter a Wizard name, start the first game state, see a top-right logical update counter, and press SPACE to pause or unpause updates. This proves the app can transition out of menus and run deterministic game logic before adding exploration, maps, combat, resources, or saves.

## Progress

- [x] (2026-05-08 01:45Z) Created this ExecPlan for the accepted main game update loop plan.
- [x] (2026-05-08 01:48Z) Added `internal/game` state, rendering, pause input, and tests.
- [x] (2026-05-08 01:50Z) Updated menu Start behavior and Wizard name length to 32 runes.
- [x] (2026-05-08 01:52Z) Updated `cmd/td` app-mode routing and tests.
- [x] (2026-05-08 01:55Z) Updated screenshots and durable control documents.
- [x] (2026-05-08 01:56Z) Ran formatting, tests, screenshot capture, launch validation, whitespace checks, git status, and the required hand-written Go file line-count review.

## Surprises & Discoveries

- Observation: A 32-rune Wizard name can fit in the New Game field if the field is widened and the field text uses a smaller fixed pixel font than menu buttons.
  Evidence: `TestWizardNameMaxLengthFitsField` measures 32 wide runes against the final name field width and passes.

- Observation: Screenshot capture can switch between menu mode and game mode in one gated test run.
  Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` wrote main menu, New Game, running game, and paused game PNG files under `plans/05-main-game-update-loop/screenshots/`.

## Decision Log

- Decision: Treat 32 characters as 32 printable Go runes.
  Rationale: The current Wizard name input already uses rune-aware append and Backspace behavior, so this preserves existing Unicode semantics.
  Date/Author: 2026-05-08 / Codex

- Decision: Use an explicit app mode in `cmd/td` to switch between menu and game behavior.
  Rationale: Ebitengine input polling and screen routing belong at the executable boundary, while menu and game packages should own their own state interpretation.
  Date/Author: 2026-05-08 / Codex

- Decision: Make SPACE toggle pause without incrementing the logical update counter on the toggle frame.
  Rationale: Pause control is app/game state input, not a logical gameplay step. This makes the counter easy to reason about in tests.
  Date/Author: 2026-05-08 / Codex

## Outcomes & Retrospective

Implementation completed the first app-level game transition and logical update loop. The New Game screen now accepts Wizard names up to 32 printable runes. Start remains disabled while the name is empty and becomes active after name entry. Clicking Start switches the app from menu mode to game mode and constructs `internal/game` state with the Wizard name. The game screen renders placeholder field geometry, the Wizard name, a top-right `Updates: N` counter, and a `PAUSED` label while paused. SPACE toggles pause without incrementing the logical update counter on the toggle frame, and paused frames do not increment the counter.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.788s

    file plans/05-main-game-update-loop/screenshots/*.png
    plans/05-main-game-update-loop/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/05-main-game-update-loop/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/05-main-game-update-loop/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/05-main-game-update-loop/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

    git status --short
    M ARCHITECTURE.md
    M PRODUCT.md
    M README.md
    M ROADMAP.md
    M cmd/td/main.go
    M cmd/td/main_test.go
    M internal/menu/menu.go
    M internal/menu/menu_test.go
    ?? internal/game/
    ?? plans/05-main-game-update-loop.md
    ?? plans/05-main-game-update-loop/

Final hand-written Go file line-count review:

    87 internal/game/game_test.go
    130 cmd/td/main_test.go
    136 cmd/td/main.go
    149 internal/game/game.go
    282 internal/menu/menu_test.go
    475 internal/menu/menu.go
    1259 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local Go/Ebitengine tower-defense prototype. Ebitengine owns the desktop window, input callbacks, drawing surface, and game loop. The executable in `cmd/td/main.go` sets the window title and size, creates app state, polls Ebitengine input, delegates update and draw behavior, and implements `Layout`, which returns the current drawable size.

The current `internal/menu` package owns menu state, rendering, hit testing, disabled-target handling, screen routing, and Wizard name input. Its New Game screen currently has a disabled Start button and caps Wizard names at 16 runes. This plan changes that package so Start becomes active once the Wizard name is non-empty and Wizard names cap at 32 runes.

This plan creates `internal/game`, a new package that owns the first game state and logical game update rules. A logical update is a project-owned game-state step that should only happen while unpaused. Ebitengine will continue calling `Update`, but the game state's logical update counter will not advance while paused.

The root control documents constrain the work. `PRODUCT.md` must change because starting a game, pausing, and the 32-character Wizard name limit become current user-visible behavior. `ARCHITECTURE.md` must change because `internal/game` becomes real and owns game state. `README.md` must change to mention starting the game and SPACE pause. `ROADMAP.md` should be adjusted only enough to reflect this foundation progress. `DESIGN.md` changes only if implementation introduces a durable visual rule. `CODESTYLE.md` requires `gofmt`, doc comments for every Go function and method, tests for pure behavior, and a final line-count review for hand-written code files.

## Plan of Work

First, add `internal/game`. Define an `Input` type with a `TogglePause bool` field. Define a `State` struct that stores the Wizard name, drawable width and height, logical update count, pause flag, and font faces needed to render the counter. Add `New(wizardName string, width, height int) (*State, error)`, `Resize(width, height int)`, `Update(input Input)`, `Draw(screen *ebiten.Image)`, `Updates() int`, `Paused() bool`, and `WizardName() string`. `Update` should toggle pause and return without incrementing when `TogglePause` is true. If not toggling and not paused, it should increment the counter once. `Draw` should render a simple placeholder game screen, Wizard name copy, a top-right `Updates: N` counter, and a small `PAUSED` label near the counter only while paused.

Next, update `internal/menu`. Increase `wizardNameMaxRunes` from 16 to 32. Add `ActionStart`. Recompute New Game buttons so Start is disabled only when the Wizard name is empty, and active with `ActionStart` once at least one printable rune has been entered. Widen the New Game panel and name field enough for 32 runes at the current UI scale while keeping the field's dimensions stable during typing.

Then, update `cmd/td`. Add an app mode enum for menu mode and game mode. Store the current drawable width and height on the app. In menu mode, continue polling pointer, typed runes, and Backspace for `internal/menu`; translate `ActionQuit` to `ebiten.Termination`; translate `ActionStart` into constructing a new `game.State` with the current Wizard name and switching to game mode. In game mode, poll SPACE with `inpututil.IsKeyJustPressed(ebiten.KeySpace)`, pass it to `game.State.Update`, and draw the game instead of the menu. `Layout` should keep returning the current drawable size and resize whichever state exists.

Finally, update tests, screenshot capture, and docs. Add focused `internal/game` tests for counter and pause behavior. Update menu tests for 32-rune caps and Start gating. Add app tests for transition to game mode where practical without opening a graphics window. Update gated screenshot capture to save main menu, New Game configuration, running game, and paused game screenshots under this plan's screenshot directory. Update `PRODUCT.md`, `ARCHITECTURE.md`, `README.md`, and `ROADMAP.md`.

## Concrete Steps

From the repository root, inspect the current state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Add `internal/game/game.go` and `internal/game/game_test.go`. Run:

    gofmt -w internal/game/game.go internal/game/game_test.go
    go test ./internal/game

Edit `internal/menu/menu.go` and `internal/menu/menu_test.go` for 32-rune names and active Start behavior. Run:

    gofmt -w internal/menu/menu.go internal/menu/menu_test.go
    go test ./internal/menu

Edit `cmd/td/main.go` and `cmd/td/main_test.go` for app-mode routing and screenshots. Run:

    gofmt -w cmd/td/main.go cmd/td/main_test.go
    go test ./...

Update `PRODUCT.md`, `ARCHITECTURE.md`, `README.md`, and `ROADMAP.md`. Then capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected screenshot artifacts:

    plans/05-main-game-update-loop/screenshots/main-menu.png
    plans/05-main-game-update-loop/screenshots/new-game-configuration.png
    plans/05-main-game-update-loop/screenshots/running-game.png
    plans/05-main-game-update-loop/screenshots/paused-game.png

Run final validation:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/menu/menu.go internal/menu/menu_test.go internal/game/game.go internal/game/game_test.go
    go test ./...
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, screenshot capture writes all four PNG files under this plan directory, and `go run ./cmd/td` opens without startup errors. A human should observe that entering a Wizard name enables Start, clicking Start closes the menu and begins the game, the top-right counter increments while running, pressing SPACE shows `PAUSED`, and the counter stops advancing until SPACE is pressed again.

Existing menu behavior should remain intact: New opens the New Game configuration screen, Cancel returns to the main menu, Settings opens its placeholder screen, Back returns to the main menu, disabled Load does nothing, and Quit closes the app from the main menu.

The documentation is accepted when `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, and `ROADMAP.md` describe the implemented current behavior without claiming exploration, base-building, resources, tower-defense combat, saves, real settings, campaign, assets, or packaging exist.

## Idempotence and Recovery

All edits are local to the executable, menu package, new game package, tests, documentation, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when the graphics environment is available.

## Artifacts and Notes

Important planned artifacts:

    plans/05-main-game-update-loop.md
    plans/05-main-game-update-loop/screenshots/main-menu.png
    plans/05-main-game-update-loop/screenshots/new-game-configuration.png
    plans/05-main-game-update-loop/screenshots/running-game.png
    plans/05-main-game-update-loop/screenshots/paused-game.png
    cmd/td/main.go
    cmd/td/main_test.go
    internal/game/game.go
    internal/game/game_test.go
    internal/menu/menu.go
    internal/menu/menu_test.go
    PRODUCT.md
    README.md
    ARCHITECTURE.md
    ROADMAP.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

At the end of this plan, `internal/game` should expose:

    type Input struct {
        TogglePause bool
    }

    func New(wizardName string, width, height int) (*State, error)
    func (s *State) Resize(width, height int)
    func (s *State) Update(input Input)
    func (s *State) Draw(screen *ebiten.Image)
    func (s *State) Updates() int
    func (s *State) Paused() bool
    func (s *State) WizardName() string

The `internal/menu` package should expose a new `ActionStart` action, keep `WizardName()` available for the executable, and cap Wizard names at 32 printable runes.
