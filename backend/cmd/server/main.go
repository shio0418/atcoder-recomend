package main

import (
	"fmt"
	"log"

	"atcoder-recomend/internal/client"
)

func main() {
	submissions, err := client.FetchSubmissions("uwaaaaaa", 0)
	if err != nil {
		log.Fatalf("取得失敗: %v", err)
	}

	fmt.Printf("件数: %d\n", len(submissions))
	if len(submissions) > 0 {
		fmt.Printf("最初の1件: %+v\n", submissions[0])
	}
}
