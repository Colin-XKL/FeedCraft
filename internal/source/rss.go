package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"fmt"
)

func init() {
	Register(constant.SourceRSS, rssSourceFactory)
}

func rssSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.HttpFetcher == nil {
		return nil, fmt.Errorf("http_fetcher config is required for rss source")
	}

	return &PipelineSource{
		Config:  cfg, // Inject the full config
		Fetcher: &fetcher.HttpFetcher{Config: cfg.HttpFetcher},
		Parser:  &parser.RssParser{},
	}, nil
}
