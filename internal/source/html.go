package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"fmt"
)

func init() {
	Register(constant.SourceHTML, htmlSourceFactory)
}

func htmlSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.HttpFetcher == nil {
		return nil, fmt.Errorf("http_fetcher config is required for html source")
	}
	if cfg.HtmlParser == nil {
		return nil, fmt.Errorf("html_parser config is required for html source")
	}

	return &PipelineSource{
		Config:  cfg,
		Fetcher: &fetcher.HttpFetcher{Config: cfg.HttpFetcher},
		Parser:  &parser.HtmlParser{Config: cfg.HtmlParser},
	}, nil
}
