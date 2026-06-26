# Add a House building that grants Peasant population

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root and is saved at `plans/33-house-building.md`.

## Purpose / Big Picture

After this change, a player can build a new House from the existing building bar during calm play. The House costs 20 Wood, requires no staff, can be placed on the same empty grass-like Tiles as towers, and immediately grants 2 Peasants by increasing both available and total Peasant population. This creates the first normal-play path from the current `0/0` Peasant state toward staffing Catapult Towers.

The House is a population-provider structure, not a combat tower. It does not fire projectiles, reserve inhabitants, recruit over time, create a settlement-management screen, or release population if removed. Removal is not implemented in the current prototype.

## Progress

- [x] (2026-06-26T00:40:32Z) Inspected current structure templates, tile features, building bar placement, population state, asset catalog, selection panels, tests, screenshots, and control documents.
- [x] (2026-06-26T00:40:32Z) Created a generated 64x64 House sprite asset under `assets/sprites/structures/house.png`.
- [x] (2026-06-26T00:40:32Z) Created this ExecPlan with the accepted narrow gameplay slice.
- [x] (2026-06-26T00:51:26Z) Added House template metadata, asset catalog loading, tile feature rendering, building-bar display, placement mapping, and Peasant population grant behavior.
- [x] (2026-06-26T00:51:26Z) Added focused tests for House catalog data, bar ordering, placement cost, population grant, invalid drop safety, selection, and panel rows.
- [x] (2026-06-26T00:51:26Z) Updated screenshot capture output to this plan and captured House visual evidence.
- [x] (2026-06-26T00:51:26Z) Updated `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, `README.md`, and `ART.md` to match implemented House behavior and asset provenance.
- [x] (2026-06-26T00:51:26Z) Ran formatting, automated tests, screenshot capture, whitespace, ownership, status, and hand-written code-file line-count checks.

When this plan changes code, the final progress entry must include the required line-count review for hand-written code files and report any file over the 600-line preference before proposing user-approved follow-up work.

## Surprises & Discoveries

- Observation: The current building bar only distinguishes resource costs and staff requirements; it has no way to display positive population effects.
  Evidence: `internal/game/building_bar.go` renders a cost row and one staffing row below each icon.

- Observation: A fourth building item fits within the existing full-height left building bar at the default 1920x1080 size.
  Evidence: The current item block is 104 pixels tall plus a 12-pixel gap. With four items starting below the 86-pixel top bar, the fourth block ends around screen Y 554, well above the 1080-pixel bottom edge.

- Observation: The House screenshot evidence confirms both placement and the population/resource mutation.
  Evidence: `plans/33-house-building/screenshots/placed-house.png` shows a placed House, Wood reduced to `80`, and Peasants increased to `2/2`.

- Observation: One existing test file is now very close to the 600-line code-size preference.
  Evidence: The final line-count review reported `internal/game/game_test.go` at 597 lines.

## Decision Log

- Decision: Treat "20 woods" as 20 Wood, the existing singular resource name used by the HUD and structure costs.
  Rationale: The implemented resource is named `Wood`; no separate `woods` resource exists.
  Date/Author: 2026-06-26 / Codex

- Decision: The House requires no staffing and grants 2 Peasants immediately on successful placement by increasing both available and total counts.
  Rationale: The user's requested benefit is population, and increasing only total would not help the current staffing gate. Immediate grant keeps the first population source observable without adding recruitment timers or assignment screens.
  Date/Author: 2026-06-26 / Codex

- Decision: The House is not a combat tower and has no range, damage, fire interval, projectile speed, projectile sprite, or area-damage behavior.
  Rationale: The requested building is an economy/base-building structure. Combat behavior should remain limited to the three tower templates.
  Date/Author: 2026-06-26 / Codex

- Decision: The building bar will use the existing metadata row for either staff requirements or population grants, showing House as `+2` beside the Peasant icon.
  Rationale: This exposes the user-visible effect without expanding the compact bar into a larger command surface.
  Date/Author: 2026-06-26 / Codex

## Outcomes & Retrospective

Implemented the House as the first population-provider structure. The typed asset catalog now embeds `assets/sprites/structures/house.png`. `internal/game` now has a House template that costs 20 Wood, requires no staffing, grants 2 Peasants, places as `featureHouse`, renders with the generated House sprite, appears first in the building bar with `20` Wood and `+2` Peasants metadata, and shows non-combat selection rows for cost and Peasant grant. Successful House placement changes a new game from 100 Wood and `0/0` Peasants to 80 Wood and `2/2` Peasants. Invalid drops leave resources and populations unchanged.

Screenshot evidence was captured under `plans/33-house-building/screenshots/`. The key files are `house-icon.png`, which shows the new bar item, and `placed-house.png`, which shows a placed House and the resulting Peasant population.

Validation completed:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

All tests passed. Whitespace and ownership checks produced no output.

Final hand-written Go file line-count review found no file over 600 lines. `internal/game/game_test.go` is at 597 lines, close enough that the next likely game-state test addition may exceed the preference. Recommended response: split map/default-state tests out of `game_test.go` into a focused file such as `map_test.go` or `status_test.go` during the next related test change. No unplanned split was performed in this plan.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. The current game starts with a Sanctum-only home Plot, 100 Wood, 50 Stone, 20 Metal, and `0/0` Apprentices, Soldiers, and Peasants. Existing construction goes through `internal/game/building_bar.go`: left-dragging an eligible building-bar icon during calm play and releasing over an empty grass-like Tile places a structure. Eligibility currently requires enough resources and enough available staff. Successful placement deducts resources, reserves required staff, and writes a `tileFeature` into the home Plot.

`internal/game/structures.go` defines structure templates. `internal/game/map.go` defines Tile features. `internal/game/scene.go` renders features. `internal/game/hud.go` owns private resource and population counts. `internal/game/selection_panel.go` formats selected structure details. `assets/catalog.go` embeds and loads the runtime PNG assets. `cmd/td/main_test.go` captures opt-in screenshots under the active plan directory.

The House should use these same boundaries. It belongs in the structure catalog and asset catalog, appears in the building bar, uses the existing placement rules, and is selectable like other structures. Population count mutation belongs in `internal/game/hud.go` or a small method on `State` because that file already owns population availability and reservation.

## Plan of Work

First, wire the House asset into `assets/catalog.go` by adding `sprites/structures/house.png` to the embed list, adding a `House *ebiten.Image` field to `StructureSprites`, loading it in `NewCatalog`, and adding a catalog test that proves it is a 64x64 sprite.

Next, extend gameplay data. Add a `House` template to `StructureCatalog` with name `House`, sprite `assetCatalog.Sprite.Structure.House`, cost `Resources{Wood: 20}`, no staffing, and a new population grant value of 2 Peasants. Add a `featureHouse` tile feature. Render `featureHouse` with `drawStructureSprite`. Keep `combatTowerTemplate` unchanged except that its default path continues to ignore non-combat structures.

Then update construction. Add House as the first building-bar template before the combat towers. Map the first building-bar item to `featureHouse`; the combat towers follow in their prior relative order. On successful placement, deduct cost, reserve any staffing requirements, apply any population grant, and then write the tile feature. For the House, this spends 20 Wood, reserves no staff, and changes Peasants from `0/0` to `2/2` in a new game. Invalid releases must leave both resources and populations unchanged.

Expose the effect in UI. The building bar should show the House cost as `20` in Wood color and show `+2` beside the Peasant icon in the metadata row. Selection should show a simple structure panel for House with rows for `Structure: House`, `Cost: 20 Wood`, and `Grants Peasants: 2`; do not show tower combat stats for House.

Finally, update tests, screenshots, and durable control documents. Tests should prove the template values, bar ordering, House placement, population grant, invalid release safety, selection support, and panel rows. Screenshot capture should write under `plans/33-house-building/screenshots/` and include at least one evidence frame showing the House icon in the bar and one showing a placed House with Peasants increased to `2/2`. Update `PRODUCT.md` and `GAME.md` because current user-visible behavior and gameplay design change. Update `ROADMAP.md` because the near-term population-seeding workflow is no longer entirely missing. Update `ARCHITECTURE.md` because structure ownership now includes population-provider structures. Update `README.md` if its current-state text mentions no normal population path. Update `ART.md` only if the asset creation notes need to record the House prompt/source pattern.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Inspect state before editing:

    git status --short
    rg -n "CatapultTower|featureCatapultTower|buildingBarItems|Population|TD_CAPTURE_SCREENSHOT" internal assets cmd README.md PRODUCT.md GAME.md ROADMAP.md ARCHITECTURE.md ART.md

After editing Go code, format and test:

    gofmt -w assets/*.go internal/game/*.go cmd/td/main_test.go
    go test ./...

Capture visual evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Run final repository checks:

    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short

End with the required hand-written code-file line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written code file exceeds 600 lines, report the file and line count in `Outcomes & Retrospective`, recommend a concrete response such as a responsibility-based split or helper extraction, and ask the user to approve extra work before implementing it unless this plan already included that split.

## Validation and Acceptance

`go test ./...` must pass. Focused tests must prove that the House template is named `House`, costs 20 Wood, uses the loaded House sprite, requires no staff, and grants 2 Peasants. Building-bar tests must prove House appears before the three combat towers, its cost and population grant metadata are exposed, and all four items fit in the bar. Placement tests must prove a new-game House drag can start with zero staff, a valid drop places `featureHouse`, Wood decreases from 100 to 80, Peasants increase from `0/0` to `2/2`, and invalid drops spend neither Wood nor population effects.

Selection tests must prove House tiles can be selected. Selection-panel tests must prove House shows basic structure and population-grant rows rather than tower combat rows. Screenshot evidence under `plans/33-house-building/screenshots/` must show the House icon in the building bar and a placed House with Peasant population increased to `2/2`.

Documentation must match the implemented behavior: current product truth should say House construction is the first normal way to add Peasants, game design should record the House cost and grant, roadmap should distinguish this population source from future recruitment, and architecture should continue keeping population mutation inside `internal/game`.

## Idempotence and Recovery

Formatting, tests, and screenshot capture are safe to repeat. Screenshot capture overwrites only files under this plan's screenshot directory. The generated source image remains under `/home/dave/.codex/generated_images/019f015d-4662-7db0-be1e-41f0b4c9fa11/`, and the project-bound final asset is `assets/sprites/structures/house.png`. If the asset needs replacement later, create a new file or intentionally replace `house.png` in a separate reviewed change.

## Artifacts and Notes

The House sprite was generated with the built-in image generation tool, then converted from a flat green chroma-key background to alpha and resized to 64x64 PNG. The final project asset is:

    assets/sprites/structures/house.png

The source generated image remains at:

    /home/dave/.codex/generated_images/019f015d-4662-7db0-be1e-41f0b4c9fa11/ig_05cb9feef7948c5a016a3dca3b30a88193b2f732eb76d9493b.png

Final prompt summary: a centered low-detail pixel-art-inspired rustic medieval timber-and-stone peasant house on a flat `#00ff00` chroma-key background, with no scenery, people, text, watermark, shadow, or green in the subject.

## Interfaces and Dependencies

No new Go module dependency is required. `StructureCatalog` gains a `House StructureTemplate`. `StructureTemplate` gains a small population-grant value, shaped like:

    type PopulationGrant struct {
        Apprentices int
        Soldiers    int
        Peasants    int
    }

    type StructureTemplate struct {
        ...
        PopulationGrant PopulationGrant
    }

The field is zero for Sanctum and combat towers. For House it is `PopulationGrant{Peasants: 2}`. The existing `StaffingRequirements` field remains the source of construction staffing requirements; do not reuse it for population grants because requirements and grants are opposite gameplay concepts.
