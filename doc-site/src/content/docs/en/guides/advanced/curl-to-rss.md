---
title: CURL to RSS
description: Convert any JSON API response into an RSS feed using jq selectors.
---

FeedCraft includes a **CURL to RSS** tool that allows you to fetch data from JSON APIs and transform it into an RSS feed using `jq` selectors.

## Overview

The CURL to RSS tool helps you:

1.  **Fetch** JSON data from an API endpoint (supporting custom headers and methods).
2.  **Parse** the JSON structure using `jq` syntax to map fields to RSS items.
3.  **Preview** the generated feed to verify your selectors.

## How to use

Navigate to **Tools > CURL to RSS** in the admin dashboard.

### Step 1: Request Configuration

You need to define how to fetch the JSON data.

- **Import from Curl**: You can paste a `curl` command to automatically populate the URL, method, headers, and body. This is useful if you copy the request from your browser's Developer Tools.
- **Method**: Select `GET` or `POST`.
- **URL**: The API endpoint URL.
- **Headers**: Add any necessary headers (e.g., `Authorization`, `Content-Type`).
- **Request Body**: For POST requests, provide the JSON body.

Click **Fetch JSON** to retrieve the data.

### Step 2: JQ Parsing Rules

Once the JSON is fetched, you will see the raw response in the left panel. You can now define selectors to extract feed items.

The tool uses **[jq](https://jqlang.github.io/jq/)** syntax for querying JSON.

- **List Selector**: The path to the array of items.
  - Example: `.items[]` or `.data.posts[]` or just `.` if the root is an array.
- **Title Selector**: The path to the item's title _relative to the item object_.
  - Example: `.title` or `.attributes.name`.
- **Link Selector**: The path to the item's URL.
  - Example: `.url` or `.permalink`.
- **Date Selector**: (Optional) Path to the publication date.
- **Content Selector**: (Optional) Path to the full content or summary.

### Step 3: Preview

Click **Preview RSS** to see how your selectors work. The parsed items will appear in the list below.

## Saving Your Recipe

Currently, the CURL to RSS is a tool for **finding and testing** the correct selectors.

To save your configuration as a permanent feed:

1.  Copy your **URL** and **Selector** values.
2.  Go to **Channels** > **Create**.
3.  Select **Source Type**: `JSON`.
4.  Paste your configuration into the **Source Config** JSON format:

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
