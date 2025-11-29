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

	// 1. Load new standard env vars
	llmApiType := envClient.GetString("LLM_API_TYPE")
	if llmApiType == "" {
		llmApiType = "openai"
	}
	llmApiBase := envClient.GetString("LLM_API_BASE")
	llmApiKey := envClient.GetString("LLM_API_KEY")
	llmApiModel := model
	if llmApiModel == "" {
		llmApiModel = envClient.GetString("LLM_API_MODEL")
	}

	// 2. Load legacy env vars for compatibility
	legacyEndpoint := envClient.GetString("OPENAI_ENDPOINT")
	legacyAuthKey := envClient.GetString("OPENAI_AUTH_KEY")
	legacyModel := envClient.GetString("OPENAI_DEFAULT_MODEL")

	// 3. Fallback logic with deprecation warnings
	if llmApiBase == "" && legacyEndpoint != "" {
		llmApiBase = legacyEndpoint
		logrus.Warn("FC_OPENAI_ENDPOINT is deprecated, please migrate to FC_LLM_API_BASE")
	}

	if llmApiKey == "" && legacyAuthKey != "" {
		llmApiKey = legacyAuthKey
		logrus.Warn("FC_OPENAI_AUTH_KEY is deprecated, please migrate to FC_LLM_API_KEY")
	}

	if llmApiModel == "" && legacyModel != "" {
		llmApiModel = legacyModel
		logrus.Warn("FC_OPENAI_DEFAULT_MODEL is deprecated, please migrate to FC_LLM_API_MODEL")
	}

	// 4. Configure client based on type
	if llmApiType == "ollama" {
		logrus.Debug("Using Ollama API compatibility mode")
		if llmApiKey == "" {
			llmApiKey = "ollama" // Ollama doesn't require a key, but library might check
		}
	}

	conf := openai.DefaultConfig(llmApiKey)
	if llmApiBase != "" {
		conf.BaseURL = llmApiBase
	} else {
		if llmApiType == "openai" {
			logrus.Info("using default openai endpoint")
		}
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
			Model: llmApiModel,
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
