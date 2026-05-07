# Expand Main Menu Actions

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan must be maintained according to `PLANS.md` in the repository root. It is saved at `plans/01-expanded-main-menu.md` because `plans/00-initial-ebitengine-menu.md` already exists.

## Purpose / Big Picture

This plan expands the first runnable shell from a one-button menu into a small menu flow. After the work is complete, a contributor can run `go run ./cmd/td`, see `New`, `Load`, `Settings`, and `Quit` on the main menu, click `New` or `Settings` to open placeholder screens, click `Back` to return to the main menu, and click `Quit` to close cleanly.

The plan intentionally keeps this slice small. `Load` is a disabled placeholder because no save system exists. `New` and `Settings` prove explicit screen transitions without introducing gameplay, real settings, a save/load implementation, a scene framework, an asset pipeline, campaign structure, release packaging, or CI.

## Progress

- [x] (2026-05-07T20:19Z) Confirmed the working tree was clean and `plans/01-expanded-main-menu.md` was the next unused ordered plan path.
- [x] (2026-05-07T20:19Z) Created this self-contained ExecPlan from the accepted plan.
- [x] (2026-05-07T20:21Z) Updated `internal/menu/` with new menu actions, disabled-button handling, and tests.
- [x] (2026-05-07T20:22Z) Updated `cmd/td/` with explicit main-menu, new-game placeholder, and settings placeholder modes.
- [x] (2026-05-07T20:23Z) Captured visual evidence under `plans/01-expanded-main-menu/screenshots/`.
- [x] (2026-05-07T20:24Z) Updated `PRODUCT.md`, `README.md`, and `ARCHITECTURE.md` to match the implemented menu flow.
- [x] (2026-05-07T20:25Z) Ran final validation: `go test ./...`, screenshot capture, launch check, `git diff --check`, and `git status --short`.
- [x] (2026-05-07T20:25Z) Checked hand-written Go file line counts; no file exceeded the 600-line preference.

## Surprises & Discoveries

- Observation: The first expanded-menu screenshot had cramped vertical spacing.
  Evidence: `plans/01-expanded-main-menu/screenshots/main-menu.png` initially showed the subtitle too close to the title and first button. The final layout moved the buttons down, reduced the title size, and regenerated the screenshot.

