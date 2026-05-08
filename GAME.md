# Game Design: td

`GAME.md` is the durable source of truth for intended game design in `td`, regardless of how far implementation has come. Update it whenever the game adopts, revises, or rejects a meaningful gameplay decision.

This file describes the game the prototype is trying to become. It may include planned behavior that does not exist yet. For current implemented behavior, use `PRODUCT.md`. For sequencing and priorities, use `ROADMAP.md`. For visual and interaction presentation, use `DESIGN.md`.

## Design Status

The game design is intentionally early. The current implementation is a runnable Go/Ebitengine shell with menus, Wizard name entry, a placeholder game screen, pause behavior, and an in-game overlay menu. The actual exploration, resource, base-building, and tower-defense systems have not been implemented.

Treat sections below as living intent. Decisions marked as open should not be silently assumed by implementation plans; they should be resolved in `GAME.md` when design work makes them concrete.

## Terminology

Use these terms consistently when describing the intended game:

- `Sanctum`: the tower at the center of the wizard's Domain.
- `Domain`: the territory claimed by the wizard.
- `Plot`: an 11x11 group of Tiles used as the main unit of exploration, Domain expansion, and Chapter progress.
- `Fable`: a playthrough of the game as a wizard wherein the wizard attempts to overcome their Nemesis. A Fable is composed of a series of Chapters.
- `Chapter`: a section of a Fable wherein the wizard attempts to overcome a Rival.
- `Nemesis`: a powerful antagonist to the wizard who influences an entire Fable and commands several of the wizard's Rivals.
- `Rival`: a significant antagonist to the wizard who operates as an agent of the Nemesis.
- `Raid`: an assault on the wizard's Domain by the Rival of the current Chapter.

## Player Fantasy

The player is a wizard establishing and defending a vulnerable Sanctum within a hostile Domain. The intended experience is not only placing towers during an attack, but also preparing the conditions that make a defense possible: scouting the surroundings, gathering useful resources, shaping the Domain, and then watching those preparation choices matter when enemies arrive.

The game should feel strategic, readable, and slightly arcane. The player should understand why they survived or failed, and should see a practical path to improve through better exploration, resource use, construction, and defense layout.

## Design Pillars

- Preparation and defense are connected. Exploration, gathering, and building should directly influence tower-defense outcomes.
- Readability comes before complexity. The player should be able to inspect threats, defenses, paths, and resource constraints without decoding hidden rules.
- The fantasy is medieval wizardry. Spellcraft, enchanted structures, hostile wilds, old stone, runes, and magical resources should shape gameplay names and decisions.
- Prototype systems should stay deterministic where practical. Predictable behavior makes early balance and tests easier while the core loop is still uncertain.
- One useful workflow should be proven at a time. Avoid adding campaign, save, production art, or large framework assumptions before the core loop is playable.

## Intended Core Loop

The mature game should revolve around a repeatable loop:

1. Explore a local area to reveal space, threats, routes, and resources.
2. Gather or claim resources that constrain what can be built or improved.
3. Build and adjust the Domain with magical defenses around the Sanctum.
4. Defend against enemies in a tower-defense encounter.
5. Recover from the result, learn what worked, and improve the next preparation phase.

The exact boundaries between exploration, building, and defense are not decided yet. They may become separate phases, a continuous loop, or a hybrid model if prototyping shows that works better.

## Intended Systems

### Exploration

Exploration should give the player information and access. It may reveal terrain, enemy routes, resource nodes, buildable areas, hazards, or future attack pressure.

During the calm phase, the wizard can spend resources to explore another Plot adjacent to the current Domain. Once a Plot has been explored, the wizard can begin building structures there. Exploration does not transition directly into tower-defense encounters. Instead, exploration expands the wizard's Domain, allowing the wizard to build structures across a greater area and defend along a longer path.

Open decisions include whether exploration uses direct player movement, camera-based inspection, tile reveal, or another interaction model, which resources are spent to explore Plots, and how newly explored Plots are connected to enemy paths and roads.

### Map

The map is built on a grid. Each grid square is called a Tile. Tiles are grouped into 11x11 Plots. A Plot is the main strategic map unit: the player explores Plots, expands the Domain by claiming Plots, discovers resources within Plots, and eventually reaches the Rival's Lair by pushing the Domain outward through adjacent Plots.

Each Tile has a terrain type, a height, and sometimes a feature. A feature can be a structure, such as a tower, a resource node, or a road.

At the beginning of a new Fable, the wizard's Domain is a single Plot. The Sanctum is at the center of that starting Plot, and a road leaves the Sanctum.

