# Add Asset Catalog

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/11-asset-catalog.md` because `plans/00-initial-ebitengine-menu.md` through `plans/10-basic-camera.md` already exist.

## Purpose / Big Picture

The prototype already has a Sanctum PNG under `assets/sprites/structures/sanctum.png`, but the running game still draws the Sanctum with procedural shapes and text. After this change, the game has a typed asset catalog that loads required runtime assets when a new game starts, stores that catalog on the game state, and renders the Sanctum from the catalog. A contributor can see the change by starting a new game and observing the sprite-based Sanctum in the home Plot.

## Progress

- [x] (2026-05-13 00:00Z) Inspected the current game construction path, scene rendering, tests, asset file, and control documents.
- [x] (2026-05-13 00:00Z) Created this ExecPlan before implementation.
- [x] (2026-05-13 03:46Z) Added the asset catalog package and tests.
- [x] (2026-05-13 03:46Z) Wired the catalog into `internal/game.State` and rendered the Sanctum sprite.
- [x] (2026-05-13 03:46Z) Updated screenshot capture and control documents.
- [x] (2026-05-13 03:46Z) Ran full validation, screenshot capture, bounded launch check, whitespace check, git status, and final hand-written Go file line-count review.

## Surprises & Discoveries

- Observation: The bounded launch check exits with code 124 because the app successfully stays open until `timeout` stops it.
  Evidence: `timeout 5s go run ./cmd/td` produced no startup error and exited after the timeout window.

## Decision Log

- Decision: Create the catalog package at `assets` with import path `td/assets`.
  Rationale: The existing PNG already lives under the repository root `assets/` directory, and Go's `embed` package can only embed files from the package directory or its subdirectories. Using `assets` also keeps asset-loading details near the runtime files.
  Date/Author: 2026-05-13 / Codex

- Decision: Store ready-to-draw `*ebiten.Image` values in the catalog.
  Rationale: The first consumer is Ebitengine rendering code, so converting once during catalog creation keeps per-frame rendering simple and avoids decoding assets during draw calls.
  Date/Author: 2026-05-13 / Codex

- Decision: Fail `game.New` if the required catalog cannot load.
  Rationale: The Sanctum sprite is a required asset for the current game scene. A startup error is clearer than silently drawing missing or inconsistent content.
  Date/Author: 2026-05-13 / Codex

## Outcomes & Retrospective

Implementation completed the asset catalog slice. New game construction now creates a typed runtime asset catalog, stores it on `internal/game.State`, and renders the centered Sanctum from `assets/sprites/structures/sanctum.png` instead of the old procedural marker. Menu flow, camera controls, pause behavior, top bar rendering, and in-game overlay behavior remain unchanged.

Validation results:

    go test ./assets
    ok  	td/assets	0.012s

    go test ./internal/game
    ok  	td/internal/game	0.026s

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	0.014s
    ok  	td/internal/game	0.025s
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    ok  	td/cmd/td	0.688s

    file plans/11-asset-catalog/screenshots/*.png
    plans/11-asset-catalog/screenshots/ingame-menu.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/11-asset-catalog/screenshots/main-menu.png:              PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/11-asset-catalog/screenshots/new-game-configuration.png: PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/11-asset-catalog/screenshots/paused-game.png:            PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced
    plans/11-asset-catalog/screenshots/running-game.png:           PNG image data, 1920 x 1080, 8-bit/color RGB, non-interlaced

    timeout 5s go run ./cmd/td
    Exit code: 124 after the app launched and stayed open with no startup error. This is expected for a bounded launch check.

    git diff --check
    No whitespace errors.

Final hand-written Go file line-count review:

    17 internal/ui/widgets.go
    20 assets/catalog_test.go
    20 internal/ui/colors.go
    22 internal/game/colors.go
    29 internal/ui/text.go
    53 internal/game/map.go
    57 assets/catalog.go
    59 internal/game/scene.go
    87 internal/game/camera.go
    100 internal/game/hud.go
    129 internal/menu/start.go
    157 cmd/td/main.go
    159 internal/game/ingamemenu.go
    181 cmd/td/main_test.go
    182 internal/game/game.go
    282 internal/menu/menu_test.go
    343 internal/menu/menu.go
    465 internal/game/game_test.go
    2362 total

No hand-written Go file exceeds or approaches the 600-line preference, so no extra split or user-approved refactor is needed.

## Context and Orientation

`td` is a Go/Ebitengine local PC prototype. The executable entry point in `cmd/td/main.go` owns the desktop window, app mode, and input polling. `internal/game/game.go` owns `State`, and `game.New` currently creates font faces, a default map, a camera, and the in-game UI. `internal/game/scene.go` currently renders the home Plot and draws the centered Sanctum with Ebitengine vector shapes and an `S` label. The asset file `assets/sprites/structures/sanctum.png` exists and is a 64x64 RGBA PNG.

An asset catalog is a typed struct that groups loaded runtime assets by category. In this plan, the grouping is intentionally narrow and follows the requested path `catalog.Sprite.Structure.Sanctum`. `Sprite` means image assets drawn into the game scene. `Structure` means sprites for placeable or map structures. `Sanctum` is the existing central structure in the prototype home Plot.

The root control documents constrain this work. `PLANS.md` requires this self-contained ExecPlan and final validation evidence. `CODESTYLE.md` requires `gofmt`, doc comments for Go functions and methods, tests for pure behavior when possible, and a final hand-written code-file line-count review against the 600-line preference. `ARCHITECTURE.md` must be updated because the durable package boundary for asset loading changes. `README.md` and `PRODUCT.md` must be updated because the current user-visible game scene changes from procedural-only rendering to sprite-backed Sanctum rendering.

## Plan of Work

Create `assets/catalog.go` in package `assets`. Use `//go:embed sprites/structures/sanctum.png` to include the existing PNG in the compiled binary. Define `Catalog`, `SpriteCatalog`, and `StructureSprites` so game code can access the Sanctum as `catalog.Sprite.Structure.Sanctum`. Define `NewCatalog() (Catalog, error)` to decode the embedded PNG and convert it to `*ebiten.Image`. Add a small private helper if needed to keep decoding focused and easy to test.

