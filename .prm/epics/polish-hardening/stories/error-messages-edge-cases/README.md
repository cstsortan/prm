# Error Messages & Edge Cases

Improve error messages and handle edge cases that currently cause silent failures or confusing output.

## Issues to fix

- **readJSON missing context**: `store.readJSON()` returns bare `os.ReadFile` and `json.Unmarshal` errors with no file path context. Wrap with file path.
- **Silent WriteEntity error in DeleteEntity**: `manage.go` line 198 — when deleting a child entity, the parent's children list update ignores the WriteEntity error. Must check and return it.
- **Silent entity skipping**: `List()`, `Search()`, `Tree()` all silently `continue` when an indexed entity can't be read. At minimum, warn on stderr.
- **init.go flag bug**: Line 21 reads the `"name"` flag twice instead of using the working directory basename. Fix default project name logic.
- **collectAndRemoveFromIndex silent failures**: Silently returns on read errors, leaving orphaned index entries.
- **Doc create not atomic**: `doc.go` uses `os.WriteFile` directly instead of going through the store's atomic write (temp + rename) pattern.
- **DeleteEntity parent cleanup**: If parent can't be resolved (already deleted or corrupted), child is deleted anyway but parent's children list retains the dangling slug.
