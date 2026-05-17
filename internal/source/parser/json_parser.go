package parser

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonParser struct {
	Config *config.JsonParserConfig
}

func (p *JsonParser) Parse(data []byte) (*model.CraftFeed, error) {
	if p == nil || p.Config == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	var rawData interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&rawData); err != nil {
		return nil, fmt.Errorf("invalid json data: %w", err)
	}

	feed := &model.CraftFeed{}
	items, err := ParseJSONItems(rawData, p.Config)
	if err != nil {
		return nil, err
	}

	for _, parsedFields := range items {
		item := &model.CraftArticle{
			Title:       parsedFields.Title,
			Link:        parsedFields.Link,
			Description: parsedFields.Description,
			Content:     parsedFields.Description,
		}
		if parsedFields.Date != "" {
			if parsedTime, ok := parseFlexibleTime(parsedFields.Date); ok {
				item.Created = parsedTime
				item.Updated = parsedTime
			}
		}

		feed.Articles = append(feed.Articles, item)
	}

	return feed, nil
}