- Observation: The screenshot capture path can capture multiple game screens in one gated test run.
  Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` wrote all three PNG files under `plans/01-expanded-main-menu/screenshots/`.

## Decision Log

- Decision: Add `New`, `Load`, `Settings`, and `Quit` to the main menu.
  Rationale: The user asked to expand the main menu with those options, and this is the next observable shell improvement before gameplay systems exist.
  Date/Author: 2026-05-07 / Codex

- Decision: Make `New` and `Settings` open placeholder screens with `Back` buttons.
  Rationale: The user requested screen loading for those menu items, and placeholder screens prove the app can transition between states without inventing gameplay or settings systems.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep `Load` visible but disabled.
  Rationale: The user did not ask to implement saves, and the product documents explicitly say save/load and campaign structure are deferred.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep input pointer-only in this slice.
  Rationale: The user accepted pointer-only input during plan discussion. Keyboard navigation remains useful later but is not necessary to prove these screens.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep screen transitions in `cmd/td` as a small local `screenMode` enum.
  Rationale: There are only three simple screens, so a scene framework would add structure before repeated needs exist. This follows `ARCHITECTURE.md`.
  Date/Author: 2026-05-07 / Codex

- Decision: Add `handleAction` in `cmd/td` for action routing.
  Rationale: Pointer input itself depends on Ebitengine runtime state, but menu action routing can be tested directly without window automation.
  Date/Author: 2026-05-07 / Codex

## Outcomes & Retrospective

Implementation completed the expanded menu flow. The app now shows `New`, disabled `Load`, `Settings`, and `Quit` on the main menu. `New` opens a placeholder New Game screen with a `Back` button. `Settings` opens a placeholder Settings screen with a `Back` button. `Back` returns to the main menu, and `Quit` still terminates cleanly through `ebiten.Termination`.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/menu	(cached)

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.346s

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This validates startup in the available environment, but automated pointer-click validation was not available.

    git diff --check
    No whitespace errors.

Screenshot evidence:

    plans/01-expanded-main-menu/screenshots/main-menu.png
    plans/01-expanded-main-menu/screenshots/new-game-placeholder.png
    plans/01-expanded-main-menu/screenshots/settings-placeholder.png

Final hand-written Go file line-count review:

    46 internal/menu/menu.go
    74 internal/menu/menu_test.go
    131 cmd/td/main_test.go
    247 cmd/td/main.go
    498 total

No file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local PC tower-defense prototype written in Go with Ebitengine, a Go 2D game engine that owns the desktop window, input callbacks, drawing surface, and game loop. The current runnable app lives in `cmd/td/main.go`. It sets the window title to `td`, creates one menu button labeled `Quit`, renders a medieval wizardry main menu, and returns `ebiten.Termination` when the quit button is clicked.

Pure menu behavior lives in `internal/menu/menu.go`. A `Button` has a label, rectangular bounds, and an `Action`. `Button.Contains` checks whether a point is inside the button rectangle. `ActionAt` returns the first button action containing a point, or `ActionNone` when no button is hit. Tests in `internal/menu/menu_test.go` cover button bounds and action selection without opening a graphics window.

`cmd/td/main_test.go` contains a gated screenshot-capture test. The test is skipped during normal `go test ./...`; when `TD_CAPTURE_SCREENSHOT=1` is set, it runs the Ebitengine game loop, captures a frame, and saves a PNG under the active plan's screenshot directory. This plan should extend that capture path so evidence can be saved for the expanded main menu, the new-game placeholder screen, and the settings placeholder screen.

The root control documents constrain this work. `PRODUCT.md` says the current shipped capability is only a main menu with quit behavior and that save/load, campaign structure, gameplay, and settings do not exist. `ROADMAP.md` says the current phase should add one observable gameplay-facing slice at a time without building speculative architecture. `DESIGN.md` says rendered UI should use readable text, stable hit targets, clear hover feedback, and restrained medieval wizardry styling. `CODESTYLE.md` says Go code must use `gofmt`, every Go function and method must have a doc comment, reusable code should live under `internal/`, and hand-written code files should stay below 600 lines when practical. `ARCHITECTURE.md` says Ebitengine startup belongs in `cmd/td/`, pure state transitions and hit testing should be testable without a graphics window, and larger scene frameworks should wait until repetition appears.

## Plan of Work

First, update `internal/menu/menu.go`. Add `ActionNew`, `ActionSettings`, and `ActionBack` to the `Action` constants while preserving `ActionNone` as the zero value and `ActionQuit` as the quit signal. Add a `Disabled bool` field to `Button`; the zero value should mean enabled so existing callers remain natural. Keep `Button.Contains` focused only on geometry. Update `ActionAt` so it skips disabled buttons before hit testing. This makes disabled menu items visible but non-interactive without mixing rendering policy into menu geometry.

Next, update `internal/menu/menu_test.go`. Keep the existing boundary tests. Extend action-selection tests to include `New`, `Settings`, `Back`, and `Quit`, and add a test proving that a disabled `Load` button returns `ActionNone` even when the point is inside its rectangle.

Then update `cmd/td/main.go`. Add a small unexported screen mode type near the `game` type, with values for the main menu, the new-game placeholder, and the settings placeholder. Do not add a scene manager or package-level framework. The `game` struct should track the current screen and the current hover action. `Update` should choose the active button set for the current screen, update hover action with `menu.ActionAt`, and on left-click route actions as follows: `ActionNew` switches to the new-game placeholder, `ActionSettings` switches to the settings placeholder, `ActionBack` switches to the main menu, and `ActionQuit` returns `ebiten.Termination`.

Adjust rendering in `cmd/td/main.go`. The main menu should show four vertically stacked buttons: `New`, disabled `Load`, `Settings`, and `Quit`. Disabled `Load` should use muted styling and should not receive hover styling. The new-game placeholder screen should show a concise title and message explaining that the first expedition is not implemented yet, plus a `Back` button. The settings placeholder screen should show a concise settings title and only a `Back` button, with no fake settings controls. All screens should keep the existing medieval wizardry tone, readable text, stable hit targets, and 960 by 540 logical resolution.

After the app code is updated, extend `cmd/td/main_test.go` so screenshot capture can save multiple screens. The capture test should create a new game, capture the expanded main menu, set the screen mode to the new-game placeholder and capture it, then set the screen mode to the settings placeholder and capture it. Save the images under `plans/01-expanded-main-menu/screenshots/` with clear filenames such as `main-menu.png`, `new-game-placeholder.png`, and `settings-placeholder.png`.

Finally, update control documents whose durable truth changes. `PRODUCT.md` should describe the expanded menu workflow, the disabled `Load` limitation, and the two placeholder screens. `README.md` should describe the current playable slice as a menu flow rather than a quit-only menu. `ARCHITECTURE.md` should mention the explicit local screen mode in `cmd/td/` only if it remains part of the implementation. `ROADMAP.md`, `DESIGN.md`, and `CODESTYLE.md` should not change unless implementation decisions alter their durable guidance.

## Concrete Steps

From the repository root, inspect the starting state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Edit `internal/menu/menu.go` and `internal/menu/menu_test.go` as described above, then run:

    gofmt -w internal/menu/*.go
    go test ./internal/menu

Edit `cmd/td/main.go` to add screen modes, button sets, action routing, and placeholder rendering. Then run:

    gofmt -w cmd/td/main.go
    go test ./...

Edit `cmd/td/main_test.go` to capture three screenshots when `TD_CAPTURE_SCREENSHOT=1` is set. Run:

    gofmt -w cmd/td/main_test.go
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected screenshot artifacts:

    plans/01-expanded-main-menu/screenshots/main-menu.png
    plans/01-expanded-main-menu/screenshots/new-game-placeholder.png
    plans/01-expanded-main-menu/screenshots/settings-placeholder.png

Run the app manually:

    go run ./cmd/td

Observe the expanded menu. Click `New`, then `Back`. Click `Settings`, then `Back`. Confirm `Load` is visibly disabled and does nothing. Click `Quit` and confirm the app exits without a fatal error.

Update `PRODUCT.md`, `README.md`, and if needed `ARCHITECTURE.md`. Then run final validation:

    go test ./...
    git diff --check
    git status --short

Check hand-written code-file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, report the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when all of these are true:

Running `go test ./...` from the repository root succeeds.

Running `go run ./cmd/td` opens a desktop window titled `td`.

The main menu shows `New`, `Load`, `Settings`, and `Quit`.

Clicking `New` opens a placeholder new-game screen. Clicking its `Back` button returns to the main menu.

Clicking `Settings` opens a placeholder settings screen. Clicking its `Back` button returns to the main menu.

`Load` is visible but disabled. Clicking it does not change screens and does not quit.

Clicking `Quit` closes the app cleanly.

The screenshot artifacts under `plans/01-expanded-main-menu/screenshots/` show the expanded menu, the new-game placeholder, and the settings placeholder with readable text and stable targets consistent with `DESIGN.md`.

`PRODUCT.md` and `README.md` describe the implemented menu flow and current limitations. `ARCHITECTURE.md` matches the implemented screen-state structure if that structure becomes durable truth.

The final line-count review has been recorded in `Outcomes & Retrospective`.

## Idempotence and Recovery

The changes are additive and can be retried safely. If a screenshot capture run fails before all screenshots are written, delete only the incomplete files in `plans/01-expanded-main-menu/screenshots/` and rerun the gated capture test. If `go run ./cmd/td` cannot open a window because the local graphics environment is unavailable, record the exact error in `Surprises & Discoveries`, keep `go test ./...` as the primary automated validation, and do not change project code unless the error points to project behavior.

If the menu layout feels cramped after adding four buttons, adjust button spacing and panel height within `cmd/td/main.go` only. Do not introduce scrolling, nested menus, or a UI framework for this small menu flow.

## Artifacts and Notes

Important planned artifacts:

    plans/01-expanded-main-menu.md
    cmd/td/main.go
    cmd/td/main_test.go
    internal/menu/menu.go
    internal/menu/menu_test.go
    plans/01-expanded-main-menu/screenshots/main-menu.png
    plans/01-expanded-main-menu/screenshots/new-game-placeholder.png
    plans/01-expanded-main-menu/screenshots/settings-placeholder.png

## Interfaces and Dependencies

Use the existing Go module path `td` and the existing Ebitengine dependency `github.com/hajimehoshi/ebiten/v2`.

At the end of this plan, `internal/menu` should expose these concepts:

    type Action int

    const (
        ActionNone Action = iota
        ActionNew
        ActionSettings
        ActionBack
        ActionQuit
    )

    type Button struct {
        Label    string
        X        int
        Y        int
        W        int
        H        int
        Action   Action
        Disabled bool
    }

`Button.Contains(x, y int) bool` should remain a geometry check. `ActionAt(buttons []Button, x, y int) Action` should return `ActionNone` for disabled buttons and for misses.

`cmd/td` should keep screen mode local to the executable. It should not export a scene API or add new runtime dependencies.

Revision note, 2026-05-07 / Codex: Updated the living sections after implementation to record completed work, validation evidence, the screenshot layout adjustment, and final line-count results.
