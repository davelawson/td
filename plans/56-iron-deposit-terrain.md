# Add Rare Iron Deposit Terrain

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/56-iron-deposit-terrain.md`.

## Purpose / Big Picture

After this change, the player can encounter rare Iron Deposit terrain while starting or expanding the Domain. Iron Deposits appear on 1% of generated grasslands Tiles and 3% of generated hills Tiles before roads and shared Plot edges are cleared. They use distinct terrain sprites, block construction, and can be selected to show `Iron Deposit` and the containing Plot biome. This slice does not make deposits gatherable and does not change the current Iron Mine, which can still be built on empty grass and produces Metal during Labour.

The result is observable by running the game, exploring grasslands and hills, and selecting a deposit, or by viewing the deterministic selected-deposit screenshot under `plans/56-iron-deposit-terrain/screenshots/`.

## Progress

- [x] (2026-07-14T17:26Z) Inspected terrain generation, rendering, selection, placement, assets, tests, control documents, and screenshot workflows; accepted the terrain-only behavior and exact rarity rates with the user.
- [x] (2026-07-14T17:26Z) Created this ExecPlan before implementation.
- [x] (2026-07-14T17:31Z) Generated, post-processed, visually reviewed, and installed four Iron Deposit terrain sprites.
- [x] (2026-07-14T17:34Z) Added generation, catalog, rendering, selection, and atomic placement rejection with focused tests.
- [x] (2026-07-14T17:38Z) Updated durable documentation and captured the standard suite plus focused selected-deposit evidence.
- [x] (2026-07-14T17:40Z) Ran formatting, full tests, race tests, screenshot capture, whitespace, alpha, and ownership validation.
- [x] (2026-07-14T17:40Z) Checked hand-written Go file line counts; no file exceeds 600 lines, and the largest is `internal/game/build_placement_test.go` at 564 lines.
- [x] (2026-07-14T17:51Z) Responded to visual review by replacing all four dark orange-veined deposits with gray boulder-like variants carrying sparse blue ore veins, then refreshed the art guidance and validation evidence.

## Surprises & Discoveries

- Observation: The current terrain generator already uses explicit percentage weights, independently from each Tile's persistent visual tweak.
  Evidence: `internal/game/plot_generator.go` uses `terrainWeights` and `weightedTerrain`, while `Tile.Tweak` only selects stable sprite variants and mirroring.
- Observation: Iron Mine already means a structure that produces the existing Metal resource anywhere it is placed on empty terrain.
  Evidence: `GAME.md`, `internal/game/structures.go`, and the Labour resource tests consistently use the name Iron Mine with `Resources{Metal: 10}`.
- Observation: Focused terrain screenshots can construct private game state directly without weakening runtime randomness or exposing a test API.
  Evidence: `internal/game/visual_test.go` already builds a selected Tree Tile and captures it through Ebitengine.
- Observation: The built-in image generator produced clean flat-key pixel-art sources for all four silhouette prompts without requiring a true-transparency fallback.
  Evidence: local chroma removal reported over 945,000 transparent source pixels for every variant; final validation found zero opaque bright-green pixels and fully transparent corners in all four 64x64 assets.
- Observation: Random standard captures naturally showed the intended biome contrast despite the low rates.
  Evidence: `plans/56-iron-deposit-terrain/screenshots/explored-biomes.png` shows several deposit variants among Boulder-heavy hills while keeping roads, shared joins, and empty grass readable.
- Observation: The first dark-charcoal, orange-veined treatment read as a separate kind of outcrop instead of a mineral-bearing member of the existing Boulder family.
  Evidence: User visual review requested Boulder-like forms with bluish veins; the replacement sprites now derive their silhouettes, gray-rock palette, and 54-pixel footprint directly from the four Boulder references.

## Decision Log

- Decision: Implement Iron Deposit as terrain only.
  Rationale: A mine-required or depleting node would require new extraction, placement, and economy rules beyond this small observable terrain slice.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Generate Iron Deposit at 1% in grasslands and 3% in hills.
  Rationale: Deposits remain rare everywhere and are three times more likely in hills, matching the request that they be particularly rare in grasslands.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Preserve existing Tree and Boulder rates and take the new percentages from empty terrain.
  Rationale: The new terrain adds map pressure without silently changing the established character of either biome.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Use four transparent sprite variants with tweak-derived selection and mirroring.
  Rationale: This follows the existing pine-tree and Boulder terrain convention and avoids a one-sprite visual pattern.
  Date/Author: 2026-07-14 / Codex
- Decision: Use dark rock with rusty orange-red iron seams for the deposit art.
  Rationale: High-contrast ore seams distinguished deposits from ordinary gray Boulders at 50% scale without adding labels to the map. This decision was superseded after visual review.
  Date/Author: 2026-07-14 / Codex
- Decision: Make Iron Deposits look like ordinary gray Boulders containing sparse cobalt-blue and cyan-blue ore veins.
  Rationale: The shared silhouette language makes the deposits read as natural boulders, while the cool high-contrast veins preserve recognition at 50% scale.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Implementation is complete. Grasslands now generates 6% Tree, 3% Boulder, 1% Iron Deposit, and 90% empty grass, including the home Plot before its protected road is applied. Hills generates 3% Tree, 6% Boulder, 3% Iron Deposit, and 88% empty grass. Road and shared-edge cleanup overwrite deposits through the existing map rules. Four generated transparent sprite variants load through the typed asset catalog and render with stable tweak-derived variant selection and mirroring.

Iron Deposits participate in the existing natural-obstacle workflow: they block construction without spending resources or changing population, can be selected after raiders and structures, receive the gold Tile outline, and show `Iron Deposit` plus the Plot biome in the information panel. Iron Mines retain their previous empty-terrain placement and 10 Metal Labour production; extraction and depletion remain unimplemented.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, both screenshot commands, `git diff --check`, asset alpha inspection, and the ownership check all succeeded. Focused evidence is at `plans/56-iron-deposit-terrain/screenshots/selected-iron-deposit.png`, and the standard suite contains 18 screenshots. The final hand-written Go review found no file above 600 lines; `internal/game/build_placement_test.go` is the largest at 564 lines, so no refactor approval was needed.

Post-implementation visual review replaced the original dark, orange-veined outcrops with four boulder-like gray-rock variants bearing blue ore veins. The refreshed sprites retain exact 64x64 RGBA canvases, transparent isolation, distinct silhouettes, and half-scale readability while matching the established Boulder footprint and visual family more closely. The revision re-passed `go test ./...`, both screenshot captures, `git diff --check`, ownership validation, and focused image checks for dimensions, RGBA mode, transparent corners, green-key fringe, and visible blue pixels.

## Context and Orientation

The Go/Ebitengine game represents the world as 15-by-15 Plots made of Tiles. `internal/game/map.go` defines Tile terrain values and applies authored road and explored-edge cleanup. `internal/game/plot_generator.go` generates home and explored Plots from biome-specific percentage weights. Grasslands currently generates 6% Tree, 3% Boulder, and 91% empty grass; hills generates 3% Tree, 6% Boulder, and 91% empty grass. Home uses grasslands weights before its north road and Sanctum are applied.

`assets/catalog.go` embeds four pine-tree and four Boulder sprites into typed terrain arrays. `internal/game/scene.go` uses each Tile's stable `Tweak` value to choose one of four variants and optionally mirror it. `internal/game/selection.go` allows Tree and Boulder selection, and `internal/game/selection_panel.go` maps those values to player-facing terrain names. Building placement accepts only `terrainEmpty`, so a distinct Iron Deposit terrain value is non-buildable through the existing invariant.

The control documents constrain this feature. `GAME.md` owns intended terrain and resource design, `PRODUCT.md` and `README.md` describe current player-visible behavior, `ROADMAP.md` summarizes the current implemented slice while keeping richer resource nodes later, `ARCHITECTURE.md` records map and asset ownership, `DESIGN.md` defines terrain readability and selection treatment, and `ART.md` defines 64x64 transparent pixel-art terrain guidance. `CODESTYLE.md` requires Go formatting, focused tests, and a final review against the 600-line preference.

## Plan of Work

Use the built-in image-generation path to create four separate Iron Deposit variants. Each prompt must request a single isolated low-detail pixel-art ore outcrop, with a different silhouette and arrangement, dark charcoal stone, rusty orange-red seams, no text, no shadow, and a perfectly flat `#00ff00` background. Copy each generated source into a workspace staging directory, remove the chroma key with the installed imagegen helper, crop to the subject, resize with nearest-neighbor sampling to fit an exact 64x64 transparent canvas, and install the results as `assets/sprites/terrains/iron-deposit-1.png` through `iron-deposit-4.png`. Validate dimensions, alpha, transparent corners, silhouette, and readability at half scale.

