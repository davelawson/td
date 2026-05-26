# Add Selection Panel

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds. This plan follows `PLANS.md`.

## Purpose / Big Picture

After this change, the player can inspect what they selected in the first game screen. Selecting a raider opens a bottom-right panel with its type, current health, maximum health, health percentage, movement speed, and Sanctum damage. Selecting a combat tower opens a bottom-right panel with tower type, range, attack speed, and damage. Selecting the Sanctum opens a small panel naming it. This makes existing click selection useful without adding tower commands, upgrades, placement, resource effects, or multi-select.

## Progress

- [x] (2026-05-25T23:24:00Z) Inspected the current selection, map projection, raider, tower, HUD, screenshot, and test code.
- [x] (2026-05-25T23:24:00Z) Confirmed panel defaults: Sanctum gets a basic name panel, and panel clicks are UI clicks that do not change selection.
- [x] (2026-05-25T23:31:00Z) Added selection-detail formatting, bottom-right panel rendering, debug-counter repositioning, and panel click blocking in `internal/game`.
- [x] (2026-05-25T23:31:00Z) Added focused tests for selected raider rows, tower rows, Sanctum rows, and panel click blocking.
- [x] (2026-05-25T23:35:00Z) Updated current-state, design, game-design, roadmap, and architecture documents to describe the panel.
- [x] (2026-05-25T23:39:00Z) Captured screenshot evidence for selected tower and selected raider states.
- [x] (2026-05-25T23:40:00Z) Ran validation commands and recorded outcomes.
- [x] (2026-05-25T23:40:00Z) Checked hand-written code-file line counts and reported files over the 600-line preference.

## Surprises & Discoveries

- Observation: `plans/24-selectable-game-objects.md` recorded that `internal/game/game_test.go` already exceeded the 600-line preference before this work.
  Evidence: That plan reports `internal/game/game_test.go` at 605 lines and placed selection tests in a separate file.

- Observation: The first selected-raider screenshot showed the long `Skeleton Sword-and-Shield` value overlapping the `Raider Type` label in a two-column panel layout.
  Evidence: Visual review of `plans/25-selection-panel/screenshots/selected-raider.png` after the first screenshot capture showed label/value overlap.

- Observation: A stacked label/value panel layout fixed the long-name overlap while keeping the panel within the bottom-right screen area.
  Evidence: The refreshed `selected-raider.png` screenshot shows each label and value on separate lines without overlap.

## Decision Log

- Decision: Add the selection details panel inside `internal/game`.
  Rationale: The game package already owns selected-object state, raider and tower templates, HUD rendering, and camera-projected input blocking.
  Date/Author: 2026-05-25 / Codex.

- Decision: Show only a basic `Structure: Sanctum` panel for Sanctum selection.
  Rationale: The user requested detailed fields for raiders and towers only, while the existing selection system also selects the Sanctum.
  Date/Author: 2026-05-25 / Codex, after user confirmation.

- Decision: Treat clicks inside the panel as UI clicks that do not clear, change, or pass through selection.
  Rationale: This matches player expectations for an informational UI panel and prevents accidental deselection.
  Date/Author: 2026-05-25 / Codex, after user confirmation.

- Decision: Use stacked label/value rows instead of a two-column row layout.
  Rationale: Raider names can be long enough to collide with labels at the target panel width; stacked rows keep text readable without widening the panel across the playfield.
  Date/Author: 2026-05-25 / Codex.

## Outcomes & Retrospective

Implemented the selected-object detail panel. Selecting a raider now shows type, current health, max health, health percentage, speed, and Sanctum damage. Selecting the Bow Tower or Flame Bolt Tower now shows type, range, attack speed, and damage. Selecting the Sanctum shows a basic Sanctum panel. Clicking inside the panel does not clear or pass through selection, and the debug counter moves above the panel when needed.

Validation passed on 2026-05-25:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The ownership check printed nothing. Screenshot evidence was written under `plans/25-selection-panel/screenshots/`, including:

    selected-tower.png
    selected-raider.png

Line-count review on 2026-05-25:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Only `internal/game/game_test.go` exceeds the 600-line preference at 605 lines. This was pre-existing before the implementation and was also reported in `plans/24-selectable-game-objects.md`. The new panel tests were placed in `internal/game/selection_panel_test.go` to avoid increasing it. Recommended follow-up, if the project wants to address the existing overage, is to move camera, HUD, or overlay tests out of `game_test.go` into responsibility-focused test files. No extra refactor was performed because that split was not part of this feature scope.

Revision note, 2026-05-25 / Codex: Completed implementation, validation, screenshot capture, and line-count review.

## Context and Orientation

