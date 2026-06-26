# Architecture: td

`ARCHITECTURE.md` helps contributors answer where code belongs and which boundaries should stay intact while `td` grows from a local prototype into a playable PC game.

The repository contains an early runtime shell: a Go module, a small Ebitengine executable, an asset catalog package that loads runtime sprites, HUD icons, and audio bytes, a small runtime sound manager, a menu package that owns the current menu flow, and a game package that owns the first logical game state, prototype map state, Sanctum-centered world coordinates, camera state, selected structure and raider state, selected-object detail panels, building bar with construction-cost and population metadata display, affordable-option icon hover feedback, calm-phase structure drag placement, House Peasant grants, home Plot projection and rendering, deterministic placeholder Raid state with sprite-backed skeleton and zombie enemies, first-pass tower projectile combat including Catapult Tile-area impact, raider-defeat sound events, combat-defeat resource rewards, and in-game overlay menu.

## System Overview

`td` is a local PC tower-defense game prototype. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The codebase is organized around a small Ebitengine executable in `cmd/td/` and reusable internal packages under `internal/`. Early code should keep menu state, rendering, input handling, and game-loop wiring easy to understand before larger gameplay systems are introduced.

## Codemap

- `cmd/td/` owns the executable entry point, Ebitengine window setup, callback wiring, app-mode routing between menu and game, Ebitengine input polling, runtime sound manager construction, quit termination handling, surrender-to-menu handling, pixel-sized Ebitengine layout, and process startup.
- `assets/` owns static runtime asset files and the typed asset catalog package. The catalog embeds required files, groups loaded assets by type and subtype, returns Ebitengine-ready images for game rendering, and returns raw audio bytes for runtime sound playback.
- `internal/menu/` owns menu screen state, menu rendering, resizable menu geometry, button hit testing, disabled-target handling, action selection, Wizard name input, the New Game configuration screen, and placeholder menu screens.
- `internal/game/` owns game status, structure templates, construction resources, available/total populations, staffing checks and reservations, House population grants, placement, combat, rendering, and input behavior. Staffing reservation and population grants are part of successful construction; no separate assignment subsystem exists yet.
- `internal/sound/` owns Ebitengine audio context creation, WAV decoding for one-shot effects, active audio players, and effect playback. It is runtime-facing; gameplay rules should depend only on the game package's sound sink interface.
- `internal/ui/` owns shared UI palette colors used by menu and game rendering. It should remain palette-only until repeated UI behavior justifies more shared code.
- `internal/render/` may later own shared drawing helpers when rendering code becomes reusable.
- `assets/` stores static images, audio effect files, and the first typed runtime asset catalog. It may later grow to include fonts and other runtime assets.
- `plans/` stores ordered ExecPlans. `plans/00-initial-ebitengine-menu.md` is the first implementation plan.
- `.agents/skills/` stores repo-local agent workflows.
- `.codex/config.toml` stores Codex defaults only; it is not application configuration.

Do not create packages before they have a clear responsibility. `internal/menu/` exists because the current menu now has enough state, rendering, and testable behavior to justify a menu-owned package. Do not turn it into a general scene framework before real gameplay screens create repeated needs.

## Main Flows

### Current Menu Flow

1. A contributor runs `go run ./cmd/td` from the repository root.
2. `cmd/td` configures a 1920x1080 Ebitengine window and starts the game loop.
3. Ebitengine calls `Update` for input and state changes, `Draw` for rendering, and `Layout` for drawable sizing. `Layout` follows the current window size so resizing does not stretch text as part of a fixed framebuffer.
4. `cmd/td` forwards pointer and keyboard input state to `internal/menu`.
5. The menu package renders `New`, disabled `Load`, `Settings`, and `Quit`.
6. When the user activates `New`, the menu switches to a New Game configuration screen with a focused Wizard name field, disabled `Start` button, and active `Cancel` button.
7. When the user types on the New Game screen, the focused Wizard name field updates up to 32 characters; Backspace removes the last typed character.
8. When the Wizard name is non-empty, the New Game `Start` button becomes active.
9. When the user activates `Start`, `cmd/td` constructs `internal/game` state with the Wizard name and switches from menu mode to game mode.
10. When the user activates `Settings`, the menu switches to a placeholder Settings screen with a `Back` button.
11. When the user activates `Cancel` on New Game or `Back` on Settings, the menu returns to the main menu.
12. When the user activates `Quit`, the menu reports a quit action and `cmd/td` returns Ebitengine's termination signal so the desktop app closes cleanly.

