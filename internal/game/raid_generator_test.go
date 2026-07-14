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
		{kind: raidEnemyGhoul, threshold: 6},
		{kind: raidEnemyArmouredSkeleton, threshold: 8},
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
	wantDuration := 5 + 2*math.Sqrt(want)
	if math.Abs(template.progressDurationSeconds-wantDuration) > 0.000000001 {
		t.Fatalf("duration = %.12f, want %.12f", template.progressDurationSeconds, wantDuration)
	}
}

// TestRaidProgressDurationPreservesBaselineAndCompressesLaterRaids verifies the accepted curve.
func TestRaidProgressDurationPreservesBaselineAndCompressesLaterRaids(t *testing.T) {
	tests := []struct {
		challenge float64
		want      float64
	}{
		{challenge: 4, want: 9},
		{challenge: 16, want: 13},
		{challenge: 36, want: 17},
	}
	for _, test := range tests {
		if got := raidProgressDuration(test.challenge); got != test.want {
			t.Fatalf("duration for challenge %.0f = %f, want %f", test.challenge, got, test.want)
		}
	}
}

// TestRaidProgressDurationIncreasesTempo verifies later rosters get less time per enemy.
func TestRaidProgressDurationIncreasesTempo(t *testing.T) {
	rules := []raidEnemyRule{
		{kind: raidEnemySkeletonSwordShield, threshold: 2},
		{kind: raidEnemyZombie, threshold: 4},
		{kind: raidEnemyGhoul, threshold: 6},
		{kind: raidEnemyArmouredSkeleton, threshold: 8},
	}
	challenges := []float64{4, 16, 36}
	previousDuration := 0.0
	previousSecondsPerEnemy := math.Inf(1)
	for _, challenge := range challenges {
		template := raidTemplate{
			challengeRating:         challenge,
			progressDurationSeconds: raidProgressDuration(challenge),
			enemyRules:              rules,
		}
		secondsPerEnemy := template.progressDurationSeconds / float64(template.totalEnemies())
		if template.progressDurationSeconds <= previousDuration {
			t.Fatalf("duration for challenge %.0f = %f, want more than %f", challenge, template.progressDurationSeconds, previousDuration)
		}
		if secondsPerEnemy >= previousSecondsPerEnemy {
			t.Fatalf("seconds per enemy for challenge %.0f = %f, want less than %f", challenge, secondsPerEnemy, previousSecondsPerEnemy)
		}
		previousDuration = template.progressDurationSeconds
		previousSecondsPerEnemy = secondsPerEnemy
	}
}

// TestRaidRosterCountsAtChallenges verifies the accepted four-rule totals.
func TestRaidRosterCountsAtChallenges(t *testing.T) {
	template := generateRaid(1, 0, 1)
	tests := []struct {
		challenge  float64
		wantCounts []int
		wantTotal  int
	}{
		{challenge: 8, wantCounts: []int{4, 2, 1, 1}, wantTotal: 8},
		{challenge: 16, wantCounts: []int{8, 4, 2, 2}, wantTotal: 16},
	}

	for _, test := range tests {
		template.challengeRating = test.challenge
		for i, rule := range template.enemyRules {
			if got := scheduledEnemyCount(test.challenge, rule); got != test.wantCounts[i] {
				t.Fatalf("challenge %.0f rule %d count = %d, want %d", test.challenge, i, got, test.wantCounts[i])
			}
		}
		if got := template.totalEnemies(); got != test.wantTotal {
			t.Fatalf("challenge %.0f total = %d, want %d", test.challenge, got, test.wantTotal)
		}
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
