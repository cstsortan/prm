package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/store"
)

// ShowEntity resolves an entity reference and returns it with its directory.
// Also returns a map of dependency IDs to titles.
func (svc *Service) ShowEntity(ref string) (*model.Entity, string, string, map[string]string, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, "", "", nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, "", "", nil, err
	}
	readme, _ := svc.Store.ReadEntityReadme(result.Dir)

	// Resolve dependency IDs to titles
	var depMap map[string]string
	if len(result.Entity.Dependencies) > 0 {
		depMap = make(map[string]string)
		for _, depID := range result.Entity.Dependencies {
			depResult, err := svc.Store.Resolve(idx, depID)
			if err == nil {
				depMap[depID] = depResult.Entity.Title
			}
		}
	}

	return result.Entity, result.Dir, readme, depMap, nil
}

// UpdateEntityOpts holds optional fields for updating an entity.
// Nil/empty values mean "don't change".
type UpdateEntityOpts struct {
	Title        string
	Description  string
	Body         string // README.md content
	Priority     string
	Tags         []string
	DueDate      *time.Time
	ClearDue     bool
	Severity     string
	Dependencies []string // refs to resolve to UUIDs
	ClearDepends bool
}

// UpdateEntity modifies fields on an existing entity.
func (svc *Service) UpdateEntity(ref string, opts UpdateEntityOpts) (*model.Entity, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, err
	}

	entity := result.Entity
	changed := false

	if opts.Title != "" && opts.Title != entity.Title {
		entity.Title = opts.Title
		changed = true
	}
	if opts.Description != "" {
		entity.Description = opts.Description
		changed = true
	}
	if opts.Priority != "" {
		p, ok := model.ParsePriority(opts.Priority)
		if !ok {
			return nil, fmt.Errorf("invalid priority: %s", opts.Priority)
		}
		entity.Priority = p
		changed = true
	}
	if opts.Tags != nil {
		// Filter out empty strings
		var filtered []string
		for _, t := range opts.Tags {
			t = strings.TrimSpace(t)
			if t != "" {
				filtered = append(filtered, t)
			}
		}
		entity.Tags = filtered
		changed = true
	}
	if opts.DueDate != nil {
		entity.DueDate = opts.DueDate
		changed = true
	}
	if opts.ClearDue {
		entity.DueDate = nil
		changed = true
	}
	if opts.Severity != "" {
		sv, ok := model.ParseSeverity(opts.Severity)
		if !ok {
			return nil, fmt.Errorf("invalid severity: %s", opts.Severity)
		}
		entity.Severity = sv
		changed = true
	}
	if opts.ClearDepends {
		entity.Dependencies = nil
		changed = true
	} else if opts.Dependencies != nil {
		var resolved []string
		for _, ref := range opts.Dependencies {
			r, err := svc.Store.Resolve(idx, ref)
			if err != nil {
				return nil, fmt.Errorf("resolving dependency %q: %w", ref, err)
			}
			resolved = append(resolved, r.Entity.ID)
		}
		entity.Dependencies = resolved
		changed = true
	}

	if opts.Body != "" {
		changed = true
	}

	if !changed {
		return entity, nil
	}

	entity.UpdatedAt = time.Now().UTC()
	readme := opts.Body
	if readme == "" {
		var err error
		readme, err = svc.Store.ReadEntityReadme(result.Dir)
		if err != nil {
			return nil, err
		}
	}
	if err := svc.Store.WriteEntity(result.Dir, entity, readme); err != nil {
		return nil, fmt.Errorf("writing entity: %w", err)
	}
	return entity, nil
}

