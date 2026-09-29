package recommend

import "testing"

func TestDifficultyBand(t *testing.T) {
	tests := []struct {
		name     string
		rating   int
		wantMin  int
		wantMax  int
	}{
		{"通常のレート", 800, 800, 1000},
		{"レート0", 0, 0, 200},
		{"レート200", 200, 200, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := DifficultyBand(tt.rating)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("DifficultyBand(%d) = %d, %d, wantMin %d, wantMax %d", tt.rating, gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}
