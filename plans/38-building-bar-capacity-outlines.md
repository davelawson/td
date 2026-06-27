# Building Bar Capacity Outlines

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/38-building-bar-capacity-outlines.md` and is maintained according to `PLANS.md`.

## Purpose / Big Picture

After this change, the building bar communicates build capacity with the icon slot outline as well as icon opacity. A building the player can construct has a bright green icon-square outline. A building blocked by insufficient resources, population cost, or staffing has a red icon-square outline. The slot remains unfilled, preserving the dark background and current icon presentation.

The result is visible by running `go run ./cmd/td`: on a new game the House icon slot has a green outline, while Barracks and other capacity-blocked choices have red outlines until their requirements are met.

## Progress

- [x] (2026-06-27 00:00Z) Created this ExecPlan after inspecting the current widened tabbed building bar, capacity opacity helper, color palette, and tests.
- [x] (2026-06-27 00:00Z) Implemented capacity-colored icon slot outlines.
- [x] (2026-06-27 00:00Z) Updated tests, docs, screenshot output, and validation evidence.
- [x] (2026-06-27 00:00Z) Ran `go test ./...`, `git diff --check`, screenshot capture, ownership checks, and a final hand-written Go file line-count review.

## Surprises & Discoveries

Observation: Building icon slot rendering already separates the dark fill and the outline.
Evidence: `drawBuildingBarItem()` calls `vector.FillRect(..., colors.plotBackdrop, false)` and then `vector.StrokeRect(..., colors.fieldEdge, false)`.

Observation: The opacity rule already encodes the requested build-capacity state.
Evidence: `buildingBarIconAlpha()` returns full opacity when `canConstructBuilding(item)` is true and 0.70 otherwise.

## Decision Log

Decision: Use the same capacity rule as icon opacity for outline color.
Rationale: The user wants can-build/cannot-build clarity, and the current capacity rule already excludes temporary phase, terrain, and overlay restrictions.
Date/Author: 2026-06-27 / Codex

## Outcomes & Retrospective

Implemented green/red capacity outlines for building icon squares. The icon slot background remains dark and unfilled. Buildable entries use a bright green outline, and capacity-blocked entries use a red outline while keeping the existing 70% icon opacity. The rule uses `canConstructBuilding(item)`, so active Raid, breach, terrain, pause, and overlay state do not affect the outline color.

Validation completed successfully with `go test ./...`, `git diff --check`, screenshot capture through `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot`, and the ownership check. Screenshot evidence is under `plans/38-building-bar-capacity-outlines/screenshots/`.

The final hand-written Go line-count review found no file over 600 lines. `internal/game/game_test.go` remains close at 597 lines, and `internal/game/building_bar_test.go` remains 572 lines. Future tests should continue using focused files unless those files are split.

## Context and Orientation

The current building bar lives under `internal/game`. `building_bar.go` owns icon slot rendering, hover, drag, and opacity. `building_bar_layout_test.go` already has focused capacity-opacity coverage and is the right place for focused outline tests. `colors.go` owns the game-local color palette, while `internal/ui/colors.go` owns shared palette primitives.

Root documents constrain this work. `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` describe current building-bar capacity feedback and must mention green/red outlines after implementation. `DESIGN.md` already says UI needs clear interaction states and does not need a change unless a broader reusable visual rule is introduced. `ART.md`, `ROADMAP.md`, and `CODESTYLE.md` do not need updates.

## Plan of Work

First, add game-local colors for buildable and blocked building slot outlines. Use a bright readable green and red that are distinct from existing bronze and resource colors without changing the shared UI palette unless reuse outside this bar becomes necessary.

Second, add a private helper on `State` that returns the building icon outline color from `canConstructBuilding(item)`. In `drawBuildingBarItem()`, keep the current dark fill and replace the slot `StrokeRect` color with that helper. Keep tab outlines, bar edge, metadata, icon opacity, hover brightening, and drag gating unchanged.

Third, add focused tests proving House is green at game start, Barracks is red at game start, Barracks turns green when Peasants are available, and a missing-staff tower is red while a staffed tower is green. These tests should use the helper rather than pixel assertions.

Fourth, update docs, screenshot capture output to `plans/38-building-bar-capacity-outlines/screenshots/`, and validation notes.

## Concrete Steps

Work from `/home/dave/dev/ai/td`.

Run the normal test suite after code and docs are updated:

    go test ./...

Expect all package tests to pass. Run whitespace validation:

    git diff --check

Expect no output. Capture screenshots:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot

Expect screenshots under `plans/38-building-bar-capacity-outlines/screenshots/`. Finally, check hand-written Go file line counts:

    rg --files -g '*.go' | xargs wc -l

If any hand-written Go file exceeds 600 lines, record the file, line count, and recommended split before doing unplanned refactor work.

## Validation and Acceptance

Acceptance requires tests that fail without the implementation: the helper returns green for capacity-buildable items and red for capacity-blocked items, and existing opacity, hover, drag, tab, and metadata tests continue to pass.

Manual acceptance is: start a new game, see House with a bright green unfilled icon-square outline, see Barracks with a red unfilled icon-square outline, then create enough Peasants and see Barracks become green.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` describe green/red icon-square outlines without implying new construction rules.

## Idempotence and Recovery

The changes are local to game UI, tests, docs, and screenshot output. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/38-building-bar-capacity-outlines/screenshots/`.

## Artifacts and Notes

Validation transcripts:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.021s
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/sound	[no test files]
    ?   	td/internal/ui	[no test files]

    git diff --check
    # no output

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
    ok  	td/cmd/td	1.522s

    find . -xdev ! -user dave -printf '%u:%g %p\n'
    # no output

Final line-count notes: no hand-written Go file exceeds 600 lines. The largest files are `internal/game/game_test.go` at 597 lines, `internal/game/building_bar_test.go` at 572 lines, and `internal/game/building_bar.go` at 439 lines.

## Interfaces and Dependencies

No external dependency is required. The implementation uses existing Ebitengine vector drawing and private helpers in `internal/game`.
