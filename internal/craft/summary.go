package craft

import (
	"FeedCraft/internal/constant"
	"github.com/gorilla/feeds"
)

func GetSummaryCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		processorType := "add-summary"
		processor := NewLLMTextProcessor(constant.ProcessorType(processorType), prompt)
		ret := processItemContent(item, processor)
		return ret, nil
	}
	cachedTransformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, "add summary")
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
	return GetAddIntroductionCraftOptions(prompt)
}

var summaryCraftParamTmpl = []ParamTemplate{
	{
		Key: "prompt", Description: "the llm prompt for generate introduction",
		Default: constant.DefaultPrompts[constant.ProcessorTypeSummary],
	},
}
