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

// --- Gemini apiBase 校验测试（需求 3）---

func TestGetOrCreateEmbedder_GeminiEmptyApiBase(t *testing.T) {
	cfg := embeddingConfig{
		apiType:  "gemini",
		apiBase:  "",
		apiKey:   "test-key",
		apiModel: "gemini-embedding-001",
	}
	_, err := getOrCreateEmbedder(cfg)
	assert.Error(t, err, "gemini with empty apiBase should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_BASE")
	assert.Contains(t, err.Error(), "gemini")
}

// --- 默认模型逻辑测试（需求 4）---

func TestLoadEmbeddingConfig_OllamaEmptyModel(t *testing.T) {
	// 设置环境变量模拟 ollama 类型但不设模型
	t.Setenv("FC_EMBEDDING_API_TYPE", "ollama")
	t.Setenv("FC_EMBEDDING_API_BASE", "http://localhost:11434")
	t.Setenv("FC_EMBEDDING_API_MODEL", "")
	_, err := loadEmbeddingConfig()
	assert.Error(t, err, "ollama with empty model should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_MODEL")
}

func TestLoadEmbeddingConfig_GeminiEmptyModel(t *testing.T) {
	t.Setenv("FC_EMBEDDING_API_TYPE", "gemini")
	t.Setenv("FC_EMBEDDING_API_BASE", "https://generativelanguage.googleapis.com/v1beta/openai/")
	t.Setenv("FC_EMBEDDING_API_KEY", "test-key")
	t.Setenv("FC_EMBEDDING_API_MODEL", "")
	_, err := loadEmbeddingConfig()
	assert.Error(t, err, "gemini with empty model should return error")
	assert.Contains(t, err.Error(), "FC_EMBEDDING_API_MODEL")
}

func TestLoadEmbeddingConfig_OpenAIDefaultModel(t *testing.T) {
	t.Setenv("FC_EMBEDDING_API_TYPE", "openai")
	t.Setenv("FC_EMBEDDING_API_MODEL", "")
	cfg, err := loadEmbeddingConfig()
	assert.NoError(t, err)
	assert.Equal(t, defaultEmbeddingModel, cfg.apiModel, "openai should use default model when not set")
}

// --- resolveInstruction 测试（需求 1）---

func TestResolveInstruction_ExplicitInstruction(t *testing.T) {
	cfg := embeddingConfig{instruction: "global_instruction"}
	result := resolveInstruction("explicit_instruction", cfg)
	assert.Equal(t, "explicit_instruction", result, "explicit instruction should take priority")
}

func TestResolveInstruction_FallbackToGlobal(t *testing.T) {
	cfg := embeddingConfig{instruction: "global_instruction"}
	result := resolveInstruction("", cfg)
	assert.Equal(t, "global_instruction", result, "should fallback to global instruction when empty")
}

func TestResolveInstruction_BothEmpty(t *testing.T) {
	cfg := embeddingConfig{instruction: ""}
	result := resolveInstruction("", cfg)
	assert.Equal(t, "", result, "should return empty when both are empty")
}
