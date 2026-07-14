# Add Gold and the Market economy

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept current while work proceeds. Maintain this document in accordance with `PLANS.md` from the repository root.

## Purpose / Big Picture

After this change, defeating raiders funds a new Gold economy instead of directly producing Wood, Stone, and Iron. The player can see Gold beside the other resources, build a staffed Market, select it during Management, and spend Gold through controls beside the building to buy one unit of a chosen construction resource. The implementation also resolves the existing Metal-versus-Iron terminology mismatch by consistently naming the third material Iron.

The result is visible by starting a game, building and selecting a Market, observing the three purchase controls beside it, and defeating raiders to increase the Gold total without changing material totals. Automated tests cover rewards, construction, trades, input blocking, phase rules, and layout; plan-local screenshots and an asset contact sheet provide visual evidence.

## Progress

- [x] (2026-07-14 22:27Z) Inspected the resource ledger, enemy templates, structure catalog, Economic building bar, selection flow, UI ownership, asset catalog, control documents, baseline tests, ownership, and current line counts.
- [x] (2026-07-14 22:27Z) Confirmed user decisions: rename Metal to Iron; use Gold drops 1, 2, 3, and 5 in Raid-tier order; use world-anchored one-unit purchase controls; allow trades only during Management; and generate both Gold and Market sprites.
- [x] (2026-07-14 22:59Z) Generated separate Gold and Market sources with the built-in image workflow, removed their chroma keys, normalized them to exact 64x64 transparent PNGs, and inspected both native and reduced views.
- [x] (2026-07-14 22:59Z) Implemented the Iron rename, Gold ledger and tiered combat rewards, four-resource HUD, Market catalog/placement/staffing, selection facts, and contextual atomic trades.
- [x] (2026-07-14 22:59Z) Added pure gameplay, layout, input, asset, and presentation tests in responsibility-specific files without expanding the near-limit placement test.
- [x] (2026-07-14 22:59Z) Synchronized current control documents and captured the asset contact sheet plus a selected-Market gameplay screenshot.
- [x] (2026-07-14 22:59Z) Passed formatting, normal, race, 50-run focused, screenshot, whitespace, terminology, status, and ownership validation.
- [x] (2026-07-14 22:59Z) Reviewed hand-written Go line counts. No file exceeds 600 lines; the pre-existing `internal/game/build_placement_test.go` and `internal/game/raid_test.go` remain the only near-limit files at 564 and 541 lines.

## Surprises & Discoveries

- Observation: The repository currently calls the extracted mineral and its source structures Iron, but calls the resulting spendable resource Metal.
  Evidence: `ART.md` explicitly says the Iron Mine produces Metal, while `internal/game/resources.go` and `internal/game/hud.go` use `metal` and `Metal` fields.

- Observation: Enemy rewards are already template-owned and applied only by tower-damage defeat, so the new economy does not require physical loot entities or changes to Barricade behavior.
  Evidence: `internal/game/combat.go` calls `grantEnemyResources` only after damage reduces health to zero; Sanctum contact removes enemies through the Raid path without that call.

- Observation: `internal/game/build_placement_test.go` is already 564 lines.
  Evidence: The baseline line-count review reports 564 lines, so Market placement scenarios will live in a new responsibility-specific test file rather than pushing it through the 600-line preference.

- Observation: The generated sources used slightly off-request border colors even though the prompts specified flat `#ff00ff`.
  Evidence: Automatic border sampling found `#fa05f4` for Gold and `#f904f8` for Market; chroma removal followed by local validation left no visible magenta pixels.

- Observation: The first isolated render after a fresh Ebitengine test process can transiently omit most screen-space UI.
  Evidence: The first capture showed only partial plot labels; the single recovery rerun documented by this plan produced the complete selected-Market frame and did not affect gameplay tests.

- Observation: Combining all Market prices into one selection-panel row clipped at 1080p.
  Evidence: Splitting the facts into `Buys Wood`, `Buys Stone`, and `Buys Iron` rows made every rate readable in the final full-resolution capture.

## Decision Log

