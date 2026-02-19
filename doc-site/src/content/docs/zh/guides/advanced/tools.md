---
title: 故障排除与工具
description: 了解如何使用 FeedCraft 内置的故障排除和实用工具。
sidebar:
  order: 50
---

FeedCraft 在管理后台内置了多种工具，帮助你调试、验证和监控你的 RSS Feed。

## RSS 预览 (Feed Viewer)

**RSS 预览** 工具允许你以清晰易读的格式预览任何 RSS Feed。这对于验证 Feed 是否有效以及检查其内容非常有用，而无需使用外部阅读器。

- **位置**: **工具 (Tools) > RSS 预览**
- **使用方法**:
  1.  输入你想要预览的 RSS Feed URL。
  2.  点击 **预览**。
  3.  Feed 内容将显示在下方，包括标题、描述和具体文章列表。

## Feed 对比 (Feed Compare)

**Feed 对比** 工具允许你并排比较两个 RSS Feed。这在调试 FlowCraft 时特别有帮助，可以直观地看到处理步骤如何改变了原始 Feed。

- **位置**: **工具 (Tools) > Feed 对比**
- **使用方法**:
  1.  输入原始 Feed 的 **来源链接**。
  2.  从下拉菜单中选择一个 **FlowCraft** (例如 `translate-title`)。
  3.  点击 **对比**。
  4.  工具将在左侧显示“原始 Feed”，在右侧显示“处理后的 Feed”，方便你发现差异。

## 系统健康 (System Health)

**系统健康** 仪表盘可视化展示了 FeedCraft 实例的内部依赖关系图。它显示了 Recipe (配方)、FlowCraft (流程) 和 AtomCraft (原子) 之间的关系。

- **位置**: **工具 (Tools) > 系统健康** (或 Craft 依赖检查)
- **用途**:
  - 验证自定义 Recipe 的所有组件是否存在。
  - 检测缺失的依赖项（例如，Recipe 引用了已删除的 FlowCraft）。
  - 识别可能导致问题的循环依赖。
- **使用方法**:
  1.  点击 **开始分析**。
  2.  查看树状图。健康的组件会标记其类型（Recipe, FlowCraft, AtomCraft）。缺失的组件将标记为红色。

## URL 生成器

**URL 生成器** 帮助你轻松创建 FeedCraft 订阅链接。它还具有“解析模式”，可以反向解析现有的 URL。

- **位置**: **仪表盘 > 快速开始**
- **功能**:
  - **生成**: 选择 FlowCraft 并输入来源 URL，即可获取订阅链接。
  - **解析**: 粘贴 FeedCraft URL 以查看其组成部分（使用的 FlowCraft、原始来源等）。

## 依赖服务 (Dependency Services)

若要监控外部服务（如 Redis, Browserless, LLM）的状态，请使用 **依赖服务状态** 仪表盘。

- **位置**: **设置 > 依赖服务状态**
- **另请参阅**: [高级定制](/zh/guides/advanced/customization/#依赖服务-dependency-services)
