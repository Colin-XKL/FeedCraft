---
title: HTML to RSS
description: Turn any webpage into an RSS feed using the visual selector.
---

FeedCraft includes a visual **HTML to RSS** tool that allows you to generate selectors for creating RSS feeds from websites that don't provide them natively.

> **Note:** This tool is designed for HTML pages. If you need to process a JSON API, use the [CURL to RSS](/en/guides/advanced/curl-to-rss/) instead.

## Overview

The HTML to RSS tool allows you to:

1.  **Fetch** a webpage's content.
2.  **Select** elements visually to define what constitutes a feed item (title, link, date, content).
3.  **Preview** the generated feed items immediately.
4.  **Use** the generated selectors in a Custom Recipe.

## How to use

1.  Navigate to **Tools > HTML to RSS** in the admin dashboard.

### Step 1: Target URL

1.  Enter the full URL of the webpage you want to scrape (e.g., a blog list or news site).
2.  Click **Fetch Page**.

### Step 2: Extract Rules

This step allows you to map HTML elements to RSS feed fields.

1.  **Page Preview**: The left pane shows the rendered webpage.
    - **Selection Mode**: When active (indicated by button text), clicking elements in the preview will generate selectors.
    - **Keyboard Shortcuts**: Use `Arrow Up` to select the parent element and `Arrow Down` to select the first child element. This is useful for fine-tuning your selection.
2.  **List Item Selector (Required)**:
    - Click the **Pick** button next to "List Item Selector".
    - In the preview pane, click on a container element that represents a _single article_ or item in the list.
    - The system will attempt to auto-calculate a CSS selector.
3.  **Field Selectors**:
    - Once the Item Selector is set, you can map the relative fields.
    - **Title Selector**: Click Pick and select the element containing the title text.
    - **Link Selector**: Click Pick and select the element containing the article URL (usually an `<a>` tag).
      > **Tip:** The extractor is smart! If you select an element that isn't a link (e.g., a `div` or `span`), it will automatically look for a link in the parent or child elements. If no link is found, it will try to use the text content if it looks like a URL. (See [ea6c69e](https://github.com/Colin-XKL/FeedCraft/commit/ea6c69eb55e114ccf7c32c7cc100e5a6514ea70b))
    - **Date Selector**: (Optional) Pick the date element.
    - **Content Selector**: (Optional) Pick the element containing the summary or full content.
4.  **Preview RSS Items**: Click this button to test your selectors. The parsed items will appear in the right-hand panel.

### Step 3: Use in Custom Recipe

Once you have verified the selectors in the preview:

1.  Note down the generated selectors (Item, Title, Link, Date, Content).
2.  Navigate to **Custom Recipes** > **Create**.
3.  Configure your recipe using these selectors to parse the target website.
