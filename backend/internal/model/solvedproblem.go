package model

import "time"

type SolvedProblem struct {
	ProblemID string
	ACCount   int
	LastACAt  time.Time
}