- Decision: Rename the current Metal resource and its current asset/code-facing identifiers to Iron, while leaving historical plans unchanged.
  Rationale: Terrain and production already use Iron language, and the Market requirement names Iron. One current term avoids presenting two names for the same resource.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Skeletons drop 1 Gold, Zombies 2, Ghouls 3, and Armoured Skeletons 5.
  Rationale: This follows the current Raid threshold order and gives the strongest tier the requested five-Gold maximum.
  Date/Author: 2026-07-14 / User and Codex

- Decision: One Market trade click buys exactly one unit; Wood and Stone cost 1 Gold and Iron costs 3 Gold.
  Rationale: This is the smallest observable transaction model and matches the accepted interaction choice without introducing quantities, caps, or pricing abstractions.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Market controls are fixed-size screen UI anchored beside the selected projected Market Tile, available only during Management including SPACE-paused Management.
  Rationale: This implements “next to the market,” remains readable across zoom levels, and follows the existing phase boundary for settlement instructions.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Generate a Market structure sprite in addition to the requested Gold icon.
  Rationale: Every current building uses the same sprite for the map and building bar; a generated Market preserves that established presentation instead of introducing a placeholder.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

The complete slice is implemented. Raider templates now award only deterministic Gold on tower-damage defeat: Skeleton 1, Zombie 2, Ghoul 3, and Armoured Skeleton 5. The renamed Iron resource and new Gold total flow through state, construction metadata, selection facts, UI colors, the typed asset catalog, and the four-item top bar. Barricade removal and breach still award nothing.

Market is the fourth Economic entry. It costs 30 Wood, reserves one Soldier and two Peasants atomically, produces nothing during Labour, and exposes three fixed-size camera-relative controls when selected during Management. Each click buys one material at the accepted price, insufficient-Gold controls retain their hit ownership without mutating state, pause permits trading, and Labour, Raid, breach, off-screen selection, and the overlay hide trading. Input blockers prevent Market UI clicks from reaching selection, exploration, dragging, placement, or Raid controls.

Both generated assets pass exact RGBA, 64x64, transparent-corner, nonempty-alpha, and zero-visible-magenta checks. Gold has alpha bounds `(5, 7, 59, 57)` and Market `(3, 4, 61, 60)`. Full-resolution inspection confirms the contact sheet reads at 64x64 and 32x32 and the 1920x1080 gameplay capture shows Gold, the Market, its three adjacent buttons, and unclipped selection facts.

Validation passed with `go test ./...`, `go test -race ./...`, and 50 consecutive focused runs over assets, game rules, and UI. The focused screenshot test passed, `git diff --check` is clean, the non-historical Metal search returns no matches, and the ownership scan returns no non-`dave` files. The largest hand-written Go files remain pre-existing tests at 564 and 541 lines; no file exceeds the 600-line preference. Future additions to placement or Raid tests should split by responsibility instead of extending those files.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `internal/game/resources.go` owns the private resource ledger and resource mutations. `internal/game/enemies.go` owns exported enemy templates, while `internal/game/combat.go` applies rewards after tower-damage defeats. `internal/game/hud.go` adapts resource counts and loaded icons into the top bar.

Structures are templates in `internal/game/structures.go`, map features in `internal/game/map.go`, and stable building actions mapped in `internal/game/building_bar_items.go`. `internal/game/building_bar.go` owns construction eligibility and placement, while `internal/ui/building_bar*.go` owns Economic-tab ordering, layout, hit testing, tooltips, and drawing. Selection is stored as a Tile coordinate in `internal/game/selection.go`; `internal/game/selection_panel.go` adapts selected subjects to the presentation-neutral types drawn by `internal/ui/selection_panel.go`.

`assets/catalog.go` explicitly embeds and loads every runtime PNG into typed asset structs. Existing icons and structures are exact 64x64 PNGs. New project-bound raster assets must use the built-in image generator, a removable flat chroma-key background, the installed removal helper, nearest-neighbour normalization, and plan-local visual evidence. No dependency, save system, reassignment system, resource capacity, loot entity, or asset pipeline is needed.

## Plan of Work

First generate two isolated pixel-art assets with separate built-in image calls. Use the current Wood, Stone, and Iron icons as style-only references for a Gold HUD icon containing exactly two stacks of golden coins. Use the Woodcutter, Stone Quarry, and Iron Mine as style-only references for a compact medieval open market stall. Both prompts require one centered subject, low-detail hard-cluster pixel art, generous padding, a perfectly flat `#ff00ff` background, no shadow, scenery, text, border, watermark, or key color in the subject. Remove the key with the installed helper, crop to nontransparent content, resize with nearest-neighbour sampling, and center each subject on an exact 64x64 transparent canvas. Save final assets at `assets/sprites/icons/gold.png` and `assets/sprites/structures/market.png`. Create `plans/60-gold-and-market/screenshots/market-assets-contact-sheet.png` with native and 32x32 views.

