# Inbox Source 设计方案

> 状态：待实施

## 1. 概述

Inbox 是一种新的 Source 类型。与现有 Source（RSS/HTML/JSON/Search）"系统主动拉取"不同，Inbox 是"第三方主动推送"——外部程序通过 HTTP POST 接口向 inbox 写入数据，系统将其存储并以标准 `CraftFeed` 形式提供给后续的 Craft / Recipe / Topic 流程。

类比：一个 email 收件箱。外部程序往里投递，FeedCraft 负责存储和消费。

## 2. 核心概念

| 概念           | 说明                                                                        |
| -------------- | --------------------------------------------------------------------------- |
| **Inbox**      | 一个具名的数据收件箱，有唯一 ID、标题、描述、容量上限。支持创建多个独立实例 |
| **InboxItem**  | 收件箱中的一条数据，对应一篇 `CraftArticle`                                 |
| **InboxToken** | 全局写入令牌。持有有效 token 的第三方程序可以向任意 inbox POST 数据         |

## 3. 存储方案

InboxItem 数据存放在现有的 SQLite 数据库（`feed-craft.db`）中，与 Recipe、Topic 等业务数据共用同一个库。

理由：

- Redis 在当前系统中是可选依赖（仅用于缓存），inbox 数据是不可重建的持久化业务数据，不适合放 Redis
- Inbox 数据量很小（每个 inbox 最多约 100 条），不存在性能瓶颈
- 共用数据库避免引入第二个 DB 实例、额外环境变量和 volume 挂载，保持部署简单

## 4. 数据模型

### 4.1 Inbox

```go
// internal/dao/inbox.go

type Inbox struct {
    BaseModelWithoutPK
    ID          string `gorm:"primaryKey" json:"id" binding:"required"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    MaxItems    int    `gorm:"default:100" json:"max_items"` // 滚动窗口上限
}
```

### 4.2 InboxItem

```go
// internal/dao/inbox.go

type InboxItem struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    InboxID     string    `gorm:"uniqueIndex:idx_inbox_item_id;not null" json:"inbox_id"`
    ItemID      string    `gorm:"uniqueIndex:idx_inbox_item_id;not null" json:"item_id"` // 业务唯一标识, 用于去重
    Title       string    `gorm:"not null" json:"title"`
    URL         string    `json:"url,omitempty"`
    Content     string    `gorm:"type:text" json:"content,omitempty"`
    Summary     string    `json:"summary,omitempty"`
    Author      string    `json:"author,omitempty"`
    PublishedAt time.Time `json:"published_at"`                        // 文章发布时间
    CreatedAt   time.Time `json:"created_at"`                          // DB 插入时间, GORM 自动填充
}
```

其中：

- `Content` 存储文章正文内容
- `URL` 优先保存调用方提供的原文链接；若 POST 写入时该字段为空，则服务端自动回填为 FeedCraft 内部内容访问地址：`<feedcraft-site-base-url>/:inbox_id/:article_id`
- 上述路由中的 `article_id` 对应 `InboxItem.ItemID`

`(InboxID, ItemID)` 建联合唯一索引，用于去重（见 5.6）。

### 4.3 InboxToken

```go
// internal/dao/inbox_token.go

