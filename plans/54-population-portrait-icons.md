# Replace population figures with portrait badges

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds. Maintain this document in accordance with `PLANS.md` at the repository root.

## Purpose / Big Picture

The game currently represents Apprentices, Soldiers, and Peasants with small full-body figures. After this change, every place that uses those population icons will instead show an expressive, front-facing portrait inside a round badge: a surprised apprentice in a purple wizard hat, a grimacing soldier in a metal cap, and a smiling peasant in a straw hat. A player can see the result both in the top status bar and in the population rows shown on building cards. The filenames and asset-catalog fields remain stable, so this is a visual replacement rather than a new runtime system.

## Progress

- [x] (2026-07-14 15:56Z) Inspected the existing 64x64 population sprites, their catalog wiring, the HUD and building-card consumers, current art guidance, and the active dirty worktree.
- [x] (2026-07-14 15:56Z) Created this plan and selected plan 53's running-game screenshot as the before-state evidence.
- [x] (2026-07-14 16:01Z) Generated three distinct portrait-badge candidates on flat chroma-key backgrounds, removed the backgrounds, and normalized each asset to a centered 64x64 RGBA PNG.
- [x] (2026-07-14 16:01Z) Reviewed the candidates at native size and at their 28-pixel HUD and 14-pixel building-card display sizes, then replaced the three stable asset files.
- [x] (2026-07-14 16:02Z) Updated `ART.md`, `DESIGN.md`, and the screenshot harness plan path to record and demonstrate the new visual language.
- [x] (2026-07-14 16:03Z) Captured and inspected after-state screenshots for the top HUD and House, Barracks, and Dorm population rows.
- [x] (2026-07-14 16:04Z) Ran focused asset tests, the full Go suite, whitespace validation, image validation, and ownership validation successfully.
- [x] (2026-07-14 16:04Z) Checked line counts for hand-written Go files. No file exceeds 600 lines; `internal/game/building_bar_test.go` is close at 595 lines, so a future change should move its test fixtures or responsibility groups into focused companion test files rather than extending it.

## Surprises & Discoveries

- Observation: The same three asset-catalog fields feed both the top HUD and building-card metadata, so replacing the existing files updates every intended presentation without duplicating assets or changing runtime APIs.
  Evidence: `internal/game/hud.go` and `internal/game/building_bar_metadata.go` both use `catalog.Sprite.Icon.Apprentice`, `Soldier`, and `Peasant`.

- Observation: Building-card metadata renders these icons at roughly 14 pixels while the top HUD uses 28 pixels, making silhouette and hat color more important than small facial detail in the card view.
  Evidence: The current drawing constants and visual screenshots show the two distinct display scales.

- Observation: Chroma removal followed by nearest-neighbor normalization retained the intended hard pixel edges without leaving any measurable bright-green subject pixels.
  Evidence: Final validation reported zero greenish nontransparent pixels for all three 64x64 RGBA files and transparent alpha values in all four corners.

- Observation: The screenshot harness can capture some building-hover targets while the rest of a new frame is still being drawn, but the targeted metadata rows and population badges are present and unobscured.
  Evidence: `house-icon.png` and `dorm-icon.png` contain partial peripheral UI frames, while their House and Dorm population rows remain visible; `running-game.png`, `barracks-icon.png`, and the contact sheet provide complete-scale evidence.

## Decision Log

- Decision: Replace the existing files in place at `assets/sprites/icons/apprentice.png`, `soldier.png`, and `peasant.png`.
  Rationale: The user requested the portraits everywhere, and stable paths automatically update every consumer while avoiding catalog or API churn.
  Date/Author: 2026-07-14 / Codex

- Decision: Use a shared round badge composition with a thin bronze rim, dark charcoal interior, transparent exterior, and a front-facing head-and-shoulders crop.
  Rationale: The common frame makes population status read as one visual family, while each hat, expression, and palette remains role-specific.
  Date/Author: 2026-07-14 / Codex

- Decision: Generate each role separately using its current icon as a pixel-art style reference, then use a flat `#00ff00` background and local chroma-key removal.
  Rationale: Separate prompts keep each requested expression precise. Chroma removal follows the repository's existing transparent-asset process and produces inspectable intermediate evidence.
  Date/Author: 2026-07-14 / Codex

- Decision: Keep gameplay, icon ordering, layout dimensions, catalog fields, and filenames unchanged.
  Rationale: The requested work is a visual representation change, not a population-system or interface redesign.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

