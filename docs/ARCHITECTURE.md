# PRM Architecture

## Project Structure (Go Source)

```
prm/
  bin/                    # Compiled binary output
    prm
  cmd/                    # CLI command definitions (cobra)
    root.go               # Root command, global flags, config loading
    init.go               # prm init
    epic.go               # prm epic [create|show|update|edit|delete]
    story.go              # prm story [create|show|update|edit|delete]
    task.go               # prm task [create|show|update|edit|delete]
    subtask.go            # prm subtask [create|show|update|edit|delete]
    bug.go                # prm bug [create|show|update|edit|delete]
    doc.go                # prm doc [create|list|show]
    comment.go            # prm comment <id> --text "..."
    status.go             # prm status <id> <status>
    archive.go            # prm archive <id>
    move.go               # prm move <id> --to <parent> | --standalone
    deps.go               # prm deps <id>
    search.go             # prm search <query>
    list.go               # prm list (cross-entity, with --sort/--desc/--archived)
    dashboard.go          # prm dashboard
    tree.go               # prm tree
    export.go             # prm export
    reindex.go            # prm reindex
    skill.go              # prm install-skill [--force]
    web.go                # prm web [--port] [--no-open]
    integration_test.go   # CLI integration tests
  internal/
    model/                # Data structures
      entity.go           # Base entity struct, status/priority/severity constants
      config.go           # prm.json config struct
      index.go            # Index struct and operations
    store/                # File I/O and persistence
      store.go            # Store struct, read/write meta.json + README.md, atomic writes
      index.go            # Index management (load, save, rebuild)
      slug.go             # Slug generation and collision handling
      resolve.go          # ID/slug resolution (UUID, partial UUID, slug)
    service/              # Business logic
      service.go          # Service struct, core methods (create, show, delete, status, comment)
      manage.go           # Update, move, archive, dependency management
      list.go             # List with filtering and sorting
      search.go           # Full-text search across entities
      export.go           # CSV/JSON export logic
      tree.go             # Hierarchy tree building (includes dashboard stats)
    tui/                  # Interactive mode (bubbletea)
      tui.go              # Terminal detection (IsInteractive)
      create.go           # Interactive entity creation wizard
      select.go           # Entity selector (fuzzy search)
      confirm.go          # Confirmation prompts
    render/               # Output formatting
      table.go            # Table renderer for list views
      tree.go             # Tree renderer for hierarchy view
      dashboard.go        # Dashboard renderer
      style.go            # Shared lipgloss styles
    web/                  # Web UI server
      server.go           # HTTP server setup, browser auto-open
      routes.go           # API route definitions
      handlers.go         # Request handlers (dashboard, entities, tree, search, status, comments)
      embed.go            # Embedded frontend static assets (go:embed)
  web/                    # Frontend source (React/TypeScript/Vite)
    src/
      App.tsx             # Main app component with routing
      main.tsx            # Entry point
      index.css           # Global styles
      api/
        client.ts         # HTTP API client
        types.ts          # TypeScript type definitions
      components/
        Nav.tsx           # Navigation bar
        TypeIcon.tsx      # Entity type icons
        StatusBadge.tsx   # Status display badges
        PriorityBadge.tsx # Priority display badges
      pages/
        DashboardPage.tsx # Dashboard overview
        ListPage.tsx      # Entity list with filters
        EntityDetailPage.tsx # Single entity detail view
        TreePage.tsx      # Hierarchy tree view
        SearchPage.tsx    # Search results page
    package.json
    vite.config.ts
  go.mod
  go.sum
  main.go                 # Entry point, calls cmd.Execute()
  Makefile                # Build targets (build, build-web, dev, dev-web, test, clean, sync-skill)
  CLAUDE.md               # Instructions for Claude Code
  README.md               # Repository README
  docs/
    REQUIREMENTS.md
    ARCHITECTURE.md
```

## Data Flow

### CLI Path

```
CLI Input (flags / interactive)
  -> cmd/ (cobra command handler)
    -> service/ (business logic, validation)
      -> store/ (filesystem read/write)
        -> .prm/ (JSON + MD files on disk)
      -> model/ (data structs)
    -> render/ (format output)
  -> Terminal Output
```

### Web Path

```
Browser Request
  -> internal/web/ (HTTP server, routes, handlers)
    -> service/ (same business logic as CLI)
      -> store/ (same filesystem layer)
        -> .prm/ (JSON + MD files on disk)
    -> JSON API Response
  -> React SPA (web/src/)
```

The web UI reuses the same `service` and `store` layers as the CLI, ensuring consistent behavior. The React frontend is compiled by Vite and embedded into the Go binary via `go:embed`.

## Key Design Decisions

### 1. meta.json + README.md Split

Every entity stores structured data in `meta.json` and long-form human content in `README.md`. This separation means:
- Programmatic operations only touch `meta.json`
- Users can edit `README.md` freely without risking data corruption
- Diffs in version control are clean and meaningful

### 2. index.json for Fast Lookups

Rather than traversing the directory tree for every operation, `index.json` provides O(1) lookup by ID. It maps `UUID -> relative path`. The index is:
- Updated on every create/delete/move operation
- Rebuildable from scratch via `prm reindex`
- The source of truth is always the file tree (index is a cache)

### 3. Slug-based Directories

