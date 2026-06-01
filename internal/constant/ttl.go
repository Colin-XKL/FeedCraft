package constant

import "time"

const FEED_EXPIRE = 3 * time.Minute
const WebContentExpire = 7 * 24 * time.Hour
const SearchCacheExpire = 10 * time.Minute

// SubFeedCacheTTL is the TTL for cached sub-feed results under topic aggregation.
// When a sub-feed fails to fetch, the cached result is used as a fallback for up to this duration.
const SubFeedCacheTTL = 14 * 24 * time.Hour

// SubFeedHealthTTL is the TTL for per-URI health metadata stored in Redis.
const SubFeedHealthTTL = 30 * 24 * time.Hour
