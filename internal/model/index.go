package model

// Index maps entity UUIDs to their relative paths within .prm/.
// This enables O(1) lookups by ID without traversing the directory tree.
type Index struct {
	Entries map[string]string `json:"entries"`
}

// NewIndex creates an empty index.
func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]string),
	}
}

// Set adds or updates an entry in the index.
func (idx *Index) Set(id, path string) {
	idx.Entries[id] = path
}

// Get returns the path for a given ID. Returns empty string if not found.
func (idx *Index) Get(id string) string {
	return idx.Entries[id]
}

// Delete removes an entry from the index.
func (idx *Index) Delete(id string) {
	delete(idx.Entries, id)
}

// FindByPrefix returns all entries whose ID starts with the given prefix.
func (idx *Index) FindByPrefix(prefix string) map[string]string {
	results := make(map[string]string)
	for id, path := range idx.Entries {
		if len(id) >= len(prefix) && id[:len(prefix)] == prefix {
			results[id] = path
		}
	}
	return results
}