type InboxToken struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    Label     string    `json:"label,omitempty"` // 用途说明，如 "github-webhook"
    CreatedAt time.Time `json:"created_at"`
}
```

## 5. POST 写入接口

### 5.1 端点

```
POST /api/inbox/:inbox_id/items
Header: Authorization: Bearer <inbox_token>
```

### 5.2 请求体

**统一使用 JSON 数组格式**，支持批量写入。单条数据即长度为 1 的数组。单次请求上限 100 条。

```jsonc
[
  {
    "title": "服务器 CPU 告警", // 必填
    "url": "https://example.com/alert", // 可选, 原文链接
    "content": "<p>详细内容</p>", // 可选, 正文 HTML
    "summary": "CPU 使用率超过 90%", // 可选, 摘要文本
    "id": "alert-20250411-001", // 可选, 唯一标识 (用于去重)
    "author": "monitoring-bot", // 可选, 作者名
    "timestamp": 1744365600 // 可选, Unix 秒级时间戳
  }
]
```

最简用法：

```json
[{ "title": "hello" }]
```

如果 `url` 为空或省略，服务端会在写入时自动回填：`<feedcraft-site-base-url>/:inbox_id/:article_id`。

### 5.3 字段说明

| 字段        | 类型   | 必填   | 说明                                                              |
| ----------- | ------ | ------ | ----------------------------------------------------------------- |
| `title`     | string | **是** | 标题                                                              |
| `url`       | string | 否     | 原文链接；若为空或省略，则服务端自动填充为 Inbox 文章内容访问地址 |
| `content`   | string | 否     | 正文内容（支持 HTML）                                             |
| `summary`   | string | 否     | 摘要文本，不填则自动截取 content 前 200 字符                      |
| `id`        | string | 否     | 唯一标识，用于去重。不填则服务端自动生成 UUID（不参与去重）       |
| `author`    | string | 否     | 作者名                                                            |
| `timestamp` | number | 否     | Unix 秒级时间戳，不填则使用服务端当前时间                         |

### 5.4 POST 请求体与 DB 模型的映射

| POST 字段    | DB 字段 (InboxItem) | 默认值逻辑                                                                           |
| ------------ | ------------------- | ------------------------------------------------------------------------------------ |
| _(URL 路径)_ | `InboxID`           | 从 URL 路径 `/api/inbox/:inbox_id/items` 提取                                        |
| `title`      | `Title`             | 直接存储                                                                             |
| `url`        | `URL`               | 优先使用请求值；若为空则自动填充为 `<feedcraft-site-base-url>/:inbox_id/:article_id` |
| `content`    | `Content`           | 直接存储                                                                             |
| `summary`    | `Summary`           | 不填则截取 Content 前 200 字符                                                       |
| `id`         | `ItemID`            | 不填则服务端生成 UUID                                                                |
| `author`     | `Author`            | 直接存储                                                                             |
| `timestamp`  | `PublishedAt`       | 不填则 = `CreatedAt`                                                                 |
| _(无)_       | `ID`                | DB 自增主键                                                                          |
| _(无)_       | `CreatedAt`         | GORM 自动填充为入库时间                                                              |

### 5.5 写入行为

1. 验证 token 有效性（查 `inbox_tokens` 表）
2. 验证 `inbox_id` 对应的 inbox 存在
3. 校验所有条目（每条必须有 `title`，总数不超过 100）。任一条校验失败则整批拒绝，返回 400
4. 填充默认值（`id` → UUID, `timestamp` → 当前时间, `summary` → 截取 content 前 200 字符）
5. 若 `url` 为空或省略，则服务端根据 `<feedcraft-site-base-url>/:inbox_id/:article_id` 自动生成内容访问地址，其中 `article_id = ItemID`
6. 批量 upsert `InboxItem`
7. 滚动清理：查询该 inbox 当前总条数，若超过 `max_items`，删除最旧的记录使总数 = `max_items`

### 5.6 去重机制

去重依据 `(inbox_id, item_id)` 联合唯一索引，行为类似 RSS 的 `<guid>` 去重：

| 场景                                     | 行为                           |
| ---------------------------------------- | ------------------------------ |
| 调用方传了 `id`，inbox 中已存在相同 `id` | **Upsert**：用新数据覆盖旧记录 |
| 调用方传了 `id`，inbox 中不存在          | 正常插入                       |
| 调用方未传 `id`（系统自动生成 UUID）     | 永远插入（UUID 不会重复）      |

实现方式：GORM `Clauses(clause.OnConflict{...})` 执行 upsert。

### 5.7 响应格式

```json
{
  "total": 3,
  "created": 2,
  "updated": 1
}
```

- `total`：本次请求的条目总数
- `created`：新插入的条数
- `updated`：因 `id` 重复而更新的条数

### 5.8 InboxItem 到 CraftArticle 的映射

| InboxItem 字段 | CraftArticle 字段 | 说明                         |
| -------------- | ----------------- | ---------------------------- |
| `Title`        | `Title`           |                              |
| `URL`          | `Link`            |                              |
| `Content`      | `Content`         |                              |
| `Summary`      | `Description`     |                              |
| `ItemID`       | `Id`              | RSS guid / Atom id           |
| `Author`       | `AuthorName`      |                              |
| `PublishedAt`  | `Created`         | RSS pubDate / Atom published |
| _(无)_         | `Updated`         | = `PublishedAt`              |

## 6. Token 鉴权

```
第三方程序                        FeedCraft
    |                                |
    |  POST /api/inbox/:id/items     |
    |  Authorization: Bearer <token> |
    |  Body: [{...}, {...}]          |
    |------------------------------->|
    |                                | 1. 查 DB: inbox_tokens 表是否存在该 token
    |                                | 2. 验证 inbox_id 是否存在
    |                                | 3. 校验所有条目 (title 必填, 总数 ≤ 100)
    |                                | 4. 填充默认值 + 批量 upsert
    |                                | 5. 滚动清理
    |          201 Created           |
    |  {total: 3, created: 2,        |
    |   updated: 1}                  |
    |<-------------------------------|
```

Token 管理通过管理后台操作：

- 创建 token 时自动生成 UUID 作为 token 值
- token 是全局的，对所有 inbox 通用
- 用户可以为 token 添加 label 标注用途

## 7. Source 层集成

### 7.1 新增 SourceType

```go
// internal/constant/source_type.go
const SourceInbox SourceType = "inbox"
```

### 7.2 新增 SourceConfig 字段

```go
// internal/config/source_config.go

