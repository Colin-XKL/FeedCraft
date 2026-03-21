# 乐高式数据流架构设计 (Lego-based Data Flow Architecture)

## 1. 核心设计理念 (Core Design Philosophy)

在当前的系统中，`第三方原始 RSS`、`用户自定义 Recipe`、`处理后的 Feed` 以及即将引入的 `Topic (主题聚合)`，往往被视作不同的业务实体。
为了实现**“极度灵活、方便代码复用、且未来能支持类似 n8n/ComfyUI 的可视化节点连线编排”**的终极目标，我们需要将系统架构进行降维与抽象。

我们将整个 FeedCraft 系统抽象为一个**“流式数据处理管线 (Data Pipeline)”**，核心只保留两层概念：
1. **The Water (流动的液体)**：系统中流转的统一数据载体。
2. **The Pipes (处理节点/管道)**：负责生成或加工数据的模块。

无论是原始 RSS 还是 Topic，在节点网络中，它们的作用是一致的，即“输出或处理一段标准数据”。

---

## 2. 核心抽象模型 (Core Abstractions)

### 2.1 统一的数据载体：CraftFeed (The Water)
无论数据处于抓取阶段、加工阶段还是最终输出阶段，在系统内部统一使用派生自 RSS 结构的标准数据模型。它不仅包含标准的 RSS 字段，还包含 AI 加工留下的 Metadata 和流转链路追踪。

* **`CraftFeed`**: 包含 Feed 级别的元数据和 `Articles` 列表。
* **`CraftArticle`**: 单篇文章/内容，完美契合系统“全文提取、摘要、打分”的业务本质。包含扩展字段如 `Depth` (裂变层级), `QualityScore` (质量分), `OriginalSourceID` (来源追踪) 等。

### 2.2 统一的行为节点 (The Pipes)
我们将所有的业务模块抽象为两类接口：

#### 抽象一：`FeedProvider` (数据提供者 / 生成器)
**定义**：任何能够输出 `CraftFeed` 的节点。
```go
type FeedProvider interface {
    GetFeed(ctx context.Context) (*CraftFeed, error)
}
```
* **RawSource**: 基础爬虫节点 (如 `HtmlSource`, `SearchSource`, 第三方 `RssSource`)。
* **Recipe**: 复合节点（包装器）。其本质是“一个 RawSource + 一组 Processor”，对外暴露的依然是 `FeedProvider` 接口。
* **Topic**: 高级复合节点。其内部并发调用多个 `FeedProvider`，也是一个 `FeedProvider`。由于实现了统一接口，Topic 内部不仅可以嵌套 Recipe，甚至可以嵌套另一个 Topic。

#### 抽象二：`FeedProcessor` (数据加工器)
**定义**：接收一个 `CraftFeed`，对其进行修改/过滤/丰富，并输出处理后 `CraftFeed` 的节点。
```go
type FeedProcessor interface {
    Process(ctx context.Context, feed *CraftFeed) (*CraftFeed, error)
}
```
* **AtomCraft (原子加工)**：如翻译、摘要、去广告。
* **FlowCraft (加工流)**：多个 Processor 组成的链条。
* **Aggregator / Deduplicator**: 即将在 Topic 中实现的“合并、去重、截断”逻辑，本质上也是一个特殊的 Processor。
* **Expander (资讯裂变)**：未来将实现的内容延展节点。

因为接口统一，**Processor 可以被随意插拔**。例如：AI 摘要 (Summary Craft) 既可以挂载在单源 Recipe 的抓取之后，也可以挂载在 Topic 聚合去重之后的最后一步执行。

---

## 3. 未来的可视化编排蓝图 (Data Flow Map)

有了这套抽象，系统未来的运转将完全契合可视化节点连线（Node-based UI）的交互模式：

```mermaid
graph LR
    subgraph "Provider Nodes (数据提供者)"
        S1[HtmlSource: 独立博客]
        S2[RssSource: 官方RSS]
        S3[SearchSource: 关键词搜索]
    end

    subgraph "Processor Nodes (数据加工器)"
        C1[Translate Craft: 翻译为中文]
        C2[Summary Craft: AI 提取核心观点]
        C3[Filter Craft: 过滤短文]
    end

    subgraph "Router/Struct Nodes (结构节点)"
        M1[Merge Aggregator: 合并去重并保留Top 20]
    end

    %% 灵活的乐高式连线
    S1 --> C1
    S2 --> C3

    C1 --> M1
    S3 --> M1
    C3 --> M1

    M1 --> C2  %% 全局后置加工：对汇总后的高价值内容统一摘要

    C2 --> Output((最终输出: 统一订阅端点))
```

---

## 4. 现有架构的改造路径 (Refactoring Path)

当前的架构（参考 `current_architecture.md`）已经有了良好的解耦基础（Source 与 Craft 分离），但在数据类型上严重依赖第三方的 `gofeed.Feed`，导致扩展性受限。为了向“乐高架构”演进并支持 Topic MVP 的开发，我们需要按以下路径进行重构改造：

### Step 1: 引入核心载体 (引入 `CraftFeed` / `CraftArticle`)
* **现状**: 系统各个层级（Source, Craft）直接传递和修改 `*gofeed.Feed`。
* **改造**: 在 `internal/model` (或类似的领域定义目录) 中定义 `CraftFeed` 和 `CraftArticle` 结构体。提供 `FromGoFeed()` 和 `ToGoFeed()` 的转换方法以保持向后兼容。

### Step 2: 接口规范化 (标准化 Provider 与 Processor)
* **现状**: Source 层返回的是 `gofeed.Feed`，Craft 层的接口签名依赖 `gofeed.Feed` 且往往只考虑 Recipe 级别的调用。
* **改造**:
  * 将所有 `Source` 的返回值修改为 `*CraftFeed`（使其符合 `FeedProvider` 的语义）。
  * 将 `Craft` 层的 `Process` 方法签名更新为接收和返回 `*CraftFeed`（使其符合 `FeedProcessor` 的语义）。消除 Craft 内部对特定业务实体（如 Recipe 数据库模型）的强耦合。

### Step 3: 业务编排层的重构 (Recipe 与 Topic)
* **现状**: `Recipe` 的执行逻辑硬编码了 "获取 -> 遍历 Craft" 的过程。
* **改造**:
  * 将 `Recipe` 的执行过程视为一次简单的 Pipeline 构建：将对应的 `Source` 视为入口节点，将配置的 `Crafts` 包装为一条 `Processor` 链条。
  * 在开发新的 `Topic` 特性时，直接拥抱这套新抽象：让 `Topic` 作为一个包裹了多个 `Recipe (FeedProvider)` 并串联了 `Aggregator (FeedProcessor)` 的复合节点执行。

通过以上三个步骤的改造，我们能在开发 Topic 功能（MVP）的同时，顺便打通底层的数据流管线，为未来彻底的“节点编排化”奠定坚实的基础。