Plots exist in one of these high-level states:

- `Unknown`: the Plot has not been explored. Its detailed Tiles, resources, roads, hazards, and structures are not visible to the player.
- `Scouted`: the Plot has been explored enough to reveal its terrain, resources, and road exits, but it is not yet part of the Domain.
- `Claimed`: the Plot belongs to the Domain. The wizard can build structures there during calm phases, and roads within the Plot can become part of Raid paths.
- `Lair`: the Plot contains the current Chapter Rival's Lair. Claiming this Plot triggers the final Raid of the Chapter.

Exploring a Plot changes it from `Unknown` to `Scouted`. Claiming a scouted Plot changes it to `Claimed` and expands the Domain. The first prototype may collapse scouting and claiming into one action if that keeps the initial loop simpler, but the design should preserve the distinction because it creates room for future choices: reveal a risky area now, or spend more to actually incorporate it into the defended Domain.

Plots are adjacent orthogonally, not diagonally, for exploration and Domain expansion. A new Plot can only be explored if it touches at least one Claimed Plot on its north, south, east, or west edge. This keeps expansion readable and prevents disconnected pockets of Domain from appearing before there is a deliberate system for them.

Each Plot edge has one possible road connection point: the middle Tile of that edge. Roads do not need to exist on every edge, but when a road connects one Plot to an adjacent Plot, it must leave one Plot through the middle Tile of the shared edge and enter the adjacent Plot through the matching middle Tile. For an 11x11 Plot, these edge connector Tiles are the center Tile of the north, south, east, or west edge.

Roads are the primary way enemies move during Raids. When a Plot is claimed, any road segment that connects back to the Sanctum can become part of the defended route. Expanding toward richer resources or the Rival's Lair may lengthen that route, creating the intended tradeoff between growth and exposure.

Plot contents should be readable at a glance. A Plot can mix terrain types, but each Plot should have a dominant character that helps the player reason about it before inspecting every Tile, such as wooded, rocky, wetland, open meadow, hill, ruin, or lair. This dominant character is a design label, not necessarily a separate data type.

The starting Plot should be mostly buildable and forgiving. It should contain the Sanctum, at least one outgoing road, some nearby buildable Grass, and enough visible resource access to support the first build decisions. It should not start with Water or Mountain terrain blocking all useful placement around the Sanctum.

The terrain types are:

- `Grass`
- `Forest`
- `Mountain`
- `Water`

Open decisions include Tile dimensions in screen or world space, height scale, movement and build rules for each terrain type, feature exclusivity rules, exact internal road-shape behavior, whether Scouted and Claimed remain separate in the first implementation, how Plot contents are generated or authored, how the Rival's Lair is placed, and where the starting road leads.

### Resources

Resources should make preparation choices meaningful. They should constrain base-building and upgrades without becoming bookkeeping for its own sake.

The initial resources are:

- `Wood`: gathered by cutting down trees.
- `Stone`: gathered by quarrying rocks.
- `Metal`: gathered by exploiting mines.

Open decisions include whether resources are gathered manually or passively, whether resources persist across encounters, whether additional resources are needed, and how resource nodes are discovered or depleted.

### Base-Building

Base-building should let the player shape a defensible Domain around the Sanctum. Placement should matter spatially, and construction should feel connected to the wizardry theme.

Structures can only be built in explored Plots. Expanding the Domain gives the wizard more buildable area, but it can also lengthen the path the wizard must defend during Raids.

Open decisions include what counts as buildable terrain, whether structures block paths, how much rebuilding is allowed between attacks, and how explored Plots become claimed, defended, or otherwise incorporated into the Domain.

### Sanctum

The Sanctum is the wizard's home and private laboratory. It is the tower at the center of the Domain and the structure the wizard must protect during Raids.

The Sanctum is protected by an arcane barrier. This barrier has a number of charges. Every time an enemy attempts to breach the Sanctum, the barrier expends one charge, if any charges remain, and atomizes that enemy. If an enemy attempts to breach the Sanctum when the barrier has no charges remaining, the Sanctum is breached.

Open decisions include how many charges the barrier has, whether charges can be restored or increased, whether different enemies interact with the barrier differently, and what happens after the Sanctum is breached.

### Tower Defense

Tower-defense encounters should use clear enemy movement, clear defensive coverage, and visible combat results. The first combat slice should favor simple fixed paths, placeholder enemies, placeholder towers, and testable targeting rules.

Open decisions include enemy families, wave structure, win and loss conditions, damage types, and upgrade rules.

### Raids

A Raid is the assault at the end of a Day. It is launched by the Rival of the current Chapter against the wizard's Domain.

