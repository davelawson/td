# Split the Peaceful Phase into Labour and Management

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root and is stored at `plans/51-peaceful-phase-split.md`.

## Purpose / Big Picture

After this change, peaceful play has a deliberate cadence. A successful Raid advances to the next Day, workers in every placed economic building immediately perform one Labour phase and add their resources, and the game then enters Management. During Management the player can spend those resources to build structures and explore Plots for as long as desired; clicking `Next Raid` ends Management and starts combat. A new game starts in Management because no preceding Raid exists.

The change is visible by starting a game and reading `Management phase` in the top bar instead of the unused two-minute countdown. It is also observable across a full cycle: build a Woodcutter, defeat the next Raid, and see its 10 Wood arrive before Day 2 Management begins.

## Progress

- [x] (2026-07-14 11:59Z) Inspected the existing calm/Raid lifecycle, economic payouts, action gates, HUD, control documents, screenshot workflow, baseline tests, and clean working tree.
- [x] (2026-07-14 11:59Z) Confirmed the accepted cadence and created this ExecPlan.
- [x] (2026-07-14 12:07Z) Replaced calm state and its unused timer with explicit Labour, Management, and Raid lifecycle behavior.
- [x] (2026-07-14 12:07Z) Restricted exploration, construction, and Raid start controls to Management and updated player-facing production wording.
- [x] (2026-07-14 12:07Z) Added focused lifecycle and action-gating tests without expanding the near-limit building-bar test file.
- [x] (2026-07-14 12:07Z) Updated durable product, game-design, roadmap, architecture, visual-direction, and onboarding documentation.
- [x] (2026-07-14 12:08Z) Captured and reviewed 17 post-change screenshots under this plan directory.
- [x] (2026-07-14 12:08Z) Ran automated, race, repeated, whitespace, and ownership validation successfully.
- [x] (2026-07-14 12:08Z) Checked hand-written Go file line counts; no file exceeds 600 lines, while the unchanged 596-line `building_bar_test.go` remains near the preference.

## Surprises & Discoveries

- Observation: The displayed `Raid in 02:00` value never counts down.
  Evidence: `gameStatus.calmTime` is initialized and reset to 120 but no update path decrements it.
- Observation: Economic production already happens atomically at successful Raid completion, so this slice changes lifecycle ownership and player-facing semantics without requiring a timer or new economic calculation.
  Evidence: `completeRaid` calls `grantEconomicBuildingResources` once before returning to the current calm state.
- Observation: The baseline test suite passes before implementation.
  Evidence: `go test ./...` passed for every package on 2026-07-14.
- Observation: The instantaneous Labour transition needs no new rendered state to remain understandable in the current UI.
  Evidence: `running-game.png` clearly shows `Management phase`, while `woodcutter-tooltip.png` explains both that the Peasant works during Labour and that production occurs during each Labour phase.

## Decision Log

- Decision: Treat Peaceful as the design-level umbrella and Labour, Management, and Raid as the concrete runtime phases.
  Rationale: Labour and Management have distinct permissions and ordering, while no runtime behavior needs a separate nested Peaceful state.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Resolve Labour immediately at the start of each post-Raid Day, then enter Management in the same call.
  Rationale: The user specified that workers perform Labour immediately and Management begins immediately afterward; a timer, animation, button, or dismissible report would contradict that cadence.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Start Day 1 directly in Management.
  Rationale: Labour begins following a completed Raid, and a new game has no preceding Raid whose work needs resolution.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Use changed resource totals as the only Labour feedback in this slice.
  Rationale: This preserves an atomic transition and avoids introducing notification or report UI before there is a broader settlement-event presentation design.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Keep current economic-building yields as the only implemented Labour work.
  Rationale: Manual gathering, worker reassignment, recruitment, resource nodes, and additional jobs remain separate gameplay slices.
  Date/Author: 2026-07-14 / User and Codex

## Outcomes & Retrospective

