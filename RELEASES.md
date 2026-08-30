# FeedCraft Releases

> 面向读者的人性化发布说明。完整技术变更（逐条 commit）见 [CHANGELOG.md](./CHANGELOG.md)（由 release-please 自动生成；下一次正式发版后才会出现）。

本文件**只由人工维护**。release-please 不会改写它。每个版本在 merge release PR **之前**补一段「亮点 / 重要变更 / 升级注意」。

## v3.2.0 (2026-08-30)

尚未打 tag；内容来自当前 `dev` 相对 v3.1.0 的人工汇总。

### 亮点

- **收件箱 (Inbox)**：脚本、自动化或第三方平台可以把文章 POST 给 FeedCraft，再作为 RSS 订阅。适合没有 RSS 的来源（webhook、私有数据）；配套系统授权令牌、条目上限与过期清理。
- **主题订阅 (Topic Feed)**：把多个来源合成一个订阅。支持按标题、SimHash（字面近似）、Embedding（语义）跨源去重；分步向导可为每个输入源加备注并即时预览；子源抓取失败时用缓存兜底。
- **新原子工艺**：`ai-filter`（自然语言筛选）、`ai-content-process`（自定义提示词加工正文）、`re-title`（按正文重写标题）、`link-flatten`（把文内链接抽成独立条目）。
- **网页监控**：为没有 RSS 的网页建订阅，只在关注字段（价格、版本号等）变化时推送。
- **预览体验**：RSS 预览支持 HTML 渲染、卡片展开与详情弹窗；Feed 对比并入同一页面。浏览器服务改为可配置 Provider，不再限定 Browserless。

### 重要变更

- 单 feed 内文章级 LLM 调用进入全局调度器并发执行，处理延迟明显下降；缓存未命中时的同 key 请求会合并（singleflight）。
- RSS 生成增强：HTML 转 RSS 支持抓取前点击/翻页等导航，并可指定条目图标；JSON 转 RSS 支持 `$` 模板。
- 可限制单次从源读取的条目数，大源不必每次全量处理。
- `fulltext-plus`：单篇渲染超时不再拖垮整个 Feed；浏览器取不到内容时回退普通全文提取。
- 上游抓取失败返回更准确的状态码；配方名重复给出明确提示。
- AtomCraft 创建后不再允许改名，避免 FlowCraft / 配方引用失效。

### 升级注意

- **破坏性变更**：Embedding 模型必须显式配置 `FC_EMBEDDING_API_MODEL`，不再提供默认值。
- 升级后需重启容器。
- 无数据库迁移要求。

## v3.1.0 (2026-03-21)

### 亮点

- LLM 并发控制，以及 fulltext 类工艺的域名级限流，大源更不容易把上游打满。
- Search-to-RSS 增强模式；RSS 向导用 `limax` 自动生成 recipe ID。
- 新增系统健康检查页，错误提示更可读。

### 重要变更

- 修复 `priority dispatcher` 潜在死锁。
- 修复搜索提供商配置中的 jq 表达式、API key 清除、读取已有配置时的错误处理，以及 active check 超时。

### 升级注意

- 无破坏性变更；升级后需重启容器。

## v3.0.0 (2026-01-19)

### 亮点

- 内置 RSS 生成器工具集：HTML / JSON / 搜索转 RSS，以及快速开始 URL 生成器。
- 多项新工艺：`time-limit`、`beautify-content`、`article-summary`、`immersive-translate`、通用 `llm-filter`；`fulltext-plus` 增加 `wait` 与 `mode`。
- 支持 Ollama；可配置多个 LLM 模型并自动重试。
- 文档站点迁到 Astro Starlight，并增加繁体中文。

### 重要变更

- 服务依赖状态页、搜索提供商设置页；Craft Flow 编辑器改为列表拖拽排序。
- UI / 文档术语：Craft Atom → AtomCraft（原子工艺），Craft Flow → FlowCraft（组合工艺）。

### 升级注意

- **破坏性变更（LLM 环境变量）**：
  - 新增 `FC_LLM_API_TYPE`（`openai` / `ollama`）
  - `FC_OPENAI_ENDPOINT` → `FC_LLM_API_BASE`
  - `FC_OPENAI_AUTH_KEY` → `FC_LLM_API_KEY`
  - `FC_OPENAI_DEFAULT_MODEL` → `FC_LLM_API_MODEL`
  - 旧变量仍暂时兼容，但会打废弃警告，请尽快迁移。
- 可用 `FC_DEFAULT_TARGET_LANG` 控制翻译默认目标语言。
- 升级后需重启容器。
