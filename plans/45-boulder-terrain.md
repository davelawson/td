# Add Boulder Terrain

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/45-boulder-terrain.md`.

## Purpose / Big Picture

After this change, explored grasslands Plots can contain sparse Boulder terrain in addition to sparse Forest. Boulder gives generated Plots another readable obstacle type without adding resource gathering or new production art. Players should see occasional stone-like terrain in newly explored grasslands, and building placement should reject those Tiles just like roads and forests.

## Progress

- [x] (2026-07-13T22:02Z) Created this ExecPlan from the accepted plan.
- [x] (2026-07-13T22:09Z) Added Boulder terrain data, generation, and rendering.
- [x] (2026-07-13T22:13Z) Added tests for Boulder generation, precedence, home-Plot exclusion, road/edge overrides, and placement rejection.
- [x] (2026-07-13T22:16Z) Updated root control documents to describe Boulder terrain.
- [x] (2026-07-13T22:20Z) Ran `go test ./...`, screenshot capture, `git diff --check`, ownership check, and Go file line-count review.

## Surprises & Discoveries

- Observation: This plan builds on uncommitted plan 44 changes.
  Evidence: `git status --short` shows `internal/game/plot_generator.go`, `plans/44-grasslands-biome-generation.md`, and related docs/tests still pending.
- Observation: The line-count review still finds one test file over the 600-line preference and one close to it.
  Evidence: final review reported `602 internal/game/game_test.go` and `596 internal/game/building_bar_test.go`.

## Decision Log

- Decision: Implement Boulder as terrain, not a feature or resource node.
  Rationale: The requested slice is a terrain type; resource behavior would require a separate gathering design.
  Date/Author: 2026-07-13 / Codex
- Decision: Generate sparse Boulders in grasslands Plots.
  Rationale: The user chose the recommended scope, making Boulder visible immediately through exploration.
  Date/Author: 2026-07-13 / Codex
- Decision: Let Boulder take precedence over Forest when a Tile's tweak matches both generation rules.
  Rationale: A deterministic precedence rule keeps generation testable and prevents ambiguous terrain assignment.
  Date/Author: 2026-07-13 / Codex

## Outcomes & Retrospective

Implementation completed Boulder terrain as a generated grasslands obstacle. `terrainBoulder` is now a map terrain value, grasslands generation creates sparse Boulders before sparse Forest so overlap is deterministic, and the starting home Plot remains open grassland without generated obstacles. Boulder renders with a simple stone-colored vector marker and remains non-buildable because placement only accepts `terrainEmpty`.

Validation passed on 2026-07-13: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded. Screenshot evidence was written under `plans/45-boulder-terrain/screenshots/`.

The final hand-written Go file line-count review found `internal/game/game_test.go` over the 600-line preference at 602 lines and `internal/game/building_bar_test.go` close to the preference at 596 lines. Recommended follow-up remains a separate approved test-organization change that splits map/default-state tests out of `game_test.go` and considers splitting building-bar rendering tests by category or behavior before adding more cases. No unplanned split was performed.

## Context and Orientation

The current working tree already contains the grasslands biome generation slice from `plans/44-grasslands-biome-generation.md`. `internal/game/map.go` defines `tileTerrain` values `terrainEmpty`, `terrainRoad`, and `terrainForest`. `internal/game/plot_generator.go` creates home and explored grasslands Plots and currently generates sparse Forest based on each Tile's random `Tweak`. `internal/game/scene.go` renders Forest with a pine sprite, and `internal/game/building_bar.go` permits placement only on `terrainEmpty` Tiles with no feature, so any new terrain value is non-buildable by default.

`GAME.md` currently lists intended terrain types and does not yet include Boulder. `PRODUCT.md`, `ARCHITECTURE.md`, `DESIGN.md`, and `ROADMAP.md` describe grasslands generation with sparse Forest only. These documents must be updated because Boulder changes current rendered map output and current placement limitations.

## Plan of Work

Add `terrainBoulder` to the terrain enum. In `internal/game/plot_generator.go`, add a sparse Boulder rule using a fixed tweak modulo. During grasslands generation, assign Boulder before Forest and only assign Forest when the tile is not already Boulder. Keep the home Plot using open grasslands so it never generates Boulder.

Update rendering by adding a Boulder color and drawing a simple vector rock shape for Boulder Tiles. Do not add assets. Keep the grid, roads, structures, and explore controls readable.

Add tests for generation and placement. Generator tests should prove grasslands can contain Boulder, Boulder wins over Forest on overlapping tweak values, home generation excludes Boulder, and map reveal road/edge cleanup overwrites generated Boulder. Placement tests should prove a Boulder tile rejects a dragged building without spending resources.

Update `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, `DESIGN.md`, and `ROADMAP.md` to describe Boulder as sparse generated grasslands terrain, non-buildable, and not yet a resource node.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go`, `internal/game/plot_generator.go`, `internal/game/colors.go`, and `internal/game/scene.go`.
2. Update tests in `internal/game/plot_generator_test.go`, `internal/game/exploration_test.go`, and `internal/game/build_placement_test.go`.
3. Update root control documents.
4. Validate:

    gofmt -w internal/game/map.go internal/game/plot_generator.go internal/game/colors.go internal/game/scene.go internal/game/*_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

If any hand-written Go file exceeds 600 lines, record it in `Outcomes & Retrospective` with a concrete recommendation. Do not do an unplanned split unless the user approves it.

## Validation and Acceptance

Acceptance is met when newly explored grasslands Plots can contain sparse Boulder terrain, Boulder renders distinctly from Forest and Road, and building placement rejects Boulder without resource or population changes. The home Plot must remain free of generated Boulder. North-chain roads and cleared shared edges must overwrite Boulder terrain. `go test ./...`, screenshot capture, `git diff --check`, ownership, and line-count checks must pass.

## Idempotence and Recovery

The change is additive. Re-running formatting, tests, and screenshot capture is safe. If generated Boulder appears on a road or cleared shared edge, preserve the existing override order rather than special-casing rendering or placement.

## Artifacts and Notes

Screenshot evidence should be written under `plans/44-grasslands-biome-generation/screenshots/` unless the screenshot harness is moved again. The active plan should record the final validation results and line-count review.

## Interfaces and Dependencies

No new external dependencies are required. Boulder uses existing map storage, existing vector rendering, and the existing placement rule that only `terrainEmpty` is buildable.
