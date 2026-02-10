---
title: 從 HTML 網頁生成 RSS
description: 使用視覺化選取器將任意網頁轉換為 RSS 訂閱源。
sidebar:
  order: 2
  badge:
    text: new
    variant: success

---

FeedCraft 內建了視覺化的 **從 HTML 網頁生成 RSS (HTML to RSS)** 工具，允許你生成選取器，以便為那些沒有提供 RSS 的網站建立訂閱源。

> **注意：** 此工具專為 HTML 頁面設計。如果你需要處理 JSON API，請使用 [從 CURL 語句生成 RSS](/zh-tw/guides/advanced/curl-to-rss/)。

## 概覽

HTML to RSS 工具允許你：

1.  **抓取 (Fetch)** 網頁內容。
2.  **視覺化選取 (Select)** 元素來定義什麼是訂閱源條目（標題、連結、日期、內容）。
3.  **元數據 (Metadata)**：定義訂閱源的標題和描述等詳情。
4.  **保存 (Save)**：直接將配置保存為自定義配方。

## 如何使用

1.  在管理後台導航至 **工作台 > HTML 轉 RSS**。

### 第一步：目標 URL (Target URL)

1.  輸入你想要抓取的網頁完整 URL（例如部落格列表或新聞網站）。
2.  **增強模式 (Enhanced Mode)**：如果網站需要 JavaScript 來載入內容（使用無頭瀏覽器），請啟用此選項。
3.  點擊 **抓取並下一步 (Fetch and Next)**。

### 第二步：提取規則 (Extract Rules)

此步驟允許你將 HTML 元素映射到 RSS 訂閱源欄位。

1.  **頁面預覽 (Page Preview)**：左側窗格顯示渲染後的網頁。
    - **選取模式 (Selection Mode)**：當處於激活狀態時，點擊預覽中的元素將生成選取器。
2.  **CSS 選取器 (列表項)**：
    - 點擊 **選取 (Pick)**。
    - 在預覽窗格中，點擊代表列表中*單篇文章*或條目的容器元素。
3.  **欄位選取器 (Field Selectors)**：
    - 設置好列表項後，你可以映射相對欄位。
    - **標題 (Title)**：點擊選取並選中包含標題文字的元素。
    - **連結 (Link)**：點擊選取並選中包含文章 URL 的元素。
    - **日期 (Date)**：（可選）選取日期元素。
    - **描述 (Description)**：（可選）選取包含摘要的元素。

點擊 **執行預覽 (Run Preview)** 測試你的選取器。如果滿意，點擊 **下一步 (Next Step)**。

### 第三步：訂閱源元數據 (Feed Metadata)

工具會嘗試從頁面自動提取元數據。

- **訂閱源標題 (Feed Title)**：如有必要請調整。
- **描述 (Description)**：訂閱源描述。
- **網站連結 (Site Link)**：源 URL。

### 第四步：保存配方 (Save Recipe)

審查你的配置並將其保存為永久配方。

- **配方唯一 ID (Recipe Unique ID)**：此訂閱源配置的唯一識別碼（例如 `tech-blog-feed`）。如果留空，將自動根據訂閱源標題生成。
- **內部描述 (Internal Description)**：關於此配方的備註。

點擊 **確認並保存 (Confirm and Save)**。工具將自動建立一個包含你的配置的新自定義配方，你可以在 **自定義配方 (Custom Recipes)** 儀表板中管理它。
