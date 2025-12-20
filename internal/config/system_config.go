package config

type SearchProviderConfig struct {
	APIUrl   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	Provider string `json:"provider"`
}
