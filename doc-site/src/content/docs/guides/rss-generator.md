---
title: RSS Generator
description: Turn any webpage into an RSS feed using the visual selector.
---

FeedCraft includes a visual **RSS Generator Wizard** that allows you to create RSS feeds from websites that don't provide them natively.

## Overview

The RSS Generator allows you to:
1.  **Fetch** a webpage's content.
2.  **Select** elements visually to define what constitutes a feed item (title, link, date, content).
3.  **Preview** the generated feed items immediately.
4.  **Save** the configuration as a reusable Recipe.

## How to use

1.  Navigate to **Tools > RSS Generator** in the admin dashboard.

### Step 1: Target URL

1.  Enter the full URL of the webpage you want to scrape (e.g., a blog list or news site).
2.  **(Optional) Enhanced Mode**: Toggle "Enhanced Mode (Browserless)" if the target site relies heavily on JavaScript to render content. This uses a headless browser to fetch the page.
3.  Click **Fetch & Next**.

### Step 2: Extract Rules

This step allows you to map HTML elements to RSS feed fields.

1.  **Page Preview**: The left pane shows the rendered webpage.
    *   **Selection Mode**: When active (indicated by a blue tag), clicking elements in the preview will generate selectors.
    *   **Keyboard Shortcuts**: Use `Arrow Up` to select the parent element and `Arrow Down` to select the first child element. This is useful for fine-tuning your selection.
2.  **List Item Selector (Required)**:
    -   Click the **Pick** button next to "CSS Selector".
    -   In the preview pane, click on a container element that represents a *single article* or item in the list.
    -   The system will attempt to auto-calculate a CSS selector.
3.  **Field Selectors**:
    -   Once the Item Selector is set, you can map the relative fields.
    -   **Title**: Click Pick and select the element containing the title text.
    -   **Link**: Click Pick and select the element containing the article URL (usually an `<a>` tag).
    -   **Date**: (Optional) Pick the date element.
    -   **Description**: (Optional) Pick the element containing the summary or full content.
4.  **Run Preview**: Click this button to test your selectors. The parsed items will appear in the right-hand panel.

### Step 3: Feed Metadata

Provide the metadata for your new RSS feed:
-   **Feed Title**: The name of the feed.
-   **Feed Description**: A short description.
-   **Site Link**: The URL of the original website.
-   **Author Info**: (Optional) Author name and email.

### Step 4: Save Recipe

1.  **Recipe Unique ID**: Enter a unique identifier for this recipe (e.g., `tech-blog-daily`). This will be part of your feed URL.
2.  **Internal Description**: Notes for your reference.
3.  Click **Confirm & Save**.

## Accessing Your Feed

Once saved, the recipe is stored as a **Custom Recipe**. You can manage it in the **Custom Recipes** dashboard.

Your new RSS feed will be available at a URL typically formatted as:
`http://your-feedcraft-instance/rss/custom/{recipe-unique-id}`

You can now subscribe to this URL in your favorite RSS reader.