// SetStatus changes an entity's status with automatic timestamp tracking.
func (svc *Service) SetStatus(ref string, newStatus string) (*model.Entity, error) {
	status, ok := model.ParseStatus(newStatus)
	if !ok {
		return nil, fmt.Errorf("invalid status: %s (valid: backlog, todo, in-progress, review, done, cancelled, archived)", newStatus)
	}

	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, err
	}

	entity := result.Entity
	now := time.Now().UTC()
	entity.Status = status
	entity.UpdatedAt = now

	// Track timestamps
	if status == model.StatusInProgress && entity.StartedAt == nil {
		entity.StartedAt = &now
	}
	if status == model.StatusDone || status == model.StatusCancelled {
		entity.CompletedAt = &now
	}

	readme, err := svc.Store.ReadEntityReadme(result.Dir)
	if err != nil {
		return nil, err
	}
	if err := svc.Store.WriteEntity(result.Dir, entity, readme); err != nil {
		return nil, fmt.Errorf("writing entity: %w", err)
	}
	return entity, nil
}

// AddComment adds a comment to an entity.
func (svc *Service) AddComment(ref, author, text string) (*model.Entity, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("comment text cannot be empty")
	}
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("comment author cannot be empty")
	}

	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, err
	}

	entity := result.Entity
	now := time.Now().UTC()
	entity.Comments = append(entity.Comments, model.Comment{
		Author:    author,
		Text:      text,
		CreatedAt: now,
	})
	entity.UpdatedAt = now

	readme, err := svc.Store.ReadEntityReadme(result.Dir)
	if err != nil {
		return nil, err
	}
	if err := svc.Store.WriteEntity(result.Dir, entity, readme); err != nil {
		return nil, fmt.Errorf("writing entity: %w", err)
	}
	return entity, nil
}

// DepRef holds a resolved dependency reference for display.
type DepRef struct {
	ID      string
	ShortID string
	Title   string
	Status  model.Status
}

// DepsResult holds the full dependency information for an entity.
type DepsResult struct {
	Entity     *model.Entity
	DependsOn  []DepRef
	DependedBy []DepRef
}

// GetDependencies returns what an entity depends on and what depends on it.
func (svc *Service) GetDependencies(ref string) (*DepsResult, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, err
	}

	entity := result.Entity
	deps := &DepsResult{Entity: entity}

	// Forward deps: what this entity depends on
	for _, depID := range entity.Dependencies {
		depResult, err := svc.Store.Resolve(idx, depID)
		if err != nil {
			deps.DependsOn = append(deps.DependsOn, DepRef{ID: depID, ShortID: depID[:8], Title: "(not found)"})
			continue
		}
		deps.DependsOn = append(deps.DependsOn, DepRef{
			ID:      depResult.Entity.ID,
			ShortID: depResult.Entity.ShortID(),
			Title:   depResult.Entity.Title,
			Status:  depResult.Entity.Status,
		})
	}

	// Reverse deps: what depends on this entity
	for _, path := range idx.Entries {
		dir := svc.Store.EntityDir(path)
		other, err := svc.Store.ReadEntity(dir)
		if err != nil || other.ID == entity.ID {
			continue
		}
		for _, depID := range other.Dependencies {
			if depID == entity.ID {
				deps.DependedBy = append(deps.DependedBy, DepRef{
					ID:      other.ID,
					ShortID: other.ShortID(),
					Title:   other.Title,
					Status:  other.Status,
				})
				break
			}
		}
	}

	return deps, nil
}

