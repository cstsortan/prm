package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/store"
)

// Service provides business logic for PRM operations.
type Service struct {
	Store *store.Store
}

// New creates a new Service backed by the given store.
func New(s *store.Store) *Service {
	return &Service{Store: s}
}

// Init initializes a new PRM project.
func (svc *Service) Init(projectName string) error {
	cfg := model.DefaultConfig(projectName)
	return svc.Store.Init(cfg)
}

// CreateEntityOpts holds parameters for creating a new entity.
type CreateEntityOpts struct {
	Type        model.EntityType
	Title       string
	Description string
	Priority    model.Priority
	Status      model.Status
	Tags        []string
	DueDate     *time.Time
	ParentID    string

	// README.md content (overrides auto-generated if non-empty)
	Body string

	// Dependencies (refs that will be resolved to UUIDs)
	Dependencies []string

	// Bug-specific
	Severity         model.Severity
	StepsToReproduce string
}

// CreateEntity creates a new entity, writes it to disk, and updates the index.
func (svc *Service) CreateEntity(opts CreateEntityOpts) (*model.Entity, error) {
	if strings.TrimSpace(opts.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}

	// Determine parent directory and resolve parent UUID
	parentDir, resolvedParentID, err := svc.resolveParentDir(idx, opts.Type, opts.ParentID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	slug := store.UniqueSlug(parentDir, store.GenerateSlug(opts.Title))

	// Resolve dependency refs to UUIDs
	var deps []string
	for _, ref := range opts.Dependencies {
		result, err := svc.Store.Resolve(idx, ref)
		if err != nil {
			return nil, fmt.Errorf("resolving dependency %q: %w", ref, err)
		}
		deps = append(deps, result.Entity.ID)
	}

	entity := &model.Entity{
		ID:               uuid.New().String(),
		Type:             opts.Type,
		Slug:             slug,
		Title:            opts.Title,
		Description:      opts.Description,
		Status:           opts.Status,
		Priority:         opts.Priority,
		Tags:             opts.Tags,
		CreatedAt:        now,
		UpdatedAt:        now,
		DueDate:          opts.DueDate,
		Dependencies:     deps,
		Comments:         nil,
		ParentID:         resolvedParentID,
		Severity:         opts.Severity,
		StepsToReproduce: opts.StepsToReproduce,
	}

	if opts.Status == model.StatusInProgress {
		entity.StartedAt = &now
	}

	entityDir := filepath.Join(parentDir, slug)
	readme := opts.Body
	if readme == "" {
		readme = svc.generateReadme(entity)
	}

	if err := svc.Store.WriteEntity(entityDir, entity, readme); err != nil {
		return nil, fmt.Errorf("writing entity: %w", err)
	}

	// Create child directory for non-terminal entities
	if childDir := entity.ChildDir(); childDir != "" {
		childPath := filepath.Join(entityDir, childDir)
		if err := mkdirIfNotExists(childPath); err != nil {
			return nil, fmt.Errorf("creating child dir: %w", err)
		}
	}

	// Update parent's children list
	if resolvedParentID != "" {
		if err := svc.addChildToParent(idx, resolvedParentID, slug); err != nil {
			return nil, fmt.Errorf("updating parent: %w", err)
		}
	}

	// Update index
	relPath, err := svc.Store.RelPath(entityDir)
	if err != nil {
		return nil, fmt.Errorf("computing relative path: %w", err)
	}
	idx.Set(entity.ID, relPath)
	if err := svc.Store.WriteIndex(idx); err != nil {
		return nil, fmt.Errorf("writing index: %w", err)
	}

	return entity, nil
}

// resolveParentDir determines the directory where a new entity should be created.
// Returns the parent directory path and the resolved parent UUID (empty for top-level entities).
func (svc *Service) resolveParentDir(idx *model.Index, entityType model.EntityType, parentRef string) (string, string, error) {
	switch entityType {
	case model.EntityEpic:
		return filepath.Join(svc.Store.Root(), "epics"), "", nil
	case model.EntityTask:
		return filepath.Join(svc.Store.Root(), "tasks"), "", nil
	case model.EntityBug:
		return filepath.Join(svc.Store.Root(), "bugs"), "", nil
	case model.EntityStory:
		if parentRef == "" {
			return "", "", fmt.Errorf("story requires --epic parent")
		}
		result, err := svc.Store.Resolve(idx, parentRef)
		if err != nil {
			return "", "", fmt.Errorf("resolving epic: %w", err)
		}
		if result.Entity.Type != model.EntityEpic {
			return "", "", fmt.Errorf("parent %q is a %s, not an epic", parentRef, result.Entity.Type)
		}
		return filepath.Join(result.Dir, "stories"), result.Entity.ID, nil
	case model.EntitySubTask:
		if parentRef == "" {
			return "", "", fmt.Errorf("sub-task requires --story parent")
		}
		result, err := svc.Store.Resolve(idx, parentRef)
		if err != nil {
			return "", "", fmt.Errorf("resolving story: %w", err)
		}
		if result.Entity.Type != model.EntityStory {
			return "", "", fmt.Errorf("parent %q is a %s, not a story", parentRef, result.Entity.Type)
		}
		return filepath.Join(result.Dir, "tasks"), result.Entity.ID, nil
	default:
		return "", "", fmt.Errorf("unknown entity type: %s", entityType)
	}
}

// addChildToParent reads the parent entity, appends the child slug, and rewrites it.
func (svc *Service) addChildToParent(idx *model.Index, parentID, childSlug string) error {
	result, err := svc.Store.Resolve(idx, parentID)
	if err != nil {
		return err
	}
	parent := result.Entity
	parent.Children = append(parent.Children, childSlug)
	parent.UpdatedAt = time.Now().UTC()

	readme, err := svc.Store.ReadEntityReadme(result.Dir)
	if err != nil {
		return err
	}
	return svc.Store.WriteEntity(result.Dir, parent, readme)
}

// generateReadme creates initial README.md content for an entity.
func (svc *Service) generateReadme(entity *model.Entity) string {
	content := fmt.Sprintf("# %s\n\n", entity.Title)
	if entity.Description != "" {
		content += entity.Description + "\n"
	}
	return content
}

func mkdirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}
