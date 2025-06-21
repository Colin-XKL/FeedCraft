package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

/*
add an introduction for article, using LLM power
*/

func getIntroductionForArticle(prompt, article string) (string, error) {
	return adapter.CallLLMUsingContext(prompt, article)
}

const promptGenerateIntroduction = `
你是一位专业的文章导读撰写专家，擅长用简洁的语言吸引读者注意力并概括文章核心内容。请根据用户提供的文章内容，生成一段言简意赅、引人入胜的文章导读。导读需满足以下要求：

吸引注意力：根据文章中的内容, 通过提问、引用数据、讲述故事或制造悬念等方式，激发读者的兴趣。
概括核心：用1-2句话清晰传达文章的主题或核心观点。
引导阅读：暗示文章的价值或结构，鼓励读者继续阅读。
语言风格：简洁有力，避免冗长或复杂表达。

输出要求：
语言要求：使用简体中文
字数限制：中文不超过120字
语言风格：简洁、生动、引人入胜。口语化但专业，避免术语堆砌。

请根据以下文章内容生成导读：

`

func combineIntroductionAndArticle(article, intro string) string {
	introInHtml := util.Markdown2HTML(intro)
	return fmt.Sprintf(`<div>%s<div><hr/><br/>%s</div>`, introInHtml, article)
}

func addIntroductionUsingLLM(item *feeds.Item, prompt string) string {
	//TODO handle description and content field separately and correctly
	finalArticleContent := ""
	originalContent := item.Content
	originalTitle := item.Title
	if len(originalContent) == 0 {
		if len(item.Description) == 0 {
			logrus.Warnf("empty content , both content as description field have no value. title [%s]", originalTitle)
		}
		originalContent = item.Description
		logrus.Warnf("empty content, use description field val as fallback")
	}

	domain, _ := util.ParseDomainFromUrl(item.Link.Href)
	cleanedArticleContent := util.Html2Markdown(originalContent, &domain)
	introduction := ""
	var err error
	if len(cleanedArticleContent) > 0 {
		introduction, err = getIntroductionForArticle(prompt, cleanedArticleContent)
	} else {
		introduction, err = getIntroductionForArticle(prompt, originalContent)
	}
	if err != nil {
		errMsg := fmt.Sprintf("add introduction for article failed. err: %v ", err)
		logrus.Warnf(errMsg)
		introduction = errMsg
	}

	finalArticleContent = combineIntroductionAndArticle(originalContent, introduction)
	return finalArticleContent
}

func GetAddIntroductionCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		ret := addIntroductionUsingLLM(item, prompt)
		return ret, nil
	}
	cachedTransformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "add intro")
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}

func introCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = promptGenerateIntroduction
	}
	return GetAddIntroductionCraftOptions(prompt)
}

var introCraftParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "the llm prompt for generate summary", Default: promptGenerateIntroduction},
}
