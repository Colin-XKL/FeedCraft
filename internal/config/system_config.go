package config

type SearchProviderConfig struct {
	APIUrl   string        `json:"api_url"`
	APIKey   string        `json:"api_key"`
	Provider string        `json:"provider"`
	LiteLLM  LiteLLMConfig `json:"litellm"`
	SearXNG  SearXNGConfig `json:"searxng"`
}

type LiteLLMConfig struct {
	SearchToolName string `json:"search_tool_name"`
}

type SearXNGConfig struct {
	Engines string `json:"engines"`
}

type FaviconSettings struct {
	DefaultProviderID string                  `json:"default_provider_id"`
	CustomProviders   []FaviconProviderConfig `json:"custom_providers,omitempty"`
}

type FaviconProviderConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URLTemplate string `json:"url_template"`
	Enabled     bool   `json:"enabled"`
}
