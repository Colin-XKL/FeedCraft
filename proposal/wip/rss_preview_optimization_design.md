# RSS 预览功能优化：支持 Recipe 运行期预览

核心目标: 构建便捷、易用的rss preview 页面工具

## 1. 需求场景 (User Scenarios)

* **配方调试与即时预览**  
  当用户（配方开发者/管理员）在平台上开发或修改一个 Custom Recipe 时，他们需要快速验证这个配方拉取源、解析以及 Craft 脚本处理后的最终输出效果。
* **避免非必要的公共接口暴露**  
  此前的预览功能通过跳转到类似 `/recipe/:id` 的公共 HTTP 渲染端点进行。这意味着：
  1. 预览未发布的或测试中的配方也必须通过公网/外部可访问的接口加载。
  2. 预览请求经历了不必要的 HTTP 环回开销（前端 -> 浏览器发起预览请求 -> 后端请求自己暴露的 HTTP 端点 -> 引擎再次渲染）。
* **统一的系统调试入口**  
  在“Feed 预览”调试工具中，用户既希望能够测试一个通用的外部 RSS/Atom 链接，也希望能直接点选已有的配方，查看其在真实运行期（feedruntime）环境下的执行结果。

---

## 2. 业务要求 (Requirements)

### 2.1 预览源模式切换
* **网络 URL 模式 (Raw URL Mode)**：支持用户手动输入标准的外部 `http://` 或 `https://` 链接。
* **本地配方模式 (Recipe Mode)**：用户可以在下拉菜单中直接选择一个平台已有的 Custom Recipe。

### 2.2 运行期统一解析
* 无论是输入外部 URL 还是选定内部配方，底层均应该使用 `feedruntime` 的统一输入协议（Input Spec URI）：
  * 外部 URL 自动解析并作为 `http(s)://...` 输入源运行。
  * 选定配方则映射为 `feedcraft://recipe/:id` 格式的内部 URI。
* 前后端接口应该保持高度兼容性。输入参数仍然可以传递为统一的 URI 字符串（或保持 `input_url` 参数命名以兼容现有客户端调用）。

### 2.3 自定义配方列表页联动
* 从“Custom Recipe（自定义配方）”管理列表中点击“预览”按钮时，不应再跳转到公共 RSS 订阅路径，而是携带 Recipe 标识直接导航到 `/tools/viewer` 页面，并在该页面自动加载、运行并展示该配方的运行期数据。

---

## 3. 注意事项 & 避坑指南 (Precautions & Caveats)

### 3.1 协议校验的差异化处理
* 对于用户输入的 `http://` 或 `https://` 链接，后端需对其进行常规 URL 合法性解析。
* 对于 `feedcraft://` 协议的内部 URI，后端在验证其结构合法（如解析出资源类型和 ID）后，应当直接通过 `feedruntime` 内部加载对应 Recipe 的配置进行解析，避免按常规 HTTP 流程发起外部请求。

### 3.2 避免 Craft 处理的双重叠加 (Double Crafting)
* `feedcraft://recipe/:id` URI 在 `feedruntime` 运行并加载时，已经执行了该 Recipe 定义中的 Craft 脚本。
* 预览工具（如 `/tools/compare`）有时支持额外应用一个 Craft 规则（参数中的 `craft_name`）。在预览已有的 Recipe 时，不应该再在前端向后端透传 `craft_name` 进行处理，否则会导致该 Recipe 的输出被进行二次 Craft 运算，产生不符合预期的干扰结果。

### 3.3 友好的错误提示
* 当预览的 Recipe 不存在、配方配置损坏（YAML/JSON 解析失败）或其依赖的源拉取超时时，后端应当返回清晰、可读的业务层错误信息，而不是直接发生 Panic 或返回模糊的“500 Internal Server Error”。
