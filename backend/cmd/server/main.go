package main

import (
	"context"
	"fmt"
	"log"

	"atcoder-recomend/internal/client"
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

	problems, err := client.FetchProblems()
	if err != nil {
		log.Fatalf("問題取得失敗: %v", err)
	}
	for _, p := range problems {
		if err := repository.InsertProblem(ctx, conn, p); err != nil {
			log.Fatalf("問題保存失敗 (%s): %v", p.ProblemID, err)
		}
	}
	fmt.Printf("問題 %d 件を保存しました\n", len(problems))

	submissions, err := client.FetchSubmissions("uwaaaaaa", 0)
	if err != nil {
		log.Fatalf("提出取得失敗: %v", err)
	}
	for _, s := range submissions {
		if err := repository.InsertSubmission(ctx, conn, s); err != nil {
			log.Fatalf("提出保存失敗 (submission_id=%d): %v", s.SubmissionID, err)
		}
	}
	fmt.Printf("提出 %d 件を保存しました\n", len(submissions))
}
