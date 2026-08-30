package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, articleContext string, option util.ContentProcessOption) (string, error) {
	return CallLLMUsingRequestContext(context.Background(), prompt, articleContext, option)
}

// CallLLMUsingRequestContext is CallLLMUsingContext bound to a request/task context.
func CallLLMUsingRequestContext(ctx context.Context, prompt, articleContext string, option util.ContentProcessOption) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	processedContext := util.ProcessContent(articleContext, option)
	// Remove backticks to avoid breaking the markdown code block
	processedContext = strings.ReplaceAll(processedContext, "`", "")
	processedContext = util.TruncateHeadTail(processedContext, util.LLMPromptMaxChars())

	finalPrompt := fmt.Sprintf("%s \n```\n%s\n```", prompt, processedContext)
	cacheKey := fmt.Sprintf("llm_call_%s_%s", util.GetTextContentHash(finalPrompt), llmCallTemperatureCachePart(option))
	valFunc := func() (string, error) {
		if option.Temperature != nil {
			return SimpleLLMCallWithOptionsContext(ctx, UseDefaultModel, finalPrompt, llms.WithTemperature(*option.Temperature))
		}
		return SimpleLLMCallContext(ctx, UseDefaultModel, finalPrompt)
	}
	return util.CachedFunc(cacheKey, valFunc)
}

func llmCallTemperatureCachePart(option util.ContentProcessOption) string {
	if option.Temperature == nil {
		return "default-temperature"
	}
	return fmt.Sprintf("temperature:%g", *option.Temperature)
}
