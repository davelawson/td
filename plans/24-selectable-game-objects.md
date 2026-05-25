# Add Selectable Game Objects

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds. This plan follows `PLANS.md`.

## Purpose / Big Picture

After this change, the player can click visible game objects in the first home Plot scene. A left click selects a raider when the click hits the raider's visible sprite, otherwise selects any structure tile such as the Sanctum, Bow Tower, or Flame Bolt Tower. The selected object is drawn brighter, and clicking elsewhere clears selection. This gives the prototype its first general object-selection interaction without adding inspection panels, tower commands, placement, upgrades, or multi-select.

## Progress

- [x] (2026-05-25T22:45:35Z) Inspected the current game state, camera projection, map features, raider IDs, rendering paths, click input, and file line counts.
- [x] (2026-05-25T22:45:35Z) Chose raider sprite bounds for raider hit testing, full structure tiles for structure hit testing, raider-first priority, and pause-compatible selection.
- [x] (2026-05-25T22:50:33Z) Added selected-item state, hit testing, stale raider cleanup, and bright selected rendering in `internal/game`.
- [x] (2026-05-25T22:50:33Z) Added focused selection tests outside `internal/game/game_test.go`.
- [x] (2026-05-25T22:50:33Z) Updated current-state, roadmap, design, game, and architecture documents.
- [x] (2026-05-25T22:50:33Z) Ran validation commands and captured findings.
- [x] (2026-05-25T22:50:33Z) Checked hand-written code-file line counts and reported files over the 600-line preference.

## Surprises & Discoveries

- Observation: `internal/game/game_test.go` already has 605 lines before this work.
  Evidence: `wc -l internal/game/*.go cmd/td/main.go` reported `605 internal/game/game_test.go`.

## Decision Log

- Decision: Store selection in `internal/game.State` as a private selected-item value with variants for no selection, structure tile, and raider ID.
  Rationale: The game package already owns map state, raider state, camera projection, and rendering, so selection does not need to cross the executable boundary.
  Date/Author: 2026-05-25 / Codex.

- Decision: Use visible raider sprite bounds for raider clicks, full projected structure tile bounds for structure clicks, and raider-first priority.
  Rationale: This matches what the player sees, keeps structures easy to click, and gives moving threats priority when objects overlap visually.
  Date/Author: 2026-05-25 / Codex, after user clarification.

- Decision: Selection works while SPACE-paused but remains blocked by the ESC overlay.
  Rationale: The current camera also works while paused for inspection, while the overlay intentionally blocks game-scene input.
  Date/Author: 2026-05-25 / Codex, after user clarification.

## Outcomes & Retrospective

Implemented selectable game objects as planned. The game state now stores a private selected-item value for structures and raiders, left-click hit testing selects visible raiders before structure tiles, clicks elsewhere clear selection, pause still allows selection, the overlay blocks selection, and stale raider selection clears when the selected raider leaves active state. Selected sprites render brighter.

Validation passed on 2026-05-25:

    go test ./...
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The ownership check printed nothing. `git status --short` showed only the intended modified and new files for this change.

Line-count review on 2026-05-25:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Only `internal/game/game_test.go` exceeds the 600-line preference at 605 lines. This was pre-existing before the implementation. The selection tests were placed in `internal/game/selection_test.go` to avoid increasing it. Recommended follow-up, if the project wants to address the existing overage, is to move camera, HUD, or overlay tests out of `game_test.go` into responsibility-focused test files. No extra refactor was performed because that split was not part of this feature scope.

## Context and Orientation

The repository is a Go/Ebitengine tower-defense prototype. `cmd/td/main.go` owns Ebitengine startup and input polling, then passes a `game.Input` value into `internal/game.State.Update`. `internal/game/game.go` owns the active game state. `internal/game/map.go` defines a 15x15 home Plot with a centered Sanctum and two starting towers. `internal/game/raid.go` owns active raiders with stable integer IDs during a Raid. `internal/game/scene.go` draws map tiles and structures through camera projection, and `internal/game/raidui.go` draws raiders, projectiles, and the `Next Raid` button.

