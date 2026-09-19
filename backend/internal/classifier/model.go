package classifier

type ClassificationResult struct {
	ProblemID     string             `json:"problem_id"`
	Tags          []string           `json:"tags"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
}
