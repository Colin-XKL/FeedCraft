# SourceConfig 配置模型规范 (方案 B2)

## 1. 定位与设计目标

本规范定义的是 FeedCraft 系统中 `Source`（数据源）的统一**配置模型**，用于解决多类型 Source 配置在序列化、持久化及代码解析上的复杂性问题。

需要明确：

- FeedCraft 运行时的顶层统一输入模型仍然是 `InputSpec`
- `SourceConfig` 不是整个系统的顶层输入抽象
- `SourceConfig` 主要用于承载 `InputSpec(kind=source)` 下的底层 source 配置

也就是说，在运行时装配层，推荐的主链路是：

```text
持久化结构 / 外部输入 -> InputSpec -> FeedProvider
```

其中当 `InputSpec.kind=source` 时，`InputSpec` 会携带一份 `SourceConfig`，再进入 Source 子系统的构建流程。

核心目标：

- **单一真值**：`Source` 子系统内的配置类型仅由 `SourceConfig.Type` 决定。
- **序列化友好**：不使用 `interface{}` 或自定义 JSON 反序列化，直接依赖 Go 原生 `encoding/json` 及 SQLite JSON 存储。
- **配置隔离**：不同类型的特有配置互不干扰，顶层结构保持精简。
- **组件复用**：底层的 Fetcher（获取）和 Parser（解析）配置模块可灵活组合。

## 2. 数据结构设计 (Typed Union Body)

采用“公共外壳 + 强类型聚合 Detail”的模式。

### 2.1 顶层配置外壳 (`SourceConfig`)

`SourceConfig` 是 Source 子系统内部的统一配置外壳，仅包含所有类型共用的基础元信息，以及一个包含具体类型特有配置的 `Detail` 结构体。

```go
type SourceConfig struct {
    Type     SourceType      `json:"type"`                // 唯一真值：标识 Source 类型 (rss/html/json/search)
    FeedMeta *FeedMetaConfig `json:"feed_meta,omitempty"` // 所有类型共用的 Feed 元信息配置

    // 特有配置收敛到一层，避免顶层无限变宽
    Detail   SourceDetail    `json:"detail,omitempty"`
}
```

### 2.2 类型特有配置聚合 (`SourceDetail`)

`SourceDetail` 是一个可选聚合体（Tagged Union 模拟）。在任意时刻，只有与 `Type` 对应的字段应为非空。

```go
type SourceDetail struct {
    RSS    *RSSSourceConfig    `json:"rss,omitempty"`
    HTML   *HTMLSourceConfig   `json:"html,omitempty"`
    JSON   *JSONSourceConfig   `json:"json,omitempty"`
    Search *SearchSourceConfig `json:"search,omitempty"`
}
```

### 2.3 底层组件组合 (以 JSON 和 Search 为例)

不同的 `SourceConfig` 变体通过组合基础的 Fetcher 和 Parser 组件来构建：

```go
// curl-to-rss 场景：需要动态存储 HTTP 获取和 JSON 解析规则
type JSONSourceConfig struct {
    HttpFetcher *HttpFetcherConfig `json:"http_fetcher,omitempty"`
    JsonParser  *JsonParserConfig  `json:"json_parser,omitempty"`
}

// search-to-rss 场景：只需要存储搜索规则，JSON 解析规则固定在代码中，不存入 DB
type SearchSourceConfig struct {
    SearchFetcher *SearchFetcherConfig `json:"search_fetcher,omitempty"`
}
```

## 3. JSON 表达与持久化示例

持久化到 SQLite 时，Go 会自动忽略 `omitempty` 为 nil 的字段，生成极其干净的 JSON 字符串。

### 场景 1：`curl-to-rss` (SourceJSON)

用户自定义接口抓取和 JSON 解析规则。

```json
{
  "type": "json",
  "feed_meta": {
    "title_override": "My Custom JSON Feed"
  },
  "detail": {
    "json": {
      "http_fetcher": {
        "url": "https://api.example.com/v1/posts"
      },
      "json_parser": {
        "items_path": "$.data.posts",
        "title_path": "$.title"
      }
    }
  }
}
```

### 场景 2：`search-to-rss` (SourceSearch)

用户配置搜索词，系统代码负责组装固定的 Parser 逻辑。

```json
{
  "type": "search",
  "detail": {
    "search": {
      "search_fetcher": {
        "query": "feedcraft update"
      }
    }
  }
}
```

## 4. 运行时分层与解析范式

在当前 runtime 分层中，推荐关系是：

```text
InputSpec(kind=source) -> SourceConfig -> Source / Provider builder
```

因此：

- `InputSpec` 是运行时顶层输入抽象
- `SourceConfig` 是 `kind=source` 分支下的底层配置载体
- `SourceConfig` 不应替代 `InputSpec` 成为 runtime 总入口

读取 DB 后，如果上层已经将输入归一为 `InputSpec(kind=source)`，则 Source 子系统内部可以继续利用 `SourceConfig.Type` 作为唯一分发依据，提取对应的特有配置传给底层组件：

```go
func BuildSource(cfg *SourceConfig) (Source, error) {
    switch cfg.Type {
    case SourceJSON:
        if cfg.Detail.JSON == nil {
            return nil, errors.New("json detail config is missing")
        }
        // 直接将反序列化好的 HttpFetcher 和 JsonParser 传入解析引擎
        return newJSONSource(cfg.Detail.JSON.HttpFetcher, cfg.Detail.JSON.JsonParser)

    case SourceSearch:
        if cfg.Detail.Search == nil {
            return nil, errors.New("search detail config is missing")
        }
        // Search 场景：提取 SearchFetcher 结果，在代码中 new 固定的 JsonParser 规则
        fixedJsonParser := &JsonParserConfig{ ItemsPath: "$.items" }
        return newSearchSource(cfg.Detail.Search.SearchFetcher, fixedJsonParser)

    // ... 处理其他类型
    default:
        return nil, fmt.Errorf("unsupported source type: %s", cfg.Type)
    }
}
```

这里的 `BuildSource(cfg)` 应理解为 Source 子系统内部的构建步骤，而不是整个 Feed runtime 的统一装配入口。

## 5. 方案优势总结

1. **类型安全与代码清晰**：纯 Struct 嵌套，无任何 `interface{}` 类型断言。
2. **极简持久化**：直接支持标准库的 `json.Marshal / Unmarshal`，无需编写自定义反序列化逻辑。
3. **灵活性**：解耦了数据获取 (`Fetcher`) 与数据解析 (`Parser`)。例如 `curl-to-rss` 可将 Parser 存入 DB，而 `search-to-rss` 可复用相同底层 Parser 但免于冗余存储。
4. **与 runtime 分层兼容**：保留 `SourceConfig` 的配置表达能力，同时明确 runtime 顶层统一入口仍然是 `InputSpec`。
