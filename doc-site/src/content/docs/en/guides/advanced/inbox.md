---
title: Inbox
description: Receive articles from third-party sources via HTTP push and expose them as an RSS feed.
sidebar:
  order: 6
  badge:
    text: new
    variant: success
---

import { Steps } from '@astrojs/starlight/components';

The **Inbox** feature lets external services, scripts, or automations push articles directly into FeedCraft over HTTP. Each inbox is then exposed as an RSS feed through a Custom Recipe, making it easy to subscribe in any RSS reader.

## Overview

A typical workflow looks like this:

1. **Create an Inbox** in the admin dashboard (gives you a unique inbox ID).
2. **Create a System Auth Token** (the secret that authorises push requests).
3. **Push items** from your script, automation, or third-party tool using a standard JSON HTTP POST.
4. **Create a Custom Recipe** that uses the inbox as its data source.
5. **Subscribe** to the generated RSS URL in your reader.

## Managing Inboxes

Navigate to **Worktable > Inbox Management** in the admin dashboard.

### Creating an Inbox

<Steps>
1. Click **Create Inbox**.
2. Fill in the required fields:
   - **Inbox ID**: A unique, URL-safe identifier (lowercase letters, numbers, hyphens, underscores). Cannot be changed after creation.
   - **Title**: A human-readable name for this inbox.
   - **Max Items**: Maximum number of articles to retain. Older items are automatically pruned when the limit is exceeded (default: 100).
   - **Public Access**: If enabled, article content can be fetched without authentication. If disabled, a valid System Auth Token must be provided.
3. Click **OK** to save.
</Steps>

### Editing an Inbox

Click **Edit Inbox** in the actions column to update the Title, Description, Max Items, or Public Access setting. The Inbox ID cannot be modified.

:::caution
Reducing **Max Items** on an existing inbox immediately prunes overflow articles. Pruned articles cannot be recovered.
:::

### Deleting an Inbox

Click **Delete** in the actions column. This permanently removes the inbox and **all articles** stored inside it.

## Managing System Auth Tokens

Navigate to **Settings > System Auth Token** to create API tokens that authorise push requests.

<Steps>
1. Click **Generate Token**.
2. Enter a descriptive label (e.g., "iPhone Shortcut", "Home Assistant").
3. Copy the generated token immediately — it is only shown once and cannot be retrieved later.
</Steps>

Tokens can be revoked at any time by clicking **Delete**. All integrations using the revoked token will stop working immediately.

## Pushing Articles

Use the push endpoint to send articles from any HTTP client, script, or automation platform.

### Endpoint

```
POST /api/inbox/{inbox_id}/items
```

### Authentication

Include the System Auth Token in the `Authorization` header:

```
Authorization: Bearer YOUR_SYSTEM_AUTH_TOKEN
Content-Type: application/json
```

### Request Body

Send a JSON array of article objects. Only `title` is required; all other fields are optional.

```json
[
  {
    "id": "optional-custom-unique-id",
    "title": "Article Title",
    "url": "https://example.com/article",
    "content": "<p>Full HTML body of the article.</p>",
    "summary": "A short description shown in feed previews.",
    "author": "Author Name",
    "timestamp": 1716470400
  }
]
```

| Field | Required | Description |
|-------|----------|-------------|
| `title` | ✅ | Article headline. |
| `id` | Optional | Custom stable ID. If omitted, a UUID is auto-generated. If the same `id` is pushed again, the article is **updated** (upsert). |
| `url` | Optional | Canonical link. If omitted, FeedCraft generates a link pointing to the article's stored content. |
| `content` | Optional | Full HTML body. |
| `summary` | Optional | Short description. Defaults to the first 200 characters of `content`. |
| `author` | Optional | Author name. |
| `timestamp` | Optional | Unix timestamp (seconds) for the publication date. Defaults to the current time. |

**Batch limit**: Maximum 100 items per request.

### cURL Example

```bash
curl -X POST "https://YOUR_SERVER/api/inbox/my-inbox/items" \
  -H "Authorization: Bearer YOUR_SYSTEM_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{"title": "Hello World", "content": "<p>First article pushed!</p>"}]'
```

### Response

```json
{
  "total": 1,
  "created": 1,
  "updated": 0
}
```

## Subscribing via RSS

To subscribe to inbox articles in an RSS reader, create a Custom Recipe that uses the inbox as its data source.

<Steps>
1. Navigate to **Worktable > Custom Recipe** and click **Create Recipe**.
2. Set **Source Type** to `inbox`.
3. In the **Source Config JSON** field, enter:
   ```json
   { "inbox_source": { "inbox_id": "YOUR_INBOX_ID" } }
   ```
4. Set **Craft** to `proxy` (or any other craft chain you want to apply).
5. Save the recipe and click **Copy Link** in the recipe list to get your RSS subscription URL.
</Steps>

:::tip
You can apply AtomCrafts or FlowCrafts on top of the inbox source — for example, use `translate-content` to automatically translate pushed articles into another language.
:::

## Access Control for Private Inboxes

When an inbox has **Public Access** disabled, the article content endpoint requires authentication.

Append `?token=YOUR_SYSTEM_AUTH_TOKEN` to the article URL:

```
GET /inbox/{inbox_id}/items/{article_id}/content?token=YOUR_TOKEN
```

Or use the `Authorization: Bearer YOUR_TOKEN` header.

:::note
The RSS feed itself (served via Custom Recipe) is controlled by the Recipe's own access settings, not the inbox's public flag. The private flag only affects the raw article content endpoint at `/inbox/{inbox_id}/items/{article_id}/content`.
:::

## Maintenance (GC)

FeedCraft provides garbage collection utilities accessible from the admin API:

- **GET `/api/admin/inboxes/gc/stats`** — Returns the count of total items, orphaned items (belonging to deleted inboxes), and overflow items.
- **POST `/api/admin/inboxes/gc/cleanup`** — Deletes all orphaned and overflow items in a single atomic transaction.
