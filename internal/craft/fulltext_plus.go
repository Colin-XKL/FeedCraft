package craft

import (
	"FeedCraft/internal/util"
	"net/url"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)


func getRenderedHTML2(websiteUrl string, timeout time.Duration) (string, error) {
	content, err := util.GetBrowserlessContent(websiteUrl, timeout)
	if err != nil {
		return "", err
	}

	parseUrl, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

	article, err := readability.FromReader(strings.NewReader(content), parseUrl)
	if err != nil {
		return "", err
	}
	return article.Content, err
}

func GetFulltextPlusCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return getRenderedHTML2(link, DefaultExtractFulltextTimeout)
	}

	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext plus")

	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()

	var craftOptions []CraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}
