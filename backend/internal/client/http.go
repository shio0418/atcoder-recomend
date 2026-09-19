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

		resp, err := http.Get(url)
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
