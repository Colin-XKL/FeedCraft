package source

import (
	"FeedCraft/internal/config"
	"context"
	"github.com/mmcdole/gofeed"
)

// Source is the top-level interface for all generators.
// Its sole responsibility is to produce a standard feed.
type Source interface {
	Generate(ctx context.Context) (*gofeed.Feed, error)
	BaseURL() string
}

// SourceFactory is responsible for creating a Source instance from a configuration.
// It receives a parsed SourceConfig struct, promoting type safety and cleaner factory implementations.
type SourceFactory func(config *config.SourceConfig) (Source, error)
