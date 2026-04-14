# 实施计划：Embedding 零样本分类主题过滤器

- [ ] 1. 新增 Embedding 环境变量配置与读取
  - 在 `internal/adapter/common_llm.go` 同级目录新建 `embedding.go`，实现 Embedding 配置读取逻辑
  - 使用 `util.GetEnvClient()` 读取 `FC_EMBEDDING_API_TYPE`、`FC_EMBEDDING_API_BASE`、`FC_EMBEDDING_API_KEY`、`FC_EMBEDDING_API_MODEL`、`FC_EMBEDDING_INSTRUCTION` 环境变量
  - 实现回退逻辑：当 Embedding 专用变量未配置时，回退使用 `FC_LLM_API_BASE`、`FC_LLM_API_KEY` 作为默认值
  - 更新 `.env.example` 文件，添加所有新增 Embedding 环境变量的示例和说明
  - _需求：5.1、5.2、5.3、5.4、5.5_

- [ ] 2. 实现统一的 Embedding 适配器接口（支持 OpenAI / Gemini / Ollama）
  - 在 `internal/adapter/embedding.go` 中定义 `EmbedTexts(ctx context.Context, texts []string, instruction string) ([][]float64, error)` 统一接口函数
  - 根据 `FC_EMBEDDING_API_TYPE` 分别实现 OpenAI、Gemini、Ollama 三种 Embedding 客户端的创建和调用逻辑
  - OpenAI 类型：使用 `langchaingo/embeddings/openai` 包（项目已依赖 `langchaingo v0.1.14`），配置 `WithBaseURL`、`WithToken`、`WithModel` 等选项
  - Gemini 类型：使用 `langchaingo/embeddings/googleai` 包或直接调用 Gemini REST API，支持 `task_type` 参数
  - Ollama 类型：使用 `langchaingo/embeddings/ollama` 包，配置 `WithServerURL`、`WithModel`
  - 实现 instruction 参数的传递：对支持 instruction 的模型传递该参数，不支持的静默忽略
  - 使用 `sync.Map` 缓存已创建的 Embedding 客户端实例（参考 `common_llm.go` 中 `llmClients` 的模式）
  - 使用 `retry-go/v4` 实现最多 3 次重试（指数退避），参考 `SimpleLLMCall` 中的重试模式
  - _需求：1.1、1.2、1.3、1.4、1.5、1.6、1.7、1.8、1.9_

- [ ] 3. 实现余弦相似度计算工具函数
  - 在 `internal/util/` 下新建 `vector.go`，实现 `CosineSimilarity(a, b []float64) float64` 函数
  - 处理边界情况：零向量、维度不匹配等，返回 0.0
  - 编写对应的单元测试 `vector_test.go`，覆盖正常情况和边界情况
  - _需求：3.3_

- [ ] 4. 实现锚点向量内存缓存管理
  - 在 `internal/adapter/embedding.go`（或新建 `embedding_cache.go`）中实现锚点向量缓存逻辑
  - 使用 `sync.Map` 作为内存缓存，缓存键为 `MD5(锚点文本) + 模型名称`
  - 实现 `GetOrComputeAnchorVectors(ctx, anchors []string, instruction string) ([][]float64, error)` 函数：先查缓存，未命中则调用 `EmbedTexts` 计算并缓存
  - 确保惰性加载策略：系统重启后首次使用时重新计算
  - 不依赖任何外部向量数据库
  - _需求：2.1、2.2、2.3、2.4、2.5、2.6_

- [ ] 5. 实现 Embedding 过滤器核心逻辑（CraftOption）
  - 在 `internal/craft/` 下新建 `embedding_filter.go`，实现 `OptionEmbeddingFilter(anchors []string, threshold float64, maxContentLen int, instruction string) CraftOption`
  - 实现文章文本拼接逻辑：将「标题 + 正文」拼接，正文超过 `maxContentLen` 时截取前 N 个字符
  - 使用 `parallel.Map` 并发计算所有文章的 Embedding 向量（参考 `OptionLLMFilterGeneric` 的并发模式）
  - 对每篇文章，与所有锚点向量逐一计算余弦相似度，任一锚点 ≥ 阈值即保留
  - 实现保守的错误处理策略：单篇文章 Embedding 失败时保留该文章并记录 Warn 日志；全部失败时返回原始 feed 并记录 Error 日志
  - 锚点文本为空时记录 Warn 日志并返回原始 feed
  - _需求：3.1、3.2、3.3、3.4、3.5、3.6、3.7、3.8、3.9、6.1、6.2、6.3、6.4_

- [ ] 6. 注册 CraftTemplate 并定义参数模板
  - 在 `internal/craft/embedding_filter.go` 中定义 `embeddingFilterParamTmpl`（`[]ParamTemplate`），包含 `anchors`、`threshold`、`max_content_length`、`instruction` 四个参数
  - 实现 `embeddingFilterLoadParam(m map[string]string) []CraftOption` 参数加载函数，解析参数并调用 `OptionEmbeddingFilter`
  - 处理参数校验：阈值格式错误时回退默认值 0.6 并记录 Warn 日志
  - 在 `internal/craft/entry.go` 的 `GetSysCraftTemplateDict()` 中注册 `embedding-filter` 模板
  - _需求：4.1、4.2、4.3、4.4、4.5、6.5_

- [ ] 7. 编写 Embedding 过滤器的单元测试
  - 在 `internal/craft/` 下新建 `embedding_filter_test.go`
  - 测试参数解析逻辑：默认值、阈值校验、锚点文本解析（换行分隔）
  - 测试过滤逻辑：mock Embedding 接口，验证相似度阈值过滤行为（保留/丢弃）
  - 测试错误处理：Embedding 失败时保留文章、锚点为空时不过滤
  - 在 `internal/adapter/` 下新建 `embedding_test.go`，测试 Embedding 适配器的配置回退逻辑和缓存命中逻辑
  - _需求：1.9、2.3、3.4、3.5、3.8、6.3、6.4、6.5_
