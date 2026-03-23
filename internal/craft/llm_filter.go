package craft

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"strings"
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
			safeContent := strings.ReplaceAll(content, "```", "")
			textToEvaluate := fmt.Sprintf("```markdown\nTitle: %s\n\nContent:\n%s\n```", itm.Title, safeContent)
			match, err := CheckConditionWithGenericPrompt(textToEvaluate, condition)
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
