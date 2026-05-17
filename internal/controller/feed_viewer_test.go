package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPreviewFeedViewerAllowsInternalFeedURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var upstreamMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamMethod = r.Method
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Internal Feed</title>
    <link>http://example.com/</link>
    <description>Internal preview feed</description>
    <item>
      <title>Internal Item</title>
      <link>http://example.com/item</link>
      <guid>item-1</guid>
      <description>Hello from internal network</description>
    </item>
  </channel>
</rss>`))
	}))
	defer server.Close()

	recorder := performFeedViewerPreviewRequest(t, http.MethodGet, server.URL)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if upstreamMethod != http.MethodGet {
		t.Fatalf("expected upstream request method GET, got %s", upstreamMethod)
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			Title string `json:"title"`
			Items []struct {
				Title string `json:"title"`
			} `json:"items"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Code != 0 {
		t.Fatalf("expected success code, got %d with msg %q", response.Code, response.Msg)
	}
	if response.Data.Title != "Internal Feed" {
		t.Fatalf("expected internal feed title, got %q", response.Data.Title)
	}
	if len(response.Data.Items) != 1 || response.Data.Items[0].Title != "Internal Item" {
		t.Fatalf("expected one internal item, got %+v", response.Data.Items)
	}
}

func TestPreviewFeedViewerRejectsNonGETRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstreamCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalled = true
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(`<rss version="2.0"><channel><title>Unexpected</title><link>http://example.com/</link><description>Unexpected</description></channel></rss>`))
	}))
	defer server.Close()

	recorder := performFeedViewerPreviewRequest(t, http.MethodPost, server.URL)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d: %s", http.StatusMethodNotAllowed, recorder.Code, recorder.Body.String())
	}
	if upstreamCalled {
		t.Fatal("expected non-GET request to be rejected before fetching upstream")
	}
}

func TestPreviewFeedViewerDoesNotReturnDataForInvalidFeedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body>not a feed</body></html>`))
	}))
	defer server.Close()

	recorder := performFeedViewerPreviewRequest(t, http.MethodGet, server.URL)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data,omitempty"`
		Msg  string          `json:"msg"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Code != -1 {
		t.Fatalf("expected error code -1, got %d", response.Code)
	}
	if len(response.Data) != 0 {
		t.Fatalf("expected invalid feed response to omit data, got %s", string(response.Data))
	}
	if !strings.Contains(response.Msg, "valid RSS, Atom, or JSON feed") {
		t.Fatalf("expected valid feed error message, got %q", response.Msg)
	}
}

func TestValidateFeedViewerURLAllowsPrivateIPLiteral(t *testing.T) {
	if err := validateFeedViewerURL("http://172.21.0.13/feed.xml"); err != nil {
		t.Fatalf("expected private IP literal to be allowed, got %v", err)
	}
}

func TestFormatFeedViewerValidationErrorPreservesUserFacingCapitalization(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "invalid URL",
			err:  errors.New("please enter a valid http(s) feed URL"),
			want: "Please enter a valid http(s) feed URL",
		},
		{
			name: "private IP",
			err:  errors.New("access to private IP 127.0.0.1 is forbidden"),
			want: "Access to private IP 127.0.0.1 is forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFeedViewerValidationError(tt.err)
			if got != tt.want {
				t.Fatalf("formatFeedViewerValidationError() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifyFeedViewerErrorHandlesLowercaseResolveMessage(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("unable to resolve this URL: lookup example.invalid: no such host"))

	if status != http.StatusOK {
		t.Fatalf("status = %d, want %d", status, http.StatusOK)
	}
	const want = "Unable to fetch this URL. Please check the address and try again."
	if msg != want {
		t.Fatalf("msg = %q, want %q", msg, want)
	}
}

func TestClassifyFeedViewerErrorHandlesBrowserProviderFailures(t *testing.T) {
	tests := []string{
		"browser cdp render failed: context deadline exceeded",
		"browser cdp version request failed: Get \"http://chrome/json/version\": connection refused",
		"failed to decode browser cdp version response: invalid character",
		"browser cdp version response missing webSocketDebuggerUrl",
		"unsupported browser provider \"cdp-typo\"",
	}

	for _, errMsg := range tests {
		status, msg := classifyFeedViewerError(errors.New(errMsg))

		if status != http.StatusOK {
			t.Fatalf("status = %d, want %d for %q", status, http.StatusOK, errMsg)
		}
		const want = "Browser provider failed to render the URL. Please check the address or the browser provider service."
		if msg != want {
			t.Fatalf("msg = %q, want %q for %q", msg, want, errMsg)
		}
	}
}

func performFeedViewerPreviewRequest(t *testing.T, method, inputURL string) *httptest.ResponseRecorder {
	t.Helper()

	router := gin.New()
	router.Handle(method, "/preview", PreviewFeedViewer)

	request := httptest.NewRequest(method, "/preview?input_url="+url.QueryEscape(inputURL), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}
