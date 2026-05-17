# Add Zombie Enemy

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/22-zombie-enemy.md`, using the next available two-digit prefix after `plans/21-enemy-health-bars.md`.

## Purpose / Big Picture

After this change, the first Raid no longer consists only of skeletons. The second and fourth enemies in Raid 1 are zombies: slower, tougher enemies that make the opening assault more readable and varied without introducing a full wave-composition system. A player can see the change by clicking `Next Raid` and watching the first wave spawn skeleton, zombie, skeleton, zombie, skeleton on the existing north road.

## Progress

- [x] (2026-05-17T16:54:07Z) Inspected current Raid spawning, enemy catalog, asset catalog, rendering, tests, and control documents.
- [x] (2026-05-17T16:54:07Z) Confirmed the implementation defaults: zombies have 75 health, 0.7 Tiles-per-second speed, and appear only on every second spawn in Raid 1.
- [x] (2026-05-17T16:54:07Z) Generated and added a 64x64 transparent zombie PNG under `assets/sprites/enemies/`.
- [x] (2026-05-17T16:56:00Z) Wired the zombie sprite into the asset catalog and enemy template catalog.
- [x] (2026-05-17T16:57:00Z) Added deterministic first-Raid zombie spawn selection.
- [x] (2026-05-17T16:58:00Z) Updated tests for zombie asset loading, enemy stats, Raid 1 composition, and the combat test that verifies spawned enemy health.
- [x] (2026-05-17T16:59:00Z) Updated `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to match the implemented behavior.
- [x] (2026-05-17T16:59:30Z) Captured screenshot evidence under `plans/22-zombie-enemy/screenshots/`.
- [x] (2026-05-17T17:00:00Z) Ran validation commands and recorded results.
- [x] (2026-05-17T17:00:32Z) Checked hand-written code-file line counts; no files exceed the 600-line preference.

## Surprises & Discoveries

- Observation: The current asset catalog embeds each sprite explicitly, so the zombie PNG must be added to the `//go:embed` line and loaded by `assets.NewCatalog`.
  Evidence: `assets/catalog.go` lists `sprites/enemies/skeleton-sword-shield.png` directly in the embed directive.
- Observation: Raid enemies already store a template pointer and current health, so the zombie can reuse movement, targeting, damage, health-bar, Sanctum-contact, and rendering rules.
  Evidence: `internal/game/raid.go` stores `template *EnemyTemplate` and `health int` in `raidEnemy`.
- Observation: One combat test encoded the old assumption that the second spawned Raid enemy had skeleton health.
  Evidence: The first `go test ./...` run failed in `TestSpawnRaidEnemyAssignsHealthAndStableIDs` with `second enemy health = 75, want 50`; the test now checks the spawned enemy's own template max health.

## Decision Log

- Decision: Zombie stats are 75 max health and 0.7 Tiles per second.
  Rationale: This is clearly tougher and slower than the current 50-health, 1.0-speed skeleton without making the first Raid excessively long.
  Date/Author: 2026-05-17 / Codex
- Decision: Apply the zombie alternating rule only to Raid 1.
  Rationale: The user asked for the first wave. Later Raid composition should wait for a fuller wave-composition design.
  Date/Author: 2026-05-17 / Codex
- Decision: Use the generated 64x64 transparent zombie sprite as a required runtime asset.
  Rationale: Current Raid enemies are sprite-backed, and a distinct zombie visual makes the new enemy observable in the running game.
  Date/Author: 2026-05-17 / Codex

## Outcomes & Retrospective

Completed. The project now has a sprite-backed zombie enemy template with 75 health and 0.7 Tiles-per-second speed. Raid 1 spawns skeleton, zombie, skeleton, zombie, skeleton, while later placeholder Raids remain skeleton-only. Zombies reuse the existing movement, targeting, projectile damage, health-bar, Barricade, and breach rules.

Validation passed on 2026-05-17: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, and `git diff --check` all succeeded. Screenshot evidence was written under `plans/22-zombie-enemy/screenshots/`; `active-raid.png` shows a zombie in the first Raid's second spawn position.

The hand-written Go line-count review found no files over the 600-line preference. The largest file was `internal/game/game_test.go` at 588 lines, which is below the limit but close enough that future broad game-state tests should consider a responsibility-based split.

## Context and Orientation

This repository is a Go/Ebitengine local tower-defense prototype. Runtime assets live in the `assets` package, with source PNG files under `assets/sprites/`. Game state lives in `internal/game/`. The current Raid system is deterministic: clicking `Next Raid` starts a Raid, spawns enemies one at a time from the north-center road edge, and moves them toward the centered Sanctum using world coordinates in Tile units.

