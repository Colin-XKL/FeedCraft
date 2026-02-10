---
title: 核心概念
sidebar:
  order: 2
---

在深入使用 FeedCraft 之前，了解以下三個核心概念將非常有幫助。

## 原子工藝 (AtomCraft)

**AtomCraft** 是最小的處理單元。除了內建的原子工藝（如 `translate-title`, `fulltext`），你可以基於範本建立自定義的原子工藝。

**示例：自定義翻譯 Prompt**
你可以基於 `translate-content` 範本建立一個名為 `translate-to-french` 的新原子工藝，並在參數中填入自定義的 Prompt，指示 AI 將內容翻譯成法語。

## 組合工藝 (FlowCraft)

**FlowCraft** 是多個 AtomCraft 的組合序列。這允許你將多個操作串聯起來。

**示例：全文 + 摘要 + 翻譯**
你可以定義一個名為 `digest-and-translate` 的組合工藝，包含以下步驟：

1.  `fulltext` (提取正文)
2.  `summary` (生成摘要)
3.  `translate-content` (翻譯內容)

### 管理組合工藝

在管理後台導航至 **工作台 (Worktable) > 組合工藝 (FlowCraft)** 頁面建立和管理組合工藝。
編輯器允許你添加原子工藝並安排它們的執行順序。使用箭頭按鈕 (⬆️/⬇️) 調整順序，或使用垃圾桶圖示將其從流程中移除。

## 配方 (Recipe)

**Recipe** 將特定的 RSS 訂閱源 URL 與某個 原子工藝 (AtomCraft) 或組合工藝 (FlowCraft) 綁定。這允許你建立一個持久化的、經過定制的訂閱源 URL。

**管理配方：**
在管理後台導航至 **工作台 (Worktable) > 自定義配方 (Custom Recipe)** 頁面，你可以管理所有已建立的配方。

- **建立 (Create)**：綁定新的 URL 和工藝。
- **預覽 (Preview)**：點擊預覽按鈕，直接在內建的 Feed Viewer 中查看生成的效果。
- **複製連結 (Copy Link)**：點擊複製圖示獲取完整的訂閱 URL。

**示例：**

- **輸入 URL：** `https://news.ycombinator.com/rss`
- **處理器：** `digest-and-translate` (上面建立的工作流)
- **結果：** 你會得到一個新的 FeedCraft URL，訂閱它即可獲得帶全文、摘要和翻譯的 Hacker News。
