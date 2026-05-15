# Add Basic Raids

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This file is saved at `plans/14-basic-raids.md` and must be maintained in accordance with `PLANS.md`.

## Purpose / Big Picture

After this change, a player can start the first visible tower-defense pressure slice from the game screen. The bottom-left `Next Raid` button starts a deterministic Raid immediately. Placeholder enemies spawn one at a time from the north road edge, move along the existing straight road toward the Sanctum, and are removed when they reach it. The Raid continues until all pending and active enemies are gone, or until the Sanctum is breached after Barricade charges are exhausted.

## Progress

- [x] (2026-05-15T17:36Z) Created this ExecPlan from the accepted plan.
- [x] (2026-05-15T17:45Z) Added private Raid state, deterministic spawning, movement, completion, Barricade, and breach rules under `internal/game`.
- [x] (2026-05-15T17:45Z) Added the bottom-left `Next Raid` screen-space button and rendered visible placeholder enemies.
- [x] (2026-05-15T17:45Z) Added Go tests for Raid button behavior, staggered spawning, movement, completion, Barricade spending, breach, and pause/menu blocking.
- [x] (2026-05-15T17:50Z) Updated `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `README.md` to match the implemented current behavior and design decisions.
- [x] (2026-05-15T17:52Z) Ran `go test ./...`, `git diff --check`, and a hand-written Go file line-count review.
- [x] (2026-05-15T17:56Z) Captured visual evidence under `plans/14-basic-raids/screenshots/` for calm, active Raid, and breached states.

The final line-count review must report any hand-written code file over the 600-line preference from `CODESTYLE.md`. Unplanned file splits, refactors, or new libraries require user approval before implementation unless this plan is revised to include them.

## Surprises & Discoveries

- Observation: Ebitengine `Image.ReadPixels` cannot be called before the game starts, so offscreen screenshot capture had to run through a temporary `ebiten.RunGame` helper and read pixels during Draw.
  Evidence: The first helper failed with `panic: ui: ReadPixels cannot be called before the game starts`; the replacement generated three 1920x1080 PNGs.

## Decision Log

- Decision: Implement the first Raid slice as deterministic scripted Raids with placeholder enemies, a fixed north-road path, staggered spawning, and no towers or combat.
  Rationale: This creates visible, testable Raid lifecycle behavior without expanding scope into tower targeting, damage, rewards, pathfinding, or art.
  Date/Author: 2026-05-15 / Codex

- Decision: Disable the `Next Raid` button during an active Raid instead of hiding it or queueing another Raid.
  Rationale: Disabled-but-visible UI makes the current state explicit and avoids queueing rules before encounter pacing is designed.
  Date/Author: 2026-05-15 / Codex

- Decision: Mark the Sanctum breached when an enemy reaches it with zero Barricade charges, clear the active Raid, and prevent further Raid starts.
  Rationale: This records a concrete failure state while avoiding new app-level loss routing or recovery UI.
  Date/Author: 2026-05-15 / Codex

## Outcomes & Retrospective

- Implemented the first deterministic placeholder Raid slice. The game now has a bottom-left `Next Raid` button, staggered placeholder enemy spawning on the fixed north road, enemy movement toward the Sanctum, real enemy-remaining HUD text, Barricade spending, and a terminal breached state that prevents further Raid starts.
- Validation passed with `go test ./...` and `git diff --check`.
- Visual evidence was captured at `plans/14-basic-raids/screenshots/calm-next-raid.png`, `plans/14-basic-raids/screenshots/active-raid.png`, and `plans/14-basic-raids/screenshots/sanctum-breached.png`.
- Final hand-written Go line-count review found no file over 600 lines. The largest file is `internal/game/game_test.go` at 537 lines, below the preference but close enough that future game-state tests should continue going into focused files such as `internal/game/raid_test.go`.

## Context and Orientation

The project is a Go/Ebitengine local PC prototype. `cmd/td/main.go` owns Ebitengine startup and forwards input into either `internal/menu` or `internal/game`. `internal/game.State` owns the running game state, including pause, camera, top-bar HUD, static map rendering, and the ESC overlay menu. The current home Plot is a 15x15 grid with a straight road from the north-center edge to the centered Sanctum. This plan uses that existing road as the only Raid path.

`GAME.md` defines a `Raid` as an assault on the wizard's Domain and says enemies attempt to reach the Sanctum. It also says the Sanctum's `Barricade` charges atomize breaching enemies until exhausted. `PRODUCT.md` currently says tower-defense combat is missing and must be updated when this slice exists. `ARCHITECTURE.md` currently describes `internal/game` as owning prototype game state and should be updated to include first Raid state and rendering. `ROADMAP.md` says an early deterministic defense loop is a near-term priority and should be revised once this baseline exists.

## Plan of Work

First, add private Raid model code inside `internal/game`. Add fields on `State` for Raid state and any screen-space Raid button UI state. Keep all new gameplay state private to the package. A Raid must track whether it is active, whether the Sanctum has been breached, which Raid number is current, how many enemies remain pending, how many updates until the next spawn, and active enemies with path progress. The first Raid has five enemies, and each later Raid adds two. Enemies spawn every 45 logical updates. They move down the existing north road at a fixed speed in world pixels per logical update.

Second, add update logic so normal unpaused updates advance Raid spawning and enemy movement. A click inside `Next Raid` starts the next Raid only when no Raid is active and the Sanctum is not breached. While the ESC overlay menu is open, the existing menu path must continue to block camera and Raid input. While SPACE pause is active, camera inspection still works but Raid spawning and movement do not advance.

Third, add rendering. Draw the `Next Raid` button in screen space at the bottom-left of the game screen when the ESC overlay menu is closed. It is disabled during an active Raid or after breach. Draw placeholder enemies as simple high-contrast shapes projected through the same camera as the map, so they stay attached to the road. Update the top bar so calm state shows the existing countdown-style text and Raid state shows real remaining enemies. In breached state, show a concise breached status.

Fourth, add tests. Prefer `internal/game/raid_test.go` so `internal/game/game_test.go` does not grow further toward the 600-line preference. Exercise pure state updates and hit testing without opening an Ebitengine window. Add helper functions locally in the new test file if needed.

Fifth, update durable documents. `PRODUCT.md` must describe the new button, visible enemies, Raid lifecycle, Barricade spending, and current limits. `GAME.md` must move the first concrete Raid spawn/path/breach decisions out of open-ended intent and into current decisions. `ROADMAP.md` must reflect that the deterministic Raid baseline exists and that towers/combat remain future work. `ARCHITECTURE.md` must state that `internal/game` owns first Raid state and placeholder enemy rendering.

## Concrete Steps

From `/home/dave/dev/ai/td`, edit the Go files using small, scoped changes. Run `gofmt` on changed Go files. Run:

    go test ./...

Expect all packages to pass. Then run:

    git diff --check

Expect no whitespace errors. For line counts, run:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Report any hand-written Go file over 600 lines, or close enough that the next likely change will cross the preference. Do not perform an unplanned split without approval.

Observed validation:

    go test ./...
    ok  	td/assets	(cached)
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    git diff --check
    # no output

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n
    ...
      537 internal/game/game_test.go
     3008 total

For visual evidence, start the app:

    go run ./cmd/td

Enter a Wizard name, start a game, and capture screenshots under `plans/14-basic-raids/screenshots/` showing calm state with `Next Raid`, an active Raid with disabled button and enemies on the road, and either calm state after completion or breached state.

## Validation and Acceptance

Acceptance is met when a player can start a game, click `Next Raid`, see placeholder enemies spawn from the north road edge in a staggered manner, see them move toward the Sanctum, and see the Raid end only after pending and active enemies are gone or the Sanctum is breached. The `Next Raid` button must be disabled while a Raid is active and after breach. The HUD must show real remaining enemy count during a Raid. Barricade must decrement when enemies reach the Sanctum while charges remain. At zero Barricade, the Sanctum must enter a breached state and prevent further Raids.

Tests must pass with `go test ./...`. Documentation acceptance requires `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, and `ARCHITECTURE.md` to describe the implemented behavior without implying towers, combat damage, rewards, resource changes, pathfinding, or loss routing exist.

