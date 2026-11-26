package repository

import "context"

func InsertRecommendationFeedback(ctx context.Context, conn DBTX, problem_id string, source string, feedback string) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO recommendation_feedback (problem_id, source, feedback)
		VALUES ($1, $2, $3)`,
		problem_id, source, feedback,
	)
	return err
}
