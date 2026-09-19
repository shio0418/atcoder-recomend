package model

type Problem struct {
	ID            int      `db:"id"`
	ProblemID     string   `db:"problem_id"`
	Title         string   `db:"title"`
	Difficulty    *int     `db:"difficulty"`
	Tags          []string `db:"tags"`
	ContestID     string   `db:"contest_id"`
	HasEditorial  bool     `db:"has_editorial"`
	Statement     *string  `db:"statement"`
	EditorialText *string  `db:"editorial_text"`
}
