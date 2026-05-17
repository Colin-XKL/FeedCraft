---
title: JSON to RSS
description: Convert any JSON API response into an RSS feed with jq selectors and optional templates.
sidebar:
  order: 3
  badge:
    text: new
    variant: success
---

FeedCraft includes a **JSON to RSS** tool that allows you to fetch data from JSON APIs, extract fields with `jq`, and optionally post-process them with templates before generating an RSS feed.

## Overview

The JSON to RSS tool helps you:

1.  **Fetch** JSON data from an API endpoint (supporting custom headers and methods).
2.  **Parse** the JSON structure using `jq` syntax, then optionally use templates to build the final RSS fields.
3.  **Metadata** Define feed details like title and description.
4.  **Save** the configuration as a Custom Recipe directly.

## How to use

Navigate to **Worktable > JSON to RSS** in the admin dashboard.

### Step 1: Request Configuration

You need to define how to fetch the JSON data.

- **Import from cURL**: You can paste a `curl` command to automatically populate the URL, method, headers, and body. This is useful if you copy the request from your browser's Developer Tools.
- **Method**: Select `GET` or `POST`.
- **URL**: The API endpoint URL.
- **Headers**: Add any necessary headers (e.g., `Authorization`, `Content-Type`).
- **Request Body**: For POST requests, provide the JSON body.

Click **Fetch and Next** to retrieve the data.

### Step 2: Parsing Rules

Once the JSON is fetched, you will see the raw response in the left panel (visualized as a tree). You can now define selectors to extract feed items.

The tool uses **[jq](https://jqlang.github.io/jq/)** syntax for querying JSON, and it can optionally apply Go templates to the extracted values.

- **List Selector** (Items Iterator): The path to the array of items.
  - Tip: You can click on a node in the tree view to auto-fill selectors.
- **Title Selector**: The path to the item's title _relative to the item object_.
- **Title Template**: Optional. Post-process the extracted title, for example `{{ .Fields.Title | trimSpace }}`.
- **Link Selector**: The path to the item's URL.
- **Link Template**: Optional. Useful when the API only returns an ID, for example `https://some-website.com/article/{{ .Item.id }}`.
- **Date Selector**: (Optional) Path to the publication date.
- **Content Selector**: (Optional) Path to the full content or summary.

#### Using Templates (Optional)

You can use [Go Templates](https://pkg.go.dev/text/template) to further process extracted values.

**Available Variables:**

- `.Fields`: The parsed field values (e.g., `.Fields.Title`, `.Fields.Link`, `.Fields.Date`, `.Fields.Description`).
- `.Item`: The raw JSON item object (e.g., `.Item.id`, `.Item.author.name`).

**Built-in Functions:**

- `trimSpace`: Removes leading and trailing whitespace.
- `trim`: Removes specified leading and trailing characters.
- `default`: Provides a fallback value if the field is empty.

**Examples:**

- **Clean up whitespace in title**: `{{ .Fields.Title | trimSpace }}`
- **Build absolute URLs**: `https://example.com/article/{{ .Item.id }}`
- **Remove specific prefixes**: `{{ .Fields.Description | trim "Prefix: " }}`
- **Fallback values**: `{{ default .Fields.Description "No summary available" }}`

Click **Run Preview** to verify your selectors, then click **Next Step**.

### Step 3: Feed Metadata

Configure the RSS feed details:

- **Feed Title**: The name of your new feed.
- **Description**: A short description.
- **Site Link**: The URL of the original website.
- **Author**: (Optional) Author details.

### Step 4: Save Recipe

Review your configuration and save it as a permanent recipe.

- **Recipe Unique ID**: A unique identifier for this feed configuration (e.g., `my-custom-api-feed`).
  - **Auto-Fill**: This field is automatically populated from the feed title.
  - **Format**: Only lowercase letters, numbers, and hyphens (`[a-z0-9-]`) are allowed.
  - **Refresh**: You can manually regenerate the ID from the title using the refresh button.
- **Internal Description**: Notes for yourself about this recipe.

Click **Confirm and Save**. The tool will automatically create a new Custom Recipe with your configuration, which you can manage in the **Custom Recipes** dashboard.
