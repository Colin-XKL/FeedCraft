# FeedCraft 架构深化与原生迁移规划 (Architecture Refactor & Native Migration)

> 状态：进行中 (WIP)

## 1. 重构总目标

将现有的割裂的 `Source - Craft - Recipe` 架构，全面升级为以**“数据产出物 (CraftFeed)”**为核心的流式图谱架构 (Feed Graph)。
实现概念上的极致统一：底层任意节点都能被复用和嵌套，同时维持顶层产品对用户的心智清晰。

## 2. 核心数据模型 (The Water)

彻底消除内部流水线对第三方 `gofeed.Feed` 的强依赖，系统流转的唯一“血液”是自研模型：

- **`CraftFeed`**: 包含全局元数据及 Article 列表。
- **`CraftArticle`**: 对应原来的 Item/Entry，包含标题、正文，以及用于系统加工的追踪字段（如 `Depth`, `QualityScore`, `OriginalFeedID`）。

内部提供与 `gofeed.Feed`、`feeds.Feed` 的互相转换函数，仅在与第三方库交互或最终输出 RSS XML 时进行转换。

## 3. 核心抽象与运行时模型 (The Graph)

### 3.1 统一数据源接口 (FeedProvider)

任何能产出 `CraftFeed` 的节点（无论是底层 Source，还是顶层 Recipe/Topic），在运行时都只暴露统一接口：

```go
type FeedProvider interface {
    Fetch(ctx context.Context) (*model.CraftFeed, error)
}
```

### 3.2 加工流水线：V3 架构 (Functional Options)

在 Craft 层的数据加工流水线设计上，我们放弃了笨重的 `FeedProcessor` 接口模式（V2），直接采用 **V3 架构 (Functional Options + 原生模型)**，回归 Go 语言最地道的中间件风格：

```go
// 接收原生 feed，返回处理后的新 feed（在内部执行 clone 以防止并发图数据污染）
type CraftOption func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error)
```

原生的 `Limit`, `Translate`, `Summary` 等所有加工逻辑，不再定义空结构体，而是直接返回上述闭包函数。流水线通过函数遍历嵌套执行，实现极致的轻量化。

### 3.3 运行时图谱，存储/产品隔离 (Runtime Graph, Separate DB)

关于 Recipe（单线配方）和 Topic（多线聚合）的边界：

- **产品与存储层**：继续保持 `Recipe` 和 `Topic` 的表结构与概念隔离。这符合普通用户的心智模型。
- **引擎运行时层**：在 Go 引擎运行时，Topic 和 Recipe 将被统一的 Builder 编译成一模一样的 `FeedProvider` 节点。它们在执行时完全是一张可以任意嵌套组合的 Feed Graph。

## 4. 统一顶层输入模型与解析 (InputSpec & Router)

为了消解 `SourceConfig` 和 `InputURI` 之间的割裂，统一采用多态路由架构：

### 4.1 统一多态 JSON (`InputSpec`)

顶层数据输入统一使用 `InputSpec` 结构，避免将 `SourceConfig` 做成包含排斥字段的“胖模型”：

```go
type InputSpec struct {
    Type   string          // 例如 "uri", "html", "json", "search"
    Config json.RawMessage // 对应的独立配置结构体序列化数据
}
```

### 4.2 统一 URI 路由器 (URI Router)

当 `InputSpec.Type == "uri"` 时，系统在解析层采用**统一 URI 路由器模式**。
引擎注册不同的 Resolver，通过 URI 的 Scheme 进行动态分发：

- `feedcraft://recipe/:id` -> 路由到内部 Recipe 解析器
- `feedcraft://topic/:id` -> 路由到内部 Topic 解析器
- `http(s)://...` -> **作为语法糖 (Syntactic Sugar)**，为了方便用户快速使用，它被路由到第三方网站解析器，在引擎内部会自动展开（Desugar）为一个标准的、带预设策略的 `InputSpec` 结构。

## 5. 重构实施步骤 (Implementation Steps)

1.  **Phase 1: 模型基建 (已部分完成)**
    - 创建 `CraftFeed` 和 `CraftArticle`，并提供适配函数。
2.  **Phase 2: 顶层输入模型与 Router 改造 (Next)**
    - 引入 `InputSpec` 和 `URI Router`。
    - 构建统一的 `InputSpec -> FeedProvider` builder，让 Topic 和 Recipe 共用同一套输入解析逻辑。
3.  **Phase 3: 引擎与 Craft 原生化 (V3)**
    - 把 `ProcessRecipeByID` 从旧的零散编排迁到新的 Builder runtime。
    - 将 `internal/craft` 逐步重构为 V3 的 `CraftOption` 闭包架构，移除旧的 V2 适配层。
    - 让底层 Source 直接返回 `*model.CraftFeed`。
4.  **Phase 4: 图谱收官**
    - 全面支撑 Topic 嵌套 Recipe、Topic 嵌套 Topic 的运行时逻辑。
    - 清理旧的遗留代码，完成整体平滑迁移。
