# Game Design: td

`GAME.md` is the durable source of truth for intended game design in `td`, regardless of how far implementation has come. Update it whenever the game adopts, revises, or rejects a meaningful gameplay decision.

This file describes the game the prototype is trying to become. It may include planned behavior that does not exist yet. For current implemented behavior, use `PRODUCT.md`. For sequencing and priorities, use `ROADMAP.md`. For visual and interaction presentation, use `DESIGN.md`.

## Design Status

The game design is intentionally early. The current implementation includes a Sanctum-only starting Plot, free calm-phase exploration of adjacent Plots, a House that adds Peasant population, a Barracks that converts Peasants into Soldiers, a Dorm that converts a Peasant into an Apprentice, economic buildings that produce resources after defeated Raids, and tower templates with resource and staffing requirements. Newly explored Plots are immediately usable grassland, and exploring the central north chain extends the current visible road and Raid path. Economic building and tower construction requires sufficient available staff and reserves those inhabitants. Timed recruitment, reassignment, and staff release are not implemented.

Treat sections below as living intent. Decisions marked as open should not be silently assumed by implementation plans; they should be resolved in `GAME.md` when design work makes them concrete.

## Terminology

Use these terms consistently when describing the intended game:

- `Sanctum`: the tower at the center of the wizard's Domain.
- `Domain`: the territory claimed by the wizard.
- `Plot`: a 15x15 group of Tiles used as the main unit of exploration, Domain expansion, and Chapter progress.
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

During the calm phase, the wizard can explore another Plot adjacent to the current Domain. The first implemented slice makes exploration free and immediate: clicking a magnifying-glass button on a border between explored and unexplored Plots reveals the adjacent Plot, and the revealed Plot is treated as part of the usable Domain. Once a Plot has been explored, the wizard can begin building structures there. Exploration does not transition directly into tower-defense encounters. Instead, exploration expands the wizard's Domain, allowing the wizard to build structures across a greater area and, when expanding north, defend along a longer path.

Early map inspection uses camera-based movement rather than direct wizard movement. The player can zoom the scene camera with the mouse wheel, pan with `WASD`, or press and hold the right mouse button over the game view and drag so the visible world follows the cursor. Camera inspection works while paused, but it is only inspection. It does not reveal new Plots, gather resources, or move a wizard character.

The first object-inspection interaction is left-click selection. Structure tiles, including the Sanctum, House, Barracks, Dorm, economic buildings, and towers, can be selected by clicking their Tile, and active raiders can be selected by clicking their visible sprite. A selected object is drawn brighter. The current inspection panel is informational only: raiders show their prototype combat stats, combat towers show their prototype attack stats, population buildings show their cost and population effects, economic buildings show their cost, Peasant requirement, and post-Raid production, and the Sanctum shows only its name. Selection currently has no command buttons or upgrades; it exists to establish readable object targeting for later inspection and command workflows.

Open decisions include whether later exploration uses direct player movement, camera-based inspection plus tile reveal, another interaction model, or a hybrid of these, which resources are spent to explore Plots after the free prototype slice, and how non-north explored Plots are connected to enemy paths and roads.

### Map

The map is built on a grid. Each grid square is called a Tile. Tiles are grouped into 15x15 Plots. A Plot is the main strategic map unit: the player explores Plots, expands the Domain by claiming Plots, discovers resources within Plots, and eventually reaches the Rival's Lair by pushing the Domain outward through adjacent Plots.

Each Tile has a terrain type, a height, and sometimes a feature. A feature can be a structure, such as a tower, a resource node, or a road.

At the beginning of a new Fable, the wizard's Domain is a single Plot. The Sanctum is at the center of that starting Plot, and a road leaves the Sanctum. The first rendered prototype home Plot is open grassland containing the centered Sanctum as its only initial structure and a straight road north to the Plot edge. The player must build every combat tower.

In the intended full design, Plots exist in one of these high-level states:

- `Unknown`: the Plot has not been explored. Its detailed Tiles, resources, roads, hazards, and structures are not visible to the player.
- `Scouted`: the Plot has been explored enough to reveal its terrain, resources, and road exits, but it is not yet part of the Domain.
- `Claimed`: the Plot belongs to the Domain. The wizard can build structures there during calm phases, and roads within the Plot can become part of Raid paths.
- `Lair`: the Plot contains the current Chapter Rival's Lair. Claiming this Plot triggers the final Raid of the Chapter.

