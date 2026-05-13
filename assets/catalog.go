package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed sprites/structures/sanctum.png
var spriteFiles embed.FS

// Catalog groups all loaded runtime assets by type and subtype.
type Catalog struct {
	Sprite SpriteCatalog
}

// SpriteCatalog groups loaded image sprites by game-domain subtype.
type SpriteCatalog struct {
	Structure StructureSprites
}

// StructureSprites groups loaded sprites for map and base structures.
type StructureSprites struct {
	Sanctum *ebiten.Image
}

// NewCatalog loads the runtime assets required by a new game.
func NewCatalog() (Catalog, error) {
	sanctum, err := loadSprite("sprites/structures/sanctum.png")
	if err != nil {
		return Catalog{}, err
	}
	return Catalog{
		Sprite: SpriteCatalog{
			Structure: StructureSprites{
				Sanctum: sanctum,
			},
		},
	}, nil
}

// loadSprite decodes an embedded image file and converts it for Ebitengine drawing.
func loadSprite(path string) (*ebiten.Image, error) {
	data, err := spriteFiles.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sprite %q: %w", path, err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode sprite %q: %w", path, err)
	}
	return ebiten.NewImageFromImage(img), nil
}
