# Accelerate Later Raid Tempo

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must remain current as work proceeds. Maintain this document in accordance with `PLANS.md` from the repository root.

## Purpose / Big Picture

Later Raids currently add enemies and almost the same amount of time needed to release them, so their larger rosters make assaults much longer without creating enough additional pressure. After this change, the first Raid retains its accepted nine-second release window while later Raids release their growing rosters over a sublinearly growing window. A player can see the difference by comparing a challenge-16 Raid after four seconds: the old schedule has released one enemy, while the new schedule has released two Skeletons and one Zombie.

## Progress

- [x] (2026-07-14 19:09Z) Inspected Raid generation, execution, tests, visual evidence, durable documentation, and repository state.
- [x] (2026-07-14 19:09Z) Recorded the accepted square-root duration curve and unchanged enemy-count rules.
- [x] (2026-07-14 19:12Z) Captured and reviewed challenge-16 pre-change evidence after four simulated seconds; one Skeleton is visible.
- [x] (2026-07-14 19:16Z) Implemented the square-root duration formula plus exact baseline, later-duration, density, and four-second runtime tests; `go test ./...` passes.
- [x] (2026-07-14 19:17Z) Updated `GAME.md`, `PRODUCT.md`, and `ARCHITECTURE.md` to agree on the new pacing and unchanged enemy rules.
- [x] (2026-07-14 19:17Z) Captured and reviewed matching before/after evidence and completed normal, race, repeated, whitespace, and ownership validation.
- [x] (2026-07-14 19:17Z) Checked hand-written Go line counts. No file exceeds 600 lines; the unchanged `internal/game/build_placement_test.go` remains closest at 564 lines.

## Surprises & Discoveries

- Observation: Spawn density already rises slightly because score is `progress * challenge`, but the current `5 + challenge` duration approaches linear growth and therefore preserves much of the later Raid length.
  Evidence: `internal/game/raid_generator.go` calculates duration as `5 + challenge`, while enemy totals are `floor(challenge/2) + floor(challenge/4)`.

- Observation: Ebitengine can occasionally omit screen-space UI from the first focused capture even though the game-state assertions pass.
  Evidence: The first post-change capture contained the correct map and enemies but no top bar; rerunning the isolated screenshot test produced the complete 1920x1080 frame.

## Decision Log

- Decision: Calculate the release window as `5 + 2 * sqrt(challenge)` seconds.
  Rationale: This keeps challenge 4 at nine seconds, changes challenge 16 from 21 to 13 seconds, and changes challenge 36 from 41 to 17 seconds. The player accepted this curve because it raises tempo progressively without making every Raid use one fixed duration.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Do not change challenge generation, thresholds, enemy templates, or completion behavior.
  Rationale: The reported issue is Raid length as monster counts grow. Changing only duration isolates that balance adjustment and preserves the already-tested meaning of challenge.
  Date/Author: 2026-07-14 / User and Codex

- Decision: Compare screenshots of the same challenge-16 state after four simulated seconds.
  Rationale: Timing is the user-visible change. At four seconds the old curve has crossed only score 2, while the new curve has crossed score 4, making the increased release tempo visible without changing rendering code.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Raid release duration now follows `5 + 2 * sqrt(challenge)` seconds. The challenge-4 baseline remains nine seconds with the same three enemies, while challenges 16 and 36 now use 13- and 17-second windows instead of 21 and 41. Enemy thresholds, totals, stable spawn order, templates, combat, movement, completion, breach, and post-Raid behavior are unchanged.

Focused generator and runtime tests cover the exact durations, increasing window length, decreasing seconds per enemy, and challenge-16 behavior after four seconds. `go test ./...`, `go test -race ./...`, and 50 repetitions of the generator and Raid tests pass. `git diff --check` passes, and the ownership check prints nothing. `GAME.md`, `PRODUCT.md`, and `ARCHITECTURE.md` agree on the new formula.

The deterministic screenshots under `plans/58-accelerating-raid-tempo/screenshots/` show the same cleared challenge-16 map after 240 updates. The legacy frame has one Skeleton; the accelerated frame has two Skeletons and one Zombie, with the two enemies released together at score 4 still close to the spawn point.

No hand-written Go file exceeds 600 lines. The changed files are `internal/game/raid_generator.go` at 74 lines, `raid_generator_test.go` at 131, `raid_test.go` at 429, and `visual_test.go` at 232. The unchanged `internal/game/build_placement_test.go` remains nearest the preference at 564 lines; split its housing, economic, and defense scenarios before its next substantive expansion. No unplanned refactor was needed for this change.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `internal/game/raid_generator.go` creates a private `raidTemplate` from Raid number, total settlement population, and explored Plot count. Challenge determines both enemy totals and their score thresholds. A Skeleton is scheduled at every score multiple of 2, and a Zombie is scheduled at every score multiple of 4. `internal/game/raid.go` advances normalized progress through the template duration and spawns every threshold crossed during an update, preserving rule order when multiple thresholds coincide.

The existing duration is `5 + challenge` seconds. The baseline challenge is 4, so the first Raid lasts nine seconds and contains two Skeletons and one Zombie. The accepted replacement is `5 + 2 * sqrt(challenge)`, which also yields nine seconds at challenge 4 but grows more slowly thereafter. Enemy movement, health, damage, resources, sprites, tower combat, pause behavior, Sanctum breach, post-Raid Labour, and the rule that completion waits for progress plus all enemies remain unchanged.

