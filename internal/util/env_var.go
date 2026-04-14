package util

import (
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func GetEnvClient() *viper.Viper {
	// Load .env file if it exists
	gotenv.Load()

	v := viper.New()

	// Set the name of the environment variable prefix (optional)
	v.SetEnvPrefix("FC")

	v.AutomaticEnv()
	return v
}
