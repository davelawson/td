# Add tower staffing requirements and reset the starting Domain

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` in the repository root and is saved at `plans/31-tower-staffing-requirements.md`.

## Purpose / Big Picture

After this change, a new game begins with only the Sanctum, 100 Wood, 50 Stone, 20 Metal, and no inhabitants. The building bar and selected-tower panel show the inhabitants each tower will eventually require: one Soldier for a Bow Tower, one Apprentice for a Flame Bolt Tower, and one Soldier plus two Peasants for a Catapult Tower.

Staffing is informational in this slice. A player can still build and operate a tower while all population values are `0/0`; recruitment, assignment, release, and operational enforcement remain future work.

## Progress

- [x] (2026-06-26T00:14:03Z) Inspected tower templates, starting-map creation, status initialization, building-bar layout, selection panels, screenshot fixtures, tests, and control documents.
- [x] (2026-06-26T00:14:03Z) Confirmed staffing values, informational behavior, and build-bar plus selection-panel presentation with the user.
- [x] (2026-06-26T00:14:03Z) Created this ExecPlan.
- [x] (2026-06-26T00:14:03Z) Added staffing data, the Sanctum-only starting map, and new initial resources.
- [x] (2026-06-26T00:14:03Z) Added staffing presentation and focused tests.
- [x] (2026-06-26T00:14:03Z) Updated screenshot fixtures and control documents.
- [x] (2026-06-26T00:14:03Z) Ran formatting, automated tests, screenshot capture, whitespace, ownership, and status checks.
- [x] (2026-06-26T00:14:03Z) Checked hand-written code-file line counts; no file exceeds the 600-line preference.

## Surprises & Discoveries

- Observation: Combat, selection, placement, and screenshot tests depended directly on the two authored starting towers.
  Evidence: Focused tests called helpers named `removeStartingBowTower` and `removeStartingFlameBoltTower`, and screenshot selection clicked the former Bow Tower location.

- Observation: The revised `100/50/20` starting resources make Flame Bolt exactly affordable as well as Bow.
  Evidence: Initial affordability tests failed until expectations were updated to leave only Catapult unaffordable.

- Observation: The added staffing row fits the existing 96-pixel building bar at 1920x1080 without overlap.
  Evidence: `running-game.png` shows three distinct tower blocks with resource and staffing rows, and `selected-tower.png` shows the expanded panel fully on screen.

## Decision Log

- Decision: Staffing requirements are template metadata and visible UI only.
  Rationale: Populations intentionally start at zero and no recruitment or assignment workflow exists, so enforcing staffing would make towers unusable without a much larger gameplay slice.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Bow Tower requires one Soldier, Flame Bolt Tower requires one Apprentice, and Catapult Tower requires one Soldier plus two Peasants.
  Rationale: The requirements communicate distinct tower roles and prove that templates can contain any combination of the three inhabitant types.
  Date/Author: 2026-06-26 / User and Codex

- Decision: A new game starts with only the Sanctum, 100 Wood, 50 Stone, 20 Metal, and `0/0` for every population.
  Rationale: The player should construct the first defenses instead of receiving free towers, while the revised resources make Bow and Flame Bolt immediately affordable but leave Catapult as a later purchase.
  Date/Author: 2026-06-26 / User

## Outcomes & Retrospective

Implemented tower staffing metadata and presentation. Bow requires one Soldier, Flame Bolt requires one Apprentice, and Catapult requires one Soldier plus two Peasants. The building bar uses the existing population icons to display non-zero requirements, and selected tower details show named requirement rows.

New games now start with only the Sanctum, 100 Wood, 50 Stone, 20 Metal, and `0/0` for all three populations. Tests and screenshot fixtures place towers explicitly when needed. Staffing remains informational, and a placement test proves a tower can still be constructed with zero inhabitants.

Validation completed:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

All tests passed. Whitespace and ownership checks produced no output. Screenshot evidence was captured under `plans/31-tower-staffing-requirements/screenshots/`.

No hand-written Go file exceeds 600 lines. The largest is `internal/game/game_test.go` at 591 lines, reduced from its previous 602 lines by removing obsolete starting-tower assertions.

## Context and Orientation

`td` is a Go/Ebitengine tower-defense prototype. `internal/game/structures.go` owns tower templates and construction costs. `internal/game/map.go` creates the 15x15 starting Plot. `internal/game/hud.go` initializes resources and population status. `internal/game/building_bar.go` renders tower choices and handles placement. `internal/game/selection_panel.go` formats selected-object details. `cmd/td/main_test.go` captures opt-in visual evidence.

The existing top bar already loads Apprentice, Soldier, and Peasant icons and shows each population as available/total. The new staffing rows must reuse those assets. `PRODUCT.md` and `README.md` must describe current behavior, `GAME.md` must record the gameplay decision, `ARCHITECTURE.md` must describe ownership and invariants, and `ROADMAP.md` must continue distinguishing metadata from future assignment mechanics.

## Plan of Work

Add a small `StaffingRequirements` value type to `internal/game/structures.go` and store it on every `StructureTemplate`. Populate the three tower templates with the accepted values and leave the Sanctum at zero requirements.

Remove authored Bow and Flame Bolt features from `NewDefaultHomePlot`. Change prototype resource initialization to 100 Wood, 50 Stone, and 20 Metal while retaining zero available and total population counts.

Extend building-bar items with staffing data. Render one compact row under resource costs using existing population icons and numeric counts in Apprentice, Soldier, Peasant order, omitting zero roles. Increase vertical item spacing so the row does not overlap the next tower. Keep icon hit targets and affordability behavior unchanged.

Append non-zero staffing rows to selected tower panels. Refactor tests that assumed starting towers so each test explicitly places the structure it exercises. Update screenshot capture so selected-tower evidence first builds and then selects a Bow Tower.

Update current-product, game-design, roadmap, onboarding, and architecture documents. Advance screenshot output to this plan directory, capture the standard evidence set, and review the starting map, building bar, placed tower, and selected tower.

## Concrete Steps

Run all commands from `/home/dave/dev/ai/td`.

Format and validate the implementation with:

    gofmt -w internal/game/*.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    git status --short

End with the required line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs -r wc -l | sort -n

Report files over 600 lines with a concrete recommendation. Do not perform an unplanned split, refactor, or library addition without user approval.

## Validation and Acceptance

`go test ./...` must pass. Tests must prove the starting Plot has only the Sanctum, resources start at `100/50/20`, populations remain `0/0`, each template has the exact staffing requirements, mixed-role requirements retain stable ordering, zero roles are omitted from display, and construction remains possible with zero inhabitants.

`plans/31-tower-staffing-requirements/screenshots/running-game.png` must show no starting towers and must show staffing beneath all three building-bar choices. `placed-tower.png` must show a player-built Bow Tower. `selected-tower.png` must show that Bow Tower's required Soldier count.

Documentation must state that staffing is informational and that recruitment, assignment, and enforcement are unimplemented. Whitespace and ownership checks must produce no output.

## Idempotence and Recovery

Go formatting, tests, and screenshot capture are safe to repeat. Screenshot capture overwrites only files under this plan's screenshot directory. Source and documentation edits are ordinary version-controlled changes.

## Artifacts and Notes

The primary evidence files are `plans/31-tower-staffing-requirements/screenshots/running-game.png`, `placed-tower.png`, and `selected-tower.png`.

## Interfaces and Dependencies

No new dependency or public package is required. `internal/game.StructureTemplate` gains a staffing value:

    type StaffingRequirements struct {
        Apprentices int
        Soldiers    int
        Peasants    int
    }

    type StructureTemplate struct {
        ...
        Staffing StaffingRequirements
    }

Staffing remains private gameplay-package data in practical use because `internal/game` cannot be imported outside the module's internal boundary.

Revision note, 2026-06-26: Updated the living plan after implementation with completed progress, discoveries, validation evidence, screenshots, and final line-count results.
