# 领域二：Topic 运行时解析与执行装配

> 状态：已完成

## 1. 背景

`TopicFeed` 当前已经有运行时引擎 `internal/engine/topic.go`，但这个引擎需要的是已经构造好的 `Inputs []FeedProvider` 和 `Aggregator FeedProcessor`。

数据库里保存的是配置态信息：

- `input_uris`
- `aggregator_config`

这意味着 Topic 真正可运行之前，还缺少一层“装配器”：

- 把配置态输入翻译成统一的输入引用模型
- 把输入引用解析为具体 provider
- 把聚合配置编译成 processor 链
- 递归构建嵌套 Topic

如果这层没有完成，TopicFeed 就只能停留在“可编辑配置”，而不能真正执行。

进一步看，这个问题并不是 Topic 独有的。

`RecipeFeed` 当前虽然可以运行，但它的输入层仍然主要依赖 `SourceConfig`。这适合表达底层抓取和解析细节，却不适合作为整个系统统一的顶层输入模型。另一方面，单纯的 `InputURI` 又不足以承载 HTML、JSON、search、curl-to-rss 这类复杂源配置。

因此，这个领域的设计目标不应只是“补一个 URI resolver”，而应该建立一层统一输入抽象：

- 顶层统一使用 `InputSpec`
- `InputSpec(kind=uri)` 用于轻量引用
- `InputSpec(kind=source)` 用于承载复杂 `SourceConfig`

这样 `TopicFeed` 和 `RecipeFeed` 才能复用同一套输入解析和运行时装配能力。

## 2. 需求场景

用户创建 Topic 或 Recipe 时，通常不会手工实例化代码对象，而是通过配置来描述输入来源。

典型轻量引用：

- `feedcraft://recipe/hn-tech`
- `feedcraft://recipe/reddit-open-source`
- `feedcraft://topic/dev-tools`
- `https://example.com/some-page`

典型复杂输入：

- 一个 HTML SourceConfig，带 selector
- 一个 JSON SourceConfig，带 jq selector 和 template
- 一个 search SourceConfig，带 query 和增强模式

系统需要自动识别这些输入，并构建出最终可执行的运行时拓扑：

- 内部 Recipe -> 包装为 provider
- 内部 Topic -> 递归构建子 Topic
- 外部 `http(s)` -> 视为第三方网站 RawFeed
- 内嵌 SourceConfig -> 包装为原始 provider

同时，用户还会配置：

- `deduplicate`
- `sort`
- `limit`

这些步骤也需要自动编译成执行链。

## 3. 当前进展

- `engine.TopicFeed` 已定义统一的运行时结构，见 `internal/engine/topic.go`
- `FeedProvider` / `FeedProcessor` 接口已存在，见 `internal/engine/interfaces.go`
- 去重、排序、截断处理器已存在，见 `internal/engine/aggregator.go`
- 适配旧 Source 的 `LegacySourceAdapter` 已存在，说明架构上允许先通过适配器接入外部源，见 `internal/source/adapter.go`
- `RecipeFeed` 现有执行入口已经存在，见 `internal/recipe/custom_recipe.go`

当前缺口：

- 没有统一的 `InputSpec -> FeedProvider` builder
- 没有 `RecipeFeed` 和 `TopicFeed` 共用的输入解析层
- 没有把 `dao.TopicFeed.AggregatorConfig` 转换为 processor pipeline 的 builder
- 没有嵌套 Topic 的递归装配逻辑
- 没有循环依赖保护

## 4. 未来待办

### 4.1 引入统一输入模型 `InputSpec`

统一顶层输入抽象，不直接在运行时层面二选一使用 `InputURI` 或 `SourceConfig`。

建议的首版模型：

- `InputSpec(kind=uri, uri=...)`
- `InputSpec(kind=source, source_config=...)`

语义分层：

- `InputSpec` 是顶层输入模型
- `InputURI` 是 `kind=uri` 的一种具体内容
- `SourceConfig` 是 `kind=source` 的底层叶子配置

### 4.2 实现统一输入解析工厂

创建统一 builder，例如：

- `BuildProviderFromInput(ctx, spec, stack)`
- `BuildTopicProvider(ctx, topicID)`
- `BuildTopic(ctx, daoTopic)`

首版 `kind=uri` 支持：

- `feedcraft://recipe/:id`
- `feedcraft://topic/:id`
- `http://...`
- `https://...`

其中 `http(s)` 的默认解释语义为第三方网站 RawFeed：

- 参考 Recipe 里的处理方式
- 底层复用现有 Source 能力
- 仅应用 `proxy` craft 语义
- 不将其收窄为“必须是标准 RSS 地址”

### 4.3 实现 Topic runtime builder

- 从 `dao.TopicFeed` 构造 `engine.TopicFeed`
- 首版继续读取现有 `input_uris`
- 在运行时把 `input_uris` 映射为 `[]InputSpec(kind=uri)`
- 递归构建 provider graph
- 将 Topic 自身元信息，如 `ID`、`Title`、`Description`，一并注入运行时对象

### 4.4 实现聚合配置 builder

- 将 `AggregatorConfig` 编译为 processor pipeline
- 支持：
  - `deduplicate`
  - `sort`
  - `limit`
- 对非法配置提供明确报错，例如：
  - 未知 step type
  - 缺少必填 option
  - option 值格式不合法

### 4.5 增加循环依赖保护

- 处理 `topic -> topic -> topic` 的递归构建
- 对 `A -> B -> A` 这类环路做检测并快速失败
- 错误信息要明确指向循环链路，便于后台排查

### 4.6 规划两阶段统一路径

第一阶段：

- 只统一运行时
- Topic 继续存 `input_uris`
- Recipe 继续存 `SourceConfig`
- 读取时统一映射为 `InputSpec`

第二阶段：

- 再考虑让 Recipe 和 Topic 的持久化模型逐步收敛到统一输入结构
- 避免一开始就做大规模 schema 重构

### 4.7 增加装配层测试

- `InputSpec(kind=uri)` 解析测试
- `InputSpec(kind=source)` 构建测试
- `AggregatorConfig` 编译测试
- 嵌套 Topic 构建测试
- 循环依赖测试

## 5. 追踪建议

这个领域完成的标准是：

- Topic 和 Recipe 可以共享同一套输入解析能力
- Topic 配置可以被完整翻译成运行时对象
- 所有输入引用都能得到确定的执行语义
- 错误配置和递归环路能被及时拦截，而不是拖到运行时随机失败
