package game

import "math"

const (
	raidChallengeGrowthBase            = 1.2
	raidProgressDurationBase           = 5.0
	raidProgressDurationChallengeScale = 2.0
)

type raidEnemyKind int

const (
	raidEnemySkeletonSwordShield raidEnemyKind = iota
	raidEnemyZombie
	raidEnemyGhoul
	raidEnemyArmouredSkeleton
)

type raidEnemyRule struct {
	kind      raidEnemyKind
	threshold int
}

type raidTemplate struct {
	challengeRating         float64
	progressDurationSeconds float64
	enemyRules              []raidEnemyRule
}

// generateRaid creates a deterministic Raid template from progression and settlement state.
func generateRaid(raidNumber, settlementPopulation, plotsExplored int) raidTemplate {
	if raidNumber < 1 {
		raidNumber = 1
	}
	if settlementPopulation < 0 {
		settlementPopulation = 0
	}
	if plotsExplored < 1 {
		plotsExplored = 1
	}

	challengeRating := math.Pow(raidChallengeGrowthBase, float64(raidNumber-1))*
		math.Pow(raidChallengeGrowthBase, float64(plotsExplored-1))*
		(1+float64(settlementPopulation)/10) + 3
	return raidTemplate{
		challengeRating:         challengeRating,
		progressDurationSeconds: raidProgressDuration(challengeRating),
		enemyRules: []raidEnemyRule{
			{kind: raidEnemySkeletonSwordShield, threshold: 2},
			{kind: raidEnemyZombie, threshold: 4},
			{kind: raidEnemyGhoul, threshold: 6},
			{kind: raidEnemyArmouredSkeleton, threshold: 8},
		},
	}
}

// raidProgressDuration returns the release window for a generated challenge.
func raidProgressDuration(challengeRating float64) float64 {
	return raidProgressDurationBase + raidProgressDurationChallengeScale*math.Sqrt(challengeRating)
}

// scheduledEnemyCount returns how many enemies one rule has scheduled at a challenge score.
func scheduledEnemyCount(score float64, rule raidEnemyRule) int {
	if score <= 0 || rule.threshold <= 0 {
		return 0
	}
	return int(math.Floor(score / float64(rule.threshold)))
}

// totalEnemies returns the number of enemies scheduled by the full Raid template.
func (t raidTemplate) totalEnemies() int {
	total := 0
	for _, rule := range t.enemyRules {
		total += scheduledEnemyCount(t.challengeRating, rule)
	}
	return total
}
