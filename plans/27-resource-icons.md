# Add resource icon sprites to the top bar

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/27-resource-icons.md`.

## Purpose / Big Picture

The current top bar writes the prototype resource counts as the words `Wood`, `Stone`, and `Metal`. After this change, the player sees compact sprite-backed icons for the three resources, followed by their counts, while the fixed prototype resource values remain unchanged. This improves visual readability and starts using the asset catalog for economy-facing HUD art without adding gathering, spending, rewards, costs, or an asset pipeline.

The result is visible by starting a game with `go run ./cmd/td`: the top bar still shows Chapter, Day, phase status, and Barricade text, but the resource labels are replaced by Wood, Stone, and Metal icons.

## Progress

- [x] (2026-05-26T00:00:00Z) Inspected the current asset catalog, top-bar HUD code, tests, art guidance, and documentation references.
- [x] (2026-05-26T00:00:00Z) Created this ExecPlan at `plans/27-resource-icons.md`.
- [x] (2026-05-26T00:00:00Z) Added the three 64x64 resource icon PNG files under `assets/sprites/icons/`.
- [x] (2026-05-26T00:00:00Z) Extended the typed asset catalog and catalog tests for resource icons.
- [x] (2026-05-26T00:00:00Z) Updated the top-bar HUD to render icon-plus-count resource items and added focused tests for resource item order and values.
- [x] (2026-05-26T00:00:00Z) Updated user-visible and architecture documentation for sprite-backed resource icons.
- [x] (2026-05-26T00:00:00Z) Captured screenshot evidence under `plans/27-resource-icons/screenshots/`.
- [x] (2026-05-26T00:00:00Z) Ran `gofmt`, `go test ./...`, `git diff --check`, ownership checks, and recorded the validation results.
- [x] (2026-05-26T00:00:00Z) Checked hand-written code-file line counts against the 600-line preference and reported the existing overage.
- [x] (2026-05-26T00:00:00Z) Revised `assets/sprites/icons/metal.png` to read as three stacked, end-on bluish ingots and refreshed screenshot evidence.

## Surprises & Discoveries

- Observation: The existing top bar formats all right-side resource and barricade data as one measured string.
  Evidence: `internal/game/hud.go` uses `resourcesAndBarricadeText()` and right-aligns it with one `text.Measure` call.

- Observation: `ART.md` already defines `assets/sprites/icons/` as the correct location and 64x64 PNG as the correct size for icons.
  Evidence: `ART.md` lists `Icon: 64x64` and says generated graphical assets should live under `assets/sprites` in purpose-named subfolders including `icons`.

- Observation: The captured game screenshot shows the resource labels have been replaced by icons and counts without overlapping the centered phase text.
  Evidence: `plans/27-resource-icons/screenshots/running-game.png` shows the right side of the top bar as icon/count pairs followed by `| Barricade 3`.

- Observation: The revised metal icon remains readable at HUD size as a stack of bluish ingots.
  Evidence: `plans/27-resource-icons/screenshots/running-game.png` shows the metal count beside a small three-ingot stack.

## Decision Log

- Decision: Add a new `IconSprites` group to `assets.SpriteCatalog` rather than placing resource icons under terrain, structure, or projectile sprites.
  Rationale: Resource icons are HUD-facing assets, not map terrain or structures, and `ART.md` names `icons` as an expected asset category.
  Date/Author: 2026-05-26 / Codex

- Decision: Keep Barricade as text and replace only the words for Wood, Stone, and Metal.
  Rationale: The user requested icons for resources specifically, and the current game design names Barricade as a separate Sanctum defense status rather than a resource.
  Date/Author: 2026-05-26 / Codex

- Decision: Use deterministic 64x64 pixel-art PNGs generated locally for this prototype slice.
  Rationale: The icons are simple, need exact dimensions, and should be reproducible. This follows `ART.md` constraints without adding an art pipeline.
  Date/Author: 2026-05-26 / Codex

- Decision: Represent Metal as three stacked, end-on ingots while preserving the bluish tint.
  Rationale: The user requested this specific read for the icon, and the bluish tint differentiates metal from the gray Stone icon in the compact HUD.
  Date/Author: 2026-05-26 / Codex

## Outcomes & Retrospective

Implemented the requested resource icon display. The runtime asset catalog now embeds Wood, Stone, and Metal icons under `assets/sprites/icons/`, exposes them through `Catalog.Sprite.Icon`, and verifies all three load as 64x64 images. The game top bar now draws each resource as an icon plus count and keeps Barricade as text. Documentation now describes sprite-backed resource icons as current behavior without implying resource gathering or spending exists. The Metal icon was revised after initial implementation to show three stacked, end-on bluish ingots.

Validation completed:

    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/hud.go internal/game/hud_test.go internal/game/game_test.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    go test ./assets ./internal/game
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.(go)$' | xargs -r wc -l | sort -n

`go test ./...` passed. Screenshot evidence was written under `plans/27-resource-icons/screenshots/`, including `running-game.png`. `git diff --check` produced no output. The ownership check produced no output.

Line-count review found one hand-written code file above the 600-line preference: `internal/game/game_test.go` at 602 lines. This overage was pre-existing in recent plans and this implementation reduced it slightly by moving new HUD tests into `internal/game/hud_test.go`. Recommended follow-up is a responsibility-based split of remaining tests in `game_test.go`, such as moving general HUD or camera tests to focused files. No extra refactor was performed because that split is outside this feature scope and was not approved as part of this plan.

## Context and Orientation

The Go module is a local Ebitengine game prototype. `assets/catalog.go` embeds runtime files and exposes them through a typed `assets.Catalog`. `internal/game.State` owns the loaded catalog and renders the top bar from `internal/game/hud.go`. The top bar currently draws left Chapter/Day text, centered phase text, and right-aligned resources plus Barricade text. The prototype resource values are fixed in `setPrototypeGameStatus()` until economy systems exist.

Root documents constrain the change. `ART.md` says icons are 64x64 PNG pixel art, visible at 50% size, bright, sharp, low detail, and placed under `assets/sprites/icons/`. `GAME.md` identifies Wood, Stone, and Metal as the initial resources and says the top bar should show their counts. `README.md` and `PRODUCT.md` must describe the current user-visible top bar accurately. `ARCHITECTURE.md` must describe durable asset catalog and HUD ownership accurately. This change does not implement resource gathering, spending, tower costs, placement, rewards, saves, or a sprite atlas.

## Implementation Plan

First, create `assets/sprites/icons/wood.png`, `assets/sprites/icons/stone.png`, and `assets/sprites/icons/metal.png` as 64x64 RGBA PNG files. Use simple high-contrast pixel art on transparent backgrounds: wood should read as a small bundle/log, stone as a rock or quarry chunk, and metal as an ingot or ore. Keep each icon readable at about 28x28 in the HUD.

Second, update `assets/catalog.go`. Add the three icon paths to the `//go:embed` directive. Add `Icon IconSprites` to `SpriteCatalog`, define `IconSprites` with `Wood`, `Stone`, and `Metal` fields, load each file with `loadSprite`, and assign them into the returned `Catalog`. Do not let game code decode files by path.

