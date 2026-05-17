# Add Enemy Health Bars

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root.

Save this file at `plans/21-enemy-health-bars.md` in the repository root. It uses the next available two-digit prefix after `plans/20-bow-tower-projectile-combat.md`.

## Purpose / Big Picture

After this change, active Raid enemies show visible health bars above their heads. A full-health enemy has a green bar as wide as its sprite. As towers damage the enemy, the filled bar shrinks in proportion to remaining health and shifts toward red, making combat damage readable without opening any debug view.

## Progress

- [x] (2026-05-17T15:37:59Z) Confirmed current enemy rendering lives in `internal/game/raidui.go` and active enemies already store current health plus a template with max health.
- [x] (2026-05-17T15:37:59Z) Created this ExecPlan before code edits.
- [x] (2026-05-17T15:41:01Z) Added pure health fraction and color helpers with focused tests.
- [x] (2026-05-17T15:41:01Z) Rendered health bars above sprite-backed and fallback Raid enemies.
- [x] (2026-05-17T15:41:01Z) Updated `README.md` and `PRODUCT.md` current-state documentation.
- [x] (2026-05-17T15:41:01Z) Captured screenshot evidence under `plans/21-enemy-health-bars/screenshots/`.
- [x] (2026-05-17T15:41:01Z) Ran `go test ./...`, screenshot capture, `git diff --check`, and `git status --short`.
- [x] (2026-05-17T15:41:01Z) Checked hand-written Go code-file line counts; no files exceed the 600-line preference from `CODESTYLE.md`.

## Surprises & Discoveries

- Observation: Health bars could be implemented without changing Raid simulation or enemy data.
  Evidence: `raidEnemy` already stores `health`, and its template already exposes `MaxHealth`.

- Observation: `internal/game/game_test.go` remains close to the 600-line preference, but this change did not increase it.
  Evidence: Final line-count review reports `588 ./internal/game/game_test.go`, with no Go file over 600 lines.

## Decision Log

- Decision: Draw health bars in screen space using the projected enemy rectangle as the anchor.
  Rationale: The user asked for bar length to match sprite width. The projected rectangle already represents the current screen-space sprite width after camera zoom.
  Date/Author: 2026-05-17 / Codex

- Decision: Clamp health fractions to `[0, 1]` and treat enemies without a valid template max health as full health.
  Rationale: This prevents divide-by-zero and keeps fallback or malformed prototype enemies readable instead of hiding or corrupting the bar.
  Date/Author: 2026-05-17 / Codex

- Decision: Use a pure green-to-red RGB interpolation, with no blue component.
  Rationale: The requested behavior says full is green and zero is completely red, so a direct red/green blend is the simplest faithful mapping and easy to test.
  Date/Author: 2026-05-17 / Codex

## Outcomes & Retrospective

Completed. Active Raid enemies now draw a health bar above the enemy marker or sprite. A full-health enemy shows a full-width green bar. Damaged enemies show a shorter filled bar, with the fill color shifting toward red as health approaches zero. The rendering uses the projected enemy rectangle, so the full bar width matches the currently drawn sprite width.

Validation passed: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `git status --short` all completed. Screenshot evidence was written under `plans/21-enemy-health-bars/screenshots/`.

The hand-written Go line-count review found no files over 600 lines. The largest file is `internal/game/game_test.go` at 588 lines. It was already near the preference and was not changed by this plan; future broad game tests should continue going into focused files rather than growing it.

## Context and Orientation

The repository is a Go/Ebitengine prototype for a medieval wizardry tower-defense game. Runtime game state lives in `internal/game/`. Active Raid simulation lives in `internal/game/raid.go`, where `raidEnemy` stores an `health int` and an optional `template *EnemyTemplate`. Enemy templates live in `internal/game/enemies.go`, where the current skeleton has `MaxHealth: 20`. Enemy rendering lives in `internal/game/raidui.go`; `drawRaidEnemy` projects a world-space enemy position into a screen-space rectangle and draws either the skeleton sprite or a fallback circle.

World positions use Tile units. Screen-space drawing uses Ebitengine and `github.com/hajimehoshi/ebiten/v2/vector` for simple filled rectangles and circles. The health bar should be a rendered UI cue attached to each enemy, not a new gameplay rule.

