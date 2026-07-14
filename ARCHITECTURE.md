# Architecture: td

`ARCHITECTURE.md` helps contributors answer where code belongs and which boundaries should stay intact while `td` grows from a local prototype into a playable PC game.

The repository contains an early runtime shell: a Go module, a small Ebitengine executable, an asset catalog package that loads runtime sprites, HUD icons, and audio bytes, a small runtime sound manager, a menu package that owns the current menu flow, a UI package that owns shared presentation helpers plus the selected-object panel and building-bar presentation, and a game package that owns the first logical game state, prototype explored-Plot map state, camera, selection, construction eligibility and placement, phase flow, population and resources, exploration, raids, combat, and in-game overlay behavior.

## System Overview

`td` is a local PC tower-defense game prototype. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The codebase is organized around a small Ebitengine executable in `cmd/td/` and reusable internal packages under `internal/`. Early code should keep menu state, rendering, input handling, and game-loop wiring easy to understand before larger gameplay systems are introduced.

## Codemap

- `cmd/td/` owns the executable entry point, Ebitengine window setup, callback wiring, app-mode routing between menu and game, Ebitengine input polling, runtime sound manager construction, quit termination handling, surrender-to-menu handling, pixel-sized Ebitengine layout, and process startup.
- `assets/` owns static runtime asset files and the typed asset catalog package. The catalog embeds required files, groups loaded assets by type and subtype, returns Ebitengine-ready images for game rendering, and returns raw audio bytes for runtime sound playback.
- `internal/menu/` owns menu screen state, menu rendering, resizable menu geometry, button hit testing, disabled-target handling, action selection, Wizard name input, the New Game configuration screen, and placeholder menu screens.
- `internal/game/` owns game status, the Labour-Management-Raid lifecycle, Wood, Stone, Iron, and Gold counts, enemy Gold drops, structure templates, populations, staffing reservations, terrain-consuming Labour production, Market trades, selected-object state, construction, combat, rendering, and input behavior. No separate assignment subsystem exists yet.
- `internal/sound/` owns Ebitengine audio context creation, WAV decoding for one-shot effects, active audio players, and effect playback. It is runtime-facing; gameplay rules should depend only on the game package's sound sink interface.
- `internal/ui/` owns shared UI palette colors, font-size constants, text drawing helpers, selected-object panel presentation, building-bar presentation, and the Market control component's fixed-pixel layout, edge fallback, hit testing, availability styling, and drawing. It receives presentation-neutral facts and does not mutate gameplay resources.
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
3. The game package renders explored 15x15 Plots from its stored prototype map data. A new game starts in Management with only the home Plot explored. The home Plot uses ordinary grasslands generation at 6% Tree, 3% Boulder, 1% Iron Deposit, and 90% empty grass, then overwrites generated terrain with the centered Sanctum's straight road north to the Plot edge; the Sanctum is its only initial structure. When an unexplored Plot first becomes orthogonally adjacent to explored land, the map independently assigns and stores grasslands or hills with equal probability. Magnifying-glass buttons appear on those borders during Management with the assigned biome name rendered outward into unexplored space; only the circular button is clickable. Clicking it generates the already assigned biome with mostly buildable grass plus sparse Tree, Boulder, and Iron Deposit Tiles. Hills uses 3% Tree, 6% Boulder, 3% Iron Deposit, and 88% empty grass. Adjacent explored Plots have no plot-level frame, gutter, or padding between them; shared-edge cleanup can clear generated home terrain, and north-chain Plots also continue the straight center road. All combat towers must be built by the player.
4. `cmd/td` polls mouse-wheel input, held `W`, `A`, `S`, and `D` keys, and right mouse button press, hold, and release state, then passes those values to `internal/game`.
5. The game package updates a private camera for map inspection. Mouse-wheel input changes zoom around the scene viewport center, `WASD` changes the camera center, and right-drag panning grabs the visible world when the drag starts over the game view rather than screen-space UI. The camera has a tiny minimum zoom for technical safety but no maximum zoom and no pan bounds.
6. Left-clicking prioritizes a visible raider, then a structure Tile, then Tree, Boulder, or Iron Deposit terrain. Empty grass, Road, and other clicks clear selection. The selected subject is stored in `internal/game.State`; raiders and structures render brighter, terrain receives a gold Tile outline, and supported subjects are adapted into UI-facing detail data for a bottom-right panel. A selected Market also anchors three fixed-size trade buttons beside its projected Tile during Management. Selection panels and Market buttons own their clicks, selection works while SPACE-paused, and the overlay blocks selection and trading.
7. The game package renders Wood, Stone, Iron, Gold, available/total populations, and centered phase-plus-challenge status. The Management-only building bar groups `Housing`, `Economic`, and `Defenses`; Market follows Iron Mine in Economic. Eligibility requires sufficient resources and every required role. Placement applies cost, population effects, and staff reservation atomically. Producers reserve one Peasant and consume matching terrain during later Labour. Market costs 30 Wood, reserves two Peasants and one Soldier, performs no Labour production, and buys one Wood or Stone for 1 Gold or one Iron for 3 Gold. A new game starts with 0 Gold and `0/0` for every role, so House is initially eligible while staffed structures are not.
8. Left-dragging an eligible building-bar icon during Management starts a private build drag. A half-sized copy of the icon follows the cursor. Releasing over an empty `terrainEmpty` Tile with no feature in any explored Plot places the matching structure feature and applies its construction cost, population cost, staff reservation, and population grant; releasing over occupied Tiles, roads, Trees, Boulders, Iron Deposits, unexplored space, during Labour, active Raids, or after breach cancels without spending resources or population effects. SPACE-paused Management still permits placement, while the in-game overlay blocks building input. Economic buildings remain placeable only on empty grass; their required terrain may be anywhere in the explored Domain.
9. While unpaused, each Ebitengine update advances the logical update counter. Active Raids advance deterministic progress, enemy spawning, tower combat, and north-road movement. Skeletons, Zombies, Ghouls, and Armoured Skeletons release at progress-score multiples of 2, 4, 6, and 8. Their templates own movement, health, sprite, and `GoldDrop`; tower-damage defeats grant 1, 2, 3, or 5 Gold respectively and report the defeat sound. No Wood, Stone, or Iron comes from combat. Enemies removed by Barricade or breach grant nothing. Successful Raid completion resolves terrain-consuming Labour in deterministic map order before Management opens.
10. Clicking `Next Raid` during unpaused Management starts the next deterministic scaling Raid immediately at zero progress. The button is disabled during Labour, pause, an active Raid, and after breach.
11. When the user presses SPACE, `cmd/td` passes pause input to `internal/game`, which toggles pause without incrementing the counter on that frame.
12. While paused, the game renders a `PAUSED` label and does not increment the logical update counter or advance Raid or combat simulation, but camera input, selection, Management exploration, construction, and Market trades still update. Pause disables `Next Raid`.
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
- Keep explored Plot storage, unexplored frontier-biome assignments, and reveal orchestration in `internal/game/map.go`, biome-backed Plot generation in `internal/game/plot_generator.go`, and the first camera/projection behavior, including right-drag camera state, in `internal/game/camera.go` until richer multi-map behavior or multiple gameplay scenes create a reason for separate packages.
- Keep gameplay positions in Sanctum-centered Tile-unit world coordinates inside `internal/game`: the Sanctum center is `(0, 0)`, one Tile is one unit, and positive Y points north. Rendering code should convert this model to Ebitengine screen coordinates at the projection boundary.
- Keep in-game overlay behavior inside `internal/game` while it is tightly coupled to game pause state and game rendering.
- Keep pure state transitions, hit testing, object selection, and simple menu text input testable without opening a graphics window.
- Keep the current menu and game transition explicit until there are enough real non-menu screens to justify a shared scene abstraction.
- Do not let rendering helpers own gameplay rules. UI helpers may format and draw presentation-neutral data supplied by game systems, but should not inspect maps, raids, structures, or other gameplay state directly.
- Do not introduce save, campaign, networking, or distribution architecture during the first menu slice.
- Keep `.codex/config.toml` limited to agent configuration.

