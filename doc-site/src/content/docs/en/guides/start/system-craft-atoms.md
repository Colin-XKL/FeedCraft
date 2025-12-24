---
title: System Craft Atoms
description: Reference guide for built-in system craft atoms in FeedCraft.
---

FeedCraft comes with a set of built-in "Craft Atoms" that perform specific processing steps on your feeds. You can chain these atoms together in a Craft Flow to create powerful pipelines.

## Content Acquisition & Repair

These atoms help you fetch full content or fix common feed issues.

### `fulltext`
Extracts the full content of the article from the original webpage.
- **Use case:** When the RSS feed only provides a summary or snippet.
- **Mechanism:** Uses a standard HTTP client to fetch the page and an algorithm to extract the main content. Fast and lightweight.

### `fulltext-plus`
Extracts full content using a headless browser (Puppeteer).
- **Use case:** For websites that require JavaScript to render content or have strong anti-bot protections.
- **Mechanism:** Connects to the configured Browserless/Puppeteer service to render the page. Slower but more robust.

### `proxy`
Simple proxy for the feed.
- **Use case:** When you just want to forward the original feed without modification, or use FeedCraft as a central gateway.

### `guid-fix`
Replaces the RSS item GUID with an MD5 hash of the item's content.
- **Use case:** Some feeds change their GUIDs frequently even when content hasn't changed, causing duplicate unread items in readers. This atom stabilizes the GUID based on content.

### `relative-link-fix`
Converts relative links (e.g., `<a href="/about">`) in the content to absolute links (e.g., `<a href="https://example.com/about">`).
- **Use case:** Essential when extracting full content, as relative links will break when viewed in an RSS reader.

### `cleanup`
Cleans up the HTML content to remove clutter.
- **Use case:** Improving readability by removing classes, styles, and empty tags.

---

## Filtering

Control which items make it into your final feed.

### `limit`
Limits the number of items in the feed.
- **Parameters:**
  - `num` (default: `10`): The maximum number of items to keep.

### `time-limit`
Filters out items that are older than a specific number of days.
- **Parameters:**
  - `days` (default: `30`): Max age of articles in days.

### `keyword`
Filters items based on keywords in the title or content.
- **Parameters:**
  - `keyword`: The Regex pattern to match.
  - `mode`: `keep` (default) to keep matching items, or `block` to remove them.
  - `target`: `title`, `content`, or `all` (default).

---

## AI Enhancement

Use Large Language Models (LLM) to transform and enrich your content.
:::note
These atoms require LLM configuration (API Key, Base URL, etc.) in your environment variables.
:::

### `translate-title`
Translates the article title to your target language.
- **Parameters:**
  - `prompt`: Custom prompt. Defaults to a standard translation prompt. Supports `{{.TargetLang}}` placeholder.

### `translate-content`
Translates the entire article content, replacing the original.
- **Parameters:**
  - `prompt`: Custom prompt. Supports `{{.TargetLang}}`.

### `translate-content-immersive`
Bilingual translation. Appends the translation after each paragraph of the original text.
- **Parameters:**
  - `prompt`: Custom prompt.

### `summary`
Generates a summary of the article and prepends it to the content.
- **Parameters:**
  - `prompt`: Custom prompt for summarization.

### `introduction`
Generates a short introduction or "lead-in" for the article.
- **Parameters:**
  - `prompt`: Custom prompt.

### `beautify-content`
Re-formats the article using LLM to fix layout, remove ads, and standardizing Markdown, then converts back to clean HTML.
- **Parameters:**
  - `prompt`: Instructions for the "editor" persona.

---

## AI Filtering

Advanced filtering using semantic understanding.

### `ignore-advertorial`
Uses LLM to detect if an article is an advertorial or soft advertisement and removes it.
- **Parameters:**
  - `prompt-for-exclude`: A prompt that should return `true` if the item is an ad.

### `llm-filter`
Generic LLM-based filter. You define the condition for **exclusion**.
- **Parameters:**
  - `filter_condition`: A natural language question/condition. If the LLM answers "yes" (true), the item is **removed**.
  - *Example:* "Is this article about sports?" (Removes sports articles).