Exploring a Plot changes it from `Unknown` to `Scouted`. Claiming a scouted Plot changes it to `Claimed` and expands the Domain. The first implemented exploration prototype collapses scouting and claiming into one action: clicking an explore button immediately reveals the adjacent Plot and makes it buildable. The design should preserve the distinction for later versions because it creates room for future choices: reveal a risky area now, or spend more to actually incorporate it into the defended Domain.

Plots are adjacent orthogonally, not diagonally, for exploration and Domain expansion. A new Plot can only be explored if it touches at least one Claimed Plot on its north, south, east, or west edge. This keeps expansion readable and prevents disconnected pockets of Domain from appearing before there is a deliberate system for them.

Each Plot edge has one possible road connection point: the middle Tile of that edge. Roads do not need to exist on every edge, but when a road connects one Plot to an adjacent Plot, it must leave one Plot through the middle Tile of the shared edge and enter the adjacent Plot through the matching middle Tile. For a 15x15 Plot, these edge connector Tiles are the center Tile of the north, south, east, or west edge.

Roads are the primary way enemies move during Raids. When a Plot is claimed, any road segment that connects back to the Sanctum can become part of the defended route. Expanding toward richer resources or the Rival's Lair may lengthen that route, creating the intended tradeoff between growth and exposure.

World positions use Tile units. The center of the Sanctum Tile is the origin `(0, 0)`, one Tile is one world unit, floating point positions are allowed, and positive Y points north. Plots use the same world-coordinate model: the home Plot is at Plot coordinate `(0, 0)`, the Plot directly north is `(0, 1)`, the Plot directly east is `(1, 0)`, and so on. For example, `(0, 1.5)` is the common edge between the first and second road Tiles north of the Sanctum. Enemy positions should use this world-coordinate model instead of storing only path progress.

Plot contents should be readable at a glance. A Plot can mix terrain types, but each Plot should have a dominant character that helps the player reason about it before inspecting every Tile, such as wooded, rocky, wetland, open meadow, hill, ruin, or lair. This dominant character is a design label, not necessarily a separate data type.

The starting Plot should be mostly buildable and forgiving. It should contain the Sanctum, at least one outgoing road, some nearby buildable Grass, and enough visible resource access to support the first build decisions. It should not start with Water or Mountain terrain blocking all useful placement around the Sanctum.

The terrain types are:

- `Grass`
- `Forest`
- `Mountain`
- `Water`

Open decisions include height scale, movement and build rules for each terrain type, feature exclusivity rules, exact internal road-shape behavior beyond the first north road, whether Scouted and Claimed remain separate in the first implementation, how richer Plot contents are generated or authored, and how the Rival's Lair is placed.

### Resources

Resources should make preparation choices meaningful. They should constrain base-building and upgrades without becoming bookkeeping for its own sake.

The initial resources are:

- `Wood`: gathered by cutting down trees.
- `Stone`: gathered by quarrying rocks.
- `Metal`: gathered by exploiting mines.

The first implemented resource-production slice uses economic buildings rather than direct gathering. A Woodcutter produces 10 Wood, a Stone Quarry produces 10 Stone, and an Iron Mine produces 10 Metal after each successfully defeated Raid. Failed Raids that breach the Sanctum do not pay economic building resources. The Iron Mine uses the existing Metal resource; there is no separate Iron resource.

Open decisions include whether later resources are gathered manually or passively during calm play, whether resources persist across longer campaign structures, whether additional resources are needed, and how resource nodes are discovered or depleted.

The in-game top bar should show the current Wood, Stone, and Metal counts as resource icons plus numbers once the economy exists. Prototype HUD values may be fixed until resource gathering and spending systems are implemented.

### Inhabitants

The wizard's Domain has three inhabitant groups:

- `Apprentices`: magical students available for future spellcraft, research, or arcane assignments.
- `Soldiers`: trained defenders available for future military and guard assignments.
- `Peasants`: general workers available for future gathering, construction, farming, and settlement assignments.

Each group has an available count and a total count. Available means inhabitants not currently committed to an assignment; total means every inhabitant of that type in the Domain. Available cannot be negative or exceed total.

