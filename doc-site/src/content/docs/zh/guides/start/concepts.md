---
title: 核心概念
sidebar:
  order: 2
---

在深入使用 FeedCraft 之前，了解以下三个核心概念将非常有帮助。

## 原子工艺 (AtomCraft)

**AtomCraft** 是最小的处理单元。除了内置的原子工艺（如 `translate-title`, `fulltext`），你可以基于模版创建自定义的原子工艺。

**示例：自定义翻译 Prompt**
你可以基于 `translate-content` 模版创建一个名为 `translate-to-french` 的新原子工艺，并在参数中填入自定义的 Prompt，指示 AI 将内容翻译成法语。

## 组合工艺 (FlowCraft)

**FlowCraft** 是多个 AtomCraft 的组合序列。这允许你将多个操作串联起来。

**示例：全文 + 摘要 + 翻译**
你可以定义一个名为 `digest-and-translate` 的组合工艺，包含以下步骤：

1.  `fulltext` (提取正文)
2.  `summary` (生成摘要)
3.  `translate-content` (翻译内容)

### 管理组合工艺

在管理后台导航至 **工作台 (Worktable) > 组合工艺 (FlowCraft)** 页面创建和管理组合工艺。
编辑器允许你添加原子工艺并安排它们的执行顺序。使用箭头按钮 (⬆️/⬇️) 调整顺序，或使用垃圾桶图标将其从流程中移除。

## 配方 (Recipe)

**Recipe** 将特定的 RSS 源 URL 与某个 原子工艺 (AtomCraft) 或组合工艺 (FlowCraft) 绑定。这允许你创建一个持久化的、经过定制的订阅源 URL。

**管理配方：**
在管理后台导航至 **工作台 (Worktable) > 自定义配方 (Custom Recipe)** 页面，你可以管理所有已创建的食谱。

- **创建 (Create)**：绑定新的 URL 和工艺。
- **预览 (Preview)**：点击预览按钮，直接在内置的 Feed Viewer 中查看生成的效果。
- **复制链接 (Copy Link)**：点击复制图标获取完整的订阅 URL。

**示例：**

- **输入 URL：** `https://news.ycombinator.com/rss`
- **处理器：** `digest-and-translate` (上面创建的工作流)
- **结果：** 你会得到一个新的 FeedCraft URL，订阅它即可获得带全文、摘要和翻译的 Hacker News。
