package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchProblemStatement fetches a problem's task page and extracts the
// Japanese statement text. Unlike editorials, every problem has one.
func FetchProblemStatement(contestID, problemID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", contestID, problemID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	statement := doc.Find("#task-statement span.lang-ja")
	return strings.TrimSpace(statement.Text()), nil
}
