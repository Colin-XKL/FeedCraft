package adapter

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

/*
 * Embedding 适配器层
 * 提供统一的 Embedding 接口，支持 OpenAI、Gemini（通过 OpenAI 兼容接口）、Ollama 三种后端。
 * 锚点向量缓存在进程内存中，不依赖外部向量数据库。
 */

const (
	defaultEmbeddingModel = "text-embedding-3-small"
	embeddingCallTimeout  = 2 * time.Minute
)

var (
	// embeddingClients 缓存已创建的 Embedding 客户端实例（惰性单例）
	embeddingClients sync.Map

	// anchorVectorCache 锚点向量内存缓存，key = MD5(锚点文本) + "|" + 模型名称
	anchorVectorCache sync.Map
)

// embeddingConfig 从环境变量读取的 Embedding 配置
type embeddingConfig struct {
	apiType     string
	apiBase     string
	apiKey      string
	apiModel    string
	instruction string // 全局默认 instruction
}

// loadEmbeddingConfig 读取 Embedding 环境变量，未配置时回退使用 LLM 配置
func loadEmbeddingConfig() embeddingConfig {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}

	cfg := embeddingConfig{}

	// 1. 读取 Embedding 专用环境变量
	cfg.apiType = envClient.GetString("EMBEDDING_API_TYPE")
	cfg.apiBase = envClient.GetString("EMBEDDING_API_BASE")
	cfg.apiKey = envClient.GetString("EMBEDDING_API_KEY")
	cfg.apiModel = envClient.GetString("EMBEDDING_API_MODEL")
	cfg.instruction = envClient.GetString("EMBEDDING_INSTRUCTION")

	// 2. 回退逻辑：未配置时使用 LLM 配置
	if cfg.apiType == "" {
		cfg.apiType = "openai"
	}

	if cfg.apiBase == "" {
		cfg.apiBase = envClient.GetString("LLM_API_BASE")
		if cfg.apiBase != "" {
			logrus.Debug("FC_EMBEDDING_API_BASE not set, falling back to FC_LLM_API_BASE")
		}
	}

	if cfg.apiKey == "" {
		cfg.apiKey = envClient.GetString("LLM_API_KEY")
		if cfg.apiKey != "" {
			logrus.Debug("FC_EMBEDDING_API_KEY not set, falling back to FC_LLM_API_KEY")
		}
	}

	if cfg.apiModel == "" {
		cfg.apiModel = defaultEmbeddingModel
		logrus.Debugf("FC_EMBEDDING_API_MODEL not set, using default: %s", defaultEmbeddingModel)
	}

	return cfg
}

// getOrCreateEmbedder 获取或创建 Embedding 客户端（带缓存）
func getOrCreateEmbedder(cfg embeddingConfig) (*embeddings.EmbedderImpl, error) {
	cacheKey := fmt.Sprintf("emb|%s|%s|%s|%s", cfg.apiType, cfg.apiBase, cfg.apiKey, cfg.apiModel)

	if cached, ok := embeddingClients.Load(cacheKey); ok {
		return cached.(*embeddings.EmbedderImpl), nil
	}

	var client embeddings.EmbedderClient
	var err error

	switch cfg.apiType {
	case "ollama":
		logrus.Debug("Creating Ollama embedding client")
		if cfg.apiBase == "" {
			return nil, fmt.Errorf("FC_EMBEDDING_API_BASE (or FC_LLM_API_BASE) must be set when using FC_EMBEDDING_API_TYPE='ollama'")
		}
		var ollamaLLM *ollama.LLM
		ollamaLLM, err = ollama.New(
			ollama.WithServerURL(cfg.apiBase),
			ollama.WithModel(cfg.apiModel),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create Ollama embedding client: %w", err)
		}
		client = ollamaLLM

	case "gemini":
		// Gemini 通过 OpenAI 兼容接口调用（Google 提供了 OpenAI 兼容的 Embedding 端点）
		// 用户需要将 FC_EMBEDDING_API_BASE 设置为 Gemini 的 OpenAI 兼容端点
		logrus.Debug("Creating Gemini embedding client (via OpenAI-compatible API)")
		opts := []openai.Option{
			openai.WithToken(cfg.apiKey),
			openai.WithEmbeddingModel(cfg.apiModel),
		}
		if cfg.apiBase != "" {
			opts = append(opts, openai.WithBaseURL(cfg.apiBase))
		}
		var openaiLLM *openai.LLM
		openaiLLM, err = openai.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini embedding client: %w", err)
		}
		client = openaiLLM

	default: // "openai" 及其他 OpenAI 兼容服务
		logrus.Debug("Creating OpenAI embedding client")
		opts := []openai.Option{
			openai.WithToken(cfg.apiKey),
			openai.WithEmbeddingModel(cfg.apiModel),
		}
		if cfg.apiBase != "" {
			opts = append(opts, openai.WithBaseURL(cfg.apiBase))
		}
		var openaiLLM *openai.LLM
		openaiLLM, err = openai.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI embedding client: %w", err)
		}
		client = openaiLLM
	}

	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	embeddingClients.Store(cacheKey, embedder)
	return embedder, nil
}

