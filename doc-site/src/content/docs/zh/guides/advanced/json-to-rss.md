---
title: 从 JSON 生成 RSS
description: 使用 jq 提取字段，并通过可选模板将任意 JSON API 响应转换为 RSS 订阅源。
sidebar:
  order: 3
  badge:
    text: new
    variant: success
---

FeedCraft 包含一个 **从 JSON 生成 RSS (JSON to RSS)** 工具，允许你从 JSON API 获取数据，先用 `jq` 提取字段，再通过可选模板将其转换为 RSS 订阅源。

## 概览

JSON RSS 生成器可以帮助你：

1.  **抓取 (Fetch)**：从 API 端点抓取 JSON 数据（支持自定义请求头和方法）。
2.  **解析 (Parse)**：使用 `jq` 语法解析 JSON 结构，并可通过模板拼接或清理 RSS 字段。
3.  **元数据 (Metadata)**：定义订阅源的标题和描述等详情。
4.  **保存 (Save)**：直接将配置保存为自定义食谱。

## 如何使用

在管理后台导航至 **工作台 > JSON 转 RSS**。

### 第一步：请求配置 (Request Configuration)

你需要定义如何获取 JSON 数据。

- **从 cURL 导入 (Import from cURL)**：你可以粘贴 `curl` 命令来自动填充 URL、方法、请求头和请求体。这在你从浏览器开发者工具复制请求时非常有用。
- **方法 (Method)**：选择 `GET` 或 `POST`。
- **URL**：API 端点 URL。
- **Headers**：添加任何必要的请求头（例如 `Authorization`, `Content-Type`）。
- **请求体 (Request Body)**：对于 POST 请求，提供 JSON 请求体。

点击 **抓取并下一步 (Fetch and Next)** 来获取数据。

### 第二步：JQ 解析规则 (Parsing Rules)

获取到 JSON 后，你将在左侧面板看到以树形可视化的响应。现在你可以定义选择器来提取订阅源条目。

该工具使用 **[jq](https://jqlang.github.io/jq/)** 语法来查询 JSON，并支持对提取结果再做一层 Go template 加工。

- **列表选择器 (Items Iterator)**：条目数组的路径。
  - 提示：你可以点击树视图中的节点来自动填充选择器。
- **标题选择器 (Title Selector)**：条目标题的路径（相对于条目对象）。
- **标题模板 (Title Template)**：（可选）对提取到的标题做进一步处理，例如 `{{ .Fields.Title | trimSpace }}`。
- **链接选择器 (Link Selector)**：条目 URL 的路径。
- **链接模板 (Link Template)**：（可选）当接口没有完整链接时，可以拼接，例如 `https://some-website.com/article/{{ .Item.id }}`。
- **日期选择器 (Date Selector)**：（可选）发布日期的路径。
- **内容选择器 (Content Selector)**：（可选）完整内容或摘要的路径。

#### 使用模板 (可选)

你可以使用 [Go Templates](https://pkg.go.dev/text/template) 语法对提取的值进行进一步处理。

**可用变量：**

- `.Fields`：已解析的字段值（例如 `.Fields.Title`, `.Fields.Link`, `.Fields.Date`, `.Fields.Description`）。
- `.Item`：原始 JSON 列表项对象（例如 `.Item.id`, `.Item.author.name`）。

**内置函数：**

- `trimSpace`：移除首尾的空白字符。
- `trim`：移除首尾指定的字符。
- `default`：如果字段为空，提供一个默认值。

**示例：**

- **清理标题空白字符**：`{{ .Fields.Title | trimSpace }}`
- **拼接完整链接**：`https://example.com/article/{{ .Item.id }}`
- **移除特定前缀**：`{{ .Fields.Description | trim "Prefix: " }}`
- **默认值兜底**：`{{ default .Fields.Description "暂无摘要" }}`

点击 **运行预览 (Run Preview)** 验证你的选择器，然后点击 **下一步 (Next Step)**。

### 第三步：订阅源元数据 (Feed Metadata)

配置 RSS 订阅源详情：

- **订阅源标题 (Feed Title)**：你的新订阅源名称。
- **描述 (Description)**：简短描述。
- **网站链接 (Site Link)**：原始网站 URL。
- **作者 (Author)**：（可选）作者详情。

### 第四步：保存食谱 (Save Recipe)

审查你的配置并将其保存为永久食谱。

- **食谱唯一 ID (Recipe Unique ID)**：此订阅源配置的唯一标识符（例如 `my-custom-api-feed`）。
  - **自动填充 (Auto-Fill)**：该字段将根据订阅源标题自动填充。
  - **格式 (Format)**：仅允许小写字母、数字和连字符 (`[a-z0-9-]`)。
  - **刷新 (Refresh)**：你可以使用刷新按钮根据标题手动重新生成 ID。
- **内部描述 (Internal Description)**：关于此食谱的备注。

点击 **确认并保存 (Confirm and Save)**。工具将自动创建一个包含你的配置的新自定义食谱，你可以在 **自定义食谱 (Custom Recipes)** 仪表板中管理它。
