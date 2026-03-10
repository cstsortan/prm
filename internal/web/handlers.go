package web

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/service"
)

// JSON response helpers

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /api/dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	stats, cfg, err := s.svc.Dashboard()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"project": cfg.ProjectName,
		"stats":   stats,
	})
}

// entityJSON is the JSON representation of an entity for the API.
type entityJSON struct {
	ID               string          `json:"id"`
	Type             model.EntityType `json:"type"`
	Slug             string          `json:"slug"`
	Title            string          `json:"title"`
	Description      string          `json:"description,omitempty"`
	Status           model.Status    `json:"status"`
	Priority         model.Priority  `json:"priority"`
	Tags             []string        `json:"tags,omitempty"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
	StartedAt        *string         `json:"started_at,omitempty"`
	CompletedAt      *string         `json:"completed_at,omitempty"`
	DueDate          *string         `json:"due_date,omitempty"`
	Dependencies     []string        `json:"dependencies,omitempty"`
	Comments         []model.Comment `json:"comments,omitempty"`
	ParentID         string          `json:"parent_id,omitempty"`
	Children         []string        `json:"children,omitempty"`
	Severity         model.Severity  `json:"severity,omitempty"`
	StepsToReproduce string          `json:"steps_to_reproduce,omitempty"`
}

func toEntityJSON(e *model.Entity) entityJSON {
	ej := entityJSON{
		ID:               e.ID,
		Type:             e.Type,
		Slug:             e.Slug,
		Title:            e.Title,
		Description:      e.Description,
		Status:           e.Status,
		Priority:         e.Priority,
		Tags:             e.Tags,
		CreatedAt:        e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Dependencies:     e.Dependencies,
		Comments:         e.Comments,
		ParentID:         e.ParentID,
		Children:         e.Children,
		Severity:         e.Severity,
		StepsToReproduce: e.StepsToReproduce,
	}
	if e.StartedAt != nil {
		t := e.StartedAt.Format("2006-01-02T15:04:05Z")
		ej.StartedAt = &t
	}
	if e.CompletedAt != nil {
		t := e.CompletedAt.Format("2006-01-02T15:04:05Z")
		ej.CompletedAt = &t
	}
	if e.DueDate != nil {
		t := e.DueDate.Format("2006-01-02T15:04:05Z")
		ej.DueDate = &t
	}
	return ej
}

// GET /api/entities?type=epic,story&status=todo,in-progress&priority=high&tags=backend&sort=priority&desc=true
func (s *Server) handleListEntities(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := service.ListFilter{
		Types:      service.ParseTypes(q.Get("type")),
		Statuses:   service.ParseStatuses(q.Get("status")),
		Priorities: service.ParsePriorities(q.Get("priority")),
		Tags:       service.ParseTags(q.Get("tags")),
	}

	if sort, ok := service.ParseSortField(q.Get("sort")); ok {
		filter.Sort = sort
	}
	filter.SortDesc = q.Get("desc") == "true"
	filter.IncludeArchived = q.Get("archived") == "true"

	results, err := s.svc.List(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	entities := make([]entityJSON, 0, len(results))
	for _, r := range results {
		entities = append(entities, toEntityJSON(r.Entity))
	}

	writeJSON(w, http.StatusOK, entities)
}

// GET /api/entities/{id}
func (s *Server) handleGetEntity(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entity, _, readme, deps, err := s.svc.ShowEntity(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Resolve children slugs to full entities
	var childEntities []entityJSON
	if len(entity.Children) > 0 {
		idx, _ := s.svc.Store.ReadIndex()
		if idx != nil {
			for _, childSlug := range entity.Children {
				result, err := s.svc.Store.Resolve(idx, childSlug)
				if err == nil {
					childEntities = append(childEntities, toEntityJSON(result.Entity))
				}
			}
		}
	}

	resp := map[string]any{
		"entity":       toEntityJSON(entity),
		"readme":       readme,
		"dependencies": deps,
	}
	if childEntities != nil {
		resp["children_entities"] = childEntities
	}
	writeJSON(w, http.StatusOK, resp)
}

// treeNodeJSON is the JSON representation of a tree node.
type treeNodeJSON struct {
	Entity   entityJSON     `json:"entity"`
	Children []treeNodeJSON `json:"children,omitempty"`
}

func toTreeNodeJSON(n *service.TreeNode) treeNodeJSON {
	node := treeNodeJSON{Entity: toEntityJSON(n.Entity)}
	for _, child := range n.Children {
		node.Children = append(node.Children, toTreeNodeJSON(child))
	}
	return node
}

// GET /api/tree or GET /api/tree/{id}
func (s *Server) handleTree(w http.ResponseWriter, r *http.Request) {
	ref := r.PathValue("id")
	nodes, err := s.svc.Tree(ref)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result := make([]treeNodeJSON, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, toTreeNodeJSON(n))
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/search?q=auth&type=task&status=todo
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	filter := service.ListFilter{
		Types:    service.ParseTypes(q.Get("type")),
		Statuses: service.ParseStatuses(q.Get("status")),
	}

	results, err := s.svc.Search(query, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type searchResultJSON struct {
		Entity entityJSON `json:"entity"`
		Score  int        `json:"score"`
	}

	out := make([]searchResultJSON, 0, len(results))
	for _, r := range results {
		out = append(out, searchResultJSON{
			Entity: toEntityJSON(r.Entity),
			Score:  r.Score,
		})
	}

	writeJSON(w, http.StatusOK, out)
}

// PATCH /api/entities/{id}/status
func (s *Server) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	entity, err := s.svc.SetStatus(id, body.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toEntityJSON(entity))
}

// POST /api/entities/{id}/comments
func (s *Server) handleAddComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if body.Author == "" {
		body.Author = "web"
	}

	entity, err := s.svc.AddComment(id, body.Author, body.Text)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toEntityJSON(entity))
}

// GET / — serve the embedded React SPA.
// For any path that doesn't match a static file, serve index.html (SPA client-side routing).
func (s *Server) handleFrontend(w http.ResponseWriter, r *http.Request) {
	fsys, err := staticFiles()
	if err != nil {
		http.Error(w, "frontend not available", http.StatusInternalServerError)
		return
	}

	// Try to serve the exact file first
	path := r.URL.Path
	if path == "/" {
		path = "index.html"
	} else {
		path = path[1:] // strip leading /
	}

	if f, err := fsys.Open(path); err == nil {
		f.Close()
		http.FileServerFS(fsys).ServeHTTP(w, r)
		return
	}

	// Fallback to index.html for SPA routing
	data, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		http.Error(w, "frontend not built (run 'make build-web')", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
