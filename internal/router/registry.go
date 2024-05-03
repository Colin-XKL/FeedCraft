package router

import (
	"FeedCraft/internal/admin"
	"FeedCraft/internal/middleware"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterRouters(router *gin.Engine) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}
	siteBaseUrl := envClient.GetString("SITE_BASE_URL")
	router.LoadHTMLFiles("web/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"SiteBaseUrl": siteBaseUrl,
		})
	})

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", admin.LoginAuth)
	}

	craftRouters := router.Group("/craft")
	{
		craftRouters.GET("/proxy", recipe.ProxyFeed)
		craftRouters.GET("/fulltext", recipe.ExtractFulltextForFeed)
		craftRouters.GET("/fulltext-plus", recipe.ExtractFulltextPlusForFeed)
		craftRouters.GET("/introduction", recipe.AddIntroductionForFeed)
	}

	// admin api
	adminApi := router.Group("/api/admin")
	adminApi.Use(middleware.JwtAuthMiddleware())
	{
		adminApi.GET("/admin-login-test", adminLoginTest)
	}
}
func adminLoginTest(c *gin.Context) {
	ret := map[string]string{
		"success": "ok",
	}
	c.JSON(http.StatusOK, ret)
}
