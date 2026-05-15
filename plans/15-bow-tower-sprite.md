# Add Bow Tower Sprite

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

## Purpose / Big Picture

After this change, the project has its first tower sprite: a 64x64 pixel-art Bow Tower PNG loaded through the runtime asset catalog. This does not add tower placement, targeting, combat, costs, or build UI. It only makes the asset available to future tower-defense slices in the same typed catalog that already loads the Sanctum, tree, and skeleton sprites.

## Progress

- [x] (2026-05-15T21:08:55Z) Confirmed scope: Bow Tower sprite, asset plus catalog integration only.
- [x] (2026-05-15T21:08:55Z) Generated a Bow Tower source image with a flat chroma-key background.
- [x] (2026-05-15T21:10:00Z) Converted the generated source into `assets/sprites/structures/bow-tower.png` as a transparent 64x64 PNG.
- [x] (2026-05-15T21:11:00Z) Added the Bow Tower sprite to the typed asset catalog and embed list.
- [x] (2026-05-15T21:11:00Z) Added an asset catalog test that verifies the Bow Tower sprite loads and is 64x64.
- [x] (2026-05-15T21:12:00Z) Ran validation commands: `go test ./...`, `git diff --check`, image inspection, and `git status --short`.
- [x] (2026-05-15T21:12:00Z) Checked hand-written code-file line counts; no Go files exceed the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: The generated source image was 1254x1254 RGB rather than 64x64 RGBA.
  Evidence: `file` reported `PNG image data, 1254 x 1254, 8-bit/color RGB`, so the final asset was cropped, scaled, and saved as a 64x64 transparent PNG.

- Observation: The generated background was close to green but not exactly `#00ff00`.
  Evidence: border pixels included values such as `(23, 240, 31)`, and the chroma helper reported `Key color: #07f811`.

## Decision Log

- Decision: Keep this slice asset-only plus catalog integration.
  Rationale: The current game has no tower placement or combat system, so displaying or simulating a tower would mix art integration with gameplay design that belongs in a separate slice.
  Date/Author: 2026-05-15 / Codex

- Decision: Use the Bow Tower as the first tower sprite.
  Rationale: `GAME.md` already records Bow Tower as the first defined tower type, so this asset supports existing design intent without introducing a new tower archetype.
  Date/Author: 2026-05-15 / Codex

## Outcomes & Retrospective

Completed. The project now has `assets/sprites/structures/bow-tower.png`, a 64x64 RGBA Bow Tower sprite with transparent corners and a non-empty subject. `assets/catalog.go` embeds and exposes it as `Catalog.Sprite.Structure.BowTower`, and `assets/catalog_test.go` verifies the sprite loads at 64x64.

Validation passed: `go test ./...`, `git diff --check`, `file assets/sprites/structures/bow-tower.png`, and alpha inspection all succeeded. The hand-written Go line-count review found no files over 600 lines; the largest was `internal/game/game_test.go` at 543 lines, so no unplanned refactor is needed.

## Context and Orientation

The repository is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. Runtime sprites live under `assets/sprites/`, and `assets/catalog.go` embeds and loads required PNG files into a typed `Catalog`. Existing structure sprites include `assets/sprites/structures/sanctum.png`; terrain and enemy sprites use the same catalog path.

`ART.md` says generated assets are 2D pixel art PNGs, 64x64, bright, high contrast, low detail, and without lighting or shadows. Structure sprites belong under `assets/sprites/structures/`. `GAME.md` defines the Bow Tower as the first tower type, but tower costs, range, damage, targeting, upgrades, and additional tower types remain future design work.

`PRODUCT.md` should not change because this asset will not become visible or usable in the current app. `GAME.md`, `DESIGN.md`, `ART.md`, `ROADMAP.md`, and `ARCHITECTURE.md` do not need updates because this plan follows existing decisions and constraints.

## Plan of Work

Generate one Bow Tower image using the built-in image generation path with a flat `#00ff00` chroma-key background. Remove the chroma-key background locally, resize the result to exactly 64x64, and save it as `assets/sprites/structures/bow-tower.png`.

Update `assets/catalog.go` so the new PNG is included in the `//go:embed` list, loaded by `NewCatalog`, and exposed as `StructureSprites.BowTower`.

Update `assets/catalog_test.go` with a focused test that constructs `NewCatalog`, confirms `catalog.Sprite.Structure.BowTower` is non-nil, and confirms its bounds are exactly 64x64.

## Concrete Steps

Run these commands from the repository root.

1. Generate a Bow Tower image with this prompt:

       Create a 64x64 2D pixel-art Bow Tower sprite for a medieval wizardry tower-defense game.
       The tower is a compact wooden archer tower with a small stone base, simple crenellations, and a clear bow/arrow identity.
       Use bright readable colors, sharp contrast, low detail, no lighting, no shadows, no text, no watermark.
       Center the full structure with generous padding on a perfectly flat solid #00ff00 chroma-key background.
       Do not use #00ff00 anywhere in the subject.

2. Remove the chroma-key background using `python3` and save the final transparent PNG at `assets/sprites/structures/bow-tower.png`.

3. Edit `assets/catalog.go` and `assets/catalog_test.go` as described in Plan of Work.

4. Run:

       go test ./...
       git diff --check
       file assets/sprites/structures/bow-tower.png
       git status --short

5. Check hand-written code-file line counts:

       find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

   Report any hand-written code file over 600 lines and recommend a concrete response. Do not perform unplanned refactors, code splits, or library additions without user approval.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, `git diff --check` reports no whitespace errors, `assets/sprites/structures/bow-tower.png` exists as a 64x64 PNG with transparency, and `NewCatalog` exposes a non-nil `StructureSprites.BowTower` image whose bounds are exactly 64x64.

## Idempotence and Recovery

Generating a new candidate image is safe because generated source images stay under `$CODEX_HOME/generated_images/`. Re-running the chroma-key removal and resize step overwrites only `assets/sprites/structures/bow-tower.png`, the intended final asset for this plan. If catalog edits fail, remove the Bow Tower fields and embed entry or restore them from git before retrying.

## Artifacts and Notes

Generated source image:

    /root/.codex/generated_images/019e2d74-1b45-7902-831c-f145c2f93720/ig_034d6e270b1c1c3d016a078b456f008198af5a3e0399cf4b31.png

Final asset path:

    assets/sprites/structures/bow-tower.png

## Interfaces and Dependencies

Use only the existing Go standard library image decoding plus Ebitengine image conversion already present in `assets/catalog.go`. The public catalog shape after this change must include:

    type StructureSprites struct {
        Sanctum  *ebiten.Image
        BowTower *ebiten.Image
    }

No new Go module dependencies are required.
