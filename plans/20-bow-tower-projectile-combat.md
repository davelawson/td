# Add Bow Tower Projectile Combat

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

## Purpose / Big Picture

After this change, the visible Bow Tower is no longer scenery-only. During an active Raid, it periodically launches the authored arrow projectile sprite at skeleton enemies within range. When a projectile reaches its original target, it deals Bow Tower damage to that enemy, and defeated enemies are removed before they can breach the Sanctum.

## Progress

- [x] (2026-05-16T22:49:33Z) Confirmed scope: automated Bow Tower targeting, projectile travel, enemy health, damage, rendering, tests, and current-state documentation.
- [x] (2026-05-16T22:57:03Z) Added real-time Bow Tower combat stats and projectile sprite references to structure templates.
- [x] (2026-05-16T22:57:03Z) Added enemy IDs, health, projectile state, targeting, cooldown, and damage rules.
- [x] (2026-05-16T22:57:03Z) Rendered active projectiles with the existing camera projection.
- [x] (2026-05-16T22:57:03Z) Added focused combat tests outside `internal/game/game_test.go`.
- [x] (2026-05-16T22:57:03Z) Updated `README.md`, `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, and `ROADMAP.md`.
- [x] (2026-05-16T22:57:03Z) Captured active-Raid screenshot evidence under this plan.
- [x] (2026-05-16T22:57:03Z) Ran validation commands: `go test ./...`, screenshot capture, `git diff --check`, and `git status --short`.
- [x] (2026-05-16T22:57:03Z) Checked hand-written code-file line counts; no file exceeds the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: The prior projectile asset/catalog work is present in the working tree but not committed.
  Evidence: `git status --short` shows `assets/sprites/structures/bow-tower-projectile.png`, `plans/19-bow-tower-projectile-sprite.md`, and catalog edits.

- Observation: `internal/game/game_test.go` is already near the 600-line code-style preference before this combat slice.
  Evidence: line-count inspection reported `572 internal/game/game_test.go`.

## Decision Log

- Decision: Store Bow Tower timing in seconds, not updates.
  Rationale: The user explicitly requested real-time units for structure timing. Runtime logic may convert through the current fixed-step update duration, but template stats and cooldown state should remain second-based.
  Date/Author: 2026-05-16 / Codex

- Decision: Use range 3.0 tiles, damage 10, fire interval 1.0 second, and projectile speed 9.0 tiles per second.
  Rationale: This moderate baseline kills the current 20-health skeleton in two hits while keeping projectile travel visible and testable.
  Date/Author: 2026-05-16 / Codex

- Decision: Target the in-range enemy closest to the Sanctum, with lower enemy ID as a deterministic tie-breaker.
  Rationale: The closest-to-Sanctum enemy is the most urgent threat on the current single road, and ID tie-breaking keeps tests deterministic.
  Date/Author: 2026-05-16 / Codex

- Decision: Remove projectiles harmlessly if their original target is gone before impact.
  Rationale: This keeps the first projectile model simple and avoids surprising retargeting behavior.
  Date/Author: 2026-05-16 / Codex

## Outcomes & Retrospective

Completed. The Bow Tower now has prototype combat stats in real-time units: 3.0-Tile range, 10 damage, 1.0-second fire interval, and 9.0-Tiles-per-second projectile speed. Active Raid enemies now have stable IDs and health. During active Raids, the Bow Tower targets the in-range enemy closest to the Sanctum, launches visible arrow projectiles, applies damage on projectile hit, and removes defeated enemies.

Validation passed: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `git status --short` all completed. Screenshot evidence was written under `plans/20-bow-tower-projectile-combat/screenshots/`, and `active-raid.png` shows an arrow projectile in flight between the Bow Tower and road enemies.

The hand-written Go line-count review found no files over 600 lines. `internal/game/game_test.go` remains the largest file at 572 lines, so future changes touching broad game behavior should continue adding focused test files rather than growing that file.

## Context and Orientation

The repository is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. Game logic lives in `internal/game/`. The fixed home Plot is created in `internal/game/map.go`; it includes a Sanctum at the center and one Bow Tower at tile `(8,5)`, one tile east of the north road. Raid simulation currently lives in `internal/game/raid.go`, and Raid rendering lives in `internal/game/raidui.go`.

World positions use Tile units. The center of the Sanctum Tile is `(0,0)`, one Tile is one world unit, and positive Y points north. The Bow Tower at tile `(8,5)` has world center `(1,2)`. Skeleton enemies currently spawn at `(0,7)` and move south toward the Sanctum along the road.

The previous plan added `assets/sprites/structures/bow-tower-projectile.png` and exposed it as `Catalog.Sprite.Projectile.BowTowerProjectile`. This combat plan depends on that asset and catalog field.

`PRODUCT.md` and `README.md` currently say tower targeting and combat damage do not exist. `GAME.md` says exact Bow Tower stats are open. `ARCHITECTURE.md` says the game package owns Raid movement and rendering but not yet combat. This change makes those statements stale, so all four documents must be updated.

## Plan of Work

Update `internal/game/structures.go` so `StructureTemplate` can describe combat-capable towers. Add fields for range in tiles, damage, fire interval in seconds, projectile speed in tiles per second, and projectile sprite. Populate only the Bow Tower with combat stats; the Sanctum remains non-combat with zero-value combat fields.

Add a small fixed-step duration constant such as `gameUpdateSeconds = 1.0 / 60.0` in `internal/game`. Use it only as the bridge from Ebitengine's current update cadence to second-based simulation. Do not store Bow Tower fire interval or cooldowns as update counts.

Update `internal/game/raid.go` so active enemies have stable IDs and health. Spawned skeletons should start at their template's `MaxHealth`. Add combat update logic that scans Bow Tower features in the home Plot, finds targets within range, spawns projectiles, advances active projectiles toward their original target, applies damage on hit, removes defeated enemies, and removes projectiles whose target is gone. Run combat after spawning and before enemy movement each unpaused Raid update.

Render projectiles from `Catalog.Sprite.Projectile.BowTowerProjectile` in `internal/game/raidui.go` using the existing camera projection. The sprite should stay readable at a small world size and appear over the map during active Raids. Rotation toward travel direction is preferred if practical without adding dependencies; a fixed diagonal sprite is acceptable for this slice.

Update tests. Add new focused tests in a dedicated file such as `internal/game/combat_test.go` instead of growing `internal/game/game_test.go`. Preserve existing Raid lifecycle tests by updating fixtures to include enemy health where needed.

Update `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to describe the first implemented Bow Tower combat slice and the prototype stats. Do not change `ROADMAP.md` unless implementation reveals a sequencing change beyond completing the already listed defense-loop priority.

