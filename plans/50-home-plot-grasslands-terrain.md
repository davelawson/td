# Generate Grasslands Terrain on the Home Plot

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/50-home-plot-grasslands-terrain.md`.

## Purpose / Big Picture

After this change, starting a new game shows the home Plot with the same sparse terrain character as any other grasslands Plot: 6% Tree, 3% Boulder, and 91% empty grass by generation weight. The centered Sanctum remains the only initial structure, and its straight road to the north edge remains clear. A contributor can run the game or inspect the new screenshot evidence to see terrain on the starting Plot before exploring.

## Progress

- [x] (2026-07-14 03:10Z) Inspected plot generation, tests, durable control documents, the existing screenshot workflow, and the clean working tree.
- [x] (2026-07-14 03:10Z) Recorded the accepted behavior and created this ExecPlan; `plans/49-frontier-biome-labels/screenshots/running-game.png` is the before-change visual baseline.
- [x] (2026-07-14 03:12Z) Reused ordinary grasslands generation for the home Plot while preserving its Sanctum and north road.
- [x] (2026-07-14 03:15Z) Replaced obsolete open-home tests, moved generator coverage out of `game_test.go`, and made fixed build destinations explicitly empty.
- [x] (2026-07-14 03:14Z) Updated product, game-design, roadmap, architecture, visual-direction, and screenshot evidence documentation.
- [x] (2026-07-14 03:17Z) Ran automated tests, race tests, 50 repeated game-package tests, screenshot capture, whitespace checks, and ownership checks.
- [x] (2026-07-14 03:17Z) Checked hand-written Go file line counts; no file exceeds 600 lines, while the unchanged 596-line `building_bar_test.go` remains near the preference.

## Surprises & Discoveries

- Observation: The home Plot already records `biomeGrasslands` but deliberately uses a separate all-empty generator.
  Evidence: `NewDefaultHomePlot` calls `newOpenGrasslandsPlotWithTweakSource`, while explored grasslands call `newGeneratedPlotWithSources` with `grasslandsTerrainWeights`.
- Observation: The existing baseline test suite passes before implementation.
  Evidence: `go test ./...` passed for every package during planning.
- Observation: `internal/game/game_test.go` is already over the preferred size and contains home-generation tests owned more naturally by `plot_generator_test.go`.
  Evidence: the initial line-count review reported 602 lines for `internal/game/game_test.go`.
- Observation: A single post-change test run can pass even though fixed placement destinations sometimes generate as obstacles.
  Evidence: the first full validation passed, but a later full and race run failed `TestHouseThenDormEnablesFlameBoltTower`, `TestHouseThenBarracksEnablesBowTower`, and `TestBuildDragPlacesCatapultTower` when their destination Tiles were not empty.
- Observation: Explicitly preparing only the fixed build destinations removes randomness without weakening terrain rules.
  Evidence: `go test ./internal/game -count=50`, `go test ./...`, and `go test -race ./...` all passed after adding `setHomeTilesEmpty` to positive build-flow fixtures; Tree, Boulder, and road rejection tests still set and exercise their own terrain.

## Decision Log

- Decision: Give the home Plot exactly the ordinary grasslands terrain weights, then overwrite the existing north road.
  Rationale: The user selected exact grasslands behavior rather than a lighter distribution or a protected clearing. Applying authored infrastructure after generated terrain preserves the established road invariant.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Do not reserve an obstacle-free clearing around the Sanctum.
  Rationale: Only the road and the Sanctum's existing Tile behavior are protected; all other home Tiles should behave like ordinary grasslands.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Keep public constructors unchanged and inject random sources only through private helpers.
  Rationale: The behavior changes, but callers do not need a new map or Plot API. Private sources keep percentage boundaries and override ordering deterministic in tests.
  Date/Author: 2026-07-14 / Codex
- Decision: Move home-generation tests from `game_test.go` to `plot_generator_test.go` as part of this feature.
  Rationale: Those tests directly cover the generator being changed, and the move addresses the existing 602-line warning without introducing an unrelated abstraction.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Implementation is complete. `NewDefaultHomePlot` now uses the shared grasslands generator with independent terrain rolls, then overwrites the existing north road and places the Sanctum. The public API and terrain semantics did not change. Deterministic tests prove the exact 0-5 Tree, 6-8 Boulder, and 9-99 empty ranges on the home Plot, independent Tile tweaks, road override ordering, grasslands metadata, and Sanctum-only initial features.

The post-change `plans/50-home-plot-grasslands-terrain/screenshots/running-game.png` visibly shows sparse Trees and Boulders across the starting Plot while the Sanctum, full north road, exploration controls, and empty buildable grass remain readable. The full screenshot suite contains 16 images in this plan's screenshot directory. `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md` now agree on the same home grasslands behavior; `GAME.md` marks both former empty-home decisions as superseded.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, `go test ./internal/game -count=50`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded. The final line-count review found no hand-written Go file over 600 lines. Moving home-generation tests reduced `internal/game/game_test.go` from 602 to 493 lines. `internal/game/building_bar_test.go` remains close to the preference at 596 lines and was not changed; a future building-bar feature should split it by behavior before adding substantial coverage. `internal/game/build_placement_test.go` is 537 lines after adding deterministic fixture setup.

## Context and Orientation

A Plot is the game's fixed 15-by-15 group of Tiles. `internal/game/plot_generator.go` creates the home, grasslands, and hills Plots. Ordinary grasslands generation assigns each Tile a visual `Tweak` and independently chooses Tree for rolls 0 through 5, Boulder for rolls 6 through 8, or empty grass for rolls 9 through 99. `NewDefaultHomePlot` currently bypasses those weights and fills every Tile with empty grass before drawing the north road and placing the Sanctum.

`internal/game/map.go` stores the home Plot and explored Plots. It clears shared edges when exploration joins two Plots, so generated terrain on a home edge can later become empty grass. `internal/game/scene.go` already renders Tree and Boulder sprites, and `internal/game/building_bar.go` already rejects placement on any terrain other than empty grass. No rendering, asset, gathering, or placement-rule additions are needed.

`PRODUCT.md` and `ARCHITECTURE.md` currently describe the home Plot as open grassland. `GAME.md` records that choice as a gameplay decision. `ROADMAP.md` summarizes the current terrain capability, and `DESIGN.md` defines readability expectations for generated terrain. These durable documents must all describe the new starting behavior without implying that Boulder becomes a resource node. `CODESTYLE.md` requires Go formatting, function comments, tests for pure behavior, and a final review of hand-written code files against a 600-line preference. `ART.md` remains unchanged because the existing Tree and Boulder sprites are reused.

## Plan of Work

Change `NewDefaultHomePlot` in `internal/game/plot_generator.go` to call a private helper accepting both a Tile-tweak source and a percentage-roll source. That helper will call the same shared grasslands generator used by explored Plots, overwrite the center column from the north edge through the Sanctum with road terrain, and place the Sanctum feature at the center. Remove the private all-empty grasslands helper once nothing uses it.

Move the shape, feature, tweak, and north-road home tests from `internal/game/game_test.go` into `internal/game/plot_generator_test.go`. Replace assertions that every non-road Tile is empty with deterministic rolls proving the home Plot can contain Tree, Boulder, and empty grass while the road overrides a generated obstacle. Any test whose purpose depends on placing a structure on a fixed home Tile must explicitly set that target Tile to `terrainEmpty`; tests should not accidentally depend on random production terrain.

Update `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md`. Mark the former open-home decision in `GAME.md` as superseded and add the accepted exact-grasslands decision. Change the screenshot destination in `cmd/td/main_test.go` to this plan's screenshot directory, generate the full existing evidence set, and review `running-game.png` specifically for visible sparse Tree and Boulder terrain with an unobstructed north road and readable Sanctum.

## Concrete Steps

Run every command from `/home/dave/dev/ai/td`.

The before-change visual baseline already exists at `plans/49-frontier-biome-labels/screenshots/running-game.png`; inspect it before code edits to confirm that the home Plot is empty aside from its road and Sanctum.

Edit `internal/game/plot_generator.go` and its tests, then run:

    gofmt -w internal/game/plot_generator.go internal/game/plot_generator_test.go internal/game/game_test.go
    go test ./...

Update fixed-Tile test setup only where the test reports a failure caused by generated home terrain. Update the five relevant control documents and `cmd/td/main_test.go`, then capture post-change evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Complete validation with:

    go test ./...
    go test -race ./...
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Report every hand-written Go file over 600 lines and files close enough to cross the preference on their next likely edit. Recommend a responsibility-based response, but do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` and `go test -race ./...` must pass. Deterministic generator tests must show that the home Plot uses the grasslands Tree, Boulder, and empty ranges; that every Tile retains its independent tweak; that the road wins over generated terrain from the north edge through the center; and that the Sanctum remains the only initial feature. Placement and exploration tests must remain deterministic despite production home terrain being random.

