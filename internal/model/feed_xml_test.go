package model

import (
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCraftFeed_ToRss_StripsIllegalXMLControlCharacters(t *testing.T) {
	feed := &CraftFeed{
		Title:       "Inbox",
		Link:        "https://example.com/inbox",
		Description: "test",
		Articles: []*CraftArticle{
			{
				Title:       "Post",
				Link:        "https://example.com/1",
				Description: "hello\x05world",
				Content:     "content\x05here",
				Id:          "item-1",
			},
		},
	}

	rss, err := feed.ToFeedsFeed().ToRss()
	require.NoError(t, err)
	assert.False(t, strings.ContainsRune(rss, 0x05))

	parsed, err := gofeed.NewParser().ParseString(rss)
	require.NoError(t, err)
	require.Len(t, parsed.Items, 1)
	assert.Equal(t, "helloworld", parsed.Items[0].Description)
	assert.Equal(t, "contenthere", parsed.Items[0].Content)
}
