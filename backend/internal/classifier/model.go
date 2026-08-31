package classifier

type ClassificationResult struct {
	ProblemID string   `json:"problem_id"`
	Tags      []string `json:"tags"`
}
