package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/store"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	s := store.New(dir)
	svc := New(s)
	if err := svc.Init("test-project"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	return svc
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	s := store.New(dir)
	svc := New(s)

	if err := svc.Init("my-project"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if !s.Exists() {
		t.Error("store does not exist after Init")
	}

	cfg, err := s.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if cfg.ProjectName != "my-project" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "my-project")
	}
}

func TestCreateEntity_Epic(t *testing.T) {
	svc := newTestService(t)

	entity, err := svc.CreateEntity(CreateEntityOpts{
		Type:     model.EntityEpic,
		Title:    "User Auth",
		Priority: model.PriorityHigh,
		Status:   model.StatusBacklog,
		Tags:     []string{"backend"},
	})
	if err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}

	if entity.Type != model.EntityEpic {
		t.Errorf("Type = %q, want epic", entity.Type)
	}
	if entity.Slug != "user-auth" {
		t.Errorf("Slug = %q, want user-auth", entity.Slug)
	}
	if entity.ParentID != "" {
		t.Errorf("ParentID = %q, want empty for epic", entity.ParentID)
	}

	// Verify files on disk
	dir := filepath.Join(svc.Store.Root(), "epics", "user-auth")
	if _, err := os.Stat(filepath.Join(dir, "meta.json")); err != nil {
		t.Errorf("meta.json not found: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); err != nil {
		t.Errorf("README.md not found: %v", err)
	}
	// Epic should have stories/ child dir
	if _, err := os.Stat(filepath.Join(dir, "stories")); err != nil {
		t.Errorf("stories/ directory not created: %v", err)
	}

	// Verify index
	idx, _ := svc.Store.ReadIndex()
	if idx.Get(entity.ID) == "" {
		t.Error("entity not in index")
	}
}

func TestCreateEntity_Task(t *testing.T) {
	svc := newTestService(t)

	entity, err := svc.CreateEntity(CreateEntityOpts{
		Type:     model.EntityTask,
		Title:    "Setup CI",
		Priority: model.PriorityMedium,
		Status:   model.StatusTodo,
	})
	if err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}

	if entity.Slug != "setup-ci" {
		t.Errorf("Slug = %q, want setup-ci", entity.Slug)
	}
	// Task is terminal, no child directory
	dir := filepath.Join(svc.Store.Root(), "tasks", "setup-ci")
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() {
			t.Errorf("terminal task should not have subdirectory: %s", e.Name())
		}
	}
}

func TestCreateEntity_Bug(t *testing.T) {
	svc := newTestService(t)

	entity, err := svc.CreateEntity(CreateEntityOpts{
		Type:     model.EntityBug,
		Title:    "Login Crash",
		Priority: model.PriorityCritical,
		Status:   model.StatusBacklog,
		Severity: model.SeverityMajor,
	})
	if err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}
	if entity.Severity != model.SeverityMajor {
		t.Errorf("Severity = %q, want major", entity.Severity)
	}
}

func TestCreateEntity_StoryUnderEpic(t *testing.T) {
	svc := newTestService(t)

	epic, err := svc.CreateEntity(CreateEntityOpts{
		Type:     model.EntityEpic,
		Title:    "Auth",
		Priority: model.PriorityHigh,
		Status:   model.StatusBacklog,
	})
	if err != nil {
		t.Fatalf("creating epic: %v", err)
	}

	story, err := svc.CreateEntity(CreateEntityOpts{
		Type:     model.EntityStory,
		Title:    "Login Flow",
		Priority: model.PriorityMedium,
		Status:   model.StatusBacklog,
		ParentID: epic.ID,
	})
	if err != nil {
		t.Fatalf("creating story: %v", err)
	}

	if story.ParentID != epic.ID {
		t.Errorf("ParentID = %q, want %q", story.ParentID, epic.ID)
	}

	// Verify directory structure
	storyDir := filepath.Join(svc.Store.Root(), "epics", "auth", "stories", "login-flow")
	if _, err := os.Stat(storyDir); err != nil {
		t.Errorf("story dir not created: %v", err)
	}

	// Verify parent has child in children list
	epicEntity, _, _, _, _ := svc.ShowEntity(epic.ID)
	found := false
	for _, c := range epicEntity.Children {
		if c == "login-flow" {
			found = true
		}
	}
	if !found {
		t.Errorf("epic children = %v, want to contain 'login-flow'", epicEntity.Children)
	}
}