The user-visible behavior affects `README.md` and `PRODUCT.md`. The completed current-phase interaction affects `ROADMAP.md`. The intended game interaction affects `GAME.md`. The brighter selected-state visual treatment affects `DESIGN.md`. The ownership of selected state affects `ARCHITECTURE.md`. `CODESTYLE.md` requires `gofmt`, doc comments for Go functions and methods, tests for pure behavior, and a final line-count review for hand-written code files. No art assets or dependencies are needed.

## Plan of Work

Add a private selection model in `internal/game`, with one selected-item kind for no selection, one for a structure tile coordinate, and one for a raider ID. Extend `State` with a private `selection` field. Implement click handling inside `State.Update` after overlay handling and camera input, before pause early-return and before Raid simulation. A click inside the `Next Raid` button should remain UI input and should not clear or change object selection.

Implement hit testing using existing projection math. Raider hit testing should compute each active raider's current projected sprite rectangle using the same size and camera viewport used by `drawRaidEnemy`; iterate from the end of the slice toward the start so later-drawn raiders win when overlapping. Structure hit testing should scan the home Plot for tiles whose `Feature` is not `featureNone`, project the full tile rectangle, and select the matching tile coordinate. Raider hit testing runs before structure hit testing.

Update rendering so selected structures and selected raiders draw brighter. For sprite-backed objects, apply an Ebitengine color scale on the `DrawImageOptions`. For fallback shape rendering, use a brighter existing palette color. Clear selected raider state whenever the selected raider ID is no longer present among active raiders.

Add focused tests in a new `internal/game/selection_test.go` file. Use helper functions that project tile or raider rectangles to screen coordinates and click their centers. Do not add tests to `internal/game/game_test.go` because it already exceeds the 600-line preference.

Update the control documents named above so current product behavior, current roadmap status, intended interaction design, visual treatment, and architecture match the implementation.

## Concrete Steps

From `/home/dave/dev/ai/td`, edit code with `apply_patch` and run `gofmt` on changed Go files. Run:

    go test ./...
    git diff --check
    git status --short
    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Report any hand-written code file over 600 lines with a concrete recommendation. Do not perform unplanned file splits or refactors without user approval.

## Validation and Acceptance

`go test ./...` should pass. New tests should prove that the Sanctum, Bow Tower, Flame Bolt Tower, and active raiders can be selected; that raider priority wins over structure priority; that empty-space clicks clear selection; that SPACE-paused selection works; that the ESC overlay blocks selection; that `Next Raid` button clicks keep their existing behavior; and that stale raider selections are cleared when the raider is removed.

Manual acceptance is: start the game, click the Sanctum or either tower and see it brighten; start a Raid, click a visible raider and see it brighten; click empty map space and see selection clear.

## Idempotence and Recovery

The code edits are local and additive. Re-running `gofmt`, tests, whitespace checks, and line-count checks is safe. If selection hit tests fail, compare the helper rectangles with `drawStructureSprite` and `drawRaidEnemy` projection logic before changing user-facing behavior. If `Next Raid` behavior regresses, inspect click handling order before changing Raid state.

## Artifacts and Notes

Validation transcripts and final line-count findings are recorded in `Outcomes & Retrospective`.

Revision note, 2026-05-25 / Codex: Completed the implementation and updated progress, validation results, and line-count review.

## Interfaces and Dependencies

No new external dependencies are required. The implementation stays in `internal/game` and uses existing Ebitengine draw options.

The final code should include private helpers equivalent to:

    type selectedItemKind int
    type selectedItem struct {
        kind selectedItemKind
        tile tileCoordinate
        raiderID int
    }

    func (s *State) updateSelection(input Input)
    func (s *State) selectedStructureAt(x, y int) (tileCoordinate, bool)
    func (s *State) selectedRaiderAt(x, y int) (int, bool)
    func (s *State) clearMissingSelectedRaider()

Revision note, 2026-05-25 / Codex: Created this plan from the accepted implementation plan and user clarification that all structures and raiders should be selectable.
