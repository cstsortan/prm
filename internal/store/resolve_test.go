package store

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cstsortan/prm/internal/model"
)

func setupResolveStore(t *testing.T) (*Store, *model.Index) {
	t.Helper()
	s := initTestStore(t)

	now := time.Now().UTC()
	entities := []struct {
		id   string
		slug string
		path string
		typ  model.EntityType
	}{
		{"aaaa1111-0000-0000-0000-000000000001", "login-flow", "epics/auth/stories/login-flow", model.EntityStory},
		{"bbbb2222-0000-0000-0000-000000000002", "auth", "epics/auth", model.EntityEpic},
		{"cccc3333-0000-0000-0000-000000000003", "fix-bug", "bugs/fix-bug", model.EntityBug},
		{"dddd4444-0000-0000-0000-000000000004", "setup-ci", "tasks/setup-ci", model.EntityTask},
		{"dddd4444-0000-0000-0000-000000000099", "deploy-cd", "tasks/deploy-cd", model.EntityTask},
	}

	idx := model.NewIndex()
	for _, e := range entities {
		entity := &model.Entity{
			ID:        e.id,
			Type:      e.typ,
			Slug:      e.slug,
			Title:     e.slug,
			Status:    model.StatusBacklog,
			Priority:  model.PriorityMedium,
			CreatedAt: now,
			UpdatedAt: now,
		}
		dir := filepath.Join(s.Root(), e.path)
		if err := s.WriteEntity(dir, entity, "# "+e.slug+"\n"); err != nil {
			t.Fatalf("WriteEntity %s: %v", e.slug, err)
		}
		idx.Set(e.id, e.path)
	}

	if err := s.WriteIndex(idx); err != nil {
		t.Fatalf("WriteIndex: %v", err)
	}

	return s, idx
}

func TestResolve_ExactUUID(t *testing.T) {
	s, idx := setupResolveStore(t)

	result, err := s.Resolve(idx, "bbbb2222-0000-0000-0000-000000000002")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result.Entity.Slug != "auth" {
		t.Errorf("Slug = %q, want %q", result.Entity.Slug, "auth")
	}
	if result.Entity.Type != model.EntityEpic {
		t.Errorf("Type = %q, want %q", result.Entity.Type, model.EntityEpic)
	}
}

func TestResolve_UUIDPrefix(t *testing.T) {
	s, idx := setupResolveStore(t)

	result, err := s.Resolve(idx, "aaaa")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result.Entity.Slug != "login-flow" {
		t.Errorf("Slug = %q, want %q", result.Entity.Slug, "login-flow")
	}
}

func TestResolve_UUIDPrefix_Ambiguous(t *testing.T) {
	s, idx := setupResolveStore(t)

	// "dddd" matches both dddd4444...001 and dddd4444...099
	_, err := s.Resolve(idx, "dddd")
	if err == nil {
		t.Fatal("expected ambiguous error, got nil")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("error = %q, want it to contain 'ambiguous'", err.Error())
	}
}

func TestResolve_Slug(t *testing.T) {
	s, idx := setupResolveStore(t)

	result, err := s.Resolve(idx, "fix-bug")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if result.Entity.ID != "cccc3333-0000-0000-0000-000000000003" {
		t.Errorf("ID = %q, want cccc3333...", result.Entity.ID)
	}
}

func TestResolve_NotFound(t *testing.T) {
	s, idx := setupResolveStore(t)

	_, err := s.Resolve(idx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent ref, got nil")
	}
	if !strings.Contains(err.Error(), "no entity found") {
		t.Errorf("error = %q, want 'no entity found'", err.Error())
	}
}

func TestResolve_ShortPrefix_TooShort(t *testing.T) {
	s, idx := setupResolveStore(t)

	// Less than 4 chars should not do prefix match, fall through to slug
	_, err := s.Resolve(idx, "aaa")
	if err == nil {
		t.Fatal("expected error for 3-char non-slug ref")
	}
}

func TestResolve_PrefersExactUUID_OverSlug(t *testing.T) {
	s := initTestStore(t)
	now := time.Now().UTC()

	// Create an entity whose slug happens to look like a UUID prefix
	entity := &model.Entity{
		ID:        "real-uuid-here-0000-000000000001",
		Type:      model.EntityTask,
		Slug:      "abcd",
		Title:     "ABCD",
		Status:    model.StatusBacklog,
		Priority:  model.PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}
	dir := filepath.Join(s.Root(), "tasks", "abcd")
	s.WriteEntity(dir, entity, "# ABCD\n")

	idx := model.NewIndex()
	idx.Set(entity.ID, "tasks/abcd")

	// Resolve by the full UUID should work
	result, err := s.Resolve(idx, "real-uuid-here-0000-000000000001")
	if err != nil {
		t.Fatalf("Resolve by UUID failed: %v", err)
	}
	if result.Entity.Slug != "abcd" {
		t.Errorf("wrong entity resolved")
	}

	// Resolve by slug "abcd" should also work
	result, err = s.Resolve(idx, "abcd")
	if err != nil {
		t.Fatalf("Resolve by slug failed: %v", err)
	}
	if result.Entity.ID != "real-uuid-here-0000-000000000001" {
		t.Errorf("wrong entity resolved by slug")
	}
}
