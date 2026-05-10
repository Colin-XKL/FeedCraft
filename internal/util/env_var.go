package util

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultFeedUserAgent = "FeedCraft/2.0"
	htmlDefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
)

func DefaultFeedUserAgent() string {
	value := strings.TrimSpace(GetEnvClient().GetString("HTTP_USER_AGENT_FEED"))
	if value == "" {
		return defaultFeedUserAgent
	}
	return value
}

func DefaultHTMLUserAgent() string {
	value := strings.TrimSpace(GetEnvClient().GetString("HTTP_USER_AGENT_HTML"))
	if value == "" {
		return htmlDefaultUserAgent
	}
	return value
}

func GetEnvClient() *viper.Viper {
	v := viper.New()

	// Set the name of the environment variable prefix (optional)
	v.SetEnvPrefix("FC")

	v.AutomaticEnv()
	return v
}
