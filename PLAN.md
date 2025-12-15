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

1.  **数据库迁移（蓝绿部署策略）**：
    - **创建新表**: 不修改现有 `custom_recipes` 表。而是创建一张新的 `custom_recipes_v2` 表，包含 `SourceType` (varchar) 和 `SourceConfig` (json/text) 等新字段。
    - **编写迁移脚本**: 编写脚本从旧的 `custom_recipes` 表中读取 `FeedURL` 数据，将其转换为 `type: "rss", config: { "url": "..." }` 的 `SourceConfig` 格式，并插入到新的 `custom_recipes_v2` 表中。此过程是非破坏性的，不触碰原表。
2.  **核心接口实现**：
    - 定义上述 `Source`, `Fetcher`, `Parser` 接口。
    - 实现 `SourceRegistry` (全局注册中心)。
3.  **实现 RSS Source**：
    - 创建 `RSSFactory`。
    - *注意*：RSS 源可以直接实现 `Source` 接口（调用 `gofeed.ParseURL`），也可以拆分为 `HttpFetcher` + `RSSParser`。为了架构统一，建议采用拆分模式。
4.  **对接 Controller**：
    - 修改业务层代码，从“读取 URL”变为“通过 Factory 创建 Source 并调用 Generate”。
    - **更新 Controller 逻辑**: `internal/recipe/custom_recipe.go` 中的 `CustomRecipe` 处理函数将不再使用旧的 `FeedURL` 或进行自循环 HTTP 请求。它将直接从 `custom_recipes_v2` 表中获取 `CustomRecipeV2` 记录，通过 `SourceFactory` 创建 `Source` 实例并调用 `Generate()` 方法获取原始 Feed，然后将该 Feed 直接传递给 `craft` 包的 `ProcessFeed` 函数进行处理。`cmd/main.go` 中的预热调度器也将更新以使用新的 V2 逻辑。

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

---

## 7. 具体实现方案 (Detailed Implementation Plan)

本章节将 `PLAN` 中定义的架构蓝图转化为具体的文件、代码和修改点。

### 7.1 新包和目录结构

我们将创建 `internal/source` 包来存放所有数据源生成相关的逻辑。

```
internal/
├───source/
│   ├─── source.go         # Source, SourceFactory 核心接口定义
│   ├─── registry.go       # 全局 Source Registry 实现
│   ├─── pipeline.go       # PipelineSource 通用实现
│   ├─── types.go          # 所有配置结构体 (SourceConfig, HTMLOptions etc.)
│   │
│   ├─── rss.go            # RSS Source 的 Factory 和实现
│   ├─── html.go           # HTML Source 的 Factory
│   ├─── json.go           # JSON Source 的 Factory
│   │
│   ├─── fetcher/
│   │    ├── fetcher.go     # Fetcher 接口定义
│   │    └── http_fetcher.go  # 基于 http.Get 的简单实现
│   │    └── curl_fetcher.go  # 基于 curl 命令的实现
│   │
│   └─── parser/
│        ├── parser.go      # Parser 接口定义
│        ├── rss_parser.go    # 将 XML/RSS byte 解析为 Feed
│        ├── html_parser.go   # 基于 goquery 的 HTML 解析器
│        └── json_parser.go   # 基于 gjson 的 JSON 解析器
│
└───dao/
    └─── recipe.go         # [修改]
    └─── migrate.go        # [修改]
└───controller/
    └─── craft_flow.go     # [修改]
```

### 7.2 核心定义实现 (Code Definitions)

