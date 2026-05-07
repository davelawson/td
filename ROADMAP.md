# Product Roadmap: td

`ROADMAP.md` is the durable source of truth for where `td` is intended to go. Update it when the long-lived product vision, intended audience, planned capabilities, sequencing assumptions, or explicit non-priorities change.

## Product Vision

`td` is intended to become a local PC tower-defense game where a player explores a fantasy environment, gathers resources, builds and improves a base, and then survives tower-defense threats. The game direction is medieval wizardry: spellcraft, enchanted structures, hostile creatures, and readable tactical spaces should carry the identity more than modern military or sci-fi motifs.

The mature game should make preparation and defense feel connected. Exploration and resource gathering should feed base-building decisions, while base layout and magical tower choices should matter during conventional wave-defense encounters.

## Intended Users and End State

The primary intended user is a single local PC player who likes strategy games with planning, spatial pressure, and progression. Success means the player can make understandable choices before an attack, watch those choices affect the defense, and improve their next attempt through better exploration, gathering, and building.

The secondary user is the developer maintaining the prototype. Success for this user means the project remains small, understandable, and easy to test while the core loop is still uncertain.

## Strategic Principles

- Prove one observable workflow at a time before expanding game systems.
- Keep the local desktop prototype simple until the core loop is fun.
- Prefer clear, deterministic gameplay code over premature framework abstractions.
- Let medieval wizardry shape visual language, naming, and later content choices.
- Treat save/load, campaign structure, and distribution as later work.

## Current Phase

The current phase is early prototype foundation. The repository can now build, test, and run a small Ebitengine app with a visible main menu and a clean quit path. The immediate product goal is still not to make a full tower-defense encounter; the next work should add one observable gameplay-facing slice at a time.

## Near-Term Priorities

The first priority, `plans/00-initial-ebitengine-menu.md`, has produced the runnable menu shell. The next likely priorities are a static prototype scene, simple input handling, a placeholder map, and an early tower-defense loop with fixed paths and placeholder shapes. These should remain separate plans so the codebase grows from verified behavior rather than speculative architecture.

## Later Opportunities

Later roadmap opportunities include exploration spaces, resource nodes, base-building placement rules, magical tower archetypes, enemy waves, combat resolution, UI panels, art assets, audio, progression, and eventually platform packaging.

Campaign structure and save/load may become important after the prototype has enough systems to preserve meaningful progress. They are intentionally not part of the first phase.

## Explicit Non-Priorities

The current roadmap window does not prioritize saving games, campaign progression, multiplayer, online services, mod support, release packaging, Steam integration, analytics, or production art pipelines.

## Relationship To Other Control Documents

- `PRODUCT.md` describes current product truth.
- `README.md` explains onboarding, status, commands, and layout.
- `DESIGN.md` captures the intended fantasy visual language.
- `ARCHITECTURE.md` describes intended code structure and boundaries.
- `CODESTYLE.md` describes source conventions.
- `PLANS.md` describes how substantial work is planned.
- `AGENTS.md` describes repository-specific coding-agent behavior.

`ROADMAP.md` may describe capabilities that do not exist yet. `PRODUCT.md` should only describe current reality.

## Open Questions

- How should exploration transition into tower-defense encounters?
- What resources should the player gather, and how should they constrain base-building?
- Should the long-term game use tile-based maps, freeform placement, or another spatial model?
