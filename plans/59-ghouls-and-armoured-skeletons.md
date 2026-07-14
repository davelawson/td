# Add Ghouls and Armoured Skeletons

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must remain current as work proceeds. Maintain this document in accordance with `PLANS.md` from the repository root.

## Purpose / Big Picture

Raids currently use only ordinary Skeletons and Zombies. After this change, higher challenge scores also release fast Ghouls and durable Armoured Skeletons, each with its own readable 64x64 sprite, combat stats, and defeat rewards. Every distinct raider type has a speed difference of at least three percent, so unlike raiders that spawn together separate after only momentary overlap. Raid 1 stays exactly two Skeletons and one Zombie because its challenge score reaches only 4; at challenge 8 a reviewer can observe all four enemy types in stable rule order.

## Progress

- [x] (2026-07-14 20:03Z) Inspected the accepted implementation plan, repository instructions, current asset and enemy catalogs, Raid generator and runtime, tests, screenshot harness, art guidance, and durable documentation.
- [x] (2026-07-14 20:03Z) Confirmed the additive design: Ghoul at score multiples of 6 and Armoured Skeleton at score multiples of 8, without changing challenge or duration formulas.
- [x] (2026-07-14 20:07Z) Generated, post-processed, normalized, and validated both transparent 64x64 enemy sprites, including a native/half-size contact sheet.
- [x] (2026-07-14 20:09Z) Extended the embedded asset catalog and game enemy catalog with exact names, stats, rewards, sprite keys, and sprite references.
- [x] (2026-07-14 20:10Z) Added ordered Raid kinds and rules, extended kind-to-template mapping, and added deterministic boundary, skipped-threshold, roster, movement, and challenge-16 tests.
- [x] (2026-07-14 20:13Z) Updated `README.md`, `PRODUCT.md`, `GAME.md`, `ART.md`, and `ARCHITECTURE.md` to match the implemented roster and asset guidance.
- [x] (2026-07-14 20:15Z) Captured and inspected focused runtime evidence showing both enemies separated on a cleared road with the Armoured Skeleton selected.
- [x] (2026-07-14 20:16Z) Ran formatting, normal, race, 50 repeated focused, image, screenshot, whitespace, status, and ownership validation successfully.
- [x] (2026-07-14 20:16Z) Checked hand-written Go file line counts. No file exceeds 600; `internal/game/build_placement_test.go` is unchanged at 564 and `internal/game/raid_test.go` is now 524, so both are reported without an unplanned split.
- [x] (2026-07-14 21:47Z) Revised the Armoured Skeleton to 0.9 Tiles per second and added pairwise speed-separation, movement-divergence, and inspection-panel assertions.
- [x] (2026-07-14 21:47Z) Updated `PRODUCT.md` and `GAME.md` with the distinct-speed behavior and reopened this plan's living sections.
- [x] (2026-07-14 21:49Z) Recaptured and inspected the focused runtime evidence with the selected Armoured Skeleton reporting 0.9 Tiles per second.
- [x] (2026-07-14 21:50Z) Reran formatting, focused repetitions, normal and race tests, whitespace, status, and ownership validation for the speed revision.
- [x] (2026-07-14 21:50Z) Repeated the hand-written Go line-count review. No file exceeds 600; `internal/game/build_placement_test.go` remains 564 lines and `internal/game/raid_test.go` is now 541 lines, so both are reported without an unplanned split.

## Surprises & Discoveries

- Observation: The current enemy source sprites already use exact 64x64 RGBA canvases with transparent corners and narrow centered subject bounds.
  Evidence: Pillow reports Skeleton bounds `(13, 11, 50, 53)` and Zombie bounds `(14, 6, 48, 58)`, providing useful scale and padding references for normalization.

- Observation: Raid rules are applied in slice order for every score interval, so adding rules after Skeleton and Zombie preserves deterministic simultaneous-threshold ordering without a new scheduler.
  Evidence: `internal/game/raid.go` iterates `s.raid.template.enemyRules` and spawns every newly reached multiple before moving to the next rule.

