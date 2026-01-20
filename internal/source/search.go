package source

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/fetcher/provider"
	"FeedCraft/internal/source/parser"
	"FeedCraft/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

func init() {
	Register(constant.SourceSearch, searchSourceFactory)
}

func searchSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.SearchFetcher == nil {
		return nil, fmt.Errorf("search_fetcher config is required for search source")
	}

	if cfg.SearchFetcher.EnhancedMode {
		return &EnhancedSearchSource{Config: cfg}, nil
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

type EnhancedSearchSource struct {
	Config *config.SourceConfig
}

func (s *EnhancedSearchSource) Generate(ctx context.Context) (*gofeed.Feed, error) {
	queries, err := expandQueryWithLLM(s.Config.SearchFetcher.Query)
	if err != nil {
		logrus.Warnf("LLM expansion failed: %v, falling back to original query", err)
		queries = []string{s.Config.SearchFetcher.Query}
	} else {
		if len(queries) == 0 {
			queries = []string{s.Config.SearchFetcher.Query}
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allItems []*gofeed.Item
	seenLinks := make(map[string]bool)

	for _, q := range queries {
		wg.Add(1)
		go func(queryStr string) {
			defer wg.Done()

			// Clone config and disable enhanced mode to avoid recursion
			newFetcherCfg := *s.Config.SearchFetcher
			newFetcherCfg.Query = queryStr
			newFetcherCfg.EnhancedMode = false

			newCfg := *s.Config
			newCfg.SearchFetcher = &newFetcherCfg

			src, err := searchSourceFactory(&newCfg)
			if err != nil {
				logrus.Warnf("Failed to create source for query %s: %v", queryStr, err)
				return
			}

			feed, err := src.Generate(ctx)
			if err != nil {
				logrus.Warnf("Generation failed for query %s: %v", queryStr, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			for _, item := range feed.Items {
				if !seenLinks[item.Link] {
					seenLinks[item.Link] = true
					allItems = append(allItems, item)
				}
			}
		}(q)
	}
	wg.Wait()

	return &gofeed.Feed{
		Title:       fmt.Sprintf("Search: %s", s.Config.SearchFetcher.Query),
		Description: fmt.Sprintf("Enhanced search results for %s", s.Config.SearchFetcher.Query),
		Items:       allItems,
	}, nil
}

func expandQueryWithLLM(query string) ([]string, error) {
	cacheKey := "search_expansion:" + query
	jsonStr, err := util.CachedFunc(cacheKey, func() (string, error) {
		prompt := fmt.Sprintf("Analyze the following search query and generate 3-5 distinct, optimized search queries to cover different aspects, languages (if implicit), and sub-topics. Original Query: %s. Return strictly a JSON array of strings, e.g., [\"query1\", \"query2\"]. Do not output any other text, markdown formatting, or code blocks. Just the raw JSON string.", query)
		return adapter.SimpleLLMCall(adapter.UseDefaultModel, prompt)
	})
	if err != nil {
		return nil, err
	}

	var queries []string
	if err := json.Unmarshal([]byte(jsonStr), &queries); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}
	return queries, nil
}

func (s *EnhancedSearchSource) BaseURL() string {
	return "search://" + s.Config.SearchFetcher.Query + "?enhanced=true"
}
