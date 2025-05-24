package craft

import (
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"time"
)

type FulltextExtractor func(url string, timeout time.Duration) (string, error)

func TrivialExtractor(url string, timeout time.Duration) (string, error) {
	article, err := readability.FromURL(url, timeout)
	return article.Content, err
}

func GetFulltextCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return TrivialExtractor(link, DefaultExtractFulltextTimeout)
	}
	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext")
	craftOptions := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)),
	}
	return craftOptions
}
