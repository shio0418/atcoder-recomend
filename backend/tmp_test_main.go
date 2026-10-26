// backend/tmp_test_main.go
package main

import (
	"atcoder-recomend/internal/client"
	"fmt"
)

func main() {
	handles, err := client.FetchResults("abc481")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("件数:", len(handles))
	for i, h := range handles {
		if i >= 10 {
			break
		}
		fmt.Println(h)
	}
}
