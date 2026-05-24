package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"fmt"
)

func init() {
	Register(constant.SourceWebMonitor, webMonitorSourceFactory)
}

func webMonitorSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.HttpFetcher == nil {
		return nil, fmt.Errorf("http_fetcher config is required for web monitor source")
	}
	if cfg.WebMonitorParser == nil {
		return nil, fmt.Errorf("web_monitor_parser config is required for web monitor source")
	}
	if cfg.HttpFetcher.Purpose == "" {
		cfg.HttpFetcher.Purpose = config.HttpFetcherPurposeHTML
	}

	return &PipelineSource{
		Config:  cfg,
		Fetcher: &fetcher.HttpFetcher{Config: cfg.HttpFetcher},
		Parser: &parser.WebMonitorParser{
			Config:  cfg.WebMonitorParser,
			PageURL: cfg.HttpFetcher.URL,
		},
	}, nil
}
