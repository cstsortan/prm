package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cstsortan/prm/internal/model"
)

// SortField defines what field to sort by.
type SortField string

const (
	SortPriority SortField = "priority"
	SortStatus   SortField = "status"
	SortTitle    SortField = "title"
	SortCreated  SortField = "created"
	SortUpdated  SortField = "updated"
	SortType     SortField = "type"
)

// ValidSortFields lists all valid sort field values.
var ValidSortFields = []SortField{
	SortPriority, SortStatus, SortTitle, SortCreated, SortUpdated, SortType,
}

// ParseSortField converts a string to a SortField, returning false if invalid.
func ParseSortField(s string) (SortField, bool) {
	sf := SortField(s)
	for _, v := range ValidSortFields {
		if v == sf {
			return sf, true
		}
	}
	return "", false
}

// ListFilter defines criteria for filtering entities.
type ListFilter struct {
	Types        []model.EntityType
	Statuses     []model.Status
	Priorities   []model.Priority
	Tags         []string
	Sort         SortField
	SortDesc     bool
	IncludeArchived bool
}

// ListResult contains an entity and its path for display.
type ListResult struct {
	Entity *model.Entity
	Path   string
}

// List returns all entities matching the given filter.
func (svc *Service) List(filter ListFilter) ([]ListResult, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}

	var results []ListResult
	for _, path := range idx.Entries {
		dir := svc.Store.EntityDir(path)
		entity, err := svc.Store.ReadEntity(dir)
		if err != nil {
			continue
		}

		if !matchesFilter(entity, filter) {
			continue
		}

		results = append(results, ListResult{Entity: entity, Path: path})
	}

	sortResults(results, filter.Sort, filter.SortDesc)

	return results, nil
}

func matchesFilter(entity *model.Entity, f ListFilter) bool {
	// Exclude archived unless explicitly requested or filtering by archived status
	if !f.IncludeArchived && entity.Status == model.StatusArchived {
		if len(f.Statuses) == 0 || !containsStatus(f.Statuses, model.StatusArchived) {
			return false
		}
	}
	if len(f.Types) > 0 && !containsType(f.Types, entity.Type) {
		return false
	}
	if len(f.Statuses) > 0 && !containsStatus(f.Statuses, entity.Status) {
		return false
	}
	if len(f.Priorities) > 0 && !containsPriority(f.Priorities, entity.Priority) {
		return false
	}
	if len(f.Tags) > 0 && !hasAnyTag(entity.Tags, f.Tags) {
		return false
	}
	return true
}

func containsType(slice []model.EntityType, val model.EntityType) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func containsStatus(slice []model.Status, val model.Status) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func containsPriority(slice []model.Priority, val model.Priority) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func hasAnyTag(entityTags, filterTags []string) bool {
	for _, ft := range filterTags {
		for _, et := range entityTags {
			if strings.EqualFold(et, ft) {
				return true
			}
		}
	}
	return false
}

func sortResults(results []ListResult, field SortField, desc bool) {
	if field == "" {
		field = SortPriority
	}

	sort.SliceStable(results, func(i, j int) bool {
		a, b := results[i].Entity, results[j].Entity
		var less bool
		switch field {
		case SortTitle:
			less = strings.ToLower(a.Title) < strings.ToLower(b.Title)
		case SortStatus:
			less = statusOrder(a.Status) < statusOrder(b.Status)
		case SortCreated:
			less = a.CreatedAt.Before(b.CreatedAt)
		case SortUpdated:
			less = a.UpdatedAt.After(b.UpdatedAt) // newest first by default
		case SortType:
			less = typeOrder(a.Type) < typeOrder(b.Type)
		default: // SortPriority
			pi, pj := priorityOrder(a.Priority), priorityOrder(b.Priority)
			if pi != pj {
				less = pi > pj // highest first by default
			} else {
				less = a.UpdatedAt.After(b.UpdatedAt)
			}
		}
		if desc {
			return !less
		}
		return less
	})
}

func statusOrder(s model.Status) int {
	switch s {
	case model.StatusInProgress:
		return 0
	case model.StatusReview:
		return 1
	case model.StatusTodo:
		return 2
	case model.StatusBacklog:
		return 3
	case model.StatusDone:
		return 4
	case model.StatusCancelled:
		return 5
	case model.StatusArchived:
		return 6
	default:
		return 7
	}
}

func typeOrder(t model.EntityType) int {
	switch t {
	case model.EntityEpic:
		return 0
	case model.EntityStory:
		return 1
	case model.EntitySubTask:
		return 2
	case model.EntityTask:
		return 3
	case model.EntityBug:
		return 4
	default:
		return 5
	}
}

func priorityOrder(p model.Priority) int {
	switch p {
	case model.PriorityCritical:
		return 4
	case model.PriorityHigh:
		return 3
	case model.PriorityMedium:
		return 2
	case model.PriorityLow:
		return 1
	default:
		return 0
	}
}

// ParseTypes parses a comma-separated string into entity types.
func ParseTypes(s string) []model.EntityType {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var types []model.EntityType
	for _, p := range parts {
		p = strings.TrimSpace(p)
		types = append(types, model.EntityType(p))
	}
	return types
}

// ParseStatuses parses a comma-separated string into statuses.
// Returns an error if any value is invalid.
func ParseStatuses(s string) []model.Status {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var statuses []model.Status
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		statuses = append(statuses, model.Status(p))
	}
	return statuses
}

// ValidateStatuses checks that all statuses in the slice are valid.
func ValidateStatuses(statuses []model.Status) error {
	for _, s := range statuses {
		if _, ok := model.ParseStatus(string(s)); !ok {
			return fmt.Errorf("invalid status: %q (valid: backlog, todo, in-progress, review, done, cancelled, archived)", s)
		}
	}
	return nil
}

// ParsePriorities parses a comma-separated string into priorities.
func ParsePriorities(s string) []model.Priority {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var priorities []model.Priority
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		priorities = append(priorities, model.Priority(p))
	}
	return priorities
}

// ValidatePriorities checks that all priorities in the slice are valid.
func ValidatePriorities(priorities []model.Priority) error {
	for _, p := range priorities {
		if _, ok := model.ParsePriority(string(p)); !ok {
			return fmt.Errorf("invalid priority: %q (valid: low, medium, high, critical)", p)
		}
	}
	return nil
}

// ParseTags parses a comma-separated string into tags.
func ParseTags(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var tags []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}

// DashboardStats holds summary statistics.
type DashboardStats struct {
	Total      int
	ByType     map[model.EntityType]int
	ByStatus   map[model.Status]int
	ByPriority map[model.Priority]int
	BySeverity map[model.Severity]int
}

// Dashboard computes summary statistics across all entities.
func (svc *Service) Dashboard() (*DashboardStats, *model.Config, error) {
	cfg, err := svc.Store.ReadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("reading config: %w", err)
	}

	all, err := svc.List(ListFilter{})
	if err != nil {
		return nil, nil, err
	}

	stats := &DashboardStats{
		Total:      len(all),
		ByType:     make(map[model.EntityType]int),
		ByStatus:   make(map[model.Status]int),
		ByPriority: make(map[model.Priority]int),
		BySeverity: make(map[model.Severity]int),
	}

	for _, r := range all {
		e := r.Entity
		stats.ByType[e.Type]++
		stats.ByStatus[e.Status]++
		stats.ByPriority[e.Priority]++
		if e.Type == model.EntityBug && e.Severity != "" {
			stats.BySeverity[e.Severity]++
		}
	}

	return stats, cfg, nil
}
