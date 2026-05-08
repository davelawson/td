# Code Style

`CODESTYLE.md` is the durable source of truth for source formatting, naming, documentation style, commenting standards, and file-size expectations in `td`.

The project contains Go runtime code for a local Ebitengine prototype.

## Scope

These rules apply to runtime code, tests, scripts, and code-bearing configuration files that support comments. Markdown documents use their headings and opening prose as their file-level explanation. For formats that do not support comments, do not add invalid comments; document purpose in a nearby Markdown file or control document.

When Go or Ebitengine conventions conflict with a generic preference in this file, prefer the ecosystem convention and record any durable exception here.

## Go Source Style

- Use `gofmt` formatting for every Go file.
- Keep package names short, lower-case, and descriptive.
- Prefer explicit, descriptive identifiers for game concepts such as `menuButton`, `screenWidth`, or `quitRequested`.
- Keep functions and methods small enough that their purpose and control flow are easy to review.
- Prefer ordinary Go errors and `errors.Is` checks where sentinel errors are involved.
- Put the executable entry point under `cmd/td/`.
- Put reusable game packages under `internal/` until there is a clear reason to expose public packages.
- Keep static assets under `assets/` once real assets exist.
- Use `go test ./...` as the canonical full test command once the Go module exists.

## File Size and Modularity

Keep each hand-written code file below 600 lines whenever practical. This is a strong maintainability preference, not a hard limit for generated files, lockfiles, vendored dependencies, or framework-owned manifests.

For hand-written source, tests, scripts, and code-bearing configuration, treat any file above 600 lines as a refactor signal that needs a conscious decision. The right response is not automatic slicing; split only when a real responsibility boundary, test boundary, or reusable abstraction makes the code easier to understand.

At the end of substantial code changes, check line counts for hand-written code files. If a file exceeds 600 lines, or is close enough that the next likely change will push it over the preference, summarize the file path, line count, and likely cause. Suggest a concrete option such as splitting by screen, separating rendering from state updates, moving fixtures, extracting test helpers, or documenting a deliberate exception. Do not perform that extra work unless the user approves it or the accepted ExecPlan already included it.

## Commenting Standard

Do not add source file header comments.

Every Go function and method must have a doc comment immediately above it. Simple functions may use one concise sentence. Complex functions need a short outline of purpose, important inputs, return values, side effects, and failure behavior.

Comments inside a function body are required only when a block is complex enough that a reader would otherwise need to reverse-engineer intent. Use comments to explain purpose and approach, not to narrate obvious statements.

## Testing Style

New behavior should ship with tests whenever it can be exercised without relying on a graphics window. Prefer package-level Go tests for game state, geometry, hit testing, menu selection, and other pure behavior. Use manual or screenshot validation for behavior that genuinely depends on a desktop Ebitengine window.

Tests may live beside the package they test when that is the idiomatic Go choice. A top-level `tests/` directory may still be used later for integration fixtures or black-box tests if the project grows enough to need it.

## Documentation Style

Keep durable truth in the file that owns it:

- `PRODUCT.md` describes current product capabilities, workflows, constraints, and user-visible limits.
- `ROADMAP.md` describes durable future direction, intended audience, planned capabilities, priorities, and non-priorities.
- `DESIGN.md` describes semantic visual direction and UI review expectations.
- `README.md` describes onboarding, commands, repository layout, and user-facing project context.
- `PLANS.md` describes how ExecPlans must be written and maintained.
- `ARCHITECTURE.md` describes code structure, boundaries, and invariants.
- `AGENTS.md` describes repository-specific instructions for coding agents.

Reserve new root-level `ALLCAPS.md` files for durable control documents, not scratch notes. When a new root control document is introduced, update `AGENTS.md`, `PLANS.md` if planning rules change, and any root document inventory that future contributors rely on.