The repository is a Go/Ebitengine tower-defense prototype. `cmd/td/main.go` owns Ebitengine startup and input polling, then passes a `game.Input` value into `internal/game.State.Update`. `internal/game/game.go` owns the active game state and draw order. `internal/game/selection.go` owns selected structure and raider state. `internal/game/raid.go` stores active raiders with stable integer IDs, current health, and pointers to `EnemyTemplate` values. `internal/game/enemies.go` defines raider template fields: `Name`, `MaxHealth`, `SpeedTilesPerSecond`, and `SanctumDamage`. `internal/game/structures.go` defines tower template fields: `Name`, `RangeTiles`, `Damage`, and `FireIntervalSeconds`. `internal/game/hud.go` and `internal/game/raidui.go` render existing game UI.

The current selection behavior is already implemented: left-clicking a visible raider selects that raider, otherwise left-clicking a structure tile selects that structure, clicking elsewhere clears selection, selection works while SPACE-paused, and the ESC overlay blocks selection. This plan adds an informational panel for that existing selected object.

The user-visible behavior affects `README.md` and `PRODUCT.md`. The ownership of panel rendering and UI click blocking affects `ARCHITECTURE.md`. `DESIGN.md` already says readable text, stable hit targets, and clear selected states matter; update it only if implementation establishes durable selection-panel visual guidance beyond that existing rule. `GAME.md`, `ROADMAP.md`, `ART.md`, and `CODESTYLE.md` do not need updates unless implementation changes gameplay decisions, product sequencing, art guidance, or source conventions.

## Plan of Work

First, add panel logic in a new `internal/game/selection_panel.go` file. Define small private types for a displayed row and a computed panel. Add a method on `State` that returns the current selected-object panel when one should be visible. It should look up selected raiders by ID in `s.raid.enemies`, selected structures by `s.selection.tile`, and map structure features through `s.structureCatalog`. Raider rows must include type, current health, max health, rounded health percentage, speed in Tiles per second with one decimal place, and Sanctum damage. Tower rows must include tower type, range with one decimal place, attack speed as `every <seconds>s`, and damage. Sanctum rows must include only `Structure: Sanctum`.

Second, add drawing in that same file. The panel should be anchored to the bottom-right of the current drawable size with a fixed margin, use a dark backing, bronze border, parchment text, muted labels or row text, and the existing `s.ui.hudFace` text face. Draw it after Raid controls and the debug counter but before the in-game overlay, so the overlay still darkens and blocks it. Move or adjust the debug counter so it does not overlap the panel.

Third, block selection clicks inside the panel. In `State.updateSelection`, treat panel bounds like the existing `Next Raid` button: if the click is inside a visible selection panel, return without changing selection.

Fourth, add focused tests in a new `internal/game/selection_panel_test.go`. Test selected raider details, Bow Tower details, Flame Bolt Tower details, Sanctum details, and that clicking inside the panel does not clear the selected object. Keep these tests separate from `internal/game/game_test.go` because that file already exceeds the 600-line preference.

Fifth, update docs. `README.md` and `PRODUCT.md` should describe the bottom-right selected-object panel in current product behavior. `ARCHITECTURE.md` should include panel ownership and the invariant that panel clicks block map selection. Update `cmd/td/main_test.go` so screenshot capture writes to `plans/25-selection-panel/screenshots/` and includes selected tower and selected raider screenshots.

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

`go test ./...` should pass. The new tests should prove that selected raider details include type, current health, max health, percentage, speed, and Sanctum damage; selected Bow Tower and Flame Bolt Tower details include type, range, attack speed, and damage; selected Sanctum details include only the basic Sanctum name row; and clicking inside the visible panel does not change selection.

Manual acceptance is: start the game, select the Bow Tower or Flame Bolt Tower, and see the bottom-right panel with tower stats. Start a Raid, select a visible raider, and see the bottom-right panel with raider stats. Select the Sanctum and see a smaller Sanctum panel. Click inside the panel and observe that the current selection remains.

Documentation acceptance is that `README.md`, `PRODUCT.md`, and `ARCHITECTURE.md` describe the implemented behavior without implying tower commands, upgrades, placement, resource effects, saves, or broader gameplay systems exist.

## Idempotence and Recovery

The code edits are local and additive. Re-running `gofmt`, tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/25-selection-panel/screenshots/`. If panel data tests fail, inspect selection-to-template lookup before changing display strings. If selection clearing regresses, inspect the input-blocking order in `updateSelection`.

## Artifacts and Notes

Record validation transcripts, screenshot artifact paths, and final line-count findings in `Outcomes & Retrospective`.

## Interfaces and Dependencies

No new external dependencies are required. The implementation stays in `internal/game` and uses existing Ebitengine drawing APIs plus `td/internal/ui`.

The final code should include private helpers equivalent to:

    type selectionPanelRow struct {
        Label string
        Value string
    }

    type selectionPanel struct {
        Title string
        Rows []selectionPanelRow
    }

    func (s *State) currentSelectionPanel() (selectionPanel, bool)
    func (s *State) selectionPanelBounds() ui.Button[int]
    func (s *State) selectionPanelContains(x, y int) bool
    func (s *State) drawSelectionPanel(screen *ebiten.Image)
