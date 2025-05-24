package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
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

const judgePrompt = `
# 如何判断一篇文章是不是营销软文？
1. 内容是否围绕某单个产品或品牌展开.
如果文章的核心内容始终围绕某单个特定品牌或产品，尤其是多次提及该品牌或产品的优点，很可能是软文。

2. 是否有明显的推广意图.
虽然软文不会直接说“买这个产品”，但会通过暗示、推荐或引导的方式，让读者对产品产生兴趣。比如，文章可能会非常强调某个产品的独特功能、优惠活动或用户好评。

3. 是否带有购买链接或引导性行动
软文会在文章末尾或中间插入购买链接、优惠码、下载链接等，引导读者采取行动。

4. 内容是否过于正面或缺乏客观性
软文通常会刻意突出产品的优点，而忽略或淡化缺点。如果一篇文章对某个产品的评价过于正面，缺乏客观分析，可能是软文。

4. 是否有夸大或煽动性语言.
软文常常使用夸张的语言或煽动性的表达，比如“颠覆性创新”“行业领先”“唯一选择”等，以吸引读者注意。
`

const promptCheckIfAdvertorial = "请阅读下面的文章, 并判断是不是广告推销软文. 如果非常确信这篇文章是营销推广文章, 请返回 'true', 不要包括其他内容. 如果不是或者没有把握确定,请返回 'false'"

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
	logrus.Infof("advertorial check: [%s]", result)
	return strings.TrimSpace(result) == "true"
}

func GetIgnoreAdvertorialCraftOptions(prompt string) []CraftOption {
	craftOptions := []CraftOption{
		OptionIgnoreAdvertorial(prompt),
	}
	return craftOptions
}

// OptionIgnoreAdvertorial option  判断一篇文章是不是推广软文和广告等
func OptionIgnoreAdvertorial(prompt string) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
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
		prompt = fmt.Sprintf("%s\n%s\n", judgePrompt, promptCheckIfAdvertorial)
	}
	return GetIgnoreAdvertorialCraftOptions(prompt)
}

var llmFilterCraftParamTmpl = []ParamTemplate{
	{Key: "prompt-for-exclude",
		Description: "apply prompt to every article item, if llm returns 'true', then article will be excluded",
		Default:     fmt.Sprintf("%s\n%s\n", judgePrompt, promptCheckIfAdvertorial)},
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
	prompt := fmt.Sprintf("%s\n%s\n", judgePrompt, promptCheckIfAdvertorial)

	result := CheckIfAdvertorial(webContent, prompt)
	ret := CheckIfAdvertorialDebugResp{
		Url:            reqBody.Url,
		IsAdvertorial:  result,
		ArticleContent: webContent,
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: ret})
}
