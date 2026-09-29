package repository

import (
	"context"
	"atcoder-recomend/internal/model"

	"github.com/jackc/pgx/v5"
)

func GetTagStats(ctx context.Context, conn DBTX) ([]model.TagStat, error) {
	rows, err := conn.Query(ctx,
		`SELECT tag, count(*) AS attempts, count(*) FILTER (WHERE result = 'AC') AS ac_count FROM submissions
		JOIN problems ON submissions.problem_id = problems.problem_id
		CROSS JOIN unnest(problems.tags) AS tag
		GROUP BY tag`,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.TagStat])
}