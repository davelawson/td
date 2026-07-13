# Add Explorable Adjacent Plots

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/43-exploring-additional-plots.md`.

## Purpose / Big Picture

After this change, the player can expand the visible Domain during calm play by clicking magnifying-glass buttons on borders between explored and unexplored orthogonal Plots. A Plot is a 15x15 group of Tiles. The home Plot remains centered on the Sanctum at world coordinate `(0, 0)`, and newly explored Plots become usable grassland space. Exploring north extends the visible road and the current deterministic Raid path, so the player can see enemies enter from the farthest explored north Plot and can build defenses across explored territory.

## Progress

- [x] (2026-07-13T20:06Z) Created the ExecPlan from the accepted plan.
- [x] (2026-07-13T20:18Z) Added plot-coordinate map state and world-coordinate helpers.
- [x] (2026-07-13T20:22Z) Added exploration button rendering, hit testing, and calm-phase reveal behavior.
- [x] (2026-07-13T20:27Z) Adapted structure placement, selection, economy, tower combat, and Raid spawning to explored Plots.
- [x] (2026-07-13T20:34Z) Updated `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, `DESIGN.md`, and `ROADMAP.md`.
- [x] (2026-07-13T20:30Z) Added and updated tests for exploration, placement, selection, economy, combat, and north-path behavior.
- [x] (2026-07-13T20:38Z) Ran `go test ./...`, screenshot capture, `git diff --check`, ownership check, and a Go file line-count review.
- [x] (2026-07-13T20:45Z) Removed per-Plot render frames and padding so adjacent explored Plots render as continuous terrain, then refreshed screenshot evidence.
- [x] (2026-07-13T21:05Z) Removed the home Plot's generated tree perimeter so its non-road Tiles are continuous buildable grassland.
- [x] (2026-07-13T21:09Z) Revalidated tests, screenshots, whitespace, ownership, and Go file line counts after removing the tree perimeter.

## Surprises & Discoveries

- Observation: Baseline tests passed before implementation.
  Evidence: `go test ./...` reported all packages passing before edits.
- Observation: `internal/game/game_test.go` was already 603 lines before this work.
  Evidence: pre-change line-count review showed `603 internal/game/game_test.go`.
- Observation: The screenshot harness can capture the new exploration frame in this environment.
  Evidence: `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1` passed and wrote `plans/43-exploring-additional-plots/screenshots/explored-north.png`.

## Decision Log

- Decision: Use an unbounded map keyed by orthogonal plot coordinates, with `(0,0)` as home and positive plot Y pointing north.
  Rationale: This supports unlimited grid expansion without introducing campaign-map boundaries.
  Date/Author: 2026-07-13 / Codex
- Decision: Collapse scouted and claimed into one explored-and-usable state for this slice.
  Rationale: The requested behavior says newly explored Plots behave as if they were part of the initial home Plot.
  Date/Author: 2026-07-13 / Codex
- Decision: Keep exploration free and available during paused calm play.
  Rationale: This matches the accepted plan and current calm building behavior.
  Date/Author: 2026-07-13 / Codex
- Decision: North-chain plots show straight center roads and extend Raid spawn distance.
  Rationale: A visible road keeps the extended enemy path readable.
  Date/Author: 2026-07-13 / Codex
- Decision: Do not draw per-Plot frames, gutters, or padding around explored Plots.
  Rationale: Once adjacent Plots are both explored, the Domain should read as one continuous field rather than separate boards.
  Date/Author: 2026-07-13 / Codex
- Decision: Generate the home Plot as open grassland without a perimeter tree layer.
  Rationale: The starting Plot should be fully usable and visually continuous with neighboring explored grassland. Forest remains a supported terrain type for future authored map content.
  Date/Author: 2026-07-13 / Codex

## Outcomes & Retrospective

Implementation completed the first adjacent-Plot exploration slice. New games start with an open-grassland home Plot explored and draw magnifying-glass buttons on unexplored orthogonal borders. Calm-phase clicks reveal connected Plots for free, paused calm play permits exploration, newly explored Plots are empty grassland, adjacent explored Plots have no frames, padding, or tree terrain ring between them, and central north exploration extends both the visible road and the deterministic Raid spawn point. Structure placement, selection, selected-structure panels, economic payouts, and tower combat now work across explored Plots.