The top bar shows the three groups in a separate population grouping after physical resources. Each group uses an icon followed by `available/total`. The first prototype initializes every group to `0/0`. House construction immediately increases Peasants by `2/2`. Barracks construction consumes 2 available and total Peasants and grants 2 available and total Soldiers. Dorm construction consumes 1 available and total Peasant and grants 1 available and total Apprentice. Economic buildings reserve one available Peasant while leaving total Peasants unchanged. Timed recruitment, reassignment, broader population growth, losses, and staff release are not implemented.

### Base-Building

Base-building should let the player shape a defensible Domain around the Sanctum. Placement should matter spatially, and construction should feel connected to the wizardry theme.

Structures can only be built in explored Plots. Expanding the Domain gives the wizard more buildable area, but it can also lengthen the path the wizard must defend during Raids.

The first build-facing UI is a 260-pixel building bar on the left side of the playable scene. It partitions structures into `Housing`, `Economic`, and `Defenses` tabs, with `Housing` selected by default. Housing shows House, Barracks, and Dorm. Economic shows Woodcutter, Stone Quarry, and Iron Mine. Defenses shows Bow Tower, Flame Bolt Tower, and Catapult Tower. Visible entries use structure sprites with colour-coded prototype construction costs and a compact population row for requirements, costs, or grants displayed to the right of each icon. Hovering a building icon shows an informational tooltip with that structure's description, cost, staffing requirement, and implemented production, population effect, or combat stats. Hovering an eligible icon also brightens that icon and emphasizes its cost row. Buildable icon squares use green outlines. Icons for buildings without sufficient resources, population, or staff use red square outlines and are drawn at 70% opacity. During calm play, left-dragging an eligible building icon attaches a half-sized copy to the cursor; releasing over an empty grass-like Tile places that structure and deducts its cost. In the current prototype, "grass-like" means the existing empty terrain Tile type, not roads or forest. Occupied Tiles, road Tiles, forest Tiles, active Raids, and breached games reject placement without spending resources or population changes. SPACE-paused calm play still allows building, but the in-game overlay blocks it.

Tower templates define staffing requirements. The Bow Tower requires one Soldier, the Flame Bolt Tower requires one Apprentice, and the Catapult Tower requires one Soldier plus two Peasants. Construction is allowed only when every required role is available. A successful build reduces each required role's available count but not its total count. Staff remain committed because tower removal and reassignment do not yet exist. Staffing does not separately disable an already-built tower.

The House is the first population-provider building. It costs 20 Wood, requires no staff, has no combat stats, and immediately grants 2 Peasants by increasing both available and total Peasant population. House population is not removed because structure removal does not yet exist.

The Barracks is the first population-conversion building. It costs 10 Wood and 10 Stone, requires no staff, has no combat stats, consumes 2 available and total Peasants, and immediately grants 2 available and total Soldiers. This creates the first normal-play Soldier source without adding timed recruitment or a general assignment system.

The Dorm is the first Apprentice-conversion building. It costs 10 Wood and 10 Stone, requires no staff, has no combat stats, consumes 1 available and total Peasant, and immediately grants 1 available and total Apprentice. This creates the first normal-play Apprentice source without adding timed recruitment or a general assignment system.

Economic buildings are the first resource-production buildings. The Woodcutter costs 10 Wood, requires one available Peasant, has no combat stats, reserves that Peasant on construction, and produces 10 Wood after each defeated Raid. The Stone Quarry costs 10 Wood and 10 Stone, requires one available Peasant, reserves that Peasant on construction, and produces 10 Stone after each defeated Raid. The Iron Mine costs 10 Wood, 10 Stone, and 10 Metal, requires one available Peasant, reserves that Peasant on construction, and produces 10 Metal after each defeated Raid.

Open decisions include richer build rules for future terrain types, whether structures block paths beyond the current fixed road rejection, how much rebuilding is allowed between attacks, whether later placement needs previews or range indicators, and how explored Plots become claimed, defended, or otherwise incorporated into the Domain.

### Sanctum

The Sanctum is the wizard's home and private laboratory. It is the tower at the center of the Domain and the structure the wizard must protect during Raids.

The Sanctum is protected by an arcane barrier. This barrier has a number of charges. Every time an enemy attempts to breach the Sanctum, the barrier expends one charge, if any charges remain, and atomizes that enemy. If an enemy attempts to breach the Sanctum when the barrier has no charges remaining, the Sanctum is breached.

The player-facing HUD label for remaining barrier charges is `Barricade`, representing the Sanctum's protective barricade charges.

Open decisions include how many charges the barrier has, whether charges can be restored or increased, whether different enemies interact with the barrier differently, and what happens after the Sanctum is breached.

