package main

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/router"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"net/http"
	"os"
)
func init() {
	logrus.Info("Starting PreheatingScheduler...")
	// 设置预热任务函数
	taskFunc := func(recipeName string) error {
		recipeData, err := recipe.QueryCustomRecipeName(recipeName)
		if err != nil {
			return err
		}
		path := recipe.GetPathForCustomRecipe(recipeData)
		_, err = recipe.RetrieveCraftRecipeUsingPath(path)
		return err
	}
	recipe.Scheduler = util.NewPreheatingScheduler(taskFunc)
	logrus.Info("Start PreheatingScheduler done.")
}

var rootCmd = &cobra.Command{
	Use:   "feedcraft",
	Short: "FeedCraft is a feed management system",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset admin password",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dao.ResetAdminPassword(); err != nil {
			logrus.Errorf("Failed to reset admin password: %v", err)
			os.Exit(1)
		}
		logrus.Info("Admin password has been reset successfully")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(resetPasswordCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer() {
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
