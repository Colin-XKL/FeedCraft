# FeedCraft 2.0 架构演进设计文档

## 1. 背景与目标 (Context & Objectives)

**现状**：FeedCraft 目前作为一个“加工中间件”，强依赖于现有的 RSS URL 输入。
**痛点**：用户无法直接将 HTML 网页、JSON API 或搜索结果转换为 RSS，限制了使用场景。
**目标**：将核心架构升级为 **通用生成器 (Universal Generator)**。
- **输入层**：支持任意来源（RSS, HTML, JSON, Search Query）。
- **处理层**：保留并复用现有的 CraftAtom 强大加工能力。
- **设计原则**：保持架构的简洁性，通过组合模式实现高扩展性。

---

## 2. 核心设计思想 (Core Philosophy)

为了兼顾灵活性与工程规范，我们采用 **混合架构 (Hybrid Architecture)**：
**Registry Factory + Compositional Pipeline**。

1.  **正交解耦 (Orthogonality)**：将“获取数据 (Fetch)”与“解析数据 (Parse)”分离。
    - *Fetcher* 负责网络 IO，屏蔽 HTTP/cURL/API 差异。
    - *Parser* 负责数据清洗，屏蔽 HTML/JSON/Text 差异。
2.  **流水线模式 (Pipeline)**：一个标准的源生成过程 = `Fetcher` + `Parser`。
3.  **配置驱动 (Configuration Driven)**：所有源的行为完全由 JSON 配置定义，便于存储和前端生成。

---

## 3. 总体架构 (Architecture Overview)

系统数据流转将变更为：

```mermaid
graph LR
    Config[Source Config] --> Factory[Source Factory]
    Factory -->|Build| Source[Source Instance]
  
    subgraph "Source Execution (Pipeline)"
        Source --> Fetcher
        Fetcher -->|[]byte| Parser
        Parser -->|Standard Feed| Output
    end
  
    Output --> CraftFlow[CraftFlow (Middleware)]
    CraftFlow --> FinalRSS
```

---

## 4. 核心接口定义 (Core Definitions)

这是系统的骨架，必须严格遵守。

### 4.1 顶层抽象 (Source Layer)

```go
// Source 是所有生成器的顶层接口
// 它的职责非常单一：产出标准 Feed
type Source interface {
    Generate(ctx context.Context) (*gofeed.Feed, error)
}

// SourceFactory 负责根据配置创建 Source 实例
// 采用注册模式，每种 type 对应一个 Factory
type SourceFactory func(config *config.SourceConfig) (Source, error)
```

### 4.2 组件抽象 (Component Layer)

实现组合式复用的关键接口。

```go
// Fetcher 专注 IO，只管拿回二进制数据
type Fetcher interface {
    Fetch(ctx context.Context) ([]byte, error)
}

// Parser 专注逻辑，将二进制数据转化为 Feed 对象
type Parser interface {
    Parse(data []byte) (*gofeed.Feed, error)
}
```

### 4.3 通用流水线实现 (The Pipeline)

这是大多数场景（HTML/JSON/Search）通用的实现类，**无需重复编码**。

```go
type PipelineSource struct {
    Config  *config.SourceConfig
    Fetcher Fetcher
    Parser  Parser
}

func (p *PipelineSource) Generate(ctx context.Context) (*gofeed.Feed, error) {
    // 1. 获取原始数据
    raw, err := p.Fetcher.Fetch(ctx)
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    // 2. 解析为 Feed
    feed, err := p.Parser.Parse(raw)
    if err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }
    // 3. 应用元数据覆盖 (可选)
    p.applyFeedMetaOverrides(feed)
    
    return feed, nil
}
```

### 4.4 配置结构 (Configuration)

数据库存储结构设计，利用 JSON 的灵活性。**已更新为实际实现结构**。

```go
// SourceConfig 是 Recipe.SourceConfig 字段的 Go 结构。
// 它组合了 Fetcher、Parser 和 Metadata 的配置。
type SourceConfig struct {
    Type constant.SourceType `json:"type"` // e.g., "rss", "html", "json"

    // Feed 元数据覆盖
    FeedMeta *FeedMetaConfig `json:"feed_meta,omitempty"`

    // Fetcher 配置 - 互斥，每次只能有一个生效
    HttpFetcher *HttpFetcherConfig `json:"http_fetcher,omitempty"`
    // CurlFetcher *CurlFetcherConfig `json:"curl_fetcher,omitempty"` // 未来扩展

    // Parser 配置 - 互斥，每次只能有一个生效
    HtmlParser *HtmlParserConfig `json:"html_parser,omitempty"`
    JsonParser *JsonParserConfig `json:"json_parser,omitempty"`
}

type HttpFetcherConfig struct {
    URL     string            `json:"url"`
    Headers map[string]string `json:"headers,omitempty"`
}

type HtmlParserConfig struct {
    ItemSelector string `json:"item_selector"`
    Title        string `json:"title"`
    Link         string `json:"link"`
    Description  string `json:"description,omitempty"`
}
```

