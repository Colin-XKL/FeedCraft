---
title: System AtomCrafts
description: Reference guide for built-in system AtomCrafts in FeedCraft.
---

FeedCraft comes with a set of built-in "AtomCrafts" that perform specific processing steps on your feeds. You can chain these AtomCrafts together in a FlowCraft to create powerful pipelines.

## Content Acquisition & Repair

These atoms help you fetch full content or fix common feed issues.

### `fulltext`

Extracts the full content of the article from the original webpage.

- **Use case:** When the RSS feed only provides a summary or snippet.
- **Mechanism:** Uses a standard HTTP client to fetch the page and an algorithm to extract the main content. Fast and lightweight.

### `fulltext-plus`

Extracts full content using a headless browser (Puppeteer).

- **Use case:** For websites that require JavaScript to render content or have strong anti-bot protections.
- **Mechanism:** Connects to the configured browser provider (`browserless-restful` or `cdp`) to render the page. Slower but more robust.
- **Parameters:**
  - `mode` (default: `networkidle2`): Wait condition.
    - `load`: Wait for the `load` event.
    - `domcontentloaded`: Wait for the `DOMContentLoaded` event.
    - `networkidle0`: Wait until there are **0** active network connections for at least 500ms.
    - `networkidle2`: Wait until there are no more than **2** active network connections for at least 500ms. (Recommended for SPAs).
  - `wait` (default: `0`): Explicit wait time in seconds (e.g., `5`).

### `proxy`

Simple proxy for the feed.

- **Use case:** When you just want to forward the original feed without modification, or use FeedCraft as a central gateway.

### `guid-fix`

Replaces the RSS item GUID with an FNV-1a hash of the item's content.

- **Use case:** Some feeds change their GUIDs frequently even when content hasn't changed, causing duplicate unread items in readers. This atom stabilizes the GUID based on content.

### `relative-link-fix`

Converts relative links (e.g., `<a href="/about">`) in the content to absolute links (e.g., `<a href="https://example.com/about">`).

- **Use case:** Essential when extracting full content, as relative links will break when viewed in an RSS reader.

### `cleanup`

Cleans up the HTML content to remove clutter.

- **Use case:** Improving readability by removing classes, styles, and empty tags.

---

## Filtering

Control which items make it into your final feed.

### `limit`

Limits the number of items in the feed.

- **Parameters:**
  - `num` (default: `10`): The maximum number of items to keep.

### `time-limit`

Filters out items that are older than a specific number of days.

- **Parameters:**
  - `days` (default: `7`): Max age of articles in days.

### `keyword`

Filters items based on keywords in the title or content.

- **Parameters:**
  - `keywords`: A comma-separated list of keywords to match (substring match, case-sensitive). Example: `ad,sell,SALE`.
  - `mode`: `include` (default) to keep matching items, or `exclude` to remove them.
  - `scope`: `title`, `content`, or `all` (default).

---

## AI Enhancement

Use Large Language Models (LLM) to transform and enrich your content.
:::note
These atoms require LLM configuration (API Key, Base URL, etc.) in your environment variables.
:::

### `translate-title`

Translates the article title to your target language.

- **Parameters:**
  - `prompt`: Custom prompt. Defaults to a standard translation prompt. Supports `{{.TargetLang}}` placeholder.

### `translate-content`

Translates the entire article content, replacing the original.

- **Parameters:**
  - `prompt`: Custom prompt. Supports `{{.TargetLang}}`.

### `translate-content-immersive`

Bilingual translation. Appends the translation after each paragraph of the original text.

- **Parameters:**
  - `prompt`: Custom prompt.

### `summary`

Generates a summary of the article and prepends it to the content.

- **Parameters:**
  - `prompt`: Custom prompt for summarization.

### `introduction`

Generates a short introduction or "lead-in" for the article.

- **Parameters:**
  - `prompt`: Custom prompt.

### `ai-content-process`

Process article content using LLM according to custom rules and insert the generated result at a specified position.

- **Parameters:**
  - `rule` (**Required**): Instruction for processing each article content. Example: "Summarize key points and list action items."
  - `extra-payload` (Default: `article_content`): Comma-separated list of extra information to send to the LLM. Supported: `article_summary`, `article_content`, `article_date`, `raw_rss_item`.
  - `placement` (Default: `prepend`): Where to write the generated content. Supported: `prepend`, `replace`, `append`.

### `beautify-content`

Re-formats the article using LLM to fix layout, remove ads, and standardizing Markdown, then converts back to clean HTML.

- **Parameters:**
  - `prompt`: Instructions for the "editor" persona.

---

## AI Filtering

Advanced filtering using semantic understanding.

### `ignore-advertorial`

Uses LLM to detect if an article is an advertorial or soft advertisement (evaluating both title and content) and removes it.

