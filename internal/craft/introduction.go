package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

/*
add introduction for article,powered by google gemini
*/

func getIntroductionForArticle(prompt, article string) (string, error) {
	return adapter.CallLLMUsingContext(prompt, article)
}

const promptGenerateIntroduction = "请阅读下面的文章并写一篇不超过200字的中文摘要, 使得读者可以快速知道文章的主题和主要结论."

func combineIntroductionAndArticle(article, intro string) string {
	introInHtml := util.Markdown2HTML(intro)
	return fmt.Sprintf(`<div>%s<div><hr/>%s</div>`, introInHtml, article)
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
		errMsg := "add introduction for article failed."
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
