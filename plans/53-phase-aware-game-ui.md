# Make the Game UI Phase-Aware and Inspectable

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/53-phase-aware-game-ui.md`.

## Purpose / Big Picture

After this change, the game gives the player UI appropriate to the current activity. Management shows construction controls and a live preview of the next Raid challenge. A Raid removes the construction bar and keeps the generated Raid's challenge visible. The player can also click a Tree or Boulder to inspect its terrain type and Plot biome, with a gold Tile outline connecting the panel to the map.

The result is visible by starting a game, observing `Management phase | Challenge 4.0`, starting a Raid, and seeing `Enemies remaining: N | Challenge 4.0` without the building bar. Clicking a Tree or Boulder shows a compact terrain panel; clicking empty grass or Road does not select it.

## Progress

- [x] (2026-07-14 14:40Z) Inspected phase transitions, building UI input/rendering, selection priority, terrain and biome storage, challenge generation, HUD layout, tests, docs, and screenshot tooling.
- [x] (2026-07-14 14:40Z) Confirmed product choices: only Trees and Boulders are selectable, terrain panels show terrain plus biome, selected terrain uses a gold Tile outline, and challenge is appended to centered phase text with one decimal place.
- [x] (2026-07-14 14:41Z) Preserved accepted Management and Raid baseline screenshots under this plan directory.
- [x] (2026-07-14 14:45Z) Implemented phase-aware building-bar rendering and input ownership with focused tests.
- [x] (2026-07-14 14:45Z) Implemented Tree and Boulder selection, panel data, and selected-Tile rendering with focused tests.
- [x] (2026-07-14 14:45Z) Implemented live next-Raid and frozen current-Raid challenge text with focused tests.
- [x] (2026-07-14 14:48Z) Updated current-state, gameplay, design, architecture, and onboarding documents.
- [x] (2026-07-14 14:50Z) Captured and reviewed updated Management, Raid, and selected-terrain screenshots.
- [x] (2026-07-14 14:51Z) Passed full validation, ownership checks, whitespace checks, and 20-run focused repeated tests.
- [x] (2026-07-14 14:51Z) Checked hand-written Go file line counts; no file exceeds 600 lines, and the unchanged 595-line building-bar test remains the only near-threshold file.

## Surprises & Discoveries

- Observation: The building bar currently draws during every phase and `buildingBarContains` blocks the entire left strip even during a Raid.
  Evidence: `State.Draw` calls `drawBuildingBar` unconditionally, and selection, camera drag, exploration, and build-drop blocking all call the unconditional bounds check.
- Observation: The active Raid already retains the exact generated challenge template needed for a stable Raid display.
  Evidence: `raidState.template.challengeRating` is set once by `startNextRaid`, while `generateRaid` can calculate the upcoming value from live Management state.
- Observation: `internal/game/building_bar_test.go` is already close to the source-size preference.
  Evidence: the baseline line-count check reports 595 lines, so new phase visibility tests belong in a focused new test file.
- Observation: Existing accepted screenshots provide a baseline for both Management and active Raid UI.
  Evidence: `plans/52-dynamic-raid-generation/screenshots/` contains `running-game.png` and `active-raid.png`.
- Observation: The full screenshot suite reproduced the known intermittent partial-frame artifact for the longer active-Raid fixture.
  Evidence: two full-suite captures omitted parts of the top HUD; a focused active-Raid capture running the same game state through its own Ebitengine loop produced the complete `plans/53-phase-aware-game-ui/screenshots/active-raid.png` used for review.

## Decision Log

- Decision: Treat Management as the visible peaceful construction phase and do not render construction UI during Raid or breach.
  Rationale: Labour is instantaneous and accepts no construction input; Management is the only observable phase in which the building bar communicates a usable action.
  Date/Author: 2026-07-14 / User and Codex
- Decision: When hidden, the building bar also stops owning pointer input and clears hover state.
  Rationale: An invisible UI surface must not block map selection or right-drag camera input.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Select only Tree and Boulder terrain, after raider and structure hit testing; empty grass and Road continue clearing selection.
  Rationale: The user explicitly excluded unoccupied grass and Road Tiles, while preserving current object priority avoids regressions.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Show terrain type and Plot biome and mark the selected Tile with a gold outline.
  Rationale: These two facts are useful player-facing context without exposing debug coordinates, and a Tile outline works for both authored terrain sprites.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Append challenge to centered phase text with one decimal place.
  Rationale: The phase area has the clearest semantic relationship and avoids increasing density in the resource/population group.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Preview the next generated challenge from live Management state and display the stored template during Raid and breach.
  Rationale: Management choices should immediately reveal their effect on the upcoming Raid, while an already-generated Raid must not change retrospectively.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Implemented phase-aware construction UI, natural-obstacle inspection, and challenge visibility. Management retains its full building workflow and now previews the next challenge. Raid and breach hide the building bar, release its former input strip, and display the stored challenge. Tree and Boulder selection works after raider and structure priority, shows terrain and biome facts, and draws a gold Tile outline; grass and Road remain ordinary deselection clicks.

Visual review passed. `running-game.png` shows the complete Management bar and `Management phase | Challenge 4.0` without overlap. The focused `active-raid.png` shows the full HUD, `Enemies remaining: 3 | Challenge 4.0`, exposed map space, and no construction bar. `selected-terrain.png` shows a readable two-row `Tree`/`Grasslands` panel and a visible gold outline without obscuring the tree sprite.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, 20 repeated runs of the focused building-bar, terrain, challenge, and selection tests, `git diff --check`, and the ownership check all succeeded. No hand-written Go file exceeds 600 lines. The unchanged `internal/game/building_bar_test.go` remains closest at 595 lines; before its next substantive expansion, move a coherent group such as tab/hover behavior into a focused test file. New phase tests were deliberately placed in `internal/game/phase_ui_test.go`, so no unplanned split was needed.

## Context and Orientation

This is a local Go/Ebitengine game. `internal/game/game.go` owns update and draw order. `internal/game/building_bar.go` owns construction-bar geometry, interaction, and drawing, while `internal/game/building_tooltip.go` owns its hover tooltip. Management is represented by `phaseManagement`; `phaseLabour` resolves immediately, and `phaseRaid` represents an active or breached assault.

`internal/game/selection.go` stores a private selected-item kind and currently prioritizes raiders over structures. Every explored map Tile stores a `tileTerrain` and `tileFeature` in `internal/game/map.go`; only `terrainTree` and `terrainBoulder` become new selectable subjects. `internal/game/selection_panel.go` adapts game-owned facts to `internal/ui.SelectionPanelData`, while `internal/ui/selection_panel.go` owns row labels, formatting, bounds, hit testing, and rendering.

`internal/game/raid_generator.go` owns the challenge formula. `internal/game/raid.go` stores the generated template in active Raid state. `internal/game/hud.go` owns top-bar text and rendering. The Management preview must call the same generator used to start a Raid with `raid.number + 1`, total settlement population, and explored Plot count. Raid and breach display must use `raid.template.challengeRating`.

`DESIGN.md` requires compact readable HUD text, plainly visible selected states, screen-space building UI, and screenshot evidence for rendered changes. `PRODUCT.md` owns current user-visible capability truth; `GAME.md` owns selection and challenge design decisions; `ARCHITECTURE.md` owns input priority and package boundaries; `README.md` summarizes the runnable workflow. `CODESTYLE.md` requires doc comments, `gofmt`, tests for pure behavior, and a strong preference for hand-written files below 600 lines.

## Plan of Work

First, preserve baseline copies of the existing Management and active-Raid screenshots under `plans/53-phase-aware-game-ui/screenshots/before/`.

Add a private `buildingBarVisible` method in `internal/game/building_bar.go` whose true state is Management without an active or breached Raid. Use it consistently for drawing, hover updates, tooltip lookup, drag start/rendering, and `buildingBarContains`. Hidden state clears item and tab hover values. Existing placement eligibility remains authoritative, and no unrelated control is repositioned. Add a new focused phase-UI test file rather than expanding the 595-line building-bar test.

Extend selected-item state with a terrain kind and Tile coordinate. After raider and structure hit testing, resolve the explored Tile under the pointer and select it only when it has no structure feature and its terrain is Tree or Boulder. Clicking grass, Road, or outside explored Tiles clears selection. Add `selectedTerrain` for drawing and render a three-pixel gold outline after Tile contents so it remains visible over Tree and Boulder sprites. Existing raider and structure brightness behavior stays unchanged.

Extend the internal UI panel model with a terrain kind plus terrain and biome names. The game adapter validates the stored Plot and Tile, maps Tree/Boulder to `Tree`/`Boulder`, maps the Plot biome through the existing `Grasslands`/`Hills` labels, and returns two rows labeled `Terrain` and `Biome`. Selection-panel clicks remain blocked as UI.

Add challenge helpers in `internal/game/hud.go`: one returns the displayed numeric rating and another formats `Challenge %.1f`. In `phaseRaid`, including breach, read the frozen Raid template. In Management or Labour, call `generateRaid` for the next Raid from live population and explored-Plot inputs. Append the formatted challenge to `phaseText` with ` | ` in the centered top-bar string. Do not change Raid generation, thresholds, duration, or state.

Update tests around selection state, game-to-UI adaptation, UI row formatting, bar visibility/input ownership, and challenge formatting. Update `cmd/td/main_test.go` to save evidence under plan 53 and add a selected-terrain target. Capture Management, Raid, and selected-terrain output and verify text, bar visibility, panel readability, outline contrast, and absence of overlap.

Update `README.md`, `PRODUCT.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md`. Record the durable gameplay decisions in `GAME.md`. Do not update `ROADMAP.md`, `ART.md`, or `CODESTYLE.md`, because this slice does not change future strategy, generated-art guidance, or source conventions.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Before production-code edits, create the evidence directory and copy the accepted baseline:

    mkdir -p plans/53-phase-aware-game-ui/screenshots/before
    cp plans/52-dynamic-raid-generation/screenshots/running-game.png plans/53-phase-aware-game-ui/screenshots/before/management.png
    cp plans/52-dynamic-raid-generation/screenshots/active-raid.png plans/53-phase-aware-game-ui/screenshots/before/raid.png

Implement focused code and tests, format changed Go files, then run:

    gofmt -w <changed Go files>
    go test ./internal/ui ./internal/game
    go test ./...

Update documents and screenshot fixtures, then capture evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureSelectedTerrainScreenshot -count=1
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureActiveRaidScreenshot -count=1

Review the plan-53 Management, active-Raid, and selected-terrain PNGs. Record findings in this plan before final acceptance.

Run final validation:

    go test ./...
    go test ./internal/game -run 'Test(BuildingBar|Terrain|Challenge|Selection)' -count=20
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'

End with the hand-written Go line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Report any file above 600 lines, or close enough that the next likely change would exceed it, with a concrete split recommendation. Do not implement an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` must pass. Focused tests must prove that Management owns the building bar, Raid and breach do not render or reserve its screen region, and hidden hover/tooltip state is cleared. Right-drag and left-click map interactions must be able to start in the former bar area during Raid.

