# Change Log

## [v3.2.0] (since v3.1.0)

### ⚠️ 破坏性变更 (Breaking Changes)

- Embedding 模型必须显式配置 `FC_EMBEDDING_API_MODEL`，不再提供默认值

### ✨ 新特性 (Features)

- **收件箱 (Inbox)**：脚本、自动化工具或第三方平台可以直接把文章 POST 给 FeedCraft，再作为 RSS 订阅。适合给本身没有 RSS 的来源（如 webhook 通知、私有数据）建一个订阅源；配套系统授权令牌、条目上限与过期清理
- **主题订阅 (Topic Feed)**：正式开放入口，把多个来源合成一个订阅。新增按标题、SimHash（字面近似）、Embedding（语义）三种跨源去重，不同媒体报道同一件事只看到一条；分步向导可为每个输入源加备注并即时预览效果，子源抓取失败时用缓存兜底而不是整个主题变空
- **新增原子工艺**：
  - `ai-filter`：用自然语言写筛选规则，让 LLM 决定文章留还是丢
  - `ai-content-process`：用自己的提示词加工正文，结果可插到文首、文末或替换原文
  - `re-title`：让 LLM 根据正文重写标题，改善标题党或含义不明的标题
  - `link-flatten`：把文章里的链接抽出来变成独立条目，适合订阅周报、导航页这类“链接列表”式内容
- **网页监控**：为没有 RSS 的网页建订阅，只在关注的字段（如价格、版本号）变化时才推送更新
- **RSS 生成增强**：HTML 转 RSS 支持抓取前先执行点击、翻页等导航动作，能取到需要交互才显示的内容，并可指定条目图标来源；JSON 转 RSS 支持 `$` 模板写法
- **预览体验**：RSS 预览支持 HTML 渲染模式、卡片展开/收起与详情弹窗，Feed 对比也并入同一页面，调 Craft 时可直接看到处理前后差异；浏览器服务改为可配置 Provider，不再限定 Browserless
- **输入条目上限**：可限制单次从源读取的条目数，订阅历史很长的大源时不必每次全量处理，明显减少耗时与 LLM 花费

### 🐛 问题修复 (Bug Fixes)

- `fulltext-plus`：个别文章渲染超时不再让整个 Feed 一起失败，浏览器取不到内容时自动回退普通全文提取
- 上游源抓取失败时返回更准确的状态码（502/504/422）；配方名重复时给出明确提示而非数据库错误
- 大量文章同时调用 LLM 时排队更平稳，正文里的 base64 图片会先清理以减少 token 消耗，缓存未命中时的重复请求会被合并
- 修复网页监控若干问题、Markdown/HTML 转换丢失段落换行、含非法 XML 字符的源解析失败、HTML 转 RSS 无法抓取内网源且错误提示不可读
- AtomCraft 创建后不再允许改名，避免引用它的 FlowCraft 和配方失效

### 📝 文档与杂项 (Documentation & Chores)

- Inbox / Topic / AI 内容处理文档（en / zh-CN / zh-TW）；开发环境指南
- 前端 TypeScript 工具链升级；CI 升级 Node 24，适配 Go 1.25 与 golangci-lint

## [v3.1.0] (since v3.0.0)

### ✨ 新特性 (Features & Refactors)

- **调度与并发优化 (Dispatcher & Concurrency)**:
  - 增加 LLM 并发控制逻辑 (concurrency control)
  - fulltext crafts 类新增 domain 级别 rate limiting
- **搜索与 RSS (Search & RSS)**:
  - Search-to-RSS: 增加 enhanced mode
  - RSS wizards: 使用 `limax` 自动生成 recipe ID
- **系统与用户体验**:
  - 增加 System Health Check 页面
  - 改进错误提示用户体验 (better error ux)

### 🐛 问题修复 (Bug Fixes)

- 修复 `priority dispatcher` 中的潜在死锁问题
- 修复 search provider configs 中的 jq expressions 错误
- 允许清除 search provider API key
- settings: 修复读取现有 search provider config 时的错误处理
- 修复 search provider active check 并增加 timeout
- monitor: 修复并暴露 search provider check 中的 db errors

### 📝 文档与杂项 (Documentation & Chores)

- **文档 (Docs)**:
  - 增加系统工具的文档 (viewer, compare, health)
  - 增加比较 FeedCraft 与其他工具的文档
  - 增加繁体中文文档 (zh-tw doc)
  - 更新 minimal docker compose example, theme 和 badge

## [v3.0.0] (since v2.1)

