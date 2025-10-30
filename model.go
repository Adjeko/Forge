package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model hält den dynamischen Zustand inkl. Terminalgröße.
// Kommentare ausschließlich auf Deutsch gemäß Leitplanken.
type model struct {
	breite int
	hoehe  int
	bereit bool
}

// initialModel erzeugt das Startmodell.
func initialModel() model { return model{} }

// Init wird beim Programmstart ausgeführt.
// Kein initialer Command notwendig.
func (m model) Init() tea.Cmd { return nil }

// Update verarbeitet Nachrichten (Events) und reagiert auf Eingaben sowie Größenänderungen.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch nachricht := msg.(type) {
	case tea.KeyMsg:
		// Beenden bei Ctrl+C oder q.
		if nachricht.Type == tea.KeyCtrlC || nachricht.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Terminalgröße speichern und Flag setzen.
		m.breite = nachricht.Width
		m.hoehe = nachricht.Height
		m.bereit = true
	}
	return m, nil
}

// View rendert die Oberfläche.
func (m model) View() string {
	// Falls Größe noch nicht verfügbar, Hinweis anzeigen.
	if !m.bereit {
		return "Initialisiere Ansicht ..."
	}

	titelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	content := titelStyle.Render("SEW Forge")
	return lipgloss.Place(m.breite, m.hoehe, lipgloss.Center, lipgloss.Center, content)
}
