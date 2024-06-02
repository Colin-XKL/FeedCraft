package craft

import (
	"FeedCraft/internal/adapter"
	"github.com/gorilla/feeds"
)

const translateArticleContentPrompt = "下面是一篇文章的内容,请将其翻译为中文. 如果文章内有图片或者链接尽量保留它们. "

const translateArticleTitlePrompt = "下面是一篇文章的标题, 请将其翻译为中文"

func translateArticleTitle(title string) (string, error) {
	return adapter.CallLLMUsingContext(translateArticleTitlePrompt, title)
}
func translateArticleContent(content string) (string, error) {
	return adapter.CallLLMUsingContext(translateArticleContentPrompt, content)
}

type ContentCacheKeyGenerator TransFunc

func cacheKeyForArticleTitle(item *feeds.Item) (string, error) {
	return getMD5Hash(item.Title), nil
}
func cacheKeyForArticleContent(item *feeds.Item) (string, error) {
	return getMD5Hash(item.Description + item.Description), nil
}

// GetTranslateTitleCraftOptions translate title
func GetTranslateTitleCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleTitle(item.Title)
	}
	transformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "translate title")
	craftOption := []CraftOption{
		OptionTransformFeedItem(
			GetArticleTitleProcessor(transformer),
		),
	}
	return craftOption
}

// GetTranslateContentCraftOptions translate article content
func GetTranslateContentCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleContent(item.Content) // TODO handle feed item content correctly
	}
	cachedTransformer := GetCommonCachedTransformer(
		cacheKeyForArticleContent, transFunc, "translate article content")
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}
