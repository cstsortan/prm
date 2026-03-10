package service

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cstsortan/prm/internal/model"
)

func seedEntities(t *testing.T, svc *Service) {
	t.Helper()

	entities := []CreateEntityOpts{
		{Type: model.EntityEpic, Title: "Auth Epic", Priority: model.PriorityHigh, Status: model.StatusInProgress, Tags: []string{"backend", "security"}},
		{Type: model.EntityEpic, Title: "UI Epic", Priority: model.PriorityMedium, Status: model.StatusBacklog, Tags: []string{"frontend"}},
		{Type: model.EntityTask, Title: "CI Setup", Priority: model.PriorityMedium, Status: model.StatusDone, Tags: []string{"devops"}},
		{Type: model.EntityTask, Title: "Docs", Priority: model.PriorityLow, Status: model.StatusBacklog, Tags: []string{"docs"}},
		{Type: model.EntityBug, Title: "Login Crash", Priority: model.PriorityCritical, Status: model.StatusTodo, Severity: model.SeverityMajor, Tags: []string{"backend"}},
	}

	for _, opts := range entities {
		if _, err := svc.CreateEntity(opts); err != nil {
			t.Fatalf("seeding entity %q: %v", opts.Title, err)
		}
	}
}

func TestList_All(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("List returned %d results, want 5", len(results))
	}
}

func TestList_FilterByType(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{Types: []model.EntityType{model.EntityBug}})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("List(type=bug) returned %d results, want 1", len(results))
	}
	if results[0].Entity.Type != model.EntityBug {
		t.Errorf("Type = %q, want bug", results[0].Entity.Type)
	}
}

func TestList_FilterByStatus(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{Statuses: []model.Status{model.StatusBacklog}})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("List(status=backlog) returned %d results, want 2", len(results))
	}
}

func TestList_FilterByPriority(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{Priorities: []model.Priority{model.PriorityCritical}})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("List(priority=critical) returned %d results, want 1", len(results))
	}
}

func TestList_FilterByTag(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{Tags: []string{"backend"}})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("List(tag=backend) returned %d results, want 2", len(results))
	}
}

func TestList_MultipleFilters(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.List(ListFilter{
		Types:    []model.EntityType{model.EntityEpic},
		Statuses: []model.Status{model.StatusInProgress},
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("List(type=epic,status=in-progress) returned %d results, want 1", len(results))
	}
}

func TestList_SortedByPriority(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, _ := svc.List(ListFilter{})
	if len(results) < 2 {
		t.Skip("not enough results to test sorting")
	}

	// First result should be critical
	if results[0].Entity.Priority != model.PriorityCritical {
		t.Errorf("first result priority = %q, want critical", results[0].Entity.Priority)
	}
}

func TestSearch(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.Search("auth", ListFilter{})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) == 0 {
		t.Error("Search('auth') returned 0 results")
	}
	// "Auth Epic" should be the top result (title match)
	if results[0].Entity.Title != "Auth Epic" {
		t.Errorf("top result = %q, want 'Auth Epic'", results[0].Entity.Title)
	}
}

func TestSearch_NoResults(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.Search("zzzznonexistent", ListFilter{})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search returned %d results, want 0", len(results))
	}
}

func TestSearch_WithTypeFilter(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	results, err := svc.Search("auth", ListFilter{Types: []model.EntityType{model.EntityBug}})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	// No bugs match "auth"
	if len(results) != 0 {
		t.Errorf("Search(auth, type=bug) returned %d results, want 0", len(results))
	}
}

func TestDashboard(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	stats, cfg, err := svc.Dashboard()
	if err != nil {
		t.Fatalf("Dashboard failed: %v", err)
	}
	if cfg.ProjectName != "test-project" {
		t.Errorf("ProjectName = %q, want test-project", cfg.ProjectName)
	}
	if stats.Total != 5 {
		t.Errorf("Total = %d, want 5", stats.Total)
	}
	if stats.ByType[model.EntityEpic] != 2 {
		t.Errorf("ByType[epic] = %d, want 2", stats.ByType[model.EntityEpic])
	}
	if stats.ByStatus[model.StatusDone] != 1 {
		t.Errorf("ByStatus[done] = %d, want 1", stats.ByStatus[model.StatusDone])
	}
	if stats.BySeverity[model.SeverityMajor] != 1 {
		t.Errorf("BySeverity[major] = %d, want 1", stats.BySeverity[model.SeverityMajor])
	}
}

func TestTree(t *testing.T) {
	svc := newTestService(t)

	epic, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "E", Priority: model.PriorityHigh, Status: model.StatusBacklog,
	})
	story, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityStory, Title: "S", Priority: model.PriorityMedium, Status: model.StatusBacklog, ParentID: epic.ID,
	})
	svc.CreateEntity(CreateEntityOpts{
		Type: model.EntitySubTask, Title: "ST", Priority: model.PriorityLow, Status: model.StatusBacklog, ParentID: story.ID,
	})

	nodes, err := svc.Tree("")
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("Tree returned %d roots, want 1", len(nodes))
	}
	if nodes[0].Entity.Title != "E" {
		t.Errorf("root title = %q, want E", nodes[0].Entity.Title)
	}
	if len(nodes[0].Children) != 1 {
		t.Fatalf("root has %d children, want 1", len(nodes[0].Children))
	}
	if len(nodes[0].Children[0].Children) != 1 {
		t.Errorf("story has %d children, want 1", len(nodes[0].Children[0].Children))
	}
}

func TestTree_SpecificEpic(t *testing.T) {
	svc := newTestService(t)

	epic1, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "Epic 1", Priority: model.PriorityHigh, Status: model.StatusBacklog,
	})
	svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "Epic 2", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	nodes, err := svc.Tree(epic1.ID)
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Errorf("Tree(epic1) returned %d nodes, want 1", len(nodes))
	}
}

func TestExportCSV(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	var buf bytes.Buffer
	if err := svc.ExportCSV(&buf); err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Header + 5 entities
	if len(lines) != 6 {
		t.Errorf("CSV has %d lines, want 6 (1 header + 5 data)", len(lines))
	}
	if !strings.HasPrefix(lines[0], "id,type,slug") {
		t.Errorf("CSV header = %q, want to start with 'id,type,slug'", lines[0])
	}
}

func TestExportJSON(t *testing.T) {
	svc := newTestService(t)
	seedEntities(t, svc)

	var buf bytes.Buffer
	if err := svc.ExportJSON(&buf); err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "[") {
		t.Error("JSON export should start with '['")
	}
	// Should contain all 5 entity IDs
	if strings.Count(output, `"id"`) != 5 {
		t.Errorf("JSON export contains %d entities, want 5", strings.Count(output, `"id"`))
	}
}

func TestParseTypes(t *testing.T) {
	got := ParseTypes("epic,story,bug")
	if len(got) != 3 {
		t.Errorf("ParseTypes returned %d, want 3", len(got))
	}
}

func TestParseTypes_Empty(t *testing.T) {
	got := ParseTypes("")
	if got != nil {
		t.Errorf("ParseTypes('') = %v, want nil", got)
	}
}

func TestParseTags(t *testing.T) {
	got := ParseTags("backend, frontend, security")
	if len(got) != 3 {
		t.Errorf("ParseTags returned %d, want 3", len(got))
	}
	if got[0] != "backend" {
		t.Errorf("first tag = %q, want 'backend'", got[0])
	}
}
