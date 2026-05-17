package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
	"strings"
)

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, context string, option util.ContentProcessOption) (string, error) {
	processedContext := util.ProcessContent(context, option)
	// Remove backticks to avoid breaking the markdown code block
	processedContext = strings.ReplaceAll(processedContext, "`", "")

	finalPrompt := fmt.Sprintf("%s \n```\n%s\n```", prompt, processedContext)
	cacheKey := fmt.Sprintf("llm_call_%s", util.GetTextContentHash(finalPrompt))
	valFunc := func() (string, error) {
		return SimpleLLMCall(UseDefaultModel, finalPrompt)
	}
	return util.CachedFunc(cacheKey, valFunc)
}
