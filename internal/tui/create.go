package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

// StepKind identifies the type of a wizard step.
type StepKind int

const (
	StepText   StepKind = iota // Free-text input
	StepSelect                 // Arrow-key enum selection
	StepFilter                 // Filterable entity selection
)

// WizardStep defines one step of the interactive create wizard.
type WizardStep struct {
	Label       string
	Kind        StepKind
	Required    bool
	Placeholder string
	Options     []SelectOption // for StepSelect
	Choices     []EntityChoice // for StepFilter
	DefaultIdx  int            // for StepSelect: initially highlighted option
}

// WizardResult holds the collected values from a wizard run.
type WizardResult struct {
	Values    []string // one value per step, in order
	Cancelled bool
}

// --- bubbletea model ---

type wizardModel struct {
	steps   []WizardStep
	current int
	values  []string

	// text step state
	textInput textinput.Model

	// select step state
	selectCursor int

	// filter step state
	filterInput  textinput.Model
	filterCursor int
	filteredIdxs []int

	done      bool
	cancelled bool
}

func newWizardModel(steps []WizardStep) wizardModel {
	ti := textinput.New()
	ti.CharLimit = 200

	fi := textinput.New()
	fi.Placeholder = "Type to filter..."
	fi.CharLimit = 100

	m := wizardModel{
		steps:       steps,
		values:      make([]string, len(steps)),
		textInput:   ti,
		filterInput: fi,
	}
	return m
}

func (m wizardModel) Init() tea.Cmd {
	return m.initCurrentStep()
}

func (m wizardModel) initCurrentStep() tea.Cmd {
	if m.current >= len(m.steps) {
		return nil
	}
	step := m.steps[m.current]
	switch step.Kind {
	case StepText:
		m.textInput.SetValue("")
		m.textInput.Placeholder = step.Placeholder
		return m.textInput.Focus()
	case StepSelect:
		m.selectCursor = step.DefaultIdx
		return nil
	case StepFilter:
		m.filterInput.SetValue("")
		m.filterCursor = 0
		m.filteredIdxs = allIndices(len(step.Choices))
		return m.filterInput.Focus()
	}
	return nil
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		}
	}

	step := m.steps[m.current]
	switch step.Kind {
	case StepText:
		return m.updateText(msg)
	case StepSelect:
		return m.updateSelect(msg)
	case StepFilter:
		return m.updateFilter(msg)
	}
	return m, nil
}

func (m wizardModel) updateText(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			val := strings.TrimSpace(m.textInput.Value())
			if m.steps[m.current].Required && val == "" {
				return m, nil // don't advance
			}
			m.values[m.current] = val
			m.textInput.Blur()
			return m.advance()
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m wizardModel) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	step := m.steps[m.current]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectCursor > 0 {
				m.selectCursor--
			}
		case "down", "j":
			if m.selectCursor < len(step.Options)-1 {
				m.selectCursor++
			}
		case "enter":
			m.values[m.current] = step.Options[m.selectCursor].Value
			return m.advance()
		}
	}
	return m, nil
}

