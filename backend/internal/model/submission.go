package model

import (
	"time"
)

type Submission struct {
	ID           int
	SubmissionID int64
	ProblemID    string
	Result       string
	Code         *string
	SubmittedAt  time.Time
}