// MoveEntity reparents an entity. For stories, moves to a different epic.
// For sub-tasks, moves to a different story. If standalone is true, promotes
// a sub-task to a standalone task.
func (svc *Service) MoveEntity(ref string, newParentRef string, standalone bool) (*model.Entity, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return nil, err
	}

	entity := result.Entity
	oldDir := result.Dir

	// Validate the move is allowed
	if standalone {
		if entity.Type != model.EntitySubTask {
			return nil, fmt.Errorf("only sub-tasks can be promoted to standalone tasks")
		}
	} else {
		if entity.Type != model.EntityStory && entity.Type != model.EntitySubTask {
			return nil, fmt.Errorf("only stories and sub-tasks can be moved (got %s)", entity.Type)
		}
	}

	// Remove from old parent's children list
	if entity.ParentID != "" {
		oldParentResult, err := svc.Store.Resolve(idx, entity.ParentID)
		if err == nil {
			parent := oldParentResult.Entity
			newChildren := make([]string, 0, len(parent.Children))
			for _, c := range parent.Children {
				if c != entity.Slug {
					newChildren = append(newChildren, c)
				}
			}
			parent.Children = newChildren
			parent.UpdatedAt = time.Now().UTC()
			readme, _ := svc.Store.ReadEntityReadme(oldParentResult.Dir)
			if err := svc.Store.WriteEntity(oldParentResult.Dir, parent, readme); err != nil {
				return nil, fmt.Errorf("updating old parent: %w", err)
			}
		}
	}

	now := time.Now().UTC()

	if standalone {
		// Promote sub-task to standalone task
		entity.Type = model.EntityTask
		entity.ParentID = ""
		entity.UpdatedAt = now

		newParentDir := svc.Store.Root() + "/tasks"
		newSlug := store.UniqueSlug(newParentDir, entity.Slug)
		entity.Slug = newSlug
		newDir := newParentDir + "/" + newSlug

		readme, _ := svc.Store.ReadEntityReadme(oldDir)
		if err := svc.Store.WriteEntity(newDir, entity, readme); err != nil {
			return nil, fmt.Errorf("writing moved entity: %w", err)
		}
		if err := svc.Store.DeleteEntity(oldDir); err != nil {
			return nil, fmt.Errorf("removing old directory: %w", err)
		}

		relPath, _ := svc.Store.RelPath(newDir)
		idx.Delete(entity.ID)
		idx.Set(entity.ID, relPath)
		if err := svc.Store.WriteIndex(idx); err != nil {
			return nil, fmt.Errorf("writing index: %w", err)
		}
		return entity, nil
	}

	// Move to new parent
	newParentResult, err := svc.Store.Resolve(idx, newParentRef)
	if err != nil {
		return nil, fmt.Errorf("resolving new parent: %w", err)
	}
	newParent := newParentResult.Entity

	// Validate parent type compatibility
	switch entity.Type {
	case model.EntityStory:
		if newParent.Type != model.EntityEpic {
			return nil, fmt.Errorf("stories can only be moved to epics (got %s)", newParent.Type)
		}
	case model.EntitySubTask:
		if newParent.Type != model.EntityStory {
			return nil, fmt.Errorf("sub-tasks can only be moved to stories (got %s)", newParent.Type)
		}
	}

	// Compute new directory
	newParentDir := newParentResult.Dir + "/" + newParent.ChildDir()
	newSlug := store.UniqueSlug(newParentDir, entity.Slug)
	entity.Slug = newSlug
	newDir := newParentDir + "/" + newSlug

	entity.ParentID = newParent.ID
	entity.UpdatedAt = now

	// Write to new location
	readme, _ := svc.Store.ReadEntityReadme(oldDir)
	if err := svc.Store.WriteEntity(newDir, entity, readme); err != nil {
		return nil, fmt.Errorf("writing moved entity: %w", err)
	}

	// Move children directories if entity has children (story with sub-tasks)
	if childDir := entity.ChildDir(); childDir != "" {
		oldChildDir := oldDir + "/" + childDir
		newChildDir := newDir + "/" + childDir
		if children, err := svc.Store.ListDirs(oldChildDir); err == nil && len(children) > 0 {
			if err := os.MkdirAll(newChildDir, 0755); err != nil {
				return nil, fmt.Errorf("creating child dir: %w", err)
			}
			for _, child := range children {
				if err := os.Rename(oldChildDir+"/"+child, newChildDir+"/"+child); err != nil {
					return nil, fmt.Errorf("moving child %s: %w", child, err)
				}
			}
		}
	}

	// Remove old directory
	if err := svc.Store.DeleteEntity(oldDir); err != nil {
		return nil, fmt.Errorf("removing old directory: %w", err)
	}

	// Add to new parent's children list
	newParent.Children = append(newParent.Children, entity.Slug)
	newParent.UpdatedAt = now
	parentReadme, _ := svc.Store.ReadEntityReadme(newParentResult.Dir)
	if err := svc.Store.WriteEntity(newParentResult.Dir, newParent, parentReadme); err != nil {
		return nil, fmt.Errorf("updating new parent: %w", err)
	}

	// Update index for entity and all children
	idx.Delete(entity.ID)
	relPath, _ := svc.Store.RelPath(newDir)
	idx.Set(entity.ID, relPath)

	// Re-index children
	if childDir := entity.ChildDir(); childDir != "" {
		svc.reindexChildren(idx, newDir+"/"+childDir)
	}

	if err := svc.Store.WriteIndex(idx); err != nil {
		return nil, fmt.Errorf("writing index: %w", err)
	}

	return entity, nil
}

