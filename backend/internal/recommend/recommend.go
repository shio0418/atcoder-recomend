package recommend

import (
	"atcoder-recomend/internal/model"
	"slices"
	"time"
)

// TODO: 苦手に含める閾値は適当に決めたので後で調整してもいいかも
const THRESHOLD = 0.5

func ReviewCandidates(solved []model.SolvedProblem, now time.Time, baseInterval float64, growthRate float64, difficultyWeight float64) []model.SolvedProblem {
	var result []model.SolvedProblem
	for _, problem := range solved {
		interval := ForgettingCurveInterval(problem.ACCount, problem.WACount, problem.Difficulty, baseInterval, growthRate, difficultyWeight)
		dueAt := problem.LastACAt.AddDate(0, 0, interval)
		if now.After(dueAt) {
			problem.DaysOverdue = now.Sub(dueAt).Hours() / 24
			result = append(result, problem)
		}
	}
	return result
}

func NewCandidates(tagStats []model.TagStat, rating int, candidates []model.Problem) []model.Problem {
	min, max := DifficultyBand(rating)
	var result []model.Problem
	weakTagSet := make(map[string]bool)
	for _, ts := range tagStats {
		if float64(ts.ACCount)/float64(ts.Attempts) < THRESHOLD {
			weakTagSet[ts.Tag] = true
		}
	}

	for _, candidate := range candidates {
		if candidate.Difficulty == nil {
			continue
		}
		if *candidate.Difficulty < min || *candidate.Difficulty > max {
			continue
		}
		hasWeakTag := slices.ContainsFunc(candidate.Tags, func(tag string) bool {
			return weakTagSet[tag]
		})
		if hasWeakTag {
			result = append(result, candidate)
		}
	}
	return result
}