- Observation: The Ebitengine screenshot harness reproduced its known transient first-frame omission of screen-space UI once.
  Evidence: The first final capture contained the map and enemies but omitted the top bar and inspection panel; rerunning the same isolated test produced a complete 1920x1080 frame with both health bars and the selected Armoured Skeleton panel.

- Observation: Using the slower raider as the percentage baseline makes the closest current pair, the 0.9 Armoured Skeleton and 1.0 ordinary Skeleton, 11.11 percent apart.
  Evidence: `TestEnemyTemplateSpeedsRemainDistinct` passes all six distinct-template comparisons against the three-percent minimum, and `TestDistinctRaiderSpeedsSeparateMomentaryOverlap` proves the closest pair diverges after one update.

## Decision Log

- Decision: Add Ghoul and Armoured Skeleton as new accepted enemy templates and append their Raid rules at thresholds 6 and 8.
  Rationale: This implements the accepted roster while leaving Skeleton threshold 2, Zombie threshold 4, the challenge formula, and the square-root duration curve unchanged.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Use built-in image generation with the existing enemy sprites as style references, a flat `#ff00ff` key, local chroma-key removal, and nearest-neighbour normalization.
  Rationale: This follows the accepted art workflow and `ART.md`, avoids a new asset pipeline, and produces isolated pixel-art sprites that remain readable at 50% size.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Create a deterministic game-package screenshot fixture rather than alter production spawning or add debug UI.
  Rationale: A test-owned state can place both new enemies apart on a cleared road and select one while exercising the normal renderer, health bar, and inspection panel.
  Date/Author: 2026-07-14 / Codex

- Decision: Reduce the Armoured Skeleton from 1.0 to 0.9 Tiles per second and enforce at least three-percent separation between every pair of distinct raider templates.
  Rationale: The previous Armoured Skeleton exactly matched the ordinary Skeleton, so simultaneous spawns could remain perfectly overlapped. A 0.9 speed fits the heavy-armour role, produces a visible ten-percent difference, and is represented honestly by the inspection panel's one-decimal display.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Completed. The runtime now embeds distinct 64x64 Ghoul and Armoured Skeleton sprites and exposes accepted templates with the requested health, movement speed, Sanctum damage, rewards, sprite keys, and image references. The Armoured Skeleton moves at 0.9 Tiles per second, and automated coverage enforces at least three-percent separation between every pair of distinct raider templates. Unlike raiders placed at identical coordinates diverge after the next movement update. Generated Raids retain Skeleton and Zombie rules at score multiples of 2 and 4, then append Ghoul and Armoured Skeleton rules at 6 and 8. Raid 1 remains two Skeletons and one Zombie. Challenge 8 produces counts `4/2/1/1`, and challenge 16 produces 16 enemies; after four seconds its existing two Skeletons and one Zombie are active with 13 pending.

Normal tests, `go test -race ./...`, and 50 focused repetitions passed after the speed revision. Final image validation reports exact PNG/RGBA 64x64 sprites, transparent corners, nonempty native and 32x32 alpha bounds, and no visible magenta key pixels. `git diff --check` passes and the ownership check prints nothing. The two evidence PNGs were visually reviewed; the refreshed runtime frame shows the separated enemies, ordinary health bars, selected brightness, and the Armoured Skeleton inspection panel reporting `0.9 tiles/s`. The in-scope documents now agree on the four-enemy roster, distinct speeds, and unchanged challenge and duration formulas.

No hand-written Go file exceeds 600 lines. `internal/game/build_placement_test.go` remains unchanged at 564 lines; its earlier recommendation still applies: split housing, economic, and defense placement scenarios before the next substantive expansion. `internal/game/raid_test.go` is now 541 lines; before another large Raid slice, move roster generation, score-boundary, and raider-speed scenarios into responsibility-named test files. Both splits are outside this accepted plan and require user approval before implementation, so no unplanned refactor was performed.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `assets/catalog.go` explicitly embeds every required PNG, loads it into Ebitengine, and exposes enemy images through `assets.EnemySprites`. `internal/game/enemies.go` owns immutable enemy facts such as name, health, speed in Tiles per second, Sanctum damage, defeat resources, sprite key, and sprite pointer. Active instances in `internal/game/raid.go` refer to those templates, so new catalog entries automatically use existing movement, combat, health-bar, selection, inspection, reward, and breach behavior.