Third, update `assets/catalog_test.go` with a test that calls `NewCatalog()` and verifies `catalog.Sprite.Icon.Wood`, `catalog.Sprite.Icon.Stone`, and `catalog.Sprite.Icon.Metal` are non-nil 64x64 images.

Fourth, update `internal/game/hud.go`. Replace the single resource text string path with a small HUD resource model that returns three items in order: Wood, Stone, Metal. Each item should include the loaded sprite and count. Draw the right side from right to left or by premeasuring a compact group so that the full group remains right-aligned to `topBarMargin`. Each resource item should draw a scaled icon around 28 pixels tall followed by the numeric count in the existing HUD font and text color. After the three resource items, draw `| Barricade 3` as text. Keep Chapter/Day and phase rendering unchanged.

Fifth, add or update tests in `internal/game`. Because `internal/game/game_test.go` is already over the 600-line preference, add focused HUD tests in a new file rather than growing that file. Test that the resource HUD items are ordered Wood, Stone, Metal, use counts 80, 45, and 12 on a new game, and use non-nil sprites from the catalog. Keep existing top-bar phase tests passing after removing the old resource text assertion.

Sixth, update documentation. `README.md` and `PRODUCT.md` should say the top bar shows fixed prototype resource counts with sprite-backed Wood, Stone, and Metal icons. `GAME.md` should say the top bar shows Wood, Stone, and Metal as icons plus counts. `ARCHITECTURE.md` should mention icon sprites in the asset catalog and sprite-backed resource display in the top bar. Do not update `ROADMAP.md`, `DESIGN.md`, `ART.md`, or `CODESTYLE.md` unless the implementation changes their durable truth.

Seventh, update the screenshot capture base path in `cmd/td/main_test.go` to `plans/27-resource-icons/screenshots/`, run the explicit screenshot capture command, and commit the resulting screenshots only if this work is being committed. The important acceptance image is `plans/27-resource-icons/screenshots/running-game.png`.

## Validation

Run these commands from `/home/dave/dev/ai/td`:

    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/hud.go internal/game/hud_test.go internal/game/game_test.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.(go)$' | xargs -r wc -l | sort -n

Acceptance is that `go test ./...` passes, the catalog test proves all three icons load as 64x64 images, the HUD test proves the order and counts are correct, `running-game.png` shows icons instead of resource words, `git diff --check` reports no whitespace errors, ownership output is empty, and the line-count review is recorded here.

## Idempotence and Recovery

Re-running asset generation should overwrite only the three icon PNG files in `assets/sprites/icons/`. Re-running screenshot capture overwrites only files under `plans/27-resource-icons/screenshots/`. If catalog loading fails, verify the icon files exist and the `//go:embed` paths match their repository-relative paths under `assets/`. If the HUD overlaps the phase text at 1920x1080, reduce icon display size or horizontal spacing while preserving icon-plus-count order.
