# Add Catapult Tower

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/29-catapult-tower.md`.

## Purpose / Big Picture

The current prototype has two projectile towers: Bow Tower and Flame Bolt Tower. After this change, the player can build a third tower, the Catapult Tower, from the left building bar. It is a long-range, slow-firing, slow-projectile, high-damage siege tower that damages every enemy in the Tile where its projectile lands. The behavior is visible by starting the game, dragging the Catapult icon onto a valid empty grass-like Tile during calm play, starting a Raid, and watching its slow boulder projectile strike a Tile on the road.

This work deliberately does not add upgrades, range previews, new targeting modes, damage types, pathfinding, a broader economy, or a starting Catapult on the default map.

## Progress

- [x] (2026-05-26T00:00:00Z) Inspected current tower templates, projectile combat, building bar, placement, asset catalog, selection panel, documentation, and tests.
- [x] (2026-05-26T00:00:00Z) Confirmed implementation choices with the user: Catapult is placeable only, uses generated pixel-art assets, and uses the long slow heavy stat line.
- [x] (2026-05-26T00:00:00Z) Created this ExecPlan at `plans/29-catapult-tower.md`.
- [x] (2026-05-26T04:34:02Z) Generated and installed Catapult structure and projectile PNG assets under `assets/sprites/structures/`.
- [x] (2026-05-26T04:34:02Z) Wired Catapult assets through `assets/catalog.go` and asset tests.
- [x] (2026-05-26T04:34:02Z) Added Catapult tower template, map feature, building-bar entry, placement mapping, selection-panel support, and rendering support.
- [x] (2026-05-26T04:34:02Z) Extended projectile combat so Catapult damages all enemies in the struck Tile while Bow and Flame Bolt remain single-target.
- [x] (2026-05-26T04:34:02Z) Updated focused Go tests for assets, structure stats, building bar, placement, selection panel, and area impact.
- [x] (2026-05-26T04:34:02Z) Updated `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to describe the new current behavior and design decision.
- [x] (2026-05-26T04:34:02Z) Captured screenshot evidence under `plans/29-catapult-tower/screenshots/`.
- [x] (2026-05-26T04:34:02Z) Ran `gofmt`, `go test ./...`, screenshot capture, `git diff --check`, ownership check, and recorded results.
- [x] (2026-05-26T04:34:02Z) Checked hand-written Go file line counts against the 600-line preference and reported the existing over-limit file.

## Surprises & Discoveries

- Observation: `go test ./...` fails before Catapult work because tests expect Bow Tower to cost 10 Metal while `internal/game/structures.go` currently has `Metal: 0`.
  Evidence: `TestNewStructureCatalogIncludesBowTower`, `TestBuildingBarItemsExposeTowerIcons`, and `TestBuildDragPlacesTowerAndDeductsResources` fail with Bow Tower cost/resource mismatches.

## Decision Log

- Decision: Add Catapult as a placeable building-bar tower but do not place one in `NewDefaultHomePlot`.
  Rationale: This makes the new tower user-accessible without changing the initial authored defense layout or giving the player a free high-damage tower.
  Date/Author: 2026-05-26 / Codex

- Decision: Use generated 64x64 pixel-art PNGs for the Catapult tower and projectile.
  Rationale: `ART.md` defines generated 2D pixel-art PNGs as the current asset style, and reusing existing tower sprites would make the new tower hard to read.
  Date/Author: 2026-05-26 / Codex

- Decision: Use prototype Catapult stats of 5.0-Tile range, 75 damage, 3.0-second fire interval, 3.0-Tiles-per-second projectile speed, and cost 40 Wood, 60 Stone, 25 Metal.
  Rationale: These values make the tower clearly long-range, slow, expensive, and high-impact compared with Bow and Flame Bolt while avoiding an extreme one-shot tower dominating the small prototype map.
  Date/Author: 2026-05-26 / Codex

