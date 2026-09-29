package repository

import (
	"context"
	"atcoder-recomend/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetSolvedProblems(ctx context.Context, conn DBTX) ([]model.SolvedProblem, error) {
	rows, err := conn.Query(ctx,
		`SELECT problem_id, count(*) AS ac_count, MAX(submitted_at) AS last_ac_at
		FROM submissions
		WHERE result = 'AC'
		GROUP BY problem_id`,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.SolvedProblem])
}