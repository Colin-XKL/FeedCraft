package config

type SearchProviderConfig struct {
	APIUrl   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	Provider string `json:"provider"`
}

func (c *SearchProviderConfig) Mask() {
	if c.APIKey != "" {
		c.APIKey = "******"
	}
}

func (c *SearchProviderConfig) Merge(newConfig SearchProviderConfig) {
	c.APIUrl = newConfig.APIUrl
	c.Provider = newConfig.Provider
	if newConfig.APIKey != "******" && newConfig.APIKey != "" {
		c.APIKey = newConfig.APIKey
	}
}