Add `terrainIronDeposit` in `internal/game/map.go`. Add an Iron Deposit field to `terrainWeights`, set grasslands to Tree 6, Boulder 3, Iron Deposit 1 and hills to Tree 3, Boulder 6, Iron Deposit 3, then extend `weightedTerrain` after Boulder. The exact grasslands ranges are Tree 0-5, Boulder 6-8, Iron Deposit 9, and empty 10-99. The exact hills ranges are Tree 0-2, Boulder 3-8, Iron Deposit 9-11, and empty 12-99. Keep home road, north-chain road, and shared-edge cleanup ordering unchanged so they overwrite any generated deposit.

Embed and load all four sprites through a new `IronDeposits [4]*ebiten.Image` field in `assets.TerrainSprites`. Render Iron Deposit in `internal/game/scene.go` through the existing general terrain-sprite helper with stable tweak-derived variant selection and horizontal mirroring. Give the underlying Tile a subtle ore-compatible tint and size the subject near the Boulder scale so it reads as terrain rather than a structure.

Extend terrain selection and panel mapping to Iron Deposit. Keep raider, structure, then natural-terrain priority unchanged. Reuse the gold selected-Tile outline and show `Iron Deposit` with the Plot's stored biome. Do not add commands, yield, depletion, gathering, or Iron Mine placement restrictions.

