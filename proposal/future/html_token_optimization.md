# HTML Token Optimization for LLM Input

> 状态：规划中

## 1. 需求背景

当前 FeedCraft 已经有统一的 LLM 调用入口与内容预处理链路，但 HTML 进入 LLM 前的清洗仍然比较粗糙。很多页面会带上大量对 RSS 场景无价值、但会显著消耗 token 的内容，例如：

- `script` / `style` / `noscript` 等无关节点
- `class` / `style` / `id` / `aria-*` / `data-*` 等低价值属性
- 很长的 `href` / `src` / `srcset`
- base64 `data:` 图片
- 过多空白、缩进与样式噪音

这些内容会带来两个问题：

1. **增加 LLM 成本**：同样的正文语义会占用更多 token。
2. **干扰模型理解**：无关 HTML 噪音会稀释正文、图片、链接等真正重要的信息。

因此需要在现有架构中增加一层简洁、稳定、可配置的 HTML 优化逻辑，在不破坏主要语义的前提下，尽量缩小送入 LLM 的内容体积。

## 2. 目标

本方案希望实现一套面向 LLM 输入的 HTML 优化机制，满足以下目标：

- 尽量复用现有 `ProcessContent` 预处理入口，不重新发明新链路
- 通过 DOM 级处理删除无意义 HTML 内容，而不是只靠正则替换
- 支持不同 Craft 使用不同优化等级
- 保持配置模型足够小，但能表达“保留多少原始内容”这种关键差异
- 让 tag / attr 规则集中定义，方便后续维护
- 为 placeholder 替换与恢复提供简单、局部的机制
- 补充单独测试，确保优化结果稳定、可预期

## 3. 适用场景与等级差异

不同 LLM 场景对 HTML 保真度的要求并不一样，因此这里不适合只用一个 `bool` 开关。

### 3.1 更激进的场景

例如：

- summary
- llm filter
- 条件判断类 craft

这类场景重点是提取正文语义或做分类判断，不需要保留太多原始结构。对于它们，可以更激进地：

- 删除更多无关节点
- 移除大部分低价值属性
- 对链接和图片做更强压缩或直接去除

### 3.2 中等保留的场景

例如：

- 常规 translate
- beautify

这些场景仍然希望模型理解原文结构，且输出内容最好保留链接、图片等信息，因此应该：

- 保留主要结构标签
- 保留关键资源属性
- 仅去掉明显无意义的噪音
- 仅替换超长 URL / `data:` URI

### 3.3 最保守的场景

例如：

- immersive-translate

这种场景强调尽量保留原始格式、链接、图片以及更多内容组织形式，因此优化应该尽量保守，只做：

- 明确无意义节点移除
- 低价值属性清理
- 极端长字段压缩
- 空白压缩

## 4. 设计原则

### 4.1 单一入口

优先把 HTML 优化纳入现有的 `internal/util/content_processor.go`，作为 `ProcessContent` 的一部分，这样：

- 大部分 LLM 调用链无需额外重构
- 现有 `RemoveLinks` / `RemoveImage` / `ConvertToMd` 能继续复用
- 逻辑集中，后续更容易维护

### 4.2 配置驱动，而不是 craft 分支硬编码

不希望在代码里到处写：

- `if summary { ... }`
- `if immersiveTranslate { ... }`

更合适的方式是：

- `ContentProcessOption` 增加一个 HTML optimize config struct
- caller 只声明“我想要哪种保留等级”
- 具体 tag/attr/placeholder 规则由 optimizer 内部统一决定

### 4.3 规则集中、可维护

attr 和 tag 规则应该集中定义，不应散落在 DOM 遍历逻辑里。后续新增规则时，最好只修改一处规则表或 helper。

### 4.4 保持 v1 简洁

这次优化是为了减少 token，不是为了做完整 HTML sanitizer，也不是为了构建一个高度可配置的通用清洗框架。v1 只处理最有收益、最确定的部分。

## 5. 推荐的数据结构

建议在 `internal/util/` 中增加一个小型配置模型，例如：

- `ContentProcessOption` 增加 `OptimizeHTML *HTMLOptimizeConfig`
- `HTMLOptimizeConfig` 只表达少数关键维度，例如：
  - 优化等级 / preservation profile
  - 是否保留链接
  - 是否保留图片
  - 长 URL / `data:` URI 的替换阈值（可内部默认）

更推荐的做法是：

- 对外暴露少量 profile 或 level
- 内部再映射成实际规则集

这样可以同时满足：

- 外部调用简洁
- 内部实现可演进
- 不会让调用方承担太多细节决策

## 6. HTML 优化的核心步骤

建议优化器按以下顺序工作：

1. 判断输入是否像 HTML；非 HTML 内容直接跳过
2. 用 DOM 解析 HTML
3. 删除无意义节点
4. 清理低价值属性
5. 根据配置决定是否保留图片、链接等结构
6. 对超长属性值做 placeholder 替换
7. 序列化 HTML
8. 压缩空白

### 6.1 建议直接删除的节点

v1 可先覆盖这些明显无价值的元素：

- `script`
- `style`
- `noscript`
- `template`
- `iframe`

