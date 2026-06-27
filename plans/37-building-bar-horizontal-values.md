# Building Bar Horizontal Values

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/37-building-bar-horizontal-values.md` and is maintained according to `PLANS.md`.

## Purpose / Big Picture

After this change, the tabbed building bar is wider and easier to scan. Each visible building keeps its 64x64 icon on the left, while construction cost and population or staffing values appear to the right of that icon instead of below it. Buildings the player lacks capacity to construct are drawn with 70% icon opacity, where capacity means the existing resource, population-cost, and staffing checks.

The result is visible by running `go run ./cmd/td`: the left bar is 260 pixels wide, `Housing` remains the default tab, House is fully opaque at the start, and Barracks or other capacity-blocked entries are visibly dimmer.

## Progress

- [x] (2026-06-27 00:00Z) Created this ExecPlan after inspecting the current tabbed building bar, metadata rendering, Raid control offset, screenshot harness, and tests.
- [x] (2026-06-27 00:00Z) Implemented the wider bar, horizontal value layout, and 70% capacity opacity.
- [x] (2026-06-27 00:00Z) Updated tests, docs, screenshot output, and validation evidence.
- [x] (2026-06-27 00:00Z) Ran `go test ./...`, `git diff --check`, screenshot capture, ownership checks, and a final hand-written Go file line-count review.

## Surprises & Discoveries

Observation: The `Next Raid` button position is derived from `buildingBarWidth`.
Evidence: `internal/game/raidui.go` defines `nextRaidButtonX = buildingBarWidth + 42`, so widening the bar automatically moves the button if the constant changes.

Observation: Hover and drag hit testing currently uses only `buildingBarItem.Bounds`, which is the icon rectangle.
Evidence: `buildingBarItemIndexAt()` checks `item.Bounds.Contains(x, y)` and metadata is not part of that rectangle.

Observation: The existing test files are close to the line-count preference.
Evidence: The final line-count review reports `internal/game/game_test.go` at 597 lines and `internal/game/building_bar_test.go` at 572 lines, so new layout-specific coverage was placed in `internal/game/building_bar_layout_test.go`.

## Decision Log

Decision: Use a fixed 260-pixel building bar width.
Rationale: The user chose this width to give comfortable room for values to the right of 64x64 icons.
Date/Author: 2026-06-27 / User and Codex

Decision: Grey only for insufficient construction capacity.
Rationale: The user chose capacity-only greying, so active Raid, breach, terrain, overlay, and pause restrictions should not control icon opacity.
Date/Author: 2026-06-27 / User and Codex

## Outcomes & Retrospective

Implemented the widened horizontal-value building bar. The bar is now 260 pixels wide. Each visible building row keeps the 64x64 icon as the hover and drag target, while resource costs and population or staffing values render to the icon's right. Metadata clicks remain non-interactive, and the whole bar still blocks map selection. Capacity-blocked building icons draw at 70% opacity based only on `canConstructBuilding(item)`.

Validation completed successfully with `go test ./...`, `git diff --check`, screenshot capture through `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot`, and the ownership check. Screenshot evidence is under `plans/37-building-bar-horizontal-values/screenshots/`.

The final hand-written Go line-count review found no file over 600 lines. `internal/game/game_test.go` remains close at 597 lines, and `internal/game/building_bar_test.go` is 572 lines. Future test additions should continue using focused files unless those files are split.

## Context and Orientation

The current tabbed building bar lives in `internal/game`. `building_bar.go` owns bar bounds, tabs, item icon bounds, hover, drag, and item rendering. `building_bar_metadata.go` owns cost and population/staffing value rendering and measurement. `building_bar_items.go` owns stable building IDs and category mapping. `raidui.go` positions the `Next Raid` button relative to `buildingBarWidth`.

Root documents constrain this work. `README.md` and `PRODUCT.md` must describe the current user-visible bar layout and capacity opacity. `GAME.md` must record the intended building-bar UI behavior. `ARCHITECTURE.md` must mention horizontal metadata and capacity-opacity feedback if it describes building-bar ownership. `DESIGN.md` already says build bars should stay compact and visually subordinate; no durable design update is needed unless implementation changes that guidance. `ART.md`, `ROADMAP.md`, and `CODESTYLE.md` do not need updates unless implementation changes their durable truth.

## Plan of Work

First, change `buildingBarWidth` to 260 and update item layout so each item row starts at the left padding, uses a 64x64 icon rectangle, and leaves the remaining row area for metadata to the right. Keep tab order, active tab state, category switching, and stable item IDs unchanged.

Second, change metadata drawing so resource costs and population/staffing values use x coordinates to the right of the icon. Preserve current value formatting, colors, icon order, and hover cost emphasis. Metadata remains non-interactive: only the icon rectangle starts hover and drag, while the full bar still blocks map selection and build drops.

Third, add an icon opacity helper based on `canConstructBuilding(item)`. If the item is constructible, draw the icon at full opacity and allow existing hover brightening. If not constructible, draw the icon with alpha 0.70 and do not brighten it. Do not change drag gating.

Fourth, update tests for width, layout fit, metadata right-of-icon anchors, hover ignoring metadata, capacity opacity, and `Next Raid` avoiding the wider bar. Update screenshot capture output to `plans/37-building-bar-horizontal-values/screenshots/`.

Fifth, update current-state docs and record validation evidence.

## Concrete Steps

Work from `/home/dave/dev/ai/td`.

Run the normal test suite after code and docs are updated:

    go test ./...

Expect all package tests to pass. Run whitespace validation:

    git diff --check

Expect no output. Capture screenshots:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot

Expect screenshots under `plans/37-building-bar-horizontal-values/screenshots/`. Finally, check hand-written Go file line counts:

    rg --files -g '*.go' | xargs wc -l

If any hand-written Go file exceeds 600 lines, record the file, line count, and recommended split before doing unplanned refactor work.

## Validation and Acceptance

Acceptance requires tests that fail without the implementation: the bar is 260 pixels wide, visible item metadata begins to the right of the icon, metadata hit areas do not trigger hover or drag, ineligible capacity-blocked icons report 0.70 opacity, eligible icons report full opacity, and `Next Raid` remains to the right of the bar.

Manual acceptance is: start a new game, see the wider left bar with House and Barracks values to the right of their icons, see House fully opaque, see Barracks greyed because Peasants are missing, switch tabs and see the same horizontal value layout for Economic and Defenses entries.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` describe the current widened horizontal-value bar and capacity opacity without implying new construction rules.

## Idempotence and Recovery

The changes are local to game UI, tests, docs, and screenshot output. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/37-building-bar-horizontal-values/screenshots/`. If bar width causes overlap, inspect `buildingBarWidth` dependents before moving unrelated UI.

## Artifacts and Notes

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
    ok  	td/cmd/td	1.388s

    find . -xdev ! -user dave -printf '%u:%g %p\n'
    # no output

Final line-count notes: no hand-written Go file exceeds 600 lines. The largest files are `internal/game/game_test.go` at 597 lines, `internal/game/building_bar_test.go` at 572 lines, and `internal/game/building_bar.go` at 431 lines.

## Interfaces and Dependencies

No external dependency is required. The implementation uses existing Ebitengine drawing APIs and private helpers in `internal/game`.
