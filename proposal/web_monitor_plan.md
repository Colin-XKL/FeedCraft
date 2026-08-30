# 网页内容变化监控 (Web Monitor) 方案设计规划

> 状态：规划中

## 1. 背景知识与需求场景

### 1.1 背景

FeedCraft 作为一个 RSS 信息流聚合与加工系统，已经具备了将静态/动态网页转换为 RSS（`html-to-rss`）、全文提取、以及基于 recipe 的可扩展抓取处理能力。

但在实际使用中，除了“网页列表新增内容”之外，还有一类非常高频、且很适合 RSS 交付的需求：**监控网页上某些特定字段是否发生变化**。

例如，用户并不关心整个页面是否改版，而只关心：

- 某个商品价格是否变化
- 某个库存状态是否从无货变为有货
- 某个排期日期是否推进
- 某个按钮文案是否从 “Coming Soon” 变成 “Buy Now”

这类需求如果用传统监控系统实现，通常会引入：

- 后端定时任务（Cron / Scheduler）
- 历史状态存储（Database / Cache）
- 对比逻辑（Diff）
- 通知分发逻辑

这会让整体架构明显变重。

FeedCraft 当前本身是 **请求驱动 + preheating** 的工作模式，而不是一个独立的定时巡检平台。因此，这个功能最合适的方向不是再造一个监控子系统，而是：**把“网页变化监控”优雅地建模成一个新的 source type，继续复用现有抓取、解析、recipe 和 RSS 输出链路。**

### 1.2 需求场景

典型场景包括但不限于：

- **电商比价**：监控某款商品价格，价格变动时通过 RSS Reader 收到提醒。
- **库存监控**：监控某个稀缺商品是否从 “Out of Stock” 变成 “In Stock” / “Add to Cart”。
- **状态追踪**：监控签证、政务、物流等页面上的状态字段是否推进。
- **模板化通知**：用户不仅想知道“变了”，还想收到可读性更好的内容，例如：
  - `【补货提醒】PS5 当前库存状态：有货`
  - `【价格更新】当前价格 $399，库存 In Stock`
- **重点字段监听**：用户可能只希望“价格变化”触发通知，但在 RSS 正文里仍然展示价格、库存、标题、链接等更丰富信息。

---

## 2. 核心设计目标

这个功能的设计目标应明确为：

1. **保持架构轻量**

   - 不新增调度子系统
   - 不引入额外的状态表或 diff 存储
   - 继续复用当前 recipe/source pipeline

2. **变化判定可由用户控制**

   - 用户应能指定“哪些字段参与变化判定”
   - 避免无关字段变化导致通知风暴

3. **阅读内容可富文本化**

   - 用户应能自定义 RSS 的标题、摘要、正文模板
   - 用于在通知之外提供更好的阅读体验

4. **与现有能力自然融合**
   - 复用 `HttpFetcher`
   - 支持 Browserless / JS 渲染页面
   - 复用 goquery selector 解析模式
   - 复用 `CustomRecipeV2` 的 JSON 存储模型
   - 复用现有 `CraftFeed` / `CraftArticle` 输出

---

## 3. 核心架构方案：无状态的值驱动模型

为了保持 FeedCraft 架构的简洁性，建议采用 **“无状态的值驱动更新 (Stateless Value-Driven Updates)”** 模型。

核心思想是：**不在后端保存“上一次监控结果”，而是通过 RSS item GUID 的稳定生成规则，把变化判定交给 RSS Reader 的去重机制。**

### 3.1 整体数据流

我们新增一个 `WebMonitorSource`，复用现有的网页抓取链路，并增加一个新的 `WebMonitorParser`：

1. **拉取网页**

   - 当 RSS Reader 请求 feed 时，执行 HTTP / Browserless 抓取。

2. **提取变量**

   - 根据用户配置的多个 CSS selector，提取多个变量值。
   - 例如：
     - `price = $399`
     - `stock = In Stock`
     - `title = PlayStation 5 Console`

3. **区分“监听字段”和“展示字段”**

   - 所有 extractor 提取出的变量都可以参与模板渲染。
   - 但只有用户显式指定的 `key_fields` 会参与变化判定和 GUID 生成。

