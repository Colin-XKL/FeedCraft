package craft

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
)

const defaultLimit = 10

func OptionLimit(n int) CraftOption {
	return func(feed *feeds.Feed) error {
		items := feed.Items
		filtered := lo.Slice(items, 0, n)
		feed.Items = filtered
		return nil
	}
}

func GetLimitHandler() func(c *gin.Context) {
	craftOptions := []CraftOption{
		OptionLimit(defaultLimit),
	}
	return func(c *gin.Context) {
		CommonCraftHandlerUsingCraftOptionList(c, craftOptions)
	}
}
