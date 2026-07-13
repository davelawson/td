# Add Grasslands Biome Plot Generation

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/44-grasslands-biome-generation.md`.

## Purpose / Big Picture

After this change, every explored Plot records the biome that generated it. The first implemented biome is grasslands. When the player explores a new Plot during calm play, the Plot is generated through biome-aware terrain logic rather than always being a completely empty grass board. The player can see sparse forest tiles appear in newly explored grasslands while roads, shared borders, building rules, and Raid path behavior remain readable.

## Progress

- [x] (2026-07-13T21:24Z) Created this ExecPlan from the accepted plan and user refinement requiring `plot_generator.go`.
- [x] (2026-07-13T21:31Z) Added biome metadata and moved plot generation logic into `internal/game/plot_generator.go`.
- [x] (2026-07-13T21:35Z) Added tests for grasslands biome metadata and generated forest terrain; updated exploration tests for biome-backed revealed Plots.
- [x] (2026-07-13T21:39Z) Updated `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, `DESIGN.md`, and `ROADMAP.md` to describe grasslands biome generation.
- [x] (2026-07-13T21:44Z) Ran `go test ./...`, screenshot capture, `git diff --check`, ownership check, and Go file line-count review.

## Surprises & Discoveries

- Observation: The working tree was clean before implementation.
  Evidence: `git status --short` produced no output.
- Observation: The next available plan number is 44.
  Evidence: the highest existing plan file was `plans/43-exploring-additional-plots.md`.
- Observation: `internal/game/game_test.go` remains over the 600-line preference and `internal/game/building_bar_test.go` is close to the preference.
  Evidence: final line-count review reported `602 internal/game/game_test.go` and `596 internal/game/building_bar_test.go`.

## Decision Log

- Decision: Put new plot generation behavior in `internal/game/plot_generator.go`.
  Rationale: The user explicitly requested that file, and it keeps generation decisions separate from exploration input and rendering.
  Date/Author: 2026-07-13 / Codex
- Decision: Use `biomeGrasslands` as the only biome value for this slice.
  Rationale: The user requested only one biome for now, while preserving an obvious extension point for later biomes.
  Date/Author: 2026-07-13 / Codex
- Decision: Generate sparse forest from each tile's random tweak value.
  Rationale: Tile tweaks already provide per-tile random data; deriving terrain from that value avoids introducing a second random stream and keeps tests deterministic with caller-provided tweak sources.
  Date/Author: 2026-07-13 / Codex
- Decision: Keep the home Plot open grassland even though it stores the grasslands biome.
  Rationale: The starting Plot is intentionally forgiving and currently documented as open grassland with a north road.
  Date/Author: 2026-07-13 / Codex

## Outcomes & Retrospective

Implementation completed the first biome-backed plot generation slice. Each Plot now stores a biome, the current biome is grasslands, and `internal/game/plot_generator.go` owns home and explored Plot generation. The home Plot remains open grassland with its existing Sanctum and north road. Newly explored Plots are generated as grasslands with mostly empty grass and sparse forest Tiles, then existing north-road and shared-edge cleanup rules are applied by map reveal orchestration.

Validation passed on 2026-07-13: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded. Screenshot evidence was written under `plans/44-grasslands-biome-generation/screenshots/`.

The final hand-written Go file line-count review found one file over the 600-line preference: `internal/game/game_test.go` at 602 lines. It also found `internal/game/building_bar_test.go` close to the limit at 596 lines. Recommended follow-up is a separate approved test-organization change that splits map/default-state tests out of `game_test.go` and considers splitting building-bar rendering tests by category or behavior before adding more cases. No unplanned split was performed.

## Context and Orientation

The current game map code lives in `internal/game/map.go`. A `Map` owns a home `Plot` and a map of explored Plots keyed by `plotCoordinate`. `revealPlot` creates a newly explored Plot, applies a north-chain road when needed, stores it, and clears shared edges with already explored neighbors. Rendering in `internal/game/scene.go` draws `terrainEmpty`, `terrainRoad`, and `terrainForest`; building placement in `internal/game/building_bar.go` only accepts `terrainEmpty` tiles with no feature.

`GAME.md` says Plots are 15x15 groups of Tiles and that Plot contents should have a dominant character. `PRODUCT.md` currently says newly explored Plots are usable grassland. `DESIGN.md` requires roads, tile boundaries, exploration affordances, and readable placement states. `ARCHITECTURE.md` says `internal/game/map.go` owns prototype map data and rendering uses stored Plot data. `CODESTYLE.md` requires Go doc comments and a final line-count review for substantial code changes.

