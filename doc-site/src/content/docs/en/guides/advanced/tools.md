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
  - **Visual Indicators**: Different colors for Recipes, Flows, Atoms, and missing components.

> **Tip:** If you encounter errors like "Craft not found", use this tool to trace the broken link in your configuration.

## Debug Tools

### LLM Debug

A sandbox for testing your LLM configuration. You can send test prompts to your configured LLM provider to verify connectivity and model response.

### Ad Check Debug

A specific tool for testing the "Ignore Advertorial" filter logic against specific content to understand why an article might be filtered out.