## Concrete Steps

Run these commands from the repository root.

1. Edit `internal/game/structures.go`, `internal/game/raid.go`, `internal/game/raidui.go`, and focused tests as described in Plan of Work.

2. Update current-state and design documents:

       README.md
       PRODUCT.md
       GAME.md
       ARCHITECTURE.md

3. Update `cmd/td/main_test.go` so screenshot evidence writes to `plans/20-bow-tower-projectile-combat/screenshots/`, then capture evidence:

       TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

4. Run:

       go test ./...
       git diff --check
       git status --short

5. Check hand-written code-file line counts:

       find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

   Report any hand-written code file over 600 lines and recommend a concrete response. Do not perform unplanned refactors, code splits, or library additions without user approval.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, screenshot capture writes active-Raid evidence under `plans/20-bow-tower-projectile-combat/screenshots/`, `git diff --check` reports no whitespace errors, and the line-count review is recorded.

Behavioral acceptance is that during an active Raid the Bow Tower fires at skeletons once they enter a 3.0-tile range, visible projectiles travel from tower to enemy, each hit deals 10 damage, two hits remove the current 20-health skeleton, and the top bar's remaining-enemy count decreases when enemies are defeated.

## Idempotence and Recovery

The code edits are deterministic and can be reapplied safely. Screenshot capture overwrites only files under `plans/20-bow-tower-projectile-combat/screenshots/`. If tests show enemies no longer breach correctly, inspect the combat update order before changing Sanctum contact rules. If projectiles never hit, inspect target lookup by enemy ID and distance-step clamping before changing balance values.

## Artifacts and Notes

Screenshot evidence:

    plans/20-bow-tower-projectile-combat/screenshots/main-menu.png
    plans/20-bow-tower-projectile-combat/screenshots/new-game-configuration.png
    plans/20-bow-tower-projectile-combat/screenshots/running-game.png
    plans/20-bow-tower-projectile-combat/screenshots/active-raid.png
    plans/20-bow-tower-projectile-combat/screenshots/paused-game.png
    plans/20-bow-tower-projectile-combat/screenshots/ingame-menu.png

Final validation commands:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    git status --short

## Interfaces and Dependencies

No new Go module dependencies are required. Use Ebitengine's existing image drawing APIs for projectile rendering.

The Bow Tower structure template must expose second-based combat values, not update-count values:

    RangeTiles                  float64
    Damage                      int
    FireIntervalSeconds         float64
    ProjectileSpeedTilesPerSecond float64
    ProjectileSprite            *ebiten.Image

Runtime tower cooldowns should also be seconds remaining. Projectile movement should use `ProjectileSpeedTilesPerSecond * gameUpdateSeconds` each unpaused logical update.
