# Building Bar Tabs

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/36-building-bar-tabs.md` and is maintained according to `PLANS.md`.

## Purpose / Big Picture

After this change, the left building bar partitions buildable structures into three tabs instead of showing all buildings at once. The tabs are `Defenses`, `Economic`, and `Housing`, with `Housing` selected by default so the first population path remains visible to a new player. Switching tabs changes only the visible build options; existing construction costs, staffing gates, placement rules, drag behavior, and UI click blocking continue to work.

The result is visible by running `go run ./cmd/td`: the left bar shows the tab row below the top HUD, House and Barracks appear under the default `Housing` tab, and the player can switch to `Economic` or `Defenses` to build those categories.

## Progress

- [x] (2026-06-27 00:00Z) Created this ExecPlan after inspecting the existing building bar, placement tests, screenshot harness, and current line counts.
- [x] (2026-06-27 00:00Z) Implemented tab/category state, rendering, hit testing, and stable placement mapping.
- [x] (2026-06-27 00:00Z) Updated tests, docs, and screenshot output for the tabbed bar.
- [x] (2026-06-27 00:00Z) Ran `go test ./...`, `git diff --check`, screenshot capture, ownership checks, and a final hand-written Go file line-count review.

## Surprises & Discoveries

Observation: `internal/game/building_bar.go` is already close to the 600-line preference.
Evidence: `wc -l internal/game/building_bar.go` reports 560 lines before this work.

Observation: Placement currently depends on building-bar list indices.
Evidence: `buildDragState.itemIndex`, `buildingBarItems()`, and `buildingFeatureForItemIndex()` all use the same flat index order.

Observation: After adding tabs directly, `internal/game/building_bar.go` would exceed the line-count preference.
Evidence: The final implementation splits stable item/category mapping into `internal/game/building_bar_items.go` and cost/population metadata rendering into `internal/game/building_bar_metadata.go`; the final line-count review reports `internal/game/building_bar.go` at 426 lines.

## Decision Log

Decision: Use text tabs and select `Housing` by default.
Rationale: The user chose text tabs with Housing selected first, which keeps the House and Barracks population path immediately visible.
Date/Author: 2026-06-27 / Codex

Decision: Replace visible-list index placement with stable building identifiers.
Rationale: Category filtering changes visible indices, so placement must not infer Tile features from the active tab position.
Date/Author: 2026-06-27 / Codex

Decision: Split building-bar support code while implementing tabs.
Rationale: The existing building-bar file is near the 600-line preference, and tab logic would otherwise push it over.
Date/Author: 2026-06-27 / Codex

## Outcomes & Retrospective

Implemented the tabbed building bar. A new game opens the building bar on `Housing`, showing House and Barracks. The `Economic` tab shows Woodcutter, Stone Quarry, and Iron Mine. The `Defenses` tab shows Bow Tower, Flame Bolt Tower, and Catapult Tower. Clicking tabs switches the visible group without starting a drag or clearing map selection. Build dragging and placement continue to use the same resource, population, staffing, phase, terrain, and UI-blocking rules.

The implementation replaced visible-list placement indices with stable building item IDs, so filtering by category does not change which Tile feature gets placed. Building-bar support code was split by responsibility to keep hand-written Go files below the 600-line preference.

Validation completed successfully with `go test ./...`, `git diff --check`, screenshot capture through `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot`, and an ownership check. Screenshot evidence is under `plans/36-building-bar-tabs/screenshots/`.

The final hand-written Go line-count review found no file over 600 lines. `internal/game/game_test.go` remains close at 597 lines and should continue to be avoided for new tests unless it is split by responsibility first.

## Context and Orientation

The Go module is a local Ebitengine prototype. `internal/game/building_bar.go` currently owns the left bar, visual item metadata, drag start and release, cost and population metadata rendering, and item-to-feature placement mapping. `internal/game/game.go` stores UI state in `gameUI`, including `buildBarHover`, and calls building-bar hover, drag, and drawing helpers. `internal/game/build_placement_test.go` and `internal/game/building_bar_test.go` exercise construction and bar hit testing.

Root documents constrain this work. `README.md` and `PRODUCT.md` must describe the current user-visible tabbed building bar. `GAME.md` must record the gameplay UI decision to group build options by category. `ARCHITECTURE.md` must describe that the game package owns category tabs and stable construction mapping. `DESIGN.md` already covers compact screen-space build bars and does not need a durable rule change unless implementation introduces one. `ART.md`, `ROADMAP.md`, and `CODESTYLE.md` do not need updates unless implementation changes their durable truth.

## Plan of Work

First, add a stable private building identifier for every buildable structure and store that identifier on each `buildingBarItem`. Add a `Feature tileFeature` or equivalent stable mapping so placement no longer depends on the item position in the currently visible category.

Second, add private category state for `Housing`, `Economic`, and `Defenses`. Initialize `gameUI` with the active category set to `Housing`. Add tab bounds under the top HUD inside the existing left bar, using text labels and selected/hover colors that fit the current restrained UI style.

Third, change `buildingBarItems()` so it returns only the active category's visible items in stable category order. The item list starts below the tab row and keeps the existing cost and population metadata rows. The categories are: Housing contains House and Barracks; Economic contains Woodcutter, Stone Quarry, and Iron Mine; Defenses contains Bow Tower, Flame Bolt Tower, and Catapult Tower.

Fourth, update input. Hover should distinguish tab hover from item hover. Clicking a tab should switch the active category, clear item hover, and not start a build drag. Clicking a visible item should keep the existing affordability, population, staffing, and phase gates. Dropping a dragged item should use the item's stable feature metadata.

Fifth, keep the line-count preference healthy by moving category/item-definition helpers or metadata helpers into focused files within `internal/game` rather than growing `building_bar.go` beyond 600 lines.

Sixth, update tests, screenshot harness output, and docs. Tests should prove default `Housing` state, each category's visible items, tab switching, tab click blocking, hover behavior, and one placement path per category.

## Concrete Steps

Work from `/home/dave/dev/ai/td`.

Run the normal test suite after code and docs are updated:

    go test ./...

Expect all package tests to pass. Run whitespace validation:

    git diff --check

Expect no output. Capture screenshots after implementation when a graphical environment is available:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot

Expect screenshots under `plans/36-building-bar-tabs/screenshots/`. Finally, check hand-written Go file line counts:

    rg --files -g '*.go' | xargs wc -l

If any hand-written code file exceeds 600 lines, report the file, line count, and a concrete recommended split before doing unplanned refactor work.

## Validation and Acceptance

Acceptance requires `go test ./...` to pass, including tests that fail without tab implementation. Manual acceptance is: start a new game, see the `Housing` tab active by default with House and Barracks visible, switch to `Economic` and see the three resource buildings, switch to `Defenses` and see the three towers, then place one buildable item from each category after meeting its existing requirements.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` describe tabbed categories without implying scrolling, upgrades, range previews, selling, worker reassignment, or any new construction rules.

## Idempotence and Recovery

The code edits are local to `internal/game`, docs, and screenshot target paths. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/36-building-bar-tabs/screenshots/`. If placement tests fail after switching tabs, inspect whether the active category is changed before pressing the requested item and verify dragged items carry stable feature metadata.

## Artifacts and Notes

Validation transcripts:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.021s
    ok  	td/internal/game	0.394s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/sound	[no test files]
    ?   	td/internal/ui	[no test files]

    git diff --check
    # no output

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
    ok  	td/cmd/td	1.450s

    find . -xdev ! -user dave -printf '%u:%g %p\n'
    # no output

Final line-count notes: no hand-written Go file exceeds 600 lines. The largest files are `internal/game/game_test.go` at 597 lines, `internal/game/building_bar_test.go` at 570 lines, and `internal/game/building_bar.go` at 426 lines after splitting metadata helpers.

## Interfaces and Dependencies

No new external dependency is required. The implementation uses existing Ebitengine text, vector, and image rendering. All new category, tab, and building identifier types remain private to `internal/game`.
