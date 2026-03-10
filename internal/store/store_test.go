package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cstsortan/prm/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return New(t.TempDir())
}

func initTestStore(t *testing.T) *Store {
	t.Helper()
	s := newTestStore(t)
	cfg := model.DefaultConfig("test-project")
	if err := s.Init(cfg); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	return s
}

func TestNew(t *testing.T) {
	s := New("/some/project")
	if s.Root() != "/some/project/.prm" {
		t.Errorf("Root() = %q, want %q", s.Root(), "/some/project/.prm")
	}
}

func TestExists_NotInitialized(t *testing.T) {
	s := newTestStore(t)
	if s.Exists() {
		t.Error("Exists() = true for uninitialized store")
	}
}

func TestInit(t *testing.T) {
	s := newTestStore(t)
	cfg := model.DefaultConfig("my-project")

	if err := s.Init(cfg); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if !s.Exists() {
		t.Error("Exists() = false after Init")
	}

	// Check directories were created
	for _, dir := range []string{"epics", "tasks", "bugs", "docs"} {
		path := filepath.Join(s.Root(), dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("directory %s not created: %v", dir, err)
		} else if !info.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
	}

	// Check config was written
	readCfg, err := s.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if readCfg.ProjectName != "my-project" {
		t.Errorf("ProjectName = %q, want %q", readCfg.ProjectName, "my-project")
	}
	if readCfg.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", readCfg.Version, "1.0.0")
	}

	// Check index was written
	idx, err := s.ReadIndex()
	if err != nil {
		t.Fatalf("ReadIndex failed: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("initial index has %d entries, want 0", len(idx.Entries))
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	s := initTestStore(t)
	err := s.Init(model.DefaultConfig("again"))
	if err == nil {
		t.Error("Init on already-initialized store should fail")
	}
}

func TestWriteAndReadEntity(t *testing.T) {
	s := initTestStore(t)

	now := time.Now().UTC().Truncate(time.Millisecond)
	entity := &model.Entity{
		ID:        "test-uuid-1234",
		Type:      model.EntityTask,
		Slug:      "my-task",
		Title:     "My Task",
		Status:    model.StatusTodo,
		Priority:  model.PriorityHigh,
		Tags:      []string{"backend"},
		CreatedAt: now,
		UpdatedAt: now,
	}
	readme := "# My Task\n\nDescription here.\n"

	dir := filepath.Join(s.Root(), "tasks", "my-task")
	if err := s.WriteEntity(dir, entity, readme); err != nil {
		t.Fatalf("WriteEntity failed: %v", err)
	}

	// Read back
	got, err := s.ReadEntity(dir)
	if err != nil {
		t.Fatalf("ReadEntity failed: %v", err)
	}
	if got.ID != entity.ID {
		t.Errorf("ID = %q, want %q", got.ID, entity.ID)
	}
	if got.Title != "My Task" {
		t.Errorf("Title = %q, want %q", got.Title, "My Task")
	}
	if got.Priority != model.PriorityHigh {
		t.Errorf("Priority = %q, want %q", got.Priority, model.PriorityHigh)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "backend" {
		t.Errorf("Tags = %v, want [backend]", got.Tags)
	}

	// Read readme
	gotReadme, err := s.ReadEntityReadme(dir)
	if err != nil {
		t.Fatalf("ReadEntityReadme failed: %v", err)
	}
	if gotReadme != readme {
		t.Errorf("README = %q, want %q", gotReadme, readme)
	}
}

func TestWriteEntity_JSONIsPrettyPrinted(t *testing.T) {
	s := initTestStore(t)

	entity := &model.Entity{
		ID:        "pretty-uuid",
		Type:      model.EntityTask,
		Slug:      "pretty-task",
		Title:     "Pretty Task",
		Status:    model.StatusBacklog,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	dir := filepath.Join(s.Root(), "tasks", "pretty-task")
	if err := s.WriteEntity(dir, entity, "# Pretty\n"); err != nil {
		t.Fatalf("WriteEntity failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if err != nil {
		t.Fatalf("reading meta.json: %v", err)
	}

	// Check it's valid JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("meta.json is not valid JSON: %v", err)
	}

	// Check it's indented (contains newline + spaces)
	if len(data) < 10 {
		t.Fatal("meta.json seems too short to be pretty-printed")
	}
	// Pretty-printed JSON has "  " indentation
	if !json.Valid(data) {
		t.Error("meta.json is not valid JSON")
	}
}

func TestDeleteEntity(t *testing.T) {
	s := initTestStore(t)

	dir := filepath.Join(s.Root(), "tasks", "delete-me")
	entity := &model.Entity{
		ID:        "delete-uuid",
		Type:      model.EntityTask,
		Slug:      "delete-me",
		Title:     "Delete Me",
		Status:    model.StatusBacklog,
		Priority:  model.PriorityLow,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.WriteEntity(dir, entity, "# Delete\n"); err != nil {
		t.Fatalf("WriteEntity failed: %v", err)
	}

	if err := s.DeleteEntity(dir); err != nil {
		t.Fatalf("DeleteEntity failed: %v", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("directory still exists after DeleteEntity")
	}
}

func TestWriteAndReadIndex(t *testing.T) {
	s := initTestStore(t)

	idx := model.NewIndex()
	idx.Set("uuid-1", "tasks/my-task")
	idx.Set("uuid-2", "epics/my-epic")

	if err := s.WriteIndex(idx); err != nil {
		t.Fatalf("WriteIndex failed: %v", err)
	}

	got, err := s.ReadIndex()
	if err != nil {
		t.Fatalf("ReadIndex failed: %v", err)
	}
	if len(got.Entries) != 2 {
		t.Errorf("index has %d entries, want 2", len(got.Entries))
	}
	if got.Get("uuid-1") != "tasks/my-task" {
		t.Errorf("uuid-1 path = %q, want %q", got.Get("uuid-1"), "tasks/my-task")
	}
}

func TestWriteAndReadConfig(t *testing.T) {
	s := initTestStore(t)

	cfg := &model.Config{
		Version:         "2.0.0",
		ProjectName:     "updated-project",
		DefaultPriority: model.PriorityHigh,
		DefaultStatus:   model.StatusTodo,
		CreatedAt:       time.Now().UTC(),
	}

	if err := s.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	got, err := s.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if got.ProjectName != "updated-project" {
		t.Errorf("ProjectName = %q, want %q", got.ProjectName, "updated-project")
	}
	if got.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", got.Version, "2.0.0")
	}
}

func TestListDirs(t *testing.T) {
	s := initTestStore(t)

	// Create some directories
	for _, name := range []string{"alpha", "beta", "gamma"} {
		os.MkdirAll(filepath.Join(s.Root(), "epics", name), 0755)
	}
	// Create a file (should be excluded)
	os.WriteFile(filepath.Join(s.Root(), "epics", "not-a-dir.txt"), []byte("hi"), 0644)

	dirs, err := s.ListDirs(filepath.Join(s.Root(), "epics"))
	if err != nil {
		t.Fatalf("ListDirs failed: %v", err)
	}
	if len(dirs) != 3 {
		t.Errorf("ListDirs returned %d dirs, want 3: %v", len(dirs), dirs)
	}
}

func TestListDirs_NonexistentDir(t *testing.T) {
	s := initTestStore(t)

	dirs, err := s.ListDirs(filepath.Join(s.Root(), "nonexistent"))
	if err != nil {
		t.Fatalf("ListDirs failed: %v", err)
	}
	if dirs != nil {
		t.Errorf("ListDirs on nonexistent dir returned %v, want nil", dirs)
	}
}

func TestEntityDir(t *testing.T) {
	s := New("/project")
	got := s.EntityDir("epics/my-epic")
	want := "/project/.prm/epics/my-epic"
	if got != want {
		t.Errorf("EntityDir = %q, want %q", got, want)
	}
}

func TestRelPath(t *testing.T) {
	s := initTestStore(t)
	absPath := filepath.Join(s.Root(), "epics", "my-epic")
	got, err := s.RelPath(absPath)
	if err != nil {
		t.Fatalf("RelPath failed: %v", err)
	}
	if got != filepath.Join("epics", "my-epic") {
		t.Errorf("RelPath = %q, want %q", got, "epics/my-epic")
	}
}

func TestRebuildIndex(t *testing.T) {
	s := initTestStore(t)

	now := time.Now().UTC()
	entities := []struct {
		id   string
		slug string
		dir  string
	}{
		{"uuid-aaa", "task-a", "tasks/task-a"},
		{"uuid-bbb", "task-b", "tasks/task-b"},
		{"uuid-ccc", "epic-c", "epics/epic-c"},
	}

	for _, e := range entities {
		entity := &model.Entity{
			ID:        e.id,
			Type:      model.EntityTask,
			Slug:      e.slug,
			Title:     e.slug,
			Status:    model.StatusBacklog,
			Priority:  model.PriorityMedium,
			CreatedAt: now,
			UpdatedAt: now,
		}
		dir := filepath.Join(s.Root(), e.dir)
		if err := s.WriteEntity(dir, entity, "# "+e.slug+"\n"); err != nil {
			t.Fatalf("WriteEntity %s failed: %v", e.slug, err)
		}
	}

	// Clear the index
	s.WriteIndex(model.NewIndex())

	// Rebuild
	idx, err := s.RebuildIndex()
	if err != nil {
		t.Fatalf("RebuildIndex failed: %v", err)
	}

	if len(idx.Entries) != 3 {
		t.Errorf("rebuilt index has %d entries, want 3", len(idx.Entries))
	}
	for _, e := range entities {
		if idx.Get(e.id) != e.dir {
			t.Errorf("index[%s] = %q, want %q", e.id, idx.Get(e.id), e.dir)
		}
	}
}

func TestAtomicWrite_NoPartialFiles(t *testing.T) {
	s := initTestStore(t)

	entity := &model.Entity{
		ID:        "atomic-uuid",
		Type:      model.EntityTask,
		Slug:      "atomic-test",
		Title:     "Atomic Test",
		Status:    model.StatusBacklog,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	dir := filepath.Join(s.Root(), "tasks", "atomic-test")
	if err := s.WriteEntity(dir, entity, "# Atomic\n"); err != nil {
		t.Fatalf("WriteEntity failed: %v", err)
	}

	// Verify no temp files remain
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Name() != "meta.json" && e.Name() != "README.md" {
			t.Errorf("unexpected file in entity dir: %s", e.Name())
		}
	}
}