func (m wizardModel) updateFilter(msg tea.Msg) (tea.Model, tea.Cmd) {
	step := m.steps[m.current]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.filterCursor > 0 {
				m.filterCursor--
			}
			return m, nil
		case "down":
			if m.filterCursor < len(m.filteredIdxs)-1 {
				m.filterCursor++
			}
			return m, nil
		case "enter":
			if len(m.filteredIdxs) > 0 {
				idx := m.filteredIdxs[m.filterCursor]
				m.values[m.current] = step.Choices[idx].Value
				m.filterInput.Blur()
				return m.advance()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	// Re-filter
	query := m.filterInput.Value()
	if query == "" {
		m.filteredIdxs = allIndices(len(step.Choices))
	} else {
		labels := make([]string, len(step.Choices))
		for i, c := range step.Choices {
			labels[i] = c.Label
		}
		matches := fuzzy.Find(query, labels)
		m.filteredIdxs = make([]int, len(matches))
		for i, match := range matches {
			m.filteredIdxs[i] = match.Index
		}
	}

	if m.filterCursor >= len(m.filteredIdxs) {
		m.filterCursor = max(0, len(m.filteredIdxs)-1)
	}

	return m, cmd
}

func (m wizardModel) advance() (tea.Model, tea.Cmd) {
	m.current++
	if m.current >= len(m.steps) {
		m.done = true
		return m, tea.Quit
	}

	step := m.steps[m.current]
	switch step.Kind {
	case StepText:
		m.textInput.SetValue("")
		m.textInput.Placeholder = step.Placeholder
		cmd := m.textInput.Focus()
		return m, cmd
	case StepSelect:
		m.selectCursor = step.DefaultIdx
		return m, nil
	case StepFilter:
		m.filterInput.SetValue("")
		m.filterCursor = 0
		m.filteredIdxs = allIndices(len(step.Choices))
		cmd := m.filterInput.Focus()
		return m, cmd
	}
	return m, nil
}

func (m wizardModel) View() string {
	var b strings.Builder

	completedLabel := lipgloss.NewStyle().Faint(true)
	completedValue := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	activeLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	activeCursor := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	dim := lipgloss.NewStyle().Faint(true)

	// Completed steps
	for i := 0; i < m.current && i < len(m.steps); i++ {
		display := m.values[i]
		if display == "" {
			display = "-"
		}
		b.WriteString(fmt.Sprintf("  %s %s\n",
			completedLabel.Render(m.steps[i].Label),
			completedValue.Render(display),
		))
	}

	if m.current >= len(m.steps) {
		return b.String()
	}

	step := m.steps[m.current]
	b.WriteString(activeLabel.Render(step.Label) + "\n")

	switch step.Kind {
	case StepText:
		b.WriteString("  " + m.textInput.View() + "\n")
		if step.Required {
			b.WriteString(dim.Render("  required • enter to continue") + "\n")
		} else {
			b.WriteString(dim.Render("  optional • enter to skip") + "\n")
		}

	case StepSelect:
		for i, opt := range step.Options {
			if i == m.selectCursor {
				b.WriteString(fmt.Sprintf("  %s %s\n", activeCursor.Render(">"), activeCursor.Render(opt.Label)))
			} else {
				b.WriteString(fmt.Sprintf("    %s\n", opt.Label))
			}
		}
		b.WriteString(dim.Render("  ↑/↓ navigate • enter select") + "\n")

	case StepFilter:
		b.WriteString("  " + m.filterInput.View() + "\n")
		maxShow := 10
		for i, idx := range m.filteredIdxs {
			if i >= maxShow {
				remaining := len(m.filteredIdxs) - maxShow
				b.WriteString(dim.Render(fmt.Sprintf("  ... and %d more", remaining)) + "\n")
				break
			}
			choice := step.Choices[idx]
			if i == m.filterCursor {
				b.WriteString(fmt.Sprintf("  %s %s\n", activeCursor.Render(">"), activeCursor.Render(choice.Label)))
			} else {
				b.WriteString(fmt.Sprintf("    %s\n", choice.Label))
			}
		}
		if len(m.filteredIdxs) == 0 {
			b.WriteString(dim.Render("  no matches") + "\n")
		}
		b.WriteString(dim.Render("  type to filter • ↑/↓ navigate • enter select") + "\n")
	}

	return b.String()
}

// RunWizard runs the interactive wizard and returns collected values.
func RunWizard(steps []WizardStep) (*WizardResult, error) {
	m := newWizardModel(steps)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}
	rm := result.(wizardModel)
	return &WizardResult{
		Values:    rm.values,
		Cancelled: rm.cancelled,
	}, nil
}

// PromptText runs a single text input prompt and returns the entered value.
func PromptText(label, placeholder string, required bool) (string, error) {
	result, err := RunWizard([]WizardStep{
		{Label: label, Kind: StepText, Required: required, Placeholder: placeholder},
	})
	if err != nil {
		return "", err
	}
	if result.Cancelled {
		return "", fmt.Errorf("cancelled")
	}
	return result.Values[0], nil
}

func allIndices(n int) []int {
	idxs := make([]int, n)
	for i := range idxs {
		idxs[i] = i
	}
	return idxs
}
