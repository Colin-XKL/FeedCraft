---
title: 从HTML网页生成RSS
description: 使用可视化选择器将任意网页转换为 RSS 订阅源。
---

FeedCraft 内置了可视化的 **RSS 生成向导 (RSS Generator Wizard)**，允许你为那些没有提供 RSS 的网站创建订阅源。

> **注意：** 此工具专为 HTML 页面设计。如果你需要处理 JSON API，请使用 [JSON RSS 生成器](/zh/guides/advanced/json-rss-generator/)。

## 概览

RSS 生成器允许你：
1.  **抓取 (Fetch)** 网页内容。
2.  **可视化选择 (Select)** 元素来定义什么是订阅源条目（标题、链接、日期、内容）。
3.  **预览 (Preview)** 生成的订阅源条目。
4.  **保存 (Save)** 配置为可复用的食谱 (Recipe)。

## 如何使用

1.  在管理后台导航至 **Tools > RSS Generator**。

### 第一步：目标 URL (Target URL)

1.  输入你想要抓取的网页完整 URL（例如博客列表或新闻网站）。
2.  点击 **Fetch & Next**。

### 第二步：提取规则 (Extract Rules)

此步骤允许你将 HTML 元素映射到 RSS 订阅源字段。

1.  **页面预览 (Page Preview)**：左侧窗格显示渲染后的网页。
    *   **增强模式 (Enhanced Mode)**：如果预览内容与你在浏览器中看到的不同（例如内容缺失），这通常意味着页面使用 JavaScript 渲染。切换“增强模式”以使用无头浏览器来获取和渲染页面。
    *   **选择模式 (Selection Mode)**：当处于激活状态时（由蓝色标签指示），点击预览中的元素将生成选择器。
    *   **键盘快捷键 (Keyboard Shortcuts)**：使用 `Arrow Up` (上箭头) 选择父元素，使用 `Arrow Down` (下箭头) 选择第一个子元素。这对于微调选择非常有用。
2.  **列表项选择器 (List Item Selector) - 必填**：
    -   点击 "CSS Selector" 旁边的 **Pick** 按钮。
    -   在预览窗格中，点击代表列表中*单个文章*或条目的容器元素。
    -   系统将尝试自动计算 CSS 选择器。
3.  **字段选择器 (Field Selectors)**：
    -   设置好列表项后，你可以映射相对字段。
    -   **标题 (Title)**：点击 Pick 并选择包含标题文本的元素。
    -   **链接 (Link)**：点击 Pick 并选择包含文章 URL 的元素（通常是 `<a>` 标签）。
    -   **日期 (Date)**：（可选）选择日期元素。
    -   **描述 (Description)**：（可选）选择包含摘要或正文的元素。
4.  **运行预览 (Run Preview)**：点击此按钮测试你的选择器。解析出的条目将显示在右侧面板中。

### 第三步：订阅源元数据 (Feed Metadata)

为你的新 RSS 订阅源提供元数据：
-   **Feed Title**：订阅源名称。
-   **Feed Description**：简短描述。
-   **Site Link**：原始网站 URL。
-   **Author Info**：（可选）作者姓名和邮箱。

### 第四步：保存食谱 (Save Recipe)

1.  **Recipe Unique ID**：输入此食谱的唯一标识符（例如 `tech-blog-daily`）。这将成为你订阅源 URL 的一部分。
2.  **Internal Description**：供你参考的备注。
3.  点击 **Confirm & Save**。

## 访问你的订阅源

保存后，该食谱将作为 **自定义食谱 (Custom Recipe)** 存储。你可以在 **Custom Recipes** 仪表板中管理它。

你的新 RSS 订阅源通常可以通过类似以下的 URL 访问：
`http://your-feedcraft-instance/rss/custom/{recipe-unique-id}`

你现在可以在你喜欢的 RSS 阅读器中订阅此 URL。
