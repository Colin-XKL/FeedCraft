package craft

import (
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"net/http"
	"time"
)

// llm-filter generic implementation

var llmFilterGenericParamTmpl = []ParamTemplate{
	{
		Key:         "filter_condition",
		Description: "Condition to filter out articles. If the article matches this condition (LLM returns true), it will be REMOVED. Example: 'Is this article about sports?'",
		Default:     "Is this content spam or low quality?",
	},
}

func llmFilterGenericLoadParam(m map[string]string) []CraftOption {
	condition, exist := m["filter_condition"]
	if !exist || len(condition) == 0 {
		condition = "Is this content spam or low quality?"
	}
	return GetLLMFilterGenericOptions(condition)
}

func GetLLMFilterGenericOptions(condition string) []CraftOption {
	return []CraftOption{
		OptionLLMFilterGeneric(condition),
	}
}

func OptionLLMFilterGeneric(condition string) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		items := feed.Items
		if len(items) == 0 {
			return nil
		}

		// 1. 并发请求 LLM 判断
		matches := parallel.Map(items, func(itm *feeds.Item, _ int) bool {
			content := itm.Content
			if len(content) == 0 {
				content = itm.Description
			}
			match, err := CheckConditionWithGenericPrompt(content, condition)
			if err != nil {
				return false
			}
			return match
		})

		// 2. 同步过滤
		feed.Items = lo.Filter(items, func(_ *feeds.Item, index int) bool {
			return !matches[index]
		})

		return nil
	}
}

// llm-filter-debug logic: Url mode
func DebugLLMFilterUrl(c *gin.Context) {
	var req struct {
		Url             string `json:"url"`
		EnhanceMode     bool   `json:"enhance_mode"`
		FilterCondition string `json:"filter_condition"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Invalid JSON format: " + err.Error()})
		return
	}

	if req.Url == "" || req.FilterCondition == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "url and filter_condition are required"})
		return
	}

	var webContent string
	var err error
	if req.EnhanceMode {
		webContent, err = getRenderedHTML2(req.Url, util.BrowserlessOptions{
			Timeout: 1 * time.Minute,
		})
	} else {
		webContent, err = TrivialExtractor(req.Url, 1*time.Minute)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to fetch content: " + err.Error()})
		return
	}
	if webContent == "" {
		c.JSON(http.StatusExpectationFailed, util.APIResponse[any]{Msg: "extract article content failed"})
		return
	}

	isFiltered, err := CheckConditionWithGenericPrompt(webContent, req.FilterCondition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "LLM Check Failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{
		Data: map[string]interface{}{
			"article_content": webContent,
			"is_filtered":     isFiltered,
		},
	})
}

// llm-filter-debug logic: Text mode
func DebugLLMFilterText(c *gin.Context) {
	var req struct {
		Text            string `json:"text"`
		FilterCondition string `json:"filter_condition"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Invalid JSON format: " + err.Error()})
		return
	}

	if req.Text == "" || req.FilterCondition == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "text and filter_condition are required"})
		return
	}

	isFiltered, err := CheckConditionWithGenericPrompt(req.Text, req.FilterCondition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "LLM Check Failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{
		Data: map[string]interface{}{
			"is_filtered": isFiltered,
		},
	})
}

// llm-filter-debug logic: Feed mode
func DebugLLMFilterFeed(c *gin.Context) {
	var req struct {
		FeedUrl         string `json:"feed_url"`
		FilterCondition string `json:"filter_condition"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Invalid JSON format: " + err.Error()})
		return
	}

	if req.FeedUrl == "" || req.FilterCondition == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "feed_url and filter_condition are required"})
		return
	}

	craftedFeed, err := NewCraftedFeedFromUrl(req.FeedUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: fmt.Sprintf("failed to fetch feed: %v", err)})
		return
	}
	if craftedFeed.OutputFeed == nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "fetched feed is nil"})
		return
	}

	type FeedItemResult struct {
		Title      string `json:"title"`
		Link       string `json:"link"`
		Content    string `json:"content"`
		IsFiltered bool   `json:"is_filtered"`
	}

	var results []FeedItemResult
	items := craftedFeed.OutputFeed.Items

	// Process items sequentially to avoid LLM rate limits.
	// If a user is debugging a feed, they expect accurate results for every item,
	// and bounded LLM concurrency is better handled by internal/util if global,
	// but sequential is safer for an interactive debug endpoint.
	for _, item := range items {
		content := item.Content
		if len(content) == 0 {
			content = item.Description
		}

		isFiltered, llmErr := CheckConditionWithGenericPrompt(content, req.FilterCondition)

		// Instead of failing the whole feed, append errors to content/title so the user knows this item failed.
		if llmErr != nil {
			content = fmt.Sprintf("[LLM ERROR: %v]\n\n%s", llmErr, content)
			isFiltered = false
		}

		link := ""
		if item.Link != nil {
			link = item.Link.Href
		}
		results = append(results, FeedItemResult{
			Title:      item.Title,
			Link:       link,
			Content:    content,
			IsFiltered: isFiltered,
		})
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{
		Data: results,
	})
}
