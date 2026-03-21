---
title: 系统工具
description: 内置的调试、对比和系统健康检查工具。
sidebar:
  order: 5
---

FeedCraft 提供了一些内置工具来帮助您调试 RSS 源并监控系统健康状况。您可以在管理后台的 **工具 (Tools)** 菜单下访问这些工具。

## RSS 预览 (RSS Viewer)

**RSS 预览** (Feed Viewer) 允许您按照 FeedCraft 的解析方式预览任何 RSS 源。

- **使用方法**:
  1. 导航至 **工具 > RSS 预览**。
  2. 输入一个 RSS/Atom 地址。
  3. 点击 **预览 (Preview)**。
- **目的**: 在设置配方 (Recipe) 之前，验证 FeedCraft 是否能够成功抓取和解析某个 Feed。
- **注意**: 预览器默认使用 `proxy` 工艺，它只是简单地抓取 Feed 而不进行修改。

## Feed 对比 (Feed Compare)

**Feed 对比** 工具让您可以直观地看到某个 Craft（Atom 或 Flow）对 Feed 的处理效果。

- **使用方法**:
  1. 导航至 **工具 > Feed 对比**。
  2. 输入原始 RSS Feed 地址。
  3. 选择一个 **FlowCraft** 或 **AtomCraft**。
  4. 点击 **对比 (Compare)**。
- **输出**: 工具显示两列：
  - **左侧**: 原始 Feed 内容。
  - **右侧**: 经过选定 Craft 处理后的 Feed 内容。
- **应用场景**: 非常适合测试新的翻译流程或摘要提示词，而无需创建永久配方。

## Craft 依赖检查 (Craft Dependencies)

**Craft 依赖检查** (System Health) 工具可视化您的 Recipes、FlowCrafts 和 AtomCrafts 之间的内部关系。

- **使用方法**:
  1. 导航至 **工具 > Craft 依赖检查**。
  2. 点击 **分析 Craft 依赖 (Analyze Craft Dependencies)**。
- **功能**:
  - 生成所有依赖关系的树状视图。
  - **健康检查**: 自动检测丢失的依赖项（例如，指向已删除 FlowCraft 的 Recipe）。
  - **视觉指示**: Recipes、Flows、Atoms 和丢失组件使用不同颜色标识。

> **提示:** 如果遇到 "Craft not found" 等错误，可以使用此工具追踪配置中的断链。

## 调试工具 (Debug Tools)

### LLM 调试 (LLM Debug)

用于测试 LLM 配置的沙盒。您可以向配置的 LLM 提供商发送测试提示，以验证连接和模型响应。

### 广告检测调试 (Ad Check Debug)

用于针对特定内容测试 "忽略软文 (Ignore Advertorial)" 过滤逻辑的专用工具，以了解文章被过滤的原因。
