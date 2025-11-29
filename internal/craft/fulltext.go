package craft

import (
	"FeedCraft/internal/util"
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"time"
)

type FulltextExtractor func(url string, timeout time.Duration) (string, error)

func TrivialExtractor(url string, timeout time.Duration) (string, error) {
	// The underlying readability library uses http.Client internally if not specified,
	// but FromURL takes a timeout. We should ensure it's robust.
	// Actually FromURL uses http.Get which uses DefaultClient but it sets a timeout on the context.
	// Wait, checking readability docs (or guessing): FromURL(url, timeout)
	// To strictly control the client, we might need to use FromReader.
	// But `readability.FromURL` signature is `func FromURL(url string, timeout time.Duration) (Article, error)`.
	// It likely handles timeout.
	// Let's rely on the passed timeout, which should be util.ExternalRequestTimeout.

	article, err := readability.FromURL(url, timeout)
	return article.Content, err
}

func GetFulltextCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return TrivialExtractor(link, util.ExternalRequestTimeout)
	}
	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext")
	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()
	var craftOptions []CraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}
