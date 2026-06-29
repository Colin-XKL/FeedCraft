package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"net/url"
)

// PipelineSource is the generic implementation for most scenarios.
// It holds the components (Fetcher, Parser) and the configuration.
type PipelineSource struct {
	Config  *config.SourceConfig
	Fetcher fetcher.Fetcher
	Parser  parser.Parser
}

// Fetch executes the fetch-parse-normalize pipeline.
func (p *PipelineSource) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	// 1. Fetch raw data
	raw, err := p.Fetcher.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	// 2. Parse into a base Feed object
	feed, err := p.Parser.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	p.normalizeCraftFeed(feed)

	return feed, nil
}

// BaseURL delegates the call to the underlying fetcher.
func (p *PipelineSource) BaseURL() string {
	if p.Fetcher == nil {
		return ""
	}
	return p.Fetcher.BaseURL()
}

func (p *PipelineSource) normalizeCraftFeed(feed *model.CraftFeed) {
	if feed == nil {
		return
	}

	baseURL := p.BaseURL()
	if baseURL != "" {
		for _, item := range feed.Articles {
			if item == nil || item.Link == "" {
				continue
			}
			if absURL, err := util.BuildAbsoluteURL(baseURL, item.Link); err == nil {
				item.Link = absURL
			}
		}
	}

	feed.Link = getAbsFeedLink(baseURL, feed.Link)
	p.applyFeedMetaOverrides(feed)
	p.applyFeedIconSource(feed, baseURL)
}

// applyFeedMetaOverrides checks for a FeedMetaConfig and uses its values
// to override the metadata of the parsed feed.
func (p *PipelineSource) applyFeedMetaOverrides(feed *model.CraftFeed) {
	if p.Config == nil || p.Config.FeedMeta == nil {
		return // No overrides provided
	}

	meta := p.Config.FeedMeta
	if meta.Title != "" {
		feed.Title = meta.Title
	}
	if meta.Link != "" {
		feed.Link = meta.Link
	}
	if meta.Description != "" {
		feed.Description = meta.Description
	}
	if meta.Copyright != "" {
		feed.Copyright = meta.Copyright
	}
	if meta.AuthorName != "" || meta.AuthorEmail != "" {
		feed.AuthorName = meta.AuthorName
		feed.AuthorEmail = meta.AuthorEmail
	}
}

func (p *PipelineSource) applyFeedIconSource(feed *model.CraftFeed, baseURL string) {
	if feed == nil {
		return
	}

	iconSource := config.FeedIconSourceAuto
	hasExplicitIconSource := p.Config != nil && p.Config.FeedMeta != nil && p.Config.FeedMeta.IconSource != ""
	if hasExplicitIconSource {
		iconSource = p.Config.FeedMeta.IconSource
	}

	if iconSource == config.FeedIconSourceFaviconService {
		if serviceURL := buildFaviconServiceURL(feed.Link, baseURL); serviceURL != "" {
			feed.ImageURL = serviceURL
			feed.ImageTitle = feed.Title
		}
		return
	}

	if feed.ImageURL != "" {
		if absURL, err := util.BuildAbsoluteURL(baseURL, feed.ImageURL); err == nil {
			feed.ImageURL = absURL
		}
		if feed.ImageTitle == "" {
			feed.ImageTitle = feed.Title
		}
		return
	}

	if hasExplicitIconSource {
		if origin := firstOrigin(feed.Link, baseURL); origin != "" {
			feed.ImageURL = origin + "/favicon.ico"
			feed.ImageTitle = feed.Title
		}
	}
}

func buildFaviconServiceURL(feedLink string, baseURL string) string {
	origin := firstOrigin(feedLink, baseURL)
	if origin == "" {
		return ""
	}
	return "https://www.google.com/s2/favicons?domain_url=" + url.QueryEscape(origin) + "&sz=64"
}

func firstOrigin(rawURLs ...string) string {
	for _, rawURL := range rawURLs {
		parsedURL, err := url.Parse(rawURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			continue
		}
		return parsedURL.Scheme + "://" + parsedURL.Host
	}
	return ""
}
