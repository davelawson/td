# Initialize Ebitengine Main Menu

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan must be maintained according to `PLANS.md` in the repository root. It is saved at `plans/00-initial-ebitengine-menu.md` because no earlier ordered plan existed when it was created.

## Purpose / Big Picture

This plan turns `td` from a documentation-only prototype into a runnable local PC game app. After the work is complete, a contributor can run `go run ./cmd/td`, see a desktop Ebitengine window with a medieval wizardry main menu, activate a quit option, and observe the app closing cleanly.

The plan intentionally stops at the smallest playable shell. It does not implement tower-defense gameplay, exploration, base-building, resource gathering, saving, campaign structure, art pipelines, release packaging, or CI.

## Progress

- [x] (2026-05-06T21:06Z) Created the initial self-contained plan during project bootstrap.
- [ ] Verify local Go and Ebitengine environment assumptions before editing code.
- [ ] Record baseline evidence that no runnable app exists before implementation.
- [ ] Initialize the Go module and add the Ebitengine dependency.
- [ ] Create the Ebitengine app entry point and main-menu behavior.
- [ ] Add tests for menu hit testing or menu action selection without opening a graphics window.
- [ ] Capture the first rendered main-menu screenshot and review it against `DESIGN.md`.
- [ ] Update control documents if implementation changes current product behavior, architecture, style, or design guidance.
- [ ] Run `go test ./...` and `go run ./cmd/td` validation.
- [ ] Check hand-written code-file line counts and report any file over the 600-line preference before proposing extra work.

## Surprises & Discoveries

- Observation: No implementation has started yet.
  Evidence: The repository has no `go.mod`, `cmd/td/`, `internal/`, or runtime source files at plan creation.

## Decision Log

- Decision: Use Go and Ebitengine for the first local PC prototype.
  Rationale: The user explicitly selected Go and Ebitengine, and the first slice only needs a lightweight 2D desktop app.
  Date/Author: 2026-05-06 / Codex

- Decision: Keep the first playable slice to a main menu with quit behavior.
  Rationale: The user explicitly scoped the first playable version to only a menu and a quit option, which proves the build, window, input, and clean termination path before gameplay systems are introduced.
  Date/Author: 2026-05-06 / Codex

- Decision: Initialize Go and Ebitengine during this plan rather than during bootstrap.
  Rationale: The user asked to put Go and engine initialization in the first ExecPlan, and the project-bootstrap workflow should not implement product feature code.
  Date/Author: 2026-05-06 / Codex

- Decision: Use Ebitengine's desktop termination path by returning `ebiten.Termination` from `Update` when the quit option is activated.
  Rationale: Current Ebitengine package documentation recommends returning `Termination` from `Update` to halt desktop execution without making `RunGame` return an error.
  Date/Author: 2026-05-06 / Codex

## Outcomes & Retrospective

No implementation outcomes exist yet. At completion, summarize the files created, the exact validation results, screenshot evidence, and the final hand-written code-file line-count review. If any hand-written code file exceeds 600 lines, list the path, line count, likely cause, recommended response, and whether the user approved, deferred, or rejected extra refactor work.

## Context and Orientation

`td` is a local PC tower-defense prototype. `PRODUCT.md` says the current repository has no runnable game yet and that the first planned workflow is a desktop app with a main menu and quit option. `ROADMAP.md` says the long-term game should combine exploration, base-building, resource gathering, and tower-defense combat in a medieval wizardry setting, but those systems are not part of this plan. `DESIGN.md` says the first screen should be the actual game menu, with readable text, stable hit boxes, clear hover or focus feedback, and restrained medieval wizardry styling. `CODESTYLE.md` says Go code must use `gofmt`, every function or method must have a doc comment, reusable code should live under `internal/`, and hand-written code files should stay below 600 lines when practical. `ARCHITECTURE.md` says Ebitengine startup belongs in `cmd/td/` and pure behavior should be testable without opening a graphics window.

Ebitengine is a Go 2D game engine. In this project it will own the desktop window, game loop, drawing surface, and input APIs. The app should implement Ebitengine's `Game` interface with `Update`, `Draw`, and `Layout` methods. `Update` changes state and handles input, `Draw` renders to the screen, and `Layout` returns the logical screen size.

The repository currently has no `README` gaps after bootstrap, but it still has no Go module or runtime code. This plan creates those pieces.

## Plan of Work

First, verify the local environment. Run `go version` and confirm it is Go 1.22 or newer, because current Ebitengine installation guidance requires Go 1.22 or later. If the system is running under WSL and a native Windows window is desired, note that Ebitengine's install guidance recommends setting `GOOS=windows` for its example run. Do not change the project for WSL unless local validation requires it.

Next, record baseline evidence before implementation. Because no app exists yet, there is no screenshotable baseline. Create `plans/00-initial-ebitengine-menu/screenshots/` and save a short text note or transcript showing that `go run ./cmd/td` cannot run before initialization. This satisfies the baseline evidence requirement without inventing a fake screenshot.

Then initialize the Go module from the repository root with module path `td`. The module path can be renamed later if the repository gains a public import path. Add Ebitengine by importing `github.com/hajimehoshi/ebiten/v2` in code and running `go mod tidy`; the Go tool should add the required dependency to `go.mod` and `go.sum`.

Create the executable under `cmd/td/main.go`. Keep the first implementation small. Define a `Game` type with menu state, a `menuButton` or similarly named small type for the quit button, and `Update`, `Draw`, and `Layout` methods. Set a desktop window title such as `td`, choose a stable logical size such as 960 by 540 or another clearly documented 16:9 size, and draw the main menu directly with Ebitengine primitives and text. The title should read `td`, and the quit option should be visually clear. On pointer activation inside the quit button, return `ebiten.Termination` from `Update`.

