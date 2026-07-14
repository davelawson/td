# Generate Scaling Raids from Settlement State

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/52-dynamic-raid-generation.md`.

## Purpose / Big Picture

After this change, Raids no longer follow a fixed five-enemy opening and two-enemy growth rule. Starting a Raid generates a deterministic challenge rating from the Raid number, the settlement's total population, and all explored Plots. Raid progress then rises from zero to 100 percent and adds Skeletons and Zombies whenever the progress-scaled challenge reaches another multiple of that enemy's threshold. Larger settlements, broader Domains, and later Raids therefore create longer assaults containing more enemies without changing individual enemy stats.

The behavior is visible by building population or exploring before clicking `Next Raid`, then observing a generated Raid whose enemies join over its progress window. Automated tests expose the exact challenge, duration, thresholds, and spawn boundaries; screenshot evidence shows generated enemies on the existing north road.

## Progress

- [x] (2026-07-14 12:30Z) Inspected the fixed Raid lifecycle, enemy catalog, population and Plot state, tests, documentation, screenshot workflow, and existing dirty worktree.
- [x] (2026-07-14 12:30Z) Confirmed the challenge formula, duration scaling, thresholds, deterministic ordering, live input meanings, and completion boundary; created this ExecPlan.
- [x] (2026-07-14 13:37Z) Added the pure private Raid generator, normalized inputs, challenge formula, duration, ordered rules, and schedule helpers in `internal/game/raid_generator.go`.
- [x] (2026-07-14 13:37Z) Replaced fixed count/stagger generation with progress-threshold spawning integrated into transient Raid state.
- [x] (2026-07-14 13:37Z) Added pure generator coverage and revised lifecycle, pause, composition, state-input, and completion tests.
- [x] (2026-07-14 13:37Z) Updated current product, game-design, roadmap, architecture, onboarding, tower-data consistency, and the 17-image visual evidence suite.
- [x] (2026-07-14 13:37Z) Passed automated, race, 50-run repeated, screenshot, whitespace, ownership, and line-count validation.

## Surprises & Discoveries

- Observation: Current Raid generation is split between a fixed count formula and a special composition rule.
  Evidence: `raidEnemyCount` returns five plus two per Raid, while `nextRaidEnemyTemplate` alternates Zombies only in Raid 1 and uses Skeletons thereafter.
- Observation: The existing HUD can continue reporting generated enemies without a new progress UI.
  Evidence: `raidEnemiesRemaining` already adds pending and active enemies, so pending can represent enemies scheduled by the generated template but not yet spawned.
- Observation: The repository contains accepted, uncommitted Labour/Management changes and user-authored tower-template changes.
  Evidence: the planning `git status --short` listed those files as modified and plans 51/phase files as untracked; this work must preserve them.
- Observation: Directly reading the Ebitengine window backbuffer produced intermittent partial screenshot frames after the longer Raid fixture.
  Evidence: repeated captures sometimes omitted HUD or building-bar regions even though gameplay rendering and tests were correct. Rendering each target into an offscreen image, reading that image before presenting it, and recapturing produced complete repeatable evidence.

## Decision Log

- Decision: Calculate challenge as `1.2^(raidNumber-1) * 1.2^(plotsExplored-1) * (1 + population/10) + 3` using floating-point division.
  Rationale: Raid sequence, Domain expansion, and population each increase pressure, with exponential progression for Raid and Plot growth and proportional settlement scaling.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Use total inhabitants across all roles and count the home Plot among explored Plots.
  Rationale: Staffing reservations should not lower the settlement strength seen by the Rival, and the initial defended territory is still part of the Domain.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Use Skeleton threshold 2 and Zombie threshold 4 in stable Skeleton-then-Zombie rule order.
  Rationale: Zombies are tougher and should join less often; stable rule order makes simultaneous crossings deterministic.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Set Raid progress duration to `5 + challenge` seconds.
  Rationale: Higher-rated Raids receive a longer spawn window while their threshold-crossing rate still increases gradually with challenge.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Spawn when a score reaches or crosses a threshold multiple and complete only after 100 percent progress and removal of all enemies.
  Rationale: Exact endpoint multiples must not disappear, and reaching the scheduling endpoint should not erase surviving combatants.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Keep generation deterministic and leave all `EnemyTemplate` stats unchanged.
  Rationale: The requested scaling applies to Raid contents, while reproducibility remains important for tests and balance comparisons.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Implemented a private `raidTemplate` generator whose exact three integer inputs are Raid number, total settlement population, and explored Plot count. The runtime now begins at zero progress with no immediate spawn, advances across a `5 + challenge` second window, releases every threshold crossing without losing skipped multiples, and waits for both 100 percent progress and an empty active roster before success. The baseline Raid has challenge 4, duration 9 seconds, two Skeletons, and one Zombie. Enemy catalog stats, combat, movement, rewards, breach, Labour, Management, and the existing HUD contract remain intact.

The five relevant control documents now describe the formula, input meanings, thresholds, stable order, completion boundary, and fixed-path limitation. The 17-image suite under `plans/52-dynamic-raid-generation/screenshots/` was captured and reviewed. `active-raid.png` shows the first generated Skeleton after its threshold, and `selected-raider.png` shows the same generated enemy selected with 50/50 health and 1.0 Tiles-per-second speed. No progress or challenge UI was added.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, `go test ./internal/game -count=50`, screenshot capture, `git diff --check`, and the ownership check all succeeded. No hand-written Go file exceeds 600 lines. The unchanged `internal/game/building_bar_test.go` remains closest at 595 lines; split its coverage by building-bar behavior before its next substantive expansion. `internal/game/build_placement_test.go` is 537 lines, while the revised `internal/game/raid_test.go` is 393 lines and the new generator has focused source and test files. No unplanned file split was necessary.

## Context and Orientation

`internal/game/raid.go` owns transient Raid state, Raid startup, fixed stagger spawning, movement, breach, completion, and remaining-enemy counts. It currently uses `firstRaidEnemyCount`, `raidEnemyGrowth`, and `raidSpawnInterval`, immediately spawns one enemy, alternates Skeletons and Zombies only in Raid 1, and uses Skeletons alone later. `internal/game/enemies.go` owns immutable shared stats and sprites for the two enemy types. `internal/game/combat.go` defines the fixed update duration as 1/60 second.

Population state lives in `gameStatus` as available/total counts for Apprentices, Soldiers, and Peasants. Challenge uses the sum of the three total counts, not availability. `Map.exploredPlotCoordinates` returns every explored Plot in deterministic order and includes home, so its length is the generator's Plot input.

The top HUD displays pending plus active enemies during a Raid. No challenge or progress display is requested. Existing pause and overlay flow stops calls to Raid update logic, and that same boundary must freeze generated progress. Combat, movement, rewards, breach, Labour, and Management remain unchanged after enemies spawn.

## Plan of Work

Create `internal/game/raid_generator.go`. Define private enemy-kind, enemy-rule, and Raid-template types. A template stores a `float64` challenge rating, a `float64` progress duration in seconds, and an ordered rule slice. Implement `generateRaid(raidNumber, settlementPopulation, plotsExplored int) raidTemplate` with exactly those three integer parameters. Normalize Raid and Plot inputs to at least one and population to at least zero, then calculate the accepted formula and duration. Include Skeleton threshold 2 followed by Zombie threshold 4. Add pure helpers that count how many enemies a rule has scheduled at a score and how many total enemies the template schedules at 100 percent.

Change `raidState` to store its generated template and normalized progress from zero through one. Remove the old spawn countdown and fixed count/composition helpers. On Raid start, increment the number, sum total population, count explored Plots, generate the template, initialize pending enemies from its full schedule, clear active enemies and IDs, and set progress to zero without immediately spawning.

During each active update, advance progress by `deltaSeconds / duration`, clamped to one. Compare the previous and current `progress * challenge` scores. For each rule in template order, calculate the difference between the counts scheduled at those scores and spawn that many mapped enemy templates. Equality counts as reached. Newly spawned enemies receive the existing stable IDs, current farthest-north spawn position, full template health, and unchanged stats. Pending count decreases for every spawn.

After progress handling, retain combat and movement. Completion requires progress at one, zero pending enemies, and no active enemies. Before 100 percent, an empty active list is a waiting interval rather than completion. At 100 percent, surviving enemies remain active until defeated or breach. The HUD continues to show pending plus active enemies.

Add `internal/game/raid_generator_test.go` for pure formula and threshold tests. Rewrite obsolete fixed-count, fixed-stagger, Raid-1 alternation, and later-Skeleton-only tests in `raid_test.go` as generator integration tests. Preserve coverage for UI start, no double start, movement, health, completion, breach, pause, overlay, and resource outcomes. Add explicit tests for live population/Plot inputs, exact threshold equality, simultaneous score-4 ordering, no early completion, final cleanup, and unchanged enemy stats.

Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, and `ARCHITECTURE.md` to describe deterministic challenge generation rather than fixed Raid composition. Mark the original fixed-placeholder Raid decision in `GAME.md` as superseded and record the accepted formula and threshold design. `DESIGN.md` and `ART.md` remain unchanged because no durable visual language or asset guidance changes.

Point `cmd/td/main_test.go` screenshot output at `plans/52-dynamic-raid-generation/screenshots/` and advance active-Raid fixtures far enough to cross the first Skeleton threshold. Capture the full suite and verify the active and selected-raider screenshots still show generated enemies on the north road without HUD or selection regressions.

## Concrete Steps

Run every command from `/home/dave/dev/ai/td`.

Implement the generator, Raid integration, and tests, then format changed Go files and run:

    gofmt -w <changed Go files>
    go test ./internal/game
    go test ./...

Update the five relevant root documents and screenshot destination, then capture evidence:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Complete validation with:

    go test ./...
    go test -race ./...
    go test ./internal/game -count=50
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Record all results here. If any hand-written Go file exceeds 600 lines or is close enough to cross on its next likely edit, record the path and line count, recommend a responsibility-based response, and request approval before any unplanned split.

## Validation and Acceptance

Pure generator tests must prove challenge 4 and duration 9 for Raid 1, population 0, and one Plot; exact larger-input formula behavior; input normalization; fixed rule order and thresholds; deterministic repeated output; and scheduled totals computed with floor division. Formula comparisons should use a small floating-point tolerance.

Runtime tests must prove a started Raid has zero progress and no active enemies, while pending count includes its full generated schedule. No enemy may spawn before score 2. Reaching score 2 spawns one Skeleton. Reaching score 4 schedules the second Skeleton and first Zombie in stable order. Pause and overlay must preserve progress. Empty intervals before 100 percent must not complete the Raid. At 100 percent, pending must be zero and surviving enemies must remain; completion occurs only after they are gone. Live total population and all explored Plots must change the generated challenge. Enemy health, speed, damage, rewards, and sprites must still come directly from the existing catalog.

The full test, race, repeated, screenshot, whitespace, and ownership checks must pass. Visual evidence must show at least one threshold-generated enemy in both active-Raid fixtures. Documentation must agree on inputs, formula, duration, thresholds, deterministic ordering, stat invariance, and completion.

## Idempotence and Recovery

Generation is pure and deterministic, and formatting, tests, and screenshot capture are safe to rerun. Screenshot capture replaces only plan 52 evidence. If a floating-point update steps over more than one multiple, spawn the full count difference rather than dropping enemies. If several rules cross together, retain template order.

Do not revert or overwrite the existing uncommitted phase, tower, documentation, test, plan 51, or screenshot changes. The generator should replace only obsolete fixed Raid generation behavior. No save format exists, so no migration or rollback data is required.

## Artifacts and Notes

The plan 51 active-Raid screenshots are the visual baseline. Post-change evidence belongs under `plans/52-dynamic-raid-generation/screenshots/`, with special review of `active-raid.png` and `selected-raider.png`.

## Interfaces and Dependencies

No exported API or external dependency changes. The private entry point is `generateRaid(raidNumber, settlementPopulation, plotsExplored int) raidTemplate`. The generator uses Go's standard `math` package. `raidState` gains private template and progress fields and loses its spawn countdown. Existing `EnemyTemplate`, `EnemyCatalog`, `Input`, `State`, and cross-package interfaces retain their signatures.

Revision note (2026-07-14): Created from the accepted implementation plan after repository inspection. It fixes the formula, duration, thresholds, deterministic crossing semantics, live state inputs, completion, documentation, screenshots, and validation.
