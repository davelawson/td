# Show the Selected Tower's Attack Range

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept current while work proceeds. Maintain this document in accordance with `PLANS.md` from the repository root.

## Purpose / Big Picture

After this change, selecting a Bow Tower, Flame Bolt Tower, or Catapult Tower shows its maximum attack coverage directly on the map. A faint gold circle fills the covered area and an opaque gold ring marks the exact center-to-center range boundary, so the player can compare tower placement with enemy paths without obscuring the field.

The result is visible by starting a game, building or arranging a combat tower, and selecting it. The circle follows camera pan and zoom, remains visible while the selection remains active in any phase or pause state, and disappears when another kind of object is selected. Automated geometry and selection tests plus a plan-local screenshot provide evidence.

## Progress

- [x] (2026-07-14 23:48Z) Inspected selection, structure templates, combat range checks, camera projection, scene draw order, visual tests, control documents, dirty worktree, baseline tests, and current line counts.
- [x] (2026-07-14 23:48Z) Confirmed the user-selected visual treatment: an outline plus a faint tint.
- [x] (2026-07-14 23:54Z) Added template-backed selected-tower range geometry, camera-projected rendering, centralized colors, and focused tests.
- [x] (2026-07-14 23:54Z) Synchronized current product, game-design, visual-design, and architecture documents without disturbing the in-progress Forest-biome changes.
- [x] (2026-07-14 23:54Z) Captured and inspected focused 1920x1080 Catapult range evidence with an in-range Ghoul and out-of-range Armoured Skeleton.
- [x] (2026-07-14 23:54Z) Passed formatting, normal, race, 50-run focused, screenshot, whitespace, status, and ownership validation.
- [x] (2026-07-14 23:54Z) Checked hand-written code-file line counts. No file exceeds 600 lines; the pre-existing placement and Raid tests remain the only near-limit files at 564 and 541 lines.

## Surprises & Discoveries

- Observation: The repository already uses the same tower-template range value for combat and the selection panel.
  Evidence: `internal/game/combat.go` passes `StructureTemplate.RangeTiles` to center-to-center target acquisition, while `internal/game/selection_panel.go` presents that field to the player.

- Observation: Drawing the range from inside one Tile render would let later Tiles cover it.
  Evidence: `internal/game/scene.go` traverses and fully draws each Tile in order, so the indicator must be drawn after all explored Tiles but before exploration controls and later scene layers.

- Observation: The baseline worktree contains uncommitted Forest-biome work, including overlapping control documents and `internal/game/visual_test.go`.
  Evidence: `git status --short` lists those files, while the tower-range source files are otherwise unchanged; edits must be additive and preserve those changes.

- Observation: `color.RGBA` made full Light Bronze channels with alpha 38 render far more strongly than the intended fifteen-percent tint because its channels are alpha-premultiplied.
  Evidence: The first 1920x1080 capture showed an opaque-looking gold interior; switching only the fill to `color.NRGBA` preserves the hue while correctly applying alpha 38.

- Observation: One isolated screenshot rerun transiently omitted most screen-space UI even though the range itself rendered correctly.
  Evidence: The single recovery rerun produced the complete 1920x1080 HUD, selection panel, scene, raiders, and range overlay, matching prior repository evidence about fresh Ebitengine capture processes.

## Decision Log

- Decision: Use a translucent Light Bronze fill and a three-pixel opaque Light Bronze outline.
  Rationale: The user chose an outline plus tint, and Light Bronze is the established gold selection accent for terrain and pause presentation.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Measure the circle from the selected tower Tile center using the tower template's exact `RangeTiles` value.
  Rationale: Combat already compares the Euclidean distance between tower and enemy centers, so this makes the visual boundary truthful rather than approximated from sprite edges.
  Date/Author: 2026-07-14 / Codex

- Decision: Show the indicator whenever a combat tower remains selected, independent of phase, pause, targets, or staffing state.
  Rationale: Selection is already phase-independent and the indicator communicates a stable authored tower fact rather than current firing eligibility.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

