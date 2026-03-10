package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Simple Title", "simple-title"},
		{"Fix Login Bug", "fix-login-bug"},
		{"UPPERCASE TITLE", "uppercase-title"},
		{"  leading and trailing spaces  ", "leading-and-trailing-spaces"},
		{"multiple   spaces", "multiple-spaces"},
		{"special!@#$%chars", "specialchars"},
		{"dashes-already-here", "dashes-already-here"},
		{"numbers 123 in title", "numbers-123-in-title"},
		{"", "untitled"},
		{"!!!@@@###", "untitled"},

		// Bug fix: ampersand and plus
		{"Polish & Hardening", "polish-and-hardening"},
		{"Rock & Roll", "rock-and-roll"},
		{"C++ Programming", "c-plus-plus-programming"},
		{"A & B & C", "a-and-b-and-c"},
		{"One+Two", "one-plus-two"},

		// Unicode
		{"Café Résumé", "cafe-resume"},
		{"über cool", "uber-cool"},
		{"naïve approach", "naive-approach"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := GenerateSlug(tt.input)
			if got != tt.want {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestUniqueSlug_NoCollision(t *testing.T) {
	dir := t.TempDir()
	got := UniqueSlug(dir, "my-task")
	if got != "my-task" {
		t.Errorf("UniqueSlug = %q, want %q", got, "my-task")
	}
}

func TestUniqueSlug_WithCollision(t *testing.T) {
	dir := t.TempDir()

	// Create colliding directories
	os.MkdirAll(filepath.Join(dir, "my-task"), 0755)
	os.MkdirAll(filepath.Join(dir, "my-task-2"), 0755)

	got := UniqueSlug(dir, "my-task")
	if got != "my-task-3" {
		t.Errorf("UniqueSlug = %q, want %q", got, "my-task-3")
	}
}

func TestUniqueSlug_MultipleCollisions(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "task"), 0755)
	os.MkdirAll(filepath.Join(dir, "task-2"), 0755)
	os.MkdirAll(filepath.Join(dir, "task-3"), 0755)

	got := UniqueSlug(dir, "task")
	if got != "task-4" {
		t.Errorf("UniqueSlug = %q, want %q", got, "task-4")
	}
}
