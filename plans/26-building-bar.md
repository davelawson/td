# Add Visual Building Bar

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds. This plan follows `PLANS.md`.

## Purpose / Big Picture

After this change, the first game screen shows a building bar along the left edge of the playable scene. The bar displays the current two tower types, Bow Tower and Flame Bolt Tower, using their existing sprites as icons. This gives the future base-building workflow a visible home without adding tower placement, build intent, resource costs, previews, upgrades, or terrain build rules yet.

The behavior is visible by starting a game and observing a dark vertical bar below the top HUD. The top HUD remains full-width. Clicking inside the bar is treated as UI interaction and does not select map objects or clear the current selection.

## Progress

- [x] (2026-05-26T01:02:13Z) Inspected current rendering, selection, Raid UI, screenshot harness, and control documents.
- [x] (2026-05-26T01:02:13Z) Confirmed decisions with the user: the building bar is visual-only, and it starts below the top HUD.
- [x] (2026-05-26T01:09:00Z) Added building-bar rendering, reused the existing tower sprites, moved `Next Raid` beside the bar, and blocked bar clicks from map selection.
- [x] (2026-05-26T01:10:00Z) Added focused tests for bar bounds, item bounds, click blocking, and `Next Raid` avoiding the bar.
- [x] (2026-05-26T01:10:00Z) Updated screenshot capture output to `plans/26-building-bar/screenshots/`.
- [x] (2026-05-26T01:15:00Z) Updated current-state, design, game-design, roadmap, and architecture documentation.
- [x] (2026-05-26T01:16:00Z) Ran validation commands and captured screenshot evidence.
- [x] (2026-05-26T01:16:00Z) Checked hand-written code-file line counts and reported files over the 600-line preference.

## Surprises & Discoveries

- Observation: The existing `Next Raid` button overlapped the new left bar's horizontal space.
  Evidence: The old button started at `X=42`, while the building bar is 96 pixels wide.

- Observation: The new screenshot shows the bar and Raid button are both visible and non-overlapping.
  Evidence: `plans/26-building-bar/screenshots/running-game.png` shows the bar on the left and `Next Raid` to its right.

## Decision Log

- Decision: Make the first building bar visual-only.
  Rationale: The user chose a visual-only slice, and build selection or placement would require additional gameplay design around costs, valid terrain, previews, and placement state.
  Date/Author: 2026-05-26 / Codex.

- Decision: Anchor the building bar below the top HUD, not over the full window height.
  Rationale: The user chose to preserve the current top status bar as a full-width HUD while the new bar covers the playable scene's left edge.
  Date/Author: 2026-05-26 / Codex.

## Outcomes & Retrospective

Implemented the visual-only building bar. The left edge of the playable scene below the top HUD now contains two tower icon slots: Bow Tower first and Flame Bolt Tower second. The bar reuses the existing structure sprites from the loaded structure catalog. Clicking the bar is treated as UI input and does not clear the current map selection. The `Next Raid` button moved to the right of the bar so it remains accessible.

Validation passed on 2026-05-26:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The ownership check printed nothing. Screenshot evidence was written under `plans/26-building-bar/screenshots/`, including:

    running-game.png
    active-raid.png
    selected-tower.png
    selected-raider.png
    paused-game.png
    ingame-menu.png

Line-count review on 2026-05-26:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Only `internal/game/game_test.go` exceeds the 600-line preference at 605 lines. This was pre-existing before the implementation and was already reported by earlier plans. The new building-bar tests were placed in `internal/game/building_bar_test.go` to avoid increasing it. Recommended follow-up, if the project wants to address the existing overage, is to split `game_test.go` by responsibility, such as moving camera, HUD, or overlay tests into focused files. No extra refactor was performed because that split was not part of this feature scope.

Revision note, 2026-05-26 / Codex: Completed implementation, validation, screenshot capture, documentation updates, and line-count review.

## Context and Orientation

The repository is a Go/Ebitengine tower-defense prototype. `cmd/td/main.go` owns Ebitengine startup and input polling, then passes a `game.Input` value into `internal/game.State.Update`. `internal/game/game.go` owns top-level game state and draw order. `internal/game/hud.go` renders the top status bar. `internal/game/raidui.go` renders the bottom-left `Next Raid` button. `internal/game/selection.go` owns selected structure and raider hit testing, and it already blocks clicks on screen-space UI such as `Next Raid` and the selection panel from changing map selection.

