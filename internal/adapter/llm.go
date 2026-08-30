package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, context string, option util.ContentProcessOption) (string, error) {
	processedContext := util.ProcessContent(context, option)
	// Remove backticks to avoid breaking the markdown code block
	processedContext = strings.ReplaceAll(processedContext, "`", "")

	finalPrompt := fmt.Sprintf("%s \n```\n%s\n```", prompt, processedContext)
	cacheKey := fmt.Sprintf("llm_call_%s_%s", util.GetTextContentHash(finalPrompt), llmCallTemperatureCachePart(option))
	valFunc := func() (string, error) {
		if option.Temperature != nil {
			return SimpleLLMCallWithOptions(UseDefaultModel, finalPrompt, llms.WithTemperature(*option.Temperature))
		}
		return SimpleLLMCall(UseDefaultModel, finalPrompt)
	}
	return util.CachedFunc(cacheKey, valFunc)
}

func llmCallTemperatureCachePart(option util.ContentProcessOption) string {
	if option.Temperature == nil {
		return "default-temperature"
	}
	return fmt.Sprintf("temperature:%g", *option.Temperature)
}
