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
- Static map scenes should keep tile boundaries, roads, and the Sanctum readable before adding decorative terrain or production art.
- Selected map objects should have a plainly brighter visual state that is easy to notice without obscuring the sprite, health bar, road, or tile context.
- Camera zoom and pan should affect the map scene only. The top HUD, debug counter, pause label, and in-game overlay should remain readable screen-space UI, and the overlay should block camera controls while it is open.
- Quit behavior must be obvious and should close the app cleanly.
- Keyboard access should be considered early, even when the first slice only requires pointer input.

## Visual Review

Any plan that changes rendered game output should define visual evidence before implementation starts. If no prior app exists, the plan should record that there is no screenshotable baseline and capture the first rendered result after implementation.

For Ebitengine desktop work, visual evidence can be a screenshot saved under the active plan directory, such as `plans/00-initial-ebitengine-menu/screenshots/main-menu.png`. Review screenshots against this file and note any usability issues in the active ExecPlan before accepting the work.

## Open Questions

- Should the game eventually use pixel art, painterly 2D, or another asset style?
- How should the early camera scaling policy evolve when the world grows beyond one static home Plot?
- Should early menus use custom font assets, or wait until gameplay systems exist?
