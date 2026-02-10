---
title: HTML to RSS
description: Turn any webpage into an RSS feed using the visual selector.
sidebar:
  order: 2
  badge:
    text: new
    variant: success

---

FeedCraft includes a visual **HTML to RSS** tool that allows you to generate selectors for creating RSS feeds from websites that don't provide them natively.

> **Note:** This tool is designed for HTML pages. If you need to process a JSON API, use the [CURL to RSS](/en/guides/advanced/curl-to-rss/) instead.

## Overview

The HTML to RSS tool allows you to:

1.  **Fetch** a webpage's content.
2.  **Select** elements visually to define what constitutes a feed item (title, link, date, content).
3.  **Metadata** Define feed details like title and description.
4.  **Save** the configuration as a Custom Recipe directly.

## How to use

1.  Navigate to **Worktable > HTML to RSS** in the admin dashboard.

### Step 1: Target URL

1.  Enter the full URL of the webpage you want to scrape (e.g., a blog list or news site).
2.  **Enhanced Mode**: Enable this if the site requires JavaScript to load content (uses headless browser).
3.  Click **Fetch and Next**.

### Step 2: Extract Rules

This step allows you to map HTML elements to RSS feed fields.

1.  **Page Preview**: The left pane shows the rendered webpage.
    - **Selection Mode**: When active, clicking elements in the preview will generate selectors.
2.  **CSS Selector (List Item)**:
    - Click **Pick**.
    - In the preview pane, click on a container element that represents a _single article_ or item in the list.
3.  **Field Selectors**:
    - Once the Item Selector is set, you can map the relative fields.
    - **Title**: Click Pick and select the element containing the title text.
    - **Link**: Click Pick and select the element containing the article URL.
    - **Date**: (Optional) Pick the date element.
    - **Description**: (Optional) Pick the element containing the summary.

Click **Run Preview** to test your selectors. If satisfied, click **Next Step**.

### Step 3: Feed Metadata

The tool attempts to auto-extract metadata from the page.

- **Feed Title**: Adjust if necessary.
- **Description**: Feed description.
- **Site Link**: The source URL.

### Step 4: Save Recipe

Review your configuration and save it as a permanent recipe.

- **Recipe Unique ID**: A unique identifier for this feed configuration (e.g., `tech-blog-feed`). If left empty, it will be automatically generated from the feed title.
- **Internal Description**: Notes for yourself about this recipe.

Click **Confirm and Save**. The tool will automatically create a new Custom Recipe with your configuration, which you can manage in the **Custom Recipes** dashboard.
