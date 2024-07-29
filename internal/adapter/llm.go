package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
)

// CallGeminiUsingArticleContext @deprecated
func CallGeminiUsingArticleContext(prompt string, article string) (string, error) {
	gemini, err := NewGeminiClient()
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("%s\n```%s```", prompt, article)
	req := GeminiReqPayload{Contents: []Content{
		{
			Parts: []Part{{Text: &text}},
			Role:  nil,
		},
	}}
	content, err := gemini.GenerateContent(req)
	if err != nil {
		return "", err
	}
	return content, nil
}

// CallLLMUsingContext using openai compatible api
func CallLLMUsingContext(prompt, context string) (string, error) {
	finalPrompt := fmt.Sprintf("%s \n `%s`", prompt, context)
	cacheKey := fmt.Sprintf("llm_call_%s", util.GetMD5Hash(finalPrompt))
	valFunc := func() (string, error) {
		return SimpleLLMCall(UseDefaultModel, finalPrompt)
	}
	return util.CachedFunc(cacheKey, valFunc)
}
