package repository

import (
	"atcoder-recomend/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertSubmission(ctx context.Context, conn *pgx.Conn, s model.Submission) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO submissions (submission_id, problem_id, result, submitted_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (submission_id) DO NOTHING`,
		s.SubmissionID, s.ProblemID, s.Result, s.SubmittedAt,
	)
	return err
}
