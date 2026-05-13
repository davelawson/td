# Add Basic Camera Zoom And Pan

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/10-basic-camera.md` because `plans/00-initial-ebitengine-menu.md` through `plans/09-static-prototype-scene.md` already exist.

## Purpose / Big Picture

The game screen currently renders a static home Plot, but the player cannot inspect it as a world. After this change, the player can use the mouse wheel to zoom the map in or out and use `W`, `A`, `S`, and `D` to pan the camera around the home Plot. The camera continues to work while paused so the player can inspect the map without advancing game logic, while the ESC in-game overlay still blocks map controls until it closes.

## Progress

- [x] (2026-05-13 01:20Z) Inspected the current game package, screenshot harness, control documents, and prior plan style.
- [x] (2026-05-13 01:30Z) Added camera state, input fields, projection rendering, and pure game-package tests.
- [x] (2026-05-13 01:30Z) Updated screenshot capture to write plan 10 evidence.
- [x] (2026-05-13 01:38Z) Updated README, PRODUCT, ROADMAP, GAME, DESIGN, and ARCHITECTURE for camera behavior.
- [x] (2026-05-13 01:44Z) Ran full validation, screenshot capture, bounded launch check, whitespace check, visual screenshot review, git status, and final hand-written Go file line-count review.

## Surprises & Discoveries

- Observation: The bounded launch check exits with code 124 because the app successfully stays open until `timeout` stops it.
  Evidence: `timeout 5s go run ./cmd/td` produced no startup error and exited after the timeout window.

## Decision Log

- Decision: Keep camera state private inside `internal/game` rather than creating a new package.
  Rationale: The camera currently serves only the prototype game scene. A separate package would imply durable reuse before additional gameplay screens or rendering systems exist.
  Date/Author: 2026-05-13 / Codex

- Decision: Treat the scene area below the HUD as the camera viewport, with zoom centered on that viewport rather than on the cursor.
  Rationale: This matches the requested centered zoom behavior and avoids adding tile hover, cursor anchoring, or selection semantics before the game needs them.
  Date/Author: 2026-05-13 / Codex

## Outcomes & Retrospective

Implementation completed the basic camera slice. Starting a new game now creates private camera state centered on the home Plot at zoom `1.0`. Mouse-wheel input zooms the map around the scene viewport center, `WASD` pans the camera without bounds, camera input works while paused, and the ESC overlay continues to block camera input while it is open. Camera controls do not mutate the stored prototype map.

Validation results:

    go test ./...
    ok  	td/cmd/td	0.017s
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.771s

    file plans/10-basic-camera/screenshots/*.png
    plans/10-basic-camera/screenshots/ingame-menu.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/10-basic-camera/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/10-basic-camera/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/10-basic-camera/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/10-basic-camera/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

Final hand-written Go file line-count review:

    17 internal/ui/widgets.go
    20 internal/ui/colors.go
    23 internal/game/colors.go
    29 internal/ui/text.go
    53 internal/game/map.go
    62 internal/game/scene.go
    87 internal/game/camera.go
    100 internal/game/hud.go
    129 internal/menu/start.go
    157 cmd/td/main.go
    159 internal/game/ingamemenu.go
    181 cmd/td/main_test.go
    182 internal/game/game.go
    282 internal/menu/menu_test.go
    343 internal/menu/menu.go
    462 internal/game/game_test.go
    2286 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local Go/Ebitengine desktop game prototype. `cmd/td/main.go` owns the Ebitengine window, app-mode routing, and input polling. `internal/game/` owns active game state, pause behavior, prototype map data, static home Plot rendering, the top-bar HUD, and the in-game overlay menu. The current map is one 15x15 `Plot` stored in `State.gameMap.Home`, with a centered `Sanctum` and a straight north road.

A camera is the game state that decides which world coordinates appear in the scene viewport. World coordinates here are unscaled scene pixels derived from map tile coordinates and a base tile size. Projection means converting a world position to a screen position. The projection for this plan is `screen = viewportCenter + (world - cameraCenter) * zoom`.

The root control documents constrain the work. `PRODUCT.md` and `README.md` must describe mouse-wheel zoom and `WASD` panning as current game-screen behavior. `ROADMAP.md` must show basic camera movement as completed or part of current basic scene interaction. `GAME.md` must record that early map inspection uses camera movement instead of wizard movement. `DESIGN.md` must note that camera controls preserve HUD readability and do not affect overlay behavior. `ARCHITECTURE.md` must record that `internal/game` owns prototype camera state and map projection. `CODESTYLE.md` requires `gofmt`, doc comments for Go functions and methods, tests for pure behavior, and a final hand-written code-file line-count review against the 600-line preference.

## Plan of Work

Extend `internal/game.Input` with `WheelY float64` and `PanUp`, `PanDown`, `PanLeft`, and `PanRight bool`. Update `cmd/td/main.go` so `updateGame` reads `_, wheelY := ebiten.Wheel()` and polls `W`, `A`, `S`, and `D` with `ebiten.IsKeyPressed`.

Add private camera state in `internal/game`, likely in `camera.go`. The camera should store `zoom`, `centerX`, and `centerY`. New games start at zoom `1.0` and center on the home Plot in world coordinates. The camera should use a small positive minimum zoom such as `0.1` to prevent zero or negative rendered sizes after repeated scroll-down input. There is no maximum zoom and no pan clamping, so the Plot can move fully off-screen.

Update `State.Update` so overlay-open input is handled first and blocks camera changes. Outside the overlay, apply camera input even when the game is paused. Toggle menu and toggle pause still do not count as logical updates. If paused after camera input and toggle handling, return without incrementing the logical update counter.

Update `internal/game/scene.go` to render the home Plot through the camera. Use the current scene area below the HUD as the camera viewport. Each Tile should have a stable world rectangle derived from map coordinates and a base tile size. Project world rectangles to screen using the camera center and zoom. Keep HUD, wizard name, update counter, pause label, and overlay rendering in screen space so zooming and panning only affects the map scene.

Add tests in `internal/game/game_test.go` for initial camera zoom and center, wheel-up zoom, wheel-down zoom floor, `WASD` direction changes, camera input while paused, overlay-open camera blocking, and camera changes not mutating map data. Update screenshot capture in `cmd/td/main_test.go` to write evidence under `plans/10-basic-camera/screenshots/`.

Update the control documents named above. Do not add exploration, map expansion, tile selection, resource rules, combat, cursor inspection, drag panning, camera bounds, cursor-centered zoom, a minimap, a scene framework, or new dependencies.

## Concrete Steps

From the repository root, inspect the working tree:

    git status --short

Edit `internal/game/game.go`, `internal/game/scene.go`, add `internal/game/camera.go`, update `internal/game/game_test.go`, and update `cmd/td/main.go` and `cmd/td/main_test.go`. Format and test the game package:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go
    go test ./internal/game

Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, `ARCHITECTURE.md`, and this plan. Then run:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, screenshot capture writes PNG evidence under `plans/10-basic-camera/screenshots/`, `go run ./cmd/td` opens without startup errors, `git diff --check` reports no whitespace errors, and the line-count review finds no hand-written Go file over the 600-line preference.

A human should be able to start a new game and see the home Plot as before. Scrolling up should zoom into the Plot; scrolling down should zoom out without ever making zoom zero or negative. `W` should pan the camera up, `S` down, `A` left, and `D` right. SPACE should still pause the logical update counter, and camera controls should still work while paused. ESC should open the in-game overlay, and while that overlay is open, camera input should not apply. Surrender should still return to the main menu.

Documentation is accepted when `PRODUCT.md` and `README.md` describe camera controls as current behavior, `ROADMAP.md` reflects that basic camera movement is part of the basic scene interaction milestone, `GAME.md` records camera inspection instead of wizard movement for early map inspection, `DESIGN.md` covers HUD and overlay expectations with camera controls, and `ARCHITECTURE.md` records `internal/game` ownership of camera state and projection.

## Idempotence and Recovery

The changes are additive and local to game input, game state, rendering, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when graphics are available.

If panning or zooming makes the map unreadable at the initial view, adjust only camera constants such as base tile size, pan speed, or the technical minimum zoom before adding any new interaction modes. If tests show camera input increments the logical update counter while paused or overlay-open input changes the camera, fix `State.Update` ordering rather than special-casing tests.

## Artifacts and Notes

Important artifacts:

    plans/10-basic-camera.md
    plans/10-basic-camera/screenshots/running-game.png
    plans/10-basic-camera/screenshots/paused-game.png
    plans/10-basic-camera/screenshots/ingame-menu.png
    internal/game/camera.go
    internal/game/scene.go
    internal/game/game.go
    internal/game/game_test.go
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

`game.Input` must include these new fields:

    WheelY float64
    PanUp bool
    PanDown bool
    PanLeft bool
    PanRight bool

`game.State` must store a private camera value. Tests may inspect private camera fields because `internal/game/game_test.go` uses package `game`. The camera must support a default zoom of `1.0`, a technical minimum zoom of `0.1`, centered zoom changes, unrestricted pan, and projection of map world coordinates into the scene viewport below the HUD.

## Revision Note

This plan was created before implementation to capture the requested basic camera controls, the overlay/pause input ordering, required documentation updates, screenshot evidence, validation commands, and the final code-file line-count review required by `PLANS.md` and `CODESTYLE.md`. It was updated after implementation to record completed validation evidence, the visual screenshot review, and the final line-count outcome.
