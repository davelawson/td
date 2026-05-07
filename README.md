# td

`td` is a local PC prototype for a 2D tower-defense game built with Go and Ebitengine. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository now contains a small runnable prototype shell: a Go/Ebitengine desktop app with a main menu, placeholder New Game and Settings screens, a disabled Load option, and a quit option. Gameplay systems, real settings, save/load behavior, and static assets have not been implemented yet.

## Current Status

- Stage: local runnable prototype foundation.
- Runtime stack: Go with Ebitengine.
- Current playable slice: a desktop app that shows `New`, `Load`, `Settings`, and `Quit`; `New` and `Settings` open placeholder screens with `Back` buttons, `Load` is disabled, and `Quit` closes the app.
- Current non-goals: saving games and campaign structure.
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
- `CODESTYLE.md` defines Go-oriented source conventions, commenting requirements, and file-size expectations.
- `cmd/td/` contains the Ebitengine executable entry point.
- `DESIGN.md` records the medieval wizardry design direction and UI review expectations.
- `go.mod` and `go.sum` define the Go module and runtime dependencies.
- `internal/menu/` contains testable menu hit-testing and action-selection behavior.
- `PLANS.md` defines how ExecPlans are written and maintained.
- `PRODUCT.md` records current user-visible product truth.
- `ROADMAP.md` records intended product direction and explicit non-priorities.
- `plans/` stores ordered ExecPlans.
- `.agents/skills/` stores repo-local agent workflows.
- `.codex/config.toml` stores project-scoped Codex defaults.

Static game assets should eventually live in `assets/`, and tests should mirror the package layout under `tests/` or live beside Go packages when idiomatic package-level tests are clearer.

## Development Notes

Do not implement product feature code during bootstrap work. Start substantial changes by reading the control documents and then creating or updating an ordered ExecPlan in `plans/`.

When adding Go code, follow `CODESTYLE.md`: keep functions focused, document every function or method with Go doc comments, prefer descriptive names, and check hand-written code file line counts at the end of substantial work.
