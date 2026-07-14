# td

`td` is a local PC prototype for a 2D tower-defense game built with Go and Ebitengine. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository now contains a small runnable prototype shell: a Go/Ebitengine desktop app with a main menu, Wizard name entry, a Sanctum-only starting Plot, camera inspection through mouse-wheel zoom, `WASD` panning, and right-drag panning, object and natural-obstacle selection, resource and population status, and a 260-pixel tabbed building bar with each building's values shown to the right of its icon. Peaceful play is split into an instantaneous Labour phase after each successful Raid and an open-ended Management phase. During Management the player can explore, construct, and trade until clicking `Next Raid`; during Labour placed producers consume terrain to make resources. The bar defaults to `Housing` for House, Barracks, and Dorm, with `Economic` for Woodcutter, Stone Quarry, Iron Mine, and Market, and `Defenses` for Bow, Flame Bolt, and Catapult Towers. A Market costs 30 Wood, reserves two Peasants and one Soldier, and exposes contextual buttons beside the selected building that buy one Wood or Stone for 1 Gold and one Iron for 3 Gold. These buttons remain available during paused Management and hide in other phases. The top bar shows Wood, Stone, Iron, and Gold. Tower-damage defeats grant only deterministic Gold—1 from a Skeleton, 2 from a Zombie, 3 from a Ghoul, and 5 from an Armoured Skeleton—while Barricade removal and breach grant nothing. Terrain-consuming Labour production, staffing-aware construction, deterministic scaling Raids, projectile combat, pause, overlay, and quit behavior remain implemented. Timed recruitment, worker reassignment, staff release, bulk trading, resource caps, and terrain regrowth are not implemented.

## Current Status

- Stage: local runnable prototype foundation.
- Runtime stack: Go with Ebitengine.
- Current playable slice: a new game starts in Management with only the Sanctum, 100 Wood, 50 Stone, 20 Iron, 0 Gold, and `0/0` for Apprentices, Soldiers, and Peasants. Housing creates the roles needed to staff producers, Markets, and defenses. Woodcutter, Stone Quarry, and Iron Mine consume matching terrain during Labour. Market costs 30 Wood, requires one Soldier and two Peasants, and trades Gold for individual materials through buttons beside a selected Market. Bow, Flame Bolt, and Catapult defenses fight deterministic Skeleton, Zombie, Ghoul, and Armoured Skeleton Raids; combat defeats award only the enemy's tiered Gold drop. Existing camera, exploration, inspection, construction, pause, overlay, and menu behavior remains available.
- Current display policy: the window is resizable, and the drawable layout follows the actual window size so text remains raw-pixel-sized instead of being stretched during upscaling.
- Current non-goals: exploration, timed or manual resource gathering, timed inhabitant recruitment and assignment, worker reassignment, staff release, tower upgrades, broader base-building, full tower-defense combat, saving games, and campaign structure.
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
- `internal/game/` contains testable game state, prototype map data, camera state, selected-object state and detail data, visual building-bar behavior, scene projection and rendering, asset-catalog ownership for active games, and logical update behavior.
- `internal/ui/` contains shared UI colors, font sizes, text helpers, widgets, and generic selected-object panel presentation.
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
