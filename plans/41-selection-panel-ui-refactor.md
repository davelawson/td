# Selection Panel UI Refactor

This ExecPlan is a living document. It follows `PLANS.md`.

## Purpose / Big Picture

This refactor moves the selected-object detail panel's UI responsibilities out of `internal/game` and into `internal/ui`. After this change, `internal/game` still owns selection state and selected-object lookup, but the generic UI package owns panel labels, formatting, bounds, hit testing, and drawing. The player-visible panel behavior should remain unchanged.

## Progress

- [x] (2026-06-28) Inspected the existing selection panel, selection input blocking, camera/build UI blocking, screenshot harness, and package docs.
- [x] (2026-06-28) Move generic selection-panel presentation into `internal/ui`.
- [x] (2026-06-28) Replace game selection-panel row formatting with a selected-object data adapter.
- [x] (2026-06-28) Update focused game/UI tests and screenshot evidence path.
- [x] (2026-06-28) Update package-boundary documentation.
- [x] (2026-06-28) Run validation and record outcomes.
- [x] (2026-06-28) Check hand-written code file line counts and record any files over the 600-line preference.

## Surprises & Discoveries

- Observation: The current screenshot harness still writes to `plans/40-right-drag-camera/screenshots/`.
  Evidence: `cmd/td/main_test.go` sets `basePath` to that plan directory.

- Observation: The UI bounds test initially used a stale expected height for a five-row tower panel.
  Evidence: `go test ./...` failed with `bounds = {X:1488 Y:738 W:390 H:300}`, and the row-count formula confirmed the implementation preserved the existing panel height.

## Decision Log

- Decision: Use a strict boundary for this slice.
  Rationale: `game` should supply raw selected-object facts, while `ui` should own selection-panel labels, formatting, bounds, hit testing, and rendering.
  Date/Author: 2026-06-28 / Codex, based on user choice.

## Plan of Work

Add a generic selected-object panel model and renderer in `internal/ui`. The UI data should contain facts like kind, name, health, resources, population counts, range, damage, and timing. It should not import or reference `internal/game`.

Change `internal/game/selection_panel.go` into a thin adapter from `State` to `ui.SelectionPanelData`. Keep `selectionPanelBounds`, `selectionPanelContains`, and `drawSelectionPanel` as game methods only because callers already use them, but make them delegate to `ui`.

Move formatting tests to `internal/ui`, keep semantic selected-object tests in `internal/game`, update screenshots to this plan directory, and update docs that describe package ownership.

## Concrete Steps

Run:

    gofmt -w internal/game/selection_panel.go internal/game/selection_panel_test.go internal/game/resources.go internal/ui/*.go cmd/td/main_test.go
    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    git status --short
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

Record final validation and line-count findings here before finishing.

## Validation and Acceptance

Automated tests should prove the game adapter returns the same selected-object facts and that UI row generation formats the same labels and values as before. Manual or screenshot acceptance is that selected raider and selected structure panels still appear in the bottom-right, remain readable, and still block map interaction beneath them.

## Outcomes & Retrospective

Implemented the selection panel package-boundary refactor. `internal/game` now adapts selected raiders and structures into `ui.SelectionPanelData`; `internal/ui` owns the panel data types, labels, formatting, bounds, hit testing, and drawing. `internal/game/resources.go` no longer contains selection-panel formatting helpers. Package-boundary docs now describe the split.

Validation passed on 2026-06-28:

    go test ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'

The ownership check printed nothing. Screenshot evidence was written under the plan directory, including `selected-sanctum.png` and `selected-raider.png`:

    plans/41-selection-panel-ui-refactor/screenshots/

Final line-count review:

    rg --files cmd internal assets | grep -E '\.go$' | xargs wc -l | sort -n

No hand-written Go file exceeds the 600-line preference. The largest files remain `internal/game/game_test.go` at 597 lines and `internal/game/building_bar_test.go` at 572 lines, so future changes touching those files should avoid adding bulk there.
