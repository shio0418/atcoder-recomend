package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// フェーズ1: まだスクレイピングしてない問題だけを対象にする（再実行時に済んだ分をやり直さない）
	toScrape, err := repository.ListUnscrapedProblems(ctx, conn, 10000)
	if err != nil {
		log.Fatalf("問題取得失敗: %v", err)
	}

	byContest := make(map[string][]model.Problem)
	for _, p := range toScrape {
		byContest[p.ContestID] = append(byContest[p.ContestID], p)
	}

	for contestID, probs := range byContest {
		links, err := client.FetchEditorialLinks(contestID)
		if err != nil {
			log.Printf("解説一覧取得失敗 (%s): %v（このコンテストは解説無しとして続行）", contestID, err)
		}

		editorialURLs := make(map[string]string)
		for _, l := range links {
			editorialURLs[l.ProblemID] = l.URL
		}

		for _, p := range probs {
			statement, err := client.FetchProblemStatement(contestID, p.ProblemID)
			if err != nil {
				log.Printf("問題文取得失敗 (%s): %v（スキップ）", p.ProblemID, err)
				continue
			}

			var editorialText *string
			hasEditorial := false
			if url, ok := editorialURLs[p.ProblemID]; ok {
				text, err := client.FetchEditorialText(url)
				if err != nil {
					log.Printf("解説取得失敗 (%s): %v（解説無しとして続行）", p.ProblemID, err)
				} else {
					editorialText = &text
					hasEditorial = true
				}
			}

			if err := repository.UpdateProblemScrapedContent(ctx, conn, p.ProblemID, statement, editorialText, hasEditorial); err != nil {
				log.Fatalf("スクレイピング内容保存失敗 (%s): %v", p.ProblemID, err)
			}

			time.Sleep(2 * time.Second)
		}
	}

	// フェーズ2: まだタグが無い問題（今回・過去にスクレイピング済みの分すべて）を分類
	problems, err := repository.ListUntaggedProblems(ctx, conn, 10000)
	if err != nil {
		log.Fatalf("問題取得失敗: %v", err)
	}

	chunkSize := 20
	for i := 0; i < len(problems); i += chunkSize {
		end := min(i+chunkSize, len(problems))
		chunk := problems[i:end]
		results, err := classifier.ClassifyBatch(ctx, chunk)
		if err != nil {
			log.Printf("分類失敗 (%d〜%d件目): %v（このチャンクはスキップ、次回再実行で対象になります）", i, end-1, err)
			continue
		}
		for _, r := range results {
			if err := repository.UpdateProblemTags(ctx, conn, r.ProblemID, r.Tags, r.Probabilities); err != nil {
				log.Fatalf("タグ保存失敗 (%s): %v", r.ProblemID, err)
			}
			fmt.Printf("%s: %v\n", r.ProblemID, r.Tags)
		}
	}
}