Directories are named by slug (derived from title) rather than UUID. This makes the file tree navigable by humans. UUIDs are used for programmatic references (dependencies, parent links) since slugs can theoretically change.

### 4. Interactive Fallback

When a required flag is missing in non-interactive mode, the command errors with a clear message. When running interactively (TTY detected), missing fields trigger a wizard prompt. This is controlled by checking `term.IsTerminal(os.Stdin)`.

### 5. Atomic Writes

All file writes use a write-to-temp-then-rename pattern to prevent corruption if the process is interrupted mid-write.

### 6. Embedded Web UI

The React frontend is built with Vite and the compiled output is copied into `internal/web/static/`. The Go binary embeds these files via `go:embed`, so `prm web` serves a fully functional dashboard from the single binary with no external file dependencies.

### 7. Archive as Status

Rather than deleting entities, the `archive` command sets status to `archived`. Archived items are hidden from `prm list` and `prm dashboard` by default, but can be included with `--archived`. This preserves history while keeping active views clean.

### 8. Move and Reparent

Entities can be moved between parents (`prm move <story> --to <epic>`) or promoted (`prm move <subtask> --standalone`). This involves updating the filesystem directory structure, parent/child references in meta.json, and the index.

## meta.json Examples

### Epic

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "epic",
  "slug": "user-authentication",
  "title": "User Authentication",
  "description": "Implement complete auth system with OAuth and session management",
  "status": "in-progress",
  "priority": "high",
  "tags": ["backend", "security"],
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-20T14:30:00Z",
  "started_at": "2024-01-16T09:00:00Z",
  "completed_at": null,
  "due_date": "2024-03-01T00:00:00Z",
  "dependencies": [],
  "comments": [
    {
      "author": "dev",
      "text": "OAuth provider decided: Google + GitHub",
      "created_at": "2024-01-16T11:00:00Z"
    }
  ],
  "children": ["login-flow", "session-management", "password-reset"]
}
```

### Bug

```json
{
  "id": "f9e8d7c6-b5a4-3210-fedc-ba0987654321",
  "type": "bug",
  "slug": "token-refresh-race-condition",
  "title": "Token refresh race condition",
  "description": "Concurrent requests can trigger multiple token refreshes",
  "status": "todo",
  "priority": "critical",
  "severity": "major",
  "tags": ["backend", "security"],
  "steps_to_reproduce": "1. Open two tabs\n2. Wait for token to expire\n3. Perform action in both tabs simultaneously",
  "created_at": "2024-01-20T14:00:00Z",
  "updated_at": "2024-01-20T14:00:00Z",
  "started_at": null,
  "completed_at": null,
  "due_date": null,
  "dependencies": ["a1b2c3d4-e5f6-7890-abcd-ef1234567890"],
  "comments": [],
  "parent_id": null
}
```

## README.md Example

```markdown
# User Authentication

Implement complete auth system with OAuth and session management.

## Goals

- Support Google and GitHub OAuth providers
- Implement secure session management with refresh tokens
- Add password reset flow via email

## Technical Notes

We'll use JWT for access tokens (15min expiry) and opaque refresh tokens
stored server-side. Session data goes in Redis.

## Acceptance Criteria

- [ ] User can sign in with Google
- [ ] User can sign in with GitHub
- [ ] Tokens refresh transparently
- [ ] User can reset password via email
```

## Build & Run

```makefile
# Makefile
.PHONY: build clean dev test sync-skill build-web dev-web

sync-skill:                           # Copy skill file to cmd/ for embedding
	cp .claude/commands/prm.md cmd/skill_prm.md

build-web:                            # Build React frontend with Vite
	cd web && npm run build
	rm -rf internal/web/static/*
	cp -r web/dist/* internal/web/static/

build: sync-skill build-web           # Full build (skill + web + Go binary)
	go build -o bin/prm .

clean:
	rm -f bin/prm
	rm -rf web/dist internal/web/static/*

dev:                                  # Run without compiling
	go run . $(ARGS)

dev-web:                              # Vite dev server (port 5173, proxies to 3141)
	cd web && npm run dev

test:
	go test ./...
```

## Implementation Status

All phases are complete and the tool is fully functional.

### Phase 1: Foundation (Complete)
- `prm init`
- Model structs and store layer (meta.json + README.md read/write)
- Index management
- Slug generation

### Phase 2: Core CRUD (Complete)
- Epic, Story, Sub-Task, Task, Bug create/show/update/edit/delete
- Comment and status commands
- ID/slug resolution
- Archive, move/reparent, dependencies

### Phase 3: Views (Complete)
- `prm list` with filtering, sorting (`--sort`, `--desc`), and `--archived`
- `prm tree`
- `prm dashboard`

### Phase 4: Interactive Mode (Complete)
- Bubbletea wizards for create commands
- Fuzzy entity selector
- Confirmation prompts

### Phase 5: Search & Export (Complete)
- Full-text search with type/status/tag filters
- CSV and JSON export
- `prm reindex`

### Phase 6: Claude Code Skill (Complete)
- `.claude/commands/prm.md` skill for agent use
- `prm install-skill` command to install/update the skill

### Phase 7: Web UI (Complete)
- Embedded React/TypeScript SPA built with Vite
- Dashboard, list, detail, tree, and search pages
- REST API backed by the same service layer as the CLI
- `prm web` command with `--port` and `--no-open` flags
