package recommend

import (
	"math"
)

func ForgettingCurveInterval(acCount int, waCount int) int {
	// TODO: d_i := 5 * 2^(i-1) - f(p, i) のfの項も実装してみる
	return 5 * int(math.Pow(2.0, max(float64(acCount)-1, 0))) / (1 + waCount)
}