Enemy templates live in `internal/game/enemies.go`. A template contains the enemy name, maximum health, movement speed in Tiles per second, Sanctum damage, sprite key, and loaded Ebitengine image pointer. Active enemies in `internal/game/raid.go` store a pointer to one of these templates plus per-instance health and position. Rendering in `internal/game/raidui.go` already draws each enemy from its template sprite and draws a health bar based on template max health.

The relevant control documents are `GAME.md`, `PRODUCT.md`, `README.md`, and `ARCHITECTURE.md`. This change affects user-visible behavior and intended gameplay, so those files must describe the new zombie enemy and first-Raid composition. `ART.md` requires generated enemy assets to be 64x64 PNG pixel art under `assets/sprites/enemies/`.

## Plan of Work

First, add `assets/sprites/enemies/zombie.png` as a 64x64 transparent PNG. Update `assets/catalog.go` so the embed directive includes the zombie path, `EnemySprites` contains a `Zombie *ebiten.Image`, and `NewCatalog` loads the zombie sprite and returns it in the catalog. Add a catalog test that verifies the zombie sprite loads and is exactly 64x64.

Next, update `internal/game/enemies.go` so `EnemyCatalog` contains a `Zombie EnemyTemplate`. In `NewEnemyCatalog`, set the zombie template to `Name: "Zombie"`, `MaxHealth: 75`, `SpeedTilesPerSecond: 0.7`, `SanctumDamage: 1`, `SpriteKey: "zombie"`, and `Sprite: assetCatalog.Sprite.Enemy.Zombie`. Add a test beside the skeleton catalog test.

Then, change Raid spawning in `internal/game/raid.go`. Keep the existing Raid counts and spawn interval. Add a private method that chooses the next enemy template for the current Raid. It should return the zombie template when `s.raid.number == 1` and `s.raid.nextEnemyID%2 == 1`, because `nextEnemyID` is the zero-based spawn position. It should return the skeleton template in all other cases. Use the selected template for both `template` and initial `health` in `spawnRaidEnemy`.

Finally, update tests and docs. Raid tests should prove that Raid 1 spawns skeleton, zombie, skeleton, zombie, skeleton and that Raid 2 still starts with skeletons. Documentation should say the first Raid alternates zombies into every second spawn and that later Raids remain simple skeleton-only placeholders for now.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

Create or update these files:

    assets/sprites/enemies/zombie.png
    assets/catalog.go
    assets/catalog_test.go
    internal/game/enemies.go
    internal/game/enemies_test.go
    internal/game/raid.go
    internal/game/raid_test.go
    cmd/td/main_test.go
    README.md
    PRODUCT.md
    GAME.md
    ARCHITECTURE.md

After implementation, format Go files:

    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/enemies.go internal/game/enemies_test.go internal/game/raid.go internal/game/raid_test.go cmd/td/main_test.go

Run the full test suite:

    go test ./...

Capture screenshot evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected result: the command passes and writes screenshots under `plans/22-zombie-enemy/screenshots/`, including `active-raid.png` with a zombie visible during the first Raid.

Check whitespace:

    git diff --check

Check hand-written Go file line counts and report any files over 600 lines before performing unplanned refactors:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

## Validation and Acceptance

The change is accepted when `go test ./...` passes, screenshot capture writes evidence under `plans/22-zombie-enemy/screenshots/`, `git diff --check` reports no whitespace errors, and the line-count review is recorded.

Behaviorally, clicking `Next Raid` starts Raid 1 with five enemies in this order: skeleton, zombie, skeleton, zombie, skeleton. Zombies render from a distinct sprite, use health bars like skeletons, move slower than skeletons, survive more damage than skeletons, and otherwise follow the same north-road movement, targeting, projectile-damage, Barricade, and breach rules. Starting Raid 2 after a successful Raid still uses skeleton-only composition.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` all describe the new zombie enemy consistently and do not imply that skeletons are still the only enemy archetype.

## Idempotence and Recovery

The code edits are deterministic and can be reapplied safely. Screenshot capture overwrites only files under `plans/22-zombie-enemy/screenshots/`. If the zombie sprite fails to load, inspect the `//go:embed` path and the `loadSprite` call before changing game logic. If Raid composition tests fail, inspect the zero-based `nextEnemyID` selection before changing Raid counts or spawn timing.

## Artifacts and Notes

The generated source image was created with the built-in image generation tool and resized to 64x64. The final project asset is `assets/sprites/enemies/zombie.png`.

## Interfaces and Dependencies

No new external dependencies are required. The public Go package surface remains minimal, but the following project-local interfaces must exist after implementation:

    type EnemySprites struct {
        SkeletonSwordShield *ebiten.Image
        Zombie              *ebiten.Image
    }

    type EnemyCatalog struct {
        SkeletonSwordShield EnemyTemplate
        Zombie              EnemyTemplate
    }

The zombie template is private to the current game package usage pattern through `State.enemyCatalog`; no new pathing API, wave API, save format, settings option, or tower behavior is introduced.
