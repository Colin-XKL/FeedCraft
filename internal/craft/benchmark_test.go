package craft

import (
	"FeedCraft/internal/util"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func BenchmarkProcessFeed_TranslateTitle_Cached(b *testing.B) {
	redis := setupTestRedis(b)
	b.Setenv("FC_DEFAULT_TARGET_LANG", "zh-CN")

	// Setup DB Path
	tmpDir := b.TempDir()
	b.Setenv("FC_DB_SQLITE_PATH", tmpDir)

	// Prepare dummy feed
	feedURL := "http://example.com/rss.xml"
	craftName := "translate-title"

	// Pre-populate Cache
	itemTitle := "Hello World"
	hash := util.GetTextContentHash(itemTitle)
	cacheKey := getCraftCacheKey("translate title", hash)
	redis.SetString(b, cacheKey, "你好世界", time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		feed := &gofeed.Feed{
			Items: []*gofeed.Item{
				{
					Title:       itemTitle,
					Description: "This is a description",
					Link:        "http://example.com/1",
					GUID:        "guid-1",
				},
			},
		}
		b.StartTimer()

		_, err := ProcessFeed(feed, feedURL, craftName)
		if err != nil {
			b.Fatalf("ProcessFeed failed: %v", err)
		}
	}
}