Implementation is complete. The former calm state and static two-minute display are gone. A new game now starts in Management. Successful Raid completion clears the active Raid, increments the Day, enters Labour, grants every economic building's yield once across all explored Plots, and immediately enters Management. Breach remains terminal in the Raid phase and does not advance the Day or resolve Labour. Construction and exploration require Management but retain their paused-Management behavior, while `Next Raid` requires unpaused Management.

Player-facing structure descriptions and tooltips now identify Labour as the production phase. `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md` agree on the same cadence, and `GAME.md` records the old timed-calm decision as superseded. Manual gathering, worker reassignment, recruitment, additional jobs, a Labour report, and a Labour timer remain intentionally unimplemented.

The 17-image suite in `plans/51-peaceful-phase-split/screenshots/` passed visual review. `running-game.png` shows the centered `Management phase` label with no HUD regression. `woodcutter-tooltip.png` shows the new Labour wording without clipping or overlap.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, `go test ./internal/game -count=50`, screenshot capture, `git diff --check`, and the ownership check all succeeded. No hand-written Go file exceeds 600 lines. The unchanged `internal/game/building_bar_test.go` remains the closest at 596 lines; its next substantive feature should split tests by behavior before adding coverage. `internal/game/build_placement_test.go` remains 537 lines, and new phase coverage lives in the focused `internal/game/phase_test.go` instead of expanding either near-limit file.

## Context and Orientation

`internal/game/hud.go` currently defines a private two-value phase type, initializes new games in `phaseCalm`, stores a fixed `calmTime`, and formats the top-center HUD as a static two-minute Raid countdown. `internal/game/raid.go` starts Raids, resolves breach, and completes successful Raids. Successful completion currently pays economic buildings, returns to calm, increments the Day, and resets the unused timer. `internal/game/resources.go` walks every explored Plot and adds each economic building's template yield.

`internal/game/exploration.go` and `internal/game/building_bar.go` allow their actions only in `phaseCalm`. `internal/game/raidui.go` renders `Next Raid`, while `canStartRaid` in `internal/game/raid.go` currently checks pause, active-Raid, and breach state but not the explicit phase. The SPACE pause behavior deliberately allows peaceful exploration and construction while stopping logical simulation; the in-game overlay blocks these inputs entirely.

Economic structure descriptions live in `internal/game/structures.go`, and their detailed production wording is composed in `internal/game/building_tooltip.go`. `cmd/td/main_test.go` captures the existing visual evidence suite and currently writes it to plan 50. The running-game screenshot is the main evidence for the top-bar phase wording; economic tooltip screenshots prove the production terminology.

The durable documentation must agree with the implemented behavior. `GAME.md` owns the intended phase design and contains a now-superseded timed-calm decision. `PRODUCT.md` and `README.md` own current user-visible behavior. `ROADMAP.md` describes current capability and future gathering work. `ARCHITECTURE.md` records state ownership and transition invariants. `DESIGN.md` names which phase owns construction and exploration presentation. `ART.md` is unaffected because no asset or art-generation guidance changes.

## Plan of Work

Define private `phaseLabour`, `phaseManagement`, and `phaseRaid` values and remove `phaseCalm` plus `gameStatus.calmTime`. Initialize a new game in Management. Format Management as `Management phase`; retain explicit Labour and Raid formatting branches even though the atomic Labour transition is normally not drawn.

Make successful Raid completion clear active combat, increment the Day, enter Labour, grant every placed economic building's yield once, and immediately enter Management. Keep breach terminal in the Raid state: it must not increment the Day, enter Labour, or pay producers. Require Management for `canStartRaid`, construction, and exploration. Preserve the existing rule that pause disables `Next Raid` but still allows exploration and construction during Management, while the overlay blocks all game input.

Rename comments and player-facing economic text from post-Raid payout language to Labour production. Tooltips should say `Production: +10 <Resource> during each Labour phase`, and descriptions should explain that the assigned Peasant works during Labour. Do not rename the general `ResourceYield` template field or add a worker-assignment abstraction; its existing data shape remains suitable.

