package classifier

import (
	"sort"
	"testing"

	"atcoder-recomend/internal/model"
)

func TestSelectTags(t *testing.T) {
	tests := []struct {
		name          string
		probabilities map[string]float64
		want          []string
	}{
		{
			name:          "閾値以上のタグだけ選ばれる",
			probabilities: map[string]float64{"dp": 0.9, "greedy": 0.2, "graph": 0.03},
			want:          []string{"dp", "greedy"},
		},
		{
			name:          "閾値ちょうど(0.05)は含まれる",
			probabilities: map[string]float64{"dp": probThreshold},
			want:          []string{"dp"},
		},
		{
			name:          "閾値未満は含まれない",
			probabilities: map[string]float64{"dp": 0.049},
			want:          nil,
		},
		{
			name:          "空のprobabilities",
			probabilities: map[string]float64{},
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectTags(tt.probabilities)
			sort.Strings(got)
			sort.Strings(tt.want)
			if len(got) != len(tt.want) {
				t.Fatalf("selectTags() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("selectTags() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBuildJevRequest(t *testing.T) {
	statement := "問題文です"

	t.Run("editorialが無い場合はstateに含めない", func(t *testing.T) {
		p := model.Problem{
			ProblemID:    "abc001_a",
			Title:        "問題A",
			Statement:    &statement,
			HasEditorial: false,
		}

		req := buildJevRequest(p)

		if req.Model != "jev-latest" {
			t.Errorf("Model = %q, want %q", req.Model, "jev-latest")
		}
		if req.State["title"] != "問題A" {
			t.Errorf("State[title] = %q, want %q", req.State["title"], "問題A")
		}
		if req.State["statement"] != statement {
			t.Errorf("State[statement] = %q, want %q", req.State["statement"], statement)
		}
		if _, ok := req.State["editorial_text"]; ok {
			t.Errorf("State[editorial_text] should not be set when HasEditorial is false")
		}
		if _, ok := req.Questions["tags"]; !ok {
			t.Errorf("Questions should contain a \"tags\" question")
		}
	})

	t.Run("editorialがある場合はstateに含める", func(t *testing.T) {
		editorial := "解説です"
		p := model.Problem{
			ProblemID:     "abc001_a",
			Title:         "問題A",
			Statement:     &statement,
			HasEditorial:  true,
			EditorialText: &editorial,
		}

		req := buildJevRequest(p)

		if req.State["editorial_text"] != editorial {
			t.Errorf("State[editorial_text] = %q, want %q", req.State["editorial_text"], editorial)
		}
	})
}
