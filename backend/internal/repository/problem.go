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