The two tower sprites already exist in the typed asset catalog and are exposed through `s.structureCatalog.BowTower.Sprite` and `s.structureCatalog.FlameBoltTower.Sprite`. This plan must reuse those images and must not decode assets by path from gameplay code.

This work changes current user-visible behavior, so `README.md` and `PRODUCT.md` must describe the visual building bar and explicitly state that tower placement is not implemented. This work also changes `internal/game` ownership of screen-space UI and input blocking, so `ARCHITECTURE.md` must mention the building bar. Because the first build-option UI is a meaningful base-building design decision, `GAME.md` must record that the current building bar is an early visual affordance only. `DESIGN.md` and `ROADMAP.md` were updated because the implementation established durable guidance for visual-only build affordances and changed the current phase description. `ART.md` and `CODESTYLE.md` do not need updates.

## Plan of Work

First, add a new `internal/game/building_bar.go` file. Define private constants for a stable vertical bar width, item size, padding, and gap. Define private helpers on `State` for `buildingBarBounds`, `buildingBarItems`, `buildingBarContains`, and `drawBuildingBar`. The bar bounds must be `X: 0`, `Y: topBarHeight`, `W: buildingBarWidth`, and `H: s.ui.height - topBarHeight`. The item list must contain exactly two items in order: Bow Tower then Flame Bolt Tower.

Second, draw the bar in `State.Draw` after the world, enemies, projectiles, and top HUD are drawn, and before Raid controls, the selection panel, and the overlay. Fill the bar with a dark surface, draw a bronze right edge, draw each item as a stable square slot, and draw each tower sprite centered inside its slot. If a sprite is nil, draw only the slot so the UI remains stable. Do not draw labels or add hover state for this first slice.

Third, update input blocking in `internal/game/selection.go`. In `updateSelection`, add `s.buildingBarContains(input.CursorX, input.CursorY)` to the existing UI guard so clicking the bar does not select map objects and does not clear the current selection. Do not add build placement state or build-intent selection.

Fourth, update the `Next Raid` button position if necessary so it no longer sits under the new left bar. Preserve its size, label, disabled behavior, and click behavior. If moved, update screenshot-target click coordinates and any tests or helpers that click it.

Fifth, add `internal/game/building_bar_test.go`. Test that the bar fills the scene's left edge below the HUD, the two item slots exist inside the bar in the expected order, and clicks inside the bar do not clear an existing structure selection or select an underlying map tile.

Sixth, update `cmd/td/main_test.go` so screenshot capture writes to `plans/26-building-bar/screenshots/` and uses any adjusted `Next Raid` click coordinate. Update `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to match the implemented visual-only behavior.

## Concrete Steps

From `/home/dave/dev/ai/td`, edit code and docs with `apply_patch`. Run `gofmt` on changed Go files after edits.

Run:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Report any hand-written code file over 600 lines with a concrete recommendation. Do not perform unplanned file splits or refactors without user approval.

## Validation and Acceptance

`go test ./...` must pass. The new tests must prove the building bar fills the left edge of the playable scene below the top HUD, exposes two tower item bounds, and blocks map selection clicks.

Manual acceptance is: start the game, see a vertical bar along the left edge below the top HUD, see Bow Tower and Flame Bolt Tower icons stacked inside it, select an existing map object, click inside the building bar, and observe that selection remains unchanged. Starting a Raid and selection-panel behavior must still work.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` describe the visual building bar without implying tower placement, resource spending, build previews, upgrades, or broader base-building systems exist.

## Idempotence and Recovery

The code edits are local and additive. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/26-building-bar/screenshots/`. If click-blocking tests fail, inspect the guard order in `updateSelection`. If Raid screenshot capture fails after moving the button, inspect the `Next Raid` button bounds and screenshot click coordinate.

## Artifacts and Notes

Record final validation transcripts, screenshot artifact paths, and line-count findings in `Outcomes & Retrospective`.

## Interfaces and Dependencies

No new external dependencies are required. The implementation stays in `internal/game` and uses existing Ebitengine drawing APIs plus `td/internal/ui`.

The final code should include private helpers equivalent to:

    type buildingBarItem struct {
        Template *StructureTemplate
        Bounds   ui.Button[int]
    }

    func (s *State) buildingBarBounds() ui.Button[int]
    func (s *State) buildingBarItems() []buildingBarItem
    func (s *State) buildingBarContains(x, y int) bool
    func (s *State) drawBuildingBar(screen *ebiten.Image)