## Plan of Work

Add a private `plotBiome` type and `biomeGrasslands` value in the game package, and add `Biome plotBiome` to `Plot`. Move plot creation helpers into a new file, `internal/game/plot_generator.go`. This file will create the home Plot, create generated grasslands Plots, create Tiles with tweaks, and provide the random tile tweak helper.

Change `Map.revealPlot` in `internal/game/map.go` so new non-home Plots are generated as grasslands. The existing road and edge cleanup functions remain in `map.go`, because they describe how a generated Plot is incorporated into the explored Domain rather than how its biome terrain is chosen.

For grasslands generation, fill every tile as empty grass with a random tweak, then mark a sparse subset as `terrainForest` based on that tweak. The home Plot uses the grasslands biome but should not use random forest terrain. After `revealPlot` generates a grasslands Plot, `applyNorthRoadIfNeeded` and `clearSharedEdges` must still override generated terrain where roads or clear shared borders are required.

Add focused tests to `internal/game`. Tests should prove the home Plot and explored Plots store `biomeGrasslands`, generated grasslands can include forest terrain with a deterministic tweak source, north-chain roads override generated forest, and shared edges are cleared after reveal. Existing tests should continue proving placement rejects forest and accepts empty terrain.

Update `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, and `DESIGN.md` so durable project truth matches the new behavior. Update `ROADMAP.md` only if the current-phase wording needs to distinguish biome-backed generation from plain empty grassland.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go` to add `plotBiome`, add `Plot.Biome`, remove generation helper implementations from this file, and call the new grasslands generator from `revealPlot`.
2. Add `internal/game/plot_generator.go` with `NewDefaultHomePlot`, `newDefaultHomePlotWithTweakSource`, `NewGrasslandsPlot`, `newGrasslandsPlotWithTweakSource`, `newTile`, and `randomTileTweak`.
3. Update or add tests under `internal/game` for biome metadata and terrain generation invariants.
4. Update the relevant root control documents.
5. Run:

    gofmt -w internal/game/map.go internal/game/plot_generator.go internal/game/*_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

If any hand-written Go file exceeds 600 lines, record it in `Outcomes & Retrospective` with a concrete recommendation. Do not do an unplanned split unless the user approves it.

## Validation and Acceptance

Acceptance is met when a new game starts with a grasslands-biome home Plot that remains open grassland with its existing Sanctum and north road, and exploring a new adjacent Plot generates a grasslands-biome Plot with sparse forest terrain. Forest must not appear on north-chain road tiles after road application or on shared edges after edge cleanup. Building placement remains valid only on empty grass and invalid on forest. `go test ./...` and `git diff --check` must pass. Screenshot evidence should show newly explored grasslands terrain if the capture path exercises exploration.

Documentation acceptance is that `PRODUCT.md` describes current explored Plot generation, `GAME.md` records the grasslands biome decision, `ARCHITECTURE.md` names `internal/game/plot_generator.go` as the plot generation owner, and `DESIGN.md` keeps readability constraints current.

## Idempotence and Recovery

The change is additive and safe to rerun. Re-running generation tests, screenshot capture, and formatting should not damage state. `revealPlot` must remain idempotent: revealing an already explored Plot returns without replacing its Tiles, biome, or structures. If validation fails because random generation creates unexpected terrain on a protected road or shared edge, fix the order of road and edge overrides rather than loosening tests.

## Artifacts and Notes

Visual evidence should be written under `plans/44-grasslands-biome-generation/screenshots/` if screenshot capture is updated for this plan. Terminal evidence should include the final `go test ./...`, `git diff --check`, ownership check, and line-count review outputs in this plan's outcome section.

## Interfaces and Dependencies

No new external dependencies are required. Use the existing Go standard library `math/rand` package and existing Ebitengine rendering paths. Keep all new interfaces private to `internal/game`.

At the end of implementation, the game package should include:

    type plotBiome int

    const (
        biomeGrasslands plotBiome = iota
    )

    type Plot struct {
        Biome plotBiome
        Tiles [plotSize][plotSize]Tile
    }

    func NewGrasslandsPlot() Plot

The exact helper names may remain private where they are only used by tests or map construction, but the generation logic must live in `internal/game/plot_generator.go`.
