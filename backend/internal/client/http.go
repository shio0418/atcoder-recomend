package client

import (
	"fmt"
	"net/http"
	"time"
)

// getWithRetry retries a GET request on transient failures (connection
// resets, timeouts, 5xx responses) with exponential backoff. A long-running
// scrape over thousands of pages will eventually hit a blip; without this,
// a single one kills the whole batch.
func getWithRetry(url string) (*http.Response, error) {
	const maxAttempts = 3
	backoff := time.Second

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		// AtCoderはデフォルトのGo User-Agent(Go-http-client/1.1)だと
		// 403を返してブロックしてくる(特に/submissions/配下)ので、
		// ブラウザっぽいUser-Agentを送る。
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; atcoder-recomend/1.0)")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			resp.Body.Close()
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
				return nil, lastErr
			}
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("giving up after %d attempts: %w", maxAttempts, lastErr)
}
