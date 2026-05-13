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

The current phase is still early prototype foundation, but the repository now has a more complete runnable shell than the first menu slice. It can build, test, and run a 1920x1080 Go/Ebitengine desktop app with a resizable pixel-sized layout, a main menu, a New Game configuration screen for entering a Wizard name, a placeholder Settings screen, disabled Load option, a first static tree-bordered home Plot scene, prototype map state, mouse-wheel map zoom, `WASD` camera panning, a prototype top-bar HUD, a debug logical update counter, SPACE pause behavior, an ESC in-game overlay menu with surrender-to-menu behavior, and a clean quit path.

The immediate product goal is not to make a full tower-defense encounter. The next work should move beyond menus by adding one observable gameplay-facing slice at a time while preserving the small, testable codebase.

## Near-Term Priorities

The completed foundation plans have established the runnable app, expanded menu flow, menu package boundary, New Game configuration screen, current 1920x1080 resize policy, first game-state package, app-mode transition, logical update counter, pause behavior, in-game overlay menu, a prototype top bar for future gameplay status, a static starting Plot scene backed by prototype map state, and basic camera zoom and pan for scene inspection. The next priorities should remain separate plans so the codebase grows from verified behavior rather than speculative architecture.

1. Add the next basic scene interaction. Camera movement is now the first implemented inspection model; a future small slice can add cursor inspection of the home Plot before adding combat or resource rules.
2. Add an early deterministic defense loop. Use fixed paths, placeholder enemies, placeholder towers, simple targeting, and testable combat rules before adding art assets or broader encounter variety.
3. Add the first resource and base-building slice. Introduce gathering, costs, and placement only after the defense loop has a visible baseline that can show why those decisions matter.

## Later Opportunities

Later roadmap opportunities include larger exploration spaces, resource nodes, base-building placement rules, magical tower archetypes, enemy waves, combat resolution, UI panels, art assets, audio, progression, and eventually platform packaging.

Campaign structure and save/load may become important after the prototype has enough systems to preserve meaningful progress. They are intentionally not part of the first phase.

## Explicit Non-Priorities

The current roadmap window does not prioritize saving games, campaign progression, multiplayer, online services, mod support, release packaging, Steam integration, analytics, or production art pipelines.

## Relationship To Other Control Documents

- `PRODUCT.md` describes current product truth.
- `GAME.md` records intended game design decisions and open gameplay questions.
- `README.md` explains onboarding, status, commands, and layout.
- `DESIGN.md` captures the intended fantasy visual language.
- `ART.md` captures guidance for generated art assets, prompt patterns, asset review criteria, and prototype asset constraints.
- `ARCHITECTURE.md` describes intended code structure and boundaries.
- `CODESTYLE.md` describes source conventions.
- `PLANS.md` describes how substantial work is planned.
- `AGENTS.md` describes repository-specific coding-agent behavior.

`ROADMAP.md` may describe capabilities that do not exist yet. `GAME.md` should hold intended gameplay design details and decisions. `PRODUCT.md` should only describe current reality.

## Open Questions

- How should exploration transition into tower-defense encounters?
- What resources should the player gather, and how should they constrain base-building?
- Should the long-term game use tile-based maps, freeform placement, or another spatial model?
- Should gameplay rendering use the menu's pixel-sized resize policy, or should the world view introduce a separate camera and scaling model?
- What is the first scene interaction that best proves the game direction: camera movement, cursor inspection, or direct player movement?
