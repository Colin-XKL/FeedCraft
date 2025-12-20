package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type SearchFetcher struct {
	Config *config.SearchFetcherConfig
}

func (f *SearchFetcher) BaseURL() string {
	if f.Config == nil {
		return "search://unknown"
	}
	return "search://" + f.Config.Query
}

func (f *SearchFetcher) Fetch(ctx context.Context) ([]byte, error) {
	if f.Config == nil || f.Config.Query == "" {
		return nil, fmt.Errorf("search fetcher is not configured with a query")
	}

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

	body := resp.Body()
	if len(body) == 0 {
		return nil, fmt.Errorf("search provider returned empty body")
	}

	// Optionally validate valid JSON
	var js json.RawMessage
	if err := json.Unmarshal(body, &js); err != nil {
		return nil, fmt.Errorf("search provider returned invalid JSON: %w", err)
	}

	return body, nil
}