`PRODUCT.md` and `README.md` describe current user-visible behavior, so they must mention health bars after the feature is implemented. `DESIGN.md` already says rendered game output should be readable and visual evidence should be captured under the active plan directory. This change does not alter durable art guidance, source conventions, architecture boundaries, roadmap priorities, or intended game-design rules, so `ART.md`, `CODESTYLE.md`, `ARCHITECTURE.md`, `ROADMAP.md`, and `GAME.md` do not need updates unless implementation reveals a broader decision.

## Plan of Work

Update `internal/game/raidui.go`. Import `image/color` for explicit bar colors. Add small constants for health-bar height and gap. Add helper functions that compute an enemy health fraction from current and max health, clamp that fraction between zero and one, and return the green-to-red fill color. Keep these helpers private to the `game` package and document them with Go doc comments.

Change `drawRaidEnemy` so it computes the projected enemy rectangle for both sprite-backed and fallback enemies, draws the enemy as it does today, and then calls a new `drawRaidEnemyHealthBar` helper. The helper should draw a subtle dark backing rectangle above `rect.y`, then draw a filled foreground rectangle whose width is `rect.w * fraction`. If `rect.w` is not positive, skip the bar. The full backing width must remain the sprite width so the player can see both current health and lost health.

Add `internal/game/raidui_test.go` with focused tests for health fraction and color interpolation. Tests should verify full, half, zero, negative, over-max, and invalid-template cases without opening an Ebitengine window.

Update `cmd/td/main_test.go` so screenshot capture writes to `plans/21-enemy-health-bars/screenshots/`. Keep the existing active-Raid capture target.

Update `README.md` and `PRODUCT.md` so current behavior says active skeleton enemies display health bars above their sprites. Do not describe future features or unimplemented enemy inspection panels.

## Concrete Steps

Run these commands from the repository root.

1. Edit `internal/game/raidui.go` and add `internal/game/raidui_test.go` as described above.

2. Update current-state docs:

       README.md
       PRODUCT.md

3. Update screenshot capture output in `cmd/td/main_test.go`, then capture visual evidence:

       TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1

4. Run validation:

       go test ./...
       git diff --check
       git status --short

5. Check hand-written Go code-file line counts:

       find . -path './.git' -prune -o -path './vendor' -prune -o -path './plans' -prune -o -name '*.go' -print | xargs wc -l

   Report any hand-written code file over 600 lines with a concrete recommendation. Do not perform unplanned refactors, code splits, or library additions without user approval.

## Validation and Acceptance

The change is accepted when `go test ./...` passes, screenshot capture writes `active-raid.png` under `plans/21-enemy-health-bars/screenshots/`, `git diff --check` reports no whitespace errors, and the line-count review is recorded.

Behavioral acceptance is that during an active Raid, each visible skeleton has a health bar above its head. At full health, the bar is green and as wide as the enemy sprite. After damage, the filled bar is shorter in proportion to remaining health and its color shifts toward red. At zero health, the helper color is pure red, though defeated enemies are normally removed from the active Raid immediately.

## Idempotence and Recovery

The code edits are deterministic and can be reapplied safely. Screenshot capture overwrites only files under `plans/21-enemy-health-bars/screenshots/`. If tests fail around fraction rounding, inspect the helper clamp and color interpolation before changing rendering. If the bar appears in the wrong place, inspect the projected enemy rectangle in `drawRaidEnemy` before changing world-coordinate math.

## Artifacts and Notes

Expected screenshot evidence:

    plans/21-enemy-health-bars/screenshots/main-menu.png
    plans/21-enemy-health-bars/screenshots/new-game-configuration.png
    plans/21-enemy-health-bars/screenshots/running-game.png
    plans/21-enemy-health-bars/screenshots/active-raid.png
    plans/21-enemy-health-bars/screenshots/paused-game.png
    plans/21-enemy-health-bars/screenshots/ingame-menu.png

Final validation commands:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    git status --short

Final line-count review:

    588 ./internal/game/game_test.go
    4253 total

No hand-written Go file exceeds the 600-line preference.

## Interfaces and Dependencies

No new Go module dependency is required. Use the existing Ebitengine `vector` package for health-bar rectangles.

At completion, `internal/game/raidui.go` should provide private helpers equivalent to:

    func raidEnemyHealthFraction(enemy raidEnemy) float64
    func raidEnemyHealthBarColor(fraction float64) color.RGBA
    func (s *State) drawRaidEnemyHealthBar(screen *ebiten.Image, rect projectedRect, enemy raidEnemy)
