# Enemy World Coordinates

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

Save this file at `plans/18-enemy-world-coordinates.md` in the repository root. The previous highest numbered plan is `plans/17-skeleton-raid-sprites.md`.

## Purpose / Big Picture

After this change, Raid enemies store their position as explicit world coordinates instead of as progress along the current path. The center of the Sanctum Tile is the world origin `(0, 0)`. One Tile is one world unit, floating point positions are allowed, and positive Y points north. The first deterministic Raid should look and play the same: skeletons enter from the north-center road, move south toward the Sanctum, spend Barricade charges on contact, and breach the Sanctum when no Barricade remains.

This is a coordinate-model migration only. It does not add pathfinding, multiple paths, collision, tower targeting, damage, new terrain rules, or a broader scene framework.

## Progress

- [x] (2026-05-16T00:00Z) Inspected current Raid enemy progress storage, rendering projection, map Tile layout, camera tests, screenshot capture, and control documents.
- [x] (2026-05-16T00:06Z) Updated Raid enemy state to store tile-space world coordinates.
- [x] (2026-05-16T00:10Z) Updated camera projection and rendering to use Sanctum-centered tile units.
- [x] (2026-05-16T00:14Z) Updated tests and screenshot capture for the new coordinate model.
- [x] (2026-05-16T00:18Z) Updated durable control documents with the coordinate decision.
- [x] (2026-05-16T00:22Z) Ran validation commands, captured visual evidence, and checked hand-written code line counts.

## Surprises & Discoveries

- Observation: Existing camera and rendering code uses pixel-space world coordinates with screen-style positive Y downward.
  Evidence: `internal/game/camera.go` initializes the camera to the pixel center of the 15x15 Plot, and `internal/game/scene.go` draws Tiles at `x * plotBaseTileSize`, `y * plotBaseTileSize`.

- Observation: `internal/game/game_test.go` is close to the 600-line code-style preference before this change.
  Evidence: `wc -l internal/game/game_test.go` reported 573 lines before implementation.

## Decision Log

- Decision: Use positive Y as north.
  Rationale: The user specified that `(0, 1.5)` is north of the Sanctum.
  Date/Author: 2026-05-16 / User and Codex

- Decision: Preserve existing Raid pacing by converting the previous pixel speed to tile units.
  Rationale: This keeps current user-visible behavior stable while changing the storage model.
  Date/Author: 2026-05-16 / Codex

- Decision: Keep the current single straight north-road contact rule.
  Rationale: The request concerns enemy position storage, not pathfinding or route definition.
  Date/Author: 2026-05-16 / Codex

## Outcomes & Retrospective

Completed. Active Raid enemies now store explicit Sanctum-centered world positions in Tile units. The current north-road Raid remains behaviorally equivalent: skeletons spawn from the north-center road, move south toward the Sanctum, spend Barricade charges on contact, and breach the Sanctum when no Barricade remains.

Validation passed on 2026-05-16: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, and `git diff --check` all succeeded. Screenshot evidence was written under `plans/18-enemy-world-coordinates/screenshots/`, including `active-raid.png`, which shows a skeleton on the north road during an active Raid after the coordinate migration.

The hand-written Go line-count review found no files over 600 lines. `internal/game/game_test.go` is the largest file at 572 lines, so it remains under the `CODESTYLE.md` preference.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. The playable game state lives in `internal/game/`. The current home Plot is a fixed 15x15 grid with the Sanctum at the center Tile and a straight road north to the Plot edge.

Raid simulation lives in `internal/game/raid.go`. Before this plan, each `raidEnemy` stored `progress float64`, movement added to progress, contact compared progress to the path length, and rendering converted progress into pixel-space coordinates on the road.

Camera and projection code lives in `internal/game/camera.go`. Before this plan, camera center and projected rectangles used pixel-space coordinates with positive Y downward.

Rendering lives in `internal/game/scene.go` and `internal/game/raidui.go`. These files should render from tile-space world coordinates after this change while preserving the current visible layout.

`GAME.md`, `PRODUCT.md`, and `ARCHITECTURE.md` are durable control documents. This change records a real gameplay-coordinate decision and changes architecture wording around Raid state and projection, so those files must be updated.

## Plan of Work

Add a private `worldPosition` type in `internal/game/raid.go` or another nearby game file. Change `raidEnemy` to store `position worldPosition` instead of `progress float64`.

Define Raid positions in tile-space units. Spawn skeletons at `worldPosition{X: 0, Y: float64(homePlotCenter)}`. Convert the existing movement speed to tile units as `raidEnemySpeedTiles = raidEnemySpeed / plotBaseTileSize` or equivalent, then move enemies south by subtracting from `position.Y`. Treat contact with the Sanctum as `position.Y <= 0` for the current straight north road.

