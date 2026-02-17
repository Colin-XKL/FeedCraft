---
title: 進階定制
sidebar:
  order: 1
---

對於進階用戶，FeedCraft 提供了一個管理後台來定制 RSS 的處理流程。

## 訪問後台

1.  使用 Docker 部署 FeedCraft（參考快速開始）。
2.  瀏覽器訪問 `http://你的伺服器IP:10088`。
3.  使用預設憑據登入：
    - 用戶名：`admin`
    - 密碼：`adminadmin`
      _(請登入後立即修改密碼)_

## 搜尋供應商配置 (Search Provider)

要使用 **搜尋轉 RSS (Search to RSS)** 功能，你必須配置搜尋供應商。

在管理後台導航至 **設定 (Settings) > Search Provider**。

### 支援的供應商

- **LiteLLM / OpenAI Compatible**
  - **API URL**: 搜尋服務的 API 端點（例如 `http://litellm-proxy:4000/v1/search`）。
  - **API Key**: 你的 API 金鑰。（留空以保留現有金鑰）
  - **Tool Name**: 特定函式調用工具名稱（如果需要，例如某些 Agent 的 `google_search`）。工具名稱將追加到 API URL 之後（例如 `.../v1/search/google_search`）。

- **SearXNG**
  - **API URL**: 你的 SearXNG 實例基礎 URL（例如 `http://my-searxng.com`）。`/search` 路徑會自動追加。
  - **Engines**: (可選) 逗號分隔的搜尋引擎列表（例如 `google,bing`）。

> **提示：** 在保存之前，你可以使用 **檢查連線 (Check Connection)** 按鈕來驗證與供應商的連線。

## 依賴服務 (Dependency Services)

**依賴服務** 儀表板 (設定 (Settings) > 依賴服務狀態 (Dependency Services)) 提供了所有連線的外部服務的健康檢查概覽。

它監控以下服務的狀態：

- **SQLite**: 資料庫連線。
- **Redis**: 快取服務連線及延遲。
- **Browserless**: 無頭瀏覽器服務可用性（全文提取功能必須）。
- **LLM Service**: 與配置的 AI 供應商的連線。
- **Search Provider**: 與配置的搜尋引擎的連線。

如果「增強模式」或「全文提取」等功能出現故障，請使用此儀表板排查連線問題。

你可以使用 **檢查連線 (Check Connection)** 按鈕來驗證 FeedCraft 是否可以成功連線到配置的搜尋供應商。

## 故障排除與工具

FeedCraft 提供了幾個內建工具來幫助你除錯和分析你的 RSS 源。這些工具可以從 **工具 (Tools)** 選單訪問。

### Feed Compare (RSS源比較工具)

**Feed Compare** 工具允許你視覺化對比原始 RSS 源和應用 AtomCraft 或 FlowCraft 後的處理結果。

- **如何使用:**
  1.  輸入原始 RSS URL。
  2.  選擇你想要測試的 Craft。
  3.  點擊 **Compare**。
- **輸出:** 該工具會並排顯示 Feed 條目的對比，突顯顯示標題、內容和其他欄位的變化。

### RSS Viewer (RSS 預覽)

**RSS Viewer** 是一個簡單的工具，用於直接在管理後台預覽任何 RSS 源的內容。

- **如何使用:** 輸入 RSS URL 並點擊 **View**。
- **使用場景:** 在將 Feed 添加到閱讀器之前，驗證其是否可訪問且格式正確。

### Craft Dependencies (Craft 依賴檢查)

**Craft Dependencies** (系統健康) 儀表板視覺化顯示了 FeedCraft 實例的內部拓撲結構。

- **目的:** 它展示了不同組件（AtomCrafts, FlowCrafts, Recipes）之間是如何相互依賴的。
- **使用場景:** 使用此工具了解修改特定 AtomCraft 的影響，或除錯循環依賴。

### Debug Tools (除錯工具)

對於進階除錯，可使用以下工具：

- **LLM Debug:** 直接測試你的 LLM 配置和 Prompt。
- **Ad Check Debug:** 針對特定內容測試 "ignore-advertorial" 邏輯，查看條目為何被標記為廣告。

## 進階配置

### Docker 環境變數

你可以在 `docker-compose.yml` 中使用環境變數配置 FeedCraft。

- **FC_PUPPETEER_HTTP_ENDPOINT**: Browserless/Chrome 實例的地址。`fulltext-plus` 功能必須。
- **FC_REDIS_URI**: Redis 連線地址。用於快取，加快處理速度並減少 AI Token 消耗。
- **FC_LLM_API_KEY**: OpenAI 或相容服務（如 DeepSeek, Gemini 等）的 API Key。
- **FC_LLM_API_MODEL**: 預設使用的模型（如 `gemini-pro`, `gpt-3.5-turbo`）。**支援多個模型：** 你可以提供一個逗號分隔的模型列表（例如 `gpt-3.5-turbo,gpt-4`）。FeedCraft 會為每個請求隨機選擇一個模型，如果調用失敗，會自動重試列表中的其他模型。
- **FC_LLM_API_BASE**: API 介面地址。如果是相容 OpenAI 的 API，通常以 `/v1` 結尾。
- **FC_LLM_API_TYPE**: (可選) `openai` (預設) 或 `ollama`.

### 外部服務

為了發揮 FeedCraft 的全部功能，建議搭配 Redis 和 Browserless 部署。

```yaml
version: "3"
services:
  app.feed-craft:
    # ... (參考快速開始)
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

服務預設監聽在 80 埠，你也可以在同一網絡下的其他容器中，使用 `http://app.feed-craft/xxx` 這樣來進行訪問(比如 RSS 閱讀器中透過這種方式來走內網通訊訂閱)。
