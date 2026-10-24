package repository

import (
	"atcoder-recomend/internal/model"
	"context"
)

func InsertInferenceCall(ctx context.Context, conn DBTX, c model.InferenceCall) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO inference_calls (call_type, latency_ms, success, item_count, token_count, error_message)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		c.CallType, c.LatencyMs, c.Success, c.ItemCount, c.TokenCount, c.ErrorMessage,
	)
	return err
}