### ⚠️ 破坏性变更 (Breaking Changes)

- **LLM 配置更新**: 重构了 LLM 集成，引入了通用的环境变量配置
  - 新增 `FC_LLM_API_TYPE` (支持 `openai`, `ollama`)
  - `FC_OPENAI_ENDPOINT` 重命名为 `FC_LLM_API_BASE`
  - `FC_OPENAI_AUTH_KEY` 重命名为 `FC_LLM_API_KEY`
  - `FC_OPENAI_DEFAULT_MODEL` 重命名为 `FC_LLM_API_MODEL`
  - 旧变量仍暂时兼容但有废弃警告
- **术语变更**: UI 和文档中的 "Craft Atom" 重命名为 "AtomCraft" (原子工艺), "Craft Flow" 重命名为 "FlowCraft" (组合工艺)

### ✨ 新特性 (Features)

- **RSS 生成器工具集**:
  - 新增 **HTML 转 RSS** 工具: 支持交互式选择器拾取、增强模式 (无头浏览器/Browserless)、富文本预览及智能选择逻辑
  - 新增 **JSON 转 RSS** 工具: 支持通过 JQ 表达式从 JSON 源生成 RSS
  - 新增 **搜索 转 RSS** 工具: 集成 SearXNG 和 LiteLLM，支持通过搜索结果生成 RSS
  - 新增 **快速开始 (URL 生成器)**: 支持生成和解析 FeedCraft URL
  - Curl 转 RSS 支持配置 HTTP 方法和请求体
- **工艺组件 (Atom/Flow)**:
  - 新增 `time-limit` (时间限制) 原子工艺
  - 新增 `beautify-content` (内容美化) 原子工艺
  - 新增 `article-summary` (文章摘要) 原子工艺
  - 新增 `immersive-translate` (沉浸式翻译) 组合工艺
  - 新增通用 `llm-filter` (LLM 过滤器) 原子工艺
  - `fulltext-plus`: 增加 `wait` (等待时间) 和 `mode` (如 `networkidle2`) 参数以更好支持动态网页
  - 支持 `DEFAULT_TARGET_LANG` 环境变量，用于控制翻译目标语言
- **用户界面与体验 (UI/UX)**:
  - 新增 **服务依赖状态** 页面，用于监控 SQLite, Redis, Browserless, LLM 等服务状态
  - 新增 **搜索提供商设置** 页面
  - 应用自定义 Arco Design 主题
  - 重构 **Craft Flow 编辑器**: 采用列表式编辑，支持拖拽排序
  - 改进 **Craft 选择器**: 模块化拆分，支持分类展示和多选
  - 自定义配方编辑器: 支持 JSON 格式化、一键复制配置
  - 添加关键操作的确认对话框 (如删除)
- **基础设施与后端**:
  - 支持 Ollama 作为 LLM 提供商
  - 支持配置多个 LLM 模型并实现自动重试逻辑
  - 优化 LLM 调用: 增加内容处理选项 (移除链接/图片) 以节省 Token
  - 构建流程: 注入版本、提交哈希等元数据到二进制文件
  - 新增 GitHub Actions CI 工作流

### 🐛 问题修复 (Bug Fixes)

- **HTML 转 RSS**:
  - 修复空响应导致的静默失败，优化错误处理
  - 修复向导中 Axios 响应未正确解包的问题
  - 优化 Fetch 逻辑，增加 User-Agent 和标准头以减少被拦截概率
- **搜索转 RSS**:
  - 修复生成失败时返回 200 状态码的问题 (现返回 500)
  - 处理数据库读取配置失败的情况
- **系统与路由**:
  - 修复缺失的 API 路由返回 HTML 的问题 (现返回 404 JSON)
  - 修复无效内存地址引用导致的 Panic
  - 验证 Browserless 服务返回的 HTTP 状态码
- **其他**:
  - 修复 RSS 生成器 CSS 预览问题
  - 修复 Docker 发布工作流中的 helper 错误

### 📝 文档与杂项 (Documentation & Chores)

- **文档**:
  - 新增关于搜索转 RSS、JSON 转 RSS、系统原子工艺的详细指南
  - 更新快速开始和自定义配置文档
  - 重构文档结构，迁移至 Astro Starlight
- **依赖与构建**:
  - 升级 Web 端 Vite 至 v5, TypeScript 至 v5
  - 升级 Go 和 Node.js 依赖 (如 gorm, axios, vue-router 等)
  - 更新 `.gitignore` 和 `Taskfile`
