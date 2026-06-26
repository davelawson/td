# Add Barracks Building

This ExecPlan adds Barracks as a placeable population-conversion building. A Barracks costs Wood and Stone, consumes Peasant population, and creates Soldier population so normal play can reach Soldier-staffed towers after building Houses.

## Progress

- [x] (2026-06-26T00:00:00Z) Inspected the existing structure catalog, building bar, population bookkeeping, placement tests, asset catalog, screenshot harness, and control documents.
- [x] (2026-06-26T00:00:00Z) Generate and embed the Barracks structure sprite.
- [x] (2026-06-26T00:00:00Z) Add Barracks structure metadata, map feature, rendering, placement, and selection-panel behavior.
- [x] (2026-06-26T00:00:00Z) Add focused tests for assets, structure metadata, building-bar ordering, population conversion, selection, and regression behavior.
- [x] (2026-06-26T00:00:00Z) Update control documents and screenshot capture evidence paths.
- [x] (2026-06-26T00:00:00Z) Run `go test ./...`, `git diff --check`, ownership check, screenshot capture, and line-count review for hand-written code files.

## Decisions

- Barracks is a non-combat population-conversion building, not a tower.
- “Consumes 2 Peasant population” means permanent conversion: deduct 2 from Peasant available and total, then add 2 to Soldier available and total.
- Barracks construction is eligible only when resources cover `10 Wood` and `10 Stone` and at least 2 Peasants are available.
- Barracks uses a generated 64x64 PNG structure sprite following `ART.md`.
- This change does not add structure removal, staff release, timed recruitment, upgrades, Apprentice generation, new placement rules, or combat behavior.

## Context

`internal/game/structures.go` owns structure templates. `internal/game/building_bar.go` creates the left-side building choices and applies resource, staffing, and population-grant effects on successful drops. `internal/game/hud.go` owns population availability and totals. `assets/catalog.go` embeds runtime sprites. `internal/game/map.go`, `scene.go`, `selection.go`, and `selection_panel.go` own feature identity, drawing, click selection, and object details.

The existing House already grants Peasants. Towers already reserve available staff without changing totals. Barracks needs a third population operation: permanent population cost plus grant.

## Concrete Steps

First, generate `assets/sprites/structures/barracks.png` as a 64x64 sharp, high-contrast, low-detail medieval barracks sprite with transparent background. Add it to the embedded asset list, `StructureSprites`, `NewCatalog`, and asset tests.

Second, add a `PopulationCost` type and field to `StructureTemplate`. Add `Barracks` to `StructureCatalog` with `Cost: Resources{Wood: 10, Stone: 10}`, `PopulationCost: PopulationCost{Peasants: 2}`, and `PopulationGrant: PopulationGrant{Soldiers: 2}`. Add helper methods in HUD/population code for checking and applying population cost together with grants.

Third, add `featureBarracks`, render it with the Barracks sprite, include it in structure click selection, and add a selection panel for non-combat population buildings that shows cost, consumed Peasants, and granted Soldiers. Keep tower panels unchanged.

Fourth, update the building bar to show item order `House`, `Barracks`, `Bow Tower`, `Flame Bolt Tower`, `Catapult Tower`. Construction eligibility must require resources, population cost, and staffing. On valid placement, deduct resources, apply population cost, reserve staff, apply population grant, and then set the Tile feature. Invalid placement must leave all counts unchanged. The compact metadata row should show population costs with negative values and grants with positive values.

Fifth, update focused tests. Cover Barracks sprite loading, catalog metadata, item order and costs, population metadata, drag gating with insufficient Peasants, successful conversion after one House, invalid-drop preservation, Soldier output enabling Bow Tower placement, Barracks selection, and selection-panel rows.

Sixth, update `README.md`, `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `ART.md`. Update `cmd/td/main_test.go` screenshot output to `plans/34-barracks-building/screenshots/` and add Barracks icon and placed-Barracks screenshots.

Seventh, run validation:

```sh
go test ./...
git diff --check
find . -xdev ! -user dave -printf '%u:%g %p\n'
TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
wc -l internal/game/*.go assets/*.go cmd/td/*.go internal/menu/*.go internal/ui/*.go internal/sound/*.go
```

For the line-count review, report any hand-written code file over 600 lines or close enough that the next likely change will push it over. Recommend a concrete follow-up and do not perform unplanned splits unless already approved.

## Acceptance

Running a new game starts with no Peasants or Soldiers. The Barracks icon appears in the building bar but is not eligible until the player builds a House. After placing a House, placing Barracks spends 10 Wood and 10 Stone, changes Peasants from `2/2` to `0/0`, and changes Soldiers from `0/0` to `2/2`. A Bow Tower can then be placed using one Soldier, reducing Soldiers to `1/2`.

`go test ./...` passes. Screenshot evidence exists under `plans/34-barracks-building/screenshots/`. Documentation states that Barracks is the current normal-play Soldier source and that Apprentices still have no normal-play source.

## Outcomes & Retrospective

Barracks is implemented as a non-combat population-conversion building. The building bar now shows House, Barracks, Bow Tower, Flame Bolt Tower, and Catapult Tower. Barracks requires 10 Wood, 10 Stone, and 2 available Peasants; successful placement consumes Peasants from available and total counts and grants Soldiers to available and total counts.

Validation passed:

```sh
go test ./...
git diff --check
find . -xdev ! -user dave -printf '%u:%g %p\n'
TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
```

Screenshot evidence was captured under `plans/34-barracks-building/screenshots/`, including `barracks-icon.png` and `placed-barracks.png`.

Line-count review found no hand-written file over 600 lines. `internal/game/game_test.go` is still close at 597 lines. This file was already near the project preference before this work, and this change did not add to it. Recommended follow-up is to split future game-state tests by responsibility before adding more cases there; no unplanned split was performed.