### Tower Defense

Tower-defense encounters should use clear enemy movement, clear defensive coverage, and visible and audible combat results. The first Raid slice uses simple fixed paths with skeleton and zombie enemies. The first combat slice adds automated Bow, Flame Bolt, and Catapult Towers with testable targeting, projectile travel, damage rules, Catapult Tile-area impact, and a prototype sound effect for tower-damage defeats without replacing the simple Raid lifecycle.

Open decisions include enemy families, additional tower targeting modes, win and loss conditions beyond the first breach state, damage types, and upgrade rules.

### Raids

A Raid is the assault at the end of a Day. It is launched by the Rival of the current Chapter against the wizard's Domain.

During a Raid, enemies attempt to reach and breach the Sanctum. The wizard's preparation during the preceding calm phase should matter: explored Plots, built structures, roads, terrain, and the length of the defended path all shape how the Raid plays out.

A Raid ends when all enemies have been defeated or when the Sanctum has been breached. If an enemy attempts to breach the Sanctum while the arcane barrier has charges remaining, the barrier spends a charge and atomizes that enemy instead. If no barrier charges remain, the Sanctum is breached.

The final Raid of a Chapter is triggered by Domain expansion. Once the wizard's Domain has expanded to include the Plot containing the Rival's Lair, the next Raid is the final Raid. If the wizard overcomes that final Raid, the Rival is defeated and the Chapter is completed successfully.

During a Raid, the in-game top bar should show how many enemies remain in the current assault. This can be formatted before enemy simulation exists, but real values should come from the Raid system once it is implemented.

The first implemented Raid behavior is deliberately deterministic. A `Next Raid` button starts the next Raid immediately during calm play. Raid 1 has five enemies in this exact spawn order: skeleton, zombie, skeleton, zombie, skeleton. Each later Raid adds two enemies and remains skeleton-only until a fuller wave-composition design exists. One enemy appears immediately, and the rest spawn one at a time on a fixed stagger. Enemies use the current starting Plot's straight north road, entering from the north-center road edge and moving south to the centered Sanctum. There are no alternate paths yet.

The first Raid slice stores each active enemy's current world position directly. On the starting road, enemies spawn at `(0, 7)` and move south by decreasing their Y coordinate until they contact the Sanctum at `Y <= 0`. Movement speed comes from `EnemyTemplate.SpeedTilesPerSecond`, is measured in Tiles per second, and is converted through the fixed logical update duration. Maximum health comes from `EnemyTemplate.MaxHealth`, and combat-defeat rewards come from `EnemyTemplate.Resources`. The current skeleton template has 50 health, moves at 1.0 Tiles per second, and rewards 5 Wood and 2 Stone. The current zombie template has 75 health, moves at 0.7 Tiles per second, and rewards 4 Wood, 3 Stone, and 1 Metal.

The first projectile-tower combat slice targets the in-range enemy closest to the Sanctum. If two enemies are equally close, the tower uses the older spawned enemy as the deterministic tie-breaker. The Bow Tower range is 3.0 Tiles, damage is 10, fire interval is 1.0 second, and projectile speed is 9.0 Tiles per second. The Flame Bolt Tower range is 2.5 Tiles, damage is 20, fire interval is 1.5 seconds, and projectile speed is 7.0 Tiles per second. The Catapult Tower range is 5.0 Tiles, damage is 75 to every active enemy in the target's Tile, fire interval is 3.0 seconds, and projectile speed is 3.0 Tiles per second. These timing and speed stats are expressed in real-time seconds rather than update counts.

If an enemy reaches the Sanctum while Barricade charges remain, the Barricade spends one charge and that enemy is removed. If an enemy reaches the Sanctum when Barricade is zero, the Sanctum is marked breached, the active Raid is cleared, and no further Raids can start until a future recovery or loss-flow design exists. A short embedded prototype sound plays and the defeated enemy's template resources are granted when tower damage defeats a raider; Barricade removal and breach clearing do not count as combat defeats and do not grant those resources. When a Raid ends with all enemies defeated, each placed economic building pays its post-Raid resource yield once. A breached Raid does not pay economic building resources.

Open decisions include enemy archetypes, whether a Raid can include multiple waves or paths, how Raid difficulty scales with Domain expansion, whether towers or resources can remove enemies before they reach the Sanctum, and what longer-term recovery or loss flow follows a breached Sanctum.

