# AI Agent Guide for doc-site Maintenance

This document provides instructions for AI agents on how to maintain and update the FeedCraft documentation site.

## Architecture Overview

- **Framework**: Astro with the [Starlight](https://starlight.astro.build/) integration.
- **Content Location**: `doc-site/src/content/docs/`.
- **Localization**: Supports English (`en`), Simplified Chinese (`zh`), and Traditional Chinese (`zh-tw`).
- **Configuration**: Main configuration is in `doc-site/astro.config.mjs`.

## Maintenance Rules

### 1. Multi-language Synchronization

FeedCraft documentation MUST be kept in sync across all supported languages (`en`, `zh`, `zh-tw`).

- **Atomic Updates**: Every single update or modification MUST be applied to ALL language versions in the same turn/task. Never leave one language behind.
- **Terminology Consistency**:
  - **Simplified Chinese (`zh`)**: Use Mainland China technical terms (e.g., 接口, 数据, 软件).
  - **Traditional Chinese (`zh-tw`)**: Use Taiwan-specific technical terms (e.g., 介面, 資料, 軟體).
  - **Common Terms**: Use "AtomCraft" (原子工藝/原子工艺), "FlowCraft" (組合工藝/组合工艺), and "Recipe" (配方) consistently.

### 2. Using Starlight Components

Since we use the [Starlight](https://starlight.astro.build/components/using-components/) framework, you should leverage its built-in components to create more readable and interactive documentation:

- **File extension**: Pages that `import` or use JSX components (`<Steps>`, `<Tabs>`, `<Card>`, etc.) MUST use the `.mdx` extension. Regular `.md` files will render `import` statements and component tags as raw text.
- **Asides**: Use `:::note`, `:::tip`, `:::caution`, or `:::danger` for callouts. These work in both `.md` and `.mdx`.
- **Steps**: Use the `<Steps>` component for multi-step instructions (e.g., deployment or configuration).
- **Tabs**: Use the `<Tabs>` and `<TabItem>` components when showing multiple options (e.g., different deployment methods or model configurations).
- **Cards**: Use `<Card>` or `<CardGrid>` for high-level overviews or navigation.
- **Badges**: Use the `badge` property in frontmatter or the `<Badge>` component to highlight status (e.g., "New", "Beta").

### 3. File Organization & Sidebar

Documents follow a four-quadrant layout. Put new pages in the group that matches the writing contract:

| Directory | Sidebar label | Writing contract |
| --- | --- | --- |
| `start/` | Getting Started / 上手 | First-time path. The reader should be able to follow it end-to-end. |
| `guides/` | Guides / 指南 | One concrete task per page. Do not mix background essays into a how-to. |
| `concepts/` | Concepts / 概念 | Background, philosophy, and design decisions. Explain *why*. |
| `reference/` | Reference / 参考 | Complete lists and parameters. Map the product; do not narrate a tutorial. |

Current files:

```text
{en,zh,zh-tw}/
├── index.mdx                 # splash + four entry cards
├── start/quick-start.md
├── guides/                   # html-to-rss, json-to-rss, search-to-rss, inbox, customization
├── concepts/                 # concepts, comparison
└── reference/                # system-craft-atoms, tools
```

- **Sidebar Order**: Use `sidebar.order` in the frontmatter to control the order of pages within a group.
- **Related docs**: Set `related` to locale-free slugs (e.g. `guides/html-to-rss`). `MarkdownContent.astro` renders them at the bottom of the page.
- **Sidebar Labels**: If you add a new category or change a folder structure, update the `sidebar` array in `doc-site/astro.config.mjs` to include the correct labels and translations.

### 4. Frontmatter Requirements

Every Markdown/MDX file must contain valid frontmatter:

```yaml
---
title: Page Title
description: Brief Content
---
```

### 5. Linking Conventions

- **Relative Links**: Prefer relative links (e.g., `../guides/customization`) when linking between documents within the same language tree.
- **Absolute Links**: Use absolute links starting with the locale (e.g., `/zh-tw/start/quick-start/`) only when necessary, such as in the index page or cross-locale references. Ensure the locale prefix matches the target file's language.

### 6. Deployment Information

- The documentation refers to a demo/public instance at `https://feed-craft.colinx.one`.
- Old URLs under `/guides/start/` and `/guides/advanced/` are redirected in `doc-site/vercel.json`.

## Common Tasks for Agents

- **Adding a Feature**: Document the feature in all three languages. Pick the correct quadrant (`guides` vs `reference` vs `concepts`). Update `reference/system-craft-atoms.md` if it's a new Atom, using the map groups: content enhancement, AI processing, filtering, input & output.
- **Translation**: When translating from `zh` to `zh-tw`, do not just convert characters; adapt the vocabulary to be natural for Traditional Chinese readers.
- **Verification**: After modifying files, verify that the links between pages are not broken.
