package main

import (
	"context"
	"fmt"
	"log"

	"atcoder-recomend/internal/classifier"
	"atcoder-recomend/internal/client"
	"atcoder-recomend/internal/model"
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

	byContest := make(map[string][]model.Problem)
	for _, p := range problems {
		byContest[p.ContestID] = append(byContest[p.ContestID], p)
	}

	var enriched []model.Problem
	for contestID, probs := range byContest {
		links, err := client.FetchEditorialLinks(contestID)
		if err != nil {
			log.Fatalf("解説一覧取得失敗 (%s): %v", contestID, err)
		}

		editorialURLs := make(map[string]string)
		for _, l := range links {
			editorialURLs[l.ProblemID] = l.URL
		}

		for _, p := range probs {
			statement, err := client.FetchProblemStatement(contestID, p.ProblemID)
			if err != nil {
				log.Fatalf("問題文取得失敗 (%s): %v", p.ProblemID, err)
			}
			p.Statement = &statement

			if url, ok := editorialURLs[p.ProblemID]; ok {
				text, err := client.FetchEditorialText(url)
				if err != nil {
					log.Fatalf("解説取得失敗 (%s): %v", p.ProblemID, err)
				}
				p.EditorialText = &text
				p.HasEditorial = true
			} else {
				p.HasEditorial = false
			}

			if err := repository.UpdateProblemScrapedContent(ctx, conn, p.ProblemID, statement, p.EditorialText, p.HasEditorial); err != nil {
				log.Fatalf("スクレイピング内容保存失敗 (%s): %v", p.ProblemID, err)
			}

			enriched = append(enriched, p)
		}
	}

	results, err := classifier.ClassifyBatch(ctx, enriched)
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
