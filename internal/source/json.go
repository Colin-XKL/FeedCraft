package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"fmt"
)

func init() {
	Register(constant.SourceJSON, jsonSourceFactory)
}

func jsonSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.HttpFetcher == nil {
		return nil, fmt.Errorf("http_fetcher config is required for json source")
	}
	if cfg.JsonParser == nil {
		return nil, fmt.Errorf("json_parser config is required for json source")
	}

	return &PipelineSource{
		Config:  cfg,
		Fetcher: &fetcher.HttpFetcher{Config: cfg.HttpFetcher},
		Parser:  &parser.JsonParser{Config: cfg.JsonParser},
	}, nil
}
