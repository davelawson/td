# Add Dorm

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/42-dorm.md`.

## Purpose / Big Picture

After this change, a player can create Apprentices in normal play. The new `Dorm` structure appears in the Housing tab, costs Wood and Stone, converts one available Peasant into one Apprentice, and makes the Flame Bolt Tower reachable without test-only setup. The behavior can be seen by starting a new game, building a House, building a Dorm, and observing Peasants decrease while Apprentices increase.

## Progress

- [x] (2026-07-13 19:21Z) Generated a Dorm source image with a flat chroma-key background.
- [x] (2026-07-13 19:25Z) Removed the chroma-key background, cropped the subject, and saved `assets/sprites/structures/dorm.png` as a 64x64 transparent PNG.
- [x] (2026-07-13 19:34Z) Wired the Dorm asset, structure template, map feature, rendering, building bar item, and selection panel route into game code.
- [x] (2026-07-13 19:39Z) Added automated tests for asset loading, template metadata, building-bar display, placement conversion, selection-panel data, and House-to-Dorm-to-Flame-Bolt workflow.
- [x] (2026-07-13 19:45Z) Updated current product, game design, roadmap, architecture, README, and art documentation.
- [x] (2026-07-13 19:47Z) Captured visual evidence in `plans/42-dorm/screenshots/`.
- [x] (2026-07-13 19:49Z) Ran `go test ./...`, `git diff --check`, ownership checks, and `git status --short`.
- [x] (2026-07-13 19:49Z) Checked hand-written Go file line counts and reported `internal/game/game_test.go` at 603 lines, over the 600-line preference. No unplanned split was performed.

## Surprises & Discoveries

- Observation: ImageMagick is not installed in this environment, so the asset was resized with Pillow after using the project-approved chroma-key removal helper.
  Evidence: `command -v magick` and `command -v convert` returned no path.

- Observation: The new tests and implementation push `internal/game/game_test.go` just over the 600-line preference.
  Evidence: the line-count review reported `603 internal/game/game_test.go`.

## Decision Log

- Decision: Name the structure `Dorm`, with filename and Go identifiers also using `dorm`/`Dorm`.
  Rationale: The user explicitly changed the proposed name from `Apprentice's Dorm` to `Dorm`.
  Date/Author: 2026-07-13 / Codex

- Decision: Use `10 Wood` and `10 Stone` as the Dorm construction cost.
  Rationale: This matches Barracks and makes the Apprentice converter a parallel Housing-path structure.
  Date/Author: 2026-07-13 / Codex

- Decision: Model Dorm as population conversion, not staffing reservation.
  Rationale: The request says it generates Apprentice population similar to Barracks converting Peasants into Soldiers, so both available and total counts should move from Peasant to Apprentice.
  Date/Author: 2026-07-13 / Codex

## Outcomes & Retrospective

Dorm is implemented as a Housing-tab population converter. It costs 10 Wood and 10 Stone, consumes one available and total Peasant, grants one available and total Apprentice, renders from `assets/sprites/structures/dorm.png`, appears in hover tooltips and selection panels, and enables the Flame Bolt Tower through normal play after building a House.

Validation passed with `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and the ownership check. Visual evidence was written under `plans/42-dorm/screenshots/`.

The final hand-written Go file line-count review found one file over the 600-line preference: `internal/game/game_test.go` at 603 lines. The likely cause is accumulated broad game-state coverage in one test file. The recommended response is a later responsibility-based split that moves focused initialization or state tests into narrower files. That split is not part of this plan and was not performed without user approval.

## Context and Orientation

`td` is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. Structure metadata lives in `internal/game/structures.go`. The map feature enum lives in `internal/game/map.go`. Building-bar stable IDs and category mapping live in `internal/game/building_bar_items.go`, while drag placement in `internal/game/building_bar.go` already applies resource cost, population cost, staff reservation, and population grant atomically. Rendering of map structures lives in `internal/game/scene.go`, and selected-structure panel data lives in `internal/game/selection_panel.go`.

Assets are embedded and loaded in `assets/catalog.go`. Structure sprites are stored under `assets/sprites/structures/` as 64x64 PNGs. The new Dorm sprite must be embedded in the asset catalog and referenced by the Dorm template.

The root control documents constrain this work. `GAME.md`, `PRODUCT.md`, `README.md`, `ROADMAP.md`, and `ARCHITECTURE.md` currently say there is no normal-play Apprentice source; they must be updated after Dorm works. `ART.md` records generated structure asset paths and prompt constraints. `CODESTYLE.md` prefers hand-written code files below 600 lines, so this plan ends with a line-count review. `DESIGN.md` does not need a rule change because the existing building-bar and sprite style are reused.

## Plan of Work

First, add `assets/sprites/structures/dorm.png` as a 64x64 transparent PNG. Update `assets/catalog.go` so the `//go:embed` directive includes it, `StructureSprites` exposes `Dorm`, `NewCatalog()` loads `sprites/structures/dorm.png`, and the returned catalog stores it.

