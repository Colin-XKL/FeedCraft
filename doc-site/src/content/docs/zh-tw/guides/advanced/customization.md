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

## 故障排除與工具

FeedCraft 提供了幾個內建工具來幫助你除錯和驗證 RSS 來源。

### RSS 預覽 (Feed Viewer)

位於 **實用工具 (Tools) > RSS 預覽 (RSS Viewer)**，此工具允許你：

- 直接在管理後台預覽任何 RSS 來源。
- 驗證 FeedCraft 伺服器是否可以存取該來源（用於除錯網路問題）。
- 檢查系統如何解析該來源的內容。

### Feed 對比 (Feed Compare)

位於 **實用工具 (Tools) > Feed 對比 (Feed Compare)**，此工具幫助你視覺化 AtomCraft 的效果。

- 輸入原始 RSS 來源地址。
- 選擇特定的 AtomCraft（例如 `translate-title` 或 `summary`）。
- 查看原始文章與處理後文章的並排對比。
- 這非常適合在建立永久 Recipe 之前測試新的 Prompt 或過濾邏輯。

### Craft 依賴檢查 (System Health)

位於 **實用工具 (Tools) > Craft 依賴檢查 (System Health)**，此視覺化工具映射出你的 Recipe、FlowCraft 和 AtomCraft 之間的內部關係。

- **目的：** 檢測損壞的引用（例如，Recipe 引用了已刪除的 FlowCraft）或循環依賴。
- **紅色節點：** 表示需要修復的缺失組件。
