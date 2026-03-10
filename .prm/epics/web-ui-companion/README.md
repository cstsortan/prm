# Web UI Companion

## Overview
A web-based dashboard bundled with the PRM binary, launched via `prm web`. Built with a Go HTTP backend serving a React SPA.

## Architecture
- **Backend**: Go `net/http` server in `internal/web/` wrapping the existing service layer
- **Frontend**: React SPA in `web/` with its own build pipeline
- **Embedding**: Go `embed` package bundles the built frontend into the binary
- **API**: JSON REST endpoints mirroring CLI capabilities
- **Launch**: `prm web [--port 8080]` starts server and opens browser

## API Endpoints
- `GET /api/dashboard` — stats summary
- `GET /api/entities` — list with query params for filtering/sorting
- `GET /api/entities/:id` — single entity detail + README
- `GET /api/tree/:id` — epic tree structure
- `GET /api/search?q=` — fuzzy search
- `PATCH /api/entities/:id/status` — update status
- `POST /api/entities/:id/comments` — add comment

## Frontend Pages
- Dashboard (stats, charts, recent activity)
- List view (filterable/sortable table)
- Entity detail (metadata, README, comments)
- Tree view (epic hierarchy)

## Build Integration
- `make build-web` — builds React app into `web/dist/`
- `make build` — depends on build-web, embeds dist into Go binary