---

## 5. 开发阶段规划 (Phased Roadmap)

### Phase 1: 基础设施重构 (Foundation) [已完成]
**状态**：核心迁移已完成，数据库已升级，API 已重构。

1.  **数据库迁移（蓝绿部署策略）**：
    - **创建新表**: 已创建 `custom_recipes_v2` 表，包含 `SourceType` 和 `SourceConfig` 字段。
    - **数据迁移**: `migrateRecipesToV2` 脚本已实现。它读取旧 `custom_recipes` 表中的 `FeedURL`，转换为 `SourceConfig` (使用 `HttpFetcher`) 并插入新表。此过程是非破坏性的。
2.  **核心接口实现**：
    - `Source`, `Fetcher`, `Parser` 接口已定义。
    - `SourceRegistry` (全局注册中心) 已实现。
3.  **实现 RSS Source**：
    - `RSSFactory` 已注册。
    - 实现了 `HttpFetcher` + `RssParser` 的组合模式。
4.  **对接 Controller**：
    - `Create/Update/Get/List` 接口已全部迁移至使用 `CustomRecipeV2`。
    - 核心处理逻辑 `ProcessRecipeByID` 已更新为通过 Factory 创建 Source 并生成 Feed。

### Phase 2: 通用能力构建 (Generic Capabilities) [进行中]
**目标**：实现 `PipelineSource` 模式，并落地 HTML 和 JSON 场景。

1.  **实现通用 Fetcher**：
    - `HttpFetcher`: 已实现。
    - `CurlFetcher`: 待实现。
2.  **实现通用 Parser**：
    - `HtmlParser`: 已实现 (基于 `goquery`)。
    - `JsonPathParser`: 待实现 (基于 `GJSON`)。
3.  **注册 Factory**：
    - `html` 类型：已注册并实现。
    - `json` 类型：待实现。

### Phase 3: 高级场景与 AI (Advanced & AI) [规划中]
**目标**：落地 Search/LLM 场景，利用 AI 能力生成 Feed。

1.  **Search Fetcher**：
    - 实现 `SearchApiFetcher`，对接 Google/Bing/Serper API。
2.  **Search Parser**：
    - 实现 `SearchResultParser`。

### Phase 4: 前端适配 (Frontend Adaptation) [规划中]
**目标**：让用户能方便地配置这些复杂的参数。

1.  **动态表单**：
    - 前端已开始适配 V2 数据结构。
    - 需要完善 HTML 选择器配置界面和 JSON 配置界面。

---

## 6. 数据库迁移检查报告 (Migration Check Report)

经过对当前代码库的详细检查，确认迁移逻辑如下：

1.  **模型一致性**：`internal/dao/recipe.go` 中定义的 `CustomRecipeV2` 结构体正确映射了新表结构，包含 JSON 类型的 `source_config` 字段。
2.  **迁移安全性**：`internal/dao/migrate.go` 中的 `migrateRecipesToV2` 函数采用了“只读旧表，写入新表”的策略。
    - 它首先检查 `custom_recipes` (旧表) 是否存在。
    - 遍历旧数据时，会先检查 `custom_recipes_v2` (新表) 中是否已存在同 ID 记录，**确保了幂等性** (Idempotency)。
    - 数据转换逻辑正确：旧的 `FeedURL` 被正确封装进 `SourceConfig` 的 `HttpFetcher` 配置中，`SourceType` 被设为 `rss`。
3.  **业务逻辑切换**：`internal/controller/custom_recipe.go` 中的所有 CRUD 操作均已切换为调用 `dao` 包中的 `*V2` 函数。旧的 `CustomRecipe` 结构体仅作为历史遗留存在，不再被新的业务逻辑使用。
4.  **结论**：数据库迁移部分已正确且安全地实施。

---

## 7. 目录结构 (Directory Structure)

```
internal/
├───source/
│   ├─── source.go         # Source, SourceFactory 核心接口
│   ├─── registry.go       # 全局 Source Registry
│   ├─── pipeline.go       # PipelineSource 通用实现
│   ├─── html.go           # HTML Source Factory [New]
│   ├─── rss.go            # RSS Source Factory
│   │
│   ├─── fetcher/
│   │    ├── fetcher.go     # Fetcher 接口
│   │    └── http_fetcher.go  # HTTP Fetcher 实现
│   │
│   └─── parser/
│        ├── parser.go      # Parser 接口
│        ├── rss_parser.go    # RSS Parser 实现
│        └── html_parser.go   # HTML Parser 实现 [New]
│
├───config/
│   └─── source_config.go  # SourceConfig 定义 [Updated]
│
└───dao/
    ├─── recipe.go         # CustomRecipeV2 定义
    └─── migrate.go        # 迁移逻辑
```