func TestCreateEntity_SubTaskUnderStory(t *testing.T) {
	svc := newTestService(t)

	epic, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "E", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})
	story, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityStory, Title: "S", Priority: model.PriorityMedium, Status: model.StatusBacklog, ParentID: epic.ID,
	})
	subtask, err := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntitySubTask, Title: "ST", Priority: model.PriorityLow, Status: model.StatusBacklog, ParentID: story.ID,
	})
	if err != nil {
		t.Fatalf("creating subtask: %v", err)
	}
	if subtask.ParentID != story.ID {
		t.Errorf("subtask ParentID = %q, want %q", subtask.ParentID, story.ID)
	}
}

func TestCreateEntity_StoryWithoutEpic_Fails(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityStory, Title: "Orphan", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})
	if err == nil {
		t.Fatal("creating story without epic should fail")
	}
}

func TestCreateEntity_SubTaskUnderEpic_Fails(t *testing.T) {
	svc := newTestService(t)

	epic, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "E", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	_, err := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntitySubTask, Title: "Wrong", Priority: model.PriorityMedium, Status: model.StatusBacklog, ParentID: epic.ID,
	})
	if err == nil {
		t.Fatal("creating subtask under epic should fail")
	}
	if !strings.Contains(err.Error(), "not a story") {
		t.Errorf("error = %q, want 'not a story'", err.Error())
	}
}

func TestCreateEntity_InProgress_SetsStartedAt(t *testing.T) {
	svc := newTestService(t)

	entity, err := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "Active", Priority: model.PriorityMedium, Status: model.StatusInProgress,
	})
	if err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}
	if entity.StartedAt == nil {
		t.Error("StartedAt should be set when created with in-progress status")
	}
}

func TestCreateEntity_DuplicateSlug_Appends_Number(t *testing.T) {
	svc := newTestService(t)

	e1, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "My Task", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})
	e2, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "My Task", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	if e1.Slug != "my-task" {
		t.Errorf("first slug = %q, want my-task", e1.Slug)
	}
	if e2.Slug != "my-task-2" {
		t.Errorf("second slug = %q, want my-task-2", e2.Slug)
	}
}

func TestSetStatus(t *testing.T) {
	svc := newTestService(t)

	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	// Move to in-progress
	updated, err := svc.SetStatus(entity.ID, "in-progress")
	if err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("Status = %q, want in-progress", updated.Status)
	}
	if updated.StartedAt == nil {
		t.Error("StartedAt should be set on in-progress")
	}

	// Move to done
	done, err := svc.SetStatus(entity.ID, "done")
	if err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}
	if done.CompletedAt == nil {
		t.Error("CompletedAt should be set on done")
	}
	// StartedAt should still be set
	if done.StartedAt == nil {
		t.Error("StartedAt should persist after done")
	}
}

func TestSetStatus_Invalid(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	_, err := svc.SetStatus(entity.ID, "invalid-status")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestSetStatus_Cancelled_SetsCompletedAt(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	updated, err := svc.SetStatus(entity.ID, "cancelled")
	if err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}
	if updated.CompletedAt == nil {
		t.Error("CompletedAt should be set on cancelled")
	}
}

func TestAddComment(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	updated, err := svc.AddComment(entity.ID, "dev", "This is a comment")
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}
	if len(updated.Comments) != 1 {
		t.Fatalf("Comments count = %d, want 1", len(updated.Comments))
	}
	if updated.Comments[0].Author != "dev" {
		t.Errorf("Author = %q, want dev", updated.Comments[0].Author)
	}
	if updated.Comments[0].Text != "This is a comment" {
		t.Errorf("Text = %q, want 'This is a comment'", updated.Comments[0].Text)
	}
}

func TestAddComment_Multiple(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	svc.AddComment(entity.ID, "alice", "First")
	updated, _ := svc.AddComment(entity.ID, "bob", "Second")

	if len(updated.Comments) != 2 {
		t.Errorf("Comments count = %d, want 2", len(updated.Comments))
	}
}

func TestUpdateEntity(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "Old Title", Priority: model.PriorityLow, Status: model.StatusBacklog,
	})

	updated, err := svc.UpdateEntity(entity.ID, UpdateEntityOpts{
		Title:    "New Title",
		Priority: "high",
		Tags:     []string{"new-tag"},
	})
	if err != nil {
		t.Fatalf("UpdateEntity failed: %v", err)
	}
	if updated.Title != "New Title" {
		t.Errorf("Title = %q, want 'New Title'", updated.Title)
	}
	if updated.Priority != model.PriorityHigh {
		t.Errorf("Priority = %q, want high", updated.Priority)
	}
	if len(updated.Tags) != 1 || updated.Tags[0] != "new-tag" {
		t.Errorf("Tags = %v, want [new-tag]", updated.Tags)
	}
}

