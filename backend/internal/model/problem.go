package model

type Problem struct {
	ID         int      `db:"id"`
	ProblemID  string   `db:"problem_id"`
	Title      string   `db:"title"`
	Difficulty *int     `db:"difficulty"`
	Tags       []string `db:"tags"`
}
