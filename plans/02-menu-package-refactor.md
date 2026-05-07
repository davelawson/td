# Move Menu Ownership Out Of main.go

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan must be maintained according to `PLANS.md` in the repository root. It is saved at `plans/02-menu-package-refactor.md` because `plans/00-initial-ebitengine-menu.md` and `plans/01-expanded-main-menu.md` already exist.

## Purpose / Big Picture

The current app works, but `cmd/td/main.go` owns menu-specific details such as button layout, colors, hover state, screen routing, and drawing helpers. After this refactor, the user-visible menu should behave the same, while the executable package becomes a thin Ebitengine shell and `internal/menu` owns the menu UI. A contributor can see the result by running `go run ./cmd/td` and using the same `New`, disabled `Load`, `Settings`, `Back`, and `Quit` flow as before.

This plan does not add gameplay, real settings, save/load behavior, a scene manager, an asset pipeline, or new dependencies. It only moves responsibility to a clearer package boundary.

## Progress

- [x] (2026-05-07T21:28Z) Created this self-contained ExecPlan from the accepted refactor plan.
- [x] (2026-05-07T21:30Z) Moved menu state, rendering, palette, and menu action routing from `cmd/td/main.go` into `internal/menu`.
- [x] (2026-05-07T21:30Z) Slimmed `cmd/td/main.go` to Ebitengine lifecycle wiring and quit termination handling.
- [x] (2026-05-07T21:31Z) Moved screen-transition tests into `internal/menu` and updated screenshot capture to use menu-owned screen control.
- [x] (2026-05-07T21:31Z) Updated `ARCHITECTURE.md` so durable ownership matches the refactor.
- [x] (2026-05-07T21:32Z) Ran `gofmt`, `go test ./...`, gated screenshot capture, launch validation, whitespace checks, git status, and the required hand-written Go file line-count review.
- [x] (2026-05-07T21:45Z) Revised the executable shell so the Ebitengine-facing top-level structure is `app`, with a `mainMenu` field it delegates to.

## Surprises & Discoveries

