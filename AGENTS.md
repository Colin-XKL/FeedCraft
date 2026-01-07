## Project Overview

FeedCraft is a Go-based RSS feed processing middleware that allows users to transform RSS feeds using various processing "processors". The application supports AI-powered feed manipulation, including translation, fulltext extraction, summarization, and content filtering.

## Architecture

### Core Components

- **Backend (Go)**: Main application in `cmd/main.go` using Gin web framework
- **Frontend (Vue.js)**: Admin panel in `web/admin/` using Vue 3 + TypeScript + Arco Design
- **Database**: SQLite with GORM ORM for data persistence
- **Cache**: Redis for caching feed processing results
- **AI Integration**: OpenAI-compatible LLM support for content processing

### Key Directories

- `cmd/`: Application entry point
- `internal/craft/`: Core feed processing logic and tool templates
- `internal/controller/`: HTTP API handlers
- `internal/dao/`: Database access objects
- `internal/adapter/`: External service adapters (LLM, etc.)
- `internal/util/`: Utility functions and helpers
- `internal/router/`: Route definitions and middleware
- `internal/recipe/`: Legacy package name for channel management
- `web/admin/`: Vue.js frontend application
- `build/`: Docker build configuration

### Core Concepts

- **Tool** (formerly CraftAtom): Individual processing units (e.g., translate, summarize, extract fulltext)
- **Blueprint** (formerly CraftFlow): Sequential combinations of multiple Tools
- **Channel** (formerly Recipe): Configuration that applies specific tools/blueprints to RSS feed URLs
- **Portable Mode**: Direct URL-based processing: `/craft/{processor-name}?input_url={rss-url}`
- **Dock Mode**: Advanced configuration through admin panel

## Development Commands

### Backend Development

```bash
# Build the application
go build -o feed-craft ./cmd/main.go

# Run the application
./feed-craft

# Run all tests
go test ./...

# Reset admin password
./feed-craft reset-password

# Run with development configuration
ENV=dev ./feed-craft
```

### Frontend Development

```bash
cd web/admin

# Install dependencies
pnpm install

# Development server
pnpm run dev

# Build for production
pnpm run build

# Type checking
pnpm run type:check

# Lint and fix
pnpm run lint-staged
```

### Hooks when you finish a task

When you finish a task, run these commands to check the code do not break a build.

```bash
# run lint
task lint
# backend build test
task backend-build
# frontend build test
task frontend-build
```

## Database Schema

The application uses SQLite with the following main entities:

- `users`: User management
- `tools`: Custom tool definitions (formerly craft_atoms)
- `blueprints`: Blueprint configurations (formerly craft_flows)
- `channels`: Feed processing channels (formerly recipes)

Database migrations are handled automatically in `dao.MigrateDatabases()`.

## Built-in Tool Templates

The system includes numerous built-in tool templates in `internal/craft/entry.go`:

- `proxy`: Simple RSS proxy
- `limit`: Limit number of items
- `fulltext`: Extract full article content
- `fulltext-plus`: Browser-based fulltext extraction
- `translate-title/content`: AI-powered translation
- `summary/introduction`: AI-generated summaries
- `ignore-advertorial`: AI content filtering
- `cleanup`: HTML content cleanup
- `keyword`: Keyword-based filtering

## Testing

- Unit tests are located alongside source files (e.g., `internal/util/language_test.go`)
- Use `go test ./...` to run all tests
- The frontend project currently does not have test scripts configured. If tests are added, they would be located in `web/admin/`.

## Common Development Tasks

### Adding New Tool Template

1. Define tool logic in `internal/craft/`
2. Add parameters template in `internal/craft/entry.go`
3. Implement option function
4. Add to `GetSysCraftTemplateDict()`

### Adding New API Endpoint

1. Create controller function in `internal/controller/`
2. Add route in `internal/router/registry.go`
3. Add authentication middleware if needed

### Frontend Development

- Vue 3 composition API with TypeScript
- Arco Design component library
- Pinia for state management
- Axios for API calls
- Vite for build tooling
