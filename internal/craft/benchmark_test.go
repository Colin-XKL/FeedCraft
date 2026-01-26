package craft

import (
	"FeedCraft/internal/util"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/mmcdole/gofeed"
)

func BenchmarkProcessFeed_TranslateTitle_Cached(b *testing.B) {
	// Setup Miniredis
	s, err := miniredis.Run()
	if err != nil {
		b.Fatalf("Could not start miniredis: %s", err)
	}
	defer s.Close()

	// Set Environment Variables
	b.Setenv("FC_REDIS_URI", fmt.Sprintf("redis://%s", s.Addr()))
	b.Setenv("FC_DEFAULT_TARGET_LANG", "zh-CN")

	// Setup DB Path
	tmpDir := b.TempDir()
	b.Setenv("FC_DB_SQLITE_PATH", tmpDir)

	// Prepare dummy feed
	feedURL := "http://example.com/rss.xml"
	craftName := "translate-title"

	// Pre-populate Cache
	itemTitle := "Hello World"
	// Calculate MD5
	md5Hash := util.GetMD5Hash(itemTitle)

	// Key format: web_content_translate title_<MD5>
	// We use the internal helper to avoid hardcoding the format
	cacheKey := getCraftCacheKey("translate title", md5Hash)

	// Set in Miniredis
	s.Set(cacheKey, "你好世界")
    // Set TTL just in case (CachedFunc sets it)
    s.SetTTL(cacheKey, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer() // Setup per iteration

		// Reset feed item because craft might modify it in place
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
