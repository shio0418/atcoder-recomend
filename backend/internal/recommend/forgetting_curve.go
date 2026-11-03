package recommend

import (
	"math"
)

func ForgettingCurveInterval(acCount int, waCount int, difficulty *int, baseInterval float64, growthRate float64, difficultyWeight float64) int {
	growth := baseInterval * math.Pow(growthRate, max(float64(acCount)-1, 0))
	difficultyFactor := 1.0
	if difficulty != nil {
		difficultyFactor = 1 + float64(*difficulty)/difficultyWeight
	}
	return int(growth / (1 + float64(waCount)) / difficultyFactor)
}