# Architecture: td

`ARCHITECTURE.md` helps contributors answer where code belongs and which boundaries should stay intact while `td` grows from a local prototype into a playable PC game.

The repository does not yet contain runtime code. The first implementation plan will initialize Go and Ebitengine.

## System Overview

`td` is a local PC tower-defense game prototype. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The planned codebase should be organized around a small Ebitengine executable in `cmd/td/` and reusable internal packages under `internal/`. Early code should keep menu state, rendering, input handling, and game-loop wiring easy to understand before larger gameplay systems are introduced.

## Codemap

- `cmd/td/` will own the executable entry point, Ebitengine window setup, and process startup.
- `internal/menu/` may own main-menu state, button hit testing, and menu actions once the first screen has enough behavior to justify a package.
- `internal/game/` may later own top-level game state and transitions between menu, exploration, base-building, and defense scenes.
- `internal/render/` may later own shared drawing helpers when rendering code becomes reusable.
- `assets/` will store static images, fonts, audio, and other runtime assets once real assets exist.
- `plans/` stores ordered ExecPlans. `plans/00-initial-ebitengine-menu.md` is the first implementation plan.
- `.agents/skills/` stores repo-local agent workflows.
- `.codex/config.toml` stores Codex defaults only; it is not application configuration.

Do not create packages before they have a clear responsibility. A single small `cmd/td/main.go` is acceptable for the first menu slice if it remains easy to read and testable behavior is factored out when needed.

## Main Flows

### First Planned App Flow

1. A contributor runs `go run ./cmd/td` from the repository root after the first implementation plan is complete.
2. `cmd/td` configures an Ebitengine window and starts the game loop.
3. Ebitengine calls `Update` for input and state changes, `Draw` for rendering, and `Layout` for logical screen sizing.
4. The main menu renders a title and a quit option.
5. When the user activates quit, `Update` returns Ebitengine's termination signal so the desktop app closes cleanly.

### Future Gameplay Flow

1. A player starts from the menu.
2. The game enters an exploration or base-management scene.
3. The player gathers resources and builds or upgrades defenses.
4. A tower-defense encounter applies enemy movement, tower targeting, damage, resource changes, and win or loss conditions.
5. The game returns results to the player through the UI and later progression systems.

This future flow is roadmap intent, not current behavior.

## Architectural Invariants

- Keep Ebitengine process startup in `cmd/td/`.
- Keep reusable game logic in `internal/` packages when it outgrows the entry point.
- Keep pure state transitions and hit testing testable without opening a graphics window.
- Do not let rendering helpers own gameplay rules.
- Do not introduce save, campaign, networking, or distribution architecture during the first menu slice.
- Keep `.codex/config.toml` limited to agent configuration.

## Boundaries and External Dependencies

The first external runtime dependency will be Ebitengine through `github.com/hajimehoshi/ebiten/v2`. It owns the desktop window, game loop, drawing surface, and input APIs. Game code should treat Ebitengine callbacks as the boundary between OS/window events and project-owned state.

Go module files will be introduced by the first ExecPlan. Until then, there are no application dependencies.

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
- What base resolution should become the long-term rendering target?
