package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
)

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, context string, option util.ContentProcessOption) (string, error) {
	processedContext := util.ProcessContent(context, option)
	finalPrompt := fmt.Sprintf("%s \n```\n%s\n```", prompt, processedContext)
	cacheKey := fmt.Sprintf("llm_call_%s", util.GetMD5Hash(finalPrompt))
	valFunc := func() (string, error) {
		return SimpleLLMCall(UseDefaultModel, finalPrompt)
	}
	return util.CachedFunc(cacheKey, valFunc)
}
