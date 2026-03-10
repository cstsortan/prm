package model

import "testing"

func TestIndex_SetAndGet(t *testing.T) {
	idx := NewIndex()
	idx.Set("uuid-1", "tasks/my-task")

	if got := idx.Get("uuid-1"); got != "tasks/my-task" {
		t.Errorf("Get = %q, want tasks/my-task", got)
	}
}

func TestIndex_Get_NotFound(t *testing.T) {
	idx := NewIndex()
	if got := idx.Get("nonexistent"); got != "" {
		t.Errorf("Get(nonexistent) = %q, want empty", got)
	}
}

func TestIndex_Delete(t *testing.T) {
	idx := NewIndex()
	idx.Set("uuid-1", "path")
	idx.Delete("uuid-1")

	if got := idx.Get("uuid-1"); got != "" {
		t.Errorf("Get after Delete = %q, want empty", got)
	}
}

func TestIndex_FindByPrefix(t *testing.T) {
	idx := NewIndex()
	idx.Set("aaaa-1111", "path1")
	idx.Set("aaaa-2222", "path2")
	idx.Set("bbbb-1111", "path3")

	matches := idx.FindByPrefix("aaaa")
	if len(matches) != 2 {
		t.Errorf("FindByPrefix('aaaa') returned %d, want 2", len(matches))
	}

	matches = idx.FindByPrefix("bbbb")
	if len(matches) != 1 {
		t.Errorf("FindByPrefix('bbbb') returned %d, want 1", len(matches))
	}

	matches = idx.FindByPrefix("cccc")
	if len(matches) != 0 {
		t.Errorf("FindByPrefix('cccc') returned %d, want 0", len(matches))
	}
}

func TestIndex_Overwrite(t *testing.T) {
	idx := NewIndex()
	idx.Set("uuid-1", "old-path")
	idx.Set("uuid-1", "new-path")

	if got := idx.Get("uuid-1"); got != "new-path" {
		t.Errorf("Get after overwrite = %q, want new-path", got)
	}
}
