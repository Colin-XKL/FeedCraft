---
title: 系統內建 AtomCraft
description: FeedCraft 系統內建處理原子工藝 (AtomCrafts) 的詳細參考指南。
---

FeedCraft 內建了一系列「原子工藝 (AtomCrafts)」，用於對訂閱源進行特定的處理。你可以將這些原子工藝組合成「組合工藝 (FlowCraft)」來構建強大的數據管道。

## 內容獲取與修復

這些原子主要用於獲取全文或修復常見的訂閱源問題。

### `fulltext` (全文提取)

從原始網頁提取文章的全文內容。

- **適用場景:** 當 RSS 訂閱源僅提供摘要或片段時。
- **機制:** 使用標準 HTTP 用戶端請求網頁，並使用演算法提取正文。速度快且輕量。

### `fulltext-plus` (瀏覽器全文提取)

使用無頭瀏覽器 (Puppeteer) 提取全文。

- **適用場景:** 針對透過 JavaScript 動態渲染內容或有較強反爬蟲措施的網站。
- **機制:** 連線到配置的 Browserless/Puppeteer 服務來渲染頁面。速度較慢但相容性更強。
- **參數:**
  - `mode` (預設: `networkidle2`): 頁面載入等待模式。
    - `load`: 等待 `load` 事件。
    - `domcontentloaded`: 等待 `DOMContentLoaded` 事件。
    - `networkidle0`: 等待直到 500ms 內沒有活躍的網絡連線。
    - `networkidle2`: 等待直到 500ms 內活躍的網絡連線數不超過 **2** 個。(推薦用於 SPA 單頁應用)。
  - `wait` (預設: `0`): 顯式等待時間（秒），例如 `5`。

### `proxy` (代理)

簡單的訂閱源代理。

- **適用場景:** 當你只想轉發原始 Feed 而不做修改，或者將 FeedCraft 作為中心閘道使用時。

### `guid-fix` (GUID 修復)

使用文章內容的 MD5 哈希值替換 RSS 條目的 GUID。

- **適用場景:** 某些訂閱源在內容未變更的情況下頻繁更改 GUID，導致閱讀器中出現重複的未讀條目。此原子可基於內容穩定 GUID。

### `relative-link-fix` (相對連結修復)

將內容中的相對連結（如 `<a href="/about">`）轉換為絕對連結（如 `<a href="https://example.com/about">`）。

- **適用場景:** 提取全文後必不可少，否則在 RSS 閱讀器中查看時連結會失效。

### `cleanup` (HTML 清理)

清理 HTML 內容以去除雜亂資訊。

- **適用場景:** 透過移除多餘的 class、style 和空標籤來提高可讀性。

---

## 過濾類 (Filtering)

控制哪些條目可以進入最終生成的 Feed。

### `limit` (數量限制)

限制 Feed 中的條目數量。

- **參數:**
  - `num` (預設: `10`): 保留的最大條目數。

### `time-limit` (時間限制)

過濾掉超過指定天數的條目。

- **參數:**
  - `days` (預設: `7`): 文章保留的最大天數。

### `keyword` (關鍵詞過濾)

根據標題或內容中的關鍵詞進行過濾。

- **參數:**
  - `keywords`: 逗號分隔的關鍵詞列表（子串匹配，區分大小寫）。例如：`ad,sell,SALE`。
  - `mode`: `include` (保留匹配項，預設) 或 `exclude` (移除匹配項)。
  - `scope`: `title` (標題), `content` (內容), 或 `all` (全部，預設)。

---

## AI 增強 (AI Enhancement)

使用大語言模型 (LLM) 來轉換和豐富你的內容。

:::note
使用此類原子需要在環境變數中配置 LLM (API Key, Base URL 等)。
:::

### `translate-title` (標題翻譯)

將文章標題翻譯為你的目標語言。

- **參數:**
  - `prompt`: 自定義提示詞。預設使用標準翻譯提示詞。支援 `{{.TargetLang}}` 佔位符。

### `translate-content` (內容翻譯)

翻譯整篇文章內容，替換原文。

- **參數:**
  - `prompt`: 自定義提示詞。支援 `{{.TargetLang}}`。

### `translate-content-immersive` (沉浸式翻譯)

雙語翻譯模式。在每一段原文後面追加翻譯後的內容。

- **參數:**
  - `prompt`: 自定義提示词。

### `summary` (AI 摘要)

生成文章摘要並將其添加到正文開頭。

- **參數:**
  - `prompt`: 用於生成摘要的自定義提示詞。

### `introduction` (AI 導讀)

為文章生成簡短的介紹或導語。

- **參數:**
  - `prompt`: 自定義提示詞。

### `beautify-content` (智能排版)

使用 LLM 重新格式化文章，修復排版錯誤，去除廣告，並標準化 Markdown 格式，最後轉換回乾淨的 HTML。

- **參數:**
  - `prompt`: 設定「編輯」角色的指令。

---

## AI 過濾 (AI Filtering)

利用語義理解進行進階過濾。

### `ignore-advertorial` (軟文過濾)

使用 LLM 檢測文章是否為軟文或廣告，並將其移除。

- **參數:**
  - `prompt-for-exclude`: 如果文章是廣告，應返回 `true` 的提示詞。

### `llm-filter` (通用 LLM 過濾)

通用的 LLM 過濾器。你可以定義**排除**條件。

- **参数:**
  - `filter_condition`: 自然語言描述的條件。如果 LLM 回答 "yes" (true)，則該條目會被**移除**。
  - _示例:_ "這篇文章是關於體育的嗎？" (移除體育類文章)。
