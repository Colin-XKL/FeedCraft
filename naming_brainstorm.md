# FeedCraft 核心概念与命名头脑风暴 (Naming Brainstorm)

## 1. 为什么我们需要这次头脑风暴？

FeedCraft 正在进行一次重大的底层架构升级。我们的目标是打造一个“乐高式”的数据流处理引擎，支持高度灵活的插件化组合，甚至未来支持类似 n8n 或 ComfyUI 的可视化节点连线编排。

在这个过程中，我们发现系统现有的一些核心名词（如 Recipe, Topic）在早期的产品形态中很容易理解，但在新的、高度统一的底层代码架构下，它们的名字已经无法准确传达它们之间的关联与本质了。

为了让系统的**产品表达（给用户看）**和**代码架构（给开发者看）**都达到高度的统一与优雅，我们需要重新审视并统一这套命名体系。

---

## 2. 核心概念回顾：过去与现状 (Background)

为了让大家在同一个语境下思考，这里先回顾一下 FeedCraft 已有的几个核心名词的来源与作用。
**FeedCraft 的项目愿景是："Craft all your feed in one place!"（在一个地方精雕细琢你所有的订阅源）。它的名称和 Logo 甚至致敬了 MineCraft，初衷是做一个简单易用、同时像沙盒游戏一样足够灵活的 RSS 工具中间件。**

基于这个愿景，我们可以把系统的能力划分为清晰的“三层架构 (Layered Architecture)”：

### 2.1 数据生成层 (RSS Generators)

FeedCraft 不仅能处理已有的 RSS，还能**无中生有**地创造 RSS。

- **生成器 (Generators)**：例如 `HTML-to-RSS` (通过可视化选择器将普通网页转为 RSS)、`Curl-to-RSS` (将 API JSON 响应转为 RSS)、`Search-to-RSS` (将搜索引擎结果转为 RSS)。
- **作用**：在底层，它们被称为 `Source`（数据源），负责在系统入口处，将各种非标准的数据抓取并统一转化为标准的 RSS 数据流。

### 2.2 加工处理层 (The Crafts)

当数据被生成或抓取后，就进入了由 `Craft` (手艺/加工/工艺) 主导的加工层：

- **概念来源**：FeedCraft 名字的后半部分。系统中所有负责修改、过滤、增强 RSS 数据的逻辑单元，统统被称为 `Craft`。
- **AtomCraft (原子工艺)**：系统中**最小的处理单元**。执行单一动作。除了内置的工艺（如 `translate-title`、`fulltext` 提取全文、`summary` AI 摘要），用户还可以基于模版创建自定义的原子工艺。
- **FlowCraft (组合工艺)**：多个 AtomCraft 的**组合序列**。这允许用户将多个操作串联起来。例如：创建一个名为 `digest-and-translate` 的 FlowCraft，包含 `提取正文 -> 生成摘要 -> 翻译内容`。

### 2.3 路由与组合层 (Recipes & Topics)

加工好的数据，最终需要被打包成一个“订阅端点 (Endpoint)”提供给用户。

- **Recipe (配方)**
  - **概念来源**：借用烹饪的隐喻。将特定的数据源 URL（比如通过 HTML-to-RSS 抓取的一块生肉）与某个 Craft（比如特定的煎烤方式）**绑定**，产生一个持久化的新订阅源。
  - **作用**：它代表了一条**“单向的数据生产线”**。比如：“从 Hacker News 的 URL 获取数据，经过 `digest-and-translate` 处理，最后输出一个带有全文和摘要的 RSS 订阅源”。
- **Topic (主题)**
  - **概念来源**：随着需求增加，用户希望把多个不同来源的信息汇聚在一起阅读。
  - **作用**：它代表了一个**“多源的数据聚合池”**。比如“AI 前沿资讯”这个 Topic，它的底层会同时并发拉取 OpenAI 的 Recipe、Anthropic 的 Recipe 以及某个推特博主的 Recipe，把它们汇聚在一起，进行合并、去重、排序截断后，统一输出成一个大 RSS 订阅源。

---

## 3. 现在的痛点是什么？(The Problem)

在最新的架构抽象中，整个系统被极度精简为两种底层能力：

1.  **Processor（加工者）**：吃进一份数据，处理后吐出新数据。
2.  **Provider（提供者）**：不需要输入，直接凭空吐出数据。

在目前的命名下，**Processor（加工者）的命名是非常和谐的**：
无论是 `Atom-Craft` 还是 `Flow-Craft`，它们都共享了 `-Craft` 这个后缀，让人一眼就能看出它们的作用是“处理数据”。

