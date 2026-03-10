package store

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cstsortan/prm/internal/model"
)

// ResolveResult contains the resolved entity and its directory path.
type ResolveResult struct {
	Entity *model.Entity
	Dir    string
}

// Resolve looks up an entity by ID (full or partial) or by slug.
// It searches in this order: exact UUID -> UUID prefix -> slug match.
// Returns an error if no match or ambiguous match.
func (s *Store) Resolve(idx *model.Index, ref string) (*ResolveResult, error) {
	// 1. Exact UUID match
	if path, ok := idx.Entries[ref]; ok {
		dir := s.EntityDir(path)
		entity, err := s.ReadEntity(dir)
		if err != nil {
			return nil, fmt.Errorf("reading entity: %w", err)
		}
		return &ResolveResult{Entity: entity, Dir: dir}, nil
	}

	// 2. UUID prefix match (minimum 4 chars)
	if len(ref) >= 4 {
		matches := idx.FindByPrefix(ref)
		if len(matches) == 1 {
			for id, path := range matches {
				_ = id
				dir := s.EntityDir(path)
				entity, err := s.ReadEntity(dir)
				if err != nil {
					return nil, fmt.Errorf("reading entity: %w", err)
				}
				return &ResolveResult{Entity: entity, Dir: dir}, nil
			}
		}
		if len(matches) > 1 {
			ids := make([]string, 0, len(matches))
			for id := range matches {
				ids = append(ids, id)
			}
			return nil, fmt.Errorf("ambiguous ID prefix %q matches: %s", ref, strings.Join(ids, ", "))
		}
	}

	// 3. Slug match - scan all entries
	var slugMatches []ResolveResult
	for _, path := range idx.Entries {
		slug := filepath.Base(path)
		if slug == ref {
			dir := s.EntityDir(path)
			entity, err := s.ReadEntity(dir)
			if err != nil {
				continue
			}
			slugMatches = append(slugMatches, ResolveResult{Entity: entity, Dir: dir})
		}
	}

	if len(slugMatches) == 1 {
		return &slugMatches[0], nil
	}
	if len(slugMatches) > 1 {
		ids := make([]string, 0, len(slugMatches))
		for _, m := range slugMatches {
			ids = append(ids, fmt.Sprintf("%s (%s)", m.Entity.ShortID(), m.Entity.Type))
		}
		return nil, fmt.Errorf("ambiguous slug %q matches: %s", ref, strings.Join(ids, ", "))
	}

	return nil, fmt.Errorf("no entity found matching %q", ref)
}
