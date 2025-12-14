这份文档旨在指导 **FeedCraft** 从单一的 RSS 中间件演进为通用的 Feed 生成平台。

---

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
type SourceFactory func(config json.RawMessage) (Source, error)
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
    fetcher Fetcher
    parser  Parser
}

func (p *PipelineSource) Generate(ctx context.Context) (*gofeed.Feed, error) {
    // 1. 获取原始数据
    raw, err := p.fetcher.Fetch(ctx)
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    // 2. 解析为 Feed
    feed, err := p.parser.Parse(raw)
    if err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }
    return feed, nil
}
```

### 4.4 配置结构 (Configuration)

数据库存储结构设计，利用 JSON 的灵活性。

```go
// 对应数据库中的 source_config 字段
type SourceConfig struct {
    Type string          `json:"type"` // e.g., "rss", "html", "json", "search"
  
    // 各个场景的专用配置，按需解析
    RSS    *RSSOptions    `json:"rss,omitempty"`
    HTML   *HTMLOptions   `json:"html,omitempty"`
    JSON   *JSONOptions   `json:"json,omitempty"`
    Search *SearchOptions `json:"search,omitempty"`
}

// 示例：HTML 场景配置
type HTMLOptions struct {
    URL       string            `json:"url"`
    Selectors map[string]string `json:"selectors"` // map: item, title, link
}

// 示例：JSON 场景配置
type JSONOptions struct {
    CurlCmd string            `json:"curl_cmd"` // 完整的 curl 命令字符串
    JQ      map[string]string `json:"jq"`       // map: iterator, title, link
}
```

---

## 5. 开发阶段规划 (Phased Roadmap)

### Phase 1: 基础设施重构 (Foundation)
**目标**：建立新的接口体系，并迁移现有的 RSS 功能，确保系统不退化。

1.  **数据库迁移**：
    - 修改 `Recipe` 表，废弃 `input_url`，新增 `source_type` (varchar) 和 `source_config` (json/text)。
    - 编写迁移脚本：将旧 `input_url` 转换为 `type: "rss", config: { "url": "..." }`。
2.  **核心接口实现**：
    - 定义上述 `Source`, `Fetcher`, `Parser` 接口。
    - 实现 `SourceRegistry` (全局注册中心)。
3.  **实现 RSS Source**：
    - 创建 `RSSFactory`。
    - *注意*：RSS 源可以直接实现 `Source` 接口（调用 `gofeed.ParseURL`），也可以拆分为 `HttpFetcher` + `RSSParser`。为了架构统一，建议采用拆分模式。
4.  **对接 Controller**：
    - 修改业务层代码，从“读取 URL”变为“通过 Factory 创建 Source 并调用 Generate”。

### Phase 2: 通用能力构建 (Generic Capabilities)
**目标**：实现 `PipelineSource` 模式，并落地 HTML 和 JSON 场景。

1.  **实现通用 Fetcher**：
    - `SimpleHttpFetcher`: 封装 `http.Get`。
    - `CurlFetcher`: 解析 curl 字符串（推荐使用第三方库如 `mattn/go-shellwords` 解析参数），构造 `http.Request`。
2.  **实现通用 Parser**：
    - `HtmlSelectorParser`: 集成 `goquery`，根据 CSS 选择器提取数据。
    - `JsonPathParser`: 集成 `GJSON`，根据路径提取数据。
3.  **注册 Factory**：
    - 注册 `html` 类型：组合 `SimpleHttpFetcher` + `HtmlSelectorParser`。
    - 注册 `json` 类型：组合 `CurlFetcher` + `JsonPathParser`。
4.  **单元测试**：重点测试 Fetcher 和 Parser 的组合逻辑。

### Phase 3: 高级场景与 AI (Advanced & AI)
**目标**：落地 Search/LLM 场景，利用 AI 能力生成 Feed。

1.  **Search Fetcher**：
    - 实现 `SearchApiFetcher`，对接 Google/Bing/Serper API。
2.  **Search Parser**：
    - 实现 `SearchResultParser`，将特定的搜索结果 JSON 映射为 Feed Item。
3.  **LLM 扩展 (可选)**：
    - 如果需要“自然语言生成 Feed 配置”，可以在这一层实现一个辅助服务：User Prompt -> LLM -> JSON Config (Selectors/JQ)。

### Phase 4: 前端适配 (Frontend Adaptation)
**目标**：让用户能方便地配置这些复杂的参数。

1.  **动态表单**：
    - 根据选择的 Source Type，动态渲染配置表单（HTML 显示选择器输入框，JSON 显示 cURL 输入框）。
2.  **可视化辅助工具**：
    - 开发简单的“可视化选择器”或“JSON 预览器”，帮助用户生成 Config JSON。

---

## 6. 关键注意事项 (Critical Notes)

1.  **错误处理**：Fetcher 的网络错误（超时、404）与 Parser 的解析错误（格式不对、字段缺失）需要有明确的区分，以便前端展示。
2.  **安全性**：
    - `CurlFetcher` 必须禁止访问内网地址（SSRF 防护）。
    - 限制 Fetcher 的最大响应体积，防止内存溢出。
3.  **兼容性**：确保 Phase 1 的迁移完全兼容旧版数据，用户无感知。
4.  **扩展性预留**：Parser 接口设计时，除了 `Parse([]byte)`，未来可能需要传入 `BaseURL` 用于处理相对路径链接（HTML 场景常见问题）。建议在 Parser 结构体初始化时注入 Context 或 Config。