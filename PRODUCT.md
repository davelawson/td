# Product State: td

`PRODUCT.md` is the durable source of truth for what `td` does for users right now. Update it whenever the repository gains, removes, or materially changes a user-visible workflow, capability, limitation, or scope boundary.

## Product Summary

`td` is planned as a PC tower-defense game prototype for a single local player. The intended game blends exploration, base-building, resource gathering, and conventional tower-defense encounters in a medieval wizardry fantasy setting.

The current repository ships a small playable shell that starts a static home Plot with only the Sanctum, 100 Wood, 50 Stone, 20 Metal, and zero inhabitants. The 260-pixel building bar groups build choices into `Housing`, `Economic`, and `Defenses` tabs, defaults to `Housing`, and shows each building's cost and population or staffing values to the right of its icon. House costs 20 Wood, requires no staff, and immediately grants 2 available and total Peasants when placed. Barracks costs 10 Wood and 10 Stone, consumes 2 available and total Peasants, and grants 2 available and total Soldiers. Woodcutter, Stone Quarry, and Iron Mine each reserve one available Peasant and add 10 of Wood, Stone, or Metal respectively after a defeated Raid. Tower construction requires both resources and available staff. A successful staffed build reserves staff by reducing available counts while preserving totals. Buildable icon squares use green outlines; buildings without enough resources, population, or staff use red outlines and draw their icons at 70% opacity. Camera inspection, selection, deterministic Raids, projectile combat, rewards, pause, overlay, and menu behavior remain implemented.

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

`Core`: The repository has a Go module, an Ebitengine executable under `cmd/td/`, typed runtime assets, menu and game packages, a Sanctum-only starting Plot, status state, camera and selection behavior, a widened tabbed building bar with right-side construction costs, population-grant metadata, staffing metadata, green/red capacity outlines, capacity opacity feedback, drag placement, deterministic Raids, and Bow, Flame Bolt, and Catapult projectile combat. Tests cover these behaviors without depending on free starting towers.

### Missing Gameplay And Operations

`Core`: There is currently no tower upgrade, exploration, timed or manual resource gathering, timed inhabitant recruitment, reassignment, staff release, broader base-building, save system, or campaign system. House placement is a calm-phase action gated by Wood that increases Peasant population immediately. Barracks placement is a calm-phase action gated by Wood, Stone, and available Peasants that converts Peasants into Soldiers. Economic building placement is a calm-phase action gated by resources and one available Peasant, and successful Raids trigger the only implemented resource production. Tower placement is a calm-phase action gated by resources and available staffing. Population totals can increase through Houses and Barracks Soldier grants, Peasant totals can decrease through Barracks conversion, and available counts can decrease through economic buildings and tower construction.

## Core Workflows

### Repository Bootstrap

A contributor opens the repository, reads the root control documents, and sees that `td` is a local Go/Ebitengine PC game prototype. The bootstrap workflow has already produced project-specific control documents and the first runnable menu shell.

### Main Menu Workflow

A contributor runs `go run ./cmd/td`, starts a new game, and sees a tower-free home Plot containing only the Sanctum. The top bar shows 100 Wood, 50 Stone, 20 Metal, and `0/0` for every inhabitant role. The left building bar defaults to the `Housing` tab and shows values to the right of each icon. The House icon is available immediately with a green outline; dragging it onto an empty grass-like Tile spends 20 Wood and changes Peasants to `2/2`. Barracks is shown with a red outline and 70% icon opacity until enough Peasants exist. Switching to `Economic` exposes Woodcutter, Stone Quarry, and Iron Mine. A Woodcutter can then be placed for 10 Wood, reserving one Peasant as `1/2`; after the next defeated Raid, Wood increases by 10 from that Woodcutter. Stone Quarry and Iron Mine follow the same pattern for Stone and Metal with their displayed costs. The Barracks icon becomes eligible once two Peasants are available; dragging it onto an empty grass-like Tile spends 10 Wood and 10 Stone, changes Peasants to `0/0`, and changes Soldiers to `2/2`. Switching to `Defenses` exposes Bow, Flame Bolt, and Catapult Towers. Although starting resources cover Bow and Flame Bolt, neither tower icon highlights or starts dragging until required staff exist. Successful staffed placement reduces the relevant available count and leaves its total unchanged.

The top bar groups Apprentice, Soldier, and Peasant icons after physical resources. Each shows `available/total`; available inhabitants can be reserved by economic buildings or tower construction, while total remains unchanged.

## Product Constraints and Known Limits

- The current target is a local prototype only.
- Distribution, release packaging, CI, license selection, and store targets are deferred.
- The current playable shell includes a Sanctum-only starting Plot and a tabbed building bar whose resource, population, and staffing requirements jointly gate placement. The same building bar can place Houses, which cost Wood and grant Peasants; Barracks, which cost Wood and Stone and convert Peasants into Soldiers; and economic buildings, which reserve Peasants and produce resources after defeated Raids. The bar shows values beside icons, uses green outlines for buildable icon squares, and uses red outlines plus 70% icon opacity for capacity-blocked buildings.
- Saving the game and campaign structure are explicit non-goals for the first phase.
- Starting a new game opens the static home Plot with `0/0` inhabitants. House construction can add Peasants, Barracks construction can convert those Peasants into Soldiers, and economic buildings can turn reserved Peasants into post-Raid Wood, Stone, or Metal income. There is still no way to add Apprentices in normal play. Staff committed to economic buildings or towers cannot currently be released because structure removal and reassignment are not implemented.
- Settings are represented only by a placeholder screen; no configurable options exist yet.
- Exploration, broader base-building, resource gathering, and broader tower-defense gameplay are intended but not implemented.

## Non-Goals

The current phase is not trying to build campaign progression, save/load, multiplayer, online services, release packaging, or a full tower-defense encounter. The current Raid slice is a deterministic enemy placeholder with skeletons, two first-Raid zombies, player-built towers, small combat-defeat resource rewards, and no multiple paths. The project is also not trying to choose final art production workflows before the game loop proves useful.

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
