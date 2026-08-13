package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRssParser_Parse_StripsIllegalXMLControlCharacters(t *testing.T) {
	rss := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Inbox</title>
    <link>https://example.com/inbox</link>
    <description>test</description>
    <item>
      <title>Post` + "\x05" + `Title</title>
      <link>https://example.com/1</link>
      <description>hello` + "\x05" + `world</description>
      <guid>item-1</guid>
    </item>
  </channel>
</rss>`

	feed, err := (&RssParser{}).Parse([]byte(rss))
	require.NoError(t, err)
	require.NotNil(t, feed)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "PostTitle", feed.Articles[0].Title)
	assert.Equal(t, "helloworld", feed.Articles[0].Description)
}
