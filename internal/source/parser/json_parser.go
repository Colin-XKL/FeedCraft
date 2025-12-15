package parser

import (
	"FeedCraft/internal/config"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
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
	// If ItemsIterator is empty or ".", use root
	var itemsNode interface{}
	if p.Config.ItemsIterator == "" || p.Config.ItemsIterator == "." {
		itemsNode = rawData
	} else {
		itemsNode = traverse(rawData, p.Config.ItemsIterator)
	}

	itemsArray, ok := itemsNode.([]interface{})
	if !ok {
		return nil, fmt.Errorf("items_iterator '%s' did not resolve to an array", p.Config.ItemsIterator)
	}

	for _, itemNode := range itemsArray {
		item := &gofeed.Item{}

		if p.Config.Title != "" {
			if val := traverse(itemNode, p.Config.Title); val != nil {
				item.Title = fmt.Sprintf("%v", val)
			}
		}
		
		if p.Config.Link != "" {
			if val := traverse(itemNode, p.Config.Link); val != nil {
				item.Link = fmt.Sprintf("%v", val)
			}
		}

		if p.Config.Date != "" {
			if val := traverse(itemNode, p.Config.Date); val != nil {
				item.Published = fmt.Sprintf("%v", val)
			}
		}

		if p.Config.Description != "" {
			if val := traverse(itemNode, p.Config.Description); val != nil {
				item.Description = fmt.Sprintf("%v", val)
				item.Content = item.Description
			}
		}

		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

// traverse walks the JSON object using dot notation "field.subfield"
// Note: This is a basic implementation. It doesn't support array indexing in path (e.g. "items.0.title").
func traverse(data interface{}, path string) interface{} {
	if path == "" {
		return data
	}
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		val, exists := m[part]
		if !exists {
			return nil
		}
		current = val
	}
	return current
}