Validation passed again after removing the tree perimeter on 2026-07-13: `go test ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and `find . -xdev ! -user dave -printf '%u:%g %p\n'` all succeeded. Refreshed screenshot evidence under `plans/43-exploring-additional-plots/screenshots/` shows an open-grassland home Plot and a seamless north expansion without perimeter trees.

The final hand-written Go file line-count review found one file over the 600-line preference: `internal/game/game_test.go` at 602 lines. This file was already 603 lines before the feature; replacing the old tree-border test with the grassland invariant reduced it by one line. Recommended follow-up remains splitting camera/map setup tests out of `internal/game/game_test.go` during a separate approved test-organization change.

## Context and Orientation

The current game code lives in `internal/game`. Before this plan, `Map` in `internal/game/map.go` stored only `Home Plot`, and many systems accessed `s.gameMap.Home.Tiles[y][x]` directly. Rendering in `internal/game/scene.go`, placement in `internal/game/building_bar.go`, selection in `internal/game/selection.go`, economy in `internal/game/resources.go`, combat in `internal/game/combat.go`, and Raid spawning in `internal/game/raid.go` all assumed the single home Plot. `GAME.md` defines Plots as orthogonally adjacent 15x15 groups of Tiles and says early prototypes may collapse scouting and claiming. `PRODUCT.md` currently says exploration is missing; that must change. `ARCHITECTURE.md` currently describes home-Plot-only projection and rendering; that must change. `DESIGN.md` requires map affordances to remain readable and plans that change rendered output to capture visual evidence.

## Plan of Work

First, change map data from a single `Home` field into a plot map while preserving a compatibility helper for tests and simple home access. Define `plotCoordinate` and update `tileCoordinate` so it includes `Plot plotCoordinate`, `X`, and `Y`. Add helpers for home tile coordinates, plot lookup, explored-plot iteration, tile world rectangles, screen hit testing, and revealing adjacent plots.

Next, add exploration state and input. The game should draw one magnifying-glass button on each border where an explored Plot has an unexplored north, south, east, or west neighbor. The button should be screen-space hit tested after projection but represent a world-space border location. A click on the button should reveal that neighbor only when the game is in calm phase, not breached, and not blocked by overlay or other UI. Paused calm play should allow exploration.

Then, adapt existing gameplay systems. Rendering should draw all explored plots. Placement should accept drops on any explored empty grass tile. Selection should find structures in any explored plot, store the plot-aware coordinate, and keep selection panels working. Economic payouts should iterate all explored plots. Combat towers and cooldown keys should use plot-aware tile coordinates. Raid spawn should use the northernmost explored plot on the central north chain, spawning at the center road edge of that plot, then continue moving straight south on `X=0`.

Finally, update control documents and tests. Tests should prove map reveal behavior, UI click behavior, north road generation, placement/selection outside home, post-Raid economic payouts outside home, tower firing outside home, and extended Raid spawn position. Validation must include `go test ./...`, `git diff --check`, manual visual review with screenshots under `plans/43-exploring-additional-plots/screenshots/`, and a line-count review.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go`, `internal/game/coordinates.go`, and `internal/game/combat.go` to introduce plot-aware coordinates and helpers.
2. Edit `internal/game/scene.go` and `internal/game/colors.go` to draw all explored Plots and the explore buttons.
3. Edit `internal/game/game.go`, `internal/game/building_bar.go`, `internal/game/selection.go`, `internal/game/resources.go`, `internal/game/raid.go`, and selection-panel code to use plot-aware tile addresses.
4. Update tests under `internal/game`.
5. Update `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, and `DESIGN.md`.
6. Validate:

    go test ./...
    git diff --check
    find internal cmd assets -path '*/vendor/*' -prune -o -type f \( -name '*.go' \) -print | xargs wc -l | sort -n | tail -n 20

If any hand-written Go file exceeds 600 lines, record it in `Outcomes & Retrospective` with a concrete recommendation. Do not do an unplanned split unless the user approves it.

## Validation and Acceptance

Acceptance is met when a new game starts with only an open-grassland home Plot containing the Sanctum and north road, with no perimeter trees. Explore buttons must appear on the four borders, clicking one during calm play must reveal the adjacent Plot, and newly revealed Plots must accept the same building placement and structure selection behavior as the home Plot. Exploring north must make the next Raid spawn farther north on the visible road. `go test ./...` must pass. `PRODUCT.md`, `GAME.md`, `ARCHITECTURE.md`, and `DESIGN.md` must describe the implemented behavior.

## Idempotence and Recovery

The work is additive and test-driven. If a step fails, rerun `go test ./...` to identify the package and restore behavior by reverting only the affected local edits, not unrelated user changes. Exploration reveal should be idempotent: revealing an already explored Plot should not replace its Tiles or structures.

## Artifacts and Notes

Visual evidence should be saved under `plans/43-exploring-additional-plots/screenshots/` after implementation. Screenshots should include the initial open-grassland home Plot with explore buttons and no perimeter trees, plus at least one explored north Plot with a visible continued road and no plot-level terrain boundary.

## Interfaces and Dependencies

No new external dependencies are planned. Use existing Ebitengine drawing helpers and `internal/ui` button patterns. New interfaces should remain inside `internal/game`.
