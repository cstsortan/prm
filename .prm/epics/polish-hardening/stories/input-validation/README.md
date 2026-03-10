# Input Validation

Ensure all user inputs are validated before reaching the store layer.

## Issues to fix

- **Whitespace-only titles**: `--title "   "` passes `!= ""` check but generates slug `untitled`. Trim and reject.
- **Slug length limit**: Very long titles produce very long slugs that can exceed filesystem path limits. Cap at 80 chars.
- **Filter flag validation**: `ParseStatuses`, `ParsePriorities`, `ParseTypes` accept invalid values silently (e.g. `--status bogus`). Validate against known enums.
- **Empty tags in update**: `UpdateEntity` assigns tags directly without filtering empty strings. Use `ParseTags` consistently.
- **Comment text/author validation**: `AddComment` in service layer accepts empty text or empty author when called programmatically. Validate at service boundary.
- **Doc path traversal**: `doc show` doesn't sanitize the slug arg — `../../../etc/passwd` could escape `.prm/docs/`. Validate slug contains no path separators.
