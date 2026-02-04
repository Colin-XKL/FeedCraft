---
title: 高级定制
sidebar:
  order: 1
---

对于高级用户，FeedCraft 提供了一个管理后台来定制 RSS 的处理流程。

## 访问后台

1.  使用 Docker 部署 FeedCraft（参考快速开始）。
2.  浏览器访问 `http://你的服务器IP:10088`。
3.  使用默认凭据登录：
    - 用户名：`admin`
    - 密码：`adminadmin`
      _(请登录后立即修改密码)_

## 核心概念

### 原子工艺 (AtomCraft)

**AtomCraft** 是最小的处理单元。除了内置的原子工艺（如 `translate-title`, `fulltext`），你可以基于模版创建自定义的原子工艺。

**示例：自定义翻译 Prompt**
你可以基于 `translate-content` 模版创建一个名为 `translate-to-french` 的新原子工艺，并在参数中填入自定义的 Prompt，指示 AI 将内容翻译成法语。

### 组合工艺 (FlowCraft)

**FlowCraft** 是多个 AtomCraft 的组合序列。这允许你将多个操作串联起来。

**示例：全文 + 摘要 + 翻译**
你可以定义一个名为 `digest-and-translate` 的组合工艺，包含以下步骤：

1.  `fulltext` (提取正文)
2.  `summary` (生成摘要)
3.  `translate-content` (翻译内容)

#### 管理组合工艺

你可以在后台的 **FlowCraft** 页面创建和管理组合工艺。
编辑器允许你添加原子工艺并安排它们的执行顺序。使用箭头按钮 (⬆️/⬇️) 调整顺序，或使用垃圾桶图标将其从流程中移除。

### 食谱 (Recipe)

**Recipe** 将特定的 RSS 源 URL 与某个 原子工艺 (AtomCraft) 或组合工艺 (FlowCraft) 绑定。这允许你创建一个持久化的、经过定制的订阅源 URL。

**管理食谱：**
在后台 **自定义食谱 (Custom Recipes)** 页面，你可以管理所有已创建的食谱。

- **创建 (Create)**：绑定新的 URL 和工艺。
- **预览 (Preview)**：点击预览按钮，直接在内置的 Feed Viewer 中查看生成的效果。
- **复制链接 (Copy Link)**：点击复制图标获取完整的订阅 URL。

**示例：**

- **输入 URL：** `https://news.ycombinator.com/rss`
- **处理器：** `digest-and-translate` (上面创建的工作流)
- **结果：** 你会得到一个新的 FeedCraft URL，订阅它即可获得带全文、摘要和翻译的 Hacker News。

## 搜索提供商配置 (Search Provider)

要使用 **搜索转 RSS (Search to RSS)** 功能，你必须配置搜索提供商。

在管理后台导航至 **设置 (Settings) > Search Provider**。

### 支持的提供商

- **LiteLLM / OpenAI Compatible**
  - **API URL**: 搜索服务的 API 端点（例如 `http://litellm-proxy:4000/v1/search`）。
  - **API Key**: 你的 API 密钥。（留空以保留现有密钥）
  - **Tool Name**: 特定函数调用工具名称（如果需要，例如某些 Agent 的 `google_search`）。工具名称将追加到 API URL 之后（例如 `.../v1/search/google_search`）。

- **SearXNG**
  - **API URL**: 你的 SearXNG 实例基础 URL（例如 `http://my-searxng.com`）。`/search` 路径会自动追加。
  - **Engines**: (可选) 逗号分隔的搜索引擎列表（例如 `google,bing`）。

> **提示：** 在保存之前，你可以使用 **检查连接 (Check Connection)** 按钮来验证与提供商的连接。

## 依赖服务 (Dependency Services)

**依赖服务** 仪表盘 (设置 (Settings) > 依赖服务状态 (Dependency Services)) 提供了所有连接的外部服务的健康检查概览。

它监控以下服务的状态：

- **SQLite**: 数据库连接。
- **Redis**: 缓存服务连接及延迟。
- **Browserless**: 无头浏览器服务可用性（全文提取功能必须）。
- **LLM Service**: 与配置的 AI 提供商的连接。
- **Search Provider**: 与配置的搜索引擎的连接。

如果“增强模式”或“全文提取”等功能出现故障，请使用此仪表盘排查连接问题。

你可以使用 **检查连接 (Check Connection)** 按钮来验证 FeedCraft 是否可以成功连接到配置的搜索提供商。

## 高级配置

### Docker 环境变量

你可以在 `docker-compose.yml` 中使用环境变量配置 FeedCraft。

- **FC_PUPPETEER_HTTP_ENDPOINT**: Browserless/Chrome 实例的地址。`fulltext-plus` 功能必须。
- **FC_REDIS_URI**: Redis 连接地址。用于缓存，加快处理速度并减少 AI Token 消耗。
- **FC_LLM_API_KEY**: OpenAI 或兼容服务（如 DeepSeek, Gemini 等）的 API Key。
- **FC_LLM_API_MODEL**: 默认使用的模型（如 `gemini-pro`, `gpt-3.5-turbo`）。**支持多个模型：** 你可以提供一个逗号分隔的模型列表（例如 `gpt-3.5-turbo,gpt-4`）。FeedCraft 会为每个请求随机选择一个模型，如果调用失败，会自动重试列表中的其他模型。
- **FC_LLM_API_BASE**: API 接口地址。如果是兼容 OpenAI 的 API，通常以 `/v1` 结尾。
- **FC_LLM_API_TYPE**: (可选) `openai` (默认) 或 `ollama`.

### 外部服务

为了发挥 FeedCraft 的全部功能，建议搭配 Redis 和 Browserless 部署。

```yaml
version: "3"
services:
  app.feed-craft:
    # ... (参考快速开始)
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000
      FC_REDIS_URI: redis://service.redis:6379/
      # ...

  service.redis:
    image: redis:6-alpine
    container_name: feedcraft_redis
    restart: always

  service.browserless:
    image: browserless/chrome
    container_name: feedcraft_browserless
    environment:
      USE_CHROME_STABLE: true
    restart: unless-stopped
```

服务默认监听在 80 端口，你也可以在同一网络下的其他容器中，使用 `http://app.feed-craft/xxx` 这样来进行访问(比如RSS 阅读器中通过这种方式来走内网通信订阅)。
