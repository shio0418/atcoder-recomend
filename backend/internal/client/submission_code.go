package client

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func FetchSubmissionCode(contestID string, submissionID int64) (string, error) {
	resp, err := getWithRetry(fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%d", contestID, submissionID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	code := doc.Find("#submission-code").Text()
	if code == "" {
		return "", fmt.Errorf("提出コードを抽出できませんでした")
	}
	return code, nil
}
