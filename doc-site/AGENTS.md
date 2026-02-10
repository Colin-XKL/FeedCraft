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
- **Asides**: Use `:::note`, `:::tip`, `:::caution`, or `:::danger` for callouts.
- **Steps**: Use the `<Steps>` component for multi-step instructions (e.g., deployment or configuration).
- **Tabs**: Use the `<Tabs>` and `<TabItem>` components when showing multiple options (e.g., different deployment methods or model configurations).
- **Cards**: Use `<Card>` or `<CardGrid>` for high-level overviews or navigation.
- **Badges**: Use the `badge` property in frontmatter or the `<Badge>` component to highlight status (e.g., "New", "Beta").

### 3. File Organization & Sidebar
- **Grouping**: Documents are organized into `guides/start` (Quick Start, Concepts) and `guides/advanced` (specific features, customization).
- **Sidebar Order**: Use `sidebar.order` in the frontmatter to control the order of pages.
    - `quick-start.md`: `order: 1`
    - `concepts.md`: `order: 2`
- **Sidebar Labels**: If you add a new category or change a folder structure, update the `sidebar` array in `doc-site/astro.config.mjs` to include the correct labels and translations.

### 3. Frontmatter Requirements
Every Markdown/MDX file must contain valid frontmatter:
```yaml
---
title: Page Title
description: Brief Content
---
```

### 4. Linking Conventions
- **Relative Links**: Prefer relative links (e.g., `../advanced/customization`) when linking between documents within the same language tree.
- **Absolute Links**: Use absolute links starting with the locale (e.g., `/zh-tw/guides/start/quick-start/`) only when necessary, such as in the index page or cross-locale references. Ensure the locale prefix matches the target file's language.

### 5. Deployment Information
- The documentation refers to a demo/public instance at `https://feed-craft.colinx.one`.

## Common Tasks for Agents

- **Adding a Feature**: Document the feature in all three languages. Update `system-craft-atoms.md` if it's a new Atom.
- **Translation**: When translating from `zh` to `zh-tw`, do not just convert characters; adapt the vocabulary to be natural for Traditional Chinese readers.
- **Verification**: After modifying files, verify that the links between pages are not broken.
