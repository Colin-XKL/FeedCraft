package craft

import (
	"FeedCraft/internal/adapter"
	"github.com/gin-gonic/gin"
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

// GetTranslateTitleHandler translate title
func GetTranslateTitleHandler() func(c *gin.Context) {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleTitle(item.Title)
	}
	transformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "translate title")
	craftOption := []CraftOption{
		OptionTransformFeedItem(
			GetArticleTitleProcessor(transformer),
		),
	}
	return func(c *gin.Context) {
		CommonCraftHandlerUsingCraftOptionList(c, craftOption)
	}
}

// GetTranslateArticleContentHandler translate article content
func GetTranslateArticleContentHandler() func(c *gin.Context) {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleContent(item.Content) // TODO handle feed item content correctly
	}
	cachedTransformer := GetCommonCachedTransformer(
		cacheKeyForArticleContent, transFunc, "translate article content")
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return func(c *gin.Context) {
		CommonCraftHandlerUsingCraftOptionList(c, craftOption)
	}
}
