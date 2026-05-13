# Art Asset Guidance: td

`ART.md` records durable advice for generating art assets for `td`. Update it when the project adopts new generated-asset guidance, prompt patterns, asset review criteria, or production constraints for art files.

This file is guidance for asset creation, not a final art bible. `DESIGN.md` remains the source of truth for the overall medieval wizardry visual direction and interaction principles. `ART.md` explains how generated art should support that direction without prematurely choosing a final medium, asset pipeline, animation system, or production workflow.

## Current Asset Status

The prototype currently uses simple drawn shapes and system text. Static runtime assets are not part of the implemented game yet, and there is no asset pipeline, packaging setup, sprite atlas, animation system, custom font, or final production-art style.

Generated art may be useful for future mockups, reference images, temporary sprites, textures, icons, terrain studies, and mood exploration. Treat those generated outputs as prototype inputs until a plan explicitly promotes them into runtime assets.

## Visual Goals

Generated assets should support a medieval wizardry tower-defense game about a wizard defending a Sanctum within a hostile Domain. Useful motifs include old stone, towers, roads, parchment, runes, spellcraft, arcane light, moss, woods, mountains, water, resource nodes, hostile wilds, and hand-built defenses.

The tone should be strategic, readable, and slightly arcane. Avoid comedic exaggeration, modern military language, science fiction machinery, corporate interface styling, and purely abstract decoration unless a plan explicitly calls for a contrast or placeholder.

## Prompting Guidance

When prompting for generated art, state the gameplay purpose before decorative details. Good prompts should identify what the asset must communicate at game scale, whether it should read as map terrain, a structure, a resource, a UI element, a reference image, or a mood study.

Prefer prompts that ask for clear silhouettes, readable shapes, restrained magical accents, and a limited palette that can coexist with the project colors in `DESIGN.md`. For map-facing assets, mention that tile boundaries, road direction, structure footprints, and enemy paths must remain easy to read.

For temporary runtime assets, request orthographic or top-down views when the asset must sit on the game map. Request transparent backgrounds only when the asset is meant to be composited over the game scene or UI. For reference art and mood studies, backgrounds are acceptable if they help evaluate atmosphere, lighting, or material choices.

## Asset Review Checklist

Review generated art before adding it to the repository or using it as a source for implementation:

- The asset supports medieval wizardry fantasy and does not drift into science fiction, modern military, generic corporate, or unrelated fantasy language.
- The asset remains readable at the size it will appear in the prototype.
- Important gameplay information is not hidden by texture noise, glow effects, heavy shadows, or excessive detail.
- The palette fits the restrained direction in `DESIGN.md` and does not make the whole interface or map read as one dominant hue family.
- Terrain, roads, structures, resources, and UI marks have distinct silhouettes.
- Lighting and perspective are consistent enough that nearby assets could plausibly share one scene.
- Temporary art is clearly treated as temporary unless an accepted plan promotes it into a longer-lived asset.

## Repository Use

Do not add generated art files casually. A plan that introduces runtime assets should specify the intended file paths under `assets/`, how those files are loaded or referenced, what visual evidence will be captured, and how the assets can be replaced later without disrupting gameplay code.

Keep source prompts or generation notes when they are useful for reproducing or revising an asset. Store those notes near the active plan evidence or in another planned location rather than scattering them through source code comments.

Generated images should not define gameplay rules. If an asset implies a new terrain type, tower, enemy, resource, interaction, or UI workflow, update the appropriate control document such as `GAME.md`, `PRODUCT.md`, `ROADMAP.md`, or `DESIGN.md` in the same planned change.

## Open Questions

- Should the final game use pixel art, painterly 2D, illustrated board-game style, or another asset style?
- What base dimensions should tiles, structures, enemies, icons, and UI art use?
- Which file formats should be preferred for runtime assets?
- How should asset filenames, prompt notes, source references, and generated variants be organized?
- When should the project introduce an asset pipeline, sprite atlas, animation workflow, or custom font?
- How should generated art be licensed, attributed, reviewed, and replaced before any public release?
