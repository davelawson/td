# Add In-Game Top Bar HUD

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/08-top-bar-hud.md` because `plans/00-initial-ebitengine-menu.md` through `plans/07-game-design-control-doc.md` already exist.

## Purpose / Big Picture

The game view currently shows a Wizard name, placeholder field, debug update counter, pause state, and in-game overlay menu, but it does not yet show the player the strategic status that will matter during calm preparation and Raids. After this change, the game view has a persistent top bar that shows the current Chapter, Day, resources, phase-specific status, and Sanctum barricade charges. The values are fixed prototype data until resource, phase, and combat systems exist, but the display shape is observable and testable now.

## Progress

- [x] (2026-05-08 22:03Z) Inspected current game rendering, tests, screenshot capture, and root control documents.
- [x] (2026-05-08 22:08Z) Added prototype HUD status data, calm and raid formatting, top-bar rendering, and game package tests.
- [x] (2026-05-08 22:12Z) Updated screenshot capture to write evidence under `plans/08-top-bar-hud/screenshots/`.
- [x] (2026-05-08 22:14Z) Updated `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, `ARCHITECTURE.md`, and `README.md` for the new visible HUD.
- [x] (2026-05-08 22:18Z) Ran full tests, screenshot capture, launch validation, whitespace check, and final line-count review.
- [x] (2026-05-08 22:40Z) Refactored HUD status fields into `game.State` so the HUD code only presents state-owned values.

## Surprises & Discoveries

- Observation: The existing screenshot harness already captures running, paused, and in-game overlay views, so the new visual evidence can reuse that path with only a plan-directory change.
  Evidence: `cmd/td/main_test.go` has `TD_CAPTURE_SCREENSHOT` gated capture targets for the relevant game states.

- Observation: The top bar remains readable in the running-game screenshot and is dimmed consistently by the existing in-game overlay layer without colliding with overlay controls.
  Evidence: Reviewed `plans/08-top-bar-hud/screenshots/running-game.png` and `plans/08-top-bar-hud/screenshots/ingame-menu.png` after capture.

## Decision Log

- Decision: Add fixed prototype game-status data rather than real countdown, resource, or raid systems.
  Rationale: The user asked to expand the UI, and the underlying systems do not exist yet. Fixed status values make the HUD observable while avoiding premature gameplay scope.
  Date/Author: 2026-05-08 / Codex

- Decision: Keep calm and raid formatting testable even though the live game starts in calm phase.
  Rationale: The top bar needs a clear rendering path for both phases, but adding a temporary phase toggle would create prototype-only input behavior.
  Date/Author: 2026-05-08 / Codex

- Decision: Put HUD status and formatting in `internal/game/hud.go`.
  Rationale: The HUD is owned by the game view, and a separate file keeps formatting and drawing easy to review without creating a new package or scene framework.
  Date/Author: 2026-05-08 / Codex

- Decision: Keep all gameplay status values shown by the HUD directly on `game.State`.
  Rationale: Chapter, Day, phase, Raid pressure, Sanctum barricade charges, and resources are game-state facts that later systems will update. HUD code should format and draw those facts, not own a separate source of truth.
  Date/Author: 2026-05-08 / Codex

## Outcomes & Retrospective

Implementation completed the prototype in-game top bar. Starting a new game now shows Chapter and Day on the left, calm phase timing in the center, and Wood, Stone, Metal, and Barricade charges on the right. The game still starts in calm phase; raid formatting is covered by tests but is not yet reachable through gameplay because raid systems do not exist. The Wizard name was moved below the bar, and the debug update counter moved to the lower-right corner.

