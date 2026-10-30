package client

import (
	"encoding/json"
	"fmt"
)

func FetchResults(contestID string) ([]string, error) {
	resp, err := getWithRetry(fmt.Sprintf("https://atcoder.jp/contests/%s/results/json", contestID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//body, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(body)[:300]) // 最初の300文字だけ見る

	var result []struct {
		UserScreenName string `json:"UserScreenName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	handles := make([]string, 0, len(result))
	for _, s := range result {
		handles = append(handles, s.UserScreenName)
	}
	return handles, nil
}
