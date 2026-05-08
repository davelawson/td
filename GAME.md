# Game Design: td

`GAME.md` is the durable source of truth for intended game design in `td`, regardless of how far implementation has come. Update it whenever the game adopts, revises, or rejects a meaningful gameplay decision.

This file describes the game the prototype is trying to become. It may include planned behavior that does not exist yet. For current implemented behavior, use `PRODUCT.md`. For sequencing and priorities, use `ROADMAP.md`. For visual and interaction presentation, use `DESIGN.md`.

## Design Status

The game design is intentionally early. The current implementation is a runnable Go/Ebitengine shell with menus, Wizard name entry, a placeholder game screen, pause behavior, and an in-game overlay menu. The actual exploration, resource, base-building, and tower-defense systems have not been implemented.

Treat sections below as living intent. Decisions marked as open should not be silently assumed by implementation plans; they should be resolved in `GAME.md` when design work makes them concrete.

## Player Fantasy

The player is a wizard establishing and defending a vulnerable magical foothold in hostile fantasy territory. The intended experience is not only placing towers during an attack, but also preparing the conditions that make a defense possible: scouting the surroundings, gathering useful resources, shaping a base, and then watching those preparation choices matter when enemies arrive.

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
3. Build and adjust a base with magical defenses.
4. Defend against enemies in a tower-defense encounter.
5. Recover from the result, learn what worked, and improve the next preparation phase.

The exact boundaries between exploration, building, and defense are not decided yet. They may become separate phases, a continuous loop, or a hybrid model if prototyping shows that works better.

## Intended Systems

### Exploration

Exploration should give the player information and access. It may reveal terrain, enemy routes, resource nodes, buildable areas, hazards, or future attack pressure.

Open decisions include whether exploration uses direct player movement, camera-based inspection, tile reveal, or another interaction model.

### Resources

Resources should make preparation choices meaningful. They should constrain base-building and upgrades without becoming bookkeeping for its own sake.

Open decisions include resource names, resource count, gathering method, whether resources are gathered manually or passively, and whether resources persist across encounters.

### Base-Building

Base-building should let the player shape a defensible magical settlement or stronghold. Placement should matter spatially, and construction should feel connected to the wizardry theme.

Open decisions include whether maps are tile-based or freeform, what counts as buildable terrain, whether structures block paths, and how much rebuilding is allowed between attacks.

### Tower Defense

Tower-defense encounters should use clear enemy movement, clear defensive coverage, and visible combat results. The first combat slice should favor simple fixed paths, placeholder enemies, placeholder towers, and testable targeting rules.

Open decisions include tower archetypes, enemy families, wave structure, win and loss conditions, damage types, and upgrade rules.

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

- How should exploration transition into tower-defense encounters?
- Should maps be tile-based, freeform, node-based, or something else?
- What resources should exist, and how many are needed for interesting choices?
- Does the wizard move as an on-map character, act through camera inspection, or both?
- What makes the base vulnerable enough that defense layout matters?
- What are the first tower and enemy archetypes?
- What are the first win and loss conditions?
- Should the early prototype use separate phases or continuous real-time play?
- When should save/load or campaign structure become meaningful enough to design?

## Decision Log

Record game design decisions here when they become durable enough to guide implementation.

- Decision: `GAME.md` is the authoritative home for intended game design, separate from current product behavior and implementation sequencing.
  Rationale: The project needs a place to accumulate gameplay decisions even before those systems exist in code.
  Date/Author: 2026-05-08 / Codex