**File: `internal/source/types.go`**
```go
package source

import (
    "encoding/json"
    "gorm.io/gorm"
)

// SourceConfig 是 Recipe.SourceConfig 字段的 Go 结构。
type SourceConfig struct {
    Type   string          `json:"type"`            // "rss", "html", "json"
    Config json.RawMessage `json:"config"`          // 特定类型的配置
}

// RSSOptions 是 Type="rss" 时的配置。
type RSSOptions struct {
    URL string `json:"url"`
}

// HTMLOptions 是 Type="html" 时的配置。
type HTMLOptions struct {
    URL          string            `json:"url"`
    ItemSelector string            `json:"item_selector"` // 列表项选择器
    Title        string            `json:"title"`         // 标题选择器
    Link         string            `json:"link"`          // 链接选择器 (可选, 默认为 item 的 href)
    Description  string            `json:"description"`   // 描述选择器 (可选)
    // ... 其他字段选择器
}

// JSONOptions 是 Type="json" 时的配置。
type JSONOptions struct {
    URL            string            `json:"url"`      // 请求 URL
    CurlCmd        string            `json:"curl_cmd"` // 或完整的 curl 命令
    ItemsIterator  string            `json:"items_iterator"` // 列表项的 gjson-path
    Title          string            `json:"title"`          // 标题的 gjson-path
    Link           string            `json:"link"`           // 链接的 gjson-path
    Description    string            `json:"description"`    // 描述的 gjson-path
    // ... 其他字段
}

// SearchOptions ...
```

**File: `internal/source/source.go`, `fetcher/fetcher.go`, `parser/parser.go`**
- 这些文件将包含 `PLAN` 第 4 节中定义的 `Source`, `SourceFactory`, `Fetcher`, `Parser` 接口的 Go 代码。

**File: `internal/source/registry.go`**
```go
package source

import (
    "fmt"
    "sync"
)

var (
    registry = make(map[string]SourceFactory)
    lock     = new(sync.RWMutex)
)

// Register a new source factory. Panics if type is already registered.
func Register(sourceType string, factory SourceFactory) {
    lock.Lock()
    defer lock.Unlock()
    if _, exists := registry[sourceType]; exists {
        panic(fmt.Sprintf("source factory for type '%s' already registered", sourceType))
    }
    registry[sourceType] = factory
}

// Get a source factory by type.
func Get(sourceType string) (SourceFactory, error) {
    lock.RLock()
    defer lock.RUnlock()
    factory, ok := registry[sourceType]
    if !ok {
        return nil, fmt.Errorf("no source factory registered for type '%s'", sourceType)
    }
    return factory, nil
}
```

### 7.3 Phase 1: 基础设施重构 (具体步骤)

1.  **数据库模型修改（蓝绿部署策略）**:
    **File: `internal/dao/recipe.go`**
    - `CustomRecipe` 结构体保持不变，对应 `custom_recipes` 表。
    - 新增 `CustomRecipeV2` 结构体，包含 `SourceType` 和 `SourceConfig` 字段，并使用 `TableName()` 方法将其映射到 `custom_recipes_v2` 表。

    ```go
    // CustomRecipe represents the original recipe structure.
    type CustomRecipe struct {
        BaseModelWithoutPK
        ID          string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
        Description string `json:"description,omitempty"`
        Craft       string `json:"craft" binding:"required"`
        FeedURL     string `json:"feed_url" binding:"required"`
    }

    // CustomRecipeV2 represents the new, refactored recipe structure.
    type CustomRecipeV2 struct {
        BaseModelWithoutPK
        ID           string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
        Description  string `json:"description,omitempty"`
        Craft        string `json:"craft" binding:"required"`
        SourceType   string `gorm:"type:varchar(50);not null;default:'rss'"`
        SourceConfig string `gorm:"type:text"`
    }

    func (CustomRecipeV2) TableName() string {
        return "custom_recipes_v2"
    }

    // 新增 GetCustomRecipeByIDV2 等 DAO 函数来操作 CustomRecipeV2
    func GetCustomRecipeByIDV2(db *gorm.DB, id string) (*CustomRecipeV2, error) { /* ... */ }
    // ...
    ```

2.  **编写迁移逻辑**:
    **File: `internal/dao/migrate.go`**
    - `MigrateDatabases` 函数将 `AutoMigrate` 新的 `CustomRecipeV2` 表。
    - 新增 `migrateRecipesToV2` 函数，该函数将从 `custom_recipes` 表读取所有 `CustomRecipe` 记录。
    - 对于每条记录，将其 `FeedURL` 转换为 `SourceConfig` 格式，并构建 `CustomRecipeV2` 记录。
    - 将新构建的 `CustomRecipeV2` 记录插入到 `custom_recipes_v2` 表中，并在插入前检查记录是否存在，以避免重复迁移。原有的修改和删除列的逻辑已被移除。

