package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width         int
	height        int
	help          help.Model
	commandOutput string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.GitStatus):
			cmd := CreateGitStatusCommand()
			return m, RunCommand(cmd)
		}
	case RunCommandMsg:
		if msg.Err != nil {
			m.commandOutput = fmt.Sprintf("Error: %v", msg.Err)
		} else {
			m.commandOutput = msg.Output
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Render the header
	header := RenderHeader(m.width)

	// Render the footer
	footer := RenderFooter(m.width, m.help)

	// Calculate content height
	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
	if contentHeight < 0 {
		contentHeight = 0
	}

	// Render content
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight)

	content := contentStyle.Render(m.commandOutput)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func main() {
	m := model{
		help: help.New(),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