4. **渲染 RSS 内容**

   - 使用 Go `text/template` 将变量渲染到：
     - `title_template`
     - `description_template`
     - `content_template`

5. **生成稳定 GUID**

   - 将 **目标 URL** + **`key_fields` 对应变量值** 按稳定顺序拼接后做 MD5。
   - 这个 hash 作为 RSS article 的 GUID。

6. **RSS Reader 接管去重**
   - 如果 key 对应值未变，则 GUID 不变，RSS Reader 认为是旧 item。
   - 如果任一 key 字段变化，则 GUID 改变，RSS Reader 将其识别为新 item。

### 3.2 为什么要把 key fields 与 content 分离

这是本次需求收敛后的关键点。

用户实际需要的是：

- **变化判定尽量精准**：例如只关心价格变化
- **阅读内容尽量丰富**：例如正文中同时展示价格、库存、标题、原始链接等

因此必须显式区分两类信息：

- `extractors`：定义“页面上提取哪些变量”
- `key_fields`：定义“哪些变量决定是否算变化”
- `content_template`：定义“最终文章正文怎么展示”

这样才能避免：

- 库存文案、页面标题或装饰性文案的小变化误触发提醒
- 但又能让 RSS 内容足够完整、可读

---

## 4. 核心配置模型

建议在 `internal/config/source_config.go` 中扩展：

```go
// WebMonitorParserConfig 网页变化监控配置
type WebMonitorParserConfig struct {
	// 提取器：变量名 -> CSS Selector
	// 示例: {"price": ".price", "stock": ".stock-status", "title": "h1"}
	Extractors map[string]string `json:"extractors"`

	// 哪些字段参与变化判定 / GUID 生成
	// 示例: ["price"] 或 ["price", "stock"]
	KeyFields []string `json:"key_fields"`

	// RSS 展示模板
	TitleTemplate       string `json:"title_template"`
	DescriptionTemplate string `json:"description_template"`
	ContentTemplate     string `json:"content_template"`
}
```

### 4.1 配置职责划分

建议职责边界如下：

- `http_fetcher`
  - 负责 URL、headers、method、browserless 等抓取行为
- `web_monitor_parser`
  - 负责 selector 提取
  - 负责 key field 变化判定
  - 负责 title / description / content 模板渲染
- `feed_meta`
  - 负责 feed 级别标题、描述等覆盖

这样结构清晰，也最符合当前项目的 source pipeline 设计。

### 4.2 模板可用变量

模板上下文建议自动注入：

- 所有 extractor 变量
- `url`：目标页面 URL

例如：

```gotemplate
TitleTemplate:
【价格更新】{{.title}}

DescriptionTemplate:
当前价格 {{.price}}，库存状态 {{.stock}}

ContentTemplate:
商品：{{.title}}

价格：{{.price}}
库存：{{.stock}}
链接：{{.url}}
```

---

## 5. 与当前项目结构的契合方式

从现有代码结构看，这个功能最自然的落点是：**新增一个 source type，而不是新增一个独立 subsystem。**

### 5.1 推荐复用点

- `internal/source/html.go`
  - 可直接作为新 source factory 的参考实现
- `internal/source/fetcher/http_fetcher.go`
  - 复用网页抓取、headers、browserless 能力
- `internal/source/parser/html_parser.go`
  - 复用 goquery selector 提取方式和文本清洗习惯
- `internal/source/pipeline.go`
  - 继续承担 fetch → parse → normalize 流程
- `internal/model/feed.go`
  - 继续使用 `CraftArticle.Id -> RSS GUID` 的输出映射
- `internal/dao/recipe.go`
  - 复用 `SourceType + SourceConfig(JSON)` 的存储模型，无需改 schema
- `web/admin/src/views/dashboard/custom_recipe/custom_recipe.vue`
  - 可作为首版最小可用入口
- `web/admin/src/views/dashboard/html_to_rss/html_to_rss.vue`
  - 后续可复用其 fetch + selector 选择体验做可视化向导

### 5.2 推荐后端接入方式

1. 在 `internal/constant/source_type.go` 新增 `SourceWebMonitor = "web_monitor"`
2. 在 `internal/config/source_config.go` 中新增 `WebMonitorParserConfig`
3. 创建 `internal/source/parser/web_monitor_parser.go`
4. 创建 `internal/source/web_monitor.go`
5. 像 `html` source 一样注册到 `internal/source/registry.go`
6. 继续使用 `PipelineSource` 组合：
   - `HttpFetcher`
   - `WebMonitorParser`

