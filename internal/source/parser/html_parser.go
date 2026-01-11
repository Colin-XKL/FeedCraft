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

			var link string
			// Try to get href from the element itself
			if href, exists := sel.Attr("href"); exists {
				link = href
			} else {
				// If not found, try to find a child 'a' tag
				childA := sel.Find("a").First()
				if childA.Length() > 0 {
					if href, exists := childA.Attr("href"); exists {
						link = href
					}
				}

				// If still not found, try to find a parent 'a' tag (closest ancestor)
				if link == "" {
					parentA := sel.Closest("a")
					if parentA.Length() > 0 {
						if href, exists := parentA.Attr("href"); exists {
							link = href
						}
					}
				}
			}

			// Fallback to text if still empty, BUT only if it looks like a URL (no spaces)
			if link == "" {
				text := strings.TrimSpace(sel.Text())
				if text != "" && !strings.ContainsAny(text, " \t\n") {
					link = text
				}
			}

			item.Link = link
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
			item.Published = dateStr
		}

		// Description (plain text)
		if p.Config.Description != "" {
			item.Description = strings.TrimSpace(getSelection(p.Config.Description).Text())
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

		feed.Items = append(feed.Items, item)

	})

	return feed, nil
}
