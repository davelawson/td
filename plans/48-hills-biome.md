# Add Hills Biome Generation

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/48-hills-biome.md`.

## Purpose / Big Picture

After this change, hills join grasslands as the second generated Plot biome. The authored home Plot remains open grasslands, while each newly explored Plot independently has an equal chance to become grasslands or hills. Hills remain mostly buildable but reverse the grasslands obstacle emphasis: 3% Tree, 6% Boulder, and 91% empty grass. The player can recognize hills by their Boulder-heavy terrain composition without adding height, new assets, gathering, or resource-node rules.

## Progress

- [x] (2026-07-14T02:29Z) Created this ExecPlan from the accepted implementation plan.
- [x] (2026-07-14T02:31Z) Implemented hills generation and equal random biome selection with deterministic test seams.
- [x] (2026-07-14T02:31Z) Updated exploration and generator tests for both biomes and retained reveal state.
- [x] (2026-07-14T02:33Z) Updated the relevant durable control documents and captured multi-Plot visual evidence.
- [x] (2026-07-14T02:34Z) Ran formatting, tests, race tests, screenshot capture, whitespace, ownership, and line-count validation.

## Surprises & Discoveries

- Observation: The working tree was clean before implementation.
  Evidence: `git status --short` printed no files.
- Observation: `internal/game/game_test.go` already exceeds the 600-line preference.
  Evidence: the pre-change line-count review reported 602 lines.
- Observation: Capturing all four adjacent generated Plots makes the biome composition contrast easier to review than a single random Plot.
  Evidence: `plans/48-hills-biome/screenshots/explored-biomes.png` shows a visibly Boulder-heavy Plot alongside Tree-heavier generated terrain while keeping the north road and joins clear.

## Decision Log

- Decision: Select grasslands for biome rolls 0 through 49 and hills for rolls 50 through 99.
  Rationale: The user chose independent 50/50 selection for every newly explored non-home Plot.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Generate hills with Tree 3%, Boulder 6%, and empty grass 91%.
  Rationale: Swapping the grasslands Tree and Boulder weights makes hills stone-biased while preserving the current sparse obstacle density and buildability.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Keep Boulder as non-buildable terrain rather than a gatherable Stone resource node.
  Rationale: The user selected the existing terrain semantics, avoiding premature gathering and depletion rules.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Reuse the current grass ground and Tree and Boulder sprites.
  Rationale: Terrain composition is enough to make this prototype biome observable; tint, height, biome labels, and new art remain outside this slice.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Implementation completed hills as the second generated Plot biome. The authored home Plot remains open grasslands. Every first-time non-home reveal now consumes one percentage roll, selecting grasslands for 0-49 or hills for 50-99, then stores the generated Plot so repeated reveal cannot reroll it. Grasslands remains Tree 6%, Boulder 3%, and empty grass 91%; hills uses Tree 3%, Boulder 6%, and empty grass 91%.

The shared generator keeps Tile tweak generation independent from biome and terrain rolls. Existing road and shared-edge cleanup still runs after biome terrain generation, and existing Tree and Boulder rendering and non-buildable placement behavior require no biome-specific branches. Durable product, game, roadmap, architecture, and design documents now describe both biomes and retain the limitation that Boulder is terrain rather than a gatherable Stone node.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and the ownership check all succeeded. Visual evidence was written under `plans/48-hills-biome/screenshots/`; the multi-Plot image shows the Boulder-heavy hills character without obscuring roads, joins, or empty building space.

The final hand-written Go line-count review found `internal/game/game_test.go` at 602 lines, unchanged by this feature, and `internal/game/building_bar_test.go` near the preference at 596 lines. A separate approved maintenance change should move default-state and map tests out of `game_test.go` and consider splitting building-bar tests by behavior. This feature did not perform that unrelated refactor.

## Context and Orientation

`internal/game/map.go` stores each Plot's biome and owns Plot reveal orchestration, north-road application, and shared-edge cleanup. `internal/game/plot_generator.go` owns the authored home Plot, generated grasslands Plots, explicit Tree and Boulder percentage weights, and random sources. `internal/game/scene.go` already renders Tree and Boulder sprites independently of biome, while building placement accepts only empty terrain. The feature therefore needs a new biome value and generation policy, not new terrain, rendering, placement, or asset behavior.

`PRODUCT.md` describes current user-visible exploration. `GAME.md` records biome and terrain design decisions. `ROADMAP.md` summarizes the current prototype phase. `ARCHITECTURE.md` records map and generator ownership. `DESIGN.md` constrains visual readability. These files currently state that every explored Plot is grasslands and must be updated with the two-biome behavior.

## Plan of Work

Add `biomeHills` beside `biomeGrasslands`. In `internal/game/plot_generator.go`, keep grasslands weights at Tree 6 and Boulder 3, add hills weights at Tree 3 and Boulder 6, and create a shared biome generator that fills Tiles using the selected weights. Keep `NewGrasslandsPlot`, add `NewHillsPlot`, and add a generated explored-Plot constructor that consumes one random percentage roll to choose grasslands for 0-49 or hills for 50-99. Private source-injection helpers must let tests control biome, Tile tweak, and terrain rolls independently.

Change `Map.revealPlot` to use the random explored-Plot constructor. Generate and store a biome only on the first reveal; revealing an already explored coordinate must still return without changing its biome or Tiles. Apply the existing north road and shared-edge cleanup after terrain generation so both biomes preserve road and join invariants.

Expand generator tests to cover hills metadata, both biome-selection boundaries, hills terrain boundaries, and representative Tree/Boulder/empty output. Update exploration tests so production reveal accepts either implemented biome and repeated reveal preserves the chosen Plot. Retain existing placement, road, and edge tests because both biomes use the same terrain enum values.

Update the five relevant root control documents. Change screenshot capture to write evidence under this plan and capture an explored scene representative of the Boulder-heavy hills composition. Do not add new art or asset catalog entries.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go`, `internal/game/plot_generator.go`, and focused tests under `internal/game`.
2. Run `gofmt -w internal/game/*.go` and `go test ./...`.
3. Update `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md`.
4. Update `cmd/td/main_test.go` only as needed to capture representative hills evidence under `plans/48-hills-biome/screenshots/`.
5. Run:

    gofmt -w internal/game/*.go cmd/td/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Record any hand-written code file over 600 lines with a concrete recommendation. Do not perform an unrelated split unless the user separately approves it; the accepted scope explicitly leaves the existing 602-line `internal/game/game_test.go` unchanged.

## Validation and Acceptance

Acceptance is met when the home Plot remains open grasslands, every first-time non-home reveal independently chooses grasslands or hills with equal probability, grasslands remains Tree 6% and Boulder 3%, and hills uses Tree 3% and Boulder 6%. Both biomes must remain 91% empty by weight. Roads and cleared shared edges override generated obstacles, repeated reveal must preserve existing biome and terrain, and Tree and Boulder must remain non-buildable and non-gatherable. Automated tests, screenshot capture, whitespace, ownership, and line-count checks must pass.

## Idempotence and Recovery

The code and document edits are additive and formatting and validation commands are safe to rerun. Plot reveal remains idempotent because it checks explored storage before generation. If generated terrain blocks protected roads or joins, fix the generation-versus-incorporation order rather than adding rendering exceptions. If random visual capture does not show a representative hills Plot, use a deterministic screenshot-only random seed; do not expose a player-facing biome-selection API.

## Artifacts and Notes

Store rendered evidence under `plans/48-hills-biome/screenshots/`. The evidence should show at least one newly explored Plot with the intended Boulder-heavy hills character. Record final validation results and line counts in this plan.

## Interfaces and Dependencies

No external dependencies or exported cross-package APIs are required. Within `internal/game`, the completed implementation includes `biomeHills`, `hillsTerrainWeights`, `NewHillsPlot() Plot`, and a constructor for a randomly selected explored Plot. Existing terrain values, sprites, building rules, and save behavior remain unchanged; no save system currently exists.
