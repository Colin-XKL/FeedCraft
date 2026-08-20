package parser

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"bytes"
	"fmt"
	"strconv"
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
	feed.Title = strings.TrimSpace(doc.Find("title").Text())
	feed.ImageURL = extractFeedIconURL(doc)
	if feed.ImageURL != "" {
		feed.ImageTitle = feed.Title
	}

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
			htmlStr, err := sel.Html()
			if err == nil && htmlStr != "" {
				item.Content = htmlStr
			}
		}

		feed.Articles = append(feed.Articles, item)

	})

	return feed, nil
}

func extractFeedIconURL(doc *goquery.Document) string {
	bestHref := ""
	bestScore := 0
	found := false
	doc.Find("link[rel]").Each(func(_ int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		kind := feedIconRelKind(rel)
		if kind == feedIconKindNone {
			return
		}
		href := strings.TrimSpace(s.AttrOr("href", ""))
		if href == "" {
			return
		}
		score := scoreFeedIcon(kind, s.AttrOr("sizes", ""), s.AttrOr("type", ""))
		if !found || score > bestScore {
			found = true
			bestScore = score
			bestHref = href
		}
	})
	return bestHref
}

type feedIconKind int

const (
	feedIconKindNone feedIconKind = iota
	feedIconKindFallback
	feedIconKindStandard
)

func hasFeedIconRel(rel string) bool {
	return feedIconRelKind(rel) != feedIconKindNone
}

func feedIconRelKind(rel string) feedIconKind {
	kind := feedIconKindNone
	for _, token := range strings.Fields(strings.ToLower(rel)) {
		switch token {
		case "icon":
			kind = feedIconKindStandard
		case "apple-touch-icon", "apple-touch-icon-precomposed", "mask-icon":
			if kind != feedIconKindStandard {
				kind = feedIconKindFallback
			}
		}
	}
	return kind
}

func scoreFeedIcon(kind feedIconKind, sizes, typ string) int {
	score := 0
	if kind == feedIconKindStandard {
		score += 1000
	}
	score += feedIconSizeScore(sizes)
	score += feedIconTypeScore(typ)
	return score
}

func feedIconSizeScore(sizes string) int {
	sizes = strings.ToLower(strings.TrimSpace(sizes))
	if sizes == "" {
		return 50
	}
	if sizes == "any" {
		return 95
	}
	best := 0
	for _, part := range strings.Fields(sizes) {
		dim, ok := parseIconSize(part)
		if !ok {
			continue
		}
		score := 20
		switch {
		case dim == 32:
			score = 100
		case dim == 16 || dim == 48:
			score = 80
		case dim >= 16 && dim <= 64:
			score = 70
		case dim <= 128:
			score = 40
		}
		if score > best {
			best = score
		}
	}
	if best == 0 {
		return 50
	}
	return best
}

func parseIconSize(part string) (int, bool) {
	widthStr, heightStr, ok := strings.Cut(strings.ToLower(part), "x")
	if !ok {
		return 0, false
	}
	width, errW := strconv.Atoi(widthStr)
	height, errH := strconv.Atoi(heightStr)
	if errW != nil || errH != nil || width <= 0 || height <= 0 {
		return 0, false
	}
	if height < width {
		return height, true
	}
	return width, true
}

func feedIconTypeScore(typ string) int {
	typ = strings.ToLower(typ)
	switch {
	case strings.Contains(typ, "svg"):
		return 25
	case strings.Contains(typ, "png"), strings.Contains(typ, "icon"):
		return 18
	case typ == "":
		return 10
	default:
		return 5
	}
}
