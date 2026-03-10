package model

import (
	"time"
)

// EntityType represents the type of a work item.
type EntityType string

const (
	EntityEpic    EntityType = "epic"
	EntityStory   EntityType = "story"
	EntitySubTask EntityType = "sub-task"
	EntityTask    EntityType = "task"
	EntityBug     EntityType = "bug"
)

// Status represents the workflow state of an entity.
type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusReview     Status = "review"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
	StatusArchived   Status = "archived"
)

// ValidStatuses lists all valid status values.
var ValidStatuses = []Status{
	StatusBacklog, StatusTodo, StatusInProgress, StatusReview, StatusDone, StatusCancelled, StatusArchived,
}

// ParseStatus converts a string to a Status, returning false if invalid.
func ParseStatus(s string) (Status, bool) {
	st := Status(s)
	for _, v := range ValidStatuses {
		if v == st {
			return st, true
		}
	}
	return "", false
}

// Priority represents the importance level of an entity.
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// ValidPriorities lists all valid priority values.
var ValidPriorities = []Priority{
	PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical,
}

// ParsePriority converts a string to a Priority, returning false if invalid.
func ParsePriority(s string) (Priority, bool) {
	p := Priority(s)
	for _, v := range ValidPriorities {
		if v == p {
			return p, true
		}
	}
	return "", false
}

// Severity represents the impact level of a bug.
type Severity string

const (
	SeverityCosmetic Severity = "cosmetic"
	SeverityMinor    Severity = "minor"
	SeverityMajor    Severity = "major"
	SeverityBlocker  Severity = "blocker"
)

// ValidSeverities lists all valid severity values.
var ValidSeverities = []Severity{
	SeverityCosmetic, SeverityMinor, SeverityMajor, SeverityBlocker,
}

// ParseSeverity converts a string to a Severity, returning false if invalid.
func ParseSeverity(s string) (Severity, bool) {
	sv := Severity(s)
	for _, v := range ValidSeverities {
		if v == sv {
			return sv, true
		}
	}
	return "", false
}

// Comment represents a timestamped note on an entity.
type Comment struct {
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// Entity is the core data structure for all work items.
// It is serialized as meta.json in each entity's directory.
type Entity struct {
	ID          string     `json:"id"`
	Type        EntityType `json:"type"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Dependencies []string  `json:"dependencies,omitempty"`
	Comments    []Comment  `json:"comments,omitempty"`
	ParentID    string     `json:"parent_id,omitempty"`

	// Hierarchy: child slugs (only for epics and stories)
	Children []string `json:"children,omitempty"`

	// Bug-specific fields
	Severity         Severity `json:"severity,omitempty"`
	StepsToReproduce string   `json:"steps_to_reproduce,omitempty"`
}

// ShortID returns the first 8 characters of the UUID for display.
func (e *Entity) ShortID() string {
	if len(e.ID) >= 8 {
		return e.ID[:8]
	}
	return e.ID
}

// IsTerminal returns true if the entity cannot have children.
func (e *Entity) IsTerminal() bool {
	return e.Type == EntitySubTask || e.Type == EntityTask || e.Type == EntityBug
}

// ChildDir returns the subdirectory name where children are stored.
// Returns empty string for terminal entities.
func (e *Entity) ChildDir() string {
	switch e.Type {
	case EntityEpic:
		return "stories"
	case EntityStory:
		return "tasks"
	default:
		return ""
	}
}
