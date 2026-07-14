# Preassign and Display Frontier Biomes

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan follows `PLANS.md` from the repository root. Save this file as `plans/49-frontier-biome-labels.md`.

## Purpose / Big Picture

After this change, the player can see whether each available exploration destination is Grasslands or Hills before choosing it. A Plot receives its permanent biome when it first becomes orthogonally adjacent to explored land, while its Tiles, terrain, roads, and features remain hidden and ungenerated until exploration. The biome name appears beside the existing magnifying-glass control, turning the current random reveal into an informed spatial choice without adding a scouting screen or biome-selection UI.

## Progress

- [x] (2026-07-14T02:41Z) Created this ExecPlan from the accepted follow-on plan.
- [x] (2026-07-14T02:43Z) Stored stable biome assignments for the current unexplored frontier and used them during reveal.
- [x] (2026-07-14T02:44Z) Rendered directional biome labels while preserving the circular button as the only hit target.
- [x] (2026-07-14T02:45Z) Added deterministic state, generation, geometry, and interaction tests.
- [x] (2026-07-14T02:46Z) Updated durable control documents and captured initial and expanded visual evidence.
- [x] (2026-07-14T02:48Z) Ran formatting, tests, race tests, screenshot capture, whitespace, ownership, and line-count validation.

## Surprises & Discoveries

- Observation: This work begins on top of the uncommitted but validated hills biome implementation from `plans/48-hills-biome.md`.
  Evidence: `git status --short` lists the hills code, documentation, plan, and screenshot evidence as pending changes.
- Observation: `internal/game/game_test.go` already exceeds the 600-line preference.
  Evidence: the prior plan's final review reported 602 lines; this feature does not need to edit that file.
- Observation: Deterministic frontier-expansion tests need to control rolls made after reveal without storing a function inside `Map`.
  Evidence: the first focused test exposed production randomness for newly adjacent coordinates; `revealPlotWithBiomeSource` now provides the test seam while `Map` stores only gameplay data.
- Observation: Fixed-size labels remain readable when the map is zoomed out and make biome differences visible before terrain generation.
  Evidence: `plans/49-frontier-biome-labels/screenshots/running-game.png` shows all four initial labels, and `explored-biomes.png` shows labels around the expanded frontier.

## Decision Log

- Decision: Assign a biome when a Plot first joins the visible exploration frontier.
  Rationale: The user chose frontier-time assignment over coordinate derivation for every possible distant Plot. This fixes the choice before exploration without introducing a world seed or infinite precomputed map.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Store frontier biomes separately from generated explored Plots.
  Rationale: Unexplored Plots need only a biome preview; generating their Tiles early would reveal or allocate state beyond the requested slice. Once explored, `Plot.Biome` remains authoritative.
  Date/Author: 2026-07-14 / Codex
- Decision: Keep only the circular magnifying-glass icon clickable.
  Rationale: The user selected an informational label and the current stable map-space hit target.
  Date/Author: 2026-07-14 / User and Codex
- Decision: Place labels outward into unexplored space.
  Rationale: North labels above, east labels right, south labels below, and west labels left preserve the readable joins and avoid covering explored terrain.
  Date/Author: 2026-07-14 / Codex

## Outcomes & Retrospective

Implementation completed preassigned frontier biomes and explore-button labels. A new map assigns grasslands or hills to its four initial unexplored neighbors in north, east, south, and west order. Revealing a Plot consumes its stored assignment, generates terrain for that biome, removes the coordinate from frontier storage, and assigns only previously unseen neighbors. Shared frontier coordinates retain their first assignment, repeated reveals remain idempotent, and direct reveal of an unassigned coordinate now does nothing.

Every explore button carries the stored biome and renders `Grasslands` or `Hills` in compact gold text outward into unexplored space. North labels appear above, east labels to the right, south labels below, and west labels to the left. Label geometry is intentionally separate from the existing projected circular hitbox, and an interaction test proves clicking the text alone does not reveal terrain.

Durable product, game-design, roadmap, architecture, and visual-direction documents now describe frontier-time assignment, previews, and button-only interaction. Visual evidence under `plans/49-frontier-biome-labels/screenshots/` confirms both the initial labeled choices and the labeled expanded frontier remain readable alongside roads, terrain, and screen-space UI.

Validation passed on 2026-07-14: `go test ./...`, `go test -race ./...`, `TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1`, `git diff --check`, and the ownership check all succeeded. The final hand-written Go line-count review found `internal/game/game_test.go` at 602 lines, unchanged by this feature, and `internal/game/building_bar_test.go` near the preference at 596 lines. A separate approved maintenance change should extract default-state and map tests from `game_test.go` and consider splitting building-bar tests by behavior; no unrelated refactor was performed here.

## Context and Orientation

`internal/game/map.go` stores generated explored Plots in `Map.Plots` and currently chooses a random biome inside `Map.revealPlot`. `internal/game/exploration.go` derives visible explore buttons from explored Plot borders, while `internal/game/scene.go` projects and draws the circular buttons. The hills implementation in `internal/game/plot_generator.go` already knows how to generate either biome with the accepted terrain weights: grasslands uses 6% Tree and 3% Boulder, hills uses 3% Tree and 6% Boulder, and both are 91% empty grass.

