package main

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/router"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	sentryDsn := os.Getenv("SENTRY_DSN")
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "prod"
	}
	if len(sentryDsn) > 0 {
		logrus.Info("initializing sentry...")
		// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
		err := sentry.Init(sentry.ClientOptions{
			Dsn:           sentryDsn,
			EnableTracing: true,
			Environment:   env,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			logrus.Warnf("sentry initialization failed: %v\n", err)
		} else {
			logrus.Info("sentry initialized.")
		}
	}

	r := gin.Default()

	if len(sentryDsn) > 0 {
		r.Use(sentrygin.New(sentrygin.Options{}))
	}

	router.RegisterRouters(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Migrate the schema
	dao.MigrateDatabases()
	localDefaultPort := util.GetLocalPort()
	listenAddr := os.Getenv("LISTEN_ADDR")
	go func() {
		_ = r.Run(fmt.Sprintf("localhost:%d", localDefaultPort))
	}()
	_ = r.Run(listenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
