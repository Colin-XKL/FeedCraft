package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	reqBody := map[string]interface{}{
		"query": f.Config.Query,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", providerConfig.APIUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if providerConfig.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+providerConfig.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search provider returned status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
