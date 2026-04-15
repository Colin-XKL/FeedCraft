package util

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var loadEnvOnce sync.Once

func GetEnvClient() *viper.Viper {
	// 确保 .env 文件仅加载一次，避免高频调用时的重复文件 I/O
	loadEnvOnce.Do(func() {
		if err := gotenv.Load(); err != nil {
			logrus.Debugf("gotenv.Load() skipped or failed: %v (this is normal if no .env file exists)", err)
		}
	})

	v := viper.New()

	// Set the name of the environment variable prefix (optional)
	v.SetEnvPrefix("FC")

	v.AutomaticEnv()
	return v
}
