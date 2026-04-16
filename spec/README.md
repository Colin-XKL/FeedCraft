# Spec 目录

`spec/` 用于存放 FeedCraft 的核心概念、运行时模型、关键数据流与设计约定。

这里的文档不是临时 proposal，而是后续开发可直接参考的实现基线。

当前文档：

- [Feed Runtime 核心概念与数据流](./feed_runtime_core_concepts.md)
  - 顶层优先级最高的运行时说明
  - 明确 `InputSpec` 是系统顶层统一输入模型
- [SourceConfig 配置模型规范](./source_config.md)
  - 描述 `InputSpec(kind=source)` 下的 Source 子系统配置模型
  - 优先级低于 `feed_runtime_core_concepts.md`；如有概念冲突，以 runtime core concepts 为准