这些元素通常不会帮助 RSS 抽取、摘要、翻译或分类，保留它们只会浪费 token。

### 6.2 建议清理的属性

低价值属性建议集中通过规则清理，例如：

- `class`
- `id`
- `style`
- `aria-*`
- `data-*`
- `on*` 事件属性

同时保留可能有语义价值的属性，具体是否保留也可受 profile 影响，例如：

- `href`
- `src`
- `srcset`
- `alt`
- `title`

### 6.3 超长字段替换

对以下字段做可逆 placeholder 压缩：

- 很长的 `href`
- 很长的 `src`
- 很长的 `srcset`
- base64 `data:` URI

例如替换为：

- `__FC_PH_URL_0001__`
- `__FC_PH_DATA_0002__`

并维护 request-scoped map。这样可以避免：

- 把极长字符串直接发给 LLM
- 使用全局共享状态
- 后续恢复时依赖不透明上下文

## 7. placeholder 恢复策略

v1 不需要把恢复逻辑强行塞进所有 LLM 流程里。

推荐做法：

- 优化器返回优化后的 HTML 和 placeholder map
- 只有确实需要恢复的流程才使用 restore helper
- placeholder map 生命周期限定在单次调用上下文内

这样能保持实现足够简单，同时给像 `beautify` 这种对原始 URL 保真度更敏感的流程留下扩展空间。

## 8. 与现有代码的集成点

### 8.1 统一入口

优先复用：

- `internal/adapter/llm.go`
- `internal/util/content_processor.go`

也就是继续通过 `ProcessContent` 统一处理大部分 LLM 输入。

### 8.2 主要调用方

后续实现时，至少应考虑这些现有路径：

- `internal/craft/common_llm_logic.go`
  - 使用更激进的 profile
  - 适合 filter / condition / summary 一类语义型任务

- `internal/craft/translate.go`
  - 使用更保守的 profile
  - 常规翻译需要保留更多结构与资源

- immersive translate
  - 使用最保守的 profile
  - 尽量保留链接、图片和格式

- `internal/craft/beautify.go`
  - 当前会直接把原始 HTML 拼进 prompt
  - 需要在 prompt 构造前引入同一套优化逻辑

## 9. 规则维护方式

为了让代码保持简洁、清晰、优雅，建议使用以下形式管理规则：

- 一组集中定义的 removable tags
- 一组集中定义的 removable attrs
- 一组 prefix-based removable attr rules
- 一组 profile-aware preserved attrs

DOM 遍历时只调用这些 helper，不在遍历逻辑里堆积大量 if/else。

这样做的优点：

- 可读性更好
- 修改规则时影响范围小
- 容易为不同 profile 扩展行为
- 更容易测试每条规则的预期

## 10. v1 明确不做的事情

为了保持范围收敛，这个方案暂时不处理：

- 完整 HTML 安全清洗
- 复杂 CSS 可见性推断
- 语义重排或正文重写
- 过度细粒度的用户可配置规则
- 所有 URL 一律替换
- 对所有 LLM 返回结果做全局 placeholder 自动恢复

## 11. 测试策略

这部分必须单独补充测试，不能只依赖集成路径顺带覆盖。

### 11.1 独立 optimizer 单测

建议增加单独测试文件，例如：

- `internal/util/content_processor_test.go`
- 或 `internal/util/html_optimize_test.go`

重点验证：

- 无意义节点是否被移除
- 低价值属性是否被清理
- 关键属性是否按 profile 保留
- whitespace 是否被正确压缩
- 长 URL / `data:` URI 是否被替换为 placeholder
- placeholder 恢复是否正确
- 非 HTML 输入是否安全透传

### 11.2 profile 差异测试

需要专门验证不同优化等级的行为差异，而不是只测单一输出。

例如：

- aggressive profile 是否删除更多噪音
- preserve profile 是否保留图片和链接
- immersive profile 是否比 preserve profile 保留更多原始结构

### 11.3 `ProcessContent` 组合测试

继续验证它和现有逻辑的组合行为：

- optimize + `ConvertToMd`
- optimize + `RemoveImage`
- optimize + `RemoveLinks`
- 多个步骤叠加时顺序是否稳定

### 11.4 调用路径 smoke test

对主要 craft 调用路径补充轻量 smoke coverage，确保它们确实选用了正确 profile。

## 12. 预期收益

如果实现得当，这套机制会带来以下收益：

- 减少 LLM 输入 token 消耗
- 提高 prompt 中有效语义密度
- 让 summary / filter 等场景更稳定
- 让 translate / immersive-translate 在保真和成本之间更可控
- 为后续更细致的 HTML 内容优化打下可扩展基础

## 13. 总结

这项工作的关键不在于“删得越多越好”，而在于：

- 用统一入口做预处理
- 用小而清晰的 config 表达不同保留等级
- 用集中规则保持 tag / attr 逻辑可维护
- 用独立测试确保优化结果稳定

如果这几个点把握好，FeedCraft 就能在不明显增加架构复杂度的前提下，让 HTML 进入 LLM 前变得更轻、更干净、更适配不同 craft 的需求。