`internal/game/raid_generator.go` generates a deterministic private `raidTemplate`. Challenge still comes from Raid number, settlement population, and explored Plots, and progress duration remains `5 + 2 * sqrt(challenge)` seconds. Each `raidEnemyRule` schedules one enemy at every multiple of its threshold. Runtime progress in `internal/game/raid.go` compares previous and current challenge score for every ordered rule, which supports exact equality and updates that cross several thresholds at once.

The relevant control documents are `README.md`, `PRODUCT.md`, `GAME.md`, `ART.md`, and `ARCHITECTURE.md`. They must describe the expanded current roster, thresholds, enemy facts, generated sprites, and code ownership. `ROADMAP.md` and `DESIGN.md` stay unchanged because the accepted product direction and medieval wizardry design language do not change. `CODESTYLE.md` stays unchanged because no source convention changes. Historical plans remain immutable records.

## Plan of Work

First, generate two isolated pixel-art characters with the built-in image-generation workflow. Use `assets/sprites/enemies/skeleton-sword-shield.png` and `assets/sprites/enemies/zombie.png` as style references, request a perfectly flat `#ff00ff` background, prohibit shadows, scenery, text, and the key color inside the subject, and keep generous padding. The Ghoul is a lean, hunched, corpse-pale runner with long claws and a narrow silhouette. The Armoured Skeleton is visibly skeletal but much broader, with heavy dark-steel helmet, pauldrons, breastplate, and greaves. Copy generated sources into a plan-local working directory, remove the key with the installed imagegen helper, crop nontransparent bounds, resize with Pillow nearest-neighbour sampling, and center each result on an exact 64x64 transparent canvas. Save only the accepted runtime PNGs in `assets/sprites/enemies/`. Create `plans/59-ghouls-and-armoured-skeletons/screenshots/enemy-sprites-contact-sheet.png` showing both at 64x64 and 32x32 without smoothing.

Next, extend `assets/catalog.go` and `assets/catalog_test.go`. Add required embedded paths, `EnemySprites.Ghoul`, `EnemySprites.ArmouredSkeleton`, loading calls, catalog assignment, and exact 64x64 load assertions. Extend `internal/game/enemies.go` with matching `EnemyCatalog` fields. Ghoul must be named `Ghoul`, have 40 health, 1.5 Tiles per second, 1 Sanctum damage, rewards of 4 Wood and 1 Metal, key `ghoul`, and the loaded Ghoul sprite. Armoured Skeleton must be named `Armoured Skeleton`, have 125 health, 0.9 Tiles per second, 1 Sanctum damage, rewards of 5 Stone and 2 Metal, key `armoured-skeleton`, and the loaded Armoured Skeleton sprite. Add focused template tests for every field and sprite identity. Compare each pair of distinct templates using `(faster - slower) / slower` and require at least `0.03` separation so a later balance edit cannot silently restore equal or near-equal speeds.

Then, append `raidEnemyGhoul` and `raidEnemyArmouredSkeleton` after the existing kinds in `internal/game/raid_generator.go`, and append rules at thresholds 6 and 8 after the existing threshold-2 Skeleton and threshold-4 Zombie rules. Extend `State.enemyTemplateForRaidKind` in `internal/game/raid.go`. Tests must prove Raid 1 still has two Skeletons and one Zombie even though all four rules exist; exact score 6 and score 8 boundaries spawn the new types; one skipped update crossing both boundaries preserves rule order; challenge 8 schedules four Skeletons, two Zombies, one Ghoul, and one Armoured Skeleton; challenge 16 schedules 16 total enemies; and four seconds into challenge 16 the same early two Skeletons and one Zombie are active while 13 remain pending. Movement tests must prove Ghoul speed is 1.5, Armoured Skeleton speed is 0.9, and an ordinary and Armoured Skeleton placed at the same coordinates no longer overlap after one movement update.

