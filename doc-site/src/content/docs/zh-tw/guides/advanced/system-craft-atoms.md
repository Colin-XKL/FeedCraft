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
- **機制:** 連線到配置的瀏覽器提供方（`browserless-restful` 或 `cdp`）來渲染頁面。速度較慢但相容性更強。
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

使用文章內容的 FNV-1a 哈希值替換 RSS 條目的 GUID。

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

### `ai-content-process` (AI 內容處理)

根據自訂規則使用大型語言模型處理文章內容，並將生成結果插入到指定位置。

- **參數:**
  - `rule` (**必填**): 針對每篇文章內容的處理指令。例如："總結文章的關鍵觀點並列出行動建議"。
  - `extra-payload` (預設: `article_content`): 逗號分隔的附加資訊列表，可發送給 LLM。支援：`article_summary` (AI 生成的摘要), `article_content` (文章內容), `article_date` (文章日期), `raw_rss_item` (原始 RSS 節點資料)。
  - `placement` (預設: `prepend`): 生成內容的寫入位置。支援：`prepend` (在原文前追加), `replace` (替換原文), `append` (在原文後追加)。

### `beautify-content` (智能排版)

使用 LLM 重新格式化文章，修復排版錯誤，去除廣告，並標準化 Markdown 格式，最後轉換回乾淨的 HTML。

- **參數:**
  - `prompt`: 設定「編輯」角色的指令。

---

## AI 過濾 (AI Filtering)

利用語義理解進行進階過濾。

### `ignore-advertorial` (軟文過濾)

使用 LLM 檢測文章是否為軟文或廣告（綜合評估標題和內容），並將其移除。

- **參數:**
  - `prompt-for-exclude`: 如果文章是廣告，應返回 `true` 的提示詞。

### `llm-filter` (通用 LLM 過濾)

通用的 LLM 過濾器。你可以定義**排除**條件。LLM 會結合文章標題和內容進行評估。

- **參數:**
  - `filter_condition`: 自然語言描述的條件。如果 LLM 回答 "yes" (true)，則該條目會被**移除**。
  - _範例:_ "這篇文章是關於體育的嗎？" (移除體育類文章)。

### `embedding-filter` (Embedding 語意過濾)

基於 Embedding 模型的語意主題過濾器。它不會讓聊天模型逐條判斷文章，而是把「主題錨點」和文章內容都轉換成向量，再用餘弦相似度判斷是否匹配。

:::tip
當你需要快速、穩定地做主題過濾時，優先使用 `embedding-filter`，例如「只保留 AI 基礎設施新聞」或「移除體育文章」。如果規則需要複雜推理、政策判斷或結構化決策，再使用 `llm-filter`。
:::

#### 環境變數

建議單獨設定 Embedding 服務：

```bash
FC_EMBEDDING_API_TYPE=openai
FC_EMBEDDING_API_BASE=https://api.openai.com/v1
FC_EMBEDDING_API_KEY=sk-your-api-key
FC_EMBEDDING_API_MODEL=text-embedding-3-small # 必填
FC_EMBEDDING_BATCH_SIZE=5
FC_EMBEDDING_MAX_INPUT_CHARS=8000
```

`FC_EMBEDDING_API_TYPE` 支援：

- `openai`: OpenAI 或 OpenAI 相容的 Embedding 介面。
- `gemini`: 透過 Gemini 的 OpenAI 相容 Embedding 介面呼叫。請明確設定 `FC_EMBEDDING_API_BASE` 和 `FC_EMBEDDING_API_MODEL`。
- `ollama`: 本機 Ollama Embedding 模型。請設定 `FC_EMBEDDING_API_BASE`，例如 `http://localhost:11434`，並使用 `nomic-embed-text` 或 `bge-m3` 這類 Embedding 模型。

如果 `FC_EMBEDDING_API_TYPE`、`FC_EMBEDDING_API_BASE`、`FC_EMBEDDING_API_KEY` 都沒有設定，FeedCraft 會回退到對應的 `FC_LLM_API_TYPE`、`FC_LLM_API_BASE`、`FC_LLM_API_KEY`。Embedding 模型名是獨立的：**你必須顯式設定 `FC_EMBEDDING_API_MODEL` 為有效的 Embedding 模型。** FeedCraft 不再提供預設 Embedding 模型，也不會複用 `FC_LLM_API_MODEL`，因為它通常是聊天模型。

`FC_EMBEDDING_MAX_INPUT_CHARS` 是傳送給 Embedding 服務前的最終安全上限，包含 `instruction` 前綴。它是字元預算，不是精確的 tokenizer token 數。建議按模型 token 視窗設定保守值，例如 8k token 的 Embedding 模型可從 `8000` 開始。

#### 參數

- `anchors` (**必填**): 每行一個主題錨點。錨點應描述你想匹配的主題，例如：

  ```text
  人工智慧基礎設施
  機器學習研究
  大語言模型部署
  ```

- `threshold` (預設: `0.6`): 餘弦相似度閾值，範圍 `0` 到 `1`。值越高越嚴格。建議從 `0.6` 開始；漏掉相關文章時調低，混入無關文章時調高。
- `mode` (預設: `include`): `include` 保留匹配項；`exclude` 移除匹配項。
- `max_content_length` (預設: `2000`): 目前 AtomCraft 使用的文章正文最大字元數；最終傳送前還會受 `FC_EMBEDDING_MAX_INPUT_CHARS` 保護。
- `instruction` (可選): 會作為文字前綴拼接到每條 Embedding 輸入前。除非你的模型確實需要固定任務前綴，否則建議留空。

#### 管理後台使用流程

1. 打開 **工作台 → AtomCraft**。
2. 新增一個 AtomCraft，例如 `ai-news-only`。
3. 模板選擇 `embedding-filter`。
4. 在 `anchors` 中每行填寫一個主題。
5. 使用 `mode=include` 保留匹配文章，或切換到 `exclude` 移除匹配文章。
6. 儲存 AtomCraft。
7. 在 FlowCraft、Recipe、Feed Compare 中使用它，或直接存取：

```text
/craft/ai-news-only?input_url=https%3A%2F%2Fexample.com%2Ffeed.xml
```

#### 常見問題

- **"anchors parameter is required"**: `anchors` 至少需要一行非空內容。
- **"FC_EMBEDDING_API_MODEL must be set"**: 請設定單一 Embedding 模型。聊天模型不適合這個功能。
- **所有文章都被移除**: 降低 `threshold`，增加更寬泛的錨點，或增大 `max_content_length`。
- **無關文章仍然保留**: 提高 `threshold`，讓錨點更具體，或用代表性短語替代寬泛分類名。
- **服務商提示輸入過長**: 降低 `FC_EMBEDDING_MAX_INPUT_CHARS`，降低 `max_content_length`，或縮短 `instruction`。
