package repository

import (
	"atcoder-recomend/internal/model"
	"context"
)

func InsertSubmission(ctx context.Context, conn DBTX, s model.Submission) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO submissions (submission_id, problem_id, result, submitted_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (submission_id) DO NOTHING`,
		s.SubmissionID, s.ProblemID, s.Result, s.SubmittedAt,
	)
	return err
}