func TestUpdateEntity_DueDate(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	due := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	updated, _ := svc.UpdateEntity(entity.ID, UpdateEntityOpts{DueDate: &due})
	if updated.DueDate == nil || !updated.DueDate.Equal(due) {
		t.Errorf("DueDate = %v, want %v", updated.DueDate, due)
	}

	// Clear due date
	cleared, _ := svc.UpdateEntity(entity.ID, UpdateEntityOpts{ClearDue: true})
	if cleared.DueDate != nil {
		t.Errorf("DueDate = %v, want nil after clear", cleared.DueDate)
	}
}

func TestUpdateEntity_InvalidPriority(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	_, err := svc.UpdateEntity(entity.ID, UpdateEntityOpts{Priority: "super-high"})
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

func TestUpdateEntity_NoChanges(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "T", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	updated, err := svc.UpdateEntity(entity.ID, UpdateEntityOpts{})
	if err != nil {
		t.Fatalf("UpdateEntity with no changes failed: %v", err)
	}
	if updated.UpdatedAt != entity.UpdatedAt {
		t.Error("UpdatedAt should not change when nothing is modified")
	}
}

func TestDeleteEntity(t *testing.T) {
	svc := newTestService(t)
	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "Delete Me", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	if err := svc.DeleteEntity(entity.ID); err != nil {
		t.Fatalf("DeleteEntity failed: %v", err)
	}

	// Should not be in index
	idx, _ := svc.Store.ReadIndex()
	if idx.Get(entity.ID) != "" {
		t.Error("deleted entity still in index")
	}

	// Directory should be gone
	dir := filepath.Join(svc.Store.Root(), "tasks", "delete-me")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("entity directory still exists")
	}
}

func TestDeleteEntity_RemovesFromParent(t *testing.T) {
	svc := newTestService(t)

	epic, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "Epic", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})
	story, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityStory, Title: "Story", Priority: model.PriorityMedium, Status: model.StatusBacklog, ParentID: epic.ID,
	})

	// Verify parent has child
	epicBefore, _, _, _, _ := svc.ShowEntity(epic.ID)
	if len(epicBefore.Children) != 1 {
		t.Fatalf("epic should have 1 child before delete, has %d", len(epicBefore.Children))
	}

	if err := svc.DeleteEntity(story.ID); err != nil {
		t.Fatalf("DeleteEntity failed: %v", err)
	}

	// Parent should no longer list the child
	epicAfter, _, _, _, _ := svc.ShowEntity(epic.ID)
	if len(epicAfter.Children) != 0 {
		t.Errorf("epic children after delete = %v, want empty", epicAfter.Children)
	}
}

func TestDeleteEntity_CascadesIndex(t *testing.T) {
	svc := newTestService(t)

	epic, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityEpic, Title: "E", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})
	story, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityStory, Title: "S", Priority: model.PriorityMedium, Status: model.StatusBacklog, ParentID: epic.ID,
	})
	subtask, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntitySubTask, Title: "ST", Priority: model.PriorityLow, Status: model.StatusBacklog, ParentID: story.ID,
	})

	// Delete the epic (should cascade remove story and subtask from index)
	if err := svc.DeleteEntity(epic.ID); err != nil {
		t.Fatalf("DeleteEntity failed: %v", err)
	}

	idx, _ := svc.Store.ReadIndex()
	if idx.Get(epic.ID) != "" {
		t.Error("epic still in index")
	}
	if idx.Get(story.ID) != "" {
		t.Error("story still in index after parent delete")
	}
	if idx.Get(subtask.ID) != "" {
		t.Error("subtask still in index after grandparent delete")
	}
}

func TestShowEntity(t *testing.T) {
	svc := newTestService(t)

	entity, _ := svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "Show Me", Description: "A description", Priority: model.PriorityHigh, Status: model.StatusTodo,
	})

	got, dir, readme, _, err := svc.ShowEntity(entity.ID)
	if err != nil {
		t.Fatalf("ShowEntity failed: %v", err)
	}
	if got.Title != "Show Me" {
		t.Errorf("Title = %q, want 'Show Me'", got.Title)
	}
	if dir == "" {
		t.Error("Dir should not be empty")
	}
	if !strings.Contains(readme, "# Show Me") {
		t.Errorf("README = %q, want to contain '# Show Me'", readme)
	}
}

func TestShowEntity_BySlug(t *testing.T) {
	svc := newTestService(t)

	svc.CreateEntity(CreateEntityOpts{
		Type: model.EntityTask, Title: "Find By Slug", Priority: model.PriorityMedium, Status: model.StatusBacklog,
	})

	got, _, _, _, err := svc.ShowEntity("find-by-slug")
	if err != nil {
		t.Fatalf("ShowEntity by slug failed: %v", err)
	}
	if got.Title != "Find By Slug" {
		t.Errorf("Title = %q, want 'Find By Slug'", got.Title)
	}
}