Update camera and projection to use tile-space coordinates with positive Y north. The initial camera center should be `(0, 0)`. Panning should preserve visual feel by dividing pixel pan speed by `plotBaseTileSize * zoom`; `PanUp` increases camera center Y, `PanDown` decreases it, `PanLeft` decreases X, and `PanRight` increases X. Projection should convert a world rectangle described by west edge, north edge, width, and height in tile units into a screen rectangle using `plotBaseTileSize * camera.zoom`.

Update map rendering so each Tile is drawn from Sanctum-centered grid coordinates. Tile `(homePlotCenter, homePlotCenter)` should draw centered on `(0, 0)`, the first north road Tile should draw centered on `(0, 1)`, and the common edge between the first and second north road Tiles should be at `Y == 1.5`.

Update Raid rendering to draw sprites and nil-sprite fallback markers from `enemy.position`. Keep the existing visual sprite and marker sizes by expressing them as `raidEnemySpriteSize / plotBaseTileSize` and `raidEnemyRadius / plotBaseTileSize` in world units.

Update tests. Existing Raid tests should assert position changes rather than progress changes. Add focused pure tests for grid-to-world Tile centers and the `Y == 1.5` north-road edge example. Update camera tests for the new world origin, panning signs, and tile-unit pan distance. Avoid large additions to `internal/game/game_test.go`; if needed, create a dedicated `internal/game/coordinates_test.go`.

Update `cmd/td/main_test.go` so screenshot evidence writes under `plans/18-enemy-world-coordinates/screenshots/`.

Update `GAME.md`, `PRODUCT.md`, and `ARCHITECTURE.md` to describe enemy world coordinates and keep current user-visible Raid behavior clear.

## Concrete Steps

Run these commands from the repository root, `/home/dave/dev/ai/td`.

Format edited Go files:

    gofmt -w internal/game/camera.go internal/game/scene.go internal/game/raid.go internal/game/raidui.go internal/game/raid_test.go internal/game/game_test.go internal/game/coordinates_test.go cmd/td/main_test.go

Run the full test suite:

    go test ./...

Expected result: all packages pass.

Capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected result: the command passes and writes screenshots under `plans/18-enemy-world-coordinates/screenshots/`, including an active-Raid screenshot with skeletons still on the north road.

Check whitespace and working tree state:

    git diff --check
    git status --short

Expected result: no whitespace errors. `git status --short` should list only intentional modified files, the new ExecPlan, and screenshot artifacts.

Check hand-written Go file line counts:

    find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

If any hand-written code file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and do not perform an unplanned split, refactor, or library addition unless the user approves it.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, `git diff --check` reports no whitespace errors, screenshot capture writes evidence under `plans/18-enemy-world-coordinates/screenshots/`, and Raid behavior remains visibly equivalent except for the internal coordinate model.

Behavioral acceptance is that clicking `Next Raid` still starts the deterministic Raid immediately, the first enemy still appears at the north-center road, later enemies still spawn on the same stagger, skeletons move south toward the Sanctum, remaining enemy counts still include pending plus active enemies, and reaching skeletons still spend Barricade charges or breach the Sanctum.

Coordinate acceptance is that enemies store world positions, the Sanctum center is `(0, 0)`, one Tile is one world unit, positive Y points north, and `(0, 1.5)` is the common edge between the first and second road Tiles north of the Sanctum.

## Idempotence and Recovery

The code edits are deterministic and can be reapplied safely. Screenshot capture overwrites only files under `plans/18-enemy-world-coordinates/screenshots/`. If the visible map is shifted or inverted, first inspect `projectRect` and Tile coordinate helpers before changing Raid logic. If Raid lifecycle tests fail, inspect spawn and contact thresholds before changing rendering.

## Artifacts and Notes

Expected screenshot artifacts:

    plans/18-enemy-world-coordinates/screenshots/main-menu.png
    plans/18-enemy-world-coordinates/screenshots/new-game-configuration.png
    plans/18-enemy-world-coordinates/screenshots/running-game.png
    plans/18-enemy-world-coordinates/screenshots/active-raid.png
    plans/18-enemy-world-coordinates/screenshots/paused-game.png
    plans/18-enemy-world-coordinates/screenshots/ingame-menu.png

Validation transcript:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	0.065s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	1.578s

    git diff --check
    # no output

    find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l | sort -n
    # largest file: ./internal/game/game_test.go at 572 lines

## Interfaces and Dependencies

No new third-party dependency is required. Use the existing Ebitengine dependency already present in `go.mod`.

The internal enemy state after this change is:

    type raidEnemy struct {
        template *EnemyTemplate
        position worldPosition
    }

This plan does not add a public coordinate API, path API, Wave API, placement API, targeting API, combat API, resource API, or new enemy archetype.

## Revision Notes

2026-05-16: Created this plan for the enemy world-coordinate migration.
