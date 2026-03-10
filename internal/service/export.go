package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportCSV writes all entities as CSV to the given writer.
func (svc *Service) ExportCSV(w io.Writer) error {
	all, err := svc.List(ListFilter{})
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"id", "type", "slug", "title", "description", "status", "priority",
		"tags", "created_at", "updated_at", "started_at", "completed_at",
		"due_date", "parent_id", "severity",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	for _, r := range all {
		e := r.Entity
		startedAt := ""
		if e.StartedAt != nil {
			startedAt = e.StartedAt.Format("2006-01-02T15:04:05Z")
		}
		completedAt := ""
		if e.CompletedAt != nil {
			completedAt = e.CompletedAt.Format("2006-01-02T15:04:05Z")
		}
		dueDate := ""
		if e.DueDate != nil {
			dueDate = e.DueDate.Format("2006-01-02")
		}

		record := []string{
			e.ID,
			string(e.Type),
			e.Slug,
			e.Title,
			e.Description,
			string(e.Status),
			string(e.Priority),
			strings.Join(e.Tags, ";"),
			e.CreatedAt.Format("2006-01-02T15:04:05Z"),
			e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			startedAt,
			completedAt,
			dueDate,
			e.ParentID,
			string(e.Severity),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("writing CSV record: %w", err)
		}
	}

	return nil
}

// ExportJSON writes all entities as a JSON array to the given writer.
func (svc *Service) ExportJSON(w io.Writer) error {
	all, err := svc.List(ListFilter{})
	if err != nil {
		return err
	}

	entities := make([]interface{}, len(all))
	for i, r := range all {
		entities[i] = r.Entity
	}

	data, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	data = append(data, '\n')

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("writing JSON: %w", err)
	}
	return nil
}
