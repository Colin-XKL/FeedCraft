package router

import (
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
	public := router.Group("/craft")
	{
		public.GET("/proxy", recipe.ProxyFeed)
		public.GET("/fulltext", recipe.ExtractFulltextForFeed)
	}
}
