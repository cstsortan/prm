package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cstsortan/prm/internal/model"
)

// RebuildIndex walks the .prm/ directory tree and rebuilds index.json
// from all meta.json files found.
func (s *Store) RebuildIndex() (*model.Index, error) {
	idx := model.NewIndex()

	err := filepath.WalkDir(s.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != metaFile {
			return nil
		}

		entity, readErr := s.ReadEntity(filepath.Dir(path))
		if readErr != nil {
			return fmt.Errorf("reading %s: %w", path, readErr)
		}

		relPath, relErr := s.RelPath(filepath.Dir(path))
		if relErr != nil {
			return fmt.Errorf("computing relative path for %s: %w", path, relErr)
		}

		idx.Set(entity.ID, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory tree: %w", err)
	}

	if err := s.WriteIndex(idx); err != nil {
		return nil, fmt.Errorf("writing rebuilt index: %w", err)
	}

	return idx, nil
}