- Decision: Interpret "damages all enemies in the tile it strikes" as damaging every living enemy whose current world position maps to the same Plot Tile as the original target's current impact position.
  Rationale: The current projectile model tracks an original target; using that target's current Tile at impact preserves deterministic targeting while making the area effect match the map grid.
  Date/Author: 2026-05-26 / Codex

- Decision: Restore Bow Tower's documented and tested Metal cost of 10 as part of this change.
  Rationale: The repository's current tests and `GAME.md` agree that Bow Tower costs 30 Wood, 10 Stone, and 10 Metal, so the code is the outlier and must be corrected before the new feature can validate cleanly.
  Date/Author: 2026-05-26 / Codex

## Outcomes & Retrospective

Implemented Catapult Tower as a third build-bar tower with generated 64x64 structure and projectile sprites. The asset catalog now embeds and loads Catapult sprites. The structure catalog exposes Catapult Tower with a 5.0-Tile range, 75 damage, 3.0-second fire interval, 3.0-Tiles-per-second projectile speed, and 40 Wood, 60 Stone, 25 Metal construction cost. The building bar shows Bow Tower, Flame Bolt Tower, and Catapult Tower in order, and placement maps Catapult to its own map feature when resources are sufficient.

Combat now preserves existing single-target Bow and Flame Bolt behavior while allowing Catapult projectiles to damage every living enemy in the Tile occupied by their original target when they land. Focused tests cover asset loading, structure stats, building-bar ordering and affordability, placement, selection-panel rows, Catapult firing stats, Tile-area impact, adjacent-Tile exclusion, and area defeat removal.

Validation completed:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go assets/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

`go test ./...` passed. Screenshot evidence was written under `plans/29-catapult-tower/screenshots/`: `main-menu.png`, `new-game-configuration.png`, `running-game.png`, `placed-tower.png`, `selected-tower.png`, `active-raid.png`, `selected-raider.png`, `paused-game.png`, and `ingame-menu.png`. `git diff --check` and the ownership check produced no output.

The final hand-written Go line-count review found one file over the 600-line preference: `internal/game/game_test.go` at 602 lines. This overage was pre-existing and this implementation avoided adding to it by placing new tests in focused files. Recommended follow-up is a responsibility-based split of `internal/game/game_test.go`, for example moving remaining state, camera, or map tests into focused files. No extra refactor was performed because that split is outside this plan.

## Context and Orientation

This repository is a Go/Ebitengine tower-defense prototype. Ebitengine is the 2D game library used by the executable under `cmd/td/`. Runtime game logic lives under `internal/game/`. Runtime sprites are embedded and loaded by `assets/catalog.go`; gameplay packages use the typed asset catalog and must not decode asset paths directly.

The current home Plot is a 15x15 Tile grid created in `internal/game/map.go`. The Sanctum is at the center. The road runs north from the Sanctum. The default map starts with one Bow Tower east of the road and one Flame Bolt Tower west of the road. The new Catapult Tower must not be added to this starting map.

Tower templates live in `internal/game/structures.go`. Each `StructureTemplate` already has construction cost, range, damage, fire interval, projectile speed, and projectile sprite fields. Combat lives in `internal/game/combat.go`: during active Raids, `fireCombatTowers()` scans all home Plot Tiles, finds features that map to combat tower templates, targets the in-range enemy closest to the Sanctum, and launches projectiles. `updateProjectiles()` moves each projectile toward its original target and applies damage on impact. Catapult should reuse this system and add only the minimum area-impact behavior needed.

The building bar lives in `internal/game/building_bar.go`. It currently returns two hardcoded items and maps item indices 0 and 1 to Bow and Flame Bolt features. This change should make the item list include Catapult as the third item and map index 2 to `featureCatapultTower`.

Selection details live in `internal/game/selection_panel.go`. The selected tower panel displays tower type, range, attack speed, and damage. Catapult should use the same panel rows.

