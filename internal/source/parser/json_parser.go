package parser

import (
	"FeedCraft/internal/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/itchyny/gojq"
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

	// Navigate to items array
	var itemsArray []interface{}

	// If ItemsIterator is empty or ".", use root
	if p.Config.ItemsIterator == "" || p.Config.ItemsIterator == "." {
		if arr, ok := rawData.([]interface{}); ok {
			itemsArray = arr
		} else {
			itemsArray = []interface{}{rawData}
		}
	} else {
		query, err := gojq.Parse(p.Config.ItemsIterator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse items_iterator '%s': %w", p.Config.ItemsIterator, err)
		}
		iter := query.Run(rawData)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return nil, fmt.Errorf("jq execution failed: %w", err)
			}
			if arr, ok := v.([]interface{}); ok {
				itemsArray = append(itemsArray, arr...)
			} else {
				itemsArray = append(itemsArray, v)
			}
		}
	}

	// Compile field selectors
	compile := func(s string) (*gojq.Query, error) {
		if s == "" {
			return nil, nil
		}
		return gojq.Parse(s)
	}

	titleQ, err := compile(p.Config.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid title selector: %w", err)
	}

	linkQ, err := compile(p.Config.Link)
	if err != nil {
		return nil, fmt.Errorf("invalid link selector: %w", err)
	}

	dateQ, err := compile(p.Config.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date selector: %w", err)
	}

	descQ, err := compile(p.Config.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid description selector: %w", err)
	}

	runQuery := func(q *gojq.Query, data interface{}) interface{} {
		if q == nil {
			return nil
		}
		iter := q.Run(data)
		v, ok := iter.Next()
		if !ok {
			return nil
		}
		if _, ok := v.(error); ok {
			return nil
		}
		return v
	}

	for _, itemNode := range itemsArray {
		item := &gofeed.Item{}

		if val := runQuery(titleQ, itemNode); val != nil {
			item.Title = fmt.Sprintf("%v", val)
		}

		if val := runQuery(linkQ, itemNode); val != nil {
			item.Link = fmt.Sprintf("%v", val)
		}

		if val := runQuery(dateQ, itemNode); val != nil {
			dateStr := fmt.Sprintf("%v", val)
			item.Published = dateStr
			// Attempt to parse into PublishedParsed
			if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
				item.PublishedParsed = &t
			} else if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				item.PublishedParsed = &t
			}
		}

		if val := runQuery(descQ, itemNode); val != nil {
			item.Description = fmt.Sprintf("%v", val)
			item.Content = item.Description
		}

		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}
