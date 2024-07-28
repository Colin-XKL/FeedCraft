package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/constant"
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
	logrus.Infof("generate introduction using gemini for article [%s]", originalTitle)

	hashVal := getMD5Hash(originalContent)
	craftName := "introduction"
	cachedIntroduction, err := util.CacheGetString(getCacheKey(craftName, hashVal))

	combineIntroductionAndArticle := func(article, intro string) string {
		introInHtml := util.Markdown2HTML(intro)
		return fmt.Sprintf(`<div>%s<div>%s</div>`, introInHtml, article)
	}

	if err != nil || cachedIntroduction == "" {
		//articleStr, err := extractor(url, DefaultExtractFulltextTimeout)
		introduction, err := getIntroductionForArticle(prompt, originalContent)
		if err != nil {
			logrus.Warnf("failed to generate introduction for article [%s], %v\n", originalTitle, err)
		} else {
			finalArticleContent = combineIntroductionAndArticle(originalContent, introduction)
			cacheErr := util.CacheSetString(getCacheKey(craftName, hashVal), introduction, constant.WebContentExpire)
			if cacheErr != nil {
				logrus.Warnf("failed to cache generated introduction for article [%s], %v\n", originalTitle, cacheErr)
			}
		}
	} else {
		finalArticleContent = combineIntroductionAndArticle(originalContent, cachedIntroduction)
	}
	return finalArticleContent
}

func GetAddIntroductionCraftOptions(prompt string) []CraftOption {
	//todo 后续在将原文发送到LLM之前, 默认去掉无效的html属性和css以节省token
	transFunc := func(item *feeds.Item) (string, error) {
		ret := addIntroductionUsingLLM(item, prompt)
		return ret, nil
	}
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(transFunc)),
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
