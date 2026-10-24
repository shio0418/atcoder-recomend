package recommend

import "testing"

func TestForgettingCurveInterval(t *testing.T) {
	tests := []struct {
		name    string
		acCount int
		waCount int
		want    int
	}{
		{"未AC(0回)、WA無し", 0, 0, 5},
		{"1回目のAC、WA無し", 1, 0, 5},
		{"2回目のAC、WA無し", 2, 0, 10},
		{"3回目のAC、WA無し", 3, 0, 20},
		{"5回目のAC、WA無し", 5, 0, 80},
		{"1回目のAC、WA1回", 1, 1, 2},
		{"1回目のAC、WA4回", 1, 4, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForgettingCurveInterval(tt.acCount, tt.waCount)
			if got != tt.want {
				t.Errorf("ForgettingCurveInterval(%d, %d) = %d, want %d", tt.acCount, tt.waCount, got, tt.want)
			}
		})
	}
}
