# Change Log

## [Unreleased] (since v2.1)

### ⚠️ 破坏性变更 (Breaking Changes)

*   **LLM 配置更新**: 重构了 LLM 集成，引入了通用的环境变量配置。
    *   新增 `FC_LLM_API_TYPE` (支持 `openai`, `ollama`)。
    *   `FC_OPENAI_ENDPOINT` 重命名为 `FC_LLM_API_BASE`。
    *   `FC_OPENAI_AUTH_KEY` 重命名为 `FC_LLM_API_KEY`。
    *   `FC_OPENAI_DEFAULT_MODEL` 重命名为 `FC_LLM_API_MODEL`。
    *   旧变量仍暂时兼容但有废弃警告。
*   **术语变更**: UI 和文档中的 "Craft Atom" 重命名为 "AtomCraft" (原子工艺), "Craft Flow" 重命名为 "FlowCraft" (组合工艺)。

### ✨ 新特性 (Features)

*   **RSS 生成器工具集**:
    *   新增 **HTML 转 RSS** 工具: 支持交互式选择器拾取、增强模式 (无头浏览器/Browserless)、富文本预览及智能选择逻辑。
    *   新增 **JSON 转 RSS** 工具: 支持通过 JQ 表达式从 JSON 源生成 RSS。
    *   新增 **搜索 转 RSS** 工具: 集成 SearXNG 和 LiteLLM，支持通过搜索结果生成 RSS。
    *   新增 **快速开始 (URL 生成器)**: 支持生成和解析 FeedCraft URL。
    *   Curl 转 RSS 支持配置 HTTP 方法和请求体。
*   **工艺组件 (Atom/Flow)**:
    *   新增 `time-limit` (时间限制) 原子工艺。
    *   新增 `beautify-content` (内容美化) 原子工艺。
    *   新增 `article-summary` (文章摘要) 原子工艺。
    *   新增 `immersive-translate` (沉浸式翻译) 组合工艺。
    *   新增通用 `llm-filter` (LLM 过滤器) 原子工艺。
    *   `fulltext-plus`: 增加 `wait` (等待时间) 和 `mode` (如 `networkidle2`) 参数以更好支持动态网页。
    *   支持 `DEFAULT_TARGET_LANG` 环境变量，用于控制翻译目标语言。
*   **用户界面与体验 (UI/UX)**:
    *   新增 **服务依赖状态** 页面，用于监控 SQLite, Redis, Browserless, LLM 等服务状态。
    *   新增 **搜索提供商设置** 页面。
    *   应用自定义 Arco Design 主题。
    *   重构 **Craft Flow 编辑器**: 采用列表式编辑，支持拖拽排序。
    *   改进 **Craft 选择器**: 模块化拆分，支持分类展示和多选。
    *   自定义配方编辑器: 支持 JSON 格式化、一键复制配置。
    *   添加关键操作的确认对话框 (如删除)。
*   **基础设施与后端**:
    *   支持 Ollama 作为 LLM 提供商。
    *   支持配置多个 LLM 模型并实现自动重试逻辑。
    *   优化 LLM 调用: 增加内容处理选项 (移除链接/图片) 以节省 Token。
    *   构建流程: 注入版本、提交哈希等元数据到二进制文件。
    *   新增 GitHub Actions CI 工作流。

### 🐛 问题修复 (Bug Fixes)

*   **HTML 转 RSS**:
    *   修复空响应导致的静默失败，优化错误处理。
    *   修复向导中 Axios 响应未正确解包的问题。
    *   优化 Fetch 逻辑，增加 User-Agent 和标准头以减少被拦截概率。
*   **搜索转 RSS**:
    *   修复生成失败时返回 200 状态码的问题 (现返回 500)。
    *   处理数据库读取配置失败的情况。
*   **系统与路由**:
    *   修复缺失的 API 路由返回 HTML 的问题 (现返回 404 JSON)。
    *   修复无效内存地址引用导致的 Panic。
    *   验证 Browserless 服务返回的 HTTP 状态码。
*   **其他**:
    *   修复 RSS 生成器 CSS 预览问题。
    *   修复 Docker 发布工作流中的 helper 错误。

### 📝 文档与杂项 (Documentation & Chores)

*   **文档**:
    *   新增关于搜索转 RSS、JSON 转 RSS、系统原子工艺的详细指南。
    *   更新快速开始和自定义配置文档。
    *   重构文档结构，迁移至 Astro Starlight。
*   **依赖与构建**:
    *   升级 Web 端 Vite 至 v5, TypeScript 至 v5。
    *   升级 Go 和 Node.js 依赖 (如 gorm, axios, vue-router 等)。
    *   更新 `.gitignore` 和 `Taskfile`。
