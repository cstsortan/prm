package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

// SelectOption represents a choice in a selection prompt.
type SelectOption struct {
	Label string
	Value string
}

// EntityChoice represents an entity in a filterable selector.
type EntityChoice struct {
	Label string
	Value string
}

// --- Simple select (arrow keys, no filter) ---

type selectModel struct {
	label     string
	options   []SelectOption
	cursor    int
	selected  string
	cancelled bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.options[m.cursor].Value
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	var b strings.Builder

	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Render(m.label)
	b.WriteString(label + "\n")

	active := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	for i, opt := range m.options {
		if i == m.cursor {
			b.WriteString(fmt.Sprintf("  %s %s\n", active.Render(">"), active.Render(opt.Label)))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", opt.Label))
		}
	}

	b.WriteString(lipgloss.NewStyle().Faint(true).Render("  ↑/↓ navigate • enter select • esc cancel") + "\n")
	return b.String()
}

// PromptSelect shows a selection list and returns the chosen value.
func PromptSelect(label string, options []SelectOption, defaultIdx int) (string, error) {
	if defaultIdx < 0 || defaultIdx >= len(options) {
		defaultIdx = 0
	}
	p := tea.NewProgram(selectModel{
		label:   label,
		options: options,
		cursor:  defaultIdx,
	})
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	rm := result.(selectModel)
	if rm.cancelled {
		return "", fmt.Errorf("cancelled")
	}
	return rm.selected, nil
}

// --- Filterable entity select (text input + fuzzy matching) ---

type entitySelectModel struct {
	label        string
	choices      []EntityChoice
	filteredIdxs []int
	cursor       int
	input        textinput.Model
	selected     string
	cancelled    bool
}

func newEntitySelectModel(label string, choices []EntityChoice) entitySelectModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.CharLimit = 100

	idxs := make([]int, len(choices))
	for i := range choices {
		idxs[i] = i
	}

	return entitySelectModel{
		label:        label,
		choices:      choices,
		filteredIdxs: idxs,
		input:        ti,
	}
}

func (m entitySelectModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m entitySelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down":
			if m.cursor < len(m.filteredIdxs)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if len(m.filteredIdxs) > 0 {
				idx := m.filteredIdxs[m.cursor]
				m.selected = m.choices[idx].Value
				return m, tea.Quit
			}
			return m, nil
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Re-filter based on input
	query := m.input.Value()
	if query == "" {
		m.filteredIdxs = make([]int, len(m.choices))
		for i := range m.choices {
			m.filteredIdxs[i] = i
		}
	} else {
		labels := make([]string, len(m.choices))
		for i, c := range m.choices {
			labels[i] = c.Label
		}
		matches := fuzzy.Find(query, labels)
		m.filteredIdxs = make([]int, len(matches))
		for i, match := range matches {
			m.filteredIdxs[i] = match.Index
		}
	}

	if m.cursor >= len(m.filteredIdxs) {
		m.cursor = max(0, len(m.filteredIdxs)-1)
	}

	return m, cmd
}

func (m entitySelectModel) View() string {
	var b strings.Builder

	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Render(m.label)
	b.WriteString(label + "\n")
	b.WriteString(m.input.View() + "\n")

	active := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	dim := lipgloss.NewStyle().Faint(true)

	maxShow := 10
	for i, idx := range m.filteredIdxs {
		if i >= maxShow {
			remaining := len(m.filteredIdxs) - maxShow
			b.WriteString(dim.Render(fmt.Sprintf("  ... and %d more", remaining)) + "\n")
			break
		}
		choice := m.choices[idx]
		if i == m.cursor {
			b.WriteString(fmt.Sprintf("  %s %s\n", active.Render(">"), active.Render(choice.Label)))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", choice.Label))
		}
	}

	if len(m.filteredIdxs) == 0 {
		b.WriteString(dim.Render("  no matches") + "\n")
	}

	b.WriteString(dim.Render("  type to filter • ↑/↓ navigate • enter select • esc cancel") + "\n")
	return b.String()
}

// PromptEntitySelect shows a filterable entity selector and returns the chosen value.
func PromptEntitySelect(label string, choices []EntityChoice) (string, error) {
	p := tea.NewProgram(newEntitySelectModel(label, choices))
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	rm := result.(entitySelectModel)
	if rm.cancelled {
		return "", fmt.Errorf("cancelled")
	}
	return rm.selected, nil
}