Validation results:

    go test ./...
    ok  	td/cmd/td	0.015s
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.792s

    file plans/08-top-bar-hud/screenshots/*.png
    plans/08-top-bar-hud/screenshots/ingame-menu.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/08-top-bar-hud/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/08-top-bar-hud/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/08-top-bar-hud/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/08-top-bar-hud/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

Final hand-written Go file line-count review:

    17 internal/ui/widgets.go
    19 internal/ui/colors.go
    29 internal/ui/text.go
    111 internal/game/hud.go
    129 internal/menu/start.go
    151 cmd/td/main.go
    163 internal/game/ingamemenu.go
    181 cmd/td/main_test.go
    193 internal/game/game.go
    234 internal/game/game_test.go
    282 internal/menu/menu_test.go
    343 internal/menu/menu.go
    1852 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

Follow-up refactor validation:

    go test ./...
    ok  	td/cmd/td	0.015s
    ok  	td/internal/game	0.019s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    git diff --check
    No whitespace errors.

Updated hand-written Go file line-count review after moving HUD status into `game.State`:

    17 internal/ui/widgets.go
    19 internal/ui/colors.go
    29 internal/ui/text.go
    99 internal/game/hud.go
    129 internal/menu/start.go
    151 cmd/td/main.go
    163 internal/game/ingamemenu.go
    181 cmd/td/main_test.go
    199 internal/game/game.go
    243 internal/game/game_test.go
    282 internal/menu/menu_test.go
    343 internal/menu/menu.go
    1855 total

No hand-written Go file exceeds or approaches the 600-line preference after the refactor.

## Context and Orientation

`td` is a local Go/Ebitengine desktop prototype. `cmd/td/main.go` owns Ebitengine startup, input polling, app-mode routing, and window layout. `internal/game/game.go` owns the active game state, placeholder field rendering, Wizard name display, pause behavior, debug update counter, and in-game overlay drawing. `internal/game/ingamemenu.go` owns the ESC overlay menu. This plan adds `internal/game/hud.go` for the top-bar status display.

The root control documents constrain the work. `PRODUCT.md` describes current user-visible behavior and must mention the new top bar. `GAME.md` records intended gameplay design and must note that the top bar shows resources, Chapter and Day, calm time before Raid, raid enemies remaining, and Sanctum barricade charges. `DESIGN.md` must capture the durable HUD readability principle. `ARCHITECTURE.md` must describe that `internal/game` owns the prototype top-bar display. `ROADMAP.md` and `README.md` should remain accurate about the current runnable shell. `CODESTYLE.md` requires `gofmt`, Go doc comments, tests for pure behavior where possible, and a final hand-written code-file line-count review.

## Plan of Work

Add a small amount of prototype gameplay status to `internal/game.State`. Define a private phase type with calm and raid values and a resource count struct for Wood, Stone, and Metal. Store Chapter name, Day, calm seconds before Raid, enemies remaining, barricade charges, and resource values directly on `State`. The HUD must not own a separate status object. The live game starts with fixed prototype values: `Chapter I: The Ashen Copse`, Day 1, Wood 80, Stone 45, Metal 12, calm phase, `Raid in 02:00`, raid enemy formatter value 12, and Barricade 3.

Render a full-width top bar at the top of the game view. The left segment shows `Chapter I: The Ashen Copse | Day 1`, the center segment shows `Raid in 02:00` in calm phase or `Enemies remaining: N` in raid phase, and the right segment shows `Wood 80  Stone 45  Metal 12 | Barricade 3`. Move the Wizard name downward and move the debug update counter away from the primary top-bar area so the HUD does not overlap existing text.

Add package tests for HUD formatting and initial state. Update the screenshot harness to save running game, paused game, and in-game menu screenshots under `plans/08-top-bar-hud/screenshots/`. Update the root control documents named above so future contributors understand which values are real current behavior and which are prototype placeholders.

## Concrete Steps

From the repository root, inspect the working tree:

    git status --short

Edit `internal/game/game.go`, add `internal/game/hud.go`, and update `internal/game/game_test.go`. Format and test the package:

    gofmt -w internal/game/game.go internal/game/hud.go internal/game/game_test.go
    go test ./internal/game

Update `cmd/td/main_test.go` so `TD_CAPTURE_SCREENSHOT=1` writes visual evidence under `plans/08-top-bar-hud/screenshots/`.

Update `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, `ARCHITECTURE.md`, and `README.md`. Then run:

    gofmt -w cmd/td/*.go internal/game/*.go internal/menu/*.go internal/ui/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, `git diff --check` reports no whitespace errors, screenshot capture writes PNG evidence under `plans/08-top-bar-hud/screenshots/`, and `go run ./cmd/td` opens without startup errors. A human should be able to start a new game and see a readable top bar showing Chapter, Day, resources, calm time before Raid, and Barricade charges. Pressing SPACE should still pause the game and show `PAUSED`; pressing ESC should still show the in-game overlay over the game view without the HUD colliding with overlay controls.

Documentation is accepted when `PRODUCT.md` and `README.md` describe the current visible top bar, `GAME.md` records the intended HUD information and the `Barricade` label, `DESIGN.md` records HUD readability expectations, `ARCHITECTURE.md` records `internal/game` ownership of the top bar, and `ROADMAP.md` remains accurate about the prototype foundation.

## Idempotence and Recovery

The changes are additive and local to game rendering, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when graphics are available.

If the HUD text overlaps at 1920x1080, reduce only the HUD font size or horizontal margins before changing gameplay data or introducing layout abstractions. Do not create a scene manager, asset pipeline, resource system, raid system, or timer system for this slice.

## Artifacts and Notes

Important artifacts:

    plans/08-top-bar-hud.md
    plans/08-top-bar-hud/screenshots/running-game.png
    plans/08-top-bar-hud/screenshots/paused-game.png
    plans/08-top-bar-hud/screenshots/ingame-menu.png
    internal/game/hud.go
    internal/game/game.go
    internal/game/game_test.go
    cmd/td/main_test.go
    PRODUCT.md
    ROADMAP.md
    GAME.md
    DESIGN.md
    ARCHITECTURE.md
    README.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

The new HUD-adjacent types and formatting methods are private to `internal/game` for now:

    type phase int
    const phaseCalm phase = ...
    const phaseRaid phase = ...
    type resourceCounts struct { wood, stone, metal int }
    func (s *State) setPrototypeGameStatus()
    func (s *State) chapterDayText() string
    func (s *State) phaseText() string
    func (s *State) resourcesAndBarricadeText() string

`State.Draw` should call the new top-bar drawing method before drawing the Wizard name, debug counter, and in-game overlay.
