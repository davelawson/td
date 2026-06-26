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

The House structure sprite lives at `assets/sprites/structures/house.png`. It was generated as a centered low-detail medieval timber-and-stone cottage on a flat chroma-key background, converted to alpha, and resized to 64x64 PNG. Future structure prompts should continue to ask for a single isolated readable subject, no scenery, no text, no watermark, no cast shadow, and no background after processing.

## Open Questions

- When should the project introduce an asset pipeline, sprite atlas, animation workflow, or custom font?
- How should generated art be licensed, attributed, reviewed, and replaced before any public release?
