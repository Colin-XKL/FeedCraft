---
title: Core Concepts
sidebar:
  order: 2
---

Before diving deep into FeedCraft, it's helpful to understand these three core concepts.

## AtomCraft

**AtomCraft** is the smallest processing unit. In addition to built-in AtomCrafts (like `translate-title`, `fulltext`), you can create custom AtomCrafts based on templates.

**Example: Custom Translation Prompt**
You can create a new AtomCraft named `translate-to-french` based on the `translate-content` template and fill in a custom Prompt in the parameters to instruct the AI to translate content into French.

## FlowCraft

**FlowCraft** is a combined sequence of multiple AtomCrafts. This allows you to chain multiple operations together.

**Example: Fulltext + Summary + Translation**
You can define a FlowCraft named `digest-and-translate` containing the following steps:

1.  `fulltext` (Extract content)
2.  `summary` (Generate summary)
3.  `translate-content` (Translate content)

### Managing FlowCraft

Navigate to **Worktable > FlowCraft** to create and manage FlowCrafts.
The editor allows you to add AtomCrafts and arrange their execution order. Use arrow buttons (⬆️/⬇️) to adjust the order, or the trash icon to remove them from the flow.

## Recipe

**Recipe** binds a specific RSS source URL to an AtomCraft or FlowCraft. This allows you to create a persistent, customized feed URL.

**Managing Recipes:**
Navigate to **Worktable > Custom Recipe** to manage all your created recipes.

- **Create**: Bind a new URL and craft.
- **Preview**: Click the preview button to verify the output directly in the built-in Feed Viewer.
- **Copy Link**: Click the copy icon to get the full subscription URL.

**Example:**

- **Input URL:** `https://news.ycombinator.com/rss`
- **Processor:** `digest-and-translate` (the workflow created above)
- **Result:** You get a new FeedCraft URL. Subscribe to it to get Hacker News with full text, summary, and translation.
