package repository

import (
	"atcoder-recomend/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertProblem(ctx context.Context, conn *pgx.Conn, p model.Problem) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO problems (problem_id, title, difficulty, tags)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (problem_id) DO NOTHING`,
		p.ProblemID, p.Title, p.Difficulty, p.Tags,
	)
	return err
}

func ListUntaggedProblems(ctx context.Context, conn *pgx.Conn, limit int) ([]model.Problem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, title, difficulty FROM problems WHERE tags IS NULL LIMIT $1`,
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

