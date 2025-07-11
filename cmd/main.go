package main

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/router"
	"FeedCraft/internal/util"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

func init() {
	logrus.Info("Preheating scheduler starting...")
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
	logrus.Info("Preheating scheduler started.")
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
		logrus.Info("Admin password reset successfully.")
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
	if len(env) == 0 { // set env to `prod` or `dev`
		env = "prod"
	}
	isProd := env == "prod"
	if !isProd {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if len(sentryDsn) > 0 {
		logrus.Info("Initializing Sentry...")
		sampledRate := 1.0
		if isProd {
			sampledRate = 0.1
		}
		// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
		err := sentry.Init(sentry.ClientOptions{
			Dsn:           sentryDsn,
			EnableTracing: true,
			Environment:   env,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: sampledRate,
		})
		if err != nil {
			logrus.Warnf("Sentry initialization failed: %v", err)
		} else {
			logrus.Info("Sentry initialized.")
		}
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	if len(sentryDsn) > 0 {
		r.Use(sentrygin.New(sentrygin.Options{}))
	}

	router.RegisterRouters(r)
	dao.MigrateDatabases()
	logrus.Info("Database migration done.")

	localDefaultPort := util.GetLocalPort() // 让gin额外监听的一个端口,用于向自身发送请求时使用
	listenAddr := os.Getenv("LISTEN_ADDR")
	go func() {
		_ = r.Run(fmt.Sprintf("localhost:%d", localDefaultPort))
	}()

	if !isProd {
		logrus.Info("Pprof server starting on :6060...")
		go func() {
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				logrus.Errorf("pprof server failed to start: %v", err)
			}
		}()
	}
	fmt.Print("=================================================================================")
	fmt.Print(`
==  ███████╗███████╗███████╗██████╗  ██████╗██████╗  █████╗ ███████╗████████╗
==  ██╔════╝██╔════╝██╔════╝██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝
==  █████╗  █████╗  █████╗  ██║  ██║██║     ██████╔╝███████║█████╗     ██║   
==  ██╔══╝  ██╔══╝  ██╔══╝  ██║  ██║██║     ██╔══██╗██╔══██║██╔══╝     ██║   
==  ██║     ███████╗███████╗██████╔╝╚██████╗██║  ██║██║  ██║██║        ██║   
==  ╚═╝     ╚══════╝╚══════╝╚═════╝  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝        ╚═╝   
==  
==                           Welcome to FeedCraft!
== Project Homepage: https://github.com/Colin-XKL/FeedCraft
`)
	fmt.Println("== Server listen at ", listenAddr)
	fmt.Println("== Admin Default User: admin\n== Default Password: adminadmin")
	fmt.Println("== Enjoy!")
	fmt.Println("=================================================================================")
	_ = r.Run(listenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
