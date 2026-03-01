# LLM 与网页抓取并发控制及重试机制改造方案

## 1. 原始核心需求与场景痛点

本次改造旨在解决 FeedCraft 在启用并发处理文章后，由于“并发水桶效应”导致的两大资源瓶颈问题：

1. **支持连接池与重试，优先复用连接**：LLM API 调用需要复用 TCP 连接，并在遇到限流时能够自动重试且不阻塞正常排队。
2. **支持外层大并发，底层细粒度限流**：当开启 `craft` 层的全量并发后，如果遇到包含多篇文章的 Feed：
   - **痛点 A (LLM 场景)**：并发请求瞬间打满 LLM 后端，容易触发频率限制（Rate Limit）。
   - **痛点 B (网页抓取场景)**：`fulltext` 功能瞬间对单一第三方目标域名（如 `nytimes.com`）发起大量并发，容易造成对方高负载甚至直接封禁我们的 IP (产生 HTTP 429 或 403 错误)。

我们需要在保持代码简洁优雅的前提下，对这两种不同特性的“重资源”分别实施保护。

---

## 2. 方案的核心特色与场景对比 (Why & How)

针对上述痛点，本方案采取**“双漏斗”设计策略**。虽然两者都需要限流，但由于场景特性的巨大差异，必须采用不同的底层模型：

### 2.1 典型场景一：LLM API 调用 (全局固定限流 + 优先级队列)

- **场景特点**：
  - 目标单一：绝大多数情况下，所有请求都打向同一个或固定的几个 LLM 代理网关端点。
  - 失败成本高且需重试：由于网络抖动和后端限流，失败率相对较高，需要引入重试退避机制。
  - 积压问题：如果失败重试的请求排在所有新请求末尾，极易导致该上层 Web 请求直接超时。
- **解决方案**：引入泛型的 `PriorityDispatcher` 模块。
  - 内部预启动固定数量的 Worker 协程（如 5 个），这天然实现了**全局绝对最大并发限制**。
  - 采用普通队列 (`normalQueue`) 和高优队列 (`urgentQueue`) 双通道设计。当请求因失败而触发 `retry-go` 退避休眠后，将其唤醒投入高优队列，Worker 一旦空闲即可“插队”接手，确保旧请求能尽早闭环释放连接。

### 2.2 典型场景二：Fulltext 网页抓取 (按域名动态限流 + 信号量)

- **场景特点**：
  - 目标发散：目标域名可能有成千上万个 (`a.com`, `b.com`...)。
  - 互不干扰：对 `a.com` 的高并发限制，绝对不能阻塞对 `b.com` 的正常抓取。
  - 禁用 Worker 模型：如果我们为每个域名启动后台 Worker，会导致严重的 Goroutine 泄露。
- **解决方案**：引入 `KeyedLimiter` 模块 (基于官方 `golang.org/x/sync/semaphore`)。
  - 根据 `url.Host` 提取目标域名作为 Key。
  - 使用 `sync.Map` 懒加载维护按域名的独立信号量。
  - **绝妙之处**：信号量仅消耗微小内存且无后台驻留。同时，由于 `sem.Acquire` 支持传递 `Context`，一旦前端取消了对 FeedCraft 的访问，所有排队等待配额的抓取请求会瞬间 `Canceled` 退出，绝不白白浪费服务器流量。

### 2.3 优雅的第三方库集成

- 避免重复造轮子，引入三个流行库支撑上述架构：
  - **`github.com/sourcegraph/conc`**：在上层业务逻辑 (`craft`) 负责“尽情发散并发”，替换冗长的 `sync.WaitGroup`。
  - **`github.com/avast/retry-go/v4`**：负责优雅重试与指数退避，与 `PriorityDispatcher` 完美结合。
  - **`golang.org/x/sync/semaphore`**：构建按域名并发的底层支撑。

---

## 3. 详细实施计划 (Plan V6)

### 3.1 引入第三方依赖

执行以下命令安装依赖：

```bash
go get github.com/sourcegraph/conc
go get github.com/avast/retry-go/v4
go get golang.org/x/sync/semaphore
```

### 3.2 核心基建 1：实现 `PriorityDispatcher` (服务于 LLM)

- **路径**: `internal/util/priority_dispatcher.go`
- **实现细节**:
  - 定义泛型结构体 `PriorityDispatcher[R any]`。
  - 启动固定数量的 Worker，读取 `FC_LLM_MAX_CONCURRENCY`。
  - 提供带 `urgent bool` 参数的 `Execute` 方法供外部阻塞式调用。

### 3.3 核心基建 2：实现 `KeyedLimiter` (服务于网页抓取)

- **路径**: `internal/util/keyed_limiter.go`
- **实现细节**:
  - 封装基于 `sync.Map` 和 `semaphore.Weighted` 的动态限流器。
  - 提供 `Acquire(ctx context.Context, key string) (release func(), err error)` 方法。

### 3.4 LLM 客户端连接复用与重试组装

- **路径**: `internal/adapter/common_llm.go`
- **实现细节**:
  - 增加 `llmClientCache sync.Map` 缓存已初始化的客户端，复用 TCP 握手。
  - 在 `SimpleLLMCall` 中，利用 `retry.DoWithData` 包裹执行逻辑，根据尝试次数决定优先级，向 `PriorityDispatcher` 提交任务。

### 3.5 Fulltext 抓取并发限制

- **路径**: (假定爬取核心点在 `internal/craft/fulltext.go` 或独立的 extractor 中)
- **实现细节**:
  - 初始化全局 `domainLimiter`，读取配置如 `FC_DOMAIN_MAX_CONCURRENCY` (默认例如 3)。
  - 在发起 HTTP 真实抓取动作前，解析 `URL` 提取 `Host`，通过 `domainLimiter.Acquire` 拿锁，成功结束后释放。

### 3.6 Craft 处理全面并发化

- **路径**: `internal/craft/option.go`, `internal/craft/advertorial.go`, `internal/craft/llm_filter.go`
- **实现细节**:
  - 使用 `conc/iter.MapErr` 替换原有的串行 `for` 循环处理。
  - 对涉及到 `lo.Filter` 的同步过滤，拆分为“并发推断 + 同步过滤”两阶段执行，保证原文顺序和最高执行效率。
