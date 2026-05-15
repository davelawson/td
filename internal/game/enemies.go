package game

// EnemyTemplate describes shared stats for every instance of one enemy type.
type EnemyTemplate struct {
	Name          string
	MaxHealth     int
	Speed         float64
	SanctumDamage int
	SpriteKey     string
}

// EnemyCatalog groups every enemy template available to game systems.
type EnemyCatalog struct {
	SkeletonSwordShield EnemyTemplate
}

// NewEnemyCatalog creates the default enemy template catalog.
func NewEnemyCatalog() EnemyCatalog {
	return EnemyCatalog{
		SkeletonSwordShield: EnemyTemplate{
			Name:          "Skeleton Sword-and-Shield",
			MaxHealth:     20,
			Speed:         3.0,
			SanctumDamage: 1,
			SpriteKey:     "skeleton-sword-shield",
		},
	}
}
