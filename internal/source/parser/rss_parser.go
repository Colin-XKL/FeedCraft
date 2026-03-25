package parser

import (
	"bytes"
	"github.com/mmcdole/gofeed"
)

// RssParser parses XML/RSS data into a gofeed.Feed.
type RssParser struct{}

func (p *RssParser) Parse(data []byte) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return feed, nil
}
