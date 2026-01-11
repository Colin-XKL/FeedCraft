package fetcher

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"time"
)

// Fetcher handles the I/O, just retrieving the raw binary data.
type Fetcher interface {
	Fetch(ctx context.Context) ([]byte, error)
	BaseURL() string
}

type CachedFetcher struct {
	Internal Fetcher
	Expire   time.Duration
}

func (f *CachedFetcher) BaseURL() string {
	return f.Internal.BaseURL()
}

func (f *CachedFetcher) Fetch(ctx context.Context) ([]byte, error) {
	cacheKey := fmt.Sprintf("%s:%s", constant.PrefixSearchSource, f.BaseURL())

	cached, err := util.CacheGetString(cacheKey)
	if err == nil && cached != "" {
		return []byte(cached), nil
	}

	data, err := f.Internal.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	_ = util.CacheSetString(cacheKey, string(data), f.Expire)
	return data, nil
}