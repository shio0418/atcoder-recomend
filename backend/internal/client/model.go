package client

type SubmissionResponse struct {
	ID int64 `json:"id"`
	ProblemID string `json:"problem_id"`
	Result string `json:"result"`
	EpochSecond int64 `json:"epoch_second"`
}

type ProblemResponse struct {
	ID        string `json:"id"`
	ContestID string `json:"contest_id"`
	Title     string `json:"title"`
}

type ProblemModel struct {
	Difficulty *int `json:"difficulty"`
}
