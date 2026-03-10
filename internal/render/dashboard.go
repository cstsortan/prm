package render

import (
	"fmt"
	"strings"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/service"
)

// Dashboard renders summary statistics.
func Dashboard(stats *service.DashboardStats, cfg *model.Config) string {
	var b strings.Builder

	b.WriteString(Title.Render(fmt.Sprintf("Project: %s", cfg.ProjectName)))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("=", 50))
	b.WriteString("\n\n")

	// By type
	typeOrder := []model.EntityType{
		model.EntityEpic, model.EntityStory, model.EntitySubTask,
		model.EntityTask, model.EntityBug,
	}
	for _, t := range typeOrder {
		count := stats.ByType[t]
		if count > 0 {
			b.WriteString(fmt.Sprintf("  %-12s %s\n", TypeStyle(t)+":", fmt.Sprintf("%d", count)))
		}
	}
	b.WriteString(fmt.Sprintf("\n  Total: %s\n", Bold.Render(fmt.Sprintf("%d", stats.Total))))

	// By status
	b.WriteString("\n")
	b.WriteString(Bold.Render("  By Status:"))
	b.WriteString("\n")
	for _, s := range model.ValidStatuses {
		count := stats.ByStatus[s]
		if count > 0 {
			b.WriteString(fmt.Sprintf("    %-15s %d\n", StatusStyle(s), count))
		}
	}

	// By priority
	b.WriteString("\n")
	b.WriteString(Bold.Render("  By Priority:"))
	b.WriteString("\n")
	for _, p := range model.ValidPriorities {
		count := stats.ByPriority[p]
		if count > 0 {
			b.WriteString(fmt.Sprintf("    %-15s %d\n", PriorityStyle(p), count))
		}
	}

	// By severity (bugs only)
	if len(stats.BySeverity) > 0 {
		b.WriteString("\n")
		b.WriteString(Bold.Render("  Bug Severity:"))
		b.WriteString("\n")
		for _, sv := range model.ValidSeverities {
			count := stats.BySeverity[sv]
			if count > 0 {
				b.WriteString(fmt.Sprintf("    %-15s %d\n", string(sv), count))
			}
		}
	}

	return b.String()
}