During a Raid, enemies attempt to reach and breach the Sanctum. The wizard's preparation during the preceding calm phase should matter: explored Plots, built structures, roads, terrain, and the length of the defended path all shape how the Raid plays out.

A Raid ends when all enemies have been defeated or when the Sanctum has been breached. If an enemy attempts to breach the Sanctum while the arcane barrier has charges remaining, the barrier spends a charge and atomizes that enemy instead. If no barrier charges remain, the Sanctum is breached.

The final Raid of a Chapter is triggered by Domain expansion. Once the wizard's Domain has expanded to include the Plot containing the Rival's Lair, the next Raid is the final Raid. If the wizard overcomes that final Raid, the Rival is defeated and the Chapter is completed successfully.

Open decisions include enemy spawn rules, enemy path selection, whether a Raid can include multiple waves, how Raid difficulty scales with Domain expansion, and what happens when the Sanctum is breached.

### Tower Types

Tower types define the defensive structures the wizard can build in the Domain.

- `Bow Tower`: a tower replete with magically automated bows that fire arrows at enemies within range. It costs only Wood to build. The Bow Tower is a general-purpose tower that deals moderate damage at moderate range.

Open decisions include exact tower costs, ranges, damage values, firing rates, targeting behavior, upgrade paths, and what other tower types exist.

### Chapters

A Chapter is the wizard's attempt to defeat a Rival. It lasts as long as it takes for the wizard to expand the Domain to the Rival's Lair and overcome the final Raid that follows.

A Chapter is composed of a series of Days. Every Day begins with a calm phase and ends with a Raid. The calm phase lasts for a set amount of time, giving the wizard a bounded window to explore, gather, build, repair, or otherwise prepare the Domain before the attack begins.

Once the wizard's Domain has expanded to include the Plot containing the Rival's Lair, the next Raid is the final Raid. If that Raid is overcome, the Chapter has been completed successfully.

Open decisions include how the Rival's Lair is revealed, what the wizard can do during the calm phase, whether time can be paused or accelerated, and what happens when the Sanctum is breached.

### Progression

Progression is intended later, after the base loop proves useful. It may include improved spells, new structures, persistent settlement growth, campaign maps, or encounter-to-encounter unlocks.

Open decisions include whether progression is run-based, campaign-based, scenario-based, or deferred entirely until the prototype has a satisfying local loop.

## Current Game Design Decisions

- The game targets a single local PC player.
- The genre blend is exploration, resource gathering, base-building, and tower-defense combat.
- The setting is medieval wizardry fantasy, not modern military or science fiction.
- The player identity is a wizard, currently represented by Wizard name entry in the New Game screen.
- Save/load, campaign structure, multiplayer, online services, production art pipelines, and release packaging are not part of the current prototype phase.
- The next gameplay-facing work should move beyond menus with one observable slice at a time, starting with a static prototype scene before combat or resource rules.

## Open Game Design Questions

- Which resources are spent to explore an adjacent Plot?
- How do newly explored Plots connect to roads and enemy paths?
- How is the Rival's Lair revealed or signaled to the player?
- What movement, build, resource, and road rules apply to each terrain type?
- Should scouting and claiming be separate actions in the first playable implementation, or should early exploration immediately claim a Plot?
- How are Plot dominant characters generated, selected, or presented to the player?
- How many road exits can a Plot have, and can roads branch inside a Plot after entering through edge-center connector Tiles?
- Are Wood, Stone, and Metal enough for interesting choices, or are additional resources needed?
- Does the wizard move as an on-map character, act through camera inspection, or both?
- How many arcane barrier charges does the Sanctum have, and can those charges be restored or increased?
- How does the Domain expand, contract, or change over time?
- What enemy archetypes, spawn rules, and pathing rules should Raids use first?
- What are the first win and loss conditions?
- Should the early prototype use separate phases or continuous real-time play?
- How long is a calm phase, and can the player pause or accelerate it?
- What happens when the Sanctum is breached during a Raid?
- When should save/load or campaign structure become meaningful enough to design?

## Decision Log

Record game design decisions here when they become durable enough to guide implementation.

- Decision: `GAME.md` is the authoritative home for intended game design, separate from current product behavior and implementation sequencing.
  Rationale: The project needs a place to accumulate gameplay decisions even before those systems exist in code.
  Date/Author: 2026-05-08 / Codex

- Decision: Use `Sanctum` for the tower at the center of the wizard's Domain, and `Domain` for the territory claimed by the wizard.
  Rationale: These terms give future design work stable names for the central defended structure and the player-claimed territory around it.
  Date/Author: 2026-05-08 / Codex

