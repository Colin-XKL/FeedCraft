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
		public.GET("/craft-flows", controller.ListCraftFlows)
		public.GET("/sys-craft-atoms", controller.ListSysCraftAtoms)
		public.GET("/craft-atoms", controller.ListCraftAtoms)
	}

	craftRouters := router.Group("/craft")
	{
		craftRouters.GET("/:craft-name", craft.Entry)
	}
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

		adminApi.POST("/recipes", controller.CreateCustomRecipe)
		adminApi.GET("/recipes", controller.ListCustomRecipe)
		adminApi.GET("/recipes/:id", controller.GetCustomRecipe)
		adminApi.PUT("/recipes/:id", controller.UpdateCustomRecipe)
		adminApi.DELETE("/recipes/:id", controller.DeleteCustomRecipe)

		adminApi.POST("/users", controller.CreateUser)
		adminApi.GET("/users", controller.ListUsers)
		adminApi.GET("/users/:username", controller.GetUser)
		adminApi.PUT("/users/:username", controller.UpdateUser)
		adminApi.DELETE("/users/:username", controller.DeleteUser)

		adminApi.GET("/craft-flows", controller.ListCraftFlows)
		adminApi.POST("/craft-flows", controller.CreateCraftFlow)
		adminApi.GET("/craft-flows/:name", controller.GetCraftFlow)
		adminApi.PUT("/craft-flows/:name", controller.UpdateCraftFlow)
		adminApi.DELETE("/craft-flows/:name", controller.DeleteCraftFlow)

		adminApi.GET("/sys-craft-atoms", controller.ListSysCraftAtoms)
		adminApi.GET("/craft-templates", controller.ListCraftTemplates)

		adminApi.GET("/craft-atoms", controller.ListCraftAtoms)
		adminApi.POST("/craft-atoms", controller.CreateCraftAtom)
		adminApi.GET("/craft-atoms/:name", controller.GetCraftAtom)
		adminApi.PUT("/craft-atoms/:name", controller.UpdateCraftAtom)
		adminApi.DELETE("/craft-atoms/:name", controller.DeleteCraftAtom)
	}

}

func AdminUserInfoHandler(c *gin.Context) {
	resp := map[string]string{
		"name": "admin",
		"role": "admin",
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: resp})
}
