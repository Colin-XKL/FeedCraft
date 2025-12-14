---
title: High Level Customization
---

For power users, FeedCraft offers a dashboard to customize how feeds are processed.

## Accessing the Dashboard

1.  Deploy FeedCraft using Docker (see Quick Start).
2.  Navigate to `http://your-server-ip:10088`.
3.  Log in with the default credentials:
    -   Username: `admin`
    -   Password: `adminadmin`
    *(Please change the password immediately after logging in)*

## Core Concepts

### Craft Atom

A **Craft Atom** is a single processing unit. While there are built-in atoms (like `translate-title`, `fulltext`), you can create your own based on templates.

**Example: Custom Translation Prompt**
You can create a new Atom named `translate-to-french` based on the `translate-content` template, but with a custom prompt instructing the AI to translate to French.

### Craft Flow

A **Craft Flow** is a sequence of Craft Atoms. This allows you to chain multiple operations together.

**Example: Fulltext + Summary + Translation**
You can define a flow named `digest-and-translate` that consists of:
1.  `fulltext` (Extract content)
2.  `summary` (Generate summary)
3.  `translate-content` (Translate the result)

### Recipe

A **Recipe** binds a specific RSS feed URL to a Craft or Craft Flow. This allows you to create a persistent, customized feed URL.

**Example:**
-   **Input URL:** `https://news.ycombinator.com/rss`
-   **Processor:** `digest-and-translate` (The flow created above)
-   **Result:** You get a new FeedCraft URL that serves Hacker News with full text, summaries, and translation.

## Advanced Configuration

### Docker Environment Variables

You can configure FeedCraft using environment variables in your `docker-compose.yml`.

-   **FC_PUPPETEER_HTTP_ENDPOINT**: URL of a browserless/chrome instance. Required for `fulltext-plus`.
-   **FC_REDIS_URI**: Redis connection string. Used for caching to speed up processing and reduce AI costs.
-   **FC_LLM_API_KEY**: API Key for OpenAI or compatible services (DeepSeek, Gemini, etc.).
-   **FC_LLM_API_MODEL**: The default model to use (e.g., `gemini-pro`, `gpt-3.5-turbo`). **Supports multiple models:** You can provide a comma-separated list of models (e.g., `gpt-3.5-turbo,gpt-4`). FeedCraft will randomly select one for each request and automatically retry with others if a call fails.
-   **FC_LLM_API_BASE**: The API endpoint. Must end with `/v1` if using OpenAI compatible APIs.
-   **FC_LLM_API_TYPE**: (Optional) `openai` (default) or `ollama`.

### External Services

To fully utilize FeedCraft, you should deploy it alongside Redis and Browserless.

```yaml
version: "3"
services:
  app.feed-craft:
    # ... (see Quick Start)
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000
      FC_REDIS_URI: redis://service.redis:6379/
      # ...

  service.redis:
    image: redis:6-alpine
    container_name: feedcraft_redis
    restart: always

  service.browserless:
    image: browserless/chrome
    container_name: feedcraft_browserless
    environment:
      USE_CHROME_STABLE: true
    restart: unless-stopped
```
