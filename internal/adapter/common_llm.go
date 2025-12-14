package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
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

	// 4. Validate configuration for Ollama
	if llmApiType == "ollama" {
		logrus.Debug("Using Ollama API provider")
		if llmApiBase == "" {
			return "", fmt.Errorf("FC_LLM_API_BASE must be set when using FC_LLM_API_TYPE='ollama'")
		}
	} else {
		if llmApiBase == "" {
			logrus.Info("using default openai endpoint")
		}
	}

	modelList := strings.Split(llmApiModel, ",")
	rand.Shuffle(len(modelList), func(i, j int) {
		modelList[i], modelList[j] = modelList[j], modelList[i]
	})

	var lastErr error
	for _, currentModel := range modelList {
		currentModel = strings.TrimSpace(currentModel)
		if currentModel == "" {
			continue
		}

		var llm llms.Model
		var err error

		if llmApiType == "ollama" {
			llm, err = ollama.New(
				ollama.WithServerURL(llmApiBase),
				ollama.WithModel(currentModel),
			)
		} else {
			opts := []openai.Option{
				openai.WithToken(llmApiKey),
				openai.WithModel(currentModel),
			}
			if llmApiBase != "" {
				opts = append(opts, openai.WithBaseURL(llmApiBase))
			}
			llm, err = openai.New(opts...)
		}

		if err != nil {
			lastErr = err
			logrus.Warnf("Failed to initialize LLM client for model %s: %v", currentModel, err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), llmCallTimeout)

		content := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, promptInput),
		}

		resp, err := llm.GenerateContent(ctx, content)
		cancel()

		if err == nil {
			if len(resp.Choices) > 0 {
				return resp.Choices[0].Content, nil
			}
			err = fmt.Errorf("no choices in response")
		}

		lastErr = err
		logrus.Warnf("LLM call failed with model %s: %v", currentModel, err)
	}

	return "", fmt.Errorf("all models failed, last error: %v", lastErr)
}
