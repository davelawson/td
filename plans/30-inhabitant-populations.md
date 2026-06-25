# Add inhabitant populations to the top bar

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root. Save this file at `plans/30-inhabitant-populations.md`.

## Purpose / Big Picture

The wizard's Domain currently exposes physical resources but no inhabitants. After this change, the top bar shows Apprentices, Soldiers, and Peasants as a separate compact group beside Wood, Stone, and Metal. Each inhabitant type has its own icon and an `available/total` value so later recruitment and assignment systems have a visible status surface.

The initial prototype values are `0/0` for every inhabitant type. This slice does not add recruitment, assignment, population growth, consumption, or gameplay effects.

## Progress

- [x] (2026-06-25T00:00:00Z) Inspected the current HUD, game-status state, asset catalog, screenshot harness, tests, and relevant control documents.
- [x] (2026-06-25T00:00:00Z) Confirmed icon-plus-`available/total` presentation and zero initial populations with the user.
- [x] (2026-06-25T00:00:00Z) Created this ExecPlan.
- [x] (2026-06-25T23:38:59Z) Generated, processed, reviewed, and installed three inhabitant icons.
- [x] (2026-06-25T23:38:59Z) Added population state, asset-catalog wiring, grouped HUD rendering, and focused tests.
- [x] (2026-06-25T23:38:59Z) Updated current-product, game-design, roadmap, onboarding, and architecture documents.
- [x] (2026-06-25T23:38:59Z) Captured and reviewed rendered screenshot evidence.
- [x] (2026-06-25T23:38:59Z) Ran formatting, tests, whitespace, ownership, and repository-status checks.
- [x] (2026-06-25T23:38:59Z) Checked hand-written code-file line counts and recorded the pre-existing file over the 600-line preference.

## Surprises & Discoveries

- Observation: The host has `python3` but no `python` alias, matching the repository's `ART.md` instruction.
  Evidence: The first chroma-removal command returned `command not found: python`; rerunning the installed helper with `python3` succeeded.

- Observation: ImageMagick is not installed on the host.
  Evidence: `magick` and `identify` returned `command not found`, so the already-installed Pillow library performed nearest-neighbor trimming, scaling, centering, and PNG writing without adding a project dependency.

- Observation: The expanded status remains clear at 1920x1080 in calm, active-Raid, and overlay states.
  Evidence: `running-game.png`, `active-raid.png`, and `ingame-menu.png` show distinct resource and population groups without overlap with the centered phase text.

## Decision Log

- Decision: Show Apprentices, Soldiers, and Peasants as 64x64 source icons scaled to the existing 28-pixel HUD icon size, followed by compact `available/total` text.
  Rationale: This matches the established resource presentation while keeping the expanded top bar compact.
  Date/Author: 2026-06-25 / User and Codex

- Decision: Start all three population types at `0/0`.
  Rationale: Zero values avoid implying that recruitment, assignment, or inhabitant simulation already exists.
  Date/Author: 2026-06-25 / User

- Decision: Keep inhabitants as private prototype game-status data rather than introducing a package or service.
  Rationale: This slice only needs state and display; broader abstractions would have no implemented behavior to own.
  Date/Author: 2026-06-25 / Codex

## Outcomes & Retrospective

Implemented the requested inhabitants display. The game status now owns fixed prototype Apprentice, Soldier, and Peasant population values, each initialized to `0/0`. The top bar renders them as a separate icon-backed group after Wood, Stone, and Metal, with padded separators before the population group and before Barricade. Shared icon/value measurement keeps the entire right-side status right-aligned.

The typed asset catalog embeds and loads three new transparent 64x64 PNG icons. Focused tests cover icon dimensions, population item order, zero initial values, available-before-total formatting, sprite wiring, and full grouped width measurement. Documentation now defines available and total while keeping recruitment, assignment, growth, losses, and gameplay effects out of scope.

Validation completed:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

All Go tests passed. Screenshot capture produced the standard nine PNG files under `plans/30-inhabitant-populations/screenshots/`. Whitespace and ownership checks produced no output.

The final hand-written Go line-count review found one file over the 600-line preference: `internal/game/game_test.go` at 602 lines. This overage was pre-existing, and this change kept new HUD tests in `internal/game/hud_test.go`. The recommended follow-up remains a responsibility-based split of `game_test.go`, such as moving remaining state, camera, or map tests into focused files. No extra refactor was performed because it is outside this accepted plan.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `internal/game/hud.go` owns the top bar and the private `gameStatus` data used by that display. The current right-hand group draws Wood, Stone, and Metal icons with counts, followed by textual Barricade status. `assets/catalog.go` embeds 64x64 PNG assets and exposes HUD icons through `assets.Catalog.Sprite.Icon`. `internal/game/hud_test.go` verifies the resource item order, values, sprites, and colors. `cmd/td/main_test.go` contains an opt-in screenshot harness whose output directory must be advanced to this plan.

The accepted display order is Wood, Stone, Metal, a visible group separator, Apprentice, Soldier, Peasant, another separator, and Barricade. Each population number is formatted as `available/total`. “Available” means inhabitants not committed to a future assignment; “total” means every inhabitant of that type. The invariant is `0 <= available <= total`, although this slice only creates fixed zero values.

`ART.md` requires generated icons to be 64x64 PNG pixel art with sharp contrast, low detail, bright colors, no lighting, and no shadows. `DESIGN.md` requires compact, high-contrast HUD text that does not obscure the field. `PRODUCT.md` and `README.md` must describe the newly visible current behavior without implying population mechanics. `GAME.md` must record the three inhabitant types, value meanings, and future-system boundary. `ROADMAP.md` must recognize inhabitants as intended future gameplay work. `ARCHITECTURE.md` must continue to assign prototype status and HUD ownership to `internal/game` and asset loading to `assets`.

