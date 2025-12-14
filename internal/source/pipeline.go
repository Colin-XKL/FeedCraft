package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
)

// PipelineSource is the generic implementation for most scenarios.
// It holds the components (Fetcher, Parser) and the configuration.
type PipelineSource struct {
	Config  *config.SourceConfig
	Fetcher fetcher.Fetcher
	Parser  parser.Parser
}

// Generate executes the fetch-parse-override pipeline.
func (p *PipelineSource) Generate(ctx context.Context) (*gofeed.Feed, error) {
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

	// 3. Apply metadata overrides from config
	p.applyFeedMetaOverrides(feed)

	return feed, nil
}

// BaseURL delegates the call to the underlying fetcher.
func (p *PipelineSource) BaseURL() string {
	if p.Fetcher == nil {
		return ""
	}
	return p.Fetcher.BaseURL()
}

// applyFeedMetaOverrides checks for a FeedMetaConfig and uses its values
// to override the metadata of the parsed feed.
func (p *PipelineSource) applyFeedMetaOverrides(feed *gofeed.Feed) {
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
		feed.Author = &gofeed.Person{
			Name:  meta.AuthorName,
			Email: meta.AuthorEmail,
		}
	}
}
