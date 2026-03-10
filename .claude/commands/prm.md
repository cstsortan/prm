# PRM Project Management Skill

You are managing project work items using the PRM CLI tool.

## Before Using

1. Check if `prm` is available on PATH by running `which prm`. If found, use `prm` for all commands below.
2. If not on PATH, install it: `curl -sSL https://raw.githubusercontent.com/cstsortan/prm/main/install.sh | sh`. Alternatively, build from source: `make build` (binary at `./bin/prm`).
3. Check if `.prm/` directory exists. If not, initialize: `prm init`

## How to Use

Run PRM commands using the Bash tool in **non-interactive mode** (always provide all required flags).

### Creating Work Items

```bash
# Epic (top-level feature)
prm epic create --title "Title" --priority high --tags "tag1,tag2" --due "2026-03-01"

# Story (belongs to an epic)
prm story create --epic <epic-slug-or-id> --title "Title" --priority medium

# Sub-task (belongs to a story)
prm subtask create --story <story-slug-or-id> --title "Title"

# Standalone task
prm task create --title "Title" --priority medium --tags "tag1"

# Bug
prm bug create --title "Title" --severity major --priority high

# With detailed description (README.md content)
prm epic create --title "Title" --body "# Title

## Goals
- Goal 1
- Goal 2"
```

### Managing Work Items

```bash
# Change status (backlog -> todo -> in-progress -> review -> done | cancelled | archived)
prm status <id-or-slug> in-progress

# Archive an item (hides from list/dashboard by default)
prm archive <id-or-slug>

# Add a comment
prm comment <id-or-slug> --text "Progress update or note"

# Update fields
prm <entity> update <id-or-slug> --title "New Title" --priority critical --tags "new,tags"

# Update detailed description (README.md content)
prm <entity> update <id-or-slug> --body "# New detailed content..."

# Set dependencies
prm <entity> update <id-or-slug> --depends "other-slug,another-slug"
prm <entity> update <id-or-slug> --clear-depends

# Move/reparent entities
prm move <story-slug> --to <new-epic-slug>       # Move story to different epic
prm move <subtask-slug> --to <new-story-slug>    # Move sub-task to different story
prm move <subtask-slug> --standalone              # Promote sub-task to standalone task

# Edit README.md in $EDITOR (interactive)
prm <entity> edit <id-or-slug>

# Delete (prompts for confirmation when interactive)
prm <entity> delete <id-or-slug>
```

### Viewing Work

```bash
# Dashboard overview
prm dashboard

# List with filters and sorting
prm list --status todo,in-progress --priority high,critical --type epic,story
prm list --tag backend
prm list --sort title                  # Sort by: priority, status, title, created, updated, type
prm list --sort status --desc          # Reverse sort order
prm list --archived                    # Include archived items

# Hierarchy tree
prm tree                    # All epics
prm tree <epic-slug>        # Specific epic

# Search
prm search "keyword" --type task,bug

# Show details of one item
prm <entity> show <id-or-slug>

# View dependencies (what it depends on + what depends on it)
prm deps <id-or-slug>
```

### Docs

```bash
prm doc create --title "Design Decision"
prm doc create --title "Guide" --body "# Setup Guide\n\nContent here..."
prm doc list
prm doc show <slug>
```

### Export

```bash
prm export --format csv --output report.csv
prm export --format json --output backup.json
```

### Install Skill

```bash
# Install the PRM Claude Code skill into the current project's .claude/commands/prm.md
prm install-skill

# Overwrite an existing skill file
prm install-skill --force
```

## Workflow Patterns

### Starting work on a feature
1. Create an epic: `prm epic create --title "Feature" --priority high`
2. Break into stories: `prm story create --epic <slug> --title "Part 1"`
3. Add sub-tasks to each story as needed
4. Move items to `in-progress` as you work on them
5. Add comments to track progress
6. Move items to `done` when complete

### Triaging a bug
1. Create bug: `prm bug create --title "..." --severity major`
2. Add steps to reproduce in the comment or edit the README.md directly
3. Update status as you investigate and fix

### Quick status check
Run `prm dashboard` to see overall project health, then `prm list --status in-progress` to see active work.

### Web UI

```bash
# Launch the web dashboard (opens browser automatically)
prm web

# Custom port
prm web --port 8080

# Don't open browser
prm web --no-open
```

The web UI provides a dashboard, entity list (with filters), detail view, tree view, and search — all backed by the same service layer as the CLI.

**Development**: Run `make dev-web` for the Vite dev server (port 5173, proxies API to port 3141). Run `prm web --no-open` in parallel for the API backend.

**Build**: `make build-web` builds the React app and copies it to the Go embed directory. `make build` includes this step automatically.

## Important Notes

- Always use `--text` flag for comments (not positional args)
- Entity types: `epic`, `story`, `subtask`, `task`, `bug`, `doc`
- Statuses: `backlog`, `todo`, `in-progress`, `review`, `done`, `cancelled`, `archived`
- Priorities: `low`, `medium`, `high`, `critical`
- Severities (bugs only): `cosmetic`, `minor`, `major`, `blocker`
- You can reference items by slug (e.g., `user-auth`) or partial UUID (min 4 chars)
