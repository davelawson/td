# Add tower drag placement

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/28-tower-drag-placement.md`.

## Purpose / Big Picture

The current building bar shows Bow Tower and Flame Bolt Tower icons and costs, but the player cannot use it to build. After this change, the player can left-drag an affordable tower icon from the building bar, see a half-sized copy follow the cursor, and release over an eligible Tile to spend resources and place that tower. This makes the first resource-spending base-building action observable without adding gathering, upgrades, selling, build previews, range indicators, pathfinding, saves, or broader economy systems.

The behavior is visible by starting the game with `go run ./cmd/td`: drag the Bow Tower icon from the left building bar onto an empty grass-like Tile during calm play, release, and the Tile becomes a Bow Tower while the top-bar resources decrease by 30 Wood, 10 Stone, and 10 Metal. With the default starting resources, the Flame Bolt Tower cannot be dragged because it costs 20 Metal and the player starts with 12 Metal.

## Progress

- [x] (2026-05-26T03:24:46Z) Inspected current building-bar, resource HUD, map, selection, Raid, input, screenshot, and control-document code.
- [x] (2026-05-26T03:24:46Z) Confirmed design choices with the user: placement is calm grass only, and building remains allowed while SPACE-paused.
- [x] (2026-05-26T03:24:46Z) Created this ExecPlan at `plans/28-tower-drag-placement.md`.
- [x] (2026-05-26T03:30:00Z) Implemented mouse held/released input, build-drag state, placement validation, resource spending, and dragged icon rendering.
- [x] (2026-05-26T03:34:00Z) Added focused automated tests for drag start, unaffordable drag blocking, successful placement, invalid placement, active-Raid blocking, and paused calm placement.
- [x] (2026-05-26T03:38:00Z) Updated screenshot capture to `plans/28-tower-drag-placement/screenshots/` and captured visual evidence.
- [x] (2026-05-26T03:36:00Z) Updated `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` for the new current behavior and design decision.
- [x] (2026-05-26T03:39:00Z) Ran `gofmt`, `go test ./...`, screenshot capture, `git diff --check`, ownership check, and recorded the results.
- [x] (2026-05-26T03:39:00Z) Checked hand-written Go file line counts against the 600-line preference and reported the existing overage.

## Surprises & Discoveries

- Observation: The building bar already exposes tower sprites and construction costs and already gates hover highlighting by affordability.
  Evidence: `internal/game/building_bar.go` defines `buildingBarItems()`, `canAffordBuildingCost()`, and `buildingBarItemHighlighted()`.

- Observation: Combat already scans placed `tileFeature` values across the home Plot, so newly placed towers can join combat without a separate tower-instance registry.
  Evidence: `internal/game/combat.go` iterates every Tile in `fireCombatTowers()` and resolves tower stats through `combatTowerTemplate()`.

- Observation: Screenshot evidence confirms the placed tower appears and resources are deducted.
  Evidence: `plans/28-tower-drag-placement/screenshots/placed-tower.png` shows a second Bow Tower east of the road and top-bar resources of 50 Wood, 35 Stone, and 2 Metal after placement.

## Decision Log

- Decision: Treat "grass-like" placement as the existing `terrainEmpty` Tile terrain.
  Rationale: The current map model does not have a separate Grass enum; `terrainEmpty` is the existing playable open ground distinct from `terrainRoad` and `terrainForest`.
  Date/Author: 2026-05-26 / Codex

- Decision: Allow placement only during calm play while blocking active Raids and breached state.
  Rationale: The user chose calm grass only; the current calm predicate is `status.phase == phaseCalm`, `raid.active == false`, and `raid.breached == false`.
  Date/Author: 2026-05-26 / Codex

- Decision: Allow build placement while SPACE-paused if the game is otherwise calm.
  Rationale: The user chose paused placement as allowed, matching existing paused inspection behavior.
  Date/Author: 2026-05-26 / Codex

- Decision: Cancel invalid drops without spending resources or adding extra invalid-drop UI in this slice.
  Rationale: This keeps the first placement workflow small and testable while still making invalid releases safe and understandable.
  Date/Author: 2026-05-26 / Codex

## Outcomes & Retrospective

Implemented the requested tower drag placement slice. The game now tracks held and released left mouse input, starts a build drag only for affordable tower icons during calm play, renders a half-sized dragged tower icon, validates drops against empty grass-like home Plot Tiles, deducts costs, and writes the corresponding tower feature into map state. Newly placed towers participate in the existing combat scan because combat already reads tower features from the home Plot. The Flame Bolt Tower remains unaffordable with default resources and cannot be dragged from the bar.

Validation completed:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.(go)$' | xargs -r wc -l | sort -n

`go test ./...` passed. Screenshot evidence was written under `plans/28-tower-drag-placement/screenshots/`, including `placed-tower.png`. `git diff --check` produced no output. The ownership check produced no output.

Line-count review found one hand-written code file above the 600-line preference: `internal/game/game_test.go` at 602 lines. This overage was pre-existing from recent plans and this implementation avoided increasing it by adding new placement tests to `internal/game/build_placement_test.go`. Recommended follow-up is a responsibility-based split of `game_test.go`, such as moving remaining general state, camera, or map tests into focused files. No extra refactor was performed because that split is outside this feature scope and was not approved as part of this plan.

## Context and Orientation

The Go module is a local Ebitengine tower-defense prototype. `cmd/td/main.go` polls Ebitengine input and passes a `game.Input` value to `internal/game.State.Update`. `internal/game/game.go` owns top-level game state and draw order. `internal/game/building_bar.go` owns the left building bar, tower icon bounds, cost rendering, and affordability hover. `internal/game/map.go` stores the 15x15 home Plot as Tiles, where each Tile has a `Terrain` and a `Feature`. A feature is a placed object such as `featureSanctum`, `featureBowTower`, or `featureFlameBoltTower`; `featureNone` means the Tile has no feature. `internal/game/combat.go` scans placed tower features during Raids, so changing a Tile feature is enough for newly built towers to fight later.

The current starting resources are fixed in `internal/game/hud.go`: 80 Wood, 45 Stone, and 12 Metal. The Bow Tower costs 30 Wood, 10 Stone, and 10 Metal. The Flame Bolt Tower costs 30 Stone and 20 Metal. With those defaults, the Bow Tower is affordable and the Flame Bolt Tower is not.

Root control documents constrain the work. `GAME.md` currently says the building bar is visual-only and that build rules are open; this change resolves the first placement rule as calm-only, empty-feature, grass-like Tile placement. `PRODUCT.md` and `README.md` must describe the implemented user-visible behavior. `DESIGN.md` must no longer warn that the building bar is visual-only and should mention compact drag feedback. `ARCHITECTURE.md` must describe building-bar ownership of drag state, placement validation, and resource spending. `ROADMAP.md`, `ART.md`, and `CODESTYLE.md` do not need updates unless implementation reveals a durable change to their guidance.

## Implementation Plan

First, extend input. In `internal/game/game.go`, add `MouseDown` and `Released` boolean fields to `Input` while keeping `Clicked` for existing just-pressed click behavior. In `cmd/td/main.go`, set `MouseDown` from `ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)`, set `Clicked` from `inpututil.IsMouseButtonJustPressed`, and set `Released` from `inpututil.IsMouseButtonJustReleased`. Existing selection and button code can keep using `Clicked`; new placement code should use the held and released states.

Second, add build-drag state in `internal/game`. Add a small private struct to `gameUI` or `State` that records whether a build drag is active and which building-bar item index is being dragged. Keep this state private to the game package. Add helpers to resolve the dragged item, convert a building-bar item to a `tileFeature`, subtract a `ResourceCost`, check whether the game can place during calm play, and check whether a Tile can receive a tower.

Third, update building-bar input flow. At the start of `State.Update`, after overlay handling and before camera, selection, or Raid control handling, update hover and then update build dragging. When the mouse is just pressed over a building-bar icon and `canAffordBuildingCost()` is true, start dragging that icon. If the icon is not affordable, do nothing so it cannot be dragged off the bar. While dragging, normal selection and Raid button handling should not consume that mouse input. When the mouse is released, attempt placement at the cursor and always clear the drag state after the release.

Fourth, implement placement validation and mutation. Convert the release cursor to a home Plot Tile by using the current camera projection and the Tile rectangles already used by selection hit testing. A release can place only if the target Tile is inside the home Plot, has `Terrain == terrainEmpty`, has `Feature == featureNone`, the game is calm and not breached, and resources still cover the dragged tower cost. On success, subtract resources, set the Tile feature, and leave selection unchanged. On failure, do not spend resources and do not change the map.

Fifth, render the dragged icon. In `State.Draw`, draw the building bar as today, then draw the dragged icon before higher-priority panels and overlays. The dragged icon should be screen-space, centered at the cursor, and half the normal building-bar icon size. If the sprite is nil, draw nothing for the dragged copy but keep the drag logic safe.

Sixth, add tests. Place focused tests in `internal/game/building_bar_test.go` or a new `internal/game/build_placement_test.go` so `internal/game/game_test.go` does not grow further. Tests should use existing projected Tile helpers where possible and should directly construct `Input` frames for press, hold, and release. Cover affordable Bow Tower drag start, unaffordable Flame Bolt Tower drag blocking, successful Bow Tower placement and cost deduction, occupied Tile drop, road drop, forest drop, active-Raid drop, paused calm drop, and invalid-drop drag cleanup.

Seventh, update visual evidence and docs. Change the screenshot capture base path in `cmd/td/main_test.go` to `plans/28-tower-drag-placement/screenshots/` and add a screenshot target that places a Bow Tower before capture. Update `README.md`, `PRODUCT.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` so they describe current placement behavior without implying broader base-building systems exist.

## Validation

Run these commands from `/home/dave/dev/ai/td`:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.(go)$' | xargs -r wc -l | sort -n

Acceptance is that `go test ./...` passes, the new placement tests prove affordable towers can be dragged and built while unaffordable towers cannot be dragged, `plans/28-tower-drag-placement/screenshots/` contains visual evidence of a newly placed tower and updated resources, `git diff --check` reports no whitespace errors, ownership output is empty, and this plan records the line-count review.

Manual acceptance is: start the game, drag the Bow Tower icon from the left building bar to an empty grass-like Tile, release, and observe the new tower on the Tile and resources reduced from 80/45/12 to 50/35/2. Try to drag the Flame Bolt Tower with default resources and observe that no dragged icon leaves the bar. Start a Raid and observe that releasing a dragged tower does not place during the active Raid.

## Idempotence and Recovery

The code edits are local and additive. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/28-tower-drag-placement/screenshots/`. If placement tests fail, inspect the order in `State.Update`: build-drag release must be handled before selection and Raid controls can consume the same mouse release. If newly placed towers do not fire during Raids, inspect whether the Tile feature was set to the existing `featureBowTower` or `featureFlameBoltTower` constants.