这样 builder 层无需新增特殊逻辑，整体最简洁。

---

## 6. 关键实现细节

### 6.1 GUID 生成规则

GUID 的生成应满足：

- 稳定
- 可预测
- 只由用户关心的变化驱动

建议规则：

- 输入：`url + key_fields 对应变量值`
- 对 `key_fields` 做字典序排序
- 按稳定顺序拼接为字符串
- 最后计算 MD5

注意：

- `Extractors` 是 `map`，天然无序，不能直接遍历参与 hash
- `KeyFields` 也应在实际 hash 前进行稳定排序

### 6.2 文本清洗

提取到的文本必须统一做 `strings.TrimSpace()`，避免：

- 换行
- 多余空格
- 页面排版细节
  导致的误触发。

### 6.3 selector 未命中时的处理

当 selector 找不到内容时，建议：

- 返回空字符串
- 不作为 parser error

原因：

- 页面结构变化本身也是一种值得感知的状态
- 只要 key field 对应值变为空，依然可以触发一次提醒
- 用户也可以在 content template 中看出“字段丢失”这一状态

### 6.4 模板错误处理

模板配置错误（语法错误、引用问题）建议：

- **直接返回 error**
- 不做静默 fallback

这样更符合当前项目整体风格：

- 配置错误尽早失败
- 页面数据变化正常产出结果

### 6.5 输出策略

建议首版 parser 产出 **单条 article**：

- `Id`：GUID hash
- `Link`：目标 URL
- `Title`：`TitleTemplate` 渲染结果
- `Description`：`DescriptionTemplate` 渲染结果
- `Content`：`ContentTemplate` 渲染结果
- `Created/Updated`：当前抓取时间

这样更符合“监控一个目标页面当前状态”的产品语义。

---

## 7. 风险与注意事项

1. **页面结构不稳定**

   - CSS selector 依赖 DOM 结构，网站改版可能导致提取失败。
   - 应对方式：未命中返回空字符串，并让结果自然体现在 feed 中。

2. **无关字段误触发**

   - 如果不区分 key fields，任何展示字段的小变动都可能触发新 item。
   - 应对方式：必须支持 `key_fields`。

3. **GUID 稳定性**

   - `map` 遍历无序，若不排序会造成错误通知风暴。
   - 应对方式：hash 前对字段做稳定排序。

4. **页面文本噪音**

   - 空格、换行、缩进等小变化可能导致误触发。
   - 应对方式：统一 `TrimSpace()`。

5. **首版不是后台主动监控系统**
   - 当前架构是请求驱动，不是 cron 巡检。
   - 应对方式：首版明确定位为“请求时重新抓取并生成变化驱动的 feed”。

---

## 8. 分阶段实现规划

### Phase 1：核心后端能力（优先实现）

目标：先把能力跑通，保证架构正确、数据模型稳定。

包含内容：

- 新增 `web_monitor` source type
- 扩展 `SourceConfig` 和 `WebMonitorParserConfig`
- 实现 `WebMonitorParser`
- 实现 `WebMonitorSource`
- 接入现有 pipeline / registry
- 编写核心单元测试

重点验证：

- key fields 驱动 GUID
- 非 key 字段变化不触发新 item
- title / description / content template 渲染正确

### Phase 2：Custom Recipe 最小入口

目标：让高级用户可以在不新增复杂 UI 的前提下直接使用。

包含内容：

- 在 `Custom Recipe` 页面新增 `Web Monitor` source type
- 支持用户直接填写相应 `source_config` JSON
- 提供基础示例配置，降低使用门槛

这样可以最快把后端能力交付给真实用户验证。

### Phase 3：可视化配置体验增强

目标：降低使用门槛，提升配置体验。

包含内容：

- 复用 `html-to-rss` 页面已有的抓取和 DOM 预览能力
- 支持可视化添加多个 extractor
- 支持从 extractor 中多选 `key_fields`
- 提供 title / description / content 模板输入框
- 提示可用变量（如 `{{.price}}`, `{{.stock}}`, `{{.url}}`）