### Current Game Flow

1. A contributor starts a game from the New Game screen after entering a Wizard name.
2. `cmd/td` routes Ebitengine updates and drawing to `internal/game`.
3. The game package renders a static 15x15 home Plot from its stored prototype map data. The Plot contains the centered Sanctum as its only initial structure, a straight road north to the Plot edge, and pine trees around the Plot edge except where the road exits. All combat towers must be built by the player.
4. `cmd/td` polls mouse-wheel input and held `W`, `A`, `S`, and `D` keys, then passes those values to `internal/game`.
5. The game package updates a private camera for map inspection. Mouse-wheel input changes zoom around the scene viewport center, and `WASD` changes the camera center. The camera has a tiny minimum zoom for technical safety but no maximum zoom and no pan bounds.
6. Left-clicking a visible raider selects that raider; otherwise, left-clicking a structure tile selects that structure; clicking elsewhere clears selection. The selected object is stored in `internal/game.State`, rendered brighter, and shown in a bottom-right detail panel when supported. Panel clicks are treated as UI input and do not clear or pass through selection. Selection works while SPACE-paused, and the in-game overlay blocks selection changes.
7. The game package renders resources and available/total populations. Building-bar eligibility requires sufficient resources and all required available roles. Successful tower placement deducts resources and reserves staff atomically; successful House placement deducts Wood and grants Peasants by increasing both available and total counts. A new game starts at `0/0` for every role, so House is initially eligible but no tower is initially eligible.
8. Left-dragging an eligible building-bar icon during calm play starts a private build drag. A half-sized copy of the icon follows the cursor. Releasing over an empty `terrainEmpty` Tile with no feature places the matching structure feature and applies its construction cost, staff reservation, and population grant; releasing over occupied Tiles, roads, forests, during active Raids, or after breach cancels without spending resources or population effects. SPACE-paused calm play still permits placement, while the in-game overlay blocks building input.
9. While unpaused, each Ebitengine update advances the logical update counter by one. If a Raid is active, the game package advances deterministic enemy spawning, tower targeting and projectile combat, and enemy movement along the fixed north road. Raid 1 spawns skeleton, zombie, skeleton, zombie, skeleton, while later placeholder Raids remain skeleton-only. Enemies store Sanctum-centered world coordinates in Tile units, read Tiles-per-second movement speed and defeat-resource rewards from their enemy template, have health, and render from their template sprite in the active asset catalog. Combat towers store timing stats in seconds, track cooldowns as seconds remaining, launch projectiles at in-range enemies, and damage enemies on projectile impact. Bow and Flame Bolt projectiles damage only their original target; Catapult projectiles damage every active enemy in the Tile occupied by the original target when the projectile lands. Tower-damage defeats grant the defeated raider's template resources and report a raider-defeated sound event to the runtime sound sink. Enemies that reach the Sanctum spend Barricade charges and are removed without rewards, or breach the Sanctum when no charges remain. Selected raider state is cleared when that raider is no longer active.
10. Clicking `Next Raid` while no Raid is active and the Sanctum is not breached starts the next deterministic placeholder Raid immediately. The button is disabled during an active Raid and after breach.
11. When the user presses SPACE, `cmd/td` passes pause input to `internal/game`, which toggles pause without incrementing the counter on that frame.
12. While paused, the game renders a `PAUSED` label and does not increment the logical update counter or advance Raid or combat simulation, but camera input, selection, and calm building still update so the map can be inspected and adjusted.
13. When the user presses ESC, `cmd/td` passes overlay-menu input to `internal/game`.
14. The game package opens a centered in-game menu, pauses the game, draws it over the still-visible game scene, darkens the rest of the scene by about 50%, and blocks camera, selection, building, and Raid input while the overlay remains open.
15. When the user presses ESC again or clicks `Resume`, the game package closes the overlay and restores the pause state from before the overlay opened.
16. When the user clicks `Surrender`, `internal/game` returns a surrender action to `cmd/td`, and `cmd/td` discards the active game state and returns to the top-level main menu.

### Future Gameplay Flow

1. A player starts from the menu.
2. The game enters an exploration or base-management scene.
3. The player gathers resources and builds or upgrades defenses.
4. A tower-defense encounter applies enemy movement, tower targeting, damage, resource changes, and win or loss conditions.
5. The game returns results to the player through the UI and later progression systems.

