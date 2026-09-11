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
		ON CONFLICT (problem_id) DO NOTHING`,
		p.ProblemID, p.Title, p.Difficulty, p.Tags, p.ContestID,
	)
	return err
}

func ListUntaggedProblems(ctx context.Context, conn *pgx.Conn, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, title, difficulty, contest_id FROM problems WHERE tags IS NULL LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Problem])
}

func UpdateProblemTags(ctx context.Context, conn *pgx.Conn, problemID string, tags []string) error {
	_, err := conn.Exec(ctx,
		`UPDATE problems SET tags = $1 WHERE problem_id = $2`,
		tags, problemID,
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

