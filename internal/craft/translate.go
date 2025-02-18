package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"github.com/gorilla/feeds"
)

const translateArticleContentPrompt = "下面是一篇文章的内容,请将其翻译为中文. 如果文章内有图片或者链接尽量保留它们, 对于专有名词也请保持原样. 注意只需要输出翻译后的文章内容即可，不要包括其他无关的内容。"

const translateArticleTitlePrompt = "下面是一篇文章的标题, 请将其翻译为中文. 对于专有名词等请保持原样。注意只需要输出一句翻译后的内容即可，不要包括其他无关的内容。"

func translateArticleTitle(title string, prompt string) (string, error) {
	return adapter.CallLLMUsingContext(prompt, title)
}
func translateArticleContent(content string, prompt string) (string, error) {
	return adapter.CallLLMUsingContext(prompt, content)
}

type ContentCacheKeyGenerator TransFunc

func cacheKeyForArticleTitle(item *feeds.Item) (string, error) {
	return util.GetMD5Hash(item.Title), nil
}
func cacheKeyForArticleContent(item *feeds.Item) (string, error) {
	return util.GetMD5Hash(item.Description + item.Description), nil
}

// =======================================
// translate article title
// ===

// GetTranslateTitleCraftOptions translate title
func GetTranslateTitleCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleTitle(item.Title, prompt)
	}
	transformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "translate title")
	craftOption := []CraftOption{
		OptionTransformFeedItem(
			GetArticleTitleProcessor(transformer),
		),
	}
	return craftOption
}

func transTitleCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = translateArticleTitlePrompt
	}
	return GetTranslateTitleCraftOptions(prompt)
}

var transTitleParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "prompt for using llm do translate job", Default: translateArticleTitlePrompt},
}

// =======================================
// translate article content
// ===

//todo 后续添加mode字段, 支持将摘要放在文章开头/文章结尾/替换掉原文

// GetTranslateContentCraftOptions translate article content
func GetTranslateContentCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		return translateArticleContent(item.Content, prompt) // TODO handle feed item content correctly
	}
	cachedTransformer := GetCommonCachedTransformer(
		cacheKeyForArticleContent, transFunc, "translate article content")
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}
func transContentCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = translateArticleContentPrompt
	}
	return GetTranslateContentCraftOptions(prompt)
}

var transContentParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "prompt for using llm do translate job", Default: translateArticleContentPrompt},
}
