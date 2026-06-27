# Economic Buildings

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/35-economic-buildings.md` and is maintained according to `PLANS.md`.

## Purpose / Big Picture

After this change, a player can add the first resource-producing buildings to the home Plot. Woodcutters, Stone Quarries, and Iron Mines appear in the building bar, reserve one available Peasant when built, and add 10 Wood, Stone, or Metal respectively after each successfully defeated Raid. This makes the current resource economy visible and repeatable without adding timed gathering, worker reassignment, building removal, or a fourth resource.

The result is visible by running `go run ./cmd/td`, building Houses to create Peasants, placing economic buildings, then completing a Raid and seeing the matching top-bar resource counts increase.

## Progress

- [x] (2026-06-27 00:00Z) Created this ExecPlan after inspecting the existing structure catalog, building bar, placement, selection, asset catalog, and Raid completion code.
- [x] (2026-06-27 00:00Z) Generated Woodcutter, Stone Quarry, and Iron Mine sprites on chroma-key backgrounds and converted them to 64x64 transparent PNGs under `assets/sprites/structures/`.
- [x] (2026-06-27 11:40Z) Added economic building templates, map features, asset catalog fields, rendering, selection panels, building-bar entries, placement mapping, and Raid-completion payouts.
- [x] (2026-06-27 11:40Z) Updated tests for catalog loading, templates, building-bar ordering and gating, placement, selection, and Raid completion rewards.
- [x] (2026-06-27 11:40Z) Updated `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `ARCHITECTURE.md`, and `ART.md` so durable docs match the new current behavior.
- [x] (2026-06-27 11:40Z) Captured visual evidence under `plans/35-economic-buildings/screenshots/`.
- [x] (2026-06-27 11:40Z) Ran `go test ./...`, `git diff --check`, ownership checks, and a final hand-written Go file line-count review.

## Surprises & Discoveries

- Observation: The existing third resource is consistently named `Metal` in code, docs, assets, and UI.
  Evidence: `internal/game/structures.go` defines `Resources.Metal`, `internal/game/hud.go` renders Metal, and `assets/sprites/icons/metal.png` is loaded by the catalog.
- Observation: Raid completion has a single success path.
  Evidence: `internal/game/raid.go` calls `completeRaid()` only when the Raid is active, no enemies remain, and no pending enemies remain. Breach clears the Raid through `applySanctumContact()` instead.
- Observation: The expanded building bar still fits the 1920x1080 prototype viewport.
  Evidence: `plans/35-economic-buildings/screenshots/placed-woodcutter.png` shows all eight icons and their metadata visible above the bottom edge.
- Observation: `internal/game/game_test.go` is already close to the 600-line preference and remains unchanged by this plan.
  Evidence: final line-count review reports `internal/game/game_test.go` at 597 lines.

## Decision Log

- Decision: Keep using the existing `Metal` resource and make "Iron Mine" the producer name for Metal.
  Rationale: Renaming or adding Iron would be a broad resource-model change beyond this slice. The user accepted the recommended Metal mapping.
  Date/Author: 2026-06-27 / Codex
- Decision: Treat "occupy one peasant population" as reserving one available Peasant while leaving total Peasants unchanged.
  Rationale: This matches current tower staffing semantics and avoids permanent population loss for a worker assignment.
  Date/Author: 2026-06-27 / Codex
- Decision: Economic buildings pay out once at successful Raid completion.
  Rationale: The user specified that each building adds 10 points of its resource following the defeat of each Raid.
  Date/Author: 2026-06-27 / Codex

## Outcomes & Retrospective

Implemented the first economic buildings. The building bar now includes Woodcutter, Stone Quarry, and Iron Mine between Barracks and combat towers. Each building uses a generated 64x64 structure sprite, costs the requested resources, reserves one available Peasant, can be selected for an informational panel, and produces 10 Wood, Stone, or Metal after each successful Raid completion. Breached Raids do not trigger economic payouts.

Validation completed successfully with `go test ./...`, `git diff --check`, screenshot capture through `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot`, and an ownership check. The screenshot evidence is under `plans/35-economic-buildings/screenshots/`.

The final hand-written Go line-count review found no file over 600 lines. Two files are close enough to watch: `internal/game/game_test.go` is 597 lines and was not touched by this plan; `internal/game/building_bar.go` is 560 lines after adding economic buildings. Recommended follow-up is to avoid growing `game_test.go` further and, when the next building-bar feature arrives, split building-bar rendering, metadata formatting, or test helpers by responsibility before adding substantial new behavior. No unplanned split was performed.

## Context and Orientation

The Go module is a local Ebitengine prototype. `assets/catalog.go` embeds and loads required sprites. `internal/game/structures.go` defines structure templates and construction metadata. `internal/game/map.go` defines tile features. `internal/game/building_bar.go` lists buildable structures, gates dragging by resources and population availability, places features on empty Tiles, deducts resources, reserves staff, and applies population costs or grants. `internal/game/scene.go` renders placed structure sprites. `internal/game/selection_panel.go` renders selected structure information. `internal/game/raid.go` owns Raid lifecycle and calls `completeRaid()` after a successful Raid.

