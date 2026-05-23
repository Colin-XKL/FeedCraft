---
title: 收件箱 (Inbox)
description: 透過 HTTP 推送從第三方來源接收文章，並將其轉為 RSS 訂閱源。
sidebar:
  order: 6
  badge:
    text: new
    variant: success
---

import { Steps } from '@astrojs/starlight/components';

**收件箱**功能允許外部服務、腳本或自動化工具透過 HTTP 主動將文章推送到 FeedCraft。每個收件箱都可以透過自訂配方產生 RSS 訂閱位址，方便在任何 RSS 閱讀器中訂閱。

## 整體流程

典型的使用流程如下：

1. 在管理面板中**建立收件箱**（取得唯一的 inbox ID）。
2. **建立系統授權令牌**（用於鑑權推送請求的密鑰）。
3. 透過腳本、自動化工具或第三方平台，使用標準 JSON HTTP POST **推送文章**。
4. **建立自訂配方**，將該收件箱作為資料來源。
5. 在閱讀器中**訂閱**產生的 RSS 位址。

## 管理收件箱

在管理面板中，前往 **工作台 > 收件箱管理**。

### 建立收件箱

<Steps>
1. 點擊**新建收件箱**。
2. 填寫必要欄位：
   - **收件箱 ID**：唯一的 URL 安全識別符（只允許小寫字母、數字、連字號、底線）。建立後無法修改。
   - **標題**：收件箱的可讀名稱。
   - **最大保存數**：最多保留的文章數量，超出後按**建立時間**從舊到新自動刪除（預設 100）。若設為 `0` 會立即刪除該收件箱的所有條目，請用大數值代替「無限制」。
   - **公開可見性**：開啟後，任何人可直接拉取文章內容；關閉後需提供系統授權令牌。
3. 點擊**確定**儲存。
</Steps>

### 編輯收件箱

點擊操作欄中的**編輯收件箱**，可修改標題、描述、最大保存數和公開可見性。收件箱 ID 不可修改。

:::caution
降低已有收件箱的**最大保存數**將立即刪除超出限制的文章（按建立時間從舊到新刪除），且無法復原。
:::

### 刪除收件箱

點擊操作欄中的**刪除**。此操作將永久刪除該收件箱及其內**所有文章**。

## 管理系統授權令牌

前往 **設定 > 系統授權令牌**，建立用於鑑權推送請求的 API 令牌。

<Steps>
1. 點擊**產生新令牌**。
2. 輸入描述性標籤（例如「iPhone 捷徑」、「Home Assistant」）。
3. 立即複製產生的令牌——它只會顯示一次，之後無法再次查看。
</Steps>

隨時可點擊**刪除**撤銷令牌。使用該令牌的所有整合將立即失效。

## 推送文章

使用推送介面，從任何 HTTP 用戶端、腳本或自動化平台推送文章。

### 介面位址

```
POST /api/inbox/{inbox_id}/items
```

### 鑑權

在 `Authorization` 請求標頭中攜帶系統授權令牌：

```
Authorization: Bearer YOUR_SYSTEM_AUTH_TOKEN
Content-Type: application/json
```

### 請求體

傳送一個 JSON 陣列，每個元素代表一篇文章。只有 `title` 是必填欄位。

```json
[
  {
    "id": "可選的自訂唯一ID",
    "title": "文章標題",
    "url": "https://example.com/article",
    "content": "<p>文章正文 HTML 內容。</p>",
    "summary": "簡短描述，在訂閱源預覽中顯示。",
    "author": "作者名",
    "timestamp": 1716470400
  }
]
```

| 欄位        | 必填 | 說明                                                                                   |
| ----------- | ---- | -------------------------------------------------------------------------------------- |
| `title`     | ✅   | 文章標題。                                                                             |
| `id`        | 可選 | 自訂穩定 ID。省略時自動產生 UUID。若相同 `id` 再次推送，則**更新**（upsert）已有文章。 |
| `url`       | 可選 | 文章原始連結。省略時 FeedCraft 自動產生指向儲存內容的連結。                            |
| `content`   | 可選 | 文章完整 HTML 正文。                                                                   |
| `summary`   | 可選 | 簡短描述，預設取 `content` 前 200 個 **Unicode 字元**（rune）。                        |
| `author`    | 可選 | 作者名。                                                                               |
| `timestamp` | 可選 | 發布時間的 Unix 時間戳記（秒）。預設為當前時間。                                       |

**批次限制**：每次請求最多 100 筆。

### cURL 範例

```bash
curl -X POST "https://YOUR_SERVER/api/inbox/my-inbox/items" \
  -H "Authorization: Bearer YOUR_SYSTEM_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{"title": "Hello World", "content": "<p>第一篇推送文章！</p>"}]'
```

### 回應格式

```json
{
  "total": 1,
  "created": 1,
  "updated": 0
}
```

## 透過 RSS 訂閱

### 直接訂閱（最簡方式）

每個收件箱都內建了一個可以直接訂閱的 RSS 位址，無需建立自訂配方：

```
GET /inbox/{inbox_id}/rss
```

將此位址直接貼到 RSS 閱讀器中即可訂閱。對於**私有收件箱**，需在 URL 後附加 Token：

```
/inbox/{inbox_id}/rss?token=YOUR_SYSTEM_AUTH_TOKEN
```

### 透過自訂配方訂閱（進階）

若需要對收件箱內容進行 Craft 處理（如 AI 翻譯、摘要產生、內容過濾），則建立自訂配方。

<Steps>
1. 前往**工作台 > 自訂配方**，點擊**新建配方**。
2. 將**資料來源類型 (Source Type)** 設定為 `inbox`。
3. 在 **Source Config JSON** 欄位中輸入：
   ```json
   { "inbox_source": { "inbox_id": "YOUR_INBOX_ID" } }
   ```
4. 將 **Craft** 設定為所需的處理鏈（例如 `translate-content`、`summary`）。
5. 儲存配方後，在配方列表中點擊**複製連結**即可取得 RSS 訂閱位址。
</Steps>

:::tip
簡單訂閱場景直接使用**直接訂閱位址**即可。只有在需要對推送內容進行 AI 加工（翻譯、摘要、過濾）時，才需要建立自訂配方。
:::

## 私有收件箱的存取控制

當收件箱關閉**公開可見性**後，文章內容介面需要進行身份驗證。

在文章 URL 後加上 `?token=YOUR_SYSTEM_AUTH_TOKEN`：

```
GET /inbox/{inbox_id}/items/{article_id}/content?token=YOUR_TOKEN
```

或使用 `Authorization: Bearer YOUR_TOKEN` 請求標頭。

:::note
RSS 訂閱源本身（透過自訂配方提供）的存取權限由配方自身的設定控制，與收件箱的公開標誌無關。私有標誌僅影響 `/inbox/{inbox_id}/items/{article_id}/content` 這一原始文章內容介面。
:::

## 垃圾回收 (GC)

FeedCraft 提供垃圾回收工具，可透過管理 API 使用：

- **GET `/api/admin/inboxes/gc/stats`** — 回傳總條目數、孤兒條目數（屬於已刪除收件箱的條目）和溢出條目數。
- **POST `/api/admin/inboxes/gc/cleanup`** — 透過單一原子交易刪除所有孤兒和溢出條目。
