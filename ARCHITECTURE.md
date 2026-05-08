# Architecture: td

`ARCHITECTURE.md` helps contributors answer where code belongs and which boundaries should stay intact while `td` grows from a local prototype into a playable PC game.

The repository contains an early runtime shell: a Go module, a small Ebitengine executable, a menu package that owns the current menu flow, and a game package that owns the first logical game state and in-game overlay menu.

## System Overview

`td` is a local PC tower-defense game prototype. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The codebase is organized around a small Ebitengine executable in `cmd/td/` and reusable internal packages under `internal/`. Early code should keep menu state, rendering, input handling, and game-loop wiring easy to understand before larger gameplay systems are introduced.

## Codemap

- `cmd/td/` owns the executable entry point, Ebitengine window setup, callback wiring, app-mode routing between menu and game, Ebitengine input polling, quit termination handling, surrender-to-menu handling, pixel-sized Ebitengine layout, and process startup.
- `internal/menu/` owns menu screen state, menu rendering, resizable menu geometry, button hit testing, disabled-target handling, action selection, Wizard name input, the New Game configuration screen, and placeholder menu screens.
- `internal/game/` owns the first top-level game state, Wizard name storage, pause state, logical update counting, prototype top-bar status display, in-game overlay menu behavior, and placeholder game rendering. It may later grow into exploration, base-building, and defense scene state when those systems exist.
- `internal/ui/` owns shared UI palette colors used by menu and game rendering. It should remain palette-only until repeated UI behavior justifies more shared code.
- `internal/render/` may later own shared drawing helpers when rendering code becomes reusable.
- `assets/` will store static images, fonts, audio, and other runtime assets once real assets exist.
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
3. The game package renders a placeholder field, the Wizard name, a top bar with prototype Chapter, Day, resources, phase, and Sanctum barricade status, and a debug logical update counter.
4. While unpaused, each Ebitengine update advances the logical update counter by one.
5. When the user presses SPACE, `cmd/td` passes pause input to `internal/game`, which toggles pause without incrementing the counter on that frame.
6. While paused, the game renders a `PAUSED` label and does not increment the logical update counter.
7. When the user presses ESC, `cmd/td` passes overlay-menu input to `internal/game`.
8. The game package opens a centered in-game menu, pauses the game, draws it over the still-visible game scene, and darkens the rest of the scene by about 50%.
9. When the user presses ESC again or clicks `Resume`, the game package closes the overlay and restores the pause state from before the overlay opened.
10. When the user clicks `Surrender`, `internal/game` returns a surrender action to `cmd/td`, and `cmd/td` discards the active game state and returns to the top-level main menu.

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
- Keep the current display policy as a pixel-sized drawable layout: the initial window is 1920x1080, resizes update menu geometry, and text remains raw-pixel-sized rather than stretched by framebuffer scaling.
- Keep reusable game logic in `internal/` packages when it outgrows the entry point.
- Keep in-game overlay behavior inside `internal/game` while it is tightly coupled to game pause state and game rendering.
- Keep pure state transitions, hit testing, and simple menu text input testable without opening a graphics window.
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

There are no real assets yet. Early prototypes may draw shapes and text directly. When assets arrive, store them under `assets/` and keep loading code separated from gameplay rules.

### Testing

Use `go test ./...` after the Go module exists. Prefer tests for pure behavior such as button hit testing, menu action selection, state transitions, map rules, and combat calculations.

### Accessibility and Usability

Menus should have readable text, clear interaction states, and stable targets. Keyboard navigation should be considered once the menu has more than the first quit action.

## How To Extend The System Safely

To add the first executable app, follow `plans/00-initial-ebitengine-menu.md` instead of improvising from chat history.

To add a new screen later, keep the transition logic explicit and avoid building a large scene framework before there are at least two or three real screens with shared needs.

To add gameplay systems, start with pure data and functions that can be tested by `go test ./...`, then connect them to Ebitengine rendering and input.

To add assets, place files under `assets/`, document source and licensing, and avoid mixing asset-loading details into gameplay rules.

## Open Questions

- What package boundaries will be useful after the first menu screen exists?
- Should the project use a custom scene manager, or keep explicit state transitions until repetition appears?
- Should later gameplay use the same pixel-sized layout policy as the current menus, or introduce a separate world camera?
