package craft

import (
	"FeedCraft/internal/util"
	"github.com/gorilla/feeds"
)

func CleanupContent(htmlContent string, domain string) (string, error) {
	// First convert HTML to Markdown to strip unnecessary elements
	markdown := util.HTMLToMarkdown(htmlContent, domain)

	cleanHtml := util.MarkdownToHTML(markdown)
	return cleanHtml, nil
}

func GetCleanupCraftOptions() []LegacyCraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		domain, _ := util.ParseDomainFromUrl(item.Link.Href)
		return CleanupContent(item.Content, domain)
	}
	cachedTransFunc := GetCommonCachedTransformer(
		cacheKeyForArticleContent,
		transFunc,
		"cleanup article content",
	)
	craftOptions := []LegacyCraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)),
	}
	return craftOptions
}
