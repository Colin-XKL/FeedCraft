package parser

import (
	"FeedCraft/internal/config"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"strings"
)

type HtmlParser struct {
	Config *config.HtmlParserConfig
}

func (p *HtmlParser) Parse(data []byte) (*gofeed.Feed, error) {
	if p == nil || p.Config == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	feed := &gofeed.Feed{}

	// Basic feed metadata (can be overridden by FeedMetaConfig later)
	// For now, we might try to extract title from <title> if not provided via overrides
	feed.Title = doc.Find("title").Text()

	doc.Find(p.Config.ItemSelector).Each(func(i int, s *goquery.Selection) {
		item := &gofeed.Item{}

		// Title
		if p.Config.Title != "" {
			item.Title = strings.TrimSpace(s.Find(p.Config.Title).Text())
		}

		// Link
		if p.Config.Link != "" {
			linkSel := s.Find(p.Config.Link)
			if href, exists := linkSel.Attr("href"); exists {
				item.Link = href
			} else {
				item.Link = strings.TrimSpace(linkSel.Text())
			}
		}

		// Date
		if p.Config.Date != "" {
			dateSel := s.Find(p.Config.Date)
			dateStr := strings.TrimSpace(dateSel.Text())
			if dateStr == "" {
				if val, exists := dateSel.Attr("datetime"); exists {
					dateStr = val
				}
			}
			item.Published = dateStr
		}

		// Description
		if p.Config.Description != "" {
			item.Description = strings.TrimSpace(s.Find(p.Config.Description).Text())
			// fall back to content if description is empty, or maybe we want Content specifically?
			// gofeed.Item has Description and Content.
			item.Content = item.Description 
		}

		feed.Items = append(feed.Items, item)
	})

	return feed, nil
}
