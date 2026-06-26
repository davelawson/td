# Enforce tower staffing availability

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root and is saved at `plans/32-enforce-tower-staffing.md`.

## Purpose / Big Picture

Tower staffing requirements currently appear in the UI but do not affect construction. After this change, a tower can be built only when every required inhabitant role has enough available people. Successful construction reserves those inhabitants by reducing available counts while leaving total counts unchanged.

New games still start with `0/0` Apprentices, Soldiers, and Peasants, so no tower can initially be built. Recruitment, reassignment, tower removal, and staff release remain outside this slice.

## Progress

- [x] (2026-06-26T00:00:00Z) Inspected population state, tower staffing metadata, affordability highlighting, drag placement, tests, screenshot capture, and control documents.
- [x] (2026-06-26T00:00:00Z) Confirmed with the user that successful construction reserves available staff without reducing totals.
- [x] (2026-06-26T00:00:00Z) Created this ExecPlan.
- [x] (2026-06-26T00:00:00Z) Added staffing availability and reservation rules to construction.
- [x] (2026-06-26T00:00:00Z) Updated focused tests and screenshot fixtures.
- [x] (2026-06-26T00:00:00Z) Updated durable control documents.
- [x] (2026-06-26T00:00:00Z) Ran formatting, tests, screenshot capture, whitespace, ownership, and status checks.
- [x] (2026-06-26T00:00:00Z) Checked hand-written code-file line counts; no file exceeds the 600-line preference.

## Surprises & Discoveries

- Observation: The screenshot harness cannot directly seed private population state from `cmd/td`.
  Evidence: Runtime population fields are intentionally private to `internal/game`, so visual evidence should demonstrate the zero-staff disabled state rather than add a public testing API.

- Observation: Resource affordability and staffing eligibility can diverge at game start.
  Evidence: The initial 100 Wood, 50 Stone, and 20 Metal cover Bow and Flame Bolt costs, but their icons remain unhighlighted because Soldier and Apprentice availability are both zero.

- Observation: Existing placement tests needed explicit population setup to keep testing their original placement constraints.
  Evidence: Focused helpers now seed matching available/total values before road, forest, pause, occupied-tile, and successful-placement scenarios.

## Decision Log

- Decision: Successful tower construction reserves required inhabitants by reducing `available` counts while preserving `total`.
  Rationale: A check without reservation would allow one inhabitant to satisfy unlimited towers, while reducing totals would incorrectly treat staffing as permanent population consumption.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Staffing gates construction but does not gate combat after a tower exists.
  Rationale: This slice has no reassignment, loss, loading reconciliation, or tower-removal workflow, so construction is the only complete and observable boundary.
  Date/Author: 2026-06-26 / Codex

## Outcomes & Retrospective

Implemented staffing-gated construction. Building-bar highlighting, drag start, and release validation now share one eligibility rule requiring both sufficient resources and every required available inhabitant role. Successful placement deducts resources, reserves staff by reducing available counts, preserves totals, and then installs the tower.

Focused tests prove exact staffing permits construction, one missing role blocks Catapult, successful Bow and Catapult builds reserve the correct roles, repeated Bow construction is blocked after its Soldier is committed, release rechecks staffing, and invalid drops preserve both resources and populations.

Validation completed:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

All tests passed. Whitespace and ownership checks produced no output. Screenshot capture produced eight files under `plans/32-enforce-tower-staffing/screenshots/`, including `insufficient-staff.png`, which shows zero populations and the resource-affordable Bow option remaining unhighlighted.

No hand-written Go file exceeds 600 lines. The largest is `internal/game/game_test.go` at 591 lines.

## Context and Orientation

`internal/game/hud.go` owns private population counts with available and total values. `internal/game/structures.go` defines each tower's `StaffingRequirements`. `internal/game/building_bar.go` owns affordability highlighting, drag start, placement validation, resource deduction, and final feature placement. Construction currently checks resources twice: when drag starts and when the icon is released.

The staffing check must follow the same pattern. A tower option highlights and starts dragging only when both resources and staff are sufficient. Release must recheck both because state may change during a drag. Invalid placement must leave resources and populations unchanged.

## Plan of Work

Add focused methods that compare all three staffing requirements against current available populations and reserve all three roles together. Use one construction-eligibility helper in building-bar highlighting, drag start, and release validation so resource and staffing behavior cannot diverge.

On a valid drop, deduct resources and reserve staff immediately before installing the tower feature. Tests will seed explicit populations for successful build scenarios because production starts at zero. Add coverage for exact availability, mixed-role failure, repeated builds, preserved totals, and invalid drops.

Update screenshot capture to replace tower-building fixtures with a cursor-hover target over a resource-affordable but staff-unavailable Bow Tower. Update `README.md`, `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, and `ARCHITECTURE.md` to make the new construction gate and current zero-population limitation explicit.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

    gofmt -w internal/game/*.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

Report files over 600 lines with a recommendation. Do not perform an unplanned split, refactor, or dependency addition without user approval.

## Validation and Acceptance

`go test ./...` must pass. Tests must prove exact available staffing permits a build, one missing Catapult role blocks construction, successful placement reserves staff but preserves totals, a second build is blocked after staffing is exhausted, and invalid drops spend neither resources nor staff. Hover emphasis and drag start must require both resources and staffing.

Screenshot evidence must show the initial `0/0` populations and a non-highlighted Bow Tower while the cursor is over its otherwise resource-affordable icon. Documentation must state that no tower is initially buildable until a future population source exists.

## Idempotence and Recovery

Formatting, tests, and screenshot capture are safe to repeat. Screenshot capture overwrites only this plan's screenshot directory. All source and documentation changes are ordinary version-controlled edits.

## Artifacts and Notes

The primary visual artifact is `plans/32-enforce-tower-staffing/screenshots/insufficient-staff.png`.

## Interfaces and Dependencies

No new dependency or public API is required. Staffing checks and reservations remain private methods on `internal/game.State`. `StructureTemplate.Staffing` remains the single source of tower requirements.

Revision note, 2026-06-26: Updated the living plan after implementation with completed progress, discoveries, validation results, screenshot evidence, and final line-count results.
