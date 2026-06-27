# Building Menu Tooltips

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/39-building-menu-tooltips.md` and is maintained according to `PLANS.md`.

## Purpose / Big Picture

After this change, a player can move the mouse over any building icon in the left building menu and see a compact tooltip with the building's name, description, construction cost, staffing requirement, and current implemented effect or combat stats. This makes the existing House, Barracks, economic buildings, and defensive towers understandable without requiring the player to already know what each icon means.

The result is visible by running `go run ./cmd/td`, starting a new game, and hovering over entries in the `Housing`, `Economic`, and `Defenses` tabs. A parchment-colored tooltip appears to the right of the building bar. The tooltip is informational only; it does not add upgrades, selling, range previews, reassignment, or new construction rules.

## Progress

- [x] (2026-06-27 00:00Z) Created this ExecPlan after inspecting the building bar, structure templates, root control documents, screenshot harness, and line counts.
- [x] (2026-06-27 00:00Z) Added structure descriptions and tooltip formatting/rendering helpers.
- [x] (2026-06-27 00:00Z) Added focused tests for tooltip content, hover targeting, bounds, and drag suppression without growing existing large test files.
- [x] (2026-06-27 00:00Z) Updated user-visible and architectural documentation.
- [x] (2026-06-27 00:00Z) Captured screenshot evidence under this plan directory.
- [x] (2026-06-27 00:00Z) Ran `go test ./...`, `git diff --check`, ownership checks, and a final hand-written Go file line-count review.

## Surprises & Discoveries

Observation: The building bar already tracks hover by visible item index and deliberately ignores metadata rows.
Evidence: `updateBuildingBarHover()` writes `s.ui.buildBarHover = s.buildingBarItemIndexAt(input.CursorX, input.CursorY)`, and `buildingBarItemIndexAt()` only checks icon bounds.

Observation: `internal/game/building_bar_test.go` is close to the 600-line preference.
Evidence: `wc -l internal/game/building_bar_test.go` reported 572 lines before this plan, so tooltip tests should live in a new focused file.

Observation: Existing hover screenshot targets were enough to verify the new tooltip visually once their output path was moved to this plan.
Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot` produced `plans/39-building-menu-tooltips/screenshots/house-icon.png`, which shows the House tooltip beside the building bar.

## Decision Log

- Decision: Show tooltips for every build-menu entry, not only combat towers.
  Rationale: The user asked for mousing over a building in the building menu, and the menu contains housing, economic, and defensive building entries. A consistent tooltip helps all entries and avoids making non-tower icons less discoverable.
  Date/Author: 2026-06-27 / Codex

- Decision: Reuse existing building-icon hover state and keep metadata rows non-interactive.
  Rationale: Current hover and drag behavior is already icon-based. Reusing it keeps tooltips aligned with the existing interaction target and avoids changing placement or selection behavior.
  Date/Author: 2026-06-27 / Codex

- Decision: Put tooltip implementation and tests in focused new files.
  Rationale: `building_bar.go` and `building_bar_test.go` are already responsible for many behaviors, while this work is mostly presentation formatting. New focused files keep line counts and responsibilities manageable.
  Date/Author: 2026-06-27 / Codex

## Outcomes & Retrospective

Implemented building-menu hover tooltips for every buildable entry across the `Housing`, `Economic`, and `Defenses` tabs. Tooltips show each structure's name, description, construction cost, staffing requirement, and its implemented population effect, post-Raid production, or combat stats. Tooltips are render-only and reuse existing icon hover state, so metadata rows, tabs, empty bar areas, placement, selection, camera input, pause, and overlay behavior are unchanged.

Validation completed successfully with `go test ./...`, `git diff --check`, screenshot capture through `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot`, and the ownership check. Screenshot evidence is under `plans/39-building-menu-tooltips/screenshots/`.

The final hand-written Go line-count review found no file over 600 lines. `internal/game/game_test.go` remains close at 597 lines, and `internal/game/building_bar_test.go` remains close at 572 lines. This plan avoided growing those files by adding `internal/game/building_tooltip_test.go`; future changes touching those large test files should consider a responsibility-based split or another focused test file.

## Context and Orientation

The Go module is a local Ebitengine prototype. The executable under `cmd/td` polls input and calls `internal/game.State.Update` and `internal/game.State.Draw`. The game package owns the active map, resources, populations, structure templates, building bar, placement rules, selection, Raids, and rendering. The left building menu is called the building bar in code.

`internal/game/structures.go` defines `StructureTemplate`, which stores shared metadata for structures such as name, sprite, construction cost, staffing requirements, population effects, economic resource yield, and combat stats. `internal/game/building_bar.go` renders the tabbed 260-pixel left building bar, tracks hover, starts drags, and places buildings. `internal/game/building_bar_items.go` maps stable building IDs to templates and tabs. `internal/game/building_bar_metadata.go` renders compact cost and population values beside icons.

Root documents constrain this work. `README.md` and `PRODUCT.md` describe current user-visible behavior and must mention hover tooltips after implementation. `GAME.md` captures intended game-design decisions and should record that build-facing UI now exposes descriptions, costs, staffing, and implemented effects through hover tooltips. `DESIGN.md` says building bars should remain compact, readable, and should only imply implemented actions. `ARCHITECTURE.md` documents package ownership and should mention tooltip metadata/rendering under `internal/game`. `ROADMAP.md`, `ART.md`, and `CODESTYLE.md` do not need updates unless implementation changes product direction, art guidance, or coding conventions.

