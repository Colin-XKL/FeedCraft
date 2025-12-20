---
title: JSON RSS 生成器
description: 使用 jq 选择器将任意 JSON API 响应转换为 RSS 订阅源。
---

FeedCraft 包含一个 **JSON RSS 生成器** 工具，允许你从 JSON API 获取数据并使用 `jq` 选择器将其转换为 RSS 订阅源。

## 概览

JSON RSS 生成器可以帮助你：
1.  **抓取 (Fetch)**：从 API 端点抓取 JSON 数据（支持自定义请求头和方法）。
2.  **解析 (Parse)**：使用 `jq` 语法解析 JSON 结构，将字段映射到 RSS 条目。
3.  **预览 (Preview)**：预览生成的订阅源以验证你的选择器。

## 如何使用

在管理后台导航至 **Tools > JSON RSS Generator**。

### 第一步：请求配置 (Request Configuration)

你需要定义如何获取 JSON 数据。

*   **Import from Curl**：你可以粘贴 `curl` 命令来自动填充 URL、方法、请求头和请求体。这在你从浏览器开发者工具复制请求时非常有用。
*   **Method**：选择 `GET` 或 `POST`。
*   **URL**：API 端点 URL。
*   **Headers**：添加任何必要的请求头（例如 `Authorization`, `Content-Type`）。
*   **Request Body**：对于 POST 请求，提供 JSON 请求体。

点击 **Fetch JSON** 来获取数据。

### 第二步：JQ 解析规则 (JQ Parsing Rules)

获取到 JSON 后，你将在左侧面板看到原始响应。现在你可以定义选择器来提取订阅源条目。

该工具使用 **[jq](https://jqlang.github.io/jq/)** 语法来查询 JSON。

*   **列表选择器 (List Selector)**：条目数组的路径。
    *   例如：`.items[]` 或 `.data.posts[]`，如果根就是数组则为 `.`。
*   **标题选择器 (Title Selector)**：条目标题的路径（相对于条目对象）。
    *   例如：`.title` 或 `.attributes.name`。
*   **链接选择器 (Link Selector)**：条目 URL 的路径。
    *   例如：`.url` 或 `.permalink`。
*   **日期选择器 (Date Selector)**：（可选）发布日期的路径。
*   **内容选择器 (Content Selector)**：（可选）完整内容或摘要的路径。

### 第三步：预览 (Preview)

点击 **Preview RSS** 查看你的选择器效果。解析出的条目将显示在下方的列表中。

## 保存你的食谱

目前，JSON RSS 生成器是一个用于 **查找和测试** 正确选择器的工具。

要将你的配置保存为永久订阅源：

1.  复制你的 **URL** 和 **Selector** 值。
2.  前往 **Custom Recipes** > **Create**。
3.  选择 **Source Type**：`JSON`。
4.  将你的配置粘贴到 **Source Config** JSON 格式中：

```json
{
  "http_fetcher": {
    "url": "https://api.example.com/posts",
    "method": "GET",
    "headers": {
      "Authorization": "Bearer token"
    }
  },
  "json_parser": {
    "list_selector": ".data[]",
    "title_selector": ".title",
    "link_selector": ".url",
    "content_selector": ".body"
  }
}
```
