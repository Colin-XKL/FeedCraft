package examplefeed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

var svgAssets = map[string]string{
	"picture-wide.svg":     svgFixture("900", "320", "#0f766e", "wide picture source"),
	"picture-medium.svg":   svgFixture("640", "320", "#2563eb", "medium picture source"),
	"picture-fallback.svg": svgFixture("480", "320", "#7c3aed", "fallback img source"),
}

func RegisterRoutes(router *gin.Engine) {
	router.GET(CatalogAPIPath, CatalogHandler)
	router.GET(RoutePrefix+"/:slug", RSSHandler)
	router.GET(AssetRoutePrefix+"/:name", AssetHandler)
}

func CatalogHandler(c *gin.Context) {
	c.JSON(http.StatusOK, util.APIResponse[[]FeedMeta]{Data: Catalog()})
}

func RSSHandler(c *gin.Context) {
	def, ok := findDefinitionByPathName(c.Param("slug"))
	if !ok {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Example RSS feed not found"})
		return
	}

	feed, err := Build(def.Slug, time.Now(), requestBaseURL(c))
	if err != nil {
		if errors.Is(err, ErrUnknownFeed) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Example RSS feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	body, contentType, err := renderFeed(def, feed.ToFeedsFeed())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to render example feed: " + err.Error()})
		return
	}
	c.Data(http.StatusOK, contentType, []byte(body))
}

func AssetHandler(c *gin.Context) {
	name := c.Param("name")
	body, ok := svgAssets[name]
	if !ok || strings.Contains(name, "/") {
		c.String(http.StatusNotFound, "asset not found")
		return
	}
	c.Data(http.StatusOK, "image/svg+xml; charset=utf-8", []byte(body))
}

func svgFixture(width string, height string, color string, label string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 %s %s" role="img" aria-label="%s"><rect width="100%%" height="100%%" fill="%s"/><circle cx="80" cy="80" r="48" fill="rgba(255,255,255,.35)"/><text x="50%%" y="50%%" dominant-baseline="middle" text-anchor="middle" font-family="Arial, sans-serif" font-size="28" fill="white">FeedCraft %s</text></svg>`, width, height, width, height, label, color, label)
}

func requestBaseURL(c *gin.Context) string {
	if configured := strings.TrimSpace(util.GetEnvClient().GetString("SITE_BASE_URL")); configured != "" {
		return normalizeBaseURL(configured)
	}

	scheme := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = c.Request.Host
	}
	return normalizeBaseURL(scheme + "://" + host)
}

func renderFeed(def feedDefinition, feed *feeds.Feed) (string, string, error) {
	switch def.outputFormat {
	case outputRSS1:
		return renderRSS1(feed), "application/rdf+xml; charset=utf-8", nil
	case outputRSS092:
		return renderRSS092(feed), "application/rss+xml; charset=utf-8", nil
	case outputAtom:
		body, err := feed.ToAtom()
		return body, "application/atom+xml; charset=utf-8", err
	case outputJSONFeed:
		body, err := feed.ToJSON()
		return body, "application/feed+json; charset=utf-8", err
	default:
		body, err := feed.ToRss()
		return body, "application/rss+xml; charset=utf-8", err
	}
}

func renderRSS092(feed *feeds.Feed) string {
	link := ""
	if feed.Link != nil {
		link = feed.Link.Href
	}
	itemTitle := "Format support sample"
	itemLink := link + "#format-support"
	itemDescription := "Format support sample"
	if len(feed.Items) > 0 && feed.Items[0] != nil {
		item := feed.Items[0]
		itemTitle = item.Title
		if item.Id != "" {
			itemLink = item.Id
		} else if item.Link != nil && item.Link.Href != "" {
			itemLink = item.Link.Href
		}
		if item.Description != "" {
			itemDescription = item.Description
		}
		if item.Content != "" {
			itemDescription = item.Content
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="0.92">
  <channel>
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <item>
      <title>%s</title>
      <link>%s</link>
      <description><![CDATA[%s]]></description>
    </item>
  </channel>
</rss>`, escapeXML(feed.Title), escapeXML(link), escapeXML(feed.Description), escapeXML(itemTitle), escapeXML(itemLink), escapeXML(itemDescription))
}

func renderRSS1(feed *feeds.Feed) string {
	link := ""
	if feed.Link != nil {
		link = feed.Link.Href
	}
	itemTitle := "Format support sample"
	itemLink := link + "#format-support"
	itemResource := itemLink
	itemDescription := "Format support sample"
	itemContent := formatFixture
	if len(feed.Items) > 0 && feed.Items[0] != nil {
		item := feed.Items[0]
		itemTitle = item.Title
		if item.Id != "" {
			itemResource = item.Id
		}
		if item.Link != nil && item.Link.Href != "" {
			itemLink = item.Link.Href
		}
		if item.Description != "" {
			itemDescription = item.Description
		}
		if item.Content != "" {
			itemContent = item.Content
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel rdf:about="%s">
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <items>
      <rdf:Seq>
        <rdf:li rdf:resource="%s"/>
      </rdf:Seq>
    </items>
  </channel>
  <item rdf:about="%s">
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <content:encoded><![CDATA[%s]]></content:encoded>
  </item>
</rdf:RDF>`, escapeXML(link), escapeXML(feed.Title), escapeXML(link), escapeXML(feed.Description), escapeXML(itemResource), escapeXML(itemResource), escapeXML(itemTitle), escapeXML(itemLink), escapeXML(itemDescription), itemContent)
}

func escapeXML(value string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(value)); err != nil {
		return value
	}
	return buf.String()
}
