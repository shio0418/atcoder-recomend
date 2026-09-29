package model

type Problem struct {
	ID            int      `db:"id" json:"id"`
	ProblemID     string   `db:"problem_id" json:"problemId"`
	Title         string   `db:"title" json:"title"`
	Difficulty    *int     `db:"difficulty" json:"difficulty"`
	Tags          []string `db:"tags" json:"tags"`
	ContestID     string   `db:"contest_id" json:"contestId"`
	HasEditorial  bool     `db:"has_editorial" json:"hasEditorial"`
	Statement     *string  `db:"statement" json:"statement,omitempty"`
	EditorialText *string  `db:"editorial_text" json:"editorialText,omitempty"`
}