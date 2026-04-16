package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/model"
	"context"
)

// Source is the top-level interface for all source providers.
// Its responsibility is to fetch and parse upstream data into a CraftFeed.
type Source interface {
	Fetch(ctx context.Context) (*model.CraftFeed, error)
	BaseURL() string
}

// SourceFactory is responsible for creating a Source instance from a configuration.
// It receives a parsed SourceConfig struct, promoting type safety and cleaner factory implementations.
type SourceFactory func(config *config.SourceConfig) (Source, error)