Root control documents constrain the implementation. `GAME.md` records durable game-design decisions and must add Catapult Tower as a tower type. `PRODUCT.md` and `README.md` must describe current user-visible behavior after the change. `ARCHITECTURE.md` must describe the new asset catalog entries, tower feature, building-bar entry, and tile-area impact. `ART.md` already defines the relevant generated asset guidance and does not need a guidance change unless implementation changes that guidance. `DESIGN.md`, `ROADMAP.md`, and `CODESTYLE.md` do not need updates unless implementation introduces new durable visual rules, product sequencing changes, or source conventions.

## Plan of Work

First, install the Catapult art assets. Generate a structure sprite and a projectile sprite using the current generated pixel-art style: 64x64 PNG, sharp contrast, low detail, bright colors, no shadows, and no text. Put the finished project assets at `assets/sprites/structures/catapult-tower.png` and `assets/sprites/structures/catapult-tower-projectile.png`.

Second, update the asset catalog. In `assets/catalog.go`, add both new PNG paths to the `//go:embed` directive, add `CatapultTowerProjectile` to `ProjectileSprites`, add `CatapultTower` to `StructureSprites`, load both images in `NewCatalog`, and assign them in the returned catalog. Update `assets/catalog_test.go` to assert both new sprite fields are non-nil.

Third, update tower data. In `internal/game/map.go`, add `featureCatapultTower` after `featureFlameBoltTower`. In `internal/game/structures.go`, add `CatapultTower StructureTemplate` to `StructureCatalog`, restore Bow Tower cost to `Resources{Wood: 30, Stone: 10, Metal: 10}`, and populate Catapult with name `Catapult Tower`, sprite and projectile sprite from the asset catalog, cost `Resources{Wood: 40, Stone: 60, Metal: 25}`, `RangeTiles: 5.0`, `Damage: 75`, `FireIntervalSeconds: 3.0`, and `ProjectileSpeedTilesPerSecond: 3.0`. Add a boolean field to `StructureTemplate`, such as `DamageAllEnemiesInTargetTile`, and set it true only for Catapult.

Fourth, update building and selection. In `internal/game/building_bar.go`, make `buildingBarItems()` include Catapult as the third item below Flame Bolt with stable spacing, and update `buildingFeatureForItemIndex()` to map index 2 to `featureCatapultTower`. Keep existing affordability hover behavior and placement restrictions. In `internal/game/scene.go`, update feature rendering so a placed Catapult draws with the Catapult sprite. In `internal/game/selection_panel.go`, map `featureCatapultTower` to `towerSelectionPanel(s.structureCatalog.CatapultTower)`.

Fifth, update combat. In `internal/game/combat.go`, add an area-impact flag to `combatProjectile`. When a tower fires, copy the template's area flag into the projectile. Keep target selection unchanged. Update `combatTowerTemplate()` so `featureCatapultTower` returns the Catapult template. In `updateProjectiles()`, when a projectile reaches its original target and the area flag is false, keep the existing single-target damage path. When the area flag is true, compute the Tile coordinate for the original target's current world position and damage every living enemy currently in that same Tile. Use a helper that converts a world `coord` to a home Plot `tileCoordinate` by reversing `tileWorldCenter`: `X = floor(position.X + homePlotCenter + 0.5)`, `Y = floor(homePlotCenter - position.Y + 0.5)`, reject coordinates outside `[0, plotSize)`. Iterate enemies carefully because `damageEnemy()` removes defeated enemies; a simple index loop that only increments when the current enemy survives is acceptable.

Sixth, update tests. Keep new tests in focused files rather than growing `internal/game/game_test.go`, which is already over the 600-line preference. Update current Bow Tower cost assertions to match restored `Metal: 10`. Add tests for Catapult structure stats, asset wiring, building bar item order and cost, placement feature mapping and resource deduction, selected Catapult panel rows, Catapult firing/cooldown/projectile speed/damage, Catapult Tile-area damage, adjacent-Tile exclusion, and vanished-target behavior.

