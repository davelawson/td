# Add New Game Configuration Screen

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan must be maintained according to `PLANS.md` in the repository root. It is saved at `plans/03-new-game-configuration.md` because `plans/00-initial-ebitengine-menu.md`, `plans/01-expanded-main-menu.md`, and `plans/02-menu-package-refactor.md` already exist.

## Purpose / Big Picture

The current app has a New button, but it only opens placeholder copy. After this change, clicking New opens the first real new-game configuration screen. A player can type a Wizard name, see a disabled Start button reserved for future gameplay, and use Cancel to return to the main menu. The feature remains deliberately small: it adds configuration UI only, not game startup, save data, campaign structure, or gameplay.

## Progress

- [x] (2026-05-07 23:00Z) Created this self-contained ExecPlan from the accepted implementation plan.
- [x] (2026-05-07 23:01Z) Replaced the New Game placeholder with a Wizard name configuration screen.
- [x] (2026-05-07 23:01Z) Added focused Wizard name text input behavior and tests.
- [x] (2026-05-07 23:01Z) Added disabled Start and active Cancel buttons, with Cancel returning to the main menu.
- [x] (2026-05-07 23:02Z) Updated product-facing documentation and screenshot capture paths.
- [x] (2026-05-07 23:03Z) Ran formatting, automated tests, screenshot capture, launch validation, whitespace checks, git status, and the required hand-written code file line-count review.

## Surprises & Discoveries

- Observation: Side-by-side New Game buttons exposed that button labels were centered across the whole screen, which only worked before because every button was horizontally centered.
  Evidence: `internal/menu/menu.go` now uses button-bounds centering in `drawCenteredButtonText`, so `Cancel` and `Start` render centered inside their own rectangles.

