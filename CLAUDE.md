# PRM - Project Resource Manager

## What is this?

A CLI project management tool written in Go. It stores all data as `.json` and `.md` files in a `.prm/` directory within a codebase. Install globally with `go install github.com/cstsortan/prm@latest`, or build locally with `make build` (outputs to `./bin/prm`).

## Quick Reference

- **Requirements**: `docs/REQUIREMENTS.md`
- **Architecture**: `docs/ARCHITECTURE.md`
- **Build**: `make build` (outputs to `bin/prm`)
- **Run dev**: `go run . <command>`
- **Test**: `go test ./...`

## Project Structure

```
cmd/           # Cobra command definitions (one file per command)
internal/
  model/       # Data structs (entity, epic, story, task, bug, comment, config, index)
  store/       # Filesystem persistence (read/write JSON+MD, index, slug, ID resolution)
  service/     # Business logic (CRUD, search, export, dashboard, tree)
  tui/         # Interactive mode (bubbletea wizards)
  render/      # Output formatting (tables, trees, dashboard, lipgloss styles)
main.go        # Entry point
```

## Coding Conventions

- **Go style**: Follow standard Go conventions. Use `gofmt`. No exported symbols unless needed outside the package.
- **Error handling**: Wrap errors with context using `fmt.Errorf("operation: %w", err)`. Never silently ignore errors.
- **File I/O**: All writes go through `store/` which uses atomic write (temp file + rename). Never write directly to `.prm/` files from `cmd/` or `service/`.
- **JSON**: Always pretty-print with 2-space indent. Use `json.MarshalIndent(v, "", "  ")`.
- **Naming**: Files use lowercase with underscores for multi-word. Packages use single lowercase words.
- **Tests**: Place tests in `_test.go` files next to the code. Use table-driven tests. Use `t.TempDir()` for filesystem tests.

## Entity Hierarchy

```
Epic -> Story -> Sub-Task
Task (standalone, no parent)
Bug (standalone, no parent)
Doc (standalone markdown file)
```

## .prm/ Directory Structure

```
.prm/
  prm.json          # Project config
  index.json         # UUID -> path lookup cache
  epics/<slug>/meta.json + README.md
    stories/<slug>/meta.json + README.md
      tasks/<slug>/meta.json + README.md      # (sub-tasks)
  tasks/<slug>/meta.json + README.md           # (standalone)
  bugs/<slug>/meta.json + README.md
  docs/<slug>.md
```

- `meta.json` = structured data (never hand-edit except in emergencies)
- `README.md` = human-readable description (freely editable)
- `index.json` = rebuild with `prm reindex` if corrupted

## CLI Pattern

Every command follows: `prm <entity> <action> [flags]`

Commands accept IDs as: full UUID, partial UUID (min 4 chars), or slug.

Non-interactive: all params via flags. Interactive: missing required fields trigger TUI prompts.

## Dependencies

- `cobra` for CLI
- `bubbletea` + `bubbles` for TUI
- `lipgloss` for styling
- `google/uuid` for IDs
- `sahilm/fuzzy` for fuzzy search

## Task Tracking

When the user requests multiple features, a batch of work, or any non-trivial implementation effort, you MUST break the work into tasks using `prm task create` (or epics/stories as appropriate) BEFORE writing any code. Update task status with `prm status <slug> in-progress` when starting and `prm status <slug> done` when finished. This is not optional — plan first, then execute.

## Keeping the Skill Up to Date

After adding, changing, or removing CLI commands, flags, or statuses, you MUST update `.claude/commands/prm.md` to reflect the changes. This is the skill file that teaches future Claude Code sessions how to use PRM — if it's stale, the agent will use outdated commands.

## When Building Features

1. Start with the **model** struct in `internal/model/`
2. Add **store** methods in `internal/store/` for persistence
3. Add **service** methods in `internal/service/` for business logic
4. Add the **cobra command** in `cmd/`
5. Add **render** formatting if the command has output
6. Add **tui** components if the command needs interactive mode
7. Write tests for store and service layers

## Using PRM as an Agent

You can manage project work items using the `prm` binary. Common operations:

```bash
# Initialize (first time only)
prm init

# Create items
prm epic create --title "Feature Name" --priority high --tags "backend,security"
prm story create --epic <epic-slug> --title "Story Name"
prm task create --title "Standalone Task" --priority medium
prm bug create --title "Bug Title" --severity major

# Create with detailed description (README.md content)
prm epic create --title "Feature Name" --body "# Feature Name

## Goals
- Goal 1
- Goal 2"

# Update status
prm status <id-or-slug> in-progress
prm status <id-or-slug> done

# Update detailed description
prm epic update <id-or-slug> --body "# Updated content..."

# Edit README.md in $EDITOR
prm epic edit <id-or-slug>

# Add comments
prm comment <id-or-slug> --text "Completed the implementation"

# View work
prm list --status todo,in-progress
prm tree <epic-slug>
prm dashboard
prm search "auth"
```
