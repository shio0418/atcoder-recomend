package recommend

import (
	"testing"
	"time"

	"atcoder-recomend/internal/model"
)

func intPtr(v int) *int {
	return &v
}

func TestReviewCandidates(t *testing.T) {
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	solved := []model.SolvedProblem{
		{
			// 1回AC、間隔は5日。19日前にACしてるので復習対象。
			ProblemID: "due",
			ACCount:   1,
			LastACAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			// 1回AC、間隔は5日。2日前にACしたばかりなのでまだ対象外。
			ProblemID: "not_due",
			ACCount:   1,
			LastACAt:  time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC),
		},
	}

	got := ReviewCandidates(solved, now, 5, 2, 1000, 90)

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1 (got: %+v)", len(got), got)
	}
	if got[0].ProblemID != "due" {
		t.Errorf("got[0].ProblemID = %q, want %q", got[0].ProblemID, "due")
	}
}

func TestNewCandidates(t *testing.T) {
	tagStats := []model.TagStat{
		{Tag: "dp", Attempts: 10, ACCount: 2},      // AC率20% → 弱点
		{Tag: "greedy", Attempts: 10, ACCount: 8},  // AC率80% → 弱点ではない
	}
	rating := 800 // バンドは800〜1000になる

	candidates := []model.Problem{
		{ProblemID: "in_band_weak", Difficulty: intPtr(900), Tags: []string{"dp"}},
		{ProblemID: "in_band_not_weak", Difficulty: intPtr(900), Tags: []string{"greedy"}},
		{ProblemID: "out_of_band_weak", Difficulty: intPtr(1200), Tags: []string{"dp"}},
		{ProblemID: "no_difficulty", Difficulty: nil, Tags: []string{"dp"}},
		{ProblemID: "band_boundary_weak", Difficulty: intPtr(800), Tags: []string{"greedy", "dp"}},
	}

	got := NewCandidates(tagStats, rating, candidates)

	gotIDs := make(map[string]bool)
	for _, p := range got {
		gotIDs[p.ProblemID] = true
	}

	wantIDs := []string{"in_band_weak", "band_boundary_weak"}
	if len(got) != len(wantIDs) {
		t.Fatalf("len(got) = %d, want %d (got: %+v)", len(got), len(wantIDs), got)
	}
	for _, id := range wantIDs {
		if !gotIDs[id] {
			t.Errorf("expected %q to be in result, but it wasn't", id)
		}
	}
}