Then rename the active Metal terminology to Iron. Move `assets/sprites/icons/metal.png` to `assets/sprites/icons/iron.png`; rename private counts, exported `Resources` and UI `ResourceAmounts` fields, catalog fields, colors, labels, costs, yields, tests, and current documentation. Add Gold to the private and exported resource representations, initialize it to zero, and render resources in Wood, Stone, Iron, Gold order. Current construction costs and Iron Mine yield keep their numeric values.

Replace `EnemyTemplate.Resources` with `GoldDrop int`. Set exact drops to 1, 2, 3, and 5 for Skeleton, Zombie, Ghoul, and Armoured Skeleton. A tower-damage defeat must add only Gold; Barricade removal and breach clearing remain unrewarded. Add the Gold-drop amount to raider selection details so tier value is inspectable.

Add a `Market` structure template with description, generated sprite, `Resources{Wood: 30}` construction cost, and `StaffingRequirements{Peasants: 2, Soldiers: 1}`. Add a map feature, stable `BuildingBarMarket` action after Iron Mine in the Economic tab, placement/rendering/selection mappings, and catalog loading. Market placement follows existing empty-grass, resource, population, phase, pause, and overlay rules. It reserves staff but does not produce during Labour or consume terrain. Give it a dedicated selection-panel kind so it displays Structure, Cost, Required Soldier, Required Peasants, and the three trade rates without a false “Produces Nothing” row.

Create a dedicated Market UI component under `internal/ui/` with stable actions Buy Wood, Buy Stone, and Buy Iron. Its presentation model contains label, resource icon, Gold price, enabled state, and hover state; it owns fixed-pixel layout, hit testing, and drawing, while `internal/game/` owns state eligibility and mutations. Lay the vertical control group immediately to the right of the selected Market Tile, fall back to its left if the group would overflow, and clamp below the top bar and inside the drawable width and height. Hide the group if the Market Tile center is outside the scene viewport. Draw labels `+1 Wood · 1 Gold`, `+1 Stone · 1 Gold`, and `+1 Iron · 3 Gold`, with disabled styling when current Gold is below the price.

Update Market input before general selection input. An enabled click subtracts the price and adds exactly one target resource in one state transition. Disabled clicks change nothing but are still consumed. Treat visible Market controls as screen UI for selection, exploration, camera-drag initiation, building drops, and any other map click blocker. The overlay menu already short-circuits normal game input. Controls are available in Management even while SPACE-paused, and hidden during Labour, Raid, breach, or when a non-Market subject is selected.

Finally update `README.md`, `PRODUCT.md`, `GAME.md`, `DESIGN.md`, `ART.md`, and `ARCHITECTURE.md`. Mark the old Wood/Stone/Metal decision in `GAME.md` as superseded and record the new Gold/Market economy. Describe current behavior rather than rewriting historical plan files. `ROADMAP.md` and `CODESTYLE.md` remain unchanged because product direction and coding conventions do not change.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`. After each built-in generation, copy its output from the Codex generated-images location to a plan-local working directory, then process it with:

    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input <source.png> --out <alpha.png> --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill

Use a small plan-local Pillow command or script to crop, nearest-neighbour resize, center, create the contact sheet, and validate RGBA mode, exact 64x64 size, transparent corners, nonempty alpha bounds, no visible `#ff00ff`, and readable 32x32 previews. Keep only accepted runtime assets and focused evidence; remove intermediate generated sources when finished.

Implement source and test changes incrementally and format all changed Go files. Run:

    gofmt -w <changed Go files>
    go test ./...
    go test -race ./...
    go test ./assets ./internal/game ./internal/ui -run 'Gold|Iron|Market|Enemy|Combat' -count=50
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureSelectedMarketScreenshot -count=1
    rg -n 'Metal|metal' --glob '!plans/**'
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The non-historical Metal search must return no matches. End with the required line-count review:

    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 25

