package render

import (
	"fmt"
	"strings"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/service"
)

// Table renders a list of entities as a formatted table.
func Table(results []service.ListResult) string {
	if len(results) == 0 {
		return Dim.Render("No items found.")
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf("  %-10s %-10s %-40s %-13s %-10s %s",
		"ID", "TYPE", "TITLE", "STATUS", "PRIORITY", "TAGS")
	b.WriteString(Bold.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 100))
	b.WriteString("\n")

	for _, r := range results {
		e := r.Entity
		title := e.Title
		if len(title) > 38 {
			title = title[:35] + "..."
		}
		tags := ""
		if len(e.Tags) > 0 {
			tags = TagStyle.Render(strings.Join(e.Tags, ","))
		}

		line := fmt.Sprintf("  %-10s %-10s %-40s %-13s %-10s %s",
			IDStyle.Render(e.ShortID()),
			TypeStyle(e.Type),
			title,
			StatusStyle(e.Status),
			PriorityStyle(e.Priority),
			tags,
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// DepInfo holds resolved dependency information for display.
type DepInfo struct {
	ID    string
	Title string
}

// EntityDetail renders a single entity in detail view.
// resolvedDeps maps dependency IDs to titles (nil if not resolved).
func EntityDetail(entity *model.Entity, readme string, resolvedDeps map[string]string) string {
	var b strings.Builder

	b.WriteString(Title.Render(entity.Title))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  ID:        %s\n", IDStyle.Render(entity.ID)))
	b.WriteString(fmt.Sprintf("  Type:      %s\n", TypeStyle(entity.Type)))
	b.WriteString(fmt.Sprintf("  Status:    %s\n", StatusStyle(entity.Status)))
	b.WriteString(fmt.Sprintf("  Priority:  %s\n", PriorityStyle(entity.Priority)))

	if entity.Severity != "" {
		b.WriteString(fmt.Sprintf("  Severity:  %s\n", SeverityStyle(entity.Severity)))
	}

	if len(entity.Tags) > 0 {
		b.WriteString(fmt.Sprintf("  Tags:      %s\n", TagStyle.Render(strings.Join(entity.Tags, ", "))))
	}

	b.WriteString(fmt.Sprintf("  Created:   %s\n", entity.CreatedAt.Format("2006-01-02 15:04")))
	b.WriteString(fmt.Sprintf("  Updated:   %s\n", entity.UpdatedAt.Format("2006-01-02 15:04")))

	if entity.StartedAt != nil {
		b.WriteString(fmt.Sprintf("  Started:   %s\n", entity.StartedAt.Format("2006-01-02 15:04")))
	}
	if entity.CompletedAt != nil {
		b.WriteString(fmt.Sprintf("  Completed: %s\n", entity.CompletedAt.Format("2006-01-02 15:04")))
	}
	if entity.DueDate != nil {
		b.WriteString(fmt.Sprintf("  Due:       %s\n", entity.DueDate.Format("2006-01-02")))
	}
	if entity.ParentID != "" {
		b.WriteString(fmt.Sprintf("  Parent:    %s\n", IDStyle.Render(entity.ParentID)))
	}
	if len(entity.Dependencies) > 0 {
		var depStrs []string
		for _, depID := range entity.Dependencies {
			if resolvedDeps != nil {
				if title, ok := resolvedDeps[depID]; ok {
					short := depID
					if len(short) >= 8 {
						short = short[:8]
					}
					depStrs = append(depStrs, fmt.Sprintf("%s (%s)", title, IDStyle.Render(short)))
					continue
				}
			}
			short := depID
			if len(short) >= 8 {
				short = short[:8]
			}
			depStrs = append(depStrs, IDStyle.Render(short))
		}
		b.WriteString(fmt.Sprintf("  Depends:   %s\n", strings.Join(depStrs, ", ")))
	}
	if len(entity.Children) > 0 {
		b.WriteString(fmt.Sprintf("  Children:  %s\n", strings.Join(entity.Children, ", ")))
	}

	if readme != "" {
		b.WriteString("\n")
		b.WriteString(Dim.Render("--- README.md ---"))
		b.WriteString("\n")
		b.WriteString(readme)
	}

	if len(entity.Comments) > 0 {
		b.WriteString("\n")
		b.WriteString(Dim.Render("--- Comments ---"))
		b.WriteString("\n")
		for _, c := range entity.Comments {
			b.WriteString(fmt.Sprintf("  [%s] %s: %s\n",
				Dim.Render(c.CreatedAt.Format("2006-01-02 15:04")),
				Bold.Render(c.Author),
				c.Text,
			))
		}
	}

	return b.String()
}
