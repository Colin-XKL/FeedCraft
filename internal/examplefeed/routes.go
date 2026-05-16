package examplefeed

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
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
	slug, ok := strings.CutSuffix(c.Param("slug"), ".xml")
	if !ok {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Example RSS feed not found"})
		return
	}

	feed, err := Build(slug, time.Now(), requestBaseURL(c))
	if err != nil {
		if errors.Is(err, ErrUnknownFeed) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Example RSS feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	rss, err := feed.ToFeedsFeed().ToRss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to render RSS: " + err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rss))
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
