package craft

import (
	"FeedCraft/internal/util"
	"github.com/gorilla/feeds"
)

func CleanupContent(htmlContent string, domain string) (string, error) {
	// First convert HTML to Markdown to strip unnecessary elements
	markdown := util.Html2Markdown(htmlContent, &domain)

	// Then convert back to HTML for clean, semantic markup
	cleanHtml := util.Markdown2HTML(markdown)
	return cleanHtml, nil
}

func GetCleanupCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		domain, _ := util.ParseDomainFromUrl(item.Link.Href)
		return CleanupContent(item.Content, domain)
	}
	cachedTransFunc := GetCommonCachedTransformer(
		cacheKeyForArticleContent,
		transFunc,
		"cleanup article content",
	)
	craftOptions := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)),
	}
	return craftOptions
}
