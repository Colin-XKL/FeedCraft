package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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

func TestPreviewFeedViewerDetectsExternalFeedTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		contentType string
		body        string
		wantType    string
		wantTitle   string
		wantItem    string
	}{
		{
			name:        "rss",
			contentType: "application/rss+xml",
			body: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>RSS Preview Feed</title>
    <link>https://example.com/rss</link>
    <description>RSS feed description</description>
    <item>
      <title>RSS Preview Item</title>
      <link>https://example.com/rss/item</link>
      <guid>rss-item-1</guid>
      <description>RSS item description</description>
    </item>
  </channel>
</rss>`,
			wantType:  "rss",
			wantTitle: "RSS Preview Feed",
			wantItem:  "RSS Preview Item",
		},
		{
			name:        "atom",
			contentType: "application/atom+xml",
			body: `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Atom Preview Feed</title>
  <link href="https://example.com/atom" rel="alternate"/>
  <link href="https://example.com/atom.xml" rel="self"/>
  <id>https://example.com/atom</id>
  <updated>2026-06-29T00:00:00Z</updated>
  <entry>
    <title>Atom Preview Item</title>
    <link href="https://example.com/atom/item"/>
    <id>atom-item-1</id>
    <updated>2026-06-29T00:00:00Z</updated>
    <content type="html">&lt;p&gt;Atom item content&lt;/p&gt;</content>
  </entry>
</feed>`,
			wantType:  "atom",
			wantTitle: "Atom Preview Feed",
			wantItem:  "Atom Preview Item",
		},
		{
			name:        "json feed",
			contentType: "application/feed+json",
			body: `{
  "version": "https://jsonfeed.org/version/1.1",
  "title": "JSON Preview Feed",
  "home_page_url": "https://example.com/json",
  "feed_url": "https://example.com/json-feed.json",
  "items": [
    {
      "id": "json-item-1",
      "url": "https://example.com/json/item",
      "title": "JSON Preview Item",
      "content_html": "<p>JSON item content</p>",
      "date_published": "2026-06-29T00:00:00Z"
    }
  ]
}`,
			wantType:  "json",
			wantTitle: "JSON Preview Feed",
			wantItem:  "JSON Preview Item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			recorder := performFeedViewerPreviewRequest(t, http.MethodGet, server.URL)

			var response struct {
				Code int `json:"code"`
				Data struct {
					FeedType string `json:"feedType"`
					Title    string `json:"title"`
					Items    []struct {
						Title string `json:"title"`
					} `json:"items"`
				} `json:"data"`
				Msg string `json:"msg"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if recorder.Code != http.StatusOK || response.Code != 0 {
				t.Fatalf("expected preview success, got status %d code %d msg %q body %s", recorder.Code, response.Code, response.Msg, recorder.Body.String())
			}
			if response.Data.FeedType != tt.wantType {
				t.Fatalf("feedType = %q, want %q", response.Data.FeedType, tt.wantType)
			}
			if response.Data.Title != tt.wantTitle {
				t.Fatalf("title = %q, want %q", response.Data.Title, tt.wantTitle)
			}
			if len(response.Data.Items) != 1 || response.Data.Items[0].Title != tt.wantItem {
				t.Fatalf("expected one item titled %q, got %+v", tt.wantItem, response.Data.Items)
			}
		})
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

func TestPreviewFeedViewerSupportsRecipeURI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := topicFeedTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}, &dao.TopicFeed{}))
	recipeID := uniqueTestID("preview-recipe")
	createTopicTestRecipe(t, db, recipeID)

	recorder := performFeedViewerPreviewRequest(t, http.MethodGet, "feedcraft://recipe/"+recipeID)

	assertFeedViewerPreviewSuccess(t, recorder, "Test Feed", "Hello Topic")
}

