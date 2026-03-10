# Body & Edit for README.md

Added two ways to set detailed README.md content on entities:

## --body flag (non-interactive)
- `prm epic create --title 'X' --body '# Content...'` sets README.md on create
- `prm epic update <ref> --body '# New content'` replaces README.md on update
- Works for all entity types: epic, story, task, subtask, bug

## edit subcommand (interactive)
- `prm epic edit <ref>` opens README.md in $EDITOR (falls back to vi)
- Works for all entity types