package game

import (
	"math"
	"reflect"
	"testing"
)

// TestGenerateRaidBaseline verifies the initial challenge, duration, thresholds, and total.
func TestGenerateRaidBaseline(t *testing.T) {
	template := generateRaid(1, 0, 1)

	if template.challengeRating != 4 {
		t.Fatalf("challenge = %f, want 4", template.challengeRating)
	}
	if template.progressDurationSeconds != 9 {
		t.Fatalf("duration = %f, want 9 seconds", template.progressDurationSeconds)
	}
	wantRules := []raidEnemyRule{
		{kind: raidEnemySkeletonSwordShield, threshold: 2},
		{kind: raidEnemyZombie, threshold: 4},
	}
	if !reflect.DeepEqual(template.enemyRules, wantRules) {
		t.Fatalf("enemy rules = %+v, want %+v", template.enemyRules, wantRules)
	}
	if template.totalEnemies() != 3 {
		t.Fatalf("total enemies = %d, want 3", template.totalEnemies())
	}
}

// TestGenerateRaidChallengeFormula verifies exponential and population scaling.
func TestGenerateRaidChallengeFormula(t *testing.T) {
	template := generateRaid(3, 20, 4)
	want := math.Pow(1.2, 2)*math.Pow(1.2, 3)*(1+20.0/10.0) + 3

	if math.Abs(template.challengeRating-want) > 0.000000001 {
		t.Fatalf("challenge = %.12f, want %.12f", template.challengeRating, want)
	}
	if math.Abs(template.progressDurationSeconds-(5+want)) > 0.000000001 {
		t.Fatalf("duration = %.12f, want %.12f", template.progressDurationSeconds, 5+want)
	}
}

// TestGenerateRaidNormalizesInvalidInputs verifies internal callers cannot reduce the baseline.
func TestGenerateRaidNormalizesInvalidInputs(t *testing.T) {
	got := generateRaid(0, -10, 0)
	want := generateRaid(1, 0, 1)

	if got.challengeRating != want.challengeRating || got.progressDurationSeconds != want.progressDurationSeconds ||
		!reflect.DeepEqual(got.enemyRules, want.enemyRules) {
		t.Fatalf("normalized template = %+v, want %+v", got, want)
	}
}

// TestScheduledEnemyCountUsesReachedMultiples verifies exact boundaries and invalid rules.
func TestScheduledEnemyCountUsesReachedMultiples(t *testing.T) {
	rule := raidEnemyRule{kind: raidEnemySkeletonSwordShield, threshold: 2}
	tests := []struct {
		score float64
		want  int
	}{
		{score: -1, want: 0},
		{score: 0, want: 0},
		{score: 1.999, want: 0},
		{score: 2, want: 1},
		{score: 5.9, want: 2},
	}
	for _, test := range tests {
		if got := scheduledEnemyCount(test.score, rule); got != test.want {
			t.Fatalf("scheduledEnemyCount(%f) = %d, want %d", test.score, got, test.want)
		}
	}
	if got := scheduledEnemyCount(10, raidEnemyRule{}); got != 0 {
		t.Fatalf("zero-threshold count = %d, want 0", got)
	}
}

// TestGenerateRaidIsDeterministic verifies identical inputs return identical content.
func TestGenerateRaidIsDeterministic(t *testing.T) {
	first := generateRaid(7, 13, 5)
	second := generateRaid(7, 13, 5)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("generated templates differ: %+v != %+v", first, second)
	}
}