这一阶段主要是体验增强，不改变核心数据模型。

### Phase 4：后续增强（可选）

仅在真实使用反馈证明有必要时再考虑：

- 更丰富的模板 helper
- selector 命中预览 / 校验
- 配置向导中的实时渲染预览
- 更友好的错误提示

不建议在首版就引入：

- 独立 scheduler
- 历史状态表
- 独立通知系统

---

## 9. 验证方案

### 9.1 后端测试

重点覆盖：

- 相同 key values -> GUID 不变
- 非 key 字段变化 -> GUID 不变
- 任一 key 字段变化 -> GUID 改变
- `Extractors` / `KeyFields` 顺序不影响 GUID
- 文本会 `TrimSpace`
- 模板渲染正确
- 模板错误能返回明确 error
- selector 未命中时仍能产出结果

### 9.2 手动联调验证

- 通过 `Custom Recipe` 创建 `web_monitor` recipe
- 访问 `/recipe/:id`，确认能正常输出 RSS
- 修改测试页面中的 key 字段值，确认 GUID 改变
- 修改非 key 字段，确认 GUID 不变
- 验证 RSS content 是否能展示更丰富的模板化正文

---

## 9.5 可选的 AI 判定（扩展能力）

在核心“值驱动”模型之上，新增一个**可选的 AI 判定**步骤，用于把“原始文本变化”升级为“语义状态变化”。

### 9.5.1 设计动机

纯文本 key field 有时过于敏感：

- 库存文案可能在 `Sold Out` / `Currently unavailable` / `缺货` 之间反复横跳，但语义上都还是“不可购买”。
- 用户真正关心的是“能不能买”，而不是文案本身。

因此引入一个 AI 判定步骤：让 LLM 根据提取到的字段值产出一个**简短、稳定的判定标签**（如 `available` / `unavailable`），再把这个标签当作派生变量参与模板与 GUID。

### 9.5.2 配置模型

在 `WebMonitorParserConfig` 中新增可选字段：

```go
type WebMonitorAIJudgeConfig struct {
	Enabled     bool   `json:"enabled,omitempty"`
	Prompt      string `json:"prompt,omitempty"`        // 判定指令，启用时必填
	OutputField string `json:"output_field,omitempty"`  // 默认 "ai_verdict"
	Model       string `json:"model,omitempty"`         // 可选模型覆盖
}
```

### 9.5.3 运行流程

1. 按 extractors 提取所有字段值。
2. 若 `ai_judge.enabled`：
   - 把字段值按 key 排序后拼成确定性上下文，调用 LLM（温度 0）。
   - 结果做 `TrimSpace` 并仅保留首行，写入 `values[output_field]`。
3. AI 判定**先于** key field 校验执行，因此判定结果本身可以被选为 key field。
4. 用户通常把 `ai_verdict` 勾选为 key field：GUID 由语义判定驱动，仅当判定结果变化时 RSS Reader 才识别为新 item。

### 9.5.4 稳定性要点

- **必须复用 LLM 缓存**：判定调用复用 `adapter.CallLLMUsingContext`，按 prompt+context 哈希缓存。请求驱动模型下 RSS Reader 会反复轮询，缓存保证相同页面状态得到相同判定，避免 GUID 抖动造成通知风暴。
- **温度固定为 0**：约束模型输出确定。
- **prompt 约束输出**：要求只输出简短标签、相同输入必须给出相同结果。
- **判定失败直接报错**：与现有“配置错误尽早失败”风格一致，不静默 fallback。

### 9.5.5 前端入口

在 Web Monitor 向导第 2 步新增可折叠的「AI 判定（可选）」卡片：开关启用、填写判定指令、可自定义输出变量名与模型、一键勾选“将判定结果用作关键字段”。预览结果中以标签形式展示 AI 判定值。

## 10. 最终建议

建议这次实现优先完成：

- **Phase 1：核心后端能力**
- **Phase 2：Custom Recipe 最小入口**

原因：

- 架构改动最小
- 能快速验证真实需求
- 数据模型可以一次设计正确（`extractors + key_fields + content_template`）
- 不会因为提前做复杂 UI 或调度系统而放大范围

后续如果用户反馈良好，再进入可视化配置体验增强阶段。
