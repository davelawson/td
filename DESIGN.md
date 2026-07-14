# Design Direction: td

`DESIGN.md` records the durable visual and interaction direction for `td`. Update it when the product adopts a new art direction, UI density, accessibility expectation, or visual review requirement.

## Visual Identity

`td` should read as medieval wizardry fantasy. The visual language should support spellcraft, towers, old stone, parchment, runes, magical light, hostile wilds, and hand-built defenses. Avoid sci-fi, modern military, generic corporate UI, and purely abstract placeholder language once the prototype moves beyond the first technical slice.

The first implementation may use simple drawn shapes and system text, but it should still gesture toward the intended identity through color, naming, and composition.

## Tone

The game should feel strategic, readable, and slightly arcane rather than loud or comedic. Menus should be clear and usable first, with fantasy styling applied in a restrained way that does not obscure interaction.

## Initial Palette Guidance

Use this as a starting point, not a final art bible:

- Deep charcoal or near-black for the main background.
- Muted parchment or warm off-white for primary text.
- Desaturated gold for focused or highlighted actions.
- Deep green, moss, or cool stone neutrals for secondary surfaces.
- Violet or blue magical accents sparingly, not as the dominant palette.

Avoid letting the whole interface become a single purple, beige, or dark-blue palette.

## UI and Interaction Principles

- The first screen should be the actual game menu, not a marketing or explanation page.
- Menu choices should have stable hit boxes and clear hover or focus feedback.
- Text must remain readable at the target window size. The current desktop target opens at 1920x1080, and menu text should remain raw-pixel-sized when the window is enlarged instead of stretching with the window.
- In-game HUD text should be compact, high-contrast, and restrained. It should expose essential play status without covering the defended field or competing with the in-game overlay menu.
- Population indicators should use a consistent family of expressive, front-facing portrait badges rather than full-body figures. A thin bronze round rim and dark interior should unify the set, while each role remains identifiable through its face, hat shape, and palette. Expressions should read in the 28-pixel top HUD; the role silhouette and colors should remain distinct in 14-pixel building metadata.
- Static map scenes should keep tile boundaries, roads, and the Sanctum readable before adding decorative terrain or production art. The home Plot and explored grasslands use the same sparse Tree and Boulder generation, so the starting view should read as grasslands without obscuring its Sanctum, protected north road, exploration controls, or buildable empty grass. Generated hills should read as more Boulder-heavy than grasslands while remaining sparse and mostly buildable. Shared explored-Plot joins should remain visually open, and Boulder terrain should use transparent terrain sprites that feel consistent with the pine-tree terrain assets rather than vector placeholder shapes.
- Exploration affordances should be legible as map controls. The current magnifying-glass buttons sit on borders between explored and unexplored orthogonal Plots, and they should remain visually distinct from Tiles, structures, roads, and screen-space UI. Each button's preassigned `Grasslands` or `Hills` label should use compact fixed-size text placed outward into unexplored space: above north, right of east, below south, and left of west. The label is informational; only the circular magnifying-glass is clickable.
- Adjacent explored Plots should render as one continuous field. Do not draw per-Plot frames, gutters, or padding once both sides of a border are explored.
- Selected raiders and structures should have a plainly brighter visual state that is easy to notice without obscuring the sprite, health bar, road, or Tile context. Selected Tree and Boulder terrain should use a gold Tile outline so the terrain sprite and biome remain readable.
- Selected-object detail panels should stay compact, anchored in screen space, and readable without becoming command surfaces before commands exist.
- Building and command bars should stay screen-space, compact, and visually subordinate to the map. Build affordances must only imply actions that are implemented; the current building bar appears only during Management and supports informational hover tooltips, affordability-aware hover, and a simple half-sized dragged icon for placement, but not upgrades, selling, range previews, or broader construction commands. During Raid and breach, the hidden bar must not reserve an invisible input region. The top HUD should name Management plainly and append the next challenge preview; Raid status should retain the generated challenge. The instantaneous Labour phase needs no transient panel or animation in the current slice.
- Camera zoom and pan should affect the map scene only. Right-drag panning should start only from the game view, not from the top HUD, building bar, command buttons, selection panel, or overlay. The top HUD, debug counter, pause label, and in-game overlay should remain readable screen-space UI, and the overlay should block camera controls while it is open.
- Quit behavior must be obvious and should close the app cleanly.
- Keyboard access should be considered early, even when the first slice only requires pointer input.

## Visual Review

Any plan that changes rendered game output should define visual evidence before implementation starts. If no prior app exists, the plan should record that there is no screenshotable baseline and capture the first rendered result after implementation.

For Ebitengine desktop work, visual evidence can be a screenshot saved under the active plan directory, such as `plans/00-initial-ebitengine-menu/screenshots/main-menu.png`. Review screenshots against this file and note any usability issues in the active ExecPlan before accepting the work.

## Open Questions

- Should the game eventually use pixel art, painterly 2D, or another asset style?
- How should the early camera scaling policy evolve as exploration creates larger Domains?
- Should early menus use custom font assets, or wait until gameplay systems exist?
