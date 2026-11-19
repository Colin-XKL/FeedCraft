package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, context string, option util.ContentProcessOption) (string, error) {
	if option.RemoveImage {
		context = util.RemoveBase64Images(context)
	}

	compactor := util.NewLLMContentCompactor()
	compactor.Reserve(prompt)
	compactedContext := compactor.Compact(context)

	processedContext := util.ProcessContent(compactedContext, option)
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
	result, err := util.CachedFunc(cacheKey, valFunc)
	if err != nil {
		return "", err
	}
	return compactor.Restore(result), nil
}

func llmCallTemperatureCachePart(option util.ContentProcessOption) string {
	if option.Temperature == nil {
		return "default-temperature"
	}
	return fmt.Sprintf("temperature:%g", *option.Temperature)
}
