package craft

import (
	"github.com/gorilla/feeds"
	"internal/util"
)

func CleanupContent(htmlContent string) (string, error) {
	// First convert HTML to Markdown to strip unnecessary elements
	markdown := util.Html2Markdown(htmlContent, nil)
	
	// Then convert back to HTML for clean, semantic markup
	cleanHtml := util.Markdown2HTML(markdown)
	return cleanHtml, nil
}

func GetCleanupCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		return CleanupContent(item.Content)
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
