# Starting Bow Tower On Path

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

## Purpose / Big Picture

After this change, every new game starts with one visible Bow Tower on the east side of the home Plot's straight north road. The tower is authored map content that helps the first defense scene look like a tower-defense prototype, but it does not target enemies, deal damage, cost resources, block movement, or introduce player placement.

## Progress

- [x] (2026-05-15T22:20Z) Confirmed the Bow Tower sprite and structure catalog already exist in the working tree.
- [x] (2026-05-15T22:22Z) Added `featureBowTower` and placed it at tile `(8,5)` in the default home Plot.
- [x] (2026-05-15T22:24Z) Rendered the Bow Tower feature using the existing Bow Tower structure template sprite.
- [x] (2026-05-15T22:26Z) Added tests for the authored Bow Tower location and adjacency to the road.
- [x] (2026-05-15T22:28Z) Updated current-state and architecture documents to describe the visible non-combat tower.
- [x] (2026-05-15T22:31Z) Ran validation commands: `go test ./...`, screenshot capture, `git diff --check`, screenshot file inspection, and `git status --short`.
- [x] (2026-05-15T22:32Z) Checked hand-written code-file line counts against the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: The working tree already had uncommitted structure catalog files and Bow Tower catalog tests.
  Evidence: `git status --short` showed modified game files plus untracked `internal/game/structures.go` and `internal/game/structures_test.go`; inspection showed they expose `StructureCatalog.BowTower`.

## Decision Log

- Decision: Place the Bow Tower at tile `(8,5)`, one tile east of the fixed north road.
  Rationale: This makes the tower visibly flank the route enemies use without crowding the Sanctum or changing the road.
  Date/Author: 2026-05-15 / User and Codex

- Decision: Keep this slice visible-only.
  Rationale: Tower targeting, combat damage, costs, build UI, and placement rules are separate gameplay systems that need their own small plans.
  Date/Author: 2026-05-15 / Codex

## Outcomes & Retrospective

Completed. New games now create a default home Plot with a Bow Tower feature at tile `(8,5)`, one tile east of the fixed north road. Rendering draws the Bow Tower from the existing structure template sprite, and the tower is visible in screenshot evidence without changing Raid movement or combat rules.

Validation passed: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and screenshot file inspection all succeeded. The running-game screenshot shows the Bow Tower east of the path.

The hand-written Go line-count review found no files over 600 lines. `internal/game/game_test.go` is 570 lines, which is close enough to the 600-line preference that the next meaningful game-test addition should consider splitting tests by responsibility, such as moving map tests to `internal/game/map_test.go` or camera tests to `internal/game/camera_test.go`. No split was performed because this plan did not include that extra refactor.

## Context and Orientation

`td` is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. The current home Plot is created in `internal/game/map.go` as a fixed 15x15 grid with a centered Sanctum, a straight north road, and a pine-tree border. Rendering lives in `internal/game/scene.go`, which projects map tiles through the camera and draws sprites from loaded assets.

`GAME.md` defines the Bow Tower as the first tower type, but exact cost, range, damage, firing rate, targeting, and upgrades are still open decisions. `PRODUCT.md` describes current user-visible behavior, so it must mention that the tower is visible but non-combat. `ARCHITECTURE.md` says `internal/game` owns prototype map data and rendering, so this change stays inside that package and does not create new packages.

## Plan of Work

Add `featureBowTower` to the private tile feature enum in `internal/game/map.go`. In `newDefaultHomePlotWithTweakSource`, set tile `(8,5)` to `featureBowTower` after the road and Sanctum are authored. Leave its terrain as empty ground.

Update `internal/game/scene.go` so `drawHomePlotTile` draws the Bow Tower when it sees `featureBowTower`. Use `s.structureCatalog.BowTower.Sprite` and scale it inside the tile similarly to the Sanctum, but slightly smaller so it reads as a defensive tower beside the path.

Update tests in `internal/game/game_test.go` to assert that new states contain the starting tower, that the default map places it at `(8,5)`, that tile `(7,5)` remains road beside it, and that the existing interior-empty feature test allows the authored tower.

Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, and `ARCHITECTURE.md` so current behavior says the starting home Plot has one visible Bow Tower beside the road and that it does not yet affect combat.

## Concrete Steps

1. Edit `internal/game/map.go`, `internal/game/scene.go`, and `internal/game/game_test.go` as described above.
2. Update `cmd/td/main_test.go` so screenshot evidence writes to `plans/16-starting-bow-tower/screenshots/`.
3. Update the control documents listed in Plan of Work.
4. Run from the repository root:

       gofmt -w cmd/td/main_test.go internal/game/game_test.go internal/game/map.go internal/game/scene.go internal/game/structures.go internal/game/structures_test.go
       go test ./...
       TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
       git diff --check
       git status --short

5. Check hand-written code-file line counts:

       find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

   If any hand-written code file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and do not perform an unplanned split unless the user approves it.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, `git diff --check` reports no whitespace errors, screenshot evidence exists under `plans/16-starting-bow-tower/screenshots/`, and the rendered running-game screenshot shows one Bow Tower east of the north road. Documentation is accepted when it accurately says the tower is visible current behavior but not a combat system.

## Idempotence and Recovery

The map edit is deterministic. Re-running tests and screenshot capture is safe; screenshot capture overwrites only files under `plans/16-starting-bow-tower/screenshots/`. If the tower is hard to see at the chosen scale, adjust only the Bow Tower draw scale before considering larger rendering changes.

## Artifacts and Notes

Expected screenshot artifacts:

    plans/16-starting-bow-tower/screenshots/main-menu.png
    plans/16-starting-bow-tower/screenshots/new-game-configuration.png
    plans/16-starting-bow-tower/screenshots/running-game.png
    plans/16-starting-bow-tower/screenshots/paused-game.png
    plans/16-starting-bow-tower/screenshots/ingame-menu.png

## Interfaces and Dependencies

No new dependencies are required. The private map feature enum gains:

    featureBowTower

The existing structure catalog remains:

    type StructureCatalog struct {
        Sanctum  StructureTemplate
        BowTower StructureTemplate
    }

This plan does not add placement APIs, tower targeting APIs, combat rules, costs, or resource changes.
