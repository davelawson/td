# Replace Boulder Vector With Terrain PNG Sprites

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/46-boulder-terrain-sprites.md`.

## Purpose / Big Picture

After this change, Boulder terrain uses sprite art that matches the existing pine-tree terrain assets instead of a vector placeholder. Newly explored grasslands Plots still generate sparse Boulder terrain, but those Tiles render with one of several 64x64 transparent PNG boulder variants, selected and horizontally flipped from the Tile tweak in the same way trees are.

## Progress

- [x] (2026-07-13T22:44Z) Created this ExecPlan from the accepted plan.
- [x] (2026-07-13T22:52Z) Generated and post-processed four 64x64 transparent Boulder terrain PNGs.
- [x] (2026-07-13T22:57Z) Loaded Boulder sprites through the asset catalog.
- [x] (2026-07-13T23:00Z) Rendered Boulder terrain with sprite variants and tweak-based horizontal flipping.
- [x] (2026-07-13T23:04Z) Updated tests and root control documents.
- [x] (2026-07-13T23:09Z) Ran asset validation, `go test ./...`, screenshot capture, `git diff --check`, ownership check, and Go file line-count review.

## Surprises & Discoveries

- Observation: This plan builds on uncommitted plans 44 and 45.
  Evidence: `git status --short` shows uncommitted grasslands generation and Boulder terrain files.
- Observation: Existing terrain assets are four 64x64 RGBA pine-tree PNGs.
  Evidence: `file assets/sprites/terrains/*.png` reports the pine tree files as `64 x 64, 8-bit/color RGBA`.
- Observation: Final Boulder sprites have transparent corners and non-empty alpha subject bounds.
  Evidence: local alpha validation reported transparent corner alpha values of zero for all four final 64x64 Boulder PNGs.
- Observation: The line-count review still finds one test file over the 600-line preference and one close to it.
  Evidence: final review reported `602 internal/game/game_test.go` and `596 internal/game/building_bar_test.go`.

## Decision Log

- Decision: Use four Boulder PNG variants.
  Rationale: This matches the existing pine-tree terrain variant count and the user's request for multiple PNG versions.
  Date/Author: 2026-07-13 / Codex
- Decision: Use runtime horizontal flipping from `Tile.Tweak` rather than storing flipped duplicate files.
  Rationale: This matches the existing tree behavior and keeps asset count focused on distinct variants.
  Date/Author: 2026-07-13 / Codex
- Decision: Generate transparent project-bound PNGs with the built-in image generation path plus local chroma-key removal.
  Rationale: The imagegen skill recommends this path for simple transparent assets without requiring CLI fallback.
  Date/Author: 2026-07-13 / Codex

## Outcomes & Retrospective

Implementation completed the sprite-backed Boulder terrain slice. Four generated Boulder terrain PNGs now live under `assets/sprites/terrains/`, the asset catalog embeds and loads them into `TerrainSprites.Boulders`, and Boulder terrain renders with a selected variant plus high-bit horizontal flipping from `Tile.Tweak`. The old vector Boulder drawing path was removed. Tree rendering now uses the same generalized terrain sprite helper functions.

Validation passed on 2026-07-13: `file assets/sprites/terrains/boulder-*.png`, `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded. Screenshot evidence was written under `plans/46-boulder-terrain-sprites/screenshots/`.

The final hand-written Go file line-count review found `internal/game/game_test.go` over the 600-line preference at 602 lines and `internal/game/building_bar_test.go` close to the preference at 596 lines. Recommended follow-up remains a separate approved test-organization change that splits map/default-state tests out of `game_test.go` and considers splitting building-bar rendering tests by category or behavior before adding more cases. No unplanned split was performed.

## Context and Orientation

The current working tree already has `terrainBoulder` from `plans/45-boulder-terrain.md`. It renders as a vector shape in `internal/game/scene.go`. Existing Forest terrain uses four pine-tree PNGs under `assets/sprites/terrains/`, loaded by `assets/catalog.go` into `TerrainSprites.PineTrees`, and rendered by `drawPineTree` in `internal/game/scene.go`. The helper functions `pineTreeSpriteIndex` and `treeSpriteFlipped` choose variants and horizontal mirroring from `Tile.Tweak`.

`ART.md` says generated assets should be 2D pixel-art PNGs, 64x64 for terrain, clearly visible at 50% scale, sharp contrast, low detail, bright colors, and no lighting or shadows. New terrain assets belong under `assets/sprites/terrains/`.

## Plan of Work

Generate four isolated boulder terrain sprites on a flat chroma-key background, remove the background locally, validate alpha and dimensions, and save them as `assets/sprites/terrains/boulder-1.png` through `boulder-4.png`. The sprites should be gray stone boulders or small boulder clusters, low-detail and high-contrast, visually aligned with the existing pine-tree terrain sprites.

Update `assets/catalog.go` to embed and load the Boulder sprites into `TerrainSprites.Boulders [4]*ebiten.Image`. Update asset catalog tests to verify all four Boulder sprites load and are 64x64.

Update `internal/game/scene.go` so Boulder terrain draws with a sprite instead of vector circles. Rename the tree-specific tweak helpers to terrain-general helpers and use them for both trees and Boulders. Keep the same high-bit horizontal flip behavior.

Update `ART.md`, `ARCHITECTURE.md`, and `DESIGN.md` to say Boulder terrain now uses generated transparent PNG sprites. Product and game behavior docs do not need semantic changes because Boulder behavior remains the same.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Generate four Boulder source images using built-in image generation with existing pine-tree sprites as style references.
2. Copy the generated images into a temporary project folder, remove chroma-key backgrounds with:

    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input <source> --out <dest> --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill

3. Resize or normalize final images to 64x64 RGBA PNGs if needed.
4. Update asset catalog, scene rendering, tests, and docs.
5. Validate:

    file assets/sprites/terrains/boulder-*.png
    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/scene.go internal/game/game_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

If any hand-written Go file exceeds 600 lines, record it in `Outcomes & Retrospective` with a concrete recommendation. Do not do an unplanned split unless the user approves it.

## Validation and Acceptance

Acceptance is met when all four Boulder PNGs exist as 64x64 RGBA terrain sprites, asset loading tests pass, Boulder terrain renders through the sprite catalog, and the same tweak-derived variant and horizontal flipping behavior applies to both Forest and Boulder terrain. Screenshot evidence must be refreshed under `plans/46-boulder-terrain-sprites/screenshots/`.

## Idempotence and Recovery

Generated source images may be discarded after final transparent PNGs are saved. Re-running catalog loading, tests, and screenshot capture is safe. If chroma-key removal leaves visible fringe, retry once with a slightly stronger local removal option before accepting the asset.

## Artifacts and Notes

Final workspace assets should live at:

    assets/sprites/terrains/boulder-1.png
    assets/sprites/terrains/boulder-2.png
    assets/sprites/terrains/boulder-3.png
    assets/sprites/terrains/boulder-4.png

Screenshot evidence should live under `plans/46-boulder-terrain-sprites/screenshots/`.

## Interfaces and Dependencies

No new runtime dependencies are required. Asset generation uses the built-in image generation tool and local post-processing only during development. Runtime loading continues through the existing embedded asset catalog.
