package recommend

import "testing"

func TestForgettingCurveInterval(t *testing.T) {
	tests := []struct {
		name    string
		acCount int
		want    int
	}{
		{"未AC(0回)", 0, 5},
		{"1回目のAC", 1, 5},
		{"2回目のAC", 2, 10},
		{"3回目のAC", 3, 20},
		{"5回目のAC", 5, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForgettingCurveInterval(tt.acCount)
			if got != tt.want {
				t.Errorf("ForgettingCurveInterval(%d) = %d, want %d", tt.acCount, got, tt.want)
			}
		})
	}
}
