# PRM - Project Resource Manager

A CLI project management tool that lives inside your codebase. PRM stores all data as `.json` and `.md` files in a `.prm/` directory, making project state human-readable, editable, and version-controllable.

## Features

- **Hierarchical work items**: Epics → Stories → Sub-Tasks, plus standalone Tasks and Bugs
- **File-based storage**: Everything in `.prm/` as JSON metadata + Markdown descriptions — no database needed
- **Git-friendly**: All data is plain text, perfect for version control and code review
- **Flexible CLI**: Non-interactive mode with flags, or interactive TUI wizards
- **Web dashboard**: Built-in web UI for browsing and searching
- **Claude Code integration**: Ships with a skill file so AI agents can manage work items
- **Search & export**: Fuzzy search across all entities, export to CSV or JSON

## Install

```bash
curl -sSL https://raw.githubusercontent.com/cstsortan/prm/main/install.sh | sh
```

Or with Go:

```bash
go install github.com/cstsortan/prm@latest
```

### Build from source

Requires Go 1.24+ and Node.js (for the web UI).

```bash
git clone https://github.com/cstsortan/prm.git
cd prm
make build
```

The binary is built to `./bin/prm`.

## Quick Start

```bash
# Initialize PRM in your project
prm init

# Create an epic
prm epic create --title "User Authentication" --priority high

# Break it into stories
prm story create --epic user-authentication --title "Login flow"
prm story create --epic user-authentication --title "Password reset"

# Track a standalone task
prm task create --title "Set up CI pipeline" --priority medium

# Log a bug
prm bug create --title "Login timeout too short" --severity major

# See the big picture
prm dashboard
prm tree user-authentication

# Update progress
prm status login-flow in-progress
prm comment login-flow --text "Implemented OAuth provider"
prm status login-flow done
```

## Entity Types

| Type | Description | Parent |
|------|-------------|--------|
| **Epic** | Large body of work | — |
| **Story** | Feature or deliverable | Epic |
| **Sub-Task** | Granular work item | Story |
| **Task** | Standalone work item | — |
| **Bug** | Defect report | — |
| **Doc** | Free-form markdown document | — |

## Storage Layout

```
.prm/
  prm.json              # Project config
  index.json             # UUID → path lookup cache
  epics/<slug>/
    meta.json + README.md
    stories/<slug>/
      meta.json + README.md
      tasks/<slug>/      # Sub-tasks
  tasks/<slug>/          # Standalone tasks
  bugs/<slug>/
  docs/<slug>.md
```

- `meta.json` — structured data (status, priority, tags, timestamps, comments)
- `README.md` — human-readable description, freely editable

## Commands

```
prm init                Initialize a new project
prm epic                Manage epics (create, show, update, edit, delete)
prm story               Manage stories
prm subtask             Manage sub-tasks
prm task                Manage standalone tasks
prm bug                 Manage bugs
prm doc                 Manage documents
prm status <id> <s>     Change an item's status
prm comment <id>        Add a comment
prm move <id>           Reparent an entity
prm archive <id>        Archive an entity
prm deps <id>           Show dependency graph
prm list                List entities with filters and sorting
prm tree                Show hierarchy tree
prm search              Fuzzy search across all entities
prm dashboard           Project summary overview
prm export              Export to CSV or JSON
prm web                 Launch the web dashboard
prm reindex             Rebuild the index from the file tree
prm install-skill       Install Claude Code skill into the project
```

Items can be referenced by slug (e.g., `user-auth`) or partial UUID (min 4 chars).

## Web UI

```bash
prm web           # Opens browser on port 3141
prm web --port 8080 --no-open
```

## Claude Code Integration

PRM includes a skill file that teaches Claude Code how to manage work items:

```bash
prm install-skill
```

This copies the skill to `.claude/commands/prm.md`. Claude Code can then create epics, track tasks, update statuses, and more — all through the CLI.

## License

MIT
