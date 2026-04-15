package util

import (
	"os"
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
			// 区分"文件不存在"（正常）和"文件存在但加载失败"（异常）
			if _, statErr := os.Stat(".env"); os.IsNotExist(statErr) {
				logrus.Debugf("gotenv.Load() skipped: .env file not found (this is normal)")
			} else {
				logrus.Warnf("gotenv.Load() failed: .env file exists but could not be loaded: %v", err)
			}
		}
	})

	v := viper.New()

	// Set the name of the environment variable prefix (optional)
	v.SetEnvPrefix("FC")

	v.AutomaticEnv()
	return v
}
