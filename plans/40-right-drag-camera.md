# Add Right-Drag Camera Panning

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/40-right-drag-camera.md` because `plans/39-building-menu-tooltips.md` was the highest existing prefixed plan when this work began.

## Purpose / Big Picture

The game already supports map inspection through mouse-wheel zoom and `WASD` panning, but mouse-only map movement is missing. After this change, a player can press the right mouse button over the game view, keep it held, and drag to move the viewport as if grabbing the visible world. Dragging is proportional to the current zoom because each cursor-pixel delta is converted through the current world-to-screen scale.

This is inspection only. It does not reveal Plots, move a wizard character, mutate map data, select objects, place structures, or interact with UI surfaces.

## Progress

- [x] (2026-06-27 00:00Z) Inspected current camera, input polling, UI hitboxes, tests, control documents, and prior camera plan.
- [x] (2026-06-27 00:00Z) Chose grab-world drag behavior with the user.
- [x] (2026-06-27 00:00Z) Added right-button game input, private camera drag state, UI-blocked drag starts, zoom-scaled drag panning, and executable input polling.
- [x] (2026-06-27 00:00Z) Added focused right-drag camera tests in a new test file to avoid expanding `internal/game/game_test.go` beyond the 600-line preference.
- [x] (2026-06-27 00:00Z) Updated README, PRODUCT, ROADMAP, GAME, DESIGN, and ARCHITECTURE for the new current behavior and design decision.
- [x] (2026-06-27 00:00Z) Ran full validation commands, captured screenshots, reviewed line counts, and recorded outcomes.

## Surprises & Discoveries

- Observation: `internal/game/game_test.go` was already near the 600-line preference before this work.
  Evidence: The initial line-count review showed `internal/game/game_test.go` at 597 lines.

## Decision Log

- Decision: Use grab-world right-drag behavior.
  Rationale: Dragging the map content with the cursor matches common strategy-map and editor behavior, and the user selected this behavior before implementation.
  Date/Author: 2026-06-27 / Codex

- Decision: Start a camera drag only when the initial right-button press is over the game view and outside screen-space UI, but allow a valid drag to continue over UI until release.
  Rationale: UI should not unexpectedly start camera movement, while an active drag should stay stable if the cursor crosses a UI panel during the gesture.
  Date/Author: 2026-06-27 / Codex

- Decision: Keep right-drag camera state private in `internal/game/camera.go`.
  Rationale: The camera is still specific to the prototype game view, and a new package would add indirection before there is another camera owner.
  Date/Author: 2026-06-27 / Codex

## Outcomes & Retrospective

Right-drag camera panning is implemented. The game now accepts distinct right mouse input, starts camera drags only from the game view, keeps valid drags active across later UI overlap, converts cursor deltas through the current zoom scale, and preserves existing wheel zoom, `WASD` pan, left-click selection, left-drag building placement, pause, and overlay behavior.

Validation results:

    go test ./internal/game
    ok  	td/internal/game	0.483s

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.023s
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/sound	[no test files]
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	1.546s

    timeout 5s go run ./cmd/td
    Exit code 124 after the app launched and stayed open until timeout stopped it. This is expected for the bounded launch check.

    git diff --check
    No whitespace errors.

Screenshot evidence was written under `plans/40-right-drag-camera/screenshots/`; all captured files are 1920x1080 PNG images.

Final hand-written Go file line-count review found no file over 600 lines. The largest files are:

    572 internal/game/building_bar_test.go
    597 internal/game/game_test.go

No extra split is required by the 600-line preference. `internal/game/game_test.go` remains close to the preference, so this plan kept the new coverage in `internal/game/camera_drag_test.go` rather than expanding it further.

## Context and Orientation

`td` is a local Go/Ebitengine desktop game prototype. `cmd/td/main.go` owns Ebitengine input polling and forwards a compact `game.Input` value to `internal/game`. `internal/game/game.go` owns game state and update ordering. `internal/game/camera.go` owns camera zoom, keyboard panning, scene viewport geometry, and world-to-screen projection. The current map is one 15x15 home Plot rendered below the top HUD.

World positions are measured in Tile units, with the Sanctum at `(0, 0)` and positive Y pointing north. Projection uses `plotBaseTileSize * camera.zoom` pixels per Tile. For right-drag panning, the cursor delta in screen pixels must be divided by that same scale so panning moves fewer world Tiles when zoomed in and more world Tiles when zoomed out.

`PRODUCT.md` and `README.md` describe current user-visible behavior. `GAME.md` records intended gameplay design decisions. `DESIGN.md` records interaction expectations. `ARCHITECTURE.md` records code ownership and boundaries. `ROADMAP.md` records product sequencing. This change affects those documents because it adds a current player-facing camera interaction.

## Plan of Work

Extend `game.Input` with distinct right mouse fields so right-drag camera movement cannot collide with existing left-click selection or left-drag building placement. Poll those fields in `cmd/td/main.go` using Ebitengine's right mouse button APIs.

Add a private `cameraDragState` to `internal/game`. When `State.applyCameraInput` runs, continue applying existing wheel and `WASD` behavior, then apply right-drag behavior. On a valid right-button press, store the cursor as the drag anchor. On later held frames, compute cursor delta from the previous frame and update camera center using grab-world semantics: `centerX -= deltaX / scale` and `centerY += deltaY / scale`, where `scale` is `plotBaseTileSize * camera.zoom`. Clear drag state when the right button is released or no longer down.

A valid drag start is inside `sceneViewport()` and outside `buildingBarContains`, `nextRaidButtonContains`, and `selectionPanelContains`. The top bar is blocked by the viewport check because `sceneViewport()` starts below `topBarHeight`. The in-game overlay blocks all camera input through existing `State.Update` ordering.

Add focused tests in `internal/game/camera_drag_test.go` for valid starts, grab-world direction, zoom scaling, UI-blocked starts, continuation over UI after a valid start, release cleanup, paused operation without logical updates, overlay blocking, and map immutability.

Update the screenshot capture base path in `cmd/td/main_test.go` to `plans/40-right-drag-camera/screenshots`. Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md`. Do not update `ART.md` or `CODESTYLE.md`, and do not add dependencies, camera bounds, cursor-centered zoom, minimaps, selection changes, or new gameplay systems.

