# 内容裂变与延展 (Content Fission & Expansion)

## 1. 需求场景 (Scenario)
在常规的 RSS 聚合之上，系统可以通过 AI 赋予订阅流“自我繁衍”的能力，我们称之为“资讯裂变”或“内容延展”。
当用户订阅了一个 Topic（如 `popular-selfhost-app`）时，系统不仅返回该 Topic 原本抓取到的文章，还能基于这些文章的内容，自动去发现和补充更深度的信息。

**典型场景：**
1. **关键词衍生搜索 (Keyword Expansion)**：
   * 原文章提到了一款有趣的软件：“最近尝试了 *some-awesome-app*，体验很好”。
   * 系统通过 LLM 识别出核心实体 `some-awesome-app`。
   * 系统自动复用底层的 `Search-to-RSS` 模块，去搜索引擎抓取关于该软件的最新动态或官网介绍，并作为新文章追加到当前订阅流中。
2. **内链深度提取 (Link Extraction)**：
   * 原文章是一篇“每周精华汇总”，里面包含了多个外部文章链接。
   * 系统自动解析 HTML，提取高价值的外链。
   * 系统自动复用底层的 `HTML-to-RSS` 或 `Fulltext` 模块，把这些外链对应的完整文章抓取下来，追加到当前订阅流中。

## 2. 核心机制设计 (Core Mechanism)
在现有的 Source-Craft 架构下，这种功能将作为一种特殊的 **Expander Craft（延展处理单元）** 存在。与常规的 `AtomCraft`（1对1修改属性）不同，它是一个 **1对多** 的生成器。

* **触发点**：读取现有 `CraftItem` 的正文或标题。
* **处理逻辑**：调用 LLM 分析或正则解析提取目标（关键词/链接）。
* **内部调用**：实例化对应的 `Source` (如 `EnhancedSearchSource` 或 `HtmlSource`) 执行抓取。
* **结果合并**：将抓取到的新 `CraftItem` 集合附加到当前的 `CraftFeed.Items` 中。

## 3. 架构适配与关键挑战 (Architecture Fit & Challenges)

当前的系统架构在结构上完全支持该功能，但在未来实现时需严格处理以下边界问题：

### 3.1 防止无限死循环 (Infinite Fission Loop)
* **挑战**：文章 A 裂变出文章 B，文章 B 的内链又指向文章 A，导致系统陷入无限抓取死循环。
* **架构支撑**：依赖系统内部统一的 `CraftItem` 模型。未来将在 `CraftItem` 中引入 `Depth` (裂变层级) 和 `ParentID` 字段。
* **策略**：在 Expander Craft 的入口处严格限制 `if item.Depth > 0 { return }`，即只允许对原生内容（Depth=0）进行一次裂变。

### 3.2 性能雪崩与执行位置 (Performance & Pipeline Position)
* **挑战**：裂变过程包含 LLM 调用和二次网络抓取，耗时极长。如果对 Topic 下合并的 100 篇文章全部执行裂变，会导致接口严重超时。
* **架构支撑**：依赖当前的 **Branch -> Trunk 两级流水线** 设计。
* **策略**：
  * **必须作为 Topic 层级的全局 Craft 执行**。
  * **必须在 Truncation (截断) 之后执行**。即先将多源数据合并、去重，截断保留 Top 10，然后再对这最有价值的 10 篇文章（甚至只挑 Top 2）执行裂变。避免大量裂变出来的数据在最后一步被 Limit 截断丢弃。
  * 由于耗时较长，未来实现该功能时，Topic 的更新策略可能需要从 MVP 的“同步实时拉取”升级为“后台异步队列预热 (Async Background Preheating)”。

### 3.3 来源追溯与用户体验 (Context Linking)
* 为了避免用户在 RSS 阅读器中看到突然冒出来的衍生文章感到困惑，新生成的 `CraftItem` 必须在正文顶部/底部注入来源声明。
* 例如：`<blockquote>🤖 FeedCraft AI 延展阅读：本文由你订阅的《[原文章标题]》提及的关键词 [some-awesome-app] 触发衍生搜索生成。</blockquote>`
