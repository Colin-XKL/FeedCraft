package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type SearchFetcher struct {
	Config *config.SearchFetcherConfig
}

func (f *SearchFetcher) BaseURL() string {
	return "search://" + f.Config.Query
}

func (f *SearchFetcher) Fetch(ctx context.Context) ([]byte, error) {
	// 1. Get Provider Config
	db := util.GetDatabase()
	var providerConfig config.SearchProviderConfig
	if err := dao.GetJsonValue(db, constant.KeySearchProviderConfig, &providerConfig); err != nil {
		return nil, fmt.Errorf("failed to load search provider config: %w", err)
	}

	if providerConfig.APIUrl == "" {
		return nil, fmt.Errorf("search provider API URL is not configured")
	}

	// 2. Prepare Request for LiteLLM Proxy (or similar)
	client := resty.New()
	client.SetTimeout(constant.GlobalHttpTimeout)

	req := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"query": f.Config.Query,
		})

	if providerConfig.APIKey != "" {
		req.SetHeader("Authorization", "Bearer "+providerConfig.APIKey)
	}

	resp, err := req.Post(providerConfig.APIUrl)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("search provider returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