- Observation: Screenshot capture succeeded for the new active plan directory.
  Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` completed successfully and wrote three PNG files under `plans/03-new-game-configuration/screenshots/`.

## Decision Log

- Decision: Implement Wizard name entry as an editable keyboard text field that is focused automatically when the New Game screen opens.
  Rationale: This matches the requested option to select the Wizard name while keeping the first slice small and testable.
  Date/Author: 2026-05-08 / Codex

- Decision: Render Start as disabled while it has no behavior.
  Rationale: A disabled button communicates future intent without creating a clickable control that appears broken.
  Date/Author: 2026-05-08 / Codex

- Decision: Keep this feature inside `internal/menu` rather than introducing a scene manager or game package.
  Rationale: The current screen is still menu-owned configuration UI. A broader scene abstraction would be premature before gameplay screens exist.
  Date/Author: 2026-05-08 / Codex

## Outcomes & Retrospective

Implementation completed the first New Game configuration screen. The main menu still owns the same top-level flow, but `New` now opens a configuration panel with a focused Wizard name field, disabled `Start`, and active `Cancel`. The text field accepts printable typed runes, removes the final rune on Backspace, and caps the name at 16 runes. `Start` is visible but inert because gameplay startup is not implemented. `Cancel` returns to the main menu, and `Quit` still closes the app through `cmd/td`.

Validation results:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/menu	(cached)

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.679s

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

    git status --short
    M ARCHITECTURE.md
    M PRODUCT.md
    M README.md
    M cmd/td/main.go
    M cmd/td/main_test.go
    M internal/menu/menu.go
    M internal/menu/menu_test.go
    ?? plans/03-new-game-configuration.md
    ?? plans/03-new-game-configuration/

Final hand-written Go file line-count review:

    70 cmd/td/main.go
    91 cmd/td/main_test.go
    185 internal/menu/menu_test.go
    408 internal/menu/menu.go
    754 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local PC tower-defense prototype written in Go with Ebitengine. Ebitengine owns the desktop window, game loop, drawing surface, and input APIs. The current executable is `cmd/td/main.go`; it configures the Ebitengine window, creates an `app`, and delegates update and drawing work to `internal/menu`.

The `internal/menu` package owns the current menu screen state, menu rendering, button hit testing, disabled-target handling, action selection, and placeholder New Game and Settings screens. Its tests live in `internal/menu/menu_test.go` and exercise pure behavior without opening a graphics window. Screenshot evidence is captured by a gated test in `cmd/td/main_test.go` when `TD_CAPTURE_SCREENSHOT=1` is set.

Several root control documents constrain this work. `PRODUCT.md` records current user-visible behavior and must change because New Game is no longer a placeholder. `README.md` describes the current runnable slice and should match the new user-visible workflow. `ARCHITECTURE.md` describes the current menu flow and should mention the configuration screen. `DESIGN.md` requires readable text, stable hit targets, clear hover feedback, and restrained medieval wizardry styling. `CODESTYLE.md` requires `gofmt`, Go doc comments for functions and methods, focused functions, package-level tests for pure behavior, and a final line-count review for hand-written code files.

## Plan of Work

First, update `internal/menu/menu.go` so `Menu.Update` accepts a small `Input` value rather than separate cursor and click arguments. The executable will still poll Ebitengine for the cursor, mouse click, typed runes, and backspace key state, but the menu package will own how those inputs affect menu state. The `Input` value should include cursor coordinates, whether the left mouse button was just clicked, the runes typed this frame, and whether Backspace was just pressed.

Next, replace the New Game placeholder with a real configuration screen. Add menu-owned state for the Wizard name and whether the Wizard name field is focused. Opening the New Game screen through the New button should focus the field automatically. The field starts empty. While the field is focused, printable typed runes append to the name, and Backspace removes one rune. Cap the Wizard name at 16 runes so it fits the current UI. Do not validate empty names yet because Start is disabled and cannot launch gameplay.

Then, update the New Game buttons. The New Game screen should have an active Cancel button using the existing back-to-main action behavior, and a visible disabled Start button that returns `ActionNone` when clicked. The Settings screen can keep its existing Back button. `cmd/td/main.go` should continue to translate only `menu.ActionQuit` into `ebiten.Termination`.

After behavior is implemented, update tests. Keep the existing button bounds, action selection, disabled button, screen routing, and quit tests. Add focused tests for entering the New Game screen, typing a Wizard name, Backspace removal, 16-rune truncation, Cancel returning to main, and disabled Start doing nothing. Update the screenshot test to write artifacts under `plans/03-new-game-configuration/screenshots/` and capture at least the main menu, New Game configuration screen, and Settings screen.

Finally, update `PRODUCT.md`, `README.md`, and `ARCHITECTURE.md` so their current behavior descriptions match the implementation. This feature does not change `ROADMAP.md`, `CODESTYLE.md`, or `DESIGN.md`.

## Concrete Steps

From the repository root, inspect the current state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Edit `internal/menu/menu.go` to add the input struct, Wizard name state, New Game buttons, text input handling, and New Game configuration rendering. Edit `cmd/td/main.go` so it builds `menu.Input` from Ebitengine APIs, including typed characters and Backspace. Run:

    gofmt -w internal/menu/menu.go cmd/td/main.go
    go test ./...

Edit `internal/menu/menu_test.go` and `cmd/td/main_test.go` for the new update signature, new behavior tests, and new screenshot paths. Run:

    gofmt -w internal/menu/menu_test.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected screenshot artifacts:

    plans/03-new-game-configuration/screenshots/main-menu.png
    plans/03-new-game-configuration/screenshots/new-game-configuration.png
    plans/03-new-game-configuration/screenshots/settings-placeholder.png

Update `PRODUCT.md`, `README.md`, and `ARCHITECTURE.md`. Then run final validation:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/menu/menu.go internal/menu/menu_test.go
    go test ./...
    go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds and `go run ./cmd/td` opens the app with the updated menu workflow. A human should observe that clicking New opens a New Game screen with a Wizard name field, disabled Start button, and Cancel button. Typing while the field is focused should update the Wizard name. Backspace should remove the last typed rune. Clicking Cancel should return to the main menu. Clicking disabled Start should not leave the screen or report a start action. Clicking Quit from the main menu should still close the app cleanly.

The documentation is accepted when `PRODUCT.md`, `README.md`, and `ARCHITECTURE.md` describe New Game as a configuration screen rather than a placeholder. The visual evidence is accepted when the screenshot capture test writes the three PNG files under `plans/03-new-game-configuration/screenshots/`.

## Idempotence and Recovery

All edits are local to menu code, the executable input adapter, tests, documentation, screenshots, and this plan. Re-running `gofmt`, `go test ./...`, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries` and rely on automated tests plus manual launch validation where possible.

## Artifacts and Notes

Important planned artifacts:

    plans/03-new-game-configuration.md
    plans/03-new-game-configuration/screenshots/main-menu.png
    plans/03-new-game-configuration/screenshots/new-game-configuration.png
    plans/03-new-game-configuration/screenshots/settings-placeholder.png
    cmd/td/main.go
    cmd/td/main_test.go
    internal/menu/menu.go
    internal/menu/menu_test.go
    PRODUCT.md
    README.md
    ARCHITECTURE.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

At the end of this plan, `internal/menu` should continue exposing the existing `Action`, `Button`, `Button.Contains`, `ActionAt`, `Screen`, `Menu`, `New`, `Draw`, `Screen`, and `SetScreenForTest` concepts. `Menu.Update` should accept one input value, similar to:

    type Input struct {
        CursorX int
        CursorY int
        Clicked bool
        Typed   []rune
        Backspace bool
    }

The exact field names can follow local Go style, but the data flow must remain the same: `cmd/td` polls Ebitengine and `internal/menu` owns interpretation.
