package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Modell repräsentiert den UI-Zustand.
// Deutsche Kommentare gemäß Leitplanken.
type model struct{}

// initialModel erzeugt das Startmodell.
func initialModel() model { return model{} }

// Init wird beim Programmstart ausgeführt.
// Kein initialer Command notwendig.
func (m model) Init() tea.Cmd { return nil }

// Update verarbeitet Nachrichten (Events).
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Beenden bei Ctrl+C oder q.
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View rendert die Oberfläche.
func (m model) View() string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	return style.Render("SEW Forge") + "\n"
}

func main() {
	// Kontext mit Signal-Abbruch (Ctrl+C) zusätzlich zur Bubble Tea eigenen Behandlung.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	p := tea.NewProgram(initialModel(), tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		// Fehler wird klar ausgegeben, kein Panic.
		fmt.Fprintf(os.Stderr, "Fehler beim Ausführen der TUI: %v\n", err)
		os.Exit(1)
	}
}
