package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	prompt    string
	confirmed bool
	done      bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			m.done = true
			return m, tea.Quit
		case "n", "N", "enter", "esc", "ctrl+c":
			m.confirmed = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render(m.prompt)
	hint := lipgloss.NewStyle().Faint(true).Render(" (y/N) ")
	return fmt.Sprintf("%s%s", prompt, hint)
}

// Confirm shows a y/N confirmation prompt and returns the user's choice.
func Confirm(prompt string) (bool, error) {
	p := tea.NewProgram(confirmModel{prompt: prompt})
	result, err := p.Run()
	if err != nil {
		return false, err
	}
	fmt.Println()
	return result.(confirmModel).confirmed, nil
}
