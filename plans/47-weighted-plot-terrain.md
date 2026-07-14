# Add Weighted Plot Terrain Generation

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/47-weighted-plot-terrain.md`.

## Purpose / Big Picture

After this change, generated grasslands Plots choose Tree, Boulder, or empty grass from explicit percentage weights instead of deriving terrain from each Tile's visual tweak value. The current sparse feel remains: Tree 6%, Boulder 3%, and empty grass 91%. The existing Forest terrain concept is renamed to Tree so the implementation matches the intended tile language.

## Progress

- [x] (2026-07-14T00:00Z) Created this ExecPlan from the accepted plan.
- [x] (2026-07-14T00:07Z) Replaced tweak-modulo terrain generation with weighted Tree/Boulder selection.
- [x] (2026-07-14T00:10Z) Renamed Forest terrain to Tree in code, tests, and control documents.
- [x] (2026-07-14T00:12Z) Ran `gofmt -w internal/game/*.go` and `go test ./...`.
- [x] (2026-07-14T00:14Z) Ran `git diff --check`, ownership check, and Go file line-count review.

## Surprises & Discoveries

- Observation: The working tree was clean before implementation.
  Evidence: `git status --short` printed no files.

## Decision Log

- Decision: Keep sparse grasslands density at Tree 6%, Boulder 3%, and empty grass 91%.
  Rationale: The user chose to keep the current generation density while simplifying the implementation.
  Date/Author: 2026-07-14 / Codex
- Decision: Keep `Tile.Tweak` for terrain sprite variation only.
  Rationale: Terrain type selection should be explicit and readable, while stable sprite variant selection still needs per-Tile data.
  Date/Author: 2026-07-14 / Codex
- Decision: Rename Forest terrain to Tree.
  Rationale: The user requested Tree terminology while refactoring terrain generation.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Implementation completed the terrain generation refactor. Grasslands generation now uses explicit percentage weights of Tree 6%, Boulder 3%, and empty grass 91%. A Tile receives only one generated terrain type from a single percentage roll. `Tile.Tweak` remains on map data for stable terrain sprite variant selection and horizontal flipping, but it no longer controls generated terrain.

The Forest terrain enum was renamed to Tree throughout game code, tests, and root control documents. Pine-tree asset names remain unchanged because they describe the concrete art files.

Validation passed on 2026-07-14: `gofmt -w internal/game/*.go`, `go test ./...`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded.

The final hand-written Go file line-count review found `internal/game/game_test.go` over the 600-line preference at 602 lines. It also found `internal/game/building_bar_test.go` close to the preference at 596 lines. Recommended follow-up remains a separate test-organization change that splits default-state/map tests out of `game_test.go` and considers splitting building-bar tests by category or behavior before adding more cases. No unplanned split was performed.

## Context and Orientation

`internal/game/plot_generator.go` currently generates grasslands obstacles by checking modulo rules against each Tile's `Tweak`. `Tile.Tweak` is also used by `internal/game/scene.go` to choose terrain sprite variants and horizontal flipping. The generation refactor should remove tweak-driven terrain choice without removing the tweak data used by rendering.

`internal/game/map.go` defines terrain enum values, including Tree and Boulder. Rendering draws Tree with pine-tree terrain sprites and Boulder with Boulder terrain sprites. Building placement already rejects every non-empty terrain because it only accepts `terrainEmpty` Tiles with no feature.

## Plan of Work

Rename the former Forest terrain enum to `terrainTree` throughout game code and tests. Keep the existing pine-tree asset names because they describe the concrete sprite files.

Replace grasslands modulo constants and helper functions with a `terrainWeights` struct that stores percentage chances for Tree and Boulder. Add a private weighted selection helper that accepts a single roll from 0 through 99 and returns exactly one terrain: Tree first, Boulder second, otherwise empty grass. Use Tree 6 and Boulder 3 as the grasslands weights.

Keep `NewDefaultHomePlot` authored and open, without weighted generated obstacles. Keep `NewGrasslandsPlot` as the public generated Plot constructor. For tests, provide a private helper that accepts both a deterministic tweak source and a deterministic terrain-roll source so terrain selection and visual variation are independently testable.

Update `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md` to describe sparse Tree and Boulder terrain. Update the `GAME.md` decision log to say generated biome terrain uses explicit percentage weights, while `Tile.Tweak` remains visual variation only.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go`, `internal/game/plot_generator.go`, `internal/game/scene.go`, `internal/game/colors.go`, and related tests.
2. Update root control documents that mention generated Tree terrain.
3. Validate:

    gofmt -w internal/game/*.go
    go test ./...
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

If any hand-written Go file exceeds 600 lines, record it in `Outcomes & Retrospective` with a concrete recommendation. Do not do an unplanned split unless the user approves it.

## Validation and Acceptance

Acceptance is met when grasslands Plot generation is driven by explicit percentage weights, generated Plots can contain Tree and Boulder terrain while keeping mostly empty grass, home Plots remain free of generated obstacles, roads and shared edges still overwrite generated terrain, building placement rejects Tree and Boulder, and `Tile.Tweak` remains only visual variation data. `go test ./...`, `git diff --check`, ownership, and line-count checks must pass.

## Idempotence and Recovery

The change is a private refactor plus terminology cleanup. Re-running formatting and tests is safe. If weighted generation causes terrain to appear on roads or shared edges, preserve the existing override order by fixing generation incorporation rather than special-casing rendering.
