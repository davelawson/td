# Product State: td

`PRODUCT.md` is the durable source of truth for what `td` does for users right now. Update it whenever the repository gains, removes, or materially changes a user-visible workflow, capability, limitation, or scope boundary.

## Product Summary

`td` is planned as a PC tower-defense game prototype for a single local player. The intended game blends exploration, base-building, resource gathering, and conventional tower-defense encounters in a medieval wizardry fantasy setting.

The current repository does not yet ship a playable game. It is in bootstrap state: the project has a defined direction, Go/Ebitengine has been selected as the runtime stack, and the first implementation plan describes how to create the smallest playable desktop app.

## Users and Jobs To Be Done

The current user is the developer-player validating whether the game can become a playable local prototype. Their immediate job is to turn the repository into a runnable desktop application with a visible game window and a working quit path.

Future players are expected to want a strategy game where they explore, build a base, gather resources, and defend against threats with tower-defense mechanics. Those systems do not exist yet and belong to the roadmap rather than current product truth.

## Current Capabilities

### Project Direction

`Core`: The repository records the product name, intended genre, target platform phase, art direction, runtime stack, and first implementation slice. This gives contributors enough context to plan work without relying on chat history.

### Agent Runtime Defaults

`Core`: `.codex/config.toml` keeps Codex project defaults for trusted local work. These settings are agent configuration only and are not application runtime configuration.

### Planning Workflow

`Core`: Substantial work must use an ordered ExecPlan under `plans/`, following `PLANS.md`. The first plan is `plans/00-initial-ebitengine-menu.md`, which initializes the Go module and Ebitengine app and creates the first menu screen.

### Missing Runtime

`Core`: There is currently no Go module, executable game, test suite, asset pipeline, save system, campaign system, CI pipeline, license, or release packaging.

## Core Workflows

### Repository Bootstrap

A contributor opens the repository, reads the root control documents, and sees that `td` is a local Go/Ebitengine PC game prototype. The workflow ends with a repository that has project-specific control documents and a first implementation plan, but no product feature code yet.

### First Planned Playable Workflow

After `plans/00-initial-ebitengine-menu.md` is implemented, a contributor should be able to run `go run ./cmd/td`, see a desktop window with a medieval wizardry main menu, click a quit option, and observe the app closing cleanly.

## Product Constraints and Known Limits

- The current target is a local prototype only.
- Distribution, release packaging, CI, license selection, and store targets are deferred.
- The first playable version intentionally includes only a main menu and quit behavior.
- Saving the game and campaign structure are explicit non-goals for the first phase.
- Exploration, base-building, resource gathering, and tower-defense gameplay are intended but not implemented.

## Non-Goals

The current phase is not trying to build campaign progression, save/load, multiplayer, online services, release packaging, or a full tower-defense encounter. It is also not trying to choose final art production workflows before the game loop proves useful.

## Relationship To Other Control Documents

- `README.md` explains repository status, commands, and layout.
- `ROADMAP.md` explains intended future direction.
- `DESIGN.md` captures the fantasy design direction and UI review expectations.
- `CODESTYLE.md` defines source conventions and file-size expectations.
- `ARCHITECTURE.md` explains intended code structure and boundaries.
- `PLANS.md` defines ExecPlan requirements.
- `AGENTS.md` explains repository-specific coding-agent instructions.

When these files disagree about current user-visible behavior, treat this file as the source of truth and update the mismatch in the same change.

## Open Questions

- What license should the project use?
- Which platforms beyond local desktop prototypes should matter after the first playable slice?
- Should tests live beside Go packages, under `tests/`, or use a mixed strategy once the code layout exists?
