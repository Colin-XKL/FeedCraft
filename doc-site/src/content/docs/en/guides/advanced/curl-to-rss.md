---
title: CURL to RSS
description: Convert any JSON API response into an RSS feed using jq selectors.
---

FeedCraft includes a **CURL to RSS** tool that allows you to fetch data from JSON APIs and transform it into an RSS feed using `jq` selectors.

## Overview

The CURL to RSS tool helps you:

1.  **Fetch** JSON data from an API endpoint (supporting custom headers and methods).
2.  **Parse** the JSON structure using `jq` syntax to map fields to RSS items.
3.  **Metadata** Define feed details like title and description.
4.  **Save** the configuration as a Custom Recipe directly.

## How to use

Navigate to **Tools > CURL to RSS** in the admin dashboard.

### Step 1: Request Configuration

You need to define how to fetch the JSON data.

- **Import from Curl**: You can paste a `curl` command to automatically populate the URL, method, headers, and body. This is useful if you copy the request from your browser's Developer Tools.
- **Method**: Select `GET` or `POST`.
- **URL**: The API endpoint URL.
- **Headers**: Add any necessary headers (e.g., `Authorization`, `Content-Type`).
- **Request Body**: For POST requests, provide the JSON body.

Click **Fetch and Next** to retrieve the data.

### Step 2: Parsing Rules

Once the JSON is fetched, you will see the raw response in the left panel (visualized as a tree). You can now define selectors to extract feed items.

The tool uses **[jq](https://jqlang.github.io/jq/)** syntax for querying JSON.

- **List Selector** (Items Iterator): The path to the array of items.
  - Tip: You can click on a node in the tree view to auto-fill selectors.
- **Title Selector**: The path to the item's title _relative to the item object_.
- **Link Selector**: The path to the item's URL.
- **Date Selector**: (Optional) Path to the publication date.
- **Content Selector**: (Optional) Path to the full content or summary.

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
  - **Tip**: If you leave this empty, it will be automatically generated from the feed title.
- **Internal Description**: Notes for yourself about this recipe.

Click **Confirm and Save**. The tool will automatically create a new Custom Recipe with your configuration, which you can manage in the **Custom Recipes** dashboard.
