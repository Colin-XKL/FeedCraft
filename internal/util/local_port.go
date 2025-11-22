package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// GetLocalPort 这里的local port 是gin会额外监听的一个端口,用于向自身发送请求时使用
func GetLocalPort() int {
	const localDefaultPort = 1025
	localPortEnv := os.Getenv("LOCAL_PORT")
	if localPortEnv == "" {
		return localDefaultPort
	}
	localPort, err := strconv.Atoi(localPortEnv)
	if err != nil {
		logrus.Errorf("invalid value for LOCAL_PORT environment variable: %v", err)
		logrus.Warnf("using default port %d", localDefaultPort)
		localPort = localDefaultPort
	}
	return localPort
}
