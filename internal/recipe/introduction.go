package recipe

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

/*
add introduction for article,powered by google gemini
*/

func getIntroductionForArticle(prompt, article string) (string, error) {
	gemini, err := adapter.NewGeminiClient()
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("%s\n```%s```", prompt, article)
	req := adapter.GeminiReqPayload{Contents: []adapter.Content{
		{
			Parts: []adapter.Part{{Text: &text}},
			Role:  nil,
		},
	}}
	content, err := gemini.GenerateContent(req)
	if err != nil {
		return "", err
	}
	return content, nil
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

const promptGenerateIntroduction = "please summary the following article."

func addIntroductionUsingGemini(item *gofeed.Item) string {
	finalArticleContent := ""
	originalContent := item.Content
	originalTitle := item.Title
	if len(originalContent) == 0 {
		return ""
	}
	hashVal := getMD5Hash(originalContent)
	cachedIntroduction, err := util.CacheGetString(GetCacheKeyForWebContent(hashVal))
	if err != nil || cachedIntroduction == "" {
		//articleStr, err := extractor(url, DefaultTimeout)
		introduction, err := getIntroductionForArticle(promptGenerateIntroduction, originalContent)
		if err != nil {
			logrus.Warnf("failed to generate introduction for article %s, %v\n", originalTitle, err)
		} else {
			//TODO use html template rendering
			finalArticleContent = fmt.Sprintf(`<div><b>Introduction<b><br/><p>%s</p>`, introduction)
			cacheErr := util.CacheSetString(GetCacheKeyForWebContent(hashVal), introduction, constant.WebContentExpire)
			if cacheErr != nil {
				logrus.Warnf("failed to cache generated introduction for article %s, %v\n", originalTitle, cacheErr)
			}
		}
	} else {
		finalArticleContent = cachedIntroduction
	}
	return finalArticleContent
}

func AddIntroductionForFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(400, "empty feed url")
		return
	}
	fp := gofeed.NewParser()
	parsedFeed, _ := fp.ParseURL(feedUrl)

	ret := TransformFeed(parsedFeed, addIntroductionUsingGemini)

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
