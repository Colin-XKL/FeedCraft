package router

import (
	"FeedCraft/internal/admin"
	"FeedCraft/internal/controller"
	"FeedCraft/internal/craft"
	"FeedCraft/internal/middleware"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/util"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	corsMiddleware := cors.New(corsConfig)
	//corsMiddleware := cors.Default()
	router.Use(corsMiddleware)

	router.Static("/assets", "./web/assets")

	//siteBaseUrl := envClient.GetString("SITE_BASE_URL")
	router.LoadHTMLFiles("web/index.html")
	router.LoadHTMLFiles("web/start.html")
	router.StaticFile("/start.html", "web/start.html")
	//router.GET("/start.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "start.html", gin.H{
	//		"SiteBaseUrl": siteBaseUrl,
	//	})
	//})

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Page not found"})
			return
		}
		c.File("./web/index.html")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", controller.LoginAuth)
		public.GET("/list-all-craft", controller.ListAllCraft)
	}

	craftRouters := router.Group("/craft")
	{
		craftRouters.GET("/:craft-name", craft.Entry)
	}
	
	// channel routes (New standard)
	channelRoutes := router.Group("/channel")
	{
		channelRoutes.GET("/:id", recipe.CustomRecipe)
	}

	// recipe routes (Legacy compatibility)
	recipeRoutes := router.Group("/recipe")
	{
		recipeRoutes.GET("/:id", recipe.CustomRecipe)
	}

	// admin api
	adminApi := router.Group("/api/admin")
	adminApi.Use(middleware.JwtAuthMiddleware(), corsMiddleware)
	{
		adminApi.POST("/user/info", AdminUserInfoHandler)
		adminApi.POST("/user/change-password", controller.ChangePassword)

		adminApi.POST("/craft-debug/advertorial", craft.DebugCheckIfAdvertorial)
		adminApi.POST("/craft-debug/common-llm-call-test", admin.LLMDebug)

		// Channels (formerly recipes)
		adminApi.POST("/channels", controller.CreateChannel)
		adminApi.GET("/channels", controller.ListChannels)
		adminApi.GET("/channels/:id", controller.GetChannel)
		adminApi.PUT("/channels/:id", controller.UpdateChannel)
		adminApi.DELETE("/channels/:id", controller.DeleteChannel)

		// Legacy recipe routes aliases (optional, but good for admin panel compatibility if frontend not fully updated instantly)
		// But plan says "Action: Update internal/router/registry.go and corresponding internal/controller files."
		// And Frontend refactor is Phase 3.
		// If I change Admin API routes now, the Frontend (Phase 3) will break until updated.
		// However, Phase 2 goal includes updating Admin API.
		// So I will assume Frontend will be broken or I should keep aliases?
		// "Phase 3: Frontend Refactor... Update API clients".
		// This implies I should BREAK the old API or PROVIDE ALIASES.
		// Ideally, keep aliases for Admin API too?
		// The plan doesn't explicitly say "keep aliases for Admin API", but it's safer.
		// However, the instructions say: "GET /api/admin/recipes -> /api/admin/channels".
		// This usually means REPLACE.
		// But since I am an agent, I should follow instructions. Instructions say "->", usually implies replacement.
		// But I will keep the old ones commented out or just replace them.
		// Wait, if I replace them, the current frontend (Phase 3 pending) will break.
		// "Suggested in a separate branch". I am working on the main codebase effectively.
		// I'll add the new ones AND keep the old ones as deprecated/aliases pointing to the NEW controllers.
		// This ensures the frontend doesn't break before Phase 3.
		
		// Legacy Admin API Aliases (Mapped to new Controllers)
		adminApi.POST("/recipes", controller.CreateChannel)
		adminApi.GET("/recipes", controller.ListChannels)
		adminApi.GET("/recipes/:id", controller.GetChannel)
		adminApi.PUT("/recipes/:id", controller.UpdateChannel)
		adminApi.DELETE("/recipes/:id", controller.DeleteChannel)


		adminApi.POST("/users", controller.CreateUser)
		adminApi.GET("/users", controller.ListUsers)
		adminApi.GET("/users/:username", controller.GetUser)
		adminApi.PUT("/users/:username", controller.UpdateUser)
		adminApi.DELETE("/users/:username", controller.DeleteUser)

		// Blueprints (formerly craft-flows)
		adminApi.GET("/blueprints", controller.ListBlueprints)
		adminApi.POST("/blueprints", controller.CreateBlueprint)
		adminApi.GET("/blueprints/:name", controller.GetBlueprint)
		adminApi.PUT("/blueprints/:name", controller.UpdateBlueprint)
		adminApi.DELETE("/blueprints/:name", controller.DeleteBlueprint)

		// Legacy Admin API Aliases
		adminApi.GET("/craft-flows", controller.ListBlueprints)
		adminApi.POST("/craft-flows", controller.CreateBlueprint)
		adminApi.GET("/craft-flows/:name", controller.GetBlueprint)
		adminApi.PUT("/craft-flows/:name", controller.UpdateBlueprint)
		adminApi.DELETE("/craft-flows/:name", controller.DeleteBlueprint)

		adminApi.GET("/sys-tools", controller.ListSysTools)
		adminApi.GET("/tool-templates", controller.ListToolTemplates)
		
		// Legacy Admin API Aliases
		adminApi.GET("/sys-craft-atoms", controller.ListSysTools)
		adminApi.GET("/craft-templates", controller.ListToolTemplates)


		// Tools (formerly craft-atoms)
		adminApi.GET("/tools", controller.ListTools)
		adminApi.POST("/tools", controller.CreateTool)
		adminApi.GET("/tools/:name", controller.GetTool)
		adminApi.PUT("/tools/:name", controller.UpdateTool)
		adminApi.DELETE("/tools/:name", controller.DeleteTool)

		// Legacy Admin API Aliases
		adminApi.GET("/craft-atoms", controller.ListTools)
		adminApi.POST("/craft-atoms", controller.CreateTool)
		adminApi.GET("/craft-atoms/:name", controller.GetTool)
		adminApi.PUT("/craft-atoms/:name", controller.UpdateTool)
		adminApi.DELETE("/craft-atoms/:name", controller.DeleteTool)


		adminApi.POST("/tools/fetch", controller.HtmlFetch)
		adminApi.POST("/tools/parse", controller.HtmlParse)

		adminApi.POST("/tools/json/fetch", controller.CurlFetch)
		adminApi.POST("/tools/json/parse", controller.CurlParse)
		adminApi.POST("/tools/json/parse_curl", controller.CurlParseCmd)

		adminApi.POST("/tools/search/preview", controller.SearchPreview)

		// Settings Routes
		adminApi.GET("/settings/search-provider", controller.GetSearchProviderConfig)
		adminApi.POST("/settings/search-provider", controller.SaveSearchProviderConfig)

		adminApi.GET("/dependencies", controller.GetDependencyStatus)
		adminApi.POST("/dependencies/check", controller.CheckDependencyStatus)
	}

}

func AdminUserInfoHandler(c *gin.Context) {
	resp := map[string]string{
		"name": "admin",
		"role": "admin",
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: resp})
}