Tree and Boulder clicks must create terrain selections in home and additional explored Plots. Their UI rows must be exactly `Terrain: Tree|Boulder` and `Biome: Grasslands|Hills`. Empty grass and Road must not remain selected. Raider and structure overlap tests must preserve their higher priority. The selected Tile outline must be visible in screenshot evidence.

The baseline new game must format `Management phase | Challenge 4.0`. Management changes to total population or explored Plot count must change the preview according to the existing formula. Starting the Raid must freeze that value, and later state mutation must not change its display while phase is Raid or breached. Enemy counts and challenge generation tests must remain unchanged.

The updated control documents must consistently describe Management-only construction UI, Tree/Boulder inspection, selection priority, terrain panel contents, gold outline, live next-Raid preview, and frozen current-Raid challenge without implying grass/Road selection or new terrain gameplay actions.

## Idempotence and Recovery

All code and documentation edits are local. Formatting, tests, screenshot capture, and checks can be repeated. Screenshot capture overwrites only plan-53 evidence. Preserve unrelated working-tree changes; if a changed hunk overlaps existing work, integrate rather than discard it. No data migration, destructive command, dependency update, or rollback procedure is needed.

## Artifacts and Notes

Before evidence lives under `plans/53-phase-aware-game-ui/screenshots/before/`. The accepted updated evidence is `running-game.png`, `active-raid.png`, and `selected-terrain.png` beside it. The remaining PNGs are the refreshed existing visual-regression suite.

## Interfaces and Dependencies

No new dependency or externally consumable API is required. Private game state gains a terrain selected-item kind and helpers for bar visibility and displayed challenge. The internal `ui.SelectionPanelKind` gains `SelectionPanelTerrain`, and `ui.SelectionPanelData` gains `TerrainName` and `BiomeName` strings. Gameplay state supplies raw facts; `internal/ui` continues to own labels, formatting, bounds, hit testing, and drawing.

Revision note, 2026-07-14 / Codex: Created the plan after repository inspection and user confirmation of terrain scope, highlight treatment, and challenge presentation.

Revision note, 2026-07-14 / Codex: Completed implementation, documentation synchronization, focused screenshot recovery, validation, and line-count review.
