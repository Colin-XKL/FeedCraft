package craft

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func TestNegotiateFeedFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		query  string
		accept string
		want   string
	}{
		{
			name:   "explicit html wins over rss accept",
			query:  "html",
			accept: "application/rss+xml",
			want:   "html",
		},
		{
			name:   "explicit json",
			query:  "json",
			accept: "text/html",
			want:   "json",
		},
		{
			name:   "explicit xml aliases to rss",
			query:  "xml",
			accept: "text/html",
			want:   "rss",
		},
		{
			name:   "unknown format falls back to accept",
			query:  "nope",
			accept: "text/html",
			want:   "html",
		},
		{
			name:   "curl default accept stays rss",
			query:  "",
			accept: "*/*",
			want:   "rss",
		},
		{
			name:   "empty accept stays rss",
			query:  "",
			accept: "",
			want:   "rss",
		},
		{
			name:   "chrome-like accept renders html",
			query:  "",
			accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			want:   "html",
		},
		{
			name:   "rss reader accept stays rss",
			query:  "",
			accept: "application/rss+xml, application/rdf+xml;q=0.8, application/atom+xml;q=0.6, application/xml;q=0.4, text/xml;q=0.4, */*;q=0.2",
			want:   "rss",
		},
		{
			name:   "tie between html and rss prefers rss",
			query:  "",
			accept: "text/html, application/rss+xml",
			want:   "rss",
		},
		{
			name:   "json clients get json feed",
			query:  "",
			accept: "application/json, text/plain, */*",
			want:   "json",
		},
		{
			name:   "generic xml without html stays rss",
			query:  "",
			accept: "application/xml",
			want:   "rss",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := negotiateFeedFormat(tt.query, tt.accept); got != tt.want {
				t.Fatalf("negotiateFeedFormat(%q, %q) = %q, want %q", tt.query, tt.accept, got, tt.want)
			}
		})
	}
}

func TestCommonCraftHandlerServesRSSForReaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := newTestRSSServer(t)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL, nil)
	c.Request.Header.Set("Accept", "application/rss+xml, application/xml;q=0.9, */*;q=0.8")

	CommonCraftHandlerUsingCraftOptionList(c, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "rss+xml") && !strings.Contains(contentType, "xml") {
		t.Fatalf("expected RSS XML content type, got %q", contentType)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "<rss") {
		t.Fatalf("expected RSS XML body, got %s", body)
	}
	if !strings.Contains(body, "Browser Preview Item") {
		t.Fatalf("expected feed item in RSS body, got %s", body)
	}
}

func TestCommonCraftHandlerRendersHTMLForBrowsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := newTestRSSServer(t)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL, nil)
	c.Request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	CommonCraftHandlerUsingCraftOptionList(c, nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected HTML content type for browsers, got %q", contentType)
	}
	body := recorder.Body.String()
	if strings.Contains(body, "<rss") {
		t.Fatalf("browser preview should not be raw RSS XML, got %s", body)
	}
	if !strings.Contains(body, "Browser Preview Item") {
		t.Fatalf("expected feed title in HTML preview, got %s", body)
	}
	if !strings.Contains(body, "format=rss") {
		t.Fatalf("expected a link to the raw RSS representation, got %s", body)
	}
}

func TestCommonCraftHandlerFormatQueryOverridesAccept(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := newTestRSSServer(t)

	t.Run("format=rss forces xml for browsers", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL+"&format=rss", nil)
		c.Request.Header.Set("Accept", "text/html")

		CommonCraftHandlerUsingCraftOptionList(c, nil)

		if !strings.Contains(recorder.Header().Get("Content-Type"), "xml") {
			t.Fatalf("expected XML content type, got %q", recorder.Header().Get("Content-Type"))
		}
		if !strings.Contains(recorder.Body.String(), "<rss") {
			t.Fatalf("expected RSS body, got %s", recorder.Body.String())
		}
	})

	t.Run("format=json returns json feed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/craft/proxy?input_url="+source.URL+"&format=json", nil)
		c.Request.Header.Set("Accept", "text/html")

		CommonCraftHandlerUsingCraftOptionList(c, nil)

		contentType := recorder.Header().Get("Content-Type")
		if !strings.Contains(contentType, "json") {
			t.Fatalf("expected JSON content type, got %q", contentType)
		}
		body := recorder.Body.String()
		if !strings.Contains(body, `"items"`) && !strings.Contains(body, `"title"`) {
			t.Fatalf("expected JSON Feed document, got %s", body)
		}
		if !strings.Contains(body, "Browser Preview Item") {
			t.Fatalf("expected feed item in JSON body, got %s", body)
		}
	})
}

func TestRenderFeedHTMLEscapesUntrustedContent(t *testing.T) {
	t.Parallel()

	feed := &feeds.Feed{
		Title:       `<script>alert("feed")</script>`,
		Description: `<img src=x onerror=alert(1)>`,
		Items: []*feeds.Item{
			{
				Title:       `<script>alert("item")</script>`,
				Description: `<a href="javascript:alert(1)">click</a>`,
				Link:        &feeds.Link{Href: "https://example.com/item"},
			},
		},
	}

	html, err := renderFeedHTML(feed, "https://feedcraft.example/craft/proxy?input_url=https://example.com/feed")
	if err != nil {
		t.Fatalf("renderFeedHTML returned error: %v", err)
	}
	if strings.Contains(html, `<script>alert("feed")</script>`) ||
		strings.Contains(html, `<script>alert("item")</script>`) ||
		strings.Contains(html, "onerror=") ||
		strings.Contains(html, "javascript:") {
		t.Fatalf("HTML preview leaked unsanitized markup: %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatalf("expected escaped title content, got %s", html)
	}
}

func newTestRSSServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Preview Source</title>
    <link>https://example.com/</link>
    <description>Test feed</description>
    <item>
      <title>Browser Preview Item</title>
      <link>https://example.com/item-1</link>
      <description>Hello from the test feed</description>
    </item>
  </channel>
</rss>`))
	}))
	t.Cleanup(server.Close)
	return server
}