The three full-body population figures were replaced in place by a cohesive portrait-badge set. The Apprentice now reads through a surprised face and pointed purple hat, the Soldier through a tooth-baring grimace and steel cap, and the Peasant through a smile and straw hat. Because the catalog paths stayed stable, the new art appears both in the top HUD and in House, Barracks, and Dorm population metadata with no gameplay, ordering, spacing, dependency, or exported-API change.

The generated sources used the built-in image-generation path and a flat green background. Local chroma removal, alpha-bounds cropping, nearest-neighbor resizing, and centered 64x64 canvases produced exact RGBA files with transparent corners and no detected green fringe. `ART.md` now records the reusable prompt and processing guidance, and `DESIGN.md` records the portrait family and two-scale legibility requirement.

`go test ./assets` and `go test ./...` passed. `git diff --check` and the ownership scan produced no output. The hand-written Go line-count review found no file over 600 lines. `internal/game/building_bar_test.go` is close to the preference at 595 lines because it groups many building-bar behaviors; a future change that needs to extend it should move population or metadata cases into a focused companion test file. That refactor was not needed or approved for this asset change, so it was not performed.

## Context and Orientation

`assets/sprites/icons/apprentice.png`, `assets/sprites/icons/soldier.png`, and `assets/sprites/icons/peasant.png` are embedded 64x64 RGBA PNGs. `assets/catalog.go` loads them into the `IconSprites` fields with the same role names. `internal/game/hud.go` draws those fields in the top population status group, while `internal/game/building_bar_metadata.go` uses the same fields for population costs and grants on building cards. `assets/catalog_test.go` already proves that all three files load and remain 64x64.

The working tree also contains uncommitted plan 53 UI work. Those edits belong to the user and must remain intact. The only hand-written Go edit needed here is changing the visual-evidence destination in `cmd/td/main_test.go` from plan 53 to plan 54; the application rendering code does not need a new branch or option.

