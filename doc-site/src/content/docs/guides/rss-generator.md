---
title: RSS Generator
description: Turn any website into an RSS feed using the visual RSS Generator tool.
---

The **RSS Generator** is a powerful tool built into FeedCraft that allows you to turn almost any website into an RSS feed. It provides a visual interface to select elements on a webpage and map them to RSS fields like Title, Link, Date, and Content.

## How to Access

You can access the RSS Generator from the FeedCraft dashboard under **Tools** > **RSS Generator**.

## Workflow

### 1. Fetch the Page

1.  Enter the URL of the website you want to convert into an RSS feed.
2.  Click **Fetch Page**.
3.  The page content will be loaded in the "Page Preview" area on the left.

:::note
The tool uses a backend proxy to fetch the page content, ensuring it handles some basic bot protection. However, extremely complex dynamic sites might still require more advanced configuration.
:::

### 2. Selection Mode

Once the page is loaded, you can use the **Selection Mode** to visually pick elements.

-   **Toggle Mode**: Click "Selection Mode" / "Preview Mode" to switch between interacting with the page (Preview) and picking elements (Selection).
-   **Visual Picking**: In Selection Mode, hovering over elements will highlight them. Clicking an element will select it for the active field.

### 3. Configure Selectors

The configuration process involves two main steps:

#### Step 1: Define List Item
First, you need to identify the repeating element that represents a single article or item in the feed.

1.  Click the **Pick** button inside the "List Item Selector" input field.
2.  In the preview area, click on one of the article cards or list items.
3.  The tool will automatically generate a CSS selector (e.g., `.article-card`) and show how many items match that selector.

#### Step 2: Map Fields
Once the List Item is defined, you can map the specific fields relative to each item.

-   **Title Selector**: Click **Pick** and select the title element inside one of the list items.
-   **Link Selector**: Click **Pick** and select the link element (usually the title or a "Read More" button).
-   **Date Selector**: (Optional) Select the date element.
-   **Content Selector**: (Optional) Select the summary or content element.

:::tip[Keyboard Shortcuts]
When in Selection Mode, you can use the **Arrow Up** and **Arrow Down** keys to navigate the DOM tree and select parent or child elements precisely.
:::

### 4. Preview RSS Items

After configuring the selectors:

1.  Click **Preview RSS Items**.
2.  The "Extracted Items" panel on the right will show the data extracted from the page.
3.  Verify that the Titles, Links, and other fields are correct.

## Advanced Usage

-   **Manual CSS Selectors**: You can manually type or edit the CSS selectors if the auto-generated ones are not precise enough.
-   **Relative Selection**: The Title, Link, Date, and Content selectors are *relative* to the List Item selector. This ensures that the data is correctly grouped for each item.

## Troubleshooting

-   **Fetch Failed**: If the page fails to load, check if the URL is correct. Some sites may block the fetcher.
-   **No Items Found**: detailed check your "List Item Selector". Ensure it matches the HTML structure of the page.
-   **Wrong Data**: If the title or link is incorrect, try re-selecting the specific element or using the arrow keys to select the parent container.
