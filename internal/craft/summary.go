package craft

import (
	"FeedCraft/internal/constant"
	"github.com/gorilla/feeds"
)

func GetSummaryCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		processorType := constant.ProcessorTypeSummary
		processor := NewLLMTextProcessor(processorType, prompt)
		ret := processItemContent(item, processor)
		return ret, nil
	}
	cachedTransformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, string(constant.ProcessorTypeSummary))
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}
func summaryCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = constant.DefaultPrompts[constant.ProcessorTypeSummary]
	}
	return GetSummaryCraftOptions(prompt)
}

var summaryCraftParamTmpl = []ParamTemplate{
	{
		Key: "prompt", Description: "the llm prompt for generate summary",
		Default: constant.DefaultPrompts[constant.ProcessorTypeSummary],
	},
}
