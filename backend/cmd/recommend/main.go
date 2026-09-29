package main

import (
	"atcoder-recomend/internal/client"
	"atcoder-recomend/internal/recommend"
	"atcoder-recomend/internal/repository"
	"context"
	"fmt"
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

	rating, err := client.FetchRating("uwaaaaaaa")
	if err != nil {
		log.Fatalf("rating取得失敗: %v", err)
	}

	tagStats, err := repository.GetTagStats(ctx, conn)
	if err != nil {
		log.Fatalf("tagStats取得失敗: %v", err)
	}

	solvedProblems, err := repository.GetSolvedProblems(ctx, conn)
	if err != nil {
		log.Fatalf("solvedProblems取得失敗: %v", err)
	}

	candidates, err := repository.ListCandidateProblems(ctx, conn)
	if err != nil {
		log.Fatalf("candidates取得失敗: %v", err)
	}

	reviewCandidates := recommend.ReviewCandidates(solvedProblems, time.Now())
	newCandidates := recommend.NewCandidates(tagStats, rating, candidates)

	for _, reviewCandidate := range reviewCandidates {
		fmt.Println(reviewCandidate)
	}
	for _, newCandidate := range newCandidates {
		fmt.Println(newCandidate)
	}
}
