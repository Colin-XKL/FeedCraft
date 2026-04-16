# 领域一：Topic 路由与 RSS API 暴露

> 状态：已完成

## 1. 背景

`TopicFeed` 的目标不是停留在后台配置页面，而是要像 `RecipeFeed` 一样，成为一个真正可订阅、可对外暴露的 RSS 节点。

在当前设计里，用户会在后台配置一个 Topic：

- 填写 `input_uris`，引用多个上游数据源
- 配置聚合规则，如去重、排序、截断
- 最终期望通过统一的公开地址订阅，例如 `/topic/self-hosted-apps`

因此，Topic 的公开访问链路需要完成两件事：

- 从数据库配置动态构建可执行的运行时 Topic
- 将执行结果转换为标准 RSS XML 对外返回

在新的统一方向里，这条公开链路不应只服务 Topic，而应建立在统一输入模型之上：

- 顶层输入统一使用 `InputSpec`
- Topic 当前存储的 `input_uris` 只是在运行时映射成 `[]InputSpec(kind=uri)`
- 后续 Recipe 也可以复用同一套输入解析语义

## 2. 需求场景

典型场景是“围绕某一主题聚合多个来源并统一订阅”。

例如用户创建一个 `popular-selfhost-app` Topic：

- 输入 1：`feedcraft://recipe/official-blog`
- 输入 2：`feedcraft://recipe/search-self-hosted`
- 输入 3：`feedcraft://topic/self-hosted-tools`

用户预期的行为是：

- 访问 `/topic/popular-selfhost-app`
- 系统自动解析 Topic 配置并执行聚合
- 返回一个可被 RSS Reader 正常消费的标准 RSS 订阅结果

这里的“自动解析”不仅是解析 URI 字符串，还包括将输入统一解释为 `InputSpec`，并通过统一 builder 构建 provider graph。

## 3. 当前进展

- Topic 的数据库模型和自动迁移已经完成，见 `internal/dao/topic.go`
- Admin API 的 Topic CRUD 已经完成，见 `internal/controller/topic_feed.go`
- Admin 前端管理页已经完成，见 `web/admin/src/views/dashboard/topic_feed/topic_feed.vue`
- Topic 聚合核心引擎已经完成，见 `internal/engine/topic.go`
- `CraftFeed -> RSS` 的输出能力在系统中已经存在，`RecipeFeed` 和系统通知都已经走通类似路径

当前缺口：

- 公开路由 `GET /topic/:id` 尚未实现
- Topic 从数据库配置到运行时实例的动态装配尚未完成
- Topic 的 RSS 输出链路尚未接通
- Topic 公开路由尚未建立在统一输入模型之上

## 4. 未来待办

### 4.1 新增公开订阅路由

- 在 `internal/router/registry.go` 中注册 `GET /topic/:id`
- 与 `/recipe/:id` 保持一致，作为公开订阅入口，不放在 admin 路由下

### 4.2 实现公开控制器

- 在 `internal/controller/topic_feed.go` 中新增公开访问 handler
- 逻辑包括：
  - 根据 `:id` 加载 Topic 配置
  - 将 Topic 配置映射到统一输入模型
  - 调用 Topic runtime builder 构造可执行 `engine.TopicFeed`
  - 执行 `Fetch(ctx)` 获取 `CraftFeed`
  - 转换为 RSS XML 并返回 `application/rss+xml`

### 4.3 统一错误与返回行为

- Topic 不存在时返回 `404`
- Topic 构建失败或执行失败时返回明确的 `500` 错误
- 对公开路由的错误文案和日志格式做统一收敛，避免后续调试困难

### 4.4 与统一输入模型对齐

- Topic 公开路由不直接绑定 `input_uris -> provider` 的一次性逻辑
- 应通过统一的 `InputSpec -> FeedProvider` builder 构建运行时对象
- `http(s)` 输入的默认语义应与未来 Recipe 保持一致：
  - 视为第三方网站 RawFeed
  - 参考最小 Recipe + `proxy` craft 的处理方式
  - 不特殊收缩为“仅支持外部 RSS”

### 4.5 补齐端到端测试

- 新增 `/topic/:id` 的接口测试
- 覆盖正常返回、Topic 不存在、上游部分失败、全部失败等情况

## 5. 追踪建议

这个领域完成的标准是：

- 用户可以通过 `/topic/:id` 直接订阅 Topic
- 后台创建的 Topic 配置可以真正驱动公开输出
- Topic 的公开访问链路建立在统一输入模型之上，而不是孤立实现
- Topic 对外行为在接口层与 `RecipeFeed` 保持一致的可用性和可预期性