## Idempotence and Recovery

All steps are additive or replacement edits to tracked files. If a test fails, inspect the failing test, fix the local behavior, and rerun `go test ./...`. If manual screenshot capture is not possible in the current environment, record that limitation in this plan's `Outcomes & Retrospective` and still provide automated test evidence.

## Artifacts and Notes

- `plans/14-basic-raids/screenshots/calm-next-raid.png` shows calm state with the enabled bottom-left `Next Raid` button.
- `plans/14-basic-raids/screenshots/active-raid.png` shows an active Raid with visible placeholder enemies on the north road, the disabled `Next Raid` button, and real remaining enemy count in the top bar.
- `plans/14-basic-raids/screenshots/sanctum-breached.png` shows the breached HUD state with Barricade at zero and no active enemies.

## Interfaces and Dependencies

Use the existing Ebitengine dependency only. Do not add external libraries. Keep app-level `game.Input` unchanged because it already carries cursor and click data. Keep Raid concepts private to `internal/game` unless a future feature needs public access.

At completion, `internal/game.State` must support:

- starting a Raid from a bottom-left screen-space button,
- updating staggered spawns and active enemy movement during unpaused logical updates,
- computing real enemies remaining for the top bar,
- drawing placeholder enemies in world space,
- disabling new Raids during active Raids and after breach.
