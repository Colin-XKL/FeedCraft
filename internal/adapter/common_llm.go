package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

/*
*
Handle LLM calling and processing, support OpenAI and all compatible services.
*/

const UseDefaultModel = ""

var llmCallTimeout = 10 * time.Minute

func SimpleLLMCall(model string, promptInput string) (string, error) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}

	openAIEndpoint := envClient.GetString("OPENAI_ENDPOINT")
	openAIAuthKey := envClient.GetString("OPENAI_AUTH_KEY")
	openAIModel := model
	if openAIModel == "" {
		openAIModel = envClient.GetString("OPENAI_DEFAULT_MODEL")
		//logrus.Info("model not specified, using default model from env:", openAIModel)
	}

	conf := openai.DefaultConfig(openAIAuthKey)
	if openAIEndpoint != "" {
		//logrus.Info("using custom openai endpoint ", openAIEndpoint)
		conf.BaseURL = openAIEndpoint
	} else {
		logrus.Info("using default openai endpoint ")
	}

	client := openai.NewClientWithConfig(conf)
	if client == nil {
		return "", fmt.Errorf("new openai client error")
	}
	ctx, cancel := context.WithTimeout(context.Background(), llmCallTimeout)
	defer cancel()
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openAIModel,
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
