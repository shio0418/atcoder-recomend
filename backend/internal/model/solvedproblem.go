package model

import "time"

type SolvedProblem struct {
	ProblemID string    `json:"problemId"`
	ACCount   int       `json:"acCount"`
	LastACAt  time.Time `json:"lastAcAt"`
}