3.  **实现 RSS Source**：
    **File: `internal/source/fetcher/http_fetcher.go`** -> `HttpFetcher` 实现。
    **File: `internal/source/parser/rss_parser.go`** -> `RssParser` 实现，内部使用 `gofeed.NewParser().Parse(bytes.NewReader(data))`。
    **File: `internal/source/rss.go`**

    ```go
    package source

    import (
        "encoding/json"
        // ... 其他 imports
    )

    func init() {
        Register("rss", rssSourceFactory)
    }

    func rssSourceFactory(config json.RawMessage) (Source, error) {
        var opts RSSOptions
        if err := json.Unmarshal(config, &opts); err != nil {
            return nil, fmt.Errorf("invalid rss config: %w", err)
        }

        return &PipelineSource{
            Fetcher: &fetcher.HttpFetcher{URL: opts.URL},
            Parser:  &parser.RssParser{},
        }, nil
    }
    ```

4.  **对接 Controller**:
    **File: `internal/recipe/custom_recipe.go`**
    - `CustomRecipe` 处理函数将修改为通过 `dao.GetCustomRecipeByIDV2` 从新的 `custom_recipes_v2` 表中获取 `CustomRecipeV2` 记录。
    - 然后，它将解析 `SourceConfig`，通过 `SourceFactory` 创建 `Source` 实例，并调用 `Generate()` 方法生成原始 `*gofeed.Feed`。
    - 接着，该原始 `*gofeed.Feed` 将直接传递给 `craft.ProcessFeed` 函数进行加工。
    - 原有通过自循环 HTTP 请求来触发 CraftFlow 的低效机制已被完全移除。
    - `cmd/main.go` 中的预热调度器也将更新以使用新的 V2 逻辑。

### 7.4 Phase 2: 通用能力构建 (具体步骤)

1.  **实现 `HtmlSelectorParser`**:
    **File: `internal/source/parser/html_parser.go`**

    ```go
    package parser
    import "github.com/PuerkitoBio/goquery"
    
    type HtmlParser struct {
        BaseURL string // 用于补全相对链接
        Opts    source.HTMLOptions
    }

    func (p *HtmlParser) Parse(data []byte) (*gofeed.Feed, error) {
        doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
        if err != nil { /* ... */ }

        feed := &gofeed.Feed{Title: doc.Find("title").First().Text()}
        var items []*gofeed.Item

        doc.Find(p.Opts.ItemSelector).Each(func(i int, s *goquery.Selection) {
            title := s.Find(p.Opts.Title).Text()
            link, _ := s.Find(p.Opts.Link).Attr("href")
            // ... 处理相对链接: url.Parse(p.BaseURL).ResolveReference(link)
            
            items = append(items, &gofeed.Item{Title: title, Link: link, ...})
        })

        feed.Items = items
        return feed, nil
    }
    ```

2.  **实现 `JsonPathParser`**:
    **File: `internal/source/parser/json_parser.go`**
    - 推荐使用 `github.com/tidwall/gjson`。
    - 解析逻辑与 `HtmlParser` 类似，但使用 `gjson.Get(jsonData, p.Opts.ItemsIterator)` 来遍历 `items`，并用 `item.Get(p.Opts.Title)` 来提取字段。

3.  **注册新的 Factory**:
    **File: `internal/source/html.go`** & **`internal/source/json.go`**
    - 仿照 `rss.go`，创建 `htmlSourceFactory` 和 `jsonSourceFactory`。
    - `htmlSourceFactory` 将组合 `fetcher.HttpFetcher` 和 `parser.HtmlParser`。
    - `jsonSourceFactory` 将组合 `fetcher.CurlFetcher` (或 `HttpFetcher`) 和 `parser.JsonPathParser`。
    - 在各自的 `init()` 函数中调用 `source.Register(...)`。

至此，一个清晰、可执行的开发路线图已经建立。后续的 Phase 3 和 4 可在此基础上继续展开。