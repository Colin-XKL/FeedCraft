# Quick Start

## Introduction

FeedCraft is a powerful tool to process your RSS feeds as a middleware. You can use it to translate your feed, extract full text, emulate a browser to render JS-heavy pages, use LLM such as Google Gemini to generate briefs for your RSS articles, use natural language to filter your RSS feed, and more!

## Portable Mode

You can quickly start using FeedCraft by modifying the URL of your RSS feed. This is called "Portable Mode".

The URL format is:
`https://feed-craft.colinx.one/craft/{craft_atom}?input_url={input_rss_url}`

Where:
- `{craft_atom}` is the name of the processing step you want to apply.
- `{input_rss_url}` is the original URL of the RSS feed you want to process.

**Note:** You may need to URL-encode the `{input_rss_url}` if your RSS reader does not handle it automatically.

### Common Craft Atoms

Here are some basic Craft Atoms you can use immediately:

- **proxy**: Simple RSS proxy, no processing.
- **limit**: Limits the number of articles (default latest 10).
- **fulltext**: Extracts full text from the article link.
- **fulltext-plus**: Emulates a browser to render and extract full text (useful for JS-heavy sites).
- **introduction**: Uses AI to generate a brief introduction at the beginning of the article.
- **summary**: Uses AI to summarize the main content of the article.
- **translate-title**: Uses AI to translate the article title.
- **translate-content**: Uses AI to translate the article content (replaces original).
- **translate-content-immersive**: Uses AI to translate content in immersive mode (paragraph by paragraph).
- **ignore-advertorial**: Uses AI to filter out advertorials.

### Examples

1.  **Translate titles of a feed:**
    `https://feed-craft.colinx.one/craft/translate-title?input_url=https://feeds.feedburner.com/visualcapitalist`

2.  **Get full text of a feed:**
    `https://feed-craft.colinx.one/craft/fulltext?input_url=https://feeds.feedburner.com/visualcapitalist`

## Basic Deployment

You can deploy your own instance using Docker Compose.

### minimal `docker-compose.yml`

```yaml
version: "3"
services:
  app.feed-craft:
    image: ghcr.io/colin-xkl/feed-craft
    container_name: feed-craft
    restart: always
    ports:
      - "10088:80"
    volumes:
      - ./feed-craft-db:/usr/local/feed-craft/db
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000
      FC_REDIS_URI: redis://service.redis:6379/
      FC_OPENAI_AUTH_KEY: skxxxxxx
      FC_OPENAI_DEFAULT_MODEL: gemini-pro
      FC_OPENAI_ENDPOINT: https://xxxxxx
```

Save this as `docker-compose.yml` and run `docker-compose up -d`.
Visit `http://localhost:10088` to access the dashboard.
Default login: `admin` / `adminadmin`.
