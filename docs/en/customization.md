# High Level Customization

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

### Output Format

You can specify the output format by adding the `output_type` parameter to the URL.

**Supported Formats:**
-   `rss`: Standard RSS 2.0 format (default)
-   `atom`: Atom format
-   `json`: JSON Feed format

**Example URLs:**
-   RSS (default): `http://localhost:10088/craft/my-flow?input_url=...`
-   Atom: `http://localhost:10088/craft/my-flow?input_url=...&output_type=atom`
-   JSON Feed: `http://localhost:10088/craft/my-flow?input_url=...&output_type=json`

For **Recipes**, you can also append `&output_type=json` to the generated URL to get the output in JSON Feed format.

## Advanced Configuration

### Docker Environment Variables

You can configure FeedCraft using environment variables in your `docker-compose.yml`.

-   **FC_PUPPETEER_HTTP_ENDPOINT**: URL of a browserless/chrome instance. Required for `fulltext-plus`.
-   **FC_REDIS_URI**: Redis connection string. Used for caching to speed up processing and reduce AI costs.
-   **FC_OPENAI_AUTH_KEY**: API Key for OpenAI or compatible services (DeepSeek, Gemini, etc.).
-   **FC_OPENAI_DEFAULT_MODEL**: The default model to use (e.g., `gemini-pro`, `gpt-3.5-turbo`).
-   **FC_OPENAI_ENDPOINT**: The API endpoint. Must end with `/v1` if using OpenAI compatible APIs.

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
