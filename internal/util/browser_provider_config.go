package util

import (
	"strings"

	"github.com/spf13/viper"
)

func IsSupportedBrowserProvider(provider string) bool {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case BrowserProviderBrowserlessRESTful, "browserless", BrowserProviderCDP:
		return true
	default:
		return false
	}
}

func ResolveBrowserProviderConfig(env *viper.Viper) BrowserProviderConfig {
	provider := strings.ToLower(strings.TrimSpace(env.GetString("BROWSER_PROVIDER")))
	endpoint := strings.TrimSpace(env.GetString("BROWSER_ENDPOINT"))
	if endpoint == "" {
		endpoint = strings.TrimSpace(env.GetString("PUPPETEER_HTTP_ENDPOINT"))
	}
	if provider == "" {
		provider = BrowserProviderBrowserlessRESTful
	}
	return BrowserProviderConfig{
		Provider: provider,
		Endpoint: endpoint,
	}
}

func resolveBrowserProviderConfig(env *viper.Viper) BrowserProviderConfig {
	return ResolveBrowserProviderConfig(env)
}