`ART.md` requires generated game art to use 2D pixel art, sharp edges, high contrast, low detail, bright readable colors, exact 64x64 icon canvases, and transparent backgrounds. `DESIGN.md` prioritizes readability at the target window size. This change establishes a durable population-icon convention, so both files must record the badge composition and small-scale readability criteria. `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `CODESTYLE.md`, and `ARCHITECTURE.md` do not need edits because no workflow, gameplay decision, source convention, or module boundary changes.

## Plan of Work

First, preserve `plans/53-phase-aware-game-ui/screenshots/running-game.png` under `plans/54-population-portrait-icons/screenshots/before/running-game.png`. Generate one candidate per role with the built-in image-generation workflow. Each prompt asks for a single isolated, front-facing pixel-art bust inside a thin bronze circle on a perfectly flat green background, with no props, body, scenery, text, shadow, or watermark. The peasant must smile beneath a straw hat and simple brown collar. The soldier must grimace beneath a simple steel skullcap and armor collar, leaving the face visible. The apprentice must have wide surprised eyes, a small open mouth, and a pointed purple wizard hat.

Copy generated sources into `tmp/imagegen/54-population-portrait-icons/`. Run the installed chroma-removal helper with automatic border keying, soft matte, and despill. Crop each alpha-bounded subject and resize it with nearest-neighbor sampling to fit about 60x60 pixels centered on an exact 64x64 transparent canvas. Inspect the resulting badge at native size and in a contact sheet that includes 28x and 14x renderings. If a subject is not recognizable or has a visible green fringe, make one targeted generation or edge-cleanup retry before promoting it.

After review, overwrite only the three approved icon paths. Update `ART.md` with the concrete asset descriptions, shared prompt pattern, processing method, and small-scale review rules. Update `DESIGN.md` so expressive round portrait badges become the stable visual direction for population indicators. Point `cmd/td/main_test.go` at the plan 54 screenshot directory, capture the normal running view plus the House, Barracks, and Dorm building-card views, and visually confirm both display scales.

Finally, run all asset and repository tests, validate PNG dimensions and transparency, check whitespace and ownership, and review hand-written Go file line counts. Preserve all unrelated working-tree changes.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Create the evidence and working directories and preserve the prior running view before replacing assets:

    mkdir -p plans/54-population-portrait-icons/screenshots/before tmp/imagegen/54-population-portrait-icons
    cp plans/53-phase-aware-game-ui/screenshots/running-game.png plans/54-population-portrait-icons/screenshots/before/running-game.png

Use one built-in image-generation call for each role, with its current file as the style reference. Copy each generated source from the tool's generated-image directory into the plan-specific working directory. Remove the flat green background with:

    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input tmp/imagegen/54-population-portrait-icons/<role>-source.png --out tmp/imagegen/54-population-portrait-icons/<role>-alpha.png --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill

Normalize each alpha image to the final 64x64 form with Pillow using nearest-neighbor resampling. Copy the reviewed finals over the stable role paths and validate them with:

    file assets/sprites/icons/{apprentice,soldier,peasant}.png
    go test ./assets

Update `ART.md`, `DESIGN.md`, and `cmd/td/main_test.go`. Capture visual evidence in an environment with an X display:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Expect screenshots in `plans/54-population-portrait-icons/screenshots/`, including `running-game.png`, `house-icon.png`, `barracks-icon.png`, and `dorm-icon.png`. Inspect those four images and the portrait contact sheet before acceptance.

Run final repository validation:

    go test ./...
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short

The test suite should pass, the whitespace and ownership checks should print nothing, and the status should retain both this work and all pre-existing user changes.

End with the required hand-written Go file line-count review:

    rg --files -g '*.go' -g '!vendor/**' | xargs wc -l | sort -n

Report any file at or above 600 lines with its count and a responsibility-based split recommendation. Do not implement an unplanned refactor without user approval.

## Validation and Acceptance

The three final asset files must each be an exact 64x64 RGBA PNG with transparent corners, centered nonempty content, and no conspicuous green fringe. At 28 pixels, the peasant's smile and straw hat, the soldier's grimace and metal cap, and the apprentice's surprised expression and purple wizard hat should be readable. At 14 pixels, each role must remain distinguishable through its hat shape, palette, and silhouette.

`go test ./assets` must prove that the embedded catalog still loads all three stable files at 64x64. `go test ./...` must pass without a runtime or gameplay regression. In `running-game.png`, all three new portraits must appear in the top HUD in the existing order with unchanged counts and spacing. In `house-icon.png`, `barracks-icon.png`, and `dorm-icon.png`, the relevant portrait must appear in the existing population cost or grant row without clipping or layout drift.

`ART.md` must describe the actual prompts and chroma-key-to-alpha workflow used for the final assets. `DESIGN.md` must state the durable portrait-badge language and its two-scale legibility expectation. No new dependency, exported API, gameplay rule, or broader asset pipeline is accepted as part of this plan.

## Idempotence and Recovery

The generation and normalization steps write named intermediates and can be repeated without accumulating runtime files. The generated sources were kept under `tmp/imagegen/54-population-portrait-icons/` through acceptance and then removed; the durable outputs are the three runtime assets and plan screenshots. Before promoting candidates, the original files remain recoverable from version control; do not use destructive Git commands because the worktree contains unrelated user edits. Re-running screenshot capture overwrites only plan 54 evidence. If green fringes appear in a future regeneration, rerun chroma removal once with `--edge-contract 1`, inspect again, and regenerate only if cleanup harms the badge rim or face.

## Artifacts and Notes

The before-state screenshot is `plans/54-population-portrait-icons/screenshots/before/running-game.png`. After-state evidence is in `plans/54-population-portrait-icons/screenshots/running-game.png`, `house-icon.png`, `barracks-icon.png`, and `dorm-icon.png`. `plans/54-population-portrait-icons/screenshots/population-portrait-contact-sheet.png` shows every final at native, HUD, and building-card scales.

The common generation prompt must retain these constraints: 2D pixel-art game icon; isolated front-facing head-and-shoulders portrait; face dominant; thin bronze circular rim; dark charcoal medieval-fantasy badge interior; perfectly flat solid `#00ff00` exterior background for removal; crisp edges and generous padding; no green in the subject; no full body, props, scenery, text, shadow, reflection, or watermark. Add only the role-specific face, hat, and collar instructions described above.

## Interfaces and Dependencies

No new Go interfaces, exported types, external services, or module dependencies are introduced. The final runtime interface remains `assets.IconSprites` with the existing `Apprentice`, `Soldier`, and `Peasant` fields. The implementation uses the existing embedded PNG paths, the built-in image-generation tool, the installed chroma-key helper, and Pillow for deterministic local resizing and contact-sheet assembly.

Revision note (2026-07-14): Created the implementation plan after inspecting the current asset wiring and resolving the intended everywhere-use, round-badge composition, and two display scales. Updated it after implementation with the final artifacts, validation evidence, screenshot-harness observation, and line-count outcome.
