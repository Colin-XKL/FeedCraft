package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"github.com/gorilla/feeds"
)

const translateArticleContentPrompt = "下面是一篇文章的内容,请将其翻译为中文. 如果文章内有图片或者链接尽量保留它们, 对于专有名词也请保持原样. 注意只需要输出翻译后的文章内容即可，不要包括其他无关的内容。"

const translateArticleTitlePrompt = "下面是一篇文章的标题, 请将其翻译为中文. 对于专有名词等请保持原样。注意只需要输出一句翻译后的内容即可，不要包括其他无关的内容。"

const immersiveTranslatePrompt = `
你是一名精通多语言的翻译专家。请将输入的文章翻译为简体中文，按段落逐段处理，输出时每段原文后紧跟对应的中文译文，原文与译文之间留一空行。  

- **语言范围**：任意语言的文章均可接受；若输入已经是简体中文，则直接原样输出，不进行翻译。  
- **资源保留**：代码块、内嵌图片、视频、音频等非文本资源保持原样，不进行翻译，且位置不变。  
- **表格处理**：保留原文中的表格原样显示；在每个原始表格下方立即添加该表格的中文译本，保持相同的排版结构（行列、对齐、边框等）。  
- **格式保留**：完整保留原文的风格、语气以及所有排版格式（标题、子标题、项目符号列表、编号列表、代码块、表格等），使译文在版面上尽量与原文一致。  
- **专有名词与术语**：保持专有名词、技术术语和领域特定表达不变，除非该词已有广泛使用的中文译法。  
- **输出要求**：不添加任何额外标签、说明或评论  
- **长度限制**：不设长度上限，全文一次性输出。
`

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
func cacheKeyForArticleLink(item *feeds.Item) (string, error) {
	uniqLinkStr := item.Title
	uniqLinkStr += item.Id
	if item.Link != nil {
		uniqLinkStr += item.Link.Href
	} else if item.Source != nil {
		uniqLinkStr += item.Source.Href
	}
	return util.GetMD5Hash(uniqLinkStr), nil
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

// =======================================
// immersive translate article content
// ===

func immersiveTranslateLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = immersiveTranslatePrompt
	}
	return GetTranslateContentCraftOptions(prompt)
}

var immersiveTranslateParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "prompt for using llm do translate job", Default: immersiveTranslatePrompt},
}
