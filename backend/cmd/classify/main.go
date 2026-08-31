package main

import (
	"context"
	"fmt"
	"log"

	"atcoder-recomend/internal/classifier"
	"atcoder-recomend/internal/repository"

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

	problems, err := repository.ListUntaggedProblems(ctx, conn, 5)
	if err != nil {
		log.Fatalf("問題取得失敗: %v", err)
	}

	results, err := classifier.ClassifyBatch(ctx, problems)
	if err != nil {
		log.Fatalf("分類失敗: %v", err)
	}

	for _, r := range results {
		if err := repository.UpdateProblemTags(ctx, conn, r.ProblemID, r.Tags); err != nil {
			log.Fatalf("タグ保存失敗 (%s): %v", r.ProblemID, err)
		}
		fmt.Printf("%s: %v\n", r.ProblemID, r.Tags)
	}
}
