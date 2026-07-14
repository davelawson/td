# Add Terrain-Consuming Resource Production

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept current as work proceeds. Maintain this document in accordance with `PLANS.md` from the repository root.

## Purpose / Big Picture

Economic buildings currently create resources from nothing after every successful Raid. After this change, production visibly depletes the explored landscape: a Woodcutter needs one Tree, a Stone Quarry needs one Boulder, and an Iron Mine needs one Iron Deposit for each Labour payout. A structure that cannot find its matching terrain produces nothing and tries again after the next successful Raid. A player can observe both the resource total increasing and the nearest matching terrain becoming grass.

## Progress

- [x] (2026-07-14 18:10Z) Inspected the phase, map, coordinate, production, selection, documentation, and screenshot-evidence implementations.
- [x] (2026-07-14 18:10Z) Recorded the accepted Domain-wide range, nearest-tile selection, deterministic ordering, and biome-default replacement rules.
- [x] (2026-07-14 18:18Z) Implemented terrain-backed Labour production and selected-terrain cleanup.
- [x] (2026-07-14 18:18Z) Added focused tests for matching, range, distance, deterministic allocation, exhaustion, replacement, and lifecycle behavior; `go test ./...` passes.
- [x] (2026-07-14 18:21Z) Updated player-facing descriptions plus `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, and `ARCHITECTURE.md`.
- [x] (2026-07-14 18:21Z) Captured and visually reviewed focused before-and-after evidence; completed normal, race, whitespace, and ownership validation.
- [x] (2026-07-14 18:21Z) Checked hand-written Go line counts. No file exceeds 600 lines; the existing `internal/game/build_placement_test.go` remains the largest at 564 lines.

## Surprises & Discoveries

- Observation: The current `terrainEmpty` value is the rendered grass terrain for both implemented biomes.
  Evidence: `internal/game/plot_generator.go` leaves the remainder of both terrain weight tables as `terrainEmpty`, and placement/documentation call that Tile grass.

- Observation: Cross-Plot world coordinates already provide integer-centered Tile positions.
  Evidence: `plotTileWorldCenter` in `internal/game/coordinates.go` incorporates each 15-by-15 Plot offset, so squared distance can be compared without a new coordinate model or square roots.

- Observation: Ebitengine permits only one `RunGame` lifecycle in a Go test process, so a regex that runs both focused screenshot tests panics when the second fixture loads images.
  Evidence: The combined capture failed with `ebiten: NewImage cannot be called after RunGame finishes`; the repository's existing evidence workflow likewise runs each focused capture separately.

## Decision Log

- Decision: Let each economic building draw from any matching Tile in any explored Plot.
  Rationale: The user chose whole-Domain production range; unexplored Plots remain unavailable because they have no usable map state.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Process buildings in deterministic explored-Plot, row, then column order and choose the nearest currently available matching Tile by squared world-space distance.
  Rationale: This implements the requested nearest-tile behavior without floating-point square roots. Keeping the first candidate in canonical map order provides a stable tie-break and makes scarce-node allocation reproducible.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Preserve a consumed Tile's feature and `Tweak` while changing only its terrain to the biome default.
  Rationale: Consumption replaces terrain rather than rebuilding the Tile. The visual tweak remains suitable for the resulting grass, and unrelated structure state must not be damaged.
  Date/Author: 2026-07-14 / Codex

- Decision: Clear selection when the selected terrain Tile is consumed.
  Rationale: Grass is not selectable, so retaining a terrain selection would leave invisible stale state after Labour.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Implementation is complete. During post-Raid Labour, Woodcutters consume Trees for Wood, Stone Quarries consume Boulders for Stone, and Iron Mines consume Iron Deposits for Metal. Each producer searches every explored Plot, selects the nearest current match by squared Tile-center distance, and uses canonical map order for producer processing and ties. A producer without a match grants nothing and remains available for later terrain. Successful consumption changes only the terrain to the Plot biome's default grass, preserves feature and tweak data, and clears selection when the consumed natural Tile was selected.

Focused tests cover all three mappings, Domain-wide nearest selection, tie-breaking, scarce terrain, delayed production after a missing input, both biome defaults, Tile-data preservation, selection cleanup, successful-Raid Labour, Management construction timing, and breach behavior. `go test ./...` and `go test -race ./...` pass. `git diff --check` passes, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` prints nothing.

