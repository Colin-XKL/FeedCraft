package parser

import (
	"FeedCraft/internal/model"
	"bytes"
	"github.com/mmcdole/gofeed"
)

// RssParser parses XML/RSS data into a CraftFeed.
type RssParser struct{}

func (p *RssParser) Parse(data []byte) (*model.CraftFeed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return model.FromGofeed(feed), nil
}
