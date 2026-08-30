package craft

import (
	"github.com/gin-gonic/gin"
)

func GetProxyFeedHandler() func(c *gin.Context) {
	var craftOptions []LegacyCraftOption
	return func(c *gin.Context) {
		CommonCraftHandlerUsingCraftOptionList(c, craftOptions)
	}
}
