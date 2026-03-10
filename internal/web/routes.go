package web

import (
	"net/http"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("GET /api/dashboard", s.handleDashboard)
	mux.HandleFunc("GET /api/entities", s.handleListEntities)
	mux.HandleFunc("GET /api/entities/{id}", s.handleGetEntity)
	mux.HandleFunc("GET /api/tree", s.handleTree)
	mux.HandleFunc("GET /api/tree/{id}", s.handleTree)
	mux.HandleFunc("GET /api/search", s.handleSearch)
	mux.HandleFunc("PATCH /api/entities/{id}/status", s.handleUpdateStatus)
	mux.HandleFunc("POST /api/entities/{id}/comments", s.handleAddComment)

	// Serve embedded frontend (will be added later)
	mux.HandleFunc("GET /", s.handleFrontend)
}
