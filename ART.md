# Art Asset Guidance: td

`ART.md` records durable advice for generating art assets for `td`. Update it when the project adopts new generated-asset guidance, prompt patterns, asset review criteria, or production constraints for art files.

This file is guidance for asset creation, not a final art bible. `DESIGN.md` remains the source of truth for the overall medieval wizardry visual direction and interaction principles. `ART.md` explains how generated art should support that direction without prematurely choosing a final medium, asset pipeline, animation system, or production workflow.

## Visual Goals

run python3 rather than python
Generated assets are 2D pixel art.  
PNG format.
Clearly visible, even if scaled down to 50% size.
Very sharp contrast, low detail, and bright colours.  
No lighting or shadows.

## Dimensions

Terrain: 64x64
Structure: 64x64
Icon: 64x64
Enemy: 64x64

## Location

Generated graphical assets should be placed in a subfolder under the `assets/sprites` folder.  The name of the subfolder should reflect the purpose of the asset: terrains, structures, enemies, or icons.

## Current Generated Asset Notes

The House structure sprite lives at `assets/sprites/structures/house.png`. It was generated as a centered low-detail medieval timber-and-stone cottage on a flat chroma-key background, converted to alpha, and resized to 64x64 PNG.

The Barracks structure sprite lives at `assets/sprites/structures/barracks.png`. It was generated as a centered low-detail medieval timber-and-stone barracks on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized to 64x64 PNG.

The Dorm structure sprite lives at `assets/sprites/structures/dorm.png`. It was generated as a centered low-detail medieval wizardry dormitory with timber-and-stone construction, blue roof, lit windows, and small apprentice study details on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized to 64x64 PNG.

The Woodcutter structure sprite lives at `assets/sprites/structures/woodcutter.png`. It was generated as a centered low-detail medieval lumber-camp hut with stacked logs and woodcutting props on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized to 64x64 PNG.

The Stone Quarry structure sprite lives at `assets/sprites/structures/stone-quarry.png`. It was generated as a centered low-detail medieval quarry worksite with a timber hoist, shed roof, and stone blocks on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized to 64x64 PNG.

The Iron Mine structure sprite lives at `assets/sprites/structures/iron-mine.png`. It was generated as a centered low-detail medieval mine entrance with timber supports, a winch, and bluish metal ore on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized to 64x64 PNG. The building produces the existing Metal resource even though the structure name uses Iron.

Future structure prompts should continue to ask for a single isolated readable subject, no scenery, no text, no watermark, no cast shadow, and no background after processing.

The population portraits live at `assets/sprites/icons/apprentice.png`, `assets/sprites/icons/soldier.png`, and `assets/sprites/icons/peasant.png`. They are front-facing head-and-shoulders pixel-art portraits inside a shared thin bronze circular rim with a dark charcoal badge interior and transparent exterior. The Apprentice has a surprised face and purple pointed wizard hat, the Soldier has a grimacing face and steel cap, and the Peasant has a smiling face and straw hat. These portraits replace full-body figures everywhere the population icons appear.

Future population portrait prompts should preserve that shared round-badge composition, keep the face dominant, and use the role's hat shape, palette, and expression as the primary identifiers. Ask for a perfectly flat solid `#00ff00` background outside the badge, no green in the subject, no full body, props, scenery, text, shadow, reflection, or watermark. Remove the chroma key with automatic border sampling, a soft matte, and despill; crop to the nontransparent bounds; resize with nearest-neighbor sampling to fit about 60x60; and center the result on an exact 64x64 transparent PNG canvas. Review the expression at the 28-pixel top-HUD size and the hat, palette, and silhouette at the 14-pixel building-metadata size before accepting a portrait.

The Boulder terrain sprites live at `assets/sprites/terrains/boulder-1.png` through `assets/sprites/terrains/boulder-4.png`. They were generated as isolated low-detail gray pixel-art boulders on a flat chroma-key background, converted to alpha, cropped to subject bounds, and resized onto 64x64 transparent PNG canvases. Future terrain obstacle sprites should match the existing pine-tree and Boulder terrain style: isolated subject, transparent background, clear silhouette, no ground patch, no cast shadow, and readable at 50% scale.

The Iron Deposit terrain sprites live at `assets/sprites/terrains/iron-deposit-1.png` through `assets/sprites/terrains/iron-deposit-4.png`. They were generated from the four Boulder silhouettes as isolated low-detail gray pixel-art rocks with sparse, high-contrast cobalt-blue and cyan-blue ore veins on flat chroma-key backgrounds. The backgrounds were removed locally, and each subject was cropped, resized with nearest-neighbor sampling to the same 54-pixel-wide footprint as the Boulder family, and centered on an exact 64x64 transparent PNG canvas. Future mineral terrain should stay visibly related to ordinary Boulder terrain while using readable blue mineral seams as its distinguishing feature at 50% scale. Preserve the isolated natural-rock treatment and avoid mine structures, gathering props, protruding crystals, shadows, ground patches, or magical glow.

## Open Questions

- When should the project introduce an asset pipeline, sprite atlas, animation workflow, or custom font?
- How should generated art be licensed, attributed, reviewed, and replaced before any public release?
