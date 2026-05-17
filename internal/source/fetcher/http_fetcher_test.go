package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/util"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpFetcherUsesDefaultFeedUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	fetcher := &HttpFetcher{Config: &config.HttpFetcherConfig{URL: server.URL}}
	_, err := fetcher.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if gotUA != util.DefaultFeedUserAgent() {
		t.Fatalf("expected default feed user agent, got %q", gotUA)
	}
}

func TestHttpFetcherAllowsHeaderOverride(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	fetcher := &HttpFetcher{Config: &config.HttpFetcherConfig{
		URL: server.URL,
		Headers: map[string]string{
			"User-Agent": "SourceSpecific/1.2.3",
		},
	}}
	_, err := fetcher.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if gotUA != "SourceSpecific/1.2.3" {
		t.Fatalf("expected source header override, got %q", gotUA)
	}
}

func TestHttpFetcherHTMLPurposeUsesDefaultHTMLUserAgent(t *testing.T) {
	var gotUA, gotAccept string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		gotAccept = r.Header.Get("Accept")
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	fetcher := &HttpFetcher{Config: &config.HttpFetcherConfig{
		URL:     server.URL,
		Purpose: config.HttpFetcherPurposeHTML,
	}}
	_, err := fetcher.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if gotUA != util.DefaultHTMLUserAgent() {
		t.Fatalf("expected html default user agent, got %q", gotUA)
	}
	if gotAccept == "" {
		t.Fatal("expected html accept header to be set")
	}
}

func TestHttpFetcherHTMLPurposeRetriesRetryableStatus(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	fetcher := &HttpFetcher{Config: &config.HttpFetcherConfig{
		URL:     server.URL,
		Purpose: config.HttpFetcherPurposeHTML,
	}}
	_, err := fetcher.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestHttpFetcherFeedPurposeDoesNotRetryOnRetryableStatus(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	fetcher := &HttpFetcher{Config: &config.HttpFetcherConfig{
		URL:     server.URL,
		Purpose: config.HttpFetcherPurposeFeed,
	}}
	_, err := fetcher.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected fetch to fail")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}
