# 错误处理与可观测性 (Error Handling & Observability)

## 1. 需求场景

在单源 Recipe 模式下，如果抓取失败，往往只影响一个订阅。但在 Topic 聚合模式下，一个 Topic 可能底层依赖 5-10 个不同渠道的 Recipe。随着时间推移，某些渠道可能会由于网站改版（HTML-to-RSS 选择器失效）、反爬策略或 API 速率限制而默默失效。如果只在后台记录日志（Log），管理员和用户往往很难及时发现“部分数据源已经不再产出内容”，导致信息逐渐产生盲区。

## 2. 当前已有的一些想法

未来需要为 FeedCraft 引入一套完善的可观测性和告警机制：

- **DB 持久化记录**: 底层 Recipe 在异步抓取（或缓存刷新）失败时，不再仅仅打印 Warn 日志，而是将错误详情（时间、Recipe ID、错误类型如超时、解析失败、无数据等）记录到一张专门的 `ExecutionLog` 数据库表中。
- **Admin Dashboard 面板**:
  - 在管理后台提供一个健康状态监控面板，直观展示每个 Topic 和 Recipe 的健康度（Health Status）。
  - 若某个 Recipe 连续失败 N 次，将其标记为“故障/下线 (Offline)”，并在 UI 上高亮警示，提示管理员介入修复。
- **通知与分发 (System Notification)**:
  - 系统自身生成一个专属的 `System Notification RSS Feed`。当核心信源失效时，往这个 Feed 写入一条警报文章，管理员只需订阅这个内部 Feed，就能在日常使用的 RSS 阅读器中随时收到系统故障警报。
  - 可选集成 Webhook（例如推送到 Slack、飞书或 Telegram）。
- **优雅降级 (Graceful Degradation)**: 在错误积累到一定程度时，自动暂停（Pause）对该失效信源的无效拉取，节约系统资源，并等待管理员手动修复后重启。

## 3. 与当前架构关联的地方

- **隔离聚合逻辑与监控逻辑**: 核心的聚合流水线依然保持轻量（遇到错误忽略并合并成功的部分），但在触发错误的节点，通过事件总线（Event Bus）或 Go Channels 发出错误事件，交由独立的可观测性服务（Observability Service）异步消费并落库，不阻塞用户的订阅响应。
- **调度器增强**: 预热调度器（Preheating Scheduler）需要感知各个 Recipe 的健康状态，避免对已永久失效的信源继续盲目重试。
- **前端改造**: Admin 接口和前端 Workbench 需要增加新的路由和页面，用于查询和展示错误日志与健康度报表。
