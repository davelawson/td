package game

import (
	"image/color"
	"testing"
)

// TestRaidEnemyHealthFractionClampsToCurrentHealth verifies health-bar proportions.
func TestRaidEnemyHealthFractionClampsToCurrentHealth(t *testing.T) {
	template := &EnemyTemplate{MaxHealth: 20}
	tests := []struct {
		name  string
		enemy raidEnemy
		want  float64
	}{
		{
			name:  "full health",
			enemy: raidEnemy{template: template, health: 20},
			want:  1,
		},
		{
			name:  "half health",
			enemy: raidEnemy{template: template, health: 10},
			want:  0.5,
		},
		{
			name:  "zero health",
			enemy: raidEnemy{template: template, health: 0},
			want:  0,
		},
		{
			name:  "negative health",
			enemy: raidEnemy{template: template, health: -5},
			want:  0,
		},
		{
			name:  "over max health",
			enemy: raidEnemy{template: template, health: 25},
			want:  1,
		},
		{
			name:  "missing template",
			enemy: raidEnemy{health: 10},
			want:  1,
		},
		{
			name:  "invalid max health",
			enemy: raidEnemy{template: &EnemyTemplate{MaxHealth: 0}, health: 10},
			want:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := raidEnemyHealthFraction(test.enemy); got != test.want {
				t.Fatalf("raidEnemyHealthFraction() = %f, want %f", got, test.want)
			}
		})
	}
}

// TestRaidEnemyHealthBarColorInterpolatesGreenToRed verifies bar color endpoints.
func TestRaidEnemyHealthBarColorInterpolatesGreenToRed(t *testing.T) {
	tests := []struct {
		name     string
		fraction float64
		want     color.RGBA
	}{
		{
			name:     "full health green",
			fraction: 1,
			want:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			name:     "half health yellow",
			fraction: 0.5,
			want:     color.RGBA{R: 128, G: 128, B: 0, A: 255},
		},
		{
			name:     "zero health red",
			fraction: 0,
			want:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "below zero clamps red",
			fraction: -1,
			want:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "above one clamps green",
			fraction: 2,
			want:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := raidEnemyHealthBarColor(test.fraction); got != test.want {
				t.Fatalf("raidEnemyHealthBarColor() = %#v, want %#v", got, test.want)
			}
		})
	}
}
