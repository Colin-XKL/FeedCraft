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

## Tools & Utilities

FeedCraft includes several built-in tools to help you debug and manage your feeds.

### RSS Viewer

Located at **Tools > RSS Viewer**, this tool allows you to preview the content of any RSS feed. You can paste an RSS URL to see how FeedCraft parses the items, which is useful for verifying feed validity before processing.

### Feed Compare

Located at **Tools > Feed Compare**, this tool lets you compare an original RSS feed against a processed version. By selecting a transformation workflow (AtomCraft or FlowCraft), you can visualize exactly how the content is modified (e.g., filtered articles, added summaries).

### Craft Dependencies

Located at **Tools > Craft Dependencies**, the **System Health** dashboard visualizes the internal dependency graph of your FeedCraft instance. It shows the relationships between Recipes, FlowCrafts, and AtomCrafts, and highlights any missing components or broken references.

> **Note:** This is different from the **Dependency Services** dashboard, which monitors external infrastructure like Redis or Browserless.
