# Product State: td

`PRODUCT.md` is the durable source of truth for what `td` does for users right now. Update it whenever the repository gains, removes, or materially changes a user-visible workflow, capability, limitation, or scope boundary.

## Product Summary

`td` is planned as a PC tower-defense game prototype for a single local player. The intended game blends exploration, base-building, resource gathering, and conventional tower-defense encounters in a medieval wizardry fantasy setting.

The current repository ships a small playable shell: a local Go/Ebitengine desktop app that opens a 1920x1080 main menu, can navigate to a New Game configuration screen and a placeholder Settings screen, accepts a Wizard name up to 32 characters, starts a first game screen after a name is entered, shows a static 15x15 home Plot scene with the centered sprite-backed Sanctum, a straight road north to the Plot edge, sprite-backed Bow and Flame Bolt Towers across the road from each other, and a pine-tree border around the Plot edge except at the road opening, supports mouse-wheel map zoom and `WASD` camera panning, lets the player left-click structures and visible raiders to select them with a brighter selected state and bottom-right selection panel, shows a top bar with prototype Chapter, Day, sprite-backed resource counts, phase, and Sanctum barricade information, shows a left building bar with Bow Tower and Flame Bolt Tower icons, construction costs, affordable-option hover highlighting with hover-only cost emphasis, and calm-phase drag placement for affordable towers on empty grass-like Tiles, starts deterministic placeholder Raids with sprite-backed skeleton and zombie enemies and visible enemy health bars from a `Next Raid` button, fires Bow Tower arrows and Flame Bolt Tower flame bolts at enemies within range, plays a short embedded sound and grants small template resources when tower damage defeats a raider, tracks logical game updates in a debug counter, toggles pause with SPACE, opens an in-game overlay menu with ESC, and can quit cleanly from the main menu. Its resizable window uses a pixel-sized drawable layout so text does not stretch when the window is enlarged. It does not yet include exploration, resource gathering, tower upgrades, broader base-building, saves, real settings, campaign structure, an asset pipeline, release packaging, or CI.

## Users and Jobs To Be Done

The current user is the developer-player validating whether the game can become a playable local prototype. Their immediate job is to run the desktop application, confirm that a visible game window appears, navigate the small menu flow, and confirm that the quit path closes cleanly.

Future players are expected to want a strategy game where they explore, build a base, gather resources, and defend against threats with tower-defense mechanics. Only the first limited tower placement action exists today; the broader versions of those systems belong to the roadmap rather than current product truth.

## Current Capabilities

### Project Direction

`Core`: The repository records the product name, intended genre, target platform phase, art direction, runtime stack, and first implementation slice. This gives contributors enough context to plan work without relying on chat history.

### Agent Runtime Defaults

`Core`: `.codex/config.toml` keeps Codex project defaults for trusted local work. These settings are agent configuration only and are not application runtime configuration.

### Planning Workflow

`Core`: Substantial work must use an ordered ExecPlan under `plans/`, following `PLANS.md`. `plans/00-initial-ebitengine-menu.md` initialized the Go module and Ebitengine app. `plans/01-expanded-main-menu.md` expanded the main menu flow. `plans/03-new-game-configuration.md` adds the Wizard name configuration screen. `plans/04-resolution-and-pixel-text-scaling.md` sets the 1920x1080 default window and current resize policy. `plans/05-main-game-update-loop.md` adds the first game state, update counter, and pause behavior. `plans/06-ingame-menu-overlay.md` adds the ESC in-game overlay menu. `plans/10-basic-camera.md` adds mouse-wheel zoom and `WASD` camera panning for the game scene.

### Runtime Shell

`Core`: The repository has a Go module, an Ebitengine executable under `cmd/td/`, a typed runtime asset catalog under `assets/`, a small runtime sound manager under `internal/sound/`, pure menu behavior under `internal/menu/`, pure game-update behavior under `internal/game/`, static prototype map data and camera state under `internal/game/`, selected structure and raider state and detail-panel behavior under `internal/game/`, building-bar hover and drag-placement behavior under `internal/game/`, deterministic placeholder Raid state and Bow Tower and Flame Bolt Tower projectile combat under `internal/game/`, and Go tests for asset loading, menu hit testing, disabled menu targets, action selection, screen routing, Wizard name input behavior, app mode startup, game update counting, pause behavior, camera zoom and pan behavior, selected-object hit testing and clearing, selected-object detail-panel formatting and click blocking, building-bar bounds, construction-cost display data, hover targeting, affordability-gated hover highlighting, hover cost fit, click blocking, tower drag placement, in-game menu behavior, top-bar status formatting, resource icon catalog wiring, Raid starts, skeleton and zombie enemy sprite and reward catalog wiring, staggered enemy spawning, first-Raid zombie composition, enemy world-coordinate movement, tower targeting, projectile hits, enemy health, raider-defeat sound events, combat-defeat resource rewards, enemy health-bar proportions and colors, Barricade spending, breach handling, and default home Plot invariants including the road opening through the tree border and the two starting towers beside the path.

