# PRM - Project Resource Manager: Requirements

## Overview

PRM is a CLI project management tool written in Go that lives inside a codebase. It stores all data as `.json` and `.md` files in a `.prm/` directory, making project state human-readable, editable, and version-controllable.

## Entity Model

### Entity Types

| Type | Description | Can have children |
|------|-------------|-------------------|
| **Epic** | Large body of work, top of hierarchy | Stories |
| **Story** | A feature or deliverable within an Epic | Sub-Tasks |
| **Sub-Task** | Granular work item within a Story | No |
| **Task** | Standalone work item (not part of hierarchy) | No |
| **Bug** | Defect report | No |
| **Doc** | Free-form markdown document | No |

### Entity Fields (common to all work items)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | UUID v4 | Yes | Unique identifier |
| `type` | enum | Yes | `epic`, `story`, `sub-task`, `task`, `bug` |
| `slug` | string | Yes | URL-safe identifier derived from title, used as directory name |
| `title` | string | Yes | Human-readable title |
| `description` | string | No | Short summary (full description lives in README.md) |
| `status` | enum | Yes | `backlog`, `todo`, `in-progress`, `review`, `done`, `cancelled`, `archived` |
| `priority` | enum | Yes | `low`, `medium`, `high`, `critical` |
| `tags` | string[] | No | Labels (e.g., `backend`, `frontend`, `security`, `reviewer`) |
| `created_at` | ISO 8601 | Yes | Auto-set on creation |
| `updated_at` | ISO 8601 | Yes | Auto-set on any change |
| `started_at` | ISO 8601 | No | Set when status moves to `in-progress` |
| `completed_at` | ISO 8601 | No | Set when status moves to `done` |
| `due_date` | ISO 8601 | No | Target completion date |
| `dependencies` | string[] | No | List of entity IDs this item depends on |
| `comments` | Comment[] | No | Array of `{author, text, created_at}` |
| `parent_id` | UUID | No | ID of parent entity (for hierarchy navigation) |

### Additional fields by type

- **Epic**: `children` object mapping to story slugs
- **Story**: `children` object mapping to sub-task slugs, `epic_id`
- **Sub-Task**: `story_id`
- **Bug**: `severity` (enum: `cosmetic`, `minor`, `major`, `blocker`), `steps_to_reproduce` (string)

## File & Storage Design

### Directory Layout

```
project-root/
  .prm/
    prm.json                          # Root config (project name, default settings)
    epics/
      <epic-slug>/
        meta.json                     # Epic metadata + children mapping
        README.md                     # Human-readable epic description
        stories/
          <story-slug>/
            meta.json                 # Story metadata + children mapping
            README.md                 # Human-readable story description
            tasks/                    # Sub-tasks of this story
              <subtask-slug>/
                meta.json
                README.md
    tasks/                            # Standalone tasks (not in hierarchy)
      <task-slug>/
        meta.json
        README.md
    bugs/
      <bug-slug>/
        meta.json
        README.md
    docs/
      <doc-slug>.md                   # Free-form documents
    index.json                        # Global index for fast lookups (id->path mapping)
```

### Conventions

- **`meta.json`**: Every entity directory contains exactly one `meta.json` with structured metadata.
- **`README.md`**: Every entity directory contains a `README.md` with the full human-readable description. This is the file users edit directly for long-form content.
- **`index.json`**: A flat map of `{id: relative_path}` for all entities. Rebuilt/updated on every write operation. Enables O(1) lookup by ID without traversing the tree.
- **Slugs**: Auto-generated from title (`"Fix Login Bug"` -> `fix-login-bug`). Must be unique within their parent directory. If collision, append `-2`, `-3`, etc.

### prm.json (Root Config)

```json
{
  "version": "1.0.0",
  "project_name": "My Project",
  "default_priority": "medium",
  "default_status": "backlog",
  "tags": ["backend", "frontend", "security", "reviewer", "devops"],
  "created_at": "2024-01-01T00:00:00Z"
}
```

## CLI Interface

### Binary Location

The compiled binary lives at `./bin/prm` within the repository.

### Modes

1. **Non-Interactive (flag-based)**: All parameters passed as flags. Suitable for scripting and agent use.
2. **Interactive (wizard)**: Guided prompts when required flags are omitted. Uses arrow-key selection for enums, fuzzy search for existing entities.

### Command Structure

