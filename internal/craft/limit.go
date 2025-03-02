package craft

import (
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strconv"
)

const defaultLimit = 10

func OptionLimit(n int) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		items := feed.Items
		filtered := lo.Slice(items, 0, n)
		feed.Items = filtered
		return nil
	}
}

func GetLimitCraftOption(num int) []CraftOption {
	craftOptions := []CraftOption{
		OptionLimit(num),
	}
	return craftOptions
}

func limitCraftLoadParams(m map[string]string) []CraftOption {
	numStr, exist := m["num"]
	if !exist {
		numStr = "10"
	}
	n, err := strconv.Atoi(numStr)
	if err != nil {
		logrus.Warnf("invalid param [num] for craft template [limit]")
		n = defaultLimit
	}
	return GetLimitCraftOption(n)
}

var limitCraftParamTmpl = []ParamTemplate{
	{Key: "num", Description: "limit article to $num"},
}
