package craft

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

const (
	feedFormatHTML = "html"
	feedFormatJSON = "json"
	feedFormatAtom = "atom"
	feedFormatRSS  = "rss"
)

// ServeCraftedFeed writes a processed feed using the representation the client asked for.
// RSS readers keep getting RSS XML. Browsers get a readable HTML preview. `format`
// can force rss, html, json, or atom.
func ServeCraftedFeed(c *gin.Context, feed *feeds.Feed) {
	if c == nil {
		return
	}
	if feed == nil {
		c.String(http.StatusInternalServerError, "empty feed")
		return
	}

	switch negotiateFeedFormat(c.Query("format"), c.GetHeader("Accept")) {
	case feedFormatHTML:
		body, err := renderFeedHTML(feed, requestPageURL(c))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(body))
	case feedFormatJSON:
		body, err := feed.ToJSON()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/feed+json; charset=utf-8", []byte(body))
	case feedFormatAtom:
		body, err := feed.ToAtom()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/atom+xml; charset=utf-8", []byte(body))
	default:
		body, err := feed.ToRss()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(body))
	}
}

func negotiateFeedFormat(formatQuery, acceptHeader string) string {
	switch strings.ToLower(strings.TrimSpace(formatQuery)) {
	case feedFormatHTML:
		return feedFormatHTML
	case feedFormatJSON, "jsonfeed", "json-feed":
		return feedFormatJSON
	case feedFormatAtom:
		return feedFormatAtom
	case feedFormatRSS, "xml":
		return feedFormatRSS
	}
	return formatFromAccept(acceptHeader)
}

func formatFromAccept(accept string) string {
	if strings.TrimSpace(accept) == "" {
		return feedFormatRSS
	}

	var (
		htmlQ, rssQ, atomQ, jsonQ         float64
		hasHTML, hasRSS, hasAtom, hasJSON bool
	)

	for _, part := range strings.Split(accept, ",") {
		media, q := parseAcceptPart(part)
		switch media {
		case "text/html", "application/xhtml+xml":
			htmlQ, hasHTML = maxQuality(htmlQ, q, hasHTML)
		case "application/rss+xml", "application/rdf+xml":
			rssQ, hasRSS = maxQuality(rssQ, q, hasRSS)
		case "application/atom+xml":
			atomQ, hasAtom = maxQuality(atomQ, q, hasAtom)
		case "application/feed+json", "application/json":
			jsonQ, hasJSON = maxQuality(jsonQ, q, hasJSON)
		}
	}

	bestFeed := feedFormatRSS
	bestFeedQ := -1.0
	hasFeed := false
	if hasRSS {
		bestFeed, bestFeedQ, hasFeed = feedFormatRSS, rssQ, true
	}
	if hasAtom && atomQ > bestFeedQ {
		bestFeed, bestFeedQ, hasFeed = feedFormatAtom, atomQ, true
	}
	if hasJSON && jsonQ > bestFeedQ {
		bestFeed, bestFeedQ, hasFeed = feedFormatJSON, jsonQ, true
	}

	if hasFeed && (!hasHTML || bestFeedQ >= htmlQ) {
		return bestFeed
	}
	if hasHTML {
		return feedFormatHTML
	}
	return feedFormatRSS
}

func parseAcceptPart(part string) (string, float64) {
	part = strings.TrimSpace(part)
	media := part
	q := 1.0
	if i := strings.Index(part, ";"); i >= 0 {
		media = strings.TrimSpace(part[:i])
		for _, param := range strings.Split(part[i+1:], ";") {
			param = strings.TrimSpace(param)
			if len(param) < 2 || !strings.EqualFold(param[:2], "q=") {
				continue
			}
			if v, err := strconv.ParseFloat(strings.TrimSpace(param[2:]), 64); err == nil {
				q = v
			}
		}
	}
	return strings.ToLower(media), q
}

func maxQuality(current, candidate float64, seen bool) (float64, bool) {
	if !seen || candidate > current {
		return candidate, true
	}
	return current, true
}

type htmlFeedView struct {
	Title       string
	Description string
	PageURL     string
	RSSURL      string
	JSONURL     string
	Items       []htmlFeedItem
}

type htmlFeedItem struct {
	Title   string
	Link    string
	Date    string
	Snippet string
}

