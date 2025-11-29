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

var llmCallTimeout = util.LLMRequestTimeout

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
		if llmApiBase == "" {
			return "", fmt.Errorf("FC_LLM_API_BASE must be set when using FC_LLM_API_TYPE='ollama'")
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

	// Retry logic for 429
	var resp openai.ChatCompletionResponse
	var err error

	for i := 0; i <= util.MaxRetries; i++ {
		if i > 0 {
			time.Sleep(util.RetryWaitTime * time.Duration(i))
			logrus.Warnf("Retrying LLM request (attempt %d/%d)", i, util.MaxRetries)
		}

		ctx, cancel := context.WithTimeout(context.Background(), llmCallTimeout)
		resp, err = client.CreateChatCompletion(
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
		cancel()

		if err == nil {
			if len(resp.Choices) > 0 {
				return resp.Choices[0].Message.Content, nil
			}
			return "", fmt.Errorf("ChatCompletion error: no choices in response")
		}

		// Check for 429 Too Many Requests
		if apiErr, ok := err.(*openai.APIError); ok {
			if apiErr.HTTPStatusCode == 429 {
				logrus.Warnf("LLM API Rate limit exceeded (429), retrying...")
				continue
			}
		}

		// Check if it's a timeout, we might want to retry on timeout too if desired,
        // but requirements emphasized 429. Standard HTTP errors often don't get exposed as *APIError by this lib easily for net errors.
		// However, given the context, we simply break for non-429 errors unless we want robust retry for everything.
        // Let's stick to 429 specific handling as requested + maybe 5xx if possible to detect.
        // go-openai returns *APIError for API responses.

        if apiErr, ok := err.(*openai.APIError); ok {
            if apiErr.HTTPStatusCode >= 500 {
                 logrus.Warnf("LLM API Server Error (%d), retrying...", apiErr.HTTPStatusCode)
                 continue
            }
        }

		// Stop retrying for other errors
		break
	}

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}
	return "", fmt.Errorf("ChatCompletion error: unexpected flow")
}
