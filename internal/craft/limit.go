package craft

import (
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

func GetLimitCraftOption() []CraftOption {
	craftOptions := []CraftOption{
		OptionLimit(defaultLimit),
	}
	return craftOptions
}
