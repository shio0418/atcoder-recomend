package client

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type EditorialLink struct {
	ProblemID string
	URL       string
}

// FetchEditorialLinks scrapes a contest's editorial index page and returns,
// for each problem that has one, the URL of its official editorial
// (preferring the Japanese version when both exist). Unofficial editorials
// (別解/ユーザ解説) and external links (YouTube, blogs) are skipped since
// they either add noise or can't be scraped as plain text.
func FetchEditorialLinks(contestID string) ([]EditorialLink, error) {
	resp, err := getWithRetry(fmt.Sprintf("https://atcoder.jp/contests/%s/editorial?lang=ja", contestID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	taskLinkPattern := regexp.MustCompile(`^/contests/` + regexp.QuoteMeta(contestID) + `/tasks/([a-z0-9_]+)$`)
	editorialLinkPattern := regexp.MustCompile(`^/contests/` + regexp.QuoteMeta(contestID) + `/editorial/\d+$`)

	var links []EditorialLink
	doc.Find("h3").Each(func(_ int, h3 *goquery.Selection) {
		taskHref, ok := h3.Find("a.small").Attr("href")
		if !ok {
			return
		}
		m := taskLinkPattern.FindStringSubmatch(taskHref)
		if m == nil {
			return
		}
		problemID := m[1]

		section := h3.Next()
		if !section.HasClass("editorial-section") {
			return
		}

		var chosenHref string
		section.Find("li").Each(func(_ int, li *goquery.Selection) {
			if li.Find(".label-default").Length() == 0 {
				return
			}
			a := li.Find("a").First()
			href, ok := a.Attr("href")
			if !ok || !editorialLinkPattern.MatchString(href) {
				return
			}
			if chosenHref == "" || strings.TrimSpace(a.Text()) == "解説" {
				chosenHref = href
			}
		})

		if chosenHref != "" {
			links = append(links, EditorialLink{
				ProblemID: problemID,
				URL:       "https://atcoder.jp" + chosenHref,
			})
		}
	})

	return links, nil
}

// FetchEditorialText fetches an editorial page and extracts its body text.
func FetchEditorialText(editorialURL string) (string, error) {
	resp, err := getWithRetry(editorialURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	body := doc.Find("hr.mt-1").Next()
	body.Find("details").Remove()
	body.Find("pre").Remove()
	return strings.TrimSpace(body.Text()), nil
}
