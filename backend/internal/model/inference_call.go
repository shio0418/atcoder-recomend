package model

type InferenceCall struct {
	CallType     string
	LatencyMs    int
	Success      bool
	ItemCount    *int
	TokenCount   *int
	ErrorMessage *string
}