Add a gated visual test in `internal/game/visual_test.go`. Construct a normal game state, clear natural terrain, start a Raid-owned state with one Ghoul and one Armoured Skeleton placed at separated north-road positions, reduce current health for visible health bars, and select one through `selectedItemRaider`. Capture `plans/59-ghouls-and-armoured-skeletons/screenshots/ghoul-and-armoured-skeleton.png` through the existing Ebitengine harness. Assert the fixture still exposes both templates and the selected-raider panel before capture.

Finally, update the five in-scope documents, run all validation, inspect both evidence images at full resolution, normalize ownership only if required, and update this living plan. Do not introduce dependencies, save changes, pathing changes, animation, exported behavioral APIs, or an asset pipeline.

## Concrete Steps

Run every command from `/home/dave/dev/ai/td`. After built-in image generation, use the installed helper and Pillow-based local normalization, then validate the deliverables with commands equivalent to:

    python3 "${CODEX_HOME:-$HOME/.codex}/skills/.system/imagegen/scripts/remove_chroma_key.py" --input <generated-source.png> --out <alpha.png> --auto-key border --soft-matte --transparent-threshold 12 --opaque-threshold 220 --despill
    python3 <plan-local normalization command or script>
    python3 <plan-local validation command or script>

The image validation must report RGBA mode, exact 64x64 dimensions, fully transparent corners, nonempty alpha bounds, no visible `#ff00ff` pixels, and readable nearest-neighbour 32x32 previews. Keep temporary generated sources outside runtime asset directories and remove them after final assets and evidence exist.

Format and test the implementation with:

    gofmt -w assets/catalog.go assets/catalog_test.go internal/game/enemies.go internal/game/enemies_test.go internal/game/raid.go internal/game/raid_generator.go internal/game/raid_generator_test.go internal/game/raid_test.go internal/game/visual_test.go
    go test ./...
    go test -race ./...
    go test ./assets ./internal/game -run 'TestNewCatalogLoads(Ghoul|ArmouredSkeleton)Sprite|TestNewEnemyCatalogIncludes(Ghoul|ArmouredSkeleton)|TestEnemyTemplateSpeedsRemainDistinct|TestGenerateRaid|TestRaidProgress|TestLaterRaid|TestGhoul|TestArmouredSkeleton|TestDistinctRaiderSpeedsSeparateMomentaryOverlap' -count=50
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureGhoulAndArmouredSkeletonScreenshot -count=1
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'

End validation with the required line-count review:

    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Report every hand-written Go file above or close to 600 lines with a concrete recommendation. Do not perform an unplanned split, refactor, or library addition without user approval. The unchanged `internal/game/build_placement_test.go`, currently 564 lines, is explicitly outside this feature's scope.

## Validation and Acceptance

Automated acceptance requires exact 64x64 embedded RGBA sprites, catalog templates with all accepted values and sprite references, and passing normal, race-enabled, and repeated focused tests. Raid 1 must remain three enemies: Skeleton, Skeleton, then Zombie at the score-4 simultaneous boundary. A challenge-8 template must total eight enemies across counts 4, 2, 1, and 1. A challenge-16 template must total 16; after four seconds it must have the existing three early enemies active and 13 pending.

Exact-boundary and skipped-threshold tests must show that Ghoul releases at score 6 and Armoured Skeleton at score 8, with ordered rules deciding simultaneous spawns. The Ghoul movement result must use 1.5 Tiles per second and the Armoured Skeleton must use 0.9 Tiles per second. Every pair of distinct raider templates must differ by at least three percent relative to the slower speed, and an ordinary and Armoured Skeleton starting together must have different positions after one movement update. Existing targeting, projectiles, defeat rewards, selection, health-bar, Sanctum contact, completion, duration, and challenge formulas must continue to pass unchanged.

