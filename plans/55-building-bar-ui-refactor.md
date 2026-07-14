# Building Bar UI Package Refactor

This ExecPlan is a living document. It follows `PLANS.md`.

## Purpose / Big Picture

This refactor moves the building bar's presentation responsibilities from `internal/game` to `internal/ui`. The UI package will own categories, visible ordering, layout, hit testing, hover styling, metadata and tooltip formatting, and drawing. The game package will continue to own structure templates, construction eligibility, drag lifecycle, placement, and gameplay effects. Player-visible behavior and rendered output should remain unchanged.

## Progress

- [x] (2026-07-14) Inspected the existing building bar, tooltip, game/UI package boundary, tests, screenshot harness, and control documents.
- [x] (2026-07-14) Add the UI-owned building-bar model, geometry, formatting, hit testing, and rendering.
- [x] (2026-07-14) Replace game-owned presentation code with a structure-template adapter and gameplay interaction delegates.
- [x] (2026-07-14) Move pure building-bar unit tests beside `internal/ui` and retain game integration coverage.
- [x] (2026-07-14) Update architecture ownership and screenshot evidence.
- [x] (2026-07-14) Run validation and record outcomes.
- [x] (2026-07-14) Check hand-written code file line counts and record any files over the 600-line preference.

## Surprises & Discoveries

- Observation: `internal/game/building_bar_test.go` is already 595 lines, immediately below the 600-line preference.
  Evidence: the baseline line-count command reports 595 lines, reinforcing the value of relocating presentation tests instead of expanding that file.

- Observation: the screenshot harness currently writes to `plans/54-population-portrait-icons/screenshots/`.
  Evidence: `cmd/td/main_test.go` sets its screenshot base path to that plan directory.

- Observation: full screenshots are nondeterministic because newly constructed game maps regenerate terrain.
  Evidence: baseline/current PNG differences begin in the map at x=438 or x=556, while the complete x=0..259 building-bar region is pixel-identical in all five comparison images. Tooltip formatting and bounds are also covered by pure UI tests and the new tooltip capture was visually inspected.

## Decision Log

- Decision: Use a full-presentation boundary while keeping the UI package stateless.
  Rationale: Categories, ordering, geometry, hit testing, formatting, and rendering are UI responsibilities, while host-owned selected-category and hover values avoid introducing a stateful widget abstraction.
  Date/Author: 2026-07-14 / User and Codex.

- Decision: Let `internal/ui` define stable building-bar actions and let `internal/game` map those actions to structure templates and tile features.
  Rationale: The action identifies a widget choice across the package boundary; construction meaning and effects remain gameplay-owned.
  Date/Author: 2026-07-14 / Codex.

## Context and Orientation

`internal/game/building_bar.go`, `building_bar_metadata.go`, and `building_tooltip.go` currently mix gameplay rules with layout, formatting, hit testing, and Ebitengine drawing. `internal/game/building_bar_items.go` defines the stable item/category mapping. `internal/ui` already owns shared palette primitives, text helpers, widgets, and selection-panel presentation. `ARCHITECTURE.md` permits UI helpers to format and draw presentation-neutral data but forbids them from inspecting gameplay state. `CODESTYLE.md` requires doc comments, `gofmt`, package-level tests for pure behavior, and a final line-count check.

## Plan of Work

Add UI-owned building-bar files defining actions, categories, presentation-neutral item/model data, layout, hit testing, metadata and tooltip formatting, availability styling, and drawing. Reuse `ui.ResourceAmounts` and `ui.PopulationAmounts` to avoid duplicate cross-package value types.

Change `internal/game` to build a `ui.BuildingBarModel` from all nine structure templates. Each item carries its game-computed `Buildable` value. Keep phase visibility, action-to-template and action-to-feature mapping, drag lifecycle, placement validation, spending, staffing, and population effects in the game package. Delegate bounds, hit testing, bar/tooltip drawing, and dragged-icon drawing to `internal/ui`.

Move pure layout, formatting, ordering, hit-testing, tooltip, and visual-state tests to `internal/ui`. Keep game tests that verify template adaptation, visibility, selection blocking, menu state, drag behavior, and construction effects. Update `ARCHITECTURE.md` and point screenshot output at this plan directory.

## Concrete Steps

From `/home/dave/dev/ai/td`, preserve the current baseline with:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Implement the package move, then run:

    gofmt -w internal/ui/*.go internal/game/*.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    # Compare the building-bar crop of running-game.png, house-icon.png,
    # barracks-icon.png, dorm-icon.png, and woodcutter-tooltip.png against plan 54.
    # Ignore nondeterministic map-terrain pixels outside the bar.
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | rg '\.go$' | xargs wc -l | sort -n

Report files over 600 lines with a concrete response. Do not perform unplanned refactors, code splits, or library additions without user approval.

## Validation and Acceptance

`go test ./...` must pass for every package. UI tests must cover category ordering, bounds, hit testing, metadata, availability styling, and tooltip formatting/bounds. Game tests must prove that catalog templates and eligibility adapt correctly and that category selection, visibility, drag, placement, and input blocking remain intact. The five building-bar screenshots listed above should compare byte-for-byte with the baseline; if encoding differs, inspect them and document any pixel differences before acceptance. `ARCHITECTURE.md` must describe the new ownership boundary accurately.

## Idempotence and Recovery

All validation commands are safe to repeat. Screenshot capture overwrites evidence files in the active plan directory. If a package move causes compile failures, keep the action/model API narrow and restore behavior by delegating one presentation concern at a time; do not change construction rules to resolve UI boundary issues.

## Artifacts and Notes

Baseline validation on 2026-07-14:

    ok td/assets
    ok td/cmd/td
    ok td/internal/game
    ok td/internal/menu
    ?  td/internal/sound [no test files]
    ok td/internal/ui

The ownership check printed nothing before implementation.

## Interfaces and Dependencies

`internal/ui` will expose `BuildingBarAction`, `BuildingBarCategory`, `BuildingBarItem`, and `BuildingBarModel`; geometry and hit-testing functions for the bar, tabs, and items; and drawing functions for the bar, tooltip, and drag sprite. The data model uses Ebitengine images plus presentation-neutral resource, population, description, production, and tower-stat facts. No new dependency is required.

## Outcomes & Retrospective

Implemented the full-presentation boundary. `internal/ui` now owns building-bar actions, categories, ordering, geometry, hit testing, metadata and tooltip formatting, availability styling, and all bar/tooltip/drag rendering. `internal/game` adapts structure templates and game-computed eligibility into the UI model, consumes UI actions, and retains drag and placement rules with eligibility revalidation at both drag start and drop.

Pure presentation tests now live beside `internal/ui`; focused game tests cover catalog adaptation, live eligibility, host category state, selection blocking, overlay clearing, visibility, drag, and placement. `ARCHITECTURE.md` records the new ownership boundary. Screenshot evidence is stored in `plans/55-building-bar-ui-refactor/screenshots/`; the building-bar crop is pixel-identical to the plan-54 baseline, and the tooltip capture was visually reviewed.

Validation passed on 2026-07-14:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The ownership check printed nothing. The final hand-written Go line-count check found no files above 600 lines. The largest file is `internal/game/build_placement_test.go` at 541 lines; the relocated `internal/ui/building_bar.go` is 376 lines and the former 595-line game building-bar test has been reduced to focused integration coverage.