## Boundaries and External Dependencies

The first external runtime dependency will be Ebitengine through `github.com/hajimehoshi/ebiten/v2`. It owns the desktop window, game loop, drawing surface, and input APIs. Game code should treat Ebitengine callbacks as the boundary between OS/window events and project-owned state.

Go module files are checked in. The current runtime dependency is Ebitengine and the Go support libraries required by the module.

## Cross-Cutting Concerns

### Configuration

There is no application configuration system yet. If configuration becomes necessary, prefer explicit Go constants for prototype-only values before adding config files.

### Assets

The runtime sprites include exact 64x64 structure PNGs through Market, enemy PNGs for the four-raider roster, terrain families, and Wood, Stone, Iron, Gold, Apprentice, Soldier, and Peasant icons. `assets/catalog.go` embeds and exposes them through typed groups. The first runtime audio asset remains the raider-defeated WAV. Gameplay rules do not decode files or know asset paths.

### Testing

Use `go test ./...` after the Go module exists. Prefer tests for pure behavior such as button hit testing, menu action selection, state transitions, map rules, and combat calculations.

### Accessibility and Usability

Menus should have readable text, clear interaction states, and stable targets. Keyboard navigation should be considered once the menu has more than the first quit action.

## How To Extend The System Safely

To add the first executable app, follow `plans/00-initial-ebitengine-menu.md` instead of improvising from chat history.

To add a new screen later, keep the transition logic explicit and avoid building a large scene framework before there are at least two or three real screens with shared needs.

To add gameplay systems, start with pure state and geometry that `go test ./...` can exercise. `internal/game/resources.go` owns resource mutations and Labour production, `internal/game/market.go` owns trade eligibility and economy mutations, and `internal/game/combat.go` applies template Gold drops. `internal/ui/market_controls.go` owns only Market layout, hit testing, styling, and drawing. A future removal or reassignment system must return reserved staff without duplicating structure-template requirements.

To add assets, place files under `assets/`, document source and licensing, add them to the typed catalog only when game code needs them, and avoid mixing asset-loading details into gameplay rules.

## Open Questions

- What package boundaries will be useful after the first menu screen exists?
- Should the project use a custom scene manager, or keep explicit state transitions until repetition appears?
- When should the prototype camera be split from `internal/game` into reusable rendering or scene infrastructure?
