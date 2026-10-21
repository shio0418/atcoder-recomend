package main

import (
	"atcoder-recomend/internal/embedding"
	"atcoder-recomend/internal/model"
	"atcoder-recomend/internal/repository"
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env読み込み失敗: %v", err)
	}
	ctx := context.Background()
	conn, err := repository.Connect(ctx)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	toEmbed, err := repository.ListUnembeddedProblems(ctx, conn, 20)
	if err != nil {
		log.Fatalf("問題取得失敗: %v", err)
	}

	texts := make(map[string]string)
	for _, p := range toEmbed {
		texts[p.ProblemID] = *p.Statement
	}

	start := time.Now()
	embedResp, err := embedding.EmbedBatch(ctx, texts)
	latency_ms := int(time.Since(start).Milliseconds())
	if err != nil {
		log.Fatalf("embedding失敗: %v", err)
	}

	item_count := len(embedResp)
	for pid, resp := range embedResp {
		if err := repository.UpdateProblemsEmbedding(ctx, conn, pid, resp.Values, resp.ModelVersion); err != nil {
			log.Fatalf("db書き込み失敗: %v", err)
		}
	}
	if err := repository.InsertInferenceCall(ctx, conn, model.InferenceCall{
		CallType:  "embedding",
		LatencyMs: latency_ms,
		Success:   true,
		ItemCount: &item_count,
	}); err != nil {
		log.Fatalf("ログ記録失敗: %v", err)
	}
}