The feature is complete. Selecting any placed Bow, Flame Bolt, or Catapult Tower now draws a camera-projected Light Bronze circle centered on the tower Tile. Its radius comes directly from the same `StructureTemplate.RangeTiles` value used by combat, its alpha-38 `color.NRGBA` fill keeps covered Tiles readable, and its three-pixel opaque edge gives the exact maximum center-to-center boundary. The overlay follows pan and zoom and remains available through Labour, Management, Raid, pause, and breach while the tower remains selected. Non-towers, missing Plots, invalid Tiles, and nonpositive ranges safely draw nothing.

Focused tests cover all three authored radii, independent camera projection math, non-home explored Plots, phase and pause independence, and rejection cases. `go test ./...`, `go test -race ./...`, and fifty focused runs pass. Full-resolution evidence shows a selected Catapult Tower, readable selection facts, an in-range Ghoul, an out-of-range Armoured Skeleton, and unobscured HUD and map context. `git diff --check` and the ownership scan are clean.

No hand-written Go file exceeds 600 lines. The largest remain the pre-existing `internal/game/build_placement_test.go` at 564 lines and `internal/game/raid_test.go` at 541 lines; future tests for those responsibilities should move into dedicated files rather than extend them. This slice correctly placed its 158-line focused coverage in `internal/game/tower_range_test.go`, so no extra refactor needs approval.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `internal/game/selection.go` stores a selected structure as a `tileCoordinate`, which includes the Plot and local Tile indices. `internal/game/structures.go` owns `StructureTemplate.RangeTiles` for the three combat towers. `internal/game/combat.go` treats that range as a circle in world-space Tile units around `plotTileWorldCenter`; an enemy is eligible when its center lies at or inside the radius.

`internal/game/camera.go` projects world-space Tiles onto the scene viewport using `plotBaseTileSize * camera.zoom`. `internal/game/scene.go` draws explored terrain and structures. The range must be projected using this same camera scale, drawn after all explored Tiles so it is not overwritten, and drawn before magnifying-glass exploration controls, raiders, projectiles, and screen-space UI so those elements stay readable.

`PRODUCT.md` describes current user-visible selection capability. `GAME.md` records tower ranges and currently leaves later range indicators open. `DESIGN.md` defines gold outlines as the terrain-selection language and requires readable tower-defense coverage. `ARCHITECTURE.md` summarizes selection and scene rendering ownership. This feature changes those current descriptions, but does not change `ROADMAP.md`, `ART.md`, `CODESTYLE.md`, or `README.md` because it adds no strategic scope, generated art, convention, dependency, or run instruction.

## Plan of Work

First, add a private `towerRangeIndicator` projected-geometry type and a helper in `internal/game/scene.go`. The helper will require a structure selection, validate its Plot and Tile, resolve the selected Tile's feature through the existing `towerTemplate`, reject non-towers or nonpositive ranges, project the Tile center with the current viewport and camera, and return screen-space center coordinates plus a radius of `RangeTiles * plotBaseTileSize * camera.zoom`.

Add a draw helper that fills that circle with non-premultiplied Light Bronze at alpha 38, approximately fifteen percent opacity, and strokes it with opaque Light Bronze at three pixels. Call it in `drawExploredPlots` after every Plot Tile has rendered and before exploration buttons. Natural screen and later UI drawing clips or covers portions outside their appropriate regions; a selected tower whose center is off-screen may therefore leave a truthful partial arc visible.

Add `internal/game/tower_range_test.go` with pure tests for exact Bow, Flame Bolt, and Catapult radii; camera zoom and pan projection; a tower in a non-home Plot; and rejection of no selection, raiders, terrain, ordinary buildings, stale Tiles, and invalid coordinates. Keep these tests out of the existing near-limit placement and Raid test files. Extend `internal/game/visual_test.go` with one opt-in Catapult screenshot arranged on clear terrain and selected by state, saving to `plans/62-selected-tower-range/screenshots/selected-tower-range.png`.

