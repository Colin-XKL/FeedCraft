package examplefeed

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogListsExampleFeeds(t *testing.T) {
	items := Catalog()

	require.Len(t, items, 4)
	assert.Equal(t, "html-elements", items[0].Slug)
	assert.Equal(t, "/example-rss-feeds/html-elements.xml", items[0].Path)
	assert.Equal(t, "html-styling", items[1].Slug)
	assert.Equal(t, "media-picture", items[2].Slug)
	assert.Equal(t, "all-in-one", items[3].Slug)
}

func TestWindowUUIDIsStableWithinFourHours(t *testing.T) {
	now := time.Date(2026, 5, 16, 8, 15, 0, 0, time.UTC)

	first := WindowUUID("html-elements", now)
	second := WindowUUID("html-elements", now.Add(3*time.Hour+29*time.Minute))
	next := WindowUUID("html-elements", now.Add(4*time.Hour))

	assert.Equal(t, first, second)
	assert.NotEqual(t, first, next)
}

func TestBuildFeedIncludesRotatingGUIDsAndHTMLFixtures(t *testing.T) {
	now := time.Date(2026, 5, 16, 8, 15, 0, 0, time.UTC)

	feed, err := Build("all-in-one", now, "https://feedcraft.example")

	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "FeedCraft Example RSS Feeds - All in One", feed.Title)
	assert.Equal(t, "https://feedcraft.example/example-rss-feeds/all-in-one.xml", feed.Id)
	require.Len(t, feed.Articles, 3)
	assert.Equal(t, "https://feedcraft.example/example-rss-feeds/all-in-one.xml#html-elements-"+WindowUUID("all-in-one", now), feed.Articles[0].Id)
	assert.Contains(t, feed.Articles[0].Id, WindowUUID("all-in-one", now))
	assert.Contains(t, feed.Articles[0].Content, "<details>")
	assert.Contains(t, feed.Articles[1].Content, "display: flex")
	assert.Contains(t, feed.Articles[2].Content, "<picture>")
	assert.Contains(t, feed.Articles[2].Content, "https://feedcraft.example/example-rss-feeds/assets/picture-wide.svg")
}

func TestBuildReturnsUnknownSlugError(t *testing.T) {
	feed, err := Build("missing", time.Now(), "https://feedcraft.example")

	assert.Nil(t, feed)
	assert.ErrorIs(t, err, ErrUnknownFeed)
}

func TestRegisterRoutesServesRSSAndCatalog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router)

	t.Run("rss", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/example-rss-feeds/html-elements.xml", nil)
		req.Host = "feedcraft.example"
		req.Header.Set("X-Forwarded-Proto", "https")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/rss+xml; charset=utf-8", recorder.Header().Get("Content-Type"))
		parsed, err := gofeed.NewParser().ParseString(recorder.Body.String())
		require.NoError(t, err)
		assert.Equal(t, "FeedCraft Example RSS Feeds - HTML Elements", parsed.Title)
		require.Len(t, parsed.Items, 1)
		assert.Equal(t, "https://feedcraft.example/example-rss-feeds/html-elements.xml", parsed.Link)
		assert.Contains(t, parsed.Items[0].GUID, "https://feedcraft.example/example-rss-feeds/html-elements.xml#html-elements-")
		assert.Contains(t, parsed.Items[0].Content, "<h1>")
	})

	t.Run("catalog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/example-rss-feeds", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"slug":"html-elements"`)
		assert.Contains(t, recorder.Body.String(), `"path":"/example-rss-feeds/html-elements.xml"`)
	})

	t.Run("asset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/example-rss-feeds/assets/picture-wide.svg", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "image/svg+xml; charset=utf-8", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "<svg")
	})

	t.Run("missing slug", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/example-rss-feeds/missing.xml", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
