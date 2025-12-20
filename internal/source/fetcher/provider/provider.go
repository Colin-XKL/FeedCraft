package provider

import (
	"FeedCraft/internal/config"
	"context"
	"fmt"
)

// SearchProvider defines the interface for different search provider implementations.
type SearchProvider interface {
	// Fetch executes the search request and returns the raw response body.
	Fetch(ctx context.Context, query string) ([]byte, error)

	// GetDefaultParserConfig returns the default JSON parser configuration
	// suitable for parsing the response from this provider.
	GetDefaultParserConfig() *config.JsonParserConfig
}

type FactoryFunc func(cfg *config.SearchProviderConfig) SearchProvider

var registry = make(map[string]FactoryFunc)

// Register registers a new search provider factory.
func Register(name string, factory FactoryFunc) {
	registry[name] = factory
}

// Get creates a new search provider instance by name.
// If name is empty, it defaults to "litellm" (for backward compatibility or default behavior).
func Get(name string, cfg *config.SearchProviderConfig) (SearchProvider, error) {
	if name == "" {
		name = "litellm"
	}
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("search provider '%s' not found", name)
	}
	return factory(cfg), nil
}
