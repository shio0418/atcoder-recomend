package recommend

import (
	"atcoder-recomend/internal/model"
)

func MasteryScores(tagStats []model.TagStat) map[string]float64 {
	scores := make(map[string]float64)
	for _, ts := range tagStats {
		if ts.Attempts == 0 {
			continue
		}
		scores[ts.Tag] = float64(ts.ACCount) / float64(ts.Attempts)
	}
	return scores
}