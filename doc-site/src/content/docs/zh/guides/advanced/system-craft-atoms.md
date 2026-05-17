---
title: 系统内置 AtomCraft
description: FeedCraft 系统内置处理原子工艺 (AtomCrafts) 的详细参考指南。
---

FeedCraft 内置了一系列“原子工艺 (AtomCrafts)”，用于对订阅源进行特定的处理。你可以将这些原子工艺组合成“组合工艺 (FlowCraft)”来构建强大的数据管道。

## 内容获取与修复

这些原子主要用于获取全文或修复常见的订阅源问题。

### `fulltext` (全文提取)

从原始网页提取文章的全文内容。

- **适用场景:** 当 RSS 订阅源仅提供摘要或片段时。
- **机制:** 使用标准 HTTP 客户端请求网页，并使用算法提取正文。速度快且轻量。

### `fulltext-plus` (浏览器全文提取)

使用无头浏览器 (Puppeteer) 提取全文。

- **适用场景:** 针对通过 JavaScript 动态渲染内容或有较强反爬虫措施的网站。
- **机制:** 连接到配置的浏览器提供方（`browserless-restful` 或 `cdp`）来渲染页面。速度较慢但兼容性更强。
- **参数:**
  - `mode` (默认: `networkidle2`): 页面加载等待模式。
    - `load`: 等待 `load` 事件。
    - `domcontentloaded`: 等待 `DOMContentLoaded` 事件。
    - `networkidle0`: 等待直到 500ms 内没有活跃的网络连接。
    - `networkidle2`: 等待直到 500ms 内活跃的网络连接数不超过 **2** 个。(推荐用于 SPA 单页应用)。
  - `wait` (默认: `0`): 显式等待时间（秒），例如 `5`。

### `proxy` (代理)

简单的订阅源代理。

- **适用场景:** 当你只想转发原始 Feed 而不做修改，或者将 FeedCraft 作为中心网关使用时。

### `guid-fix` (GUID 修复)

使用文章内容的 FNV-1a 哈希值替换 RSS 条目的 GUID。

- **适用场景:** 某些订阅源在内容未变更的情况下频繁更改 GUID，导致阅读器中出现重复的未读条目。此原子可基于内容稳定 GUID。

### `relative-link-fix` (相对链接修复)

将内容中的相对链接（如 `<a href="/about">`）转换为绝对链接（如 `<a href="https://example.com/about">`）。

- **适用场景:** 提取全文后必不可少，否则在 RSS 阅读器中查看时链接会失效。

### `cleanup` (HTML 清理)

清理 HTML 内容以去除杂乱信息。

- **适用场景:** 通过移除多余的 class、style 和空标签来提高可读性。

---

## 过滤类 (Filtering)

控制哪些条目可以进入最终生成的 Feed。

### `limit` (数量限制)

限制 Feed 中的条目数量。

- **参数:**
  - `num` (默认: `10`): 保留的最大条目数。

### `time-limit` (时间限制)

过滤掉超过指定天数的条目。

- **参数:**
  - `days` (默认: `7`): 文章保留的最大天数。

### `keyword` (关键词过滤)

根据标题或内容中的关键词进行过滤。

- **参数:**
  - `keywords`: 逗号分隔的关键词列表（子串匹配，区分大小写）。例如：`ad,sell,SALE`。
  - `mode`: `include` (保留匹配项，默认) 或 `exclude` (移除匹配项)。
  - `scope`: `title` (标题), `content` (内容), 或 `all` (全部，默认)。

---

## AI 增强 (AI Enhancement)

使用大语言模型 (LLM) 来转换和丰富你的内容。

:::note
使用此类原子需要在环境变量中配置 LLM (API Key, Base URL 等)。
:::

### `translate-title` (标题翻译)

将文章标题翻译为你的目标语言。

- **参数:**
  - `prompt`: 自定义提示词。默认使用标准翻译提示词。支持 `{{.TargetLang}}` 占位符。

### `translate-content` (内容翻译)

翻译整篇文章内容，替换原文。

- **参数:**
  - `prompt`: 自定义提示词。支持 `{{.TargetLang}}`。

### `translate-content-immersive` (沉浸式翻译)

双语翻译模式。在每一段原文后面追加翻译后的内容。

- **参数:**
  - `prompt`: 自定义提示词。

### `summary` (AI 摘要)

生成文章摘要并将其添加到正文开头。

- **参数:**
  - `prompt`: 用于生成摘要的自定义提示词。

### `introduction` (AI 导读)

为文章生成简短的介绍或导语。

- **参数:**
  - `prompt`: 自定义提示词。

### `beautify-content` (智能排版)

使用 LLM 重新格式化文章，修复排版错误，去除广告，并标准化 Markdown 格式，最后转换回干净的 HTML。

- **参数:**
  - `prompt`: 设定“编辑”角色的指令。

---

## AI 过滤 (AI Filtering)

利用语义理解进行高级过滤。

### `ignore-advertorial` (软文过滤)

使用 LLM 检测文章是否为软文或广告（综合评估标题和内容），并将其移除。

- **参数:**
  - `prompt-for-exclude`: 如果文章是广告，应返回 `true` 的提示词。

