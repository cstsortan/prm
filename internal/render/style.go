package render

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/cstsortan/prm/internal/model"
)

var (
	Bold      = lipgloss.NewStyle().Bold(true)
	Dim       = lipgloss.NewStyle().Faint(true)
	Title     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	Success   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	Warning   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	Error     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	IDStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	TagStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
)

// StatusStyle returns a styled string for a status.
func StatusStyle(s model.Status) string {
	var style lipgloss.Style
	switch s {
	case model.StatusBacklog:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	case model.StatusTodo:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	case model.StatusInProgress:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	case model.StatusReview:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("177"))
	case model.StatusDone:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	case model.StatusCancelled:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	case model.StatusArchived:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Faint(true)
	default:
		style = lipgloss.NewStyle()
	}
	return style.Render(string(s))
}

// PriorityStyle returns a styled string for a priority.
func PriorityStyle(p model.Priority) string {
	var style lipgloss.Style
	switch p {
	case model.PriorityCritical:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	case model.PriorityHigh:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	case model.PriorityMedium:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	case model.PriorityLow:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	default:
		style = lipgloss.NewStyle()
	}
	return style.Render(string(p))
}

// SeverityStyle returns a styled string for a severity.
func SeverityStyle(s model.Severity) string {
	var style lipgloss.Style
	switch s {
	case model.SeverityBlocker:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	case model.SeverityMajor:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	case model.SeverityMinor:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	case model.SeverityCosmetic:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	default:
		style = lipgloss.NewStyle()
	}
	return style.Render(string(s))
}

// TypeStyle returns a styled string for an entity type.
func TypeStyle(t model.EntityType) string {
	var style lipgloss.Style
	switch t {
	case model.EntityEpic:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("177")).Bold(true)
	case model.EntityStory:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	case model.EntitySubTask:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	case model.EntityTask:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	case model.EntityBug:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	default:
		style = lipgloss.NewStyle()
	}
	return style.Render(string(t))
}