### Missing Gameplay And Operations

`Core`: There is currently no tower upgrade, exploration, resource gathering, broader base-building, asset pipeline, music, volume controls, save system, settings implementation, campaign system, CI pipeline, license, or release packaging. Resource rewards exist only as small enemy-template payouts when tower damage defeats a raider. Tower placement exists only as a first calm-phase building-bar drag action for affordable Bow and Flame Bolt Towers on empty grass-like Tiles.

## Core Workflows

### Repository Bootstrap

A contributor opens the repository, reads the root control documents, and sees that `td` is a local Go/Ebitengine PC game prototype. The bootstrap workflow has already produced project-specific control documents and the first runnable menu shell.

### Main Menu Workflow

A contributor runs `go run ./cmd/td` and sees a 1920x1080 desktop window titled `td` with a medieval wizardry main menu. Resizing the window recenters the menu in the current drawable area while text keeps its raw pixel size. The menu offers `New`, `Load`, `Settings`, and `Quit`. Clicking `New` opens a New Game configuration screen with a focused Wizard name field, disabled `Start` button, and active `Cancel` button. Typing edits the Wizard name up to 32 characters and Backspace removes the last typed character. Once the name is non-empty, `Start` becomes active. Clicking `Start` closes the menu and opens the first game screen. The game screen shows a static 15x15 home Plot with the centered Sanctum rendered from a loaded sprite, a straight road north to the Plot edge, a Bow Tower and Flame Bolt Tower rendered from loaded sprites across the road from each other, and pine trees around the Plot edge except where the road exits at the north-center edge. Mouse-wheel input zooms the map in and out, and `W`, `A`, `S`, and `D` pan the camera without clamping the Plot to the screen. Left-clicking the Sanctum, Bow Tower, Flame Bolt Tower, or visible raider sprites selects that object and draws it brighter; left-clicking elsewhere clears selection. Raider selection takes priority when a raider overlaps a structure. Selecting a raider shows a bottom-right panel with raider type, current health, max health, health percentage, movement speed, and Sanctum damage. Selecting a combat tower shows tower type, range, attack speed, and damage. Selecting the Sanctum shows a basic Sanctum panel. Clicking inside the panel does not clear or change selection. It also shows a top bar with fixed prototype values for Chapter, Day, Wood, Stone, and Metal counts shown beside resource icons, calm time before the next Raid, and Sanctum barricade charges. A dark building bar fills the left edge of the playable scene below the top bar and shows Bow Tower and Flame Bolt Tower icons with colour-coded construction costs underneath; hovering an icon brightens it and makes only its matching cost row bolder when the current resources cover that tower's cost. Left-dragging an affordable tower icon attaches a half-sized copy of that icon to the cursor. Releasing over an empty grass-like Tile during calm play places the tower and deducts its displayed cost from Wood, Stone, and Metal. With default resources, Bow Tower is affordable and Flame Bolt Tower is not, so Flame Bolt Tower cannot be dragged from the bar until enough Metal exists. Occupied Tiles, road Tiles, forest Tiles, active Raids, and breached games reject placement without spending resources. SPACE-paused calm play still allows placement. A `Next Raid` button sits beside the building bar near the bottom-left of the scene, and a debug logical update counter appears near the lower-right. Clicking `Next Raid` starts a deterministic placeholder Raid immediately: Raid 1 spawns skeleton, zombie, skeleton, zombie, skeleton from the north-center road edge, with later placeholder Raids still spawning skeletons only. Enemies store Sanctum-centered world coordinates in Tile units, render from their loaded enemy sprite with a health bar above the sprite, move down the road toward the Sanctum, and are removed when defeated or when they reach the Sanctum. The health bar is green and sprite-width at full health; as health is lost, the filled bar shrinks and shifts toward red. During a Raid, towers target the in-range enemy closest to the Sanctum. The Bow Tower fires visible arrow projectiles about once per second and deals 10 damage on hit. The Flame Bolt Tower fires visible flame bolts every 1.5 seconds and deals 20 damage on hit. When tower damage defeats a raider, the game plays a short embedded sound effect. Current enemy health and movement speed come from each enemy template, with skeletons at 50 health and 1.0 Tiles per second and zombies at 75 health and 0.7 Tiles per second. During a Raid, the top bar shows the real remaining enemy count and the `Next Raid` button is disabled. Each reaching enemy spends one Barricade charge if any remain. If an enemy reaches the Sanctum when Barricade is zero, the top bar reports `Sanctum breached`, the active Raid is cleared, and new Raids cannot start. Pressing SPACE toggles pause; while paused, a `PAUSED` label appears and the logical update counter, Raid simulation, and combat stop advancing, but camera zoom, pan, object selection, and calm building still work. Pressing ESC in the game opens a centered in-game menu over the still-visible game view, darkens the rest of the scene by about 50%, pauses the game, blocks camera controls, selection, building, and Raid interaction, and offers `Resume` and `Surrender`. Pressing ESC again or clicking `Resume` closes the overlay and restores the previous pause state. Clicking `Surrender` leaves the game and returns to the main menu. Clicking `Cancel` on the New Game screen returns to the main menu. Clicking `Settings` opens a placeholder Settings screen with a `Back` button. Clicking `Back` returns to the main menu. `Load` is visibly disabled and does nothing because saving and loading do not exist yet. Clicking `Quit` closes the app cleanly from the main menu.

