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

The current phase is still early prototype foundation. Peaceful play is split into immediate post-Raid Labour and player-ended Management. A new game begins in Management on sparse generated grasslands, where exploration presents the preassigned grasslands or hills biome beside each free adjacent-Plot reveal button. Home and explored grasslands are biased toward sparse Trees, hills is biased toward sparse Boulders, and northward exploration extends the current Raid path. Tower and economic building construction requires and reserves available inhabitants. House construction provides the first population-seeding path by converting Wood into Peasants, Barracks construction converts Peasants into Soldiers, Dorm construction converts a Peasant into an Apprentice, and Woodcutter, Stone Quarry, and Iron Mine provide the first Labour-phase resource production. Raids now scale deterministically with Raid number, total settlement population, and every explored Plot, tying expansion and growth to defensive pressure. Broader recruitment remains a missing capability for richer population growth, but every current tower staffing role now has a normal-play source.

The immediate product goal is not to make a full tower-defense encounter. The next work should extend the basic Raid baseline with one observable gameplay-facing slice at a time while preserving the small, testable codebase.

## Near-Term Priorities

The completed foundation plans have established the runnable app, expanded menu flow, menu package boundary, New Game configuration screen, current 1920x1080 resize policy, first game-state package, app-mode transition, logical update counter, pause behavior, in-game overlay menu, a prototype top bar for future gameplay status, a static starting Plot scene backed by prototype map state, basic camera zoom plus keyboard and right-drag pan for scene inspection, first selected-object state for structures and raiders, and a first selected-object detail panel. The next priorities should remain separate plans so the codebase grows from verified behavior rather than speculative architecture.

1. Add the next basic scene interaction. Camera movement, selected-object state, selected-object detail panels, and building-bar drag placement are now the first implemented inspection and command-surface models; future interaction slices should stay small and avoid upgrades, selling, range previews, or broader resource rules until they have a visible reason to exist.
2. Extend the early deterministic defense loop. Fixed-path scaling Raids and testable tower combat now exist; the next defense slice should improve path or encounter variety, balance feedback, or defensive choices as one observable addition at a time.
3. Extend the first resource and base-building slice. Costs, Management placement, and Labour-phase economic building income now exist; future work should add gathering, richer build rules, worker reassignment, or placement feedback only as separate observable slices.
4. Extend recruitment or population-seeding beyond the first House, Barracks, and Dorm slices only when richer growth or assignment choices are needed. Reassignment and staff release should remain separate later slices tied to tower removal or management.

## Later Opportunities

Later roadmap opportunities include richer exploration spaces, exploration costs, scouting-versus-claiming choices, resource nodes, inhabitant recruitment and work assignments, base-building placement rules, magical tower archetypes, richer enemy waves, combat resolution, UI panels, art assets, broader audio such as music and volume controls, progression, and eventually platform packaging.

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
