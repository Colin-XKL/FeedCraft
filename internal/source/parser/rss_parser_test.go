package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRssParserParseRemovesInvalidXMLControlCharactersFromExternalFeed(t *testing.T) {
	rss := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>External feed</title>
    <link>https://example.com/feed</link>
    <description>Example feed</description>
    <item>
      <title>Example item</title>
      <link>https://example.com/item</link>
      <description>before` + "\x05" + `after</description>
      <content:encoded><![CDATA[content` + "\x05" + `body]]></content:encoded>
      <guid>item-1</guid>
    </item>
  </channel>
</rss>`)

	feed, err := (&RssParser{}).Parse(rss)

	require.NoError(t, err)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "beforeafter", feed.Articles[0].Description)
	assert.Equal(t, "contentbody", feed.Articles[0].Content)
}
