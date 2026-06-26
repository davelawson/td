# td

`td` is a local PC prototype for a 2D tower-defense game built with Go and Ebitengine. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository now contains a small runnable prototype shell: a Go/Ebitengine desktop app with a main menu, Wizard name entry, a Sanctum-only starting Plot, camera inspection, object selection, resource and population status, and a building bar for Bow, Flame Bolt, and Catapult Towers. Tower construction requires both the displayed resources and available staff; successful construction reserves staff by reducing available counts while totals remain unchanged. Deterministic Raids, projectile combat, Catapult Tile-area damage, combat rewards, pause, an in-game overlay, and quit behavior are also implemented. Recruitment and staff release are not implemented.

## Current Status

- Stage: local runnable prototype foundation.
- Runtime stack: Go with Ebitengine.
- Current playable slice: a new game starts with only the Sanctum, 100 Wood, 50 Stone, 20 Metal, and `0/0` for Apprentices, Soldiers, and Peasants. Bow requires one available Soldier, Flame Bolt requires one available Apprentice, and Catapult requires one available Soldier plus two available Peasants. An icon highlights and starts dragging only when both resources and staff are sufficient. Successful placement deducts resources and reduces the required roles' available counts without changing totals. Because every population currently starts at zero and recruitment is absent, no tower can initially be built. Existing camera, selection, Raid, combat, pause, overlay, and menu behavior remains available.
- Current display policy: the window is resizable, and the drawable layout follows the actual window size so text remains raw-pixel-sized instead of being stretched during upscaling.
- Current non-goals: exploration, resource gathering, inhabitant recruitment and assignment, tower upgrades, broader base-building, full tower-defense combat, saving games, and campaign structure.
- Repository operations such as license, CI, release packaging, and distribution are intentionally deferred.

## Commands

Use these commands from the repository root:

- `go test ./...` runs all tests.
- `go run ./cmd/td` starts the local prototype.
- `go mod tidy` reconciles Go module dependencies after dependency or import changes.
- `git diff --check`
- `git status --short`

## Repository Layout

- `AGENTS.md` defines repository-specific instructions for coding agents.
- `ARCHITECTURE.md` describes intended code ownership, boundaries, and extension points.
- `ART.md` records guidance for generated art assets, prompt patterns, and asset review.
- `assets/` contains the typed runtime asset catalog package and static sprite and audio files.
- `CODESTYLE.md` defines Go-oriented source conventions, commenting requirements, and file-size expectations.
- `cmd/td/` contains the Ebitengine executable entry point.
- `DESIGN.md` records the medieval wizardry design direction and UI review expectations.
- `GAME.md` records intended game design decisions and open gameplay questions regardless of implementation state.
- `go.mod` and `go.sum` define the Go module and runtime dependencies.
- `internal/game/` contains testable game state, prototype map data, camera state, selected-object state and detail-panel behavior, visual building-bar behavior, scene projection and rendering, asset-catalog ownership for active games, and logical update behavior.
- `internal/menu/` contains testable menu hit-testing, action-selection, screen-routing, and Wizard name input behavior.
- `PLANS.md` defines how ExecPlans are written and maintained.
- `PRODUCT.md` records current user-visible product truth.
- `ROADMAP.md` records intended product direction and explicit non-priorities.
- `plans/` stores ordered ExecPlans.
- `.agents/skills/` stores repo-local agent workflows.
- `.codex/config.toml` stores project-scoped Codex defaults.

Additional static game assets should live in `assets/`, and tests should mirror the package layout under `tests/` or live beside Go packages when idiomatic package-level tests are clearer.

## Development Notes

Do not implement product feature code during bootstrap work. Start substantial changes by reading the control documents and then creating or updating an ordered ExecPlan in `plans/`.

When adding Go code, follow `CODESTYLE.md`: keep functions focused, document every function or method with Go doc comments, prefer descriptive names, and check hand-written code file line counts at the end of substantial work.