## Concrete Steps

Run commands from the repository root, `/home/dave/dev/ai/td`.

Inspect the tree:

    git status --short

After edits, format Go files:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go

Run focused and full tests:

    go test ./internal/game
    go test ./...

Capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Run a bounded launch check:

    timeout 5s go run ./cmd/td

Check whitespace and pending files:

    git diff --check
    git status --short

Check hand-written Go file line counts:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

Files over 600 lines should be reported in `Outcomes & Retrospective` with a recommendation. Unplanned refactors, code splits, or library additions require user approval before implementation.

## Validation and Acceptance

Automated acceptance requires `go test ./...` to pass, screenshot capture to write PNG evidence under `plans/40-right-drag-camera/screenshots/`, `git diff --check` to report no whitespace errors, and the line-count review to be recorded.

Manual acceptance is: start the app with `go run ./cmd/td`, enter a Wizard name, start a game, right-press over the map, and drag. The map should follow the cursor. The same drag distance should move less world space when zoomed in and more world space when zoomed out. Starting a right-drag over the top HUD, building bar, Next Raid button, selection panel, or overlay should not move the camera. A drag that starts on the map may continue if the cursor passes over UI. SPACE-paused play should still allow camera inspection, and ESC overlay should block it.

Documentation acceptance is: `PRODUCT.md` and `README.md` describe right-drag panning as current behavior; `GAME.md` records it as inspection rather than wizard movement; `DESIGN.md` says UI surfaces should not start camera drag; `ARCHITECTURE.md` records input polling and camera ownership; and `ROADMAP.md` reflects keyboard plus right-drag camera panning as part of completed scene inspection.

## Idempotence and Recovery

The change is additive and local. Re-running `gofmt`, tests, screenshot capture, and line-count checks is safe. If screenshot capture fails because the local graphics environment cannot create an Ebitengine window, record the exact error in this plan and use automated tests plus a later manual launch as the available validation. If drag direction feels inverted in manual testing, adjust only the two camera-center delta signs and update the direction tests and decision log.

## Artifacts and Notes

Important artifacts:

    plans/40-right-drag-camera.md
    plans/40-right-drag-camera/screenshots/
    internal/game/camera.go
    internal/game/camera_drag_test.go
    internal/game/game.go
    cmd/td/main.go
    cmd/td/main_test.go
    README.md
    PRODUCT.md
    ROADMAP.md
    GAME.md
    DESIGN.md
    ARCHITECTURE.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

`game.Input` must include these fields:

    RightPressed bool
    RightDown bool
    RightReleased bool

`internal/game.State` must store private `cameraDragState`. The drag state must not be exported. Tests may inspect it because package-level tests use package `game`.

## Revision Note

This plan was created during implementation of right-drag camera panning to capture the selected drag semantics, UI-blocking behavior, tests, documentation updates, validation commands, and required final line-count review.