Visual evidence is stored at `plans/57-terrain-consuming-resource-production/screenshots/terrain-production-before.png` and `terrain-production-after.png`. The first frame shows Day 1, 100 Wood, a Woodcutter, and a selected Tree; the second shows Day 2, 110 Wood, grass where the Tree stood, and no stale terrain panel.

The final hand-written Go review found no file above 600 lines. The changed production files are `internal/game/resources.go` at 148 lines, `resource_production_test.go` at 174 lines, and `visual_test.go` at 173 lines. The repository's largest hand-written Go file is the pre-existing `internal/game/build_placement_test.go` at 564 lines. It did not grow in this change; before future placement-test additions, split its housing, economic, and defense scenarios into responsibility-focused files. No unplanned refactor was needed for this feature.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. Its world is stored as explored 15-by-15 Plots in `internal/game/map.go`; each Tile has terrain, an optional structure feature, and a stable visual `Tweak`. Grasslands and hills currently use the same default grass representation, `terrainEmpty`, alongside Tree, Boulder, Iron Deposit, and Road terrain. `internal/game/coordinates.go` converts a Plot coordinate and local Tile coordinate to one continuous Sanctum-centered world position.

After a successful Raid, `beginPostRaidDay` in `internal/game/phase.go` advances the Day, briefly enters Labour, calls `grantEconomicBuildingResources` in `internal/game/resources.go`, and immediately enters Management. Production currently scans every explored Tile and grants the template yield for every Woodcutter, Stone Quarry, and Iron Mine without inspecting terrain. A breached Raid never calls this transition. Day 1 starts in Management, so no producer pays immediately after construction.

Structures retain their existing yields: Woodcutter produces 10 Wood, Stone Quarry produces 10 Stone, and Iron Mine produces 10 Metal. Each already reserves one Peasant when built. This plan changes only the condition for each Labour payout and the corresponding terrain depletion; it does not add placement restrictions, worker reassignment, regrowth, accumulated production, animations, or a Labour results screen.

Selection state is stored in `internal/game/selection.go` as a Tile coordinate. Tree, Boulder, and Iron Deposit terrain is selectable, while grass is not. When production consumes a selected terrain Tile, selection must be cleared as part of the same resolution.

