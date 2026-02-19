---
title: Troubleshooting & Tools
description: Learn how to use the built-in troubleshooting and utility tools in FeedCraft.
sidebar:
  order: 50
---

FeedCraft provides several built-in tools within the admin dashboard to help you debug, verify, and monitor your feeds.

## Feed Viewer

The **Feed Viewer** allows you to preview any RSS feed rendered in a clean, readable format. This is useful for verifying that a feed is valid and checking its content without needing an external RSS reader.

- **Location**: **Tools > Feed Viewer**
- **How to use**:
  1.  Enter the URL of the RSS feed you want to preview.
  2.  Click **Preview**.
  3.  The feed content will be displayed below, showing the title, description, and individual items.

## Feed Compare

The **Feed Compare** tool lets you compare two RSS feeds side-by-side. This is particularly helpful when debugging a FlowCraft to see how the processing steps have altered the original feed.

- **Location**: **Tools > Feed Compare**
- **How to use**:
  1.  Enter the **Source URL** of the original feed.
  2.  Select a **FlowCraft** from the dropdown menu (e.g., `translate-title`).
  3.  Click **Compare**.
  4.  The tool will display the "Original Feed" on the left and the "Craft Applied Feed" on the right, allowing you to easily spot differences.

## System Health

The **System Health** dashboard visualizes the internal dependency graph of your FeedCraft instance. It shows the relationships between Recipes, FlowCrafts, and AtomCrafts.

- **Location**: **Tools > System Health**
- **Purpose**:
  - Verify that all components of your custom recipes exist.
  - Detect missing dependencies (e.g., a recipe referencing a deleted FlowCraft).
  - Identify circular dependencies that could cause issues.
- **How to use**:
  1.  Click **Analyze**.
  2.  Review the tree graph. Healthy components are marked with their type (Recipe, FlowCraft, AtomCraft). Missing components are marked in red.

## URL Generator

The **URL Generator** helps you create FeedCraft subscription URLs easily. It also features a "Parsing Mode" to reverse-engineer existing URLs.

- **Location**: **Dashboard > Quick Start**
- **Features**:
  - **Generate**: Select a FlowCraft and input a source URL to get a subscription link.
  - **Parse**: Paste a FeedCraft URL to see its components (FlowCraft used, original source, etc.).

## Dependency Services

For monitoring external services (like Redis, Browserless, LLM), use the **Dependency Services** dashboard.

- **Location**: **Settings > Dependency Services**
- **See Also**: [Advanced Customization](/en/guides/advanced/customization/#dependency-services)