Visual acceptance requires the contact sheet to show distinct, readable silhouettes at both 64x64 and 32x32. The focused 1920x1080 runtime screenshot must show both enemies separated on a cleared road, with one selected and the normal health bar plus inspection panel reporting `0.9 tiles/s`. Documentation acceptance requires `PRODUCT.md` and `GAME.md` to describe the 0.9 Armoured Skeleton speed and three-percent distinct-type guarantee while `README.md`, `ART.md`, and `ARCHITECTURE.md` retain their accurate roster, asset, and ownership statements.

## Idempotence and Recovery

Formatting, tests, screenshot capture, contact-sheet generation, and validation are safe to rerun. Screenshot generation overwrites only this plan's evidence files. Preserve the generated source until each sprite passes alpha and readability review; if chroma removal leaves a fringe, retry once with the helper's edge contraction before considering another generation. If a test exposes an invalid kind, inspect the kind-to-template switch rather than changing pending counts. If visual capture omits UI on its first Ebitengine frame, rerun the isolated capture as prior evidence shows that can be transient.

## Artifacts and Notes

Final runtime sprites belong at:

    assets/sprites/enemies/ghoul.png
    assets/sprites/enemies/armoured-skeleton.png

Deterministic evidence belongs at:

    plans/59-ghouls-and-armoured-skeletons/screenshots/enemy-sprites-contact-sheet.png
    plans/59-ghouls-and-armoured-skeletons/screenshots/ghoul-and-armoured-skeleton.png

At completion, record the final built-in prompts, image validation summary, test results, and visual review in this plan.

The final built-in prompt set used the existing Skeleton and Zombie PNGs as style-only references. The Ghoul prompt requested one isolated full-body, lean, hunched, corpse-pale runner with long separated claws, a narrow dynamic silhouette, low-detail hard-cluster pixel art, no shadow or props, and a perfectly uniform `#ff00ff` key. The Armoured Skeleton prompt requested one isolated full-body, visibly skeletal but substantially broader figure with a dark-steel helmet, large pauldrons, breastplate, gauntlets, and greaves, bright bone and steel contrast, no weapon or shield, no shadow or props, and the same uniform key. Both prompts required clear 64x64 and 32x32 readability, generous padding, no key color in the subject, and no scenery, text, border, or watermark.

Final validation recorded Ghoul alpha bounds `(10, 6, 53, 58)` and half-size bounds `(5, 3, 26, 29)`. Armoured Skeleton bounds are `(11, 4, 53, 60)` and half-size bounds `(5, 2, 26, 30)`. Both have zero-alpha corners and zero visible magenta pixels.

## Interfaces and Dependencies

No external dependency, serialized format, pathing, targeting, animation, or exported behavioral API changes. The existing public project types gain fields only:

    type EnemySprites struct {
        SkeletonSwordShield *ebiten.Image
        Zombie              *ebiten.Image
        Ghoul               *ebiten.Image
        ArmouredSkeleton    *ebiten.Image
    }

    type EnemyCatalog struct {
        SkeletonSwordShield EnemyTemplate
        Zombie              EnemyTemplate
        Ghoul               EnemyTemplate
        ArmouredSkeleton    EnemyTemplate
    }

`raidEnemyGhoul` and `raidEnemyArmouredSkeleton` remain private game-package kinds. Pillow is used only as an already-available development-time image tool and is not added to `go.mod` or the runtime.

Revision note (2026-07-14): Created from the accepted implementation plan after fresh repository and workflow inspection. It fixes exact art, catalog, Raid, test, evidence, documentation, ownership, and line-count requirements before implementation.

Revision note (2026-07-14): Completed implementation and updated every living section with the built-in prompt set, normalized sprite evidence, deterministic roster results, transient capture recovery, full validation, documentation status, and final line-count review.

Revision note (2026-07-14): Reopened the completed plan to incorporate the user's distinct-raider-speed requirement. The revision changes the Armoured Skeleton to 0.9 Tiles per second, defines the three-percent comparison, adds regression and movement-divergence coverage, and requires refreshed visual and documentation evidence.

Revision note (2026-07-14): Completed the distinct-speed revision, refreshed the runtime evidence, synchronized `PRODUCT.md` and `GAME.md`, and recorded successful normal, race, repeated, whitespace, ownership, and line-count validation.
