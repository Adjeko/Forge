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
// model hält den dynamischen Zustand inkl. Terminalgröße.
type model struct {
	breite int
	hoehe  int
	bereit bool // Flag ob erste Größe empfangen wurde
}

// initialModel erzeugt das Startmodell.
func initialModel() model { return model{} }

// Init wird beim Programmstart ausgeführt.
// Kein initialer Command notwendig.
func (m model) Init() tea.Cmd { return nil }

// Update verarbeitet Nachrichten (Events).
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
	// Inhalt vorbereiten.
	content := titelStyle.Render("SEW Forge")

	// Vollbildplatzierung mit Zentrierung.
	return lipgloss.Place(m.breite, m.hoehe, lipgloss.Center, lipgloss.Center, content)
}

func main() {
	// Kontext mit Signal-Abbruch (Ctrl+C) zusätzlich zur Bubble Tea eigenen Behandlung.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Alt-Screen aktivieren um gesamten Terminalbereich exklusiv zu nutzen.
	p := tea.NewProgram(initialModel(), tea.WithContext(ctx), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		// Fehler wird klar ausgegeben, kein Panic.
		fmt.Fprintf(os.Stderr, "Fehler beim Ausführen der TUI: %v\n", err)
		os.Exit(1)
	}
}
