---
title: 系統工具
description: 內建的除錯、比對和系統健康檢查工具。
sidebar:
  order: 5
---

FeedCraft 提供了一些內建工具來幫助您除錯 RSS 來源並監控系統健康狀況。您可以在管理後台的 **工具 (Tools)** 選單下訪問這些工具。

## RSS 預覽 (RSS Viewer)

**RSS 預覽** (Feed Viewer) 允許您按照 FeedCraft 的解析方式預覽任何 RSS 來源。

- **使用方法**:
  1. 導航至 **工具 > RSS 預覽**。
  2. 輸入一個 RSS/Atom 地址。
  3. 點擊 **預覽 (Preview)**。
- **目的**: 在設定配方 (Recipe) 之前，驗證 FeedCraft 是否能夠成功抓取和解析某個 Feed。
- **注意**: 預覽器默認使用 `proxy` 工藝，它只是簡單地抓取 Feed 而不進行修改。

## Feed 比對 (Feed Compare)

**Feed 比對** 工具讓您可以直觀地看到某個 Craft（Atom 或 Flow）對 Feed 的處理效果。

- **使用方法**:
  1. 導航至 **工具 > Feed 比對**。
  2. 輸入原始 RSS Feed 地址。
  3. 選擇一個 **FlowCraft** 或 **AtomCraft**。
  4. 點擊 **比對 (Compare)**。
- **輸出**: 工具顯示兩列：
  - **左側**: 原始 Feed 內容。
  - **右側**: 經過選定 Craft 處理後的 Feed 內容。
- **應用場景**: 非常適合測試新的翻譯流程或摘要提示詞，而無需創建永久配方。

## Craft 依賴檢查 (Craft Dependencies)

**Craft 依賴檢查** (System Health) 工具可視化您的 Recipes、FlowCrafts 和 AtomCrafts 之間的內部關係。

- **使用方法**:
  1. 導航至 **工具 > Craft 依賴檢查**。
  2. 點擊 **分析 Craft 依賴 (Analyze Craft Dependencies)**。
- **功能**:
  - 生成所有依賴關係的樹狀視圖。
  - **健康檢查**: 自動檢測丟失的依賴項（例如，指向已刪除 FlowCraft 的 Recipe），並在儀表板頂部醒目地突出顯示缺失的 Craft。
  - **視覺指示**: Recipes、Flows、Atoms 和丟失組件使用不同顏色標識。

:::tip
如果遇到 "Craft not found" 等錯誤，可以使用此工具追蹤配置中的斷鏈。
:::

## 系統運行狀態 (System Runtime)

**系統運行狀態** (Observability) 工具提供了一個全面的儀表板，用於監控資源的健康狀況和執行狀態。

- **使用方法**:
  1. 導航至 **工具 > 系統運行狀態**。
- **功能**:
  - **資源健康 (Resource Health)**: 查看配方及其他組件的當前狀態（健康、降級、暫停），包括連續失敗次數。
  - **執行日誌 (Execution Logs)**: 追蹤詳細的執行歷史、成功率以及每次運行的具體錯誤類型（例如：超時、網路錯誤、解析錯誤）。
  - **系統通知 (System Notifications)**: 查看關於資源狀態轉換的自動警報（例如當配方降級時）。您還可以透過內建的 RSS 來源 `/system/notifications/rss` 訂閱這些警報。

:::tip
如果配方反覆失敗並變為「暫停 (Paused)」狀態，您可以在解決根本問題後，透過系統運行狀態儀表板手動將其「恢復 (Resume)」。
:::

## 除錯工具 (Debug Tools)

### LLM 除錯 (LLM Debug)

用於測試 LLM 配置的沙盒。您可以向配置的 LLM 提供商發送測試提示，以驗證連接和模型響應。

### 廣告檢測除錯 (Ad Check Debug)

用於針對特定內容測試 "忽略軟文 (Ignore Advertorial)" 過濾邏輯的專用工具，以了解文章被過濾的原因。
