package craft

import (
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"time"
)

type FulltextExtractor func(url string, timeout time.Duration) (string, error)

func TrivialExtractor(url string, timeout time.Duration) (string, error) {
	article, err := readability.FromURL(url, timeout)
	return article.Content, err
}

func GetFulltextHandler() func(c *gin.Context) {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return TrivialExtractor(link, DefaultExtractFulltextTimeout)
	}
	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "extract fulltext")
	craftOptions := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)),
	}
	return func(c *gin.Context) {
		CommonCraftHandlerUsingCraftOptionList(c, craftOptions)
	}
}