Update `PRODUCT.md` to describe the selected-tower overlay. Update `GAME.md` to make selected-tower range visualization a current decision while retaining placement preview as an open question. Update `DESIGN.md` with the faint-gold-fill and exact-ring visual rule. Update the current selection-rendering sentence in `ARCHITECTURE.md`. Merge only these focused statements into the existing uncommitted Forest changes.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`. The baseline already passed:

    go test ./...

Implement the private geometry, draw path, focused tests, screenshot capture, and documentation with narrow patches. Format only changed Go files:

    gofmt -w internal/game/scene.go internal/game/tower_range_test.go internal/game/visual_test.go

Validate normal behavior, concurrency safety, and focused stability:

    go test ./...
    go test -race ./...
    go test ./internal/game -run 'SelectedTowerRange|TowerRangeIndicator' -count=50

Capture the post-implementation visual evidence and inspect it at full resolution:

    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureSelectedTowerRangeScreenshot -count=1

Then run repository checks:

    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'

End with the required hand-written code line-count review:

    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 25

Report every file near or above 600 lines with its count and a responsibility-based recommendation. Do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

Automated acceptance requires exact projected radii corresponding to Bow 3.0 Tiles, Flame Bolt 2.5 Tiles, and Catapult 5.0 Tiles at multiple camera positions and zoom levels. It must prove the center is the selected Tile's projected world center, explored non-home Plot coordinates work, and every non-tower or stale selection suppresses the indicator. All repository tests and race tests must pass, and fifty focused runs must remain stable.

Visual acceptance requires `plans/62-selected-tower-range/screenshots/selected-tower-range.png` at 1920x1080 to show a selected Catapult Tower, its selection panel, a readable faint gold covered area, and an unambiguous gold maximum-range boundary without hiding Tiles, structures, or UI. Inspect the PNG at full resolution.

Documentation acceptance requires `PRODUCT.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` to agree that only selected built combat towers show the overlay, it uses the combat template's exact range, and placement preview remains outside this slice. Existing Forest-biome edits must remain intact.

## Idempotence and Recovery

Formatting, tests, screenshot capture, and validation commands are safe to rerun. Screenshot capture overwrites only this plan's focused evidence. If a test reveals invalid selected coordinates, suppress the indicator instead of indexing map state unsafely. If visual inspection shows the tint is too strong, adjust only the plan-recorded alpha while retaining the accepted outline-plus-tint treatment and document the evidence. Do not discard or rewrite unrelated dirty-worktree changes.

## Artifacts and Notes

The focused evidence belongs at:

    plans/62-selected-tower-range/screenshots/selected-tower-range.png

The pre-implementation baseline and final `go test ./...` passed. Final validation also passed `go test -race ./...` and fifty focused runs. The accepted screenshot is a 1920x1080 RGB PNG. Final line counts show no hand-written Go file above 600 lines; `internal/game/build_placement_test.go` at 564 and `internal/game/raid_test.go` at 541 remain the closest and did not absorb these tests.

## Interfaces and Dependencies

No exported interface changes. Add a private geometry representation and helper in `internal/game/scene.go` equivalent to:

    type towerRangeIndicator struct {
        centerX float32
        centerY float32
        radius  float32
    }

    func (s *State) selectedTowerRangeIndicator(viewport sceneViewport) (towerRangeIndicator, bool)

Use Ebitengine's existing `vector.FillCircle` and `vector.StrokeCircle`; do not add libraries, assets, serialization, configuration, input actions, or gameplay state.

Revision note (2026-07-14): Created from the accepted implementation plan after repository, workflow, baseline-test, dirty-worktree, and line-count inspection. Updated after implementation with completed progress, the corrected non-premultiplied tint decision, screenshot-capture evidence, validation results, and final line-count findings.
