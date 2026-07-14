# Add the Forest Biome

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/61-forest-biome.md`.

## Purpose / Big Picture

After this change, exploration can offer a third biome named Forest. Forest Plots visibly contain many more Trees than existing biomes, a few Boulders, and no Iron Deposits, making a Forest preview an informed choice for Wood production rather than Iron. The home Plot remains Grasslands, while newly exposed frontier Plots are assigned Grasslands, Hills, or Forest with near-equal probability.

## Progress

- [x] (2026-07-14T23:28Z) Confirmed the working tree is clean, `plans/61-forest-biome.md` is the next plan path, and `go test ./...` passes before implementation.
- [x] (2026-07-14T23:28Z) Recorded the accepted Forest weights and biome odds in this ExecPlan; `plans/56-iron-deposit-terrain/screenshots/explored-biomes.png` is the latest full-map visual baseline.
- [x] (2026-07-14T23:31Z) Added Forest generation, assignment, labels, and terrain-consumption behavior.
- [x] (2026-07-14T23:32Z) Added deterministic generator, frontier, selection, resource, and visual-fixture tests.
- [x] (2026-07-14T23:34Z) Updated durable control documents and captured post-change screenshots.
- [x] (2026-07-14T23:35Z) Ran formatting, full tests, race tests, screenshot capture, whitespace, ownership, and line-count validation.

## Surprises & Discoveries

- Observation: The main screenshot suite still writes to the older Iron Deposit plan directory, while the most recent Gold/Market plan captured only focused Market evidence.
  Evidence: `cmd/td/main_test.go` points at `plans/56-iron-deposit-terrain/screenshots`; `plans/60-gold-and-market/screenshots` contains only Market images.
- Observation: A deterministic sequential percentage source made the first focused Forest screenshot arrange terrain in visible bands.
  Evidence: The first `forest-biome.png` grouped the 18 Tree rolls and 3 Boulder rolls at the start of each 100-value cycle; multiplying the deterministic index by 37 modulo 100 preserved the same weights while distributing them naturally.

## Decision Log

- Decision: Generate Forest with 18% Tree, 3% Boulder, 0% Iron Deposit, and 79% empty grass.
  Rationale: This makes Forest three times as tree-heavy as Grasslands while retaining enough buildable terrain for the prototype.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Assign frontier biomes using 33% Grasslands, 33% Hills, and 34% Forest.
  Rationale: A percentage roll cannot divide three ways exactly; Forest receives the single remainder slot through ranges 0-32, 33-65, and 66-99.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Keep the home Plot Grasslands and reuse existing Tree and Boulder assets and behavior.
  Rationale: The request adds an exploration biome, not a new starting-biome choice or art family.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Implementation is complete. Frontier biome rolls now assign Grasslands for 0-32, Hills for 33-65, and Forest for 66-99. `NewForestPlot` and the shared biome dispatch generate Forest at 18% Tree, 3% Boulder, 0% Iron Deposit, and 79% empty grass. Forest labels appear beside exploration controls and in selected-terrain details, while existing Woodcutters and Stone Quarries can consume its Trees and Boulders and leave empty grass. The home Plot remains Grasslands.

Focused tests prove exact biome and terrain boundaries, Forest's lack of Iron Deposits, stored preview/reveal consistency, selection text, and terrain consumption. The existing exploration, production, placement, road, join, combat, Market, pause, and overlay suites remain green. `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `DESIGN.md`, and `ARCHITECTURE.md` now describe the same three-biome behavior and resource tradeoffs.

Visual evidence contains the restored 17-image application suite plus `plans/61-forest-biome/screenshots/forest-biome.png`. The focused image shows a revealed Tree-heavy Forest with occasional Boulders, no Iron Deposits, a clear shared Plot join, a visible `Forest` frontier label, and a selected Tree panel naming the biome. The full `explored-biomes.png` also shows Forest among the expanded frontier choices without obscuring terrain or UI.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, focused and full screenshot capture, `git diff --check`, and the ownership check all succeeded. The final hand-written Go line-count review found no file over 600 lines. `internal/game/build_placement_test.go` remains closest at 564 lines and was unchanged; a future placement feature should split its cases by placement eligibility versus population/resource effects before adding substantial coverage. No unplanned refactor was performed.

## Context and Orientation

