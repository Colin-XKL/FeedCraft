package parser

import (
	"fmt"
	"strings"
	"testing"
)

func generateDummyRSS(itemCount int) []byte {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
 <title>RSS Title</title>
 <description>This is an example of an RSS feed</description>
 <link>http://www.example.com/main.html</link>
 <lastBuildDate>Mon, 06 Sep 2010 00:01:00 +0000 </lastBuildDate>
 <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
 <ttl>1800</ttl>
`)

	for i := 0; i < itemCount; i++ {
		sb.WriteString(fmt.Sprintf(`
 <item>
  <title>Example entry %d</title>
  <description>Here is some text containing an interesting description %d.</description>
  <link>http://www.example.com/blog/post/%d</link>
  <guid isPermaLink="false">%d</guid>
  <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
 </item>
`, i, i, i, i))
	}

	sb.WriteString(`
</channel>
</rss>
`)
	return []byte(sb.String())
}

func BenchmarkRssParser(b *testing.B) {
	parser := &RssParser{}
	data := generateDummyRSS(100) // 100 items

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}
