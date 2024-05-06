package adapter

import (
	"fmt"
)

func CallGemini(prompt string, article string) (string, error) {
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