- Observation: The existing screenshot capture test still works after moving screen state into `internal/menu`.
  Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` completed successfully and rewrote the three PNG files under `plans/01-expanded-main-menu/screenshots/`.

## Decision Log

- Decision: Move menu colors, button layout, screen state, hover tracking, and menu rendering into `internal/menu`.
  Rationale: These concepts describe the menu, not the executable package. Keeping them together makes the first runtime shell easier to extend without creating a general scene framework.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep Ebitengine startup, window configuration, and `ebiten.Termination` translation in `cmd/td`.
  Rationale: The executable package is still the process boundary, while `internal/menu` should remain reusable menu UI logic rather than process control.
  Date/Author: 2026-05-07 / Codex

- Decision: Preserve the existing user-visible menu behavior.
  Rationale: The request is a refactor, so successful implementation should be demonstrated by unchanged menu workflow plus cleaner package ownership.
  Date/Author: 2026-05-07 / Codex

- Decision: Name the Ebitengine-facing top-level structure `app` and give it a `mainMenu` field.
  Rationale: The application boundary should be explicit. Ebitengine calls `Update`, `Draw`, and `Layout` on the app, and the app delegates to the appropriate package-owned UI.
  Date/Author: 2026-05-07 / Codex

## Outcomes & Retrospective

Implementation completed the refactor without changing the intended menu workflow. `cmd/td/main.go` now only configures Ebitengine, constructs an `app`, delegates Ebitengine callbacks through that app, returns the fixed logical layout, and translates `menu.ActionQuit` into `ebiten.Termination`. The `app` has a `mainMenu` field and delegates menu input and drawing to it. `internal/menu` now owns menu screen state, palette, button layout, hover state, rendering, and non-quit menu action routing.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/menu	(cached)

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.707s

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for the bounded launch check.

    git diff --check
    No whitespace errors.

    git status --short
    M ARCHITECTURE.md
    M cmd/td/main.go
    M cmd/td/main_test.go
    M internal/menu/menu.go
    M internal/menu/menu_test.go
    M plans/01-expanded-main-menu/screenshots/main-menu.png
    M plans/01-expanded-main-menu/screenshots/new-game-placeholder.png
    M plans/01-expanded-main-menu/screenshots/settings-placeholder.png
    ?? plans/02-menu-package-refactor.md

Final hand-written Go file line-count review:

    64 cmd/td/main.go
    91 cmd/td/main_test.go
    118 internal/menu/menu_test.go
    280 internal/menu/menu.go
    553 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

Revision note, 2026-05-07 / Codex: Updated the living sections after implementation to record completed work, validation evidence, screenshot capture success, launch behavior, and final line-count results.

## Context and Orientation

`td` is a local PC tower-defense prototype written in Go with Ebitengine. Ebitengine owns the desktop window, game loop, drawing surface, and input APIs. The current executable lives in `cmd/td/main.go`; it sets up a 960 by 540 window, handles Ebitengine callbacks, and currently also owns menu palette, menu button layout, hover state, screen mode, and drawing code.

Pure menu hit testing currently lives in `internal/menu/menu.go`. It defines `Action`, `Button`, `Button.Contains`, and `ActionAt`. Tests in `internal/menu/menu_test.go` cover button bounds, disabled-target handling, and action selection without opening a graphics window. `cmd/td/main_test.go` currently tests screen action routing and contains a gated screenshot-capture test that writes PNG evidence under `plans/01-expanded-main-menu/screenshots/` when `TD_CAPTURE_SCREENSHOT=1` is set.

The root control documents constrain this work. `CODESTYLE.md` requires `gofmt`, doc comments for every Go function and method, focused functions, reusable packages under `internal/`, and a final line-count review for hand-written code files. `DESIGN.md` requires readable text, stable hit targets, clear hover feedback, and restrained medieval wizardry styling. `ARCHITECTURE.md` currently says `cmd/td` owns menu rendering and local screen-mode transitions; this refactor changes that durable boundary, so `ARCHITECTURE.md` must be updated. `PRODUCT.md` and `README.md` describe the current menu workflow; they should not change unless the implementation changes user-visible behavior.

## Plan of Work

First, expand `internal/menu/menu.go` from pure hit testing into the menu UI package. Keep the existing `Action`, `Button`, `Button.Contains`, and `ActionAt` behavior. Add an unexported screen mode type for the menu screens. Add a `Menu` type that stores width, height, current screen, button sets, hover action, and font faces. Move the palette values, backdrop drawing, panel drawing, button drawing, and centered text drawing from `cmd/td/main.go` into methods on `Menu`. Add `New(width, height int) (*Menu, error)` to construct fonts and menu layout. Add `Update(cursorX, cursorY int, clicked bool) Action` to update hover state, route non-quit menu actions internally, and return the selected action so the caller can terminate on `ActionQuit`. Add `Draw(screen *ebiten.Image)` to render the current menu state. Add a small exported `Screen` type and `SetScreenForTest(screen Screen)` only for screenshot capture and package tests; keep it clearly menu-scoped, not a generic scene API.

Next, slim `cmd/td/main.go`. The `game` struct should hold only a `*menu.Menu`. `newGame` should call `menu.New(screenWidth, screenHeight)`. `Update` should read cursor state and click state from Ebitengine, call `g.menu.Update`, and return `ebiten.Termination` only when the returned action is `menu.ActionQuit`. `Draw` should delegate to `g.menu.Draw`. `Layout` should continue returning the fixed logical resolution. Remove menu palette, button arrays, screen mode, font fields, drawing helpers, and action routing from `cmd/td/main.go`.

Then update tests. Move the current `TestHandleActionRoutesMenuScreens` behavior into `internal/menu/menu_test.go`, using `Menu.Update` and the test screen accessor or exposed screen values to verify that `ActionNew`, `ActionSettings`, `ActionBack`, `ActionNone`, and `ActionQuit` behave correctly. Update `cmd/td/main_test.go` so screenshot targets use `menu.ScreenMain`, `menu.ScreenNewGame`, and `menu.ScreenSettings`, and set the target through `game.menu.SetScreenForTest`. Keep the screenshot directory as `plans/01-expanded-main-menu/screenshots/` because the rendered output is expected to stay equivalent and this refactor is not a visual feature.

Finally, update `ARCHITECTURE.md`. Its codemap should say `cmd/td/` owns executable entry point, Ebitengine window setup, callback wiring, and process startup. It should say `internal/menu/` owns menu state, rendering, button hit testing, disabled-target handling, action selection, and placeholder menu screens. Its invariants should still say Ebitengine process startup belongs in `cmd/td/`, pure state transitions and hit testing remain testable without a graphics window, and no shared scene framework should be introduced yet.

## Concrete Steps

From the repository root, inspect the starting state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Edit `internal/menu/menu.go` to introduce `Menu`, menu screen state, rendering, and update routing. Edit `cmd/td/main.go` to delegate to the menu package. Run:

    gofmt -w internal/menu/menu.go cmd/td/main.go
    go test ./...

Edit `internal/menu/menu_test.go` and `cmd/td/main_test.go` to move state-transition coverage and adapt screenshot capture. Run:

    gofmt -w internal/menu/menu_test.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected screenshot artifacts:

    plans/01-expanded-main-menu/screenshots/main-menu.png
    plans/01-expanded-main-menu/screenshots/new-game-placeholder.png
    plans/01-expanded-main-menu/screenshots/settings-placeholder.png

Run the app manually or with a short timeout:

    go run ./cmd/td

Observe that the menu workflow is unchanged: `New` opens the New Game placeholder, `Settings` opens the Settings placeholder, `Back` returns to the main menu, `Load` stays disabled, and `Quit` exits cleanly.

Update `ARCHITECTURE.md`, then run final validation:

    go test ./...
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, report the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, `go run ./cmd/td` opens the same desktop menu, and the visible menu workflow is unchanged from the expanded main menu: `New` and `Settings` open their placeholders, `Back` returns to the main menu, `Load` is disabled, and `Quit` closes the app cleanly.

The package boundary is accepted when `cmd/td/main.go` no longer contains menu colors, button layout, menu screen modes, menu rendering helpers, or menu action routing other than translating `ActionQuit` to `ebiten.Termination`. `internal/menu` owns those concepts and has tests for pure state transitions and hit testing.

`ARCHITECTURE.md` must match the new ownership boundary. The final line-count review must be recorded in `Outcomes & Retrospective`.

## Idempotence and Recovery

The refactor is local to `cmd/td`, `internal/menu`, `ARCHITECTURE.md`, and this plan. If a test fails during migration, keep the old user-visible behavior as the source of truth and adjust package boundaries without changing menu workflow. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries` and rely on `go test ./...` plus launch validation where possible.

## Artifacts and Notes

Important planned artifacts:

    plans/02-menu-package-refactor.md
    cmd/td/main.go
    cmd/td/main_test.go
    internal/menu/menu.go
    internal/menu/menu_test.go
    ARCHITECTURE.md

## Interfaces and Dependencies

Use the existing Go module path `td` and the existing Ebitengine dependency `github.com/hajimehoshi/ebiten/v2`. Do not add new dependencies.

At the end of this plan, `internal/menu` should expose the existing `Action`, `Button`, `Button.Contains`, and `ActionAt` API, plus menu-owned runtime methods similar to:

    type Screen int

    const (
        ScreenMain Screen = iota
        ScreenNewGame
        ScreenSettings
    )

    type Menu struct { ... }

    func New(width, height int) (*Menu, error)
    func (m *Menu) Update(cursorX, cursorY int, clicked bool) Action
    func (m *Menu) Draw(screen *ebiten.Image)
    func (m *Menu) Screen() Screen
    func (m *Menu) SetScreenForTest(screen Screen)

`cmd/td` should not export a scene API or add new runtime dependencies.
