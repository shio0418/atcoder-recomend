package client

type SubmissionResponse struct {
	ProblemID string `json:"problem_id"`
	Result string `json:"result"`
	EpochSecond int64 `json:"epoch_second"`
}