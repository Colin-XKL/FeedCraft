# 网页内容变化监控 (Web Monitor) 方案设计规划

> 状态：未开始

## 1. 背景知识与需求场景

### 1.1 背景

FeedCraft 作为一个强大的 RSS 信息流聚合与加工系统，已经具备了将静态/动态网页转换为标准 RSS（`html-to-rss`）、全文提取、甚至通过 LLM 进行信息重构的能力。
然而，在实际应用中，用户除了“订阅列表更新”之外，还有一种极其常见的长尾需求：**监控网页上特定数值/状态的变化**。

### 1.2 需求场景

典型的应用场景包括但不限于：

- **电商比价**：监控某款特定商品（如 Amazon 上的特定型号）的价格，价格变动时收到通知。
- **库存监控**：监控某个稀缺商品（如演唱会门票、限量手办）是否从“Out of Stock”变成了“Add to Cart”。
- **状态追踪**：监控某个政务/签证页面的“当前处理日期”是否推进。
- **通知模版定制**：用户不仅想知道“内容变了”，还希望收到可读性强的通知，例如：“【补货提醒】您关注的 PS5 当前库存状态：有货，价格为 $399”。

传统上，这种监控需要后端引入复杂的定时任务（Cron）、历史状态存储（Database）以及对比逻辑（Diff），架构会变得极其沉重。我们需要在 FeedCraft 中以最优雅、最轻量的方式实现这一能力。

## 2. 核心架构方案：无状态的值驱动模型

为了保持 FeedCraft 架构的简洁性，我们将采用 **“无状态的值驱动更新 (Stateless Value-Driven Updates)”** 模型。其核心精髓在于：**利用 RSS 阅读器自带的去重机制，通过巧妙的 GUID 哈希策略来替代后端的数据库状态存储。**

### 2.1 整体数据流

我们将新增一个 `WebMonitorSource`，它复用现有的请求层（支持 Browserless 突破反爬和渲染 JS），并串联一个全新的 `WebMonitorParser`。

1. **拉取网页**：当用户的 RSS Reader 请求 Feed 时，触发 HTTP/Browserless 拉取最新网页内容。
2. **多变量提取**：根据用户配置的多个 CSS 选择器，提取出对应的文本值（如 `price=$399`, `stock=In Stock`）。
3. **渲染模板**：将提取到的变量注入到用户自定义的标题/描述模板中（使用 Go 的 `text/template`）。
4. **生成稳定哈希 (防抖核心)**：将 **目标 URL** 与 **所有提取到的变量值** 拼接后进行 MD5 哈希，作为这篇文章的唯一 `GUID`。
5. **RSS 阅读器接管**：
   - 如果数值未变，生成的哈希 `GUID` 相同，RSS 阅读器认为是旧内容，不提示。
   - 如果任何一个监控的数值发生变化，生成的 `GUID` 将改变，RSS 阅读器立刻将其识别为一篇**全新**的文章推送给用户。

### 2.2 核心配置模型

在 `internal/config/source_config.go` 中扩展：

```go
// WebMonitorParserConfig 网页变化监控配置
type WebMonitorParserConfig struct {
	// 提取器：定义变量名和对应的 CSS 选择器。
	// 示例: {"price": ".current-price", "stock": ".stock-status"}
	Extractors map[string]string `json:"extractors"`

	// 通知模版：支持使用 {{.varName}} 语法注入提取到的变量
	// 示例: "【监控更新】价格变动为 {{.price}}，当前状态: {{.stock}}"
	TitleTemplate       string `json:"title_template"`
	DescriptionTemplate string `json:"description_template"`
}
```

## 3. 实施步骤规划

### Phase 1: 核心引擎支持 (Backend Engine)

1. **定义常量与配置**：
   - 在 `internal/constant/source_type.go` 中新增 `SourceWebMonitor = "web_monitor"`。
   - 在 `internal/config/source_config.go` 的 `SourceConfig` 中加入 `WebMonitorParser *WebMonitorParserConfig`。
2. **实现 Parser**：
   - 创建 `internal/source/parser/web_monitor_parser.go`。
   - 依赖 `goquery` 实现基于 Selector 的变量提取。
   - 实现 Map Key 排序后拼接变量以生成稳定 MD5 哈希的防抖逻辑。
   - 实现基于 `text/template` 的模板渲染逻辑（需处理模板执行失败的 fallback，例如 fallback 返回模板原始字符串）。
3. **注册 Source**：
   - 创建 `internal/source/web_monitor.go`。
   - 仿照 `html.go`，使用 `PipelineSource` 将现成的 `HttpFetcher` 与新的 `WebMonitorParser` 组合，并注册到引擎工厂中。
4. **编写单元测试**：
   - 重点验证：相同变量提取结果生成相同的 GUID，变量一变 GUID 立即变动。
   - 验证模板渲染的正确性和容错性。

### Phase 2: API 与前端联调支持 (API & Web UI) - 后续跟进

1. **无需改动数据库**：得益于 `CustomRecipeV2` 的 JSON 存储设计，后端 API 层面无需任何 Schema 修改。
2. **前端页面适配**：
   - 在 `Custom Recipe` 创建页新增类型 "Web Monitor"。
   - 复用现有的 `HttpFetcher` 配置 UI（URL 输入、Browserless 开关）。
   - **升级可视化选择器**：允许用户在一个页面上点选多个元素，并为每个选中的 CSS Selector 分配一个变量名（如 `price`, `status`）。
   - 新增模板输入框（Title 和 Description），并在旁边提供可用变量的提示（如 `{{.price}}`, `{{.url}}` 等内置变量）。

## 4. 需要注意的细节与风险

1. **页面结构的不稳定性**：CSS Selector 依赖页面 DOM 结构，如果目标网站改版，监控可能会失效（提取不到值）。
   - **应对策略**：如果 `doc.Find(selector)` 没有匹配到元素，可以返回空字符串或特定的提示（如 `[Not Found]`）。只要空字符串这个状态与之前的值不同，依然会触发一次提醒，告知用户结构可能变了。
2. **Hash 生成的稳定性**：Go 语言中 `map` 的遍历是无序的。
   - **强制要求**：在拼接哈希字符串之前，必须对 `Extractors` 的 Key 进行字典序排序，确保同样的变量集在任何时候生成的哈希串绝对一致，否则会导致疯狂的错误通知风暴。
3. **首尾空格处理**：网页上提取的文本经常包含大量不可见的换行或空格。
   - **强制要求**：提取到元素的 `Text()` 后，必须做 `strings.TrimSpace()` 清理，避免不可见字符的微小变化导致误触发。
4. **内置可用变量**：除了用户自定义提取的变量外，渲染模板时应自动注入 `{{.url}}` 作为兜底变量，方便用户在模板里直接贴上原链接。
