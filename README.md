# td

`td` is a local PC prototype for a 2D tower-defense game built with Go and Ebitengine. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository now contains a small runnable prototype shell: a Go/Ebitengine desktop app with a main menu, a New Game configuration screen for entering a Wizard name, a first game screen with a static 15x15 home Plot scene, a sprite-backed Sanctum, sprite-backed Bow and Flame Bolt Towers across the road from each other, mouse-wheel map zoom, `WASD` camera panning, click selection and a bottom-right selection panel for structures and raiders, a prototype status top bar with sprite-backed resource icons, a visual-only left building bar with Bow Tower and Flame Bolt Tower icons, construction costs, and affordable-option hover highlighting, deterministic placeholder Raids with sprite-backed skeleton and zombie enemies and visible enemy health bars started by a `Next Raid` button, first-pass tower projectile combat with a short sound when a raider is defeated, a logical update counter, an ESC in-game overlay menu, a placeholder Settings screen, disabled Load option, SPACE pause behavior, and a quit option. Exploration, base-building, tower placement, broader tower-defense systems, real resource changes, real settings, save/load behavior, and an asset pipeline have not been implemented yet.

## Current Status

- Stage: local runnable prototype foundation.
- Runtime stack: Go with Ebitengine.
- Current playable slice: a 1920x1080 desktop app that shows `New`, `Load`, `Settings`, and `Quit`; `New` opens a configuration screen with Wizard name entry up to 32 characters, active `Start` after a name is entered, and active `Cancel`; `Start` opens the first game screen with a static 15x15 home Plot containing the centered sprite-backed Sanctum, a straight road north to the Plot edge, sprite-backed Bow and Flame Bolt Towers across the road from each other, and a pine-tree border around the Plot edge except at the road opening; mouse-wheel input zooms the map, `WASD` pans the camera, and these camera controls keep working while paused; left-clicking structures or visible raiders selects them and draws the selected object brighter, while clicking elsewhere clears selection; selected raiders and combat towers show a bottom-right information panel with their current prototype stats, and the Sanctum shows a basic name panel; a top bar shows fixed prototype Chapter, Day, resource counts with Wood, Stone, and Metal icons, phase, and Sanctum barricade values, and a left building bar below the top bar shows visual-only Bow Tower and Flame Bolt Tower build icons with colour-coded construction costs that become bolder only while hovering the matching icon and only when current resources cover that tower's cost; clicking `Next Raid` starts the first Raid with skeleton, zombie, skeleton, zombie, skeleton spawns and health bars above their sprites, and the two towers fire arrow and flame-bolt projectiles at enemies within range; combat defeats play a short embedded sound effect; SPACE toggles pause; ESC opens an in-game overlay menu with `Resume` and `Surrender` and blocks camera and selection input while open; `Settings` opens a placeholder screen with `Back`; `Load` is disabled; and `Quit` closes the app.
- Current display policy: the window is resizable, and the drawable layout follows the actual window size so text remains raw-pixel-sized instead of being stretched during upscaling.
- Current non-goals: exploration, base-building, full tower-defense combat, saving games, and campaign structure.
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
