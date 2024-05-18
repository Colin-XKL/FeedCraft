package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

/**
通过LLM 判断并排除广告软文 advertorial
*/

const prompt = "请阅读下面的文章, 并判断是不是广告推销软文. 如果非常确信这篇文章是营销推广文章, 请返回 'true', 如果不是或者没有把握确定,请返回 'false'"

// CheckIfAdvertorial 判断是否为软文, 非常有把握则返回true, 如果不是或者不确定或是发生错误则返回false
func CheckIfAdvertorial(content string) bool {
	const MinContentLength = 20
	if len(strings.TrimSpace(content)) < MinContentLength {
		return false
	}
	//result, err := adapter.CallGeminiUsingArticleContext(prompt, content)
	result, err := adapter.CallLLMUsingContext(prompt, content)
	if err != nil {
		logrus.Errorf("Error checking advertorial: %v", err)
		return false
	}
	logrus.Info("advertorial check: ", result)
	return result == "true"
}

func OptionIgnoreAdvertorial() CraftOption {
	return func(feed *feeds.Feed) error {
		items := feed.Items
		filtered := lo.Filter(items, func(item *feeds.Item, index int) bool {
			content := item.Content //TODO handle description and content field correctly
			return CheckIfAdvertorial(content)
		})
		feed.Items = filtered
		return nil
	}
}

func IgnoreAdvertorialArticle(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(400, "empty feed url")
		return
	}
	craftedFeed, err := NewCraftedFeedFromUrl(feedUrl, OptionIgnoreAdvertorial())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	rssStr, err := craftedFeed.OutputFeed.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}

type CheckIfAdvertorialDebugReq struct {
	Url string `json:"url"` // article url
}
type CheckIfAdvertorialDebugResp struct {
	Url            string `json:"url"` // url for orignial article
	ArticleContent string `json:"article_content"`
	IsAdvertorial  bool   `json:"is_advertorial"`
}

// DebugCheckIfAdvertorial input: article url, output: article content text and a bool represent if this article is marked as ad
func DebugCheckIfAdvertorial(c *gin.Context) {
	reqBody := &CheckIfAdvertorialDebugReq{}
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	webContent, err := TrivialExtractor(reqBody.Url, 1*time.Minute)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	result := CheckIfAdvertorial(webContent)
	ret := CheckIfAdvertorialDebugResp{
		Url:            reqBody.Url,
		IsAdvertorial:  result,
		ArticleContent: webContent,
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: ret})
}