- Decision: Use `Fable`, `Chapter`, `Nemesis`, `Rival`, and `Raid` for the game's playthrough and antagonist structure.
  Rationale: These terms establish a hierarchy where a wizard's Fable is shaped by a Nemesis, each Chapter focuses on a Rival, and Raids are the Rival's assaults on the Domain.
  Date/Author: 2026-05-08 / Codex

- Decision: Structure each Chapter as a series of Days, where each Day has a timed calm phase followed by a Raid.
  Rationale: This creates a repeated preparation-and-assault cadence inside each Rival-focused Chapter while keeping exact timing, Raid counts, and breach consequences open for later design.
  Date/Author: 2026-05-08 / Codex

- Decision: Protect the Sanctum with an arcane barrier that atomizes breaching enemies by spending charges.
  Rationale: This gives the central tower a clear defensive rule and creates a concrete failure threshold when enemies reach it after the barrier is exhausted.
  Date/Author: 2026-05-08 / Codex

- Decision: Use Wood, Stone, and Metal as the initial resources.
  Rationale: These resources give base-building a simple physical economy tied to visible world objects: trees, rocks, and mines.
  Date/Author: 2026-05-08 / Codex

- Decision: Use the Bow Tower as the first defined tower type.
  Rationale: A wood-only, moderate-damage, moderate-range tower gives the design a simple general-purpose baseline before specialized magical towers are defined.
  Date/Author: 2026-05-08 / Codex

- Decision: Build maps from grid Tiles grouped into 11x11 Plots.
  Rationale: A tile grid gives exploration, building, resources, roads, terrain, and height a shared spatial model while Plots provide a larger grouping for generation, reveal, and Domain-scale reasoning. The odd Plot dimension gives each Plot a natural center Tile for the Sanctum, lairs, landmarks, and authored road composition.
  Date/Author: 2026-05-08 / Codex

- Decision: Treat Plots as the main strategic map unit for exploration, Domain expansion, resource discovery, and Rival Lair progress.
  Rationale: This keeps player-facing map decisions at a readable scale while still allowing Tile-level terrain, road, resource, and structure rules inside each Plot.
  Date/Author: 2026-05-08 / Codex

- Decision: Define Plot states as `Unknown`, `Scouted`, `Claimed`, and `Lair`.
  Rationale: These states separate hidden information, revealed information, buildable Domain territory, and Chapter-ending territory without requiring the first implementation to expose every distinction immediately.
  Date/Author: 2026-05-08 / Codex

- Decision: Use orthogonal Plot adjacency for exploration and Domain expansion.
  Rationale: North, south, east, and west adjacency keeps expansion readable and avoids disconnected Domain pockets before the game deliberately supports them.
  Date/Author: 2026-05-08 / Codex

- Decision: Roads connect adjacent Plots only through the middle Tile of the shared Plot edge.
  Rationale: A fixed edge-center connector rule keeps inter-Plot paths readable, takes advantage of the 11x11 Plot centerline, and lets internal road shapes vary without making Plot-to-Plot connectivity ambiguous.
  Date/Author: 2026-05-08 / Codex

- Decision: Start each Fable with a one-Plot Domain containing the centered Sanctum and an outgoing road.
  Rationale: This gives every playthrough a clear initial defended space and a road hook for future exploration, enemy routing, and Domain expansion.
  Date/Author: 2026-05-08 / Codex

- Decision: Let the wizard spend resources during calm phases to explore new Plots and unlock building there.
  Rationale: This connects resources, calm-phase planning, Domain expansion, and base-building into one concrete preparation loop.
  Date/Author: 2026-05-08 / Codex

- Decision: Require Plot exploration to be adjacent to the current Domain and make exploration expand the defended path rather than start encounters directly.
  Rationale: This makes exploration a spatial commitment: the wizard gains buildable area, but the Domain can expose a longer route that must be defended during Raids.
  Date/Author: 2026-05-08 / Codex

- Decision: Define Raids as the Rival's end-of-Day assaults on the wizard's Domain.
  Rationale: This gives the defense phase a dedicated design home and ties Raids to Chapters, Rivals, enemy pressure, the defended path, and Sanctum breach conditions.
  Date/Author: 2026-05-08 / Codex

- Decision: Let Chapters last until Domain expansion reaches the Rival's Lair, making the next Raid the final Raid.
  Rationale: This ties Chapter completion to exploration and Domain expansion instead of a fixed day count, so the player controls when to provoke the decisive assault by reaching the Rival's territory.
  Date/Author: 2026-05-08 / Codex
