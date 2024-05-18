package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"log"
)

/**
Handle LLM calling and processing, support OpenAI and all compatible services.
*/

func SimpleLLMCall(model string, promptInput string) (string, error) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}

	openAIEndpoint := envClient.GetString("OPENAI_ENDPOINT")
	openAIAuthKey := envClient.GetString("OPENAI_AUTH_KEY")

	conf := openai.DefaultConfig(openAIAuthKey)
	if openAIEndpoint != "" {
		logrus.Info("using custom openai endpoint ", openAIEndpoint)
		conf.BaseURL = openAIEndpoint
	} else {
		logrus.Info("using default openai endpoint ")
	}
	//conf.HTTPClient=http.Client{}.

	client := openai.NewClientWithConfig(conf)
	if client == nil {
		return "", fmt.Errorf("new openai client error")
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: promptInput,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}
	return resp.Choices[0].Message.Content, nil
}
