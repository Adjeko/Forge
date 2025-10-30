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
	header headerModel // eingebettetes Header-Model für den 4-zeiligen Kopfbereich
}

// initialModel erzeugt das Startmodell.
// initialModel erzeugt das Hauptmodell und initialisiert das Header-Model.
func initialModel() model { return model{header: newHeaderModel()} }

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
		// Header über Größe informieren, damit er seine Breite kennt.
		m.header, _ = m.header.Update(nachricht)
	}
	return m, nil
}

// View rendert die Oberfläche.
// View rendert zuerst den Header und darunter den eigentlichen Inhaltsbereich.
func (m model) View() string {
	// Falls Größe noch nicht verfügbar, Hinweis anzeigen.
	if !m.bereit {
		return "Initialisiere Ansicht ..."
	}

	// Header oben platzieren, darunter (zentriert) den bisherigen Inhalt als Platzhalter.
	head := m.header.View()
	titelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	content := titelStyle.Render("SEW Forge")

	// Restliche Höhe nach Abzug des Headers bestimmen (4 feste Zeilen Header).
	verfuegbareHoehe := m.hoehe - 4
	if verfuegbareHoehe < 1 {
		verfuegbareHoehe = 1
	}
	body := lipgloss.Place(m.breite, verfuegbareHoehe, lipgloss.Center, lipgloss.Center, content)

	return head + "\n" + body
}
