package client

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchProblemStatement fetches a problem's task page and extracts the
// Japanese statement text. Unlike editorials, every problem has one.
func FetchProblemStatement(contestID, problemID string) (string, error) {
	resp, err := getWithRetry(fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", contestID, problemID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Modern contests wrap the statement in a lang-specific span so both
	// languages can coexist in the DOM. Older/single-language contests
	// don't have that split at all, so fall back to the whole block.
	text := doc.Find("#task-statement span.lang-ja")
	if text.Length() == 0 {
		text = doc.Find("#task-statement").First()
	}
	text.Find("h3:contains('例')").Parent().Remove()
	result := strings.TrimSpace(text.Text())
	if result == "" {
		return "", fmt.Errorf("問題文を抽出できませんでした")
	}
	return result, nil
}
