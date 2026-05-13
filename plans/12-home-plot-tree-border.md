# Add a Tree Border to the Home Plot

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/12-home-plot-tree-border.md` and follows `PLANS.md` from the repository root.

## Purpose / Big Picture

Starting a new game should show a home Plot that feels more like a small defended clearing instead of an empty board. After this change, the default 15x15 home Plot keeps the centered Sanctum and north road, and every edge Tile is bordered by pine trees except the road opening at the north edge. A contributor can see this by starting the game or by capturing the running game screenshot under this plan directory.

## Progress

- [x] (2026-05-13T04:28Z) Inspected the current map model, renderer, asset catalog, tests, and control documents.
- [x] (2026-05-13T04:28Z) Added forest terrain to the default home Plot edge while preserving the north road opening.
- [x] (2026-05-13T04:28Z) Added existing pine tree sprites to the typed runtime asset catalog.
- [x] (2026-05-13T04:28Z) Updated home Plot rendering to draw deterministic pine sprites for forest Tiles.
- [x] (2026-05-13T04:28Z) Updated pure tests for the tree border and catalog asset loading.
- [x] (2026-05-13T04:30Z) Updated root control documents and screenshot capture target for the new current behavior.
- [x] (2026-05-13T04:30Z) Ran `gofmt`, `go test ./...`, `git diff --check`, screenshot capture, and the final hand-written code-file line-count review.

## Surprises & Discoveries

- Observation: The repository already had four untracked 64x64 pine tree PNGs under `assets/sprites/terrains/`.
  Evidence: `file assets/sprites/terrains/*.png` reported four 64 x 64 RGBA PNG files.

## Decision Log

- Decision: Model the tree border as `terrainForest` in `internal/game/map.go`, not as render-only decoration.
  Rationale: The map state should describe what exists on the Plot so tests and later gameplay rules can reason about edge trees consistently.
  Date/Author: 2026-05-13 / Codex

- Decision: Keep the existing north road as authoritative over the edge tree border.
  Rationale: The user explicitly requested trees around the edge except where the road exists, and the current road exits through the north-center Tile `(7,0)`.
  Date/Author: 2026-05-13 / Codex

- Decision: Use the existing terrain PNGs and a deterministic coordinate-based sprite choice.
  Rationale: The assets are already in the intended `assets/sprites/terrains/` location, and deterministic variation improves the visual border without adding random map generation.
  Date/Author: 2026-05-13 / Codex

## Outcomes & Retrospective

The home Plot now has a data-backed forest terrain border around the edge, with the existing north road preserving the road opening through the border. The asset catalog embeds and loads four pine tree terrain sprites, and the renderer draws a deterministic pine variant for each forest Tile. Root documentation now describes the tree-bordered home Plot as current behavior.

Validation passed on 2026-05-13: `go test ./...`, `git diff --check`, and `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` all succeeded. Screenshot evidence was written under `plans/12-home-plot-tree-border/screenshots/`.

The final hand-written Go file line-count review reported no files over the 600-line preference. The largest file was `internal/game/game_test.go` at 485 lines, so no split or exception is needed.

## Context and Orientation

`td` is a local Go/Ebitengine tower-defense prototype. `cmd/td/main.go` owns the desktop app and routes input to `internal/game`. `internal/game/map.go` creates the default home Plot, `internal/game/scene.go` renders that Plot, and `assets/catalog.go` embeds runtime sprites into a typed catalog used by game state.

The current home Plot is a 15x15 grid. Tile coordinates are represented by array indexes where `plot.Tiles[y][x]` addresses Tile `(x,y)`. `homePlotCenter` is `7`, so the Sanctum is at `(7,7)`. The existing road runs north from the Sanctum through `(7,0)`. The tree border must occupy all edge Tiles where `x == 0`, `y == 0`, `x == 14`, or `y == 14`, except the existing road Tiles.

Root control documents constrain this work. `PRODUCT.md` and `README.md` must describe the new visible current behavior. `GAME.md` must record that the first rendered home Plot now has a pine/tree border while still avoiding broader terrain systems. `ARCHITECTURE.md` must continue to say `internal/game` owns prototype map data and rendering, and `assets/` owns runtime sprite loading. `DESIGN.md` requires roads, Tile boundaries, and the Sanctum to remain readable.

## Plan of Work

First, update `internal/game/map.go` by adding a new `tileTerrain` value named `terrainForest`. In `NewDefaultHomePlot`, mark all perimeter Tiles as forest, then apply the existing north road loop after that so road Tiles overwrite forest where the road exists. Leave the Sanctum feature at the center.

Next, update `assets/catalog.go` so the existing `assets/sprites/terrains/pine-tree-1.png` through `pine-tree-4.png` files are embedded and loaded into a typed terrain sprite group. Keep asset loading centralized in `assets`; game rendering should not decode files directly.

Then update `internal/game/scene.go` so forest Tiles fill with a forest base color and draw one pine tree sprite. Choose the pine sprite deterministically from the Tile coordinates so the same map always renders the same border. Keep the Sanctum and road rendering unchanged except for layering over the new terrain type.

Update tests in `internal/game/game_test.go` for the new default home Plot invariants: edge Tiles are forest, the north road remains continuous, the north road edge Tile is not forest, and non-edge interior Tiles remain empty except for the road and Sanctum. Update `assets/catalog_test.go` to verify the four pine sprites load as 64x64 images.

Update `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to describe the home Plot tree border as current behavior. Update the screenshot capture target in `cmd/td/main_test.go` to write under `plans/12-home-plot-tree-border/screenshots/`.

## Concrete Steps

From the repository root, inspect and edit:

    pwd
    # /home/dave/dev/ai/td

    sed -n '1,120p' internal/game/map.go
    sed -n '1,140p' internal/game/scene.go
    sed -n '1,160p' assets/catalog.go

After editing, format and validate:

    gofmt -w assets/catalog.go assets/catalog_test.go cmd/td/main_test.go internal/game/colors.go internal/game/game_test.go internal/game/map.go internal/game/scene.go
    go test ./...
    git diff --check

If the local environment supports Ebitengine screenshot capture, run:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expect screenshots in `plans/12-home-plot-tree-border/screenshots/`, especially `running-game.png`, showing a 15x15 Plot with pines around the edge and the road opening at the north-center edge.

End with a line-count review for hand-written Go files:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

If any hand-written code file exceeds 600 lines, record the file and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user before doing extra refactor work unless it was already included in this plan.

## Validation and Acceptance

`go test ./...` must pass. The new asset test must prove all four pine tree sprites load as 64x64 images. The new map tests must prove the edge forest border exists, the north road still reaches the Plot edge, and interior non-road Tiles remain empty.

Visual acceptance is that starting a new game shows pine trees around the Plot edge except at the north road opening. The road, Tile grid, and sprite-backed Sanctum must remain readable. Existing camera zoom, `WASD` panning, pause behavior, overlay menu behavior, and surrender flow must keep working.

Documentation acceptance is that `README.md` and `PRODUCT.md` describe the tree-bordered home Plot as current behavior, `GAME.md` records the starting Plot decision, and `ARCHITECTURE.md` reflects that the prototype map now includes forest edge terrain and terrain sprites.

## Idempotence and Recovery

The map generation change is deterministic and safe to rerun. Re-running `gofmt`, tests, screenshot capture, and line-count checks is safe. If asset loading fails, verify the four pine PNGs exist under `assets/sprites/terrains/` and that `assets/catalog.go` embeds exactly those paths. If visual capture fails because the environment cannot open an Ebitengine window, record the failure and rely on `go test ./...` plus manual `go run ./cmd/td` validation.

## Artifacts and Notes

Expected important artifacts:

    plans/12-home-plot-tree-border/screenshots/running-game.png

Key behavior evidence should come from the map tests and asset catalog tests added in this change.

## Interfaces and Dependencies

In `internal/game/map.go`, `tileTerrain` must include:

    terrainEmpty
    terrainRoad
    terrainForest

In `assets/catalog.go`, define a terrain sprite group reachable from the existing catalog, with four pine tree images loaded from `assets/sprites/terrains/pine-tree-1.png` through `pine-tree-4.png`.

No new external dependencies are required. Continue using Ebitengine for rendering and the existing embedded asset catalog for runtime image loading.
