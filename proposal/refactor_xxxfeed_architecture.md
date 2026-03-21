# FeedCraft 底层架构重构规划 (Code Refactoring Plan)

## 1. 重构目标

将现有的 `Source - Craft - Recipe` 架构，全面升级为以**“数据产出物 (Feed)”**为核心的流式架构，并引入 `TopicFeed` 聚合能力。
核心命名体系将全面向 `xxxFeed` 靠拢，实现概念上的极致统一。

## 2. 核心模型与接口定义 (internal/model & engine)

### 2.1 统一数据载体 (The Water)
消除对第三方 `gofeed.Feed` 的强依赖，在内部使用自定义模型：
*   **`CraftFeed`**: 包含全局元数据及 Article 列表。
*   **`CraftArticle`**: 对应原来的 Item/Entry，包含标题、正文，以及新增的追踪字段（如 `Depth`, `QualityScore`, `OriginalFeedID`）。

### 2.2 核心接口定义 (The Interfaces)
在 `internal/engine` 中定义两个极简接口：

1.  **FeedProvider**: 任何能产出 `CraftFeed` 的节点。
    ```go
    type FeedProvider interface {
        // 使用 Fetch 替代 GetFeed，避免 RecipeFeed.GetFeed() 的冗余感
        Fetch(ctx context.Context) (*model.CraftFeed, error)
    }
    ```
2.  **FeedProcessor**: 任何能加工 `CraftFeed` 的节点。
    ```go
    type FeedProcessor interface {
        Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error)
    }
    ```

## 3. 具体业务对象的实现映射

### 3.1 RawFeed (原 Source 层)
*   **重构动作**：将现有的 `HtmlSource`, `SearchSource`, `RssSource` 等模块的返回值，从 `gofeed.Feed` 改为 `model.CraftFeed`。
*   **定位**：它们实现了 `FeedProvider` 接口，是流水线的起点。

### 3.2 Crafts (加工层保持原样，仅改签名)
*   **重构动作**：将所有的 `AtomCraft`（如 translate, summary）的接口签名改为接收和返回 `model.CraftFeed`。
*   **FlowCraft**：本质上是一个包含 `[]FeedProcessor` 的结构体。

### 3.3 RecipeFeed (单线配方源)
*   **重构动作**：重构引擎层，将数据库中的 `Recipe` 模型在运行时构建为一个实现 `FeedProvider` 接口的对象。
*   **执行逻辑**：
    ```go
    type RecipeFeed struct {
        Base  FeedProvider     // 例如 RssSource (RawFeed)
        Craft FeedProcessor    // 例如 FlowCraft
    }
    func (r *RecipeFeed) Fetch(ctx) (*CraftFeed, error) {
        rawFeed, _ := r.Base.Fetch(ctx)
        return r.Craft.Process(ctx, rawFeed)
    }
    ```

### 3.4 TopicFeed (多线聚合源)
*   **重构动作**：新增 `TopicFeed` 结构，它同样实现 `FeedProvider` 接口，支持嵌套组合。
*   **执行逻辑**：
    ```go
    type TopicFeed struct {
        Inputs     []FeedProvider  // 可以是 RecipeFeed，也可以是其他 TopicFeed
        Aggregator FeedProcessor   // 负责去重、排序、截断的特殊处理器
    }
    func (t *TopicFeed) Fetch(ctx) (*CraftFeed, error) {
        // 1. 并发调用所有的 Inputs.Fetch() 获取多个子 Feed
        // 2. 将数据合并为一个大的临时 CraftFeed
        // 3. 调用 t.Aggregator.Process() 进行过滤和截断
        // 4. 返回最终的 CraftFeed
    }
    ```

## 4. 重构实施步骤 (Implementation Steps)

1.  **Phase 1: 模型基建**
    *   在 `internal/model` 创建 `CraftFeed` 和 `CraftArticle`。
    *   提供 `ToGofeed()` 和 `FromGofeed()` 转换函数，保证现有前端阅读器和第三方库的兼容。
2.  **Phase 2: 接口改造**
    *   修改所有 `Source` 抓取器的返回值。
    *   修改所有 `Craft` 插件的入参和出参。
    *   *此时，编译会报大量错误，需要逐个修复引擎层的调用。*
3.  **Phase 3: 引擎升级 (RecipeFeed)**
    *   重写现有的执行流，使其符合 `Provider -> Processor` 的新接口模式。测试所有的基础单线流是否正常工作。
4.  **Phase 4: 引入 TopicFeed**
    *   实现并发抓取逻辑 (使用 `errgroup`)。
    *   实现 `Aggregator` 处理器（包含基础的 URL 去重和条数 Limit）。
    *   暴露新的路由 `GET /topic/:id`。

## 5. 遗留问题与风险点
*   在 Phase 2 和 Phase 3 期间，整个应用可能处于不可用状态。建议在单独的 `feature/refactor-lego-flow` 分支上进行开发。
*   缓存层的适配：原来的 Recipe 是直接缓存 XML，现在可能需要考虑是缓存最终的 XML 还是中间的 `CraftFeed` 对象（建议依然缓存最终输出的 XML 文本，保持简单）。