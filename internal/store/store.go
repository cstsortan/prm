package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cstsortan/prm/internal/model"
)

const (
	prmDir     = ".prm"
	configFile = "prm.json"
	indexFile  = "index.json"
	metaFile   = "meta.json"
	readmeFile = "README.md"
)

// Store handles all filesystem operations for PRM data.
type Store struct {
	root string // absolute path to .prm/ directory
}

// New creates a Store rooted at the given project directory.
// It does not verify the directory exists yet.
func New(projectDir string) *Store {
	return &Store{
		root: filepath.Join(projectDir, prmDir),
	}
}

// Root returns the absolute path to the .prm/ directory.
func (s *Store) Root() string {
	return s.root
}

// Exists checks whether .prm/ directory has been initialized.
func (s *Store) Exists() bool {
	info, err := os.Stat(s.root)
	return err == nil && info.IsDir()
}

// Init creates the .prm/ directory structure and writes initial config and index.
func (s *Store) Init(cfg *model.Config) error {
	if s.Exists() {
		return fmt.Errorf("prm already initialized in %s", s.root)
	}

	dirs := []string{
		s.root,
		filepath.Join(s.root, "epics"),
		filepath.Join(s.root, "tasks"),
		filepath.Join(s.root, "bugs"),
		filepath.Join(s.root, "docs"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	if err := s.WriteConfig(cfg); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	if err := s.WriteIndex(model.NewIndex()); err != nil {
		return fmt.Errorf("writing index: %w", err)
	}

	return nil
}

// WriteConfig writes prm.json atomically.
func (s *Store) WriteConfig(cfg *model.Config) error {
	return writeJSON(filepath.Join(s.root, configFile), cfg)
}

// ReadConfig reads prm.json.
func (s *Store) ReadConfig() (*model.Config, error) {
	var cfg model.Config
	if err := readJSON(filepath.Join(s.root, configFile), &cfg); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return &cfg, nil
}

// WriteIndex writes index.json atomically.
func (s *Store) WriteIndex(idx *model.Index) error {
	return writeJSON(filepath.Join(s.root, indexFile), idx)
}

// ReadIndex reads index.json.
func (s *Store) ReadIndex() (*model.Index, error) {
	var idx model.Index
	if err := readJSON(filepath.Join(s.root, indexFile), &idx); err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	if idx.Entries == nil {
		idx.Entries = make(map[string]string)
	}
	return &idx, nil
}

// WriteEntity writes meta.json and README.md for an entity.
// The dir is the absolute path to the entity's directory.
func (s *Store) WriteEntity(dir string, entity *model.Entity, readmeContent string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating entity dir: %w", err)
	}

	if err := writeJSON(filepath.Join(dir, metaFile), entity); err != nil {
		return fmt.Errorf("writing meta.json: %w", err)
	}

	if err := WriteFileAtomic(filepath.Join(dir, readmeFile), []byte(readmeContent)); err != nil {
		return fmt.Errorf("writing README.md: %w", err)
	}

	return nil
}

// ReadEntity reads meta.json from the given directory.
func (s *Store) ReadEntity(dir string) (*model.Entity, error) {
	var entity model.Entity
	if err := readJSON(filepath.Join(dir, metaFile), &entity); err != nil {
		return nil, fmt.Errorf("reading entity from %s: %w", dir, err)
	}
	return &entity, nil
}

// ReadEntityReadme reads the README.md from the given directory.
func (s *Store) ReadEntityReadme(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, readmeFile))
	if err != nil {
		return "", fmt.Errorf("reading README.md from %s: %w", dir, err)
	}
	return string(data), nil
}

// DeleteEntity removes an entity directory and all its contents.
func (s *Store) DeleteEntity(dir string) error {
	return os.RemoveAll(dir)
}

// EntityDir returns the absolute path for an entity given its relative path in the index.
func (s *Store) EntityDir(relPath string) string {
	return filepath.Join(s.root, relPath)
}

// RelPath returns the path relative to .prm/ for a given absolute path.
func (s *Store) RelPath(absPath string) (string, error) {
	return filepath.Rel(s.root, absPath)
}

// ListDirs returns the names of subdirectories in the given directory.
func (s *Store) ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs, nil
}

// writeJSON marshals v to indented JSON and writes it atomically.
func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	data = append(data, '\n')
	return WriteFileAtomic(path, data)
}

// readJSON reads a JSON file and unmarshals it into v.
func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filepath.Base(path), err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}
	return nil
}

// WriteFileAtomic writes data to a temp file then renames it to the target path.
func WriteFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".prm-tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}