### Tower Types

Tower types define the defensive structures the wizard can build in the Domain.

- `Bow Tower`: costs 30 Wood, 10 Stone, and 10 Metal; requires one Soldier; has 3.0-Tile range; deals 10 damage; fires every 1.0 second; and launches projectiles at 9.0 Tiles per second.
- `Flame Bolt Tower`: costs 30 Stone and 20 Metal; requires one Apprentice; has 2.5-Tile range; deals 20 damage; fires every 1.5 seconds; and launches projectiles at 7.0 Tiles per second.
- `Catapult Tower`: costs 40 Wood, 60 Stone, and 25 Metal; requires one Soldier and two Peasants; has 5.0-Tile range; deals 75 damage to every active enemy in the struck Tile; fires every 3.0 seconds; and launches projectiles at 3.0 Tiles per second.

Open decisions include upgrade paths, specialized targeting modes, what other tower types exist, and whether these first prototype costs remain balanced once resource gathering and spending exist.

### Chapters

A Chapter is the wizard's attempt to defeat a Rival. It lasts as long as it takes for the wizard to expand the Domain to the Rival's Lair and overcome the final Raid that follows.

A Chapter is composed of a series of Days. Every Day begins with a calm phase and ends with a Raid. The calm phase lasts for a set amount of time, giving the wizard a bounded window to explore, gather, build, repair, or otherwise prepare the Domain before the attack begins.

Once the wizard's Domain has expanded to include the Plot containing the Rival's Lair, the next Raid is the final Raid. If that Raid is overcome, the Chapter has been completed successfully.

Open decisions include how the Rival's Lair is revealed, what the wizard can do during the calm phase, whether time can be paused or accelerated, and what happens when the Sanctum is breached.

The in-game top bar should always show the current Chapter name and Day number. During calm phases it should also show the time remaining before the next Raid.

### Progression

Progression is intended later, after the base loop proves useful. It may include improved spells, new structures, persistent settlement growth, campaign maps, or encounter-to-encounter unlocks.

Open decisions include whether progression is run-based, campaign-based, scenario-based, or deferred entirely until the prototype has a satisfying local loop.

## Current Game Design Decisions

- The game targets a single local PC player.
- The genre blend is exploration, resource gathering, base-building, and tower-defense combat.
- The setting is medieval wizardry fantasy, not modern military or science fiction.
- The player identity is a wizard, currently represented by Wizard name entry in the New Game screen.
- Save/load, campaign structure, multiplayer, online services, production art pipelines, and release packaging are not part of the current prototype phase.
- The first gameplay-facing rendered slice starts with an open-grassland home Plot containing only the centered Sanctum as an initial structure, a straight road north, free calm-phase exploration buttons on unexplored adjacent Plot borders, and a widened tabbed building bar listing Housing, Economic, and Defenses groups with right-side costs, population costs, population grants or staffing requirements, informational hover tooltips, green/red capacity outlines, capacity opacity, and calm-phase drag placement.
- Early map inspection uses camera zoom and pan, not wizard-character movement. Mouse-wheel zoom, `WASD` panning, and right-drag panning are inspection controls only and do not change map data.
- The first exploration slice uses free magnifying-glass border buttons during calm play, including paused calm play, to reveal orthogonally adjacent Plots. Revealed Plots are immediately buildable grassland. The first implementation collapses scouting and claiming into one action.
- Northward exploration along Plot `X=0` extends the visible center road and moves deterministic Raid spawning to the farthest explored north road. Non-north explored Plots do not add Raid paths yet.
- The first Raid slice uses deterministic sprite-backed skeleton and zombie enemies on the starting Plot's straight north road. Player-built towers fire at in-range enemies; starting a Raid without first building defenses leaves only the Barricade protecting the Sanctum.
- A new game starts with 100 Wood, 50 Stone, and 20 Metal. Resources cover House, Barracks, Dorm, all three economic buildings, Bow, and Flame Bolt, but zero starting population blocks Barracks, Dorm, economic buildings, and staffed towers until House creates Peasants.
- The first staffing slice uses available populations to gate tower and economic building construction and reserves staff on successful placement. New games still start at `0/0`, House can add Peasants, Barracks can convert Peasants into Soldiers, Dorm can convert a Peasant into an Apprentice, economic buildings can produce post-Raid resources, and timed recruitment, reassignment, release, broader growth, and losses are not implemented.
- Woodcutter, Stone Quarry, and Iron Mine are the first economic buildings. Each reserves one available Peasant and produces 10 of its matching resource after each defeated Raid. The Iron Mine produces Metal, not a separate Iron resource.