```
prm init                                          # Initialize .prm/ in current directory
prm <entity> <action> [flags]                     # General pattern

# Entity CRUD
prm epic create --title "..." [--priority ...] [--tags ...] [--due ...] [--body "..."]
prm epic show <id-or-slug>
prm epic update <id-or-slug> [--title ...] [--priority ...] [--tags ...] [--body "..."]
prm epic update <id-or-slug> [--depends "slug1,slug2"] [--clear-depends] [--clear-due]
prm epic edit <id-or-slug>                         # Open README.md in $EDITOR
prm epic delete <id-or-slug>

prm story create --epic <epic-id-or-slug> --title "..." [flags]
prm task create --title "..." [flags]               # Standalone task
prm subtask create --story <story-id-or-slug> --title "..." [flags]
prm bug create --title "..." [--severity ...] [flags]

# Cross-entity operations
prm comment <id-or-slug> --text "..."              # Add comment to any entity
prm status <id-or-slug> <new-status>               # Quick status change
prm archive <id-or-slug>                           # Set status to archived
prm move <story-slug> --to <new-epic-slug>         # Move story to different epic
prm move <subtask-slug> --to <new-story-slug>      # Move sub-task to different story
prm move <subtask-slug> --standalone                # Promote sub-task to standalone task
prm deps <id-or-slug>                              # Show dependency graph
prm search <query> [--type ...] [--status ...] [--tag ...]

# Views
prm list [--type ...] [--status ...] [--priority ...] [--tag ...]
prm list [--sort priority|status|title|created|updated|type] [--desc]
prm list [--archived]                               # Include archived items
prm dashboard                                       # Summary stats
prm tree [<epic-id-or-slug>]                        # Show hierarchy tree

# Docs
prm doc create --title "..." [--body "..."]         # Creates .md file
prm doc list
prm doc show <slug>

# Data & Maintenance
prm export [--format csv|json] [--output file]      # Export data
prm reindex                                         # Rebuild index.json
prm install-skill [--force]                         # Install Claude Code skill

# Web UI
prm web [--port 3141] [--no-open]                   # Launch web dashboard
```

### ID Resolution

Commands accept either a full UUID, a partial UUID prefix (minimum 4 chars), or a slug. The resolver checks in order: exact UUID match -> UUID prefix match -> slug match. Ambiguous matches produce an error listing candidates.

## Views & Reports

### `prm list`

Tabular output with columns: ID (short), Type, Title, Status, Priority, Tags. Supports filtering and sorting.

### `prm dashboard`

```
Project: My Project
==================
Epics:     3 total (1 in-progress, 1 todo, 1 done)
Stories:   12 total (4 in-progress, 5 todo, 3 done)
Tasks:     8 total (2 in-progress, 3 todo, 3 done)
Bugs:      2 total (1 major, 1 minor)

By Priority:
  Critical: 2  |  High: 5  |  Medium: 10  |  Low: 3

Recently Updated:
  [IN-PROGRESS] Fix auth token refresh (bug) - 2h ago
  [TODO] Add user profile page (story) - 5h ago
```

### `prm tree`

```
Epic: User Authentication [IN-PROGRESS]
  Story: Login Flow [DONE]
    Sub-Task: Design login page [DONE]
    Sub-Task: Implement OAuth [DONE]
  Story: Session Management [IN-PROGRESS]
    Sub-Task: Token refresh [IN-PROGRESS]
    Sub-Task: Logout cleanup [TODO]
```

## Search

Full-text search across titles, descriptions, README.md content, comments, and tags. Returns ranked results. Supports filters:

```
prm search "auth" --type epic,story --status todo,in-progress
```

## Export

```
prm export --format csv --output tasks.csv
prm export --format json --output backup.json
```

CSV includes all fields flattened. JSON exports the full entity tree.

## Web UI

PRM includes a built-in web dashboard served from an embedded React single-page application.

```
prm web                     # Opens browser on default port (3141)
prm web --port 8080         # Custom port
prm web --no-open           # Don't auto-open browser
```

### Pages

- **Dashboard** — Overview stats (same data as `prm dashboard`)
- **List** — Filterable entity list with status/type/priority filters
- **Detail** — Full entity view with README content, comments, and status updates
- **Tree** — Interactive hierarchy view
- **Search** — Full-text search with results

### API Endpoints

```
GET   /api/dashboard              # Dashboard stats
GET   /api/entities               # List entities (supports query filters)
GET   /api/entities/{id}          # Single entity detail
GET   /api/tree                   # Full hierarchy tree
GET   /api/tree/{id}              # Tree from specific entity
GET   /api/search?q=...           # Search
PATCH /api/entities/{id}/status   # Update status
POST  /api/entities/{id}/comments # Add comment
```

### Development

- **Frontend dev server**: `make dev-web` (Vite on port 5173, proxies API to 3141)
- **Backend API**: `./bin/prm web --no-open`
- **Production build**: `make build-web` compiles the React app and embeds it into the Go binary

## Go Tech Stack

| Dependency | Purpose |
|-----------|---------|
| `github.com/spf13/cobra` | CLI framework, subcommands |
| `github.com/charmbracelet/bubbletea` | Interactive TUI / wizard mode |
| `github.com/charmbracelet/bubbles` | TUI components (text input, selection, tables) |
| `github.com/charmbracelet/lipgloss` | Terminal styling and layout |
| `github.com/google/uuid` | UUID generation |
| `github.com/sahilm/fuzzy` | Fuzzy matching for search and ID resolution |

## Non-Functional Requirements

- **Performance**: All operations should complete in <100ms for projects with up to 1000 entities.
- **Portability**: Single binary, no external dependencies at runtime.
- **Human-readable**: All `.json` files are pretty-printed (indented). All `.md` files follow standard markdown.
- **Idempotent index**: `index.json` can always be rebuilt from the file tree via `prm reindex`.
- **Graceful errors**: Clear error messages. Never corrupt existing files on failure (write to temp, then rename).
