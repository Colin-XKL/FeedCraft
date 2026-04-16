---
title: 從搜尋結果生成 RSS
description: 使用 AI 供應商透過搜尋查詢生成 RSS 訂閱源。
sidebar:
  order: 4
  badge:
    text: beta
    variant: note
---

## 前提條件

在使用搜尋轉 RSS 功能之前，您需要在管理員設定中配置搜尋供應商。請參閱 [搜尋供應商配置指南](/zh-tw/guides/advanced/customization) 獲取設定說明。

FeedCraft 包含一個 **搜尋轉 RSS (Search to RSS)** 工具，允許你將搜尋查詢轉換為 RSS 訂閱源。這對於使用配置的搜尋供應商（例如 SearXNG, Bing, Google）追蹤新聞、話題或品牌提及非常有。

## 如何使用

1.  在管理後台導航至 **工作台 > 搜尋 轉 RSS**。

### 第一步：搜尋查詢 (Search Query)

1.  輸入你的 **搜尋查詢**（例如 `latest AI news` 或 `SpaceX launches`）。
2.  **Enhanced Mode**: (可選) 開啟此選項以使用 AI (LLM) 生成多個優化的搜尋查詢。這透過擴充你的原始查詢來幫助發現更多相關內容。
3.  點擊 **預覽結果 (Preview Results)** 獲取結果。

### 第二步：預覽結果 (Preview Results)

系統將使用配置的搜尋供應商獲取結果。

- 查看找到的項目列表（標題、日期、連結、描述）。
- 如果結果正確，點擊 **下一步 (Next Step)**。

### 第三步：Feed 元數據 (Feed Metadata)

自定義此 Feed 在 RSS 閱讀器中的顯示方式：

- **Feed 標題**：預設為 "搜尋: [查詢詞]"。
- **Feed 描述**：簡短描述。
- **站點連結**：搜尋結果頁面的連結。

### 第四步：保存配方 (Save Recipe)

1.  **配方 ID (Recipe ID)**：為此配方選擇一個唯一的識別碼（例如 `search-ai-news`）。這將成為你訂閱源 URL 的一部分。
2.  **內部描述 (Internal Description)**：關於此配方的備註。
3.  點擊 **確認並保存 (Confirm & Save)**。

## 訪問你的訂閱源

保存後，該配方將作為 **自定義配方 (Custom Recipe)** 存儲。你可以在 **Custom Recipes** 儀表板中管理它。

你的新訂閱源將透過以下地址訪問：
`http://your-feedcraft-instance/rss/custom/{recipe-unique-id}`
