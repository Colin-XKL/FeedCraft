---
title: 快速開始
sidebar:
  order: 1
---

## 簡介

FeedCraft 是一個強大的 RSS 訂閱源處理工具，作為中介軟體 (Middleware) 運行。你可以用它來翻譯訂閱源、提取全文、模擬瀏覽器渲染重 JS 網頁、使用 Google Gemini 等 LLM 生成文章摘要、透過自然語言篩選 RSS 內容等！

## 便攜模式 (Portable Mode)

你可以透過修改 RSS 訂閱源的 URL 來快速開始使用 FeedCraft，這被稱為「便攜模式」。

URL 格式為：
`https://feed-craft.colinx.one/craft/{craft_atom}?input_url={input_rss_url}`

其中：

- `{craft_atom}` 是你想要使用的處理步驟名稱（原子工藝）。你可以在 [系統內建 AtomCraft](../../advanced/system-craft-atoms) 指南中查閱完整列表。
- `{input_rss_url}` 是原始的 RSS 訂閱源 URL。

為了更深入地了解 FeedCraft 的工作原理，請參閱[核心概念](../concepts)指南。

**注意：** 如果你的 RSS 閱讀器不會自動處理 URL 編碼，你可能需要手動對 `{input_rss_url}` 進行 URL 編碼。

### URL 生成器

為了方便使用，你可以利用 Web 介面中內建的 **URL 生成器** 來輕鬆構建這些 URL。

- 訪問獨立的生成器頁面：`/start.html` ([公共實例地址](https://feed-craft.colinx.one/start.html))
- 或者在管理後台使用 "URL Generator" 工具 (位於 **儀表板 (Dashboard) > 快速開始 (Quick Start)**)。

**URL 解析模式**
URL 生成器現在支援「解析模式」。你可以貼上一個現有的 FeedCraft URL，工具會反向解析出它使用的工藝（Craft）、原始來源 URL 以及其他參數。這對於除錯或理解複雜的 FeedCraft 連結非常有用。

### 常用原子工藝 (AtomCrafts)

以下是一些你可以直接使用的基礎原子工藝：

- **proxy**: 簡易 RSS 代理，不作任何處理。
- **limit**: 限制文章數量（預設最新 10 篇）。
- **fulltext**: 提取文章全文。
- **fulltext-plus**: 模擬瀏覽器渲染並提取全文（適用於重 JS 網站）。
- **introduction**: 調用 AI 為文章生成導讀，附加在開頭。
- **summary**: 調用 AI 總結文章主要內容。
- **translate-title**: 調用 AI 翻譯文章標題。
- **translate-content**: 調用 AI 翻譯文章內容（替換原文）。
- **translate-content-immersive**: 沉浸式翻譯模式，每段原文後附加譯文。
- **ignore-advertorial**: 使用 AI 篩選並過濾掉行銷軟文。

### 示例

1.  **翻譯訂閱源標題：**
    `https://feed-craft.colinx.one/craft/translate-title?input_url=https://feeds.feedburner.com/visualcapitalist`

2.  **獲取訂閱源全文：**
    `https://feed-craft.colinx.one/craft/fulltext?input_url=https://feeds.feedburner.com/visualcapitalist`

## 基礎部署

你可以使用 Docker Compose 部署自己的實例。

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

保存為 `docker-compose.yml` 並執行 `docker-compose up -d`。
訪問 `http://localhost:10088` 進入後台管理介面。
預設帳號：`admin` / `adminadmin`。