The change must preserve the existing home Plot, road and shared-edge overrides, building rejection for Tree and Boulder, free calm-phase reveal, pause behavior, camera behavior, and Raid path extension. `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md` describe the current exploration and biome behavior and must be updated in the same change.

## Plan of Work

Add private frontier-biome storage to `Map`. `NewDefaultMap` initializes the home Plot, then assigns biomes to its north, east, south, and west unexplored neighbors using one 0-99 roll each. Grasslands remains rolls 0-49 and hills remains 50-99. A private source-injection constructor lets tests control those rolls. Assignments are made only when absent, so a coordinate that becomes adjacent from multiple explored Plots retains its first value.

Change reveal so it requires a stored frontier biome, generates the requested biome rather than rolling a new one, moves biome authority into the generated `Plot`, and then assigns biomes to newly exposed frontier coordinates. A reveal without a frontier assignment does nothing. Terrain is still generated only at reveal time, after which north roads and shared-edge cleanup override obstacles as before.

Add the assigned biome to each `exploreButton`. Format the values as `Grasslands` and `Hills`. In `internal/game/scene.go`, measure the compact bold cost font and place a fixed-screen-size label with an eight-pixel gap outside the circular button: above north, right of east, below south, and left of west. Draw labels in the existing explore-control colour. Do not add their rectangles to explore, selection, camera-blocking, or building hit tests.

Add focused tests for deterministic initial assignments, assignment stability, new-frontier expansion, missing-assignment reveal rejection, button biome data, reveal matching the preview, and idempotence. Add geometry tests for all four label directions and an interaction test proving a label click does not reveal the Plot. Existing exploration, placement, road, edge, camera, and Raid tests remain part of regression validation.

Update the five relevant root control documents. Change screenshot capture to write under this plan and preserve evidence of both the initial labeled frontier and the expanded labeled frontier after exploration.

## Concrete Steps

Run commands from `/home/dave/dev/ai/td`.

1. Edit `internal/game/map.go` and `internal/game/plot_generator.go` to add frontier assignment and biome-specific generation.
2. Edit `internal/game/exploration.go` and `internal/game/scene.go` to carry and render labels without changing the icon hitbox.
3. Add focused tests under `internal/game`, then run `gofmt -w internal/game/*.go` and `go test ./...`.
4. Update `PRODUCT.md`, `GAME.md`, `ROADMAP.md`, `ARCHITECTURE.md`, and `DESIGN.md`.
5. Update `cmd/td/main_test.go` to store evidence under `plans/49-frontier-biome-labels/screenshots/` and run:

    gofmt -w internal/game/*.go cmd/td/*.go
    go test ./...
    go test -race ./...
    TD_CAPTURE_SCREENSHOT=1 go test ./cmd/td -run TestCaptureMainMenuScreenshot -count=1
    git diff --check
    find . -xdev ! -user dave -printf '%u:%g %p\n'
    find internal cmd assets -path '*/vendor/*' -prune -o -type f -name '*.go' -print | xargs wc -l | sort -n | tail -n 20

Record any hand-written code file over 600 lines with its cause and a concrete recommendation. Do not perform an unrelated split without separate user approval. The accepted scope leaves the existing 602-line `internal/game/game_test.go` and near-limit `internal/game/building_bar_test.go` unchanged.

## Validation and Acceptance

Acceptance is met when a new game shows `Grasslands` or `Hills` beside all four initial explore buttons, each label remains stable through redraws, pause, camera movement, and zoom, and clicking the circular icon reveals a Plot whose stored biome matches the preview. Revealing a Plot must expose labeled choices on the newly expanded frontier without rerolling shared coordinates. Clicking label text alone must not explore. Terrain must remain hidden until reveal, and existing road, join, placement, Raid, and biome-weight behavior must continue to pass automated tests.

Visual evidence must show the initial labeled frontier and at least one expanded labeled frontier without obscuring roads, joins, empty terrain, or screen-space UI. `go test ./...`, `go test -race ./...`, screenshot capture, whitespace, ownership, and line-count checks must pass.

## Idempotence and Recovery

Frontier assignment is idempotent because existing coordinates are never overwritten. Plot reveal remains idempotent because explored coordinates return before generation. Formatting, tests, and screenshot capture are safe to rerun. If a label overlaps its button or explored terrain, correct the directional layout helper rather than expanding the hitbox. If a test constructs a zero-value `Map`, initialize storage without silently assigning random frontier state during drawing.

## Artifacts and Notes

Store screenshots under `plans/49-frontier-biome-labels/screenshots/`. Record final test output, visual review, ownership results, and line counts in this plan's `Outcomes & Retrospective`.

## Interfaces and Dependencies

No external package API, dependency, asset, save format, terrain type, or resource behavior changes. Private game-package additions include frontier-biome storage on `Map`, deterministic frontier assignment helpers, a biome-specific Plot generator, `exploreButton.Biome`, biome label formatting, and directional label geometry. The existing circular explore-button rectangle remains the only reveal hit target.
