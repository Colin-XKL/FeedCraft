---
title: Search to RSS
description: Generate RSS feeds from search queries using AI providers.
sidebar:
  order: 4
  badge:
    text: beta
    variant: note
---

## Prerequisites

Before using the Search to RSS feature, you need to configure a search provider in the admin settings. See the [Search Provider Configuration Guide](/docs/guides/advanced/customization) for setup instructions.

FeedCraft includes a **Search to RSS** tool that allows you to turn search queries into RSS feeds. This is useful for tracking news, topics, or brand mentions using configured search providers (e.g., SearXNG, Bing, Google).

## How to use

1.  Navigate to **Worktable > Search to RSS** in the admin dashboard.

### Step 1: Search Query

1.  Enter your **Search Query** (e.g., `latest AI news` or `SpaceX launches`).
2.  **Enhanced Mode**: (Optional) Enable this to use AI (LLM) to generate multiple optimized search queries. This helps discover more relevant content by expanding your original query.
3.  Click **Preview Results** to fetch results.

### Step 2: Preview Results

The system will fetch results using the configured search provider.

- Review the list of found items (Title, Date, Link, Description).
- If the results look correct, click **Next Step**.

### Step 3: Feed Metadata

Customize how this feed appears in your RSS reader:

- **Feed Title**: Defaults to "Search: [Query]".
- **Feed Description**: A brief description.
- **Site Link**: Link to the search results page (e.g. Google Search URL).

### Step 4: Save Recipe

1.  **Recipe Unique ID**: Choose a unique identifier for this recipe (e.g., `search-ai-news`).
    - **Auto-fill**: Automatically generated from the feed title using transliteration.
    - **Format**: Strictly enforces lowercase letters, numbers, and hyphens (`[a-z0-9-]`).
    - **Refresh**: Use the refresh button to regenerate the ID if you change the title.
2.  **Internal Description**: Notes for yourself about this recipe.
3.  Click **Confirm & Save**.

## Accessing Your Feed

Once saved, the recipe is stored as a **Custom Recipe**. You can manage it in the **Custom Recipes** dashboard.

Your new feed will be available at:
\`http://your-feedcraft-instance/rss/custom/{recipe-unique-id}\`