## Plan of Work

First, add a short `Description string` field to `StructureTemplate` in `internal/game/structures.go` and fill it for House, Barracks, Woodcutter, Stone Quarry, Iron Mine, Bow Tower, Flame Bolt Tower, and Catapult Tower. The descriptions must describe only implemented behavior. Sanctum may remain without a description because it is not a build-menu entry.

Second, add a focused tooltip renderer in `internal/game/building_tooltip.go`. It should build tooltip content from the hovered visible building item, using existing costs, staffing requirements, population effects, economic yields, and combat stats. The tooltip should anchor to the right of the building bar near the hovered icon, clamp inside the drawable area, use existing fonts and palette colors, and draw after the building bar but before build-drag icons, Raid controls, selection panels, and the in-game overlay. It should return no tooltip when no building icon is hovered or while a build drag is active.

Third, update `State.Draw` or `drawBuildingBar` so the tooltip appears above the map and bar without changing input routing. Do not change `buildingBarItemIndexAt`, drag gating, placement, tab switching, camera controls, selection, or pause behavior.

Fourth, add `internal/game/building_tooltip_test.go` with pure tests around tooltip content and bounds. Tests should verify all visible categories produce the correct building-specific tooltip content, metadata rows do not produce tooltips, and tooltip bounds stay to the right of the bar and within the drawable area.

Fifth, update `cmd/td/main_test.go` screenshot output to `plans/39-building-menu-tooltips/screenshots/` so visual evidence for this rendered change lands under the active plan. Then update `README.md`, `PRODUCT.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` to match the implemented behavior.

## Concrete Steps

Work from `/home/dave/dev/ai/td`.

Edit the plan first, then code and tests:

    gofmt -w internal/game/structures.go internal/game/building_bar.go internal/game/building_tooltip.go internal/game/building_tooltip_test.go

Run the normal test suite:

    go test ./...

Expect all packages to pass. Run whitespace validation:

    git diff --check

Expect no output. Capture screenshots:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot

Expect screenshots under `plans/39-building-menu-tooltips/screenshots/`. If the local graphics environment cannot capture screenshots, record the exact error here and rely on automated tests plus manual launch validation.

Check ownership from the repository root:

    find . -xdev ! -user dave -printf '%u:%g %p\n'

Expect no output. Finally, check hand-written Go file line counts:

    rg --files -g '*.go' | xargs wc -l

If any hand-written Go file exceeds 600 lines, record the file, line count, and recommended split before doing unplanned refactor work.

## Validation and Acceptance

Acceptance requires automated tests that fail before the tooltip implementation and pass after it. The tests should prove that hovering a build-menu icon yields a tooltip with title, description, cost, staffing, and effect or combat stats; that non-icon areas do not yield a tooltip; and that tooltip bounds are clamped within the screen.

Manual acceptance is: start the app with `go run ./cmd/td`, start a new game, hover over House and Barracks in `Housing`, hover over Woodcutter, Stone Quarry, and Iron Mine in `Economic`, and hover over Bow, Flame Bolt, and Catapult in `Defenses`. Each hovered icon should show a compact tooltip to the right of the building bar. The tooltip should not appear over metadata rows, should not block dragging, and should not imply unimplemented commands.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` describe hover tooltips without implying new build rules, upgrades, selling, range previews, staff reassignment, or staff release.

## Idempotence and Recovery

The changes are local to game UI, structure metadata, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/39-building-menu-tooltips/screenshots/`.

If a tooltip layout test fails because the tooltip is too tall for a very small drawable height, clamp its bounds and keep the visual content readable at the default 1920x1080 target rather than shrinking fonts. If tooltip copy becomes too long, shorten descriptions before widening the building bar or changing the existing compact metadata.

## Artifacts and Notes

Validation transcripts:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.022s
    ok  	td/internal/game	0.487s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/sound	[no test files]
    ?   	td/internal/ui	[no test files]

    git diff --check
    # no output

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot
    ok  	td/cmd/td	1.525s

    find . -xdev ! -user dave -printf '%u:%g %p\n'
    # no output

Final line-count notes: no hand-written Go file exceeds 600 lines. The largest files are `internal/game/game_test.go` at 597 lines, `internal/game/building_bar_test.go` at 572 lines, and `internal/game/building_bar.go` at 439 lines.

Screenshot evidence:

    plans/39-building-menu-tooltips/screenshots/house-icon.png
    plans/39-building-menu-tooltips/screenshots/barracks-icon.png
    plans/39-building-menu-tooltips/screenshots/running-game.png

## Interfaces and Dependencies

No external dependency is required. The implementation uses existing Ebitengine image, text, and vector drawing APIs already used by the game package.

At completion, `StructureTemplate` in `internal/game/structures.go` has this additional field:

    Description string

The tooltip helpers are private to `internal/game`; no exported API is added.

Revision note: This plan was updated during implementation to record completed tooltip code, tests, documentation, screenshot evidence, validation transcripts, and line-count results. The final tooltip copy uses "after each defeated Raid" for economic production so the wording matches existing game terminology.
