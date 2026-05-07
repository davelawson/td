# Set 1920x1080 Resolution And Pixel-Stable Text

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/04-resolution-and-pixel-text-scaling.md` because `plans/00-initial-ebitengine-menu.md`, `plans/01-expanded-main-menu.md`, `plans/02-menu-package-refactor.md`, and `plans/03-new-game-configuration.md` already exist.

## Purpose / Big Picture

The app currently opens with a 960 by 540 fixed logical layout. When the user enlarges the window, Ebitengine scales the whole rendered image, which also stretches menu text. After this change, a contributor can run `go run ./cmd/td`, see a 1920 by 1080 desktop window, resize it, and observe that menu geometry recenters while text keeps the same raw pixel size instead of being stretched.

## Progress

- [x] (2026-05-07 23:16Z) Created this ExecPlan for the accepted resolution and pixel-text scaling change.
- [x] (2026-05-07 23:18Z) Updated the Ebitengine window defaults and layout callback to use the current drawable size.
- [x] (2026-05-07 23:19Z) Added menu resize behavior and recomputed button geometry from the current drawable size.
- [x] (2026-05-07 23:20Z) Updated menu tests and screenshot capture paths for the 1920 by 1080 default.
- [x] (2026-05-07 23:23Z) Updated durable control documents for the new display policy.
- [x] (2026-05-07 23:24Z) Ran screenshot capture, launch validation, whitespace checks, git status, and the required hand-written code file line-count review.

## Surprises & Discoveries

- Observation: The screenshot test still referenced the old `screenWidth` and `screenHeight` constants after the executable constants were renamed.
  Evidence: `go test ./...` initially failed with `cmd/td/main_test.go:49:23: undefined: screenWidth` and `cmd/td/main_test.go:49:36: undefined: screenHeight`; updating the test to use the new default window constants fixed the failure.

## Decision Log

- Decision: Use Ebitengine `Layout` to return the actual outside window size, with a fallback to 1920 by 1080 for non-positive dimensions.
  Rationale: Returning the current size prevents Ebitengine from stretching a fixed framebuffer, while the fallback keeps startup and unusual platform reports deterministic.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep menu font sizes fixed in raw pixels and recompute only geometry when the window changes.
  Rationale: The user specifically requested text not stretch while upscaling. Fixed font face sizes satisfy that requirement and avoid introducing a premature UI scale setting.
  Date/Author: 2026-05-07 / Codex

- Decision: Keep resize behavior inside `internal/menu` rather than introducing a shared renderer or scene abstraction.
  Rationale: The current rendered behavior is still menu-owned, and the architecture documents warn against adding broader frameworks before gameplay screens create repeated needs.
  Date/Author: 2026-05-07 / Codex

## Outcomes & Retrospective

Implementation completed the resolution and resize-policy change. The app now opens with a 1920 by 1080 default window. `Layout` returns the current drawable size instead of a fixed 960 by 540 framebuffer, so Ebitengine does not stretch the complete rendered image when the window is enlarged. The menu package now recomputes centered button and panel geometry on resize while keeping text faces at fixed raw pixel sizes.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/menu	(cached)

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.686s

    file plans/04-resolution-and-pixel-text-scaling/screenshots/*.png
    plans/04-resolution-and-pixel-text-scaling/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/04-resolution-and-pixel-text-scaling/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/04-resolution-and-pixel-text-scaling/screenshots/settings-placeholder.png:   PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

    git status --short
    M ARCHITECTURE.md
    M DESIGN.md
    M PRODUCT.md
    M README.md
    M cmd/td/main.go
    M cmd/td/main_test.go
    M internal/menu/menu.go
    M internal/menu/menu_test.go
    ?? plans/04-resolution-and-pixel-text-scaling.md
    ?? plans/04-resolution-and-pixel-text-scaling/

Final hand-written Go file line-count review:

    75 cmd/td/main.go
    91 cmd/td/main_test.go
    243 internal/menu/menu_test.go
    461 internal/menu/menu.go
    870 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local Go/Ebitengine tower-defense prototype. Ebitengine is the Go 2D game engine that owns the desktop window, input callbacks, drawing surface, and game loop. The executable in `cmd/td/main.go` sets the window title and size, creates the app state, forwards input to `internal/menu`, delegates drawing, and implements `Layout`, the Ebitengine callback that chooses drawable dimensions.

The menu package in `internal/menu/menu.go` owns menu state, rendering, hit testing, disabled-target handling, screen routing, and Wizard name input. Before this plan, it stored button rectangles once at construction time from the original 960 by 540 size. Because this plan makes the drawable size follow the actual window, the menu must update its geometry when dimensions change.

The root control documents constrain the work. `PRODUCT.md` must mention the new current display behavior because default resolution and resizing are user-visible. `ARCHITECTURE.md` must replace references to a fixed logical layout with the pixel-sized layout policy. `DESIGN.md` must record that text remains raw-pixel-sized when windows are enlarged. `README.md` must describe the 1920 by 1080 default and resize policy for contributors. `CODESTYLE.md` requires `gofmt`, doc comments for Go functions and methods, tests for pure geometry behavior, and a final line-count review for hand-written code files.

## Plan of Work

First, edit `cmd/td/main.go`. Replace the old 960 by 540 constants with default window constants set to 1920 and 1080. Keep `ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)`. Change `newApp` to construct the menu with those default dimensions. Change `Layout` so it returns the actual `outsideWidth` and `outsideHeight`, falling back to the default dimensions if either value is non-positive, and call a menu resize method before returning.

Next, edit `internal/menu/menu.go`. Add a `Resize(width, height int)` method that ignores invalid dimensions, returns early when the size is unchanged, stores the new dimensions, recomputes button rectangles, and clears stale hover state. Replace constructor-time hard-coded button rectangles with a small `layoutButtons` helper that centers main, settings, and new-game buttons from the current menu dimensions. Update panel, field, backdrop accent, and text anchor calculations so they are centered in the current drawable area while keeping font sizes fixed.

Then, update tests. In `internal/menu/menu_test.go`, use 1920 by 1080 as the default test size and stop depending on absolute 960 by 540 click coordinates. Use active button rectangles for click tests, and add a resize test that verifies a 2560 by 1440 menu recenters the New button and still routes the click. In `cmd/td/main_test.go`, write screenshots under `plans/04-resolution-and-pixel-text-scaling/screenshots/` and capture a 1920 by 1080 frame.

Finally, update `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, and `DESIGN.md` to describe the new default resolution and pixel-stable resize policy. Do not update `ROADMAP.md` because this change does not alter future product sequencing or priorities. Do not add dependencies, runtime configuration, gameplay code, a settings UI, a scene manager, or a renderer package.