Create `assets/catalog_test.go`. The test should call `NewCatalog`, fail on error, verify `catalog.Sprite.Structure.Sanctum` is not nil, and verify the image bounds are 64 by 64. This proves the embedded asset path, PNG decoding, and catalog grouping work without opening a game window.

Update `internal/game/game.go`. Import `td/assets`, add a private `assetCatalog assets.Catalog` field to `State`, call `assets.NewCatalog()` inside `New`, and store the resulting catalog in the returned state. Preserve the existing font-source initialization, default map, camera, UI layout, status setup, and error-return behavior.

Update `internal/game/scene.go`. Replace the procedural Sanctum circle and `S` label with sprite drawing from `s.assetCatalog.Sprite.Structure.Sanctum`. Scale the 64x64 sprite to fit comfortably inside the projected tile rectangle, centered in the tile. Keep the rest of the home Plot, tile colors, road rendering, camera projection, top bar, pause label, and in-game overlay behavior unchanged.

Update tests in `internal/game/game_test.go` to assert that a new state stores a non-nil Sanctum sprite in its asset catalog. Update `cmd/td/main_test.go` screenshot capture output from `plans/10-basic-camera/screenshots` to `plans/11-asset-catalog/screenshots`, because rendered output changes. Update `README.md`, `PRODUCT.md`, and `ARCHITECTURE.md` to describe the current asset catalog and sprite-backed Sanctum. No new dependency is needed.

## Concrete Steps

From the repository root, inspect the working tree:

    git status --short

Create and edit:

    plans/11-asset-catalog.md
    assets/catalog.go
    assets/catalog_test.go
    internal/game/game.go
    internal/game/scene.go
    internal/game/game_test.go
    cmd/td/main_test.go
    README.md
    PRODUCT.md
    ARCHITECTURE.md

Format and test incrementally:

    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/game.go internal/game/scene.go internal/game/game_test.go cmd/td/main_test.go
    go test ./assets
    go test ./internal/game

Run full validation and capture visual evidence:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    file plans/11-asset-catalog/screenshots/*.png
    timeout 5s go run ./cmd/td
    git diff --check
    git status --short

Check hand-written Go file line counts at the end:

    rg --files cmd internal assets 2>/dev/null | grep -E '\.go$' | xargs -r wc -l | sort -n

If any hand-written Go file exceeds 600 lines, record the path and line count in `Outcomes & Retrospective`, recommend a concrete response, and ask the user for approval before implementing an unplanned split, refactor, or library addition.

## Validation and Acceptance

The implementation is accepted when `go test ./...` succeeds, screenshot capture writes PNG evidence under `plans/11-asset-catalog/screenshots/`, `go run ./cmd/td` opens without startup errors, `git diff --check` reports no whitespace errors, and the line-count review finds no hand-written Go file over the 600-line preference.

A human should be able to start a new game and see the same static home Plot, road, camera controls, top bar, pause behavior, and in-game overlay behavior as before, with the centered Sanctum now rendered from the PNG sprite instead of the old procedural circle and letter marker.

Documentation is accepted when `README.md` and `PRODUCT.md` describe the sprite-backed Sanctum as current behavior, and `ARCHITECTURE.md` records that the root `assets` package owns typed runtime asset loading while `internal/game` owns gameplay state and rendering decisions.

## Idempotence and Recovery

The changes are additive and local to asset loading, game construction, scene rendering, tests, docs, screenshots, and this plan. Re-running `gofmt`, tests, and screenshot capture is safe. If asset decoding fails, confirm that `assets/sprites/structures/sanctum.png` still exists and remains a valid PNG. If screenshot capture fails because the local graphics environment cannot open an Ebitengine window, record the exact error in `Surprises & Discoveries`, keep automated tests as validation, and use manual launch validation when graphics are available.

## Artifacts and Notes

Important artifacts:

    plans/11-asset-catalog.md
    assets/catalog.go
    assets/catalog_test.go
    assets/sprites/structures/sanctum.png
    plans/11-asset-catalog/screenshots/running-game.png
    internal/game/game.go
    internal/game/scene.go
    internal/game/game_test.go
    cmd/td/main_test.go
    README.md
    PRODUCT.md
    ARCHITECTURE.md

## Interfaces and Dependencies

Use only the existing Go module, the standard library packages `embed` and `image/png`, and the existing Ebitengine dependency.

In `assets/catalog.go`, define:

    type Catalog struct {
        Sprite SpriteCatalog
    }

    type SpriteCatalog struct {
        Structure StructureSprites
    }

    type StructureSprites struct {
        Sanctum *ebiten.Image
    }

    func NewCatalog() (Catalog, error)

`internal/game.State` must store a private `assetCatalog assets.Catalog`. `game.New` must instantiate a fresh catalog for every new game state. `internal/game/scene.go` must read the Sanctum sprite from the stored state catalog during rendering.

## Revision Note

This plan was created before implementation to capture the requested asset catalog grouping, the new-game catalog instantiation behavior, the first visible sprite-backed Sanctum use, required documentation updates, screenshot evidence, validation commands, and the final code-file line-count review required by `PLANS.md` and `CODESTYLE.md`.