// EmbedTexts 统一的 Embedding 接口，将文本列表编码为向量
// instruction 参数：对于支持 instruction 的模型会拼接到文本前面，不支持的静默忽略
// 返回 [][]float64 以便后续余弦相似度计算
func EmbedTexts(ctx context.Context, texts []string, instruction string) ([][]float64, error) {
	cfg := loadEmbeddingConfig()

	// 如果有 instruction 且非空，拼接到每条文本前面
	// 这是一种通用的 instruction 传递方式，适用于大多数模型
	globalInstruction := cfg.instruction
	if instruction != "" {
		globalInstruction = instruction
	}

	processedTexts := texts
	if globalInstruction != "" {
		processedTexts = make([]string, len(texts))
		for i, t := range texts {
			processedTexts[i] = globalInstruction + ": " + t
		}
	}

	embedder, err := getOrCreateEmbedder(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedder: %w", err)
	}

	// 带重试的 Embedding 调用
	var result32 [][]float32
	result32, err = retry.DoWithData(
		func() ([][]float32, error) {
			callCtx, cancel := context.WithTimeout(ctx, embeddingCallTimeout)
			defer cancel()
			return embedder.EmbedDocuments(callCtx, processedTexts)
		},
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(1*time.Second),
		retry.MaxDelay(5*time.Second),
		retry.OnRetry(func(n uint, err error) {
			logrus.Warnf("Retrying embedding call (attempt %d), err: %v", n+1, err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("embedding call failed after retries: %w", err)
	}

	// 将 float32 转换为 float64（余弦相似度计算使用 float64 精度更高）
	result := make([][]float64, len(result32))
	for i, vec32 := range result32 {
		vec64 := make([]float64, len(vec32))
		for j, v := range vec32 {
			vec64[j] = float64(v)
		}
		result[i] = vec64
	}

	return result, nil
}

// GetOrComputeAnchorVectors 获取或计算锚点向量（带内存缓存）
// 缓存键 = MD5(锚点文本) + "|" + 模型名称
// 惰性加载：系统重启后首次使用时重新计算
func GetOrComputeAnchorVectors(ctx context.Context, anchors []string, instruction string) ([][]float64, error) {
	if len(anchors) == 0 {
		return nil, nil
	}

	cfg := loadEmbeddingConfig()
	modelName := cfg.apiModel

	// 检查哪些锚点已缓存，哪些需要计算
	result := make([][]float64, len(anchors))
	var uncachedIndices []int
	var uncachedTexts []string

	for i, anchor := range anchors {
		cacheKey := fmt.Sprintf("anchor|%s|%s", util.GetMD5Hash(anchor), modelName)
		if cached, ok := anchorVectorCache.Load(cacheKey); ok {
			result[i] = cached.([]float64)
		} else {
			uncachedIndices = append(uncachedIndices, i)
			uncachedTexts = append(uncachedTexts, anchor)
		}
	}

	// 如果所有锚点都已缓存，直接返回
	if len(uncachedTexts) == 0 {
		logrus.Debugf("All %d anchor vectors found in cache", len(anchors))
		return result, nil
	}

	logrus.Infof("Computing embeddings for %d uncached anchors (out of %d total)", len(uncachedTexts), len(anchors))

	// 批量计算未缓存的锚点向量
	newVectors, err := EmbedTexts(ctx, uncachedTexts, instruction)
	if err != nil {
		return nil, fmt.Errorf("failed to compute anchor vectors: %w", err)
	}

	// 将新计算的向量写入缓存和结果
	for j, idx := range uncachedIndices {
		cacheKey := fmt.Sprintf("anchor|%s|%s", util.GetMD5Hash(anchors[idx]), modelName)
		anchorVectorCache.Store(cacheKey, newVectors[j])
		result[idx] = newVectors[j]
	}

	logrus.Infof("Anchor vector computation complete. Cached %d new vectors.", len(uncachedTexts))
	return result, nil
}
