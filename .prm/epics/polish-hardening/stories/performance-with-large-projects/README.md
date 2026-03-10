# Performance with Large Projects

Ensure all operations complete in <100ms for projects with up to 1000 entities (per non-functional requirements).

## Potential issues

- **UniqueSlug unbounded loop**: `UniqueSlug()` increments a counter with no upper bound. With thousands of slug collisions, this loops indefinitely. Add a max collision count.
- **Full index scan on every List/Search**: Every `List()` call reads all entities from disk via the index. For 1000+ entities this may be slow. Consider caching or lazy loading.
- **Tree builds by walking filesystem**: `Tree()` reads entities one at a time. Could be batched.
- **No index validation on startup**: Stale index entries pointing to deleted files cause per-entry filesystem stat calls that fail.
