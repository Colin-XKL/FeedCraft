package craft

import (
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
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
		filtered := lo.Filter(items, func(item *feeds.Item, index int) bool {
			content := item.Content
			if len(content) == 0 {
				content = item.Description
			}
			// We use the generic prompt wrapper which asks LLM to return true if it matches the condition
			match, err := CheckConditionWithGenericPrompt(content, condition)
			if err != nil {
				// On error, we keep the item (fail open)
				return true
			}
			// If matches condition (true), we want to EXCLUDE it.
			// So return false to Filter.
			return !match
		})
		feed.Items = filtered
		return nil
	}
}