// reindexChildren walks child directories and updates the index.
func (svc *Service) reindexChildren(idx *model.Index, dir string) {
	children, err := svc.Store.ListDirs(dir)
	if err != nil {
		return
	}
	for _, child := range children {
		childDir := dir + "/" + child
		entity, err := svc.Store.ReadEntity(childDir)
		if err != nil {
			continue
		}
		relPath, _ := svc.Store.RelPath(childDir)
		idx.Delete(entity.ID)
		idx.Set(entity.ID, relPath)

		if subDir := entity.ChildDir(); subDir != "" {
			svc.reindexChildren(idx, childDir+"/"+subDir)
		}
	}
}

// DeleteEntity removes an entity and updates the parent's children list and the index.
func (svc *Service) DeleteEntity(ref string) error {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return fmt.Errorf("reading index: %w", err)
	}
	result, err := svc.Store.Resolve(idx, ref)
	if err != nil {
		return err
	}

	entity := result.Entity

	// Remove from parent's children list
	if entity.ParentID != "" {
		parentResult, err := svc.Store.Resolve(idx, entity.ParentID)
		if err == nil {
			parent := parentResult.Entity
			newChildren := make([]string, 0, len(parent.Children))
			for _, c := range parent.Children {
				if c != entity.Slug {
					newChildren = append(newChildren, c)
				}
			}
			parent.Children = newChildren
			parent.UpdatedAt = time.Now().UTC()
			readme, _ := svc.Store.ReadEntityReadme(parentResult.Dir)
			if err := svc.Store.WriteEntity(parentResult.Dir, parent, readme); err != nil {
				return fmt.Errorf("updating parent children list: %w", err)
			}
		}
	}

	// Collect all child entity IDs to remove from index
	svc.collectAndRemoveFromIndex(idx, result.Dir)

	// Remove entity directory
	if err := svc.Store.DeleteEntity(result.Dir); err != nil {
		return fmt.Errorf("deleting entity directory: %w", err)
	}

	// Remove from index and save
	idx.Delete(entity.ID)
	if err := svc.Store.WriteIndex(idx); err != nil {
		return fmt.Errorf("writing index: %w", err)
	}

	return nil
}

// collectAndRemoveFromIndex recursively removes all child entities from the index.
func (svc *Service) collectAndRemoveFromIndex(idx *model.Index, dir string) {
	entity, err := svc.Store.ReadEntity(dir)
	if err != nil {
		return
	}
	idx.Delete(entity.ID)

	childDir := entity.ChildDir()
	if childDir == "" {
		return
	}

	children, err := svc.Store.ListDirs(dir + "/" + childDir)
	if err != nil {
		return
	}
	for _, child := range children {
		svc.collectAndRemoveFromIndex(idx, dir+"/"+childDir+"/"+child)
	}
}
