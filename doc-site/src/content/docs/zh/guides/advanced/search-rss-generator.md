---
title: 搜索转 RSS
description: 使用 AI 提供商通过搜索查询生成 RSS 订阅源。
---

## 前提条件

在使用搜索转 RSS 功能之前，您需要在管理员设置中配置搜索提供商。请参阅 [搜索提供商配置指南](/docs/zh/guides/advanced/customization) 获取设置说明。

FeedCraft 包含一个 **搜索转 RSS (Search to RSS)** 工具，允许你将搜索查询转换为 RSS 订阅源。这对于使用配置的搜索提供商（例如 SearXNG, Bing, Google）追踪新闻、话题或品牌提及非常有用。

## 如何使用

1.  在管理后台导航至 **工作台 > 搜索 转 RSS**。

### 第一步：搜索查询 (Search Query)

1.  输入你的 **搜索查询**（例如 `latest AI news` 或 `SpaceX launches`）。
2.  **Enhanced Mode**: (可选) 开启此选项以使用 AI (LLM) 生成多个优化的搜索查询。这通过扩展你的原始查询来帮助发现更多相关内容。
3.  点击 **预览结果 (Preview Results)** 获取结果。

### 第二步：预览结果 (Preview Results)

系统将使用配置的搜索提供商获取结果。

- 查看找到的项目列表（标题、日期、链接、描述）。
- 如果结果正确，点击 **下一步 (Next Step)**。

### 第三步：Feed 元数据 (Feed Metadata)

自定义此 Feed 在 RSS 阅读器中的显示方式：

- **Feed 标题**：默认为 "搜索: [查询词]"。
- **Feed 描述**：简短描述。
- **站点链接**：搜索结果页面的链接。

### 第四步：保存配方 (Save Recipe)

1.  **配方 ID (Recipe ID)**：为此配方选择一个唯一的标识符（例如 `search-ai-news`）。如果留空，将自动根据订阅源标题生成 (kebab-case)。这将成为你订阅源 URL 的一部分。
2.  **内部描述 (Internal Description)**：关于此配方的备注。
3.  点击 **确认并保存 (Confirm & Save)**。

## 访问你的订阅源

保存后，该配方将作为 **自定义配方 (Custom Recipe)** 存储。你可以在 **Custom Recipes** 仪表板中管理它。

你的新订阅源将通过以下地址访问：
\`http://your-feedcraft-instance/rss/custom/{recipe-unique-id}\`
