package util

import (
	"github.com/spf13/viper"
)

func GetEnvClient() *viper.Viper {
	v := viper.New()

	// Set the name of the environment variable prefix (optional)
	v.SetEnvPrefix("FC")

	v.AutomaticEnv()
	return v
}
