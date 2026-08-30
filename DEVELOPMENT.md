# Development Guide

This guide captures the practical setup and testing notes for developing FeedCraft locally or in Cursor Cloud.

## Prerequisites

- Go, matching the version in `go.mod`.
- Node.js and `pnpm` for the admin UI in `web/admin/`.
- `task` for repository-level build, lint, and format commands.
- Redis-compatible cache for local runtime testing. FeedCraft uses Redis for craft and LLM result caching.
- SQLite storage directory for local development.
- Optional OpenAI-compatible LLM endpoint for AI crafts.

## Dependency Setup

Install Go dependencies:

```bash
go mod download
```

Install frontend dependencies:

```bash
cd web/admin
pnpm install
```

Do not hand-edit generated dependency lock entries. Use the package manager when adding or updating dependencies.

## Backend Development

Build the backend:

```bash
go build -o feed-craft ./cmd/main.go
```

Run the backend in development mode:

```bash
ENV=dev \
LISTEN_ADDR=:8080 \
FC_DB_SQLITE_PATH=/tmp/feedcraft-dev-db \
FC_REDIS_URI=redis://127.0.0.1:6379 \
FC_LLM_API_BASE=http://127.0.0.1:18080/v1 \
FC_LLM_API_MODEL=mock-model \
FC_LLM_API_KEY=mock-key \
./feed-craft
```

Notes:

- `FC_DB_SQLITE_PATH` is a directory; FeedCraft creates `feed-craft.db` inside it.
- Default admin credentials for a fresh dev database are `admin` / `adminadmin`.
- If you do not need AI crafts, LLM variables can be omitted, but AI-related flows will fail when invoked.
- Use `./feed-craft reset-password` if the local admin password is unknown.

## Frontend Development

Start the admin UI:

```bash
cd web/admin
pnpm run dev -- --host 0.0.0.0
```

The dev server proxies `/api` to `http://localhost:8080`. Open:

```text
http://localhost:5173
```

Useful routes:

- AtomCraft management: `http://localhost:5173/worktable/craft_atom`
- Feed Compare: `http://localhost:5173/tools/feed_compare`

## Mocking LLM and Feeds for Manual Testing

For AI craft UI testing, use an OpenAI-compatible mock server that implements `POST /v1/chat/completions` and returns deterministic JSON. This keeps tests repeatable and avoids spending tokens.

Example `ai-filter` decision response:

```json
{
  "choices": [
    {
      "message": {
        "content": "{\"reason\":\"deterministic mock decision\",\"result\":\"keep\"}"
      }
    }
  ]
}
```

When testing Feed Compare, avoid `127.0.0.1` or private IP feed URLs. `validateFeedViewerURL()` blocks loopback and private IPs by design. Use a public-looking hostname such as `http://example.com/feed.xml` and, if needed, route backend fetches through an HTTP proxy that serves local test RSS content.

## Verification Commands

Run focused Go tests while iterating:

```bash
go test ./internal/craft
```

Run all Go tests:

```bash
go test ./...
```

Run frontend checks:

```bash
cd web/admin
pnpm run format
pnpm run type:check
pnpm run lint
pnpm run build
```

Run repository checks before submitting:

```bash
task fix
go test ./...
task backend-build
task frontend-build
```

If `task fix` fails, inspect the exact lint output. During recent development it was blocked by pre-existing staticcheck warnings in `internal/controller/feed_viewer.go` about capitalized error strings; do not silently ignore failures.

## Development Lessons from AI Craft Work

- Keep legacy `CraftOption` behavior and native processor behavior separate unless the task explicitly asks to touch both.
- AI filters should default to keeping articles when LLM calls, JSON parsing, or result validation fail. This avoids accidental data loss.
- Cache both intermediate AI payloads, such as generated summaries, and final AI decisions when the final prompt and payload are stable.
- Browser testing is valuable for AtomCraft changes because template parameters are rendered dynamically from backend metadata.
- For UI parameter widgets, keep backend storage simple. The `ai-filter` `extra-payload` field is edited as a multi-select in the UI but is still saved as a comma-separated string.
