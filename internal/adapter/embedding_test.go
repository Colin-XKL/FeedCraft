package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- loadEmbeddingConfig 测试 ---

func TestLoadEmbeddingConfig_Defaults(t *testing.T) {
	cfg, err := loadEmbeddingConfig()
	assert.NoError(t, err)
	// 默认 apiType 应为 "openai"
	assert.Equal(t, "openai", cfg.apiType)
	// 默认模型应为 defaultEmbeddingModel
	assert.Equal(t, defaultEmbeddingModel, cfg.apiModel)
}

func TestLoadEmbeddingConfig_ApiTypeDefault(t *testing.T) {
	cfg, err := loadEmbeddingConfig()
	assert.NoError(t, err)
	// 未设置 FC_EMBEDDING_API_TYPE 时应默认为 "openai"
	assert.Equal(t, "openai", cfg.apiType)
}

func TestLoadEmbeddingConfig_ReturnsNoFatal(t *testing.T) {
	// 确保 loadEmbeddingConfig 不会 panic 或 fatal
	cfg, err := loadEmbeddingConfig()
	assert.NoError(t, err)
	assert.NotEmpty(t, cfg.apiType)
}

// --- getOrCreateEmbedder 测试 ---

func TestGetOrCreateEmbedder_OllamaWithoutBase(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "ollama",
		apiBase:  "",
		apiKey:   "",
		apiModel: "nomic-embed-text",
	}
	_, err := getOrCreateEmbedder(cfg)
	assert.Error(t, err, "ollama without apiBase should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_BASE")
}

func TestGetOrCreateEmbedder_OpenAIDefault(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "openai",
		apiBase:  "",
		apiKey:   "test-key",
		apiModel: "text-embedding-3-small",
	}
	embedder, err := getOrCreateEmbedder(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, embedder)
}

func TestGetOrCreateEmbedder_GeminiType(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "gemini",
		apiBase:  "https://generativelanguage.googleapis.com/v1beta/openai/",
		apiKey:   "test-key",
		apiModel: "gemini-embedding-001",
	}
	embedder, err := getOrCreateEmbedder(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, embedder)
}

func TestGetOrCreateEmbedder_ClientCaching(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "openai",
		apiBase:  "https://test.example.com/v1",
		apiKey:   "cache-test-key",
		apiModel: "text-embedding-3-small",
	}

	// 第一次创建
	embedder1, err := getOrCreateEmbedder(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, embedder1)

	// 第二次应该从缓存获取（同一个实例）
	embedder2, err := getOrCreateEmbedder(cfg)
	assert.NoError(t, err)
	assert.Equal(t, embedder1, embedder2, "should return cached embedder instance")
}

// --- API Key 校验测试 ---

func TestGetOrCreateEmbedder_OpenAIEmptyApiKey(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "openai",
		apiBase:  "https://api.openai.com/v1",
		apiKey:   "",
		apiModel: "text-embedding-3-small",
	}
	_, err := getOrCreateEmbedder(cfg)
	assert.Error(t, err, "openai with empty apiKey should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_KEY")
}

func TestGetOrCreateEmbedder_GeminiEmptyApiKey(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "gemini",
		apiBase:  "https://generativelanguage.googleapis.com/v1beta/openai/",
		apiKey:   "",
		apiModel: "gemini-embedding-001",
	}
	_, err := getOrCreateEmbedder(cfg)
	assert.Error(t, err, "gemini with empty apiKey should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_KEY")
}

func TestGetOrCreateEmbedder_OllamaNoApiKeyRequired(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "ollama",
		apiBase:  "http://localhost:11434",
		apiKey:   "",
		apiModel: "bge-m3",
	}
	embedder, err := getOrCreateEmbedder(cfg)
	assert.NoError(t, err, "ollama should not require apiKey")
	assert.NotNil(t, embedder)
}