If the button hit testing or menu action selection can be separated without creating an awkward abstraction, put that pure behavior in an internal package such as `internal/menu/` and add Go tests for it. If the first implementation remains entirely in `cmd/td/main.go`, still keep any pure helper functions easy to test or explain why no useful non-window test exists yet in `Outcomes & Retrospective`.

After the app runs, capture a screenshot of the main menu and save it under `plans/00-initial-ebitengine-menu/screenshots/main-menu.png`. Review the screenshot against `DESIGN.md`: text should be readable, the quit target should be obvious, the composition should feel like a game menu rather than a documentation page, and the palette should gesture toward medieval wizardry without becoming muddy or monochrome.

Finally, update control documents if the implementation differs from this plan. `PRODUCT.md` should say the menu app exists once it does. `README.md` should list working commands. `ARCHITECTURE.md`, `CODESTYLE.md`, and `DESIGN.md` should be updated only if implementation choices change durable structure, style, or design guidance.

## Concrete Steps

From the repository root, inspect the starting state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short
    go version

Expect `git status --short` to be empty before implementation begins. Expect `go version` to report Go 1.22 or newer.

Create baseline evidence:

    mkdir -p plans/00-initial-ebitengine-menu/screenshots
    go run ./cmd/td

Before implementation, `go run ./cmd/td` should fail because `cmd/td` does not exist. Save a concise transcript or note at `plans/00-initial-ebitengine-menu/screenshots/baseline-no-app.txt`, for example:

    Before implementation there is no runnable app and no screenshotable baseline.
    Command: go run ./cmd/td
    Expected result before implementation: package path ./cmd/td does not exist.

Initialize the module and dependency:

    go mod init td
    mkdir -p cmd/td

Create `cmd/td/main.go` with the Ebitengine app. Use `github.com/hajimehoshi/ebiten/v2` for the game loop and input. Use Ebitengine drawing and text utilities that keep the first app simple. Run:

    go mod tidy
    gofmt -w cmd/td/main.go

If pure menu behavior is factored into `internal/menu/`, create that package and its tests, then run:

    gofmt -w internal/menu/*.go
    go test ./...

Run the app:

    go run ./cmd/td

Observe a desktop window titled `td` with a main menu. Click the quit option and confirm the app closes cleanly. If `RunGame` returns `ebiten.Termination`, handle it as a clean shutdown rather than logging a fatal error.

Capture visual evidence after implementation. Save the screenshot as:

    plans/00-initial-ebitengine-menu/screenshots/main-menu.png

Review the screenshot against `DESIGN.md` and record notes in `Outcomes & Retrospective`.

Run final validation:

    go test ./...
    go run ./cmd/td
    git diff --check
    git status --short

Check hand-written code-file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written code file exceeds 600 lines, report the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when all of these are true:

Running `go test ./...` from the repository root succeeds.

Running `go run ./cmd/td` opens a desktop window titled `td`.

The window shows a main menu, including the visible title `td` and an obvious quit option.

Activating the quit option closes the app cleanly. The command should return without a fatal error caused by intentional quit behavior.

The screenshot at `plans/00-initial-ebitengine-menu/screenshots/main-menu.png` shows readable text, a stable quit target, and a visual direction consistent with `DESIGN.md`.

`PRODUCT.md` is updated to reflect that the menu app exists. `README.md` lists the now-working commands. Other control documents are updated only if implementation choices changed their durable truth.

The final line-count review has been recorded in `Outcomes & Retrospective`.

## Idempotence and Recovery

The plan is mostly additive. If `go mod init td` has already run, do not run it again; inspect the existing `go.mod` and continue. If `go mod tidy` changes dependency versions, keep the resulting `go.mod` and `go.sum` together. If the app fails to open because local graphics support is missing, record the exact error and validate tests first; then follow Ebitengine's official install guidance for the local operating system before changing project code.

If the first design attempt is visually unclear, adjust only the menu drawing and colors needed for readability and rerun `go test ./...` plus `go run ./cmd/td`. Do not expand into gameplay systems to make the menu feel more complete.

## Artifacts and Notes

Important planned artifacts:

    go.mod
    go.sum
    cmd/td/main.go
    internal/menu/ (optional, if pure menu behavior is factored out)
    plans/00-initial-ebitengine-menu/screenshots/baseline-no-app.txt
    plans/00-initial-ebitengine-menu/screenshots/main-menu.png

Official references checked during plan creation:

    Ebitengine install guidance says Go 1.22 or later is required and `go mod tidy` adds dependencies after imports exist.
    Ebitengine package documentation says desktop apps should return `ebiten.Termination` from `Update` to terminate cleanly.

## Interfaces and Dependencies

Use Go with module path `td`.

Use Ebitengine through:

    github.com/hajimehoshi/ebiten/v2

The executable should provide a type that satisfies Ebitengine's `Game` interface:

    Update() error
    Draw(screen *ebiten.Image)
    Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

If a separate menu package is created, it should expose only pure behavior needed by the executable. A small shape such as this is sufficient:

    type Action int

    const (
        ActionNone Action = iota
        ActionQuit
    )

    type Button struct {
        Label string
        X     int
        Y     int
        W     int
        H     int
        Action Action
    }

    func (b Button) Contains(x, y int) bool

Do not add a UI widget library, scene framework, ECS, asset manager, save system, campaign system, or packaging dependency in this plan unless a blocker appears and the user approves the added scope.
