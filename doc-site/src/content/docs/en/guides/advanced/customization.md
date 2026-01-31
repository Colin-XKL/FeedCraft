---
title: Advanced Customization
sidebar:
  order: 1
---

For advanced users, FeedCraft provides an admin dashboard to customize the RSS processing workflow.

## Accessing the Dashboard

1.  Deploy FeedCraft using Docker (see Quick Start).
2.  Open your browser and visit `http://YOUR_SERVER_IP:10088`.
3.  Log in with default credentials:
    - Username: `admin`
    - Password: `adminadmin`
      _(Please change the password immediately after logging in)_

## Core Concepts

### AtomCraft

**AtomCraft** is the smallest processing unit. In addition to built-in AtomCrafts (like `translate-title`, `fulltext`), you can create custom AtomCrafts based on templates.

**Example: Custom Translation Prompt**
You can create a new AtomCraft named `translate-to-french` based on the `translate-content` template and fill in a custom Prompt in the parameters to instruct the AI to translate content into French.

### FlowCraft

**FlowCraft** is a combined sequence of multiple AtomCrafts. This allows you to chain multiple operations together.

**Example: Fulltext + Summary + Translation**
You can define a FlowCraft named `digest-and-translate` containing the following steps:

1.  `fulltext` (Extract content)
2.  `summary` (Generate summary)
3.  `translate-content` (Translate content)

#### Managing FlowCraft

Navigate to **Worktable > FlowCraft** to create and manage FlowCrafts.
The editor allows you to add AtomCrafts and arrange their execution order. Use arrow buttons (⬆️/⬇️) to adjust the order, or the trash icon to remove them from the flow.

### Recipe

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

## Search Provider Configuration

To use the **Search to RSS** feature, you must configure a search provider.

Navigate to **Settings > Search Provider** in the admin dashboard.

### Supported Providers

- **LiteLLM / OpenAI Compatible**
  - **API URL**: The search endpoint of your provider (e.g., `http://litellm-proxy:4000/v1/search`).
  - **API Key**: Your API key. (Leave empty to keep existing key)
  - **Tool Name**: The specific function calling tool name if required (e.g., `google_search` for some agents). The tool name is appended to the API URL (e.g. `.../v1/search/google_search`).

- **SearXNG**
  - **API URL**: The base URL of your SearXNG instance (e.g., `http://my-searxng.com`). The `/search` path is automatically appended.
  - **Engines**: (Optional) Comma-separated list of engines to use (e.g., `google,bing`).

> **Tip:** You can use the **Check Connection** button to verify connectivity with your provider before saving.

## Dependency Services

The **Dependency Services** dashboard (Settings > Dependency Services) provides a health check overview of all connected external services.

It monitors the status of:

- **SQLite**: Database connectivity.
- **Redis**: Cache service connectivity and latency.
- **Browserless**: Headless browser service availability (required for fulltext extraction).
- **LLM Service**: Connectivity to the configured AI provider.
- **Search Provider**: Connectivity to the configured search engine.

Use this dashboard to troubleshoot connectivity issues if features like "Enhanced Mode" or "Fulltext Extraction" are failing.

You can use the **Check Connection** button to verify if FeedCraft can successfully connect to the search provider with the provided credentials.

## Advanced Configuration

### Docker Environment Variables

You can configure FeedCraft using environment variables in `docker-compose.yml`.

- **FC_PUPPETEER_HTTP_ENDPOINT**: Address of the Browserless/Chrome instance. Required for `fulltext-plus`.
- **FC_REDIS_URI**: Redis connection address. Used for caching to speed up processing and reduce AI token consumption.
- **FC_LLM_API_KEY**: API Key for OpenAI or compatible services (like DeepSeek, Gemini, etc.).
- **FC_LLM_API_MODEL**: Default model to use (e.g., `gemini-pro`, `gpt-3.5-turbo`). **Multiple Models Support:** You can provide a comma-separated list of models (e.g., `gpt-3.5-turbo,gpt-4`). FeedCraft will randomly select a model for each request and automatically retry with others if a call fails.
- **FC_LLM_API_BASE**: API endpoint address. For OpenAI-compatible APIs, usually ends with `/v1`.
- **FC_LLM_API_TYPE**: (Optional) `openai` (default) or `ollama`.

### External Services

To leverage the full power of FeedCraft, it is recommended to deploy with Redis and Browserless.

```yaml
version: "3"
services:
  app.feed-craft:
    # ... (Refer to Quick Start)
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

The service listens on port 80 by default. You can also access it from other containers in the same network using `http://app.feed-craft/xxx` (e.g., for internal communication with an RSS reader).