## Open Game Design Questions

- Which resources, if any, are spent to explore an adjacent Plot after the free prototype reveal?
- How do non-north newly explored Plots connect to roads and enemy paths?
- How is the Rival's Lair revealed or signaled to the player?
- What movement, build, resource, and road rules apply to each terrain type?
- Should the entire map be defined at the start of a Fable, or should Plots and Tiles be generated as the wizard explores?
- Should the map be able to expand indefinitely, or should each Fable use bounded map dimensions?
- Should the durable map model be one 2D Tile array, with Plots acting only as convenience groupings over that array?
- If maps can expand after creation, how should the map data structure grow while preserving existing Tile, Plot, road, and Domain state?
- When should scouting and claiming become separate actions after the first implementation collapsed them into one reveal-and-claim action?
- How are Plot dominant characters generated, selected, or presented to the player?
- How many road exits can a Plot have, and can roads branch inside a Plot after entering through edge-center connector Tiles?
- Are Wood, Stone, and Metal enough for interesting choices, or are additional resources needed?
- How are Apprentices, Soldiers, and Peasants recruited, assigned, released from assignments, and affected by Raids?
- Should later exploration add an on-map wizard character, remain camera-inspection driven, or combine both?
- How many arcane barrier charges does the Sanctum have, and can those charges be restored or increased?
- How does the Domain expand, contract, or change over time?
- What enemy archetypes, spawn rules, and pathing rules should Raids use after the first placeholder north-road slice?
- How should future tower stats, damage types, and targeting modes evolve beyond the first Bow Tower baseline?
- What are the first full win and loss conditions beyond basic Raid completion and Sanctum breach?
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
  Rationale: A moderate-damage, moderate-range tower gives the design a simple general-purpose baseline before specialized magical towers are defined.
  Date/Author: 2026-05-08 / Codex

- Decision: Use first prototype construction costs of 30 Wood, 10 Stone, and 10 Metal for the Bow Tower, 30 Stone and 20 Metal for the Flame Bolt Tower, and 40 Wood, 60 Stone, and 25 Metal for the Catapult Tower.
  Rationale: These costs make current tower options visible in the resource economy and now drive the first limited tower placement action, with Catapult positioned as an expensive siege option.
  Date/Author: 2026-05-26 / Codex

- Decision: Use second-based Bow Tower combat stats for the first automated tower slice: 3.0-Tile range, 10 damage, 1.0-second fire interval, and 9.0-Tiles-per-second projectile speed.
  Rationale: Real-time units keep structure stats independent of frame or update counts, and the chosen values made then-current 20-health skeletons die in two hits while keeping projectile travel visible.
  Date/Author: 2026-05-16 / Codex

- Decision: The first Bow Tower targets the in-range enemy closest to the Sanctum and drops projectiles harmlessly if their original target is gone before impact.
  Rationale: Closest-to-Sanctum targeting prioritizes the most urgent threat on the current single road, while no-retarget projectiles keep the first combat model deterministic and easy to test.
  Date/Author: 2026-05-16 / Codex

- Decision: Add the Flame Bolt Tower as the second defined tower type and place one west of the starting road across from the Bow Tower.
  Rationale: A short-range, slower-firing magical damage tower makes the starting defense visibly more wizardly while keeping the first targeting and projectile rules shared and testable.
  Date/Author: 2026-05-16 / Codex

- Decision: Supersede the free starting-tower layout so a new Domain contains only the Sanctum, with 100 Wood, 50 Stone, and 20 Metal available for construction.
  Rationale: Requiring the player to construct the first defense makes the building workflow meaningful. Resources cover Bow and Flame Bolt at the start, while staffing now remains an additional requirement.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Store and display tower staffing requirements without enforcing them until recruitment and assignment exist.
  Rationale: Template metadata establishes the intended requirements and UI language without making every tower unusable while all populations start at `0/0`.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Supersede informational-only staffing by requiring available staff for construction and reserving that staff after a successful build.
  Rationale: Staffing requirements should constrain the number and combination of towers in the same way resources constrain construction. Reducing available counts prevents one inhabitant from staffing unlimited towers, while preserving totals accurately represents assignment rather than population loss.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Add House as the first population-provider building: it costs 20 Wood, requires no staff, has no combat stats, and immediately grants 2 available and total Peasants.
  Rationale: The current staffing gate needs a small normal-play path to create at least one inhabitant role. An immediate Peasant grant proves population-producing structures without adding timed recruitment, reassignment, or removal rules.

