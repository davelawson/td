# Add Bow Tower Projectile Sprite

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

## Purpose / Big Picture

After this change, the project has a 64x64 transparent PNG arrow projectile sprite for future Bow Tower combat work. The sprite is not rendered or simulated yet; it is available through the runtime asset catalog so a later tower-defense slice can draw projectiles without first adding asset plumbing.

## Progress

- [x] (2026-05-16T22:37:11Z) Confirmed scope: Bow Tower projectile sprite, asset plus catalog integration only.
- [x] (2026-05-16T22:37:11Z) Generated an arrow projectile source image with a flat chroma-key background.
- [x] (2026-05-16T22:37:11Z) Converted the generated source into `assets/sprites/structures/bow-tower-projectile.png` as a transparent 64x64 PNG.
- [x] (2026-05-16T22:37:11Z) Added the projectile sprite to the typed asset catalog and embed list.
- [x] (2026-05-16T22:38:31Z) Added an asset catalog test that verifies the projectile sprite loads and is 64x64.
- [x] (2026-05-16T22:38:31Z) Ran validation commands: `go test ./...`, `git diff --check`, image inspection, and `git status --short`.
- [x] (2026-05-16T22:38:31Z) Checked hand-written code-file line counts; no file exceeds the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: The generated source image was 1254x1254 RGB rather than 64x64 RGBA.
  Evidence: `python3` inspection reported `RGB (1254, 1254)`, so the final asset was cropped, scaled, and saved as a 64x64 transparent PNG.

- Observation: The generated background was close to green but not exactly `#00ff00`.
  Evidence: the chroma helper reported `Key color: #15f611`, and border pixels included values such as `(42, 238, 31)`.

## Decision Log

- Decision: Keep the projectile asset in `assets/sprites/structures/` beside `bow-tower.png`.
  Rationale: The user asked for `bow-tower-projectile.png` to be placed alongside `bow-tower.png`, and this keeps tower-related prototype art together until a broader projectile asset layout is needed.
  Date/Author: 2026-05-16 / Codex

- Decision: Expose the asset through a new `ProjectileSprites` catalog group.
  Rationale: Future combat code should depend on gameplay role rather than storage location, and the user selected this catalog shape during planning.
  Date/Author: 2026-05-16 / Codex

## Outcomes & Retrospective

Completed. The project now has `assets/sprites/structures/bow-tower-projectile.png`, a 64x64 RGBA arrow projectile sprite with transparent corners and a non-empty subject. `assets/catalog.go` embeds and exposes it as `Catalog.Sprite.Projectile.BowTowerProjectile`, and `assets/catalog_test.go` verifies the sprite loads at 64x64.

Validation passed: `go test ./...`, `git diff --check`, `file assets/sprites/structures/bow-tower-projectile.png`, alpha inspection, and `git status --short` all completed successfully. The hand-written Go line-count review found no files over 600 lines. `internal/game/game_test.go` is now 572 lines, so the next change that touches that file should consider splitting test helpers or scenarios before adding more bulk.

## Context and Orientation

The repository is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. Runtime sprites live under `assets/sprites/`, and `assets/catalog.go` embeds and loads required PNG files into a typed `Catalog`. `assets/sprites/structures/bow-tower.png` already exists as the Bow Tower sprite. The new projectile PNG should sit next to it at `assets/sprites/structures/bow-tower-projectile.png`.

`ART.md` says generated assets are 2D pixel art PNGs, 64x64, bright, high contrast, low detail, and without lighting or shadows. `PRODUCT.md` should not change because this asset will not become visible or usable in the current app. `GAME.md`, `DESIGN.md`, `ART.md`, `ROADMAP.md`, `CODESTYLE.md`, and `ARCHITECTURE.md` do not need updates because this plan follows existing decisions and constraints.

## Plan of Work

Generate one arrow projectile image using the built-in image generation path with a flat chroma-key background. Remove the chroma-key background locally, resize the result to exactly 64x64, and save it as `assets/sprites/structures/bow-tower-projectile.png`.

Update `assets/catalog.go` so the new PNG is included in the `//go:embed` list, loaded by `NewCatalog`, and exposed as `Catalog.Sprite.Projectile.BowTowerProjectile`.

Update `assets/catalog_test.go` with a focused test that constructs `NewCatalog`, confirms `catalog.Sprite.Projectile.BowTowerProjectile` is non-nil, and confirms its bounds are exactly 64x64.

## Concrete Steps

Run these commands from the repository root.

1. Generate a Bow Tower projectile image with this prompt:

       Create a small arrow projectile sprite for a Bow Tower in a medieval wizardry tower-defense game.
       The projectile is a single wooden arrow with a steel arrowhead and simple feather fletching, angled diagonally from lower-left to upper-right as if flying.
       Use 2D pixel art, bright readable colors, very sharp contrast, low detail, no lighting, no shadows, no text, no watermark, no extra objects.
       Center the full arrow with generous padding on a perfectly flat solid #00ff00 chroma-key background.
       Do not use #00ff00 anywhere in the subject.

2. Remove the chroma-key background and save the final transparent PNG at `assets/sprites/structures/bow-tower-projectile.png`.

3. Edit `assets/catalog.go` and `assets/catalog_test.go` as described in Plan of Work.

4. Run:

       go test ./...
       git diff --check
       file assets/sprites/structures/bow-tower-projectile.png
       git status --short

5. Check hand-written code-file line counts:

       find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

   Report any hand-written code file over 600 lines and recommend a concrete response. Do not perform unplanned refactors, code splits, or library additions without user approval.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, `git diff --check` reports no whitespace errors, `assets/sprites/structures/bow-tower-projectile.png` exists as a 64x64 PNG with transparency, and `NewCatalog` exposes a non-nil `ProjectileSprites.BowTowerProjectile` image whose bounds are exactly 64x64.

## Idempotence and Recovery

Generating a new candidate image is safe because generated source images stay under `$CODEX_HOME/generated_images/`. Re-running the chroma-key removal and resize step overwrites only `assets/sprites/structures/bow-tower-projectile.png`, the intended final asset for this plan. If catalog edits fail, remove the projectile fields and embed entry or restore them from git before retrying.

## Artifacts and Notes

Generated source image:

    /root/.codex/generated_images/019e32ea-ab35-7202-8cf0-f0c36b8885c0/ig_0d4b2ddd4459eb80016a08f137af4c819abd4d4023d54db633.png

Final asset path:

    assets/sprites/structures/bow-tower-projectile.png

Image inspection after processing reported:

    PNG image data, 64 x 64, 8-bit/color RGBA, non-interlaced
    RGBA (64, 64) bbox (6, 6, 58, 58)
    corners [0, 0, 0, 0]

## Interfaces and Dependencies

Use only the existing Go standard library image decoding plus Ebitengine image conversion already present in `assets/catalog.go`. The public catalog shape after this change must include:

    type SpriteCatalog struct {
        Enemy      EnemySprites
        Projectile ProjectileSprites
        Structure  StructureSprites
        Terrain    TerrainSprites
    }

    type ProjectileSprites struct {
        BowTowerProjectile *ebiten.Image
    }

No new Go module dependencies are required.
