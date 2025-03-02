package craft

import (
	"github.com/gorilla/feeds"
)

// GetRelativeLinkFixCraftOptions 将feed项中的相对链接转换为绝对路径
func GetRelativeLinkFixCraftOptions() []CraftOption {
	wrapper := func(feed *feeds.Feed, payload ExtraPayload) error {
		feedLinkAttr := feed.Link.Href
		processor := func(feedItem *feeds.Item, payload ExtraPayload) error {
			absUrl := getAbsLinkForFeedItem(payload.originalFeedUrl, feedLinkAttr, feedItem.Link.Href)
			if len(absUrl) > 0 {
				feedItem.Link.Href = absUrl
			}
			return nil
		}

		transformer := OptionTransformFeedItem(processor)
		return transformer(feed, payload)
	}
	craftOptions := []CraftOption{
		wrapper,
	}
	return craftOptions
}
