package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
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
var (
	// llmClients acts as a lazy-loaded singleton registry, NOT a traditional acquire/release connection pool.
	// It maps configuration keys to a single llms.Model instance.
	// We reuse these instances so the underlying http.Transport can naturally maintain TCP Keep-Alive connections.
	llmClients sync.Map

	// llmDispatcher is a global priority queue that limits the maximum concurrent requests sent to LLM APIs.
	// This prevents rate-limit errors (429) while allowing retried requests to jump the queue.
	llmDispatcher *util.PriorityDispatcher[string]
	llmDispOnce   sync.Once
)

type llmRetryConfig struct {
	attemptsPerModel uint
	delay            time.Duration
	maxDelay         time.Duration
}

var defaultLLMRetryConfig = llmRetryConfig{
	attemptsPerModel: 3,
	delay:            2 * time.Second,
	maxDelay:         10 * time.Second,
}

func getLLMDispatcher() *util.PriorityDispatcher[string] {
	llmDispOnce.Do(func() {
		envClient := util.GetEnvClient()
		if envClient == nil {
			log.Fatalf("get env client error.")
		}
		concurrency := envClient.GetInt("LLM_MAX_CONCURRENCY")
		if concurrency <= 0 {
			concurrency = 3
		}
		llmDispatcher = util.NewPriorityDispatcher[string](concurrency)
		// Fallback timeout to prevent tasks from sticking forever
		// llmCallTimeout is 10 min, we set fallback to 11 min
		llmDispatcher.MaxTaskDuration = llmCallTimeout + time.Minute
		logrus.Infof("LLM Global Priority Dispatcher initialized with max concurrency: %d", concurrency)
	})
	return llmDispatcher
}

type llmModelClient struct {
	model string
	llm   llms.Model
}

func SimpleLLMCall(model string, promptInput string) (string, error) {
	return simpleLLMCall(model, promptInput, defaultLLMRetryConfig)
}

func simpleLLMCall(model string, promptInput string, retryConfig llmRetryConfig) (string, error) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}
	if retryConfig.attemptsPerModel == 0 {
		retryConfig.attemptsPerModel = defaultLLMRetryConfig.attemptsPerModel
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

	modelList := make([]string, 0)
	seenModels := make(map[string]struct{})
	for _, currentModel := range strings.Split(llmApiModel, ",") {
		currentModel = strings.TrimSpace(currentModel)
		if currentModel == "" {
			continue
		}
		if _, ok := seenModels[currentModel]; ok {
			continue
		}
		seenModels[currentModel] = struct{}{}
		modelList = append(modelList, currentModel)
	}
	rand.Shuffle(len(modelList), func(i, j int) {
		modelList[i], modelList[j] = modelList[j], modelList[i]
	})

	var lastErr error
	modelClients := make([]llmModelClient, 0, len(modelList))
	for _, currentModel := range modelList {
		cacheKey := fmt.Sprintf("%s|%s|%s|%s", llmApiType, llmApiBase, llmApiKey, currentModel)
		var llm llms.Model

		if cached, ok := llmClients.Load(cacheKey); ok {
			llm = cached.(llms.Model)
		} else {
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
			llmClients.Store(cacheKey, llm)
		}

		modelClients = append(modelClients, llmModelClient{model: currentModel, llm: llm})
	}

	if len(modelClients) == 0 {
		return "", fmt.Errorf("all models failed, last error: %v", lastErr)
	}

	dispatcher := getLLMDispatcher()
	isUrgent := false // Will be set to true on retry
	attemptIndex := 0
	currentModel := ""
	maxAttempts := uint(len(modelClients)) * retryConfig.attemptsPerModel

	result, err := retry.DoWithData(
		func() (string, error) {
			modelClient := modelClients[attemptIndex%len(modelClients)]
			attemptIndex++
			currentModel = modelClient.model

			return dispatcher.Execute(context.Background(), isUrgent, func(ctx context.Context) (string, error) {
				ctx, cancel := context.WithTimeout(ctx, llmCallTimeout)
				defer cancel()

				content := []llms.MessageContent{
					llms.TextParts(llms.ChatMessageTypeHuman, promptInput),
				}

				resp, err := modelClient.llm.GenerateContent(ctx, content)
				if err != nil {
					return "", err
				}
				if len(resp.Choices) > 0 {
					return resp.Choices[0].Content, nil
				}
				return "", fmt.Errorf("empty llm call response")
			})
		},
		retry.Attempts(maxAttempts),
		retry.DelayType(modelRotationBackoffDelay(len(modelClients))),
		retry.Delay(retryConfig.delay),
		retry.MaxDelay(retryConfig.maxDelay),
		retry.OnRetry(func(n uint, err error) {
			isUrgent = true // Elevate priority on retry
			nextModel := modelClients[attemptIndex%len(modelClients)].model
			logrus.Warnf("Retrying LLM call (attempt %d) with model %s, err from previous model %s: %v", n+1, nextModel, currentModel, err)
		}),
	)

	if err == nil {
		return result, nil
	}

	lastErr = err
	logrus.Warnf("LLM call failed after retries: %v", err)
	return "", fmt.Errorf("all models failed, last error: %v", lastErr)
}

func modelRotationBackoffDelay(modelCount int) retry.DelayTypeFunc {
	return func(n uint, err error, config *retry.Config) time.Duration {
		if modelCount <= 1 {
			return retry.BackOffDelay(n, err, config)
		}
		if n%uint(modelCount) != 0 {
			return 0
		}
		return retry.BackOffDelay(n/uint(modelCount), err, config)
	}
}