type InboxSourceConfig struct {
    InboxID string `json:"inbox_id"`
}

type SourceConfig struct {
    // ...现有字段...
    InboxSource *InboxSourceConfig `json:"inbox_source,omitempty"`
}
```

### 7.3 InboxSource 实现

```go
// internal/source/inbox.go

func init() {
    Register(constant.SourceInbox, inboxSourceFactory)
}

type InboxSource struct {
    InboxID string
}

func (s *InboxSource) Fetch(ctx context.Context) (*model.CraftFeed, error) {
    // 1. 从 DB 读取 inbox 元信息
    // 2. 读取该 inbox 下的所有 InboxItem (按 created_at DESC)
    // 3. 转换为 CraftFeed + []CraftArticle (见 5.6 映射表)
    // 4. 返回
}

func (s *InboxSource) BaseURL() string { return "" }
```

InboxSource 不需要 Fetcher 和 Parser，直接从数据库读取并转换为 `CraftFeed`。

## 8. 与现有系统的集成

### 8.1 Recipe 集成

用户在管理后台创建 Recipe 时选择 `source_type: "inbox"`，配置 `inbox_source.inbox_id`：

```
Recipe "my-inbox-feed" = InboxSource(inbox_id="alerts") + Craft("summary")
```

访问 `/recipe/my-inbox-feed` 即可得到经过 AI 加工的 inbox RSS。

### 8.2 Topic 集成

Inbox 可以通过 `feedcraft://recipe/:id` 间接被 Topic 引用：先建一个 inbox recipe，再把该 recipe 加入 topic 的输入列表。

### 8.3 URI 协议扩展（未来考虑）

未来可扩展 `feedcraft://inbox/:id` 让 Topic 直接引用 inbox，MVP 阶段不需要。

## 9. 路由汇总

### 9.1 第三方写入（Inbox Token 鉴权）

| 方法     | 路径                         | 说明                                   |
| -------- | ---------------------------- | -------------------------------------- |
| **POST** | `/api/inbox/:inbox_id/items` | 批量写入数据（JSON 数组，上限 100 条） |

### 9.2 内容读取（公开访问）

| 方法    | 路径                                              | 说明                              |
| ------- | ------------------------------------------------- | --------------------------------- |
| **GET** | `<feedcraft-site-base-url>/:inbox_id/:article_id` | 返回对应文章的 `content` 字段内容 |

该路由中的 `article_id` 对应 `InboxItem.ItemID`。如果某条数据在 POST 写入时未提供 `url`，系统会自动将 `url` 回填为这个内容访问地址。

### 9.3 管理后台（JWT 鉴权）

| 方法   | 路径                          | 说明               |
| ------ | ----------------------------- | ------------------ |
| GET    | `/api/admin/inboxes`          | 列出所有 inbox     |
| POST   | `/api/admin/inboxes`          | 创建 inbox         |
| GET    | `/api/admin/inboxes/:id`      | 获取 inbox 详情    |
| PUT    | `/api/admin/inboxes/:id`      | 更新 inbox         |
| DELETE | `/api/admin/inboxes/:id`      | 删除 inbox         |
| GET    | `/api/admin/inbox-tokens`     | 列出所有写入 token |
| POST   | `/api/admin/inbox-tokens`     | 创建新 token       |
| DELETE | `/api/admin/inbox-tokens/:id` | 删除 token         |

## 10. 文件变动清单

### 10.1 改动现有文件

| 文件                               | 改动                                                              |
| ---------------------------------- | ----------------------------------------------------------------- |
| `internal/constant/source_type.go` | 新增 `SourceInbox` 常量                                           |
| `internal/config/source_config.go` | 新增 `InboxSourceConfig` 结构体和 `SourceConfig.InboxSource` 字段 |
| `internal/dao/migrate.go`          | `AutoMigrate` 中注册 `Inbox`、`InboxItem`、`InboxToken` 三个表    |
| `internal/router/registry.go`      | 注册 inbox 相关路由                                               |

### 10.2 新增文件

| 文件                                 | 职责                                                  |
| ------------------------------------ | ----------------------------------------------------- |
| `internal/dao/inbox.go`              | `Inbox` + `InboxItem` 模型定义及 CRUD                 |
| `internal/dao/inbox_token.go`        | `InboxToken` 模型定义及 CRUD                          |
| `internal/source/inbox.go`           | `InboxSource` 实现，`init()` 中注册到 source registry |
| `internal/controller/inbox.go`       | 管理后台 Inbox CRUD 接口                              |
| `internal/controller/inbox_push.go`  | 第三方 POST 写入接口                                  |
| `internal/middleware/inbox_token.go` | inbox token 鉴权中间件                                |