The current behavior is described in `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, and `ARCHITECTURE.md`. Those documents presently say every economic building pays each Labour and that Iron Deposits do not affect Iron Mines. They must be updated together. No source-convention, art-generation, or stable design-language change is needed, so `CODESTYLE.md`, `ART.md`, and `DESIGN.md` remain unchanged.

## Plan of Work

First, revise `internal/game/resources.go` so an economic feature resolves to both its existing yield and its matching required terrain. For each economic building found in canonical map order, scan every currently explored Tile, calculate squared distance between Tile centers using `plotTileWorldCenter`, and retain the closest matching terrain. Because explored Plot coordinates and local Tiles are scanned deterministically, retain the first candidate when distances are equal. If no match exists, skip both depletion and yield for that building. If a match exists, replace only its `Terrain` with the result of a biome-default helper, grant the yield, and clear selection if that exact terrain Tile was selected.

The biome-default helper should explicitly handle grasslands and hills and return `terrainEmpty` for both. Keep a safe `terrainEmpty` fallback for unknown future enum values. The three required mappings are Woodcutter to Tree, Stone Quarry to Boulder, and Iron Mine to Iron Deposit. Do not use random selection or restrict search to the producer's Plot.

Add focused tests in a new resource-production test file rather than growing an already large unrelated test file. Prove all three feature mappings and yields; cross-Plot nearest selection; equal-distance tie-breaking; deterministic producer allocation when terrain is scarce; no payout with no matching terrain; independence of nonmatching terrain; grass replacement for grasslands and hills; preservation of feature and tweak; and selected-terrain cleanup. Update existing phase and Raid tests so their fixtures include the terrain needed for expected payouts, while retaining coverage that Management construction waits and breach skips Labour.

Update economic structure descriptions so tooltips and selection details tell the player which terrain is consumed from the explored Domain. Update the five relevant durable documents so current-state, workflow, roadmap, game-design decisions, and architecture all describe terrain-backed production and depletion consistently. Remove the obsolete limitation that Iron Deposits do not influence Iron Mine behavior.

Add two focused optional screenshot tests in `internal/game/visual_test.go` using the same deterministic state: one before Labour with a selected matching terrain Tile and one after Labour with the producer's resource increased and that Tile changed to grass. Save both under `plans/57-terrain-consuming-resource-production/screenshots/` and inspect the pair for a readable producer, visible depletion, and updated HUD total.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`. Edit the production resolver, tests, descriptions, visual fixture, and control documents described above. Keep this living plan updated after each milestone. Format changed Go files and run:

    gofmt -w internal/game/*.go
    go test ./...
    go test -race ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureTerrainProductionBeforeScreenshot -count=1
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureTerrainProductionAfterScreenshot -count=1
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

The ownership command must print nothing. The final line-count step must report any hand-written Go file over 600 lines, or close enough that the next likely change will cross the preference, with a concrete split recommendation. Do not perform an unplanned refactor without user approval.

## Validation and Acceptance

Automated acceptance requires all package tests and race-enabled tests to pass. A Woodcutter, Stone Quarry, or Iron Mine must pay exactly once during Labour only when one matching Tile is available anywhere in explored map state, and that exact Tile must become the Plot's default grass. Two producers cannot consume the same Tile. Nearest selection and ties must be deterministic. Missing terrain must leave resources, map terrain, staffing, and the producer unchanged so later explored terrain can support a future Labour.

The existing lifecycle remains intact: construction in Management does not pay immediately, successful Raid completion advances the Day and resolves one Labour, and breach resolves none. The screenshot pair must make the before/after terrain change visible and show the corresponding resource increase. Documentation acceptance requires all five updated control documents to agree with these rules and no longer describe Iron Deposits as inspection-only or unrelated to Iron Mine production.

## Idempotence and Recovery

Tests, formatting, and screenshot capture are safe to rerun; screenshot files are overwritten only inside this plan's evidence directory. No migration or external service is involved. If a deterministic test fails, inspect canonical Plot and Tile scan ordering before changing the accepted nearest-distance rule. If screenshot capture cannot open the graphics environment, record the exact failure and retain automated state-transition tests as the behavioral evidence. Preserve unrelated user work if the worktree changes during implementation.

## Artifacts and Notes

Focused visual evidence belongs at `plans/57-terrain-consuming-resource-production/screenshots/terrain-production-before.png` and `terrain-production-after.png`. Record final test output, evidence interpretation, documentation status, and line-count results in the living sections above.

## Interfaces and Dependencies

No external dependency or exported API is added. `StructureTemplate.ResourceYield`, `Resources`, map types, and phase entry points retain their signatures. New helpers remain private to `internal/game`; they map economic features to required terrain, locate the nearest explored match, return biome-default terrain, and resolve one building atomically.

Revision note (2026-07-14): Created from the accepted implementation plan after repository inspection. It fixes whole-Domain range, nearest-distance selection, deterministic ordering, per-building consumption, grass replacement, selection cleanup, evidence, documentation, and validation requirements.

Revision note (2026-07-14): Completed implementation and updated every living section with the Ebitengine capture constraint, automated and visual validation, documentation results, evidence paths, and final line-count review.
