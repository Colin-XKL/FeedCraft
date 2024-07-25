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

const promptCheckIfAdvertorial = "请阅读下面的文章, 并判断是不是广告推销软文. 如果非常确信这篇文章是营销推广文章, 请返回 'true', 如果不是或者没有把握确定,请返回 'false'"

// CheckIfAdvertorial 判断是否为软文, 非常有把握则返回true, 如果不是或者不确定或是发生错误则返回false
func CheckIfAdvertorial(content string, prompt string) bool {
	const MinContentLength = 20
	if len(strings.TrimSpace(content)) < MinContentLength {
		return false
	}
	result, err := adapter.CallLLMUsingContext(prompt, content)
	if err != nil {
		logrus.Errorf("Error checking advertorial: %v", err)
		return false
	}
	logrus.Info("advertorial check: ", result)
	return result == "true"
}

func GetIgnoreAdvertorialCraftOptions(prompt string) []CraftOption {
	craftOptions := []CraftOption{
		OptionIgnoreAdvertorial(prompt),
	}
	return craftOptions
}

// OptionIgnoreAdvertorial option  判断一篇文章是不是推广软文和广告等
func OptionIgnoreAdvertorial(prompt string) CraftOption {
	return func(feed *feeds.Feed) error {
		items := feed.Items
		filtered := lo.Filter(items, func(item *feeds.Item, index int) bool {
			content := item.Content //TODO handle description and content field correctly
			return CheckIfAdvertorial(content, prompt)
		})
		feed.Items = filtered
		return nil
	}
}

func llmFilterCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt-for-exclude"]
	if !exist || len(prompt) == 0 {
		prompt = promptCheckIfAdvertorial
	}
	return GetIgnoreAdvertorialCraftOptions(prompt)
}

var llmFilterCraftParamTmpl = []ParamTemplate{
	{Key: "prompt-for-exclude",
		Description: "apply prompt to every article item, if llm returns 'true', then article will be excluded",
		Default:     promptCheckIfAdvertorial},
}

// ===============
// api for debug

type CheckIfAdvertorialDebugReq struct {
	Url         string `json:"url"  binding:"required,url" ` // article url
	EnhanceMode bool   `json:"enhance_mode"`
}
type CheckIfAdvertorialDebugResp struct {
	Url            string `json:"url"` // url for original article
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
	var webContent string
	if reqBody.EnhanceMode {
		webContent, err = getRenderedHTML2(reqBody.Url, 1*time.Minute)
	} else {
		webContent, err = TrivialExtractor(reqBody.Url, 1*time.Minute)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	if webContent == "" {
		c.JSON(http.StatusExpectationFailed, util.APIResponse[any]{Msg: "extract article content failed"})
		return
	}
	result := CheckIfAdvertorial(webContent, promptCheckIfAdvertorial)
	ret := CheckIfAdvertorialDebugResp{
		Url:            reqBody.Url,
		IsAdvertorial:  result,
		ArticleContent: webContent,
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: ret})
}
