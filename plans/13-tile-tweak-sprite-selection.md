# Add Tile Tweak Sprite Selection

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/13-tile-tweak-sprite-selection.md` and follows `PLANS.md` from the repository root.

## Purpose / Big Picture

Tile sprite variation should belong to map data instead of being derived from screen or grid position. After this change, every new home Plot Tile stores a random unsigned 16-bit `Tweak` value. Forest Tile rendering uses that value to choose the pine sprite variant, and the highest bit flips tree sprites horizontally. A contributor can see the result by starting a new game or by capturing screenshots under this plan directory.

## Progress

- [x] (2026-05-13T05:13Z) Inspected the current Tile model, home Plot creation, scene rendering, tests, and plan order.
- [x] (2026-05-13T05:13Z) Added `Tile.Tweak uint16` and assigned random tweaks during default home Plot creation.
- [x] (2026-05-13T05:13Z) Updated forest rendering to choose pine variants from tweak low bits and flip from the high bit.
- [x] (2026-05-13T05:13Z) Added tests for deterministic tweak assignment, sprite variant selection, and flip-bit behavior.
- [x] (2026-05-13T05:13Z) Updated root control documents and screenshot capture target for the new map-data behavior.
- [x] (2026-05-13T05:14Z) Ran `gofmt`, `go test ./...`, `git diff --check`, screenshot capture, and the final hand-written code-file line-count review.

## Surprises & Discoveries

- Observation: Forest sprite choice was previously deterministic but derived from Plot coordinates in `internal/game/scene.go`.
  Evidence: `drawPineTree` selected `trees[(plotX*3+plotY*5)%len(trees)]` before this plan.

## Decision Log

- Decision: Store the tweak as exported `Tile.Tweak uint16`.
  Rationale: `Tile` already exposes `Terrain` and `Feature`, and the user explicitly requested a field named `tweak`; Go style for exported struct fields uses `Tweak`.
  Date/Author: 2026-05-13 / Codex

- Decision: Use production randomness only in `NewDefaultHomePlot`, with a private deterministic helper for tests.
  Rationale: The game should assign random tweak values, while tests should verify assignment without relying on exact random output or probabilistic assertions.
  Date/Author: 2026-05-13 / Codex

- Decision: Use bit `15`, mask `0x8000`, to flip tree sprites horizontally, and use bits `0..14` for variant selection.
  Rationale: The user specified the highest order bit for tree flipping; ignoring that bit for variant selection keeps flip state from changing the selected art variant.
  Date/Author: 2026-05-13 / Codex

## Outcomes & Retrospective

The map model now stores `Tweak uint16` on every Tile. The default home Plot assigns a random tweak to every Tile in production, while tests use a deterministic private tweak source to verify construction without depending on random output. Forest rendering now chooses pine sprite variants from the lower 15 tweak bits and flips tree sprites horizontally when the high bit is set.

Validation passed on 2026-05-13: `go test ./...`, `git diff --check`, and `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` all succeeded. Screenshot evidence was written under `plans/13-tile-tweak-sprite-selection/screenshots/`.

The final hand-written Go file line-count review reported no files over the 600-line preference. The largest file was `internal/game/game_test.go` at 535 lines, so no split or exception is needed.

## Context and Orientation

`td` is a local Go/Ebitengine tower-defense prototype. `internal/game/map.go` defines `Map`, `Plot`, and `Tile`, and creates the default home Plot. `internal/game/scene.go` renders the Plot. `assets/catalog.go` loads the current Sanctum and pine tree sprites. The current home Plot is a 15x15 grid with a centered Sanctum, a north road, and a pine-tree border except at the road opening.

`Tile` is the data record for one Plot square. `Tweak` is a prototype visual-variation value stored on each Tile. In this plan it affects only forest Tile rendering because empty and road Tiles are still colored rectangles, and the Sanctum is a feature sprite drawn over its Tile.

Root control documents constrain this work. `CODESTYLE.md` requires doc comments for Go functions and a final line-count review. `ARCHITECTURE.md` should mention that prototype map Tiles store visual variation data used by rendering. `GAME.md` should record that Tile tweak is map-backed visual variation, not a gameplay rule. `PRODUCT.md` and `README.md` do not need user-facing wording unless the visible behavior changes beyond variation in the existing tree border.

## Plan of Work

Update `internal/game/map.go`. Add `Tweak uint16` to `Tile`. Add a private `randomTileTweak` helper that returns `uint16(rand.Intn(1 << 16))`. Create a private home Plot helper that accepts a tweak source function, instantiates all 225 Tiles through a `newTile` helper, then applies the existing terrain and feature overrides. Keep `NewDefaultHomePlot` as the production entrypoint and have it pass `randomTileTweak`.

Update `internal/game/scene.go`. Replace coordinate-based tree selection with helpers that decode `Tile.Tweak`. `pineTreeSpriteIndex(tweak, variants)` should use `int(tweak & 0x7fff) % variants`. `treeSpriteFlipped(tweak)` should return true when `tweak & 0x8000 != 0`. When flipped, draw the tree with a negative X scale and translate it so it remains centered in the Tile.

Update tests in `internal/game/game_test.go`. Use the private deterministic plot helper to assert tweaks are assigned to every Tile from the supplied source. Add pure tests for pine sprite index wrapping, high-bit independence, and flip detection. Keep existing map invariants unchanged.

Update `cmd/td/main_test.go` screenshot capture output to `plans/13-tile-tweak-sprite-selection/screenshots/`. Update `GAME.md` and `ARCHITECTURE.md` with one concise note each about `Tile.Tweak`.

## Concrete Steps

From the repository root:

    pwd
    # /home/dave/dev/ai/td

Edit these files:

    internal/game/map.go
    internal/game/scene.go
    internal/game/game_test.go
    cmd/td/main_test.go
    GAME.md
    ARCHITECTURE.md

After editing, format and validate:

    gofmt -w cmd/td/main_test.go internal/game/game_test.go internal/game/map.go internal/game/scene.go
    go test ./...
    git diff --check

If the local environment supports Ebitengine screenshot capture, run:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expect screenshots in `plans/13-tile-tweak-sprite-selection/screenshots/`, especially `running-game.png`, showing the existing tree border with sprite variation controlled by stored Tile tweaks.

End with a line-count review for hand-written Go files:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

If any hand-written code file exceeds 600 lines, record the file and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user before doing extra refactor work unless it was already included in this plan.

## Validation and Acceptance

`go test ./...` must pass. Tests must prove every Tile is instantiated with a tweak from the construction source, tree sprite selection wraps through available variants using low bits, the high bit does not alter variant choice, and the high bit controls horizontal flipping.

Visual acceptance is that the home Plot still shows the centered Sanctum, north road, and tree border, with tree sprite variants now selected from Tile data. Existing camera zoom, `WASD` panning, pause behavior, overlay menu behavior, and surrender flow must keep working.

Documentation acceptance is that `GAME.md` and `ARCHITECTURE.md` describe the tweak field at the right level without turning it into a gameplay mechanic.

## Idempotence and Recovery

The changes are local and safe to rerun. Re-running formatting, tests, screenshot capture, and line-count checks is safe. If a test fails because an exact random value was assumed, rewrite the test to use the private deterministic tweak source instead of testing production random output. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the failure and rely on automated tests plus manual launch validation when graphics are available.

## Artifacts and Notes

Expected important artifacts:

    plans/13-tile-tweak-sprite-selection/screenshots/running-game.png

The most important code evidence is in the pure tests for tweak assignment and tree sprite decoding.

## Interfaces and Dependencies

In `internal/game/map.go`, `Tile` must include:

    Tweak uint16

In `internal/game/scene.go`, tree rendering must interpret `Tile.Tweak` as follows:

    bits 0..14: pine tree sprite variant source
    bit 15: horizontal flip when set

No new external dependencies are required. Use the Go standard library for random number generation.