- **Parameters:**
  - `prompt-for-exclude`: A prompt that should return `true` if the item is an ad.

### `llm-filter`

Generic LLM-based filter. You define the condition for **exclusion**. The LLM evaluates both the article title and content against this condition.

- **Parameters:**
  - `filter_condition`: A natural language question/condition. If the LLM answers "yes" (true), the item is **removed**.
  - _Example:_ "Is this article about sports?" (Removes sports articles).

### `embedding-filter`

Semantic topic filtering powered by an embedding model. Instead of asking a chat model to judge every item, FeedCraft converts both your topic anchors and each article into vectors, compares them with cosine similarity, and then keeps or removes matching items.

:::tip
Use `embedding-filter` when you need fast, repeatable topic filtering such as "keep AI infrastructure news" or "remove sports articles". Use `llm-filter` when the rule needs nuanced reasoning, policy interpretation, or structured decisions.
:::

#### Environment variables

Set a dedicated embedding endpoint when possible:

```bash
FC_EMBEDDING_API_TYPE=openai
FC_EMBEDDING_API_BASE=https://api.openai.com/v1
FC_EMBEDDING_API_KEY=sk-your-api-key
FC_EMBEDDING_API_MODEL=text-embedding-3-small
FC_EMBEDDING_BATCH_SIZE=5
FC_EMBEDDING_MAX_INPUT_CHARS=8000
```

Supported `FC_EMBEDDING_API_TYPE` values:

- `openai`: OpenAI or OpenAI-compatible embedding endpoint.
- `gemini`: Gemini through its OpenAI-compatible embedding endpoint. Set `FC_EMBEDDING_API_BASE` and `FC_EMBEDDING_API_MODEL` explicitly.
- `ollama`: Local Ollama embedding model. Set `FC_EMBEDDING_API_BASE`, for example `http://localhost:11434`, and an embedding model such as `nomic-embed-text` or `bge-m3`.

If `FC_EMBEDDING_API_TYPE`, `FC_EMBEDDING_API_BASE`, and `FC_EMBEDDING_API_KEY` are not set, FeedCraft falls back to the matching `FC_LLM_API_TYPE`, `FC_LLM_API_BASE`, and `FC_LLM_API_KEY` values. The embedding model name is independent: set `FC_EMBEDDING_API_MODEL` to a real embedding model, or FeedCraft uses its OpenAI default when the API type is `openai`. FeedCraft does not reuse `FC_LLM_API_MODEL` because that value is usually a chat model.

`FC_EMBEDDING_MAX_INPUT_CHARS` is a final safety cap applied to every text sent to the embedding service, including any `instruction` prefix. It is a character budget, not an exact tokenizer count. Keep it at or below a conservative value for your model's token window, for example `8000` for an 8k-token embedding model.

#### Parameters

- `anchors` (**required**): One topic anchor per line. Each anchor should describe what you want to match, for example:

  ```text
  artificial intelligence infrastructure
  machine learning research
  large language model deployment
  ```

- `threshold` (default: `0.6`): Cosine similarity threshold from `0` to `1`. Higher values are stricter. Start with `0.6`, lower it if relevant articles are missing, and raise it if unrelated articles slip through.
- `mode` (default: `include`): `include` keeps matching items; `exclude` removes matching items.
- `max_content_length` (default: `2000`): Maximum article content characters used by this AtomCraft before the final `FC_EMBEDDING_MAX_INPUT_CHARS` safety cap is applied.
- `instruction` (optional): Text prefix prepended to every embedding input. Leave it empty unless your model benefits from a fixed task prefix.

#### Admin UI workflow

1. Open **Worktable → AtomCraft**.
2. Create a new AtomCraft, for example `ai-news-only`.
3. Select template `embedding-filter`.
4. Fill `anchors` with one topic per line.
5. Keep `mode=include` to keep matching items, or switch to `exclude` to remove matching items.
6. Save the AtomCraft.
7. Use it in a FlowCraft, Recipe, Feed Compare, or directly:

```text
/craft/ai-news-only?input_url=https%3A%2F%2Fexample.com%2Ffeed.xml
```

#### Troubleshooting

- **"anchors parameter is required"**: Add at least one non-empty line to `anchors`.
- **"FC_EMBEDDING_API_MODEL must be set"**: Configure a single embedding model. Chat models are not suitable here.
- **All items are removed**: Lower `threshold`, add broader anchors, or increase `max_content_length`.
- **Too many unrelated items remain**: Raise `threshold`, make anchors more specific, or switch from broad category names to representative phrases.
- **Provider says the input is too long**: Lower `FC_EMBEDDING_MAX_INPUT_CHARS`, lower `max_content_length`, or shorten `instruction`.
