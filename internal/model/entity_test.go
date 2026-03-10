package model

import "testing"

func TestShortID(t *testing.T) {
	e := &Entity{ID: "abcdefgh-1234-5678-9012-345678901234"}
	if got := e.ShortID(); got != "abcdefgh" {
		t.Errorf("ShortID = %q, want abcdefgh", got)
	}
}

func TestShortID_Short(t *testing.T) {
	e := &Entity{ID: "abc"}
	if got := e.ShortID(); got != "abc" {
		t.Errorf("ShortID = %q, want abc", got)
	}
}

func TestIsTerminal(t *testing.T) {
	tests := []struct {
		typ  EntityType
		want bool
	}{
		{EntityEpic, false},
		{EntityStory, false},
		{EntitySubTask, true},
		{EntityTask, true},
		{EntityBug, true},
	}
	for _, tt := range tests {
		e := &Entity{Type: tt.typ}
		if got := e.IsTerminal(); got != tt.want {
			t.Errorf("IsTerminal(%s) = %v, want %v", tt.typ, got, tt.want)
		}
	}
}

func TestChildDir(t *testing.T) {
	tests := []struct {
		typ  EntityType
		want string
	}{
		{EntityEpic, "stories"},
		{EntityStory, "tasks"},
		{EntitySubTask, ""},
		{EntityTask, ""},
		{EntityBug, ""},
	}
	for _, tt := range tests {
		e := &Entity{Type: tt.typ}
		if got := e.ChildDir(); got != tt.want {
			t.Errorf("ChildDir(%s) = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"backlog", true},
		{"todo", true},
		{"in-progress", true},
		{"review", true},
		{"done", true},
		{"cancelled", true},
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		_, ok := ParseStatus(tt.input)
		if ok != tt.valid {
			t.Errorf("ParseStatus(%q) valid = %v, want %v", tt.input, ok, tt.valid)
		}
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"low", true},
		{"medium", true},
		{"high", true},
		{"critical", true},
		{"urgent", false},
		{"", false},
	}
	for _, tt := range tests {
		_, ok := ParsePriority(tt.input)
		if ok != tt.valid {
			t.Errorf("ParsePriority(%q) valid = %v, want %v", tt.input, ok, tt.valid)
		}
	}
}

func TestParseSeverity(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"cosmetic", true},
		{"minor", true},
		{"major", true},
		{"blocker", true},
		{"severe", false},
	}
	for _, tt := range tests {
		_, ok := ParseSeverity(tt.input)
		if ok != tt.valid {
			t.Errorf("ParseSeverity(%q) valid = %v, want %v", tt.input, ok, tt.valid)
		}
	}
}