Record any file near or above 600 lines and recommend a responsibility-based response without performing extra unapproved refactoring.

## Validation and Acceptance

Automated acceptance requires exact template drops and proof that combat kills change only Gold. Tests must cover all four enemy types, multiple kills, and no Gold from Sanctum contact. HUD tests must prove the four-resource order, initial counts `100/50/20/0`, exact names, non-nil icons, and width calculations.

Market tests must prove Economic-tab ordering, its exact cost and staffing, failure when each prerequisite is missing, atomic placement and staff reservation, one unit per click, the three exact prices, exact-balance success, insufficient-Gold no-op, and independent staffing for multiple Markets. Interaction tests must prove paused Management access, hiding in every other phase and breach, overlay blocking, selection preservation, no map pass-through, camera-relative anchoring, right-to-left fallback, viewport clamping, and off-screen hiding.

Visual acceptance requires the contact sheet to show a readable two-stack coin icon and distinct medieval Market at both sizes. The focused 1920x1080 screenshot must show a selected Market, its three controls beside the Tile, its selection details, and Gold in the top bar. Inspect both PNGs at full resolution.

Documentation acceptance requires all current control documents to agree on Iron, Gold drops, starting totals, Market cost/staffing/trades, phase availability, contextual UI, and asset ownership. Historical plans remain untouched.

## Idempotence and Recovery

Formatting, tests, screenshot capture, contact-sheet creation, and validation are safe to rerun. Screenshot generation overwrites only this plan's evidence. Preserve generated source images until normalized assets pass validation. If chroma removal leaves a fringe, retry once with `--edge-contract 1`; do not switch to the CLI transparency fallback without user confirmation. If a transaction test fails, inspect the game-owned action mapping before changing UI layout. If a visual capture transiently omits screen-space UI, rerun the isolated capture once, matching prior repository evidence.

## Artifacts and Notes

Final assets belong at:

    assets/sprites/icons/gold.png
    assets/sprites/structures/market.png

Visual evidence belongs at:

    plans/60-gold-and-market/screenshots/market-assets-contact-sheet.png
    plans/60-gold-and-market/screenshots/selected-market.png

The final built-in Gold prompt requested exactly two stacks of bright golden coins, one tall and one short, as one centered low-detail hard-cluster pixel-art HUD subject; it used the existing Wood, Stone, and Iron icons as style-only references and prohibited loose coins, extra objects, text, border, shadow, scenery, watermark, and key-color pixels in the subject against flat `#ff00ff`.

The final built-in Market prompt requested one centered compact medieval timber market stall with a readable blue cloth canopy and contained trade goods in the economic-building sprite style; it used Woodcutter, Stone Quarry, and Iron Mine as style-only references and prohibited people, text, signage, border, shadow, scenery, watermark, and key-color pixels in the subject against flat `#ff00ff`.

Gold's final alpha bounds are `(5, 7, 59, 57)` and Market's are `(3, 4, 61, 60)`. Both are exact RGBA 64x64 images with transparent corners and no visible magenta pixels. `market-assets-contact-sheet.png` shows the accepted native and 32x32 previews; `selected-market.png` is the accepted full gameplay view. Intermediate generated sources were removed after acceptance.

## Interfaces and Dependencies

The existing exported project types change as follows:

    type Resources struct {
        Wood  int
        Stone int
        Iron  int
        Gold  int
    }

    type EnemyTemplate struct {
        // existing combat fields
        GoldDrop int
        // existing sprite fields
    }

`assets.IconSprites` replaces `Metal` with `Iron` and adds `Gold`; `assets.StructureSprites` and `game.StructureCatalog` add `Market`. The UI package gains stable Market trade actions and presentation-neutral models for controls and selection facts. Market kinds and map features remain private except for the UI action identifiers already used across the `internal/game` and `internal/ui` boundary.

No external runtime dependency, serialized format, migration, networking, save behavior, reassignment behavior, resource capacity, loot entity, animation, or asset pipeline is added. Pillow remains a development-time image normalization tool and is not added to `go.mod`.

Revision note (2026-07-14): Created from the accepted implementation plan after fresh repository, workflow, baseline-test, ownership, and line-count inspection. Updated after implementation with completed progress, generated-art details, visual-review discoveries, validation evidence, and final line-count findings.