Add deterministic tests for asset loading and dimensions, both biome percentage boundaries, generated terrain coverage, home generation, road and explored-edge overwrites, selection and panel data, and build rejection without resource or population changes. Add a focused `internal/game/visual_test.go` capture that places and selects an Iron Deposit and writes `plans/56-iron-deposit-terrain/screenshots/selected-iron-deposit.png`. Point the standard screenshot suite at this plan and retain an explored-biomes image for composition review.

Update `README.md`, `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, `DESIGN.md`, and `ART.md` so they agree on rates, visual behavior, selection, non-buildability, and the absence of mining behavior. Record the new gameplay decisions in `GAME.md`; keep resource nodes as a later opportunity.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

Generate the four source images with one built-in image generation call per variant. After copying sources into a workspace staging area, remove chroma key with:

    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input <source> --out <alpha-output> --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill

Resize, center, and save each accepted sprite as an exact 64x64 PNG. Then edit the asset catalog, game code, focused tests, screenshot paths, and control documents. Format changed Go files and run:

    gofmt -w assets/*.go internal/game/*.go cmd/td/*.go
    go test ./...
    go test -race ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureSelectedIronDepositScreenshot -count=1
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

The ownership command must print nothing. Report every hand-written Go file over 600 lines with a maintainable split recommendation, and request approval before doing an unplanned split, refactor, or library addition.

## Validation and Acceptance

Automated tests must prove the exact generation boundaries and all obstacle invariants. At runtime, grasslands must be 6% Tree, 3% Boulder, 1% Iron Deposit, and 90% empty before authored overrides. Hills must be 3% Tree, 6% Boulder, 3% Iron Deposit, and 88% empty. Roads and explored joins must remain clear. Clicking a deposit must show a gold Tile outline and a panel labeled `Iron Deposit` with `Grasslands` or `Hills`. Dropping any building on it must cancel without changing resources, inhabitants, or Tile features.

All four PNGs must load as transparent 64x64 terrain sprites. `selected-iron-deposit.png` must show a visually distinct ore outcrop, the selection outline, and readable terrain/biome details. `explored-biomes.png` must retain readable roads, Plot joins, controls, and mostly buildable land. The seven updated control documents must consistently say that Iron Deposits are terrain only and do not yet alter Iron Mine behavior.

## Idempotence and Recovery

Formatting, tests, chroma-key processing, resizing, and screenshot capture are safe to rerun. Keep generated sources in a staging directory until all four installed assets pass visual and alpha review. If a generated subject has green fringe, retry the helper once with `--edge-contract 1`; do not switch to a true-transparency CLI model without user approval. If a road or shared edge retains a deposit, fix generation-versus-override order rather than adding a rendering exception. Preserve unrelated user changes if the worktree changes during implementation.

## Artifacts and Notes

Final runtime assets belong at `assets/sprites/terrains/iron-deposit-1.png` through `iron-deposit-4.png`. Visual evidence belongs under `plans/56-iron-deposit-terrain/screenshots/`. Record concise final validation and line-count evidence in the living sections above.

## Interfaces and Dependencies

No external runtime dependency or exported gameplay API is added. The internal game package gains `terrainIronDeposit` and `terrainWeights.IronDeposit`. The exported asset catalog gains `assets.TerrainSprites.IronDeposits [4]*ebiten.Image`. Existing constructors, `Resources`, structure templates, and Iron Mine behavior retain their signatures and semantics.

Revision note (2026-07-14): Created the ExecPlan from the accepted implementation plan after repository inspection. It fixes the rarity boundaries, terrain-only scope, sprite workflow, control-document updates, screenshot evidence, and final validation requirements.

Revision note (2026-07-14): Updated every living section after implementation. Recorded the asset-processing and visual-review findings, completed all progress items, and added final test, alpha, ownership, screenshot, and line-count outcomes.

Revision note (2026-07-14): Replaced the first Iron Deposit art direction after user review. The final sprites now follow the Boulder silhouettes and gray-rock palette with sparse blue ore veins; `ART.md` and `DESIGN.md` were updated to preserve this direction.
