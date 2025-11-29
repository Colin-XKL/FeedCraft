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

	openAIEndpoint := envClient.GetString("OPENAI_ENDPOINT")
	openAIAuthKey := envClient.GetString("OPENAI_AUTH_KEY")
	openAIModel := model
	if openAIModel == "" {
		openAIModel = envClient.GetString("OPENAI_DEFAULT_MODEL")
	}

	conf := openai.DefaultConfig(openAIAuthKey)
	if openAIEndpoint != "" {
		conf.BaseURL = openAIEndpoint
	} else {
		logrus.Info("using default openai endpoint ")
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
				Model: openAIModel,
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
			return resp.Choices[0].Message.Content, nil
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
	return resp.Choices[0].Message.Content, nil
}
