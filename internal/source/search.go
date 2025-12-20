package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/source/fetcher"
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
	// We might use providerConfig.Provider to switch mapping logic in future
	db := util.GetDatabase()
	var providerConfig config.SearchProviderConfig
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &providerConfig)

	// Determine Parser
	var p parser.Parser
	if cfg.JsonParser != nil {
		p = &parser.JsonParser{Config: cfg.JsonParser}
	} else {
		// Default mapping for LiteLLM / generic search
		defaultConfig := &config.JsonParserConfig{
			ItemsIterator: "data",
			Title:         "title",
			Link:          "url",
			Description:   "content",
		}
		p = &parser.JsonParser{Config: defaultConfig}
	}

	return &PipelineSource{
		Config:  cfg,
		Fetcher: &fetcher.SearchFetcher{Config: cfg.SearchFetcher},
		Parser:  p,
	}, nil
}
