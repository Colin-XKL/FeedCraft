package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/fetcher/provider"
	"FeedCraft/internal/source/parser"
	"FeedCraft/internal/util"
	"fmt"
)

func init() {
	Register(constant.SourceSearch, searchSourceFactory)
}

func searchSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.SearchFetcher == nil {
		return nil, fmt.Errorf("search_fetcher config is required for search source")
	}

	// Load global provider config to decide default parser
	db := util.GetDatabase()
	var providerConfig config.SearchProviderConfig
	if err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &providerConfig); err != nil {
		return nil, fmt.Errorf("failed to load search provider config: %w", err)
	}

	// Determine Parser
	var p parser.Parser
	if cfg.JsonParser != nil {
		p = &parser.JsonParser{Config: cfg.JsonParser}
	} else {
		// Get default parser config from provider
		// Use the same provider factory logic as the fetcher to ensure consistency
		prv, err := provider.Get(providerConfig.Provider, &providerConfig)
		if err == nil {
			defaultConfig := prv.GetDefaultParserConfig()
			if defaultConfig != nil {
				p = &parser.JsonParser{Config: defaultConfig}
			}
		}

		// Fallback if provider retrieval fails or returns nil (shouldn't happen with valid provider)
		if p == nil {
			defaultConfig := &config.JsonParserConfig{
				ItemsIterator: "data",
				Title:         "title",
				Link:          "url",
				Description:   "content",
			}
			p = &parser.JsonParser{Config: defaultConfig}
		}
	}

	return &PipelineSource{
		Config: cfg,
		Fetcher: &fetcher.CachedFetcher{
			Internal: &fetcher.SearchFetcher{
				Config:   cfg.SearchFetcher,
				Provider: providerConfig.Provider,
			},
			Expire: constant.SearchCacheExpire,
		},
		Parser: p,
	}, nil
}
