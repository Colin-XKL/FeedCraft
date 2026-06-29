package craft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAbsFeedLink(t *testing.T) {
	tests := []struct {
		name         string
		feedURL      string
		feedLinkAttr string
		want         string
	}{
		{
			name:         "absolute feed link is preserved",
			feedURL:      "https://example.com/rss.xml",
			feedLinkAttr: "https://blog.example.com/",
			want:         "https://blog.example.com/",
		},
		{
			name:         "relative feed link resolved against feed url",
			feedURL:      "https://example.com/path/rss.xml",
			feedLinkAttr: "/home",
			want:         "https://example.com/home",
		},
		{
			name:         "empty feed link falls back to feed url origin",
			feedURL:      "https://example.com/path/rss.xml",
			feedLinkAttr: "",
			want:         "https://example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, getAbsFeedLink(tc.feedURL, tc.feedLinkAttr))
		})
	}
}

func TestGetAbsLinkForFeedItem(t *testing.T) {
	tests := []struct {
		name         string
		feedURL      string
		feedLinkAttr string
		feedItemURL  string
		want         string
	}{
		{
			name:         "item already absolute is kept",
			feedURL:      "https://example.com/rss.xml",
			feedLinkAttr: "https://example.com/",
			feedItemURL:  "https://cdn.example.com/article-1",
			want:         "https://cdn.example.com/article-1",
		},
		{
			name:         "relative item resolved against absolute feed link",
			feedURL:      "https://aggregator.example.net/rss.xml",
			feedLinkAttr: "https://blog.example.com/",
			feedItemURL:  "/posts/hello",
			want:         "https://blog.example.com/posts/hello",
		},
		{
			name:         "relative item resolved against feed url when feed link is relative",
			feedURL:      "https://example.com/section/rss.xml",
			feedLinkAttr: "/relative-link",
			feedItemURL:  "article-2",
			want:         "https://example.com/section/article-2",
		},
		{
			name:         "relative item resolved against feed url when feed link empty",
			feedURL:      "https://example.com/section/rss.xml",
			feedLinkAttr: "",
			feedItemURL:  "article-3",
			want:         "https://example.com/section/article-3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, getAbsLinkForFeedItem(tc.feedURL, tc.feedLinkAttr, tc.feedItemURL))
		})
	}
}