`GAME.md` records the intended Raid mechanics and their decision history. `PRODUCT.md` records the current user-visible prototype behavior. `ARCHITECTURE.md` records generator ownership and the runtime lifecycle. These three documents must agree with the implemented formula. `README.md` does not state a duration formula, and `ROADMAP.md`, `DESIGN.md`, `ART.md`, and `CODESTYLE.md` do not require changes because this work does not alter setup, strategic scope, visual language, assets, or source conventions. Historical ExecPlans remain unchanged as records of their completed work.

## Plan of Work

First, add focused optional screenshot fixtures to `internal/game/visual_test.go`. Give both states a cleared map with fixed frontier labels and 120 total Peasants, then start Raid 1 to produce challenge 16 and advance exactly 240 logical updates. The before fixture substitutes the legacy 21-second duration and verifies one active Skeleton; the after fixture uses the generated 13-second duration and verifies three active enemies. Save them under the matching `before` and `after` evidence directories.

Next, revise `internal/game/raid_generator.go`. Retain the five-second base and add a duration scale of two. Put the calculation in a private `raidProgressDuration(challengeRating float64) float64` helper and have `generateRaid` use it. The helper returns `5 + 2 * math.Sqrt(challengeRating)` for the generated challenge domain; no exported API or dependency is added because `math` is already used.

Expand `internal/game/raid_generator_test.go` to verify the baseline, exact challenge-16 and challenge-36 durations, monotonic duration growth, and falling duration-per-scheduled-enemy as challenge increases. Add or adapt a runtime test in `internal/game/raid_test.go` to prove that after four seconds a challenge-16 Raid has spawned two Skeletons and one Zombie in stable order, with the remaining scheduled enemies still pending. Preserve all existing lifecycle and threshold tests.

Update the implemented Raid description and its decision-log entry in `GAME.md`, the current capability summary in `PRODUCT.md`, and the lifecycle invariant in `ARCHITECTURE.md`. Describe the square-root formula and its purpose accurately without suggesting that enemy totals or individual stats changed.

Finally, capture `plans/58-accelerating-raid-tempo/screenshots/after/raid-tempo.png` from the matching fixture and confirm that three raiders are represented at the same four-second boundary. Run the complete validation suite, review the diff, normalize ownership if necessary, update all living sections of this plan, and perform the required line-count audit.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`. Create the baseline fixture, capture it before editing the generator, then implement and format the accepted changes. Use:

    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureRaidTempoBeforeScreenshot -count=1
    gofmt -w internal/game/raid_generator.go internal/game/raid_generator_test.go internal/game/raid_test.go internal/game/visual_test.go
    go test ./...
    go test -race ./...
    go test ./internal/game -run 'TestGenerateRaid|TestRaid|TestLaterRaid' -count=50
    TD_CAPTURE_SCREENSHOT=1 go test ./internal/game -run TestCaptureRaidTempoAfterScreenshot -count=1
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

The ownership command must print nothing. The final line-count command must be reviewed for every hand-written Go file above or close to 600 lines. Record each such file and a concrete recommendation in `Outcomes & Retrospective`; do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

Automated acceptance requires all package tests and race-enabled tests to pass. Challenge 4 must retain a 9-second window and three scheduled enemies. Challenges 16 and 36 must use 13- and 17-second windows. The challenge-16 runtime fixture must release one Skeleton at score 2, then a second Skeleton followed by one Zombie when score 4 is reached. Four elapsed seconds must be enough to release those three enemies under the new formula but not under the old formula.

No existing challenge inputs, enemy totals, template stats, pause boundaries, threshold equality rules, skipped-threshold handling, completion conditions, breach behavior, or economic results may change. The before/after screenshots must use the same state and elapsed time and visibly show the denser new release. Documentation acceptance requires `GAME.md`, `PRODUCT.md`, and `ARCHITECTURE.md` to agree on the formula and unchanged boundaries.

## Idempotence and Recovery

Formatting, tests, repeated tests, and screenshot capture are safe to rerun. Screenshot capture overwrites only this plan's evidence files. If the screenshot fixture breaches unexpectedly, verify that exactly four seconds were simulated and that no enemy speed or path state was changed. If floating-point assertions fail, compare with a small tolerance except for perfect-square challenges whose expected durations are exact. Preserve unrelated work if the worktree changes during implementation.

## Artifacts and Notes

Visual evidence belongs at `plans/58-accelerating-raid-tempo/screenshots/before/raid-tempo.png` and `plans/58-accelerating-raid-tempo/screenshots/after/raid-tempo.png`. The before image should show one spawned Skeleton; the after image should show two Skeletons and one Zombie after the same four simulated seconds. Record final command results and image review in the living sections above.

## Interfaces and Dependencies

No exported API, serialized data, external service, or new module dependency changes. `generateRaid(raidNumber, settlementPopulation, plotsExplored int) raidTemplate` and `raidTemplate` retain their signatures and fields. Add only the private helper `raidProgressDuration(challengeRating float64) float64` and a private duration scale constant in `internal/game/raid_generator.go`.

Revision note (2026-07-14): Created from the accepted implementation plan after repository inspection. It fixes the square-root curve, unchanged gameplay boundaries, focused timing evidence, documentation, validation, and line-count requirements.

Revision note (2026-07-14): Completed implementation and updated every living section with the deterministic legacy comparison fixture, Ebitengine capture behavior, automated and visual results, documentation status, and final line-count review.