Seventh, update screenshot capture and documentation. Change screenshot evidence output in `cmd/td/main_test.go` to `plans/29-catapult-tower/screenshots/` and capture at least the standard screenshots. Add a Catapult placement screenshot if the screenshot harness can place one without making the test brittle. Update `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` to describe the current placeable Catapult and area impact without implying upgrades, range previews, new damage types, or broader base-building exists.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

Before implementation, confirm the baseline and record known failures:

    go test ./...

Expected pre-change failure is limited to Bow Tower cost/resource assertions.

After code edits, format and validate:

    gofmt -w cmd/td/main.go cmd/td/main_test.go internal/game/*.go assets/*.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

End with the required hand-written Go line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path, line count, likely cause, and recommended response in `Outcomes & Retrospective`. Do not perform an unplanned split or refactor unless the user approves it.

## Validation and Acceptance

Automated acceptance is that `go test ./...` passes, screenshot capture writes evidence under `plans/29-catapult-tower/screenshots/`, `git diff --check` has no output, ownership check has no output, and the line-count review is recorded in this plan.

Behavioral acceptance is that the building bar shows Bow Tower, Flame Bolt Tower, and Catapult Tower in that order; Catapult uses its own sprite and cost row; Catapult can be placed only on valid empty grass-like Tiles during calm play when affordable; placed Catapults fire slow projectiles during Raids; and a Catapult projectile impact damages all active enemies in the struck Tile while not damaging enemies in neighboring Tiles.

Documentation acceptance is that `README.md`, `PRODUCT.md`, `GAME.md`, and `ARCHITECTURE.md` all describe Catapult Tower as a current placeable tower with slow high-damage Tile-area impact, and they still make clear that upgrades, range previews, broader base-building, exploration, and resource gathering are not implemented.

## Idempotence and Recovery

The edits are additive except for restoring the Bow Tower Metal cost to the documented/tested value. Re-running asset loading tests, Go tests, screenshot capture, whitespace checks, ownership checks, and line-count checks is safe. Screenshot capture overwrites only files under `plans/29-catapult-tower/screenshots/`.

If Catapult projectiles damage only one enemy, inspect whether the projectile copied the template area-impact flag. If Catapult damages enemies in neighboring Tiles, inspect the world-to-Tile conversion helper and add boundary tests. If a placed Catapult does not render or fire, inspect feature-to-template and feature-to-sprite switch statements first.

## Artifacts and Notes

The generated source images are retained under `/home/dave/.codex/generated_images/019e6280-a4a3-7741-8e06-4fdbb67030b7/`. The project-bound processed PNGs live at `assets/sprites/structures/catapult-tower.png` and `assets/sprites/structures/catapult-tower-projectile.png`.

Revision note, 2026-05-26: Updated the living plan after implementation with completed progress, validation results, screenshot paths, and the required line-count review.

## Interfaces and Dependencies

No new Go module dependencies are required. Use the existing Ebitengine image APIs and the existing `assets.Catalog` pattern.

At completion, these interfaces must exist:

    type StructureCatalog struct {
        Sanctum        StructureTemplate
        BowTower       StructureTemplate
        FlameBoltTower StructureTemplate
        CatapultTower  StructureTemplate
    }

    type StructureTemplate struct {
        Name                          string
        Sprite                        *ebiten.Image
        Cost                          Resources
        RangeTiles                    float64
        Damage                        int
        FireIntervalSeconds           float64
        ProjectileSpeedTilesPerSecond float64
        ProjectileSprite              *ebiten.Image
        DamageAllEnemiesInTargetTile  bool
    }

    type ProjectileSprites struct {
        BowTowerProjectile       *ebiten.Image
        FlameBoltTowerProjectile *ebiten.Image
        CatapultTowerProjectile  *ebiten.Image
    }

    type StructureSprites struct {
        Sanctum        *ebiten.Image
        BowTower       *ebiten.Image
        FlameBoltTower *ebiten.Image
        CatapultTower  *ebiten.Image
    }
