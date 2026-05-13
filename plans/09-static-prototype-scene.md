# Add Static Home Plot Scene

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/09-static-prototype-scene.md` because `plans/00-initial-ebitengine-menu.md` through `plans/08-top-bar-hud.md` already exist.

## Purpose / Big Picture

The game view currently has menu flow, a top-bar HUD, pause behavior, and an in-game overlay, but it still needs a real gameplay-facing scene rather than a freehand placeholder field. After this change, starting a new game creates prototype map state and renders a static 15x15 home Plot. The Plot is empty except for the centered Sanctum and a straight road north to the Plot edge. Logical ticks may still advance the debug update counter, but they do not change the map or scene.

## Progress

- [x] (2026-05-13 00:25Z) Inspected the current game package, screenshot harness, and control documents.
- [x] (2026-05-13 00:31Z) Added `internal/game/map.go` with prototype `Map`, `Plot`, `Tile`, and default home Plot creation.
- [x] (2026-05-13 00:34Z) Stored the default map in `game.State` when a new game starts and replaced placeholder field rendering with static Plot rendering from map state.
- [x] (2026-05-13 00:37Z) Added pure tests for home Plot dimensions, centered Sanctum, north road, otherwise-empty tiles, and update ticks not mutating map state.
- [x] (2026-05-13 00:40Z) Updated screenshot capture to write evidence under `plans/09-static-prototype-scene/screenshots/`.
- [x] (2026-05-13 00:48Z) Updated `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` for the static home Plot scene.
- [x] (2026-05-13 00:56Z) Ran full tests, screenshot capture, launch validation, whitespace check, screenshot file inspection, visual review, git status, and final hand-written code-file line-count review.
- [x] (2026-05-13 01:10Z) Revised the prototype Plot size from 11x11 to 15x15 and updated code, docs, tests, and this plan to match.

## Surprises & Discoveries

- Observation: The existing in-game overlay still depended on the previous field color aliases after placeholder field rendering was removed.
  Evidence: `go test ./internal/game` initially failed with undefined `fieldColor`, `fieldAccentColor`, and `clearingColor` from `internal/game/ingamemenu.go`; restoring those aliases kept overlay styling intact.

## Decision Log

- Decision: Add map data in `internal/game/map.go` rather than creating a new map package.
  Rationale: The current need is local to the first game scene. A separate package would imply broader map ownership before exploration, generation, or multi-map behavior exists.
  Date/Author: 2026-05-13 / Codex

- Decision: Model the first Plot as a fixed 15x15 array.
  Rationale: The user revised the Plot size from 11x11 to 15x15. A fixed array keeps the default home Plot deterministic and easy to test while preserving a natural center Tile.
  Date/Author: 2026-05-13 / Codex

- Decision: Keep the home Plot empty except for road and Sanctum.
  Rationale: The milestone is meant to prove map state and static scene rendering without smuggling in terrain, resources, build rules, or generation.
  Date/Author: 2026-05-13 / Codex

- Decision: Let logical updates continue to advance only the existing debug counter while leaving map state unchanged.
  Rationale: The user explicitly wanted time to pass without logical ticks updating the static scene.
  Date/Author: 2026-05-13 / Codex

## Outcomes & Retrospective

Implementation completed the static home Plot scene. Starting a new game now creates a prototype `Map` with a default home `Plot`, stores it on `game.State`, and renders the scene from that stored map data. The home Plot is 15x15, empty except for the centered Sanctum and a straight north road to the Plot edge. Logical ticks still advance the existing debug counter while leaving the map unchanged.

Validation results:

    go test ./...
    ok  	td/cmd/td	0.018s
    ok  	td/internal/game	0.018s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.747s

    file plans/09-static-prototype-scene/screenshots/*.png
    plans/09-static-prototype-scene/screenshots/ingame-menu.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/09-static-prototype-scene/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/09-static-prototype-scene/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/09-static-prototype-scene/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/09-static-prototype-scene/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

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
    84 internal/game/scene.go
    100 internal/game/hud.go
    129 internal/menu/start.go
    151 cmd/td/main.go
    159 internal/game/ingamemenu.go
    174 internal/game/game.go
    181 cmd/td/main_test.go
    282 internal/menu/menu_test.go
    314 internal/game/game_test.go
    343 internal/menu/menu.go
    2059 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a local Go/Ebitengine desktop game prototype. `cmd/td/main.go` owns the Ebitengine window, app-mode routing, and input polling. `internal/menu/` owns the main menu, New Game configuration screen, and Settings placeholder. `internal/game/` owns the active game state, pause behavior, top-bar HUD, in-game overlay menu, and now the prototype map and static home Plot scene. `internal/ui/` owns shared palette values.

The key game-design term for this plan is `Plot`: a 15x15 group of map `Tile`s. The `Sanctum` is the wizard's central tower. For this milestone, the home Plot is not a full gameplay map; it is a deterministic static scene with a centered Sanctum and a straight road north to the Plot edge.

The root control documents constrain the work. `PRODUCT.md` and `README.md` must describe the new visible scene as current behavior. `ROADMAP.md` must move past the static-scene milestone toward basic scene interaction. `GAME.md` must record the north road and empty home Plot decisions. `DESIGN.md` must preserve readability requirements for the static map. `ARCHITECTURE.md` must say `internal/game` owns the prototype map data and static scene rendering. `CODESTYLE.md` requires `gofmt`, doc comments for Go functions and methods, tests for pure behavior, and a final hand-written code-file line-count review against the 600-line preference.

## Plan of Work

Add `internal/game/map.go`. Define a `Map` struct that owns the current prototype home `Plot`. Define `Plot` as a 15x15 fixed grid of `Tile` values. Define each `Tile` with private terrain and feature values. For now, the only terrain values are empty ground and road, and the only feature values are no feature and Sanctum. Add `NewDefaultMap` and `NewDefaultHomePlot`; the default home Plot should put the Sanctum at Tile `(7, 7)`, mark road terrain from `(7, 7)` straight north through `(7, 0)`, and leave every other Tile empty with no feature.

Update `internal/game/game.go` so `State` stores the new `Map`. `game.New` should initialize that field with `NewDefaultMap()` before any rendering or updates happen. Remove the freehand placeholder field renderer from `game.go`.

Add `internal/game/scene.go` for static scene rendering. The renderer should calculate a centered square board below the top HUD, draw each Tile from `State.gameMap.Home`, and visually distinguish empty ground, road, and Sanctum with simple shapes. It should not introduce camera movement, zoom, animation, resource hints, terrain variety, pathfinding, or gameplay rules.

Update `internal/game/game_test.go` with pure tests for the default home Plot and for the invariant that logical updates do not mutate map state. Update `cmd/td/main_test.go` so screenshot evidence is saved under `plans/09-static-prototype-scene/screenshots/`.

Update the control documents named above. Do not add dependencies. Do not introduce an asset pipeline, save system, campaign system, scene framework, exploration rules, base-building rules, or combat rules.

## Concrete Steps

From the repository root, inspect the working tree:

    git status --short

Edit `internal/game/map.go`, `internal/game/game.go`, `internal/game/scene.go`, `internal/game/colors.go`, `internal/game/game_test.go`, and `cmd/td/main_test.go`. Format and test the game package:

    gofmt -w cmd/td/main_test.go internal/game/*.go
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

The implementation is accepted when `go test ./...` succeeds, screenshot capture writes PNG evidence under `plans/09-static-prototype-scene/screenshots/`, `go run ./cmd/td` opens without startup errors, `git diff --check` reports no whitespace errors, and the line-count review finds no hand-written Go file over the 600-line preference.

A human should be able to start a new game and see a 15x15 static home Plot below the HUD. The center Tile should show the Sanctum, and the road should run straight north from the Sanctum to the top edge of the Plot. All other Tiles should be empty ground. SPACE should still pause the logical update counter, ESC should still open the in-game overlay over the scene, and Surrender should still return to the main menu.

Documentation is accepted when `PRODUCT.md` and `README.md` describe the static home Plot as current behavior, `ROADMAP.md` lists basic scene interaction as the next priority, `GAME.md` records the empty home Plot and north road decisions, `DESIGN.md` covers static map readability, and `ARCHITECTURE.md` records `internal/game` ownership of map state and static scene rendering.

## Idempotence and Recovery

The changes are additive and local to game state, rendering, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when graphics are available.

If the Plot appears too large or overlaps HUD/debug text at 1920x1080, adjust only the static scene layout constants in `internal/game/scene.go` before adding any new layout abstraction. If the map tests fail, fix `NewDefaultHomePlot` rather than special-casing the renderer.

## Artifacts and Notes

Important artifacts:

    plans/09-static-prototype-scene.md
    plans/09-static-prototype-scene/screenshots/running-game.png
    plans/09-static-prototype-scene/screenshots/paused-game.png
    plans/09-static-prototype-scene/screenshots/ingame-menu.png
    internal/game/map.go
    internal/game/scene.go
    internal/game/game.go
    internal/game/game_test.go
    cmd/td/main_test.go
    README.md
    PRODUCT.md
    ROADMAP.md
    GAME.md
    DESIGN.md
    ARCHITECTURE.md

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add dependencies.

The new map-facing interfaces are in `internal/game/map.go`:

    type Map struct { Home Plot }
    type Plot struct { Tiles [plotSize][plotSize]Tile }
    type Tile struct { Terrain tileTerrain; Feature tileFeature }
    func NewDefaultMap() Map
    func NewDefaultHomePlot() Plot

`game.State` stores this value in a private `gameMap Map` field. `State.Draw` calls the static scene renderer before drawing the top bar, Wizard name, debug counter, and in-game overlay.

## Revision Note

This plan was updated during implementation to record the completed validation evidence, the restored overlay color aliases discovered by testing, and the final line-count review required by `PLANS.md` and `CODESTYLE.md`. It was revised again when the prototype Plot size changed from 11x11 to 15x15, so future readers see the current implemented size rather than the earlier design.