- Decision: Add Barracks as the first population-conversion building: it costs 10 Wood and 10 Stone, requires no staff, consumes 2 available and total Peasants, and immediately grants 2 available and total Soldiers.
  Rationale: Soldier staffing needs a normal-play source after House without adding a broad recruitment or assignment system. Treating the effect as conversion makes the population totals match the visible inhabitant roles.
  Date/Author: 2026-06-26 / User and Codex

- Decision: Add Dorm as the first Apprentice-conversion building: it costs 10 Wood and 10 Stone, requires no staff, consumes 1 available and total Peasant, and immediately grants 1 available and total Apprentice.
  Rationale: Apprentice staffing needs a normal-play source after House without adding a broad recruitment or assignment system. Treating the effect as conversion keeps Peasant and Apprentice totals aligned with the visible inhabitant roles.
  Date/Author: 2026-07-13 / User and Codex

- Decision: Add Woodcutter, Stone Quarry, and Iron Mine as the first economic buildings. Each reserves one available Peasant and produces 10 of its matching resource after each defeated Raid. The Iron Mine produces Metal, not a separate Iron resource.
  Rationale: This creates the first repeatable resource production while keeping worker assignment, timing, and resource taxonomy small enough for the current prototype.
  Date/Author: 2026-06-27 / User and Codex

- Decision: Partition the building bar into `Defenses`, `Economic`, and `Housing` tabs, with `Housing` selected by default.
  Rationale: The full building list is no longer convenient in one vertical view, and the default Housing tab keeps the first Peasant and Soldier workflow immediately visible.
  Date/Author: 2026-06-27 / User and Codex

- Decision: Expand the building bar to 260 pixels, show building values to the right of each icon, and draw capacity-blocked icons at 70% opacity.
  Rationale: Horizontal value rows make each building easier to scan, and opacity shows why an option cannot currently be built without changing the construction rules.
  Date/Author: 2026-06-27 / User and Codex

- Decision: Use green outlines for buildable building icon squares and red outlines for capacity-blocked icon squares.
  Rationale: The outline gives a stronger can-build/cannot-build cue while preserving the unfilled slot style and existing placement rules.
  Date/Author: 2026-06-27 / User and Codex

- Decision: Add the Catapult Tower as the third defined tower type and make it buildable from the building bar without placing one for free in the starting Plot.
  Rationale: A long-range, slow-firing, high-damage area tower creates the first clear tower-role contrast beyond single-target projectiles while preserving the existing authored starting defense.
  Date/Author: 2026-05-26 / Codex

- Decision: Catapult projectiles damage every active enemy in the Tile occupied by the original target when the projectile lands.
  Rationale: Tile-area impact matches the grid-based map model and keeps the existing deterministic target selection rule instead of adding a new targeting mode.
  Date/Author: 2026-05-26 / Codex

- Decision: Use deterministic placeholder Raids as the first enemy-wave slice.
  Rationale: A fixed enemy count, fixed stagger, and fixed north-road path make Raid behavior visible and testable before adding pathfinding, rewards, or enemy variety.
  Date/Author: 2026-05-15 / Codex

- Decision: Add zombies as the second enemy archetype in Raid 1.
  Rationale: Spawning zombies in the second and fourth first-Raid positions adds visible enemy variety while keeping the current deterministic north-road Raid model. Zombies are slower and tougher than skeletons, with 75 health and 0.7 Tiles-per-second speed.
  Date/Author: 2026-05-17 / Codex

- Decision: A zero-Barricade Sanctum breach clears the active Raid and prevents further Raid starts in the first slice.
  Rationale: This gives the prototype a concrete failure state without prematurely designing game-over routing, recovery, or campaign consequences.
  Date/Author: 2026-05-15 / Codex

- Decision: Build maps from grid Tiles grouped into 15x15 Plots.
  Rationale: A tile grid gives exploration, building, resources, roads, terrain, and height a shared spatial model while Plots provide a larger grouping for generation, reveal, and Domain-scale reasoning. The 15x15 size preserves an odd dimension with a natural center Tile while giving each Plot more space than the initial 11x11 prototype direction.
  Date/Author: 2026-05-13 / Codex

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
  Rationale: A fixed edge-center connector rule keeps inter-Plot paths readable, takes advantage of the 15x15 Plot centerline, and lets internal road shapes vary without making Plot-to-Plot connectivity ambiguous.
  Date/Author: 2026-05-08 / Codex

