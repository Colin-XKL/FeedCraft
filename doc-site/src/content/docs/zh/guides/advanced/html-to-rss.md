---
title: 从HTML网页生成RSS
description: 使用可视化选择器将任意网页转换为 RSS 订阅源。
---

FeedCraft 内置了可视化的 **HTML to RSS** 工具，允许你生成选择器，以便为那些没有提供 RSS 的网站创建订阅源。

> **注意：** 此工具专为 HTML 页面设计。如果你需要处理 JSON API，请使用 [从CURL语句生成RSS](/zh/guides/advanced/curl-to-rss/)。

## 概览

RSS 生成器允许你：

1.  **抓取 (Fetch)** 网页内容。
2.  **可视化选择 (Select)** 元素来定义什么是订阅源条目（标题、链接、日期、内容）。
3.  **预览 (Preview)** 生成的订阅源条目。
4.  **使用 (Use)** 生成的选择器在自定义频道 (Channel) 中。

## 如何使用

1.  在管理后台导航至 **Tools > HTML to RSS**。

### 第一步：目标 URL (Target URL)

1.  输入你想要抓取的网页完整 URL（例如博客列表或新闻网站）。
2.  点击 **Fetch Page**。

### 第二步：提取规则 (Extract Rules)

此步骤允许你将 HTML 元素映射到 RSS 订阅源字段。

1.  **页面预览 (Page Preview)**：左侧窗格显示渲染后的网页。
    - **选择模式 (Selection Mode)**：当处于激活状态时，点击预览中的元素将生成选择器。
    - **键盘快捷键 (Keyboard Shortcuts)**：使用 `Arrow Up` (上箭头) 选择父元素，使用 `Arrow Down` (下箭头) 选择第一个子元素。这对于微调选择非常有用。
2.  **列表项选择器 (List Item Selector) - 必填**：
    - 点击 "List Item Selector" 旁边的 **Pick** 按钮。
    - 在预览窗格中，点击代表列表中*单个文章*或条目的容器元素。
    - 系统将尝试自动计算 CSS 选择器。
3.  **字段选择器 (Field Selectors)**：
    - 设置好列表项后，你可以映射相对字段。
    - **标题选择器 (Title Selector)**：点击 Pick 并选择包含标题文本的元素。
    - **链接选择器 (Link Selector)**：点击 Pick 并选择包含文章 URL 的元素（通常是 `<a>` 标签）。
    - **日期选择器 (Date Selector)**：（可选）选择日期元素。
    - **内容选择器 (Content Selector)**：（可选）选择包含摘要或正文的元素。
4.  **预览 RSS 条目 (Preview RSS Items)**：点击此按钮测试你的选择器。解析出的条目将显示在右侧面板中。

### 第三步：在自定义频道中使用 (Use in Channel)

一旦你在预览中验证了选择器：

1.  记下生成的选择器（列表项、标题、链接、日期、内容）。
2.  导航至 **Channels** > **Create**。
3.  使用这些选择器配置你的频道以解析目标网站。