## Concrete Steps

From the repository root, inspect the current state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Edit `cmd/td/main.go` and `internal/menu/menu.go` for default 1920 by 1080 window size, current-size `Layout`, and menu-owned resize geometry. Run:

    gofmt -w cmd/td/main.go internal/menu/menu.go
    go test ./...

Edit `internal/menu/menu_test.go` and `cmd/td/main_test.go` for dynamic hit coordinates, resize coverage, and the new screenshot directory. Run:

    gofmt -w cmd/td/main_test.go internal/menu/menu_test.go
    go test ./...

Update `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, and `DESIGN.md`. Capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected screenshot artifacts:

    plans/04-resolution-and-pixel-text-scaling/screenshots/main-menu.png
    plans/04-resolution-and-pixel-text-scaling/screenshots/new-game-configuration.png
    plans/04-resolution-and-pixel-text-scaling/screenshots/settings-placeholder.png

Run final validation:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/menu/menu.go internal/menu/menu_test.go
    go test ./...
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, the gated screenshot capture writes all three 1920 by 1080 PNG files under this plan directory, and `go run ./cmd/td` opens a window titled `td` without startup errors. A human should observe that the default window is 1920 by 1080, that resizing the window keeps the menu centered in the current drawable area, and that text keeps the same raw pixel size instead of stretching with the window.

Existing menu behavior must remain unchanged: `New` opens the New Game configuration screen, typing edits the focused Wizard name field, `Cancel` returns to the main menu, `Settings` opens a placeholder screen, `Back` returns to the main menu, disabled `Load` and `Start` do nothing, and `Quit` closes the app cleanly.

The documentation is accepted when `PRODUCT.md`, `README.md`, `ARCHITECTURE.md`, and `DESIGN.md` describe the current 1920 by 1080 default and pixel-stable resize policy without claiming that settings, gameplay, save/load, or a configurable resolution menu exists.

## Idempotence and Recovery

All edits are local to the executable, menu package, tests, documentation, screenshots, and this plan. Re-running `gofmt`, `go test ./...`, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when the graphics environment is available.

## Artifacts and Notes

Important planned artifacts:

    plans/04-resolution-and-pixel-text-scaling.md
    plans/04-resolution-and-pixel-text-scaling/screenshots/main-menu.png
    plans/04-resolution-and-pixel-text-scaling/screenshots/new-game-configuration.png
    plans/04-resolution-and-pixel-text-scaling/screenshots/settings-placeholder.png
    cmd/td/main.go
    cmd/td/main_test.go
    internal/menu/menu.go
    internal/menu/menu_test.go
    PRODUCT.md
    README.md
    ARCHITECTURE.md
    DESIGN.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

At the end of this plan, `cmd/td/main.go` should define default window dimensions of 1920 by 1080 and `Layout(outsideWidth, outsideHeight int) (int, int)` should return the current window dimensions after updating menu geometry. The `internal/menu` package should expose:

    func New(width, height int) (*Menu, error)
    func (m *Menu) Resize(width, height int)
    func (m *Menu) Update(input Input) Action
    func (m *Menu) Draw(screen *ebiten.Image)

Existing `Action`, `Button`, `Button.Contains`, `ActionAt`, `Screen`, `SetScreenForTest`, `WizardName`, and `WizardNameFocused` behavior should remain available.
