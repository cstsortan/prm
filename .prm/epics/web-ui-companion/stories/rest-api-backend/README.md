# REST API Backend

Go HTTP server with JSON API endpoints wrapping the existing service layer.

## Scope
- Router setup and middleware (CORS, JSON content-type, error handling)
- All read endpoints (dashboard, list, detail, tree, search)
- Write endpoints (status update, comments)
- Server lifecycle (graceful shutdown)
- Package: `internal/web/`