package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/source/fetcher/provider"
	"FeedCraft/internal/util"
	"context"
	"fmt"
)

type SearchFetcher struct {
	Config   *config.SearchFetcherConfig
	Provider string
}

func (f *SearchFetcher) BaseURL() string {
	return "search://" + f.Provider + "/" + f.Config.Query
}

func (f *SearchFetcher) Fetch(ctx context.Context) ([]byte, error) {
	// 1. Get Provider Config
	db := util.GetDatabase()
	var providerConfig config.SearchProviderConfig
	if err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &providerConfig); err != nil {
		return nil, fmt.Errorf("failed to load search provider config: %w", err)
	}

	// 2. Get Provider Instance
	p, err := provider.Get(providerConfig.Provider, &providerConfig)
	if err != nil {
		return nil, err
	}

	// 3. Delegate Fetch
	return p.Fetch(ctx, f.Config.Query)
}