## Plan of Work

First, use the repository's image-generation workflow to create one apprentice icon, one soldier icon, and one peasant icon. Generate each opaque subject on a flat removable chroma-key background, remove that background locally, resize or crop to an exact 64x64 RGBA PNG if needed, and install the final files under `assets/sprites/icons/`. Review each icon at both source size and approximately 28x28 display size. The apprentice should read as a robed novice with a wand or scroll, the soldier as an armored guard with helmet or shield, and the peasant as a worker with a farming tool or grain. Do not include text, scenery, shadows, or floor planes.

Second, extend `assets/catalog.go` and `assets/catalog_test.go`. Embed and load the three files, expose them as `Apprentice`, `Soldier`, and `Peasant` fields on `IconSprites`, and test that all required icons load as 64x64 images.

Third, update `internal/game/hud.go`. Add private population count values to `gameStatus`, initialize all three to zero, and expose a focused population HUD item model in the stable order Apprentice, Soldier, Peasant. Refactor the existing right-side status measurement and rendering so resources and populations remain separate groups with explicit separators. Population items should reuse the current icon sizing and gaps but format their value as `available/total`. Keep top-bar height, Chapter/Day text, centered phase text, resource behavior, and Barricade behavior unchanged.

Fourth, update `internal/game/hud_test.go` with focused tests for population order, zero values, sprite wiring, formatting, and grouped width measurement. Keep tests in this focused file rather than adding to the already over-limit `internal/game/game_test.go`.

Fifth, update `README.md`, `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, and `ARCHITECTURE.md`. Describe the display as current behavior, define the intended population categories and value meanings in `GAME.md`, and state clearly that recruitment, assignment, and population changes are not implemented.

Sixth, change the screenshot output base path in `cmd/td/main_test.go` to `plans/30-inhabitant-populations/screenshots/`, capture the standard screenshots, and review `running-game.png` as the primary evidence. If the expanded right-side status overlaps the centered phase text at 1920x1080, tighten group spacing while retaining readable separators and 28-pixel icons.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Generate three icon sources with the built-in image-generation tool. Copy the selected outputs into a temporary workspace folder, remove their chroma-key backgrounds with:

    python "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" \
      --input <source> \
      --out <final.png> \
      --auto-key border \
      --soft-matte \
      --transparent-threshold 12 \
      --opaque-threshold 220 \
      --despill

After implementation, format and validate:

    gofmt -w cmd/td/main_test.go internal/game/hud.go internal/game/hud_test.go assets/catalog.go assets/catalog_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short

End with the required hand-written code line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

Record every file over 600 lines with its line count, likely cause, and a concrete recommendation. Do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` must pass. Asset tests must prove the apprentice, soldier, and peasant icons load at 64x64. HUD tests must prove the order is Apprentice, Soldier, Peasant; all initial values are `0/0`; available is formatted before total; and every item uses the corresponding non-nil icon.

The screenshot `plans/30-inhabitant-populations/screenshots/running-game.png` must show three resource icon/count pairs followed by a visually distinct group of three inhabitant icon/`0/0` pairs and then Barricade status. The group must remain readable at 1920x1080 without overlapping the centered phase text. Other standard screenshots must continue to render correctly.

Documentation must distinguish current display-only state from future inhabitant mechanics. Whitespace and ownership checks must produce no output, and the line-count result must be recorded in this plan.

## Idempotence and Recovery

Asset generation may be repeated into temporary or versioned files without overwriting accepted project assets until review. Chroma-key removal and screenshot capture are safe to repeat; screenshot capture overwrites only files under this plan's screenshot directory. Catalog, state, rendering, and documentation edits are ordinary source changes recoverable through version control.

If an icon has visible chroma spill, rerun removal once with `--edge-contract 1`. If a HUD overlap appears, reduce inter-item or inter-group horizontal gaps before changing icon size or font size. If catalog loading fails, verify the embedded path and exact PNG dimensions.

## Artifacts and Notes

The primary visual evidence is `plans/30-inhabitant-populations/screenshots/running-game.png`. The prior `plans/29-catapult-tower/screenshots/running-game.png` provides the pre-change baseline.

The built-in image-generation tool wrote source images under `/home/dave/.codex/generated_images/019f00f8-b23f-7b20-992d-006d864ac3a9/`. The three prompts requested centered, full-figure, low-detail pixel-art characters on a flat green chroma-key background: a purple-blue robed apprentice with wand and scroll, a steel-helmeted red-tunic soldier with round shield, and a brown-clothed peasant with a golden wheat sheaf. Each prompt prohibited scenery, shadows, text, watermarks, and green clothing. The installed chroma-key helper removed the background, and Pillow created centered 64x64 RGBA project assets using nearest-neighbor scaling.

## Interfaces and Dependencies

No new Go module dependency or public API is required. `assets.IconSprites` gains three exported image fields because it is the typed runtime asset catalog:

    type IconSprites struct {
        Wood       *ebiten.Image
        Stone      *ebiten.Image
        Metal      *ebiten.Image
        Apprentice *ebiten.Image
        Soldier    *ebiten.Image
        Peasant    *ebiten.Image
    }

Population types and HUD formatting remain private to `internal/game`. The implementation should use small structs containing semantic name, available count, total count, sprite, and display color. Rendering helpers must measure the exact text they draw so right alignment remains stable.

Revision note, 2026-06-25: Updated the living plan after implementation with completed progress, environmental discoveries, asset-generation details, validation results, screenshot evidence, and the required line-count review.
