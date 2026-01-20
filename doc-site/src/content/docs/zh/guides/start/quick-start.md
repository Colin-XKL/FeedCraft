---
title: 快速开始
---

## 简介

FeedCraft 是一个强大的 RSS 源处理工具，作为中间件运行。你可以用它来翻译订阅源、提取全文、模拟浏览器渲染重 JS 网页、使用 Google Gemini 等 LLM 生成文章摘要、通过自然语言筛选 RSS 内容等！

## 便携模式 (Portable Mode)

你可以通过修改 RSS 源的 URL 来快速开始使用 FeedCraft，这被称为“便携模式”。

URL 格式为：
`https://feed-craft.colinx.one/craft/{craft_atom}?input_url={input_rss_url}`

其中：

- `{craft_atom}` 是你想要使用的处理步骤名称（原子工艺）。
- `{input_rss_url}` 是原始的 RSS 订阅源 URL。

**注意：** 如果你的 RSS 阅读器不会自动处理 URL 编码，你可能需要手动对 `{input_rss_url}` 进行 URL 编码。

### URL 生成器

为了方便使用，你可以利用 Web 界面中内置的 **URL 生成器** 来轻松构建这些 URL。

- 访问独立的生成器页面：`/start.html`
- 或者在管理后台使用 "URL Generator" 工具。

**URL 解析模式**
URL 生成器现在支持“解析模式”。你可以粘贴一个现有的 FeedCraft URL，工具会反向解析出它使用的工艺（Craft）、原始来源 URL 以及其他参数。这对于调试或理解复杂的 FeedCraft 链接非常有用。

### 常用原子工艺 (AtomCrafts)

以下是一些你可以直接使用的基础原子工艺：

- **proxy**: 简易 RSS 代理，不作任何处理。
- **limit**: 限制文章数量（默认最新 10 篇）。
- **fulltext**: 提取文章全文。
- **fulltext-plus**: 模拟浏览器渲染并提取全文（适用于重 JS 网站）。
- **introduction**: 调用 AI 为文章生成导读，附加在开头。
- **summary**: 调用 AI 总结文章主要内容。
- **translate-title**: 调用 AI 翻译文章标题。
- **translate-content**: 调用 AI 翻译文章内容（替换原文）。
- **translate-content-immersive**: 沉浸式翻译模式，每段原文后附加译文。
- **ignore-advertorial**: 使用 AI 筛选并过滤掉营销软文。

### 示例

1.  **翻译订阅源标题：**
    `https://feed-craft.colinx.one/craft/translate-title?input_url=https://feeds.feedburner.com/visualcapitalist`

2.  **获取订阅源全文：**
    `https://feed-craft.colinx.one/craft/fulltext?input_url=https://feeds.feedburner.com/visualcapitalist`

## 基础部署

你可以使用 Docker Compose 部署自己的实例。

### 最小化 `docker-compose.yml`

```yaml
version: "3"
services:
  app.feed-craft:
    image: ghcr.io/colin-xkl/feed-craft
    container_name: feed-craft
    restart: always
    ports:
      - "10088:80"
    volumes:
      - ./feed-craft-db:/usr/local/feed-craft/db
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000
      FC_REDIS_URI: redis://service.redis:6379/
      FC_LLM_API_KEY: skxxxxxx
      FC_LLM_API_MODEL: gemini-pro
      FC_LLM_API_BASE: https://xxxxxx
```

保存为 `docker-compose.yml` 并运行 `docker-compose up -d`。
访问 `http://localhost:10088` 进入后台管理界面。
默认账号：`admin` / `adminadmin`。
