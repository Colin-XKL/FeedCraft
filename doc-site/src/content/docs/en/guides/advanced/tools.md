---
title: System Tools
description: Built-in tools for debugging feeds, comparing outputs, and checking system health.
sidebar:
  order: 5
---

FeedCraft provides several built-in tools to help you debug your RSS feeds and monitor the system's health. You can access these tools under the **Tools** menu in the admin dashboard.

## RSS Viewer

The **RSS Viewer** (Feed Viewer) allows you to preview any RSS feed as FeedCraft sees it.

- **Usage**:
  1. Navigate to **Tools > RSS Viewer**.
  2. Enter an RSS/Atom URL.
  3. Click **Preview**.
- **Purpose**: Verify if FeedCraft can successfully fetch and parse a feed before setting up a recipe.
- **Note**: The viewer uses the `proxy` craft by default, which simply fetches the feed without modification.

## Example RSS Feeds

The **Example RSS Feeds** page provides built-in subscriptions for testing how RSS readers render HTML, CSS, and media content, plus whether they support RSS 1.0, Atom, and JSON Feed documents.

- **Usage**:
  1. Navigate to **Tools > Example RSS Feeds**.
  2. Copy one of the subscription URLs, such as `/example-rss-feeds/html-elements.xml`.
  3. Subscribe to it in your RSS reader.
- **Available feeds**:
  - `html-elements.xml`: headings, lists, tables, blockquotes, code blocks, details/summary, figures, and other common HTML5 elements.
  - `html-styling.xml`: inline color, background, border, spacing, typography, flex, and grid styles.
  - `media-picture.xml`: `picture`, `source`, `srcset`, `sizes`, fallback images, alt text, and captions.
  - `all-in-one.xml`: combines the HTML, styling, and media fixtures into one feed.
  - `rss-1-0.rdf`: a simple RSS 1.0/RDF document.
  - `atom.xml`: a simple Atom document.
  - `json-feed.json`: a simple JSON Feed 1.1 document.
- **Refresh behavior**: The subscription URL stays stable, while item GUIDs rotate every 4 hours so readers can fetch a fresh copy.

## Feed Compare

The **Feed Compare** tool lets you visualize the effect of a Craft (Atom or Flow) on a feed.

- **Usage**:
  1. Navigate to **Tools > Feed Compare**.
  2. Enter the original RSS feed URL.
  3. Select a **FlowCraft** or **AtomCraft** to apply.
  4. Click **Compare**.
- **Output**: The tool displays two columns:
  - **Left**: The original feed content.
  - **Right**: The feed content after being processed by the selected Craft.
- **Use Case**: Great for testing new translation flows or summarization prompts without creating a permanent recipe.

## Craft Dependencies

The **Craft Dependencies** (System Health) tool visualizes the internal relationships between your Recipes, FlowCrafts, and AtomCrafts.

- **Usage**:
  1. Navigate to **Tools > Craft Dependencies**.
  2. Click **Analyze Craft Dependencies**.
- **Features**:
  - Generates a tree view of all dependencies.
  - **Health Check**: Automatically detects missing dependencies (e.g., a Recipe pointing to a deleted FlowCraft).
  - **Missing Crafts Panel**: Highlights explicitly which Crafts are missing at the top of the view.
  - **Visual Indicators**: Different colors for Recipes, Flows, Atoms, and missing components.

:::tip
If you encounter errors like "Craft not found", use this tool to trace the broken link in your configuration.
:::

## System Runtime

The **System Runtime** (Observability) tool provides a comprehensive dashboard for monitoring the health and execution status of your resources.

- **Usage**:
  1. Navigate to **Tools > System Runtime**.
- **Features**:
  - **Resource Health**: View the current status (Healthy, Degraded, Paused) of Recipes and other components, including consecutive failures.
  - **Execution Logs**: Track detailed execution history, success rates, and specific error types (e.g., Timeout, Network, Parse) across all runs.
  - **System Notifications**: Review automated alerts regarding resource state transitions (e.g., when a Recipe becomes degraded). You can also subscribe to these alerts via the built-in RSS feed at `/system/notifications/rss`.

:::tip
If a Recipe fails repeatedly and becomes "Paused", you can use the System Runtime dashboard to manually "Resume" it after fixing the underlying issue.
:::

## Debug Tools

### LLM Debug

A sandbox for testing your LLM configuration. You can send test prompts to your configured LLM provider to verify connectivity and model response.

### Ad Check Debug

A specific tool for testing the "Ignore Advertorial" filter logic against specific content to understand why an article might be filtered out.
