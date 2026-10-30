package repository

import (
	"atcoder-recomend/internal/model"
	"context"
)

func InsertExternalSolve(ctx context.Context, conn DBTX, s model.ExternalSolve) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO external_solves (user_handle, problem_id)
		VALUES ($1, $2)
		ON CONFLICT (user_handle, problem_id) DO NOTHING`,
		s.UserHandle, s.ProblemID,
	)
	return err
}
