package repository

import (
	"context"
	"atcoder-recomend/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetSolvedProblems(ctx context.Context, conn DBTX) ([]model.SolvedProblem, error) {
	rows, err := conn.Query(ctx,
		`SELECT submissions.problem_id,
			count(*) FILTER (WHERE result = 'AC') AS ac_count,
			count(*) FILTER (WHERE result != 'AC') AS wa_count,
			MAX(submitted_at) FILTER (WHERE result = 'AC') AS last_ac_at,
			problems.difficulty
		FROM submissions
		JOIN problems ON submissions.problem_id = problems.problem_id
		GROUP BY submissions.problem_id, problems.difficulty
		HAVING count(*) FILTER (WHERE result = 'AC') > 0`,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.SolvedProblem])
}