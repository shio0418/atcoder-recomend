package client

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestGetWithRetry(t *testing.T) {
	t.Run("200は1回で成功する", func(t *testing.T) {
		var requests int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		resp, err := getWithRetry(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
		if requests != 1 {
			t.Errorf("requests = %d, want 1", requests)
		}
	})

	t.Run("404は即失敗しリトライしない", func(t *testing.T) {
		var requests int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := getWithRetry(server.URL)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if requests != 1 {
			t.Errorf("requests = %d, want 1 (should not retry on 404)", requests)
		}
	})

	t.Run("429はリトライして最終的に成功する", func(t *testing.T) {
		var requests int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := atomic.AddInt32(&requests, 1)
			if n == 1 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		resp, err := getWithRetry(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
		if requests != 2 {
			t.Errorf("requests = %d, want 2", requests)
		}
	})

	t.Run("500が続くとリトライ上限で諦める", func(t *testing.T) {
		var requests int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, err := getWithRetry(server.URL)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if requests != 3 {
			t.Errorf("requests = %d, want 3 (maxAttempts)", requests)
		}
	})
}