Second, update gameplay metadata. Add a `Dorm StructureTemplate` field to `StructureCatalog` and initialize it with name `Dorm`, a short conversion description, `Resources{Wood: 10, Stone: 10}`, `PopulationCost{Peasants: 1}`, and `PopulationGrant{Apprentices: 1}`. It has no staffing, resource yield, or projectile combat stats.

Third, add placement and rendering support. Add `featureDorm` to the map feature enum. Add `buildingBarDormIndex`, place it in the Housing category after Barracks, map it to `featureDorm`, and map it to `s.structureCatalog.Dorm`. Render `featureDorm` using the Dorm sprite and return `populationBuildingSelectionPanel(s.structureCatalog.Dorm)` when a Dorm tile is selected.

Fourth, update tests. Add or extend tests for asset loading, structure template metadata, Housing tab item order and metadata, placement conversion, invalid-release atomicity, selected Dorm panel data, and a House-to-Dorm-to-Flame-Bolt scenario.

Fifth, update docs and screenshots. Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `ARCHITECTURE.md`, and `ART.md` so they describe Dorm as the normal-play Apprentice source. Extend the screenshot harness in `cmd/td/main_test.go` to capture `dorm-icon.png` and `placed-dorm.png` under `plans/42-dorm/screenshots/`.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Use the built-in image generation tool to create a flat chroma-key source image for a medieval wizardry dorm, then run:

    mkdir -p tmp/imagegen
    cp /home/dave/.codex/generated_images/019f5ccb-8807-7ed0-aa07-876ab610a551/call_hMzPMmDU2pGBoRR8qQ9J7ORD.png tmp/imagegen/dorm-source.png
    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input tmp/imagegen/dorm-source.png --out tmp/imagegen/dorm-alpha.png --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill

Resize and center the result into `assets/sprites/structures/dorm.png`, then patch the code and docs described above.

Run:

    go test ./...
    git diff --check
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short

Finally, check hand-written Go file line counts:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

Report any hand-written Go file above 600 lines with a concrete recommendation. Do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` must pass. New tests must fail before Dorm is wired in and pass after implementation. Asset tests must prove `Dorm` loads as a 64x64 sprite. Structure tests must prove the template values. Placement tests must prove the Dorm consumes one Peasant, grants one Apprentice, deducts resources, places the right feature, and preserves state on invalid release. Workflow tests must prove House then Dorm enables Flame Bolt Tower placement.

Manual acceptance is: start the game, build a House, build a Dorm on an empty grass-like tile, and observe Peasants change from `2/2` to `1/1` and Apprentices change from `0/0` to `1/1`. The Flame Bolt Tower icon should become eligible when resources are sufficient and the Apprentice is available.

Documentation acceptance is: `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `ARCHITECTURE.md`, and `ART.md` no longer claim there is no normal-play Apprentice source, and they do not imply timed recruitment, reassignment, removal, upgrades, or a broader population system exists.

## Idempotence and Recovery

The code changes are additive and safe to reapply from git. If asset loading fails, confirm the `//go:embed` path exactly matches `assets/sprites/structures/dorm.png`. If a visual evidence capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries` and keep automated tests as validation. If line-count review finds an oversized file, report it and wait for user approval before expanding scope.

## Artifacts and Notes

The generated source image is at:

    /home/dave/.codex/generated_images/019f5ccb-8807-7ed0-aa07-876ab610a551/call_hMzPMmDU2pGBoRR8qQ9J7ORD.png

The final project asset is:

    assets/sprites/structures/dorm.png

## Interfaces and Dependencies

Use the existing Go module and Ebitengine dependency. Do not add a new runtime dependency. Keep Dorm within the existing `StructureTemplate`, `PopulationCost`, and `PopulationGrant` interfaces:

    StructureCatalog.Dorm StructureTemplate
    assets.StructureSprites.Dorm *ebiten.Image
    featureDorm tileFeature
    buildingBarDormIndex buildingBarItemID

The final behavior must reuse `placeDraggedBuilding` so resource deduction, population cost, staff reservation, population grant, and tile placement remain one atomic operation.
