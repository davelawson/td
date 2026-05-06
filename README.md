# td

`td` is a local PC prototype for a 2D tower-defense game built with Go and Ebitengine. The intended game combines exploration, base-building, resource gathering, and conventional tower-defense combat in a medieval wizardry fantasy setting.

The repository is currently in bootstrap state. It contains project control documents and the first implementation plan, but it does not yet contain a Go module, game executable, source code, tests, or assets.

## Current Status

- Stage: local prototype planning.
- Runtime stack: Go with Ebitengine, to be initialized by the first ExecPlan.
- First playable slice: a desktop app that shows a main menu with a quit option; pressing quit closes the app.
- Current non-goals: saving games and campaign structure.
- Repository operations such as license, CI, release packaging, and distribution are intentionally deferred.

## Planned Commands

The Go toolchain is not initialized yet. After `plans/00-initial-ebitengine-menu.md` is executed, these will be the canonical local commands:

- `go test ./...` runs all tests.
- `go run ./cmd/td` starts the local prototype.
- `go mod tidy` reconciles Go module dependencies after dependency or import changes.

Until that plan is implemented, validate repository-only changes with:

- `rg --files --hidden -g '!.git/**'`
- `git diff --check`
- `git status --short`

## Repository Layout

- `AGENTS.md` defines repository-specific instructions for coding agents.
- `ARCHITECTURE.md` describes intended code ownership, boundaries, and extension points.
- `CODESTYLE.md` defines Go-oriented source conventions, commenting requirements, and file-size expectations.
- `DESIGN.md` records the medieval wizardry design direction and UI review expectations.
- `PLANS.md` defines how ExecPlans are written and maintained.
- `PRODUCT.md` records current user-visible product truth.
- `ROADMAP.md` records intended product direction and explicit non-priorities.
- `plans/` stores ordered ExecPlans.
- `.agents/skills/` stores repo-local agent workflows.
- `.codex/config.toml` stores project-scoped Codex defaults.

The first implementation plan will create `go.mod`, `go.sum`, `cmd/td/`, and package directories under `internal/` as needed. Static game assets should eventually live in `assets/`, and tests should mirror the package layout under `tests/` or live beside Go packages when idiomatic package-level tests are clearer.

## Development Notes

Do not implement product feature code during bootstrap work. Start substantial changes by reading the control documents and then creating or updating an ordered ExecPlan in `plans/`.

When adding Go code, follow `CODESTYLE.md`: keep functions focused, document every function or method with Go doc comments, prefer descriptive names, and check hand-written code file line counts at the end of substantial work.