## Product Constraints and Known Limits

- The current target is a local prototype only.
- Distribution, release packaging, CI, license selection, and store targets are deferred.
- The current playable shell intentionally includes only a small menu flow, Wizard name entry, a first game screen with a static tree-bordered home Plot scene, sprite-backed Sanctum, automated Bow and Flame Bolt Towers beside the road, a building bar with two tower icons, affordable-option hover feedback, construction costs, and limited calm-phase drag placement, basic camera zoom and pan, left-click selection and information panels for structures and raiders, deterministic placeholder Raids with sprite-backed skeleton and zombie enemies and health bars on the fixed north road, first-pass projectile combat, a logical update counter, pause behavior, an in-game overlay menu, surrender-to-menu behavior, and quit behavior.
- Saving the game and campaign structure are explicit non-goals for the first phase.
- Starting a new game only opens the first static home Plot scene; exploration, resource gathering, tower upgrades, and broader base-building are not implemented. The left building bar can spend resources only for this first tower placement action, and tower-damage raider defeats can add small prototype resource rewards. The top bar still uses fixed prototype Chapter, calm-time, and Barricade starting values until those systems exist, but resource labels are represented by icons, resource counts change after building or combat-defeat rewards, and it uses real remaining enemy counts during active Raids.
- Settings are represented only by a placeholder screen; no configurable options exist yet.
- Exploration, broader base-building, resource gathering, and broader tower-defense gameplay are intended but not implemented.

## Non-Goals

The current phase is not trying to build campaign progression, save/load, multiplayer, online services, release packaging, or a full tower-defense encounter. The current Raid slice is a deterministic enemy placeholder with skeletons, two first-Raid zombies, two automated starting towers, small combat-defeat resource rewards, and no multiple paths. The project is also not trying to choose final art production workflows before the game loop proves useful.

## Relationship To Other Control Documents

- `README.md` explains repository status, commands, and layout.
- `ROADMAP.md` explains intended future direction.
- `GAME.md` records intended game design decisions and open gameplay questions, including planned behavior that may not exist yet.
- `DESIGN.md` captures the fantasy design direction and UI review expectations.
- `ART.md` captures guidance for generated art assets, prompt patterns, asset review criteria, and prototype asset constraints.
- `CODESTYLE.md` defines source conventions and file-size expectations.
- `ARCHITECTURE.md` explains intended code structure and boundaries.
- `PLANS.md` defines ExecPlan requirements.
- `AGENTS.md` explains repository-specific coding-agent instructions.

When these files disagree about current user-visible behavior, treat this file as the source of truth and update the mismatch in the same change.

## Open Questions

- What license should the project use?
- Which platforms beyond local desktop prototypes should matter after the first playable slice?
- Should tests live beside Go packages, under `tests/`, or use a mixed strategy once the code layout exists?