Add a focused phase test file for lifecycle and permission coverage rather than growing `internal/game/building_bar_test.go`, which is already 596 lines. Update existing initial-state, Raid-completion, breach, exploration, construction, structure-description, and tooltip assertions where their terminology or expected phase changes. Tests must show that all producers across explored Plots pay once, a producer built during Management waits until the following Labour resolution, successful completion ends in Management on the next Day, and breach pays nothing.

Update `README.md`, `PRODUCT.md`, `ROADMAP.md`, `GAME.md`, `DESIGN.md`, and `ARCHITECTURE.md`. In `GAME.md`, replace open and decided timed-calm language with the accepted instantaneous Labour and player-ended Management cadence, and add a superseding decision-log entry rather than deleting history. Change the screenshot destination to `plans/51-peaceful-phase-split/screenshots`, capture the complete existing suite, and inspect `running-game.png` plus economic tooltip screenshots for the new wording.

## Concrete Steps

Run every command from `/home/dave/dev/ai/td`.

Edit the phase, Raid, resource, action-gating, structure-description, tooltip, and test files described above. Format all changed Go files, then run:

    gofmt -w <changed Go files>
    go test ./...

Update the six relevant control and onboarding documents plus `cmd/td/main_test.go`. Capture post-change evidence with:

    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

Review the generated screenshots, especially `plans/51-peaceful-phase-split/screenshots/running-game.png` and the economic-building tooltip images. Complete validation with:

    go test ./...
    go test -race ./...
    go test ./internal/game -count=50
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Record results in this plan. Report every hand-written Go file over 600 lines and any file close enough to cross the preference on its next likely edit. Recommend a responsibility-based response but do not perform an unplanned split without user approval.

## Validation and Acceptance

`go test ./...`, `go test -race ./...`, and 50 repeated `internal/game` runs must pass. A new game must be Day 1 Management with `Management phase` in the HUD. `Next Raid` must work only during unpaused Management. Exploration and construction must work during Management, including SPACE-paused Management, and fail during Labour, Raid, and breach.

A successful Raid must advance the Day exactly once, resolve each placed economic building exactly once across all explored Plots, and finish in Management. A structure constructed during Management must not produce immediately and must first pay in the Labour phase following the next successful Raid. A breach must leave resources and Day unchanged and keep future Raids disabled.

The screenshot suite must show the Management label and Labour wording without layout regression. The six updated documents must consistently describe Peaceful as Labour followed by Management, identify Management as the only current construction/exploration phase, and retain manual gathering, reassignment, recruitment, and resource nodes as unimplemented.

## Idempotence and Recovery

Formatting, tests, race tests, repeated tests, and screenshot capture are safe to rerun. Screenshot capture replaces evidence only inside this plan's directory. Labour production is invoked only from successful Raid completion, so repeated ordinary update frames cannot duplicate a payout. If tests expose an accidental way to call the transition twice, guard the successful completion boundary rather than adding a general event or transaction framework.

If visual capture cannot open Ebitengine in the available environment, record the exact failure and retain automated validation; do not change graphics code to bypass an environment-specific failure. Existing user changes must remain untouched, and ownership must remain `dave:dave`.

## Artifacts and Notes

The before-change running-game screenshot is `plans/50-home-plot-grasslands-terrain/screenshots/running-game.png`, which shows the obsolete `Raid in 02:00` label. Post-change evidence is under `plans/51-peaceful-phase-split/screenshots/`; `running-game.png` proves the Management HUD and `woodcutter-tooltip.png` proves the Labour production wording.

## Interfaces and Dependencies

No exported Go API, save format, asset catalog, external dependency, or cross-package interface changes. The private phase type gains Labour and Management values, `gameStatus` loses the unused calm timer, and internal permission methods switch to Management. `StructureTemplate.ResourceYield` remains the source of economic production values.

Revision note (2026-07-14): Created the ExecPlan from the accepted implementation plan after repository inspection and baseline validation. It fixes phase order, Day 1 behavior, action permissions, production timing, documentation ownership, visual evidence, and line-count handling.

Revision note (2026-07-14): Completed the implementation, documented the final state and limitations, recorded visual evidence, and added full validation plus line-count results.
