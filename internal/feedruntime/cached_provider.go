package feedruntime

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"FeedCraft/internal/constant"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/sirupsen/logrus"
)

// SubFeedHealth records health metadata for a single topic sub-feed URI, persisted in Redis.
type SubFeedHealth struct {
	URI           string     `json:"uri"`
	LastSuccessAt *time.Time `json:"last_success_at"`
	LastFailureAt *time.Time `json:"last_failure_at"`
	LastError     string     `json:"last_error"`
	// CachedAt is when the last successful result was written to the cache.
	CachedAt *time.Time `json:"cached_at"`
}

// CachedFeedProvider wraps any FeedProvider with optimistic caching for topic aggregation.
// On a successful fetch the result is stored in Redis (SubFeedCacheTTL = 14 days).
// On failure the last cached result is returned instead so the parent TopicFeed stays
// available, and the failure is recorded in the health metadata.
type CachedFeedProvider struct {
	// URI is used as the stable identifier for cache key derivation and health tracking.
	URI   string
	Inner engine.FeedProvider
}

// Fetch implements engine.FeedProvider.
func (p *CachedFeedProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	feed, err := p.Inner.Fetch(ctx)
	if err == nil {
		p.cacheResult(feed)
		p.recordSuccess()
		return feed, nil
	}

	// Live fetch failed — try the cached snapshot.
	cached := p.loadCachedResult()
	p.recordFailure(err)
	if cached != nil {
		logrus.Infof("CachedFeedProvider [%s]: live fetch failed (%v), returning cached snapshot", p.URI, err)
		return cached, nil
	}

	// No cache available; propagate the original error.
	return nil, err
}

// cacheResult writes feed to Redis with SubFeedCacheTTL.
func (p *CachedFeedProvider) cacheResult(feed *model.CraftFeed) {
	data, err := json.Marshal(feed)
	if err != nil {
		logrus.Warnf("CachedFeedProvider [%s]: marshal failed: %v", p.URI, err)
		return
	}
	if err := util.TryCacheSetString(subFeedResultKey(p.URI), string(data), constant.SubFeedCacheTTL); err != nil {
		logrus.Warnf("CachedFeedProvider [%s]: cache write failed: %v", p.URI, err)
	}
}

// loadCachedResult reads the most recent cached feed from Redis.
func (p *CachedFeedProvider) loadCachedResult() *model.CraftFeed {
	raw, err := util.TryCacheGetString(subFeedResultKey(p.URI))
	if err != nil || raw == "" {
		return nil
	}
	var feed model.CraftFeed
	if err := json.Unmarshal([]byte(raw), &feed); err != nil {
		logrus.Warnf("CachedFeedProvider [%s]: unmarshal cached feed failed: %v", p.URI, err)
		return nil
	}
	return &feed
}

// recordSuccess updates the health metadata after a live successful fetch.
func (p *CachedFeedProvider) recordSuccess() {
	now := time.Now()
	h := p.loadHealth()
	h.URI = p.URI
	h.LastSuccessAt = &now
	h.CachedAt = &now
	h.LastError = ""
	p.saveHealth(h)
}

// recordFailure updates the health metadata after a fetch error.
func (p *CachedFeedProvider) recordFailure(err error) {
	now := time.Now()
	h := p.loadHealth()
	h.URI = p.URI
	h.LastFailureAt = &now
	h.LastError = err.Error()
	p.saveHealth(h)
}

func (p *CachedFeedProvider) loadHealth() SubFeedHealth {
	raw, err := util.TryCacheGetString(subFeedHealthKey(p.URI))
	if err != nil || raw == "" {
		return SubFeedHealth{URI: p.URI}
	}
	var h SubFeedHealth
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return SubFeedHealth{URI: p.URI}
	}
	return h
}

func (p *CachedFeedProvider) saveHealth(h SubFeedHealth) {
	data, err := json.Marshal(h)
	if err != nil {
		return
	}
	if err := util.TryCacheSetString(subFeedHealthKey(p.URI), string(data), constant.SubFeedHealthTTL); err != nil {
		logrus.Warnf("CachedFeedProvider [%s]: health write failed: %v", p.URI, err)
	}
}

// GetSubFeedHealth returns the stored health metadata for a given sub-feed URI.
// Returns an empty SubFeedHealth (URI only) if no data has been recorded yet.
func GetSubFeedHealth(uri string) SubFeedHealth {
	raw, err := util.TryCacheGetString(subFeedHealthKey(uri))
	if err != nil || raw == "" {
		return SubFeedHealth{URI: uri}
	}
	var h SubFeedHealth
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return SubFeedHealth{URI: uri}
	}
	return h
}

func subFeedResultKey(uri string) string {
	h := sha256.Sum256([]byte(uri))
	return fmt.Sprintf("%s:%x", constant.PrefixSubFeedResult, h)
}

func subFeedHealthKey(uri string) string {
	h := sha256.Sum256([]byte(uri))
	return fmt.Sprintf("%s:%x", constant.PrefixSubFeedHealth, h)
}
