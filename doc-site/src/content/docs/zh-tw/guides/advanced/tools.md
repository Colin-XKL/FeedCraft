---
title: 故障排除與工具
description: 了解如何使用 FeedCraft 內置的故障排除和實用工具。
sidebar:
  order: 50
---

FeedCraft 在管理後台內置了多種工具，幫助你調試、驗證和監控你的 RSS Feed。

## RSS 預覽 (Feed Viewer)

**RSS 預覽** 工具允許你以清晰易讀的格式預覽任何 RSS Feed。這對於驗證 Feed 是否有效以及檢查其內容非常有幫助，而無需使用外部閱讀器。

- **位置**: **工具 (Tools) > RSS 預覽**
- **使用方法**:
  1.  輸入你想要預覽的 RSS Feed URL。
  2.  點擊 **預覽**。
  3.  Feed 內容將顯示在下方，包括標題、描述和具體文章列表。

## Feed 對比 (Feed Compare)

**Feed 對比** 工具允許你並排比較兩個 RSS Feed。這在調試 FlowCraft 時特別有幫助，可以直觀地看到處理步驟如何改變了原始 Feed。

- **位置**: **工具 (Tools) > Feed 對比**
- **使用方法**:
  1.  輸入原始 Feed 的 **來源連結**。
  2.  從下拉選單中選擇一個 **FlowCraft** (例如 `translate-title`)。
  3.  點擊 **對比**。
  4.  工具將在左側顯示「原始 Feed」，在右側顯示「處理後的 Feed」，方便你發現差異。

## 系統健康 (System Health)

**系統健康** 儀表板可視化展示了 FeedCraft 實例的內部依賴關係圖。它顯示了 Recipe (配方)、FlowCraft (流程) 和 AtomCraft (原子) 之間的關係。

- **位置**: **工具 (Tools) > 系統健康** (或 Craft 依賴檢查)
- **用途**:
  - 驗證自訂 Recipe 的所有組件是否存在。
  - 檢測缺失的依賴項（例如，Recipe 引用了已刪除的 FlowCraft）。
  - 識別可能導致問題的循環依賴。
- **使用方法**:
  1.  點擊 **開始分析**。
  2.  查看樹狀圖。健康的組件會標記其類型（Recipe, FlowCraft, AtomCraft）。缺失的組件將標記為紅色。

## URL 生成器

**URL 生成器** 幫助你輕鬆創建 FeedCraft 訂閱連結。它還具有「解析模式」，可以反向解析現有的 URL。

- **位置**: **儀表板 > 快速開始**
- **功能**:
  - **生成**: 選擇 FlowCraft 並輸入來源 URL，即可獲取訂閱連結。
  - **解析**: 貼上 FeedCraft URL 以查看其組成部分（使用的 FlowCraft、原始來源等）。

## 依賴服務 (Dependency Services)

若要監控外部服務（如 Redis, Browserless, LLM）的狀態，請使用 **依賴服務狀態** 儀表板。

- **位置**: **設定 > 依賴服務狀態**
- **另請參閱**: [進階定製](/zh-tw/guides/advanced/customization/#依賴服務-dependency-services)
