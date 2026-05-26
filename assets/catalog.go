package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed sprites/enemies/skeleton-sword-shield.png sprites/enemies/zombie.png sprites/icons/metal.png sprites/icons/stone.png sprites/icons/wood.png sprites/structures/bow-tower-projectile.png sprites/structures/bow-tower.png sprites/structures/flame-bolt-tower-projectile.png sprites/structures/flame-bolt-tower.png sprites/structures/sanctum.png sprites/terrains/pine-tree-1.png sprites/terrains/pine-tree-2.png sprites/terrains/pine-tree-3.png sprites/terrains/pine-tree-4.png
var spriteFiles embed.FS

//go:embed audio/raider-defeated.wav
var audioFiles embed.FS

// Catalog groups all loaded runtime assets by type and subtype.
type Catalog struct {
	Sprite SpriteCatalog
	Audio  AudioCatalog
}

// AudioCatalog groups loaded audio effect bytes by game-domain event.
type AudioCatalog struct {
	RaiderDefeated []byte
}

// SpriteCatalog groups loaded image sprites by game-domain subtype.
type SpriteCatalog struct {
	Enemy      EnemySprites
	Icon       IconSprites
	Projectile ProjectileSprites
	Structure  StructureSprites
	Terrain    TerrainSprites
}

// EnemySprites groups loaded sprites for enemy units.
type EnemySprites struct {
	SkeletonSwordShield *ebiten.Image
	Zombie              *ebiten.Image
}

// IconSprites groups loaded sprites for HUD icons.
type IconSprites struct {
	Wood  *ebiten.Image
	Stone *ebiten.Image
	Metal *ebiten.Image
}

// ProjectileSprites groups loaded sprites for projectiles fired by combat structures.
type ProjectileSprites struct {
	BowTowerProjectile       *ebiten.Image
	FlameBoltTowerProjectile *ebiten.Image
}

// StructureSprites groups loaded sprites for map and base structures.
type StructureSprites struct {
	Sanctum        *ebiten.Image
	BowTower       *ebiten.Image
	FlameBoltTower *ebiten.Image
}

// TerrainSprites groups loaded sprites for map terrain.
type TerrainSprites struct {
	PineTrees [4]*ebiten.Image
}

// NewCatalog loads the runtime assets required by a new game.
func NewCatalog() (Catalog, error) {
	audioCatalog, err := NewAudioCatalog()
	if err != nil {
		return Catalog{}, err
	}
	skeletonSwordShield, err := loadSprite("sprites/enemies/skeleton-sword-shield.png")
	if err != nil {
		return Catalog{}, err
	}
	zombie, err := loadSprite("sprites/enemies/zombie.png")
	if err != nil {
		return Catalog{}, err
	}
	woodIcon, err := loadSprite("sprites/icons/wood.png")
	if err != nil {
		return Catalog{}, err
	}
	stoneIcon, err := loadSprite("sprites/icons/stone.png")
	if err != nil {
		return Catalog{}, err
	}
	metalIcon, err := loadSprite("sprites/icons/metal.png")
	if err != nil {
		return Catalog{}, err
	}
	sanctum, err := loadSprite("sprites/structures/sanctum.png")
	if err != nil {
		return Catalog{}, err
	}
	bowTower, err := loadSprite("sprites/structures/bow-tower.png")
	if err != nil {
		return Catalog{}, err
	}
	bowTowerProjectile, err := loadSprite("sprites/structures/bow-tower-projectile.png")
	if err != nil {
		return Catalog{}, err
	}
	flameBoltTower, err := loadSprite("sprites/structures/flame-bolt-tower.png")
	if err != nil {
		return Catalog{}, err
	}
	flameBoltTowerProjectile, err := loadSprite("sprites/structures/flame-bolt-tower-projectile.png")
	if err != nil {
		return Catalog{}, err
	}
	var pineTrees [4]*ebiten.Image
	for i := range pineTrees {
		path := fmt.Sprintf("sprites/terrains/pine-tree-%d.png", i+1)
		pineTrees[i], err = loadSprite(path)
		if err != nil {
			return Catalog{}, err
		}
	}
	return Catalog{
		Audio: audioCatalog,
		Sprite: SpriteCatalog{
			Enemy: EnemySprites{
				SkeletonSwordShield: skeletonSwordShield,
				Zombie:              zombie,
			},
			Icon: IconSprites{
				Wood:  woodIcon,
				Stone: stoneIcon,
				Metal: metalIcon,
			},
			Projectile: ProjectileSprites{
				BowTowerProjectile:       bowTowerProjectile,
				FlameBoltTowerProjectile: flameBoltTowerProjectile,
			},
			Structure: StructureSprites{
				Sanctum:        sanctum,
				BowTower:       bowTower,
				FlameBoltTower: flameBoltTower,
			},
			Terrain: TerrainSprites{
				PineTrees: pineTrees,
			},
		},
	}, nil
}

// NewAudioCatalog loads embedded audio assets required by runtime sound playback.
func NewAudioCatalog() (AudioCatalog, error) {
	raiderDefeated, err := loadAudio("audio/raider-defeated.wav")
	if err != nil {
		return AudioCatalog{}, err
	}
	return AudioCatalog{
		RaiderDefeated: raiderDefeated,
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

// loadAudio reads embedded audio bytes for later decoding by the sound package.
func loadAudio(path string) ([]byte, error) {
	data, err := audioFiles.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read audio %q: %w", path, err)
	}
	return data, nil
}