### `llm-filter` (通用 LLM 过滤)

通用的 LLM 过滤器。你可以定义**排除**条件。LLM 会结合文章标题和内容进行评估。

- **参数:**
  - `filter_condition`: 自然语言描述的条件。如果 LLM 回答 "yes" (true)，则该条目会被**移除**。
  - _示例:_ "这篇文章是关于体育的吗？" (移除体育类文章)。

### `embedding-filter` (Embedding 语义过滤)

基于 Embedding 模型的语义主题过滤器。它不会让聊天模型逐条判断文章，而是把“主题锚点”和文章内容都转换成向量，再用余弦相似度判断是否匹配。

:::tip
当你需要快速、稳定地做主题过滤时，优先使用 `embedding-filter`，例如“只保留 AI 基础设施新闻”或“移除体育文章”。如果规则需要复杂推理、政策判断或结构化决策，再使用 `llm-filter`。
:::

#### 环境变量

推荐单独配置 Embedding 服务：

```bash
FC_EMBEDDING_API_TYPE=openai
FC_EMBEDDING_API_BASE=https://api.openai.com/v1
FC_EMBEDDING_API_KEY=sk-your-api-key
FC_EMBEDDING_API_MODEL=text-embedding-3-small
FC_EMBEDDING_BATCH_SIZE=5
FC_EMBEDDING_MAX_INPUT_CHARS=8000
```

`FC_EMBEDDING_API_TYPE` 支持：

- `openai`: OpenAI 或 OpenAI 兼容的 Embedding 接口。
- `gemini`: 通过 Gemini 的 OpenAI 兼容 Embedding 接口调用。请显式设置 `FC_EMBEDDING_API_BASE` 和 `FC_EMBEDDING_API_MODEL`。
- `ollama`: 本地 Ollama Embedding 模型。请设置 `FC_EMBEDDING_API_BASE`，例如 `http://localhost:11434`，并使用 `nomic-embed-text` 或 `bge-m3` 这类 Embedding 模型。

如果 `FC_EMBEDDING_API_TYPE`、`FC_EMBEDDING_API_BASE`、`FC_EMBEDDING_API_KEY` 都没有设置，FeedCraft 会回退到对应的 `FC_LLM_API_TYPE`、`FC_LLM_API_BASE`、`FC_LLM_API_KEY`。Embedding 模型名是独立的：请将 `FC_EMBEDDING_API_MODEL` 设置为真正的 Embedding 模型；如果 API 类型是 `openai` 且未设置该变量，FeedCraft 会使用默认 Embedding 模型。FeedCraft 不会复用 `FC_LLM_API_MODEL`，因为它通常是聊天模型。

`FC_EMBEDDING_MAX_INPUT_CHARS` 是发送给 Embedding 服务前的最终安全上限，包含 `instruction` 前缀。它是字符预算，不是精确的 tokenizer token 数。建议按模型 token 窗口设置保守值，例如 8k token 的 Embedding 模型可从 `8000` 开始。

#### 参数

- `anchors` (**必填**): 每行一个主题锚点。锚点应描述你想匹配的主题，例如：

  ```text
  人工智能基础设施
  机器学习研究
  大语言模型部署
  ```

- `threshold` (默认: `0.6`): 余弦相似度阈值，范围 `0` 到 `1`。值越高越严格。建议从 `0.6` 开始；漏掉相关文章时调低，混入无关文章时调高。
- `mode` (默认: `include`): `include` 保留匹配项；`exclude` 移除匹配项。
- `max_content_length` (默认: `2000`): 当前 AtomCraft 使用的文章正文最大字符数；最终发送前还会受 `FC_EMBEDDING_MAX_INPUT_CHARS` 保护。
- `instruction` (可选): 会作为文本前缀拼接到每条 Embedding 输入前。除非你的模型确实需要固定任务前缀，否则建议留空。

#### 管理后台使用流程

1. 打开 **工作台 → AtomCraft**。
2. 新建一个 AtomCraft，例如 `ai-news-only`。
3. 模板选择 `embedding-filter`。
4. 在 `anchors` 中每行填写一个主题。
5. 使用 `mode=include` 保留匹配文章，或切换到 `exclude` 移除匹配文章。
6. 保存 AtomCraft。
7. 在 FlowCraft、Recipe、Feed Compare 中使用它，或直接访问：

```text
/craft/ai-news-only?input_url=https%3A%2F%2Fexample.com%2Ffeed.xml
```

#### 常见问题

- **"anchors parameter is required"**: `anchors` 至少需要一行非空内容。
- **"FC_EMBEDDING_API_MODEL must be set"**: 请配置单个 Embedding 模型。聊天模型不适合这个功能。
- **所有文章都被移除**: 降低 `threshold`，增加更宽泛的锚点，或增大 `max_content_length`。
- **无关文章仍然保留**: 提高 `threshold`，让锚点更具体，或用代表性短语替代宽泛分类名。
- **服务商提示输入过长**: 降低 `FC_EMBEDDING_MAX_INPUT_CHARS`，降低 `max_content_length`，或缩短 `instruction`。
