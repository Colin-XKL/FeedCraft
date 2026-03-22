package parser

import (
	"FeedCraft/internal/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type JsonParser struct {
	Config *config.JsonParserConfig
}

func (p *JsonParser) Parse(data []byte) (*gofeed.Feed, error) {
	if p == nil || p.Config == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	var rawData interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("invalid json data: %w", err)
	}

	feed := &gofeed.Feed{}
	items, err := ParseJSONItems(rawData, p.Config)
	if err != nil {
		return nil, err
	}

	for _, parsedFields := range items {
		item := &gofeed.Item{}

		item.Title = parsedFields.Title
		item.Link = parsedFields.Link

		if parsedFields.Date != "" {
			item.Published = parsedFields.Date
			// Attempt to parse into PublishedParsed
			if t, err := time.Parse(time.RFC3339, parsedFields.Date); err == nil {
				item.PublishedParsed = &t
			} else if t, err := time.Parse("2006-01-02", parsedFields.Date); err == nil {
				item.PublishedParsed = &t
			}
		}

		item.Description = parsedFields.Description
		item.Content = parsedFields.Description
		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}
