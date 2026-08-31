package model

type Problem struct {
	ID        int
	ProblemID string 
	Title string 
	Difficulty *int
	Tags []string
}