This future flow is roadmap intent, not current behavior.

## Architectural Invariants

- Keep Ebitengine process startup in `cmd/td/`.
- Keep app-mode routing in `cmd/td/`; reusable game state and rules belong in `internal/game`.
- Keep runtime asset loading in `assets/` so gameplay rules do not decode files directly.
- Keep Ebitengine audio context ownership in `internal/sound` and app startup; gameplay rules may report sound events but should not create players or decode audio.
- Keep the current display policy as a pixel-sized drawable layout: the initial window is 1920x1080, resizes update menu geometry, and text remains raw-pixel-sized rather than stretched by framebuffer scaling.
- Keep reusable game logic in `internal/` packages when it outgrows the entry point.
- Keep the first map data in `internal/game/map.go` and the first camera/projection behavior in `internal/game/camera.go` until real exploration, generation, multi-map behavior, or multiple gameplay scenes create a reason for separate packages.
- Keep gameplay positions in Sanctum-centered Tile-unit world coordinates inside `internal/game`: the Sanctum center is `(0, 0)`, one Tile is one unit, and positive Y points north. Rendering code should convert this model to Ebitengine screen coordinates at the projection boundary.
- Keep in-game overlay behavior inside `internal/game` while it is tightly coupled to game pause state and game rendering.
- Keep pure state transitions, hit testing, object selection, and simple menu text input testable without opening a graphics window.
- Keep the current menu and game transition explicit until there are enough real non-menu screens to justify a shared scene abstraction.
- Do not let rendering helpers own gameplay rules.
- Do not introduce save, campaign, networking, or distribution architecture during the first menu slice.
- Keep `.codex/config.toml` limited to agent configuration.

## Boundaries and External Dependencies

The first external runtime dependency will be Ebitengine through `github.com/hajimehoshi/ebiten/v2`. It owns the desktop window, game loop, drawing surface, and input APIs. Game code should treat Ebitengine callbacks as the boundary between OS/window events and project-owned state.

Go module files are checked in. The current runtime dependency is Ebitengine and the Go support libraries required by the module.

## Cross-Cutting Concerns

### Configuration

There is no application configuration system yet. If configuration becomes necessary, prefer explicit Go constants for prototype-only values before adding config files.

### Assets

The first runtime sprites are the 64x64 Sanctum, House, Bow Tower, Flame Bolt Tower, Catapult Tower, and tower projectile PNGs under `assets/sprites/structures/`, the skeleton and zombie enemy PNGs under `assets/sprites/enemies/`, pine tree PNGs under `assets/sprites/terrains/`, and Wood, Stone, Metal, Apprentice, Soldier, and Peasant HUD icon PNGs under `assets/sprites/icons/`. The first runtime audio asset is the raider-defeated WAV under `assets/audio/`. The `assets` package embeds and loads required runtime assets into a typed catalog. Early prototypes may still draw simple shapes and text directly when no sprite exists, but gameplay rules should not decode files or know asset paths.

### Testing

Use `go test ./...` after the Go module exists. Prefer tests for pure behavior such as button hit testing, menu action selection, state transitions, map rules, and combat calculations.

### Accessibility and Usability

Menus should have readable text, clear interaction states, and stable targets. Keyboard navigation should be considered once the menu has more than the first quit action.

## How To Extend The System Safely

To add the first executable app, follow `plans/00-initial-ebitengine-menu.md` instead of improvising from chat history.

To add a new screen later, keep the transition logic explicit and avoid building a large scene framework before there are at least two or three real screens with shared needs.

To add gameplay systems, start with pure data and functions that can be tested by `go test ./...`. `internal/game/structures.go` owns staffing requirements and population-grant template metadata, `internal/game/hud.go` owns population availability, reservation, and grants, and `internal/game/building_bar.go` applies resource, staffing, and population-grant effects at drag start and release. A future removal or reassignment system must return reserved staff without duplicating template requirements, and must define how removing a population-provider affects inhabitants before changing House behavior.

To add assets, place files under `assets/`, document source and licensing, add them to the typed catalog only when game code needs them, and avoid mixing asset-loading details into gameplay rules.

## Open Questions

- What package boundaries will be useful after the first menu screen exists?
- Should the project use a custom scene manager, or keep explicit state transitions until repetition appears?
- When should the prototype camera be split from `internal/game` into reusable rendering or scene infrastructure?