**但是，Provider（提供者）的命名出现了严重的割裂**：
在底层引擎看来，`Recipe` 和 `Topic` 的本质是完全一样的——它们都是**“可以被调用的、最终能吐出 RSS 数据的节点 (Provider)”**。

- 既然它们本质一样，那么**理论上它们应该可以无限嵌套组合**。比如：Topic 里面不仅可以塞入 Recipe，还可以塞入另一个 Topic。
- 但 `Recipe (配方)` 和 `Topic (主题)` 这两个词在语义上是平行且割裂的。它们无法体现出“两者皆是数据流节点”的统一感，也无法体现出这种奇妙的“分形与套娃”关系。

---

## 4. 目标与决定：转向“xxxFeed”命名体系

经过前期的初步探讨，我们决定放弃为这两个“节点”寻找抽象的架构名词（如 Stream、Source、Branch 等），而是**转换视角，以最终的数据产出物（即 Feed）来重新定义它们**。

我们提议使用 **`xxxFeed`** 的后缀体系：

- **RawFeed (原始源)**：对应底层无加工的抓取结果，或者直接抓取到的原生 RSS。
- **RecipeFeed (配方源)**：取代原有的 `Recipe`。代表经过加工管线（FlowCraft）处理后产出的 Feed。
- **TopicFeed (主题源)**：取代原有的 `Topic`。代表经过聚合、去重后产出的宏大 Feed。

### 为什么选择 `xxxFeed` 方案？

1.  **极低的用户认知门槛**：用户使用 FeedCraft 的终极目的就是为了获取一个“订阅源 (Feed)”。无论是创建一个配方，还是建立一个主题，最终交付给用户的都是一个 `Feed`。这种命名在产品 UI 上会显得极其自然。
2.  **完美的统一感**：它彻底解决了“提供者(Provider) 命名割裂”的痛点。不管数据是怎么来的，只要带有 `Feed` 后缀，它在系统里就是一个可以被消费、被订阅的同类实体。
3.  **顺畅的嵌套逻辑**：“我把 3 个 `RawFeed` 和 2 个 `RecipeFeed` 合并在一起，塞进了一个 `TopicFeed` 里。”——这种关于多源聚合嵌套的口头表达，变得无比顺滑，没有任何逻辑违和感。

---

## 5. 本次头脑风暴需要探讨的细节

虽然大方向（`xxxFeed`）已经确定，但在具体落地和架构适配上，依然有几个关键细节需要大家集思广益：

### 议题一：处理动作（Craft）的命名是否需要跟进微调？

目前我们有 `AtomCraft` 和 `FlowCraft`。
在新的视角下，整个系统的工作流变成了：

> **`RawFeed` -> 经过 `Craft` 加工 -> 变成 `RecipeFeed`**

这个“动静结合”的隐喻（Craft 作用于 Feed）非常自洽。大家认为 `Craft` 相关的命名是否需要保持现状，还是有更匹配 `xxxFeed` 体系的叫法？（比如叫 FeedCraft, FlowCraft 保持不变？）

### 议题二：底层代码中的接口命名

在后端引擎代码中，我们需要一个接口来统领 `RawFeed`、`RecipeFeed` 和 `TopicFeed`。
以前我们管这个接口叫 `FeedProvider`（或者想叫 `Stream`、`Source`）。
如果业务概念叫 `RecipeFeed`，在底层代码中：

```go
// 接口名应该叫什么？FeedProvider 还是 FeedStream 还是就是 FeedBuilder?
type FeedProvider interface {
    // 方法名应该叫什么？GetFeed()? Fetch()? 构建一个 Feed 的 Feed 听起来有点怪
    Fetch(ctx context.Context) (*CraftFeed, error)
}
```

请大家从代码的优雅度出发，对接口的方法签名提出建议。

### 议题三：是否有比 Recipe 和 Topic 更好的前缀？

虽然确定了 `xxxFeed` 的后缀，但 `RecipeFeed` 和 `TopicFeed` 这两个前缀词是否是最终完美的选择？
比如：

- 单线流：叫 `RecipeFeed` 还是 `SoloFeed` 还是 `PipeFeed`？
- 聚合流：叫 `TopicFeed` 还是 `MergeFeed` 还是 `HubFeed`？
  考虑到我们已经沉淀了部分用户心智，保留 `Recipe` 和 `Topic` 前缀是一个稳妥的选择，但大家也可以提出更惊艳的组合。
