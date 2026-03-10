# CLI Command & Build Integration

## Scope
- `cmd/web.go` — cobra command with --port flag
- Auto-open browser on launch
- Graceful shutdown on Ctrl+C
- Makefile targets: `build-web`, update `build` to depend on it
- Dev mode: proxy Vite dev server during development