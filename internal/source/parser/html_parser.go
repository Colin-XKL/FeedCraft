package parser

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HtmlParser struct {
	Config *config.HtmlParserConfig
}

func (p *HtmlParser) Parse(data []byte) (*model.CraftFeed, error) {
	if p == nil || p.Config == nil {
		return nil, fmt.Errorf("parser config is nil")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	feed := &model.CraftFeed{}

	// Basic feed metadata (can be overridden by FeedMetaConfig later)
	// For now, we might try to extract title from <title> if not provided via overrides
	feed.Title = doc.Find("title").Text()

	doc.Find(p.Config.ItemSelector).Each(func(i int, s *goquery.Selection) {
		item := &model.CraftArticle{}

		// Helper to extract selection based on selector
		// If selector is "." or empty, use current 's'
		// Otherwise find descendant
		getSelection := func(selector string) *goquery.Selection {
			trimmedSelector := strings.TrimSpace(selector)
			if trimmedSelector == "" || trimmedSelector == "." {
				return s
			}
			return s.Find(trimmedSelector)
		}

		// Title
		if p.Config.Title != "" {
			item.Title = strings.TrimSpace(getSelection(p.Config.Title).Text())
		}

		// Link
		if p.Config.Link != "" {
			sel := getSelection(p.Config.Link)
			item.Link = util.ExtractLinkFromSelection(sel)
		}

		// Date
		if p.Config.Date != "" {
			sel := getSelection(p.Config.Date)
			dateStr := strings.TrimSpace(sel.Text())
			if dateStr == "" {
				if val, exists := sel.Attr("datetime"); exists {
					dateStr = val
				}
			}
			if dateStr != "" {
				if parsedTime, ok := parseFlexibleTime(dateStr); ok {
					item.Created = parsedTime
					item.Updated = parsedTime
				}
			}
		}

		// Description (plain text)
		if p.Config.Description != "" {
			item.Description = strings.TrimSpace(getSelection(p.Config.Description).Text())
			if item.Content == "" {
				item.Content = item.Description
			}
		}

		// Content (rich HTML)
		if p.Config.Content != "" {
			sel := getSelection(p.Config.Content)
			html, err := sel.Html()
			if err != nil {
				// Log error but don't fail, just leave content empty
				// logrus.Warnf("Failed to extract content for item: %v", err)
				item.Content = ""
			} else {
				item.Content = html
			}
		}

		feed.Articles = append(feed.Articles, item)

	})

	return feed, nil
}
