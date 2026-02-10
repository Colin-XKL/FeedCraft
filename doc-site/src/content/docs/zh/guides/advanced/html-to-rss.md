---
title: 从HTML网页生成RSS
description: 使用可视化选择器将任意网页转换为 RSS 订阅源。
sidebar:
  order: 2
  badge:
    text: new
    variant: success

---

FeedCraft 内置了可视化的 **从HTML网页生成RSS (HTML to RSS)** 工具，允许你生成选择器，以便为那些没有提供 RSS 的网站创建订阅源。

> **注意：** 此工具专为 HTML 页面设计。如果你需要处理 JSON API，请使用 [从CURL语句生成RSS](/zh/guides/advanced/curl-to-rss/)。

## 概览

HTML to RSS 工具允许你：

1.  **抓取 (Fetch)** 网页内容。
2.  **可视化选择 (Select)** 元素来定义什么是订阅源条目（标题、链接、日期、内容）。
3.  **元数据 (Metadata)**：定义订阅源的标题和描述等详情。
4.  **保存 (Save)**：直接将配置保存为自定义食谱。

## 如何使用

1.  在管理后台导航至 **工作台 > HTML 转 RSS**。

### 第一步：目标 URL (Target URL)

1.  输入你想要抓取的网页完整 URL（例如博客列表或新闻网站）。
2.  **增强模式 (Enhanced Mode)**：如果网站需要 JavaScript 来加载内容（使用无头浏览器），请启用此选项。
3.  点击 **抓取并下一步 (Fetch and Next)**。

### 第二步：提取规则 (Extract Rules)

此步骤允许你将 HTML 元素映射到 RSS 订阅源字段。

1.  **页面预览 (Page Preview)**：左侧窗格显示渲染后的网页。
    - **选择模式 (Selection Mode)**：当处于激活状态时，点击预览中的元素将生成选择器。
2.  **CSS 选择器 (列表项)**：
    - 点击 **选择 (Pick)**。
    - 在预览窗格中，点击代表列表中*单个文章*或条目的容器元素。
3.  **字段选择器 (Field Selectors)**：
    - 设置好列表项后，你可以映射相对字段。
    - **标题 (Title)**：点击选择并选中包含标题文本的元素。
    - **链接 (Link)**：点击选择并选中包含文章 URL 的元素。
    - **日期 (Date)**：（可选）选择日期元素。
    - **描述 (Description)**：（可选）选择包含摘要的元素。

点击 **运行预览 (Run Preview)** 测试你的选择器。如果满意，点击 **下一步 (Next Step)**。

### 第三步：订阅源元数据 (Feed Metadata)

工具会尝试从页面自动提取元数据。

- **订阅源标题 (Feed Title)**：如有必要请调整。
- **描述 (Description)**：订阅源描述。
- **网站链接 (Site Link)**：源 URL。

### 第四步：保存食谱 (Save Recipe)

审查你的配置并将其保存为永久食谱。

- **食谱唯一 ID (Recipe Unique ID)**：此订阅源配置的唯一标识符（例如 `tech-blog-feed`）。如果留空，将自动根据订阅源标题生成。
- **内部描述 (Internal Description)**：关于此食谱的备注。

点击 **确认并保存 (Confirm and Save)**。工具将自动创建一个包含你的配置的新自定义食谱，你可以在 **自定义食谱 (Custom Recipes)** 仪表板中管理它。
