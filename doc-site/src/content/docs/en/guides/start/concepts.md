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

## Topic Feed

**Topic Feed** is an aggregation unit that combines multiple input sources (like `RawFeed`s or other `Recipe`s) into a single, unified RSS feed. It solves information overload by bringing disparate sources into one place.

You can configure aggregation rules for a Topic Feed to automatically process the merged articles:

- **Deduplicate**: Removes duplicate articles across sources. Five strategies are available:
  - `by_link` — exact URL match (default)
  - `by_id` — exact GUID match
  - `by_title` — exact title match
  - `by_simhash` _(Lightweight)_ — character-based fuzzy deduplication using SimHash. Compares the literal text of title and content; no AI service required. Best for filtering reposts and lightly reworded copies. Threshold is a **difference score** (`0.0` = identical only, `1.0` = everything, default `0.05`; lower is stricter).
  - `by_embedding` _(Precise)_ — semantic deduplication via LLM Embedding vectors. Understands meaning, not just wording — catches paraphrased duplicates and different outlets covering the same story. Requires a configured Embedding API; fails open (keeps all articles) when unavailable. Threshold is a **cosine similarity score** (`0.0`–`1.0`, default `0.9`; higher is stricter).
- **Sort**: Orders the combined articles by date or quality score.
- **Limit**: Keeps only the specified number of most recent items.

:::tip
Rules run in the order you define them. A typical pipeline is: **Sort** (newest first) → **Deduplicate** → **Limit**.
:::

**Managing Topic Feeds:**

Navigate to **Worktable > Topic Feed** to create and manage topics.

- **Create**: Define a title, add multiple input sources using the type picker (supports external RSS URLs, custom recipes, or nesting other topics), and configure your aggregation rules.
- **Public Access**: Your new topic feed will be available without authentication at `http://your-feedcraft-instance/topic/{id}`.
