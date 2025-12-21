package craft

import (
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func OptionTimeLimit(days int) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		items := feed.Items
		if len(items) == 0 {
			return nil
		}

		// Check if there is at least one "normal" date
		// Normal date: !IsZero and Year > 1970
		hasNormalDate := false
		for _, item := range items {
			if !item.Created.IsZero() && item.Created.Year() > 1970 {
				hasNormalDate = true
				break
			}
		}

		// If all dates are abnormal (Zero or <= 1970), keep all
		if !hasNormalDate {
			logrus.Warnf("All items in feed have abnormal dates (<= 1970 or empty). Skipping time limit filter.")
			return nil
		}

		cutoff := time.Now().AddDate(0, 0, -days)

		filtered := lo.Filter(items, func(item *feeds.Item, index int) bool {
			t := item.Created
			// Always keep items with empty dates (though they might count as abnormal in the check above)
			// Note: If we are here, it means there IS some normal date in the feed.
			// So mixed case: Zero + 2024.
			// We keep Zero.
			if t.IsZero() {
				return true
			}

			// Filter out abnormal dates (<= 1970) if we are in filtering mode
			if t.Year() <= 1970 {
				return false // Drop 1970s if we have normal dates
			}

			return t.After(cutoff)
		})

		feed.Items = filtered
		return nil
	}
}

func GetTimeLimitCraftOption(days int) []CraftOption {
	return []CraftOption{
		OptionTimeLimit(days),
	}
}

func timeLimitCraftLoadParams(m map[string]string) []CraftOption {
	daysStr, exist := m["days"]
	if !exist {
		daysStr = "7"
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		logrus.Warnf("invalid param [days] for craft template [time-limit]")
		days = 7
	}
	if days < 0 {
		days = 7
	}
	return GetTimeLimitCraftOption(days)
}

var timeLimitCraftParamTmpl = []ParamTemplate{
	{Key: "days", Description: "Limit articles to the last N days", Default: "7"},
}
