package model

import "time"

type Submission struct {
	ID          int
	ProblemID   string
	Result      string
	code        string
	SubmittedAt time.Time
}
