# Product State: td

`PRODUCT.md` is the durable source of truth for what `td` does for users right now. Update it whenever the repository gains, removes, or materially changes a user-visible workflow, capability, limitation, or scope boundary.

## Product Summary

`td` is planned as a PC tower-defense game prototype for a single local player. The intended game blends exploration, base-building, resource gathering, and conventional tower-defense encounters in a medieval wizardry fantasy setting.

The current repository ships a small playable shell: a local Go/Ebitengine desktop app that opens a 1920x1080 main menu, can navigate to a New Game configuration screen and a placeholder Settings screen, accepts a Wizard name up to 32 characters, starts a first game screen after a name is entered, shows a static empty 15x15 home Plot scene with the centered Sanctum and a straight road north to the Plot edge, supports mouse-wheel map zoom and `WASD` camera panning, shows a top bar with prototype Chapter, Day, resource, phase, and Sanctum barricade information, tracks logical game updates in a debug counter, toggles pause with SPACE, opens an in-game overlay menu with ESC, and can quit cleanly from the main menu. Its resizable window uses a pixel-sized drawable layout so text does not stretch when the window is enlarged. It does not yet include tower-defense combat, exploration, base-building, real resource changes, saves, real settings, campaign structure, art assets, release packaging, or CI.

## Users and Jobs To Be Done

The current user is the developer-player validating whether the game can become a playable local prototype. Their immediate job is to run the desktop application, confirm that a visible game window appears, navigate the small menu flow, and confirm that the quit path closes cleanly.

Future players are expected to want a strategy game where they explore, build a base, gather resources, and defend against threats with tower-defense mechanics. Those systems do not exist yet and belong to the roadmap rather than current product truth.

## Current Capabilities

### Project Direction

`Core`: The repository records the product name, intended genre, target platform phase, art direction, runtime stack, and first implementation slice. This gives contributors enough context to plan work without relying on chat history.

### Agent Runtime Defaults

`Core`: `.codex/config.toml` keeps Codex project defaults for trusted local work. These settings are agent configuration only and are not application runtime configuration.

### Planning Workflow

`Core`: Substantial work must use an ordered ExecPlan under `plans/`, following `PLANS.md`. `plans/00-initial-ebitengine-menu.md` initialized the Go module and Ebitengine app. `plans/01-expanded-main-menu.md` expanded the main menu flow. `plans/03-new-game-configuration.md` adds the Wizard name configuration screen. `plans/04-resolution-and-pixel-text-scaling.md` sets the 1920x1080 default window and current resize policy. `plans/05-main-game-update-loop.md` adds the first game state, update counter, and pause behavior. `plans/06-ingame-menu-overlay.md` adds the ESC in-game overlay menu. `plans/10-basic-camera.md` adds mouse-wheel zoom and `WASD` camera panning for the game scene.

### Runtime Shell

`Core`: The repository has a Go module, an Ebitengine executable under `cmd/td/`, pure menu behavior under `internal/menu/`, pure game-update behavior under `internal/game/`, static prototype map data and camera state under `internal/game/`, and Go tests for menu hit testing, disabled menu targets, action selection, screen routing, Wizard name input behavior, app mode startup, game update counting, pause behavior, camera zoom and pan behavior, in-game menu behavior, top-bar status formatting, and default home Plot invariants.

### Missing Gameplay And Operations

`Core`: There is currently no tower-defense encounter, exploration, base-building, resource gathering, asset pipeline, save system, settings implementation, campaign system, CI pipeline, license, or release packaging.

## Core Workflows

### Repository Bootstrap

A contributor opens the repository, reads the root control documents, and sees that `td` is a local Go/Ebitengine PC game prototype. The bootstrap workflow has already produced project-specific control documents and the first runnable menu shell.

### Main Menu Workflow

A contributor runs `go run ./cmd/td` and sees a 1920x1080 desktop window titled `td` with a medieval wizardry main menu. Resizing the window recenters the menu in the current drawable area while text keeps its raw pixel size. The menu offers `New`, `Load`, `Settings`, and `Quit`. Clicking `New` opens a New Game configuration screen with a focused Wizard name field, disabled `Start` button, and active `Cancel` button. Typing edits the Wizard name up to 32 characters and Backspace removes the last typed character. Once the name is non-empty, `Start` becomes active. Clicking `Start` closes the menu and opens the first game screen. The game screen shows a static empty 15x15 home Plot with the centered Sanctum and a straight road north to the Plot edge. Mouse-wheel input zooms the map in and out, and `W`, `A`, `S`, and `D` pan the camera without clamping the Plot to the screen. It also shows a top bar with fixed prototype values for Chapter, Day, resources, calm time before the next Raid, and Sanctum barricade charges, plus a debug logical update counter. Pressing SPACE toggles pause; while paused, a `PAUSED` label appears and the logical update counter stops advancing, but camera zoom and pan still work. Pressing ESC in the game opens a centered in-game menu over the still-visible game view, darkens the rest of the scene by about 50%, pauses the game, blocks camera controls, and offers `Resume` and `Surrender`. Pressing ESC again or clicking `Resume` closes the overlay and restores the previous pause state. Clicking `Surrender` leaves the game and returns to the main menu. Clicking `Cancel` on the New Game screen returns to the main menu. Clicking `Settings` opens a placeholder Settings screen with a `Back` button. Clicking `Back` returns to the main menu. `Load` is visibly disabled and does nothing because saving and loading do not exist yet. Clicking `Quit` closes the app cleanly from the main menu.

## Product Constraints and Known Limits

- The current target is a local prototype only.
- Distribution, release packaging, CI, license selection, and store targets are deferred.
- The current playable shell intentionally includes only a small menu flow, Wizard name entry, a first game screen with a static home Plot scene, basic camera zoom and pan, a logical update counter, pause behavior, an in-game overlay menu, surrender-to-menu behavior, and quit behavior.
- Saving the game and campaign structure are explicit non-goals for the first phase.
- Starting a new game only opens the first static home Plot scene; exploration, base-building, real resource changes, and combat are not implemented. The top bar uses fixed prototype values until those systems exist, and logical ticks do not mutate the scene.
- Settings are represented only by a placeholder screen; no configurable options exist yet.
- Exploration, base-building, resource gathering, and tower-defense gameplay are intended but not implemented.

## Non-Goals

The current phase is not trying to build campaign progression, save/load, multiplayer, online services, release packaging, or a full tower-defense encounter. It is also not trying to choose final art production workflows before the game loop proves useful.

## Relationship To Other Control Documents

- `README.md` explains repository status, commands, and layout.
- `ROADMAP.md` explains intended future direction.
- `GAME.md` records intended game design decisions and open gameplay questions, including planned behavior that may not exist yet.
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