The screenshot `plans/50-home-plot-grasslands-terrain/screenshots/running-game.png` must show sparse Tree and Boulder terrain on the starting Plot while the centered Sanctum, north road, exploration controls, and buildable empty grass remain readable. `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md` must agree that home and explored grasslands use the same weights, road and shared-edge overrides remain, and terrain obstacles are still non-buildable rather than gatherable resources.

## Idempotence and Recovery

Formatting, tests, race tests, and screenshot capture are safe to rerun. Screenshot capture replaces evidence only inside this plan's directory. If random home terrain exposes a test that assumes an empty fixed Tile, make that test's setup explicit instead of weakening production generation. If a road Tile contains an obstacle, fix generation order so the road overwrites terrain rather than adding rendering exceptions.

## Artifacts and Notes

The before screenshot is `plans/49-frontier-biome-labels/screenshots/running-game.png`. Post-change screenshots belong under `plans/50-home-plot-grasslands-terrain/screenshots/`. Focus final terminal evidence on the test, race, whitespace, ownership, and line-count results.

## Interfaces and Dependencies

No external dependency or public API change is required. `NewDefaultHomePlot() Plot` and `NewDefaultMap() Map` retain their signatures. The game package should expose only the existing constructors; a private `newDefaultHomePlotWithSources(nextTweak func() uint16, nextTerrainRoll func() int) Plot` helper will provide deterministic generation for package tests. The existing `grasslandsTerrainWeights`, `newGeneratedPlotWithSources`, `weightedTerrain`, Tree and Boulder sprites, and placement rules must be reused.

Revision note (2026-07-14): Created the ExecPlan from the accepted implementation plan after repository inspection. It fixes exact grasslands weights, road-first acceptance behavior, deterministic testing, documentation ownership, visual evidence, and line-count validation.

Revision note (2026-07-14): Updated every living section after implementation. Recorded the random fixed-destination test failures and fixture response, completed validation and screenshot evidence, and documented the final line-count result.
