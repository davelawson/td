# Add Game Design Control Document

This ExecPlan is a living document. The sections `Progress`, `Surprises & Discoveries`, `Decision Log`, and `Outcomes & Retrospective` must be kept up to date as work proceeds.

This plan is maintained according to `PLANS.md` in the repository root. It is saved at `plans/07-game-design-control-doc.md` because `plans/00-initial-ebitengine-menu.md` through `plans/06-ingame-menu-overlay.md` already exist.

## Purpose / Big Picture

The repository has durable documents for current product behavior, roadmap direction, visual design, architecture, and code style, but it does not yet have a dedicated home for intended gameplay design decisions. After this change, contributors can record intended game design in `GAME.md` even when implementation has not caught up yet. Future gameplay design decisions must keep `GAME.md` synchronized instead of being left only in chat or scattered across plans.

## Progress

- [x] (2026-05-08 13:08Z) Created this ExecPlan for adding `GAME.md` as a new root control document.
- [x] (2026-05-08 13:08Z) Added `GAME.md` with initial living game-design structure, current decisions, and open questions.
- [x] (2026-05-08 13:08Z) Updated root control documents to list `GAME.md` and require synchronization when game design decisions change.
- [x] (2026-05-08 13:08Z) Ran docs and test validation.

## Surprises & Discoveries

- Observation: The repository already has later prototype work than the original `AGENTS.md` project-context text describes.
  Evidence: `README.md`, `PRODUCT.md`, `ARCHITECTURE.md`, and the Go code describe a runnable prototype with menu flow, Wizard name entry, game update counting, pause, and in-game overlay behavior.

## Decision Log

- Decision: Make `GAME.md` a living design document rather than a detailed mechanics bible.
  Rationale: The game is still in early prototype foundation, so recording pillars, intended loop, decisions, and open questions is more useful than inventing untested resource names, tower stats, wave schedules, or progression rules.
  Date/Author: 2026-05-08 / Codex

- Decision: Treat "automatically incorporated" as agent and ExecPlan process rules, not enforcement tooling.
  Rationale: The user chose process-level synchronization. This keeps the change documentation-only while making future agents responsible for updating `GAME.md` whenever game design decisions are made.
  Date/Author: 2026-05-08 / Codex

## Outcomes & Retrospective

`GAME.md` now exists as the durable source of truth for intended game design. It records the player fantasy, design pillars, intended core loop, planned gameplay systems, current game design decisions, open questions, and a game-design decision log. Root control documents now include `GAME.md` in their inventories and synchronization rules.

Validation completed successfully:

    go test ./...
    ok  	td/cmd/td	(cached)
    ok  	td/internal/game	(cached)
    ok  	td/internal/menu	(cached)
    ?   	td/internal/ui	[no test files]

    git diff --check
    No output, meaning no whitespace errors were found.

This plan did not add or change code, so the hand-written code-file line-count review requirement does not apply.

## Context and Orientation

`td` is a local Go/Ebitengine desktop prototype for a medieval wizardry tower-defense game. The current implementation is a runnable shell with menus, Wizard name entry, a placeholder game screen, pause behavior, and an in-game overlay menu. The intended future game combines exploration, resource gathering, base-building, and tower-defense combat, but those gameplay systems do not exist yet.

Before this change, the root control documents split ownership this way: `PRODUCT.md` described current user-visible behavior, `ROADMAP.md` described future direction and priorities, `DESIGN.md` described visual and interaction direction, `ARCHITECTURE.md` described code ownership, `CODESTYLE.md` described source conventions, `PLANS.md` described ExecPlan rules, and `AGENTS.md` described repository-specific agent instructions. No document specifically owned intended game design decisions independent of implementation state.

This change introduces `GAME.md` as that missing owner. `GAME.md` may describe planned gameplay that does not exist yet. `PRODUCT.md` remains the source of truth for current implemented behavior. `ROADMAP.md` remains the source of truth for sequencing, priorities, and non-priorities. `DESIGN.md` remains the source of truth for visual language and interaction presentation.

## Plan of Work

Create `GAME.md` at the repository root. Its opening prose must define it as the durable source of truth for intended game design and explain how it differs from `PRODUCT.md`, `ROADMAP.md`, and `DESIGN.md`. Seed it with the design decisions already present across the repository: single local PC player, medieval wizardry setting, wizard player identity, exploration, resource gathering, base-building, tower-defense combat, and deferred save/load and campaign structure.

Keep `GAME.md` at the right level of detail for the current prototype. Include player fantasy, design pillars, intended core loop, planned systems, current decisions, open questions, and a decision log. Do not invent detailed mechanics such as exact resource names, tower statistics, enemy schedules, progression trees, or campaign structure.

Update the root documentation inventories and synchronization rules. `AGENTS.md` should list `GAME.md` as a core document and say game design changes must update it. `PLANS.md` should tell future ExecPlans to read, update, and validate `GAME.md` when gameplay design decisions are relevant. `README.md`, `PRODUCT.md`, `ROADMAP.md`, and `CODESTYLE.md` should include `GAME.md` in their document relationship or inventory sections.

## Concrete Steps

From the repository root, inspect the current state:

    pwd
    rg --files --hidden -g '!.git/**'
    git status --short

Create `GAME.md` and this plan file at `plans/07-game-design-control-doc.md`.

Edit `AGENTS.md`, `PLANS.md`, `README.md`, `PRODUCT.md`, `ROADMAP.md`, and `CODESTYLE.md` to include `GAME.md` in their control-document descriptions and synchronization rules.

Validate with:

    go test ./...
    git diff --check
    git status --short

Because this plan changes documentation only, do not run `gofmt` and do not perform a hand-written code-file line-count review.

## Validation and Acceptance

The change is accepted when `GAME.md` exists and clearly owns intended game design independent of implementation state. A reviewer should be able to read it and understand the current player fantasy, pillars, intended loop, broad planned systems, decided facts, and unresolved design questions.

Documentation is accepted when the other root control documents consistently list `GAME.md` and state that future meaningful game design decisions must update it. `PRODUCT.md` must still only claim current implemented behavior, and `ROADMAP.md` must still focus on sequencing and priorities rather than becoming the detailed design home.

Command validation is accepted when `go test ./...` succeeds and `git diff --check` reports no whitespace errors.

## Idempotence and Recovery

The change is documentation-only and additive. Re-running tests and whitespace checks is safe. If future edits reveal duplicated or conflicting design statements, prefer `GAME.md` for intended gameplay design, `PRODUCT.md` for current implemented behavior, and `ROADMAP.md` for sequencing. Update the mismatched document in the same change.

## Artifacts and Notes

Important artifacts:

    GAME.md
    plans/07-game-design-control-doc.md
    AGENTS.md
    PLANS.md
    README.md
    PRODUCT.md
    ROADMAP.md
    CODESTYLE.md

## Interfaces and Dependencies

This plan introduces no code interfaces and no new dependencies. It introduces one documentation interface: future gameplay design work must read and update `GAME.md` when it makes, changes, or rejects a meaningful game design decision.

