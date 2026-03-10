package service

import (
	"fmt"
	"strings"

	"github.com/cstsortan/prm/internal/model"
)

// SearchResult wraps a list result with a relevance score.
type SearchResult struct {
	Entity *model.Entity
	Path   string
	Score  int
}

// Search performs full-text search across all entities.
func (svc *Service) Search(query string, filter ListFilter) ([]SearchResult, error) {
	idx, err := svc.Store.ReadIndex()
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}

	queryLower := strings.ToLower(query)
	var results []SearchResult

	for _, path := range idx.Entries {
		dir := svc.Store.EntityDir(path)
		entity, err := svc.Store.ReadEntity(dir)
		if err != nil {
			continue
		}

		if !matchesFilter(entity, filter) {
			continue
		}

		score := scoreEntity(entity, queryLower)

		// Also search README.md content
		readme, _ := svc.Store.ReadEntityReadme(dir)
		if readme != "" && strings.Contains(strings.ToLower(readme), queryLower) {
			score += 1
		}

		if score > 0 {
			results = append(results, SearchResult{
				Entity: entity,
				Path:   path,
				Score:  score,
			})
		}
	}

	// Sort by score descending
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

func scoreEntity(entity *model.Entity, queryLower string) int {
	score := 0

	// Title match (highest weight)
	if strings.Contains(strings.ToLower(entity.Title), queryLower) {
		score += 10
	}

	// Description match
	if strings.Contains(strings.ToLower(entity.Description), queryLower) {
		score += 5
	}

	// Tag match
	for _, tag := range entity.Tags {
		if strings.Contains(strings.ToLower(tag), queryLower) {
			score += 3
		}
	}

	// Comment match
	for _, c := range entity.Comments {
		if strings.Contains(strings.ToLower(c.Text), queryLower) {
			score += 2
		}
	}

	// Slug match
	if strings.Contains(entity.Slug, queryLower) {
		score += 1
	}

	return score
}
