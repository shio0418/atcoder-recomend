package main

import (
	"atcoder-recomend/internal/client"
	"atcoder-recomend/internal/model"
	"atcoder-recomend/internal/repository"
	"context"
	"log"
	"math/rand"

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

	users, err := client.FetchResults("abc481")
	if err != nil {
		log.Fatalf("ユーザー取得失敗: %v", err)
	}
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	handle := users[:1000]

	for _, user := range handle {
		submissions, err := client.FetchSubmissions(user, 0)
		if err != nil {
			log.Printf("%sの問題取得失敗: %v", user, err)
			continue
		}
		for _, s := range submissions {
			if s.Result != "AC" {
				continue
			}
			if err := repository.InsertExternalSolve(ctx, conn, model.ExternalSolve{
				UserHandle: user,
				ProblemID:  s.ProblemID,
			}); err != nil {
				log.Printf("保存失敗 (%s, %s): %v", user, s.ProblemID, err)
			}
		}
	}
	
}
