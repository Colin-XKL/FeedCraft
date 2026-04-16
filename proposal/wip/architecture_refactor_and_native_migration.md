# 领域五：架构深化与原生迁移

> 状态：进行中

## 1. 背景

TopicFeed 的长期目标不是只在现有 legacy 执行流外面再包一层，而是推动整个 FeedCraft 真正落到统一的 `FeedProvider` / `FeedProcessor` 架构上。

当前系统已经有：

- `CraftFeed` 作为统一数据载体
- `FeedProvider` / `FeedProcessor` 作为统一接口
- `TopicFeed` 作为多源聚合节点

但 Recipe 仍然大量走旧执行流，底层 Source 和 Craft 也还没有完全原生化。这会让 Topic 的实现长期停留在“混合架构”状态，增加装配复杂度和维护成本。

除了执行流本身，当前系统在“输入模型”这一层也还没有真正统一：

- Recipe 主要使用 `SourceConfig`
- Topic 主要使用 `input_uris`

这两个结构各自合理，但都不适合作为整个系统统一的顶层输入抽象：

- `SourceConfig` 太偏底层抓取和解析细节
- `InputURI` 又不足以表达复杂 HTML、JSON、search 配置

因此，长期架构收敛的关键之一，是统一顶层输入模型，而不是简单在 `InputURI` 和 `SourceConfig` 之间二选一。

## 2. 需求场景

从长期演进看，用户并不关心数据来自 Recipe 还是 Topic，他们期待的是：

- 任意节点都能被复用
- 任意节点都能嵌套
- 整个系统围绕统一的 feed graph 执行

要支撑这种能力，系统内部必须让：

- Recipe 成为真正的一等 `FeedProvider`
- 底层 Source 原生输出 `CraftFeed`
- 底层 Craft 原生处理 `CraftFeed`
- Recipe 和 Topic 使用统一的输入引用模型

## 3. 当前进展

- `CraftFeed` 及其转换函数已经存在，见 `internal/model/feed.go`
- `FeedProvider` / `FeedProcessor` 接口已经存在，见 `internal/engine/interfaces.go`
- `LegacySourceAdapter` 已存在，说明当前已采用过渡式迁移方案，见 `internal/source/adapter.go`
- Topic 聚合核心已经基于新接口工作

当前缺口：

- `ProcessRecipeByID` 仍在使用旧的 `source.Get()` + `craft.ProcessFeed()` 编排，见 `internal/recipe/custom_recipe.go`
- Recipe 还没有成为标准 runtime provider
- Source 层大多仍返回 `gofeed.Feed`
- Craft 层大多仍围绕旧数据结构工作
- 顶层输入模型尚未统一，Recipe 与 Topic 仍然分裂

## 4. 未来待办

### 4.1 统一顶层输入模型

推荐统一方向：

- 顶层统一使用 `InputSpec`
- `InputSpec(kind=uri)` 用于轻量引用
- `InputSpec(kind=source)` 用于承载复杂 `SourceConfig`

这样：

- `RecipeFeed = InputSpec + Craft`
- `TopicFeed = []InputSpec + Aggregator`

语义分层明确为：

- `InputSpec` 是顶层输入抽象
- `InputURI` 是轻量引用形式
- `SourceConfig` 是底层原始源配置

### 4.2 统一输入解析与 provider builder

- 构建统一的 `InputSpec -> FeedProvider` builder
- 让 Topic 和 Recipe 共用同一套输入解析逻辑
- 让 `feedcraft://recipe/:id`、`feedcraft://topic/:id`、`http(s)://...` 都成为统一语义的一部分

其中：

- `http(s)` 默认解释为第三方网站 RawFeed
- 语义参考“最小 source + proxy craft”
- 不是 Topic 专属特例，而是系统统一输入规则

### 4.3 重构 RecipeFeed 执行流

- 让 Recipe 在运行时成为标准 `FeedProvider`
- 将 Recipe 的 source + craft 组合收敛为统一执行对象
- 让 Topic 在解析 `feedcraft://recipe/:id` 时不需要再走特殊分支

### 4.4 推进 Source 原生迁移

- 逐步让底层 Source 直接返回 `*model.CraftFeed`
- 减少对 `gofeed.Feed` 的中间转换依赖
- 最终移除 `LegacySourceAdapter`

### 4.5 推进 Craft 原生迁移

- 逐步让 AtomCraft / FlowCraft 直接处理 `*model.CraftFeed`
- 消除旧 `feeds/gofeed` 双模型往返转换
- 最终移除旧适配层

### 4.6 统一 Topic / Recipe 的运行时模型

- 让 Topic 和 Recipe 都能以统一方式被装配、执行、缓存、观测
- 让后续内置 Feed、系统通知 Feed 等能力都复用同一套运行时抽象

### 4.7 采用两阶段迁移策略

第一阶段：

- 统一运行时
- 旧 Recipe 的 `SourceConfig` 在读取时映射为 `InputSpec(kind=source)`
- 旧 Topic 的 `input_uris` 在读取时映射为 `[]InputSpec(kind=uri)`

第二阶段：

- 再考虑统一持久化结构
- 逐步减少 Recipe / Topic 在存储层的模型差异

## 5. 追踪建议

这个领域完成的标准是：

- 顶层输入模型完成统一，不再分裂成 Recipe 专属和 Topic 专属两套语义
- Topic 和 Recipe 在引擎层面真正平级
- 新能力接入时不再区分“旧链路”与“新链路”
- 适配器层从临时迁移手段，逐步收缩到可以删除的状态
