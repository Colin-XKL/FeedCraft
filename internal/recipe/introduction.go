package recipe

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

/*
add introduction for article,powered by google gemini
*/

func getIntroductionForArticle(prompt, article string) (string, error) {
	return adapter.CallLLMUsingContext(prompt, article)
}

const promptGenerateIntroduction = "请阅读下面的文章并写一篇不超过200字的摘要,使得读者可以快速知道文章的主题和主要结论."

func addIntroductionUsingGemini(item *gofeed.Item) string {
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
		//TODO use html template rendering
		return fmt.Sprintf(`<div><b>Introduction<b><br/><p>%s</p><div>%s</div>`, intro, article)
	}

	if err != nil || cachedIntroduction == "" {
		//articleStr, err := extractor(url, DefaultTimeout)
		introduction, err := getIntroductionForArticle(promptGenerateIntroduction, originalContent)
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

//TODO handle description and content field separately and correctly

func AddIntroductionForFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(400, "empty feed url")
		return
	}
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feedUrl)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	ret := TransformFeed(parsedFeed, addIntroductionUsingGemini)

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