Root documents constrain this work. `GAME.md` records intended resource, population, building, and Raid behavior. `PRODUCT.md` and `README.md` describe current user-visible behavior and must mention economic buildings only after they work. `ROADMAP.md` should reflect that resource generation now has a first Raid-completion implementation. `ARCHITECTURE.md` should assign economic building templates, placement, and Raid-completion payouts to `internal/game`. `ART.md` should record generated asset paths and prompt constraints. `CODESTYLE.md` requires Go doc comments and a final line-count review for hand-written code files.

## Plan of Work

First, install the three generated sprites at `assets/sprites/structures/woodcutter.png`, `assets/sprites/structures/stone-quarry.png`, and `assets/sprites/structures/iron-mine.png`. Update `assets/catalog.go` so the `//go:embed` directive includes them, `StructureSprites` exposes them, `NewCatalog()` loads them, and the returned catalog stores them.

Second, update `internal/game/structures.go`. Add economic metadata to `StructureTemplate` as a `ResourceYield Resources` field. Add `Woodcutter`, `StoneQuarry`, and `IronMine` to `StructureCatalog`. Their costs are 10 Wood, 10 Wood plus 10 Stone, and 10 Wood plus 10 Stone plus 10 Metal. Their staffing is `StaffingRequirements{Peasants: 1}`. Their resource yields are `Resources{Wood: 10}`, `Resources{Stone: 10}`, and `Resources{Metal: 10}`. They have no projectile combat stats.

Third, update placed feature handling. Add `featureWoodcutter`, `featureStoneQuarry`, and `featureIronMine` to `internal/game/map.go`. Render them in `internal/game/scene.go`, map building-bar indices to them in `internal/game/building_bar.go`, and show selection panel rows for Structure, Cost, Required Peasants, and Produces.

Fourth, update building UI and placement. Extend `buildingBarItems()` to list House, Barracks, Woodcutter, Stone Quarry, Iron Mine, Bow Tower, Flame Bolt Tower, and Catapult Tower. Keep existing cost display and Peasant requirement row behavior. Placement should reserve one available Peasant through the existing `Staffing` path and should not reduce total Peasants.

Fifth, add Raid-completion resource generation. In `internal/game/raid.go`, call a helper from `completeRaid()` that scans the home Plot for economic building features and adds their `ResourceYield` to `s.status.resources`. Because breach never calls `completeRaid()`, failed Raids do not produce economic resources.

Sixth, update tests across assets and game behavior, then update docs and screenshot harness output to `plans/35-economic-buildings/screenshots/`.

## Concrete Steps

Work from `/home/dave/dev/ai/td`.

Run the normal test suite after code and docs are updated:

    go test ./...

Expect all package tests to pass. Run whitespace validation:

    git diff --check

Expect no output. Capture screenshots after implementation when a graphical environment is available:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot

Expect screenshots under `plans/35-economic-buildings/screenshots/`. Finally, check hand-written Go file line counts:

    rg --files -g '*.go' | xargs wc -l

If any file exceeds 600 lines, report the file, line count, and a concrete recommended split before doing unplanned refactor work.

## Validation and Acceptance

Acceptance requires `go test ./...` to pass, including new tests that fail without the implementation. Catalog tests must prove all three sprites load as 64x64 images. Structure tests must prove the economic templates have the requested names, costs, one-Peasant staffing, and 10-resource yields. Building-bar tests must prove all eight structures are available in order and still fit the left bar at 1920x1080. Placement tests must prove economic buildings reserve one Peasant, deduct construction resources, and map to the right features. Raid tests must prove successful Raid completion adds 10 per placed economic building and breach does not pay out.

Manual acceptance is: start the game, build enough Houses to create Peasants, place one or more economic buildings on empty grass-like Tiles, complete a Raid, and observe the top-bar resources increase by 10 for each matching economic building.

Documentation acceptance is: `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `ARCHITECTURE.md`, and `ART.md` describe the implemented behavior without implying timed gathering, worker reassignment, building removal, or an Iron resource separate from Metal.

## Idempotence and Recovery

The code edits are additive and can be rerun safely from git. The generated source images remain under `/home/dave/.codex/generated_images/019f017e-a3ea-7062-b124-746f3b39eeca/`; rerunning local post-processing should overwrite only the three intended structure PNGs. If asset loading fails, verify the `//go:embed` paths exactly match the new files under `assets/sprites/structures/`. If the building bar no longer fits vertically, reduce vertical gaps before changing icon size or removing required metadata.

## Artifacts and Notes

Generated source images were created with the built-in image generation tool on a flat `#00ff00` chroma-key background, then locally converted to transparent 64x64 PNGs. The installed project assets are:

    assets/sprites/structures/woodcutter.png
    assets/sprites/structures/stone-quarry.png
    assets/sprites/structures/iron-mine.png

Validation transcripts:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/sound	[no test files]
    ?   	td/internal/ui	[no test files]

    git diff --check
    # no output

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
    ok  	td/cmd/td	1.432s

Plan revision note: implementation completed on 2026-06-27. This plan was updated to record completed work, validation evidence, screenshot artifacts, and the final line-count review required by `PLANS.md` and `CODESTYLE.md`.

## Interfaces and Dependencies

No new external Go dependency is required. The implementation uses existing Ebitengine image rendering, the existing `assets.Catalog`, and existing `internal/game` construction and Raid state. At the end of the work, `StructureTemplate` includes `ResourceYield Resources`, and `StructureCatalog` exposes `Woodcutter`, `StoneQuarry`, and `IronMine`.
