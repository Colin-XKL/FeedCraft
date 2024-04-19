package router

import (
	"FeedCraft/internal/recipe"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine) {

	// Public routes
	public := router.Group("/craft")
	{
		public.GET("/proxy", recipe.ProxyFeed)
	}
}