`internal/game/map.go` defines the private `plotBiome` values stored on explored Plots and frontier previews. `internal/game/plot_generator.go` maps percentage rolls to biomes and generates terrain from explicit weights. `internal/game/exploration.go` formats the biome label used both beside exploration buttons and in selected-terrain details. `internal/game/resources.go` defines the terrain left after Labour consumes a natural obstacle.

Grasslands currently generates 6% Tree, 3% Boulder, 1% Iron Deposit, and 90% empty grass. Hills generates 3% Tree, 6% Boulder, 3% Iron Deposit, and 88% empty grass. Every road and shared explored edge overwrites generated obstacles. Trees and Boulders block building, can be selected, and are finite production inputs for Woodcutters and Stone Quarries.

## Plan of Work

Add `biomeForest` after the existing biome constants. Add `forestTerrainWeights`, `NewForestPlot() Plot`, a deterministic private Forest constructor, and explicit Forest dispatch in `newPlotForBiomeWithSources`. Change `biomeForRoll` to return Grasslands for 0-32, Hills for 33-65, and Forest for 66-99. Add `Forest` to `biomeLabel` and to the explicit list of biomes whose consumed natural terrain becomes empty grass.

Extend focused tests rather than adding cases to broad state tests. Prove exact Forest terrain boundaries, absence of Iron Deposits, constructor metadata, biome-roll boundaries, preview/reveal consistency, label text, selected-terrain details, and post-Labour clearing. Preserve existing road, shared-edge, placement, exploration, and production regressions.

Update the deterministic explored-biomes screenshot fixture to include a Forest Plot and Forest frontier label, then point the full screenshot suite at this plan. Update `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `DESIGN.md`, and `ARCHITECTURE.md`; no asset, dependency, `ART.md`, or source-convention change is needed.

## Concrete Steps

Run every command from `/home/dave/dev/ai/td`.

1. Edit the biome model, generator, exploration label, and terrain-consumption switch in `internal/game`.
2. Add focused deterministic tests and run `gofmt -w` on changed Go files followed by `go test ./...`.
3. Update the five durable control documents and change the screenshot fixture and destination for Forest evidence.
4. Capture and review screenshots with:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

5. Complete validation with:

    go test ./...
    go test -race ./...
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Record any hand-written code file over 600 lines, or close enough that the next likely change may cross the preference, with a concrete maintenance recommendation. Do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` and `go test -race ./...` must pass. Deterministic tests must prove Forest uses Tree rolls 0-17, Boulder rolls 18-20, and empty grass rolls 21-99, with no Iron Deposit range. Biome assignment must use Grasslands 0-32, Hills 33-65, and Forest 66-99. A frontier preview labeled `Forest` must reveal a Plot whose stored biome and generated terrain are Forest, and selected Forest terrain must report that biome. Woodcutters and Stone Quarries must consume Forest Trees and Boulders normally, while Forest offers no generated input for Iron Mines.

Post-change visual evidence must show a readable Forest label and a revealed Forest Plot with dense Trees, occasional Boulders, no Iron Deposits, open shared joins, and protected roads. Existing exploration controls, terrain selection, construction, resource production, camera, Raid, pause, overlay, and Market behavior must remain intact.

## Idempotence and Recovery

Formatting, tests, race tests, and screenshot capture are safe to rerun. Screenshot capture replaces evidence only under this plan. Biome assignment remains stable because assigned frontier coordinates are not rerolled, and repeated reveal remains idempotent. If dense Forest generation exposes a fixture that assumes an empty Tile, make that fixture explicit rather than weakening production weights.

## Artifacts and Notes

Use `plans/56-iron-deposit-terrain/screenshots/explored-biomes.png` as the pre-change full-map reference. Store all post-change screenshots under `plans/61-forest-biome/screenshots/` and review `explored-biomes.png` specifically.

## Interfaces and Dependencies

Add `NewForestPlot() Plot`, matching the existing Grasslands and Hills constructors within the internal game package. No external dependency, asset, terrain type, exploration control, resource interface, save format, or public application API changes.

Revision note (2026-07-14): Created from the accepted implementation plan after confirming the clean baseline, exact Forest weights, near-equal biome ranges, relevant control documents, and visual-test ownership.

Revision note (2026-07-14): Completed implementation, documented the deterministic visual-fixture adjustment, recorded the accepted screenshots and validation, and closed every living-plan item.
