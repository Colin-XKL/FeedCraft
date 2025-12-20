package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"fmt"
)

func init() {
	Register(constant.SourceSearch, searchSourceFactory)
}

func searchSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.SearchFetcher == nil {
		return nil, fmt.Errorf("search_fetcher config is required for search source")
	}

	// Default to JsonParser if not specified, but usually it should be specified in the recipe.
	var p parser.Parser
	if cfg.JsonParser != nil {
		p = &parser.JsonParser{Config: cfg.JsonParser}
	} else {
		// Fallback or Error?
		// We'll require JsonParser config.
		return nil, fmt.Errorf("json_parser config is required for search source")
	}

	return &PipelineSource{
		Config:  cfg,
		Fetcher: &fetcher.SearchFetcher{Config: cfg.SearchFetcher},
		Parser:  p,
	}, nil
}
