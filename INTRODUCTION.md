# FeedCraft 项目介绍

## 1. 项目主要用途

FeedCraft 是一个功能强大的 RSS 源处理中间件。它的核心理念是让用户能够在一个统一的平台上自定义和“加工 (Craft)”他们的 RSS 订阅源，从而获得更优质的阅读体验。

主要功能包括：

- **AI 智能增强**: 集成大语言模型 (LLM)，支持对 RSS 文章进行智能摘要生成、多语言翻译、以及基于自然语言的内容筛选（如自动过滤营销软文）。
- **内容深度提取**: 支持提取文章全文，并内置浏览器模拟功能 (Browserless)，能够处理需要 JavaScript 渲染的动态网页，解决 RSS 源只有摘要或无法直接抓取的问题。
- **网页转 RSS (HTML to RSS)**: 提供可视化的 HTML 转 RSS 生成器，能够将任意网页（如博客列表、新闻站点）转换为标准的 RSS 订阅源。
- **灵活的中间件架构**: 兼容所有标准的 RSS 阅读器，作为源站和阅读器之间的中间层运行。

## 2. 核心概念

理解以下三个概念是使用 FeedCraft 的关键：

### 2.1 Tool (工具, 原 CraftAtom)

这是处理 RSS 源的最小功能单元。每个 Tool 代表一种单一的处理逻辑。
**常见示例**:

- `fulltext`: 抓取并提取文章全文。
- `translate-title`: 使用 AI 翻译文章标题。
- `summary`: 使用 AI 生成文章摘要。
- `limit`: 限制输出的文章数量。
- `proxy`: 仅作为代理，不修改内容。

### 2.2 Blueprint (蓝图, 原 CraftFlow)

Blueprint 是由多个 Tool 按顺序串联而成的处理流水线。这允许用户定义复杂的处理逻辑。
**示例场景**:
你可以创建一个名为 "Clean-Read" 的 Blueprint，顺序执行以下操作：

1.  `fulltext` (获取全文)
2.  `ignore-advertorial` (AI 识别并移除软文)
3.  `translate-content` (翻译剩余的高质量文章)

### 2.3 Channel (频道, 原 Recipe)

Channel 是将“特定的 RSS 源”与“特定的处理逻辑 (Tool 或 Blueprint)”绑定在一起的持久化配置。
**作用**:
创建一个 Channel 后，系统会生成一个新的固定 RSS 订阅地址。你可以将这个地址添加到你的 RSS 阅读器中。当你访问这个地址时，FeedCraft 会自动拉取最新的源内容，并按照设定的处理逻辑进行加工。

## 3. 工作流程

FeedCraft 提供两种主要的使用模式，适应不同的使用场景：

### 3.1 便携模式 (Portable Mode) - 即用即走

适用于临时需求或简单的标准处理。无需在后台进行繁琐配置，只需拼接 URL 即可。

**使用方法**:
直接访问特定格式的 URL：
`https://{您的FeedCraft域名}/craft/{处理器名称}?input_url={原始RSS地址}`

**例如**:
想要翻译某个英文 RSS 的标题，只需在 URL 中指定 `translate-title` 工具即可获得翻译后的 RSS 流。

### 3.2 订阅模式 (Channel Mode) - 深度定制

适用于需要长期订阅、复杂处理或个性化定制的场景。

**使用步骤**:

1.  **进入管理后台**: 登录 FeedCraft 的 Web 管理界面。
2.  **定制处理器 (可选)**:
    - 创建新的 Tool：例如，自定义 AI 的 Prompt 提示词，让 AI 以特定的语气写摘要。
    - 编排 Blueprint：将多个工具组合起来。
3.  **创建 Channel**:
    - 输入原始 RSS 地址。
    - 选择要应用的处理器 (Tool 或 Blueprint)。
    - 保存并生成新的订阅地址。
4.  **订阅**: 将生成的 URL 添加到你的 RSS 阅读器中。

### 3.3 HTML to RSS 流程

对于没有提供 RSS 的网站：

1.  在 FeedCraft 中使用 "HTML to RSS" 工具。
2.  输入目标网页 URL。
3.  通过可视化选择器选取列表项、标题、链接等元素。
4.  生成 RSS 源地址，随后可直接接入上述的处理流程中。

##  technical-principles 4. 技术原理与代码逻辑

为了帮助开发者和高级用户更好地理解 FeedCraft 的工作方式，以下是其内部核心处理流程的技术解析。

### 4.1 数据流架构 (Pipeline Architecture)

FeedCraft 的处理流程遵循经典的 "Fetch-Parse-Process-Output" 管道模型：

1.  **Source (源获取)**:
    - 代码路径: `internal/source/`
    - 负责根据用户提供的 `input_url` 获取原始数据。
    - **Fetcher**: 支持 HTTP/HTTPS 请求，也支持通过 Browserless (Headless Chrome) 进行渲染以获取动态内容。
    - **Parser**: 将获取到的原始数据（XML, JSON, HTML）统一解析为标准化的 `gofeed.Feed` 内部对象。

2.  **Craft Processing (工艺处理)**:
    - 代码路径: `internal/craft/`
    - 这是 FeedCraft 的核心业务逻辑。系统获取到标准化的 Feed 对象后，会根据用户请求的处理器参数（可能是 Tool、Blueprint 或逗号分隔的列表），加载对应的 `CraftOption` 函数链。
    - 每个 `CraftOption` 都是一个闭包函数，它接收 Feed 对象，对其进行修改（如重写标题、替换内容、过滤条目），然后返回修改后的 Feed。

3.  **Output (输出)**:
    - 处理完成的 Feed 对象会被转换回标准的 RSS/Atom XML 格式，并设置正确的 Content-Type 返回给客户端。

### 4.2 核心组件代码逻辑

- **Tool Template (工具模板)**:
  - 定义在 `internal/craft/entry.go` 中。
  - 它是功能的“蓝图”，定义了某项功能（如翻译、摘要）的名称、参数结构 (`ParamTemplateDefine`) 以及如何生成处理函数 (`OptionFunc`)。
  - 系统启动时会注册所有内置的 System Tool Templates。

- **Tool (工具)**:
  - 这是 Tool Template 的具体实例。
  - 在数据库中，它存储了 Template 的名称以及具体的参数值（例如目标语言、摘要长度等）。
  - 当请求到来时，系统会查找对应的 Tool 配置，将其参数注入到 Template 的 `OptionFunc` 中，生成可以在运行时执行的逻辑。

- **Blueprint (蓝图)**:
  - 这是一个逻辑容器，它只存储了一组 Tool 的名称列表。
  - 执行时，系统会递归解析 Blueprint 中的每个 Tool，将它们展开为一个线性的 `CraftOption` 列表，按顺序对 Feed 进行处理。

### 4.3 缓存机制

为了优化性能并减少对上游源和 LLM API 的请求压力，FeedCraft 在多个层级实现了缓存：

- **内容缓存**: 对于耗时的操作（如 LLM 处理、全文提取），结果会基于内容的 Hash 值存储在 Redis 中 (`internal/craft/common.go`)。如果同一篇文章内容未变，将直接复用上次的处理结果。
- **Feed 缓存**: 处理完的最终 RSS 结果也会在短时间内缓存，避免频繁请求导致的资源浪费。