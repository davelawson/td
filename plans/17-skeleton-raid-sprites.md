# Skeleton Raid Sprites

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

Save this file at `plans/17-skeleton-raid-sprites.md` in the repository root. The previous highest numbered plan is `plans/16-starting-bow-tower.md`.

## Purpose / Big Picture

After this change, starting a Raid makes the attackers read as actual enemies instead of abstract markers. The existing deterministic Raid pacing remains unchanged: clicking `Next Raid` starts a Raid, one skeleton appears immediately, and the rest spawn one at a time on the existing stagger. The observable result is that active Raid enemies render with the loaded skeleton sword-and-shield sprite from the asset catalog while they move down the north road toward the Sanctum.

This is a visual and catalog-wiring slice only. It does not add animation, health bars, tower attacks, damage, rewards, pathfinding, new enemy archetypes, or a new Wave type.

## Progress

- [x] (2026-05-16T00:00Z) Inspected the existing Raid, enemy catalog, asset catalog, rendering path, screenshot capture, and control-document wording.
- [x] (2026-05-16T00:02Z) Confirmed with the user that the existing staggered Raid spawning should remain instead of instantiating all enemies upfront.
- [x] (2026-05-16T00:10Z) Wired `EnemyTemplate` to carry the loaded skeleton sprite from `assets.Catalog`.
- [x] (2026-05-16T00:12Z) Replaced active Raid enemy circle rendering with skeleton sprite rendering, retaining a small nil-sprite fallback.
- [x] (2026-05-16T00:14Z) Updated tests to assert skeleton sprite catalog wiring for enemy templates and spawned Raid enemies.
- [x] (2026-05-16T00:16Z) Updated screenshot capture to write to this plan directory and include an active-Raid screenshot.
- [x] (2026-05-16T00:18Z) Updated `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to describe sprite-backed skeleton Raid enemies.
- [x] (2026-05-16T00:22Z) Ran validation commands and captured screenshot evidence under this plan directory.
- [x] (2026-05-16T00:24Z) Checked hand-written code-file line counts against the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: The skeleton sprite asset was already embedded and loaded as `assets.Catalog.Sprite.Enemy.SkeletonSwordShield`, and existing Raid enemies already referenced `EnemyTemplate`.
  Evidence: `assets/catalog.go` loads `sprites/enemies/skeleton-sword-shield.png`, and `internal/game/raid.go` appends enemies with `template: &s.enemyCatalog.SkeletonSwordShield`.

- Observation: `internal/game/game_test.go` was already close to the 600-line preference before this change.
  Evidence: `wc -l internal/game/game_test.go` reported 570 lines before implementation.

## Decision Log

- Decision: Preserve the existing staggered Raid model.
  Rationale: The user chose the recommended option to keep current Raid pacing. This avoids changing the tested lifecycle while still ensuring every spawned enemy is a skeleton.
  Date/Author: 2026-05-16 / User and Codex

- Decision: Add `Sprite *ebiten.Image` to `EnemyTemplate`, matching the existing structure-template pattern.
  Rationale: Structure rendering already keeps loaded sprites on templates. Using the same pattern keeps asset-path knowledge inside `assets` and catalog construction instead of leaking it into Raid logic.
  Date/Author: 2026-05-16 / Codex

- Decision: Keep a nil-sprite circle fallback inside rendering.
  Rationale: Production state should always have the loaded sprite, but the fallback prevents a panic if a future test or incomplete catalog state creates an enemy template without a sprite.
  Date/Author: 2026-05-16 / Codex

## Outcomes & Retrospective

Completed. Active Raid enemies now render with the loaded skeleton sword-and-shield sprite instead of placeholder circles. The existing deterministic Raid lifecycle is unchanged: one skeleton appears immediately, later skeletons spawn on the existing stagger, enemy counts still include active plus pending enemies, and reaching enemies still spend Barricade charges or breach the Sanctum.

Validation passed on 2026-05-16: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, and `git diff --check` all succeeded. Screenshot evidence was written under `plans/17-skeleton-raid-sprites/screenshots/`, including `active-raid.png`, which shows a skeleton on the north road during an active Raid.

The hand-written Go line-count review found no files over 600 lines. `internal/game/game_test.go` is 573 lines, which remains close enough to the 600-line preference that the next meaningful game-test addition should consider splitting tests by responsibility, such as moving Raid-focused tests fully into `internal/game/raid_test.go` or map tests into a dedicated map test file. No split was performed because this plan did not include that extra refactor and the file remains under the limit.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. Ebitengine is the Go game library used by the executable under `cmd/td/`. The playable game state lives in `internal/game/`.

The runtime asset catalog lives in `assets/catalog.go`. It embeds PNG files under `assets/sprites/` and exposes loaded Ebitengine images through `assets.Catalog`. The relevant existing field is `assetCatalog.Sprite.Enemy.SkeletonSwordShield`, loaded from `assets/sprites/enemies/skeleton-sword-shield.png`.

Enemy template metadata lives in `internal/game/enemies.go`. A template describes shared values for enemies of one type. Before this plan, `EnemyTemplate` had a `SpriteKey` string but not the loaded sprite image itself.

Raid simulation lives in `internal/game/raid.go`. `startNextRaid` increments the Raid number, sets `pendingEnemies` from `raidEnemyCount`, clears active enemies, spawns one enemy immediately, and sets the phase to Raid. `updateRaidSpawning` later spawns more enemies after `raidSpawnInterval` updates. A Raid enemy stores a pointer to its `EnemyTemplate` plus path progress.

Raid rendering lives in `internal/game/raidui.go`. Before this plan, `drawRaidEnemies` projected each active enemy onto the road and drew a purple circle. After this plan, it should use the enemy template sprite and keep a simple circle fallback only when the sprite is nil.

`PRODUCT.md`, `README.md`, `GAME.md`, and `ARCHITECTURE.md` are durable control documents. This change affects current user-visible behavior, intended game-design wording for the first Raid slice, and the architecture description of enemy rendering, so those files must be updated in the same change.

## Plan of Work

Update `internal/game/enemies.go` so `EnemyTemplate` includes `Sprite *ebiten.Image`. Import `td/assets` and Ebitengine in this file. Change `NewEnemyCatalog` to accept `assetCatalog assets.Catalog` and assign `assetCatalog.Sprite.Enemy.SkeletonSwordShield` to the skeleton template. Keep the existing name, health, speed, damage, and `SpriteKey`.

Update `internal/game/game.go` so `New` constructs the enemy catalog with `NewEnemyCatalog(assetCatalog)` after the asset catalog has loaded. Do not load assets in gameplay rules and do not introduce asset paths outside the `assets` package.

Update `internal/game/raidui.go` so `drawRaidEnemies` delegates each active enemy to a helper that draws the enemy template sprite centered at `raidEnemyWorldPosition`. Use a world-space square around the enemy center, project that square through the existing camera projection, scale the 64x64 sprite to the projected rectangle, and draw it with `screen.DrawImage`. Preserve the current circle drawing as a fallback when `enemy.template.Sprite` is nil. Keep Raid controls unchanged.

Update `internal/game/enemies_test.go` to construct `assets.NewCatalog`, pass it into `NewEnemyCatalog`, and assert that `SkeletonSwordShield.Sprite` is non-nil and exactly references `assetCatalog.Sprite.Enemy.SkeletonSwordShield`. Update `internal/game/game_test.go` with a small assertion that a new `State` stores a skeleton enemy sprite. Update `internal/game/raid_test.go` so the first immediate spawn and staggered spawn both assert their template has a non-nil sprite. Avoid adding new tests to `internal/game/game_test.go` beyond the existing state initialization assertion because that file is close to 600 lines.

Update `cmd/td/main_test.go` so screenshot evidence writes under `plans/17-skeleton-raid-sprites/screenshots/`. Add an active-Raid screenshot target that starts a game, clicks the fixed `Next Raid` button center, and captures `active-raid.png`. Existing menu, running, paused, and in-game-menu screenshots should still be captured.

Update `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` so they describe deterministic placeholder Raids with sprite-backed skeleton enemies. Keep the documented limitations clear: the Bow Tower remains non-combat, and there is still no tower targeting, combat damage, rewards, alternate paths, or enemy variety.

## Concrete Steps

Run these commands from the repository root, `/home/dave/dev/ai/td`.

Format edited Go files:

    gofmt -w internal/game/enemies.go internal/game/game.go internal/game/raidui.go internal/game/enemies_test.go internal/game/game_test.go internal/game/raid_test.go cmd/td/main_test.go

Run the full test suite:

    go test ./...

Expected result: all packages pass.

Capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expected result: the command passes and writes screenshots under `plans/17-skeleton-raid-sprites/screenshots/`, including `active-raid.png` with a skeleton sprite visible on the road.

Check whitespace and working tree state:

    git diff --check
    git status --short

Expected result: no whitespace errors. `git status --short` should list only intentional modified files, the new ExecPlan, and new screenshot artifacts.

Check hand-written Go file line counts:

    find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

If any hand-written code file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and do not perform an unplanned split, refactor, or library addition unless the user approves it.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, `git diff --check` reports no whitespace errors, and `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` writes `plans/17-skeleton-raid-sprites/screenshots/active-raid.png` showing a skeleton sprite on the road during an active Raid.

Behavioral acceptance is that clicking `Next Raid` still starts the deterministic Raid immediately, the first enemy still spawns immediately, later enemies still spawn on the same fixed stagger, top-bar enemy counts still reflect pending plus active enemies, and reaching enemies still spend Barricade charges or breach the Sanctum exactly as before.

Documentation is accepted when `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` say current Raid enemies are sprite-backed skeletons while preserving the explicit limitations around tower targeting, combat damage, rewards, alternate paths, and enemy variety.

## Idempotence and Recovery

The code edits are deterministic and can be reapplied safely. Screenshot capture overwrites only files under `plans/17-skeleton-raid-sprites/screenshots/`. If the sprite appears too large or too small, adjust only the local `raidEnemySpriteSize` constant in `internal/game/raidui.go` and rerun tests plus screenshot capture. If any catalog-wiring change breaks initialization, revert the `NewEnemyCatalog(assetCatalog)` signature and associated tests before retrying a smaller edit.

## Artifacts and Notes

Expected screenshot artifacts:

    plans/17-skeleton-raid-sprites/screenshots/main-menu.png
    plans/17-skeleton-raid-sprites/screenshots/new-game-configuration.png
    plans/17-skeleton-raid-sprites/screenshots/running-game.png
    plans/17-skeleton-raid-sprites/screenshots/active-raid.png
    plans/17-skeleton-raid-sprites/screenshots/paused-game.png
    plans/17-skeleton-raid-sprites/screenshots/ingame-menu.png

Validation transcript:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.029s
    ok  	td/internal/game	0.066s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	1.520s

    git diff --check
    # no output

    find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l | sort -n
    # largest file: ./internal/game/game_test.go at 573 lines

## Interfaces and Dependencies

No new third-party dependency is required. Use the existing Ebitengine dependency already present in `go.mod`.

The enemy template interface after this change is:

    type EnemyTemplate struct {
        Name          string
        MaxHealth     int
        Speed         float64
        SanctumDamage int
        SpriteKey     string
        Sprite        *ebiten.Image
    }

The enemy catalog constructor after this change is:

    func NewEnemyCatalog(assetCatalog assets.Catalog) EnemyCatalog

This plan does not add a public Wave API, placement API, targeting API, combat API, resource API, or new enemy archetype.

## Revision Notes

2026-05-16: Created this plan for the skeleton Raid sprite implementation. The plan records the decision to preserve staggered Raid spawning while rendering spawned skeletons from the existing asset catalog.
