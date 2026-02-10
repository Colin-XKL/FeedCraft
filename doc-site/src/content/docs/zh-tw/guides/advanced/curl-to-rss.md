---
title: 從 CURL 語句生成 RSS
description: 使用 jq 選取器將任意 JSON API 響應轉換為 RSS 訂閱源。
sidebar:
  order: 3
  badge:
    text: new
    variant: success
---

FeedCraft 包含一個 **從 CURL 語句生成 RSS (CURL to RSS)** 工具，允許你從 JSON API 獲取資料並使用 `jq` 選取器將其轉換為 RSS 訂閱源。

## 概覽

JSON RSS 生成器可以幫助你：

1.  **抓取 (Fetch)**：從 API 端點抓取 JSON 資料（支援自定義請求標頭和方法）。
2.  **解析 (Parse)**：使用 `jq` 語法解析 JSON 結構，將欄位映射到 RSS 條目。
3.  **元數據 (Metadata)**：定義訂閱源的標題和描述等詳情。
4.  **保存 (Save)**：直接將配置保存為自定義配方。

## 如何使用

在管理後台導航至 **工作台 > Curl 轉 RSS**。

### 第一步：請求配置 (Request Configuration)

你需要定義如何獲取 JSON 資料。

- **從 Curl 匯入 (Import from Curl)**：你可以貼上 `curl` 命令來自動填充 URL、方法、請求標頭和請求體。這在你從瀏覽器開發者工具複製請求時非常有用。
- **方法 (Method)**：選擇 `GET` 或 `POST`。
- **URL**：API 端點 URL。
- **Headers**：添加任何必要的請求標頭（例如 `Authorization`, `Content-Type`）。
- **請求體 (Request Body)**：對於 POST 請求，提供 JSON 請求體。

點擊 **抓取並下一步 (Fetch and Next)** 來獲取資料。

### 第二步：JQ 解析規則 (Parsing Rules)

獲取到 JSON 後，你將在左側面板看到以樹形視覺化的響應。現在你可以定義選取器來提取訂閱源條目。

該工具使用 **[jq](https://jqlang.github.io/jq/)** 語法來查詢 JSON。

- **列表選取器 (Items Iterator)**：條目陣列的路徑。
  - 提示：你可以點擊樹視圖中的節點來自動填充選取器。
- **標題選取器 (Title Selector)**：條目標題的路徑（相對於條目對象）。
- **連結選取器 (Link Selector)**：條目 URL 的路徑。
- **日期選取器 (Date Selector)**：（可選）發布日期的路徑。
- **內容選取器 (Content Selector)**：（可選）完整內容或摘要的路徑。

點擊 **執行預覽 (Run Preview)** 驗證你的選取器，然後點擊 **下一步 (Next Step)**。

### 第三步：訂閱源元數據 (Feed Metadata)

配置 RSS 訂閱源詳情：

- **訂閱源標題 (Feed Title)**：你的新訂閱源名稱。
- **描述 (Description)**：簡短描述。
- **網站連結 (Site Link)**：原始網站 URL。
- **作者 (Author)**：（可選）作者詳情。

### 第四步：保存配方 (Save Recipe)

審查你的配置並將其保存為永久配方。

- **配方唯一 ID (Recipe Unique ID)**：此訂閱源配置的唯一識別碼（例如 `my-custom-api-feed`）。如果留空，將自動根據訂閱源標題生成。
- **內部描述 (Internal Description)**：關於此配方的備註。

點擊 **確認並保存 (Confirm and Save)**。工具將自動建立一個包含你的配置的新自定義配方，你可以在 **自定義配方 (Custom Recipes)** 儀表板中管理它。
