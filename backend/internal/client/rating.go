package client

import (
	"fmt"
	"encoding/json"
)

func FetchRating(userID string) (int, error) {
	url := fmt.Sprintf("https://atcoder.jp/users/%s/history/json", userID)
	resp, err := getWithRetry(url)

	if err != nil {
		return 0, err
	}

	type ratingHistoryEntry struct {
		NewRating int `json:"NewRating"`
	}
	var history []ratingHistoryEntry
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return 0, err
	}
	if len(history) == 0 {
		return 0, nil
	}

	ratio := history[len(history)-1].NewRating
	return ratio, nil
}