var feedPreviewTemplate = template.Must(template.New("craft-feed-preview").Parse(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} · FeedCraft</title>
  <style>
    :root { color-scheme: light; }
    body { margin: 0; font-family: ui-sans-serif, system-ui, sans-serif; background: #f8fafc; color: #0f172a; }
    main { max-width: 52rem; margin: 0 auto; padding: 2rem 1.25rem 3rem; }
    .banner { border: 1px solid #99f6e4; background: #f0fdfa; border-radius: 16px; padding: 1rem 1.1rem; margin-bottom: 1.5rem; }
    .brand { color: #0d9488; font-weight: 700; letter-spacing: .02em; }
    h1 { font-size: 1.75rem; font-weight: 650; margin: 0 0 .5rem; }
    .desc { color: #475569; margin: 0 0 1rem; line-height: 1.6; }
    .url { display: block; word-break: break-all; background: #fff; border-radius: 10px; padding: .7rem .85rem; color: #0f766e; text-decoration: none; margin: .75rem 0 1rem; }
    .actions { display: flex; flex-wrap: wrap; gap: .6rem; }
    .actions a, .actions button { appearance: none; border: 0; border-radius: 999px; padding: .45rem .9rem; font: inherit; cursor: pointer; text-decoration: none; }
    .primary { background: #0d9488; color: #fff; }
    .secondary { background: #e2e8f0; color: #0f172a; }
    article { background: #fff; border: 1px solid #e2e8f0; border-radius: 16px; padding: 1rem 1.1rem; margin: .85rem 0; }
    article h2 { font-size: 1.05rem; margin: 0 0 .35rem; }
    article h2 a { color: inherit; text-decoration: none; }
    article h2 a:hover { color: #0d9488; }
    time, .snippet { color: #64748b; font-size: .92rem; }
    .snippet { margin: .4rem 0 0; line-height: 1.55; }
    .empty { color: #64748b; }
  </style>
</head>
<body>
  <main>
    <section class="banner">
      <div class="brand">FeedCraft</div>
      <p>这是一个 RSS 订阅地址。阅读器请直接订阅当前 URL。浏览器打开时会显示可读预览，不会改变阅读器拿到的 XML。</p>
      <a class="url" href="{{.RSSURL}}">{{.PageURL}}</a>
      <div class="actions">
        <button class="primary" type="button" data-copy-url="{{.PageURL}}">复制订阅地址</button>
        <a class="secondary" href="{{.RSSURL}}">查看原始 RSS</a>
        <a class="secondary" href="{{.JSONURL}}">JSON Feed</a>
      </div>
    </section>
    <h1>{{.Title}}</h1>
    {{if .Description}}<p class="desc">{{.Description}}</p>{{end}}
    {{if .Items}}
      {{range .Items}}
        <article>
          <h2>{{if .Link}}<a href="{{.Link}}" target="_blank" rel="noopener noreferrer">{{.Title}}</a>{{else}}{{.Title}}{{end}}</h2>
          {{if .Date}}<time>{{.Date}}</time>{{end}}
          {{if .Snippet}}<p class="snippet">{{.Snippet}}</p>{{end}}
        </article>
      {{end}}
    {{else}}
      <p class="empty">这个订阅源暂时没有条目。</p>
    {{end}}
  </main>
  <script>
    document.querySelectorAll("[data-copy-url]").forEach(function (button) {
      button.addEventListener("click", function () {
        var url = button.getAttribute("data-copy-url") || "";
        if (!url || !navigator.clipboard) { return; }
        navigator.clipboard.writeText(url).then(function () {
          button.textContent = "已复制";
        }).catch(function () {});
      });
    });
  </script>
</body>
</html>`))

func renderFeedHTML(feed *feeds.Feed, pageURL string) (string, error) {
	if feed == nil {
		feed = &feeds.Feed{}
	}

	view := htmlFeedView{
		Title:       strings.TrimSpace(feed.Title),
		Description: plainTextSnippet(feed.Description, 400),
		PageURL:     pageURL,
		RSSURL:      withFeedFormat(pageURL, feedFormatRSS),
		JSONURL:     withFeedFormat(pageURL, feedFormatJSON),
		Items:       make([]htmlFeedItem, 0, len(feed.Items)),
	}
	if view.Title == "" {
		view.Title = "FeedCraft"
	}

	for _, item := range feed.Items {
		if item == nil {
			continue
		}
		content := item.Description
		if strings.TrimSpace(content) == "" {
			content = item.Content
		}
		view.Items = append(view.Items, htmlFeedItem{
			Title:   firstNonEmpty(strings.TrimSpace(item.Title), "Untitled"),
			Link:    safeHTTPURL(itemLink(item)),
			Date:    formatFeedItemTime(item.Created, item.Updated),
			Snippet: plainTextSnippet(content, 280),
		})
	}

	var buf bytes.Buffer
	if err := feedPreviewTemplate.Execute(&buf, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func requestPageURL(c *gin.Context) string {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return ""
	}
	if c.Request.URL.IsAbs() {
		return c.Request.URL.String()
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); proto != "" {
		scheme = proto
	}
	host := strings.TrimSpace(c.Request.Host)
	if host == "" {
		return c.Request.URL.RequestURI()
	}
	return scheme + "://" + host + c.Request.URL.RequestURI()
}

func withFeedFormat(rawURL, format string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed == nil {
		return rawURL
	}
	query := parsed.Query()
	query.Set("format", format)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func itemLink(item *feeds.Item) string {
	if item == nil || item.Link == nil {
		return ""
	}
	return item.Link.Href
}

func safeHTTPURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}
	return parsed.String()
}

func formatFeedItemTime(primary, fallback time.Time) string {
	when := primary
	if when.IsZero() {
		when = fallback
	}
	if when.IsZero() {
		return ""
	}
	return when.Local().Format("2006-01-02 15:04")
}

func plainTextSnippet(raw string, maxRunes int) string {
	text := strings.TrimSpace(raw)
	if text == "" {
		return ""
	}
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(text)); err == nil {
		text = strings.Join(strings.Fields(doc.Text()), " ")
	}
	return truncateRunes(text, maxRunes)
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
