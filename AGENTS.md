# Repository Guidelines

## Project Context

`td` is a local PC prototype for a Go/Ebitengine tower-defense game. The intended game blends exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository is bootstrapped and contains a small runnable prototype shell with menu flow, Wizard name entry, a placeholder game screen, pause behavior, and an in-game overlay menu. Gameplay systems such as exploration, base-building, resource gathering, and tower-defense combat are still unimplemented.

## ExecPlans

When writing complex features or significant refactors, use an ExecPlan as described in `PLANS.md`.

Create new ExecPlan files under `plans/` with a two-digit ordering prefix and a short kebab-case name, for example `plans/00-initial-ebitengine-menu.md` or `plans/01-static-prototype-scene.md`. Use the next unused number after the highest existing prefixed plan.

When creating an ExecPlan that adds or changes code, include a final `Progress` and `Concrete Steps` item to check line counts for hand-written code files. Use the `CODESTYLE.md` preference that code files stay below 600 lines. If any file exceeds that preference, list the file and line count, recommend a concrete response, and ask the user to approve extra work before implementing it unless the accepted plan already included it.

## Planning and Change Intake

Treat requested changes as suggested directions until scope, tradeoffs, and long-term consequences are understood. If a requested change would create unclear boundaries, premature abstractions, avoidable technical debt, or product scope drift, explain the concern plainly and offer concrete alternatives before proceeding.

Planning for gameplay systems should gather enough context about player workflow, future extension points, testing strategy, maintainability, and intended design in `GAME.md` before code changes begin. Keep each gameplay slice deliberately small and observable.

## Core Documents

Treat the root documentation files as durable project control documents:

- `README.md` explains what the project is, how to run it, and how to validate changes.
- `PRODUCT.md` captures current user-visible product state, workflows, capability boundaries, and important limitations.
- `ROADMAP.md` captures intended product direction, planned capabilities, strategic priorities, and explicit non-priorities.
- `GAME.md` captures intended game design decisions, player fantasy, gameplay pillars, planned systems, and open game-design questions regardless of implementation state.
- `DESIGN.md` captures medieval wizardry visual direction, interaction principles, and visual-review expectations.
- `ART.md` captures guidance for generated art assets, prompt patterns, asset review criteria, and prototype asset constraints.
- `PLANS.md` defines how ExecPlans must be written and maintained.
- `CODESTYLE.md` defines Go source formatting, naming, documentation style, commenting standards, and code-file size expectations.
- `ARCHITECTURE.md` captures high-level code ownership, boundaries, and invariants.

When a change materially affects current user-visible capabilities, workflows, scope boundaries, or important product limitations, update `PRODUCT.md` in the same change. When a change materially affects product vision, intended audience, strategic priorities, planned capabilities, sequencing assumptions, or explicit non-priorities, update `ROADMAP.md` in the same change. When a change materially affects source conventions, naming rules, documentation conventions, or commenting standards, update `CODESTYLE.md` in the same change. When a change materially affects design language, update `DESIGN.md` in the same change. When a change materially affects generated art-asset guidance, prompt patterns, asset review criteria, or prototype asset constraints, update `ART.md` in the same change. When a change materially affects structure, ownership, or system boundaries, update `ARCHITECTURE.md` in the same change.
When a change materially makes, revises, rejects, or depends on intended gameplay design decisions, update `GAME.md` in the same change.

If you introduce a new root-level `ALLCAPS.md` file, treat it as a new control document by default. Define its purpose inside the file, update this section, update `PLANS.md` if ExecPlans must read or validate it, and state what kinds of changes must keep it in sync.

## Project Structure and Module Organization

The Go module and first runtime shell already exist.

Current and expected layout:

- `cmd/td/` contains the executable entry point and Ebitengine startup.
- `internal/` contains reusable game packages once behavior outgrows the entry point.
- `assets/` contains static runtime assets once assets exist.
- `plans/` contains ordered ExecPlans and plan evidence.
- `.agents/skills/` contains repo-local agent workflows.
- `.codex/config.toml` contains Codex project defaults, not app configuration.

Keep the layout simple. Do not create a scene framework, ECS, asset pipeline, save system, campaign system, or packaging setup until a plan shows why it is needed.

## File Ownership

All files and directories in this repository should be owned by `dave:dave`. If a coding agent or command creates or modifies files as another user, normalize ownership before finishing the change:

- `chown -R dave:dave /home/dave/dev/ai/td`
- `find . -xdev ! -user dave -printf '%u:%g %p\n'` should print nothing when run from the repository root.

## Build, Test, and Development Commands

Use these commands from the repository root:

- `go test ./...` to run all tests.
- `go run ./cmd/td` to start the local prototype.
- `go mod tidy` after dependency or import changes.
- `git diff --check` checks for whitespace errors.
- `git status --short` shows pending changes.

No Makefile, CI pipeline, release process, license, or packaging workflow is defined yet.

## Testing Guidelines

New behavior should ship with tests when it can be exercised without a graphics window. For Ebitengine work, prefer tests around pure behavior such as menu hit testing, action selection, geometry, state transitions, and later gameplay rules.

Manual visual validation is acceptable for the first desktop window and menu rendering. Save visual evidence under the active plan directory when a plan changes rendered output.

## Design and Visual Review

Follow `DESIGN.md` for the medieval wizardry direction. UI changes should prioritize readable text, stable hit targets, and clear interaction states before decorative styling.

When work affects rendered game output, the ExecPlan should say what screenshot or visual evidence will be captured. If there is no existing runnable app, record that no screenshotable baseline exists before implementation and capture the first rendered result after implementation.

## Commit and Pull Request Guidelines

Follow the existing short imperative commit style, such as `Initial commit`.

Pull requests should explain:

- What changed
- Why it changed
- How it was validated

Include screenshots when the change affects rendered game output or documentation presentation.