func TestPreviewFeedViewerSupportsTopicURI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := topicFeedTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}, &dao.TopicFeed{}))
	recipeID := uniqueTestID("preview-topic-recipe")
	topicID := uniqueTestID("preview-topic")
	createTopicTestRecipe(t, db, recipeID)
	createTopicTestTopic(t, db, &dao.TopicFeed{
		ID:          topicID,
		Title:       "Preview Topic",
		Description: "Topic preview",
		Inputs:      []dao.TopicInput{{URI: "feedcraft://recipe/" + recipeID}},
	})

	recorder := performFeedViewerPreviewRequest(t, http.MethodGet, "feedcraft://topic/"+topicID)

	assertFeedViewerPreviewSuccess(t, recorder, "Preview Topic", "Hello Topic")
}

func TestPreviewFeedViewerSupportsInboxURI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := topicFeedTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:          "alerts",
		Title:       "Alerts Inbox",
		Description: "Incoming alerts",
		MaxItems:    100,
	}))
	require.NoError(t, db.Create(&dao.InboxItem{
		InboxID:     "alerts",
		ItemID:      "alert-1",
		Title:       "CPU usage high",
		URL:         "https://example.com/alerts/1",
		Content:     "<p>CPU usage exceeded threshold.</p>",
		Summary:     "CPU alert",
		Author:      "monitor",
		PublishedAt: time.Unix(1700000000, 0),
		CreatedAt:   time.Unix(1700000000, 0),
	}).Error)

	recorder := performFeedViewerPreviewRequest(t, http.MethodGet, "feedcraft://inbox/alerts")

	assertFeedViewerPreviewSuccess(t, recorder, "Alerts Inbox", "CPU usage high")
}

func TestPreviewFeedViewerRejectsCraftNameForInternalURI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := performFeedViewerPreviewRequestWithCraft(t, http.MethodGet, "feedcraft://recipe/example", "summary")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !strings.Contains(response.Msg, "craft_name is only supported for external URLs") {
		t.Fatalf("expected craft_name validation message, got %q", response.Msg)
	}
}

func TestNormalizeEmbeddingFilterPreviewRequest(t *testing.T) {
	threshold := 0.72
	maxContentLength := 1500
	cfg, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL:         "http://example.com/feed.xml",
		Anchors:          " AI infrastructure \n\n machine learning ",
		Threshold:        &threshold,
		Mode:             "EXCLUDE",
		MaxContentLength: &maxContentLength,
		Instruction:      "Represent this text for topic filtering",
	})

	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if cfg.inputURL != "http://example.com/feed.xml" {
		t.Fatalf("inputURL = %q", cfg.inputURL)
	}
	if strings.Join(cfg.anchors, "|") != "AI infrastructure|machine learning" {
		t.Fatalf("anchors = %#v", cfg.anchors)
	}
	if cfg.threshold != threshold {
		t.Fatalf("threshold = %f", cfg.threshold)
	}
	if cfg.mode != craft.EmbeddingExcludeMode {
		t.Fatalf("mode = %q", cfg.mode)
	}
	if cfg.maxContentLength != maxContentLength {
		t.Fatalf("maxContentLength = %d", cfg.maxContentLength)
	}
	if cfg.instruction != "Represent this text for topic filtering" {
		t.Fatalf("instruction = %q", cfg.instruction)
	}
}

func TestNormalizeEmbeddingFilterPreviewRequestDefaultsAndValidation(t *testing.T) {
	cfg, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL: "http://example.com/feed.xml",
		Anchors:  "AI",
	})
	if err != nil {
		t.Fatalf("expected defaults to be valid, got %v", err)
	}
	if cfg.threshold != 0.6 {
		t.Fatalf("threshold default = %f", cfg.threshold)
	}
	if cfg.mode != craft.EmbeddingIncludeMode {
		t.Fatalf("mode default = %q", cfg.mode)
	}
	if cfg.maxContentLength != 2000 {
		t.Fatalf("max content length default = %d", cfg.maxContentLength)
	}

	badThreshold := 1.5
	if _, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL:  "http://example.com/feed.xml",
		Anchors:   "AI",
		Threshold: &badThreshold,
	}); err == nil || !strings.Contains(err.Error(), "threshold") {
		t.Fatalf("expected threshold validation error, got %v", err)
	}

	if _, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL: "http://example.com/feed.xml",
		Anchors:  "  \n ",
	}); err == nil || !strings.Contains(err.Error(), "anchors") {
		t.Fatalf("expected anchors validation error, got %v", err)
	}
}