- Decision: Use Sanctum-centered Tile-unit world coordinates for gameplay positions.
  Rationale: Storing entity positions directly in world coordinates keeps enemies situated in the map instead of encoding only their progress along one current path. Positive Y points north so `(0, 1.5)` is the shared edge between the first and second Tiles north of the Sanctum.
  Date/Author: 2026-05-16 / User and Codex

- Decision: Use camera zoom and pan as the first map inspection model instead of wizard-character movement.
  Rationale: Camera inspection lets the player examine the starting Plot while keeping exploration, tile selection, resource rules, and character movement out of the first scene-interaction slice.
  Date/Author: 2026-05-13 / Codex

- Decision: Let right-drag panning grab the visible world under the cursor when the drag starts over the game view.
  Rationale: This matches common map-inspection behavior and keeps panning proportional to the current zoom without turning UI surfaces into camera controls.
  Date/Author: 2026-06-27 / Codex

- Decision: Start each Fable with a one-Plot Domain containing the centered Sanctum and an outgoing road.
  Rationale: This gives every playthrough a clear initial defended space and a road hook for future exploration, enemy routing, and Domain expansion.
  Date/Author: 2026-05-08 / Codex

- Decision: Render the first home Plot as open grassland without a tree perimeter.
  Rationale: A grass perimeter keeps the starting Plot fully usable and lets it join explored neighboring grassland without an artificial terrain boundary.
  Date/Author: 2026-07-13 / Codex

- Decision: Store prototype Tile sprite variation as a `Tweak` value on each Tile.
  Rationale: Keeping visual variation in map data lets rendering choose stable per-Tile sprite variants without turning the tweak into a gameplay mechanic.
  Date/Author: 2026-05-13 / Codex

- Decision: In the first static home Plot scene, the outgoing road runs straight north from the centered Sanctum to the north edge-center connector.
  Rationale: A single straight north road gives the first scene a clear route without implying branching, pathfinding, exploration, or combat rules.
  Date/Author: 2026-05-13 / Codex

- Decision: Keep the first rendered home Plot empty except for the Sanctum and north road.
  Rationale: The milestone is meant to prove map state and static scene rendering before terrain, resource, building, or combat rules exist.
  Date/Author: 2026-05-13 / Codex

- Decision: Let the wizard spend resources during calm phases to explore new Plots and unlock building there.
  Rationale: This connects resources, calm-phase planning, Domain expansion, and base-building into one concrete preparation loop.
  Date/Author: 2026-05-08 / Codex

- Decision: Require Plot exploration to be adjacent to the current Domain and make exploration expand the defended path rather than start encounters directly.
  Rationale: This makes exploration a spatial commitment: the wizard gains buildable area, but the Domain can expose a longer route that must be defended during Raids.
  Date/Author: 2026-05-08 / Codex

- Decision: Use the in-game top bar for Chapter, Day, resources, phase status, and Sanctum barricade charges.
  Rationale: These values are the minimum persistent status a player needs while preparing for and surviving Raids, and fixed prototype values can make the display observable before the underlying gameplay systems exist. Inhabitant populations now appear as a separate status group within the same top bar.
  Date/Author: 2026-05-08 / Codex

- Decision: Divide the Domain's inhabitants into Apprentices, Soldiers, and Peasants, and display each population as available/total.
  Rationale: The three roles create readable future pools for magical, military, and general labor assignments while the two values distinguish free capacity from the Domain's full population.
  Date/Author: 2026-06-25 / User and Codex

- Decision: Define Raids as the Rival's end-of-Day assaults on the wizard's Domain.
  Rationale: This gives the defense phase a dedicated design home and ties Raids to Chapters, Rivals, enemy pressure, the defended path, and Sanctum breach conditions.
  Date/Author: 2026-05-08 / Codex

- Decision: Let Chapters last until Domain expansion reaches the Rival's Lair, making the next Raid the final Raid.
  Rationale: This ties Chapter completion to exploration and Domain expansion instead of a fixed day count, so the player controls when to provoke the decisive assault by reaching the Rival's territory.
  Date/Author: 2026-05-08 / Codex
