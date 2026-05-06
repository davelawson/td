# Code Style

Use this document as the durable source of truth for code conventions in this template and in future projects derived from it. It covers source formatting, naming, documentation style, strict commenting standards, and code-file size expectations.

Update this file when the project adopts a formatter, linter, language-specific style rule, naming convention, documentation convention, or commenting standard that future contributors must follow.

## Scope

These rules apply to runtime code, tests, scripts, and code-bearing configuration files that support comments. Markdown documents use their heading and opening prose as their file-level explanation instead of comment syntax. For formats that do not support comments, such as strict JSON, do not add invalid comments; document the file's purpose in a nearby README, schema, or root control document.

When an ecosystem has a strong convention that conflicts with this file, prefer the ecosystem convention only when it improves readability or tooling compatibility. Record durable exceptions here so later contributors do not have to rediscover them.

## General Source Style

- Keep functions small enough that their purpose and control flow are easy to review.
- Keep hand-written code files below 600 lines whenever practical. Treat this as a strong maintainability preference, not a hard limit for generated files, lockfiles, vendored dependencies, or framework-required manifests.
- Prefer explicit, descriptive names over abbreviations. A new contributor should understand what a variable, function, file, or module owns without reading every caller first.
- Use 4 spaces for code indentation unless the language ecosystem strongly prefers another width. Use 2 spaces for Markdown, YAML, JSON, TOML, and similar structured files.
- Put runtime code in `src/`, mirror tests in `tests/`, and keep static assets in `assets/` when they are needed.
- Use `kebab-case` for ordinary Markdown filenames outside conventional root control documents.
- Keep root control documents in uppercase names such as `README.md`, `PRODUCT.md`, `ROADMAP.md`, `PLANS.md`, `CODESTYLE.md`, `ARCHITECTURE.md`, and `AGENTS.md`.
- Name ExecPlans under `plans/` with the ordered `NN-kebab-case-name.md` pattern described in `PLANS.md`.
- If you add a formatter or linter, run it before review and document the canonical command in `README.md` and `AGENTS.md`.

## File Size and Modularity

Keep each hand-written code file below 600 lines whenever practical. Long files are harder for humans and agents to review, summarize, and safely change. The right response is not automatic slicing; split only when a real responsibility boundary, test boundary, or reusable abstraction makes the code easier to understand.

Generated files, lockfiles, vendored dependencies, build output, and framework-owned manifests may exceed 600 lines. For hand-written source, tests, scripts, and code-bearing configuration, treat any file above 600 lines as a refactor signal that needs a conscious decision.

At the end of substantial code changes, check line counts for hand-written code files. If a file exceeds 600 lines, or is close enough that the next likely change will push it over the preference, summarize the file path, line count, and likely cause. Suggest one or more concrete options such as splitting by responsibility, moving fixtures or tests, extracting a shared helper, introducing a well-supported library, or accepting a documented exception. Do not perform that extra split, refactor, or library addition unless the user approves it or the accepted ExecPlan already included that work.

## Commenting Standard

Comments within a block of code are required only in situations where a particular block of code is complex enough to not be easily and quickly understood.  Comments should explain purpose and how that purpose is achieved by the code.

### File Header Comments

Source files should not include header comments.

### Function Comments

Every function must have a comment at the top, using the language's normal documentation-comment or docstring style when one exists. Simple functions may use one concise sentence. Complex functions need a short step-by-step outline.

The function comment must explain:

- the function's purpose
- the meaning of important parameters or inputs
- what the function returns or mutates
- side effects such as network calls, filesystem writes, database updates, rendering, logging, or event dispatch
- errors, exceptions, rejected promises, or failure results callers should expect
- the high-level control flow when the function has multiple branches or phases

Prefer comments that help someone understand why the function exists and how to call it safely.

## Documentation Style

Keep documentation in the file that owns the relevant truth:

- `PRODUCT.md` describes current product capabilities, workflows, constraints, and user-visible limits.
- `ROADMAP.md` describes durable future direction, intended audiences, planned capabilities, priorities, and non-priorities.
- `README.md` describes onboarding, commands, repository layout, and user-facing project context.
- `PLANS.md` describes how ExecPlans must be written and maintained.
- `CODESTYLE.md` describes code conventions, commenting standards, and code-file size expectations.
- `AGENTS.md` describes repository-specific instructions for coding agents.

Reserve new root-level `ALLCAPS.md` files for durable guidance or control documents, not scratch notes or one-off design dumps. When a new root control document is introduced, update `AGENTS.md`, `PLANS.md`, and any root document inventory that future contributors rely on.