func TestNormalizeEmbeddingFilterPreviewRequestResourceLimits(t *testing.T) {
	tooManyAnchors := make([]string, maxEmbeddingFilterPreviewAnchors+1)
	for i := range tooManyAnchors {
		tooManyAnchors[i] = "anchor"
	}
	if _, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL: "http://example.com/feed.xml",
		Anchors:  strings.Join(tooManyAnchors, "\n"),
	}); err == nil || !strings.Contains(err.Error(), "anchors") {
		t.Fatalf("expected anchors limit error, got %v", err)
	}

	longInstruction := strings.Repeat("x", maxEmbeddingFilterPreviewInstructionLength+1)
	if _, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL:    "http://example.com/feed.xml",
		Anchors:     "AI",
		Instruction: longInstruction,
	}); err == nil || !strings.Contains(err.Error(), "instruction") {
		t.Fatalf("expected instruction limit error, got %v", err)
	}

	tooLongContent := maxEmbeddingFilterPreviewContentLength + 1
	if _, err := normalizeEmbeddingFilterPreviewRequest(EmbeddingFilterPreviewReq{
		InputURL:         "http://example.com/feed.xml",
		Anchors:          "AI",
		MaxContentLength: &tooLongContent,
	}); err == nil || !strings.Contains(err.Error(), "max_content_length") {
		t.Fatalf("expected max_content_length limit error, got %v", err)
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

func TestClassifyFeedViewerErrorExplainsEmbeddingConfiguration(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("[embedding-filter] failed to compute anchor vectors: failed to load embedding config: FC_EMBEDDING_API_MODEL must be set when using FC_EMBEDDING_API_TYPE='ollama'"))

	if status != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
	}
	const want = "Embedding filter is not configured correctly: FC_EMBEDDING_API_MODEL must be set when using FC_EMBEDDING_API_TYPE='ollama'"
	if msg != want {
		t.Fatalf("msg = %q, want %q", msg, want)
	}
}

func TestClassifyFeedViewerErrorDoesNotExposeEmbeddingRuntimeError(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("[embedding-filter] all article embeddings failed: embedding call failed after retries (batch [0-1]): provider returned 500 with token detail"))

	if status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", status, http.StatusInternalServerError)
	}
	const want = "Failed to preview this feed due to an internal error."
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
	return performFeedViewerPreviewRequestWithCraft(t, method, inputURL, "")
}

func performFeedViewerPreviewRequestWithCraft(t *testing.T, method, inputURL, craftName string) *httptest.ResponseRecorder {
	t.Helper()

	router := gin.New()
	router.Handle(method, "/preview", PreviewFeedViewer)

	query := url.Values{}
	query.Set("input_url", inputURL)
	if craftName != "" {
		query.Set("craft_name", craftName)
	}
	request := httptest.NewRequest(method, "/preview?"+query.Encode(), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func assertFeedViewerPreviewSuccess(t *testing.T, recorder *httptest.ResponseRecorder, expectedTitle, expectedItemTitle string) {
	t.Helper()
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
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
	if response.Data.Title != expectedTitle {
		t.Fatalf("expected title %q, got %q", expectedTitle, response.Data.Title)
	}
	if len(response.Data.Items) != 1 || response.Data.Items[0].Title != expectedItemTitle {
		t.Fatalf("expected one item titled %q, got %+v", expectedItemTitle, response.Data.Items)
	}
}
