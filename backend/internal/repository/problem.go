package repository

import (
	"atcoder-recomend/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertProblem(ctx context.Context, conn *pgx.Conn, p model.Problem) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO problems (problem_id, title, difficulty, tags, contest_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (problem_id) DO UPDATE SET contest_id = EXCLUDED.contest_id`,
		p.ProblemID, p.Title, p.Difficulty, p.Tags, p.ContestID,
	)
	return err
}

func ListUntaggedProblems(ctx context.Context, conn *pgx.Conn, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, title, difficulty, contest_id, statement, editorial_text, has_editorial
		FROM problems WHERE tags IS NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

// ListUnscrapedProblems returns problems that still need their statement
// (and editorial, if any) scraped. Separate from ListUntaggedProblems so a
// restart after a crash mid-scrape doesn't re-scrape problems that already
// succeeded — only tags are missing for those.
func ListUnscrapedProblems(ctx context.Context, conn *pgx.Conn, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, contest_id FROM problems WHERE tags IS NULL AND statement IS NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

func UpdateProblemTags(ctx context.Context, conn *pgx.Conn, problemID string, tags []string, probabilities map[string]float64) error {
	_, err := conn.Exec(ctx,
		`UPDATE problems SET tags = $1, tag_probabilities = $2 WHERE problem_id = $3`,
		tags, probabilities, problemID,
	)
	return err
}

func UpdateProblemScrapedContent(ctx context.Context, conn *pgx.Conn, problemID string, statement string, editorialText *string, hasEditorial bool) error {
	_, err := conn.Exec(ctx,
		`UPDATE problems SET statement = $1, editorial_text = $2, has_editorial = $3 WHERE problem_id = $4`,
		statement, editorialText, hasEditorial, problemID,
	)
